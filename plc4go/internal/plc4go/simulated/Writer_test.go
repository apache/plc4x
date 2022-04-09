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
	"github.com/apache/plc4x/plc4go/internal/plc4go/s7"
	model4 "github.com/apache/plc4x/plc4go/internal/plc4go/s7/readwrite/model"
	model2 "github.com/apache/plc4x/plc4go/internal/plc4go/simulated/readwrite/model"
	model3 "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	values2 "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"reflect"
	"testing"
	"time"
)

func TestWriter_Write(t *testing.T) {
	type fields struct {
		device  *Device
		options map[string][]string
	}
	type args struct {
		fields     map[string]model.PlcField
		values     map[string]values.PlcValue
		fieldNames []string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         model.PlcWriteResponse
		newState     map[SimulatedField]*values.PlcValue
		delayAtLeast time.Duration
	}{
		{
			name: "simple state",
			fields: fields{
				device: &Device{
					Name:  "hurz",
					State: map[SimulatedField]*values.PlcValue{},
				},
				options: map[string][]string{},
			},
			args: args{
				fields: map[string]model.PlcField{
					"test": NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1),
				},
				values: map[string]values.PlcValue{
					"test": values2.NewPlcBOOL(true),
				},
				fieldNames: []string{"test"},
			},
			want: model3.NewDefaultPlcWriteResponse(nil,
				map[string]model.PlcResponseCode{
					"test": model.PlcResponseCode_OK,
				}),
			newState: map[SimulatedField]*values.PlcValue{
				NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1): ToReference(values2.NewPlcBOOL(true)),
			},
			delayAtLeast: 0,
		},
		{
			name: "simple state overwrite",
			fields: fields{
				device: &Device{
					Name: "hurz",
					State: map[SimulatedField]*values.PlcValue{
						NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1): ToReference(values2.NewPlcBOOL(true)),
					},
				},
				options: map[string][]string{},
			},
			args: args{
				fields: map[string]model.PlcField{
					"test": NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1),
				},
				values: map[string]values.PlcValue{
					"test": values2.NewPlcBOOL(false),
				},
				fieldNames: []string{"test"},
			},
			want: model3.NewDefaultPlcWriteResponse(nil,
				map[string]model.PlcResponseCode{
					"test": model.PlcResponseCode_OK,
				}),
			newState: map[SimulatedField]*values.PlcValue{
				NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1): ToReference(values2.NewPlcBOOL(false)),
			},
			delayAtLeast: 0,
		},
		{
			name: "simple state delayed",
			fields: fields{
				device: &Device{
					Name: "hurz",
					State: map[SimulatedField]*values.PlcValue{
						NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1): ToReference(values2.NewPlcBOOL(true)),
					},
				},
				options: map[string][]string{
					"writeDelay": {"1000"},
				},
			},
			args: args{
				fields: map[string]model.PlcField{
					"test": NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1),
				},
				values: map[string]values.PlcValue{
					"test": values2.NewPlcBOOL(false),
				},
				fieldNames: []string{"test"},
			},
			want: model3.NewDefaultPlcWriteResponse(nil,
				map[string]model.PlcResponseCode{
					"test": model.PlcResponseCode_OK,
				}),
			newState: map[SimulatedField]*values.PlcValue{
				NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1): ToReference(values2.NewPlcBOOL(false)),
			},
			delayAtLeast: 1000,
		},
		// Passing in a completely wrong type of field.
		{
			name: "invalid field type",
			fields: fields{
				device: &Device{
					Name: "hurz",
					State: map[SimulatedField]*values.PlcValue{
						NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1): ToReference(values2.NewPlcBOOL(true)),
					},
				},
				options: map[string][]string{},
			},
			args: args{
				fields: map[string]model.PlcField{
					"test": s7.PlcField{
						FieldType:   s7.S7Field,
						MemoryArea:  model4.MemoryArea_DATA_BLOCKS,
						BlockNumber: 1,
						ByteOffset:  1,
						BitOffset:   0,
						NumElements: 1,
						Datatype:    model4.TransportSize_BOOL,
					},
				},
				values: map[string]values.PlcValue{
					"test": values2.NewPlcBOOL(false),
				},
				fieldNames: []string{"test"},
			},
			want: model3.NewDefaultPlcWriteResponse(nil,
				map[string]model.PlcResponseCode{
					"test": model.PlcResponseCode_INVALID_ADDRESS,
				}),
			newState: map[SimulatedField]*values.PlcValue{
				NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1): ToReference(values2.NewPlcBOOL(true)),
			},
			delayAtLeast: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWriter(tt.fields.device, tt.fields.options, nil)
			writeRequest := model3.NewDefaultPlcWriteRequest(tt.args.fields, tt.args.fieldNames, tt.args.values, w, nil)
			timeBeforeWriteRequest := time.Now()
			writeResponseChannel := w.Write(writeRequest)
			select {
			case writeResponse := <-writeResponseChannel:
				timeAfterWriteRequest := time.Now()
				// If an expected delay was defined, check if closing
				// took at least this long.
				if tt.delayAtLeast > 0 {
					pingTime := timeAfterWriteRequest.Sub(timeBeforeWriteRequest)
					if pingTime < tt.delayAtLeast {
						t.Errorf("Writer.Write() completed too fast. Expected at least %v but returned after %v", tt.delayAtLeast, pingTime)
					}
				}
				if !reflect.DeepEqual(writeResponse.GetRequest(), writeRequest) {
					t.Errorf("Writer.Write() ReadRequest = %v, want %v", writeResponse.GetRequest(), writeRequest)
				}
				for _, fieldName := range writeRequest.GetFieldNames() {
					if !reflect.DeepEqual(writeResponse.GetResponse().GetResponseCode(fieldName), tt.want.GetResponseCode(fieldName)) {
						t.Errorf("Writer.Write() PlcResponse.ResponseCode = %v, want %v",
							writeResponse.GetResponse().GetResponseCode(fieldName), tt.want.GetResponseCode(fieldName))
					}
				}
				if !reflect.DeepEqual(tt.fields.device.State, tt.newState) {
					t.Errorf("Writer.Write() Device State = %v, want %v",
						tt.fields.device.State, tt.newState)
				}
			case <-time.After(3 * time.Second):
				t.Errorf("Reader.Read() got timeout")
			}
		})
	}
}
