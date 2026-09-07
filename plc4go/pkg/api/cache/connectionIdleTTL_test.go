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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// leaseWithTimeout requests a lease and waits for the outcome.
func leaseWithTimeout(t *testing.T, c *connectionContainer) *plcConnectionLease {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connChan, errChan := c.lease(ctx)
	select {
	case conn := <-connChan:
		return conn
	case err := <-errChan:
		t.Fatalf("expected a lease, got error: %v", err)
	case <-ctx.Done():
		t.Fatal("timed out waiting for lease")
	}
	return nil
}

// An idle connection older than maxIdleTime must not be handed out - the
// container has to discard it and establish a replacement. Remotes that
// silently reap idle connections leave a half-open socket that IsConnected()
// cannot detect; past the configured age the connection is not to be trusted.
func Test_connectionContainer_idleExpiredConnectionReplaced(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42")
	c.maxIdleTime = time.Minute
	c.connect(context.Background())
	c.lock.RLock()
	firstConnection := c.connection
	c.lock.RUnlock()
	require.NotNil(t, firstConnection)

	// Age the idle connection past the TTL.
	c.lock.Lock()
	c.idleSince = time.Now().Add(-2 * time.Minute)
	c.lock.Unlock()

	lease := leaseWithTimeout(t, c)
	require.NotNil(t, lease)
	assert.True(t, lease.IsConnected(), "handed-out lease must be alive")
	assert.NoError(t, lease.Close())

	c.lock.RLock()
	secondConnection := c.connection
	c.lock.RUnlock()
	assert.NotSame(t, firstConnection, secondConnection, "expired idle connection must have been replaced")
}

// A connection younger than maxIdleTime is reused as-is.
func Test_connectionContainer_idleNotExpiredKeepsConnection(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42")
	c.maxIdleTime = time.Hour
	c.connect(context.Background())
	c.lock.RLock()
	firstConnection := c.connection
	c.lock.RUnlock()
	require.NotNil(t, firstConnection)

	lease := leaseWithTimeout(t, c)
	require.NotNil(t, lease)
	assert.NoError(t, lease.Close())

	c.lock.RLock()
	secondConnection := c.connection
	c.lock.RUnlock()
	assert.Same(t, firstConnection, secondConnection, "young idle connection must be reused")
}

// maxIdleTime = 0 (the default) disables the TTL entirely - connections are
// kept forever, preserving the previous behavior.
func Test_connectionContainer_maxIdleTimeDisabledByDefault(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42")
	c.connect(context.Background())
	c.lock.RLock()
	firstConnection := c.connection
	c.lock.RUnlock()
	require.NotNil(t, firstConnection)

	// Even an ancient idle timestamp must not trigger a replacement.
	c.lock.Lock()
	c.idleSince = time.Now().Add(-24 * time.Hour)
	c.lock.Unlock()

	lease := leaseWithTimeout(t, c)
	require.NotNil(t, lease)
	assert.NoError(t, lease.Close())

	c.lock.RLock()
	secondConnection := c.connection
	c.lock.RUnlock()
	assert.Same(t, firstConnection, secondConnection, "TTL disabled: connection must be kept")
}

// Returning a lease must restart the idle clock.
func Test_connectionContainer_returnResetsIdleClock(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42")
	c.maxIdleTime = time.Minute
	c.connect(context.Background())

	// Age the container, then take and return a lease - the return must
	// refresh idleSince so the connection is trusted again.
	c.lock.Lock()
	c.idleSince = time.Now().Add(-2 * time.Minute)
	c.lock.Unlock()

	lease := leaseWithTimeout(t, c) // triggers replacement (expired)
	require.NotNil(t, lease)
	require.NoError(t, lease.Close())

	c.lock.RLock()
	idleSince := c.idleSince
	connection := c.connection
	c.lock.RUnlock()
	assert.WithinDuration(t, time.Now(), idleSince, 10*time.Second, "idle clock must restart on return")

	// An immediate follow-up lease reuses the same connection.
	lease2 := leaseWithTimeout(t, c)
	require.NotNil(t, lease2)
	assert.NoError(t, lease2.Close())
	c.lock.RLock()
	assert.Same(t, connection, c.connection)
	c.lock.RUnlock()
}

// End-to-end through the cache: WithMaxIdleTime propagates to containers and
// an aged connection is transparently replaced on the next GetConnection.
func Test_plcConnectionCache_withMaxIdleTime(t *testing.T) {
	cache := NewPlcConnectionCache(testDriverManager(t), WithMaxIdleTime(time.Minute)).(*plcConnectionCache)
	t.Cleanup(func() {
		assert.NoError(t, cache.Close())
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := cache.GetConnection(ctx, "simulated://1.2.3.4:42")
	require.NoError(t, err)
	require.NoError(t, conn.Close())

	cache.cacheLock.RLock()
	container := cache.connections["simulated://1.2.3.4:42"]
	cache.cacheLock.RUnlock()
	require.NotNil(t, container)
	assert.Equal(t, time.Minute, container.maxIdleTime, "option must propagate to containers")

	// Age the pooled connection and lease again - must succeed on a fresh one.
	container.lock.Lock()
	firstConnection := container.connection
	container.idleSince = time.Now().Add(-2 * time.Minute)
	container.lock.Unlock()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	conn2, err := cache.GetConnection(ctx2, "simulated://1.2.3.4:42")
	require.NoError(t, err)
	assert.True(t, conn2.IsConnected())
	require.NoError(t, conn2.Close())

	container.lock.RLock()
	secondConnection := container.connection
	container.lock.RUnlock()
	assert.NotSame(t, firstConnection, secondConnection, "aged connection must have been replaced")
}
