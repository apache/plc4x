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

type PlcCHAR struct {
	PlcSimpleValueAdapter
	// TODO: Why is this a byte-array?
	value string
}

func NewPlcCHAR(value string) PlcCHAR {
	return PlcCHAR{
		value: value,
	}
}

func (m PlcCHAR) IsRaw() bool {
	return true
}

func (m PlcCHAR) GetRaw() []byte {
	return []byte(m.value)
}

func (m PlcCHAR) IsString() bool {
	return true
}

func (m PlcCHAR) GetString() string {
	return m.value
}

func (m PlcCHAR) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.CHAR
}

func (m PlcCHAR) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcCHAR) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcCHAR", 8, "UTF-8", m.value)
}

func (m PlcCHAR) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), 8, m.value)
}
