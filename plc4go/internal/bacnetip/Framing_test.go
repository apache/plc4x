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
	bvlc := wrapAPDU(newWhoIsAPDU(t), false)
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
	bvlc := wrapAPDU(newWhoIsAPDU(t), false).(readWriteModel.BVLCOriginalUnicastNPDU)
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
			bvlc := wrapAPDU(newWhoIsAPDU(t), tc.expectingReply).(readWriteModel.BVLCOriginalUnicastNPDU)
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
	bvlc := wrapAPDU(newWhoIsAPDU(t), true).(readWriteModel.BVLCOriginalUnicastNPDU)
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
	bvlc := wrapAPDU(apdu, false).(readWriteModel.BVLCOriginalUnicastNPDU)
	// Same APDU identity should be reachable through the wrapper —
	// MessageCodec.Receive parses BVLC → NPDU → APDU and Reader/Writer
	// match expectations by walking that chain.
	assert.Equal(t, apdu, bvlc.GetNpdu().GetApdu())
}

func TestWrapAPDU_SerializesToValidBVLC(t *testing.T) {
	// End-to-end sanity: the wrapped message round-trips through the
	// model serializer. Catches accidental nil-required-field changes
	// in wrapAPDU that would only show up at runtime under Send().
	bvlc := wrapAPDU(newWhoIsAPDU(t), false)
	raw, err := bvlc.Serialize()
	require.NoError(t, err, "wrapAPDU output must serialize")
	// First byte is BVLC type 0x81; second is function 0x0a (OriginalUnicastNPDU).
	require.GreaterOrEqual(t, len(raw), 4)
	assert.Equal(t, byte(0x81), raw[0], "BVLC magic byte")
	assert.Equal(t, byte(0x0a), raw[1], "BVLC function = OriginalUnicastNPDU")
}
