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

package interceptors

import (
	"context"
	"errors"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

func TestNewSingleItemRequestInterceptor(t *testing.T) {
	type args struct {
		readRequestFactory   readRequestFactory
		writeRequestFactory  writeRequestFactory
		readResponseFactory  readResponseFactory
		writeResponseFactory writeResponseFactory
	}
	tests := []struct {
		name string
		args args
		want SingleItemRequestInterceptor
	}{
		{
			name: "create one",
			want: SingleItemRequestInterceptor{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSingleItemRequestInterceptor(tt.args.readRequestFactory, tt.args.writeRequestFactory, tt.args.readResponseFactory, tt.args.writeResponseFactory); !assert.Equal(t, tt.want, got) {
				t.Errorf("NewSingleItemRequestInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSingleItemRequestInterceptor_InterceptReadRequest(t *testing.T) {
	type fields struct {
		readRequestFactory   func(t *testing.T) readRequestFactory
		writeRequestFactory  func(t *testing.T) writeRequestFactory
		readResponseFactory  func(t *testing.T) readResponseFactory
		writeResponseFactory func(t *testing.T) writeResponseFactory
	}
	type args struct {
		ctx         context.Context
		readRequest apiModel.PlcReadRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		mockSetup  func(t *testing.T, fields *fields, args *args)
		wantAssert func(t *testing.T, args args, got []apiModel.PlcReadRequest) bool
	}{
		{
			name: "nil stays nil",
		},
		{
			name: "read request with no tags",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				plcReadRequest := NewMockPlcReadRequest(t)
				plcReadRequest.EXPECT().GetTagNames().Return(nil)
				args.readRequest = plcReadRequest
			},
		},
		{
			name: "read request with 1 tag",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				plcReadRequest := NewMockPlcReadRequest(t)
				plcReadRequest.EXPECT().GetTagNames().Return([]string{"a tag"})
				args.readRequest = plcReadRequest
			},
			wantAssert: func(t *testing.T, args args, got []apiModel.PlcReadRequest) bool {
				return assert.Contains(t, got, args.readRequest)
			},
		},
		{
			name: "read request with 2 tags",
			fields: fields{
				readRequestFactory: func(t *testing.T) readRequestFactory {
					return func(tags map[string]apiModel.PlcTag, tagNames []string, _ spi.PlcReader, _ ReadRequestInterceptor) apiModel.PlcReadRequest {
						plcReadRequest := NewMockPlcReadRequest(t)
						expect := plcReadRequest.EXPECT()
						expect.GetTagNames().Return(tagNames)
						expect.GetTag(mock.Anything).RunAndReturn(func(s string) apiModel.PlcTag {
							return tags[s]
						})
						return plcReadRequest
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				plcReadRequest := NewMockPlcReadRequest(t)
				expect := plcReadRequest.EXPECT()
				expect.GetTagNames().Return([]string{"1 tag", "2 tag"})
				expect.GetTag(mock.Anything).Return(nil)
				expect.GetReader().Return(nil)
				expect.GetReadRequestInterceptor().Return(nil)
				args.readRequest = plcReadRequest
			},
			wantAssert: func(t *testing.T, args args, got []apiModel.PlcReadRequest) bool {
				assert.Len(t, got, 2)
				request1 := got[0]
				assert.Len(t, request1.GetTagNames(), 1)
				assert.Equal(t, nil, request1.GetTag(request1.GetTagNames()[0]))
				request2 := got[1]
				assert.Len(t, request2.GetTagNames(), 1)
				assert.Equal(t, nil, request2.GetTag(request2.GetTagNames()[0]))
				return true
			},
		},
		{
			name: "read request with 2 tags aborted",
			fields: fields{
				readRequestFactory: func(t *testing.T) readRequestFactory {
					return func(tags map[string]apiModel.PlcTag, tagNames []string, _ spi.PlcReader, _ ReadRequestInterceptor) apiModel.PlcReadRequest {
						plcReadRequest := NewMockPlcReadRequest(t)
						expect := plcReadRequest.EXPECT()
						expect.GetTagNames().Return(tagNames)
						expect.GetTag(mock.Anything).RunAndReturn(func(s string) apiModel.PlcTag {
							return tags[s]
						})
						return plcReadRequest
					}
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					return ctx
				}(),
			},
			wantAssert: func(t *testing.T, args args, got []apiModel.PlcReadRequest) bool {
				return true
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				plcReadRequest := NewMockPlcReadRequest(t)
				plcReadRequest.EXPECT().GetTagNames().Return([]string{"1 tag", "2 tag"})
				args.readRequest = plcReadRequest
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			if tt.fields.readRequestFactory == nil {
				tt.fields.readRequestFactory = func(t *testing.T) readRequestFactory {
					return nil
				}
			}
			if tt.fields.writeRequestFactory == nil {
				tt.fields.writeRequestFactory = func(t *testing.T) writeRequestFactory {
					return nil
				}
			}
			if tt.fields.readResponseFactory == nil {
				tt.fields.readResponseFactory = func(t *testing.T) readResponseFactory {
					return nil
				}
			}
			if tt.fields.writeResponseFactory == nil {
				tt.fields.writeResponseFactory = func(t *testing.T) writeResponseFactory {
					return nil
				}
			}
			if tt.wantAssert == nil {
				tt.wantAssert = func(t *testing.T, args args, got []apiModel.PlcReadRequest) bool {
					return true
				}
			}
			m := SingleItemRequestInterceptor{
				readRequestFactory:   tt.fields.readRequestFactory(t),
				writeRequestFactory:  tt.fields.writeRequestFactory(t),
				readResponseFactory:  tt.fields.readResponseFactory(t),
				writeResponseFactory: tt.fields.writeResponseFactory(t),
			}
			if got := m.InterceptReadRequest(tt.args.ctx, tt.args.readRequest); !assert.True(t, tt.wantAssert(t, tt.args, got)) {
				t.Errorf("InterceptReadRequest() = %v", got)
			}
		})
	}
}

func TestSingleItemRequestInterceptor_InterceptWriteRequest(t *testing.T) {
	type fields struct {
		readRequestFactory   func(t *testing.T) readRequestFactory
		writeRequestFactory  func(t *testing.T) writeRequestFactory
		readResponseFactory  func(t *testing.T) readResponseFactory
		writeResponseFactory func(t *testing.T) writeResponseFactory
	}
	type args struct {
		ctx          context.Context
		writeRequest apiModel.PlcWriteRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		mockSetup  func(t *testing.T, fields *fields, args *args)
		wantAssert func(t *testing.T, args args, got []apiModel.PlcWriteRequest) bool
	}{
		{
			name: "nil stays nil",
		},
		{
			name: "write request with no tags",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				plcWriteRequest := NewMockPlcWriteRequest(t)
				plcWriteRequest.EXPECT().GetTagNames().Return(nil)
				args.writeRequest = plcWriteRequest
			},
		},
		{
			name: "write request with 1 tag",
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				plcWriteRequest := NewMockPlcWriteRequest(t)
				plcWriteRequest.EXPECT().GetTagNames().Return([]string{"a tag"})
				args.writeRequest = plcWriteRequest
			},
			wantAssert: func(t *testing.T, args args, got []apiModel.PlcWriteRequest) bool {
				return assert.Contains(t, got, args.writeRequest)
			},
		},
		{
			name: "write request with 2 tags",
			fields: fields{
				writeRequestFactory: func(t *testing.T) writeRequestFactory {
					return func(tags map[string]apiModel.PlcTag, tagNames []string, values map[string]values.PlcValue, writer spi.PlcWriter, writeRequestInterceptor WriteRequestInterceptor) apiModel.PlcWriteRequest {
						plcWriteRequest := NewMockPlcWriteRequest(t)
						expect := plcWriteRequest.EXPECT()
						expect.GetTagNames().Return(tagNames)
						expect.GetTag(mock.Anything).RunAndReturn(func(s string) apiModel.PlcTag {
							return tags[s]
						})
						return plcWriteRequest
					}
				},
			},
			args: args{
				ctx: context.Background(),
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				plcWriteRequest := NewMockPlcWriteRequest(t)
				expect := plcWriteRequest.EXPECT()
				expect.GetTagNames().Return([]string{"1 tag", "2 tag"})
				expect.GetTag(mock.Anything).Return(nil)
				expect.GetValue(mock.Anything).Return(nil)
				expect.GetWriter().Return(nil)
				expect.GetWriteRequestInterceptor().Return(nil)
				args.writeRequest = plcWriteRequest
			},
			wantAssert: func(t *testing.T, args args, got []apiModel.PlcWriteRequest) bool {
				assert.Len(t, got, 2)
				assert.Contains(t, got[0].GetTagNames(), "1 tag")
				assert.Nil(t, got[0].GetTag("1 tag"))
				assert.Contains(t, got[1].GetTagNames(), "2 tag")
				assert.Nil(t, got[1].GetTag("2 tag"))
				return true
			},
		},
		{
			name: "write request with 2 tags aborted",
			fields: fields{
				writeRequestFactory: func(t *testing.T) writeRequestFactory {
					return func(tags map[string]apiModel.PlcTag, tagNames []string, values map[string]values.PlcValue, writer spi.PlcWriter, writeRequestInterceptor WriteRequestInterceptor) apiModel.PlcWriteRequest {
						plcWriteRequest := NewMockPlcWriteRequest(t)
						expect := plcWriteRequest.EXPECT()
						expect.GetTagNames().Return(tagNames)
						expect.GetTag(mock.Anything).RunAndReturn(func(s string) apiModel.PlcTag {
							return tags[s]
						})
						return plcWriteRequest
					}
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					return ctx
				}(),
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				plcWriteRequest := NewMockPlcWriteRequest(t)
				plcWriteRequest.EXPECT().GetTagNames().Return([]string{"1 tag", "2 tag"})
				args.writeRequest = plcWriteRequest
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			if tt.fields.readRequestFactory == nil {
				tt.fields.readRequestFactory = func(t *testing.T) readRequestFactory {
					return nil
				}
			}
			if tt.fields.writeRequestFactory == nil {
				tt.fields.writeRequestFactory = func(t *testing.T) writeRequestFactory {
					return nil
				}
			}
			if tt.fields.readResponseFactory == nil {
				tt.fields.readResponseFactory = func(t *testing.T) readResponseFactory {
					return nil
				}
			}
			if tt.fields.writeResponseFactory == nil {
				tt.fields.writeResponseFactory = func(t *testing.T) writeResponseFactory {
					return nil
				}
			}
			if tt.wantAssert == nil {
				tt.wantAssert = func(t *testing.T, args args, got []apiModel.PlcWriteRequest) bool {
					return true
				}
			}
			m := SingleItemRequestInterceptor{
				readRequestFactory:   tt.fields.readRequestFactory(t),
				writeRequestFactory:  tt.fields.writeRequestFactory(t),
				readResponseFactory:  tt.fields.readResponseFactory(t),
				writeResponseFactory: tt.fields.writeResponseFactory(t),
			}
			if got := m.InterceptWriteRequest(tt.args.ctx, tt.args.writeRequest); !assert.True(t, tt.wantAssert(t, tt.args, got)) {
				t.Errorf("InterceptWriteRequest() = %v", got)
			}
		})
	}
}

func TestSingleItemRequestInterceptor_ProcessReadResponses(t *testing.T) {
	type fields struct {
		readRequestFactory   func(t *testing.T) readRequestFactory
		writeRequestFactory  func(t *testing.T) writeRequestFactory
		readResponseFactory  func(t *testing.T) readResponseFactory
		writeResponseFactory func(t *testing.T) writeResponseFactory
	}
	type args struct {
		ctx         context.Context
		readRequest apiModel.PlcReadRequest
		readResults []apiModel.PlcReadRequestResult
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		mockSetup  func(t *testing.T, fields *fields, args *args)
		wantAssert func(t *testing.T, args args, got apiModel.PlcReadRequestResult) bool
	}{
		{
			name: "no results",
			fields: fields{
				readResponseFactory: func(t *testing.T) readResponseFactory {
					return func(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) apiModel.PlcReadResponse {
						return nil
					}
				},
			},
			wantAssert: func(t *testing.T, args args, got apiModel.PlcReadRequestResult) bool {
				return assert.Equal(t, &interceptedPlcReadRequestResult{}, got)
			},
		},
		{
			name: "one result",
			fields: fields{
				readResponseFactory: func(t *testing.T) readResponseFactory {
					return func(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) apiModel.PlcReadResponse {
						return nil
					}
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				args.readResults = []apiModel.PlcReadRequestResult{
					NewMockPlcReadRequestResult(t),
				}
			},
			wantAssert: func(t *testing.T, args args, got apiModel.PlcReadRequestResult) bool {
				return assert.Equal(t, NewMockPlcReadRequestResult(t), got)
			},
		},
		{
			name: "two result (bit empty)",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				readResponseFactory: func(t *testing.T) readResponseFactory {
					return func(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) apiModel.PlcReadResponse {
						return nil
					}
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				result1 := NewMockPlcReadRequestResult(t)
				result1Expect := result1.EXPECT()
				result1Expect.GetResponse().Return(nil)
				result1Expect.GetErr().Return(nil)
				result2 := NewMockPlcReadRequestResult(t)
				result2Expect := result2.EXPECT()
				result2Expect.GetResponse().Return(nil)
				result2Expect.GetErr().Return(nil)
				args.readResults = []apiModel.PlcReadRequestResult{
					result1,
					result2,
				}
			},
			wantAssert: func(t *testing.T, args args, got apiModel.PlcReadRequestResult) bool {
				return assert.Equal(t, &interceptedPlcReadRequestResult{}, got)
			},
		},
		{
			name: "two result",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				readResponseFactory: func(t *testing.T) readResponseFactory {
					return func(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) apiModel.PlcReadResponse {
						return nil
					}
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				result1 := NewMockPlcReadRequestResult(t)
				result1Expect := result1.EXPECT()
				result1Expect.GetErr().Return(errors.New("asd"))
				result2 := NewMockPlcReadRequestResult(t)
				result2Expect := result2.EXPECT()
				result2Expect.GetResponse().Return(nil)
				result2Expect.GetErr().Return(nil)
				args.readResults = []apiModel.PlcReadRequestResult{
					result1,
					result2,
				}
			},
			wantAssert: func(t *testing.T, args args, got apiModel.PlcReadRequestResult) bool {
				return assert.Equal(t, &interceptedPlcReadRequestResult{
					Err: utils.MultiError{MainError: errors.New("while aggregating results"), Errors: []error{errors.New("asd")}},
				}, got)
			},
		},
		{
			name: "two result (canceled)",
			args: args{
				ctx: func() context.Context {
					ctx, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					return ctx
				}(),
			},
			fields: fields{
				readResponseFactory: func(t *testing.T) readResponseFactory {
					return func(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) apiModel.PlcReadResponse {
						return nil
					}
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				args.readResults = []apiModel.PlcReadRequestResult{
					NewMockPlcReadRequestResult(t),
					NewMockPlcReadRequestResult(t),
				}
			},
			wantAssert: func(t *testing.T, args args, got apiModel.PlcReadRequestResult) bool {
				return assert.Equal(t, &interceptedPlcReadRequestResult{
					Err: errors.New("context canceled"),
				}, got)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			if tt.fields.readRequestFactory == nil {
				tt.fields.readRequestFactory = func(t *testing.T) readRequestFactory {
					return nil
				}
			}
			if tt.fields.writeRequestFactory == nil {
				tt.fields.writeRequestFactory = func(t *testing.T) writeRequestFactory {
					return nil
				}
			}
			if tt.fields.readResponseFactory == nil {
				tt.fields.readResponseFactory = func(t *testing.T) readResponseFactory {
					return nil
				}
			}
			if tt.fields.writeResponseFactory == nil {
				tt.fields.writeResponseFactory = func(t *testing.T) writeResponseFactory {
					return nil
				}
			}
			if tt.wantAssert == nil {
				tt.wantAssert = func(t *testing.T, args args, got apiModel.PlcReadRequestResult) bool {
					return true
				}
			}
			m := SingleItemRequestInterceptor{
				readRequestFactory:   tt.fields.readRequestFactory(t),
				writeRequestFactory:  tt.fields.writeRequestFactory(t),
				readResponseFactory:  tt.fields.readResponseFactory(t),
				writeResponseFactory: tt.fields.writeResponseFactory(t),
			}
			if got := m.ProcessReadResponses(tt.args.ctx, tt.args.readRequest, tt.args.readResults); !assert.True(t, tt.wantAssert(t, tt.args, got)) {
				t.Errorf("ProcessReadResponses() = %v", got)
			}
		})
	}
}

func TestSingleItemRequestInterceptor_ProcessWriteResponses(t *testing.T) {
	type fields struct {
		readRequestFactory   readRequestFactory
		writeRequestFactory  writeRequestFactory
		readResponseFactory  readResponseFactory
		writeResponseFactory writeResponseFactory
	}
	type args struct {
		ctx          context.Context
		writeRequest apiModel.PlcWriteRequest
		writeResults []apiModel.PlcWriteRequestResult
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		mockSetup  func(t *testing.T, fields *fields, args *args)
		wantAssert func(t *testing.T, args args, got apiModel.PlcWriteRequestResult) bool
	}{
		{
			name: "no results",
			fields: fields{
				writeResponseFactory: func(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
					return nil
				},
			},
			wantAssert: func(t *testing.T, args args, got apiModel.PlcWriteRequestResult) bool {
				return assert.Equal(t, &interceptedPlcWriteRequestResult{}, got)
			},
		},
		{
			name: "one result",
			fields: fields{
				writeResponseFactory: func(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
					return nil
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				args.writeResults = []apiModel.PlcWriteRequestResult{
					NewMockPlcWriteRequestResult(t),
				}
			},
			wantAssert: func(t *testing.T, args args, got apiModel.PlcWriteRequestResult) bool {
				return assert.Equal(t, NewMockPlcWriteRequestResult(t), got)
			},
		},
		{
			name: "two result (bit empty)",
			fields: fields{
				writeResponseFactory: func(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
					return nil
				},
			},
			args: args{
				ctx: context.Background(),
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				result1 := NewMockPlcWriteRequestResult(t)
				result1Expect := result1.EXPECT()
				result1Expect.GetResponse().Return(nil)
				result1Expect.GetErr().Return(nil)
				result2 := NewMockPlcWriteRequestResult(t)
				result2Expect := result2.EXPECT()
				result2Expect.GetResponse().Return(nil)
				result2Expect.GetErr().Return(nil)
				args.writeResults = []apiModel.PlcWriteRequestResult{
					result1,
					result2,
				}
			},
			wantAssert: func(t *testing.T, args args, got apiModel.PlcWriteRequestResult) bool {
				return assert.Equal(t, &interceptedPlcWriteRequestResult{}, got)
			},
		},
		{
			name: "two result",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				writeResponseFactory: func(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
					return nil
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				result1 := NewMockPlcWriteRequestResult(t)
				result1Expect := result1.EXPECT()
				result1Expect.GetErr().Return(errors.New("asd"))
				result2 := NewMockPlcWriteRequestResult(t)
				result2Expect := result2.EXPECT()
				result2Expect.GetResponse().Return(nil)
				result2Expect.GetErr().Return(nil)
				args.writeResults = []apiModel.PlcWriteRequestResult{
					result1,
					result2,
				}
			},
			wantAssert: func(t *testing.T, args args, got apiModel.PlcWriteRequestResult) bool {
				return assert.Equal(t, &interceptedPlcWriteRequestResult{
					Err: utils.MultiError{MainError: errors.New("while aggregating results"), Errors: []error{errors.New("asd")}},
				}, got)
			},
		},
		{
			name: "two result (canceled)",
			args: args{
				ctx: func() context.Context {
					ctx, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					return ctx
				}(),
			},
			fields: fields{
				writeResponseFactory: func(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
					return nil
				},
			},
			mockSetup: func(t *testing.T, fields *fields, args *args) {
				args.writeResults = []apiModel.PlcWriteRequestResult{
					NewMockPlcWriteRequestResult(t),
					NewMockPlcWriteRequestResult(t),
				}
			},
			wantAssert: func(t *testing.T, args args, want apiModel.PlcWriteRequestResult) bool {
				return assert.Equal(t, &interceptedPlcWriteRequestResult{
					Err: errors.New("context canceled"),
				}, want)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup(t, &tt.fields, &tt.args)
			}
			m := SingleItemRequestInterceptor{
				readRequestFactory:   tt.fields.readRequestFactory,
				writeRequestFactory:  tt.fields.writeRequestFactory,
				readResponseFactory:  tt.fields.readResponseFactory,
				writeResponseFactory: tt.fields.writeResponseFactory,
			}
			if got := m.ProcessWriteResponses(tt.args.ctx, tt.args.writeRequest, tt.args.writeResults); !assert.True(t, tt.wantAssert(t, tt.args, got)) {
				t.Errorf("ProcessWriteResponses() = %v", got)
			}
		})
	}
}

func Test_interceptedPlcReadRequestResult_GetErr(t *testing.T) {
	type fields struct {
		Request  apiModel.PlcReadRequest
		Response apiModel.PlcReadResponse
		Err      error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "just get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &interceptedPlcReadRequestResult{
				Request:  tt.fields.Request,
				Response: tt.fields.Response,
				Err:      tt.fields.Err,
			}
			if err := d.GetErr(); (err != nil) != tt.wantErr {
				t.Errorf("GetErr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_interceptedPlcReadRequestResult_GetRequest(t *testing.T) {
	type fields struct {
		Request  apiModel.PlcReadRequest
		Response apiModel.PlcReadResponse
		Err      error
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcReadRequest
	}{
		{
			name: "just get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &interceptedPlcReadRequestResult{
				Request:  tt.fields.Request,
				Response: tt.fields.Response,
				Err:      tt.fields.Err,
			}
			if got := d.GetRequest(); !assert.Equal(t, tt.want, got) {
				t.Errorf("GetRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_interceptedPlcReadRequestResult_GetResponse(t *testing.T) {
	type fields struct {
		Request  apiModel.PlcReadRequest
		Response apiModel.PlcReadResponse
		Err      error
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcReadResponse
	}{
		{
			name: "just get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &interceptedPlcReadRequestResult{
				Request:  tt.fields.Request,
				Response: tt.fields.Response,
				Err:      tt.fields.Err,
			}
			if got := d.GetResponse(); !assert.Equal(t, tt.want, got) {
				t.Errorf("GetResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_interceptedPlcWriteRequestResult_GetErr(t *testing.T) {
	type fields struct {
		Request  apiModel.PlcWriteRequest
		Response apiModel.PlcWriteResponse
		Err      error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "just get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &interceptedPlcWriteRequestResult{
				Request:  tt.fields.Request,
				Response: tt.fields.Response,
				Err:      tt.fields.Err,
			}
			if err := d.GetErr(); (err != nil) != tt.wantErr {
				t.Errorf("GetErr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_interceptedPlcWriteRequestResult_GetRequest(t *testing.T) {
	type fields struct {
		Request  apiModel.PlcWriteRequest
		Response apiModel.PlcWriteResponse
		Err      error
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcWriteRequest
	}{
		{
			name: "just get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &interceptedPlcWriteRequestResult{
				Request:  tt.fields.Request,
				Response: tt.fields.Response,
				Err:      tt.fields.Err,
			}
			if got := d.GetRequest(); !assert.Equal(t, tt.want, got) {
				t.Errorf("GetRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_interceptedPlcWriteRequestResult_GetResponse(t *testing.T) {
	type fields struct {
		Request  apiModel.PlcWriteRequest
		Response apiModel.PlcWriteResponse
		Err      error
	}
	tests := []struct {
		name   string
		fields fields
		want   apiModel.PlcWriteResponse
	}{
		{
			name: "just get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &interceptedPlcWriteRequestResult{
				Request:  tt.fields.Request,
				Response: tt.fields.Response,
				Err:      tt.fields.Err,
			}
			if got := d.GetResponse(); !assert.Equal(t, tt.want, got) {
				t.Errorf("GetResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
