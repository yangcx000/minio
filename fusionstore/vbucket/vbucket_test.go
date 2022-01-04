/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package vbucket

import (
	"fmt"
	"os"
	"testing"

	"github.com/minio/minio/fusionstore/mgs"
	"github.com/minio/minio/fusionstore/utils"
)

var mgr *Mgr

func init() {
	mgsAddr := "192.168.58.4:8000"
	serviceTimeout := 10
	if err := mgs.NewService(mgsAddr, serviceTimeout); err != nil {
		os.Exit(1)
	}
	utils.PrettyPrint(mgr)
}

func TestNewMgr(t *testing.T) {
	var err error
	mgr, err = NewMgr()
	if err != nil {
		t.Fatal("NewMgr failed")
	}
	if len(mgr.MdsMgr.GetMdsMap()) != 2 {
		t.Errorf("mds number %d != 2", len(mgr.MdsMgr.GetMdsMap()))
	}
	if len(mgr.MdsMgr.GetMdsServices()) != 2 {
		t.Errorf("mds services %d != 2", len(mgr.MdsMgr.GetMdsServices()))
	}
}

func TestCreateVBucket(t *testing.T) {
	for i := 0; i < 100; i++ {
		vbucket, location, poolID, mdsID := fmt.Sprintf("vbucket-%d", i), "bj", "pool-beoyx", "mds-wewrp"
		err := mgr.CreateVBucket(vbucket, location, poolID, mdsID)
		if err != nil {
			t.Errorf("create vbucket %s failed", vbucket)
		}
	}
}

func TestGetVBucket(t *testing.T) {
	for i := 0; i < 100; i++ {
		vbucket := fmt.Sprintf("vbucket-%d", i)
		vb := mgr.GetVBucket(vbucket)
		if vb == nil {
			t.Fatalf("get vbucket %s failed", vbucket)
		}
		vb, exists := mgr.VBuckets[vbucket]
		if !exists {
			t.Fatalf("vbucket %s not found", vbucket)
		}
		if vb.Name != vbucket {
			t.Errorf("vbucket %s mismatch with %s", vbucket, vb.Name)
		}
	}
}

func TestListVBuckets(t *testing.T) {
	vbs, err := mgr.ListVBuckets()
	if err != nil {
		t.Fatal("list vbuckets failed")
	}
	utils.PrettyPrint(vbs)
}

func TestDeleteVBucket(t *testing.T) {
	for i := 0; i < 100; i++ {
		vbucket := fmt.Sprintf("vbucket-%d", i)
		err := mgr.DeleteVBucket(vbucket)
		if err != nil {
			t.Fatalf("delete vbucket %s failed", vbucket)
		}
		_, exists := mgr.VBuckets[vbucket]
		if exists {
			t.Errorf("vbucket %s not deleted", vbucket)
		}
	}
	TestListVBuckets(t)
}

func TestShutdown(t *testing.T) {
	mgr.Shutdown()
}
