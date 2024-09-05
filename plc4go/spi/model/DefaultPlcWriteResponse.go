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

var _ apiModel.PlcWriteResponse = &DefaultPlcWriteResponse{}

//go:generate plc4xGenerator -type=DefaultPlcWriteResponse
type DefaultPlcWriteResponse struct {
	request       apiModel.PlcWriteRequest
	responseCodes map[string]apiModel.PlcResponseCode
}

func NewDefaultPlcWriteResponse(request apiModel.PlcWriteRequest, responseCodes map[string]apiModel.PlcResponseCode) apiModel.PlcWriteResponse {
	return &DefaultPlcWriteResponse{
		request:       request,
		responseCodes: responseCodes,
	}
}

func (d *DefaultPlcWriteResponse) IsAPlcMessage() bool {
	return true
}

func (d *DefaultPlcWriteResponse) GetTagNames() []string {
	if d.request == nil {
		// Safety guard
		return nil
	}
	var tagNames []string
	// We take the tag names from the request to keep order as map is not ordered
	for _, name := range d.request.GetTagNames() {
		if _, ok := d.responseCodes[name]; ok {
			tagNames = append(tagNames, name)
		}
	}
	return tagNames
}

func (d *DefaultPlcWriteResponse) GetRequest() apiModel.PlcWriteRequest {
	return d.request
}

func (d *DefaultPlcWriteResponse) GetResponseCode(name string) apiModel.PlcResponseCode {
	code, ok := d.responseCodes[name]
	if !ok {
		return apiModel.PlcResponseCode_NOT_FOUND
	}
	return code
}
