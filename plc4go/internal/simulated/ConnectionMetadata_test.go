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

package simulated

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testConnection() *Connection {
	return &Connection{
		device:       NewDevice("hurz"),
		tagHandler:   NewTagHandler(),
		valueHandler: NewValueHandler(),
		options:      map[string][]string{},
		connected:    true,
	}
}

// assertAdvertisedCapabilityHasBuilder checks the direction that actually matters: a connection
// that advertises a capability has to hand out a working builder for it. The converse is
// deliberately not asserted - a driver may hand out a builder whose operations report
// "not implemented" while honestly advertising the capability as absent, which is what the
// simulated connection does for subscribing.
func assertAdvertisedCapabilityHasBuilder(t *testing.T, capability string, advertised bool, builder func() any) {
	t.Helper()
	if !advertised {
		return
	}
	var handedOut any
	panicked := func() (panicked bool) {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		handedOut = builder()
		return false
	}()
	assert.False(t, panicked, "%s is advertised but the builder panicked", capability)
	assert.NotNil(t, handedOut, "%s is advertised but no builder was handed out", capability)
}

func TestConnection_AdvertisedCapabilitiesHaveBuilders(t *testing.T) {
	c := testConnection()
	metadata := c.GetMetadata()

	assert.True(t, metadata.CanRead())
	assert.True(t, metadata.CanWrite())
	assert.False(t, metadata.CanSubscribe(), "the simulated subscriber is a stub")
	assert.False(t, metadata.CanBrowse(), "the simulated connection has no browser")

	assertAdvertisedCapabilityHasBuilder(t, "reading", metadata.CanRead(), func() any { return c.ReadRequestBuilder() })
	assertAdvertisedCapabilityHasBuilder(t, "writing", metadata.CanWrite(), func() any { return c.WriteRequestBuilder() })
	assertAdvertisedCapabilityHasBuilder(t, "subscribing", metadata.CanSubscribe(), func() any { return c.SubscriptionRequestBuilder() })
	assertAdvertisedCapabilityHasBuilder(t, "browsing", metadata.CanBrowse(), func() any { return c.BrowseRequestBuilder() })
}

// TestConnection_UnadvertisedSubscribingDoesNotWork pins the other half of the simulated
// connection's honesty: it hands out a subscription builder while advertising no subscribing,
// which is only truthful as long as the Subscriber behind that builder refuses every request.
// Wiring in a Subscriber that works will fail here, forcing ProvidesSubscribing to be flipped
// along with it.
func TestConnection_UnadvertisedSubscribingDoesNotWork(t *testing.T) {
	c := testConnection()
	require.False(t, c.GetMetadata().CanSubscribe())

	subscriber := NewSubscriber(c.device, c.options, c.tracer)
	subscribeResult := <-subscriber.Subscribe(context.Background(), nil)
	assert.Error(t, subscribeResult.GetErr(), "subscribing is not advertised, so it must not succeed")
	unsubscribeResult := <-subscriber.Unsubscribe(context.Background(), nil)
	assert.Error(t, unsubscribeResult.GetErr(), "subscribing is not advertised, so unsubscribing must not succeed")
}

// TestConnection_UnadvertisedBrowsingRefuses documents that browsing is not merely unadvertised
// but has no builder at all.
func TestConnection_UnadvertisedBrowsingRefuses(t *testing.T) {
	c := testConnection()
	require.False(t, c.GetMetadata().CanBrowse())
	assert.Panics(t, func() { c.BrowseRequestBuilder() })
}
