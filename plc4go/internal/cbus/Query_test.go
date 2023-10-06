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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/cbus/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUnitInfoQuery(t *testing.T) {
	type args struct {
		unitAddress readWriteModel.UnitAddress
		attribute   *readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name string
		args args
		want UnitInfoQuery
	}{
		{
			name: "empty",
			want: &unitInfoQuery{
				tagType:     0xc,
				unitAddress: nil,
				attribute:   nil,
				numElements: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewUnitInfoQuery(tt.args.unitAddress, tt.args.attribute, tt.args.numElements), "NewUnitInfoQuery(%v, %v, %v)", tt.args.unitAddress, tt.args.attribute, tt.args.numElements)
		})
	}
}

func Test_unitInfoQuery_GetArrayInfo(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		attribute   *readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []apiModel.ArrayInfo
	}{
		{
			name: "just get",
			fields: fields{
				numElements: 1,
			},
			want: []apiModel.ArrayInfo{},
		},
		{
			name: "just get",
			fields: fields{
				numElements: 2,
			},
			want: []apiModel.ArrayInfo{
				&spiModel.DefaultArrayInfo{
					LowerBound: 0,
					UpperBound: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := unitInfoQuery{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, u.GetArrayInfo(), "GetArrayInfo()")
		})
	}
}

func Test_unitInfoQuery_GetAttribute(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		attribute   *readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   *readWriteModel.Attribute
	}{
		{
			name: "just get",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := unitInfoQuery{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, u.GetAttribute(), "GetAttribute()")
		})
	}
}

func Test_unitInfoQuery_GetQueryString(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		attribute   *readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "just get",
			want: "cal/*/identify=*",
		},
		{
			name: "get with unit",
			fields: fields{
				unitAddress: readWriteModel.NewUnitAddress(13),
			},
			want: "cal/13/identify=*",
		},
		{
			name: "get with unit and attribute",
			fields: fields{
				unitAddress: readWriteModel.NewUnitAddress(13),
				attribute: func() *readWriteModel.Attribute {
					var attr = readWriteModel.Attribute_MaximumLevels
					return &attr
				}(),
			},
			want: "cal/13/identify=MaximumLevels",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := unitInfoQuery{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, u.GetQueryString(), "GetQueryString()")
		})
	}
}

func Test_unitInfoQuery_GetTagType(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		attribute   *readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   TagType
	}{
		{
			name: "just get",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := unitInfoQuery{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, u.GetTagType(), "GetTagType()")
		})
	}
}

func Test_unitInfoQuery_GetUnitAddress(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		attribute   *readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   readWriteModel.UnitAddress
	}{
		{
			name: "just get",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := unitInfoQuery{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, u.GetUnitAddress(), "GetUnitAddress()")
		})
	}
}

func Test_unitInfoQuery_GetValueType(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		attribute   *readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   values.PlcValueType
	}{
		{
			name: "just get",
			want: values.Struct,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := unitInfoQuery{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, u.GetValueType(), "GetValueType()")
		})
	}
}

func Test_unitInfoQuery_Serialize(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		attribute   *readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "just serialize",
			fields: fields{
				tagType: UNIT_INFO,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := unitInfoQuery{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			got, err := u.Serialize()
			if !tt.wantErr(t, err, fmt.Sprintf("Serialize()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Serialize()")
		})
	}
}

func Test_unitInfoQuery_SerializeWithWriteBuffer(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		attribute   *readWriteModel.Attribute
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
			name: "just get",
			fields: fields{
				tagType: UNIT_INFO,
			},
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewJsonWriteBuffer(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "just get with unit and attribute",
			fields: fields{
				tagType:     UNIT_INFO,
				unitAddress: readWriteModel.NewUnitAddress(13),
				attribute: func() *readWriteModel.Attribute {
					var attr = readWriteModel.Attribute_MaximumLevels
					return &attr
				}(),
			},
			args: args{
				ctx:         context.Background(),
				writeBuffer: utils.NewJsonWriteBuffer(),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := unitInfoQuery{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			tt.wantErr(t, u.SerializeWithWriteBuffer(tt.args.ctx, tt.args.writeBuffer), fmt.Sprintf("SerializeWithWriteBuffer(%v, %v)", tt.args.ctx, tt.args.writeBuffer))
		})
	}
}

func Test_unitInfoQuery_String(t *testing.T) {
	type fields struct {
		tagType     TagType
		unitAddress readWriteModel.UnitAddress
		attribute   *readWriteModel.Attribute
		numElements uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "just get",
			want: "<nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := unitInfoQuery{
				tagType:     tt.fields.tagType,
				unitAddress: tt.fields.unitAddress,
				attribute:   tt.fields.attribute,
				numElements: tt.fields.numElements,
			}
			assert.Equalf(t, tt.want, u.String(), "String()")
		})
	}
}
