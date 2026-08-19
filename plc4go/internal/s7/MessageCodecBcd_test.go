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
	"bytes"
	"net/url"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/s7/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// newTestCodec spins a MessageCodec up on a connected (but not running) test
// transport. The transport instance is what the test feeds raw frames into.
//
// Deliberately NOT built on testutils.EnrichOptionsWithOptionsForTesting: that
// helper calls shouldNoColor, which flips the process global color.NoColor to
// false on a developer machine (no CI/NO_COLOR env), and the box renderer then
// emits ANSI escapes that break the golden strings in TestS7MessageBytes. A
// silent logger is all this test needs.
func newTestCodec(t *testing.T) (*MessageCodec, *test.TransportInstance) {
	t.Helper()
	_options := []options.WithOption{options.WithCustomLogger(zerolog.Nop())}

	transport := test.NewTransport(_options...)
	transportInstance, err := transport.CreateTransportInstance(url.URL{Scheme: "test"}, nil, _options...)
	require.NoError(t, err)
	testTransportInstance, ok := transportInstance.(*test.TransportInstance)
	require.True(t, ok)
	require.NoError(t, testTransportInstance.Connect(t.Context()))
	t.Cleanup(func() {
		assert.NoError(t, testTransportInstance.Close())
	})

	return NewMessageCodec(testTransportInstance, _options...), testTransportInstance
}

// alarmPushPacket wraps an ALARM_S push carrying timeStamp into a complete TPKT
// frame, the way the CPU puts it on the wire.
func alarmPushPacket(t *testing.T, timeStamp readWriteModel.DateAndTime) []byte {
	t.Helper()
	message := readWriteModel.NewS7MessageUserData(
		0,
		readWriteModel.NewS7ParameterUserData([]readWriteModel.S7ParameterUserDataItem{
			// The alarm push is cpuFunctionType 0x00 / group 0x04, for which
			// s7.mspec makes dataUnitReferenceNumber, lastDataUnit and
			// errorCode absent - supplying them anyway would serialize bytes
			// the parser does not read back and shift the whole payload.
			readWriteModel.NewS7ParameterUserDataItemCPUFunctions(
				0x12, 0x00, 0x04, 0x12, 0,
				nil, nil, nil,
			),
		}),
		readWriteModel.NewS7PayloadUserData([]readWriteModel.S7PayloadUserDataItem{
			readWriteModel.NewS7PayloadAlarmS(
				// dataLength covers the AlarmMessagePushType: 8 bytes timestamp
				// + functionId + numberOfObjects, with no message objects.
				readWriteModel.DataTransportErrorCode_OK, readWriteModel.DataTransportSize_OCTET_STRING, 10,
				readWriteModel.NewAlarmMessagePushType(timeStamp, 0x11, 0, nil),
			),
		}),
	)
	data, err := readWriteModel.NewTPKTPacket(
		readWriteModel.NewCOTPPacketData(nil, message, true, 1),
	).Serialize()
	require.NoError(t, err)
	return data
}

// TestMessageCodec_ReceiveDropsAlarmWithInvalidBCDTimestamp is the driver level
// counterpart to TestAlarm_UnsetTimestampFailsTheWholeMessage in
// S7BcdCoverage_test.go, which only covers the model.
//
// Honouring encoding='"BCD"' turned a malformed timestamp from a garbage value
// into a hard parse failure, and because AlarmMessagePushType embeds DateAndTime
// as its first field, that failure now takes the WHOLE alarm message with it.
// Some CPUs signal an unset timestamp as eight 0xFF bytes, so this is reachable
// on live traffic rather than theoretical.
//
// The chosen behaviour is: reject the frame, keep the connection. Rejecting is
// the right codec behaviour (plc4j's EncodingBCD throws on the same input) and
// keeping the connection is what MessageCodec.Receive already does for any parse
// error - it logs at WARN and returns (nil, nil), so the receive worker drops the
// packet and carries on. This test pins that combination end to end, so a future
// change that turns the parse error into a fatal transport error (killing the
// connection over one bad alarm) or back into silent garbage fails here.
func TestMessageCodec_ReceiveDropsAlarmWithInvalidBCDTimestamp(t *testing.T) {
	validTimeStamp := readWriteModel.NewDateAndTime(24, 8, 19, 10, 30, 45, 123, 2)
	valid := alarmPushPacket(t, validTimeStamp)

	// The eight timestamp bytes are the first thing in the AlarmMessagePushType,
	// so locate them in the serialized frame rather than hard coding an offset.
	validTimeStampBytes, err := validTimeStamp.Serialize()
	require.NoError(t, err)
	timeStampOffset := bytes.Index(valid, validTimeStampBytes)
	require.GreaterOrEqual(t, timeStampOffset, 0, "the timestamp must be findable in the frame")

	corrupted := bytes.Clone(valid)
	for i := range validTimeStampBytes {
		corrupted[timeStampOffset+i] = 0xFF // 0xF is not a valid BCD digit
	}

	codec, transportInstance := newTestCodec(t)

	// The corrupted frame is dropped: no message, no error, connection intact.
	transportInstance.FillReadBuffer(corrupted)
	message, err := codec.Receive(t.Context())
	assert.NoError(t, err, "one unparseable alarm must not fail the codec")
	assert.Nil(t, message, "the alarm must not be surfaced with a garbage timestamp")
	assert.True(t, transportInstance.IsConnected(), "the connection must survive")

	// And the very next well formed frame still parses on the same codec.
	transportInstance.FillReadBuffer(valid)
	message, err = codec.Receive(t.Context())
	require.NoError(t, err)
	require.NotNil(t, message, "the codec must still be usable after dropping a frame")
	parsed := message.(readWriteModel.TPKTPacket).
		GetPayload().(readWriteModel.COTPPacketData).
		GetPayload().(readWriteModel.S7MessageUserData)
	alarm := parsed.GetPayload().(readWriteModel.S7PayloadUserData).
		GetItems()[0].(readWriteModel.S7PayloadAlarmS).
		GetAlarmMessage().(readWriteModel.AlarmMessagePushType)
	assert.Equal(t, validTimeStamp, alarm.GetTimeStamp())
}
