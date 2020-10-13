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
)

// The data-structure of this message
type ModbusPDUWriteMultipleHoldingRegistersRequest struct {
	startingAddress uint16
	quantity        uint16
	value           []int8
	ModbusPDU
}

// The corresponding interface
type IModbusPDUWriteMultipleHoldingRegistersRequest interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUWriteMultipleHoldingRegistersRequest) ErrorFlag() bool {
	return false
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) FunctionFlag() uint8 {
	return 0x10
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) Response() bool {
	return false
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) initialize() spi.Message {
	return m
}

func NewModbusPDUWriteMultipleHoldingRegistersRequest(startingAddress uint16, quantity uint16, value []int8) ModbusPDUInitializer {
	return &ModbusPDUWriteMultipleHoldingRegistersRequest{startingAddress: startingAddress, quantity: quantity, value: value}
}

func CastIModbusPDUWriteMultipleHoldingRegistersRequest(structType interface{}) IModbusPDUWriteMultipleHoldingRegistersRequest {
	castFunc := func(typ interface{}) IModbusPDUWriteMultipleHoldingRegistersRequest {
		if iModbusPDUWriteMultipleHoldingRegistersRequest, ok := typ.(IModbusPDUWriteMultipleHoldingRegistersRequest); ok {
			return iModbusPDUWriteMultipleHoldingRegistersRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDUWriteMultipleHoldingRegistersRequest(structType interface{}) ModbusPDUWriteMultipleHoldingRegistersRequest {
	castFunc := func(typ interface{}) ModbusPDUWriteMultipleHoldingRegistersRequest {
		if sModbusPDUWriteMultipleHoldingRegistersRequest, ok := typ.(ModbusPDUWriteMultipleHoldingRegistersRequest); ok {
			return sModbusPDUWriteMultipleHoldingRegistersRequest
		}
		return ModbusPDUWriteMultipleHoldingRegistersRequest{}
	}
	return castFunc(structType)
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Simple field (startingAddress)
	lengthInBits += 16

	// Simple field (quantity)
	lengthInBits += 16

	// Implicit Field (byteCount)
	lengthInBits += 8

	// Array field
	if len(m.value) > 0 {
		lengthInBits += 8 * uint16(len(m.value))
	}

	return lengthInBits
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUWriteMultipleHoldingRegistersRequestParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Simple Field (startingAddress)
	startingAddress, _startingAddressErr := io.ReadUint16(16)
	if _startingAddressErr != nil {
		return nil, errors.New("Error parsing 'startingAddress' field " + _startingAddressErr.Error())
	}

	// Simple Field (quantity)
	quantity, _quantityErr := io.ReadUint16(16)
	if _quantityErr != nil {
		return nil, errors.New("Error parsing 'quantity' field " + _quantityErr.Error())
	}

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	byteCount, _byteCountErr := io.ReadUint8(8)
	if _byteCountErr != nil {
		return nil, errors.New("Error parsing 'byteCount' field " + _byteCountErr.Error())
	}

	// Array field (value)
	var value []int8
	// Count array
	{
		value := make([]int8, byteCount)
		for curItem := uint16(0); curItem < uint16(byteCount); curItem++ {

			_valueVal, _err := io.ReadInt8(8)
			if _err != nil {
				return nil, errors.New("Error parsing 'value' field " + _err.Error())
			}
			value = append(value, _valueVal)
		}
	}

	// Create the instance
	return NewModbusPDUWriteMultipleHoldingRegistersRequest(startingAddress, quantity, value), nil
}

func (m ModbusPDUWriteMultipleHoldingRegistersRequest) Serialize(io spi.WriteBuffer) {

	// Simple Field (startingAddress)
	startingAddress := uint16(m.startingAddress)
	io.WriteUint16(16, (startingAddress))

	// Simple Field (quantity)
	quantity := uint16(m.quantity)
	io.WriteUint16(16, (quantity))

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	byteCount := uint8(uint8(len(m.value)))
	io.WriteUint8(8, (byteCount))

	// Array Field (value)
	if m.value != nil {
		for _, _element := range m.value {
			io.WriteInt8(8, _element)
		}
	}
}
