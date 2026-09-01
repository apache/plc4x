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

package slmp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

func TestTagHandler_ParseTag(t *testing.T) {
	tests := []struct {
		name               string
		address            string
		wantDeviceCode     readWriteModel.SlmpDeviceCode
		wantDeviceNumber   uint32
		wantDataType       DataType
		wantQuantity       uint16
		wantNumberOfPoints uint16
	}{
		{
			name: "a bare D address is one WORD", address: "D350",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_D, wantDeviceNumber: 350,
			wantDataType: DataTypeWORD, wantQuantity: 1, wantNumberOfPoints: 1,
		},
		{
			name: "D is addressed in decimal", address: "D10",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_D, wantDeviceNumber: 10,
			wantDataType: DataTypeWORD, wantQuantity: 1, wantNumberOfPoints: 1,
		},
		{
			name: "R is addressed in decimal", address: "R200[0..3]:REAL",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_R, wantDeviceNumber: 200,
			wantDataType: DataTypeREAL, wantQuantity: 4, wantNumberOfPoints: 8,
		},
		{
			// W is a link register, which MELSEC addresses in hex (SH-080008 section 8.1).
			name: "W is addressed in hex", address: "W1A[0..9]:WORD",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_W, wantDeviceNumber: 0x1A,
			wantDataType: DataTypeWORD, wantQuantity: 10, wantNumberOfPoints: 10,
		},
		{
			name: "W also takes an explicit 0x prefix", address: "W0x1A",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_W, wantDeviceNumber: 0x1A,
			wantDataType: DataTypeWORD, wantQuantity: 1, wantNumberOfPoints: 1,
		},
		{
			name: "and an uppercase one", address: "W0X1a",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_W, wantDeviceNumber: 0x1A,
			wantDataType: DataTypeWORD, wantQuantity: 1, wantNumberOfPoints: 1,
		},
		{
			// plc4j's greedy device token swallows the leading hex letters here and rejects the
			// address as device "WA" / "WAB"; the non-greedy one takes the shortest device token
			// that lets the rest parse.
			name: "a W address that starts with hex letters", address: "WAB",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_W, wantDeviceNumber: 0xAB,
			wantDataType: DataTypeWORD, wantQuantity: 1, wantNumberOfPoints: 1,
		},
		{
			name: "an all-letter W address", address: "WABCD[0..1]:INT",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_W, wantDeviceNumber: 0xABCD,
			wantDataType: DataTypeINT, wantQuantity: 2, wantNumberOfPoints: 2,
		},
		{
			name: "the device token is case insensitive", address: "d350:int",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_D, wantDeviceNumber: 350,
			wantDataType: DataTypeINT, wantQuantity: 1, wantNumberOfPoints: 1,
		},
		{
			name: "a double-word type takes two points per element", address: "D0[0..2]:DINT",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_D, wantDeviceNumber: 0,
			wantDataType: DataTypeDINT, wantQuantity: 3, wantNumberOfPoints: 6,
		},
		{
			name: "UDINT too", address: "D0:UDINT",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_D, wantDeviceNumber: 0,
			wantDataType: DataTypeUDINT, wantQuantity: 1, wantNumberOfPoints: 2,
		},
		{
			name: "UINT is a single word", address: "D0[0..7]:UINT",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_D, wantDeviceNumber: 0,
			wantDataType: DataTypeUINT, wantQuantity: 8, wantNumberOfPoints: 8,
		},
		{
			name: "the largest device number the 24-bit field can carry", address: "D16777215",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_D, wantDeviceNumber: 0xFFFFFF,
			wantDataType: DataTypeWORD, wantQuantity: 1, wantNumberOfPoints: 1,
		},
		{
			name: "the largest single-frame word transfer", address: "D0[0..959]:WORD",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_D, wantDeviceNumber: 0,
			wantDataType: DataTypeWORD, wantQuantity: 960, wantNumberOfPoints: 960,
		},
		{
			name: "which is half as many double-word elements", address: "D0[0..479]:REAL",
			wantDeviceCode: readWriteModel.SlmpDeviceCode_D, wantDeviceNumber: 0,
			wantDataType: DataTypeREAL, wantQuantity: 480, wantNumberOfPoints: 960,
		},
	}
	handler := NewTagHandler()
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			parsed, err := handler.ParseTag(testCase.address)
			require.NoError(t, err)
			tag, ok := parsed.(PlcTag)
			require.True(t, ok, "%T is not an slmp tag", parsed)
			assert.Equal(t, testCase.wantDeviceCode, tag.GetDeviceCode())
			assert.Equal(t, testCase.wantDeviceNumber, tag.GetDeviceNumber())
			assert.Equal(t, testCase.wantDataType, tag.GetDataType())
			assert.Equal(t, testCase.wantQuantity, tag.GetQuantity())
			assert.Equal(t, testCase.wantNumberOfPoints, tag.GetNumberOfPoints())
		})
	}
}

func TestTagHandler_ParseTagRejects(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{name: "nothing at all", address: ""},
		{name: "a device without an address", address: "D"},
		{name: "an address without a device", address: "350"},
		// plc4j supports the word devices D, W and R and nothing else. The bit devices the mspec's
		// device-code enum knows are read in a different point unit, which this version can't do.
		{name: "the bit device M", address: "M100"},
		{name: "the bit device X", address: "X20"},
		{name: "the bit device Y", address: "Y160"},
		{name: "the bit device B", address: "B100"},
		{name: "the timer device TN", address: "TN0"},
		{name: "a device nobody knows", address: "Z100"},
		// A hex prefix says "read the digits as hex", which can only be a mistake on a decimal
		// device.
		{name: "a 0x prefix on a decimal device", address: "D0x10"},
		{name: "a 0x prefix on R", address: "R0x10"},
		{name: "a decimal device with hex digits", address: "DAB"},
		{name: "an unknown data type", address: "D350:LREAL"},
		{name: "a data type nobody wrote", address: "D350:"},
		// There is no way to ask for zero devices any more - a range is written with the
		// indices it covers - but an inverted one is still nonsense.
		{name: "an inverted range", address: "D350[3..1]:WORD"},
		{name: "a negative quantity", address: "D350:WORD[-1]"},
		{name: "an empty quantity", address: "D350:WORD[]"},
		{name: "a device number beyond the 24-bit field", address: "D16777216"},
		// There is no request optimizer to split a bigger transfer into several frames, so a tag
		// that would need one is refused rather than sent as a frame the device rejects.
		{name: "more words than one frame carries", address: "D0[0..960]:WORD"},
		{name: "more double words than one frame carries", address: "D0[0..480]:REAL"},
		{name: "trailing junk", address: "D350:WORDx"},
	}
	handler := NewTagHandler()
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tag, err := handler.ParseTag(testCase.address)
			assert.Error(t, err)
			assert.Nil(t, tag)
		})
	}
}

// TestTagHandler_ParseQuery pins that browsing says so rather than pretending: plc4j's
// SlmpTagHandler.parseQuery throws UnsupportedOperationException.
func TestTagHandler_ParseQuery(t *testing.T) {
	query, err := NewTagHandler().ParseQuery("D*")
	assert.Error(t, err)
	assert.Nil(t, query)
}

// TestTag_AddressStringRoundTrips is load bearing rather than cosmetic: the polling subscriber
// rebuilds every poll's read request from GetAddressString, so an address that doesn't parse back
// to the same tag would silently poll something else.
func TestTag_AddressStringRoundTrips(t *testing.T) {
	addresses := []string{
		"D350", "D350:INT", "D350[0..3]:WORD", "R200[0..3]:REAL", "W1A[0..9]:WORD", "W0x1A",
		"WAB", "D0[0..2]:DINT", "D0:UDINT", "D0[0..7]:UINT", "D16777215",
	}
	handler := NewTagHandler()
	for _, address := range addresses {
		t.Run(address, func(t *testing.T) {
			parsed, err := handler.ParseTag(address)
			require.NoError(t, err)
			reparsed, err := handler.ParseTag(parsed.GetAddressString())
			require.NoError(t, err, "the address string %q didn't parse back", parsed.GetAddressString())
			assert.Equal(t, parsed, reparsed)
			// And once more, to catch a spelling that changes on every round.
			assert.Equal(t, parsed.GetAddressString(), reparsed.GetAddressString())
		})
	}
}

func TestTag_AddressStringSpelling(t *testing.T) {
	tests := []struct {
		tag  PlcTag
		want string
	}{
		{tag: NewTag(readWriteModel.SlmpDeviceCode_D, 350, DataTypeWORD, 1), want: "D350:WORD"},
		{tag: NewTag(readWriteModel.SlmpDeviceCode_R, 200, DataTypeREAL, 4), want: "R200[0..3]:REAL"},
		// A W address is written in hex, because that is how the handler reads one back.
		{tag: NewTag(readWriteModel.SlmpDeviceCode_W, 0x1A, DataTypeINT, 2), want: "W0x1A[0..1]:INT"},
	}
	for _, testCase := range tests {
		t.Run(testCase.want, func(t *testing.T) {
			assert.Equal(t, testCase.want, testCase.tag.GetAddressString())
		})
	}
}

func TestTag_Metadata(t *testing.T) {
	handler := NewTagHandler()

	scalar, err := handler.ParseTag("D350:INT")
	require.NoError(t, err)
	assert.Equal(t, apiValues.INT, scalar.GetValueType())
	assert.Empty(t, scalar.GetArrayInfo(), "a scalar tag has no array info")

	array, err := handler.ParseTag("R200[0..3]:REAL")
	require.NoError(t, err)
	assert.Equal(t, apiValues.REAL, array.GetValueType())
	require.Len(t, array.GetArrayInfo(), 1)
	// Both bounds are inclusive, so four elements are 0..3 and the size is 4.
	assert.Equal(t, uint32(0), array.GetArrayInfo()[0].GetLowerBound())
	assert.Equal(t, uint32(3), array.GetArrayInfo()[0].GetUpperBound())
	assert.Equal(t, uint32(4), array.GetArrayInfo()[0].GetSize())

	// Every slmp tag is usable as a subscription tag, because subscriptions are emulated by polling
	// the read path and the subscription builder only accepts tags which are one.
	subscriptionTag, ok := array.(apiModel.PlcSubscriptionTag)
	require.True(t, ok, "%T is not a subscription tag", array)
	assert.Equal(t, apiModel.SubscriptionChangeOfState, subscriptionTag.GetPlcSubscriptionType())
	assert.Zero(t, subscriptionTag.GetDuration())
}

func TestTag_Serialize(t *testing.T) {
	tag := NewTag(readWriteModel.SlmpDeviceCode_D, 350, DataTypeINT, 2)
	theBytes, err := tag.Serialize()
	require.NoError(t, err)
	assert.NotEmpty(t, theBytes)
	assert.NotEmpty(t, fmt.Sprintf("%s", tag))
}

func TestCastToSlmpTagFromPlcTag(t *testing.T) {
	tag := NewTag(readWriteModel.SlmpDeviceCode_D, 350, DataTypeINT, 2)
	castBack, err := castToSlmpTagFromPlcTag(tag)
	require.NoError(t, err)
	assert.Equal(t, tag, castBack)

	// A tag from another driver has to be refused rather than panicking somewhere downstream.
	foreign, err := castToSlmpTagFromPlcTag(spiModel.NewDefaultPlcSubscriptionTag(
		apiModel.SubscriptionCyclic, nil, 0))
	assert.Error(t, err)
	assert.Nil(t, foreign)
}

// An SLMP address is a device number, so a selection that starts past the declared base is
// resolved into the address itself: "D100[4..7]" is the same read as "D104[0..3]".
func TestTagHandler_ParseTag_consumesTheSelectionOffset(t *testing.T) {
	handler := NewTagHandler()

	shifted, err := handler.ParseTag("D100[4..7]:INT")
	require.NoError(t, err)
	equivalent, err := handler.ParseTag("D104[0..3]:INT")
	require.NoError(t, err)
	assert.Equal(t, equivalent, shifted)
	assert.Equal(t, "D104[0..3]:INT", shifted.GetAddressString())

	// A declared base is what the offset is measured from, so [4..7;4] shifts nothing.
	fromDeclaredBase, err := handler.ParseTag("D100[4..7;4]:INT")
	require.NoError(t, err)
	unshifted, err := handler.ParseTag("D100[0..3]:INT")
	require.NoError(t, err)
	assert.Equal(t, unshifted, fromDeclaredBase)

	// A bare index selects the device at that index, not that many devices.
	single, err := handler.ParseTag("D100[4]:INT")
	require.NoError(t, err)
	assert.Equal(t, "D104:INT", single.GetAddressString())
	assert.Empty(t, single.GetArrayInfo(), "one device is a scalar")
}

// Addresses written with the count after the type must fail, naming what to write instead.
func TestTagHandler_ParseTag_rejectsTheOldCountSuffix(t *testing.T) {
	handler := NewTagHandler()

	for _, address := range []string{"D350:WORD[2]", "R200:REAL[4]", "W1A:WORD[10]"} {
		t.Run(address, func(t *testing.T) {
			_, err := handler.ParseTag(address)
			require.Error(t, err, address)
			assert.Contains(t, err.Error(), "invalid address", address)
		})
	}
}

// A selection offset counts elements; an SLMP device number counts 16-bit words. The point count
// already scales by WordsPerElement, so an unscaled offset moved the read rather than shortening it.
func TestTagHandler_ParseTag_scalesTheOffsetByTheWordsPerElement(t *testing.T) {
	handler := NewTagHandler()

	oneWord, err := handler.ParseTag("D100[4]:INT")
	require.NoError(t, err)
	assert.Equal(t, "D104:INT", oneWord.GetAddressString())

	// The fifth DINT begins eight words along, at D108.
	twoWords, err := handler.ParseTag("D100[4]:DINT")
	require.NoError(t, err)
	assert.Equal(t, "D108:DINT", twoWords.GetAddressString())

	// A declared base is measured in elements too, so [4..7;4] shifts nothing.
	fromBase, err := handler.ParseTag("D100[4..7;4]:DINT")
	require.NoError(t, err)
	assert.Equal(t, "D100[0..3]:DINT", fromBase.GetAddressString())
}
