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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
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

func (d DefaultPlcBrowseResponse) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcBrowseResponse"); err != nil {
		return err
	}

	if serializableRequest, ok := d.request.(utils.Serializable); ok {
		if err := serializableRequest.Serialize(writeBuffer); err != nil {
			return err
		}
	} else {
		return errors.Errorf("Error serializing. Request %T doesn't implement Serializable", d.request)
	}

	if err := writeBuffer.PushContext("results"); err != nil {
		return err
	}
	for fieldName, foundFields := range d.results {
		if err := writeBuffer.PushContext(fieldName); err != nil {
			return err
		}
		for _, field := range foundFields {
			if serializableField, ok := field.(utils.Serializable); ok {
				if err := serializableField.Serialize(writeBuffer); err != nil {
					return err
				}
			} else {
				return errors.Errorf("Error serializing. Field %T doesn't implement Serializable", field)
			}
		}
		if err := writeBuffer.PopContext(fieldName); err != nil {
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
	if err := writeBuffer.WriteSerializable(d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
