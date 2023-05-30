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

package values

import (
	"fmt"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/stretchr/testify/assert"
)

func TestDefaultValueHandler_NewPlcValue(t *testing.T) {
	type args struct {
		tag   apiModel.PlcTag
		value any
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func(t *testing.T, args *args)
		want      apiValues.PlcValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "simple bool",
			args: args{
				value: true,
			},
			mockSetup: func(t *testing.T, args *args) {
				tag := NewMockPlcTag(t)
				expect := tag.EXPECT()
				expect.GetArrayInfo().Return(nil)
				expect.GetValueType().Return(apiValues.BOOL)
				args.tag = tag
			},
			want:    NewPlcBOOL(true),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.args)
			}
			m := DefaultValueHandler{}
			got, err := m.NewPlcValue(tt.args.tag, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("NewPlcValue(%v, %v)", tt.args.tag, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewPlcValue(%v, %v)", tt.args.tag, tt.args.value)
		})
	}
}

func TestDefaultValueHandler_NewPlcValueFromType(t *testing.T) {
	type args struct {
		valueType apiValues.PlcValueType
		value     any
	}
	tests := []struct {
		name    string
		args    args
		want    apiValues.PlcValue
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "null in, error out",
			wantErr: assert.Error,
		},
		{
			name: "value in, value out",
			args: args{
				value: NewPlcBOOL(true),
			},
			want:    NewPlcBOOL(true),
			wantErr: assert.NoError,
		},
		{
			name: "bool string",
			args: args{
				valueType: apiValues.BOOL,
				value:     "true",
			},
			want:    NewPlcBOOL(true),
			wantErr: assert.NoError,
		},
		{
			name: "bool",
			args: args{
				valueType: apiValues.BOOL,
				value:     true,
			},
			want:    NewPlcBOOL(true),
			wantErr: assert.NoError,
		},
		{
			name: "bool wrong",
			args: args{
				valueType: apiValues.BOOL,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "byte string",
			args: args{
				valueType: apiValues.BYTE,
				value:     "1",
			},
			want:    NewPlcBYTE(1),
			wantErr: assert.NoError,
		},
		{
			name: "byte",
			args: args{
				valueType: apiValues.BYTE,
				value:     uint8(1),
			},
			want:    NewPlcBYTE(1),
			wantErr: assert.NoError,
		},
		{
			name: "byte wrong",
			args: args{
				valueType: apiValues.BYTE,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "word string",
			args: args{
				valueType: apiValues.WORD,
				value:     "1",
			},
			want:    NewPlcWORD(1),
			wantErr: assert.NoError,
		},
		{
			name: "word",
			args: args{
				valueType: apiValues.WORD,
				value:     uint16(1),
			},
			want:    NewPlcWORD(1),
			wantErr: assert.NoError,
		},
		{
			name: "word wrong",
			args: args{
				valueType: apiValues.WORD,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "dword string",
			args: args{
				valueType: apiValues.DWORD,
				value:     "1",
			},
			want:    NewPlcDWORD(1),
			wantErr: assert.NoError,
		},
		{
			name: "dword",
			args: args{
				valueType: apiValues.DWORD,
				value:     uint32(1),
			},
			want:    NewPlcDWORD(1),
			wantErr: assert.NoError,
		},
		{
			name: "dword wrong",
			args: args{
				valueType: apiValues.DWORD,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "lword string",
			args: args{
				valueType: apiValues.LWORD,
				value:     "1",
			},
			want:    NewPlcLWORD(1),
			wantErr: assert.NoError,
		},
		{
			name: "lword",
			args: args{
				valueType: apiValues.LWORD,
				value:     uint64(1),
			},
			want:    NewPlcLWORD(1),
			wantErr: assert.NoError,
		},
		{
			name: "lword wrong",
			args: args{
				valueType: apiValues.LWORD,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "usint string",
			args: args{
				valueType: apiValues.USINT,
				value:     "1",
			},
			want:    NewPlcUSINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "usint",
			args: args{
				valueType: apiValues.USINT,
				value:     uint8(1),
			},
			want:    NewPlcUSINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "usint wrong",
			args: args{
				valueType: apiValues.USINT,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "uint string",
			args: args{
				valueType: apiValues.UINT,
				value:     "1",
			},
			want:    NewPlcUINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "uint",
			args: args{
				valueType: apiValues.UINT,
				value:     uint16(1),
			},
			want:    NewPlcUINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "uint wrong",
			args: args{
				valueType: apiValues.UINT,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "udint string",
			args: args{
				valueType: apiValues.UDINT,
				value:     "1",
			},
			want:    NewPlcUDINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "udint",
			args: args{
				valueType: apiValues.UDINT,
				value:     uint32(1),
			},
			want:    NewPlcUDINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "udint wrong",
			args: args{
				valueType: apiValues.UDINT,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "ulint string",
			args: args{
				valueType: apiValues.ULINT,
				value:     "1",
			},
			want:    NewPlcULINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "ulint",
			args: args{
				valueType: apiValues.ULINT,
				value:     uint64(1),
			},
			want:    NewPlcULINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "ulint wrong",
			args: args{
				valueType: apiValues.ULINT,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "sint string",
			args: args{
				valueType: apiValues.SINT,
				value:     "1",
			},
			want:    NewPlcSINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "sint",
			args: args{
				valueType: apiValues.SINT,
				value:     int8(1),
			},
			want:    NewPlcSINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "sint wrong",
			args: args{
				valueType: apiValues.SINT,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "int string",
			args: args{
				valueType: apiValues.INT,
				value:     "1",
			},
			want:    NewPlcINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "int",
			args: args{
				valueType: apiValues.INT,
				value:     int16(1),
			},
			want:    NewPlcINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "int wrong",
			args: args{
				valueType: apiValues.INT,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "dint string",
			args: args{
				valueType: apiValues.DINT,
				value:     "1",
			},
			want:    NewPlcDINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "dint",
			args: args{
				valueType: apiValues.DINT,
				value:     int32(1),
			},
			want:    NewPlcDINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "dint wrong",
			args: args{
				valueType: apiValues.DINT,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "lint string",
			args: args{
				valueType: apiValues.LINT,
				value:     "1",
			},
			want:    NewPlcLINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "lint",
			args: args{
				valueType: apiValues.LINT,
				value:     int64(1),
			},
			want:    NewPlcLINT(1),
			wantErr: assert.NoError,
		},
		{
			name: "lint wrong",
			args: args{
				valueType: apiValues.LINT,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "real string",
			args: args{
				valueType: apiValues.REAL,
				value:     "1",
			},
			want:    NewPlcREAL(1),
			wantErr: assert.NoError,
		},
		{
			name: "real",
			args: args{
				valueType: apiValues.REAL,
				value:     float32(1),
			},
			want:    NewPlcREAL(1),
			wantErr: assert.NoError,
		},
		{
			name: "real wrong",
			args: args{
				valueType: apiValues.REAL,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "lreal string",
			args: args{
				valueType: apiValues.LREAL,
				value:     "1",
			},
			want:    NewPlcLREAL(1),
			wantErr: assert.NoError,
		},
		{
			name: "lreal",
			args: args{
				valueType: apiValues.LREAL,
				value:     float64(1),
			},
			want:    NewPlcLREAL(1),
			wantErr: assert.NoError,
		},
		{
			name: "lreal wrong",
			args: args{
				valueType: apiValues.LREAL,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "time string",
			args: args{
				valueType: apiValues.TIME,
				value:     "1",
			},
			wantErr: assert.Error,
		},
		{
			name: "time",
			args: args{
				valueType: apiValues.TIME,
				value:     time.Duration(0),
			},
			want:    NewPlcTIME(time.Duration(0)),
			wantErr: assert.NoError,
		},
		{
			name: "time wrong",
			args: args{
				valueType: apiValues.TIME,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "date string",
			args: args{
				valueType: apiValues.DATE,
				value:     "1",
			},
			wantErr: assert.Error,
		},
		{
			name: "date",
			args: args{
				valueType: apiValues.DATE,
				value:     time.Time{},
			},
			want:    NewPlcDATE(time.Time{}),
			wantErr: assert.NoError,
		},
		{
			name: "date wrong",
			args: args{
				valueType: apiValues.DATE,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "time of day string",
			args: args{
				valueType: apiValues.TIME_OF_DAY,
				value:     "1",
			},
			wantErr: assert.Error,
		},
		{
			name: "time of day",
			args: args{
				valueType: apiValues.TIME_OF_DAY,
				value:     time.Time{},
			},
			want:    NewPlcTIME_OF_DAY(time.Time{}),
			wantErr: assert.NoError,
		},
		{
			name: "time of day  wrong",
			args: args{
				valueType: apiValues.TIME_OF_DAY,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "date and time string",
			args: args{
				valueType: apiValues.DATE_AND_TIME,
				value:     "1",
			},
			wantErr: assert.Error,
		},
		{
			name: "date and time",
			args: args{
				valueType: apiValues.DATE_AND_TIME,
				value:     time.Time{},
			},
			want:    NewPlcDATE_AND_TIME(time.Time{}),
			wantErr: assert.NoError,
		},
		{
			name: "date and time wrong",
			args: args{
				valueType: apiValues.DATE_AND_TIME,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "char wrong",
			args: args{
				valueType: apiValues.CHAR,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "char too much",
			args: args{
				valueType: apiValues.CHAR,
				value:     "12",
			},
			wantErr: assert.Error,
		},
		{
			name: "char",
			args: args{
				valueType: apiValues.CHAR,
				value:     "1",
			},
			want:    NewPlcCHAR("1"),
			wantErr: assert.NoError,
		},
		{
			name: "wchar wrong",
			args: args{
				valueType: apiValues.WCHAR,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "wchar too much",
			args: args{
				valueType: apiValues.WCHAR,
				value:     "12",
			},
			wantErr: assert.Error,
		},
		{
			name: "wchar",
			args: args{
				valueType: apiValues.WCHAR,
				value:     "1",
			},
			want:    NewPlcWCHAR("1"),
			wantErr: assert.NoError,
		},
		{
			name: "string wrong",
			args: args{
				valueType: apiValues.STRING,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "string",
			args: args{
				valueType: apiValues.STRING,
				value:     "1",
			},
			want:    NewPlcSTRING("1"),
			wantErr: assert.NoError,
		},
		{
			name: "wstring wrong",
			args: args{
				valueType: apiValues.WSTRING,
				value:     1,
			},
			wantErr: assert.Error,
		},
		{
			name: "wstring",
			args: args{
				valueType: apiValues.WSTRING,
				value:     "1",
			},
			want:    NewPlcWSTRING("1"),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := DefaultValueHandler{}
			got, err := m.NewPlcValueFromType(tt.args.valueType, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("NewPlcValueFromType(%v, %v)", tt.args.valueType, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "NewPlcValueFromType(%v, %v)", tt.args.valueType, tt.args.value)
		})
	}
}

func TestDefaultValueHandler_ParseListType(t *testing.T) {
	type args struct {
		tag       apiModel.PlcTag
		arrayInfo []apiModel.ArrayInfo
		value     any
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func(t *testing.T, args *args)
		want      apiValues.PlcValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "No array info",
			args: args{
				value: true,
			},
			mockSetup: func(t *testing.T, args *args) {
				tag := NewMockPlcTag(t)
				tag.EXPECT().GetValueType().Return(apiValues.BOOL)
				args.tag = tag
			},
			want:    NewPlcBOOL(true),
			wantErr: assert.NoError,
		},
		{
			name: "no array",
			args: args{
				value: 1,
			},
			mockSetup: func(t *testing.T, args *args) {
				args.arrayInfo = []apiModel.ArrayInfo{
					NewMockArrayInfo(t),
				}
			},
			wantErr: assert.Error,
		},
		{
			name: "bool array wrong size",
			args: args{
				value: []bool{true, true},
			},
			mockSetup: func(t *testing.T, args *args) {
				info := NewMockArrayInfo(t)
				info.EXPECT().GetSize().Return(3)
				args.arrayInfo = []apiModel.ArrayInfo{
					info,
				}
			},
			wantErr: assert.Error,
		},
		{
			name: "bool array",
			args: args{
				value: []bool{true, true},
			},
			mockSetup: func(t *testing.T, args *args) {
				{
					tag := NewMockPlcTag(t)
					tag.EXPECT().GetValueType().Return(apiValues.BOOL)
					args.tag = tag
				}
				{
					info := NewMockArrayInfo(t)
					info.EXPECT().GetSize().Return(2)
					args.arrayInfo = []apiModel.ArrayInfo{
						info,
					}
				}
			},
			want:    NewPlcList([]apiValues.PlcValue{NewPlcBOOL(true), NewPlcBOOL(true)}),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.args)
			}
			m := DefaultValueHandler{}
			got, err := m.ParseListType(tt.args.tag, tt.args.arrayInfo, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseListType(%v, %v, %v)", tt.args.tag, tt.args.arrayInfo, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseListType(%v, %v, %v)", tt.args.tag, tt.args.arrayInfo, tt.args.value)
		})
	}
}

func TestDefaultValueHandler_ParseSimpleType(t *testing.T) {
	type args struct {
		tag   apiModel.PlcTag
		value any
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func(t *testing.T, args *args)
		want      apiValues.PlcValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "fallback",
			args: args{
				value: 1,
			},
			mockSetup: func(t *testing.T, args *args) {
				tag := NewMockPlcTag(t)
				tag.EXPECT().GetValueType().Return(apiValues.BOOL)
				args.tag = tag
			},
			want:    NewPlcBOOL(true),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.args)
			}
			m := DefaultValueHandler{}
			got, err := m.ParseSimpleType(tt.args.tag, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseSimpleType(%v, %v)", tt.args.tag, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseSimpleType(%v, %v)", tt.args.tag, tt.args.value)
		})
	}
}

func TestDefaultValueHandler_ParseStructType(t *testing.T) {
	type args struct {
		in0 apiModel.PlcTag
		in1 any
	}
	tests := []struct {
		name    string
		args    args
		want    apiValues.PlcValue
		wantErr assert.ErrorAssertionFunc
	}{
		{
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := DefaultValueHandler{}
			got, err := m.ParseStructType(tt.args.in0, tt.args.in1)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseStructType(%v, %v)", tt.args.in0, tt.args.in1)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseStructType(%v, %v)", tt.args.in0, tt.args.in1)
		})
	}
}

func TestDefaultValueHandler_parseType(t *testing.T) {
	type args struct {
		tag       apiModel.PlcTag
		arrayInfo []apiModel.ArrayInfo
		value     any
	}
	tests := []struct {
		name      string
		args      args
		mockSetup func(t *testing.T, args *args)
		want      apiValues.PlcValue
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "parse list",
			args: args{
				value: []bool{true, true},
			},
			mockSetup: func(t *testing.T, args *args) {
				{
					tag := NewMockPlcTag(t)
					tag.EXPECT().GetValueType().Return(apiValues.BOOL)
					args.tag = tag
				}
				{
					info := NewMockArrayInfo(t)
					info.EXPECT().GetSize().Return(2)
					args.arrayInfo = []apiModel.ArrayInfo{
						info,
					}
				}
			},
			want:    NewPlcList([]apiValues.PlcValue{NewPlcBOOL(true), NewPlcBOOL(true)}),
			wantErr: assert.NoError,
		},
		{
			name: "parse struct",
			args: args{
				value: true,
			},
			mockSetup: func(t *testing.T, args *args) {
				tag := NewMockPlcTag(t)
				tag.EXPECT().GetValueType().Return(apiValues.Struct)
				args.tag = tag
			},
			wantErr: assert.Error,
		},
		{
			name: "parse simple",
			args: args{
				value: true,
			},
			mockSetup: func(t *testing.T, args *args) {
				tag := NewMockPlcTag(t)
				tag.EXPECT().GetValueType().Return(apiValues.BOOL)
				args.tag = tag
			},
			want:    NewPlcBOOL(true),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.args)
			}
			m := DefaultValueHandler{}
			got, err := m.parseType(tt.args.tag, tt.args.arrayInfo, tt.args.value)
			if !tt.wantErr(t, err, fmt.Sprintf("parseType(%v, %v, %v)", tt.args.tag, tt.args.arrayInfo, tt.args.value)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseType(%v, %v, %v)", tt.args.tag, tt.args.arrayInfo, tt.args.value)
		})
	}
}
