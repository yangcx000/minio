/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package fusion

// MDS metadata service
type MDS struct {
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
