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

type PlcTIME_OF_DAY struct {
	PlcSimpleValueAdapter
	value time.Time
}

func NewPlcTIME_OF_DAY(value any) PlcTIME_OF_DAY {
	var safeValue time.Time
	switch value.(type) {
	case time.Time:
		castedValue := value.(time.Time)
		safeValue = time.Date(0, 0, 0, castedValue.Hour(), castedValue.Minute(), castedValue.Second(),
			castedValue.Nanosecond(), castedValue.Location())
	case uint32:
		// Interpreted as milliseconds since midnight. The offset carries no timezone,
		// so the components are derived arithmetically instead of by resolving it as
		// an instant: going through time.Unix resolves in the host's local zone and
		// shifts GetHours (and GetMillisecondsSinceMidnight) by the local UTC offset.
		castedValue := value.(uint32)
		seconds := castedValue / 1000
		nanoseconds := (castedValue % 1000) * 1000000
		safeValue = time.Date(0, 0, 0, int(seconds/3600), int(seconds%3600/60),
			int(seconds%60), int(nanoseconds), time.UTC)
	}

	return PlcTIME_OF_DAY{
		value: safeValue,
	}
}

func NewPlcTIME_OF_DAYFromMillisecondsSinceMidnight(millisecondsSinceMidnight uint32) PlcTIME_OF_DAY {
	return NewPlcTIME_OF_DAY(millisecondsSinceMidnight)
}

func (m PlcTIME_OF_DAY) IsRaw() bool {
	return true
}

func (m PlcTIME_OF_DAY) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcTIME_OF_DAY) GetMillisecondsSinceMidnight() uint32 {
	midnight := time.Date(0, 0, 0, 0, 0, 0, 0, m.value.Location())
	return uint32(m.value.UnixMilli() - midnight.UnixMilli())
}

// GetHours returns the hour of the day: 0 .. 23.
func (m PlcTIME_OF_DAY) GetHours() uint8 {
	return uint8(m.value.Hour())
}

// GetMinutes returns the minute of the hour: 0 .. 59.
func (m PlcTIME_OF_DAY) GetMinutes() uint8 {
	return uint8(m.value.Minute())
}

// GetSeconds returns the second of the minute: 0 .. 59.
func (m PlcTIME_OF_DAY) GetSeconds() uint8 {
	return uint8(m.value.Second())
}

// GetCentiseconds returns the hundredths of a second: 0 .. 99.
//
// Protocols that transmit a time-of-day as separate BCD fields use this as the
// sub-second component (e.g. Schneider UMAS, where TOD is laid out as
// [centiseconds][seconds][minutes][hours]). It truncates rather than rounds, so
// it stays consistent with GetSeconds for any sub-second remainder.
func (m PlcTIME_OF_DAY) GetCentiseconds() uint8 {
	return uint8(m.value.Nanosecond() / 10_000_000)
}

func (m PlcTIME_OF_DAY) IsTime() bool {
	return true
}

func (m PlcTIME_OF_DAY) GetTime() time.Time {
	return m.value
}

func (m PlcTIME_OF_DAY) GetString() string {
	return m.value.Format("15:04:05.000")
}

func (m PlcTIME_OF_DAY) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.TIME_OF_DAY
}

func (m PlcTIME_OF_DAY) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcTIME_OF_DAY) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcTIME_OF_DAY", uint32(len([]rune(m.GetString()))*8), m.GetString())
}

func (m PlcTIME_OF_DAY) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(len([]rune(m.GetString()))*8), m.value)
}
