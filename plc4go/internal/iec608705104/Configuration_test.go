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
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFromOptions(t *testing.T) {
	tests := []struct {
		name              string
		connectionOptions map[string][]string
		want              Configuration
		wantErr           bool
	}{
		{
			name:              "the defaults",
			connectionOptions: map[string][]string{},
			want:              Configuration{requestTimeout: 4 * time.Second, ackThreshold: 8},
		},
		{
			name:              "a request timeout in milliseconds",
			connectionOptions: map[string][]string{"request-timeout-ms": {"1500"}},
			want:              Configuration{requestTimeout: 1500 * time.Millisecond, ackThreshold: 8},
		},
		{
			name:              "a smaller acknowledgement window",
			connectionOptions: map[string][]string{"ack-threshold": {"4"}},
			want:              Configuration{requestTimeout: 4 * time.Second, ackThreshold: 4},
		},
		{
			name:              "the largest window a sequence number can express",
			connectionOptions: map[string][]string{"ack-threshold": {"32767"}},
			want:              Configuration{requestTimeout: 4 * time.Second, ackThreshold: 32767},
		},
		{
			name:              "an option given twice takes the first",
			connectionOptions: map[string][]string{"ack-threshold": {"4", "6"}},
			want:              Configuration{requestTimeout: 4 * time.Second, ackThreshold: 4},
		},
		{
			name:              "an empty option list is no option at all",
			connectionOptions: map[string][]string{"ack-threshold": {}},
			want:              Configuration{requestTimeout: 4 * time.Second, ackThreshold: 8},
		},
		{name: "a request timeout which isn't a number", connectionOptions: map[string][]string{"request-timeout-ms": {"soon"}}, wantErr: true},
		{name: "a request timeout of zero", connectionOptions: map[string][]string{"request-timeout-ms": {"0"}}, wantErr: true},
		{name: "a window which isn't a number", connectionOptions: map[string][]string{"ack-threshold": {"lots"}}, wantErr: true},
		{name: "a window of zero", connectionOptions: map[string][]string{"ack-threshold": {"0"}}, wantErr: true},
		{name: "a window past the sequence number", connectionOptions: map[string][]string{"ack-threshold": {"32768"}}, wantErr: true},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			configuration, err := ParseFromOptions(zerolog.Nop(), testCase.connectionOptions)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, configuration)
		})
	}
}
