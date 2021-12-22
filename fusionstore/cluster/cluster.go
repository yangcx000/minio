/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package cluster

import (
	"fmt"
	"net/http"

	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/fusionstore/mgs"
	"github.com/minio/minio/fusionstore/multipart"
	"github.com/minio/minio/fusionstore/object"
	"github.com/minio/minio/fusionstore/pool"
	"github.com/minio/minio/fusionstore/sdk"
	"github.com/minio/minio/fusionstore/utils"
	"github.com/minio/minio/fusionstore/vbucket"
)

const (
	serviceTimeout = 10
)

// Cluster xxx
type Cluster struct {
	PoolMgr    *pool.Mgr
	VBucketMgr *vbucket.Mgr
	// vendor --> {pool_id --> client}
	Pools map[string]map[string]sdk.Client
}

// New xxx
func New(mgsAddr string, transport http.RoundTripper) (c *Cluster, err error) {
	if err = mgs.NewService(mgsAddr, serviceTimeout); err != nil {
		return nil, err
	}
	c = &Cluster{}
	if c.PoolMgr, err = pool.NewMgr(); err != nil {
		return nil, err
	}
	if c.VBucketMgr, err = vbucket.NewMgr(); err != nil {
		return nil, err
	}
	if err = c.initClients(transport); err != nil {
		return nil, err
	}
	return c, nil
}

// Shutdown xxx
func (c *Cluster) Shutdown() {
	mgs.GlobalService.Close()
	c.VBucketMgr.Shutdown()
}

func (c *Cluster) initClients(transport http.RoundTripper) (err error) {
	var (
		client sdk.Client
	)
	c.Pools = make(map[string]map[string]sdk.Client)
	for k, v := range c.PoolMgr.GetPoolMap() {
		_, exists := c.Pools[v.Vendor]
		if !exists {
			c.Pools[v.Vendor] = make(map[string]sdk.Client)
		}
		switch v.Vendor {
		case pool.VendorAws:
			client, err = sdk.NewS3Client(v.Endpoint, v.Creds.AccessKey, v.Creds.SecretKey, transport)
		case pool.VendorCeph:
			client, err = sdk.NewS3Client(v.Endpoint, v.Creds.AccessKey, v.Creds.SecretKey, transport)
		case pool.VendorBaidu:
			client, err = sdk.NewBosClient(v.Endpoint, v.Creds.AccessKey, v.Creds.SecretKey)
		default:
			return fmt.Errorf("pool %q has unknown vendor type %q", k, v.Vendor)
		}
		if err != nil {
			return err
		}
		c.Pools[v.Vendor][k] = client
	}
	return nil
}

// GetClient xxx
func (c *Cluster) GetClient(p *pool.Pool) sdk.Client {
	clientMap, exists := c.Pools[p.Vendor]
	if !exists {
		return nil
	}
	client, exists := clientMap[p.ID]
	if !exists {
		return nil
	}
	return client
}

// CreateVBucket xxx
func (c *Cluster) CreateVBucket(vbucket, location string) error {
	exists, err := c.VBucketMgr.VBucketExists(vbucket)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("bucket %q already exists", vbucket)
	}
	pID := c.PoolMgr.AllocPoolByVBucket(vbucket)
	mID := c.VBucketMgr.AllocMdsByVBucket(vbucket)
	if len(pID) == 0 || len(mID) == 0 {
		return fmt.Errorf("couldn't alloc pool or mds")
	}
	err = c.VBucketMgr.CreateVBucket(vbucket, location, pID, mID)
	if err != nil {
		return err
	}
	return nil
}

// GetVBucketInfo xxx
func (c *Cluster) GetVBucketInfo(vbucket string) (*minio.BucketInfo, error) {
	vb, err := c.VBucketMgr.QueryVBucket(vbucket)
	if err != nil {
		return nil, err
	}
	if vb == nil {
		return nil, nil
	}
	bi := &minio.BucketInfo{
		Name:    vb.Name,
		Created: vb.CreatedTime,
	}
	return bi, nil
}

// ListVBuckets xxx
func (c *Cluster) ListVBuckets() ([]minio.BucketInfo, error) {
	vbs, err := c.VBucketMgr.ListVBuckets()
	if err != nil {
		return nil, err
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

// DeleteVBucket xxx
func (c *Cluster) DeleteVBucket(vbucket string) error {
	return c.VBucketMgr.DeleteVBucket(vbucket)
}

// ListObjects xxx
func (c *Cluster) ListObjects(bucket, prefix, marker, delimiter string,
	maxKeys int) (loi minio.ListObjectsInfo, err error) {
	lop := &object.ListObjectsParam{
		VBucket:   bucket,
		Prefix:    prefix,
		Marker:    marker,
		Delimiter: delimiter,
		Limits:    maxKeys,
	}
	lor, err := c.VBucketMgr.ListObjects(lop)
	if err != nil {
		return minio.ListObjectsInfo{}, err
	}
	objInfos := make([]minio.ObjectInfo, len(lor.Objects))
	for i, obj := range lor.Objects {
		objInfo := minio.ObjectInfo{
			Name:            obj.Name,
			Bucket:          obj.VBucket,
			ETag:            obj.Etag,
			InnerETag:       obj.InnerEtag,
			VersionID:       obj.VersionID,
			ContentType:     obj.ContentType,
			ContentEncoding: obj.ContentEncoding,
			StorageClass:    obj.StorageClass,
			UserTags:        obj.UserTags,
			Size:            obj.Size,
			IsDir:           obj.IsDir,
			IsLatest:        obj.IsLatest,
			DeleteMarker:    obj.DeleteMarker,
			RestoreOngoing:  obj.RestoreOngoing,
			ModTime:         obj.ModTime,
			AccTime:         obj.AccTime,
			Expires:         obj.Expires,
			RestoreExpires:  obj.RestoreExpires,
		}
		objInfos[i] = objInfo
	}
	isTruncated := false
	if len(lor.NextMarker) != 0 {
		isTruncated = true
	}
	return minio.ListObjectsInfo{
		IsTruncated: isTruncated,
		Objects:     objInfos,
		NextMarker:  lor.NextMarker,
		Prefixes:    lor.CommonPrefixs,
	}, nil
}

// PutObjectMeta xxx
func (c *Cluster) PutObjectMeta(pID, pBucket string, oi minio.ObjectInfo) error {
	timeNow := utils.GetCurrentTime()
	obj := &object.Object{
		Name:            oi.Name,
		VBucket:         oi.Bucket,
		Pool:            pID,
		Bucket:          pBucket,
		Etag:            oi.ETag,
		InnerEtag:       oi.InnerETag,
		VersionID:       oi.VersionID,
		ContentType:     oi.ContentType,
		ContentEncoding: oi.ContentEncoding,
		StorageClass:    oi.StorageClass,
		UserTags:        oi.UserTags,
		Size:            oi.Size,
		IsDir:           oi.IsDir,
		IsLatest:        oi.IsLatest,
		DeleteMarker:    oi.DeleteMarker,
		RestoreOngoing:  oi.RestoreOngoing,
		ModTime:         timeNow,
		AccTime:         timeNow,
		Expires:         oi.Expires,
		RestoreExpires:  oi.RestoreExpires,
	}
	return c.VBucketMgr.PutObjectMeta(obj)
}

// GetObjectMeta xxx
func (c *Cluster) GetObjectMeta(vbucket, object string) (minio.ObjectInfo, error) {
	obj, err := c.VBucketMgr.GetObjectMeta(vbucket, object)
	if err != nil {
		return minio.ObjectInfo{}, err
	}
	objInfo := minio.ObjectInfo{
		Name:            obj.Name,
		Bucket:          obj.VBucket,
		ETag:            obj.Etag,
		InnerETag:       obj.InnerEtag,
		VersionID:       obj.VersionID,
		ContentType:     obj.ContentType,
		ContentEncoding: obj.ContentEncoding,
		StorageClass:    obj.StorageClass,
		UserTags:        obj.UserTags,
		Size:            obj.Size,
		IsDir:           obj.IsDir,
		IsLatest:        obj.IsLatest,
		DeleteMarker:    obj.DeleteMarker,
		RestoreOngoing:  obj.RestoreOngoing,
		ModTime:         obj.ModTime,
		AccTime:         obj.AccTime,
		Expires:         obj.Expires,
		RestoreExpires:  obj.RestoreExpires,
	}
	return objInfo, nil
}

// DeleteObjectMeta xxx
func (c *Cluster) DeleteObjectMeta(vbucket, object string) error {
	return c.VBucketMgr.DeleteObjectMeta(vbucket, object)
}

// GetPool xxx
func (c *Cluster) GetPool(vbucket string) (*pool.Pool, error) {
	vb, err := c.VBucketMgr.QueryVBucket(vbucket)
	if err != nil {
		return nil, err
	}
	if vb == nil {
		return nil, fmt.Errorf("bucket %q not found", vbucket)
	}
	return c.PoolMgr.GetPool(vb.Pool), nil
}

// GetObjectPoolAndBucket xxx
func (c *Cluster) GetObjectPoolAndBucket(vbucket, object string) (*pool.Pool, string, error) {
	obj, err := c.VBucketMgr.GetObjectMeta(vbucket, object)
	if err != nil {
		return nil, "", err
	}
	return c.PoolMgr.GetPool(obj.Pool), obj.Bucket, nil
}

// AllocPhysicalBucket xxx
func (c *Cluster) AllocPhysicalBucket(pID string) string {
	return c.PoolMgr.AllocBucket(pID)
}

// CreateMultipart xxx
func (c *Cluster) CreateMultipart(pBucket, pObject, bucket, object, physicalUploadID string) (string, error) {
	return c.VBucketMgr.CreateMultipart(pBucket, pObject, bucket, object, physicalUploadID)
}

// QueryMultipart xxx
func (c *Cluster) QueryMultipart(bucket, uploadID string) (*multipart.Multipart, error) {
	return c.VBucketMgr.QueryMultipart(bucket, uploadID)
}

// DeleteMultipart xxx
func (c *Cluster) DeleteMultipart(bucket, uploadID string) error {
	return c.VBucketMgr.DeleteMultipart(bucket, uploadID)
}

// MultipartCommon xxx
type MultipartCommon struct {
	Multipart *multipart.Multipart
	Client    sdk.Client
}

// GetMultipartCommon xxx
func (c *Cluster) GetMultipartCommon(bucket, uploadID string) (*MultipartCommon, error) {
	var err error
	mc := &MultipartCommon{}
	mc.Multipart, err = c.QueryMultipart(bucket, uploadID)
	if err != nil {
		return nil, fmt.Errorf("multipart %q not found", uploadID)
	}
	pool, err := c.GetPool(bucket)
	if err != nil {
		return nil, err
	}
	if pool == nil {
		return nil, fmt.Errorf("pool of bucket %q not found", bucket)
	}
	mc.Client = c.GetClient(pool)
	if mc.Client == nil {
		return nil, fmt.Errorf("client of pool %q not found", pool.ID)
	}
	return mc, nil
}
