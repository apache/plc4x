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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"testing"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPlcBrowseRequestBuilder_AddQuery(t *testing.T) {
	type fields struct {
		tagHandler   spi.PlcTagHandler
		browser      spi.PlcBrowser
		queryNames   []string
		queryStrings map[string]string
	}
	type args struct {
		name  string
		query string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields, args *args, want *apiModel.PlcBrowseRequestBuilder)
		want   apiModel.PlcBrowseRequestBuilder
	}{
		{
			name: "add one",
			fields: fields{
				queryStrings: map[string]string{},
			},
			args: args{
				name:  "a name",
				query: "a query",
			},
			setup: func(t *testing.T, fields *fields, args *args, want *apiModel.PlcBrowseRequestBuilder) {
				fields.tagHandler = NewMockPlcTagHandler(t)
				*want = &DefaultPlcBrowseRequestBuilder{
					tagHandler: NewMockPlcTagHandler(t),
					queryStrings: map[string]string{
						"a name": "a query",
					},
					queryNames: []string{"a name"},
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields, &tt.args, &tt.want)
			}
			d := &DefaultPlcBrowseRequestBuilder{
				tagHandler:   tt.fields.tagHandler,
				browser:      tt.fields.browser,
				queryNames:   tt.fields.queryNames,
				queryStrings: tt.fields.queryStrings,
			}
			assert.Equalf(t, tt.want, d.AddQuery(tt.args.name, tt.args.query), "AddQuery(%v, %v)", tt.args.name, tt.args.query)
		})
	}
}

func TestDefaultPlcBrowseRequestBuilder_Build(t *testing.T) {
	type fields struct {
		tagHandler   spi.PlcTagHandler
		browser      spi.PlcBrowser
		queryNames   []string
		queryStrings map[string]string
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func(t *testing.T, fields *fields)
		want    apiModel.PlcBrowseRequest
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "build it",
			want:    NewDefaultPlcBrowseRequest(map[string]apiModel.PlcQuery{}, nil, nil),
			wantErr: assert.NoError,
		},
		{
			name: "build it with queries",
			fields: fields{
				queryNames:   nil,
				queryStrings: map[string]string{"a": "a"},
			},
			setup: func(t *testing.T, fields *fields) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseQuery(mock.Anything).Return(nil, nil)
				fields.tagHandler = handler
			},
			want:    NewDefaultPlcBrowseRequest(map[string]apiModel.PlcQuery{"a": nil}, nil, nil),
			wantErr: assert.NoError,
		},
		{
			name: "build it parse error",
			fields: fields{
				queryNames:   nil,
				queryStrings: map[string]string{"a": "a"},
			},
			setup: func(t *testing.T, fields *fields) {
				handler := NewMockPlcTagHandler(t)
				handler.EXPECT().ParseQuery(mock.Anything).Return(nil, errors.New("nope"))
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
			d := &DefaultPlcBrowseRequestBuilder{
				tagHandler:   tt.fields.tagHandler,
				browser:      tt.fields.browser,
				queryNames:   tt.fields.queryNames,
				queryStrings: tt.fields.queryStrings,
			}
			got, err := d.Build()
			if !tt.wantErr(t, err, fmt.Sprintf("Build()")) {
				return
			}
			assert.Equalf(t, tt.want, got, "Build()")
		})
	}
}

func TestDefaultPlcBrowseRequest_Execute(t *testing.T) {
	type fields struct {
		browser    spi.PlcBrowser
		queryNames []string
		queries    map[string]apiModel.PlcQuery
	}
	tests := []struct {
		name   string
		fields fields
		setup  func(t *testing.T, fields *fields)
		want   <-chan apiModel.PlcBrowseRequestResult
	}{
		{
			name: "execute it",
			setup: func(t *testing.T, fields *fields) {
				browser := NewMockPlcBrowser(t)
				browser.EXPECT().Browse(mock.Anything, mock.Anything).Return(nil)
				fields.browser = browser
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			d := &DefaultPlcBrowseRequest{
				browser:    tt.fields.browser,
				queryNames: tt.fields.queryNames,
				queries:    tt.fields.queries,
			}
			assert.Equalf(t, tt.want, d.Execute(), "Execute()")
		})
	}
}

func TestDefaultPlcBrowseRequest_ExecuteWithContext(t *testing.T) {
	type fields struct {
		browser    spi.PlcBrowser
		queryNames []string
		queries    map[string]apiModel.PlcQuery
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields)
		want   <-chan apiModel.PlcBrowseRequestResult
	}{
		{
			name: "execute it",
			setup: func(t *testing.T, fields *fields) {
				browser := NewMockPlcBrowser(t)
				browser.EXPECT().Browse(mock.Anything, mock.Anything).Return(nil)
				fields.browser = browser
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			d := &DefaultPlcBrowseRequest{
				browser:    tt.fields.browser,
				queryNames: tt.fields.queryNames,
				queries:    tt.fields.queries,
			}
			assert.Equalf(t, tt.want, d.ExecuteWithContext(tt.args.ctx), "ExecuteWithContext(%v)", tt.args.ctx)
		})
	}
}

func TestDefaultPlcBrowseRequest_ExecuteWithInterceptor(t *testing.T) {
	type fields struct {
		browser    spi.PlcBrowser
		queryNames []string
		queries    map[string]apiModel.PlcQuery
	}
	type args struct {
		interceptor func(result apiModel.PlcBrowseItem) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields)
		want   <-chan apiModel.PlcBrowseRequestResult
	}{
		{
			name: "execute it",
			setup: func(t *testing.T, fields *fields) {
				browser := NewMockPlcBrowser(t)
				browser.EXPECT().BrowseWithInterceptor(mock.Anything, mock.Anything, mock.Anything).Return(nil)
				fields.browser = browser
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			d := &DefaultPlcBrowseRequest{
				browser:    tt.fields.browser,
				queryNames: tt.fields.queryNames,
				queries:    tt.fields.queries,
			}
			assert.Equalf(t, tt.want, d.ExecuteWithInterceptor(tt.args.interceptor), "ExecuteWithInterceptor(func(%t))", tt.args.interceptor != nil)
		})
	}
}

func TestDefaultPlcBrowseRequest_ExecuteWithInterceptorWithContext(t *testing.T) {
	type fields struct {
		browser    spi.PlcBrowser
		queryNames []string
		queries    map[string]apiModel.PlcQuery
	}
	type args struct {
		ctx         context.Context
		interceptor func(result apiModel.PlcBrowseItem) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		setup  func(t *testing.T, fields *fields)
		want   <-chan apiModel.PlcBrowseRequestResult
	}{
		{
			name: "execute it",
			setup: func(t *testing.T, fields *fields) {
				browser := NewMockPlcBrowser(t)
				browser.EXPECT().BrowseWithInterceptor(mock.Anything, mock.Anything, mock.Anything).Return(nil)
				fields.browser = browser
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, &tt.fields)
			}
			d := &DefaultPlcBrowseRequest{
				browser:    tt.fields.browser,
				queryNames: tt.fields.queryNames,
				queries:    tt.fields.queries,
			}
			assert.Equalf(t, tt.want, d.ExecuteWithInterceptorWithContext(tt.args.ctx, tt.args.interceptor), "ExecuteWithInterceptorWithContext(%v, func(%t))", tt.args.ctx, tt.args.interceptor != nil)
		})
	}
}

func TestDefaultPlcBrowseRequest_GetQuery(t *testing.T) {
	type fields struct {
		browser    spi.PlcBrowser
		queryNames []string
		queries    map[string]apiModel.PlcQuery
	}
	type args struct {
		queryName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   apiModel.PlcQuery
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcBrowseRequest{
				browser:    tt.fields.browser,
				queryNames: tt.fields.queryNames,
				queries:    tt.fields.queries,
			}
			assert.Equalf(t, tt.want, d.GetQuery(tt.args.queryName), "GetQuery(%v)", tt.args.queryName)
		})
	}
}

func TestDefaultPlcBrowseRequest_GetQueryNames(t *testing.T) {
	type fields struct {
		browser    spi.PlcBrowser
		queryNames []string
		queries    map[string]apiModel.PlcQuery
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "get it",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcBrowseRequest{
				browser:    tt.fields.browser,
				queryNames: tt.fields.queryNames,
				queries:    tt.fields.queries,
			}
			assert.Equalf(t, tt.want, d.GetQueryNames(), "GetQueryNames()")
		})
	}
}

func TestDefaultPlcBrowseRequest_IsAPlcMessage(t *testing.T) {
	type fields struct {
		browser    spi.PlcBrowser
		queryNames []string
		queries    map[string]apiModel.PlcQuery
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
			d := &DefaultPlcBrowseRequest{
				browser:    tt.fields.browser,
				queryNames: tt.fields.queryNames,
				queries:    tt.fields.queries,
			}
			assert.Equalf(t, tt.want, d.IsAPlcMessage(), "IsAPlcMessage()")
		})
	}
}

func TestNewDefaultPlcBrowseRequest(t *testing.T) {
	type args struct {
		queries    map[string]apiModel.PlcQuery
		queryNames []string
		browser    spi.PlcBrowser
	}
	tests := []struct {
		name string
		args args
		want *DefaultPlcBrowseRequest
	}{
		{
			name: "browse it",
			want: &DefaultPlcBrowseRequest{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcBrowseRequest(tt.args.queries, tt.args.queryNames, tt.args.browser), "NewDefaultPlcBrowseRequest(%v, %v, %v)", tt.args.queries, tt.args.queryNames, tt.args.browser)
		})
	}
}

func TestNewDefaultPlcBrowseRequestBuilder(t *testing.T) {
	type args struct {
		tagHandler spi.PlcTagHandler
		browser    spi.PlcBrowser
	}
	tests := []struct {
		name string
		args args
		want apiModel.PlcBrowseRequestBuilder
	}{
		{
			name: "get it",
			want: &DefaultPlcBrowseRequestBuilder{
				queryStrings: map[string]string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcBrowseRequestBuilder(tt.args.tagHandler, tt.args.browser), "NewDefaultPlcBrowseRequestBuilder(%v, %v)", tt.args.tagHandler, tt.args.browser)
		})
	}
}
