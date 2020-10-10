//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package readwrite

import (
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"reflect"
	"strconv"
)

// Constant values.
const PROTOCOLIDENTIFIER uint16 = 0x0000

// The data-structure of this message
type ModbusTcpADU struct {
	transactionIdentifier uint16
	unitIdentifier        uint8
	pdu                   ModbusPDU
}

// The corresponding interface
type IModbusTcpADU interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewModbusTcpADU(transactionIdentifier uint16, unitIdentifier uint8, pdu ModbusPDU) spi.Message {
	return &ModbusTcpADU{transactionIdentifier: transactionIdentifier, unitIdentifier: unitIdentifier, pdu: pdu}
}

func (m ModbusTcpADU) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Simple field (transactionIdentifier)
	lengthInBits += 16

	// Const Field (protocolIdentifier)
	lengthInBits += 16

	// Implicit Field (length)
	lengthInBits += 16

	// Simple field (unitIdentifier)
	lengthInBits += 8

	// Simple field (pdu)
	lengthInBits += m.pdu.LengthInBits()

	return lengthInBits
}

func (m ModbusTcpADU) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusTcpADUParse(io spi.ReadBuffer, response bool) (spi.Message, error) {

	// Simple Field (transactionIdentifier)
	var transactionIdentifier uint16 = io.ReadUint16(16)

	// Const Field (protocolIdentifier)
	var protocolIdentifier uint16 = io.ReadUint16(16)
	if protocolIdentifier != PROTOCOLIDENTIFIER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(PROTOCOLIDENTIFIER)) + " but got " + strconv.Itoa(int(protocolIdentifier)))
	}

	// Implicit Field (length) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var _ uint16 = io.ReadUint16(16)

	// Simple Field (unitIdentifier)
	var unitIdentifier uint8 = io.ReadUint8(8)

	// Simple Field (pdu)
	_pduMessage, _err := ModbusPDUParse(io, bool(response))
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'pdu'. " + _err.Error())
	}
	var pdu ModbusPDU
	pdu, _pduOk := _pduMessage.(ModbusPDU)
	if !_pduOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_pduMessage).Name() + " to ModbusPDU")
	}

	// Create the instance
	return NewModbusTcpADU(transactionIdentifier, unitIdentifier, pdu), nil
}

func (m ModbusTcpADU) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IModbusPDU); ok {

			// Simple Field (transactionIdentifier)
			var transactionIdentifier uint16 = m.transactionIdentifier
			io.WriteUint16(16, (transactionIdentifier))

			// Const Field (protocolIdentifier)
			io.WriteUint16(16, 0x0000)

			// Implicit Field (length) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
			var length uint16 = ((m.pdu.LengthInBytes()) + (1))
			io.WriteUint16(16, (length))

			// Simple Field (unitIdentifier)
			var unitIdentifier uint8 = m.unitIdentifier
			io.WriteUint8(8, (unitIdentifier))

			// Simple Field (pdu)
			var pdu ModbusPDU = m.pdu
			pdu.Serialize(io)
		}
	}
	serializeFunc(m)
}
