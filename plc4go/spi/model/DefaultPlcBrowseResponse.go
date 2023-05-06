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
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcBrowseResponse
type DefaultPlcBrowseResponse struct {
	request      apiModel.PlcBrowseRequest
	responseCode apiModel.PlcResponseCode `stringer:"true"`
	results      map[string]*DefaultPlcBrowseResponseItem
}

func NewDefaultPlcBrowseResponse(request apiModel.PlcBrowseRequest, results map[string][]apiModel.PlcBrowseItem, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcBrowseResponse {
	res := map[string]*DefaultPlcBrowseResponseItem{}
	for name, code := range responseCodes {
		value := results[name]
		res[name] = NewBrowseResponseItem(code, value)
	}
	return &DefaultPlcBrowseResponse{
		request: request,
		results: res,
	}
}

func (d *DefaultPlcBrowseResponse) IsAPlcMessage() bool {
	return true
}

func (d *DefaultPlcBrowseResponse) GetRequest() apiModel.PlcBrowseRequest {
	return d.request
}

func (d *DefaultPlcBrowseResponse) GetQueryNames() []string {
	var queryNames []string
	for queryName := range d.results {
		queryNames = append(queryNames, queryName)
	}
	return queryNames
}

func (d *DefaultPlcBrowseResponse) GetResponseCode(name string) apiModel.PlcResponseCode {
	return d.results[name].GetCode()
}

func (d *DefaultPlcBrowseResponse) GetQueryResults(queryName string) []apiModel.PlcBrowseItem {
	return d.results[queryName].GetResults()
}
