/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package vbucket

import (
	"fmt"
	"time"

	minio "github.com/minio/minio/cmd"
	"github.com/minio/minio/fusionstore/mds"
	"github.com/minio/minio/fusionstore/mgs"
	"github.com/minio/minio/fusionstore/object"
	"github.com/minio/minio/fusionstore/pool"
	"github.com/minio/minio/protos"
)

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
	// XXX(yangchunxin): fixme
	vbs, err := m.ListVBuckets()
	if err != nil {
		return err
	}
	m.VBuckets = make(map[string]*VBucket)
	for _, v := range vbs {
		m.VBuckets[v.Name] = v
	}
	return nil
}

func (m *Mgr) getVBucket(vbucket string) *VBucket {
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
	vb, err := m.queryVBucket(vbucket)
	if err != nil {
		return err
	}
	if vb != nil {
		return fmt.Errorf("bucket %q already exists", vbucket)
	}
	pool, mds := m.PoolMgr.SelectPool(vbucket), m.MdsMgr.SelectMds(vbucket)
	err = m.createVBucket(vbucket, location, pool, mds)
	if err != nil {
		return err
	}
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
	return m.queryVBucket(vbucket)
}

// ListVBuckets xxx
func (m *Mgr) ListVBuckets() ([]*VBucket, error) {
	return m.listVBuckets()
}

// GetPoolAndBucket xxx
func (m *Mgr) GetPoolAndBucket(vbucket, object string) (string, string) {
	vb := getVBucket(vbucket)
	if vb == nil {
		return "", ""
	}
	// FIXME(yangchunxin): xxx
	bucket := m.PoolMgr.AllocBucket(vb.Pool)
	return vb.Pool, bucket
}

// GetPoolAndBucket xxx
func (m *Mgr) GetPoolAndBucket(vbucket, object string) (string, string) {

	return "", ""
}

// PutObjectMeta xxx
func (m *Mgr) PutObjectMeta(pID, pBucket string, objInfo minio.ObjectInfo) error {
	return nil
}

// GetObjectMeta xxx
func (m *Mgr) GetObjectMeta(vbucket, object string) (minio.ObjectInfo, error) {
	return minio.ObjectInfo{}, nil
}

// DeleteObjectMeta xxx
func (m *Mgr) DeleteObjectMeta(pID, pBucket string, objInfo minio.ObjectInfo) error {
	return nil
}
