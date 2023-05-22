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
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/interceptors"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPlcReadRequestBuilder_AddTag(t *testing.T) {
	type fields struct {
		reader                 spi.PlcReader
		tagHandler             spi.PlcTagHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		readRequestInterceptor interceptors.ReadRequestInterceptor
	}
	type args struct {
		name string
		tag  apiModel.PlcTag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcReadRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				tagAddresses: map[string]string{},
				tags:         map[string]apiModel.PlcTag{},
			},
			want: &DefaultPlcReadRequestBuilder{
				tagNames:     []string{""},
				tagAddresses: map[string]string{},
				tags:         map[string]apiModel.PlcTag{"": nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcReadRequestBuilder{
				reader:                 tt.fields.reader,
				tagHandler:             tt.fields.tagHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				readRequestInterceptor: tt.fields.readRequestInterceptor,
			}
			assert.Equalf(t, tt.want, d.AddTag(tt.args.name, tt.args.tag), "AddTag(%v, %v)", tt.args.name, tt.args.tag)
		})
	}
}

func TestDefaultPlcReadRequestBuilder_AddTagAddress(t *testing.T) {
	type fields struct {
		reader                 spi.PlcReader
		tagHandler             spi.PlcTagHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		readRequestInterceptor interceptors.ReadRequestInterceptor
	}
	type args struct {
		name  string
		query string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcReadRequestBuilder
	}{
		{
			name: "add it",
			fields: fields{
				tagAddresses: map[string]string{},
				tags:         map[string]apiModel.PlcTag{},
			},
			want: &DefaultPlcReadRequestBuilder{
				tagNames:     []string{""},
				tagAddresses: map[string]string{"": ""},
				tags:         map[string]apiModel.PlcTag{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcReadRequestBuilder{
				reader:                 tt.fields.reader,
				tagHandler:             tt.fields.tagHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				readRequestInterceptor: tt.fields.readRequestInterceptor,
			}
			assert.Equalf(t, tt.want, d.AddTagAddress(tt.args.name, tt.args.query), "AddTagAddress(%v, %v)", tt.args.name, tt.args.query)
		})
	}
}

func TestDefaultPlcReadRequestBuilder_Build(t *testing.T) {
	type fields struct {
		reader                 spi.PlcReader
		tagHandler             spi.PlcTagHandler
		tagNames               []string
		tagAddresses           map[string]string
		tags                   map[string]apiModel.PlcTag
		readRequestInterceptor interceptors.ReadRequestInterceptor
	}
	tests := []struct {
		name      string
		fields    fields
		mockSetup func(t *testing.T, fields *fields)
		want      apiModel.PlcReadRequest
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:    "build it",
			want:    NewDefaultPlcReadRequest(nil, nil, nil, nil),
			wantErr: assert.NoError,
		},
		{
			name: "build it with tags",
			fields: fields{
				tagNames:     []string{"a"},
				tagAddresses: map[string]string{"a": ""},
				tags:         map[string]apiModel.PlcTag{},
			},
			mockSetup: func(t *testing.T, fields *fields) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseTag(mock.Anything).Return(nil, nil)
				fields.tagHandler = handler
			},
			want:    NewDefaultPlcReadRequest(map[string]apiModel.PlcTag{"a": nil}, []string{"a"}, nil, nil),
			wantErr: assert.NoError,
		},
		{
			name: "build it with parse error",
			fields: fields{
				tagNames:     []string{"a"},
				tagAddresses: map[string]string{"a": ""},
				tags:         map[string]apiModel.PlcTag{},
			},
			mockSetup: func(t *testing.T, fields *fields) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseTag(mock.Anything).Return(nil, errors.New("nope"))
				fields.tagHandler = handler
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields)
			}
			d := &DefaultPlcReadRequestBuilder{
				reader:                 tt.fields.reader,
				tagHandler:             tt.fields.tagHandler,
				tagNames:               tt.fields.tagNames,
				tagAddresses:           tt.fields.tagAddresses,
				tags:                   tt.fields.tags,
				readRequestInterceptor: tt.fields.readRequestInterceptor,
			}
			got, err := d.Build()
			if !tt.wantErr(t, err, fmt.Sprintf("Build()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Build()")
		})
	}
}

func TestDefaultPlcReadRequest_Execute(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest   *DefaultPlcTagRequest
		reader                 spi.PlcReader
		readRequestInterceptor interceptors.ReadRequestInterceptor
	}
	tests := []struct {
		name      string
		fields    fields
		mockSetup func(t *testing.T, fields *fields)
		want      <-chan apiModel.PlcReadRequestResult
	}{
		{
			name: "execute it",
			mockSetup: func(t *testing.T, fields *fields) {
				reader := NewMockPlcReader(t)
				reader.EXPECT().Read(mock.Anything, mock.Anything).Return(nil)
				fields.reader = reader
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields)
			}
			d := &DefaultPlcReadRequest{
				DefaultPlcTagRequest:   tt.fields.DefaultPlcTagRequest,
				reader:                 tt.fields.reader,
				readRequestInterceptor: tt.fields.readRequestInterceptor,
			}
			assert.Equalf(t, tt.want, d.Execute(), "Execute()")
		})
	}
}

func TestDefaultPlcReadRequest_ExecuteWithContext(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest   *DefaultPlcTagRequest
		reader                 spi.PlcReader
		readRequestInterceptor interceptors.ReadRequestInterceptor
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		mockSetup    func(t *testing.T, fields *fields)
		wantAsserter func(t *testing.T, results <-chan apiModel.PlcReadRequestResult) bool
	}{
		{
			name: "execute it",
			mockSetup: func(t *testing.T, fields *fields) {
				reader := NewMockPlcReader(t)
				reader.EXPECT().Read(mock.Anything, mock.Anything).Return(nil)
				fields.reader = reader
			},
		},
		{
			name: "execute it with interceptor with one request",
			mockSetup: func(t *testing.T, fields *fields) {
				{
					reader := NewMockPlcReader(t)
					results := make(chan apiModel.PlcReadRequestResult, 1)
					results <- NewDefaultPlcReadRequestResult(
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
						NewDefaultPlcReadResponse(nil, nil, nil),
						nil,
					)
					reader.EXPECT().Read(mock.Anything, mock.Anything).Return(results)
					fields.reader = reader
				}
				{
					interceptor := NewMockReadRequestInterceptor(t)
					interceptor.EXPECT().InterceptReadRequest(mock.Anything, mock.Anything).Return([]apiModel.PlcReadRequest{
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
					})
					fields.readRequestInterceptor = interceptor
				}
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcReadRequestResult) bool {
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
			mockSetup: func(t *testing.T, fields *fields) {
				{
					reader := NewMockPlcReader(t)
					results := make(chan apiModel.PlcReadRequestResult, 1)
					results <- NewDefaultPlcReadRequestResult(nil, nil, nil)
					reader.EXPECT().Read(mock.Anything, mock.Anything).Return(results)
					fields.reader = reader
				}
				{
					interceptor := NewMockReadRequestInterceptor(t)
					interceptor.EXPECT().InterceptReadRequest(mock.Anything, mock.Anything).Return([]apiModel.PlcReadRequest{
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
					})
					fields.readRequestInterceptor = interceptor
				}
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcReadRequestResult) bool {
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
			mockSetup: func(t *testing.T, fields *fields) {
				{
					reader := NewMockPlcReader(t)
					results := make(chan apiModel.PlcReadRequestResult, 1)
					results <- NewDefaultPlcReadRequestResult(nil, nil, nil)
					reader.EXPECT().Read(mock.Anything, mock.Anything).Return(results)
					fields.reader = reader
				}
				{
					interceptor := NewMockReadRequestInterceptor(t)
					interceptor.EXPECT().InterceptReadRequest(mock.Anything, mock.Anything).Return([]apiModel.PlcReadRequest{
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
					})
					fields.readRequestInterceptor = interceptor
				}
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcReadRequestResult) bool {
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
			mockSetup: func(t *testing.T, fields *fields) {
				{
					reader := NewMockPlcReader(t)
					results := make(chan apiModel.PlcReadRequestResult, 1)
					results <- NewDefaultPlcReadRequestResult(nil, nil, nil)
					reader.EXPECT().
						Read(mock.Anything, mock.Anything).
						Return(results)
					close(results)
					fields.reader = reader
				}
				{
					interceptor := NewMockReadRequestInterceptor(t)
					interceptor.EXPECT().InterceptReadRequest(mock.Anything, mock.Anything).Return([]apiModel.PlcReadRequest{
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
						NewDefaultPlcReadRequest(nil, nil, nil, nil),
					})
					interceptor.EXPECT().
						ProcessReadResponses(mock.Anything, mock.Anything, mock.Anything).
						Return(
							NewDefaultPlcReadRequestResult(
								NewDefaultPlcReadRequest(nil, nil, nil, nil),
								NewDefaultPlcReadResponse(nil, nil, nil),
								nil,
							),
						).Maybe()
					fields.readRequestInterceptor = interceptor
				}
			},
			wantAsserter: func(t *testing.T, results <-chan apiModel.PlcReadRequestResult) bool {
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
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields)
			}
			d := &DefaultPlcReadRequest{
				DefaultPlcTagRequest:   tt.fields.DefaultPlcTagRequest,
				reader:                 tt.fields.reader,
				readRequestInterceptor: tt.fields.readRequestInterceptor,
			}
			result := d.ExecuteWithContext(tt.args.ctx)
			if tt.wantAsserter != nil {
				assert.True(t, tt.wantAsserter(t, result), "ExecuteWithContext(%v)", tt.args.ctx)
			}
		})
	}
}

func TestDefaultPlcReadRequest_GetReadRequestInterceptor(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest   *DefaultPlcTagRequest
		reader                 spi.PlcReader
		readRequestInterceptor interceptors.ReadRequestInterceptor
	}
	tests := []struct {
		name   string
		fields fields
		want   interceptors.ReadRequestInterceptor
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcReadRequest{
				DefaultPlcTagRequest:   tt.fields.DefaultPlcTagRequest,
				reader:                 tt.fields.reader,
				readRequestInterceptor: tt.fields.readRequestInterceptor,
			}
			assert.Equalf(t, tt.want, d.GetReadRequestInterceptor(), "GetReadRequestInterceptor()")
		})
	}
}

func TestDefaultPlcReadRequest_GetReader(t *testing.T) {
	type fields struct {
		DefaultPlcTagRequest   *DefaultPlcTagRequest
		reader                 spi.PlcReader
		readRequestInterceptor interceptors.ReadRequestInterceptor
	}
	tests := []struct {
		name   string
		fields fields
		want   spi.PlcReader
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcReadRequest{
				DefaultPlcTagRequest:   tt.fields.DefaultPlcTagRequest,
				reader:                 tt.fields.reader,
				readRequestInterceptor: tt.fields.readRequestInterceptor,
			}
			assert.Equalf(t, tt.want, d.GetReader(), "GetReader()")
		})
	}
}

func TestNewDefaultPlcReadRequest(t *testing.T) {
	type args struct {
		tags                   map[string]apiModel.PlcTag
		tagNames               []string
		reader                 spi.PlcReader
		readRequestInterceptor interceptors.ReadRequestInterceptor
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcReadRequest
	}{
		{
			name: "create it",
			want: func() apiModel.PlcReadRequest {
				d := &DefaultPlcReadRequest{}
				d.DefaultPlcTagRequest = NewDefaultPlcTagRequest(nil, nil)
				return d
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcReadRequest(tt.args.tags, tt.args.tagNames, tt.args.reader, tt.args.readRequestInterceptor), "NewDefaultPlcReadRequest(%v, %v, %v, %v)", tt.args.tags, tt.args.tagNames, tt.args.reader, tt.args.readRequestInterceptor)
		})
	}
}

func TestNewDefaultPlcReadRequestBuilder(t *testing.T) {
	type args struct {
		tagHandler spi.PlcTagHandler
		reader     spi.PlcReader
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcReadRequestBuilder
	}{
		{
			name: "create it",
			want: &DefaultPlcReadRequestBuilder{
				tagNames:     []string{},
				tagAddresses: map[string]string{},
				tags:         map[string]apiModel.PlcTag{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcReadRequestBuilder(tt.args.tagHandler, tt.args.reader), "NewDefaultPlcReadRequestBuilder(%v, %v)", tt.args.tagHandler, tt.args.reader)
		})
	}
}

func TestNewDefaultPlcReadRequestBuilderWithInterceptor(t *testing.T) {
	type args struct {
		tagHandler             spi.PlcTagHandler
		reader                 spi.PlcReader
		readRequestInterceptor interceptors.ReadRequestInterceptor
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcReadRequestBuilder
	}{
		{
			name: "create it",
			want: &DefaultPlcReadRequestBuilder{
				tagNames:     []string{},
				tagAddresses: map[string]string{},
				tags:         map[string]apiModel.PlcTag{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcReadRequestBuilderWithInterceptor(tt.args.tagHandler, tt.args.reader, tt.args.readRequestInterceptor), "NewDefaultPlcReadRequestBuilderWithInterceptor(%v, %v, %v)", tt.args.tagHandler, tt.args.reader, tt.args.readRequestInterceptor)
		})
	}
}
