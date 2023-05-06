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

type PlcLTIME_OF_DAY struct {
	PlcSimpleValueAdapter
	value time.Time
}

func NewPlcLTIME_OF_DAY(value any) PlcLTIME_OF_DAY {
	var safeValue time.Time
	switch value.(type) {
	case time.Time:
		castedValue := value.(time.Time)
		safeValue = time.Date(0, 0, 0, castedValue.Hour(), castedValue.Minute(), castedValue.Second(),
			castedValue.Nanosecond(), castedValue.Location())
	case uint64:
		// Interpreted as nanoseconds since midnight
		castedValue := value.(uint64)
		seconds := castedValue / 1000000
		nanoseconds := castedValue % 1000000
		epochTime := time.Unix(int64(seconds), int64(nanoseconds))
		safeValue = time.Date(0, 0, 0, epochTime.Hour(), epochTime.Minute(), epochTime.Second(),
			epochTime.Nanosecond(), epochTime.Location())
	}

	return PlcLTIME_OF_DAY{
		value: safeValue,
	}
}

func NewPlcLTIME_OF_DAYFromNanosecondsSinceMidnight(nanosecondsSinceMidnight uint64) PlcLTIME_OF_DAY {
	return NewPlcLTIME_OF_DAY(nanosecondsSinceMidnight)
}

func (m PlcLTIME_OF_DAY) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcLTIME_OF_DAY) GetNanosecondsSinceMidnight() uint64 {
	midnight := time.Date(0, 0, 0, 0, 0, 0, 0, m.value.Location())
	return uint64(m.value.UnixNano() - midnight.UnixNano())
}

func (m PlcLTIME_OF_DAY) IsTime() bool {
	return true
}

func (m PlcLTIME_OF_DAY) GetTime() time.Time {
	return m.value
}

func (m PlcLTIME_OF_DAY) GetString() string {
	return m.value.Format("15:04:05.000")
}

func (m PlcLTIME_OF_DAY) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.TIME_OF_DAY
}

func (m PlcLTIME_OF_DAY) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcLTIME_OF_DAY) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcLTIME_OF_DAY", uint32(len([]rune(m.GetString()))*8), "UTF-8", m.GetString())
}

func (m PlcLTIME_OF_DAY) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(len([]rune(m.GetString()))*8), m.value)
}
