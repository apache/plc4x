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

package modbus

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// A caller whose context expired abandons the result channel without reading.
// The send-failure result then fills the single-slot buffer; when the still
// registered expectation later times out, its error handler must not block
// forever on the full channel — those blocked handlers pile up in the codec's
// WaitGroup and wedge Disconnect indefinitely.
func TestReader_lateTimeoutAfterFailedSendMustNotBlock(t *testing.T) {
	codec := newCaptureCodec(errors.New("send failed: broken pipe"))
	reader := NewReader(DefaultConfiguration(), codec, testTransactionManager())
	tag := NewTag(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_UINT)
	request := spiModel.NewDefaultPlcReadRequest(
		map[string]apiModel.PlcTag{"tag": tag}, []string{"tag"}, reader, nil)

	results := reader.Read(testutils.TestContext(t), request)

	var handlers capturedHandlers
	select {
	case handlers = <-codec.handlers:
	case <-time.After(time.Second):
		t.Fatal("SendRequest was never invoked")
	}

	// Wait for the send-failure result to occupy the channel buffer; the
	// abandoned caller never drains it.
	require.Eventually(t, func() bool { return len(results) == 1 },
		time.Second, time.Millisecond)

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = handlers.handleError(errors.New("timeout"))
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("late timeout handler blocked on the abandoned result channel")
	}
}

// A message that isn't a modbus ADU must never reach an unchecked type assertion; it is simply
// not ours, and if one is handed to us anyway it becomes an error result.
func TestReader_nonAduMessageIsRejected(t *testing.T) {
	codec := newCaptureCodec(nil)
	reader := NewReader(DefaultConfiguration(), codec, testTransactionManager())
	tag := NewTag(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_UINT)
	request := spiModel.NewDefaultPlcReadRequest(
		map[string]apiModel.PlcTag{"tag": tag}, []string{"tag"}, reader, nil)

	results := reader.Read(testutils.TestContext(t), request)

	var handlers capturedHandlers
	select {
	case handlers = <-codec.handlers:
	case <-time.After(time.Second):
		t.Fatal("SendRequest was never invoked")
	}

	require.False(t, handlers.acceptsMessage(notAnAdu{}), "a foreign message must not be accepted")
	require.NoError(t, handlers.handleMessage(notAnAdu{}))

	select {
	case result := <-results:
		require.Error(t, result.GetErr())
	case <-time.After(2 * time.Second):
		t.Fatal("no result delivered for the foreign message")
	}
}

func readRequestFor(t *testing.T, reader *Reader, tag apiModel.PlcTag) apiModel.PlcReadRequest {
	t.Helper()
	return spiModel.NewDefaultPlcReadRequest(
		map[string]apiModel.PlcTag{"tag": tag}, []string{"tag"}, reader, nil)
}

// capturedReadRequest runs a read and hands back the ADU that went onto the wire.
func capturedReadRequest(t *testing.T, configuration Configuration, tag apiModel.PlcTag) readWriteModel.ModbusTcpADU {
	t.Helper()
	codec := newCaptureCodec(nil)
	reader := NewReader(configuration, codec, testTransactionManager())
	_ = reader.Read(testutils.TestContext(t), readRequestFor(t, reader, tag))

	select {
	case handlers := <-codec.handlers:
		adu, ok := handlers.message.(readWriteModel.ModbusTcpADU)
		require.True(t, ok, "expected a ModbusTcpADU, got %T", handlers.message)
		return adu
	case <-time.After(time.Second):
		t.Fatal("SendRequest was never invoked")
		return nil
	}
}

// A tag may name the unit it is addressed at; without one the connection's default is used
// (plc4j ModbusTcpConnection.getUnitId).
func TestReader_resolvesTheUnitIdentifier(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.unitIdentifier = 17

	t.Run("the connection default", func(t *testing.T) {
		adu := capturedReadRequest(t, configuration, parseTag(t, "holding-register:1:INT"))
		assert.Equal(t, uint8(17), adu.GetUnitIdentifier())
	})
	t.Run("the one the tag names", func(t *testing.T) {
		adu := capturedReadRequest(t, configuration, parseTag(t, "holding-register:1:INT{unit-id: 5}"))
		assert.Equal(t, uint8(5), adu.GetUnitIdentifier())
	})
}

// A STRING(20) is ten registers long. With a literal string length of 1 in its place the request
// would ask for a single register and the value would be unreadable.
func TestReader_asksForTheRegistersAStringOccupies(t *testing.T) {
	adu := capturedReadRequest(t, DefaultConfiguration(), parseTag(t, "holding-register:1:STRING(20)"))

	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUReadHoldingRegistersRequest)
	require.True(t, ok, "expected a read-holding-registers request, got %T", adu.GetPdu())
	assert.Equal(t, uint16(10), pdu.GetQuantity())
}

// A string is decoded with the length its address declared.
func TestReader_decodesAStringOfTheDeclaredLength(t *testing.T) {
	reader := NewReader(DefaultConfiguration(), newCaptureCodec(nil), testTransactionManager())
	tag := parseTag(t, "holding-register:1:STRING(6)")
	request := readRequestFor(t, reader, tag)
	responseAdu := readWriteModel.NewModbusTcpADU(1, 1, readWriteModel.NewModbusPDUReadHoldingRegistersResponse(
		[]byte{'h', 'e', 'l', 'l', 'o', 0x00}))

	response, err := reader.ToPlc4xReadResponse(responseAdu, request)

	require.NoError(t, err)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("tag"))
	assert.Equal(t, "hello", response.GetValue("tag").GetString())
}

// The registers are decoded in the byte order the connection was configured with, or the one the
// tag names for itself (plc4j ModbusTcpConnection.getEffectiveByteOrder).
func TestReader_decodesInTheResolvedByteOrder(t *testing.T) {
	// 0x01020304 as the four byte orders lay it out.
	data := map[ByteOrder][]byte{
		BigEndianOrder:            {0x01, 0x02, 0x03, 0x04},
		LittleEndianOrder:         {0x04, 0x03, 0x02, 0x01},
		BigEndianByteSwapOrder:    {0x02, 0x01, 0x04, 0x03},
		LittleEndianByteSwapOrder: {0x03, 0x04, 0x01, 0x02},
	}
	decode := func(t *testing.T, configuration Configuration, address string, responseData []byte) uint32 {
		t.Helper()
		reader := NewReader(configuration, newCaptureCodec(nil), testTransactionManager())
		request := readRequestFor(t, reader, parseTag(t, address))
		responseAdu := readWriteModel.NewModbusTcpADU(1, 1,
			readWriteModel.NewModbusPDUReadHoldingRegistersResponse(responseData))
		response, err := reader.ToPlc4xReadResponse(responseAdu, request)
		require.NoError(t, err)
		return response.GetValue("tag").GetUint32()
	}

	for byteOrder, responseData := range data {
		t.Run("connection default "+byteOrder.String(), func(t *testing.T) {
			configuration := DefaultConfiguration()
			configuration.defaultPayloadByteOrder = byteOrder
			assert.Equal(t, uint32(0x01020304), decode(t, configuration, "holding-register:1:UDINT", responseData))
		})
		t.Run("tag override "+byteOrder.String(), func(t *testing.T) {
			// The connection says one thing, the tag another - the tag wins.
			configuration := DefaultConfiguration()
			configuration.defaultPayloadByteOrder = LittleEndianOrder
			address := "holding-register:1:UDINT{byte-order: '" + byteOrder.String() + "'}"
			assert.Equal(t, uint32(0x01020304), decode(t, configuration, address, responseData))
		})
	}
}

// addressedTag is one entry of a read request, kept in a slice so that the order the request names
// its tags in - which is the order the optimizer sees them in - is part of the test.
type addressedTag struct {
	name    string
	address string
}

// runRead performs a read of the given addresses and answers every block that goes out with the
// next of the given PDUs. It returns the result the caller ends up with and the requests that went
// onto the wire, so that a test can check both what was asked and what came back.
//
// One response per block, in order: the reader sends the blocks one at a time, so the n-th
// response answers the n-th request.
func runRead(t *testing.T, configuration Configuration, requested []addressedTag, responses ...readWriteModel.ModbusPDU) (apiModel.PlcReadRequestResult, []readWriteModel.ModbusTcpADU) {
	t.Helper()
	codec := newCaptureCodec(nil)
	reader := NewReader(configuration, codec, testTransactionManager())
	tags := make(map[string]apiModel.PlcTag, len(requested))
	tagNames := make([]string, 0, len(requested))
	for _, entry := range requested {
		tags[entry.name] = parseTag(t, entry.address)
		tagNames = append(tagNames, entry.name)
	}
	results := reader.Read(testutils.TestContext(t), spiModel.NewDefaultPlcReadRequest(tags, tagNames, reader, nil))

	sent := make([]readWriteModel.ModbusTcpADU, 0, len(responses))
	for _, response := range responses {
		select {
		case handlers := <-codec.handlers:
			requestAdu, ok := handlers.message.(readWriteModel.ModbusTcpADU)
			require.True(t, ok, "expected a ModbusTcpADU, got %T", handlers.message)
			sent = append(sent, requestAdu)
			responseAdu := readWriteModel.NewModbusTcpADU(requestAdu.GetTransactionIdentifier(), requestAdu.GetUnitIdentifier(), response)
			require.True(t, handlers.acceptsMessage(responseAdu), "the reader didn't accept the response to its own request")
			require.NoError(t, handlers.handleMessage(responseAdu))
		case <-time.After(time.Second):
			t.Fatalf("only %d of the %d expected block requests were sent", len(sent), len(responses))
		}
	}

	select {
	case result := <-results:
		return result, sent
	case <-time.After(2 * time.Second):
		t.Fatal("no result was delivered")
		return nil, nil
	}
}

func holdingRegisterRequestOf(t *testing.T, adu readWriteModel.ModbusTcpADU) readWriteModel.ModbusPDUReadHoldingRegistersRequest {
	t.Helper()
	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUReadHoldingRegistersRequest)
	require.True(t, ok, "expected a read-holding-registers request, got %T", adu.GetPdu())
	return pdu
}

// pi is 3.1415927 as a REAL occupies two registers, the value the driver testsuite reads.
var pi = []byte{0x40, 0x49, 0x0f, 0xdb}

// Tags that one request can cover are read with one request, and the single response answers all
// of them. This is what the optimizer is for: without it every tag costs a round trip, which is
// ruinous behind a TCP<->RTU gateway.
func TestReader_mergesTagsIntoOneRequest(t *testing.T) {
	result, sent := runRead(t, DefaultConfiguration(),
		[]addressedTag{{"first", "holding-register:1:REAL"}, {"second", "holding-register:3:REAL"}},
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse(append(append([]byte{}, pi...), pi...)))

	require.NoError(t, result.GetErr())
	require.Len(t, sent, 1, "both tags fit into one request")
	pdu := holdingRegisterRequestOf(t, sent[0])
	assert.Equal(t, uint16(0), pdu.GetStartingAddress())
	assert.Equal(t, uint16(4), pdu.GetQuantity())

	response := result.GetResponse()
	require.NotNil(t, response)
	for _, tagName := range []string{"first", "second"} {
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode(tagName), tagName)
		assert.InDelta(t, 3.1415927, response.GetValue(tagName).GetFloat32(), 0.0000001, tagName)
	}
}

// Tags nothing can cover with one request are read with one request each, in the order the
// optimizer emits the blocks.
func TestReader_asksOncePerBlock(t *testing.T) {
	// Two registers apart is inside the 125 register window; 500 registers apart is not.
	result, sent := runRead(t, DefaultConfiguration(),
		[]addressedTag{{"near", "holding-register:1:REAL"}, {"far", "holding-register:501:REAL"}},
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse(pi),
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse(pi))

	require.NoError(t, result.GetErr())
	require.Len(t, sent, 2)
	assert.Equal(t, uint16(0), holdingRegisterRequestOf(t, sent[0]).GetStartingAddress())
	assert.Equal(t, uint16(2), holdingRegisterRequestOf(t, sent[0]).GetQuantity())
	assert.Equal(t, uint16(500), holdingRegisterRequestOf(t, sent[1]).GetStartingAddress())
	assert.Equal(t, uint16(2), holdingRegisterRequestOf(t, sent[1]).GetQuantity())

	response := result.GetResponse()
	require.NotNil(t, response)
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("near"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("far"))
}

// Merging reads means one rejected address takes its whole block down - that is the price of the
// optimization. What it must not do is take the rest of the request down with it.
func TestReader_aFailedBlockFailsOnlyItsOwnTags(t *testing.T) {
	result, sent := runRead(t, DefaultConfiguration(),
		[]addressedTag{
			{"merged1", "holding-register:1:REAL"},
			{"merged2", "holding-register:3:REAL"},
			{"separate", "holding-register:501:REAL"},
		},
		readWriteModel.NewModbusPDUError(readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS),
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse(pi))

	require.Len(t, sent, 2)
	// The failure is reported as well as being visible per tag, the way the single-item
	// interceptor reports a partial result.
	require.Error(t, result.GetErr())

	response := result.GetResponse()
	require.NotNil(t, response)
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, response.GetResponseCode("merged1"))
	assert.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, response.GetResponseCode("merged2"))
	assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("separate"))
	assert.InDelta(t, 3.1415927, response.GetValue("separate").GetFloat32(), 0.0000001)
}

// A tag the device never answers takes its block down with a timeout, and the rest of the request
// is still delivered.
func TestReader_aBlockThatIsNeverAnsweredTimesOut(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.requestTimeout = 50 * time.Millisecond
	codec := newCaptureCodec(nil)
	reader := NewReader(configuration, codec, testTransactionManager())
	tag := parseTag(t, "holding-register:1:REAL")
	request := spiModel.NewDefaultPlcReadRequest(
		map[string]apiModel.PlcTag{"tag": tag}, []string{"tag"}, reader, nil)

	results := reader.Read(testutils.TestContext(t), request)
	select {
	case <-codec.handlers:
	case <-time.After(time.Second):
		t.Fatal("SendRequest was never invoked")
	}

	select {
	case result := <-results:
		require.Error(t, result.GetErr())
		require.NotNil(t, result.GetResponse())
		assert.Equal(t, apiModel.PlcResponseCode_REQUEST_TIMEOUT, result.GetResponse().GetResponseCode("tag"))
	case <-time.After(2 * time.Second):
		t.Fatal("the unanswered block never released the request")
	}
}

// Merging is bounded by max-registers-per-request, so a device that can't answer a full width
// request is asked for less rather than for something it rejects.
func TestReader_honoursTheConfiguredCeilings(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.maxRegistersPerRequest = 4

	result, sent := runRead(t, configuration,
		[]addressedTag{
			{"first", "holding-register:1:REAL"},
			{"second", "holding-register:3:REAL"},
			{"third", "holding-register:5:REAL"},
		},
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse(append(append([]byte{}, pi...), pi...)),
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse(pi))

	require.NoError(t, result.GetErr())
	require.Len(t, sent, 2, "four registers is all one request may carry here")
	assert.Equal(t, uint16(4), holdingRegisterRequestOf(t, sent[0]).GetQuantity())
	assert.Equal(t, uint16(2), holdingRegisterRequestOf(t, sent[1]).GetQuantity())
}

// A one-tag read asks for exactly what it asked for before anything was merged: the block the
// optimizer builds around a lone tag is that tag's own read.
func TestReader_aSingleTagRequestIsUnchanged(t *testing.T) {
	// GOLDEN BYTES, captured from the implementation as it stood BEFORE the optimizer
	// was wired in (commit 1033e8b897, via a detached worktree). They are written out
	// literally on purpose.
	//
	// The obvious way to write this test - build the expected PDU with readRequestPdu
	// and compare - is worthless, because wiring the optimizer in rewrote that function
	// to share the block-building code. Both sides of the comparison then move together,
	// so the test passes however wrong the shared code becomes. That version of this test
	// was checked by mutation: making tagSpan return one register too many for every tag
	// left all seven cases green.
	for _, tc := range []struct {
		address string
		pdu     string
	}{
		{"coil:1:BOOL", "0100000001"},
		{"coil:1[0..4]:BOOL", "0100000005"},
		{"discrete-input:1:BOOL", "0200000001"},
		{"holding-register:1:REAL", "0300000002"},
		{"holding-register:1:STRING(20)", "030000000a"},
		{"input-register:7:UDINT", "0400060002"},
		{"extended-register:9998:DINT", "1407060001270e0002"},
	} {
		t.Run(tc.address, func(t *testing.T) {
			adu := capturedReadRequest(t, DefaultConfiguration(), parseTag(t, tc.address))

			wb := utils.NewWriteBufferByteBased()
			require.NoError(t, adu.GetPdu().SerializeWithWriteBuffer(t.Context(), wb))

			assert.Equal(t, tc.pdu, hex.EncodeToString(wb.GetBytes()),
				"a single-tag read must put exactly the bytes on the wire it did before the optimizer")
		})
	}
}

// Coils are packed least significant bit first: the first coil of the response sits in bit 0 of
// the first byte, and a run of them crosses into the next byte at bit 8 rather than restarting
// there.
//
// This is a behaviour CHANGE. The unoptimized path decoded coils through ParseRegisters, which
// reads bits most significant first, so it returned the block reversed - and for a single coil it
// skipped 15 bits into a one byte response and failed outright. Both are fixed by reading the
// response the way the specification packs it.
func TestReader_readsCoilsLeastSignificantBitFirst(t *testing.T) {
	t.Run("a single coil", func(t *testing.T) {
		result, sent := runRead(t, DefaultConfiguration(),
			[]addressedTag{{"tag", "coil:1:BOOL"}},
			readWriteModel.NewModbusPDUReadCoilsResponse([]byte{0x01}))

		require.NoError(t, result.GetErr())
		require.Len(t, sent, 1)
		response := result.GetResponse()
		require.NotNil(t, response)
		assert.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("tag"))
		assert.True(t, response.GetValue("tag").GetBool())
	})

	t.Run("a run of coils", func(t *testing.T) {
		// 0x05 is coils 1, 2 and 3 reading on, off, on.
		result, _ := runRead(t, DefaultConfiguration(),
			[]addressedTag{{"tag", "coil:1[0..2]:BOOL"}},
			readWriteModel.NewModbusPDUReadCoilsResponse([]byte{0x05}))

		require.NoError(t, result.GetErr())
		response := result.GetResponse()
		require.NotNil(t, response)
		require.Equal(t, apiModel.PlcResponseCode_OK, response.GetResponseCode("tag"))
		read := response.GetValue("tag").GetList()
		require.Len(t, read, 3)
		assert.Equal(t, []bool{true, false, true},
			[]bool{read[0].GetBool(), read[1].GetBool(), read[2].GetBool()})
	})

	t.Run("coils of separate tags out of one block", func(t *testing.T) {
		// Two coils eight apart, merged into one nine-coil read: the first sits in bit 0 of the
		// first byte, the second in bit 0 of the second.
		result, sent := runRead(t, DefaultConfiguration(),
			[]addressedTag{{"low", "coil:1:BOOL"}, {"high", "coil:9:BOOL"}},
			readWriteModel.NewModbusPDUReadCoilsResponse([]byte{0x01, 0x00}))

		require.NoError(t, result.GetErr())
		require.Len(t, sent, 1)
		response := result.GetResponse()
		require.NotNil(t, response)
		assert.True(t, response.GetValue("low").GetBool())
		assert.False(t, response.GetValue("high").GetBool())
	})
}

// An exception code the specification doesn't define still says the device refused; only the way
// it spelled that is unknown. The tag has to come back as REMOTE_ERROR rather than as a code that
// claims to know more than the response does.
func TestReader_anUnmappedExceptionCodeIsARemoteError(t *testing.T) {
	result, _ := runRead(t, DefaultConfiguration(),
		[]addressedTag{{"tag", "holding-register:1:REAL"}},
		// 0x09 is one of the gaps the specification leaves between the codes it defines.
		readWriteModel.NewModbusPDUError(readWriteModel.ModbusErrorCode(9)))

	require.Error(t, result.GetErr())
	response := result.GetResponse()
	require.NotNil(t, response)
	assert.Equal(t, apiModel.PlcResponseCode_REMOTE_ERROR, response.GetResponseCode("tag"))
}

// A tag that failed still carries a value, a null one, so that a caller who reads the value before
// looking at the code gets a PlcValue rather than something to dereference (as s7's reader does).
func TestReader_aFailedTagCarriesANullValue(t *testing.T) {
	result, _ := runRead(t, DefaultConfiguration(),
		[]addressedTag{{"tag", "holding-register:1:REAL"}},
		readWriteModel.NewModbusPDUError(readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS))

	require.Error(t, result.GetErr())
	response := result.GetResponse()
	require.NotNil(t, response)
	require.Equal(t, apiModel.PlcResponseCode_INVALID_ADDRESS, response.GetResponseCode("tag"))
	value := response.GetValue("tag")
	require.NotNil(t, value)
	assert.Equal(t, apiValues.NULL, value.GetPlcValueType())
}
