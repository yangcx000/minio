/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package pool

import (
	"testing"

	"github.com/minio/minio/fusionstore/mgs"
)

const (
	mgsAddr      = "192.168.58.4:8000"
	timeoutValue = 10
)

var testGlobalMgr *Mgr

func newMgsService() error {
	return mgs.NewService(mgsAddr, timeoutValue)
}

func TestNewMgr(t *testing.T) {
	if newMgsService() != nil {
		t.Fatalf("Couldn't new mgs service %s", mgsAddr)
	}
	var err error
	testGlobalMgr, err = NewMgr()
	if err != nil {
		t.Fatalf("Init Pool mgr failed, err: %s", err)
	}
}

func TestAllocBucket(t *testing.T) {
	TestNewMgr(t)
	pool := "pool-zzbot"
	bucketNames := make(map[string]int)
	for i := 0; i < 1000; i++ {
		bucketName := testGlobalMgr.AllocBucket(pool)
		bucketNames[bucketName]++
	}
	//utils.PrettyPrint(bucketNames)
}
