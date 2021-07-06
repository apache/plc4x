//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package model

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	values2 "github.com/apache/plc4x/plc4go/internal/plc4go/spi/values"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
)

type DefaultPlcReadResponse struct {
	DefaultResponse
	request model.PlcReadRequest
	values  map[string]values.PlcValue
}

func NewDefaultPlcReadResponse(request model.PlcReadRequest, responseCodes map[string]model.PlcResponseCode, values map[string]values.PlcValue) model.PlcReadResponse {
	return DefaultPlcReadResponse{
		DefaultResponse: NewDefaultResponse(responseCodes),
		request:         request,
		values:          values,
	}
}

func (m DefaultPlcReadResponse) GetFieldNames() []string {
	var fieldNames []string
	// We take the field names from the request to keep order as map is not ordered
	for _, name := range m.request.GetFieldNames() {
		if _, ok := m.values[name]; ok {
			fieldNames = append(fieldNames, name)
		}
	}
	return fieldNames
}

func (m DefaultPlcReadResponse) GetRequest() model.PlcReadRequest {
	return m.request
}

func (m DefaultPlcReadResponse) GetValue(name string) values.PlcValue {
	return m.values[name]
}

func (m DefaultPlcReadResponse) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcReadResponse"); err != nil {
		return err
	}

	if request, ok := m.request.(utils.Serializable); ok {
		if err := request.Serialize(writeBuffer); err != nil {
			return err
		}
	}
	if err := writeBuffer.PushContext("values"); err != nil {
		return err
	}
	for _, fieldName := range m.GetFieldNames() {
		if err := writeBuffer.PushContext(fieldName); err != nil {
			return err
		}
		if err := writeBuffer.PushContext("ResponseItem"); err != nil {
			return err
		}
		codeName := m.GetResponseCode(fieldName).GetName()
		if err := writeBuffer.WriteString("result", uint8(len([]rune(codeName))*8), "UTF-8", codeName); err != nil {
			return err
		}

		valueResponse := m.GetValue(fieldName)
		if _, ok := valueResponse.(values2.PlcNULL); ok {
			// We ignore nulls
			if err := writeBuffer.PopContext("ResponseItem"); err != nil {
				return err
			}
			if err := writeBuffer.PopContext(fieldName); err != nil {
				return err
			}
			continue
		}
		if err := valueResponse.(utils.Serializable).Serialize(writeBuffer); err != nil {
			return err
		}
		if err := writeBuffer.PopContext("ResponseItem"); err != nil {
			return err
		}
		if err := writeBuffer.PopContext(fieldName); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("values"); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("PlcReadResponse"); err != nil {
		return err
	}
	return nil
}
