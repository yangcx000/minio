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

// GlobalClient xxx
var GlobalClient *Client

// Client xxx
type Client struct {
	svc     protos.MgsServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

// NewClient xxx
func NewClient(mgsAddr string, timeout int) error {
	var err error
	GlobalClient = &Client{timeout: time.Duration(timeout) * time.Second}
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	if GlobalClient.conn, err = grpc.Dial(mgsAddr, opts...); err != nil {
		return err
	}
	GlobalClient.svc = protos.NewMgsServiceClient(GlobalClient.conn)
	return nil
}

// Close xxx
func (c *Client) Close() {
	c.conn.Close()
}

// ListPools xxx
func (c *Client) ListPools() (*protos.ListPoolResponse, error) {
	req := &protos.ListPoolRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.svc.ListPool(ctx, req)
}

// ListBuckets xxx
func (c *Client) ListBuckets(poolID string) (*protos.ListBucketResponse, error) {
	req := &protos.ListBucketRequest{PoolId: poolID}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.svc.ListBucket(ctx, req)
}

// ListMds xxx
func (c *Client) ListMds() (*protos.ListMdsResponse, error) {
	req := &protos.ListMdsRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	return c.svc.ListMds(ctx, req)
}
