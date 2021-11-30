/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package fusion

import (
	"fmt"

	"github.com/minio/minio/protos"
)

// Credentials xxx
type Credentials struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
}

// Pool xxx
type Pool struct {
	ID          string      `json:"id,omitempty"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Status      string      `json:"status,omitempty"`
	Version     int         `json:"version,omitempty"`
	Endpoint    string      `json:"endpoint"`
	Creds       Credentials `json:"creds"`
	CreatedTime string      `json:"createdtime,omitempty"`
	UpdatedTime string      `json:"updatedtime,omitempty"`
}

// DecodeFromPb xxx
func DecodeFromPb(p *protos.Pool) *Pool {
	m := &Pool{
		ID:       p.GetId(),
		Name:     p.GetName(),
		Type:     p.GetType(),
		Status:   p.GetStatus(),
		Version:  int(p.GetVersion()),
		Endpoint: p.GetEndpoint(),
		Creds: Credentials{
			AccessKey: p.GetCreds().GetAccessKey(),
			SecretKey: p.GetCreds().GetSecretKey(),
		},
		CreatedTime: p.GetCreatedTime(),
		UpdatedTime: p.GetUpdatedTime(),
	}
	return m
}

// PoolMgr xxx
type PoolMgr struct {
	client *MgsClient
	Pools  map[string]*Pool
}

// NewPoolMgr xxx
func NewPoolMgr(mgsAddr string) (*PoolMgr, error) {
	pm := &PoolMgr{}
	pm.client = NewMgsClient(mgsAddr, 10)
	if pm.client == nil {
		return pm, fmt.Errorf("couldn't connect to %s", mgsAddr)
	}
	var err error
	// get pools
	pm.Pools, err = pm.client.GetPools()
	return pm, err
}

// GetPoolByBucket xxx
func (m *PoolMgr) GetPoolByBucket(bucketName string) string {
	// FIXME: design pool select algorithm
	ids := []string{}
	for k := range m.Pools {
		ids = append(ids, k)
	}
	return ids[0]
}
