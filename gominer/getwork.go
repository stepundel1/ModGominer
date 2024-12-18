// Copyright (c) 2016-2023 The Decred developers.

package main

import (
	"encoding/json"
	"fmt"

	"github.com/decred/gominer/stratum"
	"github.com/decred/gominer/work"
)

// GetPoolWork gets work from a stratum enabled pool.
func GetPoolWork(pool *stratum.Stratum) (*work.Work, error) {
	// Get Next work for stratum and mark it as used.
	if pool.PoolWork.NewWork {
		poolLog.Debug("Received new work from pool.")
		// Mark used.
		pool.PoolWork.NewWork = false

		if pool.PoolWork.JobID == "" {
			return nil, fmt.Errorf("no work available (no job id)")
		}

		err := pool.PrepWork()
		if err != nil {
			return nil, err
		}

		poolLog.Debugf("new job %q height %v", pool.PoolWork.JobID,
			pool.PoolWork.Height)

		return pool.PoolWork.Work, nil
	}

	// Return the work we already had, do not recalculate
	if pool.PoolWork.Work != nil {
		return pool.PoolWork.Work, nil
	}

	return nil, fmt.Errorf("no work available")
}

// GetPoolWorkSubmit sends the result to the stratum enabled pool.
func GetPoolWorkSubmit(pool *stratum.Stratum, worker string, jobID string, nonce string, ntimeHex string, extraNonce1 string, nonce2 uint32) (bool, error) {
	pool.Lock()
	defer pool.Unlock()

	params := []string{
		worker,
		jobID,
		fmt.Sprintf("%x", nonce),
		ntimeHex,
		extraNonce1,
		fmt.Sprintf("%x", nonce2),
	}

	sub := map[string]interface{}{
		"method": "mining.submit",
		"params": params,
	}

	// JSON encode.
	m, err := json.Marshal(sub)
	if err != nil {
		return false, err
	}

	// Send.
	poolLog.Tracef("%s", m)
	_, err = pool.Conn.Write(m)
	if err != nil {
		return false, err
	}
	_, err = pool.Conn.Write([]byte("\n"))
	if err != nil {
		return false, err
	}

	return true, nil
}
