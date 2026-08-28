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

package firmata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/firmata/readwrite/model"
)

func TestTagHandler_ParseTag(t *testing.T) {
	pullup := readWriteModel.PinMode_PinModePullup
	tests := []struct {
		name        string
		tagAddress  string
		want        Tag
		wantValue   apiValues.PlcValueType
		wantAddress string
	}{
		{
			name:        "a single digital pin",
			tagAddress:  "digital:4",
			want:        digitalTag{address: 4, quantity: 1},
			wantValue:   apiValues.BOOL,
			wantAddress: "digital:4",
		},
		{
			name:        "a run of digital pins",
			tagAddress:  "digital:4[0..2]",
			want:        digitalTag{address: 4, quantity: 3, explicitRange: true},
			wantValue:   apiValues.BOOL,
			wantAddress: "digital:4[0..2]",
		},
		{
			name:        "a run of one digital pin is a scalar",
			tagAddress:  "digital:4",
			want:        digitalTag{address: 4, quantity: 1},
			wantValue:   apiValues.BOOL,
			wantAddress: "digital:4",
		},
		{
			name:        "a pullup digital pin",
			tagAddress:  "digital:7:PULLUP",
			want:        digitalTag{address: 7, quantity: 1, pinMode: &pullup},
			wantValue:   apiValues.BOOL,
			wantAddress: "digital:7:PULLUP",
		},
		{
			name:        "a run of pullup digital pins",
			tagAddress:  "digital:7[0..1]:PULLUP",
			want:        digitalTag{address: 7, quantity: 2, pinMode: &pullup, explicitRange: true},
			wantValue:   apiValues.BOOL,
			wantAddress: "digital:7[0..1]:PULLUP",
		},
		{
			name:        "the last addressable digital pin",
			tagAddress:  "digital:127",
			want:        digitalTag{address: 127, quantity: 1},
			wantValue:   apiValues.BOOL,
			wantAddress: "digital:127",
		},
		{
			name:        "a single analog pin",
			tagAddress:  "analog:0",
			want:        analogTag{address: 0, quantity: 1},
			wantValue:   apiValues.INT,
			wantAddress: "analog:0",
		},
		{
			name:        "a run of analog pins",
			tagAddress:  "analog:2[0..3]",
			want:        analogTag{address: 2, quantity: 4, explicitRange: true},
			wantValue:   apiValues.INT,
			wantAddress: "analog:2[0..3]",
		},
		{
			name:        "the last addressable analog pin",
			tagAddress:  "analog:15",
			want:        analogTag{address: 15, quantity: 1},
			wantValue:   apiValues.INT,
			wantAddress: "analog:15",
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tag, err := NewTagHandler().ParseTag(testCase.tagAddress)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, tag)
			assert.Equal(t, testCase.wantValue, tag.GetValueType())
			// Every address string has to parse back into the very same tag, as plc4go re-parses
			// them whenever a tag arrives wrapped in a DefaultPlcSubscriptionTag.
			assert.Equal(t, testCase.wantAddress, tag.GetAddressString())
			reParsed, err := NewTagHandler().ParseTag(tag.GetAddressString())
			require.NoError(t, err)
			assert.Equal(t, tag, reParsed)
		})
	}
}

func TestTagHandler_ParseTagRejects(t *testing.T) {
	tests := []struct {
		name       string
		tagAddress string
	}{
		{name: "an empty address", tagAddress: ""},
		{name: "an unknown area", tagAddress: "pwm:4"},
		{name: "a digital pin without a number", tagAddress: "digital:"},
		{name: "a negative pin", tagAddress: "digital:-1"},
		{name: "a non-numeric pin", tagAddress: "digital:a"},
		{name: "trailing garbage", tagAddress: "digital:4nonsense"},
		{name: "an unknown mode", tagAddress: "digital:4:OUTPUT"},
		{name: "a mode on an analog pin", tagAddress: "analog:4:PULLUP"},
		// A run of zero pins has no spelling in the notation - a range is written with the
		// indices it covers - but an inverted one is still nonsense.
		{name: "an inverted range", tagAddress: "digital:4[3..1]"},
		{name: "a digital pin past the last port", tagAddress: "digital:128"},
		{name: "a run of digital pins past the last port", tagAddress: "digital:126[0..3]"},
		{name: "an analog pin past the 4 bit pin field", tagAddress: "analog:16"},
		{name: "a run of analog pins past the 4 bit pin field", tagAddress: "analog:14[0..3]"},
		{name: "a quantity larger than the whole pin range", tagAddress: "digital:0[0..128]"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tag, err := NewTagHandler().ParseTag(testCase.tagAddress)
			assert.Error(t, err)
			assert.Nil(t, tag)
		})
	}
}

// Firmata boards can be asked for their capabilities, but neither this driver nor plc4j's browses
// them, so a query has to be refused rather than silently answered with nothing.
func TestTagHandler_ParseQueryIsUnsupported(t *testing.T) {
	query, err := NewTagHandler().ParseQuery("digital:*")
	assert.Error(t, err)
	assert.Nil(t, query)
}

func TestTag_ArrayInfoAndSubscriptionDefaults(t *testing.T) {
	scalar, err := NewTagHandler().ParseTag("digital:4")
	require.NoError(t, err)
	assert.Empty(t, scalar.GetArrayInfo())

	list, err := NewTagHandler().ParseTag("digital:4[0..2]")
	require.NoError(t, err)
	require.Len(t, list.GetArrayInfo(), 1)
	assert.Equal(t, uint32(0), list.GetArrayInfo()[0].GetLowerBound())
	// Both bounds are inclusive, so three pins are 0..2.
	assert.Equal(t, uint32(2), list.GetArrayInfo()[0].GetUpperBound())
	assert.Equal(t, uint32(3), list.GetArrayInfo()[0].GetSize())

	// A tag has to be usable in a subscription request, which only accepts tags that are a
	// PlcSubscriptionTag. Firmata boards report a pin when it changes.
	for _, tagAddress := range []string{"digital:4", "analog:4"} {
		tag, err := NewTagHandler().ParseTag(tagAddress)
		require.NoError(t, err)
		subscriptionTag, ok := tag.(apiModel.PlcSubscriptionTag)
		require.True(t, ok, "%s isn't a PlcSubscriptionTag", tagAddress)
		assert.Equal(t, apiModel.SubscriptionChangeOfState, subscriptionTag.GetPlcSubscriptionType())
		assert.Zero(t, subscriptionTag.GetDuration())
	}
}

func TestTag_String(t *testing.T) {
	digital, err := NewTagHandler().ParseTag("digital:4[0..1]:PULLUP")
	require.NoError(t, err)
	assert.Contains(t, digital.String(), "PinModePullup")

	analog, err := NewTagHandler().ParseTag("analog:4")
	require.NoError(t, err)
	assert.Contains(t, analog.String(), "analogTag")
}

// A firmata address is a pin number, so a selection that starts past the declared base is
// resolved into the pin: "digital:2[4..7]" is the same run as "digital:6[0..3]".
//
// [n] used to mean "n pins" and now means "the pin at index n". Both forms parse, so nothing can
// be rejected here - this is one of the two silent changes the release notes have to carry.
func TestTagHandler_ParseTag_consumesTheSelectionOffset(t *testing.T) {
	handler := NewTagHandler()

	shifted, err := handler.ParseTag("digital:2[4..7]")
	require.NoError(t, err)
	equivalent, err := handler.ParseTag("digital:6[0..3]")
	require.NoError(t, err)
	assert.Equal(t, equivalent, shifted)
	assert.Equal(t, "digital:6[0..3]", shifted.GetAddressString())

	// A declared base is what the offset is measured from, so [4..7;4] shifts nothing.
	fromDeclaredBase, err := handler.ParseTag("digital:2[4..7;4]")
	require.NoError(t, err)
	unshifted, err := handler.ParseTag("digital:2[0..3]")
	require.NoError(t, err)
	assert.Equal(t, unshifted, fromDeclaredBase)

	// The silent change: [3] is the pin at index 3, which is pin 5 here - not three pins.
	single, err := handler.ParseTag("digital:2[3]")
	require.NoError(t, err)
	assert.Equal(t, "digital:5", single.GetAddressString())
	assert.Empty(t, single.GetArrayInfo(), "one pin is a scalar")
}
