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
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/modbus/readwrite/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/errors"
)

// modbusFlavor names the framing a modbus connection speaks. plc4j has one connection class per
// flavor (ModbusTcpConnection, ModbusRtuConnection, ...); here the drivers share Connection,
// Reader and Writer and ask the flavor's aduFactory for the wrapper around the PDU instead.
type modbusFlavor uint8

const (
	// flavorTcp is the MBAP-framed flavor. It is the zero value on purpose, so a Configuration
	// nobody told otherwise speaks Modbus TCP - which is what every caller of this package did
	// before the RTU flavor grew a framing of its own.
	flavorTcp modbusFlavor = iota
	// flavorRtu is the flavor spoken on a serial line: a station address, the PDU and a CRC.
	flavorRtu
	// flavorAscii is the other flavor spoken on a serial line: a station address, the PDU and an
	// LRC, all of it spelled out in hex characters between a colon and a CR/LF. The hex layer is
	// the codec's business; the ADU is the binary one underneath it.
	flavorAscii
)

func (f modbusFlavor) String() string {
	switch f {
	case flavorRtu:
		return "MODBUS_RTU"
	case flavorAscii:
		return "MODBUS_ASCII"
	default:
		return "MODBUS_TCP"
	}
}

// adus is the factory that builds and reads the ADUs of this flavor.
func (f modbusFlavor) adus() aduFactory {
	switch f {
	case flavorRtu:
		return rtuAduFactory{}
	case flavorAscii:
		return asciiAduFactory{}
	default:
		return tcpAduFactory{}
	}
}

// aduFactory wraps a PDU in the ADU of one modbus flavor and tells that flavor's response to a
// request apart from everything else that arrives on the wire. It is deliberately small: the
// flavors differ in framing and in how a response is correlated to its request, nothing else.
type aduFactory interface {
	// buildRequest wraps the pdu in a request ADU. The transaction identifier only exists in the
	// TCP MBAP header; the flavors without one ignore it and correlate differently.
	buildRequest(transactionIdentifier uint16, unitIdentifier uint8, pdu readWriteModel.ModbusPDU) spi.Message
	// acceptsResponse reports whether an incoming message is the response to the given request.
	// Both arguments are ADUs of this flavor; anything else must be rejected rather than panic.
	acceptsResponse(request spi.Message, response spi.Message) bool
	// extractPdu pulls the PDU out of an ADU of this flavor.
	extractPdu(message spi.Message) (readWriteModel.ModbusPDU, error)
}

// aduWithPdu is all the response conversion needs of an ADU, and all the flavors have in common -
// the generated ModbusADU parent type carries no PDU accessor of its own.
type aduWithPdu interface {
	GetPdu() readWriteModel.ModbusPDU
}

// tcpAduFactory speaks the MBAP-framed flavor, where the transaction identifier the request went
// out with is echoed back and correlates the response.
type tcpAduFactory struct{}

var _ aduFactory = tcpAduFactory{}

func (tcpAduFactory) buildRequest(transactionIdentifier uint16, unitIdentifier uint8, pdu readWriteModel.ModbusPDU) spi.Message {
	return readWriteModel.NewModbusTcpADU(transactionIdentifier, unitIdentifier, pdu)
}

func (tcpAduFactory) acceptsResponse(request spi.Message, response spi.Message) bool {
	requestAdu, ok := request.(readWriteModel.ModbusTcpADU)
	if !ok {
		return false
	}
	responseAdu, ok := response.(readWriteModel.ModbusTcpADU)
	if !ok {
		// Not something we could have asked for, so it can't be our response.
		return false
	}
	return responseAdu.GetTransactionIdentifier() == requestAdu.GetTransactionIdentifier() &&
		responseAdu.GetUnitIdentifier() == requestAdu.GetUnitIdentifier()
}

func (tcpAduFactory) extractPdu(message spi.Message) (readWriteModel.ModbusPDU, error) {
	adu, ok := message.(readWriteModel.ModbusTcpADU)
	if !ok {
		return nil, errors.Errorf("response is not a ModbusTcpADU, got %T", message)
	}
	return adu.GetPdu(), nil
}

// rtuAduFactory speaks the flavor of a serial line. An RTU frame carries no transaction
// identifier, so a response is correlated the way plc4j's ModbusRtuConnection.handleIncomingMessage
// does it: the station address has to match, and a frame that is not an exception has to answer
// with the function code that was asked for. A late response to a request that already timed out
// and happens to carry the same function code is indistinguishable from the real one - the
// protocol carries nothing to tell them apart - but requests are serialized per connection, which
// bounds that to same-function-code retries.
type rtuAduFactory struct{}

var _ aduFactory = rtuAduFactory{}

func (rtuAduFactory) buildRequest(_ uint16, unitIdentifier uint8, pdu readWriteModel.ModbusPDU) spi.Message {
	return readWriteModel.NewModbusRtuADU(unitIdentifier, pdu)
}

func (rtuAduFactory) acceptsResponse(request spi.Message, response spi.Message) bool {
	requestAdu, ok := request.(readWriteModel.ModbusRtuADU)
	if !ok {
		return false
	}
	responseAdu, ok := response.(readWriteModel.ModbusRtuADU)
	if !ok {
		// Not something we could have asked for, so it can't be our response.
		return false
	}
	if responseAdu.GetAddress() != requestAdu.GetAddress() {
		return false
	}
	responsePdu := responseAdu.GetPdu()
	if responsePdu.GetErrorFlag() {
		// An exception answers whatever was asked - it carries the function code back with the
		// error bit set, so comparing function codes would reject it.
		return true
	}
	return responsePdu.GetFunctionFlag() == requestAdu.GetPdu().GetFunctionFlag()
}

func (rtuAduFactory) extractPdu(message spi.Message) (readWriteModel.ModbusPDU, error) {
	adu, ok := message.(readWriteModel.ModbusRtuADU)
	if !ok {
		return nil, errors.Errorf("response is not a ModbusRtuADU, got %T", message)
	}
	return adu.GetPdu(), nil
}

// asciiAduFactory speaks the hex-encoded flavor of a serial line. An ASCII frame carries the same
// information an RTU one does - a station address, the PDU and a checksum - so a response is
// correlated exactly the same way plc4j's ModbusAsciiConnection.handleIncomingMessage does it: the
// station address has to match, and a frame that is not an exception has to answer with the
// function code that was asked for. The caveat about a late response of the same function code
// being indistinguishable from the real one holds here too, and is bounded the same way: requests
// are serialized per connection.
type asciiAduFactory struct{}

var _ aduFactory = asciiAduFactory{}

func (asciiAduFactory) buildRequest(_ uint16, unitIdentifier uint8, pdu readWriteModel.ModbusPDU) spi.Message {
	return readWriteModel.NewModbusAsciiADU(unitIdentifier, pdu)
}

func (asciiAduFactory) acceptsResponse(request spi.Message, response spi.Message) bool {
	requestAdu, ok := request.(readWriteModel.ModbusAsciiADU)
	if !ok {
		return false
	}
	responseAdu, ok := response.(readWriteModel.ModbusAsciiADU)
	if !ok {
		// Not something we could have asked for, so it can't be our response.
		return false
	}
	if responseAdu.GetAddress() != requestAdu.GetAddress() {
		return false
	}
	responsePdu := responseAdu.GetPdu()
	if responsePdu.GetErrorFlag() {
		// An exception answers whatever was asked - it carries the function code back with the
		// error bit set, so comparing function codes would reject it.
		return true
	}
	return responsePdu.GetFunctionFlag() == requestAdu.GetPdu().GetFunctionFlag()
}

func (asciiAduFactory) extractPdu(message spi.Message) (readWriteModel.ModbusPDU, error) {
	adu, ok := message.(readWriteModel.ModbusAsciiADU)
	if !ok {
		return nil, errors.Errorf("response is not a ModbusAsciiADU, got %T", message)
	}
	return adu.GetPdu(), nil
}
