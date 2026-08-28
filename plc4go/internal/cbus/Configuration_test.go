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

package cbus

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	spiOptions "github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func TestParseFromOptions(t *testing.T) {
	type args struct {
		options map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		want    Configuration
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "broken bool",
			args: args{
				options: map[string][]string{
					"Srchk": {"12331"},
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "broken byte",
			args: args{
				options: map[string][]string{
					"MonitoredApplication2": {"true"},
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "check no options",
			want: Configuration{
				Srchk:                 true,
				Exstat:                true,
				Pun:                   false,
				LocalSal:              true,
				Pcn:                   false,
				Idmon:                 true,
				Monitor:               true,
				Smart:                 true,
				XonXoff:               false,
				Connect:               true,
				MonitoredApplication1: 255,
				MonitoredApplication2: 255,
			},
			wantErr: assert.NoError,
		},
		{
			name: "check all options",
			args: args{
				options: map[string][]string{
					"Srchk":                 {"false"},
					"Exstat":                {"false"},
					"Pun":                   {"true"},
					"LocalSal":              {"false"},
					"Pcn":                   {"true"},
					"Idmon":                 {"false"},
					"Monitor":               {"false"},
					"Smart":                 {"false"},
					"XonXoff":               {"true"},
					"Connect":               {"false"},
					"MonitoredApplication1": {"123"},
					"MonitoredApplication2": {"123"},
				},
			},
			want: Configuration{
				Srchk:                 false,
				Exstat:                false,
				Pun:                   true,
				LocalSal:              false,
				Pcn:                   true,
				Idmon:                 false,
				Monitor:               false,
				Smart:                 false,
				XonXoff:               true,
				Connect:               false,
				MonitoredApplication1: 123,
				MonitoredApplication2: 123,
			},
			wantErr: assert.NoError,
		},
		{
			name: "check case sensitivity",
			args: args{
				options: map[string][]string{
					"srchk": {"false"},
				},
			},
			want: Configuration{
				Srchk:                 false,
				Exstat:                true,
				Pun:                   false,
				LocalSal:              true,
				Pcn:                   false,
				Idmon:                 true,
				Monitor:               true,
				Smart:                 true,
				XonXoff:               false,
				Connect:               true,
				MonitoredApplication1: 255,
				MonitoredApplication2: 255,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFromOptions(testutils.ProduceTestingLogger(t), tt.args.options)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseFromOptions(%v)", tt.args.options)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseFromOptions(%v)", tt.args.options)
		})
	}
}

// titleOptions title-cases every key in place, marker included, so the exact-match skip in
// ReportUnknown is not enough here: the title-cased duplicate must not surface as unknown.
func TestParseFromOptions_saysNothingAboutTheActiveTransportMarker(t *testing.T) {
	var logged bytes.Buffer
	log := zerolog.New(&logged)

	_, err := ParseFromOptions(log, map[string][]string{spiOptions.ActiveTransportOption: {"tcp"}})

	assert.NoError(t, err)
	assert.Empty(t, logged.String())
}

func Test_createDefaultConfiguration(t *testing.T) {
	tests := []struct {
		name string
		want Configuration
	}{
		{
			name: "default returns default",
			want: Configuration{
				Exstat:   true,
				LocalSal: true,
				Idmon:    true,
				Monitor:  true,
				Smart:    true,
				Srchk:    true,
				Connect:  true,

				MonitoredApplication1: 0xFF,
				MonitoredApplication2: 0xFF,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, createDefaultConfiguration(), "createDefaultConfiguration()")
		})
	}
}

// The per-driver getFromOptions this used to cover is gone: every driver now reads options
// through spi/options.OptionReader, and its behaviour is covered once, in OptionReader_test.go.
