/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 *
 */

package mgs

import (
	"fmt"
	"testing"
)

const (
	mgsAddr      = "192.168.58.4:8000"
	timeoutValue = 10
)

func TestNewService(t *testing.T) {
	err := NewService(mgsAddr, timeoutValue)
	if err != nil {
		t.Fatalf("Couldn't new mgs service %s", mgsAddr)
	}
}

func TestClose(t *testing.T) {
	TestNewService(t)
	GlobalService.Close()
}

func TestListPools(t *testing.T) {
	TestNewService(t)
	defer TestClose(t)

	resp, err := GlobalService.ListPools()
	if err != nil {
		t.Fatalf("ListPools rpc failed, %v", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatalf("ListPools result failed, %v", err)
	}
	if len(resp.GetPoolList()) != 1000 {
		t.Errorf("ListPools results number %d != 10", len(resp.GetPoolList()))
	}
}

func TestListBuckets(t *testing.T) {
	TestNewService(t)
	defer TestClose(t)

	poolID := "pool-xiref"
	resp, err := GlobalService.ListBuckets(poolID)
	if err != nil {
		t.Fatalf("ListBuckets rpc failed, %v", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatalf("ListBuckets result failed, %v", err)
	}
	if len(resp.GetBucketList()) != 1000 {
		t.Errorf("ListBuckets results number %d != 10", len(resp.GetBucketList()))
	}
}

func TestListMds(t *testing.T) {
	TestNewService(t)
	defer TestClose(t)

	resp, err := GlobalService.ListMds()
	if err != nil {
		t.Fatalf("ListMds rpc failed, %v", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatalf("ListMds result failed, %v", err)
	}
	if len(resp.GetMdsList()) != 1000 {
		t.Errorf("ListMds results number %d != 10", len(resp.GetMdsList()))
	}
}

/*
func TestCreateVBucket(t *testing.T) {
	TestNewService(t)
	defer TestClose(t)

	for i := 0; i < 1000; i++ {
		name, location, pool, mds := fmt.Sprintf("vbucket-%d", i), "bj", "pool-acbjx", "mds-aagka"
		resp, err := GlobalService.CreateVBucket(name, location, pool, mds)
		if err != nil {
			t.Fatalf("CreateVBucket rpc failed, %s", err)
		}
		if resp.GetStatus().GetCode() != 0 {
			t.Fatalf("CreateVBucket result failed, %v", err)
		}
	}
}
*/

func TestQueryVBuckets(t *testing.T) {
	TestNewService(t)
	defer TestClose(t)

	name := "vbucket-100"
	resp, err := GlobalService.QueryVBucket(name)
	if err != nil {
		t.Fatalf("QueryVBucket rpc failed, %v", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatalf("QueryVBucket result failed, %v", err)
	}
	if resp.GetVbucket().GetName() != name {
		t.Fatalf("QueryVBucket name %q not equal %q", resp.GetVbucket().GetName(), name)
	}
}

func TestListVBuckets(t *testing.T) {
	TestNewService(t)
	defer TestClose(t)

	resp, err := GlobalService.ListVBuckets()
	if err != nil {
		t.Fatalf("ListVBuckets rpc failed, %v", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatalf("ListVBuckets result failed, %v", err)
	}
	if len(resp.GetVbuckets()) != 1000 {
		t.Fatalf("ListVBuckets results number != %d", len(resp.GetVbuckets()))
	}
}

func TestDeleteVBucket(t *testing.T) {
	TestNewService(t)
	defer TestClose(t)

	for i := 0; i < 1000; i++ {
		name := fmt.Sprintf("vbucket-%d", i)
		resp, err := GlobalService.DeleteVBucket(name)
		if err != nil {
			t.Fatalf("DeleteVBucket rpc failed, %v", err)
		}
		if resp.GetStatus().GetCode() != 0 {
			t.Fatalf("DeleteVBucket result failed, %v", err)
		}
	}
}
