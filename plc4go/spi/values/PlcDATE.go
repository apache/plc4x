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

type PlcDATE struct {
	PlcValueAdapter
	value time.Time
}

func NewPlcDATE(value any) PlcDATE {
	var timeValue time.Time
	switch value.(type) {
	case time.Time:
		timeValue = value.(time.Time)
	case uint16:
		// In this case the date is the number of days since 1990-01-01
		// So we gotta add 7305 days to the value to have it relative to epoch
		// Then we also need to transform it from days to seconds by multiplying by 86400
		timeValue = time.Unix((int64(value.(uint16))+7305)*86400, 0)
	case uint32:
		// Interpreted as "seconds since epoch"
		timeValue = time.Unix(int64(value.(uint32)), 0)
	}
	safeValue := time.Date(timeValue.Year(), timeValue.Month(), timeValue.Day(), 0, 0, 0, 0, timeValue.Location())
	return PlcDATE{
		value: safeValue,
	}
}

func NewPlcDATEFromSecondsSinceEpoch(secondsSinceEpoch uint32) PlcDATE {
	return NewPlcDATE(time.Unix(int64(secondsSinceEpoch), 0))
}

func NewPlcDATEFromDaysSinceEpoch(daysSinceEpoch uint16) PlcDATE {
	// 86400 = 24 hours x 60 Minutes x 60 Seconds
	return NewPlcDATE(time.Unix(int64(daysSinceEpoch)*86400, 0))
}

func NewPlcDATEFromDaysSinceSiemensEpoch(daysSinceSiemensEpoch uint16) PlcDATE {
	// 86400 = 24 hours x 60 Minutes x 60 Seconds
	return NewPlcDATEFromDaysSinceEpoch(daysSinceSiemensEpoch + 7305)
}

func (m PlcDATE) IsRaw() bool {
	return true
}

func (m PlcDATE) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcDATE) GetSecondsSinceEpoch() uint32 {
	return uint32(m.value.Unix())
}

func (m PlcDATE) GetDaysSinceEpoch() uint16 {
	// Seconds to days
	return uint16(m.value.Unix() / 86400)
}

func (m PlcDATE) GetDaysSinceSiemensEpoch() uint16 {
	// Seconds to days to 1990-01-01
	return uint16((m.value.Unix() / 86400) - 7305)
}

func (m PlcDATE) IsDate() bool {
	return true
}

func (m PlcDATE) GetDate() time.Time {
	return time.Date(m.value.Year(), m.value.Month(), m.value.Day(), 0, 0, 0, 0, time.UTC)
}

func (m PlcDATE) GetString() string {
	return m.GetDate().Format("2006-01-02")
}

func (m PlcDATE) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.DATE
}

func (m PlcDATE) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcDATE) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcDATE", uint32(len([]rune(m.GetString()))*8), "UTF-8", m.GetString())
}

func (m PlcDATE) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(len([]rune(m.GetString()))*8), m.value)
}
