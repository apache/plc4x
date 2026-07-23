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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func newWhoIsAPDU(t *testing.T) readWriteModel.APDU {
	t.Helper()
	whoIs := readWriteModel.NewBACnetUnconfirmedServiceRequestWhoIs(nil, nil)
	return readWriteModel.NewAPDUUnconfirmedRequest(whoIs)
}

func TestWrapAPDU_ProducesBVLCOriginalUnicastNPDU(t *testing.T) {
	bvlc := wrapAPDU(newWhoIsAPDU(t), false, nil)
	require.NotNil(t, bvlc)
	// MessageCodec.Send type-asserts to BVLCOriginalUnicastNPDU on the
	// send-path; wrapAPDU must produce exactly that type.
	_, ok := bvlc.(readWriteModel.BVLCOriginalUnicastNPDU)
	assert.True(t, ok, "wrapAPDU must return a BVLCOriginalUnicastNPDU, got %T", bvlc)
}

func TestWrapAPDU_NPDUProtocolVersionIs1(t *testing.T) {
	// BACnet stacks reject NPDUs with a wrong protocol version. Pin it
	// to 1 (the only spec-valid value) so an accidental change shows up
	// as a test failure rather than a wire-protocol incompatibility.
	bvlc := wrapAPDU(newWhoIsAPDU(t), false, nil).(readWriteModel.BVLCOriginalUnicastNPDU)
	assert.Equal(t, uint8(1), bvlc.GetNpdu().GetProtocolVersionNumber())
}

func TestWrapAPDU_ExpectingReplyPropagatesToControl(t *testing.T) {
	cases := []struct {
		name           string
		expectingReply bool
	}{
		{"confirmed-request-sets-flag", true},
		{"unconfirmed-broadcast-clears-flag", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bvlc := wrapAPDU(newWhoIsAPDU(t), tc.expectingReply, nil).(readWriteModel.BVLCOriginalUnicastNPDU)
			control := bvlc.GetNpdu().GetControl()
			assert.Equal(t, tc.expectingReply, control.GetExpectingReply(),
				"NPDU control.expectingReply must reflect the wrapAPDU argument")
		})
	}
}

func TestWrapAPDU_LocalScope_NoRouting(t *testing.T) {
	// We only support local (same-network) addressing in Phase 6's tag layer.
	// The NPDU control fields for routing must all be off so bacpypes3/Niagara
	// don't interpret the message as routed.
	bvlc := wrapAPDU(newWhoIsAPDU(t), true, nil).(readWriteModel.BVLCOriginalUnicastNPDU)
	control := bvlc.GetNpdu().GetControl()
	assert.False(t, control.GetMessageTypeFieldPresent())
	assert.False(t, control.GetDestinationSpecified())
	assert.False(t, control.GetSourceSpecified())
	// Routing fields (DNET/DLEN/DADR + SNET/SLEN/SADR + HopCount + NLM) must
	// be nil — otherwise a peer treats it as a routed frame.
	npdu := bvlc.GetNpdu()
	assert.Nil(t, npdu.GetDestinationNetworkAddress())
	assert.Nil(t, npdu.GetDestinationLength())
	assert.Nil(t, npdu.GetDestinationAddress())
	assert.Nil(t, npdu.GetSourceNetworkAddress())
	assert.Nil(t, npdu.GetSourceLength())
	assert.Nil(t, npdu.GetSourceAddress())
	assert.Nil(t, npdu.GetHopCount())
	assert.Nil(t, npdu.GetNlm())
}

func TestWrapAPDU_PreservesAPDU(t *testing.T) {
	apdu := newWhoIsAPDU(t)
	bvlc := wrapAPDU(apdu, false, nil).(readWriteModel.BVLCOriginalUnicastNPDU)
	// Same APDU identity should be reachable through the wrapper —
	// MessageCodec.Receive parses BVLC → NPDU → APDU and Reader/Writer
	// match expectations by walking that chain.
	assert.Equal(t, apdu, bvlc.GetNpdu().GetApdu())
}

func TestWrapAPDU_SerializesToValidBVLC(t *testing.T) {
	// End-to-end sanity: the wrapped message round-trips through the
	// model serializer. Catches accidental nil-required-field changes
	// in wrapAPDU that would only show up at runtime under Send().
	bvlc := wrapAPDU(newWhoIsAPDU(t), false, nil)
	raw, err := bvlc.Serialize()
	require.NoError(t, err, "wrapAPDU output must serialize")
	// First byte is BVLC type 0x81; second is function 0x0a (OriginalUnicastNPDU).
	require.GreaterOrEqual(t, len(raw), 4)
	assert.Equal(t, byte(0x81), raw[0], "BVLC magic byte")
	assert.Equal(t, byte(0x0a), raw[1], "BVLC function = OriginalUnicastNPDU")
}

// TestWrapAPDU_RoutedDestination pins the routed framing (ASHRAE 135 clause
// 6): a connection whose target sits behind a BACnet router emits a
// destination specifier (DNET/DLEN/DADR) with a fresh hop count, and the
// frame survives a serialize/parse round trip.
func TestWrapAPDU_RoutedDestination(t *testing.T) {
	dest := &routedDestination{dnet: 3001, dadr: []uint8{192, 168, 102, 20, 0xBA, 0xC0}}
	bvlc := wrapAPDU(newWhoIsAPDU(t), false, dest).(readWriteModel.BVLCOriginalUnicastNPDU)
	npdu := bvlc.GetNpdu()

	assert.True(t, npdu.GetControl().GetDestinationSpecified(), "control.destinationSpecified")
	assert.False(t, npdu.GetControl().GetSourceSpecified(), "source must stay absent on originated frames")
	require.NotNil(t, npdu.GetDestinationNetworkAddress())
	assert.Equal(t, uint16(3001), *npdu.GetDestinationNetworkAddress())
	require.NotNil(t, npdu.GetDestinationLength())
	assert.Equal(t, uint8(6), *npdu.GetDestinationLength())
	assert.Equal(t, []uint8{192, 168, 102, 20, 0xBA, 0xC0}, npdu.GetDestinationAddress())
	require.NotNil(t, npdu.GetHopCount())
	assert.Equal(t, routedDestinationHopCount, *npdu.GetHopCount())

	// Round trip: the routed header must reparse byte-identically.
	data, err := bvlc.Serialize()
	require.NoError(t, err)
	reparsed, err := readWriteModel.BVLCParse[readWriteModel.BVLC](context.Background(), data)
	require.NoError(t, err)
	renpdu := reparsed.(readWriteModel.BVLCOriginalUnicastNPDU).GetNpdu()
	require.NotNil(t, renpdu.GetDestinationNetworkAddress())
	assert.Equal(t, uint16(3001), *renpdu.GetDestinationNetworkAddress())
	assert.Equal(t, []uint8{192, 168, 102, 20, 0xBA, 0xC0}, renpdu.GetDestinationAddress())
}

// TestParseRemoteAddress covers the two DADR syntaxes and rejects garbage.
func TestParseRemoteAddress(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []uint8
		wantErr bool
	}{
		{"bacnet-ip", "192.168.102.20:47808", []uint8{192, 168, 102, 20, 0xBA, 0xC0}, false},
		{"hex mac", "0x0C", []uint8{0x0C}, false},
		{"hex multi-octet", "0x00fF", []uint8{0x00, 0xFF}, false},
		{"ipv6 rejected", "[::1]:47808", nil, true},
		{"empty hex", "0x", nil, true},
		{"garbage", "not-an-address", nil, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseRemoteAddress(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestRoutedDestinationFromConfiguration pins the option pairing rules.
func TestRoutedDestinationFromConfiguration(t *testing.T) {
	dest, err := routedDestinationFromConfiguration(Configuration{})
	require.NoError(t, err)
	assert.Nil(t, dest, "no options -> local addressing")

	dest, err = routedDestinationFromConfiguration(Configuration{RemoteNetwork: 3001, RemoteAddress: "192.168.102.20:47808"})
	require.NoError(t, err)
	require.NotNil(t, dest)
	assert.Equal(t, uint16(3001), dest.dnet)
	assert.Len(t, dest.dadr, 6)

	_, err = routedDestinationFromConfiguration(Configuration{RemoteNetwork: 3001})
	assert.Error(t, err, "network without address must fail")
	_, err = routedDestinationFromConfiguration(Configuration{RemoteAddress: "192.168.102.20:47808"})
	assert.Error(t, err, "address without network must fail")
}

// TestRoutedOriginOptions pins the discovery-side mirror of the routed
// options: a routed I-Am's SNET/SADR comes back as the same
// RemoteNetwork/RemoteAddress keys the connection URL accepts.
func TestRoutedOriginOptions(t *testing.T) {
	mkNPDU := func(src bool, snet uint16, sadr []uint8) readWriteModel.NPDU {
		control := readWriteModel.NewNPDUControl(false, false, src, false, readWriteModel.NPDUNetworkPriority_NORMAL_MESSAGE)
		var snetP *uint16
		var slenP *uint8
		if src {
			snetP = &snet
			slen := uint8(len(sadr))
			slenP = &slen
		}
		return readWriteModel.NewNPDU(1, control, nil, nil, nil, snetP, slenP, sadr, nil, nil, newWhoIsAPDU(t))
	}

	assert.Nil(t, routedOriginOptions(mkNPDU(false, 0, nil)), "local frame -> no options")

	opts := routedOriginOptions(mkNPDU(true, 3001, []uint8{192, 168, 102, 20, 0xBA, 0xC0}))
	require.NotNil(t, opts)
	assert.Equal(t, []string{"3001"}, opts["RemoteNetwork"])
	assert.Equal(t, []string{"192.168.102.20:47808"}, opts["RemoteAddress"])

	opts = routedOriginOptions(mkNPDU(true, 5, []uint8{0x0C}))
	require.NotNil(t, opts)
	assert.Equal(t, []string{"0x0c"}, opts["RemoteAddress"], "non-B/IP MAC renders as hex")

	// Round trip: the discovery options must parse back into the same DADR.
	dadr, err := parseRemoteAddress(opts["RemoteAddress"][0])
	require.NoError(t, err)
	assert.Equal(t, []uint8{0x0C}, dadr)
}
