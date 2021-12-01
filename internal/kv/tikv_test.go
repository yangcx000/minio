/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package kv

import (
	"bytes"
	"testing"
)

var tikvClient *TiKV

func setup() error {
	pds := []string{"172.21.63.21:2379", "172.21.63.27:2379", "172.21.63.28:2379"}
	tikvClient = NewTiKV(pds)
	err := tikvClient.Open()
	if err != nil {
		return err
	}
	return nil
}

func teardown() {
	tikvClient.Close()
}

func TestPut(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatalf("setup failed, err:%s", err.Error())
	}
	key := []byte("ycx/test-001")
	value := []byte("hello 世界")
	err = tikvClient.Put(key, value)
	if err != nil {
		t.Errorf("Couldn't put key %s, err:%s", key, err.Error())
	}
	teardown()
}

func TestGet(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatalf("setup failed, err:%s", err.Error())
	}
	key := []byte("ycx/test-001")
	valueExpect := []byte("hello 世界")
	value, err := tikvClient.Get([]byte(key))
	if err != nil {
		t.Errorf("Couldn't get key %s, err:%s", key, err.Error())
		return
	}
	if !bytes.Equal(value, valueExpect) {
		t.Errorf("Get not equal, key:%s, value:%s, expected:%s", key, value, valueExpect)
	}
	teardown()
}

func TestDelete(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatalf("setup failed, err:%s", err.Error())
	}
	key := []byte("ycx/test-001")
	err = tikvClient.Delete(key)
	if err != nil {
		t.Errorf("Couldn't delete key %s, err:%s", key, err.Error())
	}
	value, err := tikvClient.Get(key)
	if err != nil {
		t.Errorf("Couldn't get key %s, err:%s", key, err.Error())
	}
	if value != nil {
		t.Errorf("Couldn't delete key %s", key)
	}
	teardown()
}

func TestBatchPut(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatalf("setup failed, err:%s", err.Error())
	}
	keys := [][]byte{[]byte("ycx/test-001"), []byte("ycx/test-002"), []byte("ycx/test-003")}
	values := [][]byte{[]byte("hello 世界1"), []byte("hello 世界2"), []byte("hello 世界3")}
	err = tikvClient.BatchPut(keys, values)
	if err != nil {
		t.Errorf("Couldn't batch put keys %+v, err:%s", keys, err.Error())
	}
	teardown()
}

func TestBatchGet(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatalf("setup failed, err:%s", err.Error())
	}
	keys := [][]byte{[]byte("ycx/test-001"), []byte("ycx/test-002"), []byte("ycx/test-003")}
	valuesExpect := [][]byte{[]byte("hello 世界1"), []byte("hello 世界2"), []byte("hello 世界3")}
	values, err := tikvClient.BatchGet(keys)
	if err != nil {
		t.Errorf("Couldn't batch get keys %+v, err:%s", keys, err.Error())
		return
	}
	for i := range values {
		if !bytes.Equal(values[i], valuesExpect[i]) {
			t.Errorf("BatchGet not equal, keys:%+v, values:%+v, expected:%+v", keys, values, valuesExpect)
		}
	}
	teardown()
}

func TestBatchDelete(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatalf("setup failed, err:%s", err.Error())
	}
	keys := [][]byte{[]byte("ycx/test-001"), []byte("ycx/test-002"), []byte("ycx/test-003")}
	err = tikvClient.BatchDelete(keys)
	if err != nil {
		t.Errorf("Couldn't batch delete keys %+v, err:%s", keys, err.Error())
	}
	values, err := tikvClient.BatchGet(keys)
	if err != nil {
		t.Errorf("Couldn't batch get keys %+v, err:%s", keys, err.Error())
	}
	if values != nil {
		for _, v := range values {
			if v != nil {
				t.Errorf("Couldn't batch delete keys %+v, v:%s", keys, v)
			}
		}
	}
	teardown()
}

func TestScanRange(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatalf("setup failed, err:%s", err.Error())
	}

	// add a...c
	keys := [][]byte{[]byte("a/test-001"), []byte("b/test-002"), []byte("c/test-003")}
	values := [][]byte{[]byte("hello 世界1"), []byte("hello 世界2"), []byte("hello 世界3")}
	err = tikvClient.BatchPut(keys, values)
	if err != nil {
		t.Errorf("Couldn't batch put keys %+v, err:%s", keys, err.Error())
	}

	// scan a..c
	skeys, svalues, err := tikvClient.Scan([]byte("a"), []byte("c"), 0)
	if err != nil {
		t.Errorf("Couldn't scan, err:%s", err.Error())
	}
	for i, key := range skeys {
		if !bytes.Equal(key, keys[i]) {
			t.Errorf("scan key not equal, key:%s, expected:%s", key, keys[i])
		}
	}
	for i, value := range svalues {
		if !bytes.Equal(value, values[i]) {
			t.Errorf("scan value not equal, value:%s, expected:%s", value, values[i])
		}
	}
	teardown()
}

func TestDeleteRange(t *testing.T) {
	err := setup()
	if err != nil {
		t.Fatalf("setup failed, err:%s", err.Error())
	}

	// delete a..c
	err = tikvClient.DeleteRange([]byte("a"), []byte("c"))
	if err != nil {
		t.Errorf("Couldn't delete range, err:%s", err.Error())
	}
	// scan a..c
	skeys, svalues, err := tikvClient.Scan([]byte("a"), []byte("c"), 0)
	if err != nil {
		t.Errorf("Couldn't scan, err:%s", err.Error())
	}
	for i := range skeys {
		if skeys[i] != nil {
			t.Errorf("Couldn't delete range, key:%s", skeys[i])
		}
	}
	for i := range svalues {
		if svalues[i] != nil {
			t.Errorf("Couldn't delete range, value:%s", svalues[i])
		}
	}
	teardown()
}
