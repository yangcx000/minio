/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package vbucket

import (
	"errors"
	"fmt"
	"time"

	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/fusionstore/mds"
	"github.com/minio/minio/fusionstore/mgs"
	"github.com/minio/minio/fusionstore/multipart"
	"github.com/minio/minio/fusionstore/object"
	"github.com/minio/minio/fusionstore/pool"
	"github.com/minio/minio/fusionstore/utils"
	"github.com/minio/minio/protos"
)

const scanLimits = 1000

// VBucket xxx
type VBucket struct {
	ID          string                    `json:"id,omitempty"`
	Name        string                    `json:"name"`
	Status      string                    `json:"status,omitempty"`
	Owner       string                    `json:"owner,omitempty"`
	Pool        string                    `json:"pool,omitempty"`
	Mds         string                    `json:"mds,omitempty"`
	Location    string                    `json:"location,omitempty"`
	Version     int                       `json:"version,omitempty"`
	Objects     map[string]*object.Object `json:"objects,omitempty"`
	CreatedTime time.Time                 `json:"created_time,omitempty"`
	UpdatedTime time.Time                 `json:"updated_time,omitempty"`
}

// DecodeFromPb xxx
func (v *VBucket) DecodeFromPb(p *protos.VBucket) {
	v.ID = p.GetId()
	v.Name = p.GetName()
	v.Status = p.GetStatus()
	v.Owner = p.GetOwner()
	v.Pool = p.GetPool()
	v.Mds = p.GetMds()
	v.Location = p.GetLocation()
	v.Version = int(p.GetVersion())
	v.CreatedTime = p.GetCreatedTime().AsTime()
	v.UpdatedTime = p.GetUpdatedTime().AsTime()
}

func (v *VBucket) getObject(object string) *object.Object {
	obj, exists := v.Objects[object]
	if exists {
		return obj
	}
	// XXX
	return nil
}

// Mgr xxx
type Mgr struct {
	PoolMgr  *pool.Mgr
	MdsMgr   *mds.Mgr
	VBuckets map[string]*VBucket
}

// NewMgr xxx
func NewMgr() (*Mgr, error) {
	mgr := &Mgr{}
	if err := mgr.init(); err != nil {
		return nil, err
	}
	return mgr, nil
}

func (m *Mgr) init() error {
	var err error
	if m.PoolMgr, err = pool.NewMgr(); err != nil {
		return err
	}
	if m.MdsMgr, err = mds.NewMgr(); err != nil {
		return err
	}
	m.loadVBuckets()
	return nil
}

func (m *Mgr) loadVBuckets() error {
	vbs, err := m.ListVBuckets()
	if err != nil {
		return err
	}
	m.VBuckets = make(map[string]*VBucket, len(vbs))
	for _, v := range vbs {
		m.VBuckets[v.Name] = v
	}
	return nil
}

func (m *Mgr) getVBucket(vbucket string) *VBucket {
	vb, exists := m.VBuckets[vbucket]
	if exists {
		return vb
	}
	vb, _ = m.queryVBucket(vbucket)
	if vb != nil {
		m.VBuckets[vb.Name] = vb
	}
	return vb
}

/* vbucket apis */
func (m *Mgr) createVBucket(vbucket, location, pool, mds string) error {
	resp, err := mgs.GlobalService.CreateVBucket(vbucket, location, pool, mds)
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	return nil
}

func (m *Mgr) queryVBucket(vbucket string) (*VBucket, error) {
	resp, err := mgs.GlobalService.QueryVBucket(vbucket)
	if err != nil {
		return nil, err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		if resp.GetStatus().Code == protos.Code_NOT_FOUND {
			return nil, nil
		}
		return nil, fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	vb := &VBucket{}
	vb.DecodeFromPb(resp.GetVbucket())
	return vb, nil
}

func (m *Mgr) deleteVBucket(vbucket string) error {
	resp, err := mgs.GlobalService.DeleteVBucket(vbucket)
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	return nil
}

func (m *Mgr) listVBuckets() ([]*VBucket, error) {
	resp, err := mgs.GlobalService.ListVBuckets()
	if err != nil {
		return nil, err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return nil, fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	vbs := make([]*VBucket, len(resp.GetVbuckets()))
	for i, v := range resp.GetVbuckets() {
		vb := &VBucket{}
		vb.DecodeFromPb(v)
		vbs[i] = vb
	}
	return vbs, nil
}

/* end vbucket apis */

// MakeVBucket xxx
func (m *Mgr) MakeVBucket(vbucket, location string) error {
	// query mds directly
	vb, err := m.queryVBucket(vbucket)
	if err != nil {
		return err
	}
	if vb != nil {
		return fmt.Errorf("bucket %q already exists", vbucket)
	}
	pool, mds := m.PoolMgr.SelectPool(vbucket), m.MdsMgr.SelectMds(vbucket)
	if len(pool) == 0 || len(mds) == 0 {
		return errors.New("couldn't alloc pool or mds")
	}
	err = m.createVBucket(vbucket, location, pool, mds)
	if err != nil {
		return err
	}
	// query and add vbucket to cache
	_ = m.getVBucket(vbucket)
	return nil
}

// DeleteVBucket xxx
func (m *Mgr) DeleteVBucket(vbucket string) error {
	err := m.deleteVBucket(vbucket)
	if err != nil {
		return err
	}
	delete(m.VBuckets, vbucket)
	return nil
}

// GetVBucketInfo xxx
func (m *Mgr) GetVBucketInfo(vbucket string) (*VBucket, error) {
	// query mds directly
	return m.queryVBucket(vbucket)
}

// ListVBuckets xxx
func (m *Mgr) ListVBuckets() ([]*VBucket, error) {
	return m.listVBuckets()
}

// GetPoolAndBucket xxx
func (m *Mgr) GetPoolAndBucket(vbucket, object string) (*pool.Pool, string) {
	pBucket := ""
	vb := m.getVBucket(vbucket)
	if vb == nil {
		return nil, pBucket
	}
	pBucket = m.PoolMgr.SelectBucket(vb.Pool)
	return m.PoolMgr.GetPool(vb.Pool), pBucket
}

// GetPool xxx
func (m *Mgr) GetPool(vbucket string) *pool.Pool {
	vb := m.getVBucket(vbucket)
	if vb == nil {
		return nil
	}
	return m.PoolMgr.GetPool(vb.Pool)
}

/* object apis */
func (m *Mgr) putObject(obj *object.Object) error {
	vb := m.getVBucket(obj.VBucket)
	srv := m.MdsMgr.GetService(vb.Mds)
	resp, err := srv.PutObject(obj)
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	return nil
}

func (m *Mgr) getObject(vbucket, objectName string) (*object.Object, error) {
	vb := m.getVBucket(vbucket)
	srv := m.MdsMgr.GetService(vb.Mds)
	resp, err := srv.QueryObject(vbucket, objectName)
	if err != nil {
		return nil, err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		if resp.GetStatus().Code == protos.Code_NOT_FOUND {
			return nil, nil
		}
		return nil, fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	obj := &object.Object{}
	obj.DecodeFromPb(resp.GetObject())
	return obj, nil
}

func (m *Mgr) deleteObject(vbucket, object string) error {
	vb := m.getVBucket(vbucket)
	srv := m.MdsMgr.GetService(vb.Mds)
	resp, err := srv.DeleteObject(vbucket, object)
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	return nil
}

func (m *Mgr) listObjects(vbucket, marker string, numObjects int) ([]*object.Object, string, error) {
	nextMarker := ""
	vb := m.getVBucket(vbucket)
	if vb == nil {
		return nil, nextMarker, errors.New("couldn't get vbucket")
	}
	srv := m.MdsMgr.GetService(vb.Mds)
	if srv == nil {
		return nil, nextMarker, errors.New("couldn't get mds service")
	}
	if numObjects > scanLimits {
		return nil, nextMarker, fmt.Errorf("object limits must less than %d", scanLimits)
	}
	resp, err := srv.ListObjects(vbucket, marker, int32(numObjects))
	if err != nil {
		return nil, nextMarker, err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return nil, nextMarker, fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	objs := make([]*object.Object, len(resp.GetObjects()))
	for i, v := range resp.GetObjects() {
		obj := &object.Object{}
		obj.DecodeFromPb(v)
		objs[i] = obj
	}
	nextMarker = resp.GetNext()
	return objs, nextMarker, nil
}

// multipart apis
func (m *Mgr) createMultipart(mp *multipart.Multipart) (string, error) {
	// FIXME(yangchunxin): multi clients upload the same object by multipart ??
	vb := m.getVBucket(mp.VBucket)
	srv := m.MdsMgr.GetService(vb.Mds)
	resp, err := srv.CreateMultipart(mp)
	if err != nil {
		return "", err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return "", fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	return resp.GetUploadId(), nil
}

func (m *Mgr) getMultipart(vbucket, uploadID string) (*multipart.Multipart, error) {
	vb := m.getVBucket(vbucket)
	srv := m.MdsMgr.GetService(vb.Mds)
	resp, err := srv.QueryMultipart(vbucket, uploadID)
	if err != nil {
		return nil, err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		if resp.GetStatus().Code == protos.Code_NOT_FOUND {
			return nil, nil
		}
		return nil, fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	mp := &multipart.Multipart{}
	mp.DecodeFromPb(resp.GetMultipart())
	return mp, nil
}

func (m *Mgr) deleteMultipart(vbucket, uploadID string) error {
	vb := m.getVBucket(vbucket)
	srv := m.MdsMgr.GetService(vb.Mds)
	resp, err := srv.DeleteMultipart(vbucket, uploadID)
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	return nil
}

func (m *Mgr) listMultiparts(vbucket, marker string, numMultiparts int) ([]*multipart.Multipart, string, error) {
	nextMarker := ""
	vb := m.getVBucket(vbucket)
	if vb == nil {
		return nil, nextMarker, errors.New("couldn't get vbucket")
	}
	srv := m.MdsMgr.GetService(vb.Mds)
	if srv == nil {
		return nil, nextMarker, errors.New("couldn't get mds service")
	}
	if numMultiparts > scanLimits {
		return nil, nextMarker, fmt.Errorf("multipart limits must less than %d", scanLimits)
	}
	resp, err := srv.ListMultiparts(vbucket, marker, int32(numMultiparts))
	if err != nil {
		return nil, nextMarker, err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return nil, nextMarker, fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	mps := make([]*multipart.Multipart, len(resp.GetMultiparts()))
	for i, v := range resp.GetMultiparts() {
		mp := &multipart.Multipart{}
		mp.DecodeFromPb(v)
		mps[i] = mp
	}
	nextMarker = resp.GetNext()
	return mps, nextMarker, nil
}

// PutObjectMeta xxx
func (m *Mgr) PutObjectMeta(pID, pBucket string, objInfo minio.ObjectInfo) error {
	timeNow := utils.GetCurrentTime()
	obj := object.Object{
		Name:            objInfo.Name,
		VBucket:         objInfo.Bucket,
		Pool:            pID,
		Bucket:          pBucket,
		Etag:            objInfo.ETag,
		InnerEtag:       objInfo.InnerETag,
		VersionID:       objInfo.VersionID,
		ContentType:     objInfo.ContentType,
		ContentEncoding: objInfo.ContentEncoding,
		StorageClass:    objInfo.StorageClass,
		UserTags:        objInfo.UserTags,
		Size:            objInfo.Size,
		IsDir:           objInfo.IsDir,
		IsLatest:        objInfo.IsLatest,
		DeleteMarker:    objInfo.DeleteMarker,
		RestoreOngoing:  objInfo.RestoreOngoing,
		ModTime:         timeNow,
		AccTime:         timeNow,
		Expires:         objInfo.Expires,
		RestoreExpires:  objInfo.RestoreExpires,
	}
	return m.putObject(&obj)
}

// GetObjectMeta xxx
func (m *Mgr) GetObjectMeta(vbucket, object string) (minio.ObjectInfo, error) {
	obj, err := m.getObject(vbucket, object)
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
func (m *Mgr) DeleteObjectMeta(vbucket, object string) error {
	return m.deleteObject(vbucket, object)
}

// ListObjects xxx
func (m *Mgr) ListObjects(vbucket, prefix, marker, delimiter string, maxKeys int) (minio.ListObjectsInfo, error) {
	objs, nextMarker, err := m.listObjects(vbucket, marker, maxKeys)
	if err != nil {
		return minio.ListObjectsInfo{}, err
	}
	objInfos := make([]minio.ObjectInfo, len(objs))
	for i, obj := range objs {
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
	return minio.ListObjectsInfo{
		Objects:    objInfos,
		NextMarker: nextMarker,
	}, nil
}

// GetObjectPoolAndBucket xxx
func (m *Mgr) GetObjectPoolAndBucket(vbucket, object string) (*pool.Pool, string) {
	obj, err := m.getObject(vbucket, object)
	if err != nil {
		return nil, ""
	}
	return m.PoolMgr.GetPool(obj.Pool), obj.Bucket
}

// GetObjectKey xxx
func (m *Mgr) GetObjectKey(vbucket, object string) string {
	// add vbucket prefix
	return fmt.Sprintf("%s/%s", vbucket, object)
}

// CreateMultipart xxx
func (m *Mgr) CreateMultipart(pBucket, pObject, bucket, object, physicalUploadID string) (string, error) {
	mp := &multipart.Multipart{
		PhysicalUploadID: physicalUploadID,
		VBucket:          bucket,
		PhysicalBucket:   pBucket,
		Object:           object,
		CreatedTime:      time.Now(),
	}
	return m.createMultipart(mp)
}

// QueryMultipart xxx
func (m *Mgr) QueryMultipart(bucket, uploadID string) (*multipart.Multipart, error) {
	return m.getMultipart(bucket, uploadID)
}

//DeleteMultipart xxx
func (m *Mgr) DeleteMultipart(bucket, uploadID string) error {
	return m.deleteMultipart(bucket, uploadID)
}
