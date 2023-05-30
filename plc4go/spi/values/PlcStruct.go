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
	"context"
	"encoding/binary"
	"fmt"
	"strings"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

type PlcStruct struct {
	values map[string]apiValues.PlcValue
	PlcValueAdapter
}

func NewPlcStruct(value map[string]apiValues.PlcValue) PlcStruct {
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

func (m PlcStruct) GetValue(key string) apiValues.PlcValue {
	if value, ok := m.values[key]; ok {
		return value
	}
	return nil
}

func (m PlcStruct) GetStruct() map[string]apiValues.PlcValue {
	return m.values
}

func (m PlcStruct) IsString() bool {
	return true
}

func (m PlcStruct) GetString() string {
	var sb strings.Builder
	sb.WriteString("PlcStruct{\n")
	for tagName, tagValue := range m.values {
		sb.WriteString("  ")
		sb.WriteString(tagName)
		sb.WriteString(": \"")
		sb.WriteString(tagValue.GetString())
		sb.WriteString("\"\n")
	}
	sb.WriteString("}")
	return sb.String()
}

func (m PlcStruct) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.Struct
}

func (m PlcStruct) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcStruct) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcStruct"); err != nil {
		return err
	}
	for tagName, tagValue := range m.values {
		if err := writeBuffer.PushContext(tagName); err != nil {
			return err
		}

		if serializablePlcValue, ok := tagValue.(utils.Serializable); ok {
			if err := serializablePlcValue.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
		} else {
			return errors.Errorf("Error serializing. %T doesn't implement Serializable", tagValue)
		}

		if err := writeBuffer.PopContext(tagName); err != nil {
			return err
		}
	}
	return writeBuffer.PopContext("PlcStruct")
}

func (m PlcStruct) String() string {
	allBits := 0
	// TODO: do we want to aggregate the bit length?
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), allBits, m.values)
}
