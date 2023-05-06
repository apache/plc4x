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
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	"github.com/stretchr/testify/assert"
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
			name: "sal monitor all",
			args: args{tagAddress: "salmonitor/*/*"},
			want: &salMonitorTag{
				tagType:     SAL_MONITOR,
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
			name: "mmi monitor all",
			args: args{tagAddress: "mmimonitor/*/*"},
			want: &mmiMonitorTag{
				tagType:     MMI_STATUS_MONITOR,
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
			assert.Equal(t, tt.want, got)
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
		{
			name: "unitQuery all",
			args: args{
				query: "info/*/*",
			},
			want: &unitInfoQuery{
				tagType:     UNIT_INFO,
				unitAddress: nil,
				attribute:   nil,
				numElements: 1,
			},
		},
		{
			name: "unitQuery",
			args: args{
				query: "info/0x13/DSIStatus",
			},
			want: &unitInfoQuery{
				tagType:     UNIT_INFO,
				unitAddress: readWriteModel.NewUnitAddress(19),
				attribute: func() *readWriteModel.Attribute {
					var attribute readWriteModel.Attribute
					attribute = readWriteModel.Attribute_DSIStatus
					return &attribute
				}(),
				numElements: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTagHandler()
			got, err := m.ParseQuery(tt.args.query)
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

func TestTagHandler_applicationIdFromArgument(t *testing.T) {
	type args struct {
		applicationIdArgument string
	}
	tests := []struct {
		name    string
		args    args
		want    readWriteModel.ApplicationIdContainer
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "empty",
			args:    args{applicationIdArgument: ""},
			wantErr: assert.Error,
		},
		{
			name:    "number",
			args:    args{applicationIdArgument: "56"},
			want:    readWriteModel.ApplicationIdContainer_LIGHTING_38,
			wantErr: assert.NoError,
		},
		{
			name:    "hex",
			args:    args{applicationIdArgument: "0x38"},
			want:    readWriteModel.ApplicationIdContainer_LIGHTING_38,
			wantErr: assert.NoError,
		},
		{
			name:    "by name",
			args:    args{applicationIdArgument: "LIGHTING"},
			want:    readWriteModel.ApplicationIdContainer_LIGHTING_38,
			wantErr: assert.NoError,
		},
		{
			name:    "by exact name",
			args:    args{applicationIdArgument: "LIGHTING_39"},
			want:    readWriteModel.ApplicationIdContainer_LIGHTING_39,
			wantErr: assert.NoError,
		},
		{
			name:    "TEMPERATURE_BROADCAST",
			args:    args{applicationIdArgument: "TEMPERATURE_BROADCAST"},
			want:    readWriteModel.ApplicationIdContainer_TEMPERATURE_BROADCAST_19,
			wantErr: assert.NoError,
		},
		{
			name:    "ROOM_CONTROL_SYSTEM",
			args:    args{applicationIdArgument: "ROOM_CONTROL_SYSTEM"},
			want:    readWriteModel.ApplicationIdContainer_ROOM_CONTROL_SYSTEM_26,
			wantErr: assert.NoError,
		},
		{
			name:    "LIGHTING",
			args:    args{applicationIdArgument: "LIGHTING"},
			want:    readWriteModel.ApplicationIdContainer_LIGHTING_38,
			wantErr: assert.NoError,
		},
		{
			name:    "VENTILATION",
			args:    args{applicationIdArgument: "VENTILATION"},
			want:    readWriteModel.ApplicationIdContainer_VENTILATION_70,
			wantErr: assert.NoError,
		},
		{
			name:    "IRRIGATION_CONTROL",
			args:    args{applicationIdArgument: "IRRIGATION_CONTROL"},
			want:    readWriteModel.ApplicationIdContainer_IRRIGATION_CONTROL_71,
			wantErr: assert.NoError,
		},
		{
			name:    "POOLS_SPAS_PONDS_FOUNTAINS_CONTROL",
			args:    args{applicationIdArgument: "POOLS_SPAS_PONDS_FOUNTAINS_CONTROL"},
			want:    readWriteModel.ApplicationIdContainer_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL_72,
			wantErr: assert.NoError,
		},
		{
			name:    "HEATING",
			args:    args{applicationIdArgument: "HEATING"},
			want:    readWriteModel.ApplicationIdContainer_HEATING_88,
			wantErr: assert.NoError,
		},
		{
			name:    "AIR_CONDITIONING",
			args:    args{applicationIdArgument: "AIR_CONDITIONING"},
			want:    readWriteModel.ApplicationIdContainer_AIR_CONDITIONING_AC,
			wantErr: assert.NoError,
		},
		{
			name:    "TRIGGER_CONTROL",
			args:    args{applicationIdArgument: "TRIGGER_CONTROL"},
			want:    readWriteModel.ApplicationIdContainer_TRIGGER_CONTROL_CA,
			wantErr: assert.NoError,
		},
		{
			name:    "ENABLE_CONTROL",
			args:    args{applicationIdArgument: "ENABLE_CONTROL"},
			want:    readWriteModel.ApplicationIdContainer_ENABLE_CONTROL_CB,
			wantErr: assert.NoError,
		},
		{
			name:    "AUDIO_AND_VIDEO",
			args:    args{applicationIdArgument: "AUDIO_AND_VIDEO"},
			want:    readWriteModel.ApplicationIdContainer_AUDIO_AND_VIDEO_CD,
			wantErr: assert.NoError,
		},
		{
			name:    "SECURITY",
			args:    args{applicationIdArgument: "SECURITY"},
			want:    readWriteModel.ApplicationIdContainer_SECURITY_D0,
			wantErr: assert.NoError,
		},
		{
			name:    "METERING",
			args:    args{applicationIdArgument: "METERING"},
			want:    readWriteModel.ApplicationIdContainer_METERING_D1,
			wantErr: assert.NoError,
		},
		{
			name:    "ACCESS_CONTROL",
			args:    args{applicationIdArgument: "ACCESS_CONTROL"},
			want:    readWriteModel.ApplicationIdContainer_ACCESS_CONTROL_D5,
			wantErr: assert.NoError,
		},
		{
			name:    "CLOCK_AND_TIMEKEEPING",
			args:    args{applicationIdArgument: "CLOCK_AND_TIMEKEEPING"},
			want:    readWriteModel.ApplicationIdContainer_CLOCK_AND_TIMEKEEPING_DF,
			wantErr: assert.NoError,
		},
		{
			name:    "TELEPHONY_STATUS_AND_CONTROL",
			args:    args{applicationIdArgument: "TELEPHONY_STATUS_AND_CONTROL"},
			want:    readWriteModel.ApplicationIdContainer_TELEPHONY_STATUS_AND_CONTROL_E0,
			wantErr: assert.NoError,
		},
		{
			name:    "MEASUREMENT",
			args:    args{applicationIdArgument: "MEASUREMENT"},
			want:    readWriteModel.ApplicationIdContainer_MEASUREMENT_E4,
			wantErr: assert.NoError,
		},
		{
			name:    "TESTING",
			args:    args{applicationIdArgument: "TESTING"},
			want:    readWriteModel.ApplicationIdContainer_TESTING_FA,
			wantErr: assert.NoError,
		},
		{
			name:    "MEDIA_TRANSPORT_CONTROL",
			args:    args{applicationIdArgument: "MEDIA_TRANSPORT_CONTROL"},
			want:    readWriteModel.ApplicationIdContainer_MEDIA_TRANSPORT_CONTROL_C0,
			wantErr: assert.NoError,
		},
		{
			name:    "ERROR_REPORTING",
			args:    args{applicationIdArgument: "ERROR_REPORTING"},
			want:    readWriteModel.ApplicationIdContainer_ERROR_REPORTING_CE,
			wantErr: assert.NoError,
		},
		{
			name:    "HVAC_ACTUATOR",
			args:    args{applicationIdArgument: "HVAC_ACTUATOR"},
			want:    readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_73,
			wantErr: assert.NoError,
		},
		{
			name:    "INFO_MESSAGES",
			args:    args{applicationIdArgument: "INFO_MESSAGES"},
			want:    readWriteModel.ApplicationIdContainer_INFO_MESSAGES,
			wantErr: assert.NoError,
		},
		{
			name:    "NETWORK_CONTROL",
			args:    args{applicationIdArgument: "NETWORK_CONTROL"},
			want:    readWriteModel.ApplicationIdContainer_NETWORK_CONTROL,
			wantErr: assert.NoError,
		},
		{
			name:    "RESERVED",
			args:    args{applicationIdArgument: "RESERVED"},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTagHandler()
			got, err := m.applicationIdFromArgument(tt.args.applicationIdArgument)
			if !tt.wantErr(t, err, fmt.Sprintf("applicationIdFromArgument(%v)", tt.args.applicationIdArgument)) {
				return
			}
			assert.Equalf(t, tt.want, got, "applicationIdFromArgument(%v)", tt.args.applicationIdArgument)
		})
	}
}

func TestTagHandler_extractBridges(t *testing.T) {
	type args struct {
		match map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []readWriteModel.BridgeAddress
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "no matches",
			args:    args{match: map[string]string{}},
			wantErr: assert.NoError,
		},
		{
			name: "empty bridges",
			args: args{match: map[string]string{
				"bridges": "",
			}},
			wantErr: assert.NoError,
		},
		{
			name: "one bridge",
			args: args{match: map[string]string{
				"bridges": "3",
			}},
			want:    []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(3)},
			wantErr: assert.NoError,
		},
		{
			name: "one bridge prefixed",
			args: args{match: map[string]string{
				"bridges": "b3",
			}},
			want:    []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(3)},
			wantErr: assert.NoError,
		},
		{
			name: "one hex bridge",
			args: args{match: map[string]string{
				"bridges": "0x03",
			}},
			want:    []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(3)},
			wantErr: assert.NoError,
		},
		{
			name: "two bridges",
			args: args{match: map[string]string{
				"bridges": "3-4",
			}},
			want:    []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(3), readWriteModel.NewBridgeAddress(4)},
			wantErr: assert.NoError,
		},
		{
			name: "two bridges prefixed",
			args: args{match: map[string]string{
				"bridges": "b3-b4",
			}},
			want:    []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(3), readWriteModel.NewBridgeAddress(4)},
			wantErr: assert.NoError,
		},
		{
			name: "two hex bridges prefixed",
			args: args{match: map[string]string{
				"bridges": "b0x03-b0x04",
			}},
			want:    []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(3), readWriteModel.NewBridgeAddress(4)},
			wantErr: assert.NoError,
		},
		{
			name: "six hex bridges prefixed",
			args: args{match: map[string]string{
				"bridges": "b0x01-b0x02-b0x03-b0x04-b0x05-b0x06",
			}},
			want:    []readWriteModel.BridgeAddress{readWriteModel.NewBridgeAddress(1), readWriteModel.NewBridgeAddress(2), readWriteModel.NewBridgeAddress(3), readWriteModel.NewBridgeAddress(4), readWriteModel.NewBridgeAddress(5), readWriteModel.NewBridgeAddress(6)},
			wantErr: assert.NoError,
		},
		{
			name: "seven hex bridges prefixed",
			args: args{match: map[string]string{
				"bridges": "b0x01-b0x02-b0x03-b0x04-b0x05-b0x06-b0x07",
			}},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTagHandler()
			got, err := m.extractBridges(tt.args.match)
			if !tt.wantErr(t, err, fmt.Sprintf("extractBridges(%v)", tt.args.match)) {
				return
			}
			assert.Equalf(t, tt.want, got, "extractBridges(%v)", tt.args.match)
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
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "empty",
			wantErr: assert.Error,
		},
		{
			name: "invalid cal type",
			args: args{match: map[string]string{
				"unitAddress": "0",
				"calType":     "invalid",
			}},
			wantErr: assert.Error,
		},
		// TODO: implement reset=
		{
			name: "recall unknown param",
			args: args{match: map[string]string{
				"unitAddress": "0",
				"calType":     "recall=",
			}},
			wantErr: assert.Error,
		},
		{
			name: "recall invalid recall count",
			args: args{match: map[string]string{
				"unitAddress":   "0",
				"calType":       "recall=",
				"recallParamNo": "1",
			}},
			wantErr: assert.Error,
		},
		{
			name: "recall",
			args: args{match: map[string]string{
				"unitAddress":   "0",
				"calType":       "recall=",
				"recallParamNo": "1",
				"recallCount":   "2",
			}},
			want: &calRecallTag{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(0),
				},
				tagType:     CAL_RECALL,
				parameter:   1,
				count:       2,
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "recall hex",
			args: args{match: map[string]string{
				"unitAddress":   "0",
				"calType":       "recall=",
				"recallParamNo": "0x30",
				"recallCount":   "2",
			}},
			want: &calRecallTag{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(0),
				},
				tagType:     CAL_RECALL,
				parameter:   readWriteModel.Parameter_INTERFACE_OPTIONS_1,
				count:       2,
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "recall name",
			args: args{match: map[string]string{
				"unitAddress":   "0",
				"calType":       "recall=",
				"recallParamNo": "INTERFACE_OPTIONS_1",
				"recallCount":   "2",
			}},
			want: &calRecallTag{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(0),
				},
				tagType:     CAL_RECALL,
				parameter:   readWriteModel.Parameter_INTERFACE_OPTIONS_1,
				count:       2,
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "identify unknown param",
			args: args{match: map[string]string{
				"unitAddress": "0",
				"calType":     "identify=",
			}},
			wantErr: assert.Error,
		},
		{
			name: "identify with number",
			args: args{match: map[string]string{
				"unitAddress":       "0",
				"calType":           "identify=",
				"identifyAttribute": "1",
			}},
			want: &calIdentifyTag{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(0),
				},
				tagType:     CAL_IDENTIFY,
				attribute:   readWriteModel.Attribute_Type,
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "identify with hex",
			args: args{match: map[string]string{
				"unitAddress":       "0",
				"calType":           "identify=",
				"identifyAttribute": "0x01",
			}},
			want: &calIdentifyTag{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(0),
				},
				tagType:     CAL_IDENTIFY,
				attribute:   readWriteModel.Attribute_Type,
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "identify with name",
			args: args{match: map[string]string{
				"unitAddress":       "0",
				"calType":           "identify=",
				"identifyAttribute": "Type",
			}},
			want: &calIdentifyTag{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(0),
				},
				tagType:     CAL_IDENTIFY,
				attribute:   readWriteModel.Attribute_Type,
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "getStatus unknown param",
			args: args{match: map[string]string{
				"unitAddress": "0",
				"calType":     "getStatus=",
			}},
			wantErr: assert.Error,
		},
		{
			name: "getStatus invalid recall count",
			args: args{match: map[string]string{
				"unitAddress":      "0",
				"calType":          "getStatus=",
				"getStatusParamNo": "1",
			}},
			wantErr: assert.Error,
		},
		{
			name: "getStatus by number",
			args: args{match: map[string]string{
				"unitAddress":      "0",
				"calType":          "getStatus=",
				"getStatusParamNo": "1",
				"getStatusCount":   "2",
			}},
			want: &calGetStatusTag{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(0),
				},
				tagType:     CAL_GETSTATUS,
				parameter:   readWriteModel.Parameter_UNKNOWN_02,
				count:       2,
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "getStatus by hex",
			args: args{match: map[string]string{
				"unitAddress":      "0",
				"calType":          "getStatus=",
				"getStatusParamNo": "0x01",
				"getStatusCount":   "2",
			}},
			want: &calGetStatusTag{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(0),
				},
				tagType:     CAL_GETSTATUS,
				parameter:   readWriteModel.Parameter_UNKNOWN_02,
				count:       2,
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "getStatus by name",
			args: args{match: map[string]string{
				"unitAddress":      "0",
				"calType":          "getStatus=",
				"getStatusParamNo": "INTERFACE_OPTIONS_1",
				"getStatusCount":   "2",
			}},
			want: &calGetStatusTag{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(0),
				},
				tagType:     CAL_GETSTATUS,
				parameter:   readWriteModel.Parameter_INTERFACE_OPTIONS_1,
				count:       2,
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
		// TODO: implement write=
		// TODO: implement identifyReply=
		// TODO: implement reply=
		// TODO: implement status=
		// TODO: implement statusExtended=
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
			if !tt.wantErr(t, err, fmt.Sprintf("handleCalPattern(%v)", tt.args.match)) {
				return
			}
			assert.Equalf(t, tt.want, got, "handleCalPattern(%v)", tt.args.match)
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
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "error",
			args: args{match: map[string]string{
				"unitAddress": "blub",
				"application": "blub",
			}},
			wantErr: assert.Error,
		},
		{
			name: "wildcards",
			args: args{match: map[string]string{
				"unitAddress": "*",
				"application": "*",
			}},
			want: &mmiMonitorTag{
				tagType:     MMI_STATUS_MONITOR,
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
		{
			name: "fixed",
			args: args{match: map[string]string{
				"unitAddress": "2",
				"application": "LIGHTING",
			}},
			want: &mmiMonitorTag{
				tagType:     MMI_STATUS_MONITOR,
				unitAddress: readWriteModel.NewUnitAddress(2),
				application: func() *readWriteModel.ApplicationIdContainer {
					var applicationId readWriteModel.ApplicationIdContainer
					applicationId = readWriteModel.ApplicationIdContainer_LIGHTING_38
					return &applicationId
				}(),
				numElements: 1,
			},
			wantErr: assert.NoError,
		},
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
			if !tt.wantErr(t, err, fmt.Sprintf("handleMMIMonitorPattern(%v)", tt.args.match)) {
				return
			}
			assert.Equalf(t, tt.want, got, "handleMMIMonitorPattern(%v)", tt.args.match)
		})
	}
}

func TestTagHandler_handleStatusRequestPattern(t *testing.T) {
	type args struct {
		match map[string]string
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
			name: "statusRequestType empty",
			args: args{match: map[string]string{
				"statusRequestType": "",
			}},
			wantErr: true,
		},
		{
			name: "statusRequestType minimal",
			args: args{match: map[string]string{
				"statusRequestType": "",
				"application":       "2",
			}},
			want: &statusTag{
				tagType:           STATUS,
				statusRequestType: StatusRequestTypeBinaryState,
				application:       readWriteModel.ApplicationIdContainer_FREE_USAGE_02,
				numElements:       1,
			},
		},
		{
			name: "statusRequestType binary",
			args: args{match: map[string]string{
				"statusRequestType": "doesn't matter",
				"binary":            "doesn't matter",
				"application":       "2",
			}},
			want: &statusTag{
				tagType:           STATUS,
				statusRequestType: StatusRequestTypeBinaryState,
				application:       readWriteModel.ApplicationIdContainer_FREE_USAGE_02,
				numElements:       1,
			},
		},
		{
			name: "statusRequestType level",
			args: args{match: map[string]string{
				"statusRequestType":         "doesn't matter",
				"startingGroupAddressLabel": "20",
				"application":               "2",
			}},
			want: &statusTag{
				tagType:           STATUS,
				statusRequestType: StatusRequestTypeLevel,
				startingGroupAddressLabel: func() *byte {
					var level byte
					level = 0x20
					return &level
				}(),
				application: readWriteModel.ApplicationIdContainer_FREE_USAGE_02,
				numElements: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTagHandler()
			got, err := m.handleStatusRequestPattern(tt.args.match)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleStatusRequestPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("handleStatusRequestPattern() got = %v, want %v", got, tt.want)
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
		{
			name: "error",
			args: args{match: map[string]string{
				"unitAddress": "blub",
				"application": "blub",
			}},
			wantErr: true,
		},
		{
			name: "wildcards",
			args: args{match: map[string]string{
				"unitAddress": "*",
				"application": "*",
			}},
			want: &salMonitorTag{
				tagType:     SAL_MONITOR,
				numElements: 1,
			},
		},
		{
			name: "fixed",
			args: args{match: map[string]string{
				"unitAddress": "2",
				"application": "LIGHTING",
			}},
			want: &salMonitorTag{
				tagType:     SAL_MONITOR,
				unitAddress: readWriteModel.NewUnitAddress(2),
				application: func() *readWriteModel.ApplicationIdContainer {
					var applicationId readWriteModel.ApplicationIdContainer
					applicationId = readWriteModel.ApplicationIdContainer_LIGHTING_38
					return &applicationId
				}(),
				numElements: 1,
			},
		},
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
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("handleSALMonitorPattern() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagHandler_handleSALPattern(t *testing.T) {
	type args struct {
		match map[string]string
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
			name: "empty sal command",
			args: args{match: map[string]string{
				"application": "LIGHTING",
			}},
			wantErr: true,
		},
		// TODO: implement ApplicationId_RESERVED
		// TODO: implement ApplicationId_FREE_USAGE
		{
			name: "TEMPERATURE_BROADCAST BROADCAST_EVENT",
			args: args{match: map[string]string{
				"application": "TEMPERATURE_BROADCAST",
				"salCommand":  "BROADCAST_EVENT",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_TEMPERATURE_BROADCAST_19,
				salCommand:  "BROADCAST_EVENT",
				numElements: 2,
			},
		},
		// TODO: implement ApplicationId_ROOM_CONTROL_SYSTEM
		{
			name: "LIGHTING",
			args: args{match: map[string]string{
				"application": "LIGHTING",
				"salCommand":  "TERMINATE_RAMP",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_LIGHTING_38,
				salCommand:  "TERMINATE_RAMP",
				numElements: 1,
			},
		},
		{
			name: "VENTILATION",
			args: args{match: map[string]string{
				"application": "VENTILATION",
				"salCommand":  "TERMINATE_RAMP",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_VENTILATION_70,
				salCommand:  "TERMINATE_RAMP",
				numElements: 1,
			},
		},
		{
			name: "IRRIGATION_CONTROL",
			args: args{match: map[string]string{
				"application": "IRRIGATION_CONTROL",
				"salCommand":  "TERMINATE_RAMP",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_IRRIGATION_CONTROL_71,
				salCommand:  "TERMINATE_RAMP",
				numElements: 1,
			},
		},
		{
			name: "POOLS_SPAS_PONDS_FOUNTAINS_CONTROL",
			args: args{match: map[string]string{
				"application": "POOLS_SPAS_PONDS_FOUNTAINS_CONTROL",
				"salCommand":  "TERMINATE_RAMP",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_POOLS_SPAS_PONDS_FOUNTAINS_CONTROL_72,
				salCommand:  "TERMINATE_RAMP",
				numElements: 1,
			},
		},
		{
			name: "HEATING",
			args: args{match: map[string]string{
				"application": "HEATING",
				"salCommand":  "TERMINATE_RAMP",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_HEATING_88,
				salCommand:  "TERMINATE_RAMP",
				numElements: 1,
			},
		},
		{
			name: "AIR_CONDITIONING",
			args: args{match: map[string]string{
				"application": "AIR_CONDITIONING",
				"salCommand":  "ZONE_HUMIDITY_PLANT_STATUS",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_AIR_CONDITIONING_AC,
				salCommand:  "ZONE_HUMIDITY_PLANT_STATUS",
				numElements: 5,
			},
		},
		{
			name: "TRIGGER_CONTROL",
			args: args{match: map[string]string{
				"application": "TRIGGER_CONTROL",
				"salCommand":  "TRIGGER_MAX",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_TRIGGER_CONTROL_CA,
				salCommand:  "TRIGGER_MAX",
				numElements: 0,
			},
		},
		{
			name: "ENABLE_CONTROL",
			args: args{match: map[string]string{
				"application": "ENABLE_CONTROL",
				"salCommand":  "SET_NETWORK_VARIABLE",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_ENABLE_CONTROL_CB,
				salCommand:  "SET_NETWORK_VARIABLE",
				numElements: 1,
			},
		},
		{
			name: "AUDIO_AND_VIDEO",
			args: args{match: map[string]string{
				"application": "AUDIO_AND_VIDEO",
				"salCommand":  "RAMP_TO_LEVEL",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_AUDIO_AND_VIDEO_CD,
				salCommand:  "RAMP_TO_LEVEL",
				numElements: 2,
			},
		},
		{
			name: "SECURITY",
			args: args{match: map[string]string{
				"application": "SECURITY",
				"salCommand":  "EVENT",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_SECURITY_D0,
				salCommand:  "EVENT",
				numElements: 255,
			},
		},
		{
			name: "METERING",
			args: args{match: map[string]string{
				"application": "METERING",
				"salCommand":  "EVENT",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_METERING_D1,
				salCommand:  "EVENT",
				numElements: 255,
			},
		},
		{
			name: "ACCESS_CONTROL",
			args: args{match: map[string]string{
				"application": "ACCESS_CONTROL",
				"salCommand":  "ACCESS_POINT_FORCED_OPEN",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_ACCESS_CONTROL_D5,
				salCommand:  "ACCESS_POINT_FORCED_OPEN",
				numElements: 0,
			},
		},
		{
			name: "CLOCK_AND_TIMEKEEPING",
			args: args{match: map[string]string{
				"application": "CLOCK_AND_TIMEKEEPING",
				"salCommand":  "REQUEST_REFRESH",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_CLOCK_AND_TIMEKEEPING_DF,
				salCommand:  "REQUEST_REFRESH",
				numElements: 0,
			},
		},
		{
			name: "TELEPHONY_STATUS_AND_CONTROL",
			args: args{match: map[string]string{
				"application": "TELEPHONY_STATUS_AND_CONTROL",
				"salCommand":  "EVENT",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_TELEPHONY_STATUS_AND_CONTROL_E0,
				salCommand:  "EVENT",
				numElements: 255,
			},
		},
		{
			name: "MEASUREMENT",
			args: args{match: map[string]string{
				"application": "MEASUREMENT",
				"salCommand":  "MEASUREMENT_EVENT",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_MEASUREMENT_E4,
				salCommand:  "MEASUREMENT_EVENT",
				numElements: 6,
			},
		},
		// TODO: implement ApplicationId_TESTING
		{
			name: "MEDIA_TRANSPORT_CONTROL",
			args: args{match: map[string]string{
				"application": "MEDIA_TRANSPORT_CONTROL",
				"salCommand":  "NEXT_PREVIOUS_SELECTION",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_MEDIA_TRANSPORT_CONTROL_C0,
				salCommand:  "NEXT_PREVIOUS_SELECTION",
				numElements: 1,
			},
		},
		{
			name: "ERROR_REPORTING",
			args: args{match: map[string]string{
				"application": "ERROR_REPORTING",
				"salCommand":  "CLEAR_MOST_SEVERE",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_ERROR_REPORTING_CE,
				salCommand:  "CLEAR_MOST_SEVERE",
				numElements: 8,
			},
		},
		{
			name: "HVAC_ACTUATOR",
			args: args{match: map[string]string{
				"application": "HVAC_ACTUATOR",
				"salCommand":  "TERMINATE_RAMP",
			}},
			want: &salTag{
				tagType:     SAL,
				application: readWriteModel.ApplicationIdContainer_HVAC_ACTUATOR_73,
				salCommand:  "TERMINATE_RAMP",
				numElements: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTagHandler()
			got, err := m.handleSALPattern(tt.args.match)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleSALPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.Equal(t, tt.want, got) {
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
		{
			name:    "empty",
			wantErr: true,
		},
		{
			name: "only unit address",
			args: args{match: map[string]string{
				"unitAddress": "2",
			}},
			wantErr: true,
		},
		{
			name: "identify all",
			args: args{match: map[string]string{
				"unitAddress":       "2",
				"identifyAttribute": "*",
			}},
			want: &unitInfoQuery{
				tagType:     UNIT_INFO,
				unitAddress: readWriteModel.NewUnitAddress(2),
				attribute:   nil,
				numElements: 1,
			},
		},
		{
			name: "identify type",
			args: args{match: map[string]string{
				"unitAddress":       "2",
				"identifyAttribute": "Type",
			}},
			want: &unitInfoQuery{
				tagType:     UNIT_INFO,
				unitAddress: readWriteModel.NewUnitAddress(2),
				attribute: func() *readWriteModel.Attribute {
					var attributeType readWriteModel.Attribute
					attributeType = readWriteModel.Attribute_Type
					return &attributeType
				}(),
				numElements: 1,
			},
		},
		{
			name: "identify type hex",
			args: args{match: map[string]string{
				"unitAddress":       "2",
				"identifyAttribute": "0x01",
			}},
			want: &unitInfoQuery{
				tagType:     UNIT_INFO,
				unitAddress: readWriteModel.NewUnitAddress(2),
				attribute: func() *readWriteModel.Attribute {
					var attributeType readWriteModel.Attribute
					attributeType = readWriteModel.Attribute_Type
					return &attributeType
				}(),
				numElements: 1,
			},
		},
		{
			name: "identify type decimal",
			args: args{match: map[string]string{
				"unitAddress":       "2",
				"identifyAttribute": "1",
			}},
			want: &unitInfoQuery{
				tagType:     UNIT_INFO,
				unitAddress: readWriteModel.NewUnitAddress(2),
				attribute: func() *readWriteModel.Attribute {
					var attributeType readWriteModel.Attribute
					attributeType = readWriteModel.Attribute_Type
					return &attributeType
				}(),
				numElements: 1,
			},
		},
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
			if !assert.Equal(t, tt.want, got) {
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
		{
			name: "empty",
			want: "STATUS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTagHandler_unitAddressFromArgument(t *testing.T) {
	type args struct {
		unitAddressArgument string
		allowWildcard       bool
	}
	tests := []struct {
		name    string
		args    args
		want    readWriteModel.UnitAddress
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "empty",
			wantErr: assert.Error,
		},
		{
			name: "wrong address",
			args: args{
				unitAddressArgument: "adsfadsf",
				allowWildcard:       false,
			},
			wantErr: assert.Error,
		},
		{
			name: "address decimal",
			args: args{
				unitAddressArgument: "1",
				allowWildcard:       false,
			},
			want:    readWriteModel.NewUnitAddress(1),
			wantErr: assert.NoError,
		},
		{
			name: "address hex",
			args: args{
				unitAddressArgument: "0x01",
				allowWildcard:       false,
			},
			want:    readWriteModel.NewUnitAddress(1),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewTagHandler()
			got, err := m.unitAddressFromArgument(tt.args.unitAddressArgument, tt.args.allowWildcard)
			if !tt.wantErr(t, err, fmt.Sprintf("unitAddressFromArgument(%v, %v)", tt.args.unitAddressArgument, tt.args.allowWildcard)) {
				return
			}
			assert.Equalf(t, tt.want, got, "unitAddressFromArgument(%v, %v)", tt.args.unitAddressArgument, tt.args.allowWildcard)
		})
	}
}
