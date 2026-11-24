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

// Alarms (IEC 61158-6-10): the device sends AlarmNotifications as acyclic real-time (RTA) frames,
// EtherType 0x8892 with frame id 0xFC01 (high priority) or 0xFE01 (low priority). Each priority is
// its own channel with its own sequence numbers. Every DATA frame with the TACK flag gets a
// transport acknowledgement (an ACK frame); the controller then answers the notification with an
// AlarmAck in a DATA frame of its own, which the device acknowledges in turn. The sequence numbers
// start at 0xFFFF and count 0..0x7FFF, as p-net and the standard have them.

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync"
	"time"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/transports/rawsocket"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

const (
	frameIDAlarmHigh = 0xFC01
	frameIDAlarmLow  = 0xFE01

	rtaTypeData = 1
	rtaTypeNack = 2
	rtaTypeAck  = 3
	rtaTypeErr  = 4
	rtaVersion  = 1
	// rtaHeaderLength is the fixed part of an RTA PDU: destination and source alarm references, PDU
	// type, flags, send and acknowledge sequence numbers and the length of what follows.
	rtaHeaderLength = 12
	rtaWindowSize   = 1
	rtaTackFlag     = 0x10

	blockTypeAlarmNotificationHigh = 0x0001
	blockTypeAlarmNotificationLow  = 0x0002
	blockTypeAlarmAckHigh          = 0x8001
	blockTypeAlarmAckLow           = 0x8002

	// controllerAlarmReference is the controller's end of the alarm CR (the AlarmCRBlockReq's
	// LocalAlarmReference); the device names its own in the AlarmCRBlockRes.
	controllerAlarmReference = 0x0000
	// rtaTimeoutFactor and rtaRetries are what the AlarmCRBlockReq asks of the device: wait
	// 100 ms times the factor for a transport acknowledgement, and retry that many times.
	rtaTimeoutFactor = 1
	rtaRetries       = 3
	// alarmAckTimeout is how long the controller waits for the device to acknowledge one of its
	// AlarmAck frames before sending it again.
	alarmAckTimeout = 300 * time.Millisecond

	minEthernetFrame = 60
)

// Alarm is one alarm notification from the device.
type Alarm struct {
	High                        bool
	Type                        uint16
	API                         uint32
	Slot, Subslot               uint16
	ModuleIdent, SubmoduleIdent uint32
	// Specifier is the AlarmSpecifier: the sequence number in bits 0-10, then the channel,
	// manufacturer, submodule and AR diagnosis flags in bits 11, 12, 13 and 15.
	Specifier uint16
	// USI is the user structure identifier of the alarm's payload, 0 when there is none, and Data
	// the payload after it.
	USI  uint16
	Data []byte
}

func (a Alarm) Sequence() uint16 { return a.Specifier & 0x07FF }

// alarmTypeNames are the AlarmType values of IEC 61158-6-10.
var alarmTypeNames = map[uint16]string{
	0x0001: "Diagnosis",
	0x0002: "Process",
	0x0003: "Pull",
	0x0004: "Plug",
	0x0005: "Status",
	0x0006: "Update",
	0x0007: "Media redundancy",
	0x0008: "Controlled by supervisor",
	0x0009: "Released",
	0x000A: "Plug wrong submodule",
	0x000B: "Return of submodule",
	0x000C: "Diagnosis disappears",
	0x000D: "Multicast communication mismatch",
	0x000E: "Port data change notification",
	0x000F: "Sync data changed notification",
	0x0010: "Isochronous mode problem",
	0x0011: "Network component problem",
	0x0012: "Time data changed notification",
	0x0013: "Dynamic frame packing problem",
	0x0014: "MRPD problem",
	0x0016: "Multiple interface mismatch",
	0x001E: "Upload and retrieval notification",
	0x001F: "Pull module",
}

func (a Alarm) TypeName() string {
	if name, ok := alarmTypeNames[a.Type]; ok {
		return name
	}
	if a.Type >= 0x0020 && a.Type <= 0x007F {
		return "Manufacturer specific"
	}
	return fmt.Sprintf("Reserved 0x%04X", a.Type)
}

func (a Alarm) Priority() string {
	if a.High {
		return "HIGH"
	}
	return "LOW"
}

var alarmFields = []string{"priority", "alarmType", "alarmTypeName", "api", "slot", "subslot", "moduleIdent",
	"submoduleIdent", "sequence", "channelDiagnosis", "manufacturerDiagnosis", "submoduleDiagnosis", "arDiagnosis",
	"usi", "data"}

// PlcValue is how a subscription event carries the alarm.
func (a Alarm) PlcValue() apiValues.PlcValue {
	return spiValues.NewPlcStructOrdered(map[string]apiValues.PlcValue{
		"priority":              spiValues.NewPlcSTRING(a.Priority()),
		"alarmType":             spiValues.NewPlcUINT(a.Type),
		"alarmTypeName":         spiValues.NewPlcSTRING(a.TypeName()),
		"api":                   spiValues.NewPlcUDINT(a.API),
		"slot":                  spiValues.NewPlcUINT(a.Slot),
		"subslot":               spiValues.NewPlcUINT(a.Subslot),
		"moduleIdent":           spiValues.NewPlcUDINT(a.ModuleIdent),
		"submoduleIdent":        spiValues.NewPlcUDINT(a.SubmoduleIdent),
		"sequence":              spiValues.NewPlcUINT(a.Sequence()),
		"channelDiagnosis":      spiValues.NewPlcBOOL(a.Specifier&(1<<11) != 0),
		"manufacturerDiagnosis": spiValues.NewPlcBOOL(a.Specifier&(1<<12) != 0),
		"submoduleDiagnosis":    spiValues.NewPlcBOOL(a.Specifier&(1<<13) != 0),
		"arDiagnosis":           spiValues.NewPlcBOOL(a.Specifier&(1<<15) != 0),
		"usi":                   spiValues.NewPlcUINT(a.USI),
		"data":                  spiValues.NewPlcRawByteArray(append([]byte(nil), a.Data...)),
	}, alarmFields)
}

// rtaPDU is the fixed part of an RTA frame and what follows it.
type rtaPDU struct {
	frameID         uint16
	dst, src        uint16
	pduType         uint8
	tack            bool
	sendSeq, ackSeq uint16
	payload         []byte
}

// parseRTA reads an RTA frame from the device; ok is false for any other frame.
func parseRTA(frame []byte, deviceMAC []byte) (rtaPDU, bool) {
	if len(frame) < 14 || string(frame[6:12]) != string(deviceMAC) {
		return rtaPDU{}, false
	}
	offset := 12
	if binary.BigEndian.Uint16(frame[offset:]) == 0x8100 {
		offset += 4
	}
	if len(frame) < offset+4+rtaHeaderLength || binary.BigEndian.Uint16(frame[offset:]) != etherTypeProfinet {
		return rtaPDU{}, false
	}
	offset += 2
	p := rtaPDU{frameID: binary.BigEndian.Uint16(frame[offset:])}
	if p.frameID != frameIDAlarmHigh && p.frameID != frameIDAlarmLow {
		return rtaPDU{}, false
	}
	h := frame[offset+2:]
	p.dst = binary.BigEndian.Uint16(h[0:])
	p.src = binary.BigEndian.Uint16(h[2:])
	if h[4]>>4 != rtaVersion {
		return rtaPDU{}, false
	}
	p.pduType = h[4] & 0x0F
	p.tack = h[5]&rtaTackFlag != 0
	p.sendSeq = binary.BigEndian.Uint16(h[6:])
	p.ackSeq = binary.BigEndian.Uint16(h[8:])
	length := int(binary.BigEndian.Uint16(h[10:]))
	if len(h) < rtaHeaderLength+length {
		return rtaPDU{}, false
	}
	p.payload = h[rtaHeaderLength : rtaHeaderLength+length]
	return p, true
}

// parseAlarmNotification reads the AlarmNotification block of a DATA frame. ok is false when the
// block is something else (an AlarmAck from the device answers an alarm of the controller's, and
// this controller sends none).
func parseAlarmNotification(high bool, block []byte) (Alarm, bool, error) {
	if len(block) < 6 {
		return Alarm{}, false, errors.New("the alarm block is shorter than its header")
	}
	blockType := binary.BigEndian.Uint16(block[0:])
	if blockType != blockTypeAlarmNotificationHigh && blockType != blockTypeAlarmNotificationLow {
		return Alarm{}, false, nil
	}
	length := int(binary.BigEndian.Uint16(block[2:]))
	if length < 2 || len(block) < 4+length {
		return Alarm{}, false, errors.Errorf("the alarm notification announces %d bytes, %d follow", length, len(block)-4)
	}
	body := block[6 : 4+length]
	if len(body) < 22-2 {
		return Alarm{}, false, errors.Errorf("the alarm notification has %d bytes, too few", len(body))
	}
	a := Alarm{
		High:           high,
		Type:           binary.BigEndian.Uint16(body[0:]),
		API:            binary.BigEndian.Uint32(body[2:]),
		Slot:           binary.BigEndian.Uint16(body[6:]),
		Subslot:        binary.BigEndian.Uint16(body[8:]),
		ModuleIdent:    binary.BigEndian.Uint32(body[10:]),
		SubmoduleIdent: binary.BigEndian.Uint32(body[14:]),
		Specifier:      binary.BigEndian.Uint16(body[18:]),
	}
	if rest := body[20:]; len(rest) >= 2 {
		a.USI = binary.BigEndian.Uint16(rest)
		a.Data = append([]byte(nil), rest[2:]...)
	}
	return a, true, nil
}

// alarmAckBlock is the AlarmAck answering an alarm: its type, API, slot, subslot and specifier,
// and PNIO status OK.
func alarmAckBlock(a Alarm) []byte {
	block := make([]byte, 22)
	blockType := uint16(blockTypeAlarmAckLow)
	if a.High {
		blockType = blockTypeAlarmAckHigh
	}
	binary.BigEndian.PutUint16(block[0:], blockType)
	binary.BigEndian.PutUint16(block[2:], uint16(len(block)-4))
	block[4], block[5] = blockVersionHigh, blockVersionLow
	binary.BigEndian.PutUint16(block[6:], a.Type)
	binary.BigEndian.PutUint32(block[8:], a.API)
	binary.BigEndian.PutUint16(block[12:], a.Slot)
	binary.BigEndian.PutUint16(block[14:], a.Subslot)
	binary.BigEndian.PutUint16(block[16:], a.Specifier)
	// block[18:22]: PNIO status 0, OK.
	return block
}

// nextSequence counts an RTA sequence number on: 0xFFFF is the start, then 0 to 0x7FFF.
func nextSequence(n uint16) uint16 {
	return (n + 1) & 0x7FFF
}

// alarmChannel is one priority's RTA state: the receiving side (expected sequence numbers) and
// the sending side for the controller's AlarmAcks.
type alarmChannel struct {
	high    bool
	frameID uint16

	mutex             sync.Mutex
	expSeq, expSeqO   uint16
	sendSeq, sendSeqO uint16
	waitingTack       bool
	tack              chan struct{}

	acks chan Alarm
}

func newAlarmChannel(high bool) *alarmChannel {
	frameID := uint16(frameIDAlarmLow)
	if high {
		frameID = frameIDAlarmHigh
	}
	return &alarmChannel{
		high: high, frameID: frameID,
		expSeq: 0xFFFF, expSeqO: 0xFFFE, sendSeq: 0xFFFF, sendSeqO: 0xFFFE,
		tack: make(chan struct{}, 1),
		acks: make(chan Alarm, 32),
	}
}

// acknowledged takes the acknowledge sequence number of a frame from the device: when it names the
// controller's outstanding DATA frame, that frame has arrived.
func (ch *alarmChannel) acknowledged(ackSeq uint16) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	if ch.waitingTack && ackSeq == ch.sendSeq {
		ch.waitingTack = false
		ch.sendSeqO = ch.sendSeq
		ch.sendSeq = nextSequence(ch.sendSeq)
		select {
		case ch.tack <- struct{}{}:
		default:
		}
	}
}

// rtaFrame is an RTA frame on a channel to the device.
func (ar *applicationRelation) rtaFrame(ch *alarmChannel, pduType uint8, tack bool, sendSeq, ackSeq uint16, payload []byte) []byte {
	frame := make([]byte, 0, max(minEthernetFrame, 16+rtaHeaderLength+len(payload)))
	frame = append(frame, ar.device.mac...)
	frame = append(frame, ar.local...)
	frame = binary.BigEndian.AppendUint16(frame, etherTypeProfinet)
	frame = binary.BigEndian.AppendUint16(frame, ch.frameID)
	frame = binary.BigEndian.AppendUint16(frame, ar.deviceAlarmReference)
	frame = binary.BigEndian.AppendUint16(frame, controllerAlarmReference)
	flags := byte(rtaWindowSize)
	if tack {
		flags |= rtaTackFlag
	}
	frame = append(frame, rtaVersion<<4|pduType, flags)
	frame = binary.BigEndian.AppendUint16(frame, sendSeq)
	frame = binary.BigEndian.AppendUint16(frame, ackSeq)
	frame = binary.BigEndian.AppendUint16(frame, uint16(len(payload)))
	frame = append(frame, payload...)
	for len(frame) < minEthernetFrame {
		frame = append(frame, 0)
	}
	return frame
}

func (ar *applicationRelation) sendRTA(ctx context.Context, frame []byte) error {
	record, err := rawsocket.EncodeRecord(frame)
	if err != nil {
		return err
	}
	writeCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return ar.station.Write(writeCtx, record)
}

// sendTransportAck acknowledges the device's DATA frames up to the channel's last received one.
func (ar *applicationRelation) sendTransportAck(ctx context.Context, ch *alarmChannel) {
	ch.mutex.Lock()
	frame := ar.rtaFrame(ch, rtaTypeAck, false, ch.sendSeqO, ch.expSeqO, nil)
	ch.mutex.Unlock()
	if err := ar.sendRTA(ctx, frame); err != nil && ctx.Err() == nil {
		ar.log.Debug().Err(err).Msg("error sending an alarm transport acknowledgement")
	}
}

// acceptAlarm handles an RTA frame of this relation; it is false for any other frame. It runs on
// the receive loop, so it only queues the AlarmAck for the channel's sender.
func (ar *applicationRelation) acceptAlarm(ctx context.Context, frame []byte) bool {
	p, ok := parseRTA(frame, ar.device.mac)
	if !ok || p.dst != controllerAlarmReference {
		return false
	}
	ch := ar.alarmLow
	if p.frameID == frameIDAlarmHigh {
		ch = ar.alarmHigh
	}
	if ch == nil {
		return true
	}
	switch p.pduType {
	case rtaTypeAck:
		ch.acknowledged(p.ackSeq)
	case rtaTypeErr:
		status := uint32(0)
		if len(p.payload) >= 4 {
			status = binary.BigEndian.Uint32(p.payload)
		}
		ar.fail(errors.Errorf("%s ended the application relation with an alarm error frame (PNIO status 0x%08X)", ar.device.name, status))
	case rtaTypeNack:
		// The device missed a frame; the sender's retry covers it.
	case rtaTypeData:
		ch.acknowledged(p.ackSeq)
		if !p.tack {
			return true
		}
		ch.mutex.Lock()
		switch p.sendSeq {
		case ch.expSeq:
			ch.expSeqO = ch.expSeq
			ch.expSeq = nextSequence(ch.expSeq)
			ch.mutex.Unlock()
		case ch.expSeqO:
			// A repeat: the device did not get the transport acknowledgement.
			ch.mutex.Unlock()
			ar.sendTransportAck(ctx, ch)
			return true
		default:
			ch.mutex.Unlock()
			ar.log.Debug().Uint16("sequence", p.sendSeq).Msg("alarm frame out of sequence, ignored")
			return true
		}
		ar.sendTransportAck(ctx, ch)
		alarm, isNotification, err := parseAlarmNotification(ch.high, p.payload)
		if err != nil {
			ar.log.Warn().Err(err).Msg("unreadable alarm notification")
			return true
		}
		if !isNotification {
			return true
		}
		ar.log.Info().Str("device", ar.device.name).Str("priority", alarm.Priority()).Str("type", alarm.TypeName()).
			Uint16("slot", alarm.Slot).Uint16("subslot", alarm.Subslot).Uint16("usi", alarm.USI).Msg("Profinet alarm")
		select {
		case ch.acks <- alarm:
		default:
			ar.log.Error().Msg("too many alarms waiting for their acknowledgement; one is not acknowledged")
		}
		if ar.onAlarm != nil {
			ar.onAlarm(alarm)
		}
	}
	return true
}

// sendAlarmAcks answers the channel's alarms with AlarmAck frames, one at a time, each sent again
// until the device acknowledges it.
func (ar *applicationRelation) sendAlarmAcks(ctx context.Context, ch *alarmChannel) {
	for {
		var alarm Alarm
		select {
		case <-ctx.Done():
			return
		case alarm = <-ch.acks:
		}
		payload := alarmAckBlock(alarm)
		acknowledged := false
		for attempt := 0; attempt <= rtaRetries && !acknowledged; attempt++ {
			ch.mutex.Lock()
			ch.waitingTack = true
			frame := ar.rtaFrame(ch, rtaTypeData, true, ch.sendSeq, ch.expSeqO, payload)
			ch.mutex.Unlock()
			if err := ar.sendRTA(ctx, frame); err != nil && ctx.Err() == nil {
				ar.log.Debug().Err(err).Msg("error sending an alarm acknowledgement")
			}
			select {
			case <-ctx.Done():
				return
			case <-ch.tack:
				acknowledged = true
			case <-time.After(alarmAckTimeout):
			}
		}
		if !acknowledged {
			ar.fail(errors.Errorf("%s did not acknowledge the acknowledgement of a %s alarm", ar.device.name, alarm.TypeName()))
			return
		}
	}
}
