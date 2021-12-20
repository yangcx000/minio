/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package object

import (
	"time"

	"github.com/minio/minio/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Object xxx
type Object struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	VBucket         string    `json:"vbucket"`
	Pool            string    `json:"pool"`
	Bucket          string    `json:"bucket"`
	Etag            string    `json:"etag"`
	InnerEtag       string    `json:"inner_etag"`
	VersionID       string    `json:"version_id"`
	ContentType     string    `json:"content_type"`
	ContentEncoding string    `json:"content_encoding"`
	StorageClass    string    `json:"storage_class"`
	UserTags        string    `json:"user_tags"`
	Size            int64     `json:"size"`
	IsDir           bool      `json:"is_dir"`
	IsLatest        bool      `json:"is_latest"`
	DeleteMarker    bool      `json:"delete_marker"`
	RestoreOngoing  bool      `json:"restore_ongoing"`
	ModTime         time.Time `json:"mod_time"`
	AccTime         time.Time `json:"acc_time"`
	Expires         time.Time `json:"expires"`
	RestoreExpires  time.Time `json:"restore_expires"`
}

// EncodeToPb xxx
func (o *Object) EncodeToPb() *protos.Object {
	v := &protos.Object{
		Id:              o.ID,
		Name:            o.Name,
		Vbucket:         o.VBucket,
		Pool:            o.Pool,
		Bucket:          o.Bucket,
		Etag:            o.Etag,
		InnerEtag:       o.InnerEtag,
		VersionId:       o.VersionID,
		ContentType:     o.ContentType,
		ContentEncoding: o.ContentEncoding,
		StorageClass:    o.StorageClass,
		UserTags:        o.UserTags,
		Size:            o.Size,
		IsDir:           o.IsDir,
		IsLatest:        o.IsLatest,
		DeleteMarker:    o.DeleteMarker,
		RestoreOngoing:  o.RestoreOngoing,
		ModTime:         timestamppb.New(o.ModTime),
		AccTime:         timestamppb.New(o.AccTime),
		Expires:         timestamppb.New(o.Expires),
		RestoreExpires:  timestamppb.New(o.RestoreExpires),
	}
	return v
}

// DecodeFromPb xxx
func (o *Object) DecodeFromPb(p *protos.Object) {
	o.ID = p.Id
	o.Name = p.Name
	o.VBucket = p.Vbucket
	o.Pool = p.Pool
	o.Bucket = p.Bucket
	o.Etag = p.Etag
	o.InnerEtag = p.InnerEtag
	o.VersionID = p.VersionId
	o.ContentType = p.ContentType
	o.ContentEncoding = p.ContentEncoding
	o.StorageClass = p.StorageClass
	o.UserTags = p.UserTags
	o.Size = p.Size
	o.IsDir = p.IsDir
	o.IsLatest = p.IsLatest
	o.DeleteMarker = p.DeleteMarker
	o.RestoreOngoing = p.RestoreOngoing
	o.ModTime = p.ModTime.AsTime()
	o.AccTime = p.AccTime.AsTime()
	o.Expires = p.Expires.AsTime()
	o.RestoreExpires = p.RestoreExpires.AsTime()
}

// ListObjectsParam xxx
type ListObjectsParam struct {
	VBucket   string
	Prefix    string
	Marker    string
	Delimiter string
	Limits    int
}

// ListObjectsResult xxx
type ListObjectsResult struct {
	Objects       []*Object
	CommonPrefixs []string
	NextMarker    string
}
