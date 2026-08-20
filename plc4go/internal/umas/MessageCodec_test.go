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

package umas

import (
	"context"
	"encoding/hex"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/umas/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// newTestCodec spins a MessageCodec up on a connected (but not running) test transport. The transport
// instance is what the tests feed packets into and read the codec's output back out of.
func newTestCodec(t *testing.T) (*MessageCodec, *test.TransportInstance) {
	t.Helper()
	_options := testutils.EnrichOptionsWithOptionsForTesting(t)

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

func mustDecodeHex(t *testing.T, hexString string) []byte {
	t.Helper()
	decoded, err := hex.DecodeString(hexString)
	require.NoError(t, err)
	return decoded
}

func mustEncodeHex(theBytes []byte) string {
	return hex.EncodeToString(theBytes)
}

// drainAll hands back everything the codec has written so far, hex encoded.
func drainAll(transportInstance *test.TransportInstance) string {
	return mustEncodeHex(transportInstance.DrainWriteBuffer(transportInstance.GetNumDrainableBytes()))
}

func TestPacketSizeFromHeader(t *testing.T) {
	tests := []struct {
		name   string
		header []byte
		want   uint32
	}{
		// The length field counts the unit identifier plus the PDU, so a frame is length + 6 bytes.
		{name: "the smallest frame there can be", header: mustDecodeHex(t, "000100000002"), want: 8},
		{name: "the length field is big endian", header: mustDecodeHex(t, "000100000100"), want: 262},
		// A uint16 total would wrap 0xFFFF + 6 back to 5, making Read consume less than a frame and
		// leaving the receive worker unable to ever resynchronize.
		{name: "the biggest length doesn't wrap the total", header: mustDecodeHex(t, "00010000FFFF"), want: 65541},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, packetSizeFromHeader(testCase.header))
		})
	}
}

func TestHeaderLooksLikeMbap(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   bool
	}{
		{name: "a plain header", header: "000100000004", want: true},
		{name: "a non zero protocol identifier is not Modbus/TCP", header: "000100010004"},
		{name: "a length below the unit identifier plus a function code", header: "000100000001"},
		{name: "a length of zero", header: "000100000000"},
		{name: "an incomplete header can't be judged", header: "0001000000"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, headerLooksLikeMbap(mustDecodeHex(t, testCase.header)))
		})
	}
}

// The wire bytes of the first request of the handshake, spelled out: MBAP header (transaction
// identifier 0001 big endian, protocol identifier 0000, length 0004, unit identifier 00) followed by
// the Modbus function code 5A, the UMAS pairing key 00 and the UMAS function key 02 (PlcIdent).
func TestMessageCodec_SendFramesARequest(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	request := readWriteModel.NewModbusTcpADU(1, 0,
		readWriteModel.NewUmasPDU(readWriteModel.NewUmasPDUPlcIdentRequest(0)))
	require.NoError(t, codec.Send(testutils.TestContext(t), "PlcIdent", request))

	assert.Equal(t, "000100000004"+"00"+"5a"+"00"+"02", drainAll(transportInstance))
}

// The unit identifier is the last byte of the MBAP header and a UMAS PDU is little endian inside a
// big endian header, so a request against unit 3 differs from the one above in exactly one byte.
func TestMessageCodec_SendCarriesTheUnitIdentifier(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	request := readWriteModel.NewModbusTcpADU(0x1234, 3,
		readWriteModel.NewUmasPDU(readWriteModel.NewUmasPDUPlcIdentRequest(0)))
	require.NoError(t, codec.Send(testutils.TestContext(t), "PlcIdent", request))

	assert.Equal(t, "123400000004"+"03"+"5a"+"00"+"02", drainAll(transportInstance))
}

func TestMessageCodec_SendRejectsForeignMessages(t *testing.T) {
	codec, _ := newTestCodec(t)
	// A bare UMAS PDU isn't an ADU, so it can't be framed.
	err := codec.Send(testutils.TestContext(t), "nonsense",
		readWriteModel.NewUmasPDU(readWriteModel.NewUmasPDUPlcIdentRequest(0)))
	assert.Error(t, err)
}

// This is the whole reason the codec exists rather than a plain parse: a UMAS response carries the
// generic function key FE and only the function key of the *request* says what kind of response it is.
// Send records it, Receive looks it up by transaction identifier.
func TestMessageCodec_ResolvesResponsesByTheRequestFunctionKey(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// A read-variable request (UMAS function key 22) for transaction 0007.
	readRequest := readWriteModel.NewModbusTcpADU(7, 0,
		readWriteModel.NewUmasPDU(readWriteModel.NewUmasPDUReadVariableRequest(0, 0, 1,
			[]readWriteModel.VariableReadRequestReference{
				readWriteModel.NewVariableReadRequestReference(0, 3, 0, 0, 0, nil),
			})))
	require.NoError(t, codec.Send(testutils.TestContext(t), "ReadVariable", readRequest))
	// The wire bytes of a read: header, unit 00, function code 5A, pairing key 00, UMAS function key
	// 22, the project CRC (all zeroes here), the variable count 01 and one reference whose isArray
	// and dataSizeIndex nibbles are 0 and 3.
	assert.Equal(t, "000700000010"+"00"+"5a"+"00"+"22"+"00000000"+"01"+"03"+"0000"+"01"+"0000"+"00",
		drainAll(transportInstance))

	// The answer: transaction 0007, length 0008, unit 00, function code 5A, pairing key 00, the
	// generic response function key FE, then four payload bytes.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "00070000000800"+"5a00fe"+"11223344"))

	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message)
	response, ok := message.(readWriteModel.ModbusTcpADU)
	require.True(t, ok, "%T is not an ADU", message)
	assert.Equal(t, uint16(7), response.GetTransactionIdentifier())
	umasPdu, ok := response.GetPdu().(readWriteModel.UmasPDU)
	require.True(t, ok, "%T is not a UMAS PDU", response.GetPdu())
	readResponse, ok := umasPdu.GetItem().(readWriteModel.UmasPDUReadVariableResponse)
	require.True(t, ok, "%T is not a read-variable response - the request function key wasn't applied", umasPdu.GetItem())
	assert.Equal(t, mustDecodeHex(t, "11223344"), readResponse.GetBlock())
}

// Without a tracked request the very same bytes can only be read as the generic success response.
// That is what an unsolicited packet, or one whose request has already timed out, looks like.
func TestMessageCodec_UntrackedResponsesFallBackToTheGenericType(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	transportInstance.FillReadBuffer(mustDecodeHex(t, "00070000000800"+"5a00fe"+"11223344"))

	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message)
	umasPdu, ok := message.(readWriteModel.ModbusTcpADU).GetPdu().(readWriteModel.UmasPDU)
	require.True(t, ok)
	_, ok = umasPdu.GetItem().(readWriteModel.UmasPDUSuccessResponse)
	assert.True(t, ok, "%T should have been the generic success response", umasPdu.GetItem())
}

// A Modbus exception sets the top bit of the function code, which the model discriminates on without
// needing the request's function key at all.
func TestMessageCodec_ReceiveParsesAModbusException(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// Transaction 0001, length 0003, unit 00, function code DA (5A with the error flag), exception
	// code 02 (illegal data address).
	transportInstance.FillReadBuffer(mustDecodeHex(t, "00010000000300"+"da02"))

	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message)
	errorPdu, ok := message.(readWriteModel.ModbusTcpADU).GetPdu().(readWriteModel.ModbusPDUError)
	require.True(t, ok, "%T is not a Modbus error", message.(readWriteModel.ModbusTcpADU).GetPdu())
	assert.Equal(t, readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS, errorPdu.GetExceptionCode())
}

// TCP hands over whatever has arrived, which is regularly less or more than one frame.
func TestMessageCodec_ReceiveFramesTheStream(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// Two status responses back to back in one chunk.
	transportInstance.FillReadBuffer(mustDecodeHex(t,
		"000100000004005a00fe"+
			"000200000005005a00fe07"))

	first, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, first)
	assert.Equal(t, uint16(1), first.(readWriteModel.ModbusTcpADU).GetTransactionIdentifier())

	second, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, second)
	assert.Equal(t, uint16(2), second.(readWriteModel.ModbusTcpADU).GetTransactionIdentifier())
}

func TestMessageCodec_ReceiveWaitsForAWholeFrame(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// A partial frame leaves the codec waiting for the rest of it, so the incomplete steps below get
	// a context of their own that runs out - which is what a real transport going quiet looks like
	// from in here.
	receivePartial := func() spi.Message {
		t.Helper()
		ctx, cancel := context.WithTimeout(testutils.TestContext(t), 250*time.Millisecond)
		defer cancel()
		message, err := codec.Receive(ctx)
		require.NoError(t, err)
		return message
	}

	// Only the first two bytes: not even the length field is there, so there is nothing to decide a
	// frame size from.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "0001"))
	assert.Nil(t, receivePartial(), "a frame whose length isn't known yet must not be framed")

	// The length is decidable now - it says ten bytes - but only eight have arrived.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "0000000400"+"5a"))
	assert.Nil(t, receivePartial(), "an incomplete frame must not be framed")

	// And now the rest of it.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "00fe"))
	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, uint16(1), message.(readWriteModel.ModbusTcpADU).GetTransactionIdentifier())
}

// A stream which is out of step has no start marker to search for, so the codec drops one byte at a
// time until a header appears. Without the check on the protocol identifier it would instead read the
// garbage as a length and wait for a frame nobody is going to send.
func TestMessageCodec_ReceiveResynchronizesOnGarbage(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	transportInstance.FillReadBuffer(append(
		mustDecodeHex(t, "ffffffffffff"),
		mustDecodeHex(t, "000100000004005a00fe")...))

	// One Receive per discarded byte, then the frame.
	var message spi.Message
	for attempt := 0; attempt < 12 && message == nil; attempt++ {
		var err error
		message, err = codec.Receive(testutils.TestContext(t))
		require.NoError(t, err)
	}
	require.NotNil(t, message, "the codec never resynchronized")
	assert.Equal(t, uint16(1), message.(readWriteModel.ModbusTcpADU).GetTransactionIdentifier())
}

// A frame whose header is fine but whose body doesn't parse is consumed anyway, so one bad packet
// costs exactly that packet.
func TestMessageCodec_ReceiveDiscardsAnUnparsableFrame(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// Function code 01 is not 5A, so there is no PDU the UMAS model can make of it.
	transportInstance.FillReadBuffer(append(
		mustDecodeHex(t, "00010000000400"+"010203"),
		mustDecodeHex(t, "000200000004005a00fe")...))

	message, err := codec.Receive(testutils.TestContext(t))
	assert.NoError(t, err)
	assert.Nil(t, message, "the unparsable frame is not a message")

	message, err = codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message, "the frame after the bad one has to still be readable")
	assert.Equal(t, uint16(2), message.(readWriteModel.ModbusTcpADU).GetTransactionIdentifier())
}
