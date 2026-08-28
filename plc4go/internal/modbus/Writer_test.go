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
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func testWriteRequest(t *testing.T, writer *Writer) apiModel.PlcWriteRequest {
	t.Helper()
	tag := NewTag(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_UINT)
	return spiModel.NewDefaultPlcWriteRequest(
		map[string]apiModel.PlcTag{"tag": tag},
		[]string{"tag"},
		map[string]apiValues.PlcValue{"tag": spiValues.NewPlcUINT(42)},
		writer,
		nil,
	)
}

// A failed send must deliver an error result to the caller instead of only
// logging it — a caller without its own deadline otherwise waits forever on
// a channel that never receives anything.
func TestWriter_failedSendDeliversErrorResult(t *testing.T) {
	codec := newCaptureCodec(errors.New("send failed: broken pipe"))
	writer := NewWriter(DefaultConfiguration(), codec)

	results := writer.Write(testutils.TestContext(t), testWriteRequest(t, writer))

	select {
	case result := <-results:
		if result.GetErr() == nil {
			t.Fatal("expected an error result for the failed send")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no result delivered for a failed send")
	}
}

// When the caller has abandoned the result channel, a duplicate error-handler
// invocation (e.g. expectation timeout racing a handled message, or the
// disconnect fan-out) must not block forever on the full single-slot buffer.
func TestWriter_duplicateErrorHandlerMustNotBlock(t *testing.T) {
	codec := newCaptureCodec(nil)
	writer := NewWriter(DefaultConfiguration(), codec)

	_ = writer.Write(testutils.TestContext(t), testWriteRequest(t, writer))

	var handlers capturedHandlers
	select {
	case handlers = <-codec.handlers:
	case <-time.After(time.Second):
		t.Fatal("SendRequest was never invoked")
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = handlers.handleError(errors.New("timeout"))
		_ = handlers.handleError(errors.New("disconnected"))
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("duplicate error handler blocked on the abandoned result channel")
	}
}

func writeRequestFor(t *testing.T, writer *Writer, tag apiModel.PlcTag, value apiValues.PlcValue) apiModel.PlcWriteRequest {
	t.Helper()
	values := map[string]apiValues.PlcValue{}
	if value != nil {
		values["tag"] = value
	}
	return spiModel.NewDefaultPlcWriteRequest(
		map[string]apiModel.PlcTag{"tag": tag},
		[]string{"tag"},
		values,
		writer,
		nil,
	)
}

func capturedRequestPdu(t *testing.T, tag apiModel.PlcTag, value apiValues.PlcValue) (readWriteModel.ModbusTcpADU, capturedHandlers) {
	t.Helper()
	return capturedRequestPduWith(t, DefaultConfiguration(), tag, value)
}

func capturedRequestPduWith(t *testing.T, configuration Configuration, tag apiModel.PlcTag, value apiValues.PlcValue) (readWriteModel.ModbusTcpADU, capturedHandlers) {
	t.Helper()
	codec := newCaptureCodec(nil)
	writer := NewWriter(configuration, codec)
	_ = writer.Write(testutils.TestContext(t), writeRequestFor(t, writer, tag, value))

	select {
	case handlers := <-codec.handlers:
		adu, ok := handlers.message.(readWriteModel.ModbusTcpADU)
		require.True(t, ok, "expected a ModbusTcpADU, got %T", handlers.message)
		return adu, handlers
	case <-time.After(time.Second):
		t.Fatal("SendRequest was never invoked")
		return nil, capturedHandlers{}
	}
}

// Every coil write goes out as FC 0x0F, one coil or many: plc4j's getWriteRequestPdu builds a
// WriteMultipleCoilsRequest for all of them, and FC 0x05 only ever turns up in the response.
func TestWriter_coilWritesUseWriteMultipleCoils(t *testing.T) {
	for _, test := range []struct {
		name     string
		quantity uint16
		value    apiValues.PlcValue
		expected []byte
	}{
		{"a single coil, on", 1, spiValues.NewPlcBOOL(true), []byte{0x01}},
		{"a single coil, off", 1, spiValues.NewPlcBOOL(false), []byte{0x00}},
		{
			// The first coil is the least significant bit of the first byte.
			"two coils", 2,
			spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcBOOL(true), spiValues.NewPlcBOOL(false)}),
			[]byte{0x01},
		},
		{
			"a coil in the second byte", 10,
			spiValues.NewPlcList([]apiValues.PlcValue{
				spiValues.NewPlcBOOL(false), spiValues.NewPlcBOOL(false), spiValues.NewPlcBOOL(false),
				spiValues.NewPlcBOOL(false), spiValues.NewPlcBOOL(false), spiValues.NewPlcBOOL(false),
				spiValues.NewPlcBOOL(false), spiValues.NewPlcBOOL(false), spiValues.NewPlcBOOL(true),
				spiValues.NewPlcBOOL(false),
			}),
			[]byte{0x00, 0x01},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			adu, _ := capturedRequestPdu(t, NewTag(Coil, 3, test.quantity, readWriteModel.ModbusDataType_BOOL), test.value)
			pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUWriteMultipleCoilsRequest)
			require.True(t, ok, "expected a ModbusPDUWriteMultipleCoilsRequest, got %T", adu.GetPdu())
			assert.Equal(t, uint16(2), pdu.GetStartingAddress())
			assert.Equal(t, test.quantity, pdu.GetQuantity())
			// The byte count the PDU announces is the length of this payload, so it has to cover
			// exactly the addressed coils.
			assert.Equal(t, test.expected, pdu.GetValue())
		})
	}
}

// A payload that doesn't cover the addressed coils would go out as a frame whose byte count and
// quantity contradict each other.
func TestWriter_coilWriteRejectsAMismatchedPayload(t *testing.T) {
	writer := NewWriter(DefaultConfiguration(), newCaptureCodec(nil))
	value := spiValues.NewPlcList([]apiValues.PlcValue{spiValues.NewPlcBOOL(true), spiValues.NewPlcBOOL(false)})
	results := writer.Write(testutils.TestContext(t), writeRequestFor(t, writer, NewTag(Coil, 3, 20, readWriteModel.ModbusDataType_BOOL), value))

	select {
	case result := <-results:
		require.Error(t, result.GetErr())
		assert.Contains(t, result.GetErr().Error(), "doesn't match the number of addressed coils")
	case <-time.After(2 * time.Second):
		t.Fatal("no result delivered for a mismatched coil payload")
	}
}

// A value that isn't a coil state at all has no place in a coil write.
func TestWriter_coilWriteRejectsANonBoolValue(t *testing.T) {
	writer := NewWriter(DefaultConfiguration(), newCaptureCodec(nil))
	results := writer.Write(testutils.TestContext(t), writeRequestFor(t, writer, NewTag(Coil, 3, 1, readWriteModel.ModbusDataType_BOOL), spiValues.NewPlcSTRING("nope")))

	select {
	case result := <-results:
		require.Error(t, result.GetErr())
		assert.Contains(t, result.GetErr().Error(), "BOOL")
	case <-time.After(2 * time.Second):
		t.Fatal("no result delivered for a non-BOOL coil value")
	}
}

// A value that fits into a single register goes out as FC 0x06.
func TestWriter_singleRegisterUsesWriteSingleRegister(t *testing.T) {
	adu, _ := capturedRequestPdu(t, NewTag(HoldingRegister, 3, 1, readWriteModel.ModbusDataType_UINT), spiValues.NewPlcUINT(42))
	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUWriteSingleRegisterRequest)
	require.True(t, ok, "expected a ModbusPDUWriteSingleRegisterRequest, got %T", adu.GetPdu())
	assert.Equal(t, uint16(2), pdu.GetAddress())
	assert.Equal(t, uint16(42), pdu.GetValue())
}

// Anything wider than one register still needs FC 0x10.
func TestWriter_multiWordValueUsesWriteMultipleHoldingRegisters(t *testing.T) {
	adu, _ := capturedRequestPdu(t, NewTag(HoldingRegister, 3, 1, readWriteModel.ModbusDataType_REAL), spiValues.NewPlcREAL(2.5))
	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUWriteMultipleHoldingRegistersRequest)
	require.True(t, ok, "expected a ModbusPDUWriteMultipleHoldingRegistersRequest, got %T", adu.GetPdu())
	assert.Equal(t, uint16(2), pdu.GetStartingAddress())
	assert.Equal(t, uint16(2), pdu.GetQuantity())
}

func TestWriter_singleWriteEchoValidation(t *testing.T) {
	writer := NewWriter(DefaultConfiguration(), newCaptureCodec(nil))
	tests := []struct {
		name     string
		request  readWriteModel.ModbusPDU
		response readWriteModel.ModbusPDU
		expected apiModel.PlcResponseCode
	}{
		{
			"coil echo matches",
			readWriteModel.NewModbusPDUWriteSingleCoilRequest(2, 0xFF00),
			readWriteModel.NewModbusPDUWriteSingleCoilResponse(2, 0xFF00),
			apiModel.PlcResponseCode_OK,
		},
		{
			"coil echoes a different value",
			readWriteModel.NewModbusPDUWriteSingleCoilRequest(2, 0xFF00),
			readWriteModel.NewModbusPDUWriteSingleCoilResponse(2, 0x0000),
			apiModel.PlcResponseCode_REMOTE_ERROR,
		},
		{
			"coil echoes a different address",
			readWriteModel.NewModbusPDUWriteSingleCoilRequest(2, 0xFF00),
			readWriteModel.NewModbusPDUWriteSingleCoilResponse(3, 0xFF00),
			apiModel.PlcResponseCode_REMOTE_ERROR,
		},
		{
			"register echo matches",
			readWriteModel.NewModbusPDUWriteSingleRegisterRequest(2, 42),
			readWriteModel.NewModbusPDUWriteSingleRegisterResponse(2, 42),
			apiModel.PlcResponseCode_OK,
		},
		{
			"register echoes a different value",
			readWriteModel.NewModbusPDUWriteSingleRegisterRequest(2, 42),
			readWriteModel.NewModbusPDUWriteSingleRegisterResponse(2, 43),
			apiModel.PlcResponseCode_REMOTE_ERROR,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := writeRequestFor(t, writer, NewTag(HoldingRegister, 3, 1, readWriteModel.ModbusDataType_UINT), spiValues.NewPlcUINT(42))
			response, err := writer.ToPlc4xWriteResponse(
				readWriteModel.NewModbusTcpADU(1, 1, test.request),
				readWriteModel.NewModbusTcpADU(1, 1, test.response),
				request,
			)
			require.NoError(t, err)
			assert.Equal(t, test.expected, response.GetResponseCode("tag"))
		})
	}
}

// A response PDU that doesn't belong to the request we sent must be an error rather than a panic
// on an unchecked type assertion.
func TestWriter_mismatchedRequestPduIsAnError(t *testing.T) {
	writer := NewWriter(DefaultConfiguration(), newCaptureCodec(nil))
	request := writeRequestFor(t, writer, NewTag(HoldingRegister, 3, 1, readWriteModel.ModbusDataType_UINT), spiValues.NewPlcUINT(42))
	response, err := writer.ToPlc4xWriteResponse(
		readWriteModel.NewModbusTcpADU(1, 1, readWriteModel.NewModbusPDUWriteSingleCoilRequest(2, 0xFF00)),
		readWriteModel.NewModbusTcpADU(1, 1, readWriteModel.NewModbusPDUWriteSingleRegisterResponse(2, 42)),
		request,
	)
	assert.Error(t, err)
	assert.Nil(t, response)
}

// A message that isn't a modbus ADU must never reach an unchecked type assertion.
func TestWriter_nonAduMessageIsRejected(t *testing.T) {
	_, handlers := capturedRequestPdu(t, NewTag(HoldingRegister, 3, 1, readWriteModel.ModbusDataType_UINT), spiValues.NewPlcUINT(42))

	assert.False(t, handlers.acceptsMessage(notAnAdu{}), "a foreign message must not be accepted")
	assert.NoError(t, handlers.handleMessage(notAnAdu{}))
}

// A panic anywhere in the request goroutine must be turned into an error result instead of
// tearing down the whole process.
func TestWriter_panicIsDeliveredAsAnErrorResult(t *testing.T) {
	codec := newCaptureCodec(nil)
	writer := NewWriter(DefaultConfiguration(), codec)
	// No value for the tag: serializing it dereferences a nil PlcValue.
	results := writer.Write(testutils.TestContext(t), writeRequestFor(t, writer, NewTag(HoldingRegister, 3, 1, readWriteModel.ModbusDataType_UINT), nil))

	select {
	case result := <-results:
		require.Error(t, result.GetErr())
		assert.Contains(t, result.GetErr().Error(), "panic-ed")
	case <-time.After(2 * time.Second):
		t.Fatal("no result delivered for the panicking request")
	}
}

// A tag may name the unit it is addressed at; without one the connection's default is used
// (plc4j ModbusTcpConnection.getUnitId).
func TestWriter_resolvesTheUnitIdentifier(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.unitIdentifier = 17

	t.Run("the connection default", func(t *testing.T) {
		adu, _ := capturedRequestPduWith(t, configuration, parseTag(t, "holding-register:3:UINT"), spiValues.NewPlcUINT(42))
		assert.Equal(t, uint8(17), adu.GetUnitIdentifier())
	})
	t.Run("the one the tag names", func(t *testing.T) {
		adu, _ := capturedRequestPduWith(t, configuration, parseTag(t, "holding-register:3:UINT{unit-id: 5}"), spiValues.NewPlcUINT(42))
		assert.Equal(t, uint8(5), adu.GetUnitIdentifier())
	})
}

// The value is encoded in the byte order the connection was configured with, or the one the tag
// names for itself (plc4j ModbusTcpConnection.getEffectiveByteOrder).
func TestWriter_encodesInTheResolvedByteOrder(t *testing.T) {
	expected := map[ByteOrder][]byte{
		BigEndianOrder:            {0x01, 0x02, 0x03, 0x04},
		LittleEndianOrder:         {0x04, 0x03, 0x02, 0x01},
		BigEndianByteSwapOrder:    {0x02, 0x01, 0x04, 0x03},
		LittleEndianByteSwapOrder: {0x03, 0x04, 0x01, 0x02},
	}
	written := func(t *testing.T, configuration Configuration, address string) []byte {
		t.Helper()
		adu, _ := capturedRequestPduWith(t, configuration, parseTag(t, address), spiValues.NewPlcUDINT(0x01020304))
		pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUWriteMultipleHoldingRegistersRequest)
		require.True(t, ok, "expected a write-multiple-holding-registers request, got %T", adu.GetPdu())
		return pdu.GetValue()
	}

	for byteOrder, bytes := range expected {
		t.Run("connection default "+byteOrder.String(), func(t *testing.T) {
			configuration := DefaultConfiguration()
			configuration.defaultPayloadByteOrder = byteOrder
			assert.Equal(t, bytes, written(t, configuration, "holding-register:3:UDINT"))
		})
		t.Run("tag override "+byteOrder.String(), func(t *testing.T) {
			// The connection says one thing, the tag another - the tag wins.
			configuration := DefaultConfiguration()
			configuration.defaultPayloadByteOrder = LittleEndianOrder
			assert.Equal(t, bytes, written(t, configuration, "holding-register:3:UDINT{byte-order: '"+byteOrder.String()+"'}"))
		})
	}
}

// A string is written with the length its address declared, padded out to it, and covers the
// registers that length needs.
func TestWriter_writesAStringOfTheDeclaredLength(t *testing.T) {
	adu, _ := capturedRequestPdu(t, parseTag(t, "holding-register:3:STRING(6)"), spiValues.NewPlcSTRING("hello"))

	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUWriteMultipleHoldingRegistersRequest)
	require.True(t, ok, "expected a write-multiple-holding-registers request, got %T", adu.GetPdu())
	assert.Equal(t, uint16(2), pdu.GetStartingAddress())
	assert.Equal(t, uint16(3), pdu.GetQuantity(), "a STRING(6) fills three registers")
	assert.Equal(t, []byte{'h', 'e', 'l', 'l', 'o', 0x00}, pdu.GetValue())
}

// The quantity a register write announces and the byte count the payload carries have to agree.
// An odd-length payload used to go out as quantity 2 with a byte count of 3, a frame no conforming
// server accepts; plc4j's getWriteRequestPdu throws instead of sending it.
func TestWriter_registerWriteRejectsAPayloadThatIsntWholeRegisters(t *testing.T) {
	writer := NewWriter(DefaultConfiguration(), newCaptureCodec(nil))
	tag := parseTag(t, "holding-register:1:STRING(3)")
	results := writer.Write(testutils.TestContext(t), writeRequestFor(t, writer, tag, spiValues.NewPlcSTRING("abc")))

	select {
	case result := <-results:
		require.Error(t, result.GetErr())
		assert.Contains(t, result.GetErr().Error(), "doesn't match the number of addressed registers")
	case <-time.After(2 * time.Second):
		t.Fatal("no result delivered for an odd-length payload")
	}
}

// An odd number of strings fills its registers exactly, so it goes out with a quantity that
// matches the byte count.
func TestWriter_writesAnOddNumberOfStrings(t *testing.T) {
	tag := parseTag(t, "holding-register:1[0..2]:STRING(4)")
	value := spiValues.NewPlcList([]apiValues.PlcValue{
		spiValues.NewPlcSTRING("ab"), spiValues.NewPlcSTRING("cd"), spiValues.NewPlcSTRING("ef"),
	})

	adu, _ := capturedRequestPdu(t, tag, value)

	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUWriteMultipleHoldingRegistersRequest)
	require.True(t, ok, "expected a ModbusPDUWriteMultipleHoldingRegistersRequest, got %T", adu.GetPdu())
	assert.Equal(t, uint16(6), pdu.GetQuantity())
	assert.Len(t, pdu.GetValue(), 12, "the byte count has to be twice the quantity")
}

// Several BOOLs share a register, so the quantity a write announces is the number of registers the
// packed bits occupy rather than one per value.
func TestWriter_writesPackedBoolsAsOneRegister(t *testing.T) {
	tag := parseTag(t, "holding-register:1[0..2]:BOOL")
	value := spiValues.NewPlcList([]apiValues.PlcValue{
		spiValues.NewPlcBOOL(true), spiValues.NewPlcBOOL(false), spiValues.NewPlcBOOL(true),
	})

	adu, _ := capturedRequestPdu(t, tag, value)

	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUWriteSingleRegisterRequest)
	require.True(t, ok, "expected a ModbusPDUWriteSingleRegisterRequest, got %T", adu.GetPdu())
	assert.Equal(t, uint16(0), pdu.GetAddress())
}
