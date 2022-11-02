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
	"encoding/binary"
	"fmt"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcWCHAR struct {
	PlcSimpleValueAdapter
	value string
}

func NewPlcWCHAR(value string) PlcWCHAR {
	return PlcWCHAR{
		value: value,
	}
}

func (m PlcWCHAR) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcWCHAR) IsString() bool {
	return true
}

func (m PlcWCHAR) GetString() string {
	return string(m.value)
}

func (m PlcWCHAR) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.WCHAR
}

func (m PlcWCHAR) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcWCHAR) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcWCHAR", uint32(16), "UTF-16", string(m.value))
}

func (m PlcWCHAR) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(8), m.value)
}
