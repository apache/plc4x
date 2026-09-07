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

package serial

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/spi/transports/serial/serialport"
)

func TestParseSerialOptions(t *testing.T) {
	tests := []struct {
		name    string
		options map[string][]string
		want    serialConfig
		wantErr string // substring; empty = expect success
	}{
		{
			name:    "empty map yields all defaults",
			options: map[string][]string{},
			want:    defaultSerialConfig(),
		},
		{
			name: "full happy path",
			options: map[string][]string{
				"serial.baud-rate":          {"19200"},
				"serial.data-bits":          {"7"},
				"serial.stop-bits":          {"2"},
				"serial.parity":             {"even"},
				"serial.flow-control":       {"rts-cts"},
				"serial.dtr":                {"true"},
				"serial.rts":                {"true"},
				"serial.read-timeout-ms":    {"250"},
				"serial.write-timeout-ms":   {"0"},
				"serial.connect-timeout-ms": {"5000"},
			},
			want: serialConfig{
				port: serialport.Config{
					BaudRate: 19200, DataBits: 7,
					StopBits: serialport.StopBitsTwo, Parity: serialport.ParityEven,
					RTSCTSFlowControl: true,
				},
				dtr: true, rts: true,
				readTimeout: 250 * time.Millisecond, writeTimeout: 0,
				connectTimeout: 5000,
			},
		},
		{
			name:    "enum values are case-insensitive and accept underscores",
			options: map[string][]string{"serial.parity": {"EVEN"}, "serial.flow-control": {"XON_XOFF"}},
			want: func() serialConfig {
				c := defaultSerialConfig()
				c.port.Parity = serialport.ParityEven
				c.port.XONXOFFFlowControl = true
				return c
			}(),
		},
		{
			name:    "mark and space parity",
			options: map[string][]string{"serial.parity": {"Mark"}},
			want: func() serialConfig {
				c := defaultSerialConfig()
				c.port.Parity = serialport.ParityMark
				return c
			}(),
		},
		{
			name:    "read-timeout-ms zero means blocking",
			options: map[string][]string{"serial.read-timeout-ms": {"0"}},
			want: func() serialConfig {
				c := defaultSerialConfig()
				c.readTimeout = 0
				return c
			}(),
		},
		{
			name:    "invalid baud rate",
			options: map[string][]string{"serial.baud-rate": {"fast"}},
			wantErr: `"serial.baud-rate"`,
		},
		{
			name:    "zero baud rate rejected",
			options: map[string][]string{"serial.baud-rate": {"0"}},
			wantErr: `"serial.baud-rate"`,
		},
		{
			name:    "data-bits out of range",
			options: map[string][]string{"serial.data-bits": {"9"}},
			wantErr: `"serial.data-bits"`,
		},
		{
			name:    "stop-bits out of range",
			options: map[string][]string{"serial.stop-bits": {"3"}},
			wantErr: `"serial.stop-bits"`,
		},
		{
			name:    "unknown parity value",
			options: map[string][]string{"serial.parity": {"strong"}},
			wantErr: `"serial.parity"`,
		},
		{
			name:    "combined flow control value rejected",
			options: map[string][]string{"serial.flow-control": {"rts-cts-xon-xoff"}},
			wantErr: `"serial.flow-control"`,
		},
		{
			name:    "invalid dtr boolean",
			options: map[string][]string{"serial.dtr": {"yes-please"}},
			wantErr: `"serial.dtr"`,
		},
		{
			name:    "invalid read-timeout-ms",
			options: map[string][]string{"serial.read-timeout-ms": {"-5"}},
			wantErr: `"serial.read-timeout-ms"`,
		},
		{
			name:    "empty value slice ignored like absent option",
			options: map[string][]string{"serial.parity": {}},
			want:    defaultSerialConfig(),
		},
		{
			name:    "unknown options are ignored",
			options: map[string][]string{"break-enabled": {"true"}, "no-such-thing": {"1"}},
			want:    defaultSerialConfig(),
		},
		{
			name:    "reuse-port accepted",
			options: map[string][]string{"serial.reuse-port": {"true"}},
			want: func() serialConfig {
				c := defaultSerialConfig()
				c.reusePort = true
				return c
			}(),
		},
		{
			name:    "interframe-delay accepted",
			options: map[string][]string{"serial.interframe-delay": {"50"}},
			want: func() serialConfig {
				c := defaultSerialConfig()
				c.interframeDelay = 50 * time.Millisecond
				return c
			}(),
		},
		{
			name:    "invalid reuse-port",
			options: map[string][]string{"serial.reuse-port": {"maybe"}},
			wantErr: `"serial.reuse-port"`,
		},
		{
			name:    "invalid interframe-delay",
			options: map[string][]string{"serial.interframe-delay": {"-1"}},
			wantErr: `"serial.interframe-delay"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSerialOptions(tt.options)
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
