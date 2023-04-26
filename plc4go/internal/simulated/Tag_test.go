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

	"github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
)

func TestNewSimulatedField(t *testing.T) {
	type args struct {
		fieldType    TagType
		name         string
		dataTypeSize model.SimulatedDataTypeSizes
		quantity     uint16
	}
	tests := []struct {
		name string
		args args
		want simulatedTag
	}{
		{
			name: "simple",
			args: args{
				fieldType:    TagRandom,
				name:         "test",
				dataTypeSize: model.SimulatedDataTypeSizes_BOOL,
				quantity:     1,
			},
			want: simulatedTag{
				TagType:      TagRandom,
				Name:         "test",
				DataTypeSize: model.SimulatedDataTypeSizes_BOOL,
				Quantity:     1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSimulatedTag(tt.args.fieldType, tt.args.name, tt.args.dataTypeSize, tt.args.quantity); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewSimulatedTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimulatedField_GetAddressString(t1 *testing.T) {
	type fields struct {
		FieldType    TagType
		Name         string
		DataTypeSize model.SimulatedDataTypeSizes
		Quantity     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple",
			fields: fields{
				FieldType:    TagRandom,
				Name:         "test",
				DataTypeSize: model.SimulatedDataTypeSizes_BOOL,
				Quantity:     1,
			},
			want: "RANDOM/test:BOOL[1]",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := simulatedTag{
				TagType:      tt.fields.FieldType,
				Name:         tt.fields.Name,
				DataTypeSize: tt.fields.DataTypeSize,
				Quantity:     tt.fields.Quantity,
			}
			if got := t.GetAddressString(); got != tt.want {
				t1.Errorf("GetAddressString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimulatedField_GetDataTypeSize(t1 *testing.T) {
	type fields struct {
		FieldType    TagType
		Name         string
		DataTypeSize model.SimulatedDataTypeSizes
		Quantity     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   model.SimulatedDataTypeSizes
	}{
		{
			name: "simple",
			fields: fields{
				FieldType:    TagRandom,
				Name:         "test",
				DataTypeSize: model.SimulatedDataTypeSizes_BOOL,
				Quantity:     1,
			},
			want: model.SimulatedDataTypeSizes_BOOL,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := simulatedTag{
				TagType:      tt.fields.FieldType,
				Name:         tt.fields.Name,
				DataTypeSize: tt.fields.DataTypeSize,
				Quantity:     tt.fields.Quantity,
			}
			if got := t.GetDataTypeSize(); got != tt.want {
				t1.Errorf("GetDataTypeSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimulatedField_GetFieldType(t1 *testing.T) {
	type fields struct {
		FieldType    TagType
		Name         string
		DataTypeSize model.SimulatedDataTypeSizes
		Quantity     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   TagType
	}{
		{
			name: "simple",
			fields: fields{
				FieldType:    TagRandom,
				Name:         "test",
				DataTypeSize: model.SimulatedDataTypeSizes_BOOL,
				Quantity:     1,
			},
			want: TagRandom,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := simulatedTag{
				TagType:      tt.fields.FieldType,
				Name:         tt.fields.Name,
				DataTypeSize: tt.fields.DataTypeSize,
				Quantity:     tt.fields.Quantity,
			}
			if got := t.GetTagType(); got != tt.want {
				t1.Errorf("GetTagType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimulatedField_GetName(t1 *testing.T) {
	type fields struct {
		FieldType    TagType
		Name         string
		DataTypeSize model.SimulatedDataTypeSizes
		Quantity     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "simple",
			fields: fields{
				FieldType:    TagRandom,
				Name:         "test",
				DataTypeSize: model.SimulatedDataTypeSizes_BOOL,
				Quantity:     1,
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := simulatedTag{
				TagType:      tt.fields.FieldType,
				Name:         tt.fields.Name,
				DataTypeSize: tt.fields.DataTypeSize,
				Quantity:     tt.fields.Quantity,
			}
			if got := t.GetName(); got != tt.want {
				t1.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*func TestSimulatedTag_GetQuantity(t1 *testing.T) {
	type fields struct {
		TagType    TagType
		Name         string
		DataTypeSize model.SimulatedDataTypeSizes
		Quantity     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   uint16
	}{
		{
			name: "simple",
			fields: fields{
				TagType:    TagRandom,
				Name:         "test",
				DataTypeSize: model.SimulatedDataTypeSizes_BOOL,
				Quantity:     1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := simulatedTag{
				TagType:    tt.fields.TagType,
				Name:         tt.fields.Name,
				DataTypeSize: tt.fields.DataTypeSize,
				Quantity:     tt.fields.Quantity,
			}
			if got := t.GetQuantity(); got != tt.want {
				t1.Errorf("GetQuantity() = %v, want %v", got, tt.want)
			}
		})
	}
}*/
