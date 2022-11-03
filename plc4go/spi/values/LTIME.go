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
	"time"
)

type PlcLTIME struct {
	PlcSimpleValueAdapter
	value uint64
}

func NewPlcLTIME(value uint64) PlcLTIME {
	return PlcLTIME{
		value: value,
	}
}

func (m PlcLTIME) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcLTIME) IsDuration() bool {
	return true
}

func (m PlcLTIME) GetDuration() time.Duration {
	return time.Duration(m.value)
}

func (m PlcLTIME) IsString() bool {
	return true
}

func (m PlcLTIME) GetString() string {
	return fmt.Sprintf("PT%0.fS", m.GetDuration().Seconds())
}

func (m PlcLTIME) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.LTIME
}

func (m PlcLTIME) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcLTIME) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcLTIME", uint32(len([]rune(m.GetString()))*8), "UTF-8", m.GetString())
}

func (m PlcLTIME) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(len([]rune(m.GetString()))*8), m.value)
}
