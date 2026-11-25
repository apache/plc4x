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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
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

func encodedOnTheWire(t *testing.T, value apiValues.PlcValue) []byte {
	t.Helper()
	tag, err := plcValueToApplicationTag(value, hintNone)
	require.NoError(t, err)
	wire, err := tag.Serialize()
	require.NoError(t, err)
	return wire
}

// A PlcDATE goes out as a BACnet Date with the day of week numbered 1=Mon .. 7=Sun.
func TestPlcValueToApplicationTag_dateIsADate(t *testing.T) {
	cases := []struct {
		name  string
		value apiValues.PlcValue
		wire  []byte
	}{
		{"2026-11-24 Tuesday", spiValues.NewPlcDATE(time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)), []byte{0xA4, 0x7E, 0x0B, 0x18, 0x02}},
		{"Sunday", spiValues.NewPlcDATE(time.Date(2026, 11, 22, 0, 0, 0, 0, time.UTC)), []byte{0xA4, 0x7E, 0x0B, 0x16, 0x07}},
		{"first year", spiValues.NewPlcDATE(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)), []byte{0xA4, 0x00, 0x01, 0x01, 0x01}},
		{"last year", spiValues.NewPlcDATE(time.Date(2154, 12, 31, 0, 0, 0, 0, time.UTC)), []byte{0xA4, 0xFE, 0x0C, 0x1F, 0x02}},
		{"local zone keeps its calendar day", spiValues.NewPlcDATE(time.Date(2026, 11, 24, 0, 0, 0, 0, time.FixedZone("UTC+14", 14*3600))), []byte{0xA4, 0x7E, 0x0B, 0x18, 0x02}},
		{"LDATE", spiValues.NewPlcLDATE(time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)), []byte{0xA4, 0x7E, 0x0B, 0x18, 0x02}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.wire, encodedOnTheWire(t, c.value))
		})
	}
}

// BACnet counts the year in one octet from 1900 and keeps 255 for "any year", so a
// date outside 1900..2154 cannot be sent and must not wrap into another year.
func TestPlcValueToApplicationTag_dateOutsideTheBACnetYearsIsRejected(t *testing.T) {
	for _, year := range []int{1899, 2155, 255} {
		t.Run(fmt.Sprint(year), func(t *testing.T) {
			_, err := plcValueToApplicationTag(spiValues.NewPlcDATE(time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)), hintNone)
			assert.Error(t, err)
		})
	}
}

// A PlcTIME_OF_DAY goes out as a BACnet Time; the sub-second part is truncated to
// hundredths, the resolution of the wire format.
func TestPlcValueToApplicationTag_timeOfDayIsATime(t *testing.T) {
	cases := []struct {
		name  string
		value apiValues.PlcValue
		wire  []byte
	}{
		{"13:45:30.25", spiValues.NewPlcTIME_OF_DAYFromMillisecondsSinceMidnight(((13*60+45)*60+30)*1000 + 250), []byte{0xB4, 0x0D, 0x2D, 0x1E, 0x19}},
		{"milliseconds truncate", spiValues.NewPlcTIME_OF_DAYFromMillisecondsSinceMidnight(((13*60+45)*60+7)*1000 + 899), []byte{0xB4, 0x0D, 0x2D, 0x07, 0x59}},
		{"midnight", spiValues.NewPlcTIME_OF_DAYFromMillisecondsSinceMidnight(0), []byte{0xB4, 0x00, 0x00, 0x00, 0x00}},
		{"last hundredth", spiValues.NewPlcTIME_OF_DAY(time.Date(0, 0, 0, 23, 59, 59, 999_999_999, time.UTC)), []byte{0xB4, 0x17, 0x3B, 0x3B, 0x63}},
		{"LTIME_OF_DAY", spiValues.NewPlcLTIME_OF_DAY(time.Date(0, 0, 0, 13, 45, 30, 250_000_000, time.UTC)), []byte{0xB4, 0x0D, 0x2D, 0x1E, 0x19}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.wire, encodedOnTheWire(t, c.value))
		})
	}
}

// A concrete Date or Time read from a device is written back with the same octets.
func TestPlcValueToApplicationTag_dateAndTimeRoundTrip(t *testing.T) {
	for _, wire := range [][]byte{
		{0xA4, 0x7E, 0x0B, 0x18, 0x02},
		{0xA4, 0x7C, 0x02, 0x1D, 0x04},
		{0xB4, 0x0D, 0x2D, 0x1E, 0x19},
		{0xB4, 0x17, 0x3B, 0x3B, 0x63},
	} {
		t.Run(fmt.Sprintf("% X", wire), func(t *testing.T) {
			assert.Equal(t, wire, encodedOnTheWire(t, appTagFromTheWire(t, wire)))
		})
	}
}

// A Date or Time written to a Date Value or Time Value object reaches the
// WriteProperty request as a Date or Time application tag.
func TestBuildServiceRequest_dateAndTimeValuesGoOutAsDateAndTime(t *testing.T) {
	date := writtenApplicationTag(t, readWriteModel.BACnetObjectType_DATE_VALUE, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, nil,
		spiValues.NewPlcDATE(time.Date(2026, 11, 24, 0, 0, 0, 0, time.UTC)))
	assert.Implements(t, (*readWriteModel.BACnetApplicationTagDate)(nil), date)
	timeOfDay := writtenApplicationTag(t, readWriteModel.BACnetObjectType_TIME_VALUE, readWriteModel.BACnetPropertyIdentifier_PRESENT_VALUE, nil,
		spiValues.NewPlcTIME_OF_DAYFromMillisecondsSinceMidnight(1000))
	assert.Implements(t, (*readWriteModel.BACnetApplicationTagTime)(nil), timeOfDay)
}
