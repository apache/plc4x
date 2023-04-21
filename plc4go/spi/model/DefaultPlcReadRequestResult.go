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

import apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcReadRequestResult
type DefaultPlcReadRequestResult struct {
	Request  apiModel.PlcReadRequest
	Response apiModel.PlcReadResponse
	Err      error
}

func NewDefaultPlcReadRequestResult(Request apiModel.PlcReadRequest, Response apiModel.PlcReadResponse, Err error) apiModel.PlcReadRequestResult {
	return &DefaultPlcReadRequestResult{Request, Response, Err}
}

func (d *DefaultPlcReadRequestResult) GetRequest() apiModel.PlcReadRequest {
	return d.Request
}

func (d *DefaultPlcReadRequestResult) GetResponse() apiModel.PlcReadResponse {
	return d.Response
}

func (d *DefaultPlcReadRequestResult) GetErr() error {
	return d.Err
}
