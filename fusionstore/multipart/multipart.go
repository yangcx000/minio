/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package multipart

import (
	"time"

	"github.com/minio/minio/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Multipart xxx
type Multipart struct {
	UploadID         string    `json:"upload_id"`
	PhysicalUploadID string    `json:"physical_upload_id"`
	VBucket          string    `json:"vbucket"`
	PhysicalBucket   string    `json:"physical_bucket"`
	Object           string    `json:"object"`
	CreatedTime      time.Time `json:"created_time"`
}

// EncodeToPb xxx
func (m *Multipart) EncodeToPb() *protos.Multipart {
	p := &protos.Multipart{
		UploadId:         m.UploadID,
		PhysicalUploadId: m.PhysicalUploadID,
		Vbucket:          m.VBucket,
		PhysicalBucket:   m.PhysicalBucket,
		Object:           m.Object,
		CreatedTime:      timestamppb.New(m.CreatedTime),
	}
	return p
}

// DecodeFromPb xxx
func (m *Multipart) DecodeFromPb(p *protos.Multipart) {
	m.UploadID = p.GetUploadId()
	m.PhysicalUploadID = p.GetPhysicalUploadId()
	m.VBucket = p.GetVbucket()
	m.PhysicalBucket = p.GetPhysicalBucket()
	m.Object = p.GetObject()
	m.CreatedTime = p.GetCreatedTime().AsTime()
}
