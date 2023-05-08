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
		want   apiModel.PlcBrowseRequestBuilder
	}{
		{
			name: "add one",
			fields: fields{
				tagHandler:   NewMockPlcTagHandler(t),
				queryStrings: map[string]string{},
			},
			args: args{
				name:  "a name",
				query: "a query",
			},
			want: &DefaultPlcBrowseRequestBuilder{
				tagHandler: NewMockPlcTagHandler(t),
				queryStrings: map[string]string{
					"a name": "a query",
				},
				queryNames: []string{"a name"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		want    apiModel.PlcBrowseRequest
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		want   <-chan apiModel.PlcBrowseRequestResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		want   <-chan apiModel.PlcBrowseRequestResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		want   <-chan apiModel.PlcBrowseRequestResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcBrowseRequest{
				browser:    tt.fields.browser,
				queryNames: tt.fields.queryNames,
				queries:    tt.fields.queries,
			}
			assert.Equalf(t, tt.want, d.ExecuteWithInterceptor(tt.args.interceptor), "ExecuteWithInterceptor(%v)", tt.args.interceptor)
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
		want   <-chan apiModel.PlcBrowseRequestResult
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultPlcBrowseRequest{
				browser:    tt.fields.browser,
				queryNames: tt.fields.queryNames,
				queries:    tt.fields.queries,
			}
			assert.Equalf(t, tt.want, d.ExecuteWithInterceptorWithContext(tt.args.ctx, tt.args.interceptor), "ExecuteWithInterceptorWithContext(%v, %v)", tt.args.ctx, tt.args.interceptor)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewDefaultPlcBrowseRequestBuilder(tt.args.tagHandler, tt.args.browser), "NewDefaultPlcBrowseRequestBuilder(%v, %v)", tt.args.tagHandler, tt.args.browser)
		})
	}
}
