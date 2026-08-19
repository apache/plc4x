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
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/testutils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

func rtuConfiguration() Configuration {
	return DefaultConfiguration().withFlavor(flavorRtu)
}

func awaitHandlers(t *testing.T, codec *captureCodec) capturedHandlers {
	t.Helper()
	select {
	case handlers := <-codec.handlers:
		return handlers
	case <-time.After(time.Second):
		t.Fatal("SendRequest was never invoked")
		return capturedHandlers{}
	}
}

// A read on an RTU connection has to put an RTU ADU on the wire. Before the flavor existed the
// reader built an MBAP-framed one unconditionally, which the serial line has no idea what to do
// with.
func TestReader_rtuFlavorSendsAnRtuAdu(t *testing.T) {
	codec := newCaptureCodec(nil)
	configuration := rtuConfiguration()
	configuration.unitIdentifier = 7
	reader := NewReader(configuration, codec)
	tag := NewTag(HoldingRegister, 1, 1, readWriteModel.ModbusDataType_UINT)
	request := spiModel.NewDefaultPlcReadRequest(
		map[string]apiModel.PlcTag{"tag": tag}, []string{"tag"}, reader, nil)

	results := reader.Read(testutils.TestContext(t), request)
	handlers := awaitHandlers(t, codec)

	requestAdu, ok := handlers.message.(readWriteModel.ModbusRtuADU)
	require.True(t, ok, "expected a ModbusRtuADU, got %T", handlers.message)
	assert.Equal(t, uint8(7), requestAdu.GetAddress())
	assert.IsType(t, readWriteModel.NewModbusPDUReadHoldingRegistersRequest(1, 1), requestAdu.GetPdu())

	// An MBAP-framed response is not something an RTU connection could have asked for.
	assert.False(t, handlers.acceptsMessage(readWriteModel.NewModbusTcpADU(0, 7,
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A}))))
	// Neither is a frame from another station on the same line.
	assert.False(t, handlers.acceptsMessage(readWriteModel.NewModbusRtuADU(8,
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A}))))

	response := readWriteModel.NewModbusRtuADU(7,
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

// The same for writes: the request goes out RTU-framed and the echo comes back RTU-framed.
func TestWriter_rtuFlavorSendsAnRtuAdu(t *testing.T) {
	codec := newCaptureCodec(nil)
	configuration := rtuConfiguration()
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

	requestAdu, ok := handlers.message.(readWriteModel.ModbusRtuADU)
	require.True(t, ok, "expected a ModbusRtuADU, got %T", handlers.message)
	assert.Equal(t, uint8(7), requestAdu.GetAddress())
	// A single register goes out as FC 0x06, which the device answers by echoing it back.
	requestPdu, ok := requestAdu.GetPdu().(readWriteModel.ModbusPDUWriteSingleRegisterRequest)
	require.True(t, ok, "expected a write-single-register request, got %T", requestAdu.GetPdu())
	echo := readWriteModel.NewModbusPDUWriteSingleRegisterResponse(requestPdu.GetAddress(), requestPdu.GetValue())

	assert.False(t, handlers.acceptsMessage(readWriteModel.NewModbusTcpADU(0, 7, echo)))

	response := readWriteModel.NewModbusRtuADU(7, echo)
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

// Ping goes through the same factory, so an RTU connection pings with an RTU ADU too.
func TestConnection_rtuFlavorPingsWithAnRtuAdu(t *testing.T) {
	codec := newCaptureCodec(nil)
	connection := NewConnection(rtuConfiguration(), codec, map[string][]string{}, NewTagHandler())

	pingDone := make(chan error, 1)
	go func() { pingDone <- connection.Ping(testutils.TestContext(t)) }()
	handlers := awaitHandlers(t, codec)

	pingAdu, ok := handlers.message.(readWriteModel.ModbusRtuADU)
	require.True(t, ok, "expected a ModbusRtuADU, got %T", handlers.message)
	assert.Equal(t, defaultUnitIdentifier, pingAdu.GetAddress())

	response := readWriteModel.NewModbusRtuADU(pingAdu.GetAddress(),
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
