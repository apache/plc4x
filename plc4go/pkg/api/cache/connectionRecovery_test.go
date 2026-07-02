/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package cache

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/simulated"
	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	"github.com/apache/plc4x/plc4go/pkg/api/config"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func testDriverManager(t *testing.T) plc4go.PlcDriverManager {
	logger := testutils.ProduceTestingLogger(t)
	driverManager := plc4go.NewPlcDriverManager(config.WithCustomLogger(logger))
	t.Cleanup(func() {
		assert.NoError(t, driverManager.Close())
	})
	driverManager.RegisterDriver(simulated.NewDriver(options.WithCustomLogger(logger)))
	return driverManager
}

// A container whose initial connect failed must not ignore later lease requests:
// it has to re-attempt the connection (and deliver an error if that fails again)
// instead of leaving the caller to time out with nothing on either channel.
func Test_connectionContainer_leaseAfterFailedConnectDeliversError(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42?connectionError=nope")

	// Initial connect fails and must leave the container in StateInvalid.
	c.connect(context.Background())
	c.lock.RLock()
	state := c.state
	c.lock.RUnlock()
	assert.Equal(t, StateInvalid, state)

	// A subsequent lease request must trigger a reconnect attempt and get the
	// (failing) result delivered - not silence.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connChan, errChan := c.lease(ctx)
	select {
	case err := <-errChan:
		assert.Error(t, err)
	case conn := <-connChan:
		t.Fatalf("expected an error, got a lease: %v", conn)
	case <-ctx.Done():
		t.Fatal("lease request on invalid container was silently ignored (pre-fix behavior)")
	}
}

// A container in StateInvalid whose endpoint has recovered must reconnect and
// hand out a fresh lease on the next request.
func Test_connectionContainer_leaseOnInvalidReconnects(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42")
	c.lock.Lock()
	c.state = StateInvalid // simulate a previously failed connect/reconnect
	c.lock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connChan, errChan := c.lease(ctx)
	select {
	case conn := <-connChan:
		require.NotNil(t, conn)
		assert.True(t, conn.IsConnected())
		assert.NoError(t, conn.Close())
	case err := <-errChan:
		t.Fatalf("expected a lease, got error: %v", err)
	case <-ctx.Done():
		t.Fatal("lease request on invalid container was silently ignored (pre-fix behavior)")
	}
}

// An idle cached connection that died while parked in the cache must not be
// handed out; the container has to establish a replacement.
func Test_connectionContainer_deadIdleConnectionReplaced(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42")
	c.connect(context.Background())
	c.lock.RLock()
	state, firstConnection := c.state, c.connection
	c.lock.RUnlock()
	require.Equal(t, StateIdle, state)
	require.NotNil(t, firstConnection)

	// Simulate the remote dropping the idle connection.
	firstConnection.Invalidate()
	require.False(t, firstConnection.IsConnected())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connChan, errChan := c.lease(ctx)
	select {
	case conn := <-connChan:
		require.NotNil(t, conn)
		assert.True(t, conn.IsConnected(), "handed-out lease must be alive")
		assert.NoError(t, conn.Close())
	case err := <-errChan:
		t.Fatalf("expected a lease, got error: %v", err)
	case <-ctx.Done():
		t.Fatal("timed out waiting for a replacement connection")
	}

	c.lock.RLock()
	secondConnection := c.connection
	c.lock.RUnlock()
	assert.NotSame(t, firstConnection, secondConnection, "dead connection must have been replaced")
}

// Queued lease requests whose caller already gave up must be skipped when a
// connection is handed out - otherwise the connection is leased to nobody and
// stays StateInUse forever.
func Test_connectionContainer_cancelledWaiterSkipped(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42")
	c.connect(context.Background())

	// Take the connection so follow-up requests queue.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connChan, _ := c.lease(ctx)
	lease := <-connChan

	// Queue a request and abandon it.
	abandonedCtx, abandonedCancel := context.WithCancel(context.Background())
	abandonedConnChan, abandonedErrChan := c.lease(abandonedCtx)
	abandonedCancel()

	// Queue a live request behind it.
	liveConnChan, liveErrChan := c.lease(ctx)

	// Returning the lease must skip the abandoned waiter (failing it explicitly)
	// and serve the live one.
	require.NoError(t, lease.Close())

	select {
	case conn := <-liveConnChan:
		require.NotNil(t, conn)
		assert.NoError(t, conn.Close())
	case err := <-liveErrChan:
		t.Fatalf("live waiter got error: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("live waiter was starved - connection was likely leased to the abandoned waiter")
	}

	// The abandoned waiter must have received an explicit error, not a lease.
	select {
	case err := <-abandonedErrChan:
		assert.Error(t, err)
	case conn := <-abandonedConnChan:
		t.Fatalf("abandoned waiter got a lease: %v", conn)
	case <-time.After(5 * time.Second):
		t.Fatal("abandoned waiter was never notified")
	}
}

// End-to-end through the cache: a connection string whose first connect fails
// must keep yielding prompt errors on every attempt - not one error and then
// silent timeouts (pre-fix behavior: the container stayed StateInitialized with
// requests queueing forever, or StateInvalid with requests ignored).
func Test_plcConnectionCache_failedFirstConnectDoesNotBrick(t *testing.T) {
	cache := NewPlcConnectionCache(testDriverManager(t))
	t.Cleanup(func() {
		assert.NoError(t, cache.Close())
	})

	for i := range 3 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		start := time.Now()
		_, err := cache.GetConnection(ctx, "simulated://1.2.3.4:42?connectionError=nope")
		cancel()
		assert.Error(t, err, "attempt %d", i+1)
		assert.NotErrorIs(t, err, context.DeadlineExceeded, "attempt %d must fail fast with the connect error, not time out", i+1)
		assert.Less(t, time.Since(start), 3*time.Second, "attempt %d must not run into the timeout", i+1)
	}
}

// Concurrent lease requests against an invalid container must result in exactly
// one reconnect attempt with everyone served (leases handed out sequentially).
func Test_connectionContainer_concurrentLeaseOnInvalid(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42")
	c.lock.Lock()
	c.state = StateInvalid
	c.lock.Unlock()

	const clients = 5
	var wg sync.WaitGroup
	for range clients {
		wg.Go(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			connChan, errChan := c.lease(ctx)
			select {
			case conn := <-connChan:
				assert.NoError(t, conn.Close())
			case err := <-errChan:
				t.Errorf("expected a lease, got error: %v", err)
			case <-ctx.Done():
				t.Error("timed out waiting for lease")
			}
		})
	}
	wg.Wait()
}
