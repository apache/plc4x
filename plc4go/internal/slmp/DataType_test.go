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

package slmp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func TestDataType_WordsPerElement(t *testing.T) {
	assert.Equal(t, uint16(1), DataTypeWORD.WordsPerElement())
	assert.Equal(t, uint16(1), DataTypeINT.WordsPerElement())
	assert.Equal(t, uint16(1), DataTypeUINT.WordsPerElement())
	assert.Equal(t, uint16(2), DataTypeDINT.WordsPerElement())
	assert.Equal(t, uint16(2), DataTypeUDINT.WordsPerElement())
	assert.Equal(t, uint16(2), DataTypeREAL.WordsPerElement())
}

func TestDataTypeByName(t *testing.T) {
	for _, name := range SupportedDataTypeNames {
		dataType, ok := DataTypeByName(name)
		require.True(t, ok, "%s should be a known type", name)
		assert.Equal(t, name, dataType.String())
		valueType, ok := apiValues.PlcValueTypeByName(name)
		require.True(t, ok)
		assert.Equal(t, valueType, dataType.GetValueType())
	}
	_, ok := DataTypeByName("LREAL")
	assert.False(t, ok, "a type plc4j's SlmpDataType doesn't have must not resolve")
}

// TestDataType_DecodeLittleEndian is the wire contract: SLMP words arrive least significant byte
// first. The D350/D351 payload is the one from the SH-080008 Batch Read worked example, which is
// also the parser-serializer testsuite's read-response case.
func TestDataType_DecodeLittleEndian(t *testing.T) {
	tests := []struct {
		name     string
		dataType DataType
		data     []byte
		quantity uint16
		assert   func(t *testing.T, value apiValues.PlcValue)
	}{
		{
			name: "a WORD is little-endian and unsigned", dataType: DataTypeWORD,
			data: []byte{0xAB, 0x56}, quantity: 1,
			assert: func(t *testing.T, value apiValues.PlcValue) {
				assert.Equal(t, uint16(0x56AB), value.GetUint16())
			},
		},
		{
			// The top bit set must not come back negative: a WORD is a bit pattern.
			name: "a WORD keeps the full unsigned range", dataType: DataTypeWORD,
			data: []byte{0xFF, 0xFF}, quantity: 1,
			assert: func(t *testing.T, value apiValues.PlcValue) {
				assert.Equal(t, uint16(0xFFFF), value.GetUint16())
			},
		},
		{
			name: "a UINT keeps the full unsigned range too", dataType: DataTypeUINT,
			data: []byte{0xFF, 0xFF}, quantity: 1,
			assert: func(t *testing.T, value apiValues.PlcValue) {
				assert.Equal(t, uint16(0xFFFF), value.GetUint16())
			},
		},
		{
			name: "an INT is two's complement", dataType: DataTypeINT,
			data: []byte{0xFF, 0xFF}, quantity: 1,
			assert: func(t *testing.T, value apiValues.PlcValue) {
				assert.Equal(t, int16(-1), value.GetInt16())
			},
		},
		{
			name: "a DINT spans two words, low word first", dataType: DataTypeDINT,
			data: []byte{0xFF, 0xFF, 0xFF, 0xFF}, quantity: 1,
			assert: func(t *testing.T, value apiValues.PlcValue) {
				assert.Equal(t, int32(-1), value.GetInt32())
			},
		},
		{
			name: "a UDINT spans two words too", dataType: DataTypeUDINT,
			data: []byte{0x78, 0x56, 0x34, 0x12}, quantity: 1,
			assert: func(t *testing.T, value apiValues.PlcValue) {
				assert.Equal(t, uint32(0x12345678), value.GetUint32())
			},
		},
		{
			name: "a REAL is IEEE754", dataType: DataTypeREAL,
			data: []byte{0x00, 0x00, 0x80, 0x3F}, quantity: 1,
			assert: func(t *testing.T, value apiValues.PlcValue) {
				assert.InDelta(t, float32(1.0), value.GetFloat32(), 0.0001)
			},
		},
		{
			// The Batch Read worked example: D350=0x56AB, D351=0x170F.
			name: "two words come back as a list", dataType: DataTypeWORD,
			data: []byte{0xAB, 0x56, 0x0F, 0x17}, quantity: 2,
			assert: func(t *testing.T, value apiValues.PlcValue) {
				require.True(t, value.IsList())
				require.Len(t, value.GetList(), 2)
				assert.Equal(t, uint16(0x56AB), value.GetList()[0].GetUint16())
				assert.Equal(t, uint16(0x170F), value.GetList()[1].GetUint16())
			},
		},
		{
			name: "a list of double words strides two words each", dataType: DataTypeDINT,
			data: []byte{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}, quantity: 2,
			assert: func(t *testing.T, value apiValues.PlcValue) {
				require.True(t, value.IsList())
				require.Len(t, value.GetList(), 2)
				assert.Equal(t, int32(1), value.GetList()[0].GetInt32())
				assert.Equal(t, int32(2), value.GetList()[1].GetInt32())
			},
		},
		{
			// A device is free to answer with more than was asked for; the points that were asked
			// for are still right.
			name: "a longer payload than needed is tolerated", dataType: DataTypeWORD,
			data: []byte{0xAB, 0x56, 0xFF, 0xFF}, quantity: 1,
			assert: func(t *testing.T, value apiValues.PlcValue) {
				assert.Equal(t, uint16(0x56AB), value.GetUint16())
			},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			value, err := testCase.dataType.Decode(testCase.data, testCase.quantity)
			require.NoError(t, err)
			require.NotNil(t, value)
			testCase.assert(t, value)
		})
	}
}

func TestDataType_DecodeRejectsShortPayloads(t *testing.T) {
	tests := []struct {
		name     string
		dataType DataType
		data     []byte
		quantity uint16
	}{
		{name: "nothing at all", dataType: DataTypeWORD, data: nil, quantity: 1},
		{name: "half a word", dataType: DataTypeWORD, data: []byte{0x01}, quantity: 1},
		{name: "one word for a double word", dataType: DataTypeDINT, data: []byte{0x01, 0x02}, quantity: 1},
		{name: "one word short of a list", dataType: DataTypeWORD, data: []byte{0x01, 0x02}, quantity: 2},
		{name: "one word short of a REAL list", dataType: DataTypeREAL, data: make([]byte, 6), quantity: 2},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			value, err := testCase.dataType.Decode(testCase.data, testCase.quantity)
			assert.Error(t, err)
			assert.Nil(t, value)
		})
	}
}

func TestDataType_Encode(t *testing.T) {
	tests := []struct {
		name     string
		dataType DataType
		value    apiValues.PlcValue
		quantity uint16
		want     []byte
	}{
		{
			name: "a WORD goes out low byte first", dataType: DataTypeWORD,
			value: spiValues.NewPlcWORD(0x56AB), quantity: 1, want: []byte{0xAB, 0x56},
		},
		{
			// plc4j is careful about exactly this: writing 0xFFFF through a signed short accessor
			// would reject it as negative and corrupt the word.
			name: "a WORD carries the full unsigned range", dataType: DataTypeWORD,
			value: spiValues.NewPlcWORD(0xFFFF), quantity: 1, want: []byte{0xFF, 0xFF},
		},
		{
			name: "a UINT carries the full unsigned range too", dataType: DataTypeUINT,
			value: spiValues.NewPlcUINT(0xFFFF), quantity: 1, want: []byte{0xFF, 0xFF},
		},
		{
			name: "an INT is two's complement", dataType: DataTypeINT,
			value: spiValues.NewPlcINT(-1), quantity: 1, want: []byte{0xFF, 0xFF},
		},
		{
			name: "a DINT spans two words", dataType: DataTypeDINT,
			value: spiValues.NewPlcDINT(-2), quantity: 1, want: []byte{0xFE, 0xFF, 0xFF, 0xFF},
		},
		{
			name: "a UDINT spans two words", dataType: DataTypeUDINT,
			value: spiValues.NewPlcUDINT(0x12345678), quantity: 1, want: []byte{0x78, 0x56, 0x34, 0x12},
		},
		{
			name: "a REAL is IEEE754", dataType: DataTypeREAL,
			value: spiValues.NewPlcREAL(1.0), quantity: 1, want: []byte{0x00, 0x00, 0x80, 0x3F},
		},
		{
			name: "a list fills the words in order", dataType: DataTypeWORD,
			value: spiValues.NewPlcList([]apiValues.PlcValue{
				spiValues.NewPlcWORD(0x56AB), spiValues.NewPlcWORD(0x170F),
			}), quantity: 2, want: []byte{0xAB, 0x56, 0x0F, 0x17},
		},
		{
			name: "a list of double words strides two words each", dataType: DataTypeDINT,
			value: spiValues.NewPlcList([]apiValues.PlcValue{
				spiValues.NewPlcDINT(1), spiValues.NewPlcDINT(2),
			}), quantity: 2, want: []byte{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00},
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			payload, err := testCase.dataType.Encode(testCase.value, testCase.quantity)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, payload)
			// The invariant the write path relies on: the payload is exactly the announced number of
			// points times two bytes, so the frame's header and its payload can't disagree.
			assert.Len(t, payload, int(testCase.quantity)*int(testCase.dataType.WordsPerElement())*2)
		})
	}
}

func TestDataType_EncodeRejects(t *testing.T) {
	tests := []struct {
		name     string
		dataType DataType
		value    apiValues.PlcValue
		quantity uint16
	}{
		{name: "no value at all", dataType: DataTypeWORD, value: nil, quantity: 1},
		{
			name: "a list for a scalar tag", dataType: DataTypeWORD, quantity: 1,
			value: spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcWORD(1)}),
		},
		{
			name: "a scalar for an array tag", dataType: DataTypeWORD, quantity: 2,
			value: spiValues.NewPlcWORD(1),
		},
		{
			name: "too few values for an array tag", dataType: DataTypeWORD, quantity: 3,
			value: spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcWORD(1), spiValues.NewPlcWORD(2)}),
		},
		{
			name: "too many values for an array tag", dataType: DataTypeWORD, quantity: 1,
			value: spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcWORD(1), spiValues.NewPlcWORD(2)}),
		},
		{
			// A negative value has no unsigned word to go into, so it is refused rather than
			// wrapping around into a huge positive one.
			name: "a negative value for a WORD", dataType: DataTypeWORD, quantity: 1,
			value: spiValues.NewPlcINT(-1),
		},
		{
			name: "a value past the signed range for an INT", dataType: DataTypeINT, quantity: 1,
			value: spiValues.NewPlcUDINT(40000),
		},
		{
			name: "a value past the unsigned 16-bit range for a UINT", dataType: DataTypeUINT, quantity: 1,
			value: spiValues.NewPlcUDINT(0x10000),
		},
		{
			name: "a negative value for a UDINT", dataType: DataTypeUDINT, quantity: 1,
			value: spiValues.NewPlcDINT(-1),
		},
		{
			name: "a string for a REAL", dataType: DataTypeREAL, quantity: 1,
			value: spiValues.NewPlcSTRING("hurz"),
		},
		{
			name: "a nil element inside a list", dataType: DataTypeWORD, quantity: 2,
			value: spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcWORD(1), nil}),
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			payload, err := testCase.dataType.Encode(testCase.value, testCase.quantity)
			assert.Error(t, err)
			assert.Nil(t, payload)
		})
	}
}

// TestDataType_EncodeDecodeRoundTrip is the symmetry the driver leans on: what the writer puts on
// the wire is what the reader takes off it.
func TestDataType_EncodeDecodeRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		dataType DataType
		value    apiValues.PlcValue
	}{
		{name: "WORD", dataType: DataTypeWORD, value: spiValues.NewPlcWORD(0xBEEF)},
		{name: "INT", dataType: DataTypeINT, value: spiValues.NewPlcINT(-12345)},
		{name: "UINT", dataType: DataTypeUINT, value: spiValues.NewPlcUINT(54321)},
		{name: "DINT", dataType: DataTypeDINT, value: spiValues.NewPlcDINT(-123456789)},
		{name: "UDINT", dataType: DataTypeUDINT, value: spiValues.NewPlcUDINT(3123456789)},
		{name: "REAL", dataType: DataTypeREAL, value: spiValues.NewPlcREAL(-273.15)},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			payload, err := testCase.dataType.Encode(testCase.value, 1)
			require.NoError(t, err)
			decoded, err := testCase.dataType.Decode(payload, 1)
			require.NoError(t, err)
			assert.Equal(t, testCase.value.GetString(), decoded.GetString())
		})
	}
}
