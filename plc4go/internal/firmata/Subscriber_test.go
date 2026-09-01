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

package firmata

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// subscribe runs a change-of-state subscription for a single tag and hands back its handle.
func subscribe(t *testing.T, connection *Connection, tagName string, tagAddress string) (apiModel.PlcSubscriptionHandle, apiModel.PlcResponseCode) {
	t.Helper()
	builder := connection.SubscriptionRequestBuilder()
	builder.AddChangeOfStateTagAddress(tagName, tagAddress)
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
func collectEvents(t *testing.T, handle apiModel.PlcSubscriptionHandle) (func() []apiModel.PlcSubscriptionEvent, apiModel.PlcConsumerRegistration) {
	t.Helper()
	events := make(chan apiModel.PlcSubscriptionEvent, 16)
	registration := handle.Register(func(event apiModel.PlcSubscriptionEvent) {
		events <- event
	})
	require.NotNil(t, registration)
	return func() []apiModel.PlcSubscriptionEvent {
		var collected []apiModel.PlcSubscriptionEvent
		for {
			select {
			case event := <-events:
				collected = append(collected, event)
			default:
				return collected
			}
		}
	}, registration
}

// Subscribing to a digital pin has to put the pin into input mode and switch reporting on for the
// port it sits in - reporting is per port of 8 pins, not per pin.
func TestSubscriber_SubscribeDigitalPin(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	handle, responseCode := subscribe(t, connection, "button", "digital:9")
	assert.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	require.NotNil(t, handle)

	assert.Equal(t, []byte{
		0xF4, 0x09, 0x00, // set pin mode: pin 9, input
		0xD1, 0x01, // report digital: port 1, on
	}, sentBytes(t, transportInstance))
}

func TestSubscriber_SubscribeDigitalPinWithPullup(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	_, responseCode := subscribe(t, connection, "button", "digital:9:PULLUP")
	assert.Equal(t, apiModel.PlcResponseCode_OK, responseCode)

	assert.Equal(t, []byte{
		0xF4, 0x09, 0x0B, // set pin mode: pin 9, pullup
		0xD1, 0x01,
	}, sentBytes(t, transportInstance))
}

// A run of pins spanning two ports needs reporting switched on for both of them, once each.
func TestSubscriber_SubscribeARunOfDigitalPinsAcrossPorts(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	_, responseCode := subscribe(t, connection, "bar", "digital:6[0..3]")
	assert.Equal(t, apiModel.PlcResponseCode_OK, responseCode)

	assert.Equal(t, []byte{
		0xF4, 0x06, 0x00, 0xF4, 0x07, 0x00, 0xF4, 0x08, 0x00, 0xF4, 0x09, 0x00,
		0xD0, 0x01, // report digital: port 0
		0xD1, 0x01, // report digital: port 1
	}, sentBytes(t, transportInstance))
}

// Analog pins need no mode: asking the board to report one is enough (plc4j
// FirmataConnection.buildSubscribeMessages).
func TestSubscriber_SubscribeAnalogPin(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	_, responseCode := subscribe(t, connection, "dial", "analog:2[0..1]")
	assert.Equal(t, apiModel.PlcResponseCode_OK, responseCode)

	assert.Equal(t, []byte{
		0xC2, 0x01, // report analog: pin 2, on
		0xC3, 0x01, // report analog: pin 3, on
	}, sentBytes(t, transportInstance))
}

// Subscribing to a pin which is already being reported on costs nothing on the wire.
func TestSubscriber_SubscribeTwiceSendsNothingTwice(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	_, responseCode := subscribe(t, connection, "button", "digital:9")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	sentBytes(t, transportInstance)

	handle, responseCode := subscribe(t, connection, "button-again", "digital:9")
	assert.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	assert.NotNil(t, handle, "a second subscription still gets a handle of its own")
	assert.Empty(t, sentBytes(t, transportInstance))
}

// A pin which is driven as an output can't be reported on.
func TestSubscriber_SubscribeRefusesAnOutputPin(t *testing.T) {
	connection, transportInstance := newTestConnection(t)
	require.NoError(t, executeWrite(t, connection, "led", "digital:9", true).GetErr())
	sentBytes(t, transportInstance)

	handle, responseCode := subscribe(t, connection, "button", "digital:9")
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, responseCode)
	assert.Nil(t, handle)
	assert.Empty(t, sentBytes(t, transportInstance))
}

// A firmata board reports a pin when it changes, not on a schedule of our choosing, and the global
// sampling interval it does have is a property of the board rather than of one tag. plc4j accepts a
// cyclic subscription and then quietly treats it as a change-of-state one.
func TestSubscriber_SubscribeRefusesCyclicTags(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	builder := connection.SubscriptionRequestBuilder()
	builder.AddCyclicTagAddress("button", "digital:9", time.Second)
	builder.AddChangeOfStateTagAddress("other", "digital:10")
	request, err := builder.Build()
	require.NoError(t, err)
	result := <-request.Execute(testutils.TestContext(t))

	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, result.GetResponse().GetResponseCode("button"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("other"))
	cyclicHandle, _ := result.GetResponse().GetSubscriptionHandle("button")
	assert.Nil(t, cyclicHandle)
	// Only the tag which was accepted was configured on the board.
	assert.Equal(t, []byte{0xF4, 0x0A, 0x00, 0xD1, 0x01}, sentBytes(t, transportInstance))
}

func TestSubscriber_DeliversDigitalChanges(t *testing.T) {
	connection, _ := newTestConnection(t)
	handle, responseCode := subscribe(t, connection, "button", "digital:9")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collect, _ := collectEvents(t, handle)

	// Port 1 covers pins 8 to 15; pin 9 is bit 1 of the port.
	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x02, 0x00}))
	events := collect()
	require.Len(t, events, 1)
	assert.Equal(t, apiModel.PlcResponseCode_OK, events[0].GetResponseCode("button"))
	assert.True(t, events[0].GetValue("button").GetBool())
	assert.Equal(t, "digital:9", events[0].GetAddress("button"))

	// The same report again isn't a change and isn't delivered.
	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x02, 0x00}))
	assert.Empty(t, collect())

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x00, 0x00}))
	events = collect()
	require.Len(t, events, 1)
	assert.False(t, events[0].GetValue("button").GetBool())
}

// A pin outside the subscribed run must not wake the consumer up.
func TestSubscriber_IgnoresPinsItDoesNotCover(t *testing.T) {
	connection, _ := newTestConnection(t)
	handle, _ := subscribe(t, connection, "button", "digital:9")
	collect, _ := collectEvents(t, handle)

	// Pin 8 changes, pin 9 doesn't.
	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x01, 0x00}))
	assert.Empty(t, collect())
}

// A tag covering several pins is answered with the value of every pin it covers, not just the one
// which changed (plc4j FirmataConnection.publishDigitalEvents).
func TestSubscriber_DeliversTheWholeRun(t *testing.T) {
	connection, _ := newTestConnection(t)
	handle, _ := subscribe(t, connection, "bar", "digital:8[0..2]")
	collect, _ := collectEvents(t, handle)

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x05, 0x00}))
	events := collect()
	require.Len(t, events, 1)
	value := events[0].GetValue("bar")
	require.True(t, value.IsList())
	require.Len(t, value.GetList(), 3)
	assert.True(t, value.GetList()[0].GetBool())
	assert.False(t, value.GetList()[1].GetBool())
	assert.True(t, value.GetList()[2].GetBool())
}

func TestSubscriber_DeliversAnalogChanges(t *testing.T) {
	connection, _ := newTestConnection(t)
	handle, responseCode := subscribe(t, connection, "dial", "analog:3")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	collect, _ := collectEvents(t, handle)

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageAnalogIO(3, []int8{0x05, 0x01}))
	events := collect()
	require.Len(t, events, 1)
	assert.Equal(t, int16(0x85), events[0].GetValue("dial").GetInt16())
	assert.Equal(t, "analog:3", events[0].GetAddress("dial"))

	// The very same sample again is no change.
	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageAnalogIO(3, []int8{0x05, 0x01}))
	assert.Empty(t, collect())
}

// A pin of a run the board hasn't sampled yet is reported as -1, the way plc4j fills the gaps.
func TestSubscriber_DeliversUnknownAnalogPinsAsMinusOne(t *testing.T) {
	connection, _ := newTestConnection(t)
	handle, _ := subscribe(t, connection, "dials", "analog:3[0..1]")
	collect, _ := collectEvents(t, handle)

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageAnalogIO(3, []int8{0x05, 0x00}))
	events := collect()
	require.Len(t, events, 1)
	value := events[0].GetValue("dials")
	require.True(t, value.IsList())
	assert.Equal(t, int16(5), value.GetList()[0].GetInt16())
	assert.Equal(t, int16(-1), value.GetList()[1].GetInt16())
}

// Unsubscribing drops the handle, which stops the delivery. Reporting stays switched on at the
// board: it is per port for digital pins and another subscription may well still want it.
func TestSubscriber_UnsubscribeStopsDelivery(t *testing.T) {
	connection, _ := newTestConnection(t)
	handle, _ := subscribe(t, connection, "button", "digital:9")
	collect, _ := collectEvents(t, handle)

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x02, 0x00}))
	require.Len(t, collect(), 1)

	builder := connection.UnsubscriptionRequestBuilder()
	builder.AddHandles(handle)
	request, err := builder.Build()
	require.NoError(t, err)
	result := <-request.Execute(testutils.TestContext(t))
	require.NoError(t, result.GetErr())
	require.NotNil(t, result.GetResponse())

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x00, 0x00}))
	assert.Empty(t, collect())
	assert.Empty(t, connection.activeSubscribers(), "a subscriber without handles stops being fed")
}

func TestSubscriber_UnregisterStopsDelivery(t *testing.T) {
	connection, _ := newTestConnection(t)
	handle, _ := subscribe(t, connection, "button", "digital:9")
	collect, registration := collectEvents(t, handle)

	registration.Unregister()

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x02, 0x00}))
	assert.Empty(t, collect())
}

// A consumer which was handed to the builder is registered by the time the response comes back.
func TestSubscriber_PreRegisteredConsumer(t *testing.T) {
	connection, _ := newTestConnection(t)

	events := make(chan apiModel.PlcSubscriptionEvent, 4)
	builder := connection.SubscriptionRequestBuilder()
	builder.AddChangeOfStateTagAddress("button", "digital:9")
	builder.AddPreRegisteredConsumer("button", func(event apiModel.PlcSubscriptionEvent) {
		events <- event
	})
	request, err := builder.Build()
	require.NoError(t, err)
	result := <-request.Execute(testutils.TestContext(t))
	require.NoError(t, result.GetErr())
	require.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("button"))

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x02, 0x00}))
	select {
	case event := <-events:
		assert.True(t, event.GetValue("button").GetBool())
	case <-time.After(20 * time.Second):
		require.FailNow(t, "the pre-registered consumer was never called")
	}
}

// A handle of another subscriber is none of this subscriber's business: an unsubscription request
// is handed to every subscriber it names a handle of.
func TestSubscriber_UnsubscribeIgnoresForeignHandles(t *testing.T) {
	connection, _ := newTestConnection(t)
	handle, _ := subscribe(t, connection, "button", "digital:9")
	collect, _ := collectEvents(t, handle)

	otherConnection, _ := newTestConnection(t)
	otherHandle, _ := subscribe(t, otherConnection, "other", "digital:9")
	require.NotNil(t, otherHandle)

	builder := connection.UnsubscriptionRequestBuilder()
	builder.AddHandles(otherHandle)
	request, err := builder.Build()
	require.NoError(t, err)
	require.NoError(t, (<-request.Execute(testutils.TestContext(t))).GetErr())

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x02, 0x00}))
	assert.Len(t, collect(), 1, "our own subscription must be untouched")
}

// Closing the connection drops the subscribers, so nothing is delivered afterwards.
func TestSubscriber_ClosedConnectionDeliversNothing(t *testing.T) {
	connection, _ := newTestConnection(t)
	handle, _ := subscribe(t, connection, "button", "digital:9")
	collect, _ := collectEvents(t, handle)

	require.NoError(t, connection.Close())

	connection.handleIncomingMessage(readWriteModel.NewFirmataMessageDigitalIO(1, []int8{0x02, 0x00}))
	assert.Empty(t, collect())
}

// The same for a subscription: reporting is only really switched on once the messages are out, so a
// claim whose messages failed has to be given back - otherwise the retry hands out a subscription
// which never delivers because the board was never told to report the pin.
func TestSubscriber_AFailedSendGivesThePinClaimBack(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())
	codec.failSends()

	handle, responseCode := subscribe(t, connection, "button", "digital:9")
	require.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, responseCode)
	assert.Nil(t, handle)
	func() {
		connection.pinMutex.Lock()
		defer connection.pinMutex.Unlock()
		assert.NotContains(t, connection.digitalPins, uint8(9), "the claim of a subscription that failed has to be given back")
	}()

	codec.allowSends()
	handle, responseCode = subscribe(t, connection, "button", "digital:9")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	require.NotNil(t, handle)
	assert.Equal(t, []byte{
		0xF4, 0x09, 0x00, // set pin mode: pin 9, input - sent again by the retry
		0xD1, 0x01, // report digital: port 1, on
	}, codec.sentBytesOf(t))
}

// An analog claim is given back just the same, and analog pins are the case where the claim is the
// only thing that is recorded: there is no set-pin-mode, just the report-analog message.
func TestSubscriber_AFailedAnalogSendGivesThePinClaimBack(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())
	codec.failSends()

	_, responseCode := subscribe(t, connection, "dial", "analog:2")
	require.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, responseCode)
	func() {
		connection.pinMutex.Lock()
		defer connection.pinMutex.Unlock()
		assert.NotContains(t, connection.analogPins, uint8(2))
	}()

	codec.allowSends()
	_, responseCode = subscribe(t, connection, "dial", "analog:2")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	assert.Equal(t, []byte{0xC2, 0x01}, codec.sentBytesOf(t))
}

// A firmata board emits no discrete events, and only changes are published, so an event subscription
// could only ever be served as a change-of-state one. plc4j quietly substitutes that behaviour;
// saying UNSUPPORTED is the honest answer and the one a caller can act on.
func TestSubscriber_SubscribeRefusesEventTags(t *testing.T) {
	connection, transportInstance := newTestConnection(t)

	builder := connection.SubscriptionRequestBuilder()
	builder.AddEventTagAddress("button", "digital:9")
	builder.AddChangeOfStateTagAddress("other", "digital:10")
	request, err := builder.Build()
	require.NoError(t, err)
	result := <-request.Execute(testutils.TestContext(t))

	require.NoError(t, result.GetErr())
	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, result.GetResponse().GetResponseCode("button"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("other"))
	eventHandle, _ := result.GetResponse().GetSubscriptionHandle("button")
	assert.Nil(t, eventHandle)
	// Only the tag which was accepted was configured on the board.
	assert.Equal(t, []byte{0xF4, 0x0A, 0x00, 0xD1, 0x01}, sentBytes(t, transportInstance))
}

// Same for a subscription which dies part-way: the set-pin-mode gets out but the report-digital
// which actually switches reporting on for the port does not, so the board says nothing about the
// pin. A claim left standing would make the retry hand out a handle that never delivers.
func TestSubscriber_APartialSendGivesThePinClaimBack(t *testing.T) {
	codec := newStubCodec()
	connection := NewConnection(DefaultConfiguration(), codec, map[string][]string{}, NewTagHandler())
	codec.failSendsAfter(1)

	handle, responseCode := subscribe(t, connection, "button", "digital:9")
	require.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, responseCode)
	assert.Nil(t, handle)
	func() {
		connection.pinMutex.Lock()
		defer connection.pinMutex.Unlock()
		assert.NotContains(t, connection.digitalPins, uint8(9), "a claim whose batch died part-way has to be given back")
	}()

	codec.allowSends()
	handle, responseCode = subscribe(t, connection, "button", "digital:9")
	require.Equal(t, apiModel.PlcResponseCode_OK, responseCode)
	require.NotNil(t, handle)
	assert.Equal(t, []byte{
		0xF4, 0x09, 0x00, // the set pin mode which did get out the first time
		0xF4, 0x09, 0x00, // and both messages again from the retry
		0xD1, 0x01,
	}, codec.sentBytesOf(t))
}
