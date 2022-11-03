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

type PlcDATE_AND_TIME struct {
	value time.Time
	PlcValueAdapter
}

func NewPlcDATE_AND_TIME(value time.Time) PlcDATE_AND_TIME {
	return PlcDATE_AND_TIME{
		value: value,
	}
}

func (m PlcDATE_AND_TIME) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcDATE_AND_TIME) IsDateTime() bool {
	return true
}
func (m PlcDATE_AND_TIME) GetDateTime() time.Time {
	return m.value
}

func (m PlcDATE_AND_TIME) GetString() string {
	return fmt.Sprintf("%v", m.GetDateTime())
}

func (m PlcDATE_AND_TIME) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.DATE_AND_TIME
}

func (m PlcDATE_AND_TIME) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcDATE_AND_TIME) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcDATE_AND_TIME", uint32(len([]rune(m.GetString()))*8), "UTF-8", m.GetString())
}

func (m PlcDATE_AND_TIME) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(len([]rune(m.GetString()))*8), m.value)
}
