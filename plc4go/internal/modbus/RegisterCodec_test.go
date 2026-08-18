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

package modbus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/apache/plc4x/plc4go/spi/values"
)

// A Modbus register is 16 bits wide, so a value narrower than that is padded when it stands alone
// and packed when it does not. Those rules used to be separate cases in the protocol description
// and now live in RegisterCodec.go, which is why they are pinned down here directly.

func parseFor(t *testing.T, data []byte, dataType readWriteModel.ModbusDataType, numberOfValues uint16, bigEndian bool) apiValues.PlcValue {
	t.Helper()
	value, err := ParseRegisters(context.Background(), utils.NewReadBufferByteBased(data), dataType, numberOfValues, bigEndian, 1)
	require.NoError(t, err)
	return value
}

func serializeFor(t *testing.T, value apiValues.PlcValue, dataType readWriteModel.ModbusDataType, numberOfValues uint16, bigEndian bool) []byte {
	t.Helper()
	data, err := SerializeRegisters(context.Background(), value, dataType, numberOfValues, bigEndian, 1)
	require.NoError(t, err)
	return data
}

// Big endian puts the padding in front of a lone byte, little endian behind it.
func TestPadsALoneByteAccordingToByteOrder(t *testing.T) {
	assert.Equal(t, uint8(0x2A), parseFor(t, []byte{0x00, 0x2A}, readWriteModel.ModbusDataType_USINT, 1, true).GetUint8())
	assert.Equal(t, uint8(0x2A), parseFor(t, []byte{0x2A, 0x00}, readWriteModel.ModbusDataType_USINT, 1, false).GetUint8())

	assert.Equal(t, []byte{0x00, 0x2A}, serializeFor(t, values.NewPlcUSINT(0x2A), readWriteModel.ModbusDataType_USINT, 1, true))
	assert.Equal(t, []byte{0x2A, 0x00}, serializeFor(t, values.NewPlcUSINT(0x2A), readWriteModel.ModbusDataType_USINT, 1, false))
}

// A signed byte is padded the same way and keeps its sign.
func TestPadsALoneSignedByte(t *testing.T) {
	assert.Equal(t, int8(-1), parseFor(t, []byte{0x00, 0xFF}, readWriteModel.ModbusDataType_SINT, 1, true).GetInt8())
	assert.Equal(t, int8(-1), parseFor(t, []byte{0xFF, 0x00}, readWriteModel.ModbusDataType_SINT, 1, false).GetInt8())
}

// A lone BOOL fills a whole register. Big endian leaves the bit last; little endian puts it after
// seven bits, with eight more behind it.
func TestPadsALoneBoolToAWholeRegister(t *testing.T) {
	assert.True(t, parseFor(t, []byte{0x00, 0x01}, readWriteModel.ModbusDataType_BOOL, 1, true).GetBool())
	assert.False(t, parseFor(t, []byte{0x00, 0x00}, readWriteModel.ModbusDataType_BOOL, 1, true).GetBool())
	assert.True(t, parseFor(t, []byte{0x01, 0x00}, readWriteModel.ModbusDataType_BOOL, 1, false).GetBool())
}

// CHAR is not padded - it never has been, unlike the other single byte types.
func TestLeavesALoneCharUnpadded(t *testing.T) {
	assert.Equal(t, uint16(1), lengthInBytes(readWriteModel.ModbusDataType_CHAR, 1, 1))
}

// Values wide enough to fill a register are never padded.
func TestNeverPadsARegisterWideValue(t *testing.T) {
	assert.Equal(t, uint16(0x1234), parseFor(t, []byte{0x12, 0x34}, readWriteModel.ModbusDataType_UINT, 1, true).GetUint16())
	assert.Equal(t, uint16(2), lengthInBytes(readWriteModel.ModbusDataType_UINT, 1, 1))
}

// Several sub-register values are packed, two bytes to a register, rather than padded each.
func TestPacksSeveralBytesWithoutPadding(t *testing.T) {
	value := parseFor(t, []byte{0x01, 0x02, 0x03, 0x04}, readWriteModel.ModbusDataType_USINT, 4, true)

	require.True(t, value.IsList())
	assert.Len(t, value.GetList(), 4)
	assert.Equal(t, uint8(1), value.GetList()[0].GetUint8())
	assert.Equal(t, uint8(4), value.GetList()[3].GetUint8())
}

// Packed BOOLs are one bit each, not one register each.
func TestPacksBoolsAsBits(t *testing.T) {
	value := parseFor(t, []byte{0xA0, 0x00}, readWriteModel.ModbusDataType_BOOL, 3, true)

	require.True(t, value.IsList())
	assert.True(t, value.GetList()[0].GetBool())
	assert.False(t, value.GetList()[1].GetBool())
	assert.True(t, value.GetList()[2].GetBool())
}

// An odd number of bytes leaves half a register, which has to be padded out when writing.
func TestPadsTheLastRegisterForAnOddCount(t *testing.T) {
	three := values.NewPlcList([]apiValues.PlcValue{values.NewPlcUSINT(1), values.NewPlcUSINT(2), values.NewPlcUSINT(3)})

	assert.Equal(t, uint16(4), lengthInBytes(readWriteModel.ModbusDataType_USINT, 3, 1))
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x00}, serializeFor(t, three, readWriteModel.ModbusDataType_USINT, 3, true))
}

// An even count fills its registers exactly and needs no trailing pad.
func TestAddsNoTrailingPadForAnEvenCount(t *testing.T) {
	two := values.NewPlcList([]apiValues.PlcValue{values.NewPlcUSINT(1), values.NewPlcUSINT(2)})

	assert.Equal(t, uint16(2), lengthInBytes(readWriteModel.ModbusDataType_USINT, 2, 1))
	assert.Equal(t, []byte{0x01, 0x02}, serializeFor(t, two, readWriteModel.ModbusDataType_USINT, 2, true))
}

// Three packed bools still occupy a whole register.
func TestPadsPackedBoolsToTheRegisterBoundary(t *testing.T) {
	threeBools := values.NewPlcList([]apiValues.PlcValue{values.NewPlcBOOL(true), values.NewPlcBOOL(false), values.NewPlcBOOL(true)})

	assert.Equal(t, uint16(2), lengthInBytes(readWriteModel.ModbusDataType_BOOL, 3, 1))
	assert.Equal(t, []byte{0xA0, 0x00}, serializeFor(t, threeBools, readWriteModel.ModbusDataType_BOOL, 3, true))
}

// What is written for a padded value has to read back as the same value, in either byte order.
func TestRoundTripsAPaddedValue(t *testing.T) {
	for _, bigEndian := range []bool{true, false} {
		written := serializeFor(t, values.NewPlcSINT(-42), readWriteModel.ModbusDataType_SINT, 1, bigEndian)

		assert.Len(t, written, 2, "a lone byte occupies a whole register")
		assert.Equal(t, int8(-42), parseFor(t, written, readWriteModel.ModbusDataType_SINT, 1, bigEndian).GetInt8())
	}
}

// And so does a packed run.
func TestRoundTripsAPackedRun(t *testing.T) {
	input := values.NewPlcList([]apiValues.PlcValue{values.NewPlcSINT(-1), values.NewPlcSINT(2), values.NewPlcSINT(-3)})
	written := serializeFor(t, input, readWriteModel.ModbusDataType_SINT, 3, true)

	read := parseFor(t, written, readWriteModel.ModbusDataType_SINT, 3, true)
	assert.Equal(t, int8(-1), read.GetList()[0].GetInt8())
	assert.Equal(t, int8(2), read.GetList()[1].GetInt8())
	assert.Equal(t, int8(-3), read.GetList()[2].GetInt8())
}
