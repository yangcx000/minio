/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package mds

import (
	"testing"

	"github.com/minio/minio/fusionstore/object"
)

const (
	mdsAddr = "192.168.58.4:9010"
)

var testGlobalSvc *Service

func TestNewService(t *testing.T) {
	var err error
	testGlobalSvc, err = NewService(mdsAddr, timeoutValue)
	if err != nil {
		t.Fatalf("NewService failed, %s", err)
	}
}

/*
func TestPutObject(t *testing.T) {
	for i := 2000; i < 2010; i++ {
		obj := &object.Object{
			ID:          fmt.Sprintf("object-%d", i),
			Name:        fmt.Sprintf("a/object-name-%d", i),
			VBucket:     "vbucket-001",
			Pool:        "pool-001",
			Bucket:      "bucket-001",
			ETag:        "etag-001",
			InnerEtag:   "inner-etag-001",
			VersionID:   "v1",
			ContentType: "stream",
		}
		resp, err := testGlobalSvc.PutObject(obj)
		if err != nil {
			t.Fatalf("PutObject failed, err:%s", err)
		}
		if resp.GetStatus().GetCode() != 0 {
			t.Fatal("PutObject code failed")
		}
	}
}

func TestQueryObject(t *testing.T) {
	vbucket, objName := "vbucket-001", "object-name-001"
	resp, err := testGlobalSvc.QueryObject(vbucket, objName)
	if err != nil {
		t.Fatalf("QueryObject failed, err:%s", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatalf("QueryObject code failed, code:%d", resp.GetStatus().GetCode())
	}
	utils.PrettyPrint(resp.GetObject())
}
*/

func TestDeleteObject(t *testing.T) {
	//vbucket, objName := "vbucket-001", "object-name-001"
	lop := &object.ListObjectsParam{
		VBucket: "vbucket-001",
		Prefix:  "",
		//Marker:    "a/b/c/d/object-name-18",
		Delimiter: "",
		Limits:    1000,
	}
	resp, err := testGlobalSvc.ListObjects(lop)
	if err != nil {
		t.Fatalf("ListObjects failed, err:%s", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatal("ListObjects code failed")
	}
	for _, obj := range resp.GetObjects() {
		resp, err := testGlobalSvc.DeleteObject(obj.Vbucket, obj.Name)
		if err != nil {
			t.Fatalf("DeleteObject failed, err:%s", err)
		}
		if resp.GetStatus().GetCode() != 0 {
			t.Fatal("DeleteObject code failed")
		}
	}
}

/*
func TestListObjects(t *testing.T) {
	lop := &object.ListObjectsParam{
		VBucket: "vbucket-001",
		Prefix:  "",
		//Marker:    "a/b/c/d/object-name-18",
		Delimiter: "",
		Limits:    300,
	}
	resp, err := testGlobalSvc.ListObjects(lop)
	if err != nil {
		t.Fatalf("ListObjects failed, err:%s", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatal("ListObjects code failed")
	}
	if len(resp.GetObjects()) != 200 {
		t.Errorf("ListObjects count %d != 200", len(resp.GetObjects()))
	}
	//utils.PrettyPrint(resp)
}

func TestCreateMultipart(t *testing.T) {
	mp := &multipart.Multipart{
		PhysicalUploadID: "phy-upload-id-00001",
		Pool:             "pool-1",
		VBucket:          "vbucket-001",
		PhysicalBucket:   "phy-bucket-001",
		Object:           "obj-0001",
		CreatedTime:      time.Now(),
	}
	resp, err := testGlobalSvc.CreateMultipart(mp)
	if err != nil {
		t.Fatalf("CreateMultipart failed, err:%s", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatal("CreateMultipart code failed")
	}
}

func TestQueryMultipart(t *testing.T) {
	vbucket := "vbucket-001"
	uploadID := "upload-id-00001"
	resp, err := testGlobalSvc.QueryMultipart(vbucket, uploadID)
	if err != nil {
		t.Fatalf("QueryMultipart failed, err:%s", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatal("QueryMultipart code failed")
	}
	utils.PrettyPrint(resp.GetMultipart())
}

func TestDeleteMultipart(t *testing.T) {
	vbucket := "vbucket-001"
	uploadID := "upload-id-00001"
	resp, err := testGlobalSvc.DeleteMultipart(vbucket, uploadID)
	if err != nil {
		t.Fatalf("DeleteMultipart failed, err:%s", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatal("DeleteMultipart code failed")
	}
}

func TestListMultiparts(t *testing.T) {
	vbucket := "vbucket-001"
	marker := ""
	limits := 0
	resp, err := testGlobalSvc.ListMultiparts(vbucket, marker, int32(limits))
	if err != nil {
		t.Fatalf("ListMultiparts failed, err:%s", err)
	}
	if resp.GetStatus().GetCode() != 0 {
		t.Fatal("ListMultiparts code failed")
	}
	utils.PrettyPrint(resp.GetMultiparts())
}
*/

func TestClose(t *testing.T) {
	testGlobalSvc.Close()
}
