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
	"strconv"
	"strings"
	"time"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type PlcTIME struct {
	PlcSimpleValueAdapter
	value time.Duration
}

func NewPlcTIME(value time.Duration) PlcTIME {
	return PlcTIME{
		value: value,
	}
}

func NewPlcTIMEFromMilliseconds(milliseconds int64) PlcTIME {
	return NewPlcTIME(time.Duration(milliseconds) * time.Millisecond)
}

func (m PlcTIME) IsRaw() bool {
	return true
}

func (m PlcTIME) GetRaw() []byte {
	theBytes, _ := m.Serialize()
	return theBytes
}

func (m PlcTIME) GetMilliseconds() int64 {
	return m.value.Milliseconds()
}

func (m PlcTIME) IsDuration() bool {
	return true
}

func (m PlcTIME) GetDuration() time.Duration {
	return m.value
}

func (m PlcTIME) IsString() bool {
	return true
}

func (m PlcTIME) GetString() string {
	return formatDurationIso8601(m.GetDuration())
}

// parseDurationIso8601 parses the time-only ISO-8601 duration form produced by
// formatDurationIso8601 ("PT24015H23M12.034002044S").
func parseDurationIso8601(value string) (time.Duration, error) {
	input := value
	negative := strings.HasPrefix(input, "-")
	input = strings.TrimPrefix(input, "-")
	rest, ok := strings.CutPrefix(input, "PT")
	if !ok || rest == "" {
		return 0, fmt.Errorf("invalid ISO-8601 duration %q", value)
	}
	var duration time.Duration
	for _, component := range []struct {
		designator string
		unit       time.Duration
	}{{"H", time.Hour}, {"M", time.Minute}, {"S", time.Second}} {
		number, remainder, found := strings.Cut(rest, component.designator)
		if !found {
			continue
		}
		amount, err := strconv.ParseFloat(number, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid ISO-8601 duration %q: %v", value, err)
		}
		duration += time.Duration(amount * float64(component.unit))
		rest = remainder
	}
	if rest != "" {
		return 0, fmt.Errorf("invalid ISO-8601 duration %q", value)
	}
	if negative {
		duration = -duration
	}
	return duration, nil
}

// formatDurationIso8601 renders a duration in ISO-8601 form with an hours/minutes/seconds
// decomposition and the sub-second fraction ("PT24015H23M12.034002044S"), the shared
// cross-implementation rendering of duration-typed PLC values.
func formatDurationIso8601(duration time.Duration) string {
	if duration == 0 {
		return "PT0S"
	}
	var sb strings.Builder
	if duration < 0 {
		sb.WriteString("-")
		duration = -duration
	}
	sb.WriteString("PT")
	totalSeconds := int64(duration / time.Second)
	nanos := int64(duration % time.Second)
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	if hours != 0 {
		fmt.Fprintf(&sb, "%dH", hours)
	}
	if minutes != 0 {
		fmt.Fprintf(&sb, "%dM", minutes)
	}
	if seconds != 0 || nanos != 0 || (hours == 0 && minutes == 0) {
		if nanos == 0 {
			fmt.Fprintf(&sb, "%dS", seconds)
		} else {
			fmt.Fprintf(&sb, "%d.%sS", seconds, strings.TrimRight(fmt.Sprintf("%09d", nanos), "0"))
		}
	}
	return sb.String()
}

func (m PlcTIME) GetPlcValueType() apiValues.PlcValueType {
	return apiValues.TIME
}

func (m PlcTIME) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m PlcTIME) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	return writeBuffer.WriteString("PlcTIME", uint32(len([]rune(m.GetString()))*8), m.GetString(), utils.WithEncoding("UTF-8"))
}

func (m PlcTIME) String() string {
	return fmt.Sprintf("%s(%dbit):%v", m.GetPlcValueType(), uint32(len([]rune(m.GetString()))*8), m.value)
}
