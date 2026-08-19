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

package knxnetip

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// groupAddress123 is the wire representation of the 3-level group address "1/2/3"
// (5 bits main, 3 bits middle, 8 bits sub).
var groupAddress123 = []byte{1<<3 | 2, 3}

// switchPayload is the group-value payload of a "switch on" telegram for a DPT_Switch tag,
// exactly as it arrives from the bus: a datapoint-type of up to 6 bits is carried inside the
// single byte holding the embedded data bits, no data bytes follow it.
var switchPayload = []byte{0x01}

// scalingPayload is the group-value payload for a DPT_Scaling (8 bit) tag: the embedded-data
// byte stays empty and the value follows as a data byte.
var scalingPayload = []byte{0x00, 0x80}

// newSubscriberTestConnection creates a Connection which is just complete enough to
// drive the subscription handling (no transport, no codec involved).
func newSubscriberTestConnection(t *testing.T) *Connection {
	t.Helper()
	return &Connection{
		options:    map[string][]string{"group-address-num-levels": {"3"}},
		tagHandler: NewTagHandler(),
		valueCache: map[uint16][]byte{},
		metadata:   &ConnectionMetadata{},
		log:        testutils.ProduceTestingLogger(t),
	}
}

// subscribe builds and executes a subscription request for the given tag addresses and
// returns the resulting handles by tag name.
func subscribe(t *testing.T, connection *Connection, subscriber *Subscriber, subscriptionType apiModel.PlcSubscriptionType, tagAddresses map[string]string) map[string]apiModel.PlcSubscriptionHandle {
	t.Helper()
	builder := spiModel.NewDefaultPlcSubscriptionRequestBuilder(connection.tagHandler, connection.valueHandler, subscriber)
	for tagName, tagAddress := range tagAddresses {
		switch subscriptionType {
		case apiModel.SubscriptionChangeOfState:
			builder.AddChangeOfStateTagAddress(tagName, tagAddress)
		case apiModel.SubscriptionEvent:
			builder.AddEventTagAddress(tagName, tagAddress)
		default:
			t.Fatalf("unsupported subscription type %v", subscriptionType)
		}
	}
	subscriptionRequest, err := builder.Build()
	require.NoError(t, err)

	subscriptionResult := <-subscriptionRequest.Execute(t.Context())
	require.NoError(t, subscriptionResult.GetErr())
	response := subscriptionResult.GetResponse()
	require.NotNil(t, response)

	handles := map[string]apiModel.PlcSubscriptionHandle{}
	for tagName := range tagAddresses {
		require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode(tagName))
		handle, err := response.GetSubscriptionHandle(tagName)
		require.NoError(t, err)
		require.NotNil(t, handle)
		handles[tagName] = handle
	}
	return handles
}

func Test_Connection_UnsubscriptionRequestBuilder(t *testing.T) {
	connection := newSubscriberTestConnection(t)
	builder := connection.UnsubscriptionRequestBuilder()
	require.NotNil(t, builder, "callers would nil-deref on this")

	unsubscriptionRequest, err := builder.Build()
	require.NoError(t, err)
	assert.NotNil(t, unsubscriptionRequest)
}

func Test_Subscriber_Unsubscribe(t *testing.T) {
	connection := newSubscriberTestConnection(t)
	subscriber := NewSubscriber(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))

	handles := subscribe(t, connection, subscriber, apiModel.SubscriptionChangeOfState, map[string]string{
		"switch": "1/2/3:DPT_Switch",
	})
	assert.Len(t, connection.getSubscribers(), 1, "the subscriber gets fed by the connection")

	var events atomic.Int32
	handles["switch"].Register(func(_ apiModel.PlcSubscriptionEvent) {
		events.Add(1)
	})

	subscriber.handleValueChange(t.Context(), groupAddress123, switchPayload, true)
	require.Equal(t, int32(1), events.Load(), "a subscribed handle gets the event")

	unsubscriptionRequest, err := connection.UnsubscriptionRequestBuilder().AddHandles(handles["switch"]).Build()
	require.NoError(t, err)
	unsubscriptionResult := <-unsubscriptionRequest.Execute(t.Context())
	require.NoError(t, unsubscriptionResult.GetErr())
	require.NotNil(t, unsubscriptionResult.GetResponse())

	subscriber.handleValueChange(t.Context(), groupAddress123, switchPayload, true)
	assert.Equal(t, int32(1), events.Load(), "an unsubscribed handle doesn't get any more events")
	assert.Empty(t, subscriber.consumers, "the orphaned registration is dropped")
	assert.Empty(t, connection.getSubscribers(), "the connection stops feeding an empty subscriber")
}

func Test_Subscriber_UnsubscribeOnlyDropsTheGivenHandles(t *testing.T) {
	connection := newSubscriberTestConnection(t)
	subscriber := NewSubscriber(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))

	handles := subscribe(t, connection, subscriber, apiModel.SubscriptionChangeOfState, map[string]string{
		"first":  "1/2/3:DPT_Switch",
		"second": "1/2/3:DPT_Switch",
	})

	var firstEvents, secondEvents atomic.Int32
	handles["first"].Register(func(_ apiModel.PlcSubscriptionEvent) { firstEvents.Add(1) })
	handles["second"].Register(func(_ apiModel.PlcSubscriptionEvent) { secondEvents.Add(1) })

	unsubscriptionRequest, err := connection.UnsubscriptionRequestBuilder().AddHandles(handles["first"]).Build()
	require.NoError(t, err)
	require.NoError(t, (<-unsubscriptionRequest.Execute(t.Context())).GetErr())

	subscriber.handleValueChange(t.Context(), groupAddress123, switchPayload, true)
	assert.Equal(t, int32(0), firstEvents.Load())
	assert.Equal(t, int32(1), secondEvents.Load())
	assert.Len(t, connection.getSubscribers(), 1, "there is still something to deliver")
}

func Test_Subscriber_UnsubscribeIgnoresForeignHandles(t *testing.T) {
	connection := newSubscriberTestConnection(t)
	subscriber := NewSubscriber(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
	otherSubscriber := NewSubscriber(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))

	handles := subscribe(t, connection, subscriber, apiModel.SubscriptionChangeOfState, map[string]string{
		"switch": "1/2/3:DPT_Switch",
	})
	otherHandles := subscribe(t, connection, otherSubscriber, apiModel.SubscriptionChangeOfState, map[string]string{
		"switch": "1/2/3:DPT_Switch",
	})

	// A request spanning both subscribers gets handed to each of them in full.
	unsubscriptionRequest, err := connection.UnsubscriptionRequestBuilder().
		AddHandles(handles["switch"], otherHandles["switch"]).
		Build()
	require.NoError(t, err)
	require.NoError(t, (<-unsubscriptionRequest.Execute(t.Context())).GetErr())

	assert.Empty(t, connection.getSubscribers())
}

func Test_Subscriber_UnsubscribeUnsupportedHandle(t *testing.T) {
	connection := newSubscriberTestConnection(t)
	subscriber := NewSubscriber(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))

	err := subscriber.unsubscribeHandles([]apiModel.PlcSubscriptionHandle{
		spiModel.NewDefaultPlcSubscriptionHandle(subscriber),
	})
	assert.ErrorContains(t, err, "is not a knx subscription handle")
}

func Test_Subscriber_handleValueChange_deliversToAllAdvertisedTypes(t *testing.T) {
	tests := []struct {
		name             string
		subscriptionType apiModel.PlcSubscriptionType
		changed          bool
		wantEvents       int32
	}{
		{"change-of-state with a changed value", apiModel.SubscriptionChangeOfState, true, 1},
		{"change-of-state with an unchanged value", apiModel.SubscriptionChangeOfState, false, 0},
		{"event with a changed value", apiModel.SubscriptionEvent, true, 1},
		// KNX is event driven, an event subscription wants every write on the bus,
		// not only the ones which changed the value.
		{"event with an unchanged value", apiModel.SubscriptionEvent, false, 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			connection := newSubscriberTestConnection(t)
			subscriber := NewSubscriber(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
			handles := subscribe(t, connection, subscriber, test.subscriptionType, map[string]string{
				"switch": "1/2/3:DPT_Switch",
			})

			var events atomic.Int32
			handles["switch"].Register(func(event apiModel.PlcSubscriptionEvent) {
				assert.Equal(t, apiModel.PlcResponseCode_OK, event.GetResponseCode("switch"))
				events.Add(1)
			})

			subscriber.handleValueChange(t.Context(), groupAddress123, switchPayload, test.changed)
			assert.Equal(t, test.wantEvents, events.Load())
		})
	}
}

// Test_Subscriber_handleValueChange_decodesThePayload is the regression test for the
// "skip the first byte" guard which used to be driven by KnxDatapointType.GetLengthInBits
// (hardcoded to 32 for every datapoint-type by the generated model, so it always fired).
// Skipping a byte of a single-byte DPT_Switch payload emptied the buffer and the parse
// failed with INVALID_DATA, and a DPT_Scaling value was consumed as the reserved byte.
func Test_Subscriber_handleValueChange_decodesThePayload(t *testing.T) {
	tests := []struct {
		name       string
		tagAddress string
		payload    []byte
		assertion  func(t *testing.T, value apiValues.PlcValue)
	}{
		{
			name:       "a datapoint-type which fits into the embedded data bits",
			tagAddress: "1/2/3:DPT_Switch",
			payload:    switchPayload,
			assertion: func(t *testing.T, value apiValues.PlcValue) {
				assert.True(t, value.GetBool())
			},
		},
		{
			name:       "a datapoint-type which needs its own data byte",
			tagAddress: "1/2/3:DPT_Scaling",
			payload:    scalingPayload,
			assertion: func(t *testing.T, value apiValues.PlcValue) {
				assert.Equal(t, uint8(0x80), value.GetUint8())
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			connection := newSubscriberTestConnection(t)
			subscriber := NewSubscriber(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
			handles := subscribe(t, connection, subscriber, apiModel.SubscriptionEvent, map[string]string{
				"tag": test.tagAddress,
			})

			var events atomic.Int32
			handles["tag"].Register(func(event apiModel.PlcSubscriptionEvent) {
				events.Add(1)
				require.Equal(t, apiModel.PlcResponseCode_OK, event.GetResponseCode("tag"))
				test.assertion(t, event.GetValue("tag"))
			})

			subscriber.handleValueChange(t.Context(), groupAddress123, test.payload, true)
			assert.Equal(t, int32(1), events.Load())
		})
	}
}

func Test_Subscriber_handleValueChange_ignoresNonMatchingAddresses(t *testing.T) {
	connection := newSubscriberTestConnection(t)
	subscriber := NewSubscriber(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
	handles := subscribe(t, connection, subscriber, apiModel.SubscriptionEvent, map[string]string{
		"switch": "1/2/3:DPT_Switch",
	})

	var events atomic.Int32
	handles["switch"].Register(func(_ apiModel.PlcSubscriptionEvent) { events.Add(1) })

	subscriber.handleValueChange(t.Context(), []byte{1<<3 | 2, 4}, switchPayload, true)
	assert.Equal(t, int32(0), events.Load())
}

// Test_Subscriber_consumerMapRace hammers the consumer map from all the goroutines which
// touch it in production: the codec worker delivering values, the api registering and
// unregistering consumers and the (un)subscription paths.
func Test_Subscriber_consumerMapRace(t *testing.T) {
	connection := newSubscriberTestConnection(t)
	subscriber := NewSubscriber(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
	handles := subscribe(t, connection, subscriber, apiModel.SubscriptionEvent, map[string]string{
		"switch": "1/2/3:DPT_Switch",
	})

	const iterations = 200
	var wg sync.WaitGroup
	wg.Go(func() {
		for range iterations {
			// This is what the codec worker does for every incoming group value write.
			connection.handleValueCacheUpdate(context.Background(), groupAddress123, switchPayload)
		}
	})
	wg.Go(func() {
		for range iterations {
			registration := handles["switch"].Register(func(_ apiModel.PlcSubscriptionEvent) {})
			registration.Unregister()
		}
	})
	wg.Go(func() {
		for range iterations {
			otherSubscriber := NewSubscriber(connection, options.WithCustomLogger(testutils.ProduceTestingLogger(t)))
			connection.addSubscriber(otherSubscriber)
			connection.removeSubscriber(otherSubscriber)
		}
	})
	wg.Go(func() {
		for range iterations {
			_ = subscriber.unsubscribeHandles(nil)
		}
	})
	wg.Wait()
}
