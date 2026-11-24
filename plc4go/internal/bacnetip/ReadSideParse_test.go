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
	"context"
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	plc4go "github.com/apache/plc4x/plc4go/pkg/api"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiTransports "github.com/apache/plc4x/plc4go/pkg/api/transports"
	model "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// readPropertyAckFrame builds the bytes a device sends for a ReadProperty
// ComplexAck: BVLC Original-Unicast-NPDU, an NPDU without addresses, and the
// APDU with the object identifier, property identifier and the value between
// opening and closing tag 3.
func readPropertyAckFrame(invokeId uint8, objectType model.BACnetObjectType, instance uint32, property model.BACnetPropertyIdentifier, value ...byte) []byte {
	objectId := uint32(objectType)<<22 | instance&0x3FFFFF
	apdu := []byte{0x30, invokeId, 0x0C, 0x0C, byte(objectId >> 24), byte(objectId >> 16), byte(objectId >> 8), byte(objectId)}
	if property < 0x100 {
		apdu = append(apdu, 0x19, byte(property))
	} else {
		apdu = append(apdu, 0x1A, byte(property>>8), byte(property))
	}
	apdu = append(apdu, 0x3E)
	apdu = append(apdu, value...)
	apdu = append(apdu, 0x3F)
	frame := []byte{0x81, 0x0A, 0, 0, 0x01, 0x00}
	frame = append(frame, apdu...)
	frame[2], frame[3] = byte(len(frame)>>8), byte(len(frame))
	return frame
}

// parseReadPropertyAckValues parses a device frame like the message codec does and
// returns the constructed data of its ReadProperty ack. A parser panic fails the test
// instead of killing the test binary.
func parseReadPropertyAckValues(t *testing.T, frame []byte) (values model.BACnetConstructedData, err error) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("parser panicked on device bytes % X: %v", frame, r)
		}
	}()
	bvlc, err := model.BVLCParse[model.BVLC](t.Context(), frame)
	if err != nil {
		return nil, err
	}
	apdu := bvlc.(interface{ GetNpdu() model.NPDU }).GetNpdu().GetApdu()
	ack, ok := apdu.(model.APDUComplexAck)
	require.True(t, ok, "expected a ComplexAck, got %T", apdu)
	rp, ok := ack.GetServiceAck().(model.BACnetServiceAckReadProperty)
	require.True(t, ok, "expected a ReadProperty ack, got %T", ack.GetServiceAck())
	return rp.GetValues(), nil
}

// Every *Tagged enum of bacnet-private-enums.mspec reads its value through
// ReadEnumGenericFailing. It panicked ("support for ... not yet implemented") for 26
// of them, among them Network Port properties a device answers when it is read.
func TestReadEnumGenericFailing_EveryPrivateEnumTaggedParses(t *testing.T) {
	cases := map[string]func(ctx context.Context, theBytes []byte) (any, error){
		"BACnetAction":                             wrapTagged(model.BACnetActionTaggedParse),
		"BACnetAccessPassbackMode":                 wrapTagged(model.BACnetAccessPassbackModeTaggedParse),
		"BACnetAccessRuleLocationSpecifier":        wrapTagged(model.BACnetAccessRuleLocationSpecifierTaggedParse),
		"BACnetAccessRuleTimeRangeSpecifier":       wrapTagged(model.BACnetAccessRuleTimeRangeSpecifierTaggedParse),
		"BACnetAccumulatorRecordAccumulatorStatus": wrapTagged(model.BACnetAccumulatorRecordAccumulatorStatusTaggedParse),
		"BACnetAuthenticationFactorType":           wrapTagged(model.BACnetAuthenticationFactorTypeTaggedParse),
		"BACnetAuthenticationStatus":               wrapTagged(model.BACnetAuthenticationStatusTaggedParse),
		"BACnetBackupState":                        wrapTagged(model.BACnetBackupStateTaggedParse),
		"BACnetBinaryPV":                           wrapTagged(model.BACnetBinaryPVTaggedParse),
		"BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriority":          wrapTagged(model.BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTaggedParse),
		"BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisable":      wrapTagged(model.BACnetConfirmedServiceRequestDeviceCommunicationControlEnableDisableTaggedParse),
		"BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilter":    wrapTagged(model.BACnetConfirmedServiceRequestGetEnrollmentSummaryAcknowledgementFilterTaggedParse),
		"BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilter":         wrapTagged(model.BACnetConfirmedServiceRequestGetEnrollmentSummaryEventStateFilterTaggedParse),
		"BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDevice": wrapTagged(model.BACnetConfirmedServiceRequestReinitializeDeviceReinitializedStateOfDeviceTaggedParse),
		"BACnetDoorSecuredStatus":        wrapTagged(model.BACnetDoorSecuredStatusTaggedParse),
		"BACnetDoorValue":                wrapTagged(model.BACnetDoorValueTaggedParse),
		"BACnetFaultType":                wrapTagged(model.BACnetFaultTypeTaggedParse),
		"BACnetFileAccessMethod":         wrapTagged(model.BACnetFileAccessMethodTaggedParse),
		"BACnetIPMode":                   wrapTagged(model.BACnetIPModeTaggedParse),
		"BACnetLiftCarDoorCommand":       wrapTagged(model.BACnetLiftCarDoorCommandTaggedParse),
		"BACnetLiftGroupMode":            wrapTagged(model.BACnetLiftGroupModeTaggedParse),
		"BACnetLightingInProgress":       wrapTagged(model.BACnetLightingInProgressTaggedParse),
		"BACnetLockStatus":               wrapTagged(model.BACnetLockStatusTaggedParse),
		"BACnetNetworkNumberQuality":     wrapTagged(model.BACnetNetworkNumberQualityTaggedParse),
		"BACnetNodeType":                 wrapTagged(model.BACnetNodeTypeTaggedParse),
		"BACnetNotifyType":               wrapTagged(model.BACnetNotifyTypeTaggedParse),
		"BACnetPolarity":                 wrapTagged(model.BACnetPolarityTaggedParse),
		"BACnetProgramRequest":           wrapTagged(model.BACnetProgramRequestTaggedParse),
		"BACnetProgramState":             wrapTagged(model.BACnetProgramStateTaggedParse),
		"BACnetProtocolLevel":            wrapTagged(model.BACnetProtocolLevelTaggedParse),
		"BACnetRouterEntryStatus":        wrapTagged(model.BACnetRouterEntryStatusTaggedParse),
		"BACnetSecurityLevel":            wrapTagged(model.BACnetSecurityLevelTaggedParse),
		"BACnetSecurityPolicy":           wrapTagged(model.BACnetSecurityPolicyTaggedParse),
		"BACnetSegmentation":             wrapTagged(model.BACnetSegmentationTaggedParse),
		"BACnetShedState":                wrapTagged(model.BACnetShedStateTaggedParse),
		"BACnetTimerState":               wrapTagged(model.BACnetTimerStateTaggedParse),
		"BACnetTimerTransition":          wrapTagged(model.BACnetTimerTransitionTaggedParse),
		"BACnetUnconfirmedServiceChoice": wrapTagged(model.BACnetUnconfirmedServiceChoiceTaggedParse),
		"BACnetWriteStatus":              wrapTagged(model.BACnetWriteStatusTaggedParse),
		"BVLCResultCode":                 wrapTagged(model.BVLCResultCodeTaggedParse),
		"NPDUNetworkPriority":            wrapTagged(model.NPDUNetworkPriorityTaggedParse),
		"SecurityResponseCode":           wrapTagged(model.SecurityResponseCodeTaggedParse),
	}
	for name, parse := range cases {
		t.Run(name, func(t *testing.T) {
			// Enumerated 0, application tag: every one of these enums defines 0.
			value, err, panicked := parseRecovering(t.Context(), parse, []byte{0x91, 0x00})
			require.Empty(t, panicked, "parser panicked on a known value")
			require.NoError(t, err)
			assert.EqualValues(t, 0, value)

			// Enumerated 0xFFFFFF: wider than any of these enums. A parse error, not a panic.
			_, err, panicked = parseRecovering(t.Context(), parse, []byte{0x93, 0xFF, 0xFF, 0xFF})
			require.Empty(t, panicked, "parser panicked on an unknown value")
			assert.Error(t, err)
		})
	}
}

// wrapTagged adapts a *TaggedParse function to the test table. GetValue returns a
// different enum type per *Tagged type, so it is read through reflection.
func wrapTagged[T any](parse func(context.Context, []byte, uint8, model.TagClass) (T, error)) func(context.Context, []byte) (any, error) {
	return func(ctx context.Context, theBytes []byte) (any, error) {
		tagged, err := parse(ctx, theBytes, 0, model.TagClass_APPLICATION_TAGS)
		if err != nil {
			return nil, err
		}
		return reflect.ValueOf(tagged).MethodByName("GetValue").Call(nil)[0].Uint(), nil
	}
}

func parseRecovering(ctx context.Context, parse func(context.Context, []byte) (any, error), theBytes []byte) (value any, err error, panicked string) {
	defer func() {
		if r := recover(); r != nil {
			panicked = fmt.Sprint(r)
		}
	}()
	value, err = parse(ctx, theBytes)
	return
}

// A Network Port object (required since protocol revision 17) answers BACnet_IP_Mode
// with an Enumerated. Parsing that reply panicked.
func TestParseReadPropertyAck_NetworkPortBacnetIpMode(t *testing.T) {
	frame := readPropertyAckFrame(1, model.BACnetObjectType_NETWORK_PORT, 1, model.BACnetPropertyIdentifier_BACNET_IP_MODE, 0x91, 0x00)
	values, err := parseReadPropertyAckValues(t, frame)
	require.NoError(t, err)
	_, ok := values.(model.BACnetConstructedDataBACnetIPMode)
	require.True(t, ok, "expected the typed BACnet_IP_Mode case, got %T", values)
	pv := constructedDataToPlcValue(values)
	require.True(t, pv.IsUint32(), "got %T", pv)
	assert.EqualValues(t, 0, pv.GetUint32())
}

// Acked_Transitions and Ack_Required are BACnetEventTransitionBits, a BIT STRING
// (application tag 8). The mspec keyed their typed case on tag 9 (Enumerated), so a
// device's correct reply never reached it and fell through to Unspecified.
func TestParseReadPropertyAck_EventTransitionBitsAreTyped(t *testing.T) {
	// BIT STRING of 3 bits, 5 unused: to-offnormal=1, to-fault=0, to-normal=1.
	bitString := []byte{0x82, 0x05, 0xA0}
	for _, tc := range []struct {
		property model.BACnetPropertyIdentifier
		check    func(model.BACnetConstructedData) (model.BACnetEventTransitionBitsTagged, bool)
	}{
		{model.BACnetPropertyIdentifier_ACKED_TRANSITIONS, func(d model.BACnetConstructedData) (model.BACnetEventTransitionBitsTagged, bool) {
			typed, ok := d.(model.BACnetConstructedDataAckedTransitions)
			if !ok {
				return nil, false
			}
			return typed.GetActualValue(), true
		}},
		{model.BACnetPropertyIdentifier_ACK_REQUIRED, func(d model.BACnetConstructedData) (model.BACnetEventTransitionBitsTagged, bool) {
			typed, ok := d.(model.BACnetConstructedDataAckRequired)
			if !ok {
				return nil, false
			}
			return typed.GetActualValue(), true
		}},
	} {
		t.Run(tc.property.String(), func(t *testing.T) {
			frame := readPropertyAckFrame(1, model.BACnetObjectType_ANALOG_INPUT, 1, tc.property, bitString...)
			values, err := parseReadPropertyAckValues(t, frame)
			require.NoError(t, err)
			bits, ok := tc.check(values)
			require.True(t, ok, "expected the typed case, got %T", values)
			assert.True(t, bits.GetToOffnormal())
			assert.False(t, bits.GetToFault())
			assert.True(t, bits.GetToNormal())
			pv := constructedDataToPlcValue(values)
			assert.Equal(t, []byte{0xA0}, pv.GetRaw())
		})
	}
}

// Interface_Value of a Multi-state Input or Output is BACnetOptionalUnsigned
// (ASHRAE 135 clauses 12.18 and 12.19). The mspec parsed it as BACnetOptionalBinaryPV:
// an Unsigned above 1 failed the whole APDU ("unmapped value"), so the read timed out,
// and 0 or 1 came back as a BinaryPV enum.
func TestParseReadPropertyAck_MultiStateInterfaceValueIsUnsigned(t *testing.T) {
	type optionalUnsigned interface {
		GetActualValue() model.BACnetOptionalUnsigned
	}
	for _, objectType := range []model.BACnetObjectType{model.BACnetObjectType_MULTI_STATE_INPUT, model.BACnetObjectType_MULTI_STATE_OUTPUT} {
		t.Run(objectType.String(), func(t *testing.T) {
			frame := readPropertyAckFrame(1, objectType, 1, model.BACnetPropertyIdentifier_INTERFACE_VALUE, 0x21, 0x03)
			values, err := parseReadPropertyAckValues(t, frame)
			require.NoError(t, err)
			typed, ok := values.(optionalUnsigned)
			require.True(t, ok, "expected the typed case, got %T", values)
			unsigned, ok := typed.GetActualValue().(model.BACnetOptionalUnsignedValue)
			require.True(t, ok, "got %T", typed.GetActualValue())
			assert.EqualValues(t, 3, unsigned.GetUnsignedValue().GetActualValue())

			frame = readPropertyAckFrame(1, objectType, 1, model.BACnetPropertyIdentifier_INTERFACE_VALUE, 0x00)
			values, err = parseReadPropertyAckValues(t, frame)
			require.NoError(t, err)
			typed, ok = values.(optionalUnsigned)
			require.True(t, ok, "expected the typed case, got %T", values)
			_, isNull := typed.GetActualValue().(model.BACnetOptionalUnsignedNull)
			assert.True(t, isNull, "a Null Interface_Value stays Null, got %T", typed.GetActualValue())
		})
	}
}

// rawReplyDevice answers each confirmed ReadProperty with reply(property, invokeId).
type rawReplyDevice struct {
	conn  *net.UDPConn
	wg    sync.WaitGroup
	t     *testing.T
	reply func(property model.BACnetPropertyIdentifier, invokeId uint8) []byte
}

func startRawReplyDevice(t *testing.T, reply func(property model.BACnetPropertyIdentifier, invokeId uint8) []byte) *rawReplyDevice {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	require.NoError(t, err)
	d := &rawReplyDevice{conn: conn, t: t, reply: reply}
	d.wg.Go(d.serve)
	t.Cleanup(func() {
		_ = d.conn.Close()
		d.wg.Wait()
	})
	return d
}

func (d *rawReplyDevice) serve() {
	buf := make([]byte, 1500)
	for {
		n, src, err := d.conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		bvlc, err := model.BVLCParse[model.BVLC](context.Background(), buf[:n])
		if err != nil {
			continue
		}
		request, ok := bvlc.(interface{ GetNpdu() model.NPDU }).GetNpdu().GetApdu().(model.APDUConfirmedRequest)
		if !ok {
			continue
		}
		rp, ok := request.GetServiceRequest().(model.BACnetConfirmedServiceRequestReadProperty)
		if !ok {
			continue
		}
		_, _ = d.conn.WriteToUDP(d.reply(rp.GetPropertyIdentifier().GetValue(), request.GetInvokeId()), src)
	}
}

// The codec recovers a parser panic, but reports it as a fatal transport error, which
// invalidates the connection: one device reply took down every read on it.
func TestNativeBacnetRead_EnumReplyKeepsConnectionUsable(t *testing.T) {
	log := testutils.ProduceTestingLogger(t)
	device := startRawReplyDevice(t, func(property model.BACnetPropertyIdentifier, invokeId uint8) []byte {
		if property == model.BACnetPropertyIdentifier_BACNET_IP_MODE {
			return readPropertyAckFrame(invokeId, model.BACnetObjectType_NETWORK_PORT, 1, property, 0x91, 0x00)
		}
		return readPropertyAckFrame(invokeId, model.BACnetObjectType_ANALOG_INPUT, 1, property, 0x44, 0x41, 0xBC, 0x00, 0x00)
	})

	dm := plc4go.NewPlcDriverManager(options.WithCustomLogger(log))
	dm.RegisterDriver(NewDriver(options.WithCustomLogger(log)))
	apiTransports.RegisterUdpTransport(dm)
	connStr := fmt.Sprintf("bacnet-ip:udp://127.0.0.1:%d?local-port=0&ApduTimeoutMs=2000", device.conn.LocalAddr().(*net.UDPAddr).Port)
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	t.Cleanup(cancel)
	conn, err := dm.GetConnection(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	rr, err := conn.ReadRequestBuilder().AddTagAddress("ipMode", "NETWORK_PORT,1/BACNET_IP_MODE").Build()
	require.NoError(t, err)
	select {
	case <-ctx.Done():
		t.Fatal("BACnet_IP_Mode read timed out")
	case res := <-rr.Execute(ctx):
		require.NoError(t, res.GetErr())
		require.Equal(t, apiModel.PlcResponseCode_OK, res.GetResponse().GetResponseCode("ipMode"))
		assert.EqualValues(t, 0, res.GetResponse().GetValue("ipMode").GetUint32())
	}

	assert.True(t, conn.IsConnected(), "the connection was invalidated by the reply")
	code, val, ok := readPresentValue(t, conn)
	require.True(t, ok, "the next read on the same connection failed")
	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.InDelta(t, 23.5, val, 0.001)
}
