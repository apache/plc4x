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
type APDUError struct {
	originalInvokeId uint8
	error            IBACnetError
	APDU
}

// The corresponding interface
type IAPDUError interface {
	IAPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m APDUError) ApduType() uint8 {
	return 0x5
}

func (m APDUError) initialize() spi.Message {
	return m
}

func NewAPDUError(originalInvokeId uint8, error IBACnetError) APDUInitializer {
	return &APDUError{originalInvokeId: originalInvokeId, error: error}
}

func CastIAPDUError(structType interface{}) IAPDUError {
	castFunc := func(typ interface{}) IAPDUError {
		if iAPDUError, ok := typ.(IAPDUError); ok {
			return iAPDUError
		}
		return nil
	}
	return castFunc(structType)
}

func CastAPDUError(structType interface{}) APDUError {
	castFunc := func(typ interface{}) APDUError {
		if sAPDUError, ok := typ.(APDUError); ok {
			return sAPDUError
		}
		return APDUError{}
	}
	return castFunc(structType)
}

func (m APDUError) LengthInBits() uint16 {
	var lengthInBits uint16 = m.APDU.LengthInBits()

	// Reserved Field (reserved)
	lengthInBits += 4

	// Simple field (originalInvokeId)
	lengthInBits += 8

	// Simple field (error)
	lengthInBits += m.error.LengthInBits()

	return lengthInBits
}

func (m APDUError) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func APDUErrorParse(io spi.ReadBuffer) (APDUInitializer, error) {

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint8(4)
		if _err != nil {
			return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
		}
		if reserved != uint8(0x00) {
			log.WithFields(log.Fields{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (originalInvokeId)
	originalInvokeId, _originalInvokeIdErr := io.ReadUint8(8)
	if _originalInvokeIdErr != nil {
		return nil, errors.New("Error parsing 'originalInvokeId' field " + _originalInvokeIdErr.Error())
	}

	// Simple Field (error)
	_errorMessage, _err := BACnetErrorParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'error'. " + _err.Error())
	}
	var error IBACnetError
	error, _errorOk := _errorMessage.(IBACnetError)
	if !_errorOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_errorMessage).Name() + " to IBACnetError")
	}

	// Create the instance
	return NewAPDUError(originalInvokeId, error), nil
}

func (m APDUError) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Reserved Field (reserved)
		io.WriteUint8(4, uint8(0x00))

		// Simple Field (originalInvokeId)
		originalInvokeId := uint8(m.originalInvokeId)
		io.WriteUint8(8, (originalInvokeId))

		// Simple Field (error)
		error := CastIBACnetError(m.error)
		error.Serialize(io)

	}
	APDUSerialize(io, m.APDU, CastIAPDU(m), ser)
}
