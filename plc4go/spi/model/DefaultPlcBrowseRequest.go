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

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/pkg/errors"
)

type DefaultPlcBrowseRequestBuilder struct {
	tagHandler spi.PlcTagHandler
	browser    spi.PlcBrowser
	// The double structure is in order to preserve the order of elements.
	queryNames   []string
	queryStrings map[string]string
}

func NewDefaultPlcBrowseRequestBuilder(tagHandler spi.PlcTagHandler, browser spi.PlcBrowser) *DefaultPlcBrowseRequestBuilder {
	return &DefaultPlcBrowseRequestBuilder{
		tagHandler:   tagHandler,
		browser:      browser,
		queryStrings: map[string]string{},
	}
}

func (d *DefaultPlcBrowseRequestBuilder) AddQuery(name string, query string) model.PlcBrowseRequestBuilder {
	d.queryNames = append(d.queryNames, name)
	d.queryStrings[name] = query
	return d
}

func (d *DefaultPlcBrowseRequestBuilder) Build() (model.PlcBrowseRequest, error) {
	queries := map[string]model.PlcQuery{}
	for name, queryString := range d.queryStrings {
		query, err := d.tagHandler.ParseQuery(queryString)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing query: %s", query)
		}
		queries[name] = query
	}
	return NewDefaultPlcBrowseRequest(queries, d.queryNames, d.browser), nil
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcBrowseRequest
type DefaultPlcBrowseRequest struct {
	DefaultPlcRequest
	browser    spi.PlcBrowser
	queryNames []string
	queries    map[string]model.PlcQuery
}

func NewDefaultPlcBrowseRequest(queries map[string]model.PlcQuery, queryNames []string, browser spi.PlcBrowser) model.PlcBrowseRequest {
	return &DefaultPlcBrowseRequest{
		DefaultPlcRequest: DefaultPlcRequest{},
		browser:           browser,
		queryNames:        queryNames,
		queries:           queries,
	}
}

func (d *DefaultPlcBrowseRequest) Execute() <-chan model.PlcBrowseRequestResult {
	return d.browser.Browse(context.TODO(), d)
}

func (d *DefaultPlcBrowseRequest) ExecuteWithContext(ctx context.Context) <-chan model.PlcBrowseRequestResult {
	return d.browser.Browse(ctx, d)
}

func (d *DefaultPlcBrowseRequest) ExecuteWithInterceptor(interceptor func(result model.PlcBrowseItem) bool) <-chan model.PlcBrowseRequestResult {
	return d.ExecuteWithInterceptorWithContext(context.TODO(), interceptor)
}

func (d *DefaultPlcBrowseRequest) ExecuteWithInterceptorWithContext(ctx context.Context, interceptor func(result model.PlcBrowseItem) bool) <-chan model.PlcBrowseRequestResult {
	return d.browser.BrowseWithInterceptor(ctx, d, interceptor)
}

func (d *DefaultPlcBrowseRequest) GetQueryNames() []string {
	return d.queryNames
}

func (d *DefaultPlcBrowseRequest) GetQuery(queryName string) model.PlcQuery {
	return d.queries[queryName]
}
