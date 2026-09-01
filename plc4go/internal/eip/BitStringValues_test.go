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
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// The CIP bit-string family - BYTE, WORD, DWORD and LWORD - are unsigned N-byte bit strings.
// All-bits-set is an entirely ordinary value for them, so the decoding must not sign extend.

func TestParsePlcValueSingleBYTE(t *testing.T) {
	value, err := parsePlcValue(mustTag(t, "%b:BYTE"), []byte{0xFF}, readWriteModel.CIPDataTypeCode_BYTE)
	require.NoError(t, err)
	assert.Equal(t, uint8(0xFF), value.GetUint8())
}

func TestParsePlcValueSingleWORD(t *testing.T) {
	raw := binary.LittleEndian.AppendUint16(nil, 0xFFFF)
	value, err := parsePlcValue(mustTag(t, "%w:WORD"), raw, readWriteModel.CIPDataTypeCode_WORD)
	require.NoError(t, err)
	assert.Equal(t, uint16(0xFFFF), value.GetUint16())
}

func TestParsePlcValueSingleDWORD(t *testing.T) {
	raw := binary.LittleEndian.AppendUint32(nil, 0xFFFFFFFF)
	value, err := parsePlcValue(mustTag(t, "%d:DWORD"), raw, readWriteModel.CIPDataTypeCode_DWORD)
	require.NoError(t, err)
	assert.Equal(t, uint32(0xFFFFFFFF), value.GetUint32())
}

// Bit 31 set is the case a signed decode turns into a negative number.
func TestParsePlcValueDWORDHighBitSet(t *testing.T) {
	raw := binary.LittleEndian.AppendUint32(nil, 0x80000000)
	value, err := parsePlcValue(mustTag(t, "%d:DWORD"), raw, readWriteModel.CIPDataTypeCode_DWORD)
	require.NoError(t, err)
	assert.Equal(t, uint32(0x80000000), value.GetUint32())
}

func TestParsePlcValueSingleLWORD(t *testing.T) {
	raw := binary.LittleEndian.AppendUint64(nil, 0xFFFFFFFFFFFFFFFF)
	value, err := parsePlcValue(mustTag(t, "%l:LWORD"), raw, readWriteModel.CIPDataTypeCode_LWORD)
	require.NoError(t, err)
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), value.GetUint64())
}

func TestParsePlcValueArrayOfDWORDs(t *testing.T) {
	raw := binary.LittleEndian.AppendUint32(nil, 0x00000001)
	raw = binary.LittleEndian.AppendUint32(raw, 0xFFFFFFFF)
	raw = binary.LittleEndian.AppendUint32(raw, 0x80000000)

	value, err := parsePlcValue(mustTag(t, "%d[0..2]:DWORD"), raw, readWriteModel.CIPDataTypeCode_DWORD)
	require.NoError(t, err)
	require.True(t, value.IsList())
	list := value.GetList()
	require.Len(t, list, 3)
	assert.Equal(t, uint32(0x00000001), list[0].GetUint32())
	assert.Equal(t, uint32(0xFFFFFFFF), list[1].GetUint32())
	assert.Equal(t, uint32(0x80000000), list[2].GetUint32())
}

func TestParsePlcValueArrayOfLWORDs(t *testing.T) {
	raw := binary.LittleEndian.AppendUint64(nil, 1)
	raw = binary.LittleEndian.AppendUint64(raw, 0xFFFFFFFFFFFFFFFF)

	value, err := parsePlcValue(mustTag(t, "%l[0..1]:LWORD"), raw, readWriteModel.CIPDataTypeCode_LWORD)
	require.NoError(t, err)
	list := value.GetList()
	require.Len(t, list, 2)
	assert.Equal(t, uint64(1), list[0].GetUint64())
	assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), list[1].GetUint64())
}

func TestParsePlcValueArrayOfBYTEsAndWORDs(t *testing.T) {
	bytesValue, err := parsePlcValue(mustTag(t, "%b[0..2]:BYTE"), []byte{0x01, 0x80, 0xFF}, readWriteModel.CIPDataTypeCode_BYTE)
	require.NoError(t, err)
	assert.Equal(t, uint8(0xFF), bytesValue.GetList()[2].GetUint8())

	raw := binary.LittleEndian.AppendUint16(nil, 0x8000)
	raw = binary.LittleEndian.AppendUint16(raw, 0xFFFF)
	wordsValue, err := parsePlcValue(mustTag(t, "%w[0..1]:WORD"), raw, readWriteModel.CIPDataTypeCode_WORD)
	require.NoError(t, err)
	assert.Equal(t, uint16(0x8000), wordsValue.GetList()[0].GetUint16())
	assert.Equal(t, uint16(0xFFFF), wordsValue.GetList()[1].GetUint16())
}

// The bit-string types are fixed size, so a reply shorter than the declared element count is
// reported as an error rather than read past the end of the buffer.
func TestParsePlcValueBitStringShortReply(t *testing.T) {
	_, err := parsePlcValue(mustTag(t, "%d[0..7]:DWORD"), binary.LittleEndian.AppendUint32(nil, 1), readWriteModel.CIPDataTypeCode_DWORD)
	assert.Error(t, err)

	_, err = parsePlcValue(mustTag(t, "%l[0..3]:LWORD"), binary.LittleEndian.AppendUint64(nil, 1), readWriteModel.CIPDataTypeCode_LWORD)
	assert.Error(t, err)
}

func TestEncodeValueBitStrings(t *testing.T) {
	for _, test := range []struct {
		name     string
		value    apiValues.PlcValue
		dataType readWriteModel.CIPDataTypeCode
		expected []byte
	}{
		{"BYTE", spiValues.NewPlcBYTE(0xFF), readWriteModel.CIPDataTypeCode_BYTE, []byte{0xFF}},
		{"WORD", spiValues.NewPlcWORD(0xFFFF), readWriteModel.CIPDataTypeCode_WORD, []byte{0xFF, 0xFF}},
		{"DWORD", spiValues.NewPlcDWORD(0x80000000), readWriteModel.CIPDataTypeCode_DWORD, []byte{0x00, 0x00, 0x00, 0x80}},
		{"LWORD", spiValues.NewPlcLWORD(0xFFFFFFFFFFFFFFFF), readWriteModel.CIPDataTypeCode_LWORD,
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

// LWORD used to share 0x00D3 with DWORD and STRINGI used to share 0x00DD with ENGUNIT, which
// dropped one of each pair from the generated lookup tables entirely.
func TestCIPDataTypeCodesAreUnique(t *testing.T) {
	seen := map[readWriteModel.CIPDataTypeCode]string{}
	for _, code := range readWriteModel.CIPDataTypeCodeValues {
		previous, duplicate := seen[code]
		assert.Falsef(t, duplicate, "%s and %s share value %#04x", previous, code, uint16(code))
		seen[code] = code.String()
	}
}

func TestBitStringCodesMatchTheCipSpecification(t *testing.T) {
	assert.Equal(t, readWriteModel.CIPDataTypeCode(0x00D1), readWriteModel.CIPDataTypeCode_BYTE)
	assert.Equal(t, readWriteModel.CIPDataTypeCode(0x00D2), readWriteModel.CIPDataTypeCode_WORD)
	assert.Equal(t, readWriteModel.CIPDataTypeCode(0x00D3), readWriteModel.CIPDataTypeCode_DWORD)
	assert.Equal(t, readWriteModel.CIPDataTypeCode(0x00D4), readWriteModel.CIPDataTypeCode_LWORD)

	assert.Equal(t, uint8(1), readWriteModel.CIPDataTypeCode_BYTE.Size())
	assert.Equal(t, uint8(2), readWriteModel.CIPDataTypeCode_WORD.Size())
	assert.Equal(t, uint8(4), readWriteModel.CIPDataTypeCode_DWORD.Size())
	assert.Equal(t, uint8(8), readWriteModel.CIPDataTypeCode_LWORD.Size())
}

// The duplicate value made this lookup fail for the dropped name, so a %tag:LWORD address was
// rejected as an unknown data type before it ever reached the decoder.
func TestTagHandlerAcceptsBitStringTypes(t *testing.T) {
	for _, name := range []string{"BYTE", "WORD", "DWORD", "LWORD"} {
		t.Run(name, func(t *testing.T) {
			code, found := readWriteModel.CIPDataTypeCodeByName(name)
			require.Truef(t, found, "%s is not resolvable by name", name)
			assert.Equal(t, code, mustTag(t, "%tag:"+name).GetType())
		})
	}
}

// --- invariants that hold for every entry in the codec table ---

// Every fixed-size type must decode a payload and encode it back to the identical bytes, in
// exactly Size() bytes. Driving this off the table itself means a type added later is covered
// without anyone remembering to extend this test.
func TestEveryFixedSizeTypeRoundTripsItsOwnBytes(t *testing.T) {
	for dataType, codec := range fixedSizeCodecs {
		if dataType == readWriteModel.CIPDataTypeCode_BOOL {
			continue // BOOL normalises any non-zero byte to 1, see TestEncodeValueBool
		}
		for _, filler := range []byte{0x00, 0x7F, 0xFF} {
			if filler == 0xFF &&
				(dataType == readWriteModel.CIPDataTypeCode_REAL || dataType == readWriteModel.CIPDataTypeCode_LREAL) {
				continue // all bits set is NaN for the floats, which is not a value to compare
			}
			t.Run(fmt.Sprintf("%s/%#02x", dataType, filler), func(t *testing.T) {
				raw := bytes.Repeat([]byte{filler}, int(dataType.Size()))

				decoded := codec.read(raw, 0)
				require.NotNil(t, decoded)

				encoded, err := encodeValue(decoded, dataType)
				require.NoError(t, err)
				assert.Len(t, encoded, int(dataType.Size()))
				assert.Equal(t, raw, encoded)
			})
		}
	}
}

// BOOL is the one type whose encoding is not the identity, so it is checked on its own.
func TestEncodeValueBool(t *testing.T) {
	encoded, err := encodeValue(spiValues.NewPlcBOOL(true), readWriteModel.CIPDataTypeCode_BOOL)
	require.NoError(t, err)
	assert.Equal(t, []byte{1}, encoded)

	encoded, err = encodeValue(spiValues.NewPlcBOOL(false), readWriteModel.CIPDataTypeCode_BOOL)
	require.NoError(t, err)
	assert.Equal(t, []byte{0}, encoded)
}

// An unsupported type must be reported, not silently encoded.
func TestEncodeValueRejectsUnsupportedType(t *testing.T) {
	_, err := encodeValue(spiValues.NewPlcDINT(1), readWriteModel.CIPDataTypeCode_TIME)
	assert.Error(t, err)
}
