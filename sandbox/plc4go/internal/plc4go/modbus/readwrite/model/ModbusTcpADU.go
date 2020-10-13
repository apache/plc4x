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
package model

import (
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"reflect"
	"strconv"
)

// Constant values.
const ModbusTcpADU_PROTOCOLIDENTIFIER uint16 = 0x0000

// The data-structure of this message
type ModbusTcpADU struct {
	transactionIdentifier uint16
	unitIdentifier        uint8
	pdu                   IModbusPDU
}

// The corresponding interface
type IModbusTcpADU interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewModbusTcpADU(transactionIdentifier uint16, unitIdentifier uint8, pdu IModbusPDU) spi.Message {
	return &ModbusTcpADU{transactionIdentifier: transactionIdentifier, unitIdentifier: unitIdentifier, pdu: pdu}
}

func CastIModbusTcpADU(structType interface{}) IModbusTcpADU {
	castFunc := func(typ interface{}) IModbusTcpADU {
		if iModbusTcpADU, ok := typ.(IModbusTcpADU); ok {
			return iModbusTcpADU
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusTcpADU(structType interface{}) ModbusTcpADU {
	castFunc := func(typ interface{}) ModbusTcpADU {
		if sModbusTcpADU, ok := typ.(ModbusTcpADU); ok {
			return sModbusTcpADU
		}
		return ModbusTcpADU{}
	}
	return castFunc(structType)
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
	transactionIdentifier, _transactionIdentifierErr := io.ReadUint16(16)
	if _transactionIdentifierErr != nil {
		return nil, errors.New("Error parsing 'transactionIdentifier' field " + _transactionIdentifierErr.Error())
	}

	// Const Field (protocolIdentifier)
	protocolIdentifier, _protocolIdentifierErr := io.ReadUint16(16)
	if _protocolIdentifierErr != nil {
		return nil, errors.New("Error parsing 'protocolIdentifier' field " + _protocolIdentifierErr.Error())
	}
	if protocolIdentifier != ModbusTcpADU_PROTOCOLIDENTIFIER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(ModbusTcpADU_PROTOCOLIDENTIFIER)) + " but got " + strconv.Itoa(int(protocolIdentifier)))
	}

	// Implicit Field (length) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	_, _lengthErr := io.ReadUint16(16)
	if _lengthErr != nil {
		return nil, errors.New("Error parsing 'length' field " + _lengthErr.Error())
	}

	// Simple Field (unitIdentifier)
	unitIdentifier, _unitIdentifierErr := io.ReadUint8(8)
	if _unitIdentifierErr != nil {
		return nil, errors.New("Error parsing 'unitIdentifier' field " + _unitIdentifierErr.Error())
	}

	// Simple Field (pdu)
	_pduMessage, _err := ModbusPDUParse(io, response)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'pdu'. " + _err.Error())
	}
	var pdu IModbusPDU
	pdu, _pduOk := _pduMessage.(IModbusPDU)
	if !_pduOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_pduMessage).Name() + " to IModbusPDU")
	}

	// Create the instance
	return NewModbusTcpADU(transactionIdentifier, unitIdentifier, pdu), nil
}

func (m ModbusTcpADU) Serialize(io spi.WriteBuffer) {

	// Simple Field (transactionIdentifier)
	transactionIdentifier := uint16(m.transactionIdentifier)
	io.WriteUint16(16, (transactionIdentifier))

	// Const Field (protocolIdentifier)
	io.WriteUint16(16, 0x0000)

	// Implicit Field (length) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	length := uint16(uint16(m.pdu.LengthInBytes()) + uint16(uint16(1)))
	io.WriteUint16(16, (length))

	// Simple Field (unitIdentifier)
	unitIdentifier := uint8(m.unitIdentifier)
	io.WriteUint8(8, (unitIdentifier))

	// Simple Field (pdu)
	pdu := CastIModbusPDU(m.pdu)
	pdu.Serialize(io)

}
