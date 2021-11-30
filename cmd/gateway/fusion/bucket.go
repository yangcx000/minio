/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package fusion

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

/*
// DecodeFromPb xxx
func DecodeFromPb(b *protos.Bucket) *Bucket {
	m := &Bucket{
		ID:          b.GetId(),
		Name:        b.GetName(),
		Status:      b.GetStatus(),
		PoolID:      b.GetPoolId(),
		Version:     int(b.GetVersion()),
		CreatedTime: b.GetCreatedTime(),
		UpdatedTime: b.GetUpdatedTime(),
	}
	return m
}
*/
