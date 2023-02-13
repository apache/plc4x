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

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcWSTRING struct {
	PlcSimpleValueAdapter
	value string
}

func NewPlcWSTRING(value string) PlcWSTRING {
	return PlcWSTRING{
		value: value,
	}
}

func (m PlcWSTRING) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcWSTRING) IsString() bool {
	return true
}

func (m PlcWSTRING) GetString() string {
	return string(m.value)
}

func (m PlcWSTRING) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.WSTRING
}

func (m PlcWSTRING) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcWSTRING) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcSTRING", uint32(len(m.value)*8), "UTF-8", string(m.value))
}

func (m PlcWSTRING) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(len(m.value)*8), m.value)
}
