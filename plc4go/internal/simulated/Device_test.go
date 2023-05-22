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
	"github.com/stretchr/testify/assert"
	"testing"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"

	"github.com/rs/zerolog/log"
)

func TestDevice_Get(t1 *testing.T) {
	type fields struct {
		Name  string
		State map[simulatedTag]*apiValues.PlcValue
	}
	type args struct {
		field        simulatedTag
		verifyOutput bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *apiValues.PlcValue
	}{
		{
			name: "simple state",
			fields: fields{
				Name: "hurz",
				State: map[simulatedTag]*apiValues.PlcValue{
					NewSimulatedTag(TagState, "boolTag", readWriteModel.SimulatedDataTypeSizes_BOOL, 1): ToReference(spiValues.NewPlcBOOL(true)),
				},
			},
			args: args{
				field:        NewSimulatedTag(TagState, "boolTag", readWriteModel.SimulatedDataTypeSizes_BOOL, 1),
				verifyOutput: true,
			},
			want: ToReference(spiValues.NewPlcBOOL(true)),
		},
		{
			name: "simple random",
			fields: fields{
				Name:  "hurz",
				State: map[simulatedTag]*apiValues.PlcValue{},
			},
			args: args{
				field:        NewSimulatedTag(TagRandom, "boolTag", readWriteModel.SimulatedDataTypeSizes_BOOL, 1),
				verifyOutput: false,
			},
			want: ToReference(spiValues.NewPlcBOOL(true)),
		},
		{
			name: "simple stdout",
			fields: fields{
				Name:  "hurz",
				State: map[simulatedTag]*apiValues.PlcValue{},
			},
			args: args{
				field:        NewSimulatedTag(TagStdOut, "boolTag", readWriteModel.SimulatedDataTypeSizes_BOOL, 1),
				verifyOutput: false,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Device{
				Name:  tt.fields.Name,
				State: tt.fields.State,
			}
			got := t.Get(tt.args.field)
			if got != nil {
				log.Debug().Msgf("Result: %v", *got)
			} else {
				log.Debug().Msg("Result: nil")
			}
			if tt.args.verifyOutput && !assert.Equal(t1, tt.want, got) {
				t1.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
 * When first playing around with random apiValues I only got "false" apiValues.
 * So I added this test in order to verify I'm actually getting random apiValues.
 */
func TestDevice_Random(t1 *testing.T) {
	type fields struct {
		Name  string
		State map[simulatedTag]*apiValues.PlcValue
	}
	type args struct {
		field   simulatedTag
		numRuns int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *apiValues.PlcValue
	}{
		{
			name: "simple random",
			fields: fields{
				Name:  "hurz",
				State: map[simulatedTag]*apiValues.PlcValue{},
			},
			args: args{
				field:   NewSimulatedTag(TagRandom, "boolTag", readWriteModel.SimulatedDataTypeSizes_BOOL, 1),
				numRuns: 1000,
			},
			want: ToReference(spiValues.NewPlcBOOL(true)),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Device{
				Name:  tt.fields.Name,
				State: tt.fields.State,
			}
			numTrue := 0
			numFalse := 0
			for i := 0; i < tt.args.numRuns; i++ {
				got := t.Get(tt.args.field)
				boolValue := (*got).GetBool()
				if boolValue {
					numTrue++
				} else {
					numFalse++
				}
			}
			if numTrue == 0 || numFalse == 0 {
				t1.Errorf("Random doesn't seem to work. In %d runs I got %d true and %d false apiValues", tt.args.numRuns, numTrue, numFalse)
			} else {
				log.Info().Msgf("In %d runs I got %d true and %d false apiValues", tt.args.numRuns, numTrue, numFalse)
			}
		})
	}
}

func TestDevice_Set(t1 *testing.T) {
	type fields struct {
		Name  string
		State map[simulatedTag]*apiValues.PlcValue
	}
	type args struct {
		field         simulatedTag
		value         *apiValues.PlcValue
		shouldBeSaved bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "simple state",
			fields: fields{
				Name:  "hurz",
				State: map[simulatedTag]*apiValues.PlcValue{},
			},
			args: args{
				field:         NewSimulatedTag(TagState, "boolTag", readWriteModel.SimulatedDataTypeSizes_BOOL, 1),
				value:         ToReference(spiValues.NewPlcBOOL(true)),
				shouldBeSaved: true,
			},
		},
		{
			name: "simple random",
			fields: fields{
				Name:  "hurz",
				State: map[simulatedTag]*apiValues.PlcValue{},
			},
			args: args{
				field:         NewSimulatedTag(TagRandom, "boolTag", readWriteModel.SimulatedDataTypeSizes_BOOL, 1),
				value:         ToReference(spiValues.NewPlcBOOL(true)),
				shouldBeSaved: false,
			},
		},
		{
			name: "simple stdout",
			fields: fields{
				Name:  "hurz",
				State: map[simulatedTag]*apiValues.PlcValue{},
			},
			args: args{
				field:         NewSimulatedTag(TagStdOut, "boolTag", readWriteModel.SimulatedDataTypeSizes_BOOL, 1),
				value:         ToReference(spiValues.NewPlcBOOL(true)),
				shouldBeSaved: false,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Device{
				Name:  tt.fields.Name,
				State: tt.fields.State,
			}
			// It shouldn't exist in the map before
			if _, ok := tt.fields.State[tt.args.field]; ok {
				t1.Errorf("Value for %v already present in map", tt.args.field)
			}
			t.Set(tt.args.field, tt.args.value)
			// It should exist in the map after
			if _, ok := tt.fields.State[tt.args.field]; tt.args.shouldBeSaved != ok {
				if tt.args.shouldBeSaved {
					t1.Errorf("Value for %v not present in map (it should)", tt.args.field)
				} else {
					t1.Errorf("Value for %v present in map (is should not)", tt.args.field)
				}
			}
		})
	}
}

func TestDevice_getRandomValue(t1 *testing.T) {
	type fields struct {
		Name  string
		State map[simulatedTag]*apiValues.PlcValue
	}
	type args struct {
		field simulatedTag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *apiValues.PlcValue
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Device{
				Name:  tt.fields.Name,
				State: tt.fields.State,
			}
			if got := t.getRandomValue(tt.args.field); !assert.Equal(t1, tt.want, got) {
				t1.Errorf("getRandomValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDevice(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *Device
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDevice(tt.args.name); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ToReference(value apiValues.PlcValue) *apiValues.PlcValue {
	return &value
}
