/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 *
 */

package mgs

import (
	"context"
	"time"

	"github.com/minio/minio/protos"
	"google.golang.org/grpc"
)

// GlobalService xxx
var GlobalService *Service

// Service xxx
type Service struct {
	svc     protos.MgsServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

// NewService xxx
func NewService(mgsAddr string, timeout int) error {
	var err error
	GlobalService = &Service{timeout: time.Duration(timeout) * time.Second}
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	if GlobalService.conn, err = grpc.Dial(mgsAddr, opts...); err != nil {
		return err
	}
	GlobalService.svc = protos.NewMgsServiceClient(GlobalService.conn)
	return nil
}

// Close xxx
func (s *Service) Close() {
	s.conn.Close()
}

// ListPools xxx
func (s *Service) ListPools() (*protos.ListPoolResponse, error) {
	req := &protos.ListPoolRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.ListPool(ctx, req)
}

// ListBuckets xxx
func (s *Service) ListBuckets(poolID string) (*protos.ListBucketResponse, error) {
	req := &protos.ListBucketRequest{PoolId: poolID}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.ListBucket(ctx, req)
}

// ListMds xxx
func (s *Service) ListMds() (*protos.ListMdsResponse, error) {
	req := &protos.ListMdsRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.ListMds(ctx, req)
}

// CreateVBucket xxx
func (s *Service) CreateVBucket(name, location, pool, mds string) (*protos.CreateVBucketResponse, error) {
	req := &protos.CreateVBucketRequest{
		Vbucket: &protos.VBucket{
			Name:     name,
			Pool:     pool,
			Mds:      mds,
			Location: location,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.CreateVBucket(ctx, req)
}

// DeleteVBucket xxx
func (s *Service) DeleteVBucket(name string) (*protos.DeleteVBucketResponse, error) {
	req := &protos.DeleteVBucketRequest{
		Name: name,
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.DeleteVBucket(ctx, req)
}

// QueryVBucket xxx
func (s *Service) QueryVBucket(name string) (*protos.QueryVBucketResponse, error) {
	req := &protos.QueryVBucketRequest{
		Name: name,
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.QueryVBucket(ctx, req)
}

// ListVBuckets xxx
func (s *Service) ListVBuckets() (*protos.ListVBucketResponse, error) {
	req := &protos.ListVBucketRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	return s.svc.ListVBucket(ctx, req)
}
