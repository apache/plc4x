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

package values

import (
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	api "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/pkg/errors"
	"strings"
)

type PlcStruct struct {
	values map[string]api.PlcValue
	PlcValueAdapter
}

func NewPlcStruct(value map[string]api.PlcValue) api.PlcValue {
	return PlcStruct{
		values: value,
	}
}

func (m PlcStruct) IsStruct() bool {
	return true
}

func (m PlcStruct) GetKeys() []string {
	var keys []string
	for k := range m.values {
		keys = append(keys, k)
	}
	return keys
}

func (m PlcStruct) HasKey(key string) bool {
	if _, ok := m.values[key]; ok {
		return true
	}
	return false
}

func (m PlcStruct) GetValue(key string) api.PlcValue {
	if value, ok := m.values[key]; ok {
		return value
	}
	return nil
}

func (m PlcStruct) GetStruct() map[string]api.PlcValue {
	return m.values
}

func (m PlcStruct) GetString() string {
	var sb strings.Builder
	sb.WriteString("PlcStruct{\n")
	for fieldName, fieldValue := range m.values {
		sb.WriteString("  ")
		sb.WriteString(fieldName)
		sb.WriteString(": \"")
		sb.WriteString(fieldValue.GetString())
		sb.WriteString("\"\n")
	}
	sb.WriteString("}")
	return sb.String()
}

func (m PlcStruct) Serialize(writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcStruct"); err != nil {
		return err
	}
	for fieldName, plcValue := range m.values {
		if err := writeBuffer.PushContext(fieldName); err != nil {
			return err
		}

		if serializablePlcValue, ok := plcValue.(utils.Serializable); ok {
			if err := serializablePlcValue.Serialize(writeBuffer); err != nil {
				return err
			}
		} else {
			return errors.Errorf("Error serializing. %T doesn't implement Serializable", plcValue)
		}

		if err := writeBuffer.PopContext(fieldName); err != nil {
			return err
		}
	}
	return writeBuffer.PopContext("PlcStruct")
}
