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

type PlcDATE_AND_LTIME struct {
	PlcValueAdapter
	value time.Time
}

func NewPlcDATE_AND_LTIME(value time.Time) PlcDATE_AND_LTIME {
	return PlcDATE_AND_LTIME{
		value: value,
	}
}

func NewPlcDATA_AND_LTIMEFromSegments(year, month, day, hour, minutes, seconds, nanoseconds uint32) PlcDATE_AND_LTIME {
	return NewPlcDATE_AND_LTIME(time.Date(int(year), time.Month(month), int(day), int(hour), int(minutes), int(seconds), int(nanoseconds), time.Local))
}

func NewPlcDATA_AND_LTIMEFromNanosecondsSinceEpoch(nanosecondsSinceEpoch uint64) PlcDATE_AND_LTIME {
	timeStamp := time.Time{}
	timeStamp.Add(time.Duration(nanosecondsSinceEpoch) * time.Nanosecond)
	return NewPlcDATE_AND_LTIME(timeStamp)
}

func (m PlcDATE_AND_LTIME) IsRaw() bool {
	return true
}

func (m PlcDATE_AND_LTIME) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcDATE_AND_LTIME) GetSecondsSinceEpoch() uint32 {
	return uint32(m.value.Unix())
}

func (m PlcDATE_AND_LTIME) GetYear() uint16 {
	return uint16(m.value.Year())
}

func (m PlcDATE_AND_LTIME) GetMonth() uint8 {
	return uint8(m.value.Month())
}

func (m PlcDATE_AND_LTIME) GetDay() uint8 {
	return uint8(m.value.Day())
}

func (m PlcDATE_AND_LTIME) GetDayOfWeek() uint8 {
	return uint8(m.value.Weekday())
}

func (m PlcDATE_AND_LTIME) GetHour() uint8 {
	return uint8(m.value.Hour())
}

func (m PlcDATE_AND_LTIME) GetMinutes() uint8 {
	return uint8(m.value.Minute())
}

func (m PlcDATE_AND_LTIME) GetSeconds() uint8 {
	return uint8(m.value.Second())
}

func (m PlcDATE_AND_LTIME) GetNanoseconds() uint32 {
	return uint32(m.value.Nanosecond())
}

func (m PlcDATE_AND_LTIME) GetMillisecondsOfSecond() uint64 {
	return uint64(time.Duration(m.GetNanoseconds()).Milliseconds())
}

func (m PlcDATE_AND_LTIME) IsDateTime() bool {
	return true
}

func (m PlcDATE_AND_LTIME) GetDateTime() time.Time {
	return m.value
}

func (m PlcDATE_AND_LTIME) GetString() string {
	return fmt.Sprintf("%v", m.GetDateTime())
}

func (m PlcDATE_AND_LTIME) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.DATE_AND_LTIME
}

func (m PlcDATE_AND_LTIME) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcDATE_AND_LTIME) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcDATE_AND_LTIME", uint32(len([]rune(m.GetString()))*8), m.GetString())
}

func (m PlcDATE_AND_LTIME) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(len([]rune(m.GetString()))*8), m.value)
}
