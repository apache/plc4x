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

package profinet

import (
	"context"
	"encoding/binary"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
	"github.com/apache/plc4x/plc4go/spi/transports/virtualeth"
)

const testDeviceAlarmReference = 0x0003

// deviceRTA builds an RTA frame the way p-net sends one: VLAN tagged, to the controller's alarm
// reference from the device's.
func deviceRTA(deviceMAC, controllerMAC []byte, frameID uint16, pduType uint8, tack bool, sendSeq, ackSeq uint16, payload []byte) []byte {
	frame := append([]byte(nil), controllerMAC...)
	frame = append(frame, deviceMAC...)
	frame = binary.BigEndian.AppendUint16(frame, 0x8100)
	frame = binary.BigEndian.AppendUint16(frame, 6<<13)
	frame = binary.BigEndian.AppendUint16(frame, etherTypeProfinet)
	frame = binary.BigEndian.AppendUint16(frame, frameID)
	frame = binary.BigEndian.AppendUint16(frame, controllerAlarmReference)
	frame = binary.BigEndian.AppendUint16(frame, testDeviceAlarmReference)
	flags := byte(1)
	if tack {
		flags |= rtaTackFlag
	}
	frame = append(frame, rtaVersion<<4|pduType, flags)
	frame = binary.BigEndian.AppendUint16(frame, sendSeq)
	frame = binary.BigEndian.AppendUint16(frame, ackSeq)
	frame = binary.BigEndian.AppendUint16(frame, uint16(len(payload)))
	return append(frame, payload...)
}

// processAlarm is the AlarmNotification p-net's sample device sends for its button: a process
// alarm on slot 2 subslot 1 with USI 0x0010 and one byte of payload.
func processAlarm(sequence uint16, payload byte) []byte {
	body := binary.BigEndian.AppendUint16(nil, 0x0002) // Process
	body = binary.BigEndian.AppendUint32(body, 0)
	body = binary.BigEndian.AppendUint16(body, 2)
	body = binary.BigEndian.AppendUint16(body, 1)
	body = binary.BigEndian.AppendUint32(body, 0x30)
	body = binary.BigEndian.AppendUint32(body, 0x130)
	body = binary.BigEndian.AppendUint16(body, sequence|1<<13)
	body = binary.BigEndian.AppendUint16(body, 0x0010)
	body = append(body, payload)
	block := binary.BigEndian.AppendUint16(nil, blockTypeAlarmNotificationLow)
	block = binary.BigEndian.AppendUint16(block, uint16(2+len(body)))
	block = append(block, 1, 0)
	return append(block, body...)
}

func TestParseAlarmNotification(t *testing.T) {
	deviceMAC, controllerMAC := virtualeth.RandomMAC(), virtualeth.RandomMAC()
	p, ok := parseRTA(deviceRTA(deviceMAC, controllerMAC, frameIDAlarmLow, rtaTypeData, true, 0xFFFF, 0xFFFE, processAlarm(5, 0x2A)), deviceMAC)
	require.True(t, ok)
	assert.Equal(t, uint16(testDeviceAlarmReference), p.src)
	assert.Equal(t, uint16(controllerAlarmReference), p.dst)
	assert.True(t, p.tack)
	assert.Equal(t, uint16(0xFFFF), p.sendSeq)
	alarm, isNotification, err := parseAlarmNotification(false, p.payload)
	require.NoError(t, err)
	require.True(t, isNotification)
	assert.Equal(t, Alarm{Type: 2, Slot: 2, Subslot: 1, ModuleIdent: 0x30, SubmoduleIdent: 0x130, Specifier: 5 | 1<<13, USI: 0x10, Data: []byte{0x2A}}, alarm)
	assert.Equal(t, "Process", alarm.TypeName())
	assert.Equal(t, uint16(5), alarm.Sequence())

	_, ok = parseRTA(deviceRTA(controllerMAC, deviceMAC, frameIDAlarmLow, rtaTypeData, true, 0, 0, nil), deviceMAC)
	assert.False(t, ok, "a frame from another station")

	ack := alarmAckBlock(alarm)
	assert.Equal(t, []byte{0x80, 0x02, 0x00, 0x12, 1, 0, 0x00, 0x02, 0, 0, 0, 0, 0x00, 0x02, 0x00, 0x01, 0x20, 0x05, 0, 0, 0, 0}, ack)
}

func TestParseAlarmTag(t *testing.T) {
	for address, want := range map[string]AlarmTag{
		"ALARM":     {Slot: -1, Subslot: -1},
		"ALARM:2":   {Slot: 2, Subslot: -1},
		"ALARM:2.1": {Slot: 2, Subslot: 1},
	} {
		tag, err := NewTagHandler().ParseTag(address)
		require.NoError(t, err, address)
		assert.Equal(t, want, tag, address)
		assert.Equal(t, address, tag.GetAddressString())
	}
	_, err := NewTagHandler().ParseTag("ALARM:99999")
	assert.Error(t, err)
	_, err = NewTagHandler().ParseTag("ALARMS")
	assert.Error(t, err)
	assert.True(t, AlarmTag{Slot: 2, Subslot: -1}.matches(Alarm{Slot: 2, Subslot: 7}))
	assert.False(t, AlarmTag{Slot: 2, Subslot: 1}.matches(Alarm{Slot: 2, Subslot: 7}))
}

// alarmPeer is a controller relation and a station standing in for the device on one segment.
type alarmPeer struct {
	ar     *applicationRelation
	device *virtualeth.TransportInstance
	alarms chan Alarm
	cancel context.CancelFunc
}

func newAlarmPeer(t *testing.T) *alarmPeer {
	segment := t.Name()
	filter := rawsocket.Filter{EtherTypes: []uint16{etherTypeProfinet}}
	controller := virtualeth.NewTransportInstance(segment, filter, virtualeth.RandomMAC())
	device := virtualeth.NewTransportInstance(segment, filter, virtualeth.RandomMAC())
	require.NoError(t, controller.Connect(t.Context()))
	require.NoError(t, device.Connect(t.Context()))
	t.Cleanup(func() { _ = controller.Close(); _ = device.Close() })
	ar := newApplicationRelation(arSettings{requestTimeout: time.Second}, device_(device), nil, controller, nil, zerolog.Nop())
	ar.deviceAlarmReference = testDeviceAlarmReference
	peer := &alarmPeer{ar: ar, device: device, alarms: make(chan Alarm, 8)}
	ar.onAlarm = func(a Alarm) { peer.alarms <- a }
	ctx, cancel := context.WithCancel(context.Background())
	peer.cancel = cancel
	ar.wg.Go(func() { ar.sendAlarmAcks(ctx, ar.alarmLow) })
	t.Cleanup(func() { cancel(); ar.wg.Wait() })
	return peer
}

func device_(st *virtualeth.TransportInstance) device {
	return device{name: "test-device", mac: st.HardwareAddress()}
}

// fromDevice delivers a frame from the device to the relation.
func (p *alarmPeer) fromDevice(pduType uint8, tack bool, sendSeq, ackSeq uint16, payload []byte) {
	p.ar.acceptAlarm(context.Background(), deviceRTA(p.device.HardwareAddress(), p.ar.local, frameIDAlarmLow, pduType, tack, sendSeq, ackSeq, payload))
}

// next is the next RTA frame the controller sent the device.
func (p *alarmPeer) next(t *testing.T) rtaPDU {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	frame, err := readFrame(ctx, p.device)
	require.NoError(t, err)
	// The device's parse: the frame comes from the controller.
	pdu, ok := parseRTA(frame, p.ar.local)
	require.True(t, ok)
	assert.Equal(t, uint16(testDeviceAlarmReference), pdu.dst)
	assert.Equal(t, uint16(controllerAlarmReference), pdu.src)
	return pdu
}

func TestAlarmAcknowledgement(t *testing.T) {
	p := newAlarmPeer(t)

	// The first notification: a transport acknowledgement, the alarm, then the AlarmAck.
	p.fromDevice(rtaTypeData, true, 0xFFFF, 0xFFFE, processAlarm(1, 0x01))
	ack := p.next(t)
	assert.Equal(t, uint8(rtaTypeAck), ack.pduType)
	assert.False(t, ack.tack)
	assert.Equal(t, uint16(0xFFFF), ack.ackSeq)
	assert.Equal(t, uint16(0xFFFE), ack.sendSeq)
	alarm := <-p.alarms
	assert.Equal(t, uint16(2), alarm.Slot)
	data := p.next(t)
	assert.Equal(t, uint8(rtaTypeData), data.pduType)
	assert.True(t, data.tack)
	assert.Equal(t, uint16(0xFFFF), data.sendSeq)
	assert.Equal(t, uint16(0xFFFF), data.ackSeq)
	assert.Equal(t, alarmAckBlock(alarm), data.payload)

	// Unacknowledged, the AlarmAck comes again with the same sequence number.
	again := p.next(t)
	assert.Equal(t, data.sendSeq, again.sendSeq)
	p.fromDevice(rtaTypeAck, false, 0xFFFE, 0xFFFF, nil)

	// A repeat of the notification (our ACK lost) is acknowledged again but not delivered twice.
	p.fromDevice(rtaTypeData, true, 0xFFFF, 0xFFFE, processAlarm(1, 0x01))
	assert.Equal(t, uint8(rtaTypeAck), p.next(t).pduType)
	assert.Empty(t, p.alarms)

	// The second notification counts on from 0, and so does the second AlarmAck.
	p.fromDevice(rtaTypeData, true, 0x0000, 0xFFFF, processAlarm(2, 0x02))
	ack = p.next(t)
	assert.Equal(t, uint16(0x0000), ack.ackSeq)
	<-p.alarms
	data = p.next(t)
	for data.pduType != rtaTypeData || data.sendSeq != 0x0000 {
		data = p.next(t) // skip repeats of the first AlarmAck still in flight
	}
	assert.Equal(t, uint16(0x0000), data.ackSeq)
	p.fromDevice(rtaTypeAck, false, 0x0000, 0x0000, nil)
	assert.Eventually(t, func() bool {
		ch := p.ar.alarmLow
		ch.mutex.Lock()
		defer ch.mutex.Unlock()
		return ch.sendSeq == 0x0001 && !ch.waitingTack
	}, time.Second, 10*time.Millisecond, "the second AlarmAck is acknowledged")
}

func TestAlarmErrorEndsRelation(t *testing.T) {
	p := newAlarmPeer(t)
	p.fromDevice(rtaTypeErr, false, 0xFFFF, 0xFFFE, []byte{0xCF, 0x81, 0xFD, 0x00})
	select {
	case err := <-p.ar.lost:
		assert.ErrorContains(t, err, "0xCF81FD00")
	case <-time.After(time.Second):
		t.Fatal("an alarm error frame does not end the relation")
	}
}

func TestAlarmSubscription(t *testing.T) {
	s := newAlarmSubscriber()
	defer s.Close()
	router := subscriberRouter{polling: nil, alarms: s}
	request, err := spiModel.NewDefaultPlcSubscriptionRequestBuilder(NewTagHandler(), nil, router).
		AddEventTagAddress("slot2", "ALARM:2").
		AddEventTagAddress("slot3", "ALARM:3").Build()
	require.NoError(t, err)
	result := <-request.Execute(t.Context())
	require.NoError(t, result.GetErr())
	response := result.GetResponse()
	require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("slot2"))
	var events []apiModel.PlcSubscriptionEvent
	for _, name := range response.GetTagNames() {
		handle, err := response.GetSubscriptionHandle(name)
		require.NoError(t, err)
		handle.Register(func(event apiModel.PlcSubscriptionEvent) { events = append(events, event) })
	}
	s.publish(Alarm{Type: 2, Slot: 2, Subslot: 1, USI: 0x10, Data: []byte{7}})
	require.Len(t, events, 1)
	value := events[0].GetValue("slot2")
	assert.Equal(t, "Process", value.GetStruct()["alarmTypeName"].GetString())
	assert.Equal(t, uint16(1), value.GetStruct()["subslot"].GetUint16())
	assert.Equal(t, []byte{7}, value.GetStruct()["data"].GetRaw())

	mixed, err := spiModel.NewDefaultPlcSubscriptionRequestBuilder(NewTagHandler(), nil, router).
		AddEventTagAddress("alarm", "ALARM").
		AddCyclicTagAddress("value", "1.1.INPUT:BYTE", time.Second).Build()
	require.NoError(t, err)
	assert.Error(t, (<-mixed.Execute(t.Context())).GetErr())
}
