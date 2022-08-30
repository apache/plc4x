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
	"fmt"
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DefaultRequest struct {
	fields     map[string]model.PlcField
	fieldNames []string
}

func (m DefaultRequest) IsAPlcMessage() bool {
	return true
}

func NewDefaultRequest(Fields map[string]model.PlcField, FieldNames []string) DefaultRequest {
	return DefaultRequest{Fields, FieldNames}
}

func (m DefaultRequest) GetFieldNames() []string {
	return m.fieldNames
}

func (m DefaultRequest) GetField(name string) model.PlcField {
	if field, ok := m.fields[name]; ok {
		return field
	}
	return nil
}

func (m DefaultRequest) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("Request"); err != nil {
		return err
	}

	if err := writeBuffer.PushContext("fields"); err != nil {
		return err
	}
	for fieldName, field := range m.fields {
		if serializableField, ok := field.(utils.Serializable); ok {
			if err := writeBuffer.PushContext(fieldName); err != nil {
				return err
			}
			if err := serializableField.Serialize(writeBuffer); err != nil {
				return err
			}
			if err := writeBuffer.PopContext(fieldName); err != nil {
				return err
			}
		} else {
			fieldString := fmt.Sprintf("%v", field)
			if err := writeBuffer.WriteString(fieldName, uint32(len(fieldString)*8), "UTF-8", fieldString); err != nil {
				return err
			}
		}
	}
	if err := writeBuffer.PopContext("fields"); err != nil {
		return err
	}
	if err := writeBuffer.PushContext("fieldNames"); err != nil {
		return err
	}
	for _, name := range m.fieldNames {
		if err := writeBuffer.WriteString("value", uint32(len(name)*8), "UTF-8", name); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("fieldNames"); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("Request"); err != nil {
		return err
	}
	return nil
}

func (m DefaultRequest) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
