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
	"strconv"
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
	result, err := b.client.InitiateMultipartUpload(pBucket, pObject, "", nil)
	if err != nil {
		return uploadID, minio.ErrorRespToObjectError(err, bucket, object)
	}
	return result.UploadId, nil
}

// PutObjectPart xxx
func (b *Bos) PutObjectPart(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, partID int, r *minio.PutObjReader, opts minio.ObjectOptions) (pi minio.PartInfo, e error) {
	data := r.Reader
	args := &api.UploadPartArgs{
		ContentMD5:    data.MD5Base64String(),
		ContentSha256: data.SHA256HexString(),
	}
	body, err := bce.NewBodyFromSizedReader(data, data.Size())
	if err != nil {
		return pi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	etag, err := b.client.UploadPart(pBucket, pObject, uploadID, partID, body, args)
	if err != nil {
		return pi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	pi = minio.PartInfo{
		PartNumber:   partID,
		LastModified: time.Now(),
		ETag:         etag,
		Size:         data.Size(),
	}
	return pi, nil
}

func toParts(lpts []api.ListPartType) []minio.PartInfo {
	pis := make([]minio.PartInfo, len(lpts))
	for i, v := range lpts {
		pis[i] = minio.PartInfo{
			PartNumber: v.PartNumber,
			//LastModified: v.LastModified,
			ETag: v.ETag,
			Size: int64(v.Size),
		}
	}
	return pis
}

func resultToPartsInfo(lpr *api.ListPartsResult) minio.ListPartsInfo {
	lpi := minio.ListPartsInfo{
		Bucket:               lpr.Bucket,
		Object:               lpr.Key,
		UploadID:             lpr.UploadId,
		StorageClass:         lpr.StorageClass,
		PartNumberMarker:     lpr.PartNumberMarker,
		NextPartNumberMarker: lpr.NextPartNumberMarker,
		MaxParts:             lpr.MaxParts,
		IsTruncated:          lpr.IsTruncated,
		Parts:                toParts(lpr.Parts),
	}
	return lpi
}

// ListObjectParts xxx
func (b *Bos) ListObjectParts(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, partNumberMarker int, maxParts int, opts minio.ObjectOptions) (lpi minio.ListPartsInfo, e error) {
	args := &api.ListPartsArgs{
		MaxParts:         maxParts,
		PartNumberMarker: strconv.Itoa(partNumberMarker),
	}
	result, err := b.client.ListParts(pBucket, pObject, uploadID, args)
	if err != nil {
		return lpi, err
	}
	lpi = resultToPartsInfo(result)
	if lpi.IsTruncated && maxParts > len(lpi.Parts) {
		partNumberMarker = lpi.NextPartNumberMarker
		for {
			args.PartNumberMarker = strconv.Itoa(partNumberMarker)
			result, err = b.client.ListParts(pBucket, pObject, uploadID, args)
			if err != nil {
				return lpi, err
			}
			nlpi := resultToPartsInfo(result)
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
func (b *Bos) AbortMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, opts minio.ObjectOptions) error {
	err := b.client.AbortMultipartUpload(pBucket, pObject, uploadID)
	return minio.ErrorRespToObjectError(err, bucket, object)
}

func toCompleteMultipartUploadArgs(parts []minio.CompletePart) *api.CompleteMultipartUploadArgs {
	args := &api.CompleteMultipartUploadArgs{
		Parts: make([]api.UploadInfoType, len(parts)),
	}
	for i, v := range parts {
		args.Parts[i] = api.UploadInfoType{
			PartNumber: v.PartNumber,
			ETag:       v.ETag,
		}
	}
	return args
}

// CompleteMultipartUpload xxx
func (b *Bos) CompleteMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, uploadedParts []minio.CompletePart, opts minio.ObjectOptions) (oi minio.ObjectInfo, e error) {
	args := toCompleteMultipartUploadArgs(uploadedParts)
	result, err := b.client.CompleteMultipartUploadFromStruct(pBucket, pObject, uploadID, args)
	if err != nil {
		return oi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	meta, err := b.client.GetObjectMeta(pBucket, pObject)
	if err != nil {
		return oi, minio.ErrorRespToObjectError(err, bucket, object)
	}
	if result.ETag != meta.ETag {
		return oi, minio.ErrorRespToObjectError(errors.New("etag not equal"), bucket, object)
	}
	oi = minio.ObjectInfo{
		Bucket: bucket,
		Name:   object,
		Size:   meta.ContentLength,
		ETag:   meta.ETag,
	}
	return oi, nil
}
