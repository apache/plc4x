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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type ModbusPDUWriteMultipleCoilsRequest struct {
	startingAddress uint16
	quantity        uint16
	value           []int8
	ModbusPDU
}

// The corresponding interface
type IModbusPDUWriteMultipleCoilsRequest interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUWriteMultipleCoilsRequest) ErrorFlag() bool {
	return false
}

func (m ModbusPDUWriteMultipleCoilsRequest) FunctionFlag() uint8 {
	return 0x0F
}

func (m ModbusPDUWriteMultipleCoilsRequest) Response() bool {
	return false
}

func (m ModbusPDUWriteMultipleCoilsRequest) initialize() spi.Message {
	return m
}

func NewModbusPDUWriteMultipleCoilsRequest(startingAddress uint16, quantity uint16, value []int8) ModbusPDUInitializer {
	return &ModbusPDUWriteMultipleCoilsRequest{startingAddress: startingAddress, quantity: quantity, value: value}
}

func CastIModbusPDUWriteMultipleCoilsRequest(structType interface{}) IModbusPDUWriteMultipleCoilsRequest {
	castFunc := func(typ interface{}) IModbusPDUWriteMultipleCoilsRequest {
		if iModbusPDUWriteMultipleCoilsRequest, ok := typ.(IModbusPDUWriteMultipleCoilsRequest); ok {
			return iModbusPDUWriteMultipleCoilsRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDUWriteMultipleCoilsRequest(structType interface{}) ModbusPDUWriteMultipleCoilsRequest {
	castFunc := func(typ interface{}) ModbusPDUWriteMultipleCoilsRequest {
		if sModbusPDUWriteMultipleCoilsRequest, ok := typ.(ModbusPDUWriteMultipleCoilsRequest); ok {
			return sModbusPDUWriteMultipleCoilsRequest
		}
		return ModbusPDUWriteMultipleCoilsRequest{}
	}
	return castFunc(structType)
}

func (m ModbusPDUWriteMultipleCoilsRequest) LengthInBits() uint16 {
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

func (m ModbusPDUWriteMultipleCoilsRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUWriteMultipleCoilsRequestParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Simple Field (startingAddress)
	var startingAddress uint16 = io.ReadUint16(16)

	// Simple Field (quantity)
	var quantity uint16 = io.ReadUint16(16)

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var byteCount uint8 = io.ReadUint8(8)

	// Array field (value)
	var value []int8
	// Count array
	{
		value := make([]int8, byteCount)
		for curItem := uint16(0); curItem < uint16(byteCount); curItem++ {

			value = append(value, io.ReadInt8(8))
		}
	}

	// Create the instance
	return NewModbusPDUWriteMultipleCoilsRequest(startingAddress, quantity, value), nil
}

func (m ModbusPDUWriteMultipleCoilsRequest) Serialize(io spi.WriteBuffer) {

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
