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

package spi

import (
	"context"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
)

type PlcBrowser interface {
	// Browse Non-Blocking request, which will return a full result as soon as the operation is finished
	Browse(ctx context.Context, browseRequest model.PlcBrowseRequest) <-chan model.PlcBrowseRequestResult

	// BrowseWithInterceptor Variant of the Browser, which allows immediately intercepting found resources
	// This is ideal, if additional information has to be queried on such found resources
	// and especially for connection-based protocols can reduce the stress on the system
	// and increase throughput. It can also be used for simple filtering.
	// If the interceptor function returns 'true' the result is added to the overall result
	// if it's 'false' is is not.
	BrowseWithInterceptor(ctx context.Context, browseRequest model.PlcBrowseRequest, interceptor func(result model.PlcBrowseEvent) bool) <-chan model.PlcBrowseRequestResult
}
