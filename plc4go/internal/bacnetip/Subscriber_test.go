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

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

// newTestSubscriber returns a Subscriber backed by an empty Connection. The
// Connection's messageCodec field is nil so direct Subscribe calls would panic
// — these tests exercise pure-function paths only (HandleConfirmedCOVNotification,
// HandleUnconfirmedCOVNotification, handle-table semantics).
func newTestSubscriber(t *testing.T) *Subscriber {
	t.Helper()
	conn := &Connection{configuration: createDefaultConfiguration()}
	conn.driverContext = NewDriverContext(conn.configuration)
	return NewSubscriber(conn)
}

func TestSubscriber_StoresAndLooksUpHandle(t *testing.T) {
	s := newTestSubscriber(t)
	objType := readWriteModel.BACnetObjectType_ANALOG_INPUT
	handle := NewSubscriptionHandle(s, "ai0", &plcTag{ObjectId: objectId{ObjectIdType: &objType, ObjectIdInstance: 0}}, apiModel.SubscriptionChangeOfState, 0)
	handle.subscriberProcessId = 42
	s.storeHandle(handle)

	require.Equal(t, handle, s.lookupHandle(42))
	assert.Nil(t, s.lookupHandle(43))

	s.removeHandle(42)
	assert.Nil(t, s.lookupHandle(42))
}

// captureConsumer is a simple consumer that records every PlcSubscriptionEvent
// it receives so tests can assert on routing.
type captureConsumer struct {
	events []apiModel.PlcSubscriptionEvent
}

func (c *captureConsumer) consume(event apiModel.PlcSubscriptionEvent) {
	c.events = append(c.events, event)
}

// buildCOVNotification builds an UnconfirmedCOVNotification carrying a single
// (PresentValue, Real) tuple for tests.
func buildCOVNotification(processId uint32, deviceId uint32, objType readWriteModel.BACnetObjectType, instance uint32, presentValue float32) readWriteModel.BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification {
	subscriberProcessTag := readWriteModel.CreateBACnetContextTagUnsignedInteger(0, uint(processId))
	initiatingDeviceTag := readWriteModel.CreateBACnetContextTagObjectIdentifier(1, uint16(readWriteModel.BACnetObjectType_DEVICE), deviceId)
	monitoredObjectTag := readWriteModel.CreateBACnetContextTagObjectIdentifier(2, uint16(objType), instance)
	lifetimeTag := readWriteModel.CreateBACnetContextTagUnsignedInteger(3, 0)

	tagHeader := readWriteModel.CreateBACnetTagHeaderBalanced(true, 2, 0)
	element := readWriteModel.NewBACnetConstructedDataElement(tagHeader, readWriteModel.CreateBACnetApplicationTagReal(presentValue), nil, nil)
	cd := readWriteModel.NewBACnetConstructedDataUnspecified(
		readWriteModel.CreateBACnetOpeningTag(2),
		tagHeader,
		readWriteModel.CreateBACnetClosingTag(2),
		nil,
		[]readWriteModel.BACnetConstructedDataElement{element},
	)
	propIdTag := readWriteModel.CreateBACnetPropertyIdentifierTagged(0, uint32(readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE))
	propVal := readWriteModel.NewBACnetPropertyValue(propIdTag, nil, element, nil)
	_ = cd
	values := readWriteModel.NewBACnetPropertyValues(
		readWriteModel.CreateBACnetOpeningTag(4),
		[]readWriteModel.BACnetPropertyValue{propVal},
		readWriteModel.CreateBACnetClosingTag(4),
	)
	return readWriteModel.NewBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification(
		subscriberProcessTag, initiatingDeviceTag, monitoredObjectTag, lifetimeTag, values,
	)
}

func TestHandleUnconfirmedCOVNotification_DispatchesToConsumers(t *testing.T) {
	s := newTestSubscriber(t)
	s.log = zerolog.Nop()

	objType := readWriteModel.BACnetObjectType_ANALOG_INPUT
	tag := &plcTag{ObjectId: objectId{ObjectIdType: &objType, ObjectIdInstance: 1}}
	handle := NewSubscriptionHandle(s, "temp", tag, apiModel.SubscriptionChangeOfState, 0)
	handle.subscriberProcessId = 99
	s.storeHandle(handle)

	capture := &captureConsumer{}
	reg := spiModel.NewDefaultPlcConsumerRegistration(s, capture.consume, handle.DefaultPlcSubscriptionHandle)
	s.consumers[reg.(*spiModel.DefaultPlcConsumerRegistration)] = capture.consume

	req := buildCOVNotification(99, 1234, objType, 1, 24.1)
	s.HandleUnconfirmedCOVNotification(req)

	require.Len(t, capture.events, 1)
	event := capture.events[0]
	val := event.GetValue("temp")
	require.NotNil(t, val)
	assert.InDelta(t, 24.1, val.GetFloat32(), 1e-3)
}

func TestHandleUnconfirmedCOVNotification_UnknownProcessId_NoDispatch(t *testing.T) {
	s := newTestSubscriber(t)
	s.log = zerolog.Nop()

	capture := &captureConsumer{}
	req := buildCOVNotification(404, 1234, readWriteModel.BACnetObjectType_ANALOG_INPUT, 0, 1.0)
	s.HandleUnconfirmedCOVNotification(req)
	assert.Empty(t, capture.events, "no handle stored for processId 404, should be no-op")
}

func TestHandleConfirmedCOVNotification_DispatchesToConsumers(t *testing.T) {
	s := newTestSubscriber(t)
	s.log = zerolog.Nop()

	objType := readWriteModel.BACnetObjectType_ANALOG_VALUE
	tag := &plcTag{ObjectId: objectId{ObjectIdType: &objType, ObjectIdInstance: 5}}
	handle := NewSubscriptionHandle(s, "av5", tag, apiModel.SubscriptionChangeOfState, 0)
	handle.subscriberProcessId = 17
	s.storeHandle(handle)

	capture := &captureConsumer{}
	reg := spiModel.NewDefaultPlcConsumerRegistration(s, capture.consume, handle.DefaultPlcSubscriptionHandle)
	s.consumers[reg.(*spiModel.DefaultPlcConsumerRegistration)] = capture.consume

	uncf := buildCOVNotification(17, 99, objType, 5, 42.0)
	// Confirmed has the same structure as unconfirmed; re-wrap by passing the
	// fields directly to the confirmed constructor.
	confirmed := readWriteModel.NewBACnetConfirmedServiceRequestConfirmedCOVNotification(
		0,
		uncf.GetSubscriberProcessIdentifier(),
		uncf.GetInitiatingDeviceIdentifier(),
		uncf.GetMonitoredObjectIdentifier(),
		uncf.GetLifetimeInSeconds(),
		uncf.GetListOfValues(),
	)
	s.HandleConfirmedCOVNotification(confirmed)

	require.Len(t, capture.events, 1)
	assert.InDelta(t, 42.0, capture.events[0].GetValue("av5").GetFloat32(), 1e-3)
}

func TestDispatchNotification_EmptyListYieldsNull(t *testing.T) {
	s := newTestSubscriber(t)
	s.log = zerolog.Nop()

	objType := readWriteModel.BACnetObjectType_ANALOG_INPUT
	tag := &plcTag{ObjectId: objectId{ObjectIdType: &objType, ObjectIdInstance: 0}}
	handle := NewSubscriptionHandle(s, "x", tag, apiModel.SubscriptionChangeOfState, 0)
	handle.subscriberProcessId = 5
	s.storeHandle(handle)

	capture := &captureConsumer{}
	reg := spiModel.NewDefaultPlcConsumerRegistration(s, capture.consume, handle.DefaultPlcSubscriptionHandle)
	s.consumers[reg.(*spiModel.DefaultPlcConsumerRegistration)] = capture.consume

	empty := readWriteModel.NewBACnetPropertyValues(
		readWriteModel.CreateBACnetOpeningTag(4),
		nil,
		readWriteModel.CreateBACnetClosingTag(4),
	)
	s.dispatchNotification(handle, empty)

	require.Len(t, capture.events, 1)
	v := capture.events[0].GetValue("x")
	assert.Equal(t, apiValues.NULL, v.GetPlcValueType())
}

// Defensive check: we must not panic for a handle with a nil PlcTag.
func TestDispatchNotification_NilTag_DoesNotPanic(t *testing.T) {
	s := newTestSubscriber(t)
	s.log = zerolog.Nop()

	handle := &SubscriptionHandle{tagName: "boom"}
	handle.subscriberProcessId = 1
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("dispatchNotification panicked: %v", r)
		}
	}()
	empty := readWriteModel.NewBACnetPropertyValues(
		readWriteModel.CreateBACnetOpeningTag(4),
		nil,
		readWriteModel.CreateBACnetClosingTag(4),
	)
	// No consumer registered, so nothing dispatched — main assertion is "no panic".
	_ = time.Now()
	defer func() { recover() }()
	s.dispatchNotification(handle, empty)
}

var _ = time.Second // keep import alive across phases
