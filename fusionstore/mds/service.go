/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package mds

import (
	"context"
	"time"

	"github.com/minio/minio/fusionstore/multipart"
	"github.com/minio/minio/fusionstore/object"
	"github.com/minio/minio/protos"
	"google.golang.org/grpc"
)

// Service xxx
type Service struct {
	svc     protos.MdsServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

// NewService xxx
func NewService(mdsAddr string, timeout int) (*Service, error) {
	var err error
	service := &Service{timeout: time.Duration(timeout) * time.Second}
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	if service.conn, err = grpc.Dial(mdsAddr, opts...); err != nil {
		return service, err
	}
	service.svc = protos.NewMdsServiceClient(service.conn)
	return service, nil
}

// Close xxx
func (s *Service) Close() {
	s.conn.Close()
}

// PutObject xxx
func (s *Service) PutObject(obj *object.Object) (*protos.PutObjectResponse, error) {
	req := &protos.PutObjectRequest{Object: obj.EncodeToPb()}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.PutObject(ctx, req)
}

// DeleteObject xxx
func (s *Service) DeleteObject(vbucket, object string) (*protos.DeleteObjectResponse, error) {
	req := &protos.DeleteObjectRequest{
		Vbucket: vbucket,
		Object:  object,
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.DeleteObject(ctx, req)
}

// QueryObject xxx
func (s *Service) QueryObject(vbucket, object string) (*protos.QueryObjectResponse, error) {
	req := &protos.QueryObjectRequest{
		Vbucket: vbucket,
		Object:  object,
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.QueryObject(ctx, req)
}

// ListObjects xxx
func (s *Service) ListObjects(vbucket, prefix, marker, delimiter string, limits int32) (*protos.ListObjectsResponse, error) {
	req := &protos.ListObjectsRequest{
		Vbucket:   vbucket,
		Prefix:    prefix,
		Marker:    marker,
		Delimiter: delimiter,
		Limits:    limits,
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.ListObjects(ctx, req)
}

// CreateMultipart xxx
func (s *Service) CreateMultipart(mp *multipart.Multipart) (*protos.CreateMultipartResponse, error) {
	req := &protos.CreateMultipartRequest{Multipart: mp.EncodeToPb()}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.CreateMultipart(ctx, req)
}

// DeleteMultipart xxx
func (s *Service) DeleteMultipart(vbucket, uploadID string) (*protos.DeleteMultipartResponse, error) {
	req := &protos.DeleteMultipartRequest{
		Vbucket:  vbucket,
		UploadId: uploadID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.DeleteMultipart(ctx, req)
}

// QueryMultipart xxx
func (s *Service) QueryMultipart(vbucket, uploadID string) (*protos.QueryMultipartResponse, error) {
	req := &protos.QueryMultipartRequest{
		Vbucket:  vbucket,
		UploadId: uploadID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.QueryMultipart(ctx, req)
}

// ListMultiparts xxx
func (s *Service) ListMultiparts(vbucket, marker string, limits int32) (*protos.ListMultipartsResponse, error) {
	req := &protos.ListMultipartsRequest{Vbucket: vbucket, Prev: marker, Limits: limits}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.ListMultiparts(ctx, req)
}
