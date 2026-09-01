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

package umas

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// UMAS has no subscription mechanism of its own, so a subscription is emulated by polling the read
// path - the same thing plc4j's PollingSubscriptionConnectionBase does for this driver. The point of
// this test is that a subscription really does turn into ordinary variable reads and that the events
// reach a registered consumer.
func TestSubscriber_PollsTheReadPath(t *testing.T) {
	connection, codec := newConnectedConnection(t)

	request, err := connection.SubscriptionRequestBuilder().
		AddCyclicTagAddress("cyclic", "g_b16", 50*time.Millisecond).
		Build()
	require.NoError(t, err)

	var subscriptionResponse apiModel.PlcSubscriptionResponse
	select {
	case result := <-request.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		subscriptionResponse = result.GetResponse()
	case <-time.After(5 * time.Second):
		t.Fatal("the subscription never delivered a result")
	}
	require.Equal(t, apiModel.PlcResponseCode_OK, subscriptionResponse.GetResponseCode("cyclic"))

	handle, err := subscriptionResponse.GetSubscriptionHandle("cyclic")
	require.NoError(t, err)

	events := make(chan apiModel.PlcSubscriptionEvent, 4)
	registration := handle.Register(
		func(event apiModel.PlcSubscriptionEvent) {
			select {
			case events <- event:
			default:
			}
		})
	t.Cleanup(registration.Unregister)

	// The poller reads the symbol like any other read, so the PLC side of this looks exactly like the
	// reader test.
	pollRequest := codec.nextRequest(t)
	readRequest, ok := pollRequest.item(t).(readWriteModel.UmasPDUReadVariableRequest)
	require.True(t, ok, "a poll cycle has to be an ordinary variable read, got %T", pollRequest.item(t))
	assert.Equal(t, uint16(2), readRequest.GetVariables()[0].GetBlock())
	codec.answerWith(t, pollRequest, readWriteModel.NewUmasPDUReadVariableResponse(0, []byte{0x2A, 0x00}))

	select {
	case event := <-events:
		assert.Equal(t, apiModel.PlcResponseCode_OK, event.GetResponseCode("cyclic"))
		assert.Equal(t, apiValues.INT, event.GetValue("cyclic").GetPlcValueType())
		assert.Equal(t, int16(42), event.GetValue("cyclic").GetInt16())
	case <-time.After(5 * time.Second):
		t.Fatal("the subscription never emitted an event")
	}

	// Unsubscribing has to stop the poller.
	unsubscribeRequest, err := connection.UnsubscriptionRequestBuilder().
		AddHandles(handle).
		Build()
	require.NoError(t, err)
	select {
	case result := <-unsubscribeRequest.Execute(testutils.TestContext(t)):
		assert.NoError(t, result.GetErr())
	case <-time.After(5 * time.Second):
		t.Fatal("the unsubscription never delivered a result")
	}
}
