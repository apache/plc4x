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
type ModbusPDUReadExceptionStatusResponse struct {
	value uint8
	ModbusPDU
}

// The corresponding interface
type IModbusPDUReadExceptionStatusResponse interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUReadExceptionStatusResponse) ErrorFlag() bool {
	return false
}

func (m ModbusPDUReadExceptionStatusResponse) FunctionFlag() uint8 {
	return 0x07
}

func (m ModbusPDUReadExceptionStatusResponse) Response() bool {
	return true
}

func (m ModbusPDUReadExceptionStatusResponse) initialize() spi.Message {
	return m
}

func NewModbusPDUReadExceptionStatusResponse(value uint8) ModbusPDUInitializer {
	return &ModbusPDUReadExceptionStatusResponse{value: value}
}

func CastIModbusPDUReadExceptionStatusResponse(structType interface{}) IModbusPDUReadExceptionStatusResponse {
	castFunc := func(typ interface{}) IModbusPDUReadExceptionStatusResponse {
		if iModbusPDUReadExceptionStatusResponse, ok := typ.(IModbusPDUReadExceptionStatusResponse); ok {
			return iModbusPDUReadExceptionStatusResponse
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDUReadExceptionStatusResponse(structType interface{}) ModbusPDUReadExceptionStatusResponse {
	castFunc := func(typ interface{}) ModbusPDUReadExceptionStatusResponse {
		if sModbusPDUReadExceptionStatusResponse, ok := typ.(ModbusPDUReadExceptionStatusResponse); ok {
			return sModbusPDUReadExceptionStatusResponse
		}
		return ModbusPDUReadExceptionStatusResponse{}
	}
	return castFunc(structType)
}

func (m ModbusPDUReadExceptionStatusResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Simple field (value)
	lengthInBits += 8

	return lengthInBits
}

func (m ModbusPDUReadExceptionStatusResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUReadExceptionStatusResponseParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Simple Field (value)
	value, _valueErr := io.ReadUint8(8)
	if _valueErr != nil {
		return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
	}

	// Create the instance
	return NewModbusPDUReadExceptionStatusResponse(value), nil
}

func (m ModbusPDUReadExceptionStatusResponse) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Simple Field (value)
		value := uint8(m.value)
		io.WriteUint8(8, (value))

	}
	ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
