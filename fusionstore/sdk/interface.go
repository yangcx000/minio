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
	GetObject(ctx context.Context, pBucket string, pObject string, bucket string, object string, startOffset int64, length int64, writer io.Writer, etag string, o minio.ObjectOptions) error
}
