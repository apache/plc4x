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

package utils

import (
	"fmt"
	"testing"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"

	"github.com/stretchr/testify/assert"
)

func TestPlcValueUint8ListToByteArray(t *testing.T) {
	type args struct {
		value apiValues.PlcValue
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "no input no output",
			want: []byte{},
		},
		{
			name: "no the right plc value",
			args: args{value: func() apiValues.PlcValue {
				value := NewMockPlcValue(t)
				expect := value.EXPECT()
				expect.IsList().Return(false)
				expect.String().Return("false").Maybe()
				return value
			}()},
			want: []byte{},
		},
		{
			name: "the right plc value",
			args: args{value: func() apiValues.PlcValue {
				value := NewMockPlcValue(t)
				expect := value.EXPECT()
				expect.IsList().Return(true)
				listValue := NewMockPlcValue(t)
				listValue.EXPECT().GetUint8().Return(255)
				expect.GetList().Return([]apiValues.PlcValue{listValue, listValue})
				expect.String().Return("false").Maybe()
				return value
			}()},
			want: []byte{0xff, 0xff},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, PlcValueUint8ListToByteArray(tt.args.value), "PlcValueUint8ListToByteArray(%v)", tt.args.value)
		})
	}
}

func TestStrToBool(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "true",
			args:    args{"true"},
			want:    true,
			wantErr: assert.NoError,
		},
		{
			name:    "false",
			args:    args{"false"},
			want:    false,
			wantErr: assert.NoError,
		},
		{
			name:    "herbert",
			args:    args{"herbert"},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrToBool(tt.args.str)
			if !tt.wantErr(t, err, fmt.Sprintf("StrToBool(%v)", tt.args.str)) {
				return
			}
			assert.Equalf(t, tt.want, got, "StrToBool(%v)", tt.args.str)
		})
	}
}

func TestStrToInt32(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    int32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "1",
			args:    args{"1"},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name:    "123456789123456789",
			args:    args{"123456789123456789"},
			wantErr: assert.Error,
		},
		{
			name:    "herbert",
			args:    args{"herbert"},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrToInt32(tt.args.str)
			if !tt.wantErr(t, err, fmt.Sprintf("StrToInt32(%v)", tt.args.str)) {
				return
			}
			assert.Equalf(t, tt.want, got, "StrToInt32(%v)", tt.args.str)
		})
	}
}

func TestStrToString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "1",
			args:    args{"1"},
			want:    "1",
			wantErr: assert.NoError,
		},
		{
			name:    "123456789123456789",
			args:    args{"123456789123456789"},
			want:    "123456789123456789",
			wantErr: assert.NoError,
		},
		{
			name:    "herbert",
			args:    args{"herbert"},
			want:    "herbert",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrToString(tt.args.s)
			if !tt.wantErr(t, err, fmt.Sprintf("StrToString(%v)", tt.args.s)) {
				return
			}
			assert.Equalf(t, tt.want, got, "StrToString(%v)", tt.args.s)
		})
	}
}

func TestStrToUint16(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "1",
			args:    args{"1"},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name:    "123456789123456789",
			args:    args{"123456789123456789"},
			wantErr: assert.Error,
		},
		{
			name:    "-1",
			args:    args{"-1"},
			want:    0xffff,
			wantErr: assert.NoError,
		},
		{
			name:    "herbert",
			args:    args{"herbert"},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrToUint16(tt.args.str)
			if !tt.wantErr(t, err, fmt.Sprintf("StrToUint16(%v)", tt.args.str)) {
				return
			}
			assert.Equalf(t, tt.want, got, "StrToUint16(%v)", tt.args.str)
		})
	}
}

func TestStrToUint32(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    uint32
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "1",
			args:    args{"1"},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name:    "123456789123456789",
			args:    args{"123456789123456789"},
			wantErr: assert.Error,
		},
		{
			name:    "-1",
			args:    args{"-1"},
			want:    0xffffffff,
			wantErr: assert.NoError,
		},
		{
			name:    "herbert",
			args:    args{"herbert"},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrToUint32(tt.args.str)
			if !tt.wantErr(t, err, fmt.Sprintf("StrToUint32(%v)", tt.args.str)) {
				return
			}
			assert.Equalf(t, tt.want, got, "StrToUint32(%v)", tt.args.str)
		})
	}
}

func TestStrToUint8(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    uint8
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "1",
			args:    args{"1"},
			want:    1,
			wantErr: assert.NoError,
		},
		{
			name:    "123456789123456789",
			args:    args{"123456789123456789"},
			wantErr: assert.Error,
		},
		{
			name:    "-1",
			args:    args{"-1"},
			want:    0xff,
			wantErr: assert.NoError,
		},
		{
			name:    "herbert",
			args:    args{"herbert"},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StrToUint8(tt.args.str)
			if !tt.wantErr(t, err, fmt.Sprintf("StrToUint8(%v)", tt.args.str)) {
				return
			}
			assert.Equalf(t, tt.want, got, "StrToUint8(%v)", tt.args.str)
		})
	}
}
