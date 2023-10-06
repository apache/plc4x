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

package model

import (
	"testing"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPlcReadResponse_GetRequest(t *testing.T) {
	type fields struct {
		request apiModel.PlcReadRequest
		values  map[string]*ResponseItem
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcReadRequest
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcReadResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetRequest(), "GetRequest()")
		})
	}
}

func TestDefaultPlcReadResponse_GetResponseCode(t *testing.T) {
	type fields struct {
		request apiModel.PlcReadRequest
		values  map[string]*ResponseItem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcResponseCode
	}{
		{
			name: "get it (not found)",
			want: apiModel.PlcResponseCode_NOT_FOUND,
		},
		{
			name: "get it",
			fields: fields{
				values: map[string]*ResponseItem{
					"something": {
						code: apiModel.PlcResponseCode_OK,
					},
				},
			},
			args: args{
				name: "something",
			},
			want: apiModel.PlcResponseCode_OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcReadResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetResponseCode(tt.args.name), "GetResponseCode(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcReadResponse_GetTagNames(t *testing.T) {
	type fields struct {
		request apiModel.PlcReadRequest
		values  map[string]*ResponseItem
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "get it",
			fields: fields{
				request: NewDefaultPlcReadRequest(nil, []string{"tag1", "tag2"}, nil, nil),
				values: map[string]*ResponseItem{
					"tag1": nil,
					"tag2": nil,
				},
			},
			want: []string{"tag1", "tag2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcReadResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetTagNames(), "GetTagNames()")
		})
	}
}

func TestDefaultPlcReadResponse_GetValue(t *testing.T) {
	type fields struct {
		request apiModel.PlcReadRequest
		values  map[string]*ResponseItem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiValues.PlcValue
	}{
		{
			name: "get it (not found)",
			want: spiValues.PlcNull{},
		},
		{
			name: "get it",
			fields: fields{
				values: map[string]*ResponseItem{
					"something": {
						code:  apiModel.PlcResponseCode_OK,
						value: spiValues.NewPlcSTRING("yes"),
					},
				},
			},
			args: args{
				name: "something",
			},
			want: spiValues.NewPlcSTRING("yes"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcReadResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetValue(tt.args.name), "GetValue(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcReadResponse_IsAPlcMessage(t *testing.T) {
	type fields struct {
		request apiModel.PlcReadRequest
		values  map[string]*ResponseItem
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "it is",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcReadResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.IsAPlcMessage(), "IsAPlcMessage()")
		})
	}
}

func TestNewDefaultPlcReadResponse(t *testing.T) {
	type args struct {
		request       apiModel.PlcReadRequest
		responseCodes map[string]apiModel.PlcResponseCode
		values        map[string]apiValues.PlcValue
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcReadResponse
	}{
		{
			name: "create it",
			want: &DefaultPlcReadResponse{values: map[string]*ResponseItem{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcReadResponse(tt.args.request, tt.args.responseCodes, tt.args.values), "NewDefaultPlcReadResponse(%v, %v, %v)", tt.args.request, tt.args.responseCodes, tt.args.values)
		})
	}
}
