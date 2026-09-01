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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestPlcDATE_AND_TIME_GetDayOfWeek pins the ISO-8601 numbering (1 == Monday ..
// 7 == Sunday) that both consumers of this getter need.
//
// It used to return time.Weekday, i.e. 0 == Sunday .. 6 == Saturday. That is
// wrong for every wire format in the tree: the KNX DPT 19.001 dayOfWeek field
// (protocols/knxnetip .. KnxDatapoint.go) reserves 0 for "no day given", and the
// S7 DATE_AND_TIME dow nibble is documented in s7.mspec as "representing 1 - 7",
// so a Sunday used to serialize as the invalid value 0. plc4j's
// PlcDATE_AND_TIME#getDayOfWeek returns LocalDateTime#getDayOfWeek#getValue,
// which is this same 1..7 numbering.
func TestPlcDATE_AND_TIME_GetDayOfWeek(t *testing.T) {
	// 2024-08-19 is a Monday, so the week below walks Monday to Sunday.
	expected := []uint8{1, 2, 3, 4, 5, 6, 7}
	weekdays := []time.Weekday{
		time.Monday, time.Tuesday, time.Wednesday, time.Thursday,
		time.Friday, time.Saturday, time.Sunday,
	}
	for i, want := range expected {
		day := time.Date(2024, time.August, 19+i, 10, 30, 45, 0, time.UTC)
		t.Run(weekdays[i].String(), func(t *testing.T) {
			assert.Equal(t, weekdays[i], day.Weekday(), "fixture sanity")
			assert.Equal(t, want, NewPlcDATE_AND_TIME(day).GetDayOfWeek())
			assert.Equal(t, want, NewPlcDATE_AND_LTIME(day).GetDayOfWeek())
		})
	}
}

// TestPlcDATE_AND_TIME_GetDayOfWeekIsAlwaysInRange guards the whole domain: no
// input may ever produce the 0 that both wire formats treat as "unset"/invalid.
func TestPlcDATE_AND_TIME_GetDayOfWeekIsAlwaysInRange(t *testing.T) {
	day := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	for i := range 400 {
		got := NewPlcDATE_AND_TIME(day.AddDate(0, 0, i)).GetDayOfWeek()
		assert.GreaterOrEqual(t, got, uint8(1))
		assert.LessOrEqual(t, got, uint8(7))
	}
}
