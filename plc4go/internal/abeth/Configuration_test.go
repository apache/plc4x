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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func TestParseFromOptions(t *testing.T) {
	tests := []struct {
		name    string
		options map[string][]string
		want    Configuration
		wantErr bool
	}{
		{
			name:    "no options at all",
			options: map[string][]string{},
			want:    Configuration{station: 0, requestTimeout: 10 * time.Second},
		},
		{
			name:    "a station",
			options: map[string][]string{"station": {"8"}},
			want:    Configuration{station: 8, requestTimeout: 10 * time.Second},
		},
		{
			name:    "a request timeout in milliseconds",
			options: map[string][]string{"request-timeout-ms": {"2500"}},
			want:    Configuration{station: 0, requestTimeout: 2500 * time.Millisecond},
		},
		{
			name:    "both options",
			options: map[string][]string{"station": {"255"}, "request-timeout-ms": {"1"}},
			want:    Configuration{station: 255, requestTimeout: time.Millisecond},
		},
		{
			name:    "an empty option value falls back to the default",
			options: map[string][]string{"station": {}},
			want:    Configuration{station: 0, requestTimeout: 10 * time.Second},
		},
		{
			// The station is the DF1 destination address, a single byte on the wire.
			name:    "a station that doesn't fit into a byte",
			options: map[string][]string{"station": {"256"}},
			wantErr: true,
		},
		{
			name:    "a non-numeric station",
			options: map[string][]string{"station": {"eight"}},
			wantErr: true,
		},
		{
			name:    "a zero request timeout",
			options: map[string][]string{"request-timeout-ms": {"0"}},
			wantErr: true,
		},
		{
			name:    "a non-numeric request timeout",
			options: map[string][]string{"request-timeout-ms": {"soon"}},
			wantErr: true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			configuration, err := ParseFromOptions(testutils.ProduceTestingLogger(t), testCase.options)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, configuration)
		})
	}
}

func TestParseFromOptionsWarnsAboutDuplicates(t *testing.T) {
	// A duplicated option isn't an error, the first value wins - the same as every other plc4go
	// driver does it.
	configuration, err := ParseFromOptions(testutils.ProduceTestingLogger(t),
		map[string][]string{"station": {"3", "4"}})
	require.NoError(t, err)
	assert.Equal(t, uint8(3), configuration.station)
}
