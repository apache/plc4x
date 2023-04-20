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
	"time"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcLDATE struct {
	PlcValueAdapter
	value time.Time
}

func NewPlcLDATE(value any) PlcLDATE {
	var timeValue time.Time
	switch value.(type) {
	case time.Time:
		timeValue = value.(time.Time)
	case uint64:
		// Interpreted as "nanoseconds since epoch"
		castedValue := value.(uint64)
		seconds := castedValue / 1000000
		nanoseconds := castedValue % 1000000
		timeValue = time.Unix(int64(seconds), int64(nanoseconds))
	}
	safeValue := time.Date(timeValue.Year(), timeValue.Month(), timeValue.Day(), 0, 0, 0, 0, timeValue.Location())
	return PlcLDATE{
		value: safeValue,
	}
}

func NewPlcLDATEFromNanosecondsSinceEpoch(nanosecondsSinceEpoch uint64) PlcDATE {
	return NewPlcDATE(nanosecondsSinceEpoch)
}

func (m PlcLDATE) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcLDATE) GetNanosecondsSinceEpoch() uint64 {
	return uint64(m.value.UnixNano())
}

func (m PlcLDATE) IsDate() bool {
	return true
}

func (m PlcLDATE) GetDate() time.Time {
	return time.Date(m.value.Year(), m.value.Month(), m.value.Day(), 0, 0, 0, 0, time.UTC)
}

func (m PlcLDATE) GetString() string {
	return m.GetDate().Format("2006-01-02")
}

func (m PlcLDATE) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.DATE
}

func (m PlcLDATE) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcLDATE) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcLDATE", uint32(len([]rune(m.GetString()))*8), "UTF-8", m.GetString())
}

func (m PlcLDATE) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(len([]rune(m.GetString()))*8), m.value)
}
