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

package knxnetip

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
)

func Test_TagHandler_ParseQuery(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		want        apiModel.PlcQuery
		wantErrText string
	}{
		{
			name:  "plain device query",
			query: "1.2.3",
			want:  NewDeviceQuery("1", "2", "3"),
		},
		{
			name:  "device query with wildcards",
			query: "*.*.*",
			want:  NewDeviceQuery("*", "*", "*"),
		},
		{
			name:  "device query with a range",
			query: "1.2.[1-9]",
			want:  NewDeviceQuery("1", "2", "[1-9]"),
		},
		{
			// Without this form Browser.executeCommunicationObjectQuery is unreachable
			// through the public browse api.
			name:  "communication object query",
			query: "1.2.3#com-obj",
			want:  NewCommunicationObjectQuery(1, 2, 3),
		},
		{
			name:        "garbage",
			query:       "not-an-address",
			wantErrText: "Invalid address format for query 'not-an-address'",
		},
		{
			name:        "communication object query doesn't accept wildcards",
			query:       "*.*.*#com-obj",
			wantErrText: "Invalid address format for query '*.*.*#com-obj'",
		},
	}
	handler := NewTagHandler()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			query, err := handler.ParseQuery(test.query)
			if test.wantErrText != "" {
				assert.EqualError(t, err, test.wantErrText)
				assert.Nil(t, query)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, query)
			assert.Equal(t, test.query, query.GetQueryString(), "the query has to round-trip")
		})
	}
}

// Test_TagHandler_ParseQuery_matchesParseTag makes sure both entry-points produce the
// same communication-object query, as the browse api goes through ParseQuery while the
// read/write api goes through ParseTag.
func Test_TagHandler_ParseQuery_matchesParseTag(t *testing.T) {
	handler := NewTagHandler()
	tag, err := handler.ParseTag("1.2.3#com-obj")
	require.NoError(t, err)
	query, err := handler.ParseQuery("1.2.3#com-obj")
	require.NoError(t, err)
	assert.Equal(t, tag, query)
}

func Test_DeviceQuery_identity(t *testing.T) {
	query := NewDeviceQuery("1", "2", "3")
	assert.Equal(t, "1.2.3", query.GetAddressString())
	assert.Equal(t, "1.2.3", query.GetQueryString())
	assert.Equal(t, "knx.DeviceQuery{1.2.3}", query.String())
	assert.Equal(t, driverModel.NewKnxAddress(1, 2, 3), query.toKnxAddress())
}

func Test_DeviceQuery_toKnxAddress_patterns(t *testing.T) {
	for _, query := range []DeviceQuery{
		NewDeviceQuery("*", "2", "3"),
		NewDeviceQuery("1", "*", "3"),
		NewDeviceQuery("1", "2", "[1-9]"),
	} {
		t.Run(query.GetAddressString(), func(t *testing.T) {
			assert.Nil(t, query.toKnxAddress(), "a pattern has to be expanded first")
		})
	}
}

func Test_CommunicationObjectQuery_identity(t *testing.T) {
	query := NewCommunicationObjectQuery(1, 2, 3)
	assert.Equal(t, "1.2.3#com-obj", query.GetAddressString())
	assert.Equal(t, "1.2.3#com-obj", query.GetQueryString())
	assert.Equal(t, "knx.CommunicationObjectQuery{1.2.3#com-obj}", query.String())
	assert.Equal(t, driverModel.NewKnxAddress(1, 2, 3), query.toKnxAddress())
}

// Test_Connection_BrowseRequestBuilder_acceptsBothQueryForms makes sure both browse
// flavours survive the public request-builder.
func Test_Connection_BrowseRequestBuilder_acceptsBothQueryForms(t *testing.T) {
	connection := newSubscriberTestConnection(t)
	builder := connection.BrowseRequestBuilder()
	builder.AddQuery("devices", "1.2.*")
	builder.AddQuery("comObjects", "1.2.3#com-obj")
	browseRequest, err := builder.Build()
	require.NoError(t, err)

	assert.IsType(t, DeviceQuery{}, browseRequest.GetQuery("devices"))
	assert.IsType(t, CommunicationObjectQuery{}, browseRequest.GetQuery("comObjects"))
}

// The two device address forms carry a real element count, so they take the shared notation.
// Nothing here covered them before, which is how the browser came to build addresses in a
// spelling the handler no longer accepts - a failure that only showed up against a device.
func TestTagHandler_DeviceAddressesUseTheSharedNotation(t *testing.T) {
	handler := NewTagHandler()

	for _, address := range []string{
		"1.2.3#4B1C:UINT",        // a memory address, one element
		"1.2.3#4B1C[0..7]:UINT",  // eight of them, the selection before the type
		"1.2.3#4B1C[0..3]:USINT", // the form the browser builds when it walks the tables
		"1.2.3#11/1/5[0..3]",     // a property address; no type suffix, so the selection ends it
		"1.2.3#3/23/5",           // a property address with no selection
	} {
		t.Run(address, func(t *testing.T) {
			tag, err := handler.ParseTag(address)
			require.NoError(t, err)
			assert.Equal(t, address, tag.GetAddressString())

			reparsed, err := handler.ParseTag(tag.GetAddressString())
			require.NoError(t, err)
			assert.Equal(t, tag, reparsed)
		})
	}

	// An omitted property index defaults to 1 and is spelled out when the tag is rendered, which
	// still re-parses to the same tag.
	implied, err := handler.ParseTag("1.2.3#11/1[0..3]")
	require.NoError(t, err)
	assert.Equal(t, "1.2.3#11/1/1[0..3]", implied.GetAddressString())
	reparsed, err := handler.ParseTag(implied.GetAddressString())
	require.NoError(t, err)
	assert.Equal(t, implied, reparsed)
}

// A property is read with a start index and a count, and the property index in the address is
// that start index - so a selection that starts past the first element moves it.
func TestTagHandler_APropertySelectionMovesTheStartIndex(t *testing.T) {
	handler := NewTagHandler()

	shifted, err := handler.ParseTag("1.2.3#11/1/1[4..7]")
	require.NoError(t, err)
	equivalent, err := handler.ParseTag("1.2.3#11/1/5[0..3]")
	require.NoError(t, err)
	assert.Equal(t, equivalent, shifted)

	// A memory address has no such mapping: what one element occupies depends on the datapoint
	// type, which is measured in bits, so an offset is reported rather than guessed at.
	_, err = handler.ParseTag("1.2.3#4B1C[4..7]:UINT")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must start at the first element")
}

// The group-address forms are not array addresses at all: their brackets hold a set of group
// addresses to match, so they are left exactly as they were.
func TestTagHandler_GroupAddressBracketsAreAMatchExpression(t *testing.T) {
	handler := NewTagHandler()

	for _, address := range []string{"1/[2-3]/4:BOOL", "[1,3]/2/[4-6]:BOOL", "*/*/*:BOOL"} {
		_, err := handler.ParseTag(address)
		assert.NoError(t, err, address)
	}
}
