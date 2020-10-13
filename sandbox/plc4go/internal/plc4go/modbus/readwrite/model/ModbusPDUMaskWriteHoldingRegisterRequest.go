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
type ModbusPDUMaskWriteHoldingRegisterRequest struct {
	referenceAddress uint16
	andMask          uint16
	orMask           uint16
	ModbusPDU
}

// The corresponding interface
type IModbusPDUMaskWriteHoldingRegisterRequest interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUMaskWriteHoldingRegisterRequest) ErrorFlag() bool {
	return false
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) FunctionFlag() uint8 {
	return 0x16
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) Response() bool {
	return false
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) initialize() spi.Message {
	return m
}

func NewModbusPDUMaskWriteHoldingRegisterRequest(referenceAddress uint16, andMask uint16, orMask uint16) ModbusPDUInitializer {
	return &ModbusPDUMaskWriteHoldingRegisterRequest{referenceAddress: referenceAddress, andMask: andMask, orMask: orMask}
}

func CastIModbusPDUMaskWriteHoldingRegisterRequest(structType interface{}) IModbusPDUMaskWriteHoldingRegisterRequest {
	castFunc := func(typ interface{}) IModbusPDUMaskWriteHoldingRegisterRequest {
		if iModbusPDUMaskWriteHoldingRegisterRequest, ok := typ.(IModbusPDUMaskWriteHoldingRegisterRequest); ok {
			return iModbusPDUMaskWriteHoldingRegisterRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDUMaskWriteHoldingRegisterRequest(structType interface{}) ModbusPDUMaskWriteHoldingRegisterRequest {
	castFunc := func(typ interface{}) ModbusPDUMaskWriteHoldingRegisterRequest {
		if sModbusPDUMaskWriteHoldingRegisterRequest, ok := typ.(ModbusPDUMaskWriteHoldingRegisterRequest); ok {
			return sModbusPDUMaskWriteHoldingRegisterRequest
		}
		return ModbusPDUMaskWriteHoldingRegisterRequest{}
	}
	return castFunc(structType)
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Simple field (referenceAddress)
	lengthInBits += 16

	// Simple field (andMask)
	lengthInBits += 16

	// Simple field (orMask)
	lengthInBits += 16

	return lengthInBits
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUMaskWriteHoldingRegisterRequestParse(io *spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Simple Field (referenceAddress)
	referenceAddress, _referenceAddressErr := io.ReadUint16(16)
	if _referenceAddressErr != nil {
		return nil, errors.New("Error parsing 'referenceAddress' field " + _referenceAddressErr.Error())
	}

	// Simple Field (andMask)
	andMask, _andMaskErr := io.ReadUint16(16)
	if _andMaskErr != nil {
		return nil, errors.New("Error parsing 'andMask' field " + _andMaskErr.Error())
	}

	// Simple Field (orMask)
	orMask, _orMaskErr := io.ReadUint16(16)
	if _orMaskErr != nil {
		return nil, errors.New("Error parsing 'orMask' field " + _orMaskErr.Error())
	}

	// Create the instance
	return NewModbusPDUMaskWriteHoldingRegisterRequest(referenceAddress, andMask, orMask), nil
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (referenceAddress)
		referenceAddress := uint16(m.referenceAddress)
		_referenceAddressErr := io.WriteUint16(16, (referenceAddress))
		if _referenceAddressErr != nil {
			return errors.New("Error serializing 'referenceAddress' field " + _referenceAddressErr.Error())
		}

		// Simple Field (andMask)
		andMask := uint16(m.andMask)
		_andMaskErr := io.WriteUint16(16, (andMask))
		if _andMaskErr != nil {
			return errors.New("Error serializing 'andMask' field " + _andMaskErr.Error())
		}

		// Simple Field (orMask)
		orMask := uint16(m.orMask)
		_orMaskErr := io.WriteUint16(16, (orMask))
		if _orMaskErr != nil {
			return errors.New("Error serializing 'orMask' field " + _orMaskErr.Error())
		}

		return nil
	}
	return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
