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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
)

// The Siemens DATE_AND_TIME wire format is a plain BCD nibble stream, most
// significant digit first (s7.mspec lines 307-316):
//
//	year/month/day/hour/minutes/seconds  uint 8  encoding='"BCD"'
//	msec                                 uint 12 encoding='"BCD"'
//	dow                                  uint 4  encoding='"BCD"'
//
// msec and dow are the only BCD fields in the whole repository whose bit length
// is not a multiple of 8, which is exactly where plc4j gets it wrong. See the
// divergence note on decodeBCD/encodeBCD in spi/utils/encodingBCD.go.
//
// dateAndTimeWireBytes decodes to 2024-08-19 14:35:07.123, day of week 4.
var dateAndTimeWireBytes = []byte{0x24, 0x08, 0x19, 0x14, 0x35, 0x07, 0x12, 0x34}

func TestDateAndTime_BCDDecode(t *testing.T) {
	dateAndTime, err := model.DateAndTimeParse(context.Background(), dateAndTimeWireBytes)
	require.NoError(t, err)
	require.NotNil(t, dateAndTime)

	assert.Equal(t, uint8(24), dateAndTime.GetYear(), "year: 0x24 is BCD 24, i.e. 2024")
	assert.Equal(t, uint8(8), dateAndTime.GetMonth(), "month: 0x08 is BCD 8")
	assert.Equal(t, uint8(19), dateAndTime.GetDay(), "day: 0x19 is BCD 19")
	assert.Equal(t, uint8(14), dateAndTime.GetHour(), "hour: 0x14 is BCD 14")
	assert.Equal(t, uint8(35), dateAndTime.GetMinutes(), "minutes: 0x35 is BCD 35")
	assert.Equal(t, uint8(7), dateAndTime.GetSeconds(), "seconds: 0x07 is BCD 7")

	// The two fields plc4j corrupts. plc4j returns msec=12 (it drops the third
	// nibble) and dow=0 (it reads the wrong nibble of 0x34) for this very input.
	assert.Equal(t, uint16(123), dateAndTime.GetMsec(),
		"msec: 12 bits, nibbles 1,2 of 0x12 plus the high nibble of 0x34 (plc4j wrongly yields 12)")
	assert.Equal(t, uint8(4), dateAndTime.GetDow(),
		"dow: 4 bits, the low nibble of 0x34 (plc4j wrongly yields 0)")
}

func TestDateAndTime_BCDRoundTripFromWire(t *testing.T) {
	dateAndTime, err := model.DateAndTimeParse(context.Background(), dateAndTimeWireBytes)
	require.NoError(t, err)

	serialized, err := dateAndTime.Serialize()
	require.NoError(t, err)
	assert.Equal(t, dateAndTimeWireBytes, serialized,
		"re-serializing the parsed DateAndTime must reproduce the original bytes")
}

func TestDateAndTime_BCDRoundTripFromValues(t *testing.T) {
	// Encoding the semantically correct field values must produce exactly the
	// wire bytes. plc4j serializes these as ...2300 instead of ...1234.
	dateAndTime := model.NewDateAndTime(24, 8, 19, 14, 35, 7, 123, 4)

	serialized, err := dateAndTime.Serialize()
	require.NoError(t, err)
	assert.Equal(t, dateAndTimeWireBytes, serialized)

	reparsed, err := model.DateAndTimeParse(context.Background(), serialized)
	require.NoError(t, err)
	assert.Equal(t, uint16(123), reparsed.GetMsec())
	assert.Equal(t, uint8(4), reparsed.GetDow())
}

func TestDateAndTime_BCDRoundTripVectors(t *testing.T) {
	tests := []struct {
		name                                          string
		bytes                                         []byte
		year, month, day, hour, minutes, seconds, dow uint8
		msec                                          uint16
	}{
		{
			name:  "all zeroes",
			bytes: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:  "reference timestamp",
			bytes: dateAndTimeWireBytes,
			year:  24, month: 8, day: 19, hour: 14, minutes: 35, seconds: 7, msec: 123, dow: 4,
		},
		{
			name:  "maximum msec and dow",
			bytes: []byte{0x99, 0x12, 0x31, 0x23, 0x59, 0x59, 0x99, 0x97},
			year:  99, month: 12, day: 31, hour: 23, minutes: 59, seconds: 59, msec: 999, dow: 7,
		},
		{
			name:  "msec below 100 keeps its leading zero nibble",
			bytes: []byte{0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x51},
			year:  1, month: 1, day: 1, msec: 5, dow: 1,
		},
		{
			name:  "msec 100 and dow 1",
			bytes: []byte{0x24, 0x08, 0x19, 0x00, 0x00, 0x00, 0x10, 0x01},
			year:  24, month: 8, day: 19, msec: 100, dow: 1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parsed, err := model.DateAndTimeParse(context.Background(), test.bytes)
			require.NoError(t, err)

			assert.Equal(t, test.year, parsed.GetYear())
			assert.Equal(t, test.month, parsed.GetMonth())
			assert.Equal(t, test.day, parsed.GetDay())
			assert.Equal(t, test.hour, parsed.GetHour())
			assert.Equal(t, test.minutes, parsed.GetMinutes())
			assert.Equal(t, test.seconds, parsed.GetSeconds())
			assert.Equal(t, test.msec, parsed.GetMsec())
			assert.Equal(t, test.dow, parsed.GetDow())

			serialized, err := parsed.Serialize()
			require.NoError(t, err)
			assert.Equal(t, test.bytes, serialized)

			// And the same values built by hand must serialize identically.
			built, err := model.NewDateAndTimeBuilder().
				WithMandatoryFields(test.year, test.month, test.day, test.hour, test.minutes, test.seconds, test.msec, test.dow).
				Build()
			require.NoError(t, err)
			builtBytes, err := built.Serialize()
			require.NoError(t, err)
			assert.Equal(t, test.bytes, builtBytes)
		})
	}
}

// TestDateAndTime_BCDRejectsInvalidNibbles guards that a non-decimal nibble is
// reported rather than silently mangled.
func TestDateAndTime_BCDRejectsInvalidNibbles(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{name: "invalid year nibble", bytes: []byte{0xA4, 0x08, 0x19, 0x14, 0x35, 0x07, 0x12, 0x34}},
		{name: "invalid msec nibble", bytes: []byte{0x24, 0x08, 0x19, 0x14, 0x35, 0x07, 0x1F, 0x34}},
		{name: "invalid dow nibble", bytes: []byte{0x24, 0x08, 0x19, 0x14, 0x35, 0x07, 0x12, 0x3B}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := model.DateAndTimeParse(context.Background(), test.bytes)
			assert.Error(t, err)
		})
	}
}

// TestAlarmMessagePushType_EmbeddedDateAndTimeBCD covers a second S7 type that
// carries the BCD DateAndTime, to prove the encoding survives being nested
// inside a complex field rather than parsed standalone.
func TestAlarmMessagePushType_EmbeddedDateAndTimeBCD(t *testing.T) {
	// 8 bytes DateAndTime, then functionId 0x00 and numberOfObjects 0x00.
	wireBytes := append(append([]byte{}, dateAndTimeWireBytes...), 0x00, 0x00)

	alarmMessage, err := model.AlarmMessagePushTypeParse(context.Background(), wireBytes)
	require.NoError(t, err)
	require.NotNil(t, alarmMessage)

	timeStamp := alarmMessage.GetTimeStamp()
	require.NotNil(t, timeStamp)
	assert.Equal(t, uint8(24), timeStamp.GetYear())
	assert.Equal(t, uint8(8), timeStamp.GetMonth())
	assert.Equal(t, uint8(19), timeStamp.GetDay())
	assert.Equal(t, uint8(14), timeStamp.GetHour())
	assert.Equal(t, uint8(35), timeStamp.GetMinutes())
	assert.Equal(t, uint8(7), timeStamp.GetSeconds())
	assert.Equal(t, uint16(123), timeStamp.GetMsec())
	assert.Equal(t, uint8(4), timeStamp.GetDow())

	assert.Equal(t, uint8(0), alarmMessage.GetFunctionId())
	assert.Equal(t, uint8(0), alarmMessage.GetNumberOfObjects())

	serialized, err := alarmMessage.Serialize()
	require.NoError(t, err)
	assert.Equal(t, wireBytes, serialized)
}
