// Copyright (C) 2019-2021, Dijets Inc, All rights reserved.
// See the file LICENSE for licensing terms.

package throttling

import (
	"context"
	"sync"
	"testing"

	"github.com/lasthyphen/DijetsNetworkBinary/ids"
	"github.com/lasthyphen/DijetsNetworkBinary/utils/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestBandwidthThrottler(t *testing.T) {
	assert := assert.New(t)
	// Assert initial state
	config := BandwidthThrottlerConfig{
		RefillRate:   8,
		MaxBurstSize: 10,
	}
	throttlerIntf, err := newBandwidthThrottler(logging.NoLog{}, "", prometheus.NewRegistry(), config)
	assert.NoError(err)
	throttler, ok := throttlerIntf.(*bandwidthThrottlerImpl)
	assert.True(ok)
	assert.NotNil(throttler.log)
	assert.NotNil(throttler.limiters)
	assert.EqualValues(throttler.RefillRate, 8)
	assert.EqualValues(throttler.MaxBurstSize, 10)
	assert.Len(throttler.limiters, 0)

	// Add a node
	nodeID1 := ids.GenerateTestNodeID()
	throttler.AddNode(nodeID1)
	assert.Len(throttler.limiters, 1)

	// Remove the node
	throttler.RemoveNode(nodeID1)
	assert.Len(throttler.limiters, 0)

	// Add the node back
	throttler.AddNode(nodeID1)
	assert.Len(throttler.limiters, 1)

	// Should be able to acquire 8
	throttler.Acquire(context.Background(), 8, nodeID1)

	// Make several goroutines that acquire bytes.
	wg := sync.WaitGroup{}
	wg.Add(int(config.MaxBurstSize) + 5)
	for i := uint64(0); i < config.MaxBurstSize+5; i++ {
		go func() {
			throttler.Acquire(context.Background(), 1, nodeID1)
			wg.Done()
		}()
	}
	wg.Wait()
}
