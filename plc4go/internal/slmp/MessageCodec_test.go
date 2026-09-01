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

package slmp

import (
	"context"
	"encoding/hex"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/slmp/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/transports/test"
)

// The 3E frames below are the reference vectors from
// plc4go/assets/testing/protocols/slmp/ParserSerializerTestsuite.xml, which derives them from the
// worked examples in SH-080008.
const (
	// batchReadRequestHex is Batch Read of D350, 2 words.
	batchReadRequestHex = "500000ffff03000c000000010400005e0100a80200"
	// batchReadResponseHex answers it with D350=0x56AB, D351=0x170F.
	batchReadResponseHex = "d00000ffff030006000000ab560f17"
	// batchWriteRequestHex is Batch Write of D350, 2 words (0x1234, 0x5678). The parser-serializer
	// testsuite has no write case, so this vector is assembled from the frame layout in the mspec:
	// subheader 50 00, access route 00 FF FF03 00, requestDataLength 0x0010 (6 + 10), monitoring
	// timer 0x0000, command 0x1401, subcommand 0x0000, then D350 / code D / 2 points / the words.
	batchWriteRequestHex = "500000ffff030010000000011400005e0100a8020034127856"
	// writeResponseHex is a Batch Write success, which carries no payload at all.
	writeResponseHex = "d00000ffff03000200" + "0000"
)

// newTestCodec spins a MessageCodec up on a connected (but not running) test transport. The
// transport instance is what the tests feed frames into and read the codec's output back out of.
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

// TestFrameSizeFromHeader pins the length field down. Reading it big-endian - the way modbus reads
// its MBAP length at a similar offset - would turn a 15 byte response into a 1545 byte one and leave
// the receive loop waiting for bytes that never come.
func TestFrameSizeFromHeader(t *testing.T) {
	tests := []struct {
		name   string
		header []byte
		want   uint32
	}{
		{
			name:   "the length field is little-endian",
			header: mustDecodeHex(t, "d00000ffff03000600"),
			want:   15,
		},
		{
			name:   "a response with no payload at all",
			header: mustDecodeHex(t, "d00000ffff03000200"),
			want:   11,
		},
		{
			// Read big-endian this would be 0x0006 -> 15; the little-endian read is 0x0600 -> 1545.
			name:   "a length whose high byte is set",
			header: mustDecodeHex(t, "d00000ffff03000006"),
			want:   1545,
		},
		{
			// A uint16 total would wrap 0xFFFF + 9 back to 8, making Read consume less than a frame
			// and leaving the receive worker unable to ever resynchronize.
			name:   "the biggest length doesn't wrap the total",
			header: mustDecodeHex(t, "d00000ffff0300ffff"),
			want:   65544,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, frameSizeFromHeader(testCase.header))
		})
	}
}

func TestMessageCodec_SendSerializesTheReferenceFrames(t *testing.T) {
	t.Run("batch read", func(t *testing.T) {
		codec, transportInstance := newTestCodec(t)
		request := readWriteModel.NewSlmpRequestFrame3E(0x0000, commandBatchRead, subCommandWordUnits,
			readWriteModel.NewSlmpReadRequest(350, readWriteModel.SlmpDeviceCode_D, 2))
		require.NoError(t, codec.Send(testutils.TestContext(t), "read", request))
		want := mustDecodeHex(t, batchReadRequestHex)
		assert.Equal(t, want, transportInstance.DrainWriteBuffer(uint32(len(want))))
	})
	t.Run("batch write", func(t *testing.T) {
		codec, transportInstance := newTestCodec(t)
		request := readWriteModel.NewSlmpRequestFrame3E(0x0000, commandBatchWrite, subCommandWordUnits,
			readWriteModel.NewSlmpWriteRequest(350, readWriteModel.SlmpDeviceCode_D, 2,
				[]byte{0x34, 0x12, 0x78, 0x56}))
		require.NoError(t, codec.Send(testutils.TestContext(t), "write", request))
		want := mustDecodeHex(t, batchWriteRequestHex)
		assert.Equal(t, want, transportInstance.DrainWriteBuffer(uint32(len(want))))
	})
}

func TestMessageCodec_SendRejectsForeignMessages(t *testing.T) {
	codec, _ := newTestCodec(t)
	// Request data on its own is not a 3E frame, so it can't be framed.
	err := codec.Send(testutils.TestContext(t), "nonsense",
		readWriteModel.NewSlmpReadRequest(350, readWriteModel.SlmpDeviceCode_D, 2))
	assert.Error(t, err)
}

func TestMessageCodec_ReceiveFramesTheStream(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// A write success (no payload at all) directly followed by a read response with a four byte
	// payload, back to back in one chunk: the codec has to cut them apart on the length field alone.
	transportInstance.FillReadBuffer(mustDecodeHex(t, writeResponseHex+batchReadResponseHex))

	writeResponse, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	typedWriteResponse, ok := writeResponse.(readWriteModel.SlmpResponseFrame3E)
	require.True(t, ok, "%T is not a 3E response frame", writeResponse)
	assert.Equal(t, uint16(0x0000), typedWriteResponse.GetEndCode())
	assert.Empty(t, typedWriteResponse.GetResponseData(), "a write success carries no data")

	readResponse, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	typedReadResponse, ok := readResponse.(readWriteModel.SlmpResponseFrame3E)
	require.True(t, ok, "%T is not a 3E response frame", readResponse)
	assert.Equal(t, uint16(0x0000), typedReadResponse.GetEndCode())
	assert.Equal(t, []byte{0xAB, 0x56, 0x0F, 0x17}, typedReadResponse.GetResponseData())
}

func TestMessageCodec_ReceiveWaitsForTheWholeFrame(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// A partial frame leaves the codec waiting for the rest of it, so the two incomplete steps below
	// get a context of their own that runs out - which is what a real transport going quiet looks
	// like from in here.
	receivePartial := func() any {
		t.Helper()
		ctx, cancel := context.WithTimeout(testutils.TestContext(t), 250*time.Millisecond)
		defer cancel()
		message, err := codec.Receive(ctx)
		require.NoError(t, err)
		return message
	}

	// Only the subheader and part of the access route: the length field isn't complete, so there is
	// nothing to decide a frame size from.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "d00000ffff"))
	assert.Nil(t, receivePartial(), "a frame whose length isn't known yet must not be framed")

	// The length is decidable now, but the payload is still missing.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "0300060000"))
	assert.Nil(t, receivePartial(), "an incomplete frame must not be framed")

	// And now the rest of it.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "00ab560f17"))
	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message)
	response, ok := message.(readWriteModel.SlmpResponseFrame3E)
	require.True(t, ok, "%T is not a 3E response frame", message)
	assert.Equal(t, []byte{0xAB, 0x56, 0x0F, 0x17}, response.GetResponseData())
}

func TestMessageCodec_ReceiveDropsAMalformedFrame(t *testing.T) {
	codec, transportInstance := newTestCodec(t)

	// 0x99 is not a subheader the type switch knows, so the frame can't be parsed. It is consumed
	// all the same, which is what keeps the stream in sync for the frame after it.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "990000ffff030006000000ab560f17"))
	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	assert.Nil(t, message)

	transportInstance.FillReadBuffer(mustDecodeHex(t, batchReadResponseHex))
	message, err = codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	require.NotNil(t, message, "the codec has to resynchronize on the frame after a bad one")
	_, ok := message.(readWriteModel.SlmpResponseFrame3E)
	assert.True(t, ok, "%T is not a 3E response frame", message)
}

// TestMessageCodec_ReceiveParsesAnErrorFrame is what an abnormal completion looks like on the wire:
// the end code is non-zero and the payload carries error information rather than device data.
func TestMessageCodec_ReceiveParsesAnErrorFrame(t *testing.T) {
	codec, transportInstance := newTestCodec(t)
	// endCode 0xC059 (a "command cannot be executed" style error), with the 11 byte error
	// information block the manual appends to an abnormal 3E response.
	transportInstance.FillReadBuffer(mustDecodeHex(t, "d00000ffff03000d0059c00000ffff03005e0100a801"))

	message, err := codec.Receive(testutils.TestContext(t))
	require.NoError(t, err)
	response, ok := message.(readWriteModel.SlmpResponseFrame3E)
	require.True(t, ok, "%T is not a 3E response frame", message)
	assert.Equal(t, uint16(0xC059), response.GetEndCode())
	assert.Len(t, response.GetResponseData(), 11)
}
