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
	"github.com/apache/plc4x/plc4go/spi/values"
)

// A Modbus register is 16 bits wide, so a value narrower than that is padded when it stands alone
// and packed when it does not. Those rules used to be separate cases in the protocol description
// and now live in RegisterCodec.go, which is why they are pinned down here directly.

func parseFor(t *testing.T, data []byte, dataType readWriteModel.ModbusDataType, numberOfValues uint16, bigEndian bool) apiValues.PlcValue {
	t.Helper()
	return parseWith(t, data, dataType, numberOfValues, byteOrderFor(bigEndian), 1)
}

func serializeFor(t *testing.T, value apiValues.PlcValue, dataType readWriteModel.ModbusDataType, numberOfValues uint16, bigEndian bool) []byte {
	t.Helper()
	return serializeWith(t, value, dataType, numberOfValues, byteOrderFor(bigEndian), 1)
}

func byteOrderFor(bigEndian bool) ByteOrder {
	if bigEndian {
		return BigEndianOrder
	}
	return LittleEndianOrder
}

func parseWith(t *testing.T, data []byte, dataType readWriteModel.ModbusDataType, numberOfValues uint16, byteOrder ByteOrder, stringLength uint16) apiValues.PlcValue {
	t.Helper()
	value, err := ParseRegisters(context.Background(), data, dataType, numberOfValues, byteOrder, stringLength)
	require.NoError(t, err)
	return value
}

func serializeWith(t *testing.T, value apiValues.PlcValue, dataType readWriteModel.ModbusDataType, numberOfValues uint16, byteOrder ByteOrder, stringLength uint16) []byte {
	t.Helper()
	data, err := SerializeRegisters(context.Background(), value, dataType, numberOfValues, byteOrder, stringLength)
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
	assert.Equal(t, uint64(1), lengthInBytes(readWriteModel.ModbusDataType_CHAR, 1, 1))
}

// Values wide enough to fill a register are never padded.
func TestNeverPadsARegisterWideValue(t *testing.T) {
	assert.Equal(t, uint16(0x1234), parseFor(t, []byte{0x12, 0x34}, readWriteModel.ModbusDataType_UINT, 1, true).GetUint16())
	assert.Equal(t, uint64(2), lengthInBytes(readWriteModel.ModbusDataType_UINT, 1, 1))
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

	assert.Equal(t, uint64(4), lengthInBytes(readWriteModel.ModbusDataType_USINT, 3, 1))
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x00}, serializeFor(t, three, readWriteModel.ModbusDataType_USINT, 3, true))
}

// An even count fills its registers exactly and needs no trailing pad.
func TestAddsNoTrailingPadForAnEvenCount(t *testing.T) {
	two := values.NewPlcList([]apiValues.PlcValue{values.NewPlcUSINT(1), values.NewPlcUSINT(2)})

	assert.Equal(t, uint64(2), lengthInBytes(readWriteModel.ModbusDataType_USINT, 2, 1))
	assert.Equal(t, []byte{0x01, 0x02}, serializeFor(t, two, readWriteModel.ModbusDataType_USINT, 2, true))
}

// Three packed bools still occupy a whole register.
func TestPadsPackedBoolsToTheRegisterBoundary(t *testing.T) {
	threeBools := values.NewPlcList([]apiValues.PlcValue{values.NewPlcBOOL(true), values.NewPlcBOOL(false), values.NewPlcBOOL(true)})

	assert.Equal(t, uint64(2), lengthInBytes(readWriteModel.ModbusDataType_BOOL, 3, 1))
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

// The four byte orders differ only in how the bytes of one value are arranged on the wire. The
// layouts pinned down here are the ones plc4j's ModbusByteOrder enum documents: [1, 2, 3, 4],
// [4, 3, 2, 1], [2, 1, 4, 3] and [3, 4, 1, 2].
func TestLaysOutAFourByteValueAccordingToTheByteOrder(t *testing.T) {
	for _, test := range []struct {
		byteOrder ByteOrder
		expected  []byte
	}{
		{BigEndianOrder, []byte{0x01, 0x02, 0x03, 0x04}},
		{LittleEndianOrder, []byte{0x04, 0x03, 0x02, 0x01}},
		{BigEndianByteSwapOrder, []byte{0x02, 0x01, 0x04, 0x03}},
		{LittleEndianByteSwapOrder, []byte{0x03, 0x04, 0x01, 0x02}},
	} {
		t.Run(test.byteOrder.String(), func(t *testing.T) {
			written := serializeWith(t, values.NewPlcUDINT(0x01020304), readWriteModel.ModbusDataType_UDINT, 1, test.byteOrder, 1)
			assert.Equal(t, test.expected, written)
			assert.Equal(t, uint32(0x01020304), parseWith(t, written, readWriteModel.ModbusDataType_UDINT, 1, test.byteOrder, 1).GetUint32())
		})
	}
}

// And the same for a value spanning four registers.
func TestLaysOutAnEightByteValueAccordingToTheByteOrder(t *testing.T) {
	const value = uint64(0x0102030405060708)
	for _, test := range []struct {
		byteOrder ByteOrder
		expected  []byte
	}{
		{BigEndianOrder, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}},
		{LittleEndianOrder, []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}},
		{BigEndianByteSwapOrder, []byte{0x02, 0x01, 0x04, 0x03, 0x06, 0x05, 0x08, 0x07}},
		{LittleEndianByteSwapOrder, []byte{0x07, 0x08, 0x05, 0x06, 0x03, 0x04, 0x01, 0x02}},
	} {
		t.Run(test.byteOrder.String(), func(t *testing.T) {
			written := serializeWith(t, values.NewPlcULINT(value), readWriteModel.ModbusDataType_ULINT, 1, test.byteOrder, 1)
			assert.Equal(t, test.expected, written)
			assert.Equal(t, value, parseWith(t, written, readWriteModel.ModbusDataType_ULINT, 1, test.byteOrder, 1).GetUint64())
		})
	}
}

// Every byte order has to read back what it wrote, for a lone value and for a packed run alike.
func TestRoundTripsThroughEveryByteOrder(t *testing.T) {
	for _, byteOrder := range []ByteOrder{BigEndianOrder, LittleEndianOrder, BigEndianByteSwapOrder, LittleEndianByteSwapOrder} {
		t.Run(byteOrder.String(), func(t *testing.T) {
			single := serializeWith(t, values.NewPlcDINT(-123456), readWriteModel.ModbusDataType_DINT, 1, byteOrder, 1)
			assert.Equal(t, int32(-123456), parseWith(t, single, readWriteModel.ModbusDataType_DINT, 1, byteOrder, 1).GetInt32())

			list := values.NewPlcList([]apiValues.PlcValue{values.NewPlcINT(1), values.NewPlcINT(-2), values.NewPlcINT(3)})
			packed := serializeWith(t, list, readWriteModel.ModbusDataType_INT, 3, byteOrder, 1)
			read := parseWith(t, packed, readWriteModel.ModbusDataType_INT, 3, byteOrder, 1)
			require.True(t, read.IsList())
			assert.Equal(t, int16(1), read.GetList()[0].GetInt16())
			assert.Equal(t, int16(-2), read.GetList()[1].GetInt16())
			assert.Equal(t, int16(3), read.GetList()[2].GetInt16())
		})
	}
}

// A byte swap exchanges the two bytes of every register and leaves an odd trailing byte alone.
func TestByteSwapLeavesAnOddTrailingByteInPlace(t *testing.T) {
	assert.Equal(t, []byte{0x02, 0x01, 0x04, 0x03, 0x05}, byteSwap([]byte{0x01, 0x02, 0x03, 0x04, 0x05}))
	assert.Empty(t, byteSwap(nil))

	// The input must not be modified - it is the response buffer of the codec.
	input := []byte{0x01, 0x02}
	_ = byteSwap(input)
	assert.Equal(t, []byte{0x01, 0x02}, input)
}

// A string is as long as its declared length says, which is the only thing that tells the codec
// how many bytes belong to it.
func TestRoundTripsAStringOfTheDeclaredLength(t *testing.T) {
	written := serializeWith(t, values.NewPlcSTRING("hi"), readWriteModel.ModbusDataType_STRING, 1, BigEndianOrder, 6)

	assert.Len(t, written, 6, "a STRING(6) occupies six bytes no matter how short the value is")
	assert.Equal(t, uint64(6), lengthInBytes(readWriteModel.ModbusDataType_STRING, 1, 6))
	assert.Equal(t, "hi", parseWith(t, written, readWriteModel.ModbusDataType_STRING, 1, BigEndianOrder, 6).GetString())
}

// A WSTRING declares its length in characters too, and every character is two bytes wide.
func TestRoundTripsAWStringOfTheDeclaredLength(t *testing.T) {
	written := serializeWith(t, values.NewPlcWSTRING("hi"), readWriteModel.ModbusDataType_WSTRING, 1, BigEndianOrder, 4)

	assert.Len(t, written, 8, "a WSTRING(4) occupies eight bytes")
	assert.Equal(t, uint64(8), lengthInBytes(readWriteModel.ModbusDataType_WSTRING, 1, 4))
	assert.Equal(t, "hi", parseWith(t, written, readWriteModel.ModbusDataType_WSTRING, 1, BigEndianOrder, 4).GetString())
}

// Several strings follow one another, each one as long as the declared length.
func TestPacksSeveralStrings(t *testing.T) {
	both := values.NewPlcList([]apiValues.PlcValue{values.NewPlcSTRING("ab"), values.NewPlcSTRING("cd")})
	written := serializeWith(t, both, readWriteModel.ModbusDataType_STRING, 2, BigEndianOrder, 4)

	assert.Equal(t, []byte{'a', 'b', 0, 0, 'c', 'd', 0, 0}, written)
	read := parseWith(t, written, readWriteModel.ModbusDataType_STRING, 2, BigEndianOrder, 4)
	require.True(t, read.IsList())
	assert.Equal(t, "ab", read.GetList()[0].GetString())
	assert.Equal(t, "cd", read.GetList()[1].GetString())
}

// An odd number of strings must not pick up a pad byte: a string is padded to a whole register by
// its own declared length, not by the single character the element width used to be computed from.
// With the pad byte in place a 'holding-register:1:STRING(20)[3]' write serialized to 61 bytes
// while the read side asked for 30 registers.
func TestPacksAnOddNumberOfStringsWithoutAPadByte(t *testing.T) {
	three := values.NewPlcList([]apiValues.PlcValue{
		values.NewPlcSTRING("ab"), values.NewPlcSTRING("cd"), values.NewPlcSTRING("ef"),
	})
	written := serializeWith(t, three, readWriteModel.ModbusDataType_STRING, 3, BigEndianOrder, 4)

	assert.Len(t, written, 12, "three STRING(4) occupy twelve bytes, six whole registers")
	assert.Equal(t, uint64(12), lengthInBytes(readWriteModel.ModbusDataType_STRING, 3, 4))
	assert.Equal(t, uint16(0), trailingPaddingBits(readWriteModel.ModbusDataType_STRING, 3, 4))
	read := parseWith(t, written, readWriteModel.ModbusDataType_STRING, 3, BigEndianOrder, 4)
	require.True(t, read.IsList())
	assert.Equal(t, "ef", read.GetList()[2].GetString())
}

// The same for wide strings, whose characters are two bytes each.
func TestPacksAnOddNumberOfWStringsWithoutAPadByte(t *testing.T) {
	three := values.NewPlcList([]apiValues.PlcValue{
		values.NewPlcWSTRING("ab"), values.NewPlcWSTRING("cd"), values.NewPlcWSTRING("ef"),
	})
	written := serializeWith(t, three, readWriteModel.ModbusDataType_WSTRING, 3, BigEndianOrder, 2)

	assert.Len(t, written, 12, "three WSTRING(2) occupy twelve bytes")
	assert.Equal(t, uint64(12), lengthInBytes(readWriteModel.ModbusDataType_WSTRING, 3, 2))
}

// A run of values narrower than a register is still rounded up to a whole one.
func TestPacksSubRegisterValuesUpToAWholeRegister(t *testing.T) {
	assert.Equal(t, uint16(8), trailingPaddingBits(readWriteModel.ModbusDataType_SINT, 3, 1))
	assert.Equal(t, uint16(0), trailingPaddingBits(readWriteModel.ModbusDataType_SINT, 4, 1))
	assert.Equal(t, uint16(13), trailingPaddingBits(readWriteModel.ModbusDataType_BOOL, 3, 1))
	// A STRING(1) is a single character wide, so an odd number of them does need the pad.
	assert.Equal(t, uint16(8), trailingPaddingBits(readWriteModel.ModbusDataType_STRING, 3, 1))
}
