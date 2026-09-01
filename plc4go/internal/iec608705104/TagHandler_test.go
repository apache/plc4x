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

package iec608705104

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// parseTag parses an address, insisting that it really is an IEC 60870-5-104 tag.
func parseTag(t *testing.T, address string) Tag {
	t.Helper()
	parsed, err := NewTagHandler(testutils.EnrichOptionsWithOptionsForTesting(t)...).ParseTag(address)
	require.NoError(t, err)
	tag, ok := parsed.(Tag)
	require.True(t, ok, "%T is not an iec608705104 tag", parsed)
	return tag
}

// The four shapes of the recovered syntax, plus what each of them resolves to. plc4j's tag handler
// answers every single one of these with ASDU 0 / IOA 0.
func TestTagHandler_ParseTag(t *testing.T) {
	tests := []struct {
		address      string
		wantAsdu     uint16
		wantIoa      uint32
		wantWildcard bool
	}{
		{address: "1/2", wantAsdu: 1, wantIoa: 2},
		{address: "0/0", wantAsdu: 0, wantIoa: 0},
		{address: "65535/16777215", wantAsdu: 65535, wantIoa: 16777215},
		// A split ASDU is the same address spelled octet by octet, low octet first.
		{address: "1/0/2", wantAsdu: 1, wantIoa: 2},
		{address: "0/1/2", wantAsdu: 256, wantIoa: 2},
		{address: "255/255/7", wantAsdu: 65535, wantIoa: 7},
		// A split IOA likewise, low octet first.
		{address: "1/2.0.0", wantAsdu: 1, wantIoa: 2},
		{address: "1/0.1.0", wantAsdu: 1, wantIoa: 256},
		{address: "1/0.0.1", wantAsdu: 1, wantIoa: 65536},
		{address: "1/255.255.255", wantAsdu: 1, wantIoa: 16777215},
		// Both parts split at once.
		{address: "1/1/2.0.1", wantAsdu: 257, wantIoa: 65538},
		// Wildcards cover many points, so they have no single address.
		{address: "*/*", wantWildcard: true},
		{address: "1/*", wantWildcard: true},
		{address: "*/2", wantWildcard: true},
		{address: "1/*/2", wantWildcard: true},
		{address: "*/2.*.0", wantWildcard: true},
	}
	for _, testCase := range tests {
		t.Run(testCase.address, func(t *testing.T) {
			tag := parseTag(t, testCase.address)
			assert.Equal(t, testCase.wantWildcard, tag.HasWildcard())
			asdu, haveAsdu := tag.GetAsduAddress()
			ioa, haveIoa := tag.GetInformationObjectAddress()
			if testCase.wantWildcard {
				assert.False(t, haveAsdu && haveIoa, "a wildcard tag has no single address")
				return
			}
			require.True(t, haveAsdu)
			require.True(t, haveIoa)
			assert.Equal(t, testCase.wantAsdu, asdu)
			assert.Equal(t, testCase.wantIoa, ioa)
		})
	}
}

// An address string has to parse back into the very same tag, because a tag which arrives wrapped in
// a DefaultPlcSubscriptionTag is re-parsed from its address string. plc4j's getAddressString returns
// null, which makes that impossible.
func TestTagHandler_AddressStringRoundTrips(t *testing.T) {
	for _, address := range []string{"1/2", "1/0/2", "1/2.0.0", "1/1/2.0.1", "*/*", "1/*", "*/2.*.0"} {
		t.Run(address, func(t *testing.T) {
			tag := parseTag(t, address)
			assert.Equal(t, address, tag.GetAddressString())
			assert.Equal(t, tag, parseTag(t, tag.GetAddressString()))
		})
	}
}

func TestTagHandler_ParseTagRejects(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{name: "nothing at all", address: ""},
		{name: "an ASDU without an IOA", address: "1"},
		{name: "four slash separated parts", address: "1/2/3/4"},
		{name: "an IOA with two octets", address: "1/2.3"},
		{name: "an IOA with four octets", address: "1/2.3.4.5"},
		{name: "an ASDU beyond two octets", address: "65536/1"},
		{name: "an IOA beyond three octets", address: "1/16777216"},
		{name: "a split ASDU octet beyond one octet", address: "256/0/1"},
		{name: "a split IOA octet beyond one octet", address: "1/256.0.0"},
		{name: "a negative ASDU", address: "-1/2"},
		{name: "a signed ASDU", address: "+1/2"},
		{name: "a hexadecimal IOA", address: "1/0x2"},
		{name: "an empty ASDU", address: "/2"},
		{name: "an empty IOA", address: "1/"},
		{name: "a partial wildcard", address: "1*/2"},
		{name: "whitespace inside an address", address: "1 / 2"},
		{name: "whitespace inside an IOA", address: "1/2. 3.4"},
		{name: "the firmata syntax", address: "digital:1"},
		// plc4j caps the digits of a component at the width of its limit (\d{1,5} for a whole ASDU,
		// \d{1,8} for a whole IOA, \d{1,3} for an octet), so an over padded number is refused there
		// too - an address should be portable in both directions, not just from Go to Java.
		{name: "a zero padded ASDU with more digits than its limit", address: "0000000001/2"},
		{name: "a zero padded IOA with more digits than its limit", address: "1/000000002"},
		{name: "a zero padded octet with more digits than its limit", address: "1/0002.3.4"},
		{name: "a zero padded split ASDU octet", address: "0001/2/3"},
	}
	handler := NewTagHandler(testutils.EnrichOptionsWithOptionsForTesting(t)...)
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tag, err := handler.ParseTag(testCase.address)
			assert.Error(t, err, "%q must not parse", testCase.address)
			assert.Nil(t, tag)
		})
	}
}

// Matching is what a wildcard tag is for: the station reports a concrete point and every handle
// whose tag covers it hears about it.
func TestTag_Matches(t *testing.T) {
	tests := []struct {
		address string
		asdu    uint16
		ioa     uint32
		want    bool
	}{
		{address: "1/2", asdu: 1, ioa: 2, want: true},
		{address: "1/2", asdu: 1, ioa: 3, want: false},
		{address: "1/2", asdu: 2, ioa: 2, want: false},
		// A split spelling addresses the same point as the whole one.
		{address: "1/0/2", asdu: 1, ioa: 2, want: true},
		{address: "1/2.0.0", asdu: 1, ioa: 2, want: true},
		{address: "0/1/2", asdu: 256, ioa: 2, want: true},
		{address: "1/0.1.0", asdu: 1, ioa: 256, want: true},
		// Full wildcards take everything.
		{address: "*/*", asdu: 4242, ioa: 999999, want: true},
		// A wildcard in one part still pins the other.
		{address: "1/*", asdu: 1, ioa: 12345, want: true},
		{address: "1/*", asdu: 2, ioa: 12345, want: false},
		{address: "*/7", asdu: 4242, ioa: 7, want: true},
		{address: "*/7", asdu: 4242, ioa: 8, want: false},
		// A wildcard octet matches only in its own position.
		{address: "1/*/9", asdu: 0x0001, ioa: 9, want: true},
		{address: "1/*/9", asdu: 0xFF01, ioa: 9, want: true},
		{address: "1/*/9", asdu: 0x0002, ioa: 9, want: false},
		{address: "1/5.*.0", asdu: 1, ioa: 0x000105, want: true},
		{address: "1/5.*.0", asdu: 1, ioa: 0x00FF05, want: true},
		{address: "1/5.*.0", asdu: 1, ioa: 0x010005, want: false},
		{address: "1/5.*.0", asdu: 1, ioa: 0x000106, want: false},
	}
	for _, testCase := range tests {
		t.Run(testCase.address, func(t *testing.T) {
			tag := parseTag(t, testCase.address)
			assert.Equal(t, testCase.want, tag.Matches(testCase.asdu, testCase.ioa))
		})
	}
}

// The address of an incoming point has to be spelled the way the handler parses it back, because it
// is what the subscription event reports as its address.
func TestPointAddressString(t *testing.T) {
	assert.Equal(t, "10/13", pointAddressString(10, 13))
	assert.Equal(t, "0/0", pointAddressString(0, 0))
	assert.Equal(t, "65535/16777215", pointAddressString(65535, 16777215))
	assert.Equal(t, parseTag(t, "10/13"), NewTag(10, 13))
}

// Surrounding whitespace is trimmed rather than refused, matching the plc4j tag, so that the two
// languages accept the same set of addresses.
func TestTagHandler_ParseTagTrimsSurroundingWhitespace(t *testing.T) {
	tag := parseTag(t, "  10/13\t")
	assert.Equal(t, "10/13", tag.GetAddressString())
	assert.Equal(t, NewTag(10, 13), tag)
}

func TestTagHandler_ParseQueryIsUnsupported(t *testing.T) {
	query, err := NewTagHandler().ParseQuery("*/*")
	assert.Error(t, err)
	assert.Nil(t, query)
}

// Zero padding up to the digits plc4j's regex allows is still accepted, and addresses the same
// point as the unpadded spelling - only the canonical GetAddressString drops the padding.
func TestTagHandler_AcceptsZeroPaddingUpToTheDigitLimit(t *testing.T) {
	tests := []struct {
		address   string
		canonical string
	}{
		{address: "00001/2", canonical: "1/2"},
		{address: "1/00000002", canonical: "1/2"},
		{address: "001/2/003.4.005", canonical: "1/2/3.4.5"},
	}
	for _, testCase := range tests {
		t.Run(testCase.address, func(t *testing.T) {
			tag := parseTag(t, testCase.address)
			assert.Equal(t, testCase.canonical, tag.GetAddressString())
			assert.Equal(t, parseTag(t, testCase.canonical), tag)
		})
	}
}
