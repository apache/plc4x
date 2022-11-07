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

package simulated

import (
	"reflect"
	"testing"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	model2 "github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
)

func TestFieldHandler_ParseQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    model.PlcField
		wantErr bool
	}{
		{
			name: "simple random",
			args: args{
				query: "RANDOM/test_random:BOOL",
			},
			want:    NewSimulatedField(FieldRandom, "test_random", model2.SimulatedDataTypeSizes_BOOL, 1),
			wantErr: false,
		},
		{
			name: "simple random array",
			args: args{
				query: "RANDOM/test_random:BOOL[10]",
			},
			want:    NewSimulatedField(FieldRandom, "test_random", model2.SimulatedDataTypeSizes_BOOL, 10),
			wantErr: false,
		},
		{
			name: "simple state",
			args: args{
				query: "STATE/test_state:BOOL",
			},
			want:    NewSimulatedField(FieldState, "test_state", model2.SimulatedDataTypeSizes_BOOL, 1),
			wantErr: false,
		},
		{
			name: "simple state array",
			args: args{
				query: "STATE/test_state:BOOL[42]",
			},
			want:    NewSimulatedField(FieldState, "test_state", model2.SimulatedDataTypeSizes_BOOL, 42),
			wantErr: false,
		},
		{
			name: "simple stdout",
			args: args{
				query: "STDOUT/test_stdout:BOOL",
			},
			want:    NewSimulatedField(FieldStdOut, "test_stdout", model2.SimulatedDataTypeSizes_BOOL, 1),
			wantErr: false,
		},
		{
			name: "simple stdout array",
			args: args{
				query: "STDOUT/test_stdout:BOOL[23]",
			},
			want:    NewSimulatedField(FieldStdOut, "test_stdout", model2.SimulatedDataTypeSizes_BOOL, 23),
			wantErr: false,
		},
		{
			name: "error invalid type",
			args: args{
				query: "HURZ/test_stdout:BOOL[23]",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error invalid name format",
			args: args{
				query: "RANDOM/test/stdout:BOOL[23]",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error invalid datatype",
			args: args{
				query: "RANDOM/test_stdout:HURZ[23]",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error invalid array size",
			args: args{
				query: "RANDOM/test_stdout:BOOL[999999999999999999999999999999999999]",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewFieldHandler()
			got, err := m.ParseField(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldType_Name(t *testing.T) {
	tests := []struct {
		name string
		e    FieldType
		want string
	}{
		{
			name: "simple random",
			e:    FieldRandom,
			want: "RANDOM",
		},
		{
			name: "simple state",
			e:    FieldState,
			want: "STATE",
		},
		{
			name: "simple stdout",
			e:    FieldStdOut,
			want: "STDOUT",
		},
		{
			name: "simple stdout",
			e:    10,
			want: "UNKNOWN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}
