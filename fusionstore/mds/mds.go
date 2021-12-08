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
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Status      string   `json:"status"`
	Region      string   `json:"region"`
	Used        uint64   `json:"used"`
	Capacity    uint64   `json:"capacity"`
	Version     int32    `json:"version"`
	Pdservers   []string `json:"pdservers,omitempty"`
	CreatedTime string   `json:"created_time"`
	UpdatedTime string   `json:"updated_time"`
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
	m.Pdservers = p.GetPdservers()
	m.Version = p.GetVersion()
	m.CreatedTime = p.GetCreatedTime()
	m.UpdatedTime = p.GetUpdatedTime()
}

// Mgr xxx
type Mgr struct {
	MdsMap      map[string]*Mds
	MdsServices map[string]*Service
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
	mdsList, err := m.ListMds()
	if err != nil {
		return err
	}
	m.MdsMap = make(map[string]*Mds, len(mdsList))
	m.MdsServices = make(map[string]*Service, len(mdsList))
	for _, v := range mdsList {
		// TODO(yangchunxin): fix it
		svc, err := NewService(v.Pdservers[0], 10)
		if err != nil {
			return err
		}
		m.MdsServices[v.ID] = svc
		m.MdsMap[v.ID] = v
	}
	return nil
}

// ListMds xxx
func (m *Mgr) ListMds() ([]*Mds, error) {
	resp, err := mgs.GlobalService.ListMds()
	if err != nil {
		return nil, err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return nil, fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	mdsList := make([]*Mds, len(resp.GetMdsList()))
	for i, v := range resp.GetMdsList() {
		p := &Mds{}
		p.DecodeFromPb(v)
		mdsList[i] = p
	}
	return mdsList, nil
}

// AllocMds xxx
func (m *Mgr) AllocMds(vbucket string) string {
	_ = vbucket
	// FIXME(yangchunxin): design algorithm
	for mdsID := range m.MdsMap {
		return mdsID
	}
	return ""
}
