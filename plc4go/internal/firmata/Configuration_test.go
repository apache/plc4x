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
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFromOptions(t *testing.T) {
	tests := []struct {
		name              string
		connectionOptions map[string][]string
		want              time.Duration
		wantErr           bool
	}{
		{
			name:              "no options at all",
			connectionOptions: map[string][]string{},
			want:              defaultRequestTimeout,
		},
		{
			// plc4j spells the timeout in milliseconds (@IntDefaultValue(10_000)).
			name:              "a request timeout in milliseconds",
			connectionOptions: map[string][]string{"request-timeout": {"250"}},
			want:              250 * time.Millisecond,
		},
		{
			name:              "an empty option falls back to the default",
			connectionOptions: map[string][]string{"request-timeout": {}},
			want:              defaultRequestTimeout,
		},
		{
			name:              "the first of several values wins",
			connectionOptions: map[string][]string{"request-timeout": {"250", "500"}},
			want:              250 * time.Millisecond,
		},
		{
			name:              "a timeout which isn't a number",
			connectionOptions: map[string][]string{"request-timeout": {"soon"}},
			wantErr:           true,
		},
		{
			name:              "a timeout of zero would never wait",
			connectionOptions: map[string][]string{"request-timeout": {"0"}},
			wantErr:           true,
		},
		{
			name:              "a negative timeout",
			connectionOptions: map[string][]string{"request-timeout": {"-1"}},
			wantErr:           true,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			configuration, err := ParseFromOptions(zerolog.Nop(), testCase.connectionOptions)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, configuration.requestTimeout)
		})
	}
}
