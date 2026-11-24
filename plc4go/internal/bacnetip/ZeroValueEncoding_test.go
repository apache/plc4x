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
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// An Enumerated, Unsigned or Signed 0 has one data octet 0x00 (BACnet clauses 20.2.4, 20.2.5, 20.2.11). A tag
// whose header says one octet but that carries none makes the device read the next octet (the closing tag) as the
// value and reject the request.
func TestTagsOfZeroCarryOneDataOctet(t *testing.T) {
	cases := []struct {
		name string
		tag  interface{ Serialize() ([]byte, error) }
		wire []byte
	}{
		{"application enumerated 0", readWriteModel.CreateBACnetApplicationTagEnumerated(0), []byte{0x91, 0x00}},
		{"context enumerated 0", readWriteModel.CreateBACnetContextTagEnumerated(2, 0), []byte{0x29, 0x00}},
		{"application unsigned 0", readWriteModel.CreateBACnetApplicationTagUnsignedInteger(0), []byte{0x21, 0x00}},
		{"context unsigned 0", readWriteModel.CreateBACnetContextTagUnsignedInteger(3, 0), []byte{0x39, 0x00}},
		{"application signed 0", readWriteModel.CreateBACnetApplicationTagSignedInteger(0), []byte{0x31, 0x00}},
		{"context signed 0", readWriteModel.CreateBACnetContextTagSignedInteger(1, 0), []byte{0x19, 0x00}},
		{"property identifier 0", readWriteModel.CreateBACnetPropertyIdentifierTagged(1, 0), []byte{0x19, 0x00}},
		{"vendor id 0", readWriteModel.CreateBACnetVendorIdApplicationTagged(0), []byte{0x21, 0x00}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			wire, err := c.tag.Serialize()
			require.NoError(t, err)
			assert.Equal(t, c.wire, wire)
		})
	}
}

// The enumerated payload reads back as 0, whether it is application or context tagged.
func TestEnumeratedZeroReadsBack(t *testing.T) {
	app, err := readWriteModel.BACnetApplicationTagParse[readWriteModel.BACnetApplicationTagEnumerated](t.Context(), []byte{0x91, 0x00})
	require.NoError(t, err)
	assert.Equal(t, uint32(0), app.GetActualValue())
	assert.Equal(t, uint32(1), app.GetHeader().GetActualLength())

	ctx, err := readWriteModel.BACnetContextTagParse[readWriteModel.BACnetContextTagEnumerated](t.Context(), []byte{0x29, 0x00}, 2, readWriteModel.BACnetDataType_ENUMERATED)
	require.NoError(t, err)
	assert.Equal(t, uint32(0), ctx.GetActualValue())
}

// Switching a binary point off: WriteProperty and WritePropertyMultiple of Present_Value := INACTIVE, as
// complete service requests.
func TestBuildServiceRequest_binaryInactiveOnTheWire(t *testing.T) {
	bo1 := makeTag(readWriteModel.BACnetObjectType_BINARY_OUTPUT, 1, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE)
	bv2 := makeTag(readWriteModel.BACnetObjectType_BINARY_VALUE, 2, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE)
	cases := []struct {
		name  string
		specs []writeTagSpec
		wire  []byte
	}{
		{
			name:  "WriteProperty bool false",
			specs: []writeTagSpec{{name: "bo", tag: bo1, value: spiValues.NewPlcBOOL(false)}},
			// service 0f, [0] binary-output,1, [1] present-value, [3] { Enumerated 0 }
			wire: []byte{0x0F, 0x0C, 0x01, 0x00, 0x00, 0x01, 0x19, 0x55, 0x3E, 0x91, 0x00, 0x3F},
		},
		{
			name:  "WriteProperty unsigned 0",
			specs: []writeTagSpec{{name: "bo", tag: bo1, value: spiValues.NewPlcUDINT(0)}},
			wire:  []byte{0x0F, 0x0C, 0x01, 0x00, 0x00, 0x01, 0x19, 0x55, 0x3E, 0x91, 0x00, 0x3F},
		},
		{
			name: "WritePropertyMultiple bool false",
			specs: []writeTagSpec{
				{name: "bo", tag: bo1, value: spiValues.NewPlcBOOL(false)},
				{name: "bv", tag: bv2, value: spiValues.NewPlcBOOL(false)},
			},
			// service 10, then per object: [0] object id, [1] { [0] present-value, [2] { Enumerated 0 } }
			wire: []byte{
				0x10,
				0x0C, 0x01, 0x00, 0x00, 0x01, 0x1E, 0x09, 0x55, 0x2E, 0x91, 0x00, 0x2F, 0x1F,
				0x0C, 0x01, 0x40, 0x00, 0x02, 0x1E, 0x09, 0x55, 0x2E, 0x91, 0x00, 0x2F, 0x1F,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req, err := newTestWriter(t).buildServiceRequest(writeRequestFor(t, c.specs))
			require.NoError(t, err)
			wire, err := req.Serialize()
			require.NoError(t, err)
			assert.Equal(t, c.wire, wire)

			parsed, err := readWriteModel.BACnetConfirmedServiceRequestParse[readWriteModel.BACnetConfirmedServiceRequest](t.Context(), wire, uint32(len(wire)))
			require.NoError(t, err, "the request must parse back")
			reserialized, err := parsed.Serialize()
			require.NoError(t, err)
			assert.Equal(t, wire, reserialized)
		})
	}
}

// A tagged enum's header length and its data octets must agree for every length, not just for 0: a proprietary
// property identifier or vendor id of 3 octets went out with a length of 3 and 4 data octets.
func TestTaggedEnumLengthMatchesItsOctets(t *testing.T) {
	cases := []struct {
		name string
		tag  interface{ Serialize() ([]byte, error) }
		wire []byte
	}{
		{"property identifier 0xFF", readWriteModel.CreateBACnetPropertyIdentifierTagged(1, 0xFF), []byte{0x19, 0xFF}},
		{"property identifier 0x100", readWriteModel.CreateBACnetPropertyIdentifierTagged(1, 0x100), []byte{0x1A, 0x01, 0x00}},
		{"property identifier 0x10000", readWriteModel.CreateBACnetPropertyIdentifierTagged(1, 0x10000), []byte{0x1B, 0x01, 0x00, 0x00}},
		{"property identifier max 22 bits", readWriteModel.CreateBACnetPropertyIdentifierTagged(1, 0x3FFFFF), []byte{0x1B, 0x3F, 0xFF, 0xFF}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			wire, err := c.tag.Serialize()
			require.NoError(t, err)
			assert.Equal(t, c.wire, wire)
		})
	}
}
