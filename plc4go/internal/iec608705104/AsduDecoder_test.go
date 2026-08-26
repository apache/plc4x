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

package iec608705104

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/iec608705104/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
// Frame builders
////////////////////////////////////////////////////////////////////////////////////////////////////

// littleEndianUint16 spells a 16 bit field the way the APCI and the ASDU do.
func littleEndianUint16(value uint16) []byte {
	return []byte{byte(value), byte(value >> 8)}
}

// informationObjectBytes is one information object: its three octet address, low octet first,
// followed by the payload of its type.
func informationObjectBytes(address uint32, payload ...byte) []byte {
	object := []byte{byte(address), byte(address >> 8), byte(address >> 16)}
	return append(object, payload...)
}

// asduBytes assembles an ASDU around a set of information objects, with the structure qualifier
// clear (each object carries its own address) and no originator.
func asduBytes(typeIdentification byte, causeOfTransmission byte, asduAddress uint16, objects ...[]byte) []byte {
	asdu := []byte{typeIdentification, byte(len(objects)), causeOfTransmission, 0x00}
	asdu = append(asdu, littleEndianUint16(asduAddress)...)
	for _, object := range objects {
		asdu = append(asdu, object...)
	}
	return asdu
}

// iFormatFrame wraps an ASDU in an I-format APCI. Both sequence numbers sit in the control field
// shifted left by one, the low bit being the format discriminator.
func iFormatFrame(sendSequenceNo uint16, receiveSequenceNo uint16, asdu []byte) []byte {
	body := append(littleEndianUint16(sendSequenceNo<<1), littleEndianUint16(receiveSequenceNo<<1)...)
	body = append(body, asdu...)
	return append([]byte{startByte, byte(len(body))}, body...)
}

// parseApdu parses a frame the way the codec would.
func parseApdu(t *testing.T, frame []byte) readWriteModel.APDU {
	t.Helper()
	apdu, err := readWriteModel.APDUParse[readWriteModel.APDU](t.Context(), frame)
	require.NoError(t, err)
	return apdu
}

// parseHexApdu parses one of the frames the parser-serializer testsuite captured off a real station.
func parseHexApdu(t *testing.T, frame string) readWriteModel.APDU {
	t.Helper()
	theBytes, err := hex.DecodeString(frame)
	require.NoError(t, err)
	return parseApdu(t, theBytes)
}

// decodeFirstPoint decodes the first (and usually only) information object of an I-format frame.
func decodeFirstPoint(t *testing.T, frame []byte) (apiValues.PlcValue, string, apiModel.PlcResponseCode) {
	t.Helper()
	return decodePointAt(t, frame, 0)
}

// decodePointAt decodes one information object of an I-format frame by its position in the ASDU.
func decodePointAt(t *testing.T, frame []byte, index int) (apiValues.PlcValue, string, apiModel.PlcResponseCode) {
	t.Helper()
	apdu := parseApdu(t, frame)
	iFormat, ok := apdu.(readWriteModel.APDUIFormat)
	require.True(t, ok, "%T is not an I-format frame", apdu)
	objects := iFormat.GetAsdu().GetInformationObjects()
	require.Greater(t, len(objects), index)
	return decodePoint(iFormat.GetAsdu(), objects[index])
}

// nested walks into a struct value along a path of keys.
func nested(t *testing.T, value apiValues.PlcValue, keys ...string) apiValues.PlcValue {
	t.Helper()
	current := value
	for _, key := range keys {
		require.True(t, current.IsStruct(), "%v is not a struct, cannot look up %q", current, key)
		require.True(t, current.HasKey(key), "no key %q in %v", key, current.GetKeys())
		current = current.GetValue(key)
	}
	return current
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// The envelope every event carries
////////////////////////////////////////////////////////////////////////////////////////////////////

// A wildcard tag says nothing about where a report came from, so every event has to carry the
// concrete point, the type identification and the cause of transmission. plc4j carries none of it.
func TestDecodePoint_Envelope(t *testing.T) {
	// A single point information at ASDU 10, information object 13, reported spontaneously.
	frame := iFormatFrame(0, 0, asduBytes(0x01, 3, 10, informationObjectBytes(13, 0x01)))
	value, _, code := decodeFirstPoint(t, frame)

	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.Equal(t, uint16(10), nested(t, value, fieldAsduAddress).GetUint16())
	assert.Equal(t, uint32(13), nested(t, value, fieldObjectAddress).GetUint32())
	assert.Equal(t, uint8(0x01), nested(t, value, fieldTypeIdentification).GetUint8())
	assert.Equal(t, "SINGLE_POINT_INFORMATION", nested(t, value, fieldTypeName).GetString())
	assert.Equal(t, "SPONTANEOUS_SPONT", nested(t, value, fieldCauseOfTransmission).GetString())
	assert.False(t, nested(t, value, fieldTest).GetBool())
	assert.False(t, nested(t, value, fieldNegative).GetBool())
}

// The test bit says the station is exercising the link rather than reporting the process. A value
// which arrives under it must be visible as such.
func TestDecodePoint_CarriesTheTestBit(t *testing.T) {
	// Bit 7 of the cause octet is the test bit, bit 6 the negative-confirmation bit.
	frame := iFormatFrame(0, 0, asduBytes(0x01, 0x80|0x40|7, 10, informationObjectBytes(13, 0x01)))
	value, _, code := decodeFirstPoint(t, frame)

	assert.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.True(t, nested(t, value, fieldTest).GetBool())
	assert.True(t, nested(t, value, fieldNegative).GetBool())
}

// An event value is a struct rather than a bare reading precisely so that the quality cannot be
// walked past by accident: asking a flagged-invalid measurement for its boolean has to fail loudly.
func TestDecodePoint_ValueIsNotMistakableForAScalar(t *testing.T) {
	frame := iFormatFrame(0, 0, asduBytes(0x01, 3, 10, informationObjectBytes(13, 0x81)))
	value, _, _ := decodeFirstPoint(t, frame)

	assert.False(t, value.IsBool(), "the event value must not look like a plain boolean")
	assert.False(t, value.IsUint8(), "nor like a number")
	assert.True(t, value.IsStruct())
	assert.True(t, nested(t, value, fieldValue).GetBool(), "the reading itself is still there")
	assert.True(t, nested(t, value, fieldQuality, "invalid").GetBool())
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Quality descriptors
////////////////////////////////////////////////////////////////////////////////////////////////////

// The SIQ quality bits, one at a time. The wire order is invalid, not-topical, substituted, blocked
// from the most significant bit down, with the status in bit 0.
func TestDecodePoint_SinglePointQuality(t *testing.T) {
	tests := []struct {
		name string
		siq  byte
		flag string
		on   bool
	}{
		{name: "clean and off", siq: 0x00, on: false},
		{name: "clean and on", siq: 0x01, on: true},
		{name: "invalid", siq: 0x80, flag: "invalid"},
		{name: "not topical", siq: 0x40, flag: "notTopical"},
		{name: "substituted", siq: 0x20, flag: "substituted"},
		{name: "blocked", siq: 0x10, flag: "blocked"},
		{name: "invalid and on", siq: 0x81, flag: "invalid", on: true},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			frame := iFormatFrame(0, 0, asduBytes(0x01, 3, 10, informationObjectBytes(13, testCase.siq)))
			value, _, code := decodeFirstPoint(t, frame)

			require.Equal(t, apiModel.PlcResponseCode_OK, code)
			assert.Equal(t, testCase.on, nested(t, value, fieldValue).GetBool())
			for _, flag := range []string{"invalid", "notTopical", "substituted", "blocked"} {
				assert.Equal(t, flag == testCase.flag, nested(t, value, fieldQuality, flag).GetBool(), "quality flag %s", flag)
			}
		})
	}
}

// The QDS of a measurement, which adds the overflow bit to the four SIQ flags.
func TestDecodePoint_MeasurementQuality(t *testing.T) {
	tests := []struct {
		name string
		qds  byte
		flag string
	}{
		{name: "clean", qds: 0x00},
		{name: "invalid", qds: 0x80, flag: "invalid"},
		{name: "not topical", qds: 0x40, flag: "notTopical"},
		{name: "substituted", qds: 0x20, flag: "substituted"},
		{name: "blocked", qds: 0x10, flag: "blocked"},
		{name: "overflow", qds: 0x01, flag: "overflow"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// A scaled measurement of 1234 with the quality descriptor under test.
			object := informationObjectBytes(13, 0xD2, 0x04, testCase.qds)
			frame := iFormatFrame(0, 0, asduBytes(0x0B, 3, 10, object))
			value, _, code := decodeFirstPoint(t, frame)

			require.Equal(t, apiModel.PlcResponseCode_OK, code)
			assert.Equal(t, int16(1234), nested(t, value, fieldValue).GetInt16())
			for _, flag := range []string{"invalid", "notTopical", "substituted", "blocked", "overflow"} {
				assert.Equal(t, flag == testCase.flag, nested(t, value, fieldQuality, flag).GetBool(), "quality flag %s", flag)
			}
		})
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Value conversions
////////////////////////////////////////////////////////////////////////////////////////////////////

// A scaled value is a signed 16 bit integer, so the negative half of its range has to come out
// negative.
func TestDecodePoint_ScaledValueIsSigned(t *testing.T) {
	object := informationObjectBytes(13, 0xFF, 0xFF, 0x00)
	frame := iFormatFrame(0, 0, asduBytes(0x0B, 3, 10, object))
	value, _, _ := decodeFirstPoint(t, frame)

	assert.Equal(t, int16(-1), nested(t, value, fieldValue).GetInt16())
}

// A normalized value is a two's complement fraction of full scale in [-1, 1). The generated model
// hands out the raw 16 bits, and plc4j passes those straight on as an unsigned integer - so a
// measurement at minus full scale reads as 32768 there instead of as -1.
func TestDecodePoint_NormalizedValueIsAFraction(t *testing.T) {
	tests := []struct {
		name  string
		nva   []byte
		want  float32
		wantR int16
	}{
		{name: "zero", nva: []byte{0x00, 0x00}, want: 0, wantR: 0},
		{name: "minus full scale", nva: []byte{0x00, 0x80}, want: -1, wantR: -32768},
		{name: "half of full scale", nva: []byte{0x00, 0x40}, want: 0.5, wantR: 16384},
		{name: "minus half of full scale", nva: []byte{0x00, 0xC0}, want: -0.5, wantR: -16384},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			object := informationObjectBytes(13, append(testCase.nva, 0x00)...)
			frame := iFormatFrame(0, 0, asduBytes(0x09, 3, 10, object))
			value, _, _ := decodeFirstPoint(t, frame)

			assert.InDelta(t, testCase.want, nested(t, value, fieldValue).GetFloat32(), 0.0001)
			assert.Equal(t, testCase.wantR, nested(t, value, fieldRawValue).GetInt16())
		})
	}
}

// A step position is a seven bit two's complement number, so a tap below the neutral position has to
// come out negative. plc4j hands out the raw byte, which turns tap -63 into 65.
func TestDecodePoint_StepPositionIsSigned(t *testing.T) {
	tests := []struct {
		name      string
		vti       byte
		want      int8
		transient bool
	}{
		{name: "the neutral position", vti: 0x00, want: 0},
		{name: "the highest position", vti: 0x3F, want: 63},
		{name: "one below neutral", vti: 0x7F, want: -1},
		{name: "the lowest position", vti: 0x40, want: -64},
		{name: "in transit", vti: 0xC1, want: -63, transient: true},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			object := informationObjectBytes(13, testCase.vti, 0x00)
			frame := iFormatFrame(0, 0, asduBytes(0x05, 3, 10, object))
			value, _, code := decodeFirstPoint(t, frame)

			require.Equal(t, apiModel.PlcResponseCode_OK, code)
			assert.Equal(t, testCase.want, nested(t, value, fieldValue).GetInt8())
			// The model's uint7 field already drops the transient bit, so the raw value is the
			// seven bits of the step position rather than the whole octet.
			assert.Equal(t, testCase.vti&0x7F, nested(t, value, fieldRawValue).GetUint8())
			assert.Equal(t, testCase.transient, nested(t, value, fieldQuality, "transientState").GetBool())
		})
	}
}

// The counter qualifier's top bit is IV - set when the reading is *invalid*. The mspec calls that
// getter GetCounterValid, and plc4j publishes it under the name "counterValid", so a consumer there
// reading counterValid=true is looking at a counter the station has just declared broken.
func TestDecodePoint_IntegratedTotalsQualityIsNotInverted(t *testing.T) {
	// A counter of 10 whose qualifier octet has IV set and a sequence number of 1.
	object := informationObjectBytes(13, 0x0A, 0x00, 0x00, 0x00, 0x81)
	frame := iFormatFrame(0, 0, asduBytes(0x0F, 3, 10, object))
	value, _, code := decodeFirstPoint(t, frame)

	require.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.Equal(t, uint32(10), nested(t, value, fieldValue).GetUint32())
	assert.True(t, nested(t, value, fieldQuality, "invalid").GetBool(), "the IV bit means invalid, not valid")
	assert.False(t, nested(t, value, fieldQuality, "counterAdjusted").GetBool())
	assert.False(t, nested(t, value, fieldQuality, "carry").GetBool())
	assert.Equal(t, uint8(1), nested(t, value, fieldQuality, "sequenceNumber").GetUint8())
}

// A double point carries a two bit state in which two of the four codes mean "the station cannot
// tell". plc4j renders it as a list of the two bits, which makes indeterminate look like a reading.
func TestDecodePoint_DoublePointStates(t *testing.T) {
	tests := []struct {
		name          string
		diq           byte
		wantState     string
		indeterminate bool
	}{
		{name: "both contacts open", diq: 0x00, wantState: "INDETERMINATE", indeterminate: true},
		{name: "off", diq: 0x01, wantState: "OFF"},
		{name: "on", diq: 0x02, wantState: "ON"},
		{name: "both contacts closed", diq: 0x03, wantState: "INDETERMINATE", indeterminate: true},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			frame := iFormatFrame(0, 0, asduBytes(0x03, 3, 10, informationObjectBytes(13, testCase.diq)))
			value, _, code := decodeFirstPoint(t, frame)

			require.Equal(t, apiModel.PlcResponseCode_OK, code)
			assert.Equal(t, testCase.diq, nested(t, value, fieldValue).GetUint8())
			assert.Equal(t, testCase.wantState, nested(t, value, fieldState).GetString())
			assert.Equal(t, testCase.indeterminate, nested(t, value, fieldQuality, "indeterminate").GetBool())
		})
	}
}

// A bit string and a status-change detection are both dropped on the floor by plc4j (its
// processBinaryStateInformation and processStatusChangeDetection return null).
func TestDecodePoint_BitStrings(t *testing.T) {
	t.Run("a 32 bit string", func(t *testing.T) {
		object := informationObjectBytes(13, 0x78, 0x56, 0x34, 0x12, 0x00)
		frame := iFormatFrame(0, 0, asduBytes(0x07, 3, 10, object))
		value, _, code := decodeFirstPoint(t, frame)

		require.Equal(t, apiModel.PlcResponseCode_OK, code)
		assert.Equal(t, uint32(0x12345678), nested(t, value, fieldValue).GetUint32())
	})
	t.Run("a status change detection", func(t *testing.T) {
		object := informationObjectBytes(13, 0x78, 0x56, 0x34, 0x12, 0x00)
		frame := iFormatFrame(0, 0, asduBytes(0x14, 3, 10, object))
		value, _, code := decodeFirstPoint(t, frame)

		require.Equal(t, apiModel.PlcResponseCode_OK, code)
		assert.Equal(t, uint32(0x12345678), nested(t, value, fieldValue).GetUint32())
		assert.Equal(t, uint16(0x5678), nested(t, value, "status").GetUint16())
		assert.Equal(t, uint16(0x1234), nested(t, value, "changeDetection").GetUint16())
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Time tags
////////////////////////////////////////////////////////////////////////////////////////////////////

// A CP56Time2a is an absolute time the station read off its own clock. It is handed over unshifted:
// the protocol carries no timezone, and plc4j's reinterpretation in the client's default timezone
// silently moves every timestamp whenever the two sit in different zones.
func TestDecodePoint_SevenOctetTimeTag(t *testing.T) {
	// Straight out of the parser-serializer testsuite: a single point information with a CP56Time2a.
	apdu := parseHexApdu(t, "68154c0008001e0103000a000d000001c75d170884070d")
	iFormat, ok := apdu.(readWriteModel.APDUIFormat)
	require.True(t, ok)
	objects := iFormat.GetAsdu().GetInformationObjects()
	require.Len(t, objects, 1)

	value, _, code := decodePoint(iFormat.GetAsdu(), objects[0])

	require.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.Equal(t, uint32(13), nested(t, value, fieldObjectAddress).GetUint32())
	assert.True(t, nested(t, value, fieldValue).GetBool())
	assert.Equal(t,
		time.Date(2013, time.July, 4, 8, 23, 24, 7*int(time.Millisecond), time.UTC),
		nested(t, value, fieldTimestamp, "dateTime").GetDateTime())
	assert.False(t, nested(t, value, fieldTimestamp, "invalid").GetBool())
	assert.False(t, nested(t, value, fieldTimestamp, "substituted").GetBool())
	assert.False(t, nested(t, value, fieldTimestamp, "daylightSaving").GetBool())
	assert.Equal(t, uint8(4), nested(t, value, fieldTimestamp, "dayOfWeek").GetUint8())
}

// A CP24Time2a carries only the minute and the millisecond within it, so there is no absolute time
// in it at all. plc4j fills the year, month, day and hour in from the client's clock, inventing four
// fields' worth of data; the parts the wire really carries are reported as they are instead.
func TestDecodePoint_ThreeOctetTimeTag(t *testing.T) {
	// A single point information with a CP24Time2a: 12345 ms into minute 42, the time flagged valid.
	object := informationObjectBytes(13, 0x01, 0x39, 0x30, 0x2A)
	frame := iFormatFrame(0, 0, asduBytes(0x02, 3, 10, object))
	value, _, code := decodeFirstPoint(t, frame)

	require.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.True(t, nested(t, value, fieldValue).GetBool())
	assert.Equal(t, uint8(42), nested(t, value, fieldTimestamp, "minutes").GetUint8())
	assert.Equal(t, uint8(12), nested(t, value, fieldTimestamp, "seconds").GetUint8())
	assert.Equal(t, uint16(12345), nested(t, value, fieldTimestamp, "milliseconds").GetUint16())
	assert.False(t, nested(t, value, fieldTimestamp, "invalid").GetBool())
	assert.False(t, value.HasKey("dateTime"), "a CP24Time2a is not an absolute time")
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// System information and commands
////////////////////////////////////////////////////////////////////////////////////////////////////

// A command mirrored back as an activation confirmation. plc4j returns null for all three command
// types, which its own caller then hands on as the value of the event.
func TestDecodePoint_SingleCommand(t *testing.T) {
	apdu := parseHexApdu(t, "680e040042002d0106000a0002000001")
	iFormat, ok := apdu.(readWriteModel.APDUIFormat)
	require.True(t, ok)
	objects := iFormat.GetAsdu().GetInformationObjects()
	require.Len(t, objects, 1)

	value, _, code := decodePoint(iFormat.GetAsdu(), objects[0])

	require.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.Equal(t, "SINGLE_COMMAND", nested(t, value, fieldTypeName).GetString())
	assert.True(t, nested(t, value, fieldValue).GetBool())
	assert.False(t, nested(t, value, fieldCommandQualifier, "select").GetBool())
	assert.Equal(t, uint8(0), nested(t, value, fieldCommandQualifier, "qualifier").GetUint8())
}

// The end-of-initialisation a station sends after a restart, taken from the testsuite.
func TestDecodePoint_EndOfInitialisation(t *testing.T) {
	apdu := parseHexApdu(t, "680e00000000460104000a0000000000")
	iFormat, ok := apdu.(readWriteModel.APDUIFormat)
	require.True(t, ok)
	objects := iFormat.GetAsdu().GetInformationObjects()
	require.Len(t, objects, 1)

	value, _, code := decodePoint(iFormat.GetAsdu(), objects[0])

	require.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.Equal(t, "END_OF_INITIALISATION", nested(t, value, fieldTypeName).GetString())
	assert.Equal(t, uint8(0), nested(t, value, fieldValue).GetUint8())
	assert.False(t, nested(t, value, "afterLocalParameterChange").GetBool())
}

// An interrogation command mirrored back, taken from the testsuite.
func TestDecodePoint_InterrogationCommand(t *testing.T) {
	apdu := parseHexApdu(t, "680e00000000640106000a0000000014")
	iFormat, ok := apdu.(readWriteModel.APDUIFormat)
	require.True(t, ok)
	objects := iFormat.GetAsdu().GetInformationObjects()
	require.Len(t, objects, 1)

	value, _, code := decodePoint(iFormat.GetAsdu(), objects[0])

	require.Equal(t, apiModel.PlcResponseCode_OK, code)
	assert.Equal(t, "INTERROGATION_COMMAND", nested(t, value, fieldTypeName).GetString())
	assert.Equal(t, "ACTIVATION_ACT", nested(t, value, fieldCauseOfTransmission).GetString())
	assert.Equal(t, uint8(20), nested(t, value, fieldValue).GetUint8())
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// What cannot be decoded says so
////////////////////////////////////////////////////////////////////////////////////////////////////

// The file transfer ASDUs parse, but every one of their payload types is an empty type in the mspec,
// so there is nothing at all to read out of them. Reporting them as OK with a null value would make
// a point which said nothing look like a point which reported nothing.
func TestDecodePoint_FileTransferIsUnsupported(t *testing.T) {
	frame := iFormatFrame(0, 0, asduBytes(0x78, 3, 10, informationObjectBytes(13)))
	value, _, code := decodeFirstPoint(t, frame)

	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, code)
	assert.Equal(t, "FILE_READY", nested(t, value, fieldTypeName).GetString())
	assert.Equal(t, apiValues.NULL, nested(t, value, fieldValue).GetPlcValueType())
}

// The mspec for the protection-equipment event ASDUs lists only their two time fields and leaves out
// the SEP octet the event state lives in, so the generated model has nowhere to keep it. What the
// model does carry is still reported, but the point is not reported as OK.
func TestDecodePoint_ProtectionEquipmentEventIsUnsupported(t *testing.T) {
	// M_EP_TD_1: the elapsed time followed by a CP56Time2a.
	object := informationObjectBytes(13, 0xE8, 0x03, 0xC7, 0x5D, 0x17, 0x08, 0x84, 0x07, 0x0D)
	frame := iFormatFrame(0, 0, asduBytes(0x26, 3, 10, object))
	value, _, code := decodeFirstPoint(t, frame)

	assert.Equal(t, apiModel.PlcResponseCode_UNSUPPORTED, code)
	assert.Equal(t, uint16(1000), nested(t, value, fieldElapsedTime).GetUint16())
	assert.Equal(t,
		time.Date(2013, time.July, 4, 8, 23, 24, 7*int(time.Millisecond), time.UTC),
		nested(t, value, fieldTimestamp, "dateTime").GetDateTime())
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// Change detection
////////////////////////////////////////////////////////////////////////////////////////////////////

// The fingerprint is what a change-of-state subscription compares reports by, so it has to be stable
// across repeated reports of an unchanged point - and a Go map iteration order is not.
func TestDecodePoint_FingerprintIsStable(t *testing.T) {
	frame := iFormatFrame(0, 0, asduBytes(0x0B, 3, 10, informationObjectBytes(13, 0xD2, 0x04, 0x00)))
	_, first, _ := decodeFirstPoint(t, frame)
	for range 20 {
		_, again, _ := decodeFirstPoint(t, frame)
		require.Equal(t, first, again)
	}
}

// A value change, a quality change and a type change are all changes. A new timestamp, a new elapsed
// time and a different cause of transmission for the same reading are not - a station re-sends the
// very same value in answer to a general interrogation.
func TestDecodePoint_FingerprintTracksStateOnly(t *testing.T) {
	spontaneous := byte(3)
	interrogated := byte(20)
	object := func(value byte, qds byte) []byte {
		return informationObjectBytes(13, value, 0x00, qds)
	}
	_, base, _ := decodeFirstPoint(t, iFormatFrame(0, 0, asduBytes(0x0B, spontaneous, 10, object(0x2A, 0x00))))

	t.Run("a different value is a change", func(t *testing.T) {
		_, other, _ := decodeFirstPoint(t, iFormatFrame(0, 0, asduBytes(0x0B, spontaneous, 10, object(0x2B, 0x00))))
		assert.NotEqual(t, base, other)
	})
	t.Run("a different quality is a change", func(t *testing.T) {
		_, other, _ := decodeFirstPoint(t, iFormatFrame(0, 0, asduBytes(0x0B, spontaneous, 10, object(0x2A, 0x80))))
		assert.NotEqual(t, base, other)
	})
	t.Run("a different cause of transmission is not", func(t *testing.T) {
		_, other, _ := decodeFirstPoint(t, iFormatFrame(0, 0, asduBytes(0x0B, interrogated, 10, object(0x2A, 0x00))))
		assert.Equal(t, base, other)
	})
	t.Run("a different time tag is not", func(t *testing.T) {
		withTime := func(milliseconds byte) []byte {
			return informationObjectBytes(13, 0x2A, 0x00, 0x00, milliseconds, 0x00, 0x2A, 0x08, 0x84, 0x07, 0x0D)
		}
		_, early, _ := decodeFirstPoint(t, iFormatFrame(0, 0, asduBytes(0x23, spontaneous, 10, withTime(0x01))))
		_, late, _ := decodeFirstPoint(t, iFormatFrame(0, 0, asduBytes(0x23, spontaneous, 10, withTime(0x02))))
		assert.Equal(t, early, late)
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// A whole multi-object ASDU off the wire
////////////////////////////////////////////////////////////////////////////////////////////////////

// One ASDU carries as many information objects as it likes, and each of them is a point of its own.
// The frame is one the parser-serializer testsuite captured: four single point informations at
// ASDU 10, addresses 1 to 4, in answer to a general interrogation.
func TestDecodePoint_EveryObjectOfAnAsduIsItsOwnPoint(t *testing.T) {
	apdu := parseHexApdu(t, "681a04000200010414000a0001000000020000000300000004000000")
	iFormat, ok := apdu.(readWriteModel.APDUIFormat)
	require.True(t, ok)
	objects := iFormat.GetAsdu().GetInformationObjects()
	require.Len(t, objects, 4)

	for index, informationObject := range objects {
		value, _, code := decodePoint(iFormat.GetAsdu(), informationObject)
		require.Equal(t, apiModel.PlcResponseCode_OK, code)
		assert.Equal(t, uint16(10), nested(t, value, fieldAsduAddress).GetUint16())
		assert.Equal(t, uint32(index+1), nested(t, value, fieldObjectAddress).GetUint32())
		assert.Equal(t, "INTERROGATED_BY_GENERAL_INTERROGATION_INROGEN", nested(t, value, fieldCauseOfTransmission).GetString())
		assert.False(t, nested(t, value, fieldValue).GetBool())
	}
	_ = testutils.TestContext(t)
}
