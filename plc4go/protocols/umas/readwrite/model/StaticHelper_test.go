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

package model

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func TestParseTerminatedStringVariableLength(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		wantValue   string
		wantConsume uint32
		wantErr     bool
	}{
		{
			name:        "normal string",
			data:        []byte("hello\x00"),
			wantValue:   "hello",
			wantConsume: 6,
		},
		{
			name:        "empty string is just the terminator",
			data:        []byte("\x00"),
			wantValue:   "",
			wantConsume: 1,
		},
		{
			name: "stops at the terminator and leaves the tail",
			// The tail must survive for the following field to parse.
			data:        []byte("hi\x00tail"),
			wantValue:   "hi",
			wantConsume: 3,
		},
		{
			name:      "multi byte runes survive",
			data:      append([]byte("äöü"), 0x00),
			wantValue: "äöü",
			// 3 two-byte runes plus the terminator.
			wantConsume: 7,
		},
		{
			name:    "missing terminator is an error",
			data:    []byte("nope"),
			wantErr: true,
		},
		{
			name:    "empty buffer is an error",
			data:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := utils.NewReadBufferByteBased(tt.data)
			value, err := ParseTerminatedString(context.Background(), rb, -1)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantConsume, rb.GetPos(), "unexpected number of bytes consumed")
		})
	}
}

func TestParseTerminatedStringFixedLength(t *testing.T) {
	tests := []struct {
		name         string
		data         []byte
		stringLength int
		wantValue    string
		wantConsume  uint32
		wantErr      bool
	}{
		{
			name:         "padded field is trimmed at the terminator",
			data:         []byte("hi\x00\x00\x00\x00\x00\x00"),
			stringLength: 8,
			wantValue:    "hi",
			// The whole fixed width field is consumed, not just the string.
			wantConsume: 8,
		},
		{
			name:         "unterminated field uses the whole width",
			data:         []byte("abcdefgh"),
			stringLength: 8,
			wantValue:    "abcdefgh",
			wantConsume:  8,
		},
		{
			name:         "terminator in the first byte yields an empty string",
			data:         []byte("\x00bcdefgh"),
			stringLength: 8,
			wantValue:    "",
			wantConsume:  8,
		},
		{
			name:         "zero width reads nothing",
			data:         []byte("abc"),
			stringLength: 0,
			wantValue:    "",
			wantConsume:  0,
		},
		{
			name:         "truncated buffer is an error",
			data:         []byte("ab"),
			stringLength: 8,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := utils.NewReadBufferByteBased(tt.data)
			value, err := ParseTerminatedString(context.Background(), rb, tt.stringLength)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantConsume, rb.GetPos(), "unexpected number of bytes consumed")
		})
	}
}

func TestParseTerminatedStringBytes(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		numberOfValues uint16
		wantValue      string
		wantConsume    uint32
		wantErr        bool
	}{
		{
			name:           "stops at the terminator",
			data:           []byte("hi\x00tail"),
			numberOfValues: 7,
			wantValue:      "hi",
			// numberOfValues is a byte count, so the whole field is consumed.
			wantConsume: 7,
		},
		{
			name:           "single byte field",
			data:           []byte("A"),
			numberOfValues: 1,
			wantValue:      "A",
			wantConsume:    1,
		},
		{
			name:           "unterminated field returns everything",
			data:           []byte("abc"),
			numberOfValues: 3,
			wantValue:      "abc",
			wantConsume:    3,
		},
		{
			name:           "zero values",
			data:           []byte("abc"),
			numberOfValues: 0,
			wantValue:      "",
			wantConsume:    0,
		},
		{
			name:           "more values than bytes is an error",
			data:           []byte("abc"),
			numberOfValues: 4,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := utils.NewReadBufferByteBased(tt.data)
			value, err := ParseTerminatedStringBytes(context.Background(), rb, tt.numberOfValues)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantConsume, rb.GetPos(), "unexpected number of bytes consumed")
		})
	}
}

func TestSerializeTerminatedString(t *testing.T) {
	tests := []struct {
		name         string
		value        any
		stringLength int
		wantBytes    []byte
		wantErr      bool
	}{
		{
			name:         "variable length appends a terminator",
			value:        "hello",
			stringLength: -1,
			wantBytes:    []byte("hello\x00"),
		},
		{
			name:         "variable length empty string is just the terminator",
			value:        "",
			stringLength: -1,
			wantBytes:    []byte{0x00},
		},
		{
			name:         "variable length keeps multi byte runes",
			value:        "äöü",
			stringLength: -1,
			wantBytes:    append([]byte("äöü"), 0x00),
		},
		{
			name:         "fixed width pads with terminators",
			value:        "hi",
			stringLength: 8,
			wantBytes:    []byte("hi\x00\x00\x00\x00\x00\x00"),
		},
		{
			name:         "fixed width empty string is all terminators",
			value:        "",
			stringLength: 4,
			wantBytes:    []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:         "fixed width exactly filled carries no terminator",
			value:        "abcd",
			stringLength: 4,
			wantBytes:    []byte("abcd"),
		},
		{
			name:         "fixed width truncates an oversized value",
			value:        "abcdefgh",
			stringLength: 4,
			wantBytes:    []byte("abcd"),
		},
		{
			name:  "fixed width truncation respects rune boundaries",
			value: "äöü",
			// Cutting at 3 bytes would split the second rune, so only the
			// first rune fits and the rest is padded.
			stringLength: 3,
			wantBytes:    []byte{0xc3, 0xa4, 0x00},
		},
		{
			name:         "zero width writes nothing",
			value:        "abc",
			stringLength: 0,
			// Nothing was written, so the buffer never allocated a backing array.
			wantBytes: nil,
		},
		{
			name:         "PlcValue input",
			value:        spiValues.NewPlcSTRING("hi"),
			stringLength: 4,
			wantBytes:    []byte("hi\x00\x00"),
		},
		{
			name:         "PlcValue input variable length",
			value:        spiValues.NewPlcSTRING("hi"),
			stringLength: -1,
			wantBytes:    []byte("hi\x00"),
		},
		{
			name:         "nil value serializes as an empty string",
			value:        nil,
			stringLength: -1,
			wantBytes:    []byte{0x00},
		},
		{
			name:         "unsupported value type is an error",
			value:        42,
			stringLength: -1,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wb := utils.NewWriteBufferByteBased()
			err := SerializeTerminatedString(context.Background(), wb, tt.value, tt.stringLength)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBytes, wb.GetBytes())
		})
	}
}

func TestTerminatedStringRoundTrip(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		stringLength int
	}{
		{name: "variable length", value: "MyVariable", stringLength: -1},
		{name: "variable length empty", value: "", stringLength: -1},
		{name: "variable length multi byte", value: "Wärmemenge", stringLength: -1},
		{name: "fixed width", value: "MyVar", stringLength: 16},
		{name: "fixed width empty", value: "", stringLength: 16},
		{name: "fixed width exactly filled", value: "abcd", stringLength: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			wb := utils.NewWriteBufferByteBased()
			require := assert.New(t)
			require.NoError(SerializeTerminatedString(ctx, wb, tt.value, tt.stringLength))

			serialized := wb.GetBytes()
			if tt.stringLength >= 0 {
				require.Len(serialized, tt.stringLength, "fixed width field must occupy its full width")
			} else {
				require.Len(serialized, len(tt.value)+1, "variable length field is the value plus a terminator")
			}

			rb := utils.NewReadBufferByteBased(serialized)
			parsed, err := ParseTerminatedString(ctx, rb, tt.stringLength)
			require.NoError(err)
			require.Equal(tt.value, parsed)
		})
	}
}

func TestParseTerminatedStringBytesRoundTrip(t *testing.T) {
	ctx := context.Background()
	const numberOfValues = uint16(8)

	wb := utils.NewWriteBufferByteBased()
	assert.NoError(t, SerializeTerminatedString(ctx, wb, spiValues.NewPlcSTRING("plc4x"), numberOfValues))
	serialized := wb.GetBytes()
	assert.Len(t, serialized, int(numberOfValues))

	rb := utils.NewReadBufferByteBased(serialized)
	parsed, err := ParseTerminatedStringBytes(ctx, rb, numberOfValues)
	assert.NoError(t, err)
	assert.Equal(t, "plc4x", parsed)
}

func TestWriteSizeIndexToByteCount(t *testing.T) {
	tests := []struct {
		name      string
		sizeIndex uint8
		want      uint16
	}{
		{name: "index 1 is one byte (BOOL, BYTE)", sizeIndex: 1, want: 1},
		{name: "index 2 is two bytes (INT, UINT, WORD)", sizeIndex: 2, want: 2},
		{name: "index 3 is four bytes (DINT, REAL, TIME)", sizeIndex: 3, want: 4},
		{name: "index 4 is eight bytes (DATE_AND_TIME)", sizeIndex: 4, want: 8},
		{name: "index 17 marks STRING, one byte per character", sizeIndex: 17, want: 1},
		{name: "index 0 falls through to itself", sizeIndex: 0, want: 0},
		{name: "unknown index falls through to itself", sizeIndex: 9, want: 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, WriteSizeIndexToByteCount(context.Background(), tt.sizeIndex))
		})
	}
}

// TestWriteSizeIndexToByteCountMatchesDataTypeSize pins the helper against the
// dataTypeSize/requestSize pairs declared for UmasDataType in the mspec.
func TestWriteSizeIndexToByteCountMatchesDataTypeSize(t *testing.T) {
	ctx := context.Background()
	for _, dataType := range UmasDataTypeValues {
		requestSize := dataType.RequestSize()
		if requestSize == 17 {
			// STRING is length delimited, its dataTypeSize is the per
			// character size and cannot be derived from the index.
			continue
		}
		assert.Equal(t, uint16(dataType.DataTypeSize()), WriteSizeIndexToByteCount(ctx, requestSize),
			"size index of %v should map back to its data type size", dataType)
	}
}
