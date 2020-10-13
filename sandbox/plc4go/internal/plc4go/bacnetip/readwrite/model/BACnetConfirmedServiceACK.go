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
type BACnetConfirmedServiceACK struct {
}

// The corresponding interface
type IBACnetConfirmedServiceACK interface {
	spi.Message
	ServiceChoice() uint8
	Serialize(io spi.WriteBuffer) error
}

type BACnetConfirmedServiceACKInitializer interface {
	initialize() spi.Message
}

func BACnetConfirmedServiceACKServiceChoice(m IBACnetConfirmedServiceACK) uint8 {
	return m.ServiceChoice()
}

func CastIBACnetConfirmedServiceACK(structType interface{}) IBACnetConfirmedServiceACK {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceACK {
		if iBACnetConfirmedServiceACK, ok := typ.(IBACnetConfirmedServiceACK); ok {
			return iBACnetConfirmedServiceACK
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceACK(structType interface{}) BACnetConfirmedServiceACK {
	castFunc := func(typ interface{}) BACnetConfirmedServiceACK {
		if sBACnetConfirmedServiceACK, ok := typ.(BACnetConfirmedServiceACK); ok {
			return sBACnetConfirmedServiceACK
		}
		return BACnetConfirmedServiceACK{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceACK) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (serviceChoice)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m BACnetConfirmedServiceACK) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceACKParse(io *spi.ReadBuffer) (spi.Message, error) {

	// Discriminator Field (serviceChoice) (Used as input to a switch field)
	serviceChoice, _serviceChoiceErr := io.ReadUint8(8)
	if _serviceChoiceErr != nil {
		return nil, errors.New("Error parsing 'serviceChoice' field " + _serviceChoiceErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer BACnetConfirmedServiceACKInitializer
	var typeSwitchError error
	switch {
	case serviceChoice == 0x03:
		initializer, typeSwitchError = BACnetConfirmedServiceACKGetAlarmSummaryParse(io)
	case serviceChoice == 0x04:
		initializer, typeSwitchError = BACnetConfirmedServiceACKGetEnrollmentSummaryParse(io)
	case serviceChoice == 0x1D:
		initializer, typeSwitchError = BACnetConfirmedServiceACKGetEventInformationParse(io)
	case serviceChoice == 0x06:
		initializer, typeSwitchError = BACnetConfirmedServiceACKAtomicReadFileParse(io)
	case serviceChoice == 0x07:
		initializer, typeSwitchError = BACnetConfirmedServiceACKAtomicWriteFileParse(io)
	case serviceChoice == 0x0A:
		initializer, typeSwitchError = BACnetConfirmedServiceACKCreateObjectParse(io)
	case serviceChoice == 0x0C:
		initializer, typeSwitchError = BACnetConfirmedServiceACKReadPropertyParse(io)
	case serviceChoice == 0x0E:
		initializer, typeSwitchError = BACnetConfirmedServiceACKReadPropertyMultipleParse(io)
	case serviceChoice == 0x1A:
		initializer, typeSwitchError = BACnetConfirmedServiceACKReadRangeParse(io)
	case serviceChoice == 0x12:
		initializer, typeSwitchError = BACnetConfirmedServiceACKConfirmedPrivateTransferParse(io)
	case serviceChoice == 0x15:
		initializer, typeSwitchError = BACnetConfirmedServiceACKVTOpenParse(io)
	case serviceChoice == 0x17:
		initializer, typeSwitchError = BACnetConfirmedServiceACKVTDataParse(io)
	case serviceChoice == 0x18:
		initializer, typeSwitchError = BACnetConfirmedServiceACKRemovedAuthenticateParse(io)
	case serviceChoice == 0x0D:
		initializer, typeSwitchError = BACnetConfirmedServiceACKRemovedReadPropertyConditionalParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func BACnetConfirmedServiceACKSerialize(io spi.WriteBuffer, m BACnetConfirmedServiceACK, i IBACnetConfirmedServiceACK, childSerialize func() error) error {

	// Discriminator Field (serviceChoice) (Used as input to a switch field)
	serviceChoice := uint8(i.ServiceChoice())
	_serviceChoiceErr := io.WriteUint8(8, (serviceChoice))
	if _serviceChoiceErr != nil {
		return errors.New("Error serializing 'serviceChoice' field " + _serviceChoiceErr.Error())
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	_typeSwitchErr := childSerialize()
	if _typeSwitchErr != nil {
		return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
	}

	return nil
}
