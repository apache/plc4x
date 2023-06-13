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
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
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
				Srchk:                 true, // Note: at the moment the fields are case-sensitive so the key above is just ignored
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

func Test_getFromOptions(t *testing.T) {
	type args struct {
		options map[string][]string
		key     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "key not found",
			args: args{
				options: map[string][]string{},
				key:     "testKey",
			},
		},
		{
			name: "key found",
			args: args{
				options: map[string][]string{"testKey": {"asd", "asd"}},
				key:     "testKey",
			},
			want: "asd", // note: multi keys not supported yet, so first one is returned
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getFromOptions(testutils.ProduceTestingLogger(t), tt.args.options, tt.args.key), "getFromOptions(%v, %v)", tt.args.options, tt.args.key)
		})
	}
}

func TestConfiguration_String(t *testing.T) {
	type fields struct {
		Srchk                 bool
		Exstat                bool
		Pun                   bool
		LocalSal              bool
		Pcn                   bool
		Idmon                 bool
		Monitor               bool
		Smart                 bool
		XonXoff               bool
		Connect               bool
		MonitoredApplication1 byte
		MonitoredApplication2 byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "string it",
			fields: fields{
				Srchk:                 true,
				Exstat:                true,
				Pun:                   true,
				LocalSal:              true,
				Pcn:                   true,
				Idmon:                 true,
				Monitor:               true,
				Smart:                 true,
				XonXoff:               true,
				Connect:               true,
				MonitoredApplication1: 2,
				MonitoredApplication2: 3,
			},
			want: "cbus.Configuration{Srchk:true, Exstat:true, Pun:true, LocalSal:true, Pcn:true, Idmon:true, Monitor:true, Smart:true, XonXoff:true, Connect:true, MonitoredApplication1:0x2, MonitoredApplication2:0x3}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Configuration{
				Srchk:                 tt.fields.Srchk,
				Exstat:                tt.fields.Exstat,
				Pun:                   tt.fields.Pun,
				LocalSal:              tt.fields.LocalSal,
				Pcn:                   tt.fields.Pcn,
				Idmon:                 tt.fields.Idmon,
				Monitor:               tt.fields.Monitor,
				Smart:                 tt.fields.Smart,
				XonXoff:               tt.fields.XonXoff,
				Connect:               tt.fields.Connect,
				MonitoredApplication1: tt.fields.MonitoredApplication1,
				MonitoredApplication2: tt.fields.MonitoredApplication2,
			}
			assert.Equalf(t, tt.want, c.String(), "String()")
		})
	}
}
