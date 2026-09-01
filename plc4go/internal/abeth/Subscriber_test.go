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

package abeth

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/abeth/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// servePolls answers every read request the stub codec sees with the two payload bytes the counter
// hands out, so a poll cycle sees a value that keeps moving. It stops when the returned function is
// called.
func servePolls(t *testing.T, codec *stubCodec, payload func() []uint8) (stop func()) {
	t.Helper()
	done := make(chan struct{})
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		for {
			select {
			case <-done:
				return
			case request := <-codec.requests:
				readPacket, ok := request.message.(readWriteModel.CIPEncapsulationReadRequest)
				if !ok {
					continue
				}
				df1Request, ok := readPacket.GetRequest().(readWriteModel.DF1CommandRequestMessage)
				if !ok {
					continue
				}
				_ = request.handleMessage(readWriteModel.NewCIPEncapsulationReadResponse(
					testSessionHandle, 0, emptySenderContext, 0,
					readWriteModel.NewDF1CommandResponseMessageProtectedTypedLogicalRead(
						df1SourceAddress, 0, 0, df1Request.GetTransactionCounter(), payload())))
			}
		}
	}()
	return func() {
		close(done)
		<-stopped
	}
}

// TestConnection_SubscribeCyclicPollsTheReadPath is the driver side of the polling emulation: ab-eth
// has no subscription of its own, so a subscription has to turn into repeated reads over the wire.
// plc4j gets the same behaviour by extending PollingSubscriptionConnectionBase.
func TestConnection_SubscribeCyclicPollsTheReadPath(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()
	connectHandshake(t, codec)
	require.NoError(t, <-connectResult)

	var counter atomic.Uint32
	stop := servePolls(t, codec, func() []uint8 {
		// A fresh value on every poll, so a change-of-state subscription would fire too.
		return []uint8{uint8(counter.Add(1)), 0}
	})
	t.Cleanup(stop)

	subscriptionRequest, err := connection.SubscriptionRequestBuilder().
		AddCyclicTagAddress("hurz", "N7:3:WORD", 20*time.Millisecond).
		Build()
	require.NoError(t, err)

	var subscriptionResponse apiModel.PlcSubscriptionResponse
	select {
	case result := <-subscriptionRequest.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		subscriptionResponse = result.GetResponse()
		require.NotNil(t, subscriptionResponse)
	case <-time.After(5 * time.Second):
		t.Fatal("the subscription didn't finish")
	}
	assert.Equal(t, apiModel.PlcResponseCode_OK, subscriptionResponse.GetResponseCode("hurz"))

	handle, err := subscriptionResponse.GetSubscriptionHandle("hurz")
	require.NoError(t, err)

	events := make(chan apiModel.PlcSubscriptionEvent, 8)
	registration := handle.Register(func(event apiModel.PlcSubscriptionEvent) {
		select {
		case events <- event:
		default:
		}
	})
	t.Cleanup(registration.Unregister)

	select {
	case event := <-events:
		assert.Equal(t, apiModel.PlcResponseCode_OK, event.GetResponseCode("hurz"))
		// The address the poller polled has to be the address the tag spells, which is what makes
		// the round-trip of GetAddressString load bearing.
		assert.Equal(t, "N7:3:WORD", event.GetAddress("hurz"))
		require.NotNil(t, event.GetValue("hurz"))
		assert.Positive(t, event.GetValue("hurz").GetUint16(), "the poll has to carry the polled value")
	case <-time.After(5 * time.Second):
		t.Fatal("no subscription event arrived")
	}

	// Closing the connection has to take the poller down with it.
	require.NoError(t, connection.Close())
	assert.Positive(t, counter.Load(), "the subscription has to have polled at least once")
}

// TestConnection_SubscribeRejectsEventTags pins that the polling emulation says so rather than
// pretending: a spontaneous-event subscription cannot be emulated by polling.
func TestConnection_SubscribeRejectsEventTags(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	connectResult := make(chan error, 1)
	go func() {
		connectResult <- connection.Connect(testutils.TestContext(t))
	}()
	connectHandshake(t, codec)
	require.NoError(t, <-connectResult)
	t.Cleanup(func() {
		assert.NoError(t, connection.Close())
	})

	subscriptionRequest, err := connection.SubscriptionRequestBuilder().
		AddEventTagAddress("hurz", "N7:3:WORD").
		Build()
	require.NoError(t, err)

	select {
	case result := <-subscriptionRequest.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, result.GetResponse().GetResponseCode("hurz"))
	case <-time.After(5 * time.Second):
		t.Fatal("the subscription didn't finish")
	}
}
