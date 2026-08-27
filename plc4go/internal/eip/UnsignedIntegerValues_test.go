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
	"encoding/binary"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// The CIP unsigned integer family - USINT, UINT, UDINT and ULINT. Values above the signed
// maximum of the same width are ordinary, so the decoding must not sign extend.

func TestParsePlcValueSingleUSINT(t *testing.T) {
	value, err := parsePlcValue(mustTag(t, "%b:USINT"), []byte{0xFF}, readWriteModel.CIPDataTypeCode_USINT)
	require.NoError(t, err)
	assert.Equal(t, uint8(0xFF), value.GetUint8())
}

func TestParsePlcValueSingleUINT(t *testing.T) {
	raw := binary.LittleEndian.AppendUint16(nil, 0xFFFF)
	value, err := parsePlcValue(mustTag(t, "%w:UINT"), raw, readWriteModel.CIPDataTypeCode_UINT)
	require.NoError(t, err)
	assert.Equal(t, uint16(0xFFFF), value.GetUint16())
}

func TestParsePlcValueSingleUDINT(t *testing.T) {
	raw := binary.LittleEndian.AppendUint32(nil, 0xFFFFFFFF)
	value, err := parsePlcValue(mustTag(t, "%d:UDINT"), raw, readWriteModel.CIPDataTypeCode_UDINT)
	require.NoError(t, err)
	assert.Equal(t, uint32(0xFFFFFFFF), value.GetUint32())
}

// The value a signed decode turns negative.
func TestParsePlcValueUDINTAboveInt32Max(t *testing.T) {
	raw := binary.LittleEndian.AppendUint32(nil, 0x80000000)
	value, err := parsePlcValue(mustTag(t, "%d:UDINT"), raw, readWriteModel.CIPDataTypeCode_UDINT)
	require.NoError(t, err)
	assert.Equal(t, uint32(0x80000000), value.GetUint32())
}

func TestParsePlcValueSingleULINT(t *testing.T) {
	raw := binary.LittleEndian.AppendUint64(nil, 0xFFFFFFFFFFFFFFFF)
	value, err := parsePlcValue(mustTag(t, "%l:ULINT"), raw, readWriteModel.CIPDataTypeCode_ULINT)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), value.GetUint64())
}

func TestParsePlcValueArrayOfUDINTs(t *testing.T) {
	raw := binary.LittleEndian.AppendUint32(nil, 0x00000001)
	raw = binary.LittleEndian.AppendUint32(raw, 0xFFFFFFFF)
	raw = binary.LittleEndian.AppendUint32(raw, 0x80000000)

	value, err := parsePlcValue(mustTag(t, "%d[0..2]:UDINT"), raw, readWriteModel.CIPDataTypeCode_UDINT)
	require.NoError(t, err)
	require.True(t, value.IsList())
	list := value.GetList()
	require.Len(t, list, 3)
	assert.Equal(t, uint32(0x00000001), list[0].GetUint32())
	assert.Equal(t, uint32(0xFFFFFFFF), list[1].GetUint32())
	assert.Equal(t, uint32(0x80000000), list[2].GetUint32())
}

func TestParsePlcValueArrayOfULINTs(t *testing.T) {
	raw := binary.LittleEndian.AppendUint64(nil, 1)
	raw = binary.LittleEndian.AppendUint64(raw, 0xFFFFFFFFFFFFFFFF)

	value, err := parsePlcValue(mustTag(t, "%l[0..1]:ULINT"), raw, readWriteModel.CIPDataTypeCode_ULINT)
	require.NoError(t, err)
	list := value.GetList()
	require.Len(t, list, 2)
	assert.Equal(t, uint64(1), list[0].GetUint64())
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), list[1].GetUint64())
}

func TestParsePlcValueArrayOfUSINTsAndUINTs(t *testing.T) {
	usints, err := parsePlcValue(mustTag(t, "%b[0..2]:USINT"), []byte{0x01, 0x80, 0xFF}, readWriteModel.CIPDataTypeCode_USINT)
	require.NoError(t, err)
	assert.Equal(t, uint8(0xFF), usints.GetList()[2].GetUint8())

	raw := binary.LittleEndian.AppendUint16(nil, 0x8000)
	raw = binary.LittleEndian.AppendUint16(raw, 0xFFFF)
	uints, err := parsePlcValue(mustTag(t, "%w[0..1]:UINT"), raw, readWriteModel.CIPDataTypeCode_UINT)
	require.NoError(t, err)
	assert.Equal(t, uint16(0x8000), uints.GetList()[0].GetUint16())
	assert.Equal(t, uint16(0xFFFF), uints.GetList()[1].GetUint16())
}

func TestParsePlcValueUnsignedIntegerShortReply(t *testing.T) {
	_, err := parsePlcValue(mustTag(t, "%d[0..7]:UDINT"), binary.LittleEndian.AppendUint32(nil, 1), readWriteModel.CIPDataTypeCode_UDINT)
	assert.Error(t, err)

	_, err = parsePlcValue(mustTag(t, "%l[0..3]:ULINT"), binary.LittleEndian.AppendUint64(nil, 1), readWriteModel.CIPDataTypeCode_ULINT)
	assert.Error(t, err)
}

func TestEncodeValueUnsignedIntegers(t *testing.T) {
	for _, test := range []struct {
		name     string
		value    apiValues.PlcValue
		dataType readWriteModel.CIPDataTypeCode
		expected []byte
	}{
		{"USINT", spiValues.NewPlcUSINT(0xFF), readWriteModel.CIPDataTypeCode_USINT, []byte{0xFF}},
		{"UINT", spiValues.NewPlcUINT(0xFFFF), readWriteModel.CIPDataTypeCode_UINT, []byte{0xFF, 0xFF}},
		{"UDINT", spiValues.NewPlcUDINT(0x80000000), readWriteModel.CIPDataTypeCode_UDINT, []byte{0x00, 0x00, 0x00, 0x80}},
		{"ULINT", spiValues.NewPlcULINT(0xFFFFFFFFFFFFFFFF), readWriteModel.CIPDataTypeCode_ULINT,
			[]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	} {
		t.Run(test.name, func(t *testing.T) {
			encoded, err := encodeValue(test.value, test.dataType)
			require.NoError(t, err)
			assert.Equal(t, test.expected, encoded)
			assert.Len(t, encoded, int(test.dataType.Size()))
		})
	}
}

func TestUnsignedIntegerCodesMatchTheCipSpecification(t *testing.T) {
	assert.Equal(t, readWriteModel.CIPDataTypeCode(0x00C6), readWriteModel.CIPDataTypeCode_USINT)
	assert.Equal(t, readWriteModel.CIPDataTypeCode(0x00C7), readWriteModel.CIPDataTypeCode_UINT)
	assert.Equal(t, readWriteModel.CIPDataTypeCode(0x00C8), readWriteModel.CIPDataTypeCode_UDINT)
	assert.Equal(t, readWriteModel.CIPDataTypeCode(0x00C9), readWriteModel.CIPDataTypeCode_ULINT)

	assert.Equal(t, uint8(1), readWriteModel.CIPDataTypeCode_USINT.Size())
	assert.Equal(t, uint8(2), readWriteModel.CIPDataTypeCode_UINT.Size())
	assert.Equal(t, uint8(4), readWriteModel.CIPDataTypeCode_UDINT.Size())
	assert.Equal(t, uint8(8), readWriteModel.CIPDataTypeCode_ULINT.Size())
}

// The signed counterparts must keep decoding to their own types.
func TestSignedCounterpartsAreUnaffected(t *testing.T) {
	value, err := parsePlcValue(mustTag(t, "%d:DINT"), binary.LittleEndian.AppendUint32(nil, 0xFFFFFFFF), readWriteModel.CIPDataTypeCode_DINT)
	require.NoError(t, err)
	assert.Equal(t, int32(-1), value.GetInt32())

	value, err = parsePlcValue(mustTag(t, "%w:INT"), binary.LittleEndian.AppendUint16(nil, 0xFFFF), readWriteModel.CIPDataTypeCode_INT)
	require.NoError(t, err)
	assert.Equal(t, int16(-1), value.GetInt16())
}

func TestTagHandlerAcceptsUnsignedIntegerTypes(t *testing.T) {
	for _, name := range []string{"USINT", "UINT", "UDINT", "ULINT"} {
		t.Run(name, func(t *testing.T) {
			code, found := readWriteModel.CIPDataTypeCodeByName(name)
			require.Truef(t, found, "%s is not resolvable by name", name)
			assert.Equal(t, code, mustTag(t, "%tag:"+name).GetType())
		})
	}
}

// --- string writes ---

// The payload has to match what the read path parses back.
func TestEncodeStringMatchesTheReadPath(t *testing.T) {
	for _, text := range []string{"", "a", "Hello, world", "Grüße"} {
		t.Run(text, func(t *testing.T) {
			encoded, err := encodeValue(spiValues.NewPlcSTRING(text), readWriteModel.CIPDataTypeCode_STRUCTURED)
			require.NoError(t, err)
			require.Len(t, encoded, int(readWriteModel.CIPDataTypeCode_STRUCTURED.Size()))

			assert.Equal(t, uint16(readWriteModel.CIPStructTypeCode_STRING), binary.LittleEndian.Uint16(encoded))
			// The length is the byte count, which differs from the rune count for "Grüße".
			assert.Equal(t, uint32(len(text)), binary.LittleEndian.Uint32(encoded[stringLenOffset:]))

			decoded, err := parsePlcValue(mustTag(t, "%s:STRUCTURED"), encoded, readWriteModel.CIPDataTypeCode_STRUCTURED)
			require.NoError(t, err)
			assert.Equal(t, text, decoded.GetString())
		})
	}
}

func TestEncodeStringLongestFitting(t *testing.T) {
	text := strings.Repeat("x", int(readWriteModel.CIPDataTypeCode_STRUCTURED.Size())-stringDataOffset)

	encoded, err := encodeValue(spiValues.NewPlcSTRING(text), readWriteModel.CIPDataTypeCode_STRUCTURED)
	require.NoError(t, err)
	decoded, err := parsePlcValue(mustTag(t, "%s:STRUCTURED"), encoded, readWriteModel.CIPDataTypeCode_STRUCTURED)
	require.NoError(t, err)
	assert.Equal(t, text, decoded.GetString())
}

func TestEncodeStringTooLongIsRejected(t *testing.T) {
	text := strings.Repeat("x", int(readWriteModel.CIPDataTypeCode_STRUCTURED.Size())-5)

	_, err := encodeValue(spiValues.NewPlcSTRING(text), readWriteModel.CIPDataTypeCode_STRUCTURED)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "83")
}

// CIPDataTypeCode_STRING declares a size of 0, so the serializer emits no payload for it at all.
func TestEncodeStringWithNoRoomIsRejected(t *testing.T) {
	require.Equal(t, uint8(0), readWriteModel.CIPDataTypeCode_STRING.Size())

	_, err := encodeValue(spiValues.NewPlcSTRING("Hello"), readWriteModel.CIPDataTypeCode_STRING)
	assert.Error(t, err)
}
