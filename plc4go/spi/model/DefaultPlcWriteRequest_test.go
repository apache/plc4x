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
	"context"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/interceptors"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPlcWriteRequestBuilder_AddTag(t *testing.T) {
	type fields struct {
		writer                  spi.PlcWriter
		tagHandler              spi.PlcTagHandler
		valueHandler            spi.PlcValueHandler
		tagNames                []string
		tagAddresses            map[string]string
		tags                    map[string]apiModel.PlcTag
		values                  map[string]any
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	type args struct {
		name  string
		tag   apiModel.PlcTag
		value any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcWriteRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				tagAddresses: map[string]string{},
				tags:         map[string]apiModel.PlcTag{},
				values:       map[string]any{},
			},
			want: &DefaultPlcWriteRequestBuilder{
				tagNames:     []string{""},
				tagAddresses: map[string]string{},
				tags:         map[string]apiModel.PlcTag{"": nil},
				values:       map[string]any{"": nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultPlcWriteRequestBuilder{
				writer:                  tt.fields.writer,
				tagHandler:              tt.fields.tagHandler,
				valueHandler:            tt.fields.valueHandler,
				tagNames:                tt.fields.tagNames,
				tagAddresses:            tt.fields.tagAddresses,
				tags:                    tt.fields.tags,
				values:                  tt.fields.values,
				writeRequestInterceptor: tt.fields.writeRequestInterceptor,
			}
			assert.Equalf(t, tt.want, m.AddTag(tt.args.name, tt.args.tag, tt.args.value), "AddTag(%v, %v, %v)", tt.args.name, tt.args.tag, tt.args.value)
		})
	}
}

func TestDefaultPlcWriteRequestBuilder_AddTagAddress(t *testing.T) {
	type fields struct {
		writer                  spi.PlcWriter
		tagHandler              spi.PlcTagHandler
		valueHandler            spi.PlcValueHandler
		tagNames                []string
		tagAddresses            map[string]string
		tags                    map[string]apiModel.PlcTag
		values                  map[string]any
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	type args struct {
		name       string
		tagAddress string
		value      any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcWriteRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				tagAddresses: map[string]string{},
				tags:         map[string]apiModel.PlcTag{},
				values:       map[string]any{},
			},
			want: &DefaultPlcWriteRequestBuilder{
				tagNames:     []string{""},
				tagAddresses: map[string]string{"": ""},
				tags:         map[string]apiModel.PlcTag{},
				values:       map[string]any{"": nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultPlcWriteRequestBuilder{
				writer:                  tt.fields.writer,
				tagHandler:              tt.fields.tagHandler,
				valueHandler:            tt.fields.valueHandler,
				tagNames:                tt.fields.tagNames,
				tagAddresses:            tt.fields.tagAddresses,
				tags:                    tt.fields.tags,
				values:                  tt.fields.values,
				writeRequestInterceptor: tt.fields.writeRequestInterceptor,
			}
			assert.Equalf(t, tt.want, m.AddTagAddress(tt.args.name, tt.args.tagAddress, tt.args.value), "AddTagAddress(%v, %v, %v)", tt.args.name, tt.args.tagAddress, tt.args.value)
		})
	}
}

func TestDefaultPlcWriteRequestBuilder_Build(t *testing.T) {
	type fields struct {
		writer                  spi.PlcWriter
		tagHandler              spi.PlcTagHandler
		valueHandler            spi.PlcValueHandler
		tagNames                []string
		tagAddresses            map[string]string
		tags                    map[string]apiModel.PlcTag
		values                  map[string]any
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func(t *testing.T, fields *fields)
		want    apiModel.PlcWriteRequest
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "build it",
			want:    NewDefaultPlcWriteRequest(nil, nil, map[string]apiValues.PlcValue{}, nil, nil),
			wantErr: assert.NoError,
		},
		{
			name: "build it with tags",
			fields: fields{
				tagNames:     []string{"a"},
				tagAddresses: map[string]string{"a": ""},
				tags:         map[string]apiModel.PlcTag{},
			},
			setup: func(t *testing.T, fields *fields) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseTag(mock.Anything).Return(nil, nil)
				fields.tagHandler = handler
				valueHandler := NewMockPlcValueHandler(t)
				valueHandler.EXPECT().NewPlcValue(mock.Anything, mock.Anything).Return(nil, nil)
				fields.valueHandler = valueHandler
			},
			want:    NewDefaultPlcWriteRequest(map[string]apiModel.PlcTag{"a": nil}, []string{"a"}, map[string]apiValues.PlcValue{"a": nil}, nil, nil),
			wantErr: assert.NoError,
		},
		{
			name: "build it with value error",
			fields: fields{
				tagNames:     []string{"a"},
				tagAddresses: map[string]string{"a": ""},
				tags:         map[string]apiModel.PlcTag{},
			},
			setup: func(t *testing.T, fields *fields) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseTag(mock.Anything).Return(nil, nil)
				fields.tagHandler = handler
				valueHandler := NewMockPlcValueHandler(t)
				valueHandler.EXPECT().NewPlcValue(mock.Anything, mock.Anything).Return(nil, errors.New("nope"))
				fields.valueHandler = valueHandler
			},
			wantErr: assert.Error,
		},
		{
			name: "build it with parse error",
			fields: fields{
				tagNames:     []string{"a"},
				tagAddresses: map[string]string{"a": ""},
				tags:         map[string]apiModel.PlcTag{},
			},
			setup: func(t *testing.T, fields *fields) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseTag(mock.Anything).Return(nil, errors.New("nope"))
				fields.tagHandler = handler
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			m := &DefaultPlcWriteRequestBuilder{
				writer:                  tt.fields.writer,
				tagHandler:              tt.fields.tagHandler,
				valueHandler:            tt.fields.valueHandler,
				tagNames:                tt.fields.tagNames,
				tagAddresses:            tt.fields.tagAddresses,
				tags:                    tt.fields.tags,
				values:                  tt.fields.values,
				writeRequestInterceptor: tt.fields.writeRequestInterceptor,
			}
			got, err := m.Build()
			if !tt.wantErr(t, err, fmt.Sprintf("Build()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Build()")
		})
	}
}

func TestDefaultPlcWriteRequestBuilder_GetWriteRequestInterceptor(t *testing.T) {
	type fields struct {
		writer                  spi.PlcWriter
		tagHandler              spi.PlcTagHandler
		valueHandler            spi.PlcValueHandler
		tagNames                []string
		tagAddresses            map[string]string
		tags                    map[string]apiModel.PlcTag
		values                  map[string]any
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	tests := []struct {
		name   string
		fields fields
		want   interceptors.WriteRequestInterceptor
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultPlcWriteRequestBuilder{
				writer:                  tt.fields.writer,
				tagHandler:              tt.fields.tagHandler,
				valueHandler:            tt.fields.valueHandler,
				tagNames:                tt.fields.tagNames,
				tagAddresses:            tt.fields.tagAddresses,
				tags:                    tt.fields.tags,
				values:                  tt.fields.values,
				writeRequestInterceptor: tt.fields.writeRequestInterceptor,
			}
			assert.Equalf(t, tt.want, m.GetWriteRequestInterceptor(), "GetWriteRequestInterceptor()")
		})
	}
}

func TestDefaultPlcWriteRequestBuilder_GetWriter(t *testing.T) {
	type fields struct {
		writer                  spi.PlcWriter
		tagHandler              spi.PlcTagHandler
		valueHandler            spi.PlcValueHandler
		tagNames                []string
		tagAddresses            map[string]string
		tags                    map[string]apiModel.PlcTag
		values                  map[string]any
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	tests := []struct {
		name   string
		fields fields
		want   spi.PlcWriter
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &DefaultPlcWriteRequestBuilder{
				writer:                  tt.fields.writer,
				tagHandler:              tt.fields.tagHandler,
				valueHandler:            tt.fields.valueHandler,
				tagNames:                tt.fields.tagNames,
				tagAddresses:            tt.fields.tagAddresses,
				tags:                    tt.fields.tags,
				values:                  tt.fields.values,
				writeRequestInterceptor: tt.fields.writeRequestInterceptor,
			}
			assert.Equalf(t, tt.want, m.GetWriter(), "GetWriter()")
		})
	}
}

func TestDefaultPlcWriteRequest_Execute(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest    *DefaultPlcTagRequest
		values                  map[string]apiValues.PlcValue
		writer                  spi.PlcWriter
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	tests := []struct {
		name   string
		fields fields
		setup  func(t *testing.T, fields *fields)
		want   <-chan apiModel.PlcWriteRequestResult
	}{
		{
			name: "execute it",
			setup: func(t *testing.T, fields *fields) {
				writer := NewMockPlcWriter(t)
				writer.EXPECT().Write(mock.Anything, mock.Anything).Return(nil)
				fields.writer = writer
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			d := &DefaultPlcWriteRequest{
				DefaultPlcTagRequest:    tt.fields.DefaultPlcTagRequest,
				values:                  tt.fields.values,
				writer:                  tt.fields.writer,
				writeRequestInterceptor: tt.fields.writeRequestInterceptor,
			}
			assert.Equalf(t, tt.want, d.Execute(), "Execute()")
		})
	}
}

func TestDefaultPlcWriteRequest_ExecuteWithContext(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest    *DefaultPlcTagRequest
		values                  map[string]apiValues.PlcValue
		writer                  spi.PlcWriter
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		setup        func(t *testing.T, fields *fields)
		wantAsserter func(t *testing.T, results <-chan apiModel.PlcWriteRequestResult) bool
	}{
		{
			name: "execute it",
			setup: func(t *testing.T, fields *fields) {
				writer := NewMockPlcWriter(t)
				writer.EXPECT().Write(mock.Anything, mock.Anything).Return(nil)
				fields.writer = writer
			},
		},
		{
			name: "execute it with interceptor with one request",
			setup: func(t *testing.T, fields *fields) {
				{
					writer := NewMockPlcWriter(t)
					results := make(chan apiModel.PlcWriteRequestResult, 1)
					results <- NewDefaultPlcWriteRequestResult(
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
						NewDefaultPlcWriteResponse(nil, nil),
						nil,
					)
					writer.EXPECT().Write(mock.Anything, mock.Anything).Return(results)
					fields.writer = writer
				}
				{
					interceptor := NewMockWriteRequestInterceptor(t)
					interceptor.EXPECT().InterceptWriteRequest(mock.Anything, mock.Anything).Return([]apiModel.PlcWriteRequest{
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
					})
					fields.writeRequestInterceptor = interceptor
				}
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcWriteRequestResult) bool {
				timeout := time.NewTimer(100 * time.Millisecond)
				t.Cleanup(func() {
					utils.CleanupTimer(timeout)
				})
				select {
				case result := <-results:
					assert.NoError(t, result.GetErr())
					assert.NotNil(t, result.GetRequest())
					assert.NotNil(t, result.GetResponse())
				case <-timeout.C:
					t.Error("timeout getting a response")
				}
				return true
			},
		},
		{
			name: "execute it with interceptor with three request (panics)",
			setup: func(t *testing.T, fields *fields) {
				{
					writer := NewMockPlcWriter(t)
					results := make(chan apiModel.PlcWriteRequestResult, 1)
					results <- NewDefaultPlcWriteRequestResult(nil, nil, nil)
					writer.EXPECT().Write(mock.Anything, mock.Anything).Return(results)
					fields.writer = writer
				}
				{
					interceptor := NewMockWriteRequestInterceptor(t)
					interceptor.EXPECT().InterceptWriteRequest(mock.Anything, mock.Anything).Return([]apiModel.PlcWriteRequest{
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
					})
					fields.writeRequestInterceptor = interceptor
				}
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcWriteRequestResult) bool {
				timeout := time.NewTimer(100 * time.Millisecond)
				t.Cleanup(func() {
					utils.CleanupTimer(timeout)
				})
				select {
				case result := <-results:
					assert.Error(t, result.GetErr())
					assert.NotNil(t, result.GetRequest())
					assert.Nil(t, result.GetResponse())
				case <-timeout.C:
					t.Error("timeout getting a response")
				}
				return true
			},
		},
		{
			name: "execute it with interceptor with three request (context done)",
			args: args{
				ctx: func() context.Context {
					timeout, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					return timeout
				}(),
			},
			setup: func(t *testing.T, fields *fields) {
				{
					writer := NewMockPlcWriter(t)
					results := make(chan apiModel.PlcWriteRequestResult, 1)
					results <- NewDefaultPlcWriteRequestResult(nil, nil, nil)
					writer.EXPECT().Write(mock.Anything, mock.Anything).Return(results)
					fields.writer = writer
				}
				{
					interceptor := NewMockWriteRequestInterceptor(t)
					interceptor.EXPECT().InterceptWriteRequest(mock.Anything, mock.Anything).Return([]apiModel.PlcWriteRequest{
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
					})
					fields.writeRequestInterceptor = interceptor
				}
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcWriteRequestResult) bool {
				timeout := time.NewTimer(100 * time.Millisecond)
				t.Cleanup(func() {
					utils.CleanupTimer(timeout)
				})
				select {
				case result := <-results:
					assert.Error(t, result.GetErr())
					assert.NotNil(t, result.GetRequest())
					assert.Nil(t, result.GetResponse())
				case <-timeout.C:
					t.Error("timeout getting a response")
				}
				return true
			},
		},
		{
			name: "execute it with interceptor with three request",
			args: args{
				ctx: func() context.Context {
					timeout, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
					t.Cleanup(cancelFunc)
					return timeout
				}(),
			},
			setup: func(t *testing.T, fields *fields) {
				{
					writer := NewMockPlcWriter(t)
					results := make(chan apiModel.PlcWriteRequestResult, 1)
					results <- NewDefaultPlcWriteRequestResult(nil, nil, nil)
					writer.EXPECT().
						Write(mock.Anything, mock.Anything).
						Return(results)
					close(results)
					fields.writer = writer
				}
				{
					interceptor := NewMockWriteRequestInterceptor(t)
					interceptor.EXPECT().InterceptWriteRequest(mock.Anything, mock.Anything).Return([]apiModel.PlcWriteRequest{
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
						NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
					})
					interceptor.EXPECT().
						ProcessWriteResponses(mock.Anything, mock.Anything, mock.Anything).
						Return(
							NewDefaultPlcWriteRequestResult(
								NewDefaultPlcWriteRequest(nil, nil, nil, nil, nil),
								NewDefaultPlcWriteResponse(nil, nil),
								nil,
							),
						).Maybe()
					fields.writeRequestInterceptor = interceptor
				}
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcWriteRequestResult) bool {
				timeout := time.NewTimer(100 * time.Millisecond)
				t.Cleanup(func() {
					utils.CleanupTimer(timeout)
				})
				select {
				case result := <-results:
					assert.NoError(t, result.GetErr())
					assert.NotNil(t, result.GetRequest())
					assert.NotNil(t, result.GetResponse())
				case <-timeout.C:
					t.Error("timeout getting a response")
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			d := &DefaultPlcWriteRequest{
				DefaultPlcTagRequest:    tt.fields.DefaultPlcTagRequest,
				values:                  tt.fields.values,
				writer:                  tt.fields.writer,
				writeRequestInterceptor: tt.fields.writeRequestInterceptor,
			}
			result := d.ExecuteWithContext(tt.args.ctx)
			if tt.wantAsserter != nil {
				assert.Truef(t, tt.wantAsserter(t, result), "ExecuteWithContext(%v)", tt.args.ctx)
			}
		})
	}
}

func TestDefaultPlcWriteRequest_GetValue(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest    *DefaultPlcTagRequest
		values                  map[string]apiValues.PlcValue
		writer                  spi.PlcWriter
		writeRequestInterceptor interceptors.WriteRequestInterceptor
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
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcWriteRequest{
				DefaultPlcTagRequest:    tt.fields.DefaultPlcTagRequest,
				values:                  tt.fields.values,
				writer:                  tt.fields.writer,
				writeRequestInterceptor: tt.fields.writeRequestInterceptor,
			}
			assert.Equalf(t, tt.want, d.GetValue(tt.args.name), "GetValue(%v)", tt.args.name)
		})
	}
}

func TestDefaultPlcWriteRequest_GetWriteRequestInterceptor(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest    *DefaultPlcTagRequest
		values                  map[string]apiValues.PlcValue
		writer                  spi.PlcWriter
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	tests := []struct {
		name   string
		fields fields
		want   interceptors.WriteRequestInterceptor
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcWriteRequest{
				DefaultPlcTagRequest:    tt.fields.DefaultPlcTagRequest,
				values:                  tt.fields.values,
				writer:                  tt.fields.writer,
				writeRequestInterceptor: tt.fields.writeRequestInterceptor,
			}
			assert.Equalf(t, tt.want, d.GetWriteRequestInterceptor(), "GetWriteRequestInterceptor()")
		})
	}
}

func TestDefaultPlcWriteRequest_GetWriter(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest    *DefaultPlcTagRequest
		values                  map[string]apiValues.PlcValue
		writer                  spi.PlcWriter
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	tests := []struct {
		name   string
		fields fields
		want   spi.PlcWriter
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcWriteRequest{
				DefaultPlcTagRequest:    tt.fields.DefaultPlcTagRequest,
				values:                  tt.fields.values,
				writer:                  tt.fields.writer,
				writeRequestInterceptor: tt.fields.writeRequestInterceptor,
			}
			assert.Equalf(t, tt.want, d.GetWriter(), "GetWriter()")
		})
	}
}

func TestNewDefaultPlcWriteRequest(t *testing.T) {
	type args struct {
		tags                    map[string]apiModel.PlcTag
		tagNames                []string
		values                  map[string]apiValues.PlcValue
		writer                  spi.PlcWriter
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcWriteRequest
	}{
		{
			name: "create it",
			want: &DefaultPlcWriteRequest{
				DefaultPlcTagRequest: NewDefaultPlcTagRequest(nil, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcWriteRequest(tt.args.tags, tt.args.tagNames, tt.args.values, tt.args.writer, tt.args.writeRequestInterceptor), "NewDefaultPlcWriteRequest(%v, %v, %v, %v, %v)", tt.args.tags, tt.args.tagNames, tt.args.values, tt.args.writer, tt.args.writeRequestInterceptor)
		})
	}
}

func TestNewDefaultPlcWriteRequestBuilder(t *testing.T) {
	type args struct {
		tagHandler   spi.PlcTagHandler
		valueHandler spi.PlcValueHandler
		writer       spi.PlcWriter
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcWriteRequestBuilder
	}{
		{
			name: "create it",
			want: &DefaultPlcWriteRequestBuilder{
				tagNames:     []string{},
				tagAddresses: map[string]string{},
				tags:         map[string]apiModel.PlcTag{},
				values:       map[string]any{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcWriteRequestBuilder(tt.args.tagHandler, tt.args.valueHandler, tt.args.writer), "NewDefaultPlcWriteRequestBuilder(%v, %v, %v)", tt.args.tagHandler, tt.args.valueHandler, tt.args.writer)
		})
	}
}

func TestNewDefaultPlcWriteRequestBuilderWithInterceptor(t *testing.T) {
	type args struct {
		tagHandler              spi.PlcTagHandler
		valueHandler            spi.PlcValueHandler
		writer                  spi.PlcWriter
		writeRequestInterceptor interceptors.WriteRequestInterceptor
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcWriteRequestBuilder
	}{
		{
			name: "create it",
			want: &DefaultPlcWriteRequestBuilder{
				tagNames:     []string{},
				tagAddresses: map[string]string{},
				tags:         map[string]apiModel.PlcTag{},
				values:       map[string]any{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcWriteRequestBuilderWithInterceptor(tt.args.tagHandler, tt.args.valueHandler, tt.args.writer, tt.args.writeRequestInterceptor), "NewDefaultPlcWriteRequestBuilderWithInterceptor(%v, %v, %v, %v)", tt.args.tagHandler, tt.args.valueHandler, tt.args.writer, tt.args.writeRequestInterceptor)
		})
	}
}
