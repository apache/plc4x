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
type ModbusPDUError struct {
	ExceptionCode uint8
	ModbusPDU
}

// The corresponding interface
type IModbusPDUError interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUError) ErrorFlag() bool {
	return true
}

func (m ModbusPDUError) FunctionFlag() uint8 {
	return 0
}

func (m ModbusPDUError) Response() bool {
	return false
}

func (m ModbusPDUError) initialize() spi.Message {
	return m
}

func NewModbusPDUError(exceptionCode uint8) ModbusPDUInitializer {
	return &ModbusPDUError{ExceptionCode: exceptionCode}
}

func CastIModbusPDUError(structType interface{}) IModbusPDUError {
	castFunc := func(typ interface{}) IModbusPDUError {
		if iModbusPDUError, ok := typ.(IModbusPDUError); ok {
			return iModbusPDUError
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDUError(structType interface{}) ModbusPDUError {
	castFunc := func(typ interface{}) ModbusPDUError {
		if sModbusPDUError, ok := typ.(ModbusPDUError); ok {
			return sModbusPDUError
		}
		return ModbusPDUError{}
	}
	return castFunc(structType)
}

func (m ModbusPDUError) LengthInBits() uint16 {
	var lengthInBits = m.ModbusPDU.LengthInBits()

	// Simple field (exceptionCode)
	lengthInBits += 8

	return lengthInBits
}

func (m ModbusPDUError) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUErrorParse(io *spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Simple Field (exceptionCode)
	exceptionCode, _exceptionCodeErr := io.ReadUint8(8)
	if _exceptionCodeErr != nil {
		return nil, errors.New("Error parsing 'exceptionCode' field " + _exceptionCodeErr.Error())
	}

	// Create the instance
	return NewModbusPDUError(exceptionCode), nil
}

func (m ModbusPDUError) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (exceptionCode)
		exceptionCode := uint8(m.ExceptionCode)
		_exceptionCodeErr := io.WriteUint8(8, exceptionCode)
		if _exceptionCodeErr != nil {
			return errors.New("Error serializing 'exceptionCode' field " + _exceptionCodeErr.Error())
		}

		return nil
	}
	return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
