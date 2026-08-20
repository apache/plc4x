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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// The type ids are the ones the mspec's UmasDataType enum assigns.
const (
	typeIdBool        = uint16(1)
	typeIdInt         = uint16(4)
	typeIdUint        = uint16(5)
	typeIdDint        = uint16(6)
	typeIdUdint       = uint16(7)
	typeIdReal        = uint16(8)
	typeIdString      = uint16(9)
	typeIdTime        = uint16(10)
	typeIdUnknown11   = uint16(11)
	typeIdDate        = uint16(14)
	typeIdTod         = uint16(15)
	typeIdDateAndTime = uint16(16)
	typeIdByte        = uint16(21)
	typeIdWord        = uint16(22)
	typeIdDword       = uint16(23)
	typeIdEbool       = uint16(25)
	// typeIdCustom is the first id the project's own types get, so it is not a primitive.
	typeIdCustom = customTypeIdBase
)

// The payload layouts are the ones plc4j's UmasConnection.parseReadResponse reads, and for the
// non-temporal types the ones the mspec's DataItem declares - little endian throughout.
func TestDecodeReadResponse(t *testing.T) {
	tests := []struct {
		name       string
		dataTypeId uint16
		block      []byte
		want       apiValues.PlcValue
		wantErr    bool
	}{
		{
			// The dataIo reads 7 reserved bits and then the bit, so the value is the low bit.
			name: "a BOOL is the low bit of one byte", dataTypeId: typeIdBool,
			block: []byte{0x01}, want: spiValues.NewPlcBOOL(true),
		},
		{
			name: "a false BOOL", dataTypeId: typeIdBool,
			block: []byte{0x00}, want: spiValues.NewPlcBOOL(false),
		},
		{
			name: "an EBOOL is read like a BOOL", dataTypeId: typeIdEbool,
			block: []byte{0x01}, want: spiValues.NewPlcBOOL(true),
		},
		{
			name: "a BYTE is one byte", dataTypeId: typeIdByte,
			block: []byte{0xAB}, want: spiValues.NewPlcBYTE(0xAB),
		},
		{
			name: "an INT is two signed bytes, little endian", dataTypeId: typeIdInt,
			block: []byte{0xFE, 0xFF}, want: spiValues.NewPlcINT(-2),
		},
		{
			name: "a UINT is two unsigned bytes, little endian", dataTypeId: typeIdUint,
			block: []byte{0x34, 0x12}, want: spiValues.NewPlcUINT(0x1234),
		},
		{
			name: "a WORD is two unsigned bytes, little endian", dataTypeId: typeIdWord,
			block: []byte{0x34, 0x12}, want: spiValues.NewPlcWORD(0x1234),
		},
		{
			name: "a DINT is four signed bytes, little endian", dataTypeId: typeIdDint,
			block: []byte{0xFF, 0xFF, 0xFF, 0xFF}, want: spiValues.NewPlcDINT(-1),
		},
		{
			name: "a UDINT is four unsigned bytes, little endian", dataTypeId: typeIdUdint,
			block: []byte{0x78, 0x56, 0x34, 0x12}, want: spiValues.NewPlcUDINT(0x12345678),
		},
		{
			name: "a DWORD is four unsigned bytes, little endian", dataTypeId: typeIdDword,
			block: []byte{0x78, 0x56, 0x34, 0x12}, want: spiValues.NewPlcDWORD(0x12345678),
		},
		{
			// 1.0f is 0x3F800000, little endian on the wire.
			name: "a REAL is an IEEE754 single, little endian", dataTypeId: typeIdReal,
			block: []byte{0x00, 0x00, 0x80, 0x3F}, want: spiValues.NewPlcREAL(1.0),
		},
		{
			name: "a STRING ends at its NUL terminator", dataTypeId: typeIdString,
			block: append([]byte("hello"), 0x00, 'j', 'u', 'n', 'k'), want: spiValues.NewPlcSTRING("hello"),
		},
		{
			name: "a STRING without a terminator is the whole buffer", dataTypeId: typeIdString,
			block: []byte("hello"), want: spiValues.NewPlcSTRING("hello"),
		},
		{
			name: "an empty STRING", dataTypeId: typeIdString,
			block: []byte{0x00, 0x00}, want: spiValues.NewPlcSTRING(""),
		},
		{
			// A uint32 of milliseconds, not BCD - the one temporal type which isn't.
			name: "a TIME is milliseconds as a uint32", dataTypeId: typeIdTime,
			block: []byte{0xE8, 0x03, 0x00, 0x00}, want: spiValues.NewPlcTIMEFromMilliseconds(1000),
		},
		{
			// day(1) + month(1) + year(2), all BCD, the low year byte first.
			name: "a DATE is BCD day, month and year", dataTypeId: typeIdDate,
			block: []byte{0x25, 0x12, 0x24, 0x20},
			want:  spiValues.NewPlcDATE(time.Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC)),
		},
		{
			name: "a BCD nibble above 9 is not a date", dataTypeId: typeIdDate,
			block: []byte{0x1F, 0x12, 0x24, 0x20}, wantErr: true,
		},
		{
			name: "month 13 is not a date", dataTypeId: typeIdDate,
			block: []byte{0x01, 0x13, 0x24, 0x20}, wantErr: true,
		},
		{
			name: "February 30th is not a date", dataTypeId: typeIdDate,
			block: []byte{0x30, 0x02, 0x24, 0x20}, wantErr: true,
		},
		{
			// centiseconds(1) + seconds(1) + minutes(1) + hours(1), all BCD.
			name: "a TIME_OF_DAY is BCD centiseconds, seconds, minutes and hours", dataTypeId: typeIdTod,
			block: []byte{0x50, 0x30, 0x15, 0x13},
			want: spiValues.NewPlcTIME_OF_DAY(time.Date(0, 1, 1, 13, 15, 30,
				500*int(time.Millisecond), time.UTC)),
		},
		{
			name: "hour 25 is not a time of day", dataTypeId: typeIdTod,
			block: []byte{0x00, 0x00, 0x00, 0x25}, wantErr: true,
		},
		{
			// reserved(1) + seconds + minutes + hour + day + month + year(2), all BCD.
			name: "a DATE_AND_TIME is BCD, seconds first, after a reserved byte", dataTypeId: typeIdDateAndTime,
			block: []byte{0x00, 0x30, 0x15, 0x13, 0x25, 0x12, 0x24, 0x20},
			want:  spiValues.NewPlcDATE_AND_TIME(time.Date(2024, time.December, 25, 13, 15, 30, 0, time.UTC)),
		},
		{
			name: "a truncated DATE_AND_TIME", dataTypeId: typeIdDateAndTime,
			block: []byte{0x00, 0x30, 0x15}, wantErr: true,
		},
		{
			// The dataIo has no case for the UNKNOWN types; the driver reads the one byte their
			// request size declares, as the type mapping already calls them BYTE.
			name: "an UNKNOWN11 is one byte", dataTypeId: typeIdUnknown11,
			block: []byte{0x07}, want: spiValues.NewPlcBYTE(0x07),
		},
		{
			// The project's own types are opaque to the driver, so the caller gets the bytes.
			name: "a custom type stays raw", dataTypeId: typeIdCustom,
			block: []byte{0x01, 0x02, 0x03}, want: spiValues.NewPlcRawByteArray([]byte{0x01, 0x02, 0x03}),
		},
		{
			name: "an empty payload is not a value", dataTypeId: typeIdDint,
			block: []byte{}, wantErr: true,
		},
		{
			name: "a payload too short for the type", dataTypeId: typeIdDint,
			block: []byte{0x01, 0x02}, wantErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			value, err := decodeReadResponse(t.Context(), testCase.dataTypeId, testCase.block)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Nil(t, value)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, value)
		})
	}
}

// The write payloads are what plc4j's UmasConnection.serializeForType builds, byte for byte.
func TestEncodeWriteValue(t *testing.T) {
	tests := []struct {
		name       string
		dataTypeId uint16
		value      apiValues.PlcValue
		want       []byte
		wantErr    bool
	}{
		{
			name: "a true BOOL is one", dataTypeId: typeIdBool,
			value: spiValues.NewPlcBOOL(true), want: []byte{0x01},
		},
		{
			name: "a false BOOL is zero", dataTypeId: typeIdBool,
			value: spiValues.NewPlcBOOL(false), want: []byte{0x00},
		},
		{
			name: "a BYTE is one byte", dataTypeId: typeIdByte,
			value: spiValues.NewPlcBYTE(0xAB), want: []byte{0xAB},
		},
		{
			name: "an INT is two signed bytes, little endian", dataTypeId: typeIdInt,
			value: spiValues.NewPlcINT(-2), want: []byte{0xFE, 0xFF},
		},
		{
			name: "a UINT is two unsigned bytes, little endian", dataTypeId: typeIdUint,
			value: spiValues.NewPlcUINT(0x1234), want: []byte{0x34, 0x12},
		},
		{
			name: "a WORD is two unsigned bytes, little endian", dataTypeId: typeIdWord,
			value: spiValues.NewPlcWORD(0x1234), want: []byte{0x34, 0x12},
		},
		{
			name: "a DINT is four signed bytes, little endian", dataTypeId: typeIdDint,
			value: spiValues.NewPlcDINT(-1), want: []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name: "a UDINT is four unsigned bytes, little endian", dataTypeId: typeIdUdint,
			value: spiValues.NewPlcUDINT(0x12345678), want: []byte{0x78, 0x56, 0x34, 0x12},
		},
		{
			name: "a DWORD is four unsigned bytes, little endian", dataTypeId: typeIdDword,
			value: spiValues.NewPlcDWORD(0x12345678), want: []byte{0x78, 0x56, 0x34, 0x12},
		},
		{
			name: "a REAL is an IEEE754 single, little endian", dataTypeId: typeIdReal,
			value: spiValues.NewPlcREAL(1.0), want: []byte{0x00, 0x00, 0x80, 0x3F},
		},
		{
			// plc4j appends exactly one NUL terminator and no padding; the buffer size only bounds
			// what the PLC keeps.
			name: "a STRING is its bytes plus a NUL terminator", dataTypeId: typeIdString,
			value: spiValues.NewPlcSTRING("hi"), want: []byte{'h', 'i', 0x00},
		},
		{
			name: "an empty STRING is just the terminator", dataTypeId: typeIdString,
			value: spiValues.NewPlcSTRING(""), want: []byte{0x00},
		},
		{
			// plc4j encodes with US_ASCII, which silently turns anything else into a '?'.
			name: "a STRING outside ASCII has no known encoding", dataTypeId: typeIdString,
			value: spiValues.NewPlcSTRING("grün"), wantErr: true,
		},
		{
			name: "a STRING can't carry its own terminator", dataTypeId: typeIdString,
			value: spiValues.NewPlcSTRING("a\x00b"), wantErr: true,
		},
		{
			name: "a TIME is milliseconds as a uint32", dataTypeId: typeIdTime,
			value: spiValues.NewPlcTIME(1500 * time.Millisecond), want: []byte{0xDC, 0x05, 0x00, 0x00},
		},
		{
			name: "a negative TIME doesn't fit an unsigned field", dataTypeId: typeIdTime,
			value: spiValues.NewPlcTIME(-time.Second), wantErr: true,
		},
		{
			// plc4j takes a plain number as milliseconds here (value.getLong()), so a caller which
			// hands over 1500 rather than a duration gets the same four bytes.
			name: "a plain number is a TIME in milliseconds", dataTypeId: typeIdTime,
			value: spiValues.NewPlcUDINT(1500), want: []byte{0xDC, 0x05, 0x00, 0x00},
		},
		{
			name: "a negative number is not a TIME", dataTypeId: typeIdTime,
			value: spiValues.NewPlcDINT(-1), wantErr: true,
		},
		{
			// 2024 becomes the BCD pair 24, 20 - the low two digits first.
			name: "a DATE is BCD day, month and year", dataTypeId: typeIdDate,
			value: spiValues.NewPlcDATE(time.Date(2024, time.December, 25, 0, 0, 0, 0, time.UTC)),
			want:  []byte{0x25, 0x12, 0x24, 0x20},
		},
		{
			name: "a TIME_OF_DAY is BCD centiseconds, seconds, minutes and hours", dataTypeId: typeIdTod,
			value: spiValues.NewPlcTIME_OF_DAY(time.Date(0, 1, 1, 13, 15, 30, 500*int(time.Millisecond), time.UTC)),
			want:  []byte{0x50, 0x30, 0x15, 0x13},
		},
		{
			name: "a DATE_AND_TIME starts with a reserved zero byte", dataTypeId: typeIdDateAndTime,
			value: spiValues.NewPlcDATE_AND_TIME(time.Date(2024, time.December, 25, 13, 15, 30, 0, time.UTC)),
			want:  []byte{0x00, 0x30, 0x15, 0x13, 0x25, 0x12, 0x24, 0x20},
		},
		{
			// A custom type is opaque, so only a raw value can be written to it.
			name: "a custom type takes the raw bytes", dataTypeId: typeIdCustom,
			value: spiValues.NewPlcRawByteArray([]byte{0x01, 0x02}), want: []byte{0x01, 0x02},
		},
		{
			name: "no value at all", dataTypeId: typeIdDint, value: nil, wantErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			data, err := encodeWriteValue(testCase.dataTypeId, testCase.value)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, data)
		})
	}
}

// A value which doesn't fit the symbol's type is refused rather than truncated: plc4j's
// value.getShort() and friends silently narrow, which turns a mistyped write into a wrong value in
// the PLC.
func TestEncodeWriteValue_RefusesValuesOutOfRange(t *testing.T) {
	tests := []struct {
		name       string
		dataTypeId uint16
		value      apiValues.PlcValue
	}{
		{name: "an INT can't hold 40000", dataTypeId: typeIdInt, value: spiValues.NewPlcDINT(40000)},
		{name: "a UINT can't hold -1", dataTypeId: typeIdUint, value: spiValues.NewPlcDINT(-1)},
		{name: "a BYTE can't hold 256", dataTypeId: typeIdByte, value: spiValues.NewPlcUINT(256)},
		{name: "a UDINT can't hold -1", dataTypeId: typeIdUdint, value: spiValues.NewPlcDINT(-1)},
		{name: "a REAL is not a string", dataTypeId: typeIdReal, value: spiValues.NewPlcSTRING("nope")},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := encodeWriteValue(testCase.dataTypeId, testCase.value)
			assert.Error(t, err)
		})
	}
}

// Whatever a read decodes has to be writable again unchanged, or a read-modify-write cycle would
// drift.
func TestValues_RoundTripThroughTheWire(t *testing.T) {
	tests := []struct {
		name       string
		dataTypeId uint16
		block      []byte
	}{
		{name: "BOOL", dataTypeId: typeIdBool, block: []byte{0x01}},
		{name: "BYTE", dataTypeId: typeIdByte, block: []byte{0xAB}},
		{name: "INT", dataTypeId: typeIdInt, block: []byte{0xFE, 0xFF}},
		{name: "UINT", dataTypeId: typeIdUint, block: []byte{0x34, 0x12}},
		{name: "WORD", dataTypeId: typeIdWord, block: []byte{0x34, 0x12}},
		{name: "DINT", dataTypeId: typeIdDint, block: []byte{0x78, 0x56, 0x34, 0x12}},
		{name: "UDINT", dataTypeId: typeIdUdint, block: []byte{0x78, 0x56, 0x34, 0x12}},
		{name: "DWORD", dataTypeId: typeIdDword, block: []byte{0x78, 0x56, 0x34, 0x12}},
		{name: "REAL", dataTypeId: typeIdReal, block: []byte{0x00, 0x00, 0x80, 0x3F}},
		{name: "TIME", dataTypeId: typeIdTime, block: []byte{0xE8, 0x03, 0x00, 0x00}},
		{name: "DATE", dataTypeId: typeIdDate, block: []byte{0x25, 0x12, 0x24, 0x20}},
		{name: "TOD", dataTypeId: typeIdTod, block: []byte{0x50, 0x30, 0x15, 0x13}},
		{name: "DATE_AND_TIME", dataTypeId: typeIdDateAndTime,
			block: []byte{0x00, 0x30, 0x15, 0x13, 0x25, 0x12, 0x24, 0x20}},
		{name: "STRING", dataTypeId: typeIdString, block: []byte{'h', 'i', 0x00}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			value, err := decodeReadResponse(t.Context(), testCase.dataTypeId, testCase.block)
			require.NoError(t, err)
			encoded, err := encodeWriteValue(testCase.dataTypeId, value)
			require.NoError(t, err)
			assert.Equal(t, testCase.block, encoded)
		})
	}
}

func TestBcdCodec(t *testing.T) {
	t.Run("a byte round trips", func(t *testing.T) {
		for value := uint8(0); value <= 99; value++ {
			encoded := encodeBcdByte(value)
			decoded, err := decodeBcdByte(encoded, "value")
			require.NoError(t, err)
			assert.Equal(t, value, decoded, "0x%02X", encoded)
		}
	})

	t.Run("25 is 0x25, not 0x19", func(t *testing.T) {
		assert.Equal(t, byte(0x25), encodeBcdByte(25))
	})

	t.Run("a nibble above 9 is not BCD", func(t *testing.T) {
		for _, notBcd := range []byte{0x0A, 0xA0, 0xFF, 0x1F} {
			_, err := decodeBcdByte(notBcd, "value")
			assert.Error(t, err, "0x%02X", notBcd)
		}
	})

	t.Run("a year has its low digits first", func(t *testing.T) {
		yearBytes, err := encodeBcdYear(2024)
		require.NoError(t, err)
		assert.Equal(t, [2]byte{0x24, 0x20}, yearBytes)
		year, err := decodeBcd16(yearBytes[0], yearBytes[1], "year")
		require.NoError(t, err)
		assert.Equal(t, uint16(2024), year)
	})

	t.Run("a year outside four digits is refused", func(t *testing.T) {
		_, err := encodeBcdYear(10000)
		assert.Error(t, err)
		_, err = encodeBcdYear(-1)
		assert.Error(t, err)
	})
}

// These are the numbers a read or write reference carries in its 4 bit dataSizeIndex field, and they
// have to agree with the byte counts the model derives from them - the PLC sizes the payload from the
// index, not from how many bytes we actually sent.
func TestDataSizeIndexAgreesWithThePayloadSize(t *testing.T) {
	tests := []struct {
		name       string
		dataTypeId uint16
		wantIndex  uint8
		wantBytes  uint16
	}{
		{name: "BOOL", dataTypeId: typeIdBool, wantIndex: 1, wantBytes: 1},
		{name: "BYTE", dataTypeId: typeIdByte, wantIndex: 1, wantBytes: 1},
		{name: "INT", dataTypeId: typeIdInt, wantIndex: 2, wantBytes: 2},
		{name: "WORD", dataTypeId: typeIdWord, wantIndex: 2, wantBytes: 2},
		{name: "DINT", dataTypeId: typeIdDint, wantIndex: 3, wantBytes: 4},
		{name: "REAL", dataTypeId: typeIdReal, wantIndex: 3, wantBytes: 4},
		{name: "TIME", dataTypeId: typeIdTime, wantIndex: 3, wantBytes: 4},
		{name: "DATE", dataTypeId: typeIdDate, wantIndex: 3, wantBytes: 4},
		{name: "TOD", dataTypeId: typeIdTod, wantIndex: 3, wantBytes: 4},
		{name: "DATE_AND_TIME", dataTypeId: typeIdDateAndTime, wantIndex: 4, wantBytes: 8},
		// A type the dictionary knows nothing about gets index 3, which is plc4j's hard-coded
		// fallback.
		{name: "a custom type", dataTypeId: typeIdCustom, wantIndex: 3, wantBytes: 4},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			index := dataSizeIndexFor(testCase.dataTypeId)
			assert.Equal(t, testCase.wantIndex, index)
			assert.Equal(t, testCase.wantBytes, readWriteModel.WriteSizeIndexToByteCount(t.Context(), index))
		})
	}
}

// STRING's request size is 17, which doesn't fit the 4 bit field - which is exactly why the driver
// reads and writes a STRING as a byte array instead.
func TestStringDoesNotFitTheSizeIndexField(t *testing.T) {
	assert.Equal(t, uint8(17), readWriteModel.UmasDataType_STRING.RequestSize())
	assert.Greater(t, readWriteModel.UmasDataType_STRING.RequestSize(), uint8(0x0F))
	assert.True(t, isStringType(typeIdString))
	assert.False(t, isStringType(typeIdDint))
	assert.False(t, isStringType(typeIdCustom))
}

func TestMapToPlcValueType(t *testing.T) {
	tests := []struct {
		dataTypeId uint16
		want       apiValues.PlcValueType
	}{
		{dataTypeId: typeIdBool, want: apiValues.BOOL},
		{dataTypeId: typeIdEbool, want: apiValues.BOOL},
		{dataTypeId: typeIdByte, want: apiValues.BYTE},
		{dataTypeId: typeIdUnknown11, want: apiValues.BYTE},
		{dataTypeId: typeIdInt, want: apiValues.INT},
		{dataTypeId: typeIdUint, want: apiValues.UINT},
		{dataTypeId: typeIdDint, want: apiValues.DINT},
		{dataTypeId: typeIdUdint, want: apiValues.UDINT},
		{dataTypeId: typeIdReal, want: apiValues.REAL},
		{dataTypeId: typeIdString, want: apiValues.STRING},
		{dataTypeId: typeIdTime, want: apiValues.TIME},
		{dataTypeId: typeIdDate, want: apiValues.DATE},
		{dataTypeId: typeIdTod, want: apiValues.TIME_OF_DAY},
		{dataTypeId: typeIdDateAndTime, want: apiValues.DATE_AND_TIME},
		{dataTypeId: typeIdWord, want: apiValues.WORD},
		{dataTypeId: typeIdDword, want: apiValues.DWORD},
		// Not a primitive, so nothing more specific can be said than "some bytes".
		{dataTypeId: typeIdCustom, want: apiValues.RAW_BYTE_ARRAY},
		{dataTypeId: 0, want: apiValues.RAW_BYTE_ARRAY},
		{dataTypeId: 1000, want: apiValues.RAW_BYTE_ARRAY},
	}
	for _, testCase := range tests {
		t.Run(testCase.want.String(), func(t *testing.T) {
			assert.Equal(t, testCase.want, mapToPlcValueType(testCase.dataTypeId))
		})
	}
}

// A value of the wrong shape has to be refused, not guessed at. Both cases here used to be silent:
// PlcValueAdapter.GetString panics, which the recover in Writer.Write turns into a stack trace for
// the whole request instead of a response code for the one tag, and PlcValueAdapter.GetBool answers
// false, which would have written a 0x00 nobody asked for to a PLC output and reported OK.
//
// A []byte handed to a write is the reachable case for both: a UMAS tag has no value type, so the
// value handler wraps it as a PlcRawByteArray rather than coercing it to the symbol's type.
func TestEncodeWriteValue_RefusesValuesOfTheWrongShape(t *testing.T) {
	tests := []struct {
		name       string
		dataTypeId uint16
		value      apiValues.PlcValue
	}{
		{name: "a byte array is not a STRING", dataTypeId: typeIdString, value: spiValues.NewPlcRawByteArray([]byte("hi"))},
		{name: "a null is not a STRING", dataTypeId: typeIdString, value: spiValues.NewPlcNULL()},
		{name: "a date is not a STRING", dataTypeId: typeIdString, value: spiValues.NewPlcDATE(time.Now())},
		{name: "a byte array is not a BOOL", dataTypeId: typeIdBool, value: spiValues.NewPlcRawByteArray([]byte{0x01})},
		{name: "a null is not a BOOL", dataTypeId: typeIdBool, value: spiValues.NewPlcNULL()},
		{name: "a date is not a BOOL", dataTypeId: typeIdBool, value: spiValues.NewPlcDATE(time.Now())},
		{name: "a byte array is not an EBOOL", dataTypeId: typeIdEbool, value: spiValues.NewPlcRawByteArray([]byte{0x01})},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Not just an error: a panic here would take a whole write request down with it.
			require.NotPanics(t, func() {
				_, err := encodeWriteValue(testCase.dataTypeId, testCase.value)
				assert.Error(t, err)
			})
		})
	}
}

// A number written to a BOOL symbol is "not zero", which is what plc4j's serializeForType gets out
// of value.getBoolean() and what the Go numeric values answer to GetBool - they just don't claim
// IsBool, so the shape check can't be IsBool alone.
func TestEncodeWriteValue_TakesANumberAsABool(t *testing.T) {
	tests := []struct {
		name  string
		value apiValues.PlcValue
		want  []byte
	}{
		{name: "a one is true", value: spiValues.NewPlcLINT(1), want: []byte{0x01}},
		{name: "a zero is false", value: spiValues.NewPlcLINT(0), want: []byte{0x00}},
		{name: "a negative number is true", value: spiValues.NewPlcDINT(-1), want: []byte{0x01}},
		{name: "an unsigned number is true", value: spiValues.NewPlcUSINT(2), want: []byte{0x01}},
		// PlcDWORD is built on the plain value adapter, so all of its Is<Integer> predicates answer
		// false while its getters work fine - the same reason asUnsigned asks it by value type.
		{name: "a DWORD is asked by value type", value: spiValues.NewPlcDWORD(0), want: []byte{0x00}},
		{name: "a non zero DWORD is true", value: spiValues.NewPlcDWORD(7), want: []byte{0x01}},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			data, err := encodeWriteValue(typeIdBool, testCase.value)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, data)
		})
	}
}
