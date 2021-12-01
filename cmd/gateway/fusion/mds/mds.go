/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package mds

import "github.com/minio/minio/cmd/gateway/fusion/mgs"

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
	mc     *mgs.Client
	MdsMap map[string]*Mds
}

// NewMgr xxx
func NewMgr(c *mgs.Client) (*Mgr, error) {
	var err error
	mgr := &Mgr{mc: c}
	//mgr.MdsMap, err = mgr.mc.GetMdss()
	return mgr, err
}

/*
// DecodeFromPb xxx
func DecodeFromPb(p *protos.Mds) *MDS {
	m := &MDS{
		ID:          p.GetId(),
		Name:        p.GetName(),
		Type:        p.GetType(),
		Status:      p.GetStatus(),
		Region:      p.GetRegion(),
		Used:        p.GetUsed(),
		Capacity:    p.GetCapacity(),
		Pdservers:   p.GetPdservers(),
		Version:     p.GetVersion(),
		CreatedTime: p.GetCreatedTime(),
		UpdatedTime: p.GetUpdatedTime(),
	}
	return m
}
*/
