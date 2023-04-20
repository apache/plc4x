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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readwriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/stretchr/testify/assert"
	"testing"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func TestNewValueHandler(t *testing.T) {
	tests := []struct {
		name string
		want ValueHandler
	}{
		{
			name: "create a new one",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewValueHandler(), "NewValueHandler()")
		})
	}
}

func TestValueHandler_NewPlcValue(t *testing.T) {
	type fields struct {
		DefaultValueHandler spiValues.DefaultValueHandler
	}
	type args struct {
		tag   apiModel.PlcTag
		value any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    apiValues.PlcValue
		wantErr assert.ErrorAssertionFunc
	}{
		/*  TODO: add tests for that
		CAL_WRITE,
		CAL_IDENTIFY_REPLY,
		CAL_STATUS,
		CAL_STATUS_EXTENDED
		*/
		/*
			TODO: add tests for that
			SAL
			ApplicationId_FREE_USAGE
		*/
		{
			name: "sal temperature broadcast event",
			args: args{
				tag: &salTag{
					tagType:     SAL,
					application: readwriteModel.ApplicationIdContainer_TEMPERATURE_BROADCAST_19,
					salCommand:  "BROADCAST_EVENT",
					numElements: 2,
				},
				value: []string{
					"1",
					"2",
				},
			},
			want: spiValues.NewPlcList([]apiValues.PlcValue{
				spiValues.NewPlcBYTE(1),
				spiValues.NewPlcBYTE(2),
			}),
			wantErr: assert.NoError,
		},
		/*
			TODO: add tests for that
			SAL
			ApplicationId_ROOM_CONTROL_SYSTEM
		*/
		/*
					note: those are all ApplicationId_LIGHTING based
			ApplicationId_VENTILATION
			ApplicationId_IRRIGATION_CONTROL
			ApplicationId_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL
			ApplicationId_HEATING
			ApplicationId_AUDIO_AND_VIDEO
			ApplicationId_HVAC_ACTUATOR
		*/
		{
			name: "sal lighting OFF",
			args: args{
				tag: &salTag{
					tagType:     SAL,
					application: readwriteModel.ApplicationIdContainer_LIGHTING_3A,
					salCommand:  "OFF",
					numElements: 1,
				},
				value: "1",
			},
			want:    spiValues.NewPlcBYTE(1),
			wantErr: assert.NoError,
		},
		{
			name: "sal lighting ON",
			args: args{
				tag: &salTag{
					tagType:     SAL,
					application: readwriteModel.ApplicationIdContainer_LIGHTING_3A,
					salCommand:  "ON",
					numElements: 1,
				},
				value: "1",
			},
			want:    spiValues.NewPlcBYTE(1),
			wantErr: assert.NoError,
		},
		{
			name: "sal lighting RAMP_TO_LEVEL",
			args: args{
				tag: &salTag{
					tagType:     SAL,
					application: readwriteModel.ApplicationIdContainer_LIGHTING_3A,
					salCommand:  "RAMP_TO_LEVEL",
					numElements: 2,
				},
				value: []string{"1", "2"},
			},
			want: spiValues.NewPlcList([]apiValues.PlcValue{
				spiValues.NewPlcBYTE(1),
				spiValues.NewPlcBYTE(2),
			}),
			wantErr: assert.NoError,
		},
		{
			name: "sal lighting TERMINATE_RAMP",
			args: args{
				tag: &salTag{
					tagType:     SAL,
					application: readwriteModel.ApplicationIdContainer_LIGHTING_3A,
					salCommand:  "TERMINATE_RAMP",
					numElements: 1,
				},
				value: "1",
			},
			want:    spiValues.NewPlcBYTE(1),
			wantErr: assert.NoError,
		},
		// TODO: implement SAL LIGHTING LABEL
		// TODO: implement remaining tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := ValueHandler{
				DefaultValueHandler: tt.fields.DefaultValueHandler,
			}
			got, err := m.NewPlcValue(tt.args.tag, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("NewPlcValue(\n%v, \n%v)", tt.args.tag, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewPlcValue(%v, %v)", tt.args.tag, tt.args.value)
		})
	}
}
