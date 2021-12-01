/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package kv

const (
	ScanLimits = 100
)

// Store xxx
type Store interface {
	Open() error
	Close() error
	Put(key, value []byte) error
	BatchPut(keys, values [][]byte) error
	Get(key []byte) ([]byte, error)
	BatchGet(keys [][]byte) ([][]byte, error)
	Delete(key []byte) error
	BatchDelete(keys [][]byte) error
	DeleteRange(start, end []byte) error
	Scan(start, end []byte, limit int) ([][]byte, [][]byte, error)
	Exists(key []byte) (bool, error)
}
