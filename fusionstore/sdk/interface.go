/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package sdk

import (
	"context"
	"io"

	minio "github.com/minio/minio/cmd"
)

// Client xxx
type Client interface {
	PutObject(ctx context.Context, pBucket string, pObject string, bucket string, object string, r *minio.PutObjReader, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error)
	DeleteObject(ctx context.Context, pBucket string, pObject string, bucket string, object string, opts minio.ObjectOptions) (minio.ObjectInfo, error)
	GetObject(ctx context.Context, pBucket string, pObject string, bucket string, object string, rs *minio.HTTPRangeSpec, writer io.Writer, etag string, o minio.ObjectOptions) error
	NewMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string, object string, o minio.ObjectOptions) (uploadID string, err error)
	PutObjectPart(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, partID int, r *minio.PutObjReader, opts minio.ObjectOptions) (pi minio.PartInfo, e error)
	ListObjectParts(ctx context.Context, pUploadID string, pBucket string, pObject string, bucket string, object string, uploadID string, partNumberMarker int, maxParts int, opts minio.ObjectOptions) (lpi minio.ListPartsInfo, e error)
	AbortMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, opts minio.ObjectOptions) error
	CompleteMultipartUpload(ctx context.Context, pBucket string, pObject string, bucket string, object string, uploadID string, uploadedParts []minio.CompletePart, opts minio.ObjectOptions) (oi minio.ObjectInfo, e error)
}
