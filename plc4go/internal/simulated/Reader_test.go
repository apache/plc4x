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
	"context"
	"github.com/apache/plc4x/plc4go/internal/s7"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	model4 "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	model2 "github.com/apache/plc4x/plc4go/protocols/simulated/readwrite/model"
	model3 "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	values2 "github.com/apache/plc4x/plc4go/spi/values"
	"reflect"
	"testing"
	"time"
)

func TestReader_Read(t *testing.T) {
	type fields struct {
		device  *Device
		options map[string][]string
	}
	type args struct {
		fields     map[string]model.PlcField
		fieldNames []string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         model.PlcReadResponse
		delayAtLeast time.Duration
	}{
		{
			name: "simple state",
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
				fieldNames: []string{"test"},
			},
			want: model3.NewDefaultPlcReadResponse(nil,
				map[string]model.PlcResponseCode{
					"test": model.PlcResponseCode_OK,
				},
				map[string]values.PlcValue{
					"test": values2.NewPlcBOOL(true),
				}),
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
					"readDelay": {"1000"},
				},
			},
			args: args{
				fields: map[string]model.PlcField{
					"test": NewSimulatedField(FieldState, "test", model2.SimulatedDataTypeSizes_BOOL, 1),
				},
				fieldNames: []string{"test"},
			},
			want: model3.NewDefaultPlcReadResponse(nil,
				map[string]model.PlcResponseCode{
					"test": model.PlcResponseCode_OK,
				},
				map[string]values.PlcValue{
					"test": values2.NewPlcBOOL(true),
				}),
			delayAtLeast: 1000,
		},
		{
			name: "state not found",
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
					"test": NewSimulatedField(FieldState, "lalala", model2.SimulatedDataTypeSizes_BOOL, 1),
				},
				fieldNames: []string{"test"},
			},
			want: model3.NewDefaultPlcReadResponse(nil,
				map[string]model.PlcResponseCode{
					"test": model.PlcResponseCode_NOT_FOUND,
				},
				map[string]values.PlcValue{
					"test": nil,
				}),
			delayAtLeast: 0,
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
					"test": s7.NewField(model4.MemoryArea_DATA_BLOCKS, 1, 1, 0, 1, model4.TransportSize_BOOL),
				},
				fieldNames: []string{"test"},
			},
			want: model3.NewDefaultPlcReadResponse(nil,
				map[string]model.PlcResponseCode{
					"test": model.PlcResponseCode_INVALID_ADDRESS,
				},
				map[string]values.PlcValue{
					"test": nil,
				}),
			delayAtLeast: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReader(tt.fields.device, tt.fields.options, nil)
			readRequest := model3.NewDefaultPlcReadRequest(tt.args.fields, tt.args.fieldNames, r, nil)
			timeBeforeReadRequest := time.Now()
			readResponseChannel := r.Read(context.TODO(), readRequest)
			timeout := time.NewTimer(3 * time.Second)
			defer utils.CleanupTimer(timeout)
			select {
			case readResponse := <-readResponseChannel:
				timeAfterReadRequest := time.Now()
				// If an expected delay was defined, check if closing
				// took at least this long.
				if tt.delayAtLeast > 0 {
					pingTime := timeAfterReadRequest.Sub(timeBeforeReadRequest)
					if pingTime < tt.delayAtLeast {
						t.Errorf("Reader.Read() completed too fast. Expected at least %v but returned after %v", tt.delayAtLeast, pingTime)
					}
				}
				if !reflect.DeepEqual(readResponse.GetRequest(), readRequest) {
					t.Errorf("Reader.Read() ReadRequest = %v, want %v", readResponse.GetRequest(), readRequest)
				}
				for _, fieldName := range readRequest.GetFieldNames() {
					if !reflect.DeepEqual(readResponse.GetResponse().GetResponseCode(fieldName), tt.want.GetResponseCode(fieldName)) {
						t.Errorf("Reader.Read() PlcResponse.ResponseCode = %v, want %v",
							readResponse.GetResponse().GetResponseCode(fieldName), tt.want.GetResponseCode(fieldName))
					}
					if !reflect.DeepEqual(readResponse.GetResponse().GetValue(fieldName), tt.want.GetValue(fieldName)) {
						t.Errorf("Reader.Read() PlcResponse.Value = %v, want %v",
							readResponse.GetResponse().GetValue(fieldName), tt.want.GetValue(fieldName))
					}
				}
			case <-timeout.C:
				t.Errorf("Reader.Read() got timeout")
			}
		})
	}
}
