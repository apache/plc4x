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

package bacnetip

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// newRefreshSubscriber returns a Subscriber whose connection sends through
// codec, with one stored handle subscribed for lifetime seconds at issuedAt.
func newRefreshSubscriber(t *testing.T, codec *captureCodec, lifetime uint32, issuedAt time.Time) (*Subscriber, *SubscriptionHandle) {
	t.Helper()
	conn := &Connection{configuration: createDefaultConfiguration(), messageCodec: codec}
	conn.driverContext = NewDriverContext(conn.configuration)
	s := NewSubscriber(conn)
	objType := readWriteModel.BACnetObjectType_ANALOG_INPUT
	handle := NewSubscriptionHandle(s, "ai1", &plcTag{ObjectId: objectId{ObjectIdType: &objType, ObjectIdInstance: 1}}, apiModel.SubscriptionChangeOfState, 0)
	handle.subscriberProcessId = 7
	handle.lifetimeSec = lifetime
	handle.lastIssuedAt = issuedAt
	s.storeHandle(handle)
	return s, handle
}

// refreshOnce runs refreshDue at now and answers a SubscribeCOV it sends with
// a SimpleAck (ack) or a timeout. It reports whether one was sent.
func refreshOnce(t *testing.T, s *Subscriber, codec *captureCodec, now time.Time, ack bool) bool {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.refreshDue(testutils.TestContext(t), now)
	}()
	select {
	case <-done:
		return false
	case h := <-codec.handlers:
		if ack {
			require.NoError(t, h.handleMessage(wrapAPDU(readWriteModel.NewAPDUSimpleAck(0, readWriteModel.BACnetConfirmedServiceChoice_SUBSCRIBE_COV), false, nil)))
		} else {
			require.NoError(t, h.handleError(errors.New("timeout")))
		}
	case <-time.After(time.Second):
		t.Fatal("refreshDue neither returned nor sent a SubscribeCOV")
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("refreshDue did not return after the reply")
	}
	return true
}

// A device drops a COV subscription when its lifetime runs out, so the
// subscriber must re-issue it at lifetime/2, and only then.
func TestSubscriber_RefreshDueReissuesAtHalfTheLifetime(t *testing.T) {
	codec := newCaptureCodec(nil)
	issued := time.Now()
	s, handle := newRefreshSubscriber(t, codec, 600, issued)

	assert.False(t, refreshOnce(t, s, codec, issued.Add(299*time.Second), true), "not due before lifetime/2")
	assert.True(t, refreshOnce(t, s, codec, issued.Add(300*time.Second), true), "due at lifetime/2")
	handle.mu.Lock()
	assert.True(t, handle.lastIssuedAt.After(issued), "a successful refresh restarts the lifetime")
	handle.mu.Unlock()
}

// A failed refresh is retried after a tenth of the lifetime, not on every tick.
func TestSubscriber_RefreshDueBacksOffAfterAFailure(t *testing.T) {
	codec := newCaptureCodec(nil)
	issued := time.Now().Add(-time.Hour)
	s, _ := newRefreshSubscriber(t, codec, 600, issued)
	now := time.Now()

	require.True(t, refreshOnce(t, s, codec, now, false), "due, sent and timed out")
	assert.False(t, refreshOnce(t, s, codec, now.Add(59*time.Second), true), "held back after the failure")
	assert.True(t, refreshOnce(t, s, codec, now.Add(60*time.Second), true), "retried after lifetime/10")
}

// Indefinite (lifetime 0) and cancelled subscriptions are never refreshed:
// refreshing a cancelled one would re-create what Unsubscribe removed.
func TestSubscriber_RefreshDueSkipsIndefiniteAndCancelled(t *testing.T) {
	codec := newCaptureCodec(nil)
	issued := time.Now().Add(-time.Hour)

	s, _ := newRefreshSubscriber(t, codec, 0, issued)
	assert.False(t, refreshOnce(t, s, codec, time.Now(), true), "lifetime 0 never expires")

	s, handle := newRefreshSubscriber(t, codec, 600, issued)
	handle.cancelled = true
	assert.False(t, refreshOnce(t, s, codec, time.Now(), true), "cancelled handle is not refreshed")
}
