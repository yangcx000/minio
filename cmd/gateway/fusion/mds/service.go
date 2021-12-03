/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package mds

import (
	"context"
	"time"

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

// AddObject xxx
func (s *Service) AddObject() (*protos.AddObjectResponse, error) {
	req := &protos.AddObjectRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.AddObject(ctx, req)
}
