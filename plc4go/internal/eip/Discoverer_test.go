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

package eip

import (
	"bytes"
	"context"
	"encoding/binary"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func TestListIdentityRequestRoundTrip(t *testing.T) {
	request := buildListIdentityRequest()
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	require.NoError(t, request.SerializeWithWriteBuffer(context.Background(), wb))
	raw := wb.GetBytes()
	// 24-byte encapsulation header, ListIdentity command 0x63, empty payload
	require.Len(t, raw, 24)
	assert.Equal(t, uint16(0x63), binary.LittleEndian.Uint16(raw[0:2]))
	assert.Equal(t, uint16(0), binary.LittleEndian.Uint16(raw[2:4]))

	rb := utils.NewReadBufferByteBased(raw, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	parsed, err := readWriteModel.EipPacketParseWithBuffer[readWriteModel.EipPacket](context.Background(), rb, false)
	require.NoError(t, err)
	_, ok := parsed.(readWriteModel.EipListIdentityRequest)
	assert.True(t, ok)
}

// TestCipIdentityProductNameEncoding pins the CipIdentity.productName field to actually
// serialize its ASCII bytes onto the wire. The generated Serialize used to write the
// productName string without a UTF8 encoding hint, which caused the byte-based writer to
// zero-fill the field instead of emitting the string's bytes.
func TestCipIdentityProductNameEncoding(t *testing.T) {
	const productName = "1756-EN2T"
	identity := readWriteModel.NewCipIdentity(
		uint16(0),           // encapsulationProtocolVersion
		uint16(0),           // socketAddressFamily
		uint16(0),           // socketAddressPort
		[]uint8{0, 0, 0, 0}, // socketAddressAddress
		uint16(0),           // vendorId
		uint16(0),           // deviceType
		uint16(0),           // productCode
		uint8(0),            // revisionMajor
		uint8(0),            // revisionMinor
		uint16(0),           // status
		uint32(0),           // serialNumber
		productName,         // productName
		uint8(0),            // state
	)

	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	require.NoError(t, identity.SerializeWithWriteBuffer(context.Background(), wb))
	raw := wb.GetBytes()

	assert.True(t, bytes.Contains(raw, []byte(productName)), "expected serialized CipIdentity to contain product name bytes %q, got % x", productName, raw)
}

// TestSubnetBroadcast pins the broadcast-address computation used to target
// the per-interface ListIdentity request (subnetBroadcast ORs in the inverted
// mask, i.e. sets all host bits).
func TestSubnetBroadcast(t *testing.T) {
	tests := []struct {
		name string
		cidr string
		want net.IP
	}{
		{
			name: "/24",
			cidr: "192.168.1.10/24",
			want: net.IPv4(192, 168, 1, 255).To4(),
		},
		{
			name: "/16",
			cidr: "10.20.30.40/16",
			want: net.IPv4(10, 20, 255, 255).To4(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, ipNet, err := net.ParseCIDR(tt.cidr)
			require.NoError(t, err)
			// ParseCIDR zeroes the host bits in ipNet.IP; use the original
			// host address the same way a real interface address would be
			// reported by net.Interface.Addrs().
			ipNet.IP = ip
			got := subnetBroadcast(ipNet)
			assert.True(t, tt.want.Equal(got), "expected %v, got %v", tt.want, got)
		})
	}
}

// TestIdentityFromPacket hand-serializes an EipListIdentityResponse carrying
// a CipIdentity item, round-trips it through the wire format, and asserts
// identityFromPacket (the helper the discoverer uses to decode ListIdentity
// replies) extracts the product name.
func TestIdentityFromPacket(t *testing.T) {
	const productName = "1756-EN2T"
	identity := readWriteModel.NewCipIdentity(
		uint16(1),                          // encapsulationProtocolVersion
		uint16(2),                          // socketAddressFamily
		uint16(EipUdpDiscoveryDefaultPort), // socketAddressPort
		[]uint8{192, 168, 1, 50},           // socketAddressAddress
		uint16(1),                          // vendorId
		uint16(14),                         // deviceType
		uint16(54),                         // productCode
		uint8(1),                           // revisionMajor
		uint8(1),                           // revisionMinor
		uint16(0),                          // status
		uint32(123456),                     // serialNumber
		productName,                        // productName
		uint8(3),                           // state
	)
	response := readWriteModel.NewEipListIdentityResponse(
		0, 0, []byte{0, 0, 0, 0, 0, 0, 0, 0}, 0,
		[]readWriteModel.CommandSpecificDataItem{identity},
	)

	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.LittleEndian))
	require.NoError(t, response.SerializeWithWriteBuffer(context.Background(), wb))
	raw := wb.GetBytes()

	rb := utils.NewReadBufferByteBased(raw, utils.WithByteOrderForReadBufferByteBased(binary.LittleEndian))
	packet, err := readWriteModel.EipPacketParseWithBuffer[readWriteModel.EipPacket](context.Background(), rb, true)
	require.NoError(t, err)

	gotIdentity, ok := identityFromPacket(packet)
	require.True(t, ok)
	assert.Equal(t, productName, gotIdentity.GetProductName())
}
