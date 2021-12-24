/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package mds

import (
	"fmt"

	"github.com/minio/minio/fusionstore/mgs"
	"github.com/minio/minio/protos"
)

// Mds metadata service
type Mds struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Region      string `json:"region"`
	Used        uint64 `json:"used"`
	Capacity    uint64 `json:"capacity"`
	Version     int32  `json:"version"`
	Endpoint    string `json:"endpoint"`
	CreatedTime string `json:"created_time"`
	UpdatedTime string `json:"updated_time"`
}

// DecodeFromPb xxx
func (m *Mds) DecodeFromPb(p *protos.Mds) {
	m.ID = p.GetId()
	m.Name = p.GetName()
	m.Type = p.GetType()
	m.Status = p.GetStatus()
	m.Region = p.GetRegion()
	m.Used = p.GetUsed()
	m.Capacity = p.GetCapacity()
	m.Endpoint = p.GetEndpoint()
	m.Version = p.GetVersion()
	m.CreatedTime = p.GetCreatedTime()
	m.UpdatedTime = p.GetUpdatedTime()
}

// Mgr xxx
type Mgr struct {
	mdsMap      map[string]*Mds
	mdsServices map[string]*Service
}

// NewMgr xxx
func NewMgr() (m *Mgr, err error) {
	m = &Mgr{}
	if err = m.loadMds(); err != nil {
		return nil, err
	}
	return m, nil
}

// Shutdown xxx
func (m *Mgr) Shutdown() {
	for _, srv := range m.mdsServices {
		srv.Close()
	}
}

// GetService xxx
func (m *Mgr) GetService(mdsID string) *Service {
	return m.mdsServices[mdsID]
}

// loadMds xxx
func (m *Mgr) loadMds() error {
	resp, err := mgs.GlobalService.ListMds()
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	m.mdsMap = make(map[string]*Mds, len(resp.GetMdsList()))
	m.mdsServices = make(map[string]*Service, len(resp.GetMdsList()))
	for _, v := range resp.GetMdsList() {
		md := &Mds{}
		md.DecodeFromPb(v)
		svc, err := NewService(md.Endpoint, 10)
		if err != nil {
			return err
		}
		m.mdsServices[md.ID] = svc
		m.mdsMap[md.ID] = md
	}
	return nil
}

// AllocateMds xxx
func (m *Mgr) AllocateMds(vbucket string) string {
	// FIXME(yangchunxin): select one mds
	return ""
}
