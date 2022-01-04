/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package pool

import (
	"fmt"
	"sync/atomic"

	"github.com/minio/minio/fusionstore/mgs"
	"github.com/minio/minio/protos"
)

// Constants of pool
const (
	// vendor type
	VendorUnknown = "unknown"
	VendorBaidu   = "bos"
	VendorAws     = "s3"
	VendorCeph    = "rgw"

	// pool status
	poolStatusUnknown = "unknown"
	poolStatusActive  = "active"
	poolStatusStandby = "standby"

	// bucket status
	bucketStatusUnknown = "unknown"
	bucketStatusActive  = "active"
	bucketStatusStandby = "standby"
)

// global bucket index
var bucketIndex uint64

// Bucket xxx
type Bucket struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Status      string `json:"status,omitempty"`
	PoolID      string `json:"pool_id,omitempty"`
	Version     int    `json:"version,omitempty"`
	CreatedTime string `json:"created_time,omitempty"`
	UpdatedTime string `json:"updated_time,omitempty"`
}

// DecodeFromPb xxx
func (b *Bucket) DecodeFromPb(p *protos.Bucket) {
	b.ID = p.GetId()
	b.Name = p.GetName()
	b.Status = p.GetStatus()
	b.PoolID = p.GetPoolId()
	b.Version = int(p.GetVersion())
	b.CreatedTime = p.GetCreatedTime()
	b.UpdatedTime = p.GetUpdatedTime()
}

// Credentials xxx
type Credentials struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// Pool xxx
type Pool struct {
	ID       string      `json:"id,omitempty"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Vendor   string      `json:"vendor"`
	Status   string      `json:"status,omitempty"`
	Version  int         `json:"version,omitempty"`
	Endpoint string      `json:"endpoint"`
	Creds    Credentials `json:"creds"`
	// id --> bucket
	Buckets     []*Bucket `json:"buckets,omitempty"`
	CreatedTime string    `json:"createdtime,omitempty"`
	UpdatedTime string    `json:"updatedtime,omitempty"`
}

// DecodeFromPb xxx
func (p *Pool) DecodeFromPb(pool *protos.Pool) {
	p.ID = pool.GetId()
	p.Name = pool.GetName()
	p.Type = pool.GetType()
	p.Vendor = pool.GetVendor()
	p.Status = pool.GetStatus()
	p.Version = int(pool.GetVersion())
	p.Endpoint = pool.GetEndpoint()
	p.Creds = Credentials{
		AccessKey: pool.GetCreds().GetAccessKey(),
		SecretKey: pool.GetCreds().GetSecretKey(),
	}
	p.CreatedTime = pool.GetCreatedTime()
	p.UpdatedTime = pool.GetUpdatedTime()
}

func (p *Pool) init() error {
	resp, err := mgs.GlobalService.ListBuckets(p.ID)
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	for _, v := range resp.GetBucketList() {
		b := &Bucket{}
		b.DecodeFromPb(v)
		// FIXME(yangchunxin): filter active buckets
		if b.Status != bucketStatusActive {
			continue
		}
		p.Buckets = append(p.Buckets, b)
	}
	return nil
}

func (p *Pool) allocBucket() string {
	// no bucket
	if len(p.Buckets) == 0 {
		return ""
	}
	// rr algorithm
	counter := atomic.AddUint64(&bucketIndex, 1)
	index := counter % uint64(len(p.Buckets))
	return p.Buckets[index].Name
}

// Mgr xxx
type Mgr struct {
	pools map[string]*Pool
}

// NewMgr xxx
func NewMgr() (m *Mgr, err error) {
	m = &Mgr{}
	if err = m.loadPools(); err != nil {
		return nil, err
	}
	return m, nil
}

// loadPools xxx
func (m *Mgr) loadPools() error {
	// TODO(yangchunxin): add marker and limits
	resp, err := mgs.GlobalService.ListPools()
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	m.pools = make(map[string]*Pool, len(resp.GetPoolList()))
	for _, v := range resp.GetPoolList() {
		p := &Pool{}
		p.DecodeFromPb(v)
		// filter active pools
		if p.Status != poolStatusActive {
			continue
		}
		if err = p.init(); err != nil {
			return err
		}
		m.pools[p.ID] = p
	}
	return nil
}

// GetPoolMap xxx
func (m *Mgr) GetPoolMap() map[string]*Pool {
	return m.pools
}

// GetPool xxx
func (m *Mgr) GetPool(pID string) *Pool {
	return m.pools[pID]
}

// AllocatePool xxx
func (m *Mgr) AllocatePool(vbucket string) string {
	// XXX(yangchunxin): disabled here, manully allocate
	return ""
}

// AllocBucket xxx
func (m *Mgr) AllocBucket(pID string) string {
	p, exists := m.pools[pID]
	if !exists {
		return ""
	}
	return p.allocBucket()
}
