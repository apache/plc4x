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
)

type DefaultPlcBrowseResponse struct {
	DefaultResponse
	request model.PlcBrowseRequest
	results map[string][]model.PlcBrowseFoundField
}

func NewDefaultPlcBrowseResponse(request model.PlcBrowseRequest, results map[string][]model.PlcBrowseFoundField, responseCodes map[string]model.PlcResponseCode) DefaultPlcBrowseResponse {
	return DefaultPlcBrowseResponse{
		DefaultResponse: DefaultResponse{responseCodes: responseCodes},
		request:         request,
		results:         results,
	}
}

func (d DefaultPlcBrowseResponse) GetFieldNames() []string {
	var fieldNames []string
	// We take the field names from the request to keep order as map is not ordered
	for _, name := range d.request.GetFieldNames() {
		if _, ok := d.results[name]; ok {
			fieldNames = append(fieldNames, name)
		}
	}
	return fieldNames
}

func (d DefaultPlcBrowseResponse) GetRequest() model.PlcBrowseRequest {
	return d.request
}

func (d DefaultPlcBrowseResponse) GetQueryNames() []string {
	var queryNames []string
	for queryName := range d.results {
		queryNames = append(queryNames, queryName)
	}
	return queryNames
}

func (d DefaultPlcBrowseResponse) GetQueryResults(queryName string) []model.PlcBrowseFoundField {
	return d.results[queryName]
}
