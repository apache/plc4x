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
type BACnetTag struct {
	typeOrTagNumber uint8
	lengthValueType uint8
	extTagNumber    *uint8
	extLength       *uint8
}

// The corresponding interface
type IBACnetTag interface {
	spi.Message
	ContextSpecificTag() uint8
	Serialize(io spi.WriteBuffer)
}

type BACnetTagInitializer interface {
	initialize(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) spi.Message
}

func BACnetTagContextSpecificTag(m IBACnetTag) uint8 {
	return m.ContextSpecificTag()
}

func CastIBACnetTag(structType interface{}) IBACnetTag {
	castFunc := func(typ interface{}) IBACnetTag {
		if iBACnetTag, ok := typ.(IBACnetTag); ok {
			return iBACnetTag
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetTag(structType interface{}) BACnetTag {
	castFunc := func(typ interface{}) BACnetTag {
		if sBACnetTag, ok := typ.(BACnetTag); ok {
			return sBACnetTag
		}
		return BACnetTag{}
	}
	return castFunc(structType)
}

func (m BACnetTag) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Simple field (typeOrTagNumber)
	lengthInBits += 4

	// Discriminator Field (contextSpecificTag)
	lengthInBits += 1

	// Simple field (lengthValueType)
	lengthInBits += 3

	// Optional Field (extTagNumber)
	if m.extTagNumber != nil {
		lengthInBits += 8
	}

	// Optional Field (extLength)
	if m.extLength != nil {
		lengthInBits += 8
	}

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m BACnetTag) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetTagParse(io spi.ReadBuffer) (spi.Message, error) {

	// Simple Field (typeOrTagNumber)
	typeOrTagNumber, _typeOrTagNumberErr := io.ReadUint8(4)
	if _typeOrTagNumberErr != nil {
		return nil, errors.New("Error parsing 'typeOrTagNumber' field " + _typeOrTagNumberErr.Error())
	}

	// Discriminator Field (contextSpecificTag) (Used as input to a switch field)
	contextSpecificTag, _contextSpecificTagErr := io.ReadUint8(1)
	if _contextSpecificTagErr != nil {
		return nil, errors.New("Error parsing 'contextSpecificTag' field " + _contextSpecificTagErr.Error())
	}

	// Simple Field (lengthValueType)
	lengthValueType, _lengthValueTypeErr := io.ReadUint8(3)
	if _lengthValueTypeErr != nil {
		return nil, errors.New("Error parsing 'lengthValueType' field " + _lengthValueTypeErr.Error())
	}

	// Optional Field (extTagNumber) (Can be skipped, if a given expression evaluates to false)
	var extTagNumber *uint8 = nil
	if bool((typeOrTagNumber) == (15)) {
		_val, _err := io.ReadUint8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'extTagNumber' field " + _err.Error())
		}

		extTagNumber = &_val
	}

	// Optional Field (extLength) (Can be skipped, if a given expression evaluates to false)
	var extLength *uint8 = nil
	if bool((lengthValueType) == (5)) {
		_val, _err := io.ReadUint8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'extLength' field " + _err.Error())
		}

		extLength = &_val
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer BACnetTagInitializer
	var typeSwitchError error
	switch {
	case contextSpecificTag == 0 && typeOrTagNumber == 0x0:
		initializer, typeSwitchError = BACnetTagApplicationNullParse(io)
	case contextSpecificTag == 0 && typeOrTagNumber == 0x1:
		initializer, typeSwitchError = BACnetTagApplicationBooleanParse(io)
	case contextSpecificTag == 0 && typeOrTagNumber == 0x2:
		initializer, typeSwitchError = BACnetTagApplicationUnsignedIntegerParse(io, lengthValueType, *extLength)
	case contextSpecificTag == 0 && typeOrTagNumber == 0x3:
		initializer, typeSwitchError = BACnetTagApplicationSignedIntegerParse(io, lengthValueType, *extLength)
	case contextSpecificTag == 0 && typeOrTagNumber == 0x4:
		initializer, typeSwitchError = BACnetTagApplicationRealParse(io, lengthValueType, *extLength)
	case contextSpecificTag == 0 && typeOrTagNumber == 0x5:
		initializer, typeSwitchError = BACnetTagApplicationDoubleParse(io, lengthValueType, *extLength)
	case contextSpecificTag == 0 && typeOrTagNumber == 0x6:
		initializer, typeSwitchError = BACnetTagApplicationOctetStringParse(io)
	case contextSpecificTag == 0 && typeOrTagNumber == 0x7:
		initializer, typeSwitchError = BACnetTagApplicationCharacterStringParse(io)
	case contextSpecificTag == 0 && typeOrTagNumber == 0x8:
		initializer, typeSwitchError = BACnetTagApplicationBitStringParse(io, lengthValueType, *extLength)
	case contextSpecificTag == 0 && typeOrTagNumber == 0x9:
		initializer, typeSwitchError = BACnetTagApplicationEnumeratedParse(io, lengthValueType, *extLength)
	case contextSpecificTag == 0 && typeOrTagNumber == 0xA:
		initializer, typeSwitchError = BACnetTagApplicationDateParse(io)
	case contextSpecificTag == 0 && typeOrTagNumber == 0xB:
		initializer, typeSwitchError = BACnetTagApplicationTimeParse(io)
	case contextSpecificTag == 0 && typeOrTagNumber == 0xC:
		initializer, typeSwitchError = BACnetTagApplicationObjectIdentifierParse(io)
	case contextSpecificTag == 1:
		initializer, typeSwitchError = BACnetTagContextParse(io, typeOrTagNumber, *extTagNumber, lengthValueType, *extLength)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(typeOrTagNumber, lengthValueType, extTagNumber, extLength), nil
}

func (m BACnetTag) Serialize(io spi.WriteBuffer) {
	iBACnetTag := CastIBACnetTag(m)

	// Simple Field (typeOrTagNumber)
	typeOrTagNumber := uint8(m.typeOrTagNumber)
	io.WriteUint8(4, (typeOrTagNumber))

	// Discriminator Field (contextSpecificTag) (Used as input to a switch field)
	contextSpecificTag := uint8(BACnetTagContextSpecificTag(iBACnetTag))
	io.WriteUint8(1, (contextSpecificTag))

	// Simple Field (lengthValueType)
	lengthValueType := uint8(m.lengthValueType)
	io.WriteUint8(3, (lengthValueType))

	// Optional Field (extTagNumber) (Can be skipped, if the value is null)
	var extTagNumber *uint8 = nil
	if m.extTagNumber != nil {
		extTagNumber = m.extTagNumber
		io.WriteUint8(8, *(extTagNumber))
	}

	// Optional Field (extLength) (Can be skipped, if the value is null)
	var extLength *uint8 = nil
	if m.extLength != nil {
		extLength = m.extLength
		io.WriteUint8(8, *(extLength))
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	iBACnetTag.Serialize(io)
}
