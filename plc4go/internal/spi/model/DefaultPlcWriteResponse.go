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
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
)

type DefaultPlcWriteResponse struct {
	DefaultResponse
	request model.PlcWriteRequest
}

func NewDefaultPlcWriteResponse(request model.PlcWriteRequest, responseCodes map[string]model.PlcResponseCode) model.PlcWriteResponse {
	return DefaultPlcWriteResponse{
		DefaultResponse: NewDefaultResponse(responseCodes),
		request:         request,
	}
}

func (m DefaultPlcWriteResponse) GetFieldNames() []string {
	var fieldNames []string
	// We take the field names from the request to keep order as map is not ordered
	for _, name := range m.request.GetFieldNames() {
		if _, ok := m.responseCodes[name]; ok {
			fieldNames = append(fieldNames, name)
		}
	}
	return fieldNames
}

func (m DefaultPlcWriteResponse) GetRequest() model.PlcWriteRequest {
	return m.request
}

func (m DefaultPlcWriteResponse) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcWriteResponse"); err != nil {
		return err
	}

	if serializableRequest, ok := m.request.(utils.Serializable); ok {
		if err := serializableRequest.Serialize(writeBuffer); err != nil {
			return err
		}
	}

	if err := writeBuffer.PushContext("fields"); err != nil {
		return err
	}
	for _, fieldName := range m.GetFieldNames() {
		responseCodeName := m.GetResponseCode(fieldName).GetName()
		if err := writeBuffer.WriteString(fieldName, uint32(len([]rune(responseCodeName))*8), "UTF-8", responseCodeName); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("fields"); err != nil {
		return err
	}

	if err := writeBuffer.PopContext("PlcWriteResponse"); err != nil {
		return err
	}
	return nil
}

func (m DefaultPlcWriteResponse) String() string {
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
