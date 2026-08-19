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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
)

// A configuration nobody told otherwise speaks Modbus TCP, which is what every caller of this
// package relied on before the RTU flavor grew a framing of its own.
func TestConfigurationFlavorDefaultsToTcp(t *testing.T) {
	assert.Equal(t, flavorTcp, DefaultConfiguration().flavor)
	assert.Equal(t, flavorTcp, Configuration{}.flavor)
	assert.IsType(t, tcpAduFactory{}, DefaultConfiguration().adus())
	assert.IsType(t, rtuAduFactory{}, DefaultConfiguration().withFlavor(flavorRtu).adus())
	// withFlavor hands back a copy rather than editing the configuration in place.
	configuration := DefaultConfiguration()
	_ = configuration.withFlavor(flavorRtu)
	assert.Equal(t, flavorTcp, configuration.flavor)
}

func TestModbusFlavorString(t *testing.T) {
	assert.Equal(t, "MODBUS_TCP", flavorTcp.String())
	assert.Equal(t, "MODBUS_RTU", flavorRtu.String())
}

func TestAduFactoryBuildsTheFlavorsFrame(t *testing.T) {
	pdu := readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 10)

	tcpRequest := tcpAduFactory{}.buildRequest(42, 3, pdu)
	tcpAdu, ok := tcpRequest.(readWriteModel.ModbusTcpADU)
	require.True(t, ok, "got %T", tcpRequest)
	assert.Equal(t, uint16(42), tcpAdu.GetTransactionIdentifier())
	assert.Equal(t, uint8(3), tcpAdu.GetUnitIdentifier())
	assert.Equal(t, pdu, tcpAdu.GetPdu())

	// RTU has no transaction identifier at all - the unit identifier is the station address.
	rtuRequest := rtuAduFactory{}.buildRequest(42, 3, pdu)
	rtuAdu, ok := rtuRequest.(readWriteModel.ModbusRtuADU)
	require.True(t, ok, "got %T", rtuRequest)
	assert.Equal(t, uint8(3), rtuAdu.GetAddress())
	assert.Equal(t, pdu, rtuAdu.GetPdu())
}

func TestTcpAduFactoryAcceptsResponse(t *testing.T) {
	request := readWriteModel.NewModbusTcpADU(42, 3, readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 1))
	response := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})

	tests := []struct {
		name     string
		response spi.Message
		want     bool
	}{
		{name: "the echoed transaction and unit identifier", response: readWriteModel.NewModbusTcpADU(42, 3, response), want: true},
		{name: "another transaction identifier", response: readWriteModel.NewModbusTcpADU(43, 3, response), want: false},
		{name: "another unit identifier", response: readWriteModel.NewModbusTcpADU(42, 4, response), want: false},
		{name: "the same frame in the wrong flavor", response: readWriteModel.NewModbusRtuADU(3, response), want: false},
		{name: "not an ADU at all", response: notAnAdu{}, want: false},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, tcpAduFactory{}.acceptsResponse(request, testCase.response))
		})
	}
}

// An RTU frame carries no transaction identifier, so a response is correlated by station address
// and function code, the way plc4j's ModbusRtuConnection.handleIncomingMessage does it.
func TestRtuAduFactoryAcceptsResponse(t *testing.T) {
	request := readWriteModel.NewModbusRtuADU(3, readWriteModel.NewModbusPDUReadHoldingRegistersRequest(0, 1))

	tests := []struct {
		name     string
		response spi.Message
		want     bool
	}{
		{
			name:     "the same address answering the same function code",
			response: readWriteModel.NewModbusRtuADU(3, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})),
			want:     true,
		},
		{
			name:     "an exception from the same address answers whatever was asked",
			response: readWriteModel.NewModbusRtuADU(3, readWriteModel.NewModbusPDUError(readWriteModel.ModbusErrorCode_ILLEGAL_DATA_ADDRESS)),
			want:     true,
		},
		{
			name:     "another station on the same line",
			response: readWriteModel.NewModbusRtuADU(4, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})),
			want:     false,
		},
		{
			name:     "a late response to a request that already timed out",
			response: readWriteModel.NewModbusRtuADU(3, readWriteModel.NewModbusPDUReadCoilsResponse([]byte{0x01})),
			want:     false,
		},
		{
			name:     "the same frame in the wrong flavor",
			response: readWriteModel.NewModbusTcpADU(42, 3, readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})),
			want:     false,
		},
		{name: "not an ADU at all", response: notAnAdu{}, want: false},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.want, rtuAduFactory{}.acceptsResponse(request, testCase.response))
		})
	}
}

func TestAduFactoryExtractPdu(t *testing.T) {
	pdu := readWriteModel.NewModbusPDUReadHoldingRegistersResponse([]byte{0x00, 0x2A})

	extracted, err := tcpAduFactory{}.extractPdu(readWriteModel.NewModbusTcpADU(1, 1, pdu))
	require.NoError(t, err)
	assert.Equal(t, pdu, extracted)

	extracted, err = rtuAduFactory{}.extractPdu(readWriteModel.NewModbusRtuADU(1, pdu))
	require.NoError(t, err)
	assert.Equal(t, pdu, extracted)

	// Handing a flavor the wrong ADU must be an error rather than a failed type assertion.
	_, err = tcpAduFactory{}.extractPdu(readWriteModel.NewModbusRtuADU(1, pdu))
	assert.ErrorContains(t, err, "ModbusTcpADU")
	_, err = rtuAduFactory{}.extractPdu(readWriteModel.NewModbusTcpADU(1, 1, pdu))
	assert.ErrorContains(t, err, "ModbusRtuADU")
	_, err = rtuAduFactory{}.extractPdu(notAnAdu{})
	assert.ErrorContains(t, err, "ModbusRtuADU")
}
