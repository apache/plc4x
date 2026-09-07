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

package abeth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
)

func TestTagHandler_ParseTag(t *testing.T) {
	tests := []struct {
		name              string
		address           string
		wantByteSize      uint8
		wantFileNumber    uint8
		wantFileType      FileType
		wantElementNumber uint8
		wantBitNumber     uint8
	}{
		{
			name:              "an integer tag takes its width from the size suffix",
			address:           "N7:3:INTEGER[2]",
			wantByteSize:      2,
			wantFileNumber:    7,
			wantFileType:      FileTypeInteger,
			wantElementNumber: 3,
		},
		{
			name:              "a word tag is always two bytes wide",
			address:           "N7:3:WORD",
			wantByteSize:      2,
			wantFileNumber:    7,
			wantFileType:      FileTypeWord,
			wantElementNumber: 3,
		},
		{
			name:              "a dword tag is always four bytes wide",
			address:           "N100:200:DWORD",
			wantByteSize:      4,
			wantFileNumber:    100,
			wantFileType:      FileTypeDword,
			wantElementNumber: 200,
		},
		{
			name:              "a singlebit tag reads two bytes and remembers the bit",
			address:           "N7:3/5:SINGLEBIT",
			wantByteSize:      2,
			wantFileNumber:    7,
			wantFileType:      FileTypeSinglebit,
			wantElementNumber: 3,
			wantBitNumber:     5,
		},
		{
			name:              "a lower case file type is accepted",
			address:           "N7:3:integer[1]",
			wantByteSize:      1,
			wantFileNumber:    7,
			wantFileType:      FileTypeInteger,
			wantElementNumber: 3,
		},
		{
			name:              "a redundant size suffix which agrees with a fixed-width type is accepted",
			address:           "N7:3:WORD[2]",
			wantByteSize:      2,
			wantFileNumber:    7,
			wantFileType:      FileTypeWord,
			wantElementNumber: 3,
		},
		{
			name:              "a bit number on a non-bit tag is kept",
			address:           "N7:3/9:INTEGER[4]",
			wantByteSize:      4,
			wantFileNumber:    7,
			wantFileType:      FileTypeInteger,
			wantElementNumber: 3,
			wantBitNumber:     9,
		},
	}
	tagHandler := NewTagHandler()
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			parsed, err := tagHandler.ParseTag(testCase.address)
			require.NoError(t, err)
			tag, ok := parsed.(PlcTag)
			require.True(t, ok, "%T is not an ab-eth tag", parsed)
			assert.Equal(t, testCase.wantByteSize, tag.GetByteSize(), "byte size")
			assert.Equal(t, testCase.wantFileNumber, tag.GetFileNumber(), "file number")
			assert.Equal(t, testCase.wantFileType, tag.GetFileType(), "file type")
			assert.Equal(t, testCase.wantElementNumber, tag.GetElementNumber(), "element number")
			assert.Equal(t, testCase.wantBitNumber, tag.GetBitNumber(), "bit number")
		})
	}
}

func TestTagHandler_ParseTagRejects(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{name: "an empty address", address: ""},
		{name: "an address without the leading N", address: "7:3:INTEGER[1]"},
		{name: "an address without an element", address: "N7:INTEGER[1]"},
		{name: "an address without a file type", address: "N7:3"},
		{name: "an unknown file type", address: "N7:3:BANANA[1]"},
		{name: "a real file type without a size", address: "N7:3:INTEGER"},
		{name: "a zero size", address: "N7:3:INTEGER[0]"},
		{name: "a file number that doesn't fit into a byte", address: "N256:3:INTEGER[1]"},
		{name: "an element number that doesn't fit into a byte", address: "N7:256:INTEGER[1]"},
		{name: "a size that doesn't fit into a byte", address: "N7:3:INTEGER[256]"},
		{name: "a bit number past the end of the element", address: "N7:3/16:SINGLEBIT"},
		{name: "trailing garbage", address: "N7:3:WORDx"},
		// A fixed-width type takes its width from its name, so a suffix which says something else
		// can't be honored - it used to be dropped on the floor, leaving a two byte WORD tag.
		{name: "a size suffix contradicting a word tag", address: "N7:3:WORD[9]"},
		{name: "a size suffix contradicting a dword tag", address: "N7:3:DWORD[2]"},
		{name: "a size suffix contradicting a singlebit tag", address: "N7:3/5:SINGLEBIT[1]"},
	}
	tagHandler := NewTagHandler()
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tag, err := tagHandler.ParseTag(testCase.address)
			assert.Error(t, err)
			assert.Nil(t, tag)
		})
	}
}

func TestTagHandler_ParseQueryIsUnsupported(t *testing.T) {
	query, err := NewTagHandler().ParseQuery("N7:*")
	assert.Error(t, err)
	assert.Nil(t, query)
}

// TestTag_AddressStringRoundTrips is what makes subscriptions work: the polling subscriber rebuilds
// its read requests from GetAddressString, so an address that doesn't parse back into the same tag
// would silently poll something else.
func TestTag_AddressStringRoundTrips(t *testing.T) {
	addresses := []string{
		"N7:3:INTEGER[2]",
		"N7:3:WORD",
		"N7:3:DWORD",
		"N7:3/5:SINGLEBIT",
		"N0:0:INTEGER[1]",
	}
	tagHandler := NewTagHandler()
	for _, address := range addresses {
		t.Run(address, func(t *testing.T) {
			tag, err := tagHandler.ParseTag(address)
			require.NoError(t, err)
			assert.Equal(t, address, tag.GetAddressString())

			reparsed, err := tagHandler.ParseTag(tag.GetAddressString())
			require.NoError(t, err)
			assert.Equal(t, tag, reparsed)
		})
	}
}

func TestFileType(t *testing.T) {
	tests := []struct {
		fileType      FileType
		wantName      string
		wantTypeCode  uint8
		wantValueType apiValues.PlcValueType
	}{
		{fileType: FileTypeStatus, wantName: "STATUS", wantTypeCode: 0x84, wantValueType: apiValues.RAW_BYTE_ARRAY},
		{fileType: FileTypeBit, wantName: "BIT", wantTypeCode: 0x85, wantValueType: apiValues.BOOL},
		{fileType: FileTypeTimer, wantName: "TIMER", wantTypeCode: 0x86, wantValueType: apiValues.TIME},
		{fileType: FileTypeCounter, wantName: "COUNTER", wantTypeCode: 0x87, wantValueType: apiValues.RAW_BYTE_ARRAY},
		{fileType: FileTypeControl, wantName: "CONTROL", wantTypeCode: 0x88, wantValueType: apiValues.RAW_BYTE_ARRAY},
		{fileType: FileTypeInteger, wantName: "INTEGER", wantTypeCode: 0x89, wantValueType: apiValues.INT},
		{fileType: FileTypeFloat, wantName: "FLOAT", wantTypeCode: 0x8A, wantValueType: apiValues.REAL},
		{fileType: FileTypeOutput, wantName: "OUTPUT", wantTypeCode: 0x8B, wantValueType: apiValues.RAW_BYTE_ARRAY},
		{fileType: FileTypeInput, wantName: "INPUT", wantTypeCode: 0x8C, wantValueType: apiValues.RAW_BYTE_ARRAY},
		{fileType: FileTypeString, wantName: "STRING", wantTypeCode: 0x8D, wantValueType: apiValues.STRING},
		{fileType: FileTypeAscii, wantName: "ASCII", wantTypeCode: 0x8E, wantValueType: apiValues.STRING},
		{fileType: FileTypeBcd, wantName: "BCD", wantTypeCode: 0x8F, wantValueType: apiValues.RAW_BYTE_ARRAY},
		// The three synthetic types all read an integer file and differ only in the decoding.
		{fileType: FileTypeWord, wantName: "WORD", wantTypeCode: 0x89, wantValueType: apiValues.WORD},
		{fileType: FileTypeDword, wantName: "DWORD", wantTypeCode: 0x89, wantValueType: apiValues.DWORD},
		{fileType: FileTypeSinglebit, wantName: "SINGLEBIT", wantTypeCode: 0x89, wantValueType: apiValues.BOOL},
	}
	for _, testCase := range tests {
		t.Run(testCase.wantName, func(t *testing.T) {
			assert.Equal(t, testCase.wantName, testCase.fileType.String())
			assert.Equal(t, testCase.wantTypeCode, testCase.fileType.GetTypeCode())
			assert.Equal(t, testCase.wantValueType, testCase.fileType.GetPlcValueType())

			byName, ok := FileTypeByName(testCase.wantName)
			require.True(t, ok)
			assert.Equal(t, testCase.fileType, byName)
		})
	}
}

func TestFileTypeByNameRejectsUnknown(t *testing.T) {
	_, ok := FileTypeByName("BANANA")
	assert.False(t, ok)
	assert.Equal(t, "UNKNOWN", FileType(200).String())
	assert.Equal(t, apiValues.NULL, FileType(200).GetPlcValueType())
}

func TestTag_ValueTypeFollowsFileType(t *testing.T) {
	tag, err := NewTagHandler().ParseTag("N7:3:WORD")
	require.NoError(t, err)
	assert.Equal(t, apiValues.WORD, tag.GetValueType())
	assert.Empty(t, tag.GetArrayInfo(), "an ab-eth tag is always a scalar")
}
