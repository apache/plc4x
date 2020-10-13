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
type APDU struct {
}

// The corresponding interface
type IAPDU interface {
	spi.Message
	ApduType() uint8
	Serialize(io spi.WriteBuffer)
}

type APDUInitializer interface {
	initialize() spi.Message
}

func APDUApduType(m IAPDU) uint8 {
	return m.ApduType()
}

func CastIAPDU(structType interface{}) IAPDU {
	castFunc := func(typ interface{}) IAPDU {
		if iAPDU, ok := typ.(IAPDU); ok {
			return iAPDU
		}
		return nil
	}
	return castFunc(structType)
}

func CastAPDU(structType interface{}) APDU {
	castFunc := func(typ interface{}) APDU {
		if sAPDU, ok := typ.(APDU); ok {
			return sAPDU
		}
		return APDU{}
	}
	return castFunc(structType)
}

func (m APDU) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (apduType)
	lengthInBits += 4

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m APDU) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func APDUParse(io spi.ReadBuffer, apduLength uint16) (spi.Message, error) {

	// Discriminator Field (apduType) (Used as input to a switch field)
	apduType, _apduTypeErr := io.ReadUint8(4)
	if _apduTypeErr != nil {
		return nil, errors.New("Error parsing 'apduType' field " + _apduTypeErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer APDUInitializer
	var typeSwitchError error
	switch {
	case apduType == 0x0:
		initializer, typeSwitchError = APDUConfirmedRequestParse(io, apduLength)
	case apduType == 0x1:
		initializer, typeSwitchError = APDUUnconfirmedRequestParse(io, apduLength)
	case apduType == 0x2:
		initializer, typeSwitchError = APDUSimpleAckParse(io)
	case apduType == 0x3:
		initializer, typeSwitchError = APDUComplexAckParse(io)
	case apduType == 0x4:
		initializer, typeSwitchError = APDUSegmentAckParse(io)
	case apduType == 0x5:
		initializer, typeSwitchError = APDUErrorParse(io)
	case apduType == 0x6:
		initializer, typeSwitchError = APDURejectParse(io)
	case apduType == 0x7:
		initializer, typeSwitchError = APDUAbortParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func APDUSerialize(io spi.WriteBuffer, m APDU, i IAPDU, childSerialize func()) {

	// Discriminator Field (apduType) (Used as input to a switch field)
	apduType := uint8(i.ApduType())
	io.WriteUint8(4, (apduType))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	childSerialize()

}
