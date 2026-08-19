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
