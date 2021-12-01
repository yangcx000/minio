/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package pool

import (
	"fmt"

	"github.com/minio/minio/cmd/gateway/fusion/mgs"
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
func (p *Pool) Init() error {
	resp, err := mgs.GlobalClient.ListBuckets(p.ID)
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	p.Buckets = make(map[string]*Bucket)
	for _, v := range resp.GetBucketList() {
		b := &Bucket{}
		b.DecodeFromPb(v)
		p.Buckets[b.ID] = b
	}
	return nil
}

// Mgr xxx
type Mgr struct {
	Pools map[string]*Pool
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
	resp, err := mgs.GlobalClient.ListPools()
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	m.Pools = make(map[string]*Pool)
	for _, v := range resp.GetPoolList() {
		p := &Pool{}
		p.DecodeFromPb(v)
		if err = p.Init(); err != nil {
			return err
		}
		m.Pools[p.ID] = p
	}
	return nil
}

// GetPoolByBucket xxx
func (m *Mgr) GetPoolByBucket(bucketName string) string {
	// FIXME: design pool select algorithm
	ids := []string{}
	for k := range m.Pools {
		ids = append(ids, k)
	}
	return ids[0]
}
