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
	"encoding/binary"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// TODO: use generator once we figured out how to render results with ast
type DefaultPlcBrowseResponse struct {
	DefaultResponse
	request      model.PlcBrowseRequest
	responseCode model.PlcResponseCode
	results      map[string]*DefaultPlcBrowseResponseItem
}

func NewDefaultPlcBrowseResponse(request model.PlcBrowseRequest, results map[string][]model.PlcBrowseItem, responseCodes map[string]model.PlcResponseCode) DefaultPlcBrowseResponse {
	res := map[string]*DefaultPlcBrowseResponseItem{}
	for name, code := range responseCodes {
		value := results[name]
		res[name] = NewBrowseResponseItem(code, value)
	}
	return DefaultPlcBrowseResponse{
		request: request,
		results: res,
	}
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

func (d DefaultPlcBrowseResponse) GetResponseCode(name string) model.PlcResponseCode {
	return d.results[name].GetCode()
}

func (d DefaultPlcBrowseResponse) GetQueryResults(queryName string) []model.PlcBrowseItem {
	return d.results[queryName].GetResults()
}

func (d DefaultPlcBrowseResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d DefaultPlcBrowseResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcBrowseResponse"); err != nil {
		return err
	}

	if serializableRequest, ok := d.request.(utils.Serializable); ok {
		if err := serializableRequest.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
			return err
		}
	} else {
		return errors.Errorf("Error serializing. Request %T doesn't implement Serializable", d.request)
	}

	if err := writeBuffer.PushContext("results"); err != nil {
		return err
	}
	for tagName, browseResponseItem := range d.results {
		if err := writeBuffer.PushContext(tagName); err != nil {
			return err
		}
		for _, tag := range browseResponseItem.GetResults() {
			if serializableTag, ok := tag.(utils.Serializable); ok {
				if err := serializableTag.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
					return err
				}
			} else {
				return errors.Errorf("Error serializing. Tag %T doesn't implement Serializable", tag)
			}
		}
		if err := writeBuffer.PopContext(tagName); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("results"); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("PlcBrowseResponse"); err != nil {
		return err
	}
	return nil
}

func (d DefaultPlcBrowseResponse) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
