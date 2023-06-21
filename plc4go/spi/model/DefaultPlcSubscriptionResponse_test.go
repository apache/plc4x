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
	"fmt"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/stretchr/testify/mock"
	"testing"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPlcSubscriptionResponse_GetRequest(t *testing.T) {
	type fields struct {
		request apiModel.PlcSubscriptionRequest
		values  map[string]*DefaultPlcSubscriptionResponseItem
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcSubscriptionRequest
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetRequest(), "GetRequest()")
		})
	}
}

func TestDefaultPlcSubscriptionResponse_GetResponseCode(t *testing.T) {
	type fields struct {
		request apiModel.PlcSubscriptionRequest
		values  map[string]*DefaultPlcSubscriptionResponseItem
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
			name: "get code",
			want: apiModel.PlcResponseCode_NOT_FOUND,
		},
		{
			name: "get code",
			fields: fields{
				values: map[string]*DefaultPlcSubscriptionResponseItem{
					"": {code: apiModel.PlcResponseCode_INVALID_DATATYPE},
				},
			},
			want: apiModel.PlcResponseCode_INVALID_DATATYPE,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetResponseCode(tt.args.name), "GetResponseCode(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionResponse_GetSubscriptionHandle(t *testing.T) {
	type fields struct {
		request apiModel.PlcSubscriptionRequest
		values  map[string]*DefaultPlcSubscriptionResponseItem
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    apiModel.PlcSubscriptionHandle
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "get it not found",
			wantErr: assert.Error,
		},
		{
			name: "get it (failed to subscribe)",
			fields: fields{
				values: map[string]*DefaultPlcSubscriptionResponseItem{
					"": {},
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "get it",
			fields: fields{
				values: map[string]*DefaultPlcSubscriptionResponseItem{
					"": {code: apiModel.PlcResponseCode_OK},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			got, err := d.GetSubscriptionHandle(tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("GetSubscriptionHandle(%v)", tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GetSubscriptionHandle(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcSubscriptionResponse_GetSubscriptionHandles(t *testing.T) {
	type fields struct {
		request apiModel.PlcSubscriptionRequest
		values  map[string]*DefaultPlcSubscriptionResponseItem
	}
	tests := []struct {
		name   string
		fields fields
		want   []apiModel.PlcSubscriptionHandle
	}{
		{
			name: "get em",
			want: []apiModel.PlcSubscriptionHandle{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetSubscriptionHandles(), "GetSubscriptionHandles()")
		})
	}
}

func TestDefaultPlcSubscriptionResponse_GetTagNames(t *testing.T) {
	type fields struct {
		request apiModel.PlcSubscriptionRequest
		values  map[string]*DefaultPlcSubscriptionResponseItem
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "get em",
		},
		{
			name: "get em with filling",
			fields: fields{
				request: NewDefaultPlcSubscriptionRequest(nil, []string{"asd"}, nil, nil, nil, nil),
				values:  map[string]*DefaultPlcSubscriptionResponseItem{"asd": nil},
			},
			want: []string{"asd"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.GetTagNames(), "GetTagNames()")
		})
	}
}

func TestDefaultPlcSubscriptionResponse_IsAPlcMessage(t *testing.T) {
	type fields struct {
		request apiModel.PlcSubscriptionRequest
		values  map[string]*DefaultPlcSubscriptionResponseItem
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "is it",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcSubscriptionResponse{
				request: tt.fields.request,
				values:  tt.fields.values,
			}
			assert.Equalf(t, tt.want, d.IsAPlcMessage(), "IsAPlcMessage()")
		})
	}
}

func TestNewDefaultPlcSubscriptionResponse(t *testing.T) {
	type args struct {
		request       apiModel.PlcSubscriptionRequest
		responseCodes map[string]apiModel.PlcResponseCode
		values        map[string]apiModel.PlcSubscriptionHandle
	}
	tests := []struct {
		name  string
		args  args
		setup func(t *testing.T, args *args, want *apiModel.PlcSubscriptionResponse)
		want  apiModel.PlcSubscriptionResponse
	}{
		{
			name: "create it",
			args: args{
				request:       NewDefaultPlcSubscriptionRequest(nil, nil, nil, nil, nil, nil),
				responseCodes: nil,
				values:        nil,
			},
			want: &DefaultPlcSubscriptionResponse{
				request: NewDefaultPlcSubscriptionRequest(nil, nil, nil, nil, nil, nil),
				values:  map[string]*DefaultPlcSubscriptionResponseItem{},
			},
		},
		{
			name: "create it with pre-registered consumers",
			args: args{
				request: NewDefaultPlcSubscriptionRequest(
					nil,
					nil,
					nil,
					nil,
					nil,
					map[string][]apiModel.PlcSubscriptionEventConsumer{
						"asd":  {nil},
						"asd2": {nil},
					},
				),
				responseCodes: map[string]apiModel.PlcResponseCode{
					"asd":  apiModel.PlcResponseCode_OK,
					"asd2": apiModel.PlcResponseCode_OK,
				},
				values: map[string]apiModel.PlcSubscriptionHandle{},
			},
			setup: func(t *testing.T, args *args, want *apiModel.PlcSubscriptionResponse) {
				handle := NewDefaultPlcSubscriptionHandle(func() spi.PlcSubscriber {
					subscriber := NewMockPlcSubscriber(t)
					subscriber.EXPECT().Register(mock.Anything, mock.Anything).Return(nil)
					return subscriber
				}())
				args.values["asd2"] = handle
				(*want).(*DefaultPlcSubscriptionResponse).values["asd2"] = NewDefaultPlcSubscriptionResponseItem(apiModel.PlcResponseCode_OK, handle)
			},
			want: &DefaultPlcSubscriptionResponse{
				request: NewDefaultPlcSubscriptionRequest(nil, nil, nil, nil, nil,
					map[string][]apiModel.PlcSubscriptionEventConsumer{
						"asd":  {nil},
						"asd2": {nil},
					},
				),
				values: map[string]*DefaultPlcSubscriptionResponseItem{
					"asd": NewDefaultPlcSubscriptionResponseItem(apiModel.PlcResponseCode_OK, nil),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.args, &tt.want)
			}
			assert.Equalf(t, tt.want, NewDefaultPlcSubscriptionResponse(tt.args.request, tt.args.responseCodes, tt.args.values), "NewDefaultPlcSubscriptionResponse(%v, %v, %v)", tt.args.request, tt.args.responseCodes, tt.args.values)
		})
	}
}
