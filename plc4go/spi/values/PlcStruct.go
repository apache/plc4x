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
	"sort"
	"strings"

	"github.com/rs/zerolog/log"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcStruct struct {
	values map[string]apiValues.PlcValue
	// order fixes the member order for GetKeys and serialization: insertion order when
	// constructed with NewPlcStructOrdered (e.g. the declaration order of a decoded PLC
	// struct), sorted key order otherwise - never Go's random map iteration order.
	order []string
	PlcValueAdapter
}

func NewPlcStruct(value map[string]apiValues.PlcValue) PlcStruct {
	order := make([]string, 0, len(value))
	for key := range value {
		order = append(order, key)
	}
	sort.Strings(order)
	return PlcStruct{
		values: value,
		order:  order,
	}
}

// NewPlcStructOrdered creates a PlcStruct whose members keep the given order (usually the
// struct's declaration order on the PLC). Keys missing from values are ignored; values not
// listed in order are appended in sorted order.
func NewPlcStructOrdered(value map[string]apiValues.PlcValue, order []string) PlcStruct {
	seen := make(map[string]struct{}, len(order))
	orderedKeys := make([]string, 0, len(value))
	for _, key := range order {
		if _, ok := value[key]; ok {
			orderedKeys = append(orderedKeys, key)
			seen[key] = struct{}{}
		}
	}
	var rest []string
	for key := range value {
		if _, ok := seen[key]; !ok {
			rest = append(rest, key)
		}
	}
	sort.Strings(rest)
	return PlcStruct{
		values: value,
		order:  append(orderedKeys, rest...),
	}
}

////
// Raw Access

func (m PlcStruct) IsRaw() bool {
	return true
}

func (m PlcStruct) GetRaw() []byte {
	if theBytes, err := m.Serialize(); err != nil {
		log.Error().Err(err).Msg("Error getting raw")
		return nil
	} else {
		return theBytes
	}
}

//
///

func (m PlcStruct) IsStruct() bool {
	return true
}

func (m PlcStruct) GetKeys() []string {
	return append([]string(nil), m.order...)
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
	for _, tagName := range m.order {
		tagValue := m.values[tagName]
		sb.WriteString("  ")
		sb.WriteString(tagName)
		sb.WriteString(": \"")
		if tagValue.IsString() {
			sb.WriteString(tagValue.GetString())
		} else {
			sb.WriteString(fmt.Sprintf("%v", tagValue))
		}
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
	for _, tagName := range m.order {
		tagValue := m.values[tagName]
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
