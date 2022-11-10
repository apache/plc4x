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

type PlcWriteRequestBuilder interface {
	AddTagAddress(tagName string, tagAddress string, value interface{}) PlcWriteRequestBuilder
	AddTag(tagName string, tag PlcTag, value interface{}) PlcWriteRequestBuilder
	Build() (PlcWriteRequest, error)
}

type PlcWriteRequestResult interface {
	GetRequest() PlcWriteRequest
	GetResponse() PlcWriteResponse
	GetErr() error
}

type PlcWriteRequest interface {
	PlcRequest
	Execute() <-chan PlcWriteRequestResult
	ExecuteWithContext(ctx context.Context) <-chan PlcWriteRequestResult

	GetTagNames() []string
	GetTag(tagName string) PlcTag
	GetValue(tagName string) values.PlcValue
}

type PlcWriteResponse interface {
	PlcResponse
	GetRequest() PlcWriteRequest
	GetTagNames() []string
	GetResponseCode(tagName string) PlcResponseCode
}
