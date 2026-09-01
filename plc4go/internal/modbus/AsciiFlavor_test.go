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
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func asciiConfiguration() Configuration {
	return DefaultConfiguration().withFlavor(flavorAscii)
}

func TestAsciiFlavorPicksTheAsciiAduFactory(t *testing.T) {
	assert.Equal(t, "MODBUS_ASCII", flavorAscii.String())
	assert.IsType(t, asciiAduFactory{}, asciiConfiguration().adus())
	// withFlavor hands back a copy rather than editing the configuration in place.
	configuration := DefaultConfiguration()
	_ = configuration.withFlavor(flavorAscii)
	assert.Equal(t, flavorTcp, configuration.flavor)
}

// ASCII has no transaction identifier either - the unit identifier is the station address.
func TestAsciiAduFactoryBuildsAnAsciiAdu(t *testing.T) {
	pdu := readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 10)

	request := asciiAduFactory{}.buildRequest(42, 3, pdu)

	adu, ok := request.(readWriteModel.ModbusAsciiADU)
	require.True(t, ok, "got %T", request)
	assert.Equal(t, uint8(3), adu.GetAddress())
	assert.Equal(t, pdu, adu.GetPdu())
}

// A response is correlated by station address and function code, the way plc4j's
// ModbusAsciiConnection.handleIncomingMessage does it.
func TestAsciiAduFactoryAcceptsResponse(t *testing.T) {
	request := readWriteModel.NewModbusAsciiADU(3, readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 1))

	tests := []struct {
		name     string
		response spi.Message
		want     bool
	}{
		{
			name:     "the same address answering the same function code",
			response: readWriteModel.NewModbusAsciiADU(3, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})),
			want:     true,
		},
		{
			name:     "an exception from the same address answers whatever was asked",
			response: readWriteModel.NewModbusAsciiADU(3, readWriteModel.NewModbusPDUError(readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS)),
			want:     true,
		},
		{
			name:     "another station on the same line",
			response: readWriteModel.NewModbusAsciiADU(4, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})),
			want:     false,
		},
		{
			name:     "a late response to a request that already timed out",
			response: readWriteModel.NewModbusAsciiADU(3, readWriteModel.NewModbusPDUReadCoilsResponse([]byte{0x01})),
			want:     false,
		},
		{
			name:     "the same frame in the RTU flavor",
			response: readWriteModel.NewModbusRtuADU(3, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})),
			want:     false,
		},
		{
			name:     "the same frame in the TCP flavor",
			response: readWriteModel.NewModbusTcpADU(42, 3, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})),
			want:     false,
		},
		{name: "not an ADU at all", response: notAnAdu{}, want: false},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, asciiAduFactory{}.acceptsResponse(request, testCase.response))
		})
	}
}

func TestAsciiAduFactoryExtractPdu(t *testing.T) {
	pdu := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})

	extracted, err := asciiAduFactory{}.extractPdu(readWriteModel.NewModbusAsciiADU(1, pdu))
	require.NoError(t, err)
	assert.Equal(t, pdu, extracted)

	// Handing the flavor the wrong ADU must be an error rather than a failed type assertion.
	_, err = asciiAduFactory{}.extractPdu(readWriteModel.NewModbusRtuADU(1, pdu))
	assert.ErrorContains(t, err, "ModbusAsciiADU")
	_, err = asciiAduFactory{}.extractPdu(readWriteModel.NewModbusTcpADU(1, 1, pdu))
	assert.ErrorContains(t, err, "ModbusAsciiADU")
	_, err = asciiAduFactory{}.extractPdu(notAnAdu{})
	assert.ErrorContains(t, err, "ModbusAsciiADU")
	// ... and an ASCII ADU is not something the other flavors accept either.
	_, err = rtuAduFactory{}.extractPdu(readWriteModel.NewModbusAsciiADU(1, pdu))
	assert.ErrorContains(t, err, "ModbusRtuADU")
	_, err = tcpAduFactory{}.extractPdu(readWriteModel.NewModbusAsciiADU(1, pdu))
	assert.ErrorContains(t, err, "ModbusTcpADU")
}

// A read on an ASCII connection has to put an ASCII ADU on the wire. Before the flavor existed the
// reader built an MBAP-framed one unconditionally, which the ASCII codec has no idea what to do
// with.
func TestReader_asciiFlavorSendsAnAsciiAdu(t *testing.T) {
	codec := newCaptureCodec(nil)
	configuration := asciiConfiguration()
	configuration.unitIdentifier = 7
	reader := NewReader(configuration, codec)
	tag := NewTag(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_UINT)
	request := spiModel.NewDefaultPlcReadRequest(
		map[string]apiModel.PlcTag{"tag": tag}, []string{"tag"}, reader, nil)

	results := reader.Read(testutils.TestContext(t), request)
	handlers := awaitHandlers(t, codec)

	requestAdu, ok := handlers.message.(readWriteModel.ModbusAsciiADU)
	require.True(t, ok, "expected a ModbusAsciiADU, got %T", handlers.message)
	assert.Equal(t, uint8(7), requestAdu.GetAddress())
	assert.IsType(t, readWriteModel.NewModbusPDUReadHoldingRegistersRequest(1, 1), requestAdu.GetPdu())

	// A frame in another flavor is not something an ASCII connection could have asked for.
	assert.False(t, handlers.acceptsMessage(readWriteModel.NewModbusTcpADU(0, 7,
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A}))))
	assert.False(t, handlers.acceptsMessage(readWriteModel.NewModbusRtuADU(7,
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A}))))
	// Neither is a frame from another station on the same line.
	assert.False(t, handlers.acceptsMessage(readWriteModel.NewModbusAsciiADU(8,
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A}))))

	response := readWriteModel.NewModbusAsciiADU(7,
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A}))
	require.True(t, handlers.acceptsMessage(response))
	require.NoError(t, handlers.handleMessage(response))

	select {
	case result := <-results:
		require.NoError(t, result.GetErr())
		assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("tag"))
		assert.Equal(t, uint16(42), result.GetResponse().GetValue("tag").GetUint16())
	case <-time.After(2 * time.Second):
		t.Fatal("no result delivered")
	}
}

// The same for writes: the request goes out ASCII-framed and the echo comes back ASCII-framed.
func TestWriter_asciiFlavorSendsAnAsciiAdu(t *testing.T) {
	codec := newCaptureCodec(nil)
	configuration := asciiConfiguration()
	configuration.unitIdentifier = 7
	writer := NewWriter(configuration, codec)
	tag := NewTag(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_UINT)
	request := spiModel.NewDefaultPlcWriteRequest(
		map[string]apiModel.PlcTag{"tag": tag},
		[]string{"tag"},
		map[string]apiValues.PlcValue{"tag": spiValues.NewPlcUINT(42)},
		writer,
		nil)

	results := writer.Write(testutils.TestContext(t), request)
	handlers := awaitHandlers(t, codec)

	requestAdu, ok := handlers.message.(readWriteModel.ModbusAsciiADU)
	require.True(t, ok, "expected a ModbusAsciiADU, got %T", handlers.message)
	assert.Equal(t, uint8(7), requestAdu.GetAddress())
	// A single register goes out as FC 0x06, which the device answers by echoing it back.
	requestPdu, ok := requestAdu.GetPdu().(readWriteModel.ModbusPDUWriteSingleRegisterRequest)
	require.True(t, ok, "expected a write-single-register request, got %T", requestAdu.GetPdu())
	echo := readWriteModel.NewModbusPDUWriteSingleRegisterResponse(requestPdu.GetAddress(), requestPdu.GetValue())

	assert.False(t, handlers.acceptsMessage(readWriteModel.NewModbusRtuADU(7, echo)))

	response := readWriteModel.NewModbusAsciiADU(7, echo)
	require.True(t, handlers.acceptsMessage(response))
	require.NoError(t, handlers.handleMessage(response))

	select {
	case result := <-results:
		require.NoError(t, result.GetErr())
		assert.Equal(t, apiModel.PlcResponseCode_OK, result.GetResponse().GetResponseCode("tag"))
	case <-time.After(2 * time.Second):
		t.Fatal("no result delivered")
	}
}

// Ping goes through the same factory, so an ASCII connection pings with an ASCII ADU too.
func TestConnection_asciiFlavorPingsWithAnAsciiAdu(t *testing.T) {
	codec := newCaptureCodec(nil)
	connection := NewConnection(asciiConfiguration(), codec, map[string][]string{}, NewTagHandler())

	pingDone := make(chan error, 1)
	go func() { pingDone <- connection.Ping(testutils.TestContext(t)) }()
	handlers := awaitHandlers(t, codec)

	pingAdu, ok := handlers.message.(readWriteModel.ModbusAsciiADU)
	require.True(t, ok, "expected a ModbusAsciiADU, got %T", handlers.message)
	assert.Equal(t, defaultUnitIdentifier, pingAdu.GetAddress())

	response := readWriteModel.NewModbusAsciiADU(pingAdu.GetAddress(),
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x00}))
	require.True(t, handlers.acceptsMessage(response))
	require.NoError(t, handlers.handleMessage(response))

	select {
	case err := <-pingDone:
		require.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("ping never returned")
	}
}
