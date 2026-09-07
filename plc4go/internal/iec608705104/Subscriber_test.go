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

package iec608705104

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// subscriptionKind is how a test asks for one of the three subscription types.
type subscriptionKind int

const (
	changeOfState subscriptionKind = iota
	eventDriven
	cyclic
)

// subscribe runs a subscription for a single tag and hands back its handle and response code. Every
// path through Subscribe has to complete the result channel, which is what the timeout catches.
func subscribe(t *testing.T, connection *Connection, kind subscriptionKind, tagName string, tagAddress string) (apiModel.PlcSubscriptionHandle, apiModel.PlcResponseCode) {
	t.Helper()
	builder := connection.SubscriptionRequestBuilder()
	switch kind {
	case eventDriven:
		builder.AddEventTagAddress(tagName, tagAddress)
	case cyclic:
		builder.AddCyclicTagAddress(tagName, tagAddress, time.Second)
	default:
		builder.AddChangeOfStateTagAddress(tagName, tagAddress)
	}
	request, err := builder.Build()
	require.NoError(t, err)
	select {
	case result := <-request.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		handle, _ := result.GetResponse().GetSubscriptionHandle(tagName)
		return handle, result.GetResponse().GetResponseCode(tagName)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the subscription never completed the result channel")
		return nil, apiModel.PlcResponseCode_INTERNAL_ERROR
	}
}

// collectEvents registers a consumer which collects everything it is handed.
func collectEvents(t *testing.T, handle apiModel.PlcSubscriptionHandle) func(want int) []apiModel.PlcSubscriptionEvent {
	t.Helper()
	events := make(chan apiModel.PlcSubscriptionEvent, 64)
	registration := handle.Register(func(event apiModel.PlcSubscriptionEvent) {
		events <- event
	})
	require.NotNil(t, registration)
	// collected waits for at least want events and then hands back everything which arrived, so a
	// test which expects nothing can pass zero and a test which expects two doesn't race the worker.
	return func(want int) []apiModel.PlcSubscriptionEvent {
		var collected []apiModel.PlcSubscriptionEvent
		deadline := time.After(5 * time.Second)
		for len(collected) < want {
			select {
			case event := <-events:
				collected = append(collected, event)
			case <-deadline:
				require.FailNowf(t, "not enough events", "wanted %d, got %d", want, len(collected))
			}
		}
		// Give anything extra a moment to show up, so an over-delivery is caught rather than missed.
		settle := time.After(100 * time.Millisecond)
		for {
			select {
			case event := <-events:
				collected = append(collected, event)
			case <-settle:
				return collected
			}
		}
	}
}

// pushAsdu hands the connection one I-format frame carrying one ASDU.
func pushAsdu(t *testing.T, codec *stubCodec, sendSequenceNo uint16, asdu []byte) {
	t.Helper()
	pushIncoming(t, codec, parseApdu(t, iFormatFrame(sendSequenceNo, 0, asdu)))
}

// singlePointAsdu is a spontaneous single point information at ASDU 10.
func singlePointAsdu(informationObjectAddress uint32, siq byte) []byte {
	return asduBytes(0x01, 3, 10, informationObjectBytes(informationObjectAddress, siq))
}

// An exact tag hears about its own point and nothing else. plc4j hands every event to every
// consumer regardless of what it subscribed to, because its tags all resolve to ASDU 0 / IOA 0.
func TestSubscriber_ExactTagFiltersByPoint(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, changeOfState, "breaker", "10/13")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	require.NotNil(t, handle)
	collected := collectEvents(t, handle)

	pushAsdu(t, codec, 0, singlePointAsdu(12, 0x01)) // a different point
	pushAsdu(t, codec, 1, singlePointAsdu(13, 0x01)) // ours
	pushAsdu(t, codec, 2, asduBytes(0x01, 3, 11, informationObjectBytes(13, 0x01)))

	events := collected(1)
	require.Len(t, events, 1, "only the subscribed point is delivered")
	assert.Equal(t, "10/13", events[0].GetAddress("breaker"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, events[0].GetResponseCode("breaker"))
	assert.True(t, nested(t, events[0].GetValue("breaker"), fieldValue).GetBool())
}

// A wildcard tag hears about every point it covers, and the event says which one fired - the tag
// address alone could not.
func TestSubscriber_WildcardTagReportsTheConcretePoint(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, eventDriven, "everything", "*/*")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collected := collectEvents(t, handle)

	pushAsdu(t, codec, 0, singlePointAsdu(13, 0x01))
	pushAsdu(t, codec, 1, asduBytes(0x01, 3, 11, informationObjectBytes(99, 0x00)))

	events := collected(2)
	require.Len(t, events, 2)
	assert.Equal(t, "10/13", events[0].GetAddress("everything"))
	assert.Equal(t, "11/99", events[1].GetAddress("everything"))
}

// A tag wildcarded in one part still pins the other.
func TestSubscriber_WildcardTagStillFilters(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, eventDriven, "station", "10/*")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collected := collectEvents(t, handle)

	pushAsdu(t, codec, 0, asduBytes(0x01, 3, 11, informationObjectBytes(13, 0x01)))
	pushAsdu(t, codec, 1, singlePointAsdu(13, 0x01))

	events := collected(1)
	require.Len(t, events, 1)
	assert.Equal(t, "10/13", events[0].GetAddress("station"))
}

// Every information object of an ASDU is a point of its own, so one frame can produce several events.
func TestSubscriber_EveryInformationObjectIsItsOwnEvent(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, eventDriven, "all", "10/*")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collected := collectEvents(t, handle)

	pushAsdu(t, codec, 0, asduBytes(0x01, 3, 10,
		informationObjectBytes(1, 0x01),
		informationObjectBytes(2, 0x00),
		informationObjectBytes(3, 0x01)))

	events := collected(3)
	require.Len(t, events, 3)
	assert.Equal(t, "10/1", events[0].GetAddress("all"))
	assert.Equal(t, "10/2", events[1].GetAddress("all"))
	assert.Equal(t, "10/3", events[2].GetAddress("all"))
}

// A change-of-state subscription reports a point when its state changes and not when the station
// simply re-sends it - which it does on every general interrogation. plc4j forwards everything.
func TestSubscriber_ChangeOfStateReportsChangesOnly(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, changeOfState, "breaker", "10/13")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collected := collectEvents(t, handle)

	pushAsdu(t, codec, 0, singlePointAsdu(13, 0x01)) // closed
	pushAsdu(t, codec, 1, singlePointAsdu(13, 0x01)) // still closed - not a change
	pushAsdu(t, codec, 2, singlePointAsdu(13, 0x00)) // opened
	pushAsdu(t, codec, 3, singlePointAsdu(13, 0x80)) // still open, but now flagged invalid

	events := collected(3)
	require.Len(t, events, 3, "the repeat is not a change, the quality change is")
	assert.True(t, nested(t, events[0].GetValue("breaker"), fieldValue).GetBool())
	assert.False(t, nested(t, events[1].GetValue("breaker"), fieldValue).GetBool())
	assert.True(t, nested(t, events[2].GetValue("breaker"), fieldQuality, "invalid").GetBool())
}

// Change detection is per point, so a wildcard subscription doesn't let one point's report suppress
// another's.
func TestSubscriber_ChangeDetectionIsPerPoint(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, changeOfState, "all", "*/*")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collected := collectEvents(t, handle)

	pushAsdu(t, codec, 0, singlePointAsdu(13, 0x01))
	pushAsdu(t, codec, 1, singlePointAsdu(14, 0x01))
	pushAsdu(t, codec, 2, singlePointAsdu(13, 0x01)) // a repeat of the first point

	events := collected(2)
	require.Len(t, events, 2)
	assert.Equal(t, "10/13", events[0].GetAddress("all"))
	assert.Equal(t, "10/14", events[1].GetAddress("all"))
}

// An event subscription gets everything the station sends, repeats included: it is the station which
// decided the report was worth sending.
func TestSubscriber_EventSubscriptionGetsRepeats(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, eventDriven, "breaker", "10/13")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collected := collectEvents(t, handle)

	pushAsdu(t, codec, 0, singlePointAsdu(13, 0x01))
	pushAsdu(t, codec, 1, singlePointAsdu(13, 0x01))

	assert.Len(t, collected(2), 2)
}

// How often a point is reported is the station's own business, so a request for a particular interval
// cannot be honoured. Saying so is better than accepting it and ignoring the interval, which is what
// plc4j does.
func TestSubscriber_CyclicIsUnsupported(t *testing.T) {
	connection, _ := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, cyclic, "breaker", "10/13")

	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, responseCode)
	assert.Nil(t, handle)
}

// An address which doesn't parse never becomes a request at all: the builder runs the tag handler,
// so the caller is told at Build time rather than by a subscription which never delivers.
func TestSubscriber_RefusesAnUnparseableAddress(t *testing.T) {
	connection, _ := newHandshakenConnection(t)

	builder := connection.SubscriptionRequestBuilder()
	builder.AddChangeOfStateTagAddress("nonsense", "not-an-address")
	request, err := builder.Build()

	assert.Error(t, err)
	assert.Nil(t, request)
}

// foreignTag is a subscription tag from some other driver. Its address string cannot be re-parsed as
// an IEC 60870-5-104 address, which is the one way a tag can reach the subscriber unusable.
type foreignTag struct{}

func (foreignTag) GetAddressString() string             { return "digital:1" }
func (foreignTag) GetValueType() apiValues.PlcValueType { return apiValues.BOOL }
func (foreignTag) GetArrayInfo() []apiModel.ArrayInfo   { return nil }
func (foreignTag) GetPlcSubscriptionType() apiModel.PlcSubscriptionType {
	return apiModel.SubscriptionEvent
}
func (foreignTag) GetDuration() time.Duration { return 0 }
func (foreignTag) String() string             { return "foreignTag" }

// A tag which isn't one of ours gets a response code rather than a handle which would never match a
// single point.
func TestSubscriber_RefusesAForeignTag(t *testing.T) {
	connection, _ := newHandshakenConnection(t)

	builder := connection.SubscriptionRequestBuilder()
	builder.AddEventTag("alien", foreignTag{})
	request, err := builder.Build()
	require.NoError(t, err)

	select {
	case result := <-request.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, result.GetResponse().GetResponseCode("alien"))
		handle, _ := result.GetResponse().GetSubscriptionHandle("alien")
		assert.Nil(t, handle)
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the subscription never completed the result channel")
	}
}

// Unsubscribing stops the delivery. Nothing goes out on the wire - STOPDT would silence every other
// subscription of the same connection.
func TestSubscriber_Unsubscribe(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, eventDriven, "breaker", "10/13")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collected := collectEvents(t, handle)
	pushAsdu(t, codec, 0, singlePointAsdu(13, 0x01))
	require.Len(t, collected(1), 1)

	sentBefore := len(codec.sentCommands())
	unsubscriptionRequest, err := connection.UnsubscriptionRequestBuilder().AddHandles(handle).Build()
	require.NoError(t, err)
	select {
	case result := <-unsubscriptionRequest.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the unsubscription never completed the result channel")
	}

	pushAsdu(t, codec, 1, singlePointAsdu(13, 0x00))
	assert.Empty(t, collected(0), "nothing is delivered after unsubscribing")
	assert.Len(t, codec.sentCommands(), sentBefore, "unsubscribing sends nothing")
}

// An ASDU whose payload the model carries nothing decodable for is still delivered, with a response
// code which says so - dropping it would leave the caller believing the point never reported.
func TestSubscriber_DeliversUndecodableAsdusWithACode(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, eventDriven, "file", "10/13")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collected := collectEvents(t, handle)

	pushAsdu(t, codec, 0, asduBytes(0x78, 3, 10, informationObjectBytes(13)))

	events := collected(1)
	require.Len(t, events, 1)
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, events[0].GetResponseCode("file"))
	assert.Equal(t, apiValues.NULL, nested(t, events[0].GetValue("file"), fieldValue).GetPlcValueType())
}

// A sequence-of-objects ASDU (structure qualifier set, several objects) carries one address for all
// of them, which the generated model cannot represent - it reads an address in front of every
// object. Publishing it would file one point's reading under another point's address, so the whole
// ASDU is refused rather than delivered wrong.
func TestSubscriber_RefusesASequenceOfObjectsAsdu(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, eventDriven, "all", "*/*")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collected := collectEvents(t, handle)

	// The structure qualifier is the top bit of the octet the number of objects sits in.
	sequenceAsdu := asduBytes(0x01, 3, 10,
		informationObjectBytes(1, 0x01),
		informationObjectBytes(2, 0x01))
	sequenceAsdu[1] |= 0x80
	pushAsdu(t, codec, 0, sequenceAsdu)

	assert.Empty(t, collected(0), "a mis-addressable ASDU is dropped rather than delivered wrong")

	// A single-object ASDU is the one case where both layouts are the same bytes, so it still gets
	// through even with the structure qualifier set.
	singleAsdu := asduBytes(0x01, 3, 10, informationObjectBytes(1, 0x01))
	singleAsdu[1] |= 0x80
	pushAsdu(t, codec, 1, singleAsdu)

	events := collected(1)
	require.Len(t, events, 1)
	assert.Equal(t, "10/1", events[0].GetAddress("all"))
}

// A closed connection delivers nothing further, and closing while a subscription is live must not
// hang.
func TestSubscriber_ClosedConnectionStopsDelivering(t *testing.T) {
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler(_options...), _options...)
	completeHandshake(t, connection, codec)

	handle, responseCode := subscribe(t, connection, eventDriven, "breaker", "10/13")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collected := collectEvents(t, handle)

	require.NoError(t, connection.Close())

	pushAsdu(t, codec, 0, singlePointAsdu(13, 0x01))
	assert.Empty(t, collected(0))
}

// Two consumers on the same handle both hear about a change - the change detection is per handle,
// not per consumer, so the first delivery must not swallow the report for the second.
func TestSubscriber_SeveralConsumersOnOneHandle(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	handle, responseCode := subscribe(t, connection, changeOfState, "breaker", "10/13")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	first := collectEvents(t, handle)
	second := collectEvents(t, handle)

	pushAsdu(t, codec, 0, singlePointAsdu(13, 0x01))

	assert.Len(t, first(1), 1)
	assert.Len(t, second(1), 1)
}

// The subscription tag has to survive a round trip through its address string, because a tag which
// arrives wrapped in a DefaultPlcSubscriptionTag is re-parsed from it.
func TestSubscriber_SubscribeByTag(t *testing.T) {
	connection, codec := newHandshakenConnection(t)

	builder := connection.SubscriptionRequestBuilder()
	builder.AddEventTag("breaker", parseTag(t, "10/13"))
	request, err := builder.Build()
	require.NoError(t, err)

	var handle apiModel.PlcSubscriptionHandle
	select {
	case result := <-request.Execute(testutils.TestContext(t)):
		require.NoError(t, result.GetErr())
		require.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("breaker"))
		handle, _ = result.GetResponse().GetSubscriptionHandle("breaker")
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the subscription never completed the result channel")
	}
	require.NotNil(t, handle)
	collected := collectEvents(t, handle)

	pushAsdu(t, codec, 0, singlePointAsdu(13, 0x01))
	assert.Len(t, collected(1), 1)
}

// A handle whose change-detection table has grown past its bound forgets rather than growing without
// limit; the cost is one duplicate report per point.
func TestSubscriptionHandle_ChangeDetectionIsBounded(t *testing.T) {
	handle := NewSubscriptionHandle(nil, "all", parseTag(t, "*/*"), apiModel.SubscriptionChangeOfState, 0)

	for point := range uint32(maxTrackedPoints) {
		require.True(t, handle.shouldPublish(1, point, "same"))
	}
	assert.Len(t, handle.lastFingerprints, maxTrackedPoints)

	// The next point pushes it over, which clears the table.
	assert.True(t, handle.shouldPublish(1, maxTrackedPoints, "same"))
	assert.Len(t, handle.lastFingerprints, 1)
}
