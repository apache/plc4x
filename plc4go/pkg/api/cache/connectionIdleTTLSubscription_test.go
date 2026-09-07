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

// stubSubscribedConnection implements tracedPlcConnection via interface
// embedding (unimplemented methods panic if called) and reports a fixed
// number of active subscription handles.
type stubSubscribedConnection struct {
	tracedPlcConnection
	activeSubscriptions int
}

func (s *stubSubscribedConnection) ActiveSubscriptionCount() int { return s.activeSubscriptions }
func (s *stubSubscribedConnection) IsConnected() bool            { return true }
func (s *stubSubscribedConnection) Close() error                 { return nil }
func (s *stubSubscribedConnection) IsTraceEnabled() bool         { return false }

// A connection carrying active subscription handles must never be reaped by
// the idle TTL: its subscription state (e.g. BACnet COV refresh timers) lives
// on the connection and would be silently destroyed, cutting off passive
// updates until the client happens to re-subscribe.
func Test_connectionContainer_idleExpiredSkipsSubscribedConnection(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42")
	c.maxIdleTime = time.Minute

	stub := &stubSubscribedConnection{activeSubscriptions: 1}
	c.lock.Lock()
	c.connection = stub
	c.state = StateIdle
	c.idleSince = time.Now().Add(-2 * time.Minute) // aged past the TTL
	c.lock.Unlock()

	c.lock.RLock()
	expired := c.idleExpired()
	c.lock.RUnlock()
	assert.False(t, expired, "aged connection with active subscriptions must not be considered expired")

	// The same aged connection without subscriptions expires as usual.
	stub.activeSubscriptions = 0
	c.lock.RLock()
	expired = c.idleExpired()
	c.lock.RUnlock()
	assert.True(t, expired, "aged connection without subscriptions must expire")
}

// Lease-level behavior: an aged-but-subscribed connection is handed out as-is
// instead of being replaced.
func Test_connectionContainer_leaseKeepsAgedSubscribedConnection(t *testing.T) {
	c := newConnectionContainer(testutils.ProduceTestingLogger(t), testDriverManager(t), "simulated://1.2.3.4:42")
	c.maxIdleTime = time.Minute

	stub := &stubSubscribedConnection{activeSubscriptions: 3}
	c.lock.Lock()
	c.connection = stub
	c.state = StateIdle
	c.idleSince = time.Now().Add(-2 * time.Minute)
	c.lock.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	connChan, errChan := c.lease(ctx)
	select {
	case lease := <-connChan:
		require.NotNil(t, lease)
		c.lock.RLock()
		assert.Same(t, stub, c.connection, "subscribed connection must be kept across the lease")
		c.lock.RUnlock()
	case err := <-errChan:
		t.Fatalf("expected a lease, got error: %v", err)
	case <-ctx.Done():
		t.Fatal("timed out waiting for lease")
	}
}
