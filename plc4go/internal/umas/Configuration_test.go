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

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The defaults are plc4j's UmasConfiguration annotations: unit-identifier 0, request-timeout 4000 ms,
// max-frame-size 65535.
func TestDefaultConfiguration(t *testing.T) {
	configuration := DefaultConfiguration()
	assert.Equal(t, uint8(0), configuration.unitIdentifier)
	assert.Equal(t, 4*time.Second, configuration.requestTimeout)
	assert.Equal(t, uint16(65535), configuration.maxFrameSize)
}

func TestParseFromOptions(t *testing.T) {
	tests := []struct {
		name    string
		options map[string][]string
		want    Configuration
		wantErr bool
	}{
		{
			name:    "no options at all is the default",
			options: map[string][]string{},
			want:    DefaultConfiguration(),
		},
		{
			name:    "the unit identifier is a single byte",
			options: map[string][]string{"unit-identifier": {"3"}},
			want:    Configuration{unitIdentifier: 3, requestTimeout: defaultRequestTimeout, maxFrameSize: defaultMaxFrameSize},
		},
		{
			name:    "a unit identifier above 255 doesn't fit the wire field",
			options: map[string][]string{"unit-identifier": {"256"}},
			wantErr: true,
		},
		{
			name:    "the request timeout is spelled in milliseconds",
			options: map[string][]string{"request-timeout": {"1500"}},
			want:    Configuration{unitIdentifier: 0, requestTimeout: 1500 * time.Millisecond, maxFrameSize: defaultMaxFrameSize},
		},
		{
			name:    "a request timeout of zero would never wait for an answer",
			options: map[string][]string{"request-timeout": {"0"}},
			wantErr: true,
		},
		{
			name:    "a non numeric request timeout is refused",
			options: map[string][]string{"request-timeout": {"soon"}},
			wantErr: true,
		},
		{
			name:    "the max frame size can be lowered",
			options: map[string][]string{"max-frame-size": {"260"}},
			want:    Configuration{unitIdentifier: 0, requestTimeout: defaultRequestTimeout, maxFrameSize: 260},
		},
		{
			// The echo request of the handshake sends maxFrameSize - 3 bytes, so anything below 4
			// has no payload to send at all.
			name:    "a max frame size below the echo overhead is refused",
			options: map[string][]string{"max-frame-size": {"3"}},
			wantErr: true,
		},
		{
			name:    "a max frame size above two bytes is refused",
			options: map[string][]string{"max-frame-size": {"65536"}},
			wantErr: true,
		},
		{
			name: "every option at once",
			options: map[string][]string{
				"unit-identifier": {"1"},
				"request-timeout": {"250"},
				"max-frame-size":  {"1024"},
			},
			want: Configuration{unitIdentifier: 1, requestTimeout: 250 * time.Millisecond, maxFrameSize: 1024},
		},
		{
			name:    "an option with no value is ignored",
			options: map[string][]string{"unit-identifier": {}},
			want:    DefaultConfiguration(),
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			configuration, err := ParseFromOptions(zerolog.Nop(), testCase.options)
			if testCase.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, testCase.want, configuration)
		})
	}
}

func TestConfiguration_String(t *testing.T) {
	assert.Contains(t, DefaultConfiguration().String(), "unitIdentifier: 0")
}
