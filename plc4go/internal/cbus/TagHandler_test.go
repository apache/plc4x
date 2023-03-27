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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/stretchr/testify/assert"
	"reflect"
	"regexp"
	"testing"
)

func TestNewTagHandler(t *testing.T) {
	tests := []struct {
		name          string
		wantValidator func(t *testing.T, handler TagHandler)
	}{
		{
			name: "simple",
			wantValidator: func(t *testing.T, handler TagHandler) {
				assert.NotNil(t, handler.statusRequestPattern)
				assert.NotNil(t, handler.calPattern)
				assert.NotNil(t, handler.salPattern)
				assert.NotNil(t, handler.salMonitorPattern)
				assert.NotNil(t, handler.mmiMonitorPattern)
				assert.NotNil(t, handler.unityQuery)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTagHandler()
			tt.wantValidator(t, got)
		})
	}
}

func TestTagHandler_ParseTag(t *testing.T) {
	type args struct {
		tagAddress string
	}
	tests := []struct {
		name    string
		args    args
		want    apiModel.PlcTag
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name:    "nonsense",
			args:    args{tagAddress: "gobblegobble"},
			wantErr: true,
		},
		{
			name: "status request",
			args: args{tagAddress: "status/binary/LIGHTING"},
			want: &statusTag{
				bridgeAddresses:           nil,
				tagType:                   STATUS,
				statusRequestType:         StatusRequestTypeBinaryState,
				startingGroupAddressLabel: nil,
				application:               readWriteModel.ApplicationIdContainer_LIGHTING_38,
				numElements:               1,
			},
		},
		{
			name: "cal get status",
			args: args{tagAddress: "cal/2/getStatus=1,2"},
			want: &calGetStatusTag{
				calTag: calTag{
					bridgeAddresses: nil,
					unitAddress:     readWriteModel.NewUnitAddress(2),
				},
				tagType:     CAL_GETSTATUS,
				parameter:   1,
				count:       2,
				numElements: 1,
			},
		},
		{
			name: "sal light on",
			args: args{tagAddress: "sal/LIGHTING/ON"},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_LIGHTING_38,
				salCommand:  "ON",
				numElements: 1,
			},
		},
		{
			name: "sal monitor",
			args: args{tagAddress: "salmonitor/2/LIGHTING"},
			want: &salMonitorTag{
				tagType:     SAL_MONITOR,
				unitAddress: readWriteModel.NewUnitAddress(2),
				application: func() *readWriteModel.ApplicationIdContainer {
					lighting_38 := readWriteModel.ApplicationIdContainer_LIGHTING_38
					return &lighting_38
				}(),
				numElements: 1,
			},
		},
		{
			name: "mmi monitor",
			args: args{tagAddress: "mmimonitor/2/LIGHTING"},
			want: &mmiMonitorTag{
				tagType:     MMI_STATUS_MONITOR,
				unitAddress: readWriteModel.NewUnitAddress(2),
				application: func() *readWriteModel.ApplicationIdContainer {
					lighting_38 := readWriteModel.ApplicationIdContainer_LIGHTING_38
					return &lighting_38
				}(),
				numElements: 1,
			},
		},
		{
			name: "bridged status 1 bridge",
			args: args{tagAddress: "status/b1/binary/LIGHTING"},
			want: &statusTag{
				bridgeAddresses:           []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(1)},
				tagType:                   STATUS,
				statusRequestType:         StatusRequestTypeBinaryState,
				startingGroupAddressLabel: nil,
				application:               readWriteModel.ApplicationIdContainer_LIGHTING_38,
				numElements:               1,
			},
		},
		{
			name: "bridged status 6 bridges",
			args: args{tagAddress: "status/b1-b2-b3-b4-b5-b6/binary/LIGHTING"},
			want: &statusTag{
				bridgeAddresses:           []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewBridgeAddress(3), readWriteModel.NewBridgeAddress(4), readWriteModel.NewBridgeAddress(5), readWriteModel.NewBridgeAddress(6)},
				tagType:                   STATUS,
				statusRequestType:         StatusRequestTypeBinaryState,
				startingGroupAddressLabel: nil,
				application:               readWriteModel.ApplicationIdContainer_LIGHTING_38,
				numElements:               1,
			},
		},
		{
			name: "bridged cal 1 bridge",
			args: args{tagAddress: "cal/b1-u2/getStatus=1,2"},
			want: &calGetStatusTag{
				calTag: calTag{
					bridgeAddresses: []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(1)},
					unitAddress:     readWriteModel.NewUnitAddress(2),
				},
				tagType:     CAL_GETSTATUS,
				parameter:   1,
				count:       2,
				numElements: 1,
			},
		},
		{
			name: "bridged cal 6 bridges",
			args: args{tagAddress: "cal/b1-b2-b3-b4-b5-b6-u2/getStatus=1,2"},
			want: &calGetStatusTag{
				calTag: calTag{
					bridgeAddresses: []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewBridgeAddress(3), readWriteModel.NewBridgeAddress(4), readWriteModel.NewBridgeAddress(5), readWriteModel.NewBridgeAddress(6)},
					unitAddress:     readWriteModel.NewUnitAddress(2),
				},
				tagType:     CAL_GETSTATUS,
				parameter:   1,
				count:       2,
				numElements: 1,
			},
		},
		{
			name: "sal light on 1 bridge",
			args: args{tagAddress: "sal/b1/LIGHTING/ON"},
			want: &salTag{
				bridgeAddresses: []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(1)},
				tagType:         SAL,
				application:     readWriteModel.ApplicationIdContainer_LIGHTING_38,
				salCommand:      "ON",
				numElements:     1,
			},
		},
		{
			name: "sal light on 6 bridges",
			args: args{tagAddress: "sal/b1-b2-b3-b4-b5-b6/LIGHTING/ON"},
			want: &salTag{
				bridgeAddresses: []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewBridgeAddress(3), readWriteModel.NewBridgeAddress(4), readWriteModel.NewBridgeAddress(5), readWriteModel.NewBridgeAddress(6)},
				tagType:         SAL,
				application:     readWriteModel.ApplicationIdContainer_LIGHTING_38,
				salCommand:      "ON",
				numElements:     1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTagHandler()
			got, err := m.ParseTag(tt.args.tagAddress)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestTagHandler_ParseQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		args    args
		want    apiModel.PlcQuery
		wantErr bool
	}{
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name:    "nonsense",
			args:    args{query: "gobblegobble"},
			wantErr: true,
		},
		// TODO: other cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTagHandler()
			got, err := m.ParseQuery(tt.args.query)
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

func TestTagHandler_handleStatusRequestPattern(t *testing.T) {
	type fields struct {
		statusRequestPattern *regexp.Regexp
		calPattern           *regexp.Regexp
		salPattern           *regexp.Regexp
		salMonitorPattern    *regexp.Regexp
		mmiMonitorPattern    *regexp.Regexp
		unityQuery           *regexp.Regexp
	}
	type args struct {
		match map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    apiModel.PlcTag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TagHandler{
				statusRequestPattern: tt.fields.statusRequestPattern,
				calPattern:           tt.fields.calPattern,
				salPattern:           tt.fields.salPattern,
				salMonitorPattern:    tt.fields.salMonitorPattern,
				mmiMonitorPattern:    tt.fields.mmiMonitorPattern,
				unityQuery:           tt.fields.unityQuery,
			}
			got, err := m.handleStatusRequestPattern(tt.args.match)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleStatusRequestPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleStatusRequestPattern() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagHandler_handleCalPattern(t *testing.T) {
	type fields struct {
		statusRequestPattern *regexp.Regexp
		calPattern           *regexp.Regexp
		salPattern           *regexp.Regexp
		salMonitorPattern    *regexp.Regexp
		mmiMonitorPattern    *regexp.Regexp
		unityQuery           *regexp.Regexp
	}
	type args struct {
		match map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    apiModel.PlcTag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TagHandler{
				statusRequestPattern: tt.fields.statusRequestPattern,
				calPattern:           tt.fields.calPattern,
				salPattern:           tt.fields.salPattern,
				salMonitorPattern:    tt.fields.salMonitorPattern,
				mmiMonitorPattern:    tt.fields.mmiMonitorPattern,
				unityQuery:           tt.fields.unityQuery,
			}
			got, err := m.handleCalPattern(tt.args.match)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleCalPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleCalPattern() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagHandler_handleMMIMonitorPattern(t *testing.T) {
	type fields struct {
		statusRequestPattern *regexp.Regexp
		calPattern           *regexp.Regexp
		salPattern           *regexp.Regexp
		salMonitorPattern    *regexp.Regexp
		mmiMonitorPattern    *regexp.Regexp
		unityQuery           *regexp.Regexp
	}
	type args struct {
		match map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    apiModel.PlcTag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TagHandler{
				statusRequestPattern: tt.fields.statusRequestPattern,
				calPattern:           tt.fields.calPattern,
				salPattern:           tt.fields.salPattern,
				salMonitorPattern:    tt.fields.salMonitorPattern,
				mmiMonitorPattern:    tt.fields.mmiMonitorPattern,
				unityQuery:           tt.fields.unityQuery,
			}
			got, err := m.handleMMIMonitorPattern(tt.args.match)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleMMIMonitorPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleMMIMonitorPattern() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagHandler_handleSALMonitorPattern(t *testing.T) {
	type fields struct {
		statusRequestPattern *regexp.Regexp
		calPattern           *regexp.Regexp
		salPattern           *regexp.Regexp
		salMonitorPattern    *regexp.Regexp
		mmiMonitorPattern    *regexp.Regexp
		unityQuery           *regexp.Regexp
	}
	type args struct {
		match map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    apiModel.PlcTag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TagHandler{
				statusRequestPattern: tt.fields.statusRequestPattern,
				calPattern:           tt.fields.calPattern,
				salPattern:           tt.fields.salPattern,
				salMonitorPattern:    tt.fields.salMonitorPattern,
				mmiMonitorPattern:    tt.fields.mmiMonitorPattern,
				unityQuery:           tt.fields.unityQuery,
			}
			got, err := m.handleSALMonitorPattern(tt.args.match)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleSALMonitorPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleSALMonitorPattern() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagHandler_handleSALPattern(t *testing.T) {
	type fields struct {
		statusRequestPattern *regexp.Regexp
		calPattern           *regexp.Regexp
		salPattern           *regexp.Regexp
		salMonitorPattern    *regexp.Regexp
		mmiMonitorPattern    *regexp.Regexp
		unityQuery           *regexp.Regexp
	}
	type args struct {
		match map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    apiModel.PlcTag
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TagHandler{
				statusRequestPattern: tt.fields.statusRequestPattern,
				calPattern:           tt.fields.calPattern,
				salPattern:           tt.fields.salPattern,
				salMonitorPattern:    tt.fields.salMonitorPattern,
				mmiMonitorPattern:    tt.fields.mmiMonitorPattern,
				unityQuery:           tt.fields.unityQuery,
			}
			got, err := m.handleSALPattern(tt.args.match)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleSALPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleSALPattern() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagHandler_handleUnitQuery(t *testing.T) {
	type fields struct {
		statusRequestPattern *regexp.Regexp
		calPattern           *regexp.Regexp
		salPattern           *regexp.Regexp
		salMonitorPattern    *regexp.Regexp
		mmiMonitorPattern    *regexp.Regexp
		unityQuery           *regexp.Regexp
	}
	type args struct {
		match map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    apiModel.PlcQuery
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := TagHandler{
				statusRequestPattern: tt.fields.statusRequestPattern,
				calPattern:           tt.fields.calPattern,
				salPattern:           tt.fields.salPattern,
				salMonitorPattern:    tt.fields.salMonitorPattern,
				mmiMonitorPattern:    tt.fields.mmiMonitorPattern,
				unityQuery:           tt.fields.unityQuery,
			}
			got, err := m.handleUnitQuery(tt.args.match)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleUnitQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handleUnitQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagType_GetName(t *testing.T) {
	tests := []struct {
		name string
		i    TagType
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_applicationIdFromArgument(t *testing.T) {
	type args struct {
		applicationIdArgument string
	}
	tests := []struct {
		name    string
		args    args
		want    readWriteModel.ApplicationIdContainer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applicationIdFromArgument(tt.args.applicationIdArgument)
			if (err != nil) != tt.wantErr {
				t.Errorf("applicationIdFromArgument() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("applicationIdFromArgument() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_c2nl(t *testing.T) {
	type args struct {
		t []CommandAndArgumentsCount
	}
	tests := []struct {
		name string
		args args
		want []CommandAndArgumentsCount
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c2nl(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("c2nl() = %v, want %v", got, tt.want)
			}
		})
	}
}
