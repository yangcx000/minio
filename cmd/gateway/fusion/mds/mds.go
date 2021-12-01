/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package mds

import (
	"fmt"

	"github.com/minio/minio/cmd/gateway/fusion/mgs"
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

// Mgr xxx
type Mgr struct {
	MdsMap map[string]*Mds
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
	resp, err := mgs.GlobalClient.ListMds()
	if err != nil {
		return err
	}
	if resp.GetStatus().Code != protos.Code_OK {
		return fmt.Errorf("%s", resp.GetStatus().GetMsg())
	}
	m.MdsMap = make(map[string]*Mds)
	for _, v := range resp.GetMdsList() {
		p := &Mds{}
		p.DecodeFromPb(v)
		m.MdsMap[p.ID] = p
	}
	return nil
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
