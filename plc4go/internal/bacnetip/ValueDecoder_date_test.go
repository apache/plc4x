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

package bacnetip

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
)

func appTagFromTheWire(t *testing.T, wire []byte) apiValues.PlcValue {
	t.Helper()
	tag, err := readWriteModel.BACnetApplicationTagParse[readWriteModel.BACnetApplicationTag](t.Context(), wire)
	require.NoError(t, err)
	return appTagToPlcValue(tag)
}

// A Date that names one calendar day is a PlcDATE; the day of week is not checked.
func TestAppTagToPlcValue_concreteDateIsADate(t *testing.T) {
	cases := []struct {
		name string
		wire []byte
		want time.Time
	}{
		{"2026-11-24 Tuesday", []byte{0xA4, 0x7E, 0x0B, 0x18, 0x02}, time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)},
		{"any day of week", []byte{0xA4, 0x7E, 0x0B, 0x18, 0xFF}, time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)},
		{"leap day", []byte{0xA4, 0x7C, 0x02, 0x1D, 0x04}, time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)},
		{"first year", []byte{0xA4, 0x00, 0x01, 0x01, 0x01}, time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"last year", []byte{0xA4, 0xFE, 0x0C, 0x1F, 0xFF}, time.Date(2154, 12, 31, 0, 0, 0, 0, time.UTC)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := appTagFromTheWire(t, c.wire)
			require.Equal(t, apiValues.DATE, got.GetPlcValueType())
			assert.Equal(t, c.want, got.GetDate())
		})
	}
}

// A Date pattern (clause 20.2.12) or an impossible day must not be normalised into a
// different real date: it comes back as the raw octets.
func TestAppTagToPlcValue_datePatternKeepsItsOctets(t *testing.T) {
	cases := []struct {
		name string
		wire []byte
	}{
		{"any year", []byte{0xA4, 0xFF, 0x0C, 0x19, 0xFF}},
		{"any month", []byte{0xA4, 0x7E, 0xFF, 0x01, 0xFF}},
		{"odd months", []byte{0xA4, 0x7E, 0x0D, 0x01, 0xFF}},
		{"even months", []byte{0xA4, 0x7E, 0x0E, 0x01, 0xFF}},
		{"any day", []byte{0xA4, 0x7E, 0x03, 0xFF, 0xFF}},
		{"last day of month", []byte{0xA4, 0x7E, 0x03, 0x20, 0xFF}},
		{"odd days", []byte{0xA4, 0x7E, 0x03, 0x21, 0xFF}},
		{"even days", []byte{0xA4, 0x7E, 0x03, 0x22, 0xFF}},
		{"everything wildcard", []byte{0xA4, 0xFF, 0xFF, 0xFF, 0xFF}},
		{"February 30", []byte{0xA4, 0x7E, 0x02, 0x1E, 0xFF}},
		{"February 29 outside a leap year", []byte{0xA4, 0x7E, 0x02, 0x1D, 0xFF}},
		{"month 0", []byte{0xA4, 0x7E, 0x00, 0x01, 0xFF}},
		{"day 0", []byte{0xA4, 0x7E, 0x01, 0x00, 0xFF}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := appTagFromTheWire(t, c.wire)
			require.Equal(t, apiValues.Struct, got.GetPlcValueType(), "got %v", got)
			assert.Equal(t, []string{"yearMinus1900", "month", "dayOfMonth", "dayOfWeek"}, got.GetKeys())
			for i, key := range got.GetKeys() {
				assert.Equal(t, c.wire[1+i], got.GetValue(key).GetUint8(), key)
			}
		})
	}
}

func TestAppTagToPlcValue_concreteTimeIsATimeOfDay(t *testing.T) {
	got := appTagFromTheWire(t, []byte{0xB4, 0x0D, 0x2D, 0x1E, 0x19})
	require.Equal(t, apiValues.TIME_OF_DAY, got.GetPlcValueType())
	assert.Equal(t, "13:45:30.250", got.GetTime().Format("15:04:05.000"))
}

// A Time with a wildcard (255) or an out-of-range octet comes back as the raw octets
// instead of a duration such as 255 hours.
func TestAppTagToPlcValue_timePatternKeepsItsOctets(t *testing.T) {
	cases := []struct {
		name string
		wire []byte
	}{
		{"any time", []byte{0xB4, 0xFF, 0xFF, 0xFF, 0xFF}},
		{"any hour", []byte{0xB4, 0xFF, 0x00, 0x00, 0x00}},
		{"noon, any second", []byte{0xB4, 0x0C, 0x00, 0xFF, 0xFF}},
		{"hour 24", []byte{0xB4, 0x18, 0x00, 0x00, 0x00}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := appTagFromTheWire(t, c.wire)
			require.Equal(t, apiValues.Struct, got.GetPlcValueType(), "got %v", got)
			assert.Equal(t, []string{"hour", "minute", "second", "fractional"}, got.GetKeys())
			for i, key := range got.GetKeys() {
				assert.Equal(t, c.wire[1+i], got.GetValue(key).GetUint8(), key)
			}
		})
	}
}
