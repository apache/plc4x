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

// TestPlcTIME_OF_DAY_Components pins the component getters the umas DataItem
// serialization needs. umas.mspec lays TOD out as four BCD fields in the order
// [centiseconds][seconds][minutes][hours], so all four have to be readable
// individually and centiseconds is hundredths of a second, matching plc4j's
// PlcTIME_OF_DAY#ofSegments.
func TestPlcTIME_OF_DAY_Components(t *testing.T) {
	tests := []struct {
		name                                  string
		in                                    time.Time
		hours, minutes, seconds, centiseconds uint8
	}{
		{"midnight", time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), 0, 0, 0, 0},
		{"last representable instant", time.Date(0, 0, 0, 23, 59, 59, 990000000, time.UTC), 23, 59, 59, 99},
		{"afternoon", time.Date(0, 0, 0, 14, 30, 45, 120000000, time.UTC), 14, 30, 45, 12},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := NewPlcTIME_OF_DAY(test.in)
			assert.Equal(t, test.hours, value.GetHours())
			assert.Equal(t, test.minutes, value.GetMinutes())
			assert.Equal(t, test.seconds, value.GetSeconds())
			assert.Equal(t, test.centiseconds, value.GetCentiseconds())
		})
	}
}

// TestPlcTIME_OF_DAY_GetCentisecondsTruncates pins that the sub second component
// truncates rather than rounds. Rounding would let 999ms report 100
// centiseconds, which is out of the 0..99 range the BCD field can hold and would
// disagree with GetSeconds, which does not carry the rounded second.
func TestPlcTIME_OF_DAY_GetCentisecondsTruncates(t *testing.T) {
	for millis, want := range map[int]uint8{0: 0, 9: 0, 10: 1, 19: 1, 999: 99} {
		value := NewPlcTIME_OF_DAY(time.Date(0, 0, 0, 1, 2, 3, millis*1000000, time.UTC))
		assert.Equal(t, want, value.GetCentiseconds(), "%dms", millis)
		assert.Equal(t, uint8(3), value.GetSeconds(), "%dms must not roll into the second", millis)
	}
}

// TestPlcTIME_OF_DAY_ComponentsFromMillisecondsAreTimezoneIndependent pins that a
// milliseconds-since-midnight offset is decoded arithmetically. It used to be
// resolved through time.Unix, i.e. in the host's local zone, so the same
// S7/ADS/KNX payload reported a different hour on every machine (and
// GetMillisecondsSinceMidnight did not round trip).
func TestPlcTIME_OF_DAY_ComponentsFromMillisecondsAreTimezoneIndependent(t *testing.T) {
	const millisSinceMidnight = 1*3600000 + 2*60000 + 3*1000 + 450
	value := NewPlcTIME_OF_DAYFromMillisecondsSinceMidnight(millisSinceMidnight)
	assert.Equal(t, uint8(1), value.GetHours())
	assert.Equal(t, uint8(2), value.GetMinutes())
	assert.Equal(t, uint8(3), value.GetSeconds())
	assert.Equal(t, uint8(45), value.GetCentiseconds())
	assert.Equal(t, uint32(millisSinceMidnight), value.GetMillisecondsSinceMidnight(), "offset must round trip")
	assert.Equal(t, "01:02:03.450", value.GetString())
}

// TestPlcTIME_OF_DAY_ComponentsCoverTheWholeDay walks every minute of the day to
// guard that no offset produces an out of range component, since each one has to
// fit a two digit BCD field.
func TestPlcTIME_OF_DAY_ComponentsCoverTheWholeDay(t *testing.T) {
	for minute := range 24 * 60 {
		value := NewPlcTIME_OF_DAYFromMillisecondsSinceMidnight(uint32(minute) * 60000)
		assert.Equal(t, uint8(minute/60), value.GetHours(), "minute %d", minute)
		assert.Equal(t, uint8(minute%60), value.GetMinutes(), "minute %d", minute)
		assert.Equal(t, uint8(0), value.GetSeconds(), "minute %d", minute)
		assert.Equal(t, uint8(0), value.GetCentiseconds(), "minute %d", minute)
	}
}
