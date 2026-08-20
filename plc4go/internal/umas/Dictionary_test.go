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

package umas

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
)

func TestParseSymbolTable(t *testing.T) {
	payload := symbolTablePayload(
		symbolRecord{name: "g_r32", dataType: typeIdReal, block: 0x0002, offset: 0x00001234, flags: 0x01, unknown4: 0x02},
		symbolRecord{name: "g_string", dataType: typeIdString, block: 0x0002, offset: 0x00001238},
		symbolRecord{name: "g_plant", dataType: typeIdCustom, block: 0x0003, offset: 0x00000000},
	)

	records, err := parseSymbolTable(t.Context(), payload)
	require.NoError(t, err)
	require.Len(t, records, 3)

	// The name has to come back whole. It is a NUL terminated field with no length prefix, so a
	// parser which reads a fixed number of bytes truncates it and then desynchronizes on the record
	// after it. The generated parser used to do exactly that, because the code generator dropped the
	// unary minus of STATIC_CALL("parseTerminatedString", readBuffer, -1) and asked for a one byte
	// fixed width field; three whole names arriving here is what says the generated call now reads
	// -(1), the helper's variable length mode.
	assert.Equal(t, "g_r32", records[0].GetValue())
	assert.Equal(t, typeIdReal, records[0].GetDataType())
	assert.Equal(t, uint16(0x0002), records[0].GetBlock())
	assert.Equal(t, uint32(0x00001234), records[0].GetOffset())
	assert.Equal(t, uint8(0x01), records[0].GetFlags())
	assert.Equal(t, uint8(0x02), records[0].GetUnknown4())

	assert.Equal(t, "g_string", records[1].GetValue())
	assert.Equal(t, uint32(0x00001238), records[1].GetOffset())

	assert.Equal(t, "g_plant", records[2].GetValue())
	assert.Equal(t, typeIdCustom, records[2].GetDataType())
}

func TestParseSymbolTable_RejectsBrokenPayloads(t *testing.T) {
	t.Run("a payload too short for the header", func(t *testing.T) {
		_, err := parseSymbolTable(t.Context(), []byte{0x00, 0x00})
		assert.Error(t, err)
	})

	t.Run("a record count larger than the payload", func(t *testing.T) {
		payload := symbolTablePayload(symbolRecord{name: "g_r32", dataType: typeIdReal})
		// Claim five records where one is present.
		payload[5] = 0x05
		_, err := parseSymbolTable(t.Context(), payload)
		assert.Error(t, err)
	})

	t.Run("a name without a terminator", func(t *testing.T) {
		payload := symbolTablePayload(symbolRecord{name: "g_r32", dataType: typeIdReal})
		// Drop the terminator, leaving a name which runs off the end of the payload. Reading it as a
		// string anyway would put every following record at the wrong offset.
		_, err := parseSymbolTable(t.Context(), payload[:len(payload)-1])
		assert.Error(t, err)
	})

	t.Run("an empty table is not an error", func(t *testing.T) {
		records, err := parseSymbolTable(t.Context(), symbolTablePayload())
		require.NoError(t, err)
		assert.Empty(t, records)
	})
}

func TestParseDatatypeNames(t *testing.T) {
	payload := datatypeDictionaryPayload(
		datatypeRecord{name: "MY_STRING", dataSize: 20, classIdentifier: 0, dataType: uint8(typeIdString)},
		datatypeRecord{name: "MY_STRUCT", dataSize: 12, classIdentifier: 2, dataType: 0},
		datatypeRecord{name: "MY_ARRAY", dataSize: 40, classIdentifier: 4, dataType: 0},
	)

	records, err := parseDatatypeNames(t.Context(), payload)
	require.NoError(t, err)
	require.Len(t, records, 3)

	assert.Equal(t, "MY_STRING", records[0].GetValue())
	assert.Equal(t, uint16(20), records[0].GetDataSize())
	assert.Equal(t, uint8(0), records[0].GetClassIdentifier())
	assert.Equal(t, uint8(typeIdString), records[0].GetDataType())

	assert.Equal(t, "MY_STRUCT", records[1].GetValue())
	assert.Equal(t, uint8(2), records[1].GetClassIdentifier())

	assert.Equal(t, "MY_ARRAY", records[2].GetValue())
	assert.Equal(t, uint16(40), records[2].GetDataSize())
}

func TestParseDatatypeNames_RejectsBrokenPayloads(t *testing.T) {
	_, err := parseDatatypeNames(t.Context(), []byte{0x00})
	assert.Error(t, err)

	payload := datatypeDictionaryPayload(datatypeRecord{name: "MY_STRING", dataSize: 20})
	_, err = parseDatatypeNames(t.Context(), payload[:len(payload)-1])
	assert.Error(t, err)
}

func TestParseUdtDefinition(t *testing.T) {
	payload := udtDefinitionPayload(
		udtMember{name: "meta", dataType: typeIdDint, offset: 0},
		udtMember{name: "r32", dataType: typeIdReal, offset: 4},
	)

	members, err := parseUdtDefinition(t.Context(), payload)
	require.NoError(t, err)
	require.Len(t, members, 2)
	assert.Equal(t, "meta", members[0].GetValue())
	assert.Equal(t, typeIdDint, members[0].GetDataType())
	assert.Equal(t, uint16(0), members[0].GetOffset())
	assert.Equal(t, "r32", members[1].GetValue())
	assert.Equal(t, typeIdReal, members[1].GetDataType())
	assert.Equal(t, uint16(4), members[1].GetOffset())
}

func TestParseUdtDefinition_RejectsBrokenPayloads(t *testing.T) {
	_, err := parseUdtDefinition(t.Context(), []byte{0x01, 0x00})
	assert.Error(t, err)

	payload := udtDefinitionPayload(udtMember{name: "meta", dataType: typeIdDint})
	_, err = parseUdtDefinition(t.Context(), payload[:len(payload)-1])
	assert.Error(t, err)
}

// The array type definition is the fourth dictionary payload, and the only one whose mspec type pins
// LITTLE_ENDIAN on every field. That is why it alone is safe to hand to the convenience parser: the
// generator emits the byte order into the parser itself, so UmasArrayTypeDefinitionParse builds a
// little endian buffer rather than the default big endian one the other three would get. It also has
// no NUL terminated string in it, so the generator bug never touched it.
func TestArrayTypeDefinitionParsesThroughTheModel(t *testing.T) {
	payload := arrayTypeDefinitionPayload(typeIdDint,
		arrayDimension{startIndex: 0, upperBound: 9},
		arrayDimension{startIndex: 1, upperBound: 3})

	definition, err := readWriteModel.UmasArrayTypeDefinitionParse(t.Context(), payload)
	require.NoError(t, err)
	assert.Equal(t, arrayClassId, definition.GetClassId())
	assert.Equal(t, typeIdDint, definition.GetElementTypeId())
	require.Len(t, definition.GetDimensions(), 2)
	assert.Equal(t, uint32(0), definition.GetDimensions()[0].GetStartIndex())
	assert.Equal(t, uint32(9), definition.GetDimensions()[0].GetUpperBound())
	assert.Equal(t, uint32(1), definition.GetDimensions()[1].GetStartIndex())
	assert.Equal(t, uint32(3), definition.GetDimensions()[1].GetUpperBound())
}

// unlocatedVariableNamesHeaderSizeForTest is range(1) + nextAddress(2) + unknown1(2) +
// noOfRecords(2), the header symbolTablePayload writes before the first record.
const unlocatedVariableNamesHeaderSizeForTest = 7

// The byte order is the one thing these three types do not carry themselves: none of their mspec
// fields pins one, so it is whatever buffer they are handed. Parsing the very same fixture through
// the convenience parser - which builds a default, big endian buffer - has to disagree, otherwise
// the little endian buffer in Dictionary.go would be decoration and a later reader could drop it.
func TestDictionaryRecordsNeedALittleEndianBuffer(t *testing.T) {
	record := symbolRecord{name: "g_r32", dataType: typeIdReal, block: 0x0002, offset: 0x00001234}
	// Only the record itself, without the table header the convenience parser knows nothing about.
	payload := symbolTablePayload(record)[unlocatedVariableNamesHeaderSizeForTest:]

	littleEndian, err := parseSymbolTable(t.Context(), symbolTablePayload(record))
	require.NoError(t, err)
	require.Len(t, littleEndian, 1)

	bigEndian, err := readWriteModel.UmasUnlocatedVariableReferenceParse(t.Context(), payload)
	require.NoError(t, err)

	assert.Equal(t, "g_r32", bigEndian.GetValue(), "a terminated string reads the same either way")
	assert.NotEqual(t, littleEndian[0].GetDataType(), bigEndian.GetDataType())
	assert.NotEqual(t, littleEndian[0].GetBlock(), bigEndian.GetBlock())
	assert.NotEqual(t, littleEndian[0].GetOffset(), bigEndian.GetOffset())
}

// A datatype record carries a reserved byte the mspec pins to 0x00. The Go spi fails the parse when
// it is something else, where plc4j's reserved field reader only logs it - so a device that puts
// anything there loses the whole dictionary rather than the one record. That is a property of the
// shared field reader, not of this driver, and this test is here so the difference is recorded
// rather than discovered.
func TestParseDatatypeNames_ANonZeroReservedByteFailsTheParse(t *testing.T) {
	payload := datatypeDictionaryPayload(datatypeRecord{name: "MY_STRING", dataSize: 20})
	records, err := parseDatatypeNames(t.Context(), payload)
	require.NoError(t, err)
	require.Len(t, records, 1)

	// The reserved byte sits after dataSize(2) + unknown1(2) + classIdentifier(1) + dataType(1) of
	// the first record, which itself follows the 6 byte header.
	payload[6+6] = 0xFF
	_, err = parseDatatypeNames(t.Context(), payload)
	assert.Error(t, err)
}
