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

package slmp

import (
	"encoding/binary"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// servePolls answers every Batch Read the stub codec sees with the word the payload function hands
// out, so a poll cycle sees a value that keeps moving. It stops when the returned function is called.
func servePolls(t *testing.T, codec *stubCodec, payload func() uint16) (stop func()) {
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
				frame, ok := request.message.(readWriteModel.SlmpRequestFrame3E)
				if !ok || frame.GetCommand() != commandBatchRead {
					continue
				}
				readRequest, ok := frame.GetRequestData().(readWriteModel.SlmpReadRequest)
				if !ok {
					continue
				}
				responseData := make([]byte, int(readRequest.GetNumberOfPoints())*2)
				binary.LittleEndian.PutUint16(responseData, payload())
				_ = request.handleMessage(readResponseFrame(0x0000, responseData))
			}
		}
	}()
	return func() {
		close(done)
		<-stopped
	}
}

// TestConnection_SubscribeCyclicPollsTheReadPath is the driver side of the polling emulation: SLMP
// has no subscription of its own, so a subscription has to turn into repeated Batch Reads over the
// wire. plc4j gets the same behaviour by extending PollingSubscriptionConnectionBase.
func TestConnection_SubscribeCyclicPollsTheReadPath(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	require.NoError(t, connection.Connect(testutils.TestContext(t)))

	var counter atomic.Uint32
	stop := servePolls(t, codec, func() uint16 {
		// A fresh value on every poll, so a change-of-state subscription would fire too.
		return uint16(counter.Add(1))
	})
	t.Cleanup(stop)

	subscriptionRequest, err := connection.SubscriptionRequestBuilder().
		AddCyclicTagAddress("hurz", "D350", 20*time.Millisecond).
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
		// the round trip of GetAddressString load bearing.
		assert.Equal(t, "D350:WORD[1]", event.GetAddress("hurz"))
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
	connection, _ := newConnectedTestConnection(t, DefaultConfiguration())

	subscriptionRequest, err := connection.SubscriptionRequestBuilder().
		AddEventTagAddress("hurz", "D350").
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

// TestConnection_UnsubscribeStopsThePoller closes the loop the other way round: taking a
// subscription back has to stop the frames going out for it.
func TestConnection_UnsubscribeStopsThePoller(t *testing.T) {
	connection, codec := newTestConnection(t, DefaultConfiguration())
	require.NoError(t, connection.Connect(testutils.TestContext(t)))
	t.Cleanup(func() {
		assert.NoError(t, connection.Close())
	})

	var counter atomic.Uint32
	stop := servePolls(t, codec, func() uint16 { return uint16(counter.Add(1)) })
	t.Cleanup(stop)

	subscriptionRequest, err := connection.SubscriptionRequestBuilder().
		AddCyclicTagAddress("hurz", "D350", 20*time.Millisecond).
		Build()
	require.NoError(t, err)
	result := <-subscriptionRequest.Execute(testutils.TestContext(t))
	require.NoError(t, result.GetErr())
	handle, err := result.GetResponse().GetSubscriptionHandle("hurz")
	require.NoError(t, err)

	// Wait for the poller to have run at least once, so the unsubscribe below really stops something.
	require.Eventually(t, func() bool { return counter.Load() > 0 }, 5*time.Second, 10*time.Millisecond)

	unsubscriptionRequest, err := connection.UnsubscriptionRequestBuilder().
		AddHandles(handle).
		Build()
	require.NoError(t, err)
	select {
	case unsubscriptionResult := <-unsubscriptionRequest.Execute(testutils.TestContext(t)):
		require.NoError(t, unsubscriptionResult.GetErr())
	case <-time.After(5 * time.Second):
		t.Fatal("the unsubscription didn't finish")
	}

	// A poll that was already in flight when the unsubscribe landed still gets answered, so let the
	// counter settle first and only then check that nothing moves it any more. The poll interval is
	// 20ms, so the second wait covers a good ten cycles a still-running poller would have used.
	time.Sleep(200 * time.Millisecond)
	settled := counter.Load()
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, settled, counter.Load(), "the poller has to stop when the subscription is taken back")
}
