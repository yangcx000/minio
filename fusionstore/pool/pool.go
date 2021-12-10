/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package pool

import (
	"fmt"

	"github.com/minio/minio/fusionstore/mgs"
	"github.com/minio/minio/protos"
)

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
	ID          string             `json:"id,omitempty"`
	Name        string             `json:"name"`
	Type        string             `json:"type"`
	Status      string             `json:"status,omitempty"`
	Version     int                `json:"version,omitempty"`
	Endpoint    string             `json:"endpoint"`
	Creds       Credentials        `json:"creds"`
	Buckets     map[string]*Bucket `json:"buckets"`
	CreatedTime string             `json:"createdtime,omitempty"`
	UpdatedTime string             `json:"updatedtime,omitempty"`
}

// DecodeFromPb xxx
func (p *Pool) DecodeFromPb(pool *protos.Pool) {
	p.ID = pool.GetId()
	p.Name = pool.GetName()
	p.Type = pool.GetType()
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

// Init xxx
func (p *Pool) init() error {
	buckets, err := p.ListBuckets()
	if err != nil {
		return err
	}
	p.Buckets = make(map[string]*Bucket)
	for _, b := range buckets {
		p.Buckets[b.ID] = b
	}
	return nil
}

// ListBuckets xxx
func (p *Pool) ListBuckets() ([]*Bucket, error) {
	resp, err := mgs.GlobalService.ListBuckets(p.ID)
	if err != nil {
		return nil, err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return nil, fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	buckets := make([]*Bucket, len(resp.GetBucketList()))
	for i, v := range resp.GetBucketList() {
		b := &Bucket{}
		b.DecodeFromPb(v)
		buckets[i] = b
	}
	return buckets, nil
}

func (p *Pool) selectBucket() string {
	// FIXME(yangchunxin): design algorithm
	for key := range p.Buckets {
		return key
	}
	return ""
}

// Mgr xxx
type Mgr struct {
	Pools map[string]*Pool
}

// NewMgr xxx
func NewMgr() (*Mgr, error) {
	mgr := &Mgr{}
	if err := mgr.init(); err != nil {
		return nil, err
	}
	return mgr, nil
}

// Init xxx
func (m *Mgr) init() error {
	pools, err := m.ListPools()
	if err != nil {
		return err
	}
	m.Pools = make(map[string]*Pool, len(pools))
	for _, v := range pools {
		m.Pools[v.ID] = v
	}
	return nil
}

// ListPools xxx
func (m *Mgr) ListPools() ([]*Pool, error) {
	resp, err := mgs.GlobalService.ListPools()
	if err != nil {
		return nil, err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return nil, fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	pools := make([]*Pool, len(resp.GetPoolList()))
	for i, v := range resp.GetPoolList() {
		p := &Pool{}
		p.DecodeFromPb(v)
		if err = p.init(); err != nil {
			return nil, err
		}
		pools[i] = p
	}
	return pools, nil
}

// SelectPool xxx
func (m *Mgr) SelectPool(vbucket string) string {
	// XXX: design algorithm
	for k := range m.Pools {
		return k
	}
	return ""
}

// SelectBucket xxx
func (m *Mgr) SelectBucket(pool string) string {
	// XXX: pool has cached
	p := m.Pools[pool]
	return p.selectBucket()
}
