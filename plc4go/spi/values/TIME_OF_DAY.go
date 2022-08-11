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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"time"
)

type PlcTIME_OF_DAY struct {
	value time.Time
	PlcSimpleValueAdapter
}

func NewPlcTIME_OF_DAY(value interface{}) PlcTIME_OF_DAY {
	var safeValue time.Time
	switch value.(type) {
	case time.Time:
		castedValue := value.(time.Time)
		safeValue = time.Date(0, 0, 0, castedValue.Hour(), castedValue.Minute(), castedValue.Second(),
			castedValue.Nanosecond(), castedValue.Location())
	case uint32:
		// Interpreted as milliseconds since midnight
		castedValue := value.(uint32)
		seconds := castedValue / 1000
		nanoseconds := (castedValue % 1000) * 1000000
		epochTime := time.Unix(int64(seconds), int64(nanoseconds))
		safeValue = time.Date(0, 0, 0, epochTime.Hour(), epochTime.Minute(), epochTime.Second(),
			epochTime.Nanosecond(), epochTime.Location())
	}

	return PlcTIME_OF_DAY{
		value: safeValue,
	}
}

func (m PlcTIME_OF_DAY) GetRaw() []byte {
	buf := utils.NewWriteBufferByteBased()
	m.Serialize(buf)
	return buf.GetBytes()
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

func (m PlcTIME_OF_DAY) Serialize(writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcTIME_OF_DAY", uint32(len([]rune(m.GetString()))*8), "UTF-8", m.GetString())
}
