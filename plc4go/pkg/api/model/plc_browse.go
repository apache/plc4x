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

	"github.com/apache/plc4x/plc4go/pkg/api/values"
)

type PlcBrowseRequestBuilder interface {
	AddQuery(name string, query string) PlcBrowseRequestBuilder
	Build() (PlcBrowseRequest, error)
}

type PlcBrowseRequest interface {
	PlcRequest
	// Execute Will not return until a potential scan is finished and will return all results in one block
	Execute() <-chan PlcBrowseRequestResult
	// ExecuteWithContext is the same as Execute but handles the Context if implemented for Driver
	ExecuteWithContext(ctx context.Context) <-chan PlcBrowseRequestResult
	// ExecuteWithInterceptor Will call the given callback for every found resource
	ExecuteWithInterceptor(interceptor func(result PlcBrowseItem) bool) <-chan PlcBrowseRequestResult
	// ExecuteWithInterceptorWithContext Will call the given callback for every found resource
	ExecuteWithInterceptorWithContext(ctx context.Context, interceptor func(result PlcBrowseItem) bool) <-chan PlcBrowseRequestResult

	GetQueryNames() []string
	GetQuery(queryName string) PlcQuery
}

type PlcBrowseResponse interface {
	PlcResponse
	GetRequest() PlcBrowseRequest
	GetQueryNames() []string
	GetResponseCode(name string) PlcResponseCode
	GetQueryResults(name string) []PlcBrowseItem
}

type PlcBrowseRequestResult interface {
	GetRequest() PlcBrowseRequest
	GetResponse() PlcBrowseResponse
	GetErr() error
}

type PlcBrowseItem interface {
	GetTag() PlcTag

	GetName() string
	IsReadable() bool
	IsWritable() bool
	IsSubscribable() bool
	GetChildren() map[string]PlcBrowseItem
	GetOptions() map[string]values.PlcValue
}
