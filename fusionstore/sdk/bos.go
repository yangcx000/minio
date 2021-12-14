/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package sdk

import (
	"context"
	"io"

	"github.com/baidubce/bce-sdk-go/services/bos"
	minio "github.com/minio/minio/cmd"
)

// Bos xxx
type Bos struct {
	client *bos.Client
}

// NewBosClient xxx
func NewBosClient(endpoint, accessKey, secretKey string) (*Bos, error) {
	var err error
	b := &Bos{}
	b.client, err = bos.NewClient(accessKey, secretKey, endpoint)
	return b, err
}

// PutObject xxx
func (b *Bos) PutObject(ctx context.Context, pBucket string, pObject string, bucket string, object string, r *minio.PutObjReader, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
	return
}

// GetObject xxx
func (b *Bos) GetObject(ctx context.Context, pBucket string, pObject string, bucket string, object string, startOffset int64, length int64, writer io.Writer, etag string, o minio.ObjectOptions) error {
	return nil
}

// DeleteObject xxx
func (b *Bos) DeleteObject(ctx context.Context, pBucket string, pObject string, bucket string, object string, opts minio.ObjectOptions) (minio.ObjectInfo, error) {
	return minio.ObjectInfo{}, nil
}
