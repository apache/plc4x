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
	"bytes"
	"context"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// An Octet String's header length is the number of content octets (ASHRAE 135 Clause 20.2.8). It used to be one
// more, so a device read the octet after the value (in a WriteProperty, the closing tag 3f) as part of it.
func TestPlcValueToApplicationTag_octetStringOnTheWire(t *testing.T) {
	cases := []struct {
		name  string
		value []byte
		want  string
	}{
		{"empty", []byte{}, "60"},
		{"one octet", []byte{0xAB}, "61ab"},
		{"IPv4 address", []byte{0xC0, 0xA8, 0x01, 0x6F}, "64c0a8016f"},
		{"MAC, extended length", []byte{0xC0, 0xA8, 0x01, 0x6F, 0xBA, 0xC0}, "6506c0a8016fbac0"},
		{"254 octets, 16-bit extended length", bytes.Repeat([]byte{0x11}, 254), "65fe00fe" + hex.EncodeToString(bytes.Repeat([]byte{0x11}, 254))},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tag, err := plcValueToApplicationTag(spiValues.NewPlcRawByteArray(c.value), hintNone)
			require.NoError(t, err)
			wire, err := tag.Serialize()
			require.NoError(t, err)
			assert.Equal(t, c.want, hex.EncodeToString(wire))

			parsed, err := readWriteModel.BACnetApplicationTagParse[readWriteModel.BACnetApplicationTag](context.Background(), wire)
			require.NoError(t, err)
			octets, ok := parsed.(readWriteModel.BACnetApplicationTagOctetString)
			require.True(t, ok, "parsed %T", parsed)
			assert.Equal(t, c.value, octets.GetPayload().GetOctets())
		})
	}
}

func TestCreateBACnetContextTagOctetString_lengthIsTheNumberOfOctets(t *testing.T) {
	wire, err := readWriteModel.CreateBACnetContextTagOctetString(2, []byte{0xAB, 0xCD}).Serialize()
	require.NoError(t, err)
	assert.Equal(t, "2aabcd", hex.EncodeToString(wire))
}

// The value of a WriteProperty to an Octet String Value's Present_Value ends right before the closing tag.
func TestBuildServiceRequest_octetStringWritePropertyOnTheWire(t *testing.T) {
	writer := newTestWriter(t)
	req := writeRequestFor(t, []writeTagSpec{
		{
			name:  "osv",
			tag:   makeTag(readWriteModel.BACnetObjectType_OCTETSTRING_VALUE, 1, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE),
			value: spiValues.NewPlcRawByteArray([]byte{0xC0, 0xA8, 0x01, 0x6F}),
		},
	})
	got, err := writer.buildServiceRequest(req)
	require.NoError(t, err)
	wire, err := got.Serialize()
	require.NoError(t, err)
	assert.Contains(t, hex.EncodeToString(wire), "3e64c0a8016f3f")
}
