/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package fusionstore

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/minio/madmin-go"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"github.com/minio/minio-go/v7/pkg/tags"
	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/fusionstore/mgs"
	"github.com/minio/minio/fusionstore/object"
	"github.com/minio/minio/fusionstore/pool"
	"github.com/minio/minio/fusionstore/sdk"
	"github.com/minio/minio/fusionstore/vbucket"
	"github.com/minio/pkg/bucket/policy"
)

const (
	serviceTimeout = 10
)

// Store implements gateway apis.
type Store struct {
	minio.GatewayUnsupported
	// vendor --> {pool_id --> client}
	Pools      map[string]map[string]sdk.Client
	VBucketMgr *vbucket.Mgr

	HTTPClient *http.Client
	Metrics    *minio.BackendMetrics
	debug      bool
}

// New xxx
func New(mgsAddr string) (*Store, error) {
	if len(mgsAddr) == 0 {
		return nil, errors.New("mgs addr empty")
	}
	metrics := minio.NewMetrics()
	t := &minio.MetricsTransport{
		Transport: minio.NewGatewayHTTPTransport(),
		Metrics:   metrics,
	}
	s := &Store{
		Metrics: metrics,
		HTTPClient: &http.Client{
			Transport: t,
		},
		//debug: true,
	}
	err := s.init(mgsAddr, t)
	return s, err
}

func (s *Store) init(mgsAddr string, transport http.RoundTripper) error {
	err := mgs.NewService(mgsAddr, serviceTimeout)
	if err != nil {
		return err
	}
	if s.VBucketMgr, err = vbucket.NewMgr(); err != nil {
		return err
	}
	if err = s.initClients(transport); err != nil {
		return err
	}
	return nil
}

func (s *Store) initClients(transport http.RoundTripper) error {
	s.Pools = make(map[string]map[string]sdk.Client)
	for k, v := range s.VBucketMgr.PoolMgr.Pools {
		_, exists := s.Pools[v.Vendor]
		if !exists {
			s.Pools[v.Vendor] = make(map[string]sdk.Client)
		}
		var (
			c   sdk.Client
			err error
		)
		switch v.Vendor {
		case pool.VendorAws:
			c, err = sdk.NewS3Client(v.Endpoint, v.Creds.AccessKey, v.Creds.SecretKey,
				transport)
		case pool.VendorCeph:
			c, err = sdk.NewS3Client(v.Endpoint, v.Creds.AccessKey, v.Creds.SecretKey,
				transport)
		case pool.VendorBaidu:
			c, err = sdk.NewBosClient(v.Endpoint, v.Creds.AccessKey, v.Creds.SecretKey)
		default:
			return fmt.Errorf("pool %q has unknown vendor type %q", k, v.Vendor)
		}
		if err != nil {
			return err
		}
		s.Pools[v.Vendor][k] = c
	}
	return nil
}

func (s *Store) getClient(pl *pool.Pool) sdk.Client {
	var client sdk.Client
	switch pl.Vendor {
	case pool.VendorAws:
		client = s.Pools[pool.VendorAws][pl.ID]
	case pool.VendorCeph:
		client = s.Pools[pool.VendorCeph][pl.ID]
	case pool.VendorBaidu:
		client = s.Pools[pool.VendorBaidu][pl.ID]
	default:
	}
	return client
}

// GetMetrics returns this gateway's metrics
func (s *Store) GetMetrics(ctx context.Context) (*minio.BackendMetrics, error) {
	return s.Metrics, nil
}

// Shutdown saves any gateway metadata to disk
// if necessary and reload upon next restart.
func (s *Store) Shutdown(ctx context.Context) error {
	return nil
}

// StorageInfo is not relevant to S3 backend.
func (s *Store) StorageInfo(ctx context.Context) (si minio.StorageInfo, _ []error) {
	// TODO(yangchunxin): check bucket exists using probe bucket on every pool
	si.Backend.Type = madmin.Gateway
	si.Backend.GatewayOnline = true
	return si, nil
}

/*****************************************Bucket Operations***************************************/

// MakeBucketWithLocation creates a new bucket on S3 backend.
func (s *Store) MakeBucketWithLocation(ctx context.Context, bucket string,
	opts minio.BucketOptions) error {
	if opts.LockEnabled || opts.VersioningEnabled {
		return minio.NotImplemented{}
	}
	if s3utils.CheckValidBucketName(bucket) != nil {
		return minio.BucketNameInvalid{Bucket: bucket}
	}
	err := s.VBucketMgr.MakeVBucket(bucket, opts.Location)
	if err != nil {
		return minio.ErrorRespToObjectError(err, bucket)
	}
	return nil
}

// GetBucketInfo gets bucket metadata.
func (s *Store) GetBucketInfo(ctx context.Context, bucket string) (bi minio.BucketInfo, e error) {
	vb, err := s.VBucketMgr.GetVBucketInfo(bucket)
	if err != nil {
		return bi, minio.ErrorRespToObjectError(err)
	}
	if vb == nil {
		return bi, minio.BucketNotFound{Bucket: bucket}
	}
	return minio.BucketInfo{
		Name:    vb.Name,
		Created: vb.CreatedTime,
	}, nil
}

// ListBuckets lists all buckets.
func (s *Store) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	vbs, err := s.VBucketMgr.ListVBuckets()
	if err != nil {
		return nil, minio.ErrorRespToObjectError(err)
	}
	bis := make([]minio.BucketInfo, len(vbs))
	for i, v := range vbs {
		bis[i] = minio.BucketInfo{
			Name:    v.Name,
			Created: v.CreatedTime,
		}
	}
	return bis, err
}

// DeleteBucket deletes one bucket.
func (s *Store) DeleteBucket(ctx context.Context, bucket string,
	opts minio.DeleteBucketOptions) error {
	err := s.VBucketMgr.DeleteVBucket(bucket)
	if err != nil {
		return minio.ErrorRespToObjectError(err, bucket)
	}
	return nil
}

/*****************************************Object Operations***************************************/

// ListObjects lists all blobs in S3 bucket filtered by prefix
func (s *Store) ListObjects(ctx context.Context, bucket string, prefix string, marker string,
	delimiter string, maxKeys int) (loi minio.ListObjectsInfo, e error) {
	// Validate bucket name.
	if err := s3utils.CheckValidBucketName(bucket); err != nil {
		return minio.ListObjectsInfo{}, err
	}
	// Validate object prefix.
	if err := s3utils.CheckValidObjectNamePrefix(prefix); err != nil {
		return minio.ListObjectsInfo{}, err
	}
	lop := &object.ListObjectsParam{
		VBucket:   bucket,
		Prefix:    prefix,
		Marker:    marker,
		Delimiter: delimiter,
		Limits:    maxKeys,
	}
	results, err := s.VBucketMgr.ListObjects(lop)
	//results, err := s.VBucketMgr.ListObjects(bucket, prefix, marker, delimiter, maxKeys)
	if err != nil {
		return loi, minio.ErrorRespToObjectError(err, bucket)
	}
	return results, nil
}

// ListObjectsV2 lists all blobs in S3 bucket filtered by prefix
func (s *Store) ListObjectsV2(ctx context.Context, bucket, prefix, continuationToken, delimiter string,
	maxKeys int, fetchOwner bool, startAfter string) (loi minio.ListObjectsV2Info, e error) {
	// FIXME(yangchunxin): why use v2?
	return minio.ListObjectsV2Info{}, minio.NotImplemented{}
}

// GetObjectNInfo - returns object info and locked object ReadCloser
func (s *Store) GetObjectNInfo(ctx context.Context, bucket, object string, rs *minio.HTTPRangeSpec, h http.Header, lockType minio.LockType,
	opts minio.ObjectOptions) (gr *minio.GetObjectReader, err error) {
	// 1. get object info
	oi, err := s.GetObjectInfo(ctx, bucket, object, opts)
	if err != nil {
		return nil, minio.ErrorRespToObjectError(err, bucket, object)
	}
	// 2. get pool and physical bucket
	pl, pBucket := s.VBucketMgr.GetObjectPoolAndBucket(bucket, object)
	if pl == nil || len(pBucket) == 0 {
		return nil, minio.ErrorRespToObjectError(errors.New("object not found"), bucket, object)
	}
	// 3. get client of pool
	client := s.getClient(pl)
	if client == nil {
		return nil, minio.ErrorRespToObjectError(errors.New("client of pool not found"), bucket, object)
	}
	// 4. get physical object name
	pObject := s.VBucketMgr.GetObjectKey(bucket, object)
	fn, off, length, err := minio.NewGetObjectReader(rs, oi, opts)
	if err != nil {
		return nil, minio.ErrorRespToObjectError(err, bucket, object)
	}
	pr, pw := io.Pipe()
	go func() {
		err := client.GetObject(ctx, pBucket, pObject, bucket, object, off, length, pw, oi.ETag, opts)
		pw.CloseWithError(err)
	}()
	// Setup cleanup function to cause the above go-routine to
	// exit in case of partial read
	pipeCloser := func() { pr.Close() }
	return fn(pr, h, pipeCloser)
}

// GetObjectInfo reads object info and replies back ObjectInfo
func (s *Store) GetObjectInfo(ctx context.Context, bucket string, object string,
	opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
	oi, err := s.VBucketMgr.GetObjectMeta(bucket, object)
	if err != nil {
		return minio.ObjectInfo{}, minio.ErrorRespToObjectError(err, bucket, object)
	}
	return oi, nil
}

// PutObject creates a new object with the incoming data,
func (s *Store) PutObject(ctx context.Context, bucket string, object string, r *minio.PutObjReader,
	opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
	// 1. get pool and physical bucket name
	pl, pBucket := s.VBucketMgr.GetPoolAndBucket(bucket, object)
	if pl == nil || len(pBucket) == 0 {
		return objInfo, minio.ErrorRespToObjectError(errors.New("physical bucket not found"), bucket, object)
	}
	// 2. get client of pool
	client := s.getClient(pl)
	if client == nil {
		return objInfo, minio.ErrorRespToObjectError(errors.New("client of pool not found"), bucket, object)
	}
	// 3. get physical object name
	pObject := s.VBucketMgr.GetObjectKey(bucket, object)
	// 4. put object
	oi, err := client.PutObject(ctx, pBucket, pObject, bucket, object, r, opts)
	if err != nil {
		return oi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	// 5. insert object meta
	err = s.VBucketMgr.PutObjectMeta(pl.ID, pBucket, oi)
	if err != nil {
		return oi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	return oi, nil
}

// CopyObject copies an object from source bucket to a destination bucket.
func (s *Store) CopyObject(ctx context.Context, srcBucket string, srcObject string, dstBucket string, dstObject string,
	srcInfo minio.ObjectInfo, srcOpts, dstOpts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
	/*
		if srcOpts.CheckPrecondFn != nil && srcOpts.CheckPrecondFn(srcInfo) {
			return minio.ObjectInfo{}, minio.PreConditionFailed{}
		}
		// Set this header such that following CopyObject() always sets the right metadata on the destination.
		// metadata input is already a trickled down value from interpreting x-amz-metadata-directive at
		// handler layer. So what we have right now is supposed to be applied on the destination object anyways.
		// So preserve it by adding "REPLACE" directive to save all the metadata set by CopyObject API.
		srcInfo.UserDefined["x-amz-metadata-directive"] = "REPLACE"
		srcInfo.UserDefined["x-amz-copy-source-if-match"] = srcInfo.ETag
		header := make(http.Header)
		if srcOpts.ServerSideEncryption != nil {
			encrypt.SSECopy(srcOpts.ServerSideEncryption).Marshal(header)
		}

		if dstOpts.ServerSideEncryption != nil {
			dstOpts.ServerSideEncryption.Marshal(header)
		}

		for k, v := range header {
			srcInfo.UserDefined[k] = v[0]
		}

		if _, err = l.Clients["test"].CopyObject(ctx, srcBucket, srcObject, dstBucket, dstObject, srcInfo.UserDefined, miniogo.CopySrcOptions{}, miniogo.PutObjectOptions{}); err != nil {
			return objInfo, minio.ErrorRespToObjectError(err, srcBucket, srcObject)
		}
		return l.GetObjectInfo(ctx, dstBucket, dstObject, dstOpts)
	*/
	return minio.ObjectInfo{}, minio.NotImplemented{}
}

// DeleteObject deletes a blob in bucket
func (s *Store) DeleteObject(ctx context.Context, bucket string, object string,
	opts minio.ObjectOptions) (minio.ObjectInfo, error) {
	// 1. get pool and physical bucket name
	pl, pBucket := s.VBucketMgr.GetObjectPoolAndBucket(bucket, object)
	if pl == nil || len(pBucket) == 0 {
		return minio.ObjectInfo{}, minio.ErrorRespToObjectError(errors.New("bucket not found"), bucket, object)
	}
	// 2. get client
	client := s.getClient(pl)
	if client == nil {
		return minio.ObjectInfo{}, minio.ErrorRespToObjectError(errors.New("client of pool not found"), bucket, object)
	}
	// 3. get physical object name
	pObject := s.VBucketMgr.GetObjectKey(bucket, object)
	oi, err := client.DeleteObject(ctx, pBucket, pObject, bucket, object, opts)
	if err != nil {
		return oi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	// 4. delete object meta
	err = s.VBucketMgr.DeleteObjectMeta(bucket, object)
	if err != nil {
		return oi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	return oi, nil
}

// DeleteObjects xxx
func (s *Store) DeleteObjects(ctx context.Context, bucket string, objects []minio.ObjectToDelete, opts minio.ObjectOptions) ([]minio.DeletedObject, []error) {
	errs := make([]error, len(objects))
	dobjects := make([]minio.DeletedObject, len(objects))
	for idx, object := range objects {
		_, errs[idx] = s.DeleteObject(ctx, bucket, object.ObjectName, opts)
		if errs[idx] == nil {
			dobjects[idx] = minio.DeletedObject{
				ObjectName: object.ObjectName,
			}
		}
	}
	return dobjects, errs
}

// ListMultipartUploads lists all multipart uploads.
func (s *Store) ListMultipartUploads(ctx context.Context, bucket string, prefix string, keyMarker string, uploadIDMarker string, delimiter string, maxUploads int) (lmi minio.ListMultipartsInfo, e error) {
	/*
		poolID := l.GetPool(bucket)
		result, err := l.Clients[poolID].ListMultipartUploads(ctx, bucket, prefix, keyMarker, uploadIDMarker, delimiter, maxUploads)
		if err != nil {
			return lmi, err
		}

		return minio.FromMinioClientListMultipartsInfo(result), nil
	*/
	fmt.Printf("----------------------ListMultipartUploads NotImplemented---------------------------\n")
	return minio.ListMultipartsInfo{}, minio.NotImplemented{}
}

// NewMultipartUpload upload object in multiple parts
func (s *Store) NewMultipartUpload(ctx context.Context, bucket string, object string, o minio.ObjectOptions) (uploadID string, err error) {
	// 1. get pool and physical bucket name
	pl, pBucket := s.VBucketMgr.GetPoolAndBucket(bucket, object)
	if pl == nil || len(pBucket) == 0 {
		return "", minio.ErrorRespToObjectError(errors.New("physical bucket not found"), bucket, object)
	}
	// 2. get client
	client := s.getClient(pl)
	if client == nil {
		return "", minio.ErrorRespToObjectError(errors.New("client of pool not found"), bucket, object)
	}
	// 3. get physical object name
	pObject := s.VBucketMgr.GetObjectKey(bucket, object)
	// 4. new multipart upload
	PhysicUploadID, err := client.NewMultipartUpload(ctx, pBucket, pObject, bucket, object, o)
	if err != nil {
		return "", minio.ErrorRespToObjectError(err, bucket, object)
	}
	// 5. create multipart meta
	uploadID, err = s.VBucketMgr.CreateMultipart(pBucket, pObject, bucket, object, PhysicUploadID)
	if err != nil {
		return "", minio.ErrorRespToObjectError(err, bucket, object)
	}
	return uploadID, nil
}

// PutObjectPart puts a part of object in bucket
func (s *Store) PutObjectPart(ctx context.Context, bucket string, object string, uploadID string, partID int,
	r *minio.PutObjReader, opts minio.ObjectOptions) (pi minio.PartInfo, e error) {
	// 1. get multipart meta
	mp, err := s.VBucketMgr.QueryMultipart(bucket, uploadID)
	if err != nil {
		return minio.PartInfo{}, minio.ErrorRespToObjectError(errors.New("bucket not found"), bucket, object)
	}
	// 2. get pool and physical bucket name
	pl, pBucket := s.VBucketMgr.GetPool(bucket), mp.PhysicalBucket
	if pl == nil || len(pBucket) == 0 {
		return minio.PartInfo{}, minio.ErrorRespToObjectError(errors.New("bucket not found"), bucket, object)
	}
	// 3. get client
	client := s.getClient(pl)
	if client == nil {
		return pi, minio.ErrorRespToObjectError(errors.New("client of pool not found"), bucket, object)
	}
	// 4. get physical object name
	pObject := s.VBucketMgr.GetObjectKey(bucket, object)
	// 5. put object part
	pi, err = client.PutObjectPart(ctx, pBucket, pObject, bucket, object, mp.PhysicalUploadID, partID, r, opts)
	if err != nil {
		return pi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	return pi, nil
}

// CopyObjectPart creates a part in a multipart upload by copying
// existing object or a part of it.
func (s *Store) CopyObjectPart(ctx context.Context, srcBucket, srcObject, destBucket, destObject, uploadID string,
	partID int, startOffset, length int64, srcInfo minio.ObjectInfo, srcOpts, dstOpts minio.ObjectOptions) (p minio.PartInfo, err error) {
	/*
		if srcOpts.CheckPrecondFn != nil && srcOpts.CheckPrecondFn(srcInfo) {
			return minio.PartInfo{}, minio.PreConditionFailed{}
		}
		srcInfo.UserDefined = map[string]string{
			"x-amz-copy-source-if-match": srcInfo.ETag,
		}
		header := make(http.Header)
		if srcOpts.ServerSideEncryption != nil {
			encrypt.SSECopy(srcOpts.ServerSideEncryption).Marshal(header)
		}

		if dstOpts.ServerSideEncryption != nil {
			dstOpts.ServerSideEncryption.Marshal(header)
		}
		for k, v := range header {
			srcInfo.UserDefined[k] = v[0]
		}

		completePart, err := l.Clients["test"].CopyObjectPart(ctx, srcBucket, srcObject, destBucket, destObject,
			uploadID, partID, startOffset, length, srcInfo.UserDefined)
		if err != nil {
			return p, minio.ErrorRespToObjectError(err, srcBucket, srcObject)
		}
		p.PartNumber = completePart.PartNumber
		p.ETag = completePart.ETag
		return p, nil
	*/
	fmt.Printf("---------------CopyObjectPart-------------\n")
	return minio.PartInfo{}, nil
}

// GetMultipartInfo returns multipart info of the uploadId of the object
func (s *Store) GetMultipartInfo(ctx context.Context, bucket, object, uploadID string,
	opts minio.ObjectOptions) (result minio.MultipartInfo, err error) {
	result.Bucket = bucket
	result.Object = object
	result.UploadID = uploadID
	return result, nil
}

// ListObjectParts returns all object parts for specified object in specified bucket
func (s *Store) ListObjectParts(ctx context.Context, bucket string, object string, uploadID string, partNumberMarker int,
	maxParts int, opts minio.ObjectOptions) (lpi minio.ListPartsInfo, e error) {
	// 1. get multipart meta
	mp, err := s.VBucketMgr.QueryMultipart(bucket, uploadID)
	if err != nil {
		return lpi, minio.ErrorRespToObjectError(errors.New("bucket not found"), bucket, object)
	}
	// 2. get pool and physical bucket
	pl, pBucket := s.VBucketMgr.GetPool(bucket), mp.PhysicalBucket
	if pl == nil || len(pBucket) == 0 {
		return lpi, minio.ErrorRespToObjectError(errors.New("bucket not found"), bucket, object)
	}
	// 3. get client
	client := s.getClient(pl)
	if client == nil {
		return lpi, minio.ErrorRespToObjectError(errors.New("client of pool not found"), bucket, object)
	}
	// 4. get physical object
	pObject := s.VBucketMgr.GetObjectKey(bucket, object)
	// 5. list object parts
	lpi, err = client.ListObjectParts(ctx, pBucket, pObject, bucket, object, mp.PhysicalUploadID, partNumberMarker, maxParts, opts)
	if err != nil {
		return lpi, err
	}
	return lpi, nil
}

// AbortMultipartUpload aborts a ongoing multipart upload
func (s *Store) AbortMultipartUpload(ctx context.Context, bucket string, object string, uploadID string, opts minio.ObjectOptions) error {
	// 1. get multipart meta
	mp, err := s.VBucketMgr.QueryMultipart(bucket, uploadID)
	if err != nil {
		return minio.ErrorRespToObjectError(errors.New("bucket not found"), bucket, object)
	}
	// 2. get pool and physical bucket
	pl, pBucket := s.VBucketMgr.GetPool(bucket), mp.PhysicalBucket
	if pl == nil || len(pBucket) == 0 {
		return minio.ErrorRespToObjectError(errors.New("bucket not found"), bucket, object)
	}
	// 3. get client
	client := s.getClient(pl)
	if client == nil {
		return minio.ErrorRespToObjectError(errors.New("client of pool not found"), bucket, object)
	}
	// 4. get physical object name
	pObject := s.VBucketMgr.GetObjectKey(bucket, object)
	// 5. abort multipart upload
	err = client.AbortMultipartUpload(ctx, pBucket, pObject, bucket, object, mp.PhysicalUploadID, opts)
	if err != nil {
		return minio.ErrorRespToObjectError(err, bucket, object)
	}
	// 6. delete multipart meta
	_ = s.VBucketMgr.DeleteMultipart(bucket, uploadID)
	return nil
}

// CompleteMultipartUpload completes ongoing multipart upload and finalizes object
func (s *Store) CompleteMultipartUpload(ctx context.Context, bucket string, object string, uploadID string,
	uploadedParts []minio.CompletePart, opts minio.ObjectOptions) (oi minio.ObjectInfo, e error) {
	// 1. get multipart meta
	mp, err := s.VBucketMgr.QueryMultipart(bucket, uploadID)
	if err != nil {
		return oi, minio.ErrorRespToObjectError(errors.New("bucket not found"), bucket, object)
	}
	// 2. get pool and physical bucket
	pl, pBucket := s.VBucketMgr.GetPool(bucket), mp.PhysicalBucket
	if pl == nil || len(pBucket) == 0 {
		return oi, minio.ErrorRespToObjectError(errors.New("bucket not found"), bucket, object)
	}
	// 3. get client
	client := s.getClient(pl)
	if client == nil {
		return oi, minio.ErrorRespToObjectError(errors.New("client of pool not found"), bucket, object)
	}
	// 4. get physical object
	pObject := s.VBucketMgr.GetObjectKey(bucket, object)
	// 5. complete multipart upload
	oi, err = client.CompleteMultipartUpload(ctx, pBucket, pObject, bucket, object, mp.PhysicalUploadID, uploadedParts, opts)
	if err != nil {
		return oi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	// 6. delete multipart meta
	_ = s.VBucketMgr.DeleteMultipart(bucket, uploadID)
	// 7. insert object meta
	err = s.VBucketMgr.PutObjectMeta(pl.ID, pBucket, oi)
	if err != nil {
		return oi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	return oi, nil
}

// SetBucketPolicy sets policy on bucket
func (s *Store) SetBucketPolicy(ctx context.Context, bucket string, bucketPolicy *policy.Policy) error {
	/*
		data, err := json.Marshal(bucketPolicy)
		if err != nil {
			// This should not happen.
			logger.LogIf(ctx, err)
			return minio.ErrorRespToObjectError(err, bucket)
		}

		poolID := l.GetPool(bucket)
		if err := l.Clients[poolID].SetBucketPolicy(ctx, bucket, string(data)); err != nil {
			return minio.ErrorRespToObjectError(err, bucket)
		}

		return nil
	*/
	return nil
}

// GetBucketPolicy will get policy on bucket
func (s *Store) GetBucketPolicy(ctx context.Context, bucket string) (*policy.Policy, error) {
	/*
		poolID := l.GetPool(bucket)
		data, err := l.Clients[poolID].GetBucketPolicy(ctx, bucket)
		if err != nil {
			return nil, minio.ErrorRespToObjectError(err, bucket)
		}

		bucketPolicy, err := policy.ParseConfig(strings.NewReader(data), bucket)
		return bucketPolicy, minio.ErrorRespToObjectError(err, bucket)
	*/
	return nil, nil
}

// DeleteBucketPolicy deletes all policies on bucket
func (s *Store) DeleteBucketPolicy(ctx context.Context, bucket string) error {
	/*
		poolID := l.GetPool(bucket)
		if err := l.Clients[poolID].SetBucketPolicy(ctx, bucket, ""); err != nil {
			return minio.ErrorRespToObjectError(err, bucket, "")
		}
		return nil
	*/
	return nil
}

// GetObjectTags gets the tags set on the object
func (s *Store) GetObjectTags(ctx context.Context, bucket string, object string, opts minio.ObjectOptions) (*tags.Tags, error) {
	/*
		var err error
		if _, err = l.GetObjectInfo(ctx, bucket, object, opts); err != nil {
			return nil, minio.ErrorRespToObjectError(err, bucket, object)
		}

		poolID := l.GetPool(bucket)
		t, err := l.Clients[poolID].GetObjectTagging(ctx, bucket, object, miniogo.GetObjectTaggingOptions{})
		if err != nil {
			return nil, minio.ErrorRespToObjectError(err, bucket, object)
		}

		return t, nil
	*/
	return nil, nil
}

// PutObjectTags attaches the tags to the object
func (s *Store) PutObjectTags(ctx context.Context, bucket, object string, tagStr string, opts minio.ObjectOptions) (minio.ObjectInfo, error) {
	/*
		tagObj, err := tags.Parse(tagStr, true)
		if err != nil {
			return minio.ObjectInfo{}, minio.ErrorRespToObjectError(err, bucket, object)
		}
		poolID := l.GetPool(bucket)
		if err = l.Clients[poolID].PutObjectTagging(ctx, bucket, object, tagObj, miniogo.PutObjectTaggingOptions{VersionID: opts.VersionID}); err != nil {
			return minio.ObjectInfo{}, minio.ErrorRespToObjectError(err, bucket, object)
		}

		objInfo, err := l.GetObjectInfo(ctx, bucket, object, opts)
		if err != nil {
			return minio.ObjectInfo{}, minio.ErrorRespToObjectError(err, bucket, object)
		}

		return objInfo, nil
	*/
	return minio.ObjectInfo{}, nil
}

// DeleteObjectTags removes the tags attached to the object
func (s *Store) DeleteObjectTags(ctx context.Context, bucket, object string, opts minio.ObjectOptions) (minio.ObjectInfo, error) {
	/*
		poolID := l.GetPool(bucket)
		if err := l.Clients[poolID].RemoveObjectTagging(ctx, bucket, object, miniogo.RemoveObjectTaggingOptions{}); err != nil {
			return minio.ObjectInfo{}, minio.ErrorRespToObjectError(err, bucket, object)
		}
		objInfo, err := l.GetObjectInfo(ctx, bucket, object, opts)
		if err != nil {
			return minio.ObjectInfo{}, minio.ErrorRespToObjectError(err, bucket, object)
		}

		return objInfo, nil
	*/
	return minio.ObjectInfo{}, nil
}

// IsCompressionSupported returns whether compression is applicable for this layer.
func (s *Store) IsCompressionSupported() bool {
	return false
}

// IsEncryptionSupported returns whether server side encryption is implemented for this layer.
func (s *Store) IsEncryptionSupported() bool {
	return minio.GlobalKMS != nil || minio.GlobalGatewaySSE.IsSet()
}

// IsTaggingSupported xxx
func (s *Store) IsTaggingSupported() bool {
	return true
}
