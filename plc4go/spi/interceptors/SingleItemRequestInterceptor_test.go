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

// TODO: replace with mock
type _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest struct {
	_String                    string
	_IsAPlcMessage             bool
	_Execute                   func() <-chan apiModel.PlcReadRequestResult
	_ExecuteWithContext        func(ctx context.Context) <-chan apiModel.PlcReadRequestResult
	_GetTagNames               []string
	_GetTag                    map[string]apiModel.PlcTag
	_GetReader                 spi.PlcReader
	_GetReadRequestInterceptor ReadRequestInterceptor
}

func (t _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest) String() string {
	return t._String
}

func (t _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest) IsAPlcMessage() bool {
	return t._IsAPlcMessage
}

func (t _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest) Execute() <-chan apiModel.PlcReadRequestResult {
	return t._Execute()
}

func (t _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest) ExecuteWithContext(ctx context.Context) <-chan apiModel.PlcReadRequestResult {
	return t._ExecuteWithContext(ctx)
}

func (t _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest) GetTagNames() []string {
	return t._GetTagNames
}

func (t _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest) GetTag(tagName string) apiModel.PlcTag {
	return t._GetTag[tagName]
}

func (t _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest) GetReader() spi.PlcReader {
	return t._GetReader
}

func (t _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest) GetReadRequestInterceptor() ReadRequestInterceptor {
	return t._GetReadRequestInterceptor
}

func TestSingleItemRequestInterceptor_InterceptReadRequest(t *testing.T) {
	type fields struct {
		readRequestFactory   readRequestFactory
		writeRequestFactory  writeRequestFactory
		readResponseFactory  readResponseFactory
		writeResponseFactory writeResponseFactory
	}
	type args struct {
		ctx         context.Context
		readRequest apiModel.PlcReadRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []apiModel.PlcReadRequest
	}{
		{
			name: "nil stays nil",
		},
		{
			name: "read request with no tags",
			args: args{
				readRequest: _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{},
			},
		},
		{
			name: "read request with 1 tag",
			args: args{
				readRequest: _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{
					_GetTagNames: []string{"a tag"},
				},
			},
			want: []apiModel.PlcReadRequest{
				_TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{
					_GetTagNames: []string{"a tag"},
				},
			},
		},
		{
			name: "read request with 2 tags",
			fields: fields{
				readRequestFactory: func(tags map[string]apiModel.PlcTag, tagNames []string, _ spi.PlcReader, _ ReadRequestInterceptor) apiModel.PlcReadRequest {
					return _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{
						_GetTagNames: tagNames,
						_GetTag:      tags,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				readRequest: _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{
					_GetTagNames: []string{"1 tag", "2 tag"},
				},
			},
			want: []apiModel.PlcReadRequest{
				_TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{
					_GetTagNames: []string{"1 tag"},
					_GetTag:      map[string]apiModel.PlcTag{"1 tag": nil},
				}, _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{
					_GetTagNames: []string{"2 tag"},
					_GetTag:      map[string]apiModel.PlcTag{"2 tag": nil},
				},
			},
		},
		{
			name: "read request with 2 tags aborted",
			fields: fields{
				readRequestFactory: func(tags map[string]apiModel.PlcTag, tagNames []string, _ spi.PlcReader, _ ReadRequestInterceptor) apiModel.PlcReadRequest {
					return _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{
						_GetTagNames: tagNames,
						_GetTag:      tags,
					}
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					return ctx
				}(),
				readRequest: _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{
					_GetTagNames: []string{"1 tag", "2 tag"},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := SingleItemRequestInterceptor{
				readRequestFactory:   tt.fields.readRequestFactory,
				writeRequestFactory:  tt.fields.writeRequestFactory,
				readResponseFactory:  tt.fields.readResponseFactory,
				writeResponseFactory: tt.fields.writeResponseFactory,
			}
			if got := m.InterceptReadRequest(tt.args.ctx, tt.args.readRequest); !assert.Equal(t, tt.want, got) {
				t.Errorf("InterceptReadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: replace with mock
type _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest struct {
	_String                     string
	_IsAPlcMessage              bool
	_Execute                    func() <-chan apiModel.PlcWriteRequestResult
	_ExecuteWithContext         func(ctx context.Context) <-chan apiModel.PlcWriteRequestResult
	_GetTagNames                []string
	_GetTag                     map[string]apiModel.PlcTag
	_GetValue                   map[string]values.PlcValue
	_GetWriter                  spi.PlcWriter
	_GetWriteRequestInterceptor WriteRequestInterceptor
}

func (t _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest) String() string {
	return t._String
}

func (t _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest) IsAPlcMessage() bool {
	return t._IsAPlcMessage
}

func (t _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest) Execute() <-chan apiModel.PlcWriteRequestResult {
	return t._Execute()
}

func (t _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest) ExecuteWithContext(ctx context.Context) <-chan apiModel.PlcWriteRequestResult {
	return t._ExecuteWithContext(ctx)
}

func (t _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest) GetTagNames() []string {
	return t._GetTagNames
}

func (t _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest) GetTag(tagName string) apiModel.PlcTag {
	return t._GetTag[tagName]
}

func (t _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest) GetValue(tagName string) values.PlcValue {
	return t._GetValue[tagName]
}

func (t _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest) GetWriter() spi.PlcWriter {
	return t._GetWriter
}

func (t _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest) GetWriteRequestInterceptor() WriteRequestInterceptor {
	return t._GetWriteRequestInterceptor
}

func TestSingleItemRequestInterceptor_InterceptWriteRequest(t *testing.T) {
	type fields struct {
		readRequestFactory   readRequestFactory
		writeRequestFactory  writeRequestFactory
		readResponseFactory  readResponseFactory
		writeResponseFactory writeResponseFactory
	}
	type args struct {
		ctx          context.Context
		writeRequest apiModel.PlcWriteRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []apiModel.PlcWriteRequest
	}{
		{
			name: "nil stays nil",
		},
		{
			name: "read request with no tags",
			args: args{
				writeRequest: _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{},
			},
		},
		{
			name: "read request with 1 tag",
			args: args{
				writeRequest: _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{
					_GetTagNames: []string{"a tag"},
				},
			},
			want: []apiModel.PlcWriteRequest{
				_TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{
					_GetTagNames: []string{"a tag"},
				},
			},
		},
		{
			name: "read request with 2 tags",
			fields: fields{
				writeRequestFactory: func(tags map[string]apiModel.PlcTag, tagNames []string, values map[string]values.PlcValue, writer spi.PlcWriter, writeRequestInterceptor WriteRequestInterceptor) apiModel.PlcWriteRequest {
					return _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{
						_GetTagNames: tagNames,
						_GetTag:      tags,
					}
				},
			},
			args: args{
				ctx: context.Background(),
				writeRequest: _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{
					_GetTagNames: []string{"1 tag", "2 tag"},
				},
			},
			want: []apiModel.PlcWriteRequest{
				_TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{
					_GetTagNames: []string{"1 tag"},
					_GetTag:      map[string]apiModel.PlcTag{"1 tag": nil},
				}, _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{
					_GetTagNames: []string{"2 tag"},
					_GetTag:      map[string]apiModel.PlcTag{"2 tag": nil},
				},
			},
		},
		{
			name: "read request with 2 tags aborted",
			fields: fields{
				writeRequestFactory: func(tags map[string]apiModel.PlcTag, tagNames []string, values map[string]values.PlcValue, writer spi.PlcWriter, writeRequestInterceptor WriteRequestInterceptor) apiModel.PlcWriteRequest {
					return _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{
						_GetTagNames: tagNames,
						_GetTag:      tags,
					}
				},
			},
			args: args{
				ctx: func() context.Context {
					ctx, cancelFunc := context.WithCancel(context.Background())
					cancelFunc()
					return ctx
				}(),
				writeRequest: _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{
					_GetTagNames: []string{"1 tag", "2 tag"},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := SingleItemRequestInterceptor{
				readRequestFactory:   tt.fields.readRequestFactory,
				writeRequestFactory:  tt.fields.writeRequestFactory,
				readResponseFactory:  tt.fields.readResponseFactory,
				writeResponseFactory: tt.fields.writeResponseFactory,
			}
			if got := m.InterceptWriteRequest(tt.args.ctx, tt.args.writeRequest); !assert.Equal(t, tt.want, got) {
				t.Errorf("InterceptWriteRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: replace with mock
type _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult struct {
	_GetRequest  apiModel.PlcReadRequest
	_GetResponse apiModel.PlcReadResponse
	_GetErr      error
}

func (t _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult) GetRequest() apiModel.PlcReadRequest {
	return t._GetRequest
}

func (t _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult) GetResponse() apiModel.PlcReadResponse {
	return t._GetResponse
}

func (t _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult) GetErr() error {
	return t._GetErr
}

// TODO: replace with mock
type _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResponse struct {
	_String          string
	_GetRequest      apiModel.PlcReadRequest
	_GetTagNames     []string
	_GetResponseCode map[string]apiModel.PlcResponseCode
	_GetValue        map[string]values.PlcValue
}

func (t _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResponse) String() string {
	return t._String
}

func (t _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResponse) IsAPlcMessage() bool {
	return true
}

func (t _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResponse) GetRequest() apiModel.PlcReadRequest {
	return t._GetRequest
}

func (t _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResponse) GetTagNames() []string {
	return t._GetTagNames
}

func (t _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResponse) GetResponseCode(tagName string) apiModel.PlcResponseCode {
	return t._GetResponseCode[tagName]
}

func (t _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResponse) GetValue(tagName string) values.PlcValue {
	return t._GetValue[tagName]
}

func TestSingleItemRequestInterceptor_ProcessReadResponses(t *testing.T) {
	type fields struct {
		readRequestFactory   readRequestFactory
		writeRequestFactory  writeRequestFactory
		readResponseFactory  readResponseFactory
		writeResponseFactory writeResponseFactory
	}
	type args struct {
		ctx         context.Context
		readRequest apiModel.PlcReadRequest
		readResults []apiModel.PlcReadRequestResult
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcReadRequestResult
	}{
		{
			name: "no results",
			fields: fields{
				readResponseFactory: func(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) apiModel.PlcReadResponse {
					return nil
				},
			},
			want: &interceptedPlcReadRequestResult{},
		},
		{
			name: "one result",
			args: args{
				readResults: []apiModel.PlcReadRequestResult{
					_TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult{},
				},
			},
			fields: fields{
				readResponseFactory: func(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) apiModel.PlcReadResponse {
					return nil
				},
			},
			want: _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult{},
		},
		{
			name: "two result (bit empty)",
			args: args{
				ctx: context.Background(),
				readResults: []apiModel.PlcReadRequestResult{
					_TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult{},
					_TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult{},
				},
			},
			fields: fields{
				readResponseFactory: func(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) apiModel.PlcReadResponse {
					return nil
				},
			},
			want: &interceptedPlcReadRequestResult{},
		},
		{
			name: "two result",
			args: args{
				ctx: context.Background(),
				readResults: []apiModel.PlcReadRequestResult{
					_TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult{
						_GetErr: errors.New("asd"),
					},
					_TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult{
						_GetRequest: _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{},
						_GetResponse: _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResponse{
							_GetRequest: _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{},
						},
					},
				},
			},
			fields: fields{
				readResponseFactory: func(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) apiModel.PlcReadResponse {
					return nil
				},
			},
			want: &interceptedPlcReadRequestResult{
				Err: utils.MultiError{MainError: errors.New("while aggregating results"), Errors: []error{errors.New("asd")}},
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
				readResults: []apiModel.PlcReadRequestResult{
					_TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult{
						_GetErr: errors.New("asd"),
					},
					_TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResult{
						_GetRequest: _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{},
						_GetResponse: _TestSingleItemRequestInterceptor_ProcessReadResponses_ReadResponse{
							_GetRequest: _TestSingleItemRequestInterceptor_InterceptReadRequestPlcReadRequest{},
						},
					},
				},
			},
			fields: fields{
				readResponseFactory: func(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) apiModel.PlcReadResponse {
					return nil
				},
			},
			want: &interceptedPlcReadRequestResult{
				Err: errors.New("context canceled"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := SingleItemRequestInterceptor{
				readRequestFactory:   tt.fields.readRequestFactory,
				writeRequestFactory:  tt.fields.writeRequestFactory,
				readResponseFactory:  tt.fields.readResponseFactory,
				writeResponseFactory: tt.fields.writeResponseFactory,
			}
			if got := m.ProcessReadResponses(tt.args.ctx, tt.args.readRequest, tt.args.readResults); !assert.Equal(t, tt.want, got) {
				t.Errorf("ProcessReadResponses() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: replace wioth mock
type _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult struct {
	_GetRequest  apiModel.PlcWriteRequest
	_GetResponse apiModel.PlcWriteResponse
	_GetErr      error
}

func (t _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult) GetRequest() apiModel.PlcWriteRequest {
	return t._GetRequest
}

func (t _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult) GetResponse() apiModel.PlcWriteResponse {
	return t._GetResponse
}

func (t _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult) GetErr() error {
	return t._GetErr
}

// TODO: replace wioth mock
type _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResponse struct {
	_String          string
	_IsAPlcMessage   bool
	_GetRequest      apiModel.PlcWriteRequest
	_GetTagNames     []string
	_GetResponseCode map[string]apiModel.PlcResponseCode
}

func (t _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResponse) String() string {
	return t._String
}

func (t _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResponse) IsAPlcMessage() bool {
	return t._IsAPlcMessage
}

func (t _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResponse) GetRequest() apiModel.PlcWriteRequest {
	return t._GetRequest
}

func (t _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResponse) GetTagNames() []string {
	return t._GetTagNames
}

func (t _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResponse) GetResponseCode(tagName string) apiModel.PlcResponseCode {
	return t._GetResponseCode[tagName]
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
		name   string
		fields fields
		args   args
		want   apiModel.PlcWriteRequestResult
	}{
		{
			name: "no results",
			fields: fields{
				writeResponseFactory: func(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
					return nil
				},
			},
			want: &interceptedPlcWriteRequestResult{},
		},
		{
			name: "one result",
			args: args{
				writeResults: []apiModel.PlcWriteRequestResult{
					_TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult{},
				},
			},
			fields: fields{
				writeResponseFactory: func(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
					return nil
				},
			},
			want: _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult{},
		},
		{
			name: "two result (bit empty)",
			args: args{
				ctx: context.Background(),
				writeResults: []apiModel.PlcWriteRequestResult{
					_TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult{},
					_TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult{},
				},
			},
			fields: fields{
				writeResponseFactory: func(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
					return nil
				},
			},
			want: &interceptedPlcWriteRequestResult{},
		},
		{
			name: "two result",
			args: args{
				ctx: context.Background(),
				writeResults: []apiModel.PlcWriteRequestResult{
					_TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult{
						_GetErr: errors.New("asd"),
					},
					_TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult{
						_GetRequest: _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{},
						_GetResponse: _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResponse{
							_GetRequest: _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{},
						},
					},
				},
			},
			fields: fields{
				writeResponseFactory: func(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
					return nil
				},
			},
			want: &interceptedPlcWriteRequestResult{
				Err: utils.MultiError{MainError: errors.New("while aggregating results"), Errors: []error{errors.New("asd")}},
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
				writeResults: []apiModel.PlcWriteRequestResult{
					_TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult{
						_GetErr: errors.New("asd"),
					},
					_TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResult{
						_GetRequest: _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{},
						_GetResponse: _TestSingleItemRequestInterceptor_ProcessWriteResponses_WriteResponse{
							_GetRequest: _TestSingleItemRequestInterceptor_InterceptWriteRequestPlcWriteRequest{},
						},
					},
				},
			},
			fields: fields{
				writeResponseFactory: func(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
					return nil
				},
			},
			want: &interceptedPlcWriteRequestResult{
				Err: errors.New("context canceled"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := SingleItemRequestInterceptor{
				readRequestFactory:   tt.fields.readRequestFactory,
				writeRequestFactory:  tt.fields.writeRequestFactory,
				readResponseFactory:  tt.fields.readResponseFactory,
				writeResponseFactory: tt.fields.writeResponseFactory,
			}
			if got := m.ProcessWriteResponses(tt.args.ctx, tt.args.writeRequest, tt.args.writeResults); !assert.Equal(t, tt.want, got) {
				t.Errorf("ProcessWriteResponses() = %v, want %v", got, tt.want)
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
