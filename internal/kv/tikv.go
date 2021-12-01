/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 */

package kv

import (
	"fmt"

	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/tikv"
)

// TiKV xxx
type TiKV struct {
	pds    []string
	client *tikv.RawKVClient
}

// NewTiKV xxxx
func NewTiKV(pds []string) *TiKV {
	return &TiKV{pds: pds}
}

// Open xxx
func (t *TiKV) Open() error {
	var err error
	t.client, err = tikv.NewRawKVClient(t.pds, config.DefaultConfig().Security)
	if err != nil {
		return fmt.Errorf("couldn't create tikv raw client, %w", err)
	}
	return nil
}

// Close xxx
func (t *TiKV) Close() error {
	err := t.client.Close()
	if err != nil {
		return fmt.Errorf("couldn't close tikv raw client, %w", err)
	}
	return nil
}

// Put xxx
func (t *TiKV) Put(key, value []byte) error {
	err := t.client.Put(key, value)
	if err != nil {
		return fmt.Errorf("couldn't put key %s, %w", key, err)
	}
	return nil
}

// BatchPut xxx
func (t *TiKV) BatchPut(keys, values [][]byte) error {
	err := t.client.BatchPut(keys, values)
	if err != nil {
		return fmt.Errorf("couldn't batch put keys, %w", err)
	}
	return nil
}

// Get xxx
func (t *TiKV) Get(key []byte) ([]byte, error) {
	value, err := t.client.Get(key)
	if err != nil {
		return nil, fmt.Errorf("couldn't get key %s, %w", key, err)
	}
	return value, nil
}

// BatchGet xxx
func (t *TiKV) BatchGet(keys [][]byte) ([][]byte, error) {
	values, err := t.client.BatchGet(keys)
	if err != nil {
		return nil, fmt.Errorf("couldn't batch get keys, %w", err)
	}
	return values, nil
}

// Delete xxx
func (t *TiKV) Delete(key []byte) error {
	err := t.client.Delete(key)
	if err != nil {
		return fmt.Errorf("couldn't delete key %s, %w", key, err)
	}
	return nil
}

// BatchDelete xxx
func (t *TiKV) BatchDelete(keys [][]byte) error {
	err := t.client.BatchDelete(keys)
	if err != nil {
		return fmt.Errorf("couldn't batch delete keys, %w", err)
	}
	return nil
}

// DeleteRange xxx
func (t *TiKV) DeleteRange(start, end []byte) error {
	err := t.client.DeleteRange(start, end)
	if err != nil {
		return fmt.Errorf("couldn't delete range, start:%s, end:%s, %w", start, end, err)
	}
	return nil
}

// Scan xxx
func (t *TiKV) Scan(start, end []byte, limit int) ([][]byte, [][]byte, error) {
	keys, values, err := t.client.Scan(start, end, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't scan, start:%s, end:%s, %w", start, end, err)
	}
	return keys, values, err
}

// Exists xxx
func (t *TiKV) Exists(key []byte) (bool, error) {
	value, err := t.Get(key)
	if err != nil {
		return false, err
	}
	if value != nil {
		return true, nil
	}
	return false, nil
}
