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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

var _ apiModel.PlcReadResponse = &DefaultPlcReadResponse{}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcReadResponse
type DefaultPlcReadResponse struct {
	request apiModel.PlcReadRequest
	values  map[string]*ResponseItem
}

func NewDefaultPlcReadResponse(request apiModel.PlcReadRequest, responseCodes map[string]apiModel.PlcResponseCode, values map[string]apiValues.PlcValue) apiModel.PlcReadResponse {
	valueMap := map[string]*ResponseItem{}
	for name, code := range responseCodes {
		value := values[name]
		valueMap[name] = NewResponseItem(code, value)
	}
	return &DefaultPlcReadResponse{
		request: request,
		values:  valueMap,
	}
}

func (d *DefaultPlcReadResponse) IsAPlcMessage() bool {
	return true
}

func (d *DefaultPlcReadResponse) GetTagNames() []string {
	if d.request == nil {
		// safety guard
		return nil
	}
	var tagNames []string
	// We take the tag names from the request to keep order as map is not ordered
	for _, name := range d.request.GetTagNames() {
		if _, ok := d.values[name]; ok {
			tagNames = append(tagNames, name)
		}
	}
	return tagNames
}

func (d *DefaultPlcReadResponse) GetRequest() apiModel.PlcReadRequest {
	return d.request
}

func (d *DefaultPlcReadResponse) GetResponseCode(name string) apiModel.PlcResponseCode {
	item, ok := d.values[name]
	if !ok {
		return apiModel.PlcResponseCode_NOT_FOUND
	}
	return item.GetCode()
}

func (d *DefaultPlcReadResponse) GetValue(name string) apiValues.PlcValue {
	item, ok := d.values[name]
	if !ok {
		return spiValues.PlcNull{}
	}
	return item.GetValue()
}
