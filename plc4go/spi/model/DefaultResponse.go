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
)

type DefaultResponse struct {
	responseCodes map[string]model.PlcResponseCode
}

func (m DefaultResponse) IsAPlcMessage() bool {
	return true
}

func (m DefaultResponse) GetResponseCode(name string) model.PlcResponseCode {
	return m.responseCodes[name]
}

func NewDefaultResponse(responseCodes map[string]model.PlcResponseCode) DefaultResponse {
	return DefaultResponse{responseCodes}
}

func (m DefaultResponse) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("Response"); err != nil {
		return err
	}

	if err := writeBuffer.PushContext("responseCodes"); err != nil {
		return err
	}
	for fieldName, code := range m.responseCodes {
		if err := writeBuffer.PushContext(fieldName); err != nil {
			return err
		}
		if err := writeBuffer.WriteUint8("code", 8, uint8(code), utils.WithAdditionalStringRepresentation(code.String())); err != nil {
			return err
		}
		if err := writeBuffer.PopContext(fieldName); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("responseCodes"); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("Response"); err != nil {
		return err
	}
	return nil
}

func (m DefaultResponse) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
