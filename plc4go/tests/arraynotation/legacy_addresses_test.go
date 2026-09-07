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

// Package arraynotation holds the cross-driver checks for the unified array notation. They live
// outside the driver packages because what they assert is a property of the whole binding: the
// release notes tell an upgrading user that every address whose meaning changed is either
// rejected with its replacement named, or listed as one of the two silent changes. That claim is
// only true if it holds for every driver at once, which is what these tests check.
package arraynotation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/internal/ads"
	"github.com/apache/plc4x/plc4go/internal/eip"
	"github.com/apache/plc4x/plc4go/internal/firmata"
	"github.com/apache/plc4x/plc4go/internal/knxnetip"
	"github.com/apache/plc4x/plc4go/internal/modbus"
	"github.com/apache/plc4x/plc4go/internal/s7"
	"github.com/apache/plc4x/plc4go/internal/simulated"
	"github.com/apache/plc4x/plc4go/internal/slmp"
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

type parseTag func(string) (apiModel.PlcTag, error)

// Every pre-migration address in the release notes must be rejected, and the rejection must name
// the address to write instead. A rejection alone is not enough: the point of moving the
// brackets was that an upgrade reports the change rather than quietly returning different data,
// and a user with a configuration full of old addresses needs to be told what to write.
func TestEveryLegacyAddressIsRejectedWithItsReplacement(t *testing.T) {
	for _, c := range []struct {
		driver      string
		parse       parseTag
		address     string
		replacement string
	}{
		{"s7", s7.NewTagHandler().ParseTag, "%M100:INT[10]", "%M100[0..9]:INT"},
		{"s7 string", s7.NewTagHandler().ParseTag, "%DB69.DBX68:WSTRING[3]", "%DB69.DBX68[0..2]:WSTRING"},
		{"modbus", modbus.NewTagHandler().ParseTag, "holding-register:1:INT[4]", "holding-register:1[0..3]:INT"},
		{"slmp", slmp.NewTagHandler().ParseTag, "D100:INT[4]", "D100[0..3]:INT"},
		{"eip", eip.NewTagHandler().ParseTag, "%rate:DINT:4", "%rate[0..3]:DINT"},
		{"simulated", simulated.NewTagHandler().ParseTag, "RANDOM/foo:INT[4]", "RANDOM/foo[0..3]:INT"},
		{"knxnetip memory", knxnetip.NewTagHandler().ParseTag, "1.2.3#4B1C:UINT[4]", "1.2.3#4B1C[0..3]:UINT"},
		{"ads direct", ads.NewTagHandler().ParseTag, "0x4020/0:DINT[4]", "0x4020/0[0..3]:DINT"},
		// The start-and-count form was a plc4go extension with no counterpart in plc4j.
		{"ads start-and-count", ads.NewTagHandler().ParseTag, "MAIN.g_arr[2:4]", "MAIN.g_arr[2..5]"},
	} {
		t.Run(c.driver, func(t *testing.T) {
			_, err := c.parse(c.address)
			require.Error(t, err, c.address)
			assert.Contains(t, err.Error(), c.replacement,
				"the rejection must name the address to write instead")
		})
	}
}

// The two addresses that parse before and after, and only change meaning. Neither can be
// rejected, so the release notes carry them - and these tests are what keeps that list honest.
func TestTheSilentChangesAreExactlyTheTwoThatAreDocumented(t *testing.T) {
	// Firmata: [n] was a run of n pins and is now the pin at index n.
	pin, err := firmata.NewTagHandler().ParseTag("digital:2[3]")
	require.NoError(t, err)
	assert.Equal(t, "digital:5", pin.GetAddressString(), "pin 5, not three pins from pin 2")
	assert.Empty(t, pin.GetArrayInfo(), "one pin is a scalar")

	// ADS: [n] was a count of n elements and is now the element at index n.
	element, err := ads.NewTagHandler().ParseTag("MAIN.g_arr[3]")
	require.NoError(t, err)
	assert.Equal(t, "MAIN.g_arr[3]", element.GetAddressString())
	assert.Empty(t, element.GetArrayInfo(), "one element is a scalar, not three elements")
}

// The same address selects the same elements in plc4go as in plc4j. The two bindings share a
// specification rather than code, so this is asserted case by case; the numbers here are the
// ones the Java parity tests assert.
func TestOneAddressMeansOneThingAcrossDrivers(t *testing.T) {
	for _, c := range []struct {
		driver  string
		parse   parseTag
		address string
	}{
		{"s7", s7.NewTagHandler().ParseTag, "%M100[0..7]:INT"},
		{"modbus", modbus.NewTagHandler().ParseTag, "holding-register:1[0..7]:INT"},
		{"slmp", slmp.NewTagHandler().ParseTag, "D100[0..7]:INT"},
		{"eip", eip.NewTagHandler().ParseTag, "%rate[0..7]:DINT"},
		{"simulated", simulated.NewTagHandler().ParseTag, "RANDOM/foo[0..7]:INT"},
		{"firmata", firmata.NewTagHandler().ParseTag, "digital:0[0..7]"},
		{"ads", ads.NewTagHandler().ParseTag, "MAIN.g_arr[0..7]"},
	} {
		t.Run(c.driver, func(t *testing.T) {
			tag, err := c.parse(c.address)
			require.NoError(t, err)

			dimensions := tag.GetArrayInfo()
			require.Len(t, dimensions, 1, "one dimension")
			assert.Equal(t, uint32(8), dimensions[0].GetSize(), "eight elements")
			assert.True(t, dimensions[0].IsRange(), "written as a range")

			reparsed, err := c.parse(tag.GetAddressString())
			require.NoError(t, err, "a rendered address must parse back")
			assert.Equal(t, tag, reparsed)
		})
	}
}
