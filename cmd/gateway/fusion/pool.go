/*
 * Copyright 2021 LiAuto authors.
 * @yangchunxin
 *
 */

package fusion

import "time"

type poolStatus int

const (
	poolStatusUnknown poolStatus = iota
	poolStatusReady
	poolStatusFull
	poolStatusFailed
)

// contains AK/SK
type poolCredentials struct {
	accessKey  string
	secretKey  string
	expiration time.Time
}

// object storage pool
type storagePool struct {
	id       string
	name     string
	endpoint string
	status   poolStatus
	creds    poolCredentials
}
