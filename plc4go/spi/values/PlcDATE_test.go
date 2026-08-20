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

// TestPlcDATE_Components pins the component getters the umas DataItem
// serialization needs: umas.mspec lays DATE out as three BCD fields
// [day][month][year], so the numbering has to be the human one (1 == January,
// 1 == first of the month) and NOT any zero based variant. This matches
// PlcDATE_AND_TIME.GetMonth/GetDay, which the neighbouring DATE_AND_TIME branch
// of the same generated switch already uses.
func TestPlcDATE_Components(t *testing.T) {
	tests := []struct {
		name       string
		in         time.Time
		year       uint16
		month, day uint8
	}{
		{"january first", time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC), 2024, 1, 1},
		{"december last", time.Date(1999, time.December, 31, 23, 59, 59, 0, time.UTC), 1999, 12, 31},
		{"leap day", time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC), 2024, 2, 29},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value := NewPlcDATE(test.in)
			assert.Equal(t, test.year, value.GetYear())
			assert.Equal(t, test.month, value.GetMonth())
			assert.Equal(t, test.day, value.GetDay())
		})
	}
}

// TestPlcDATE_ComponentsIgnoreTheTimeOfDay guards that the getters read the
// normalized date and not the instant handed in: NewPlcDATE truncates to
// midnight, so a value built from a late evening timestamp must still report the
// same day.
func TestPlcDATE_ComponentsIgnoreTheTimeOfDay(t *testing.T) {
	day := time.Date(2024, time.August, 19, 23, 59, 59, 999999999, time.UTC)
	value := NewPlcDATE(day)
	assert.Equal(t, uint16(2024), value.GetYear())
	assert.Equal(t, uint8(8), value.GetMonth())
	assert.Equal(t, uint8(19), value.GetDay())
}

// TestPlcDATE_ComponentsFromEpochAreTimezoneIndependent pins that an epoch based
// count is resolved as UTC. It used to go through time.Unix, which resolves in
// the host's local zone: west of Greenwich that reports the previous calendar
// day, and even east of it the day boundary lands in the wrong place, so the
// same S7/ADS/KNX payload decoded to a different date per machine.
func TestPlcDATE_ComponentsFromEpochAreTimezoneIndependent(t *testing.T) {
	// 1990-01-01 is the Siemens epoch, 7305 days after the Unix epoch.
	value := NewPlcDATEFromDaysSinceEpoch(7305)
	assert.Equal(t, uint16(1990), value.GetYear())
	assert.Equal(t, uint8(1), value.GetMonth())
	assert.Equal(t, uint8(1), value.GetDay())
	assert.Equal(t, uint16(7305), value.GetDaysSinceEpoch(), "day count must round trip")
	assert.Equal(t, uint16(0), value.GetDaysSinceSiemensEpoch())

	// The Unix epoch itself: the most sensitive case, since a negative local
	// offset used to push it back into 1969.
	epoch := NewPlcDATEFromSecondsSinceEpoch(0)
	assert.Equal(t, uint16(1970), epoch.GetYear())
	assert.Equal(t, uint8(1), epoch.GetMonth())
	assert.Equal(t, uint8(1), epoch.GetDay())
}
