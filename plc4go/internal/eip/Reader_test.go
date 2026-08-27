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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
)

// TestSliceServiceHappyPath exercises a valid 2-entry MultipleServiceResponse
// offset table, so the reversed-table test below isn't the only case covered.
func TestSliceServiceHappyPath(t *testing.T) {
	service1, err := readWriteModel.NewCipReadResponse(0, nil,
		readWriteModel.NewCIPData(readWriteModel.CIPDataTypeCode_DINT, []byte{0x01, 0x00, 0x00, 0x00})).Serialize()
	require.NoError(t, err)
	service2, err := readWriteModel.NewCipReadResponse(0, nil,
		readWriteModel.NewCIPData(readWriteModel.CIPDataTypeCode_DINT, []byte{0x02, 0x00, 0x00, 0x00})).Serialize()
	require.NoError(t, err)

	servicesData := append(append([]byte{}, service1...), service2...)
	offsets := []uint16{0, uint16(len(service1))}

	first, ok := sliceService(servicesData, offsets, 0, 2)
	require.True(t, ok)
	firstResponse, isReadResponse := first.(readWriteModel.CipReadResponse)
	require.True(t, isReadResponse)
	assert.Equal(t, []byte{0x01, 0x00, 0x00, 0x00}, firstResponse.GetData().GetData())

	second, ok := sliceService(servicesData, offsets, 1, 2)
	require.True(t, ok)
	secondResponse, isReadResponse := second.(readWriteModel.CipReadResponse)
	require.True(t, isReadResponse)
	assert.Equal(t, []byte{0x02, 0x00, 0x00, 0x00}, secondResponse.GetData().GetData())
}

// TestSliceServiceReversedOffsetsIsRejected guards against a lying/reversed
// offset table wrapping around in unsigned 16-bit arithmetic before the
// bounds guards run. offsets[0] sits near 0xFFFF while offsets[1] sits near
// 0, so a naive `uint16(offsets[i] - offsets[0])` subtraction would wrap to a
// small positive number instead of tripping the "offset < 0" guard.
func TestSliceServiceReversedOffsetsIsRejected(t *testing.T) {
	servicesData := make([]byte, 8)
	offsets := []uint16{0xFFF0, 2}

	service, ok := sliceService(servicesData, offsets, 1, 2)
	assert.False(t, ok)
	assert.Nil(t, service)
}

// TestToAnsi pins toAnsi's byte-for-byte output against the CIP ANSI encoding
// rules used by the Java driver (EipTcpConnection.java): every ANSI extended
// symbol segment (0x91 tag, produced by a 3-bit DataSegment discriminator
// (0x04) plus a 5-bit AnsiExtendedSymbolSegment discriminator (0x11) packed
// into one byte) is followed by a length byte, the ASCII symbol bytes, and -
// whenever the symbol has an odd length, dotted qualifier or not - a single
// 0x00 pad byte so the segment stays an even number of bytes. An array index
// (`[n]`) instead produces a 2-byte LogicalSegment/MemberID pair: 0x28 (3-bit
// PathSegment discriminator 0x01 + 3-bit LogicalSegmentType discriminator
// 0x02 + 2-bit format 0) followed by the raw index byte - MemberID has no pad
// byte of its own.
func TestToAnsi(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected []byte
	}{
		{
			// simple, even-length symbol: no pad byte.
			name:     "simple even",
			tag:      "rate",
			expected: []byte{0x91, 0x04, 0x72, 0x61, 0x74, 0x65},
		},
		{
			// simple, odd-length symbol: padded to keep the segment even.
			name:     "simple odd",
			tag:      "count",
			expected: []byte{0x91, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x00},
		},
		{
			// dotted member with an even-length name: no pad on either segment.
			name: "dotted even member",
			tag:  "tag.rate",
			expected: []byte{
				0x91, 0x03, 0x74, 0x61, 0x67, 0x00, // "tag" (odd) - padded
				0x91, 0x04, 0x72, 0x61, 0x74, 0x65, // "rate" (even) - no pad
			},
		},
		{
			// dotted member with an ODD-length name - this is the regression
			// this test pins: prior to the fix, toAnsi passed a nil pad for
			// the dotted branch regardless of the identifier's length, so
			// "count" here would have serialized without its trailing 0x00,
			// leaving the frame internally inconsistent with requestPathSize.
			name: "dotted odd member",
			tag:  "abcd.count",
			expected: []byte{
				0x91, 0x04, 0x61, 0x62, 0x63, 0x64, // "abcd" (even) - no pad
				0x91, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x00, // "count" (odd) - padded
			},
		},
		{
			// array index: LogicalSegment(MemberID) - no pad byte here, ever.
			name: "array index",
			tag:  "arr[0]",
			expected: []byte{
				0x91, 0x03, 0x61, 0x72, 0x72, 0x00, // "arr" (odd) - padded
				0x28, 0x00, // MemberID(format=0, instance=0)
			},
		},
		{
			// mixed: unqualified, dotted, indexed, dotted again - all with
			// odd-length single-character/three-character names, so every
			// ANSI segment here must carry its pad byte.
			name: "mixed",
			tag:  "a.b[2].cde",
			expected: []byte{
				0x91, 0x01, 0x61, 0x00, // "a" (odd) - padded
				0x91, 0x01, 0x62, 0x00, // "b" (odd, dotted) - padded
				0x28, 0x02, // MemberID(format=0, instance=2)
				0x91, 0x03, 0x63, 0x64, 0x65, 0x00, // "cde" (odd, dotted) - padded
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := toAnsi(tt.tag)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

// An array index rides in a MemberID, whose instance field the mspec declares as uint 8, so it
// holds 0 to 255. A larger index used to be converted with uint8(...), which wraps silently -
// a[300] addressed element 44 and produced a request that looked entirely valid. The Java
// driver rejects the same address when it serializes ("Value 300 is out of range for 8 bits"),
// so this reports it too rather than reading the wrong element.
func TestToAnsiRejectsOutOfRangeIndex(t *testing.T) {
	for _, tag := range []string{"a[256]", "a[300]", "myArray[1000]"} {
		t.Run(tag, func(t *testing.T) {
			_, err := toAnsi(tag)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "255")
		})
	}
}

func TestToAnsiAcceptsTheLargestIndex(t *testing.T) {
	actual, err := toAnsi("a[255]")
	require.NoError(t, err)
	assert.Equal(t, []byte{0x91, 0x01, 0x61, 0x00, 0x28, 0xFF}, actual)
}
