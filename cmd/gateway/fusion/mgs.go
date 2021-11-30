/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 *
 */

package fusion

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio/protos"
	"google.golang.org/grpc"
)

// MgsClient xxx
type MgsClient struct {
	svc     protos.MgsServiceClient
	conn    *grpc.ClientConn
	timeout time.Duration
}

// NewMgsClient xxx
func NewMgsClient(mgsAddr string, timeout int) *MgsClient {
	var err error
	c := &MgsClient{timeout: time.Duration(timeout) * time.Second}
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	if c.conn, err = grpc.Dial(mgsAddr, opts...); err != nil {
		return nil
	}
	c.svc = protos.NewMgsServiceClient(c.conn)
	return c
}

// Close xxx
func (c *MgsClient) Close() {
	c.conn.Close()
}

// GetPools xxx
func (c *MgsClient) GetPools() (map[string]*Pool, error) {
	pools := make(map[string]*Pool)
	req := &protos.ListPoolRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	if resp, err := c.svc.ListPool(ctx, req); err == nil {
		if resp.GetStatus().Code != protos.Code_OK {
			return nil, fmt.Errorf("%s", resp.GetStatus().GetMsg())
		}
		for _, v := range resp.GetPoolList() {
			p := DecodeFromPb(v)
			pools[p.ID] = p
		}
	}
	return pools, nil
}
