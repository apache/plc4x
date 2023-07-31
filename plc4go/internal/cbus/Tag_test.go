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
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

func TestNewCALGetStatusTag(t *testing.T) {
	type args struct {
		unitAddress     readWriteModel.UnitAddress
		bridgeAddresses []readWriteModel.BridgeAddress
		parameter       readWriteModel.Parameter
		count           uint8
		numElements     uint16
	}
	tests := []struct {
		name string
		args args
		want CALGetStatusTag
	}{
		{
			name: "just create it",
			want: &calGetStatusTag{
				tagType: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewCALGetStatusTag(tt.args.unitAddress, tt.args.bridgeAddresses, tt.args.parameter, tt.args.count, tt.args.numElements), "NewCALGetStatusTag(%v, %v, %v, %v, %v)", tt.args.unitAddress, tt.args.bridgeAddresses, tt.args.parameter, tt.args.count, tt.args.numElements)
		})
	}
}

func TestNewCALIdentifyTag(t *testing.T) {
	type args struct {
		unitAddress     readWriteModel.UnitAddress
		bridgeAddresses []readWriteModel.BridgeAddress
		attribute       readWriteModel.Attribute
		numElements     uint16
	}
	tests := []struct {
		name string
		args args
		want CALIdentifyTag
	}{
		{
			name: "just create it",
			want: &calIdentifyTag{
				tagType: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewCALIdentifyTag(tt.args.unitAddress, tt.args.bridgeAddresses, tt.args.attribute, tt.args.numElements), "NewCALIdentifyTag(%v, %v, %v, %v)", tt.args.unitAddress, tt.args.bridgeAddresses, tt.args.attribute, tt.args.numElements)
		})
	}
}

func TestNewCALRecallTag(t *testing.T) {
	type args struct {
		unitAddress     readWriteModel.UnitAddress
		bridgeAddresses []readWriteModel.BridgeAddress
		parameter       readWriteModel.Parameter
		count           uint8
		numElements     uint16
	}
	tests := []struct {
		name string
		args args
		want CALRecallTag
	}{
		{
			name: "just create it",
			want: &calRecallTag{
				tagType: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewCALRecallTag(tt.args.unitAddress, tt.args.bridgeAddresses, tt.args.parameter, tt.args.count, tt.args.numElements), "NewCALRecallTag(%v, %v, %v, %v, %v)", tt.args.unitAddress, tt.args.bridgeAddresses, tt.args.parameter, tt.args.count, tt.args.numElements)
		})
	}
}

func TestNewMMIMonitorTag(t *testing.T) {
	type args struct {
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name string
		args args
		want MMIMonitorTag
	}{
		{
			name: "just create it",
			want: &mmiMonitorTag{
				tagType: 11,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewMMIMonitorTag(tt.args.unitAddress, tt.args.application, tt.args.numElements), "NewMMIMonitorTag(%v, %v, %v)", tt.args.unitAddress, tt.args.application, tt.args.numElements)
		})
	}
}

func TestNewSALMonitorTag(t *testing.T) {
	type args struct {
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name string
		args args
		want SALMonitorTag
	}{
		{
			name: "just create it",
			want: &salMonitorTag{
				tagType: 0xa,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewSALMonitorTag(tt.args.unitAddress, tt.args.application, tt.args.numElements), "NewSALMonitorTag(%v, %v, %v)", tt.args.unitAddress, tt.args.application, tt.args.numElements)
		})
	}
}

func TestNewSALTag(t *testing.T) {
	type args struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	tests := []struct {
		name string
		args args
		want SALTag
	}{
		{
			name: "just create it",
			want: &salTag{
				tagType: 9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewSALTag(tt.args.bridgeAddresses, tt.args.application, tt.args.salCommand, tt.args.numElements), "NewSALTag(%v, %v, %v, %v)", tt.args.bridgeAddresses, tt.args.application, tt.args.salCommand, tt.args.numElements)
		})
	}
}

func TestNewStatusTag(t *testing.T) {
	type args struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name string
		args args
		want StatusTag
	}{
		{
			name: "just create it",
			want: &statusTag{
				tagType: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewStatusTag(tt.args.bridgeAddresses, tt.args.statusRequestType, tt.args.startingGroupAddressLabel, tt.args.application, tt.args.numElements), "NewStatusTag(%v, %v, %v, %v, %v)", tt.args.bridgeAddresses, tt.args.statusRequestType, tt.args.startingGroupAddressLabel, tt.args.application, tt.args.numElements)
		})
	}
}

func TestStatusRequestType_String(t *testing.T) {
	tests := []struct {
		name string
		s    StatusRequestType
		want string
	}{
		{
			name: "get a string",
			want: "StatusRequestTypeBinaryState",
		},
		{
			name: "get a string",
			s:    StatusRequestTypeLevel,
			want: "StatusRequestTypeLevel",
		},
		{
			name: "non type",
			s:    255,
			want: "StatusRequestType(255)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.s.String(), "String()")
		})
	}
}

func Test_calGetStatusTag_GetAddressString(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get a string",
			want: "cal/getstatus=UNKNOWN_01, 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calGetStatusTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetAddressString(), "GetAddressString()")
		})
	}
}

func Test_calGetStatusTag_GetArrayInfo(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []apiModel.ArrayInfo
	}{
		{
			name: "get array info",
			want: []apiModel.ArrayInfo{
				&spiModel.DefaultArrayInfo{},
			},
		},
		{
			name: "one element",
			fields: fields{
				count: 1,
			},
			want: []apiModel.ArrayInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calGetStatusTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetArrayInfo(), "GetArrayInfo()")
		})
	}
}

func Test_calGetStatusTag_GetCount(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "get a count",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calGetStatusTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetCount(), "GetCount()")
		})
	}
}

func Test_calGetStatusTag_GetParameter(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   readWriteModel.Parameter
	}{
		{
			name: "get a parameter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calGetStatusTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetParameter(), "GetParameter()")
		})
	}
}

func Test_calGetStatusTag_GetTagType(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   TagType
	}{
		{
			name: "get the type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calGetStatusTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetTagType(), "GetTagType()")
		})
	}
}

func Test_calGetStatusTag_GetValueType(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   apiValues.PlcValueType
	}{
		{
			name: "get the value type",
			want: 0x61,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calGetStatusTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetValueType(), "GetValueType()")
		})
	}
}

func Test_calGetStatusTag_Serialize(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "get the value type",
			want:    []byte{0, 0},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calGetStatusTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			got, err := c.Serialize()
			if !tt.wantErr(t, err, fmt.Sprintf("Serialize()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Serialize()")
		})
	}
}

func Test_calGetStatusTag_SerializeWithWriteBuffer(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	type args struct {
		ctx         context.Context
		writeBuffer utils.WriteBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "serialize empty",
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calGetStatusTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			tt.wantErr(t, c.SerializeWithWriteBuffer(tt.args.ctx, tt.args.writeBuffer), fmt.Sprintf("SerializeWithWriteBuffer(%v, %v)", tt.args.ctx, tt.args.writeBuffer))
		})
	}
}

func Test_calGetStatusTag_String(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "default to string",
			want: `
╔═STATUS════════════════════╗
║╔═Parameter═══════╗╔═count╗║
║║0x00 0 UNKNOWN_01║║0x00 0║║
║╚═════════════════╝╚══════╝║
╚═══════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calGetStatusTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.String(), "String()")
		})
	}
}

func Test_calIdentifyTag_GetAddressString(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		attribute   readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get the address string",
			fields: fields{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(2),
				},
			},
			want: "cal/2/identify=Manufacturer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calIdentifyTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetAddressString(), "GetAddressString()")
		})
	}
}

func Test_calIdentifyTag_GetArrayInfo(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		attribute   readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []apiModel.ArrayInfo
	}{
		{
			name: "get empty array info",
			want: []apiModel.ArrayInfo{
				&spiModel.DefaultArrayInfo{},
			},
		},
		{
			name: "one element",
			fields: fields{
				numElements: 1,
			},
			want: []apiModel.ArrayInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calIdentifyTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetArrayInfo(), "GetArrayInfo()")
		})
	}
}

func Test_calIdentifyTag_GetAttribute(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		attribute   readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   readWriteModel.Attribute
	}{
		{
			name: "get the attribute",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calIdentifyTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetAttribute(), "GetAttribute()")
		})
	}
}

func Test_calIdentifyTag_GetTagType(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		attribute   readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   TagType
	}{
		{
			name: "get the tag type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calIdentifyTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetTagType(), "GetTagType()")
		})
	}
}

func Test_calIdentifyTag_GetValueType(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		attribute   readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   apiValues.PlcValueType
	}{
		{
			name: "get the value type",
			want: 0x61,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calIdentifyTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetValueType(), "GetValueType()")
		})
	}
}

func Test_calIdentifyTag_Serialize(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		attribute   readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "serialize empty",
			want:    []byte{0},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calIdentifyTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			got, err := c.Serialize()
			if !tt.wantErr(t, err, fmt.Sprintf("Serialize()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Serialize()")
		})
	}
}

func Test_calIdentifyTag_SerializeWithWriteBuffer(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		attribute   readWriteModel.Attribute
		numElements uint16
	}
	type args struct {
		ctx         context.Context
		writeBuffer utils.WriteBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "serialize with write buffer",
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calIdentifyTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			tt.wantErr(t, c.SerializeWithWriteBuffer(tt.args.ctx, tt.args.writeBuffer), fmt.Sprintf("SerializeWithWriteBuffer(%v, %v)", tt.args.ctx, tt.args.writeBuffer))
		})
	}
}

func Test_calIdentifyTag_String(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		attribute   readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "default to string",
			want: `
╔═STATUS/Attribute╗
║0x00 0 Manufacturer║
╚═════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calIdentifyTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.String(), "String()")
		})
	}
}

func Test_calRecallTag_GetAddressString(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get the address string",
			fields: fields{
				calTag: calTag{
					unitAddress: readWriteModel.NewUnitAddress(24),
				},
			},
			want: "cal/24/recall=UNKNOWN_01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calRecallTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetAddressString(), "GetAddressString()")
		})
	}
}

func Test_calRecallTag_GetArrayInfo(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []apiModel.ArrayInfo
	}{
		{
			name: "empty array info",
			want: []apiModel.ArrayInfo{
				&spiModel.DefaultArrayInfo{},
			},
		},
		{
			name: "one element",
			fields: fields{
				count: 1,
			},
			want: []apiModel.ArrayInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calRecallTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetArrayInfo(), "GetArrayInfo()")
		})
	}
}

func Test_calRecallTag_GetCount(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "get the count",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calRecallTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetCount(), "GetCount()")
		})
	}
}

func Test_calRecallTag_GetParameter(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   readWriteModel.Parameter
	}{
		{
			name: "get the parameter",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calRecallTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetParameter(), "GetParameter()")
		})
	}
}

func Test_calRecallTag_GetTagType(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   TagType
	}{
		{
			name: "get the tag type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calRecallTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetTagType(), "GetTagType()")
		})
	}
}

func Test_calRecallTag_GetValueType(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   apiValues.PlcValueType
	}{
		{
			name: "get the value type",
			want: 0x61,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calRecallTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.GetValueType(), "GetValueType()")
		})
	}
}

func Test_calRecallTag_Serialize(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "serialize empty",
			want:    []byte{0, 0},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calRecallTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			got, err := c.Serialize()
			if !tt.wantErr(t, err, fmt.Sprintf("Serialize()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Serialize()")
		})
	}
}

func Test_calRecallTag_SerializeWithWriteBuffer(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	type args struct {
		ctx         context.Context
		writeBuffer utils.WriteBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "serialize empty",
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calRecallTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			tt.wantErr(t, c.SerializeWithWriteBuffer(tt.args.ctx, tt.args.writeBuffer), fmt.Sprintf("SerializeWithWriteBuffer(%v, %v)", tt.args.ctx, tt.args.writeBuffer))
		})
	}
}

func Test_calRecallTag_String(t *testing.T) {
	type fields struct {
		calTag      calTag
		tagType     TagType
		parameter   readWriteModel.Parameter
		count       uint8
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "default to string",
			want: `
╔═STATUS════════════════════╗
║╔═Parameter═══════╗╔═count╗║
║║0x00 0 UNKNOWN_01║║0x00 0║║
║╚═════════════════╝╚══════╝║
╚═══════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calRecallTag{
				calTag:      tt.fields.calTag,
				tagType:     tt.fields.tagType,
				parameter:   tt.fields.parameter,
				count:       tt.fields.count,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, c.String(), "String()")
		})
	}
}

func Test_calTag_GetBridgeAddresses(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		unitAddress     readWriteModel.UnitAddress
	}
	tests := []struct {
		name   string
		fields fields
		want   []readWriteModel.BridgeAddress
	}{
		{
			name: "get the bridge addresses",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				unitAddress:     tt.fields.unitAddress,
			}
			assert.Equalf(t, tt.want, c.GetBridgeAddresses(), "GetBridgeAddresses()")
		})
	}
}

func Test_calTag_GetUnitAddress(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		unitAddress     readWriteModel.UnitAddress
	}
	tests := []struct {
		name   string
		fields fields
		want   readWriteModel.UnitAddress
	}{
		{
			name: "get the unit address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				unitAddress:     tt.fields.unitAddress,
			}
			assert.Equalf(t, tt.want, c.GetUnitAddress(), "GetUnitAddress()")
		})
	}
}

func Test_calTag_Serialize(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		unitAddress     readWriteModel.UnitAddress
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "serialize empty",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				unitAddress:     tt.fields.unitAddress,
			}
			got, err := c.Serialize()
			if !tt.wantErr(t, err, fmt.Sprintf("Serialize()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Serialize()")
		})
	}
}

func Test_calTag_SerializeWithWriteBuffer(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		unitAddress     readWriteModel.UnitAddress
	}
	type args struct {
		ctx         context.Context
		writeBuffer utils.WriteBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "serialize empty",
			wantErr: assert.NoError,
		},
		{
			name: "serialize with bridges",
			fields: fields{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewBridgeAddress(2),
					readWriteModel.NewBridgeAddress(3),
					readWriteModel.NewBridgeAddress(4),
					readWriteModel.NewBridgeAddress(5),
					readWriteModel.NewBridgeAddress(6),
				},
				unitAddress: readWriteModel.NewUnitAddress(34),
			},
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				unitAddress:     tt.fields.unitAddress,
			}
			tt.wantErr(t, c.SerializeWithWriteBuffer(tt.args.ctx, tt.args.writeBuffer), fmt.Sprintf("SerializeWithWriteBuffer(%v, %v)", tt.args.ctx, tt.args.writeBuffer))
		})
	}
}

func Test_calTag_String(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		unitAddress     readWriteModel.UnitAddress
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "to string",
			want: "<nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := calTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				unitAddress:     tt.fields.unitAddress,
			}
			assert.Equalf(t, tt.want, c.String(), "String()")
		})
	}
}

func Test_mmiMonitorTag_GetAddressString(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get the address string",
			want: "mmimonitor/*/*",
		},
		{
			name: "get addressed sting",
			fields: fields{
				unitAddress: readWriteModel.NewUnitAddress(23),
				application: func() *readWriteModel.ApplicationIdContainer {
					a := readWriteModel.ApplicationIdContainer_HEATING_88
					return &a
				}(),
			},
			want: "mmimonitor/0x17/HEATING_88",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mmiMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, m.GetAddressString(), "GetAddressString()")
		})
	}
}

func Test_mmiMonitorTag_GetApplication(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   *readWriteModel.ApplicationIdContainer
	}{
		{
			name: "get application",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mmiMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, m.GetApplication(), "GetApplication()")
		})
	}
}

func Test_mmiMonitorTag_GetArrayInfo(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []apiModel.ArrayInfo
	}{
		{
			name: "mmi monitor tag",
			want: []apiModel.ArrayInfo{
				&spiModel.DefaultArrayInfo{},
			},
		},
		{
			name: "one element",
			fields: fields{
				numElements: 1,
			},
			want: []apiModel.ArrayInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mmiMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, m.GetArrayInfo(), "GetArrayInfo()")
		})
	}
}

func Test_mmiMonitorTag_GetTagType(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   TagType
	}{
		{
			name: "get tag type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mmiMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, m.GetTagType(), "GetTagType()")
		})
	}
}

func Test_mmiMonitorTag_GetUnitAddress(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   readWriteModel.UnitAddress
	}{
		{
			name: "get unit address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mmiMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, m.GetUnitAddress(), "GetUnitAddress()")
		})
	}
}

func Test_mmiMonitorTag_GetValueType(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   apiValues.PlcValueType
	}{
		{
			name: "get value type",
			want: 0x61,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mmiMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, m.GetValueType(), "GetValueType()")
		})
	}
}

func Test_mmiMonitorTag_Serialize(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "serialize empty",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mmiMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			got, err := m.Serialize()
			if !tt.wantErr(t, err, fmt.Sprintf("Serialize()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Serialize()")
		})
	}
}

func Test_mmiMonitorTag_SerializeWithWriteBuffer(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	type args struct {
		ctx         context.Context
		writeBuffer utils.WriteBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "serialize empty",
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "serialize empty",
			fields: fields{
				unitAddress: readWriteModel.NewUnitAddress(23),
				application: func() *readWriteModel.ApplicationIdContainer {
					a := readWriteModel.ApplicationIdContainer_HEATING_88
					return &a
				}(),
			},
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mmiMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			tt.wantErr(t, m.SerializeWithWriteBuffer(tt.args.ctx, tt.args.writeBuffer), fmt.Sprintf("SerializeWithWriteBuffer(%v, %v)", tt.args.ctx, tt.args.writeBuffer))
		})
	}
}

func Test_mmiMonitorTag_String(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "default to string",
			want: "<nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mmiMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, m.String(), "String()")
		})
	}
}

func Test_salMonitorTag_GetAddressString(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get address string",
			want: "salmonitor/*/*",
		},
		{
			name: "get address string with fixed values",
			fields: fields{
				unitAddress: readWriteModel.NewUnitAddress(34),
				application: func() *readWriteModel.ApplicationIdContainer {
					a := readWriteModel.ApplicationIdContainer_HEATING_88
					return &a
				}(),
			},
			want: "salmonitor/0x22/HEATING_88",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetAddressString(), "GetAddressString()")
		})
	}
}

func Test_salMonitorTag_GetApplication(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   *readWriteModel.ApplicationIdContainer
	}{
		{
			name: "get application",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetApplication(), "GetApplication()")
		})
	}
}

func Test_salMonitorTag_GetArrayInfo(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []apiModel.ArrayInfo
	}{
		{
			name: "get empty array info",
			want: []apiModel.ArrayInfo{
				&spiModel.DefaultArrayInfo{},
			},
		},
		{
			name: "one element",
			fields: fields{
				numElements: 1,
			},
			want: []apiModel.ArrayInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetArrayInfo(), "GetArrayInfo()")
		})
	}
}

func Test_salMonitorTag_GetTagType(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   TagType
	}{
		{
			name: "get tag type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetTagType(), "GetTagType()")
		})
	}
}

func Test_salMonitorTag_GetUnitAddress(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   readWriteModel.UnitAddress
	}{
		{
			name: "get unit address",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetUnitAddress(), "GetUnitAddress()")
		})
	}
}

func Test_salMonitorTag_GetValueType(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   apiValues.PlcValueType
	}{
		{
			name: "get value type",
			want: 0x61,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetValueType(), "GetValueType()")
		})
	}
}

func Test_salMonitorTag_Serialize(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "default serialize",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			got, err := s.Serialize()
			if !tt.wantErr(t, err, fmt.Sprintf("Serialize()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Serialize()")
		})
	}
}

func Test_salMonitorTag_SerializeWithWriteBuffer(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	type args struct {
		ctx         context.Context
		writeBuffer utils.WriteBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "default serialize",
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "serialize with content",
			fields: fields{
				unitAddress: readWriteModel.NewUnitAddress(34),
				application: func() *readWriteModel.ApplicationIdContainer {
					a := readWriteModel.ApplicationIdContainer_HEATING_88
					return &a
				}(),
			},
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			tt.wantErr(t, s.SerializeWithWriteBuffer(tt.args.ctx, tt.args.writeBuffer), fmt.Sprintf("SerializeWithWriteBuffer(%v, %v)", tt.args.ctx, tt.args.writeBuffer))
		})
	}
}

func Test_salMonitorTag_String(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		application *readWriteModel.ApplicationIdContainer
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "default string",
			want: "<nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salMonitorTag{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				application: tt.fields.application,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.String(), "String()")
		})
	}
}

func Test_salTag_GetAddressString(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		tagType         TagType
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get address string",
			want: "sal/RESERVED_00/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				tagType:         tt.fields.tagType,
				application:     tt.fields.application,
				salCommand:      tt.fields.salCommand,
				numElements:     tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetAddressString(), "GetAddressString()")
		})
	}
}

func Test_salTag_GetApplication(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		tagType         TagType
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   readWriteModel.ApplicationIdContainer
	}{
		{
			name: "get application",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				tagType:         tt.fields.tagType,
				application:     tt.fields.application,
				salCommand:      tt.fields.salCommand,
				numElements:     tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetApplication(), "GetApplication()")
		})
	}
}

func Test_salTag_GetArrayInfo(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		tagType         TagType
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []apiModel.ArrayInfo
	}{
		{
			name: "get empty array info",
			want: []apiModel.ArrayInfo{
				&spiModel.DefaultArrayInfo{},
			},
		},
		{
			name: "one element",
			fields: fields{
				numElements: 1,
			},
			want: []apiModel.ArrayInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				tagType:         tt.fields.tagType,
				application:     tt.fields.application,
				salCommand:      tt.fields.salCommand,
				numElements:     tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetArrayInfo(), "GetArrayInfo()")
		})
	}
}

func Test_salTag_GetBridgeAddresses(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		tagType         TagType
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []readWriteModel.BridgeAddress
	}{
		{
			name: "get bridge addresses",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				tagType:         tt.fields.tagType,
				application:     tt.fields.application,
				salCommand:      tt.fields.salCommand,
				numElements:     tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetBridgeAddresses(), "GetBridgeAddresses()")
		})
	}
}

func Test_salTag_GetSALCommand(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		tagType         TagType
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get sal command",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				tagType:         tt.fields.tagType,
				application:     tt.fields.application,
				salCommand:      tt.fields.salCommand,
				numElements:     tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetSALCommand(), "GetSALCommand()")
		})
	}
}

func Test_salTag_GetTagType(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		tagType         TagType
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   TagType
	}{
		{
			name: "get tag type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				tagType:         tt.fields.tagType,
				application:     tt.fields.application,
				salCommand:      tt.fields.salCommand,
				numElements:     tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetTagType(), "GetTagType()")
		})
	}
}

func Test_salTag_GetValueType(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		tagType         TagType
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   apiValues.PlcValueType
	}{
		{
			name: "get value type",
			want: 0x61,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				tagType:         tt.fields.tagType,
				application:     tt.fields.application,
				salCommand:      tt.fields.salCommand,
				numElements:     tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetValueType(), "GetValueType()")
		})
	}
}

func Test_salTag_Serialize(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		tagType         TagType
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "serialize empty",
			want:    []byte{0x0},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				tagType:         tt.fields.tagType,
				application:     tt.fields.application,
				salCommand:      tt.fields.salCommand,
				numElements:     tt.fields.numElements,
			}
			got, err := s.Serialize()
			if !tt.wantErr(t, err, fmt.Sprintf("Serialize()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Serialize()")
		})
	}
}

func Test_salTag_SerializeWithWriteBuffer(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		tagType         TagType
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	type args struct {
		ctx         context.Context
		writeBuffer utils.WriteBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "serialize default",
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "serialize with bridges",
			fields: fields{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewBridgeAddress(2),
					readWriteModel.NewBridgeAddress(3),
					readWriteModel.NewBridgeAddress(4),
					readWriteModel.NewBridgeAddress(5),
					readWriteModel.NewBridgeAddress(6),
				},
			},
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				tagType:         tt.fields.tagType,
				application:     tt.fields.application,
				salCommand:      tt.fields.salCommand,
				numElements:     tt.fields.numElements,
			}
			tt.wantErr(t, s.SerializeWithWriteBuffer(tt.args.ctx, tt.args.writeBuffer), fmt.Sprintf("SerializeWithWriteBuffer(%v, %v)", tt.args.ctx, tt.args.writeBuffer))
		})
	}
}

func Test_salTag_String(t *testing.T) {
	type fields struct {
		bridgeAddresses []readWriteModel.BridgeAddress
		tagType         TagType
		application     readWriteModel.ApplicationIdContainer
		salCommand      string
		numElements     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "default to string",
			want: `
╔═STATUS/ApplicationIdContainer═════╗
║        0x00 0 RESERVED_00         ║
╚═══════════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := salTag{
				bridgeAddresses: tt.fields.bridgeAddresses,
				tagType:         tt.fields.tagType,
				application:     tt.fields.application,
				salCommand:      tt.fields.salCommand,
				numElements:     tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.String(), "String()")
		})
	}
}

func Test_statusTag_GetAddressString(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get address string",
			want: "status/binary/RESERVED_00",
		},
		{
			name: "get binary state",
			fields: fields{
				statusRequestType: StatusRequestTypeBinaryState,
				application:       readWriteModel.ApplicationIdContainer_HEATING_88,
			},
			want: "status/binary/HEATING_88",
		},
		{
			name: "get level",
			fields: fields{
				statusRequestType: StatusRequestTypeLevel,
				application:       readWriteModel.ApplicationIdContainer_HEATING_88,
			},
			want: "status/level/HEATING_88",
		},
		{
			name: "get level with label",
			fields: fields{
				statusRequestType: StatusRequestTypeLevel,
				application:       readWriteModel.ApplicationIdContainer_HEATING_88,
				startingGroupAddressLabel: func() *byte {
					label := byte(5)
					return &label
				}(),
			},
			want: "status/level=0x05/HEATING_88",
		},
		{
			name: "invalid",
			fields: fields{
				statusRequestType: 255,
			},
			want: "status/invalid/RESERVED_00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetAddressString(), "GetAddressString()")
		})
	}
}

func Test_statusTag_GetApplication(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   readWriteModel.ApplicationIdContainer
	}{
		{
			name: "get application",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetApplication(), "GetApplication()")
		})
	}
}

func Test_statusTag_GetArrayInfo(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []apiModel.ArrayInfo
	}{
		{
			name: "get empty array info",
			want: []apiModel.ArrayInfo{
				&spiModel.DefaultArrayInfo{},
			},
		},
		{
			name: "one element",
			fields: fields{
				numElements: 1,
			},
			want: []apiModel.ArrayInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetArrayInfo(), "GetArrayInfo()")
		})
	}
}

func Test_statusTag_GetBridgeAddresses(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []readWriteModel.BridgeAddress
	}{
		{
			name: "get bridge addresses",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetBridgeAddresses(), "GetBridgeAddresses()")
		})
	}
}

func Test_statusTag_GetStartingGroupAddressLabel(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   *byte
	}{
		{
			name: "get starting group address label",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetStartingGroupAddressLabel(), "GetStartingGroupAddressLabel()")
		})
	}
}

func Test_statusTag_GetStatusRequestType(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   StatusRequestType
	}{
		{
			name: "get status request type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetStatusRequestType(), "GetStatusRequestType()")
		})
	}
}

func Test_statusTag_GetTagType(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   TagType
	}{
		{
			name: "get tag type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetTagType(), "GetTagType()")
		})
	}
}

func Test_statusTag_GetValueType(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   apiValues.PlcValueType
	}{
		{
			name: "get value type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.GetValueType(), "GetValueType()")
		})
	}
}

func Test_statusTag_Serialize(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "serialize empty",
			want:    []byte{0, 0},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			got, err := s.Serialize()
			if !tt.wantErr(t, err, fmt.Sprintf("Serialize()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Serialize()")
		})
	}
}

func Test_statusTag_SerializeWithWriteBuffer(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	type args struct {
		ctx         context.Context
		writeBuffer utils.WriteBuffer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "serialize empty",
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "Serialize with bridge",
			fields: fields{
				bridgeAddresses: []readWriteModel.BridgeAddress{
					readWriteModel.NewBridgeAddress(1),
					readWriteModel.NewBridgeAddress(2),
					readWriteModel.NewBridgeAddress(3),
					readWriteModel.NewBridgeAddress(4),
					readWriteModel.NewBridgeAddress(5),
					readWriteModel.NewBridgeAddress(6),
				},
				tagType:           0,
				statusRequestType: 0,
				startingGroupAddressLabel: func() *byte {
					label := byte(4)
					return &label
				}(),
				application: 0,
				numElements: 0,
			},
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewWriteBufferByteBased(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			tt.wantErr(t, s.SerializeWithWriteBuffer(tt.args.ctx, tt.args.writeBuffer), fmt.Sprintf("SerializeWithWriteBuffer(%v, %v)", tt.args.ctx, tt.args.writeBuffer))
		})
	}
}

func Test_statusTag_String(t *testing.T) {
	type fields struct {
		bridgeAddresses           []readWriteModel.BridgeAddress
		tagType                   TagType
		statusRequestType         StatusRequestType
		startingGroupAddressLabel *byte
		application               readWriteModel.ApplicationIdContainer
		numElements               uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "default to string",
			want: `
╔═STATUS══════════════════════════════════════════════════╗
║╔═statusRequestType═════════════════╗╔═application══════╗║
║║0x00 0 StatusRequestTypeBinaryState║║0x00 0 RESERVED_00║║
║╚═══════════════════════════════════╝╚══════════════════╝║
╚═════════════════════════════════════════════════════════╝`[1:],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := statusTag{
				bridgeAddresses:           tt.fields.bridgeAddresses,
				tagType:                   tt.fields.tagType,
				statusRequestType:         tt.fields.statusRequestType,
				startingGroupAddressLabel: tt.fields.startingGroupAddressLabel,
				application:               tt.fields.application,
				numElements:               tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, s.String(), "String()")
		})
	}
}
