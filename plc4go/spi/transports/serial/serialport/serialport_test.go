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

package serialport

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigNormalize(t *testing.T) {
	tests := []struct {
		name    string
		in      Config
		want    Config
		wantErr string
	}{
		{
			name: "defaults applied: zero data bits become 8",
			in:   Config{BaudRate: 9600},
			want: Config{BaudRate: 9600, DataBits: 8, StopBits: StopBitsOne, Parity: ParityNone},
		},
		{
			name: "explicit full config kept",
			in:   Config{BaudRate: 115200, DataBits: 7, StopBits: StopBitsTwo, Parity: ParityEven, RTSCTSFlowControl: true},
			want: Config{BaudRate: 115200, DataBits: 7, StopBits: StopBitsTwo, Parity: ParityEven, RTSCTSFlowControl: true},
		},
		{
			name:    "zero baud rate rejected",
			in:      Config{},
			wantErr: "baud rate",
		},
		{
			name:    "data bits below 5 rejected",
			in:      Config{BaudRate: 9600, DataBits: 4},
			wantErr: "data bits",
		},
		{
			name:    "data bits above 8 rejected",
			in:      Config{BaudRate: 9600, DataBits: 9},
			wantErr: "data bits",
		},
		{
			name:    "invalid stop bits enum rejected",
			in:      Config{BaudRate: 9600, StopBits: StopBits(99)},
			wantErr: "stop bits",
		},
		{
			name:    "invalid parity enum rejected",
			in:      Config{BaudRate: 9600, Parity: Parity(99)},
			wantErr: "parity",
		},
		{
			// Windows rule (only platform supporting 1.5): 1.5 stop bits require 5 data bits.
			name:    "1.5 stop bits with 8 data bits rejected",
			in:      Config{BaudRate: 9600, StopBits: StopBitsOnePointFive},
			wantErr: "1.5 stop bits",
		},
		{
			name: "1.5 stop bits with 5 data bits accepted",
			in:   Config{BaudRate: 9600, DataBits: 5, StopBits: StopBitsOnePointFive},
			want: Config{BaudRate: 9600, DataBits: 5, StopBits: StopBitsOnePointFive, Parity: ParityNone},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.in.normalize()
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOpenRejectsEmptyPortName(t *testing.T) {
	_, err := Open("", Config{BaudRate: 9600})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "port name")
}

func TestOpenRejectsInvalidConfigBeforeTouchingPort(t *testing.T) {
	_, err := Open("/dev/does-not-matter", Config{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "baud rate")
}
