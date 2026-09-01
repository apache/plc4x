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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
)

func TestFieldHandler_ParseQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    apiModel.PlcTag
		wantErr bool
	}{
		{
			name: "simple random",
			args: args{
				query: "RANDOM/test_random:BOOL",
			},
			want:    NewSimulatedTag(TagRandom, "test_random", readWriteModel.SimulatedDataTypeSizes_BOOL, 1),
			wantErr: false,
		},
		{
			name: "simple random array",
			args: args{
				query: "RANDOM/test_random[0..9]:BOOL",
			},
			want:    NewSimulatedTag(TagRandom, "test_random", readWriteModel.SimulatedDataTypeSizes_BOOL, 10),
			wantErr: false,
		},
		{
			name: "simple state",
			args: args{
				query: "STATE/test_state:BOOL",
			},
			want:    NewSimulatedTag(TagState, "test_state", readWriteModel.SimulatedDataTypeSizes_BOOL, 1),
			wantErr: false,
		},
		{
			name: "simple state array",
			args: args{
				query: "STATE/test_state[0..41]:BOOL",
			},
			want:    NewSimulatedTag(TagState, "test_state", readWriteModel.SimulatedDataTypeSizes_BOOL, 42),
			wantErr: false,
		},
		{
			name: "simple stdout",
			args: args{
				query: "STDOUT/test_stdout:BOOL",
			},
			want:    NewSimulatedTag(TagStdOut, "test_stdout", readWriteModel.SimulatedDataTypeSizes_BOOL, 1),
			wantErr: false,
		},
		{
			name: "simple stdout array",
			args: args{
				query: "STDOUT/test_stdout[0..22]:BOOL",
			},
			want:    NewSimulatedTag(TagStdOut, "test_stdout", readWriteModel.SimulatedDataTypeSizes_BOOL, 23),
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
				query: "RANDOM/test_stdout[0..22]:HURZ",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error invalid array size",
			args: args{
				query: "RANDOM/test_stdout[0..999999999999999999999999999999999998]:BOOL",
			},
			want:    nil,
			wantErr: true,
		},
		{
			// The count is carried as a uint16, so anything above it used to come back as the low two
			// bytes of what was asked for - a tag of 4464 elements rather than an error.
			name: "error array size past the count a tag can carry",
			args: args{
				query: "RANDOM/test_stdout[0..69999]:BOOL",
			},
			want:    nil,
			wantErr: true,
		},
		{
			// The same, one past the boundary, where the low two bytes are zero and the tag came back
			// holding nothing at all.
			name: "error array size one past the count a tag can carry",
			args: args{
				query: "RANDOM/test_stdout[0..65535]:BOOL",
			},
			want:    nil,
			wantErr: true,
		},
		{
			// A count of zero has no spelling in the notation - a range is written with the
			// indices it covers - but an inverted one is still nonsense.
			name: "error inverted range",
			args: args{
				query: "RANDOM/test_stdout[3..1]:BOOL",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "largest array size a tag can carry",
			args: args{
				query: "RANDOM/test_stdout[0..65534]:BOOL",
			},
			want:    NewSimulatedTag(TagRandom, "test_stdout", readWriteModel.SimulatedDataTypeSizes_BOOL, 65535),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTagHandler()
			got, err := m.ParseTag(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("ParseQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldType_Name(t *testing.T) {
	tests := []struct {
		name string
		e    TagType
		want string
	}{
		{
			name: "simple random",
			e:    TagRandom,
			want: "RANDOM",
		},
		{
			name: "simple state",
			e:    TagState,
			want: "STATE",
		},
		{
			name: "simple stdout",
			e:    TagStdOut,
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

// This driver addresses a named variable rather than a numeric offset, so a selection that does
// not start at the first element has nothing to apply to. Reading it as though it did would hand
// back the first elements under the impression they were the requested ones.
func TestFieldHandler_ParseTagRejectsASelectionWithAnOffset(t *testing.T) {
	handler := NewTagHandler()

	_, err := handler.ParseTag("RANDOM/test[4..7]:INT")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must start at the first element")

	// A declared base is what the offset is measured from, so [4..7;4] does start at the first
	// element and is accepted.
	fromDeclaredBase, err := handler.ParseTag("RANDOM/test[4..7;4]:INT")
	require.NoError(t, err)
	assert.Equal(t, uint16(4), fromDeclaredBase.(simulatedTag).Quantity)
}
