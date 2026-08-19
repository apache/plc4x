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

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi/errors"
	"github.com/apache/plc4x/plc4go/spi/testutils"
)

func testConnection(configuration Configuration, codec *captureCodec) *Connection {
	return NewConnection(configuration, codec, map[string][]string{}, NewTagHandler())
}

// startPing runs a ping and hands back the request that went out together with a channel carrying
// its outcome. Ping blocks until the request is answered, so it has to run on its own goroutine.
func startPing(t *testing.T, connection *Connection, codec *captureCodec) (readWriteModel.ModbusTcpADU, capturedHandlers, <-chan error) {
	t.Helper()
	pinged := make(chan error, 1)
	go func() { pinged <- connection.Ping(testutils.TestContext(t)) }()

	select {
	case handlers := <-codec.handlers:
		adu, ok := handlers.message.(readWriteModel.ModbusTcpADU)
		require.True(t, ok, "expected a ModbusTcpADU, got %T", handlers.message)
		return adu, handlers, pinged
	case <-time.After(2 * time.Second):
		t.Fatal("the ping never sent a request")
		return nil, capturedHandlers{}, nil
	}
}

func awaitPing(t *testing.T, pinged <-chan error) error {
	t.Helper()
	select {
	case err := <-pinged:
		return err
	case <-time.After(2 * time.Second):
		t.Fatal("the ping never returned")
		return nil
	}
}

// A ping reads the configured ping address instead of sending the optional diagnostic function
// code 0x08, which a good many devices don't implement (plc4j ModbusTcpConnection.onPing).
func TestConnection_pingReadsThePingAddress(t *testing.T) {
	codec := newCaptureCodec(nil)
	connection := testConnection(DefaultConfiguration(), codec)

	adu, handlers, pinged := startPing(t, connection, codec)

	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUReadHoldingRegistersRequest)
	require.True(t, ok, "expected a read of the default ping address 4x00001, got %T", adu.GetPdu())
	assert.Equal(t, uint16(0), pdu.GetStartingAddress(), "4x00001 is address 0 on the wire")
	assert.Equal(t, uint16(1), pdu.GetQuantity())
	assert.Equal(t, uint8(1), adu.GetUnitIdentifier())

	require.NoError(t, handlers.handleMessage(readWriteModel.NewModbusTcpADU(
		adu.GetTransactionIdentifier(), adu.GetUnitIdentifier(),
		readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x01}))))
	assert.NoError(t, awaitPing(t, pinged))
}

// The ping address is configurable, coil addresses included.
func TestConnection_pingUsesTheConfiguredAddress(t *testing.T) {
	configuration := DefaultConfiguration()
	configuration.pingAddress = "coil:5:BOOL"
	codec := newCaptureCodec(nil)
	connection := testConnection(configuration, codec)

	adu, handlers, pinged := startPing(t, connection, codec)

	pdu, ok := adu.GetPdu().(readWriteModel.ModbusPDUReadCoilsRequest)
	require.True(t, ok, "expected a read of a coil, got %T", adu.GetPdu())
	assert.Equal(t, uint16(4), pdu.GetStartingAddress())

	require.NoError(t, handlers.handleMessage(readWriteModel.NewModbusTcpADU(
		adu.GetTransactionIdentifier(), adu.GetUnitIdentifier(),
		readWriteModel.NewModbusPDUReadCoilsResponse([]byte{0x01}))))
	assert.NoError(t, awaitPing(t, pinged))
}

// The ping goes to the connection's unit, unless the ping address names one of its own.
func TestConnection_pingUnitIdentifier(t *testing.T) {
	t.Run("the connection default", func(t *testing.T) {
		configuration := DefaultConfiguration()
		configuration.unitIdentifier = 17
		codec := newCaptureCodec(nil)
		adu, _, _ := startPing(t, testConnection(configuration, codec), codec)
		assert.Equal(t, uint8(17), adu.GetUnitIdentifier())
	})
	t.Run("the one the ping address names", func(t *testing.T) {
		configuration := DefaultConfiguration()
		configuration.unitIdentifier = 17
		configuration.pingAddress = "4x00001:BOOL{unit-id: 5}"
		codec := newCaptureCodec(nil)
		adu, _, _ := startPing(t, testConnection(configuration, codec), codec)
		assert.Equal(t, uint8(5), adu.GetUnitIdentifier())
	})
}

// Only a response carrying the transaction and unit identifier of the request is ours.
func TestConnection_pingOnlyAcceptsItsOwnResponse(t *testing.T) {
	codec := newCaptureCodec(nil)
	connection := testConnection(DefaultConfiguration(), codec)

	adu, handlers, _ := startPing(t, connection, codec)
	response := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x01})

	assert.False(t, handlers.acceptsMessage(notAnAdu{}), "a foreign message is not a ping response")
	assert.False(t, handlers.acceptsMessage(readWriteModel.NewModbusTcpADU(adu.GetTransactionIdentifier()+1, adu.GetUnitIdentifier(), response)))
	assert.False(t, handlers.acceptsMessage(readWriteModel.NewModbusTcpADU(adu.GetTransactionIdentifier(), adu.GetUnitIdentifier()+1, response)))
	assert.True(t, handlers.acceptsMessage(readWriteModel.NewModbusTcpADU(adu.GetTransactionIdentifier(), adu.GetUnitIdentifier(), response)))
}

// A device that answers with a modbus exception is still a device that answers, so the ping
// succeeds - the ping address is a guess that a given device need not have. plc4j answers a ping
// the same way.
func TestConnection_pingSucceedsOnAnExceptionResponse(t *testing.T) {
	codec := newCaptureCodec(nil)
	connection := testConnection(DefaultConfiguration(), codec)

	adu, handlers, pinged := startPing(t, connection, codec)
	require.NoError(t, handlers.handleMessage(readWriteModel.NewModbusTcpADU(
		adu.GetTransactionIdentifier(), adu.GetUnitIdentifier(),
		readWriteModel.NewModbusPDUError(readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS))))

	assert.NoError(t, awaitPing(t, pinged))
}

// No answer at all is what a ping is there to find out about.
func TestConnection_pingFailsOnATimeout(t *testing.T) {
	codec := newCaptureCodec(nil)
	connection := testConnection(DefaultConfiguration(), codec)

	_, handlers, pinged := startPing(t, connection, codec)
	require.NoError(t, handlers.handleError(errors.New("timeout")))

	assert.Error(t, awaitPing(t, pinged))
}

func TestConnection_pingFailsWhenTheRequestCantBeSent(t *testing.T) {
	codec := newCaptureCodec(errors.New("send failed: broken pipe"))
	connection := testConnection(DefaultConfiguration(), codec)

	assert.Error(t, connection.Ping(testutils.TestContext(t)))
}

// A ping address the driver can't turn into a request has to say so rather than silently never
// pinging.
func TestConnection_pingFailsOnAnUnusablePingAddress(t *testing.T) {
	for _, test := range []struct {
		name        string
		pingAddress string
	}{
		{"unparsable", "this is not an address"},
		{"an address outside the address space", "holding-register:70000:INT"},
		{"a quantity no request can carry", "holding-register:1:INT[126]"},
	} {
		t.Run(test.name, func(t *testing.T) {
			configuration := DefaultConfiguration()
			configuration.pingAddress = test.pingAddress
			codec := newCaptureCodec(nil)

			assert.Error(t, testConnection(configuration, codec).Ping(testutils.TestContext(t)))
		})
	}
}
