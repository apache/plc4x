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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

// A caller whose context expired abandons the result channel without reading.
// The send-failure result then fills the single-slot buffer; when the still
// registered expectation later times out, its error handler must not block
// forever on the full channel — those blocked handlers pile up in the codec's
// WaitGroup and wedge Disconnect indefinitely.
func TestReader_lateTimeoutAfterFailedSendMustNotBlock(t *testing.T) {
	codec := newCaptureCodec(errors.New("send failed: broken pipe"))
	reader := NewReader(DefaultConfiguration(), codec)
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
	reader := NewReader(DefaultConfiguration(), codec)
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
	reader := NewReader(configuration, codec)
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
	reader := NewReader(DefaultConfiguration(), newCaptureCodec(nil))
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
		reader := NewReader(configuration, newCaptureCodec(nil))
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
