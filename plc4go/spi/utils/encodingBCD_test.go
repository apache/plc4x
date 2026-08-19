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

package utils

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// bcdArg is the exact option the generated code emits for `encoding='"BCD"'`.
var bcdArg = WithEncoding("BCD")

// nonBcdEncodings are all the encoding values that must NOT change behaviour.
// "bcd" and "BCD " are in here on purpose: the match is an exact string compare,
// mirroring the literal the mspec compiler emits.
var nonBcdEncodings = []string{"", "UTF8", "UTF16", "UTF16BE", "UTF16LE", "ASCII", "ISO-8859-1", "IEEE754", "unsigned-binary", "bcd", "Bcd", "BCD "}

///////////////////////////////////////////////////////////////////////////////
// pure codec level
///////////////////////////////////////////////////////////////////////////////

func TestEncodeDecodeBCD_RoundTripAllWidths(t *testing.T) {
	for digits := 1; digits <= maxBCDDigits; digits++ {
		bitLength := uint8(digits * 4)
		max := bcdPow10[digits] - 1
		values := []uint64{0, 1, 9, max}
		if digits > 1 {
			values = append(values, bcdPow10[digits-1], max/2, max-1)
		}
		for _, value := range values {
			t.Run(fmt.Sprintf("%dbits/%d", bitLength, value), func(t *testing.T) {
				raw, err := encodeBCD(value, bitLength)
				require.NoError(t, err)
				// every nibble of the packed representation must be a valid digit
				for i := 0; i < digits; i++ {
					assert.LessOrEqual(t, (raw>>uint(4*i))&0xF, uint64(9), "nibble %d of %#x is not a decimal digit", i, raw)
				}
				decoded, err := decodeBCD(raw, bitLength)
				require.NoError(t, err)
				assert.Equal(t, value, decoded)
			})
		}
	}
}

func TestEncodeBCD_PackedLayout(t *testing.T) {
	// most significant digit lands in the most significant nibble
	tests := []struct {
		value     uint64
		bitLength uint8
		want      uint64
	}{
		{value: 0, bitLength: 4, want: 0x0},
		{value: 9, bitLength: 4, want: 0x9},
		{value: 0, bitLength: 8, want: 0x00},
		{value: 7, bitLength: 8, want: 0x07},
		{value: 24, bitLength: 8, want: 0x24},
		{value: 99, bitLength: 8, want: 0x99},
		{value: 123, bitLength: 12, want: 0x123},
		{value: 999, bitLength: 12, want: 0x999},
		{value: 1234, bitLength: 16, want: 0x1234},
		{value: 12345678, bitLength: 32, want: 0x12345678},
		{value: 1234567890123456, bitLength: 64, want: 0x1234567890123456},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d@%d", tt.value, tt.bitLength), func(t *testing.T) {
			raw, err := encodeBCD(tt.value, tt.bitLength)
			require.NoError(t, err)
			assert.Equal(t, tt.want, raw)
			decoded, err := decodeBCD(tt.want, tt.bitLength)
			require.NoError(t, err)
			assert.Equal(t, tt.value, decoded)
		})
	}
}

func TestBCD_BitLengthMustBeMultipleOfFour(t *testing.T) {
	for _, bitLength := range []uint8{1, 2, 3, 5, 6, 7, 9, 11, 13, 15, 63} {
		t.Run(fmt.Sprintf("%dbits", bitLength), func(t *testing.T) {
			_, err := encodeBCD(0, bitLength)
			assert.ErrorContains(t, err, "must be a multiple of 4")
			_, err = decodeBCD(0, bitLength)
			assert.ErrorContains(t, err, "must be a multiple of 4")
		})
	}
}

func TestBCD_BitLengthBeyondUint64(t *testing.T) {
	// 17 digits no longer fit into the uint64 the bit buffers exchange
	_, err := encodeBCD(0, 68)
	assert.ErrorContains(t, err, "exceeds")
	_, err = decodeBCD(0, 68)
	assert.ErrorContains(t, err, "exceeds")
}

func TestDecodeBCD_InvalidNibble(t *testing.T) {
	tests := []struct {
		name      string
		raw       uint64
		bitLength uint8
	}{
		{name: "0xA in the only nibble", raw: 0xA, bitLength: 4},
		{name: "0xF in the only nibble", raw: 0xF, bitLength: 4},
		{name: "0xA in the high nibble", raw: 0xA1, bitLength: 8},
		{name: "0xA in the low nibble", raw: 0x1A, bitLength: 8},
		{name: "0xB in the middle nibble", raw: 0x1B3, bitLength: 12},
		{name: "0xE at the very end", raw: 0x123456789012345E, bitLength: 64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decodeBCD(tt.raw, tt.bitLength)
			assert.ErrorContains(t, err, "invalid BCD digit")
		})
	}
}

func TestDecodeBCD_AllNibbleValues(t *testing.T) {
	// exhaustive over one byte: 0x00..0x99 with both nibbles <= 9 decode, all else fail
	for raw := 0; raw <= 0xFF; raw++ {
		high, low := raw>>4, raw&0xF
		decoded, err := decodeBCD(uint64(raw), 8)
		if high <= 9 && low <= 9 {
			require.NoError(t, err, "%#02x should decode", raw)
			assert.Equal(t, uint64(high*10+low), decoded)
		} else {
			assert.Error(t, err, "%#02x should not decode", raw)
		}
	}
}

func TestEncodeBCD_ValueTooLarge(t *testing.T) {
	tests := []struct {
		value     uint64
		bitLength uint8
	}{
		{value: 10, bitLength: 4},
		{value: 100, bitLength: 8},
		{value: 1000, bitLength: 12},
		{value: 10000, bitLength: 16},
		{value: math.MaxUint64, bitLength: 64},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d@%d", tt.value, tt.bitLength), func(t *testing.T) {
			_, err := encodeBCD(tt.value, tt.bitLength)
			assert.ErrorContains(t, err, "cannot be encoded in")
		})
	}
}

func TestEncodeBCDSigned_RejectsNegative(t *testing.T) {
	for _, value := range []int64{-1, -9, math.MinInt64} {
		_, err := encodeBCDSigned(value, 8)
		assert.ErrorContains(t, err, "non-negative")
	}
	// non-negative still works
	raw, err := encodeBCDSigned(42, 8)
	require.NoError(t, err)
	assert.Equal(t, uint64(0x42), raw)
}

func TestDecodeBCDBounded(t *testing.T) {
	// 12 BCD digits can hold 999, which does not fit a uint8
	_, err := decodeBCDBounded(0x999, 12, math.MaxUint8)
	assert.ErrorContains(t, err, "too large for the target type")
	value, err := decodeBCDBounded(0x255, 12, math.MaxUint8)
	require.NoError(t, err)
	assert.Equal(t, uint64(255), value)
}

///////////////////////////////////////////////////////////////////////////////
// arg selection
///////////////////////////////////////////////////////////////////////////////

func TestArgsSelectBCD(t *testing.T) {
	tests := []struct {
		name string
		args []WithReaderWriterArgs
		want bool
	}{
		{name: "no args", args: nil, want: false},
		{name: "BCD", args: []WithReaderWriterArgs{WithEncoding("BCD")}, want: true},
		{name: "UTF8", args: []WithReaderWriterArgs{WithEncoding("UTF8")}, want: false},
		{name: "empty encoding", args: []WithReaderWriterArgs{WithEncoding("")}, want: false},
		{name: "lower case bcd", args: []WithReaderWriterArgs{WithEncoding("bcd")}, want: false},
		{name: "trailing space", args: []WithReaderWriterArgs{WithEncoding("BCD ")}, want: false},
		{name: "unrelated arg only", args: []WithReaderWriterArgs{WithRenderAsList(true)}, want: false},
		{name: "unrelated arg then BCD", args: []WithReaderWriterArgs{WithRenderAsList(true), WithEncoding("BCD")}, want: true},
		{name: "custom option then BCD", args: []WithReaderWriterArgs{NewCustomOption(), WithEncoding("BCD")}, want: true},
		{name: "additional string repr then UTF8", args: []WithReaderWriterArgs{WithAdditionalStringRepresentation("x"), WithEncoding("UTF8")}, want: false},
		// Two encoding args: the FIRST one decides, exactly as in
		// BufferCommons.ExtractEncoding. Nothing in the tree emits two encoding
		// args, but the two implementations must not disagree if it ever happens.
		{name: "UTF8 then BCD", args: []WithReaderWriterArgs{WithEncoding("UTF8"), WithEncoding("BCD")}, want: false},
		{name: "BCD then UTF8", args: []WithReaderWriterArgs{WithEncoding("BCD"), WithEncoding("UTF8")}, want: true},
		{name: "empty encoding then BCD", args: []WithReaderWriterArgs{WithEncoding(""), WithEncoding("BCD")}, want: false},
		{name: "BCD then BCD", args: []WithReaderWriterArgs{WithEncoding("BCD"), WithEncoding("BCD")}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readerArgs := make([]WithReaderArgs, len(tt.args))
			writerArgs := make([]WithWriterArgs, len(tt.args))
			for i, arg := range tt.args {
				readerArgs[i] = arg
				writerArgs[i] = arg
			}
			assert.Equal(t, tt.want, readerArgsSelectBCD(readerArgs), "readerArgsSelectBCD")
			assert.Equal(t, tt.want, writerArgsSelectBCD(writerArgs), "writerArgsSelectBCD")
			// the upcast (nested) form must agree with the plain one
			assert.Equal(t, tt.want, readerArgsSelectBCD(toReaderArgs(UpcastReaderArgs(readerArgs...))), "readerArgsSelectBCD (upcast)")
			assert.Equal(t, tt.want, writerArgsSelectBCD(toWriterArgs(UpcastWriterArgs(writerArgs...))), "writerArgsSelectBCD (upcast)")
		})
	}
}

// TestArgsSelectBCD_AgreesWithExtractEncoding pins that the allocation free fast
// path picks the same encoding as the canonical BufferCommons.ExtractEncoding.
func TestArgsSelectBCD_AgreesWithExtractEncoding(t *testing.T) {
	var commons BufferCommons
	for _, encoding := range append([]string{"BCD"}, nonBcdEncodings...) {
		t.Run(encoding, func(t *testing.T) {
			arg := WithEncoding(encoding)
			want := commons.ExtractEncoding(arg) == "BCD"
			assert.Equal(t, want, readerArgsSelectBCD([]WithReaderArgs{arg}))
			assert.Equal(t, want, writerArgsSelectBCD([]WithWriterArgs{arg}))
		})
	}
}

// TestArgsSelectBCD_AgreesWithExtractEncodingOnMultipleEncodings is the
// multi-arg half of the equivalence above.
//
// Both implementations implement "the first arg that carries an encoding
// decides" and both do so for the direct withEncoding case AND for the nested
// readerWriterArg case, so a leading non-BCD encoding shadows a trailing BCD one
// either way. The nested case is exercised explicitly because it cannot be
// reached through UpcastReaderArgs/UpcastWriterArgs (withEncoding embeds
// readerWriterArg, hence already satisfies WithReaderWriterArgs and is passed
// through unwrapped) - it is only reachable by a hand built arg, which is
// precisely why it needs a test rather than the reader's trust.
func TestArgsSelectBCD_AgreesWithExtractEncodingOnMultipleEncodings(t *testing.T) {
	var commons BufferCommons
	nest := func(encoding string) WithReaderWriterArgs {
		return readerWriterArg{WithEncoding(encoding), writerArg{}}
	}
	nestWriter := func(encoding string) WithReaderWriterArgs {
		return readerWriterArg{readerArg{}, WithEncoding(encoding)}
	}
	for _, first := range append([]string{"BCD"}, nonBcdEncodings...) {
		for _, second := range []string{"BCD", "UTF8"} {
			for shape, build := range map[string]func(string) WithReaderWriterArgs{
				"direct":       WithEncoding,
				"nestedReader": nest,
				"nestedWriter": nestWriter,
			} {
				t.Run(fmt.Sprintf("%s/%q_then_%q", shape, first, second), func(t *testing.T) {
					args := []WithReaderWriterArgs{build(first), WithEncoding(second)}
					want := commons.ExtractEncoding(args...) == "BCD"
					assert.Equal(t, want, readerArgsSelectBCD(toReaderArgs(args)), "readerArgsSelectBCD")
					assert.Equal(t, want, writerArgsSelectBCD(toWriterArgs(args)), "writerArgsSelectBCD")
				})
			}
		}
	}
}

func toReaderArgs(args []WithReaderWriterArgs) []WithReaderArgs {
	result := make([]WithReaderArgs, len(args))
	for i, arg := range args {
		result[i] = arg
	}
	return result
}

func toWriterArgs(args []WithReaderWriterArgs) []WithWriterArgs {
	result := make([]WithWriterArgs, len(args))
	for i, arg := range args {
		result[i] = arg
	}
	return result
}

///////////////////////////////////////////////////////////////////////////////
// buffer level: reads
///////////////////////////////////////////////////////////////////////////////

func TestReadBuffer_BCD_Uint(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		bitLength uint8
		want      uint64
	}{
		{name: "single digit 0", data: []byte{0x00}, bitLength: 4, want: 0},
		{name: "single digit 9", data: []byte{0x90}, bitLength: 4, want: 9},
		{name: "two digits", data: []byte{0x24}, bitLength: 8, want: 24},
		{name: "two digits max", data: []byte{0x99}, bitLength: 8, want: 99},
		{name: "three digits", data: []byte{0x12, 0x30}, bitLength: 12, want: 123},
		{name: "four digits", data: []byte{0x20, 0x24}, bitLength: 16, want: 2024},
		{name: "eight digits", data: []byte{0x20, 0x24, 0x08, 0x19}, bitLength: 32, want: 20240819},
		{name: "sixteen digits", data: []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56}, bitLength: 64, want: 1234567890123456},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want <= math.MaxUint8 {
				rb := NewReadBufferByteBased(tt.data)
				got, err := rb.ReadUint8("test", tt.bitLength, bcdArg)
				require.NoError(t, err)
				assert.Equal(t, uint8(tt.want), got)
			}
			if tt.want <= math.MaxUint16 {
				rb := NewReadBufferByteBased(tt.data)
				got, err := rb.ReadUint16("test", tt.bitLength, bcdArg)
				require.NoError(t, err)
				assert.Equal(t, uint16(tt.want), got)
			}
			if tt.want <= math.MaxUint32 {
				rb := NewReadBufferByteBased(tt.data)
				got, err := rb.ReadUint32("test", tt.bitLength, bcdArg)
				require.NoError(t, err)
				assert.Equal(t, uint32(tt.want), got)
			}
			rb := NewReadBufferByteBased(tt.data)
			got, err := rb.ReadUint64("test", tt.bitLength, bcdArg)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReadBuffer_BCD_Int(t *testing.T) {
	rb := NewReadBufferByteBased([]byte{0x42})
	got8, err := rb.ReadInt8("test", 8, bcdArg)
	require.NoError(t, err)
	assert.Equal(t, int8(42), got8)

	rb = NewReadBufferByteBased([]byte{0x12, 0x34})
	got16, err := rb.ReadInt16("test", 16, bcdArg)
	require.NoError(t, err)
	assert.Equal(t, int16(1234), got16)

	rb = NewReadBufferByteBased([]byte{0x12, 0x34, 0x56, 0x78})
	got32, err := rb.ReadInt32("test", 32, bcdArg)
	require.NoError(t, err)
	assert.Equal(t, int32(12345678), got32)

	rb = NewReadBufferByteBased([]byte{0x00, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56})
	got64, err := rb.ReadInt64("test", 64, bcdArg)
	require.NoError(t, err)
	assert.Equal(t, int64(34567890123456), got64)
}

// TestReadBuffer_BCD_SignedNarrowingMatchesPlc4j pins the divergence-free
// behaviour for signed BCD fields wider than their target type.
//
// plc4j's EncodingBCD.decodeByte bounds the decoded value at 255 and then
// narrows with a (byte) cast, and decodeShort does the same at 65535 with a
// (short) cast - so a 12 bit signed-byte field holding 0x200 yields
// (byte) 200 == -56 rather than an error. No mspec in the tree declares a signed
// BCD field today, but the read buffer must not invent a stricter rule than the
// reference implementation if one ever appears.
func TestReadBuffer_BCD_SignedNarrowingMatchesPlc4j(t *testing.T) {
	// 12 bit field holding decimal 200: fits a byte only after the cast.
	rb := NewReadBufferByteBased([]byte{0x20, 0x00})
	got8, err := rb.ReadInt8("test", 12, bcdArg)
	require.NoError(t, err)
	assert.Equal(t, int8(-56), got8, "plc4j yields (byte) 200 == -56")

	// 999 exceeds 255 even unsigned, so plc4j throws and so must this.
	rb = NewReadBufferByteBased([]byte{0x99, 0x90})
	_, err = rb.ReadInt8("test", 12, bcdArg)
	assert.ErrorContains(t, err, "too large for the target type")

	// 20 bit field holding decimal 40000: fits a short only after the cast.
	rb = NewReadBufferByteBased([]byte{0x40, 0x00, 0x00})
	got16, err := rb.ReadInt16("test", 20, bcdArg)
	require.NoError(t, err)
	assert.Equal(t, int16(-25536), got16, "plc4j yields (short) 40000 == -25536")

	// 99999 exceeds 65535 even unsigned.
	rb = NewReadBufferByteBased([]byte{0x99, 0x99, 0x90})
	_, err = rb.ReadInt16("test", 20, bcdArg)
	assert.ErrorContains(t, err, "too large for the target type")
}

func TestReadBuffer_BCD_InvalidNibble(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		bitLength uint8
	}{
		{name: "0xA high nibble", data: []byte{0xA1}, bitLength: 8},
		{name: "0xA low nibble", data: []byte{0x1A}, bitLength: 8},
		{name: "0xF single nibble", data: []byte{0xF0}, bitLength: 4},
		{name: "0xB in a 12 bit field", data: []byte{0x1B, 0x30}, bitLength: 12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := NewReadBufferByteBased(tt.data)
			_, err := rb.ReadUint8("test", tt.bitLength, bcdArg)
			assert.ErrorContains(t, err, "invalid BCD digit")

			rb = NewReadBufferByteBased(tt.data)
			_, err = rb.ReadUint16("test", tt.bitLength, bcdArg)
			assert.ErrorContains(t, err, "invalid BCD digit")

			rb = NewReadBufferByteBased(tt.data)
			_, err = rb.ReadUint32("test", tt.bitLength, bcdArg)
			assert.ErrorContains(t, err, "invalid BCD digit")

			rb = NewReadBufferByteBased(tt.data)
			_, err = rb.ReadUint64("test", tt.bitLength, bcdArg)
			assert.ErrorContains(t, err, "invalid BCD digit")
		})
	}
}

func TestReadBuffer_BCD_BitLengthNotMultipleOfFour(t *testing.T) {
	for _, bitLength := range []uint8{1, 3, 5, 7} {
		t.Run(fmt.Sprintf("%dbits", bitLength), func(t *testing.T) {
			rb := NewReadBufferByteBased([]byte{0x12, 0x34})
			_, err := rb.ReadUint8("test", bitLength, bcdArg)
			assert.ErrorContains(t, err, "must be a multiple of 4")
		})
	}
}

func TestReadBuffer_BCD_DecodedValueTooLargeForTarget(t *testing.T) {
	// 999 does not fit a uint8
	rb := NewReadBufferByteBased([]byte{0x99, 0x90})
	_, err := rb.ReadUint8("test", 12, bcdArg)
	assert.ErrorContains(t, err, "too large for the target type")

	// but it does fit a uint16
	rb = NewReadBufferByteBased([]byte{0x99, 0x90})
	got, err := rb.ReadUint16("test", 12, bcdArg)
	require.NoError(t, err)
	assert.Equal(t, uint16(999), got)
}

func TestReadBuffer_BCD_ZeroBitLengthIsUntouched(t *testing.T) {
	rb := NewReadBufferByteBased([]byte{0x24})
	got, err := rb.ReadUint8("test", 0, bcdArg)
	require.NoError(t, err)
	assert.Equal(t, uint8(0), got)
	// nothing consumed
	assert.Equal(t, uint32(0), rb.GetPos())
}

// TestReadBuffer_BCD_IgnoresLittleEndian pins that the BCD path operates on the
// raw MSB-first nibble stream regardless of the configured byte order.
func TestReadBuffer_BCD_IgnoresLittleEndian(t *testing.T) {
	for _, byteOrder := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		t.Run(fmt.Sprintf("%v", byteOrder), func(t *testing.T) {
			rb := NewReadBufferByteBased([]byte{0x20, 0x24}, WithByteOrderForReadBufferByteBased(byteOrder))
			got, err := rb.ReadUint16("test", 16, bcdArg)
			require.NoError(t, err)
			assert.Equal(t, uint16(2024), got)
		})
	}
}

func TestReadBuffer_BCD_SequentialUnalignedFields(t *testing.T) {
	// the S7 DATE_AND_TIME wire layout: 6 byte-aligned fields, a 12 bit field and
	// a 4 bit field
	rb := NewReadBufferByteBased([]byte{0x24, 0x08, 0x19, 0x10, 0x30, 0x45, 0x12, 0x32})
	year, err := rb.ReadUint8("year", 8, bcdArg)
	require.NoError(t, err)
	month, err := rb.ReadUint8("month", 8, bcdArg)
	require.NoError(t, err)
	day, err := rb.ReadUint8("day", 8, bcdArg)
	require.NoError(t, err)
	hour, err := rb.ReadUint8("hour", 8, bcdArg)
	require.NoError(t, err)
	minutes, err := rb.ReadUint8("minutes", 8, bcdArg)
	require.NoError(t, err)
	seconds, err := rb.ReadUint8("seconds", 8, bcdArg)
	require.NoError(t, err)
	msec, err := rb.ReadUint16("msec", 12, bcdArg)
	require.NoError(t, err)
	dow, err := rb.ReadUint8("dow", 4, bcdArg)
	require.NoError(t, err)

	assert.Equal(t, uint8(24), year)
	assert.Equal(t, uint8(8), month)
	assert.Equal(t, uint8(19), day)
	assert.Equal(t, uint8(10), hour)
	assert.Equal(t, uint8(30), minutes)
	assert.Equal(t, uint8(45), seconds)
	assert.Equal(t, uint16(123), msec)
	assert.Equal(t, uint8(2), dow)
}

///////////////////////////////////////////////////////////////////////////////
// buffer level: writes
///////////////////////////////////////////////////////////////////////////////

func TestWriteBuffer_BCD_Uint(t *testing.T) {
	tests := []struct {
		name      string
		value     uint64
		bitLength uint8
		want      []byte
	}{
		{name: "single digit", value: 9, bitLength: 4, want: []byte{0x90}},
		{name: "two digits", value: 24, bitLength: 8, want: []byte{0x24}},
		{name: "three digits", value: 123, bitLength: 12, want: []byte{0x12, 0x30}},
		{name: "three digits max", value: 999, bitLength: 12, want: []byte{0x99, 0x90}},
		{name: "four digits", value: 2024, bitLength: 16, want: []byte{0x20, 0x24}},
		{name: "eight digits", value: 20240819, bitLength: 32, want: []byte{0x20, 0x24, 0x08, 0x19}},
		{name: "sixteen digits", value: 1234567890123456, bitLength: 64, want: []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value <= math.MaxUint8 {
				wb := NewWriteBufferByteBased()
				require.NoError(t, wb.WriteUint8("test", tt.bitLength, uint8(tt.value), bcdArg))
				assert.Equal(t, tt.want, wb.GetBytes())
			}
			if tt.value <= math.MaxUint16 {
				wb := NewWriteBufferByteBased()
				require.NoError(t, wb.WriteUint16("test", tt.bitLength, uint16(tt.value), bcdArg))
				assert.Equal(t, tt.want, wb.GetBytes())
			}
			if tt.value <= math.MaxUint32 {
				wb := NewWriteBufferByteBased()
				require.NoError(t, wb.WriteUint32("test", tt.bitLength, uint32(tt.value), bcdArg))
				assert.Equal(t, tt.want, wb.GetBytes())
			}
			wb := NewWriteBufferByteBased()
			require.NoError(t, wb.WriteUint64("test", tt.bitLength, tt.value, bcdArg))
			assert.Equal(t, tt.want, wb.GetBytes())
		})
	}
}

func TestWriteBuffer_BCD_Int(t *testing.T) {
	wb := NewWriteBufferByteBased()
	require.NoError(t, wb.WriteInt8("test", 8, 42, bcdArg))
	assert.Equal(t, []byte{0x42}, wb.GetBytes())

	wb = NewWriteBufferByteBased()
	require.NoError(t, wb.WriteInt16("test", 16, 1234, bcdArg))
	assert.Equal(t, []byte{0x12, 0x34}, wb.GetBytes())

	wb = NewWriteBufferByteBased()
	require.NoError(t, wb.WriteInt32("test", 32, 12345678, bcdArg))
	assert.Equal(t, []byte{0x12, 0x34, 0x56, 0x78}, wb.GetBytes())

	wb = NewWriteBufferByteBased()
	require.NoError(t, wb.WriteInt64("test", 64, 34567890123456, bcdArg))
	assert.Equal(t, []byte{0x00, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56}, wb.GetBytes())
}

func TestWriteBuffer_BCD_RejectsNegative(t *testing.T) {
	wb := NewWriteBufferByteBased()
	assert.ErrorContains(t, wb.WriteInt8("test", 8, -1, bcdArg), "non-negative")
	assert.ErrorContains(t, wb.WriteInt16("test", 16, -1, bcdArg), "non-negative")
	assert.ErrorContains(t, wb.WriteInt32("test", 32, -1, bcdArg), "non-negative")
	assert.ErrorContains(t, wb.WriteInt64("test", 64, -1, bcdArg), "non-negative")
	// a rejected value must leave the buffer untouched
	assert.Equal(t, uint32(0), wb.GetPos())
	assert.Empty(t, wb.GetBytes())
}

func TestWriteBuffer_BCD_RejectsOverflow(t *testing.T) {
	wb := NewWriteBufferByteBased()
	assert.ErrorContains(t, wb.WriteUint8("test", 4, 10, bcdArg), "cannot be encoded in")
	assert.ErrorContains(t, wb.WriteUint8("test", 8, 100, bcdArg), "cannot be encoded in")
	assert.ErrorContains(t, wb.WriteUint16("test", 12, 1000, bcdArg), "cannot be encoded in")
	assert.ErrorContains(t, wb.WriteUint16("test", 16, 10000, bcdArg), "cannot be encoded in")
	assert.ErrorContains(t, wb.WriteUint32("test", 16, 10000, bcdArg), "cannot be encoded in")
	assert.ErrorContains(t, wb.WriteUint64("test", 64, math.MaxUint64, bcdArg), "cannot be encoded in")
	// a rejected value must leave the buffer untouched
	assert.Equal(t, uint32(0), wb.GetPos())
	assert.Empty(t, wb.GetBytes())
}

func TestWriteBuffer_BCD_BitLengthNotMultipleOfFour(t *testing.T) {
	for _, bitLength := range []uint8{1, 3, 5, 7} {
		t.Run(fmt.Sprintf("%dbits", bitLength), func(t *testing.T) {
			wb := NewWriteBufferByteBased()
			assert.ErrorContains(t, wb.WriteUint8("test", bitLength, 1, bcdArg), "must be a multiple of 4")
			assert.Equal(t, uint32(0), wb.GetPos())
		})
	}
}

func TestWriteBuffer_BCD_ZeroBitLengthIsUntouched(t *testing.T) {
	wb := NewWriteBufferByteBased()
	require.NoError(t, wb.WriteUint8("test", 0, 24, bcdArg))
	assert.Equal(t, uint32(0), wb.GetPos())
	assert.Empty(t, wb.GetBytes())
}

// TestWriteBuffer_BCD_IgnoresLittleEndian is the write twin of
// TestReadBuffer_BCD_IgnoresLittleEndian.
func TestWriteBuffer_BCD_IgnoresLittleEndian(t *testing.T) {
	for _, byteOrder := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		t.Run(fmt.Sprintf("%v", byteOrder), func(t *testing.T) {
			wb := NewWriteBufferByteBased(WithByteOrderForByteBasedBuffer(byteOrder))
			require.NoError(t, wb.WriteUint16("test", 16, 2024, bcdArg))
			assert.Equal(t, []byte{0x20, 0x24}, wb.GetBytes())
		})
	}
}

func TestWriteBuffer_BCD_SequentialUnalignedFields(t *testing.T) {
	wb := NewWriteBufferByteBased()
	require.NoError(t, wb.WriteUint8("year", 8, 24, bcdArg))
	require.NoError(t, wb.WriteUint8("month", 8, 8, bcdArg))
	require.NoError(t, wb.WriteUint8("day", 8, 19, bcdArg))
	require.NoError(t, wb.WriteUint8("hour", 8, 10, bcdArg))
	require.NoError(t, wb.WriteUint8("minutes", 8, 30, bcdArg))
	require.NoError(t, wb.WriteUint8("seconds", 8, 45, bcdArg))
	require.NoError(t, wb.WriteUint16("msec", 12, 123, bcdArg))
	require.NoError(t, wb.WriteUint8("dow", 4, 2, bcdArg))
	assert.Equal(t, []byte{0x24, 0x08, 0x19, 0x10, 0x30, 0x45, 0x12, 0x32}, wb.GetBytes())
}

///////////////////////////////////////////////////////////////////////////////
// read/write symmetry
///////////////////////////////////////////////////////////////////////////////

// TestBCD_BufferRoundTripAllWidthsSigned is the signed twin of
// TestBCD_BufferRoundTripAllWidths: every value the signed reads can produce
// WITHOUT overflowing the target type must be writable back to the same bytes.
func TestBCD_BufferRoundTripAllWidthsSigned(t *testing.T) {
	for digits := 1; digits <= maxBCDDigits; digits++ {
		bitLength := uint8(digits * 4)
		max := int64(bcdPow10[digits] - 1)
		if max > math.MaxInt64 {
			max = math.MaxInt64
		}
		for _, value := range []int64{0, 1, max / 2, max} {
			t.Run(fmt.Sprintf("%dbits/%d", bitLength, value), func(t *testing.T) {
				wb := NewWriteBufferByteBased()
				require.NoError(t, wb.WriteInt64("test", bitLength, value, bcdArg))
				rb := NewReadBufferByteBased(wb.GetBytes())
				got, err := rb.ReadInt64("test", bitLength, bcdArg)
				require.NoError(t, err)
				assert.Equal(t, value, got)
			})
		}
	}
}

// TestBCD_SignedNarrowingIsNotRoundTrippable pins a deliberate asymmetry that
// exists for plc4j parity and is otherwise easy to "fix" into a divergence.
//
// The signed reads bound the decode at the UNSIGNED maximum of the target type
// and then narrow with a Go cast, because that is what
// EncodingBCD.decodeByte/decodeShort do (they compare against 255/65535 and cast
// to (byte)/(short)). So a BCD field wider than its target type can decode to a
// negative number, which encodeBCDSigned then refuses to write back, since
// EncodingBCD.encodeInt rejects negative values outright.
//
// Both halves match plc4j, so the asymmetry is parity rather than a Go bug, and
// it is unreachable in practice: every BCD field in the tree is unsigned (14 in
// s7.mspec, 13 in umas.mspec). If a signed BCD field ever appears, this test is
// the place that says the behaviour was chosen, not stumbled into.
func TestBCD_SignedNarrowingIsNotRoundTrippable(t *testing.T) {
	tests := []struct {
		name      string
		bitLength uint8
		wire      []byte
		read      func(ReadBufferByteBased, uint8) (int64, error)
		write     func(WriteBufferByteBased, uint8, int64) error
		want      int64
	}{
		{
			name:      "int8 over 3 BCD digits",
			bitLength: 12,
			wire:      []byte{0x20, 0x00}, // digits 2,0,0 -> 200
			read: func(rb ReadBufferByteBased, bitLength uint8) (int64, error) {
				value, err := rb.ReadInt8("test", bitLength, bcdArg)
				return int64(value), err
			},
			write: func(wb WriteBufferByteBased, bitLength uint8, value int64) error {
				return wb.WriteInt8("test", bitLength, int8(value), bcdArg)
			},
			want: -56, // int8(200)
		},
		{
			name:      "int16 over 5 BCD digits",
			bitLength: 20,
			wire:      []byte{0x40, 0x00, 0x00}, // digits 4,0,0,0,0 -> 40000
			read: func(rb ReadBufferByteBased, bitLength uint8) (int64, error) {
				value, err := rb.ReadInt16("test", bitLength, bcdArg)
				return int64(value), err
			},
			write: func(wb WriteBufferByteBased, bitLength uint8, value int64) error {
				return wb.WriteInt16("test", bitLength, int16(value), bcdArg)
			},
			want: -25536, // int16(40000)
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := NewReadBufferByteBased(tt.wire)
			got, err := tt.read(rb, tt.bitLength)
			require.NoError(t, err, "the decode is bounded at the unsigned max, so it succeeds")
			require.Negative(t, tt.want, "the fixture must actually wrap")
			assert.Equal(t, tt.want, got)

			wb := NewWriteBufferByteBased()
			assert.ErrorContains(t, tt.write(wb, tt.bitLength, got), "non-negative",
				"writing the wrapped value back must fail loudly rather than emit different bytes")
			assert.Empty(t, wb.GetBytes(), "a rejected value must leave the buffer untouched")
		})
	}
}

func TestBCD_BufferRoundTripAllWidths(t *testing.T) {
	for digits := 1; digits <= maxBCDDigits; digits++ {
		bitLength := uint8(digits * 4)
		max := bcdPow10[digits] - 1
		for _, value := range []uint64{0, 1, max / 2, max} {
			t.Run(fmt.Sprintf("%dbits/%d", bitLength, value), func(t *testing.T) {
				wb := NewWriteBufferByteBased()
				require.NoError(t, wb.WriteUint64("test", bitLength, value, bcdArg))
				rb := NewReadBufferByteBased(wb.GetBytes())
				got, err := rb.ReadUint64("test", bitLength, bcdArg)
				require.NoError(t, err)
				assert.Equal(t, value, got)
			})
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// non-regression: anything that is not exactly "BCD" must behave as before
///////////////////////////////////////////////////////////////////////////////

func TestReadBuffer_NonBCDEncodingsAreIdenticalToNoArgs(t *testing.T) {
	data := []byte{0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x07, 0x18}
	for _, byteOrder := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		for _, bitLength := range []uint8{1, 3, 4, 5, 8, 12, 16, 24, 32, 64} {
			for _, encoding := range nonBcdEncodings {
				t.Run(fmt.Sprintf("%v/%dbits/%q", byteOrder, bitLength, encoding), func(t *testing.T) {
					arg := WithEncoding(encoding)
					newBuf := func() ReadBufferByteBased {
						return NewReadBufferByteBased(data, WithByteOrderForReadBufferByteBased(byteOrder))
					}

					wantU8, wantU8Err := newBuf().ReadUint8("t", bitLength)
					gotU8, gotU8Err := newBuf().ReadUint8("t", bitLength, arg)
					assert.Equal(t, wantU8, gotU8)
					assert.Equal(t, wantU8Err == nil, gotU8Err == nil)

					wantU16, _ := newBuf().ReadUint16("t", bitLength)
					gotU16, _ := newBuf().ReadUint16("t", bitLength, arg)
					assert.Equal(t, wantU16, gotU16)

					wantU32, _ := newBuf().ReadUint32("t", bitLength)
					gotU32, _ := newBuf().ReadUint32("t", bitLength, arg)
					assert.Equal(t, wantU32, gotU32)

					wantU64, _ := newBuf().ReadUint64("t", bitLength)
					gotU64, _ := newBuf().ReadUint64("t", bitLength, arg)
					assert.Equal(t, wantU64, gotU64)

					wantI8, _ := newBuf().ReadInt8("t", bitLength)
					gotI8, _ := newBuf().ReadInt8("t", bitLength, arg)
					assert.Equal(t, wantI8, gotI8)

					wantI16, _ := newBuf().ReadInt16("t", bitLength)
					gotI16, _ := newBuf().ReadInt16("t", bitLength, arg)
					assert.Equal(t, wantI16, gotI16)

					wantI32, _ := newBuf().ReadInt32("t", bitLength)
					gotI32, _ := newBuf().ReadInt32("t", bitLength, arg)
					assert.Equal(t, wantI32, gotI32)

					wantI64, _ := newBuf().ReadInt64("t", bitLength)
					gotI64, _ := newBuf().ReadInt64("t", bitLength, arg)
					assert.Equal(t, wantI64, gotI64)
				})
			}
		}
	}
}

// TestReadBuffer_BCDDoesNotAffectNonIntegerMethods pins that BCD stays confined
// to the integer methods, exactly like it is confined in plc4j (where every
// other EncodingBCD method is simply not implemented). It must not turn into an
// error either.
func TestReadBuffer_BCDDoesNotAffectNonIntegerMethods(t *testing.T) {
	data := []byte{0xA1, 0xB2, 0xC3, 0xD4, 0xE5, 0xF6, 0x07, 0x18}
	newBuf := func() ReadBufferByteBased { return NewReadBufferByteBased(data) }

	wantBit, wantErr := newBuf().ReadBit("t")
	gotBit, gotErr := newBuf().ReadBit("t", bcdArg)
	assert.Equal(t, wantBit, gotBit)
	assert.Equal(t, wantErr, gotErr)

	wantByte, _ := newBuf().ReadByte("t")
	gotByte, _ := newBuf().ReadByte("t", bcdArg)
	assert.Equal(t, wantByte, gotByte)

	wantBytes, _ := newBuf().ReadByteArray("t", 4)
	gotBytes, _ := newBuf().ReadByteArray("t", 4, bcdArg)
	assert.Equal(t, wantBytes, gotBytes)

	wantF32, _ := newBuf().ReadFloat32("t", 32)
	gotF32, _ := newBuf().ReadFloat32("t", 32, bcdArg)
	assert.Equal(t, wantF32, gotF32)

	wantF64, _ := newBuf().ReadFloat64("t", 64)
	gotF64, _ := newBuf().ReadFloat64("t", 64, bcdArg)
	assert.Equal(t, wantF64, gotF64)

	wantBigInt, _ := newBuf().ReadBigInt("t", 64)
	gotBigInt, _ := newBuf().ReadBigInt("t", 64, bcdArg)
	assert.Equal(t, wantBigInt, gotBigInt)

	wantBigFloat, _ := newBuf().ReadBigFloat("t", 64)
	gotBigFloat, _ := newBuf().ReadBigFloat("t", 64, bcdArg)
	assert.Equal(t, wantBigFloat, gotBigFloat)

	wantString, _ := newBuf().ReadString("t", 32)
	gotString, _ := newBuf().ReadString("t", 32, bcdArg)
	assert.Equal(t, wantString, gotString)
}

func TestWriteBuffer_NonBCDEncodingsAreIdenticalToNoArgs(t *testing.T) {
	for _, byteOrder := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		for _, bitLength := range []uint8{1, 3, 4, 5, 8, 12, 16, 24, 32, 64} {
			for _, encoding := range nonBcdEncodings {
				t.Run(fmt.Sprintf("%v/%dbits/%q", byteOrder, bitLength, encoding), func(t *testing.T) {
					arg := WithEncoding(encoding)
					newBuf := func() WriteBufferByteBased {
						return NewWriteBufferByteBased(WithByteOrderForByteBasedBuffer(byteOrder))
					}
					u := uint64(0xABCDEF0123456789)
					i := int64(-81985529216486896)

					assertSameWrite := func(name string, want, got func(WriteBufferByteBased) error) {
						t.Helper()
						wantBuf, gotBuf := newBuf(), newBuf()
						wantErr, gotErr := want(wantBuf), got(gotBuf)
						assert.Equal(t, wantErr == nil, gotErr == nil, name)
						assert.Equal(t, wantBuf.GetBytes(), gotBuf.GetBytes(), name)
						assert.Equal(t, wantBuf.GetPos(), gotBuf.GetPos(), name)
					}

					assertSameWrite("uint8",
						func(b WriteBufferByteBased) error { return b.WriteUint8("t", bitLength, uint8(u)) },
						func(b WriteBufferByteBased) error { return b.WriteUint8("t", bitLength, uint8(u), arg) })
					assertSameWrite("uint16",
						func(b WriteBufferByteBased) error { return b.WriteUint16("t", bitLength, uint16(u)) },
						func(b WriteBufferByteBased) error { return b.WriteUint16("t", bitLength, uint16(u), arg) })
					assertSameWrite("uint32",
						func(b WriteBufferByteBased) error { return b.WriteUint32("t", bitLength, uint32(u)) },
						func(b WriteBufferByteBased) error { return b.WriteUint32("t", bitLength, uint32(u), arg) })
					assertSameWrite("uint64",
						func(b WriteBufferByteBased) error { return b.WriteUint64("t", bitLength, uint64(u)) },
						func(b WriteBufferByteBased) error { return b.WriteUint64("t", bitLength, uint64(u), arg) })
					assertSameWrite("int8",
						func(b WriteBufferByteBased) error { return b.WriteInt8("t", bitLength, int8(i)) },
						func(b WriteBufferByteBased) error { return b.WriteInt8("t", bitLength, int8(i), arg) })
					assertSameWrite("int16",
						func(b WriteBufferByteBased) error { return b.WriteInt16("t", bitLength, int16(i)) },
						func(b WriteBufferByteBased) error { return b.WriteInt16("t", bitLength, int16(i), arg) })
					assertSameWrite("int32",
						func(b WriteBufferByteBased) error { return b.WriteInt32("t", bitLength, int32(i)) },
						func(b WriteBufferByteBased) error { return b.WriteInt32("t", bitLength, int32(i), arg) })
					assertSameWrite("int64",
						func(b WriteBufferByteBased) error { return b.WriteInt64("t", bitLength, i) },
						func(b WriteBufferByteBased) error { return b.WriteInt64("t", bitLength, i, arg) })
				})
			}
		}
	}
}

func TestWriteBuffer_BCDDoesNotAffectNonIntegerMethods(t *testing.T) {
	assertSameWrite := func(name string, want, got func(WriteBufferByteBased) error) {
		t.Helper()
		wantBuf, gotBuf := NewWriteBufferByteBased(), NewWriteBufferByteBased()
		wantErr, gotErr := want(wantBuf), got(gotBuf)
		assert.Equal(t, wantErr == nil, gotErr == nil, name)
		assert.Equal(t, wantBuf.GetBytes(), gotBuf.GetBytes(), name)
	}

	assertSameWrite("bit",
		func(b WriteBufferByteBased) error { return b.WriteBit("t", true) },
		func(b WriteBufferByteBased) error { return b.WriteBit("t", true, bcdArg) })
	assertSameWrite("byte",
		func(b WriteBufferByteBased) error { return b.WriteByte("t", 0xAB) },
		func(b WriteBufferByteBased) error { return b.WriteByte("t", 0xAB, bcdArg) })
	assertSameWrite("byteArray",
		func(b WriteBufferByteBased) error { return b.WriteByteArray("t", []byte{0xAB, 0xCD}) },
		func(b WriteBufferByteBased) error { return b.WriteByteArray("t", []byte{0xAB, 0xCD}, bcdArg) })
	assertSameWrite("float32",
		func(b WriteBufferByteBased) error { return b.WriteFloat32("t", 32, 1.5) },
		func(b WriteBufferByteBased) error { return b.WriteFloat32("t", 32, 1.5, bcdArg) })
	assertSameWrite("float64",
		func(b WriteBufferByteBased) error { return b.WriteFloat64("t", 64, 1.5) },
		func(b WriteBufferByteBased) error { return b.WriteFloat64("t", 64, 1.5, bcdArg) })
	assertSameWrite("string",
		func(b WriteBufferByteBased) error { return b.WriteString("t", 32, "ab") },
		func(b WriteBufferByteBased) error { return b.WriteString("t", 32, "ab", bcdArg) })
}

///////////////////////////////////////////////////////////////////////////////
// hot path guard
///////////////////////////////////////////////////////////////////////////////

// BenchmarkReadUint8WithNonBCDEncoding measures the encoding dispatch on the
// hottest path in the repo (several thousand generated numeric reads carry a
// non-BCD encoding argument).
//
// It is a stopwatch, not a guard: a Go benchmark cannot fail on a regression,
// and the allocations it reports are the variadic []WithReaderArgs slices built
// at the CALL SITES in this file, not anything the dispatch itself allocates.
// The evidence that readerArgsSelectBCD is allocation free is the escape
// analysis, which nothing here can assert:
//
//	go build -gcflags=-m ./spi/utils/ 2>&1 | grep readerArgs
//	  ... inlining call to readerArgsSelectBCD
//	  ... readerArgs does not escape
func BenchmarkReadUint8WithNonBCDEncoding(b *testing.B) {
	data := make([]byte, 1024)
	arg := WithEncoding("UTF8")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rb := NewReadBufferByteBased(data)
		for range 1024 {
			if _, err := rb.ReadUint8("t", 8, arg); err != nil {
				b.Fatal(err)
			}
		}
	}
}

// TestBCD_DivergenceFromPlc4jIsIntentional pins the exact cases where plc4j
// disagrees with this implementation, so that a future "make it match Java"
// refactor has to consciously break a test that says why it must not.
//
// plc4j is the buggy side: its ReadBufferByteBased.readBits right-aligns a
// partial read while EncodingBCD reads nibbles left-aligned from bytes[0], and
// the write path is broken symmetrically. See the note on decodeBCD.
func TestBCD_DivergenceFromPlc4jIsIntentional(t *testing.T) {
	decodeTests := []struct {
		raw         uint64
		bitLength   uint8
		want        uint64
		plc4jYields uint64
	}{
		{raw: 0x1, bitLength: 4, want: 1, plc4jYields: 0},
		{raw: 0x4, bitLength: 4, want: 4, plc4jYields: 0},
		{raw: 0x12, bitLength: 8, want: 12, plc4jYields: 12},        // agrees
		{raw: 0x123, bitLength: 12, want: 123, plc4jYields: 12},     // diverges
		{raw: 0x1234, bitLength: 16, want: 1234, plc4jYields: 1234}, // agrees
	}
	for _, tt := range decodeTests {
		t.Run(fmt.Sprintf("decode_0x%X@%d", tt.raw, tt.bitLength), func(t *testing.T) {
			got, err := decodeBCD(tt.raw, tt.bitLength)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
			if tt.bitLength%8 != 0 {
				assert.NotEqual(t, tt.plc4jYields, got,
					"bit lengths that are not a multiple of 8 must NOT reproduce the plc4j answer")
			} else {
				assert.Equal(t, tt.plc4jYields, got,
					"bit lengths that are a multiple of 8 must agree with plc4j")
			}
		})
	}

	// Write side: what plc4j actually puts on the wire is in the plc4jEmits column.
	encodeTests := []struct {
		value      uint64
		bitLength  uint8
		want       uint64
		plc4jEmits uint64
	}{
		{value: 4, bitLength: 4, want: 0x4, plc4jEmits: 0x0},
		{value: 35, bitLength: 8, want: 0x35, plc4jEmits: 0x35},        // agrees
		{value: 123, bitLength: 12, want: 0x123, plc4jEmits: 0x230},    // diverges
		{value: 1234, bitLength: 16, want: 0x1234, plc4jEmits: 0x1234}, // agrees
	}
	for _, tt := range encodeTests {
		t.Run(fmt.Sprintf("encode_%d@%d", tt.value, tt.bitLength), func(t *testing.T) {
			got, err := encodeBCD(tt.value, tt.bitLength)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
			if tt.bitLength%8 != 0 {
				assert.NotEqual(t, tt.plc4jEmits, got)
			} else {
				assert.Equal(t, tt.plc4jEmits, got)
			}
		})
	}
}
