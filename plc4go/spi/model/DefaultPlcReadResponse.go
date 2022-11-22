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
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcReadResponse
type DefaultPlcReadResponse struct {
	DefaultResponse
	request model.PlcReadRequest
	values  map[string]values.PlcValue
}

func NewDefaultPlcReadResponse(request model.PlcReadRequest, responseCodes map[string]model.PlcResponseCode, values map[string]values.PlcValue) model.PlcReadResponse {
	return &DefaultPlcReadResponse{
		DefaultResponse: NewDefaultResponse(responseCodes),
		request:         request,
		values:          values,
	}
}

func (d *DefaultPlcReadResponse) GetTagNames() []string {
	var tagNames []string
	// We take the tag names from the request to keep order as map is not ordered
	for _, name := range d.request.GetTagNames() {
		if _, ok := d.values[name]; ok {
			tagNames = append(tagNames, name)
		}
	}
	return tagNames
}

func (d *DefaultPlcReadResponse) GetRequest() model.PlcReadRequest {
	return d.request
}

func (d *DefaultPlcReadResponse) GetValue(name string) values.PlcValue {
	return d.values[name]
}
