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

type PlcDATE_AND_TIME struct {
	PlcValueAdapter
	value time.Time
}

func NewPlcDATE_AND_TIME(value time.Time) PlcDATE_AND_TIME {
	return PlcDATE_AND_TIME{
		value: value,
	}
}

func NewPlcDATA_AND_TIMEFromSegments(year, month, day, hour, minutes, seconds, nanoseconds uint32) PlcDATE_AND_TIME {
	return NewPlcDATE_AND_TIME(time.Date(int(year), time.Month(month), int(day), int(hour), int(minutes), int(seconds), int(nanoseconds), time.Local))
}

func NewPlcDATA_AND_TIMEFromSecondsSinceEpoch(secondsSinceEpoch uint32) PlcDATE_AND_TIME {
	return NewPlcDATE_AND_TIME(time.Unix(int64(secondsSinceEpoch), 0))
}

func (m PlcDATE_AND_TIME) IsRaw() bool {
	return true
}

func (m PlcDATE_AND_TIME) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcDATE_AND_TIME) GetSecondsSinceEpoch() uint32 {
	return uint32(m.value.Unix())
}

func (m PlcDATE_AND_TIME) GetYear() uint16 {
	return uint16(m.value.Year())
}

func (m PlcDATE_AND_TIME) GetMonth() uint8 {
	return uint8(m.value.Month())
}

func (m PlcDATE_AND_TIME) GetDay() uint8 {
	return uint8(m.value.Day())
}

// GetDayOfWeek returns the ISO-8601 day of week: 1 == Monday .. 7 == Sunday.
//
// This mirrors plc4j's PlcDATE_AND_TIME#getDayOfWeek, which returns
// LocalDateTime#getDayOfWeek#getValue, i.e. the same 1..7 ISO numbering. It is
// deliberately NOT time.Weekday (0 == Sunday .. 6 == Saturday): every wire
// format that carries this field numbers the days from 1, so a 0 either means
// "no day given" (KNX DPT 19.001, knxnetip.mspec dayOfWeek) or is outright
// invalid (S7 DATE_AND_TIME, s7.mspec: "one 4-bit value representing 1 - 7").
//
// Note that the S7 dow field uses a DIFFERENT 1..7 numbering (1 == Sunday ..
// 7 == Saturday). Rotating into it is a protocol level concern that this generic
// value cannot do for KNX and S7 at the same time, and the S7 data-io serializer
// does not do it yet.
func (m PlcDATE_AND_TIME) GetDayOfWeek() uint8 {
	// time.Weekday numbers Sunday 0, ISO-8601 numbers it last as 7.
	if weekday := m.value.Weekday(); weekday != time.Sunday {
		return uint8(weekday)
	}
	return 7
}

func (m PlcDATE_AND_TIME) GetHour() uint8 {
	return uint8(m.value.Hour())
}

func (m PlcDATE_AND_TIME) GetMinutes() uint8 {
	return uint8(m.value.Minute())
}

func (m PlcDATE_AND_TIME) GetSeconds() uint8 {
	return uint8(m.value.Second())
}

func (m PlcDATE_AND_TIME) GetNanoseconds() uint32 {
	return uint32(m.value.Nanosecond())
}

func (m PlcDATE_AND_TIME) GetMillisecondsOfSecond() uint64 {
	return uint64(time.Duration(m.GetNanoseconds()).Milliseconds())
}

func (m PlcDATE_AND_TIME) GetNanosecondsOfSecond() uint64 {
	return uint64(time.Duration(m.GetNanoseconds()))
}

func (m PlcDATE_AND_TIME) IsDateTime() bool {
	return true
}

func (m PlcDATE_AND_TIME) GetDateTime() time.Time {
	return m.value
}

func (m PlcDATE_AND_TIME) GetString() string {
	// DATE_AND_TIME is a wall-clock value without a time zone: render the instant's UTC
	// wall time in ISO-8601, with the fraction only when there is one.
	return m.GetDateTime().UTC().Format("2006-01-02T15:04:05.999999999")
}

func (m PlcDATE_AND_TIME) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.DATE_AND_TIME
}

func (m PlcDATE_AND_TIME) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcDATE_AND_TIME) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcDATE_AND_TIME", uint32(len([]rune(m.GetString()))*8), m.GetString(), utils.WithEncoding("UTF-8"))
}

func (m PlcDATE_AND_TIME) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(len([]rune(m.GetString()))*8), m.value)
}
