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
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/eip/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func mustTag(t *testing.T, address string) PlcTag {
	tag, err := NewTagHandler().ParseTag(address)
	require.NoError(t, err)
	return tag.(PlcTag)
}

func TestParsePlcValueSingleDINT(t *testing.T) {
	value, err := parsePlcValue(mustTag(t, "%rate:DINT"), []byte{0x16, 0x02, 0x00, 0x00}, readWriteModel.CIPDataTypeCode_DINT)
	require.NoError(t, err)
	assert.Equal(t, int32(0x216), value.GetInt32())
}

func TestParsePlcValueSingleSINT(t *testing.T) {
	value, err := parsePlcValue(mustTag(t, "%b:SINT"), []byte{0x2A}, readWriteModel.CIPDataTypeCode_SINT)
	require.NoError(t, err)
	assert.Equal(t, int8(0x2A), value.GetInt8())
}

func TestParsePlcValueSingleINT(t *testing.T) {
	value, err := parsePlcValue(mustTag(t, "%w:INT"), []byte{0x16, 0x02}, readWriteModel.CIPDataTypeCode_INT)
	require.NoError(t, err)
	assert.Equal(t, int16(0x216), value.GetInt16())
}

func TestParsePlcValueSingleLINT(t *testing.T) {
	raw := binary.LittleEndian.AppendUint64(nil, uint64(0x0102030405060708))
	value, err := parsePlcValue(mustTag(t, "%l:LINT"), raw, readWriteModel.CIPDataTypeCode_LINT)
	require.NoError(t, err)
	assert.Equal(t, int64(0x0102030405060708), value.GetInt64())
}

func TestParsePlcValueSingleLREAL(t *testing.T) {
	raw := binary.LittleEndian.AppendUint64(nil, math.Float64bits(2.5))
	value, err := parsePlcValue(mustTag(t, "%d:LREAL"), raw, readWriteModel.CIPDataTypeCode_LREAL)
	require.NoError(t, err)
	assert.Equal(t, 2.5, value.GetFloat64())
}

func TestParsePlcValueUnsupportedTypeIsError(t *testing.T) {
	_, err := parsePlcValue(mustTag(t, "%x:DINT"), []byte{0x00}, readWriteModel.CIPDataTypeCode(0xFFFF))
	require.Error(t, err)
}

func TestParsePlcValueSingleREAL(t *testing.T) {
	raw := binary.LittleEndian.AppendUint32(nil, math.Float32bits(1.5))
	value, err := parsePlcValue(mustTag(t, "%f:REAL"), raw, readWriteModel.CIPDataTypeCode_REAL)
	require.NoError(t, err)
	assert.Equal(t, float32(1.5), value.GetFloat32())
}

func TestParsePlcValueSingleBOOL(t *testing.T) {
	value, err := parsePlcValue(mustTag(t, "%b:BOOL"), []byte{0x01}, readWriteModel.CIPDataTypeCode_BOOL)
	require.NoError(t, err)
	assert.True(t, value.GetBool())
}

func TestParsePlcValueDINTArray(t *testing.T) {
	raw := []byte{0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00}
	value, err := parsePlcValue(mustTag(t, "%arr[0]:DINT:2"), raw, readWriteModel.CIPDataTypeCode_DINT)
	require.NoError(t, err)
	require.True(t, value.IsList())
	list := value.GetList()
	require.Len(t, list, 2)
	assert.Equal(t, int32(1), list[0].GetInt32())
	assert.Equal(t, int32(2), list[1].GetInt32())
}

func TestParsePlcValueShortReplyIsError(t *testing.T) {
	// 2 DINT elements requested but only 4 bytes returned: must error, not panic (GH-954 thread)
	_, err := parsePlcValue(mustTag(t, "%arr[0]:DINT:2"), []byte{0x01, 0x00, 0x00, 0x00}, readWriteModel.CIPDataTypeCode_DINT)
	require.Error(t, err)
}

func TestParsePlcValueString(t *testing.T) {
	// structured header: type 0x0FCE (STRING), len 2, 2 pad bytes, then "AB"
	raw := []byte{0xCE, 0x0F, 0x02, 0x00, 0x00, 0x00, 0x41, 0x42}
	value, err := parsePlcValue(mustTag(t, "%s:STRING"), raw, readWriteModel.CIPDataTypeCode_STRING)
	require.NoError(t, err)
	assert.Equal(t, "AB", value.GetString())
}

func TestEncodeValueDINT(t *testing.T) {
	raw, err := encodeValue(spiValues.NewPlcDINT(0x216), readWriteModel.CIPDataTypeCode_DINT)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x16, 0x02, 0x00, 0x00}, raw)
}

func TestEncodeValueSINT(t *testing.T) {
	raw, err := encodeValue(spiValues.NewPlcSINT(0x2A), readWriteModel.CIPDataTypeCode_SINT)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x2A}, raw)
}

func TestEncodeValueINT(t *testing.T) {
	raw, err := encodeValue(spiValues.NewPlcINT(0x216), readWriteModel.CIPDataTypeCode_INT)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x16, 0x02}, raw)
}

func TestEncodeValueLINT(t *testing.T) {
	raw, err := encodeValue(spiValues.NewPlcLINT(0x0102030405060708), readWriteModel.CIPDataTypeCode_LINT)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}, raw)
}

func TestEncodeValueLREAL(t *testing.T) {
	raw, err := encodeValue(spiValues.NewPlcLREAL(2.5), readWriteModel.CIPDataTypeCode_LREAL)
	require.NoError(t, err)
	assert.Equal(t, binary.LittleEndian.AppendUint64(nil, math.Float64bits(2.5)), raw)
}

func TestEncodeValueUnsupportedTypeIsError(t *testing.T) {
	_, err := encodeValue(spiValues.NewPlcDINT(1), readWriteModel.CIPDataTypeCode(0xFFFF))
	require.Error(t, err)
}

func TestEncodeValueREAL(t *testing.T) {
	raw, err := encodeValue(spiValues.NewPlcREAL(1.0), readWriteModel.CIPDataTypeCode_REAL)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x00, 0x00, 0x80, 0x3F}, raw)
}

func TestEncodeValueBOOL(t *testing.T) {
	raw, err := encodeValue(spiValues.NewPlcBOOL(true), readWriteModel.CIPDataTypeCode_BOOL)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x01}, raw)
}

func TestEncodeValueSTRING(t *testing.T) {
	raw, err := encodeValue(spiValues.NewPlcSTRING("AB"), readWriteModel.CIPDataTypeCode_STRING)
	require.NoError(t, err)
	assert.Equal(t, []byte{0x02, 0x00, 0x00, 0x00, 0x41, 0x42}, raw)
}

func TestDecodeResponseCode(t *testing.T) {
	assert.Equal(t, apiModel.PlcResponseCode_OK, decodeResponseCode(0))
	assert.Equal(t, apiModel.PlcResponseCode_INTERNAL_ERROR, decodeResponseCode(5))
}
