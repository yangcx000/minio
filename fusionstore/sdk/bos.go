/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package sdk

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/fusionstore/utils"
)

// Bos xxx
type Bos struct {
	client *bos.Client
}

// NewBosClient xxx
func NewBosClient(endpoint, accessKey, secretKey string) (*Bos, error) {
	client, err := bos.NewClient(accessKey, secretKey, endpoint)
	if err != nil {
		return nil, fmt.Errorf("couldn't new bos client %s", endpoint)
	}
	probeBucketName := utils.RandString(60, rand.NewSource(time.Now().UnixNano()), "probe-bucket-sign-")
	err = client.HeadBucket(probeBucketName)
	if err != nil {
		if realErr, ok := err.(*bce.BceServiceError); ok {
			if realErr.StatusCode == http.StatusForbidden || realErr.StatusCode == http.StatusNotFound {
				return &Bos{client: client}, nil
			}
			return nil, err
		}
	}
	return &Bos{client: client}, err
}

// PutObject xxx
func (b *Bos) PutObject(ctx context.Context, pBucket string, pObject string, bucket string, object string, r *minio.PutObjReader, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
	data := r.Reader
	args := &api.PutObjectArgs{
		UserMeta:      opts.UserDefined,
		ContentMD5:    data.MD5Base64String(),
		ContentLength: data.Size(),
		ContentSha256: data.SHA256HexString(),
	}
	etag, err := b.client.PutObjectFromStream(pBucket, pObject, data, args)
	if err != nil {
		return objInfo, minio.ErrorRespToObjectError(err, bucket, object)
	}
	objInfo = minio.ObjectInfo{
		Bucket: bucket,
		Name:   object,
		Size:   data.Size(),
		ETag:   etag,
	}
	return objInfo, nil
}

// GetObject xxx
func (b *Bos) GetObject(ctx context.Context, pBucket string, pObject string, bucket string, object string, startOffset int64, length int64, writer io.Writer, etag string, o minio.ObjectOptions) error {
	if length < 0 && length != -1 {
		return minio.ErrorRespToObjectError(minio.InvalidRange{}, bucket, object)
	}
	var (
		result *api.GetObjectResult
		err    error
	)
	start, end := startOffset, startOffset+length-1
	switch {
	case 0 < start && end == 0:
		result, err = b.client.GetObject(pBucket, pObject, nil, start)
	case 0 <= start && start <= end:
		result, err = b.client.GetObject(pBucket, pObject, nil, start, end)
	default:
		return minio.ErrorRespToObjectError(fmt.Errorf("Invalid range specified: start=%d end=%d", start, end), bucket, object)
	}
	if err != nil {
		return err
	}
	// XXX(yangchunxin): check etag
	reader := result.Body
	defer reader.Close()
	if _, err := io.Copy(writer, reader); err != nil {
		return minio.ErrorRespToObjectError(err, bucket, object)
	}
	return nil
}

// DeleteObject xxx
func (b *Bos) DeleteObject(ctx context.Context, pBucket string, pObject string, bucket string, object string, opts minio.ObjectOptions) (minio.ObjectInfo, error) {
	err := b.client.DeleteObject(pBucket, pObject)
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
func (b *Bos) NewMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string, object string, o minio.ObjectOptions) (uploadID string, err error) {
	return "", nil
}

// PutObjectPart xxx
func (b *Bos) PutObjectPart(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, partID int, r *minio.PutObjReader, opts minio.ObjectOptions) (pi minio.PartInfo, e error) {
	return minio.PartInfo{}, nil
}

// ListObjectParts xxx
func (b *Bos) ListObjectParts(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, partNumberMarker int, maxParts int, opts minio.ObjectOptions) (lpi minio.ListPartsInfo, e error) {
	return minio.ListPartsInfo{}, nil
}

// AbortMultipartUpload xxx
func (b *Bos) AbortMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, opts minio.ObjectOptions) error {
	return nil
}

// CompleteMultipartUpload xxx
func (b *Bos) CompleteMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, uploadedParts []minio.CompletePart, opts minio.ObjectOptions) (oi minio.ObjectInfo, e error) {
	return minio.ObjectInfo{}, nil
}
