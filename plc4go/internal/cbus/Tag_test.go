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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		want SALMonitorTag
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
