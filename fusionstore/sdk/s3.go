/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package sdk

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/minio/madmin-go"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/s3utils"
	"github.com/minio/minio-go/v7/pkg/tags"
	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/fusionstore/utils"
	xhttp "github.com/minio/minio/internal/http"
)

// S3 xxx
type S3 struct {
	client *miniogo.Core
}

// NewS3Client xxx
func NewS3Client(endpoint, accessKey, secretKey string, transport http.RoundTripper) (*S3, error) {
	cred := madmin.Credentials{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
	core, err := newCore(endpoint, cred, transport)
	if err != nil {
		return nil, err
	}
	return &S3{client: core}, nil
}

func newCore(endpoint string, creds madmin.Credentials, transport http.RoundTripper) (*miniogo.Core, error) {
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
	//if s.debug {
	//	client.TraceOn(os.Stderr)
	//}
	probeBucketName := utils.RandString(60, rand.NewSource(time.Now().UnixNano()), "probe-bucket-sign-")
	if _, err = client.BucketExists(context.Background(), probeBucketName); err != nil {
		switch miniogo.ToErrorResponse(err).Code {
		case "AccessDenied":
			// this is a good error means backend is reachable
			// and credentials are valid but credentials don't
			// have access to 'probeBucketName' which is harmless.
			return &miniogo.Core{Client: client}, nil
		default:
			return nil, err
		}
	}
	return &miniogo.Core{Client: client}, nil
}

// PutObject xxx
func (s *S3) PutObject(ctx context.Context, pBucket string, pObject string, bucket string, object string,
	r *minio.PutObjReader, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
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
	ui, err := s.client.PutObject(ctx, pBucket, pObject, data, data.Size(), data.MD5Base64String(), data.SHA256HexString(), putOpts)
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
	objInfo = minio.FromMinioClientObjectInfo(bucket, oi)
	return objInfo, nil
}

// GetObject xxx
func (s *S3) GetObject(ctx context.Context, pBucket string, pObject string, bucket string, object string,
	startOffset int64, length int64, writer io.Writer, etag string, o minio.ObjectOptions) error {
	if length < 0 && length != -1 {
		return minio.ErrorRespToObjectError(minio.InvalidRange{}, bucket, object)
	}
	opts := miniogo.GetObjectOptions{}
	opts.ServerSideEncryption = o.ServerSideEncryption
	if startOffset >= 0 && length >= 0 {
		if err := opts.SetRange(startOffset, startOffset+length-1); err != nil {
			return minio.ErrorRespToObjectError(err, bucket, object)
		}
	}
	if etag != "" {
		opts.SetMatchETag(etag)
	}
	reader, _, _, err := s.client.GetObject(ctx, pBucket, pObject, opts)
	if err != nil {
		return minio.ErrorRespToObjectError(err, bucket, object)
	}
	defer reader.Close()
	if _, err := io.Copy(writer, reader); err != nil {
		return minio.ErrorRespToObjectError(err, bucket, object)
	}
	return nil
}

// DeleteObject xxx
func (s *S3) DeleteObject(ctx context.Context, pBucket string, pObject string, bucket string,
	object string, opts minio.ObjectOptions) (minio.ObjectInfo, error) {
	err := s.client.RemoveObject(ctx, pBucket, pObject, miniogo.RemoveObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, minio.ErrorRespToObjectError(err, bucket, object)
	}
	objInfo := minio.ObjectInfo{
		Bucket: bucket,
		Name:   object,
	}
	return objInfo, nil
}

// NewMultipartUpload xxx
func (s *S3) NewMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string,
	object string, o minio.ObjectOptions) (uploadID string, err error) {
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
	uploadID, err = s.client.NewMultipartUpload(ctx, pBucket, pObject, opts)
	if err != nil {
		return uploadID, minio.ErrorRespToObjectError(err, bucket, object)
	}
	return uploadID, nil
}

// PutObjectPart xxx
func (s *S3) PutObjectPart(ctx context.Context, pBucket string, pObject string, bucket string, object string,
	uploadID string, partID int, r *minio.PutObjReader, opts minio.ObjectOptions) (pi minio.PartInfo, e error) {
	data := r.Reader
	info, err := s.client.PutObjectPart(ctx, pBucket, pObject, uploadID, partID, data, data.Size(),
		data.MD5Base64String(), data.SHA256HexString(), opts.ServerSideEncryption)
	if err != nil {
		return pi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	return minio.FromMinioClientObjectPart(info), nil
}

// ListObjectParts xxx
func (s *S3) ListObjectParts(ctx context.Context, pBucket string, pObject string, bucket string, object string,
	uploadID string, partNumberMarker int, maxParts int, opts minio.ObjectOptions) (lpi minio.ListPartsInfo, e error) {
	result, err := s.client.ListObjectParts(ctx, pBucket, pObject, uploadID, partNumberMarker, maxParts)
	if err != nil {
		return lpi, err
	}
	lpi = minio.FromMinioClientListPartsInfo(result)
	if lpi.IsTruncated && maxParts > len(lpi.Parts) {
		partNumberMarker = lpi.NextPartNumberMarker
		for {
			result, err = s.client.ListObjectParts(ctx, pBucket, pObject, uploadID, partNumberMarker, maxParts)
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
}

// AbortMultipartUpload xxx
func (s *S3) AbortMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string,
	object string, uploadID string, opts minio.ObjectOptions) error {
	err := s.client.AbortMultipartUpload(ctx, pBucket, pObject, uploadID)
	return minio.ErrorRespToObjectError(err, bucket, object)
}

// CompleteMultipartUpload xxx
func (s *S3) CompleteMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string, object string,
	uploadID string, uploadedParts []minio.CompletePart, opts minio.ObjectOptions) (oi minio.ObjectInfo, e error) {
	etag, err := s.client.CompleteMultipartUpload(ctx, pBucket, pObject, uploadID,
		minio.ToMinioClientCompleteParts(uploadedParts), miniogo.PutObjectOptions{})
	if err != nil {
		return oi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	cmdObjInfo, err := s.client.StatObject(ctx, pBucket, pObject, miniogo.StatObjectOptions{})
	if err != nil {
		return oi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	oi = minio.FromMinioClientObjectInfo(bucket, cmdObjInfo)
	// replace pobject with object
	oi.Name = object
	cmuEtag := strings.Trim(etag, "\"")
	if cmuEtag != oi.ETag {
		fmt.Printf("****************CompleteMultipartUpload Etag %s != objInfo Etag %s\n", cmuEtag, oi.ETag)
	}
	return oi, nil
}
