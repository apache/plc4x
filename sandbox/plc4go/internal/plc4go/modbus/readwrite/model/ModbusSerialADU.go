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
	log "github.com/sirupsen/logrus"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"reflect"
)

// The data-structure of this message
type ModbusSerialADU struct {
	transactionId uint16
	length        uint16
	address       uint8
	pdu           IModbusPDU
}

// The corresponding interface
type IModbusSerialADU interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewModbusSerialADU(transactionId uint16, length uint16, address uint8, pdu IModbusPDU) spi.Message {
	return &ModbusSerialADU{transactionId: transactionId, length: length, address: address, pdu: pdu}
}

func CastIModbusSerialADU(structType interface{}) IModbusSerialADU {
	castFunc := func(typ interface{}) IModbusSerialADU {
		if iModbusSerialADU, ok := typ.(IModbusSerialADU); ok {
			return iModbusSerialADU
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusSerialADU(structType interface{}) ModbusSerialADU {
	castFunc := func(typ interface{}) ModbusSerialADU {
		if sModbusSerialADU, ok := typ.(ModbusSerialADU); ok {
			return sModbusSerialADU
		}
		return ModbusSerialADU{}
	}
	return castFunc(structType)
}

func (m ModbusSerialADU) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Simple field (transactionId)
	lengthInBits += 16

	// Reserved Field (reserved)
	lengthInBits += 16

	// Simple field (length)
	lengthInBits += 16

	// Simple field (address)
	lengthInBits += 8

	// Simple field (pdu)
	lengthInBits += m.pdu.LengthInBits()

	return lengthInBits
}

func (m ModbusSerialADU) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusSerialADUParse(io spi.ReadBuffer, response bool) (spi.Message, error) {

	// Simple Field (transactionId)
	transactionId, _transactionIdErr := io.ReadUint16(16)
	if _transactionIdErr != nil {
		return nil, errors.New("Error parsing 'transactionId' field " + _transactionIdErr.Error())
	}

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint16(16)
		if _err != nil {
			return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
		}
		if reserved != uint16(0x0000) {
			log.WithFields(log.Fields{
				"expected value": uint16(0x0000),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (length)
	length, _lengthErr := io.ReadUint16(16)
	if _lengthErr != nil {
		return nil, errors.New("Error parsing 'length' field " + _lengthErr.Error())
	}

	// Simple Field (address)
	address, _addressErr := io.ReadUint8(8)
	if _addressErr != nil {
		return nil, errors.New("Error parsing 'address' field " + _addressErr.Error())
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
	return NewModbusSerialADU(transactionId, length, address, pdu), nil
}

func (m ModbusSerialADU) Serialize(io spi.WriteBuffer) {

	// Simple Field (transactionId)
	transactionId := uint16(m.transactionId)
	io.WriteUint16(16, (transactionId))

	// Reserved Field (reserved)
	io.WriteUint16(16, uint16(0x0000))

	// Simple Field (length)
	length := uint16(m.length)
	io.WriteUint16(16, (length))

	// Simple Field (address)
	address := uint8(m.address)
	io.WriteUint8(8, (address))

	// Simple Field (pdu)
	pdu := IModbusPDU(m.pdu)
	pdu.Serialize(io)
}
