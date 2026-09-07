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

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

// The accepted and rejected addresses are the cases of plc4j's SymbolicUmasTagTest, plus the shapes
// its regex documents but doesn't exercise.
func TestTagHandler_ParseTag(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{name: "a plain global variable", address: "MyVariable"},
		{name: "an identifier may start with an underscore", address: "_private"},
		{name: "a nested struct member", address: "g_plant.meta.r32"},
		{name: "an array element", address: "g_arrInt[3]"},
		{name: "mixed struct and array access", address: "g_plant.items[2].value"},
		{name: "several indices in a row", address: "g_matrix[1][2]"},
		{name: "a name can't start with a digit", address: "9bad-name", wantErr: true},
		{name: "a name can't contain spaces", address: "not a valid name", wantErr: true},
		{name: "an index has to be numeric", address: "g_arr[x]", wantErr: true},
		{name: "an index has to be closed", address: "g_arr[3", wantErr: true},
		{name: "a dot has to be followed by a member", address: "g_plant.", wantErr: true},
		{name: "an empty address is not a symbol", address: "", wantErr: true},
	}
	handler := NewTagHandler()
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			tag, err := handler.ParseTag(testCase.address)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tag)
				return
			}
			require.NoError(t, err)
			umasTag, ok := tag.(PlcTag)
			require.True(t, ok, "%T is not a UMAS tag", tag)
			assert.Equal(t, testCase.address, umasTag.GetSymbolicAddress())
			// The address string has to parse back, because the polling subscriber rebuilds its read
			// requests from it.
			assert.Equal(t, testCase.address, umasTag.GetAddressString())
			reparsed, err := handler.ParseTag(umasTag.GetAddressString())
			require.NoError(t, err)
			assert.Equal(t, tag, reparsed)
		})
	}
}

// A tag which came out of an address knows no type: only the data dictionary does, and the writer
// looks it up there. The value handler depends on this staying NULL.
func TestTag_AnAddressCarriesNoType(t *testing.T) {
	tag, err := NewTagHandler().ParseTag("g_r32")
	require.NoError(t, err)
	assert.Equal(t, apiValues.NULL, tag.GetValueType())
	assert.Empty(t, tag.GetArrayInfo())
}

func TestTag_TheBrowserCanGiveATagItsType(t *testing.T) {
	arrayInfo := []apiModel.ArrayInfo{&spiModel.DefaultArrayInfo{LowerBound: 0, UpperBound: 3}}
	tag := NewTagWithType("g_arrInt", apiValues.INT, arrayInfo)
	assert.Equal(t, apiValues.INT, tag.GetValueType())
	assert.Equal(t, arrayInfo, tag.GetArrayInfo())
	assert.Equal(t, "g_arrInt", tag.GetAddressString())
}

// The subscription-request builder only accepts tags which are subscription tags, and polling can
// only emulate cyclic and change-of-state subscriptions.
func TestTag_IsASubscriptionTag(t *testing.T) {
	tag, err := NewTagHandler().ParseTag("g_b16")
	require.NoError(t, err)
	subscriptionTag, ok := tag.(apiModel.PlcSubscriptionTag)
	require.True(t, ok, "%T is not a subscription tag", tag)
	assert.Equal(t, apiModel.SubscriptionChangeOfState, subscriptionTag.GetPlcSubscriptionType())
	assert.Equal(t, time.Duration(0), subscriptionTag.GetDuration())
}

func TestTag_Serializes(t *testing.T) {
	serialized, err := NewTagWithType("g_r32", apiValues.REAL, nil).Serialize()
	require.NoError(t, err)
	assert.NotEmpty(t, serialized)
}

func TestTagHandler_ParseQuery(t *testing.T) {
	handler := NewTagHandler()

	t.Run("an empty query is refused", func(t *testing.T) {
		query, err := handler.ParseQuery("   ")
		assert.Error(t, err)
		assert.Nil(t, query)
	})

	t.Run("a query keeps its own string", func(t *testing.T) {
		query, err := handler.ParseQuery("g_*")
		require.NoError(t, err)
		assert.Equal(t, "g_*", query.GetQueryString())
	})
}

func TestSymbolQuery_Matches(t *testing.T) {
	tests := []struct {
		query  string
		symbol string
		want   bool
	}{
		{query: "*", symbol: "g_r32", want: true},
		{query: "*", symbol: "", want: true},
		{query: "g_*", symbol: "g_r32", want: true},
		{query: "g_*", symbol: "h_r32", want: false},
		{query: "*_r32", symbol: "g_r32", want: true},
		{query: "g_r??", symbol: "g_r32", want: true},
		{query: "g_r?", symbol: "g_r32", want: false},
		{query: "g_r32", symbol: "g_r32", want: true},
		{query: "g_r32", symbol: "g_r321", want: false},
		// Symbol names are folded to lower case throughout the driver, and the query matches
		// regardless of how the user spelled it.
		{query: "G_R32", symbol: "g_r32", want: true},
		// A regex metacharacter in a query is matched literally, not as a pattern.
		{query: "g_r.2", symbol: "g_r32", want: false},
		{query: "g_r.2", symbol: "g_r.2", want: true},
	}
	for _, testCase := range tests {
		t.Run(testCase.query+" vs "+testCase.symbol, func(t *testing.T) {
			query, err := NewSymbolQuery(testCase.query)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, query.Matches(testCase.symbol))
		})
	}
}
