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

package s7

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

// This file complements DateAndTimeBcd_test.go: it covers the remaining BCD
// surface of the S7 driver, namely the hand written StaticHelper year codec,
// the second alarm type that embeds a BCD DateAndTime, the encode side range
// checks, and the downstream formatter that renders the decoded timestamp.

// TestParseSiemensYear_BCD exercises
// protocols/s7/readwrite/model/StaticHelper.go, the only place outside the
// generated DateAndTime model that passes utils.WithEncoding("BCD") to a
// numeric read. It mirrors plc4j's
// drivers/s7/.../readwrite/utils/StaticHelper.java#parseSiemensYear: the wire
// byte is two BCD digits, values below 90 are 2000s, the rest are 1900s.
func TestParseSiemensYear_BCD(t *testing.T) {
	tests := []struct {
		raw      byte
		expected uint16
	}{
		{raw: 0x00, expected: 2000},
		{raw: 0x09, expected: 2009},
		{raw: 0x24, expected: 2024},
		{raw: 0x89, expected: 2089},
		{raw: 0x90, expected: 1990},
		{raw: 0x91, expected: 1991},
		{raw: 0x99, expected: 1999},
	}
	for _, test := range tests {
		year, err := readWriteModel.ParseSiemensYear(context.Background(), utils.NewReadBufferByteBased([]byte{test.raw}))
		require.NoError(t, err, "raw 0x%02X", test.raw)
		assert.Equal(t, test.expected, year, "raw 0x%02X", test.raw)
	}

	// Without the BCD dispatch 0x90 would decode as 144 and land in the 1900s
	// branch as 2044; the table above pins the correct 1990.
	_, err := readWriteModel.ParseSiemensYear(context.Background(), utils.NewReadBufferByteBased([]byte{0xA0}))
	assert.Error(t, err, "0xA is not a valid BCD digit")
}

// TestSerializeSiemensYear_RoundTripsEveryParseableYear is the write-side twin
// of TestParseSiemensYear_BCD. Every year ParseSiemensYear can produce must be
// serializable back to the byte it came from.
//
// This is the regression guard for the `year > 2000` boundary bug plc4j still
// carries (StaticHelper.java#serializeSiemensYear): with `>` the year 2000 takes
// the 1900 branch and asks BCD to encode 100 in two digits, which the encoder
// rejects - so a timestamp that parses cannot be written back.
func TestSerializeSiemensYear_RoundTripsEveryParseableYear(t *testing.T) {
	for raw := 0; raw <= 0x99; raw++ {
		if raw&0x0F > 9 || (raw>>4)&0x0F > 9 {
			continue // not a valid BCD byte
		}
		year, err := readWriteModel.ParseSiemensYear(context.Background(), utils.NewReadBufferByteBased([]byte{byte(raw)}))
		require.NoError(t, err, "raw 0x%02X", raw)

		writeBuffer := utils.NewWriteBufferByteBased()
		dateTime := spiValues.NewPlcDATE_AND_TIME(time.Date(int(year), time.August, 19, 10, 30, 45, 0, time.UTC))
		require.NoError(t,
			readWriteModel.SerializeSiemensYear(context.Background(), writeBuffer, dateTime),
			"year %d (raw 0x%02X)", year, raw)
		assert.Equal(t, []byte{byte(raw)}, writeBuffer.GetBytes(), "year %d", year)
	}
}

// TestSerializeSiemensYear_RejectsUnrepresentableYears pins the range guard: the
// single BCD byte only spans 1990-2089, and anything outside must fail loudly
// rather than put a wrong year on the wire.
func TestSerializeSiemensYear_RejectsUnrepresentableYears(t *testing.T) {
	for _, year := range []int{1899, 1970, 1989, 2090, 2100} {
		writeBuffer := utils.NewWriteBufferByteBased()
		dateTime := spiValues.NewPlcDATE_AND_TIME(time.Date(year, time.August, 19, 10, 30, 45, 0, time.UTC))
		assert.Error(t,
			readWriteModel.SerializeSiemensYear(context.Background(), writeBuffer, dateTime),
			"year %d", year)
	}
}

// TestAlarmMessageAckPushType_EmbeddedDateAndTimeBCD is the ack twin of the
// push type test: both carry a nested BCD DateAndTime as their first field.
func TestAlarmMessageAckPushType_EmbeddedDateAndTimeBCD(t *testing.T) {
	// 8 bytes DateAndTime (2024-08-19T10:30:45.123, dow 2), functionId,
	// numberOfObjects (0, so no message objects follow).
	wireBytes := []byte{0x24, 0x08, 0x19, 0x10, 0x30, 0x45, 0x12, 0x32, 0x12, 0x00}

	ackMessage, err := readWriteModel.AlarmMessageAckPushTypeParse(context.Background(), wireBytes)
	require.NoError(t, err)
	require.NotNil(t, ackMessage)

	assert.Equal(t, uint8(0x12), ackMessage.GetFunctionId())
	assert.Equal(t, uint8(0), ackMessage.GetNumberOfObjects())
	assert.Equal(t,
		readWriteModel.NewDateAndTime(24, 8, 19, 10, 30, 45, 123, 2),
		ackMessage.GetTimeStamp())

	serialized, err := ackMessage.Serialize()
	require.NoError(t, err)
	assert.Equal(t, wireBytes, serialized)
}

// TestAlarm_UnsetTimestampFailsTheWholeMessage documents a real behaviour change
// that the strict nibble validation introduces on live traffic.
//
// Some CPUs emit an "unset" timestamp as eight 0xFF bytes. Before the encoding
// dispatch existed those nibbles were read as plain binary and decoded to
// garbage, but the enclosing S7 message still parsed. Now the first 0xF nibble
// aborts DateAndTimeParse, so the whole AlarmMessagePushType / -AckPushType
// parse fails. This matches plc4j (EncodingBCD.decodeInt throws on the same
// input), so the strictness is intentional, but the consequence is worth
// pinning: internal/s7/MessageCodec.go#Receive logs the parse error at WARN and
// returns (nil, nil), i.e. the offending packet is dropped and the connection
// stays up - it is not surfaced to the subscriber and it does not kill the
// codec loop.
func TestAlarm_UnsetTimestampFailsTheWholeMessage(t *testing.T) {
	unsetTimeStamp := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

	_, err := readWriteModel.DateAndTimeParse(context.Background(), unsetTimeStamp)
	assert.ErrorContains(t, err, "invalid BCD digit")

	pushBytes := append(append([]byte{}, unsetTimeStamp...), 0x00, 0x00)
	_, err = readWriteModel.AlarmMessagePushTypeParse(context.Background(), pushBytes)
	assert.ErrorContains(t, err, "invalid BCD digit",
		"the nested timestamp must fail the enclosing alarm message, not decode to garbage")

	ackBytes := append(append([]byte{}, unsetTimeStamp...), 0x12, 0x00)
	_, err = readWriteModel.AlarmMessageAckPushTypeParse(context.Background(), ackBytes)
	assert.ErrorContains(t, err, "invalid BCD digit")
}

// TestS7PayloadUserDataItemClk_DateAndTimeBCD covers the timestamp carried by
// the clock service responses. Those types can only be parsed through their
// discriminated parent, so the timestamp is exercised through the very same
// DateAndTimeParseWithBuffer call the generated Clk model makes, plus a check
// that the model actually holds the decoded value.
func TestS7PayloadUserDataItemClk_DateAndTimeBCD(t *testing.T) {
	timeStamp, err := readWriteModel.DateAndTimeParse(context.Background(),
		[]byte{0x24, 0x08, 0x19, 0x10, 0x30, 0x45, 0x12, 0x32})
	require.NoError(t, err)

	clkResponse := readWriteModel.NewS7PayloadUserDataItemClkResponse(
		readWriteModel.DataTransportErrorCode_OK,
		readWriteModel.DataTransportSize_OCTET_STRING,
		10,
		0xFF,
		0x20,
		timeStamp,
	)
	require.NotNil(t, clkResponse)
	assert.Equal(t, uint8(24), clkResponse.GetTimeStamp().GetYear())
	assert.Equal(t, uint8(19), clkResponse.GetTimeStamp().GetDay())
	assert.Equal(t, uint16(123), clkResponse.GetTimeStamp().GetMsec())
	assert.Equal(t, uint8(2), clkResponse.GetTimeStamp().GetDow())
}

// TestDateAndTime_BCDRejectsUnencodableValues pins the encode side range check.
// plc4j's EncodingBCD.encodeInt throws for a value that needs more digits than
// the field provides; the Go port returns an error instead of silently
// truncating.
func TestDateAndTime_BCDRejectsUnencodableValues(t *testing.T) {
	tests := []struct {
		name      string
		dateTime  readWriteModel.DateAndTime
		wantError bool
	}{
		{name: "year 99 fits two digits", dateTime: readWriteModel.NewDateAndTime(99, 8, 19, 10, 30, 45, 123, 2)},
		{name: "year 100 needs three digits", dateTime: readWriteModel.NewDateAndTime(100, 8, 19, 10, 30, 45, 123, 2), wantError: true},
		{name: "msec 999 fits three digits", dateTime: readWriteModel.NewDateAndTime(24, 8, 19, 10, 30, 45, 999, 2)},
		{name: "msec 1000 needs four digits", dateTime: readWriteModel.NewDateAndTime(24, 8, 19, 10, 30, 45, 1000, 2), wantError: true},
		{name: "dow 9 fits one digit", dateTime: readWriteModel.NewDateAndTime(24, 8, 19, 10, 30, 45, 123, 9)},
		{name: "dow 10 needs two digits", dateTime: readWriteModel.NewDateAndTime(24, 8, 19, 10, 30, 45, 123, 10), wantError: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.dateTime.Serialize()
			if test.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestDataItem_DateAndTimeIsStillHalfBCD is a characterization test for a known
// defect that this change makes visible but cannot fix from Go.
//
// s7.mspec declares all eight DATE_AND_TIME data-item fields as
// encoding='"BCD"', but protocols/s7/readwrite/model/DataItem.go is generated by
// the un-migrated data-io template, which used to drop the encoding option for
// numeric fields. Only `year` carries it, because it is a manual field routed
// through SerializeSiemensYear, which passes utils.WithEncoding("BCD") by hand.
// The result is a frame that is BCD in its first byte and raw binary in the
// other seven.
//
// The root cause was fixed in
// code-generation/language/go/src/main/java/org/apache/plc4x/language/go/GoLanguageTemplateHelper.java
// (see getEncodingOption) and DataItem.go has since been regenerated, so every
// field of the data item now carries the encoding option.
func TestDataItem_DateAndTimeSerializesAsBCD(t *testing.T) {
	dateTime := spiValues.NewPlcDATE_AND_TIME(time.Date(2024, time.August, 19, 14, 35, 7, 123000000, time.UTC))

	serialized, err := readWriteModel.DataItemSerialize(dateTime, "IEC61131_DATE_AND_TIME", readWriteModel.ControllerType_S7_1200, 0)
	require.NoError(t, err)

	// Every field is BCD now: 2024-08-19 14:35:07.123 reads straight off the wire
	// as the decimal digits it denotes. Before the regeneration only year was BCD
	// and the rest was raw binary (day 19 -> 0x13, hour 14 -> 0x0E, msec 123 -> 0x07B).
	//
	// CAUTION on the trailing nibble: it is the dow, and it still carries the
	// ISO-8601 numbering (1 == Monday .. 7 == Sunday) that PlcDATE_AND_TIME shares
	// with plc4j and KNX DPT 19.001. S7 numbers the same field 1 == Sunday ..
	// 7 == Saturday, so a real PLC expects 2 for this Monday, not the 1 below.
	// Rotating into the Siemens numbering is a separate S7 layer fix; see
	// TestDataItem_DateAndTimeDayOfWeekIsSiemensNumbered once that lands.
	want := []byte{0x24, 0x08, 0x19, 0x14, 0x35, 0x07, 0x12, 0x31}
	assert.Equal(t, want, serialized)

	// Independent of the encoding, the read side of this data item is dead: the
	// generated parser assigns all eight fields to _ and falls through to
	// "unsupported type" (DataItem.go, end of DataItemParseWithBuffer).
	_, err = readWriteModel.DataItemParse(context.Background(), serialized,
		"IEC61131_DATE_AND_TIME", readWriteModel.ControllerType_S7_1200, 0)
	assert.ErrorContains(t, err, "unsupported type")
}

// TestFormatAlarmDateAndTime_FromWireBytes closes the loop from the wire bytes
// all the way to the string the alarm subscription surfaces to the user.
func TestFormatAlarmDateAndTime_FromWireBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected string
	}{
		{
			name:     "2024",
			bytes:    []byte{0x24, 0x08, 0x19, 0x10, 0x30, 0x45, 0x12, 0x32},
			expected: "2024-08-19T10:30:45.123",
		},
		{
			// The vector Subscriber_test.go builds in memory.
			name:     "subscriber vector",
			bytes:    []byte{0x24, 0x08, 0x18, 0x10, 0x30, 0x00, 0x12, 0x31},
			expected: "2024-08-18T10:30:00.123",
		},
		{
			// >= 90 selects the 1900s branch.
			name:     "1991",
			bytes:    []byte{0x91, 0x12, 0x31, 0x23, 0x59, 0x59, 0x99, 0x97},
			expected: "1991-12-31T23:59:59.999",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dateAndTime, err := readWriteModel.DateAndTimeParse(context.Background(), test.bytes)
			require.NoError(t, err)
			assert.Equal(t, test.expected, formatAlarmDateAndTime(dateAndTime))
		})
	}
}
