package fusionstore

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/minio/madmin-go"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"github.com/minio/minio-go/v7/pkg/tags"
	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/fusionstore/mgs"
	"github.com/minio/minio/fusionstore/vbucket"
	"github.com/minio/pkg/bucket/policy"
)

const (
	serviceTimeout = 10
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz01234569"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// randString generates random names and prepends them with a known prefix.
func randString(n int, src rand.Source, prefix string) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return prefix + string(b[0:30-len(prefix)])
}

// Store xxx
type Store struct {
	minio.GatewayUnsupported
	Pools      map[string]*miniogo.Core
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
		debug: true,
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
	if err = s.setupClients(transport); err != nil {
		return err
	}
	return nil
}

func (s *Store) setupClients(transport http.RoundTripper) error {
	s.Pools = make(map[string]*miniogo.Core)
	for k, v := range s.VBucketMgr.PoolMgr.Pools {
		cred := madmin.Credentials{
			AccessKey: v.Creds.AccessKey,
			SecretKey: v.Creds.SecretKey,
		}
		core, err := s.newCore("http://"+v.Endpoint, cred, transport)
		if err != nil {
			return err
		}
		s.Pools[k] = core
	}
	return nil
}

func (s *Store) newCore(endpoint string, creds madmin.Credentials, transport http.RoundTripper) (*miniogo.Core, error) {
	if len(endpoint) == 0 {
		return nil, errors.New("endpoint empty")
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	// Override default params if the host is provided
	endpoint, secure, err := minio.ParseGatewayEndpoint(endpoint)
	if err != nil {
		return nil, err
	}
	optionsStaticCreds := &miniogo.Options{
		Creds:        credentials.NewStaticV4(creds.AccessKey, creds.SecretKey, creds.SessionToken),
		Secure:       secure,
		Region:       s3utils.GetRegionFromURL(*u),
		BucketLookup: miniogo.BucketLookupAuto,
		Transport:    transport,
	}
	client, err := miniogo.New(endpoint, optionsStaticCreds)
	if err != nil {
		return nil, err
	}
	if s.debug {
		client.TraceOn(os.Stderr)
	}
	probeBucketName := randString(60, rand.NewSource(time.Now().UnixNano()), "probe-bucket-sign-")
	if _, err = client.BucketExists(context.Background(), probeBucketName); err != nil {
		switch miniogo.ToErrorResponse(err).Code {
		case "AccessDenied":
			return &miniogo.Core{Client: client}, nil
		default:
			return nil, err
		}
	}
	return &miniogo.Core{Client: client}, nil
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
	// TODO(yangchunxin): add later
	si.Backend.Type = madmin.Gateway
	si.Backend.GatewayOnline = true
	return si, nil
}

// MakeBucketWithLocation creates a new container on S3 backend.
func (s *Store) MakeBucketWithLocation(ctx context.Context, bucket string, opts minio.BucketOptions) error {
	if opts.LockEnabled || opts.VersioningEnabled {
		return minio.NotImplemented{}
	}
	// Verify if bucket name is valid.
	// We are using a separate helper function here to validate bucket
	// names instead of IsValidBucketName() because there is a possibility
	// that certains users might have buckets which are non-DNS compliant
	// in us-east-1 and we might severely restrict them by not allowing
	// access to these buckets.
	// Ref - http://docs.aws.amazon.com/AmazonS3/latest/dev/BucketRestrictions.html
	if s3utils.CheckValidBucketName(bucket) != nil {
		return minio.BucketNameInvalid{Bucket: bucket}
	}
	err := s.VBucketMgr.MakeBucket(bucket, opts.Location)
	if err != nil {
		return minio.ErrorRespToObjectError(err, bucket)
	}
	return err
}

// GetBucketInfo gets bucket metadata..
func (s *Store) GetBucketInfo(ctx context.Context, bucket string) (bi minio.BucketInfo, e error) {
	vb, err := s.VBucketMgr.GetBucketInfo(bucket)
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

// ListBuckets lists all buckets
func (s *Store) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	vbs, err := s.VBucketMgr.ListBuckets()
	if err != nil {
		return nil, minio.ErrorRespToObjectError(err)
	}
	b := make([]minio.BucketInfo, len(vbs))
	for i, v := range vbs {
		b[i] = minio.BucketInfo{
			Name:    v.Name,
			Created: v.CreatedTime,
		}
	}
	return b, err
}

// DeleteBucket deletes a bucket
func (s *Store) DeleteBucket(ctx context.Context, bucket string, opts minio.DeleteBucketOptions) error {
	err := s.VBucketMgr.DeleteBucket(bucket)
	if err != nil {
		return minio.ErrorRespToObjectError(err, bucket)
	}
	return nil
}

// ListObjects lists all blobs in S3 bucket filtered by prefix
func (s *Store) ListObjects(ctx context.Context, bucket string, prefix string, marker string, delimiter string, maxKeys int) (loi minio.ListObjectsInfo, e error) {
	/*
		poolID := l.GetPool(bucket)
		result, err := l.Clients[poolID].ListObjects(bucket, prefix, marker, delimiter, maxKeys)
		if err != nil {
			return loi, minio.ErrorRespToObjectError(err, bucket)
		}

		return minio.FromMinioClientListBucketResult(bucket, result), nil
	*/
	return minio.ListObjectsInfo{}, nil
}

// ListObjectsV2 lists all blobs in S3 bucket filtered by prefix
func (s *Store) ListObjectsV2(ctx context.Context, bucket, prefix, continuationToken, delimiter string, maxKeys int, fetchOwner bool, startAfter string) (loi minio.ListObjectsV2Info, e error) {
	/*
		poolID := l.GetPool(bucket)
		result, err := l.Clients[poolID].ListObjectsV2(bucket, prefix, startAfter, continuationToken, delimiter, maxKeys)
		if err != nil {
			return loi, minio.ErrorRespToObjectError(err, bucket)
		}

		return minio.FromMinioClientListBucketV2Result(bucket, result), nil
	*/
	return minio.ListObjectsV2Info{}, nil
}

// GetObjectNInfo - returns object info and locked object ReadCloser
func (s *Store) GetObjectNInfo(ctx context.Context, bucket, object string, rs *minio.HTTPRangeSpec, h http.Header, lockType minio.LockType, opts minio.ObjectOptions) (gr *minio.GetObjectReader, err error) {
	/*
		var objInfo minio.ObjectInfo
		objInfo, err = l.GetObjectInfo(ctx, bucket, object, opts)
		if err != nil {
			return nil, minio.ErrorRespToObjectError(err, bucket, object)
		}

		fn, off, length, err := minio.NewGetObjectReader(rs, objInfo, opts)
		if err != nil {
			return nil, minio.ErrorRespToObjectError(err, bucket, object)
		}

		pr, pw := io.Pipe()
		go func() {
			err := l.getObject(ctx, bucket, object, off, length, pw, objInfo.ETag, opts)
			pw.CloseWithError(err)
		}()

		// Setup cleanup function to cause the above go-routine to
		// exit in case of partial read
		pipeCloser := func() { pr.Close() }
		return fn(pr, h, pipeCloser)
	*/
	return nil, nil
}

// GetObject reads an object from S3. Supports additional
// parameters like offset and length which are synonymous with
// HTTP Range requests.
//
// startOffset indicates the starting read location of the object.
// length indicates the total length of the object.
func (s *Store) getObject(ctx context.Context, bucket string, key string, startOffset int64, length int64, writer io.Writer, etag string, o minio.ObjectOptions) error {
	/*
		if length < 0 && length != -1 {
			return minio.ErrorRespToObjectError(minio.InvalidRange{}, bucket, key)
		}

		opts := miniogo.GetObjectOptions{}
		opts.ServerSideEncryption = o.ServerSideEncryption

		if startOffset >= 0 && length >= 0 {
			if err := opts.SetRange(startOffset, startOffset+length-1); err != nil {
				return minio.ErrorRespToObjectError(err, bucket, key)
			}
		}

		if etag != "" {
			opts.SetMatchETag(etag)
		}

		poolID := l.GetPool(bucket)
		object, _, _, err := l.Clients[poolID].GetObject(ctx, bucket, key, opts)
		if err != nil {
			return minio.ErrorRespToObjectError(err, bucket, key)
		}
		defer object.Close()
		if _, err := io.Copy(writer, object); err != nil {
			return minio.ErrorRespToObjectError(err, bucket, key)
		}
		return nil
	*/
	return nil
}

// GetObjectInfo reads object info and replies back ObjectInfo
func (s *Store) GetObjectInfo(ctx context.Context, bucket string, object string, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
	/*
		poolID := l.GetPool(bucket)
		oi, err := l.Clients[poolID].StatObject(ctx, bucket, object, miniogo.StatObjectOptions{
			ServerSideEncryption: opts.ServerSideEncryption,
		})
		if err != nil {
			return minio.ObjectInfo{}, minio.ErrorRespToObjectError(err, bucket, object)
		}

		return minio.FromMinioClientObjectInfo(bucket, oi), nil
	*/
	return minio.ObjectInfo{}, nil
}

// PutObject creates a new object with the incoming data,
func (s *Store) PutObject(ctx context.Context, bucket string, object string, r *minio.PutObjReader, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
	/*
		data := r.Reader
		var tagMap map[string]string
		if tagstr, ok := opts.UserDefined[xhttp.AmzObjectTagging]; ok && tagstr != "" {
			tagObj, err := tags.ParseObjectTags(tagstr)
			if err != nil {
				return objInfo, minio.ErrorRespToObjectError(err, bucket, object)
			}
			tagMap = tagObj.ToMap()
			delete(opts.UserDefined, xhttp.AmzObjectTagging)
		}
		putOpts := miniogo.PutObjectOptions{
			UserMetadata:         opts.UserDefined,
			ServerSideEncryption: opts.ServerSideEncryption,
			UserTags:             tagMap,
			// Content-Md5 is needed for buckets with object locking,
			// instead of spending an extra API call to detect this
			// we can set md5sum to be calculated always.
			SendContentMd5: true,
		}
		poolID := l.GetPool(bucket)
		ui, err := l.Clients[poolID].PutObject(ctx, bucket, object, data, data.Size(), data.MD5Base64String(), data.SHA256HexString(), putOpts)
		if err != nil {
			return objInfo, minio.ErrorRespToObjectError(err, bucket, object)
		}
		// On success, populate the key & metadata so they are present in the notification
		oi := miniogo.ObjectInfo{
			ETag:     ui.ETag,
			Size:     ui.Size,
			Key:      object,
			Metadata: minio.ToMinioClientObjectInfoMetadata(opts.UserDefined),
		}

		return minio.FromMinioClientObjectInfo(bucket, oi), nil
	*/
	return minio.ObjectInfo{}, nil
}

// CopyObject copies an object from source bucket to a destination bucket.
func (s *Store) CopyObject(ctx context.Context, srcBucket string, srcObject string, dstBucket string, dstObject string, srcInfo minio.ObjectInfo, srcOpts, dstOpts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
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
	return minio.ObjectInfo{}, nil
}

// DeleteObject deletes a blob in bucket
func (s *Store) DeleteObject(ctx context.Context, bucket string, object string, opts minio.ObjectOptions) (minio.ObjectInfo, error) {
	/*
		poolID := l.GetPool(bucket)
		err := l.Clients[poolID].RemoveObject(ctx, bucket, object, miniogo.RemoveObjectOptions{})
		if err != nil {
			return minio.ObjectInfo{}, minio.ErrorRespToObjectError(err, bucket, object)
		}

		return minio.ObjectInfo{
			Bucket: bucket,
			Name:   object,
		}, nil
	*/
	return minio.ObjectInfo{}, nil
}

// DeleteObjects xxx
func (s *Store) DeleteObjects(ctx context.Context, bucket string, objects []minio.ObjectToDelete, opts minio.ObjectOptions) ([]minio.DeletedObject, []error) {
	/*
		errs := make([]error, len(objects))
		dobjects := make([]minio.DeletedObject, len(objects))
		for idx, object := range objects {
			_, errs[idx] = l.DeleteObject(ctx, bucket, object.ObjectName, opts)
			if errs[idx] == nil {
				dobjects[idx] = minio.DeletedObject{
					ObjectName: object.ObjectName,
				}
			}
		}
		return dobjects, errs
	*/
	return nil, nil
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
	return minio.ListMultipartsInfo{}, nil
}

// NewMultipartUpload upload object in multiple parts
func (s *Store) NewMultipartUpload(ctx context.Context, bucket string, object string, o minio.ObjectOptions) (uploadID string, err error) {
	/*
		var tagMap map[string]string
		if tagStr, ok := o.UserDefined[xhttp.AmzObjectTagging]; ok {
			tagObj, err := tags.Parse(tagStr, true)
			if err != nil {
				return uploadID, minio.ErrorRespToObjectError(err, bucket, object)
			}
			tagMap = tagObj.ToMap()
			delete(o.UserDefined, xhttp.AmzObjectTagging)
		}
		// Create PutObject options
		opts := miniogo.PutObjectOptions{
			UserMetadata:         o.UserDefined,
			ServerSideEncryption: o.ServerSideEncryption,
			UserTags:             tagMap,
		}
		poolID := l.GetPool(bucket)
		uploadID, err = l.Clients[poolID].NewMultipartUpload(ctx, bucket, object, opts)
		if err != nil {
			return uploadID, minio.ErrorRespToObjectError(err, bucket, object)
		}
		return uploadID, nil
	*/
	return "", nil
}

// PutObjectPart puts a part of object in bucket
func (s *Store) PutObjectPart(ctx context.Context, bucket string, object string, uploadID string, partID int, r *minio.PutObjReader, opts minio.ObjectOptions) (pi minio.PartInfo, e error) {
	/*
		data := r.Reader
		poolID := l.GetPool(bucket)
		info, err := l.Clients[poolID].PutObjectPart(ctx, bucket, object, uploadID, partID, data, data.Size(), data.MD5Base64String(), data.SHA256HexString(), opts.ServerSideEncryption)
		if err != nil {
			return pi, minio.ErrorRespToObjectError(err, bucket, object)
		}

		return minio.FromMinioClientObjectPart(info), nil
	*/
	return minio.PartInfo{}, nil
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
	return minio.PartInfo{}, nil
}

// GetMultipartInfo returns multipart info of the uploadId of the object
func (s *Store) GetMultipartInfo(ctx context.Context, bucket, object, uploadID string, opts minio.ObjectOptions) (result minio.MultipartInfo, err error) {
	/*
		result.Bucket = bucket
		result.Object = object
		result.UploadID = uploadID
		return result, nil
	*/
	return minio.MultipartInfo{}, nil
}

// ListObjectParts returns all object parts for specified object in specified bucket
func (s *Store) ListObjectParts(ctx context.Context, bucket string, object string, uploadID string, partNumberMarker int, maxParts int, opts minio.ObjectOptions) (lpi minio.ListPartsInfo, e error) {
	/*
		poolID := l.GetPool(bucket)
		result, err := l.Clients[poolID].ListObjectParts(ctx, bucket, object, uploadID, partNumberMarker, maxParts)
		if err != nil {
			return lpi, err
		}
		lpi = minio.FromMinioClientListPartsInfo(result)
		if lpi.IsTruncated && maxParts > len(lpi.Parts) {
			partNumberMarker = lpi.NextPartNumberMarker
			for {
				poolID := l.GetPool(bucket)
				result, err = l.Clients[poolID].ListObjectParts(ctx, bucket, object, uploadID, partNumberMarker, maxParts)
				if err != nil {
					return lpi, err
				}

				nlpi := minio.FromMinioClientListPartsInfo(result)

				partNumberMarker = nlpi.NextPartNumberMarker

				lpi.Parts = append(lpi.Parts, nlpi.Parts...)
				if !nlpi.IsTruncated {
					break
				}
			}
		}
		return lpi, nil
	*/
	return minio.ListPartsInfo{}, nil
}

// AbortMultipartUpload aborts a ongoing multipart upload
func (s *Store) AbortMultipartUpload(ctx context.Context, bucket string, object string, uploadID string, opts minio.ObjectOptions) error {
	/*
		poolID := l.GetPool(bucket)
		err := l.Clients[poolID].AbortMultipartUpload(ctx, bucket, object, uploadID)
		return minio.ErrorRespToObjectError(err, bucket, object)
	*/
	return nil
}

// CompleteMultipartUpload completes ongoing multipart upload and finalizes object
func (s *Store) CompleteMultipartUpload(ctx context.Context, bucket string, object string, uploadID string, uploadedParts []minio.CompletePart, opts minio.ObjectOptions) (oi minio.ObjectInfo, e error) {
	/*
		poolID := l.GetPool(bucket)
		etag, err := l.Clients[poolID].CompleteMultipartUpload(ctx, bucket, object, uploadID, minio.ToMinioClientCompleteParts(uploadedParts), miniogo.PutObjectOptions{})
		if err != nil {
			return oi, minio.ErrorRespToObjectError(err, bucket, object)
		}

		return minio.ObjectInfo{Bucket: bucket, Name: object, ETag: strings.Trim(etag, "\"")}, nil
	*/
	return minio.ObjectInfo{}, nil
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
