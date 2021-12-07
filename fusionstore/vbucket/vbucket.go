/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package vbucket

import (
	"fmt"
	"time"

	"github.com/minio/minio/fusionstore/mds"
	"github.com/minio/minio/fusionstore/mgs"
	"github.com/minio/minio/fusionstore/pool"
	"github.com/minio/minio/protos"
)

// VBucket xxx
type VBucket struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Status      string    `json:"status,omitempty"`
	Owner       string    `json:"owner,omitempty"`
	Pool        string    `json:"pool,omitempty"`
	Mds         string    `json:"mds,omitempty"`
	Location    string    `json:"location,omitempty"`
	Version     int       `json:"version,omitempty"`
	CreatedTime time.Time `json:"created_time,omitempty"`
	UpdatedTime time.Time `json:"updated_time,omitempty"`
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

// Mgr xxx
type Mgr struct {
	VBucketMap map[string]*VBucket
	PoolMgr    *pool.Mgr
	MdsMgr     *mds.Mgr
}

// NewMgr xxx
func NewMgr() (*Mgr, error) {
	mgr := &Mgr{}
	if err := mgr.Init(); err != nil {
		return nil, err
	}
	return mgr, nil
}

// Init xxx
func (m *Mgr) Init() error {
	var err error
	if m.PoolMgr, err = pool.NewMgr(); err != nil {
		return err
	}
	if m.MdsMgr, err = mds.NewMgr(); err != nil {
		return err
	}
	resp, err := mgs.GlobalService.ListVBuckets()
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	m.VBucketMap = make(map[string]*VBucket)
	for _, v := range resp.GetVbuckets() {
		p := &VBucket{}
		p.DecodeFromPb(v)
		m.VBucketMap[p.Name] = p
	}
	return nil
}

// MakeBucket xxx
func (m *Mgr) MakeBucket(bucket, location string) error {
	// TODO(yangchunxin): select pool and mds
	pool := "xxx"
	mds := "xxx"
	resp, err := mgs.GlobalService.CreateVBucket(bucket, location, pool, mds)
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	m.VBucketMap[bucket] = &VBucket{
		Name:     bucket,
		Pool:     pool,
		Mds:      mds,
		Location: location,
	}
	return nil
}

// DeleteBucket xxx
func (m *Mgr) DeleteBucket(bucket string) error {
	resp, err := mgs.GlobalService.DeleteVBucket(bucket)
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	return nil
}

// GetBucketInfo xxx
func (m *Mgr) GetBucketInfo(bucket string) (*VBucket, error) {
	resp, err := mgs.GlobalService.QueryVBucket(bucket)
	if err != nil {
		return nil, err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return nil, fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	vb := &VBucket{}
	vb.DecodeFromPb(resp.GetVbucket())
	return vb, nil
}

// ListBuckets xxx
func (m *Mgr) ListBuckets() ([]*VBucket, error) {
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
