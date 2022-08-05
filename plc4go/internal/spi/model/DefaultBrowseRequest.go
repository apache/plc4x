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
	"github.com/apache/plc4x/plc4go/internal/spi"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
)

type DefaultPlcBrowseRequestBuilder struct {
	browser spi.PlcBrowser
	queries map[string]string
}

func NewDefaultPlcBrowseRequestBuilder(browser spi.PlcBrowser) *DefaultPlcBrowseRequestBuilder {
	return &DefaultPlcBrowseRequestBuilder{
		browser: browser,
		queries: map[string]string{},
	}
}

func (m *DefaultPlcBrowseRequestBuilder) AddItem(name string, query string) model.PlcBrowseRequestBuilder {
	m.queries[name] = query
	return m
}

func (m *DefaultPlcBrowseRequestBuilder) Build() (model.PlcBrowseRequest, error) {
	queries := m.queries
	return DefaultPlcBrowseRequest{
		queries: queries,
		browser: m.browser,
	}, nil
}

type DefaultPlcBrowseRequest struct {
	DefaultRequest
	queries map[string]string
	browser spi.PlcBrowser
}

func (d DefaultPlcBrowseRequest) GetQueryNames() []string {
	var queryNames []string
	for queryName := range d.queries {
		queryNames = append(queryNames, queryName)
	}
	return queryNames
}

func (d DefaultPlcBrowseRequest) GetQueryString(name string) string {
	return d.queries[name]
}

func (d DefaultPlcBrowseRequest) Execute() <-chan model.PlcBrowseRequestResult {
	return d.browser.Browse(d)
}

func (d DefaultPlcBrowseRequest) ExecuteWithInterceptor(interceptor func(result model.PlcBrowseEvent) bool) <-chan model.PlcBrowseRequestResult {
	return d.browser.BrowseWithInterceptor(d, interceptor)
}
