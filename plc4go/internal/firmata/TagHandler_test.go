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
			tagAddress:  "digital:4[3]",
			want:        digitalTag{address: 4, quantity: 3},
			wantValue:   apiValues.BOOL,
			wantAddress: "digital:4[3]",
		},
		{
			name:        "a run of one digital pin is a scalar",
			tagAddress:  "digital:4[1]",
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
			tagAddress:  "digital:7[2]:PULLUP",
			want:        digitalTag{address: 7, quantity: 2, pinMode: &pullup},
			wantValue:   apiValues.BOOL,
			wantAddress: "digital:7[2]:PULLUP",
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
			tagAddress:  "analog:2[4]",
			want:        analogTag{address: 2, quantity: 4},
			wantValue:   apiValues.INT,
			wantAddress: "analog:2[4]",
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
		{name: "a quantity of zero", tagAddress: "digital:4[0]"},
		{name: "a digital pin past the last port", tagAddress: "digital:128"},
		{name: "a run of digital pins past the last port", tagAddress: "digital:126[4]"},
		{name: "an analog pin past the 4 bit pin field", tagAddress: "analog:16"},
		{name: "a run of analog pins past the 4 bit pin field", tagAddress: "analog:14[4]"},
		{name: "a quantity larger than the whole pin range", tagAddress: "digital:0[129]"},
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

	list, err := NewTagHandler().ParseTag("digital:4[3]")
	require.NoError(t, err)
	require.Len(t, list.GetArrayInfo(), 1)
	assert.Equal(t, uint32(0), list.GetArrayInfo()[0].GetLowerBound())
	assert.Equal(t, uint32(3), list.GetArrayInfo()[0].GetUpperBound())

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
	digital, err := NewTagHandler().ParseTag("digital:4[2]:PULLUP")
	require.NoError(t, err)
	assert.Contains(t, digital.String(), "PinModePullup")

	analog, err := NewTagHandler().ParseTag("analog:4")
	require.NoError(t, err)
	assert.Contains(t, analog.String(), "analogTag")
}
