/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/apache/plc4x/plc4go/internal/plc4go/simulated/readwrite/model"
	values2 "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"github.com/rs/zerolog/log"
	"reflect"
	"testing"
)

func TestDevice_Get(t1 *testing.T) {
	type fields struct {
		Name  string
		State map[SimulatedField]*values.PlcValue
	}
	type args struct {
		field        SimulatedField
		verifyOutput bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *values.PlcValue
	}{
		{
			name: "simple state",
			fields: fields{
				Name: "hurz",
				State: map[SimulatedField]*values.PlcValue{
					NewSimulatedField(FieldState, "boolField", model.SimulatedDataTypeSizes_BOOL, 1): ToReference(values2.NewPlcBOOL(true)),
				},
			},
			args: args{
				field:        NewSimulatedField(FieldState, "boolField", model.SimulatedDataTypeSizes_BOOL, 1),
				verifyOutput: true,
			},
			want: ToReference(values2.NewPlcBOOL(true)),
		},
		{
			name: "simple random",
			fields: fields{
				Name:  "hurz",
				State: map[SimulatedField]*values.PlcValue{},
			},
			args: args{
				field:        NewSimulatedField(FieldRandom, "boolField", model.SimulatedDataTypeSizes_BOOL, 1),
				verifyOutput: false,
			},
			want: ToReference(values2.NewPlcBOOL(true)),
		},
		{
			name: "simple stdout",
			fields: fields{
				Name:  "hurz",
				State: map[SimulatedField]*values.PlcValue{},
			},
			args: args{
				field:        NewSimulatedField(FieldStdOut, "boolField", model.SimulatedDataTypeSizes_BOOL, 1),
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
			if tt.args.verifyOutput && !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
 * When first playing around with random values I only got "false" values.
 * So I added this test in order to verify I'm actually getting random values.
 */
func TestDevice_Random(t1 *testing.T) {
	type fields struct {
		Name  string
		State map[SimulatedField]*values.PlcValue
	}
	type args struct {
		field   SimulatedField
		numRuns int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *values.PlcValue
	}{
		{
			name: "simple random",
			fields: fields{
				Name:  "hurz",
				State: map[SimulatedField]*values.PlcValue{},
			},
			args: args{
				field:   NewSimulatedField(FieldRandom, "boolField", model.SimulatedDataTypeSizes_BOOL, 1),
				numRuns: 1000,
			},
			want: ToReference(values2.NewPlcBOOL(true)),
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
				t1.Errorf("Random doesn't seem to work. In %d runs I got %d true and %d false values", tt.args.numRuns, numTrue, numFalse)
			} else {
				log.Info().Msgf("In %d runs I got %d true and %d false values", tt.args.numRuns, numTrue, numFalse)
			}
		})
	}
}

func TestDevice_Set(t1 *testing.T) {
	type fields struct {
		Name  string
		State map[SimulatedField]*values.PlcValue
	}
	type args struct {
		field         SimulatedField
		value         *values.PlcValue
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
				State: map[SimulatedField]*values.PlcValue{},
			},
			args: args{
				field:         NewSimulatedField(FieldState, "boolField", model.SimulatedDataTypeSizes_BOOL, 1),
				value:         ToReference(values2.NewPlcBOOL(true)),
				shouldBeSaved: true,
			},
		},
		{
			name: "simple random",
			fields: fields{
				Name:  "hurz",
				State: map[SimulatedField]*values.PlcValue{},
			},
			args: args{
				field:         NewSimulatedField(FieldRandom, "boolField", model.SimulatedDataTypeSizes_BOOL, 1),
				value:         ToReference(values2.NewPlcBOOL(true)),
				shouldBeSaved: false,
			},
		},
		{
			name: "simple stdout",
			fields: fields{
				Name:  "hurz",
				State: map[SimulatedField]*values.PlcValue{},
			},
			args: args{
				field:         NewSimulatedField(FieldStdOut, "boolField", model.SimulatedDataTypeSizes_BOOL, 1),
				value:         ToReference(values2.NewPlcBOOL(true)),
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
		State map[SimulatedField]*values.PlcValue
	}
	type args struct {
		field SimulatedField
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *values.PlcValue
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Device{
				Name:  tt.fields.Name,
				State: tt.fields.State,
			}
			if got := t.getRandomValue(tt.args.field); !reflect.DeepEqual(got, tt.want) {
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
			if got := NewDevice(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ToReference(value values.PlcValue) *values.PlcValue {
	return &value
}
