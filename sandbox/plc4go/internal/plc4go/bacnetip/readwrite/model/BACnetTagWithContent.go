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
	"reflect"
	"strconv"
)

// Constant values.
const BACnetTagWithContent_OPENTAG uint8 = 0x2e
const BACnetTagWithContent_CLOSINGTAG uint8 = 0x2f

// The data-structure of this message
type BACnetTagWithContent struct {
	typeOrTagNumber    uint8
	contextSpecificTag uint8
	lengthValueType    uint8
	extTagNumber       *uint8
	extLength          *uint8
	propertyIdentifier []uint8
	value              BACnetTag
}

// The corresponding interface
type IBACnetTagWithContent interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewBACnetTagWithContent(typeOrTagNumber uint8, contextSpecificTag uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8, propertyIdentifier []uint8, value BACnetTag) spi.Message {
	return &BACnetTagWithContent{typeOrTagNumber: typeOrTagNumber, contextSpecificTag: contextSpecificTag, lengthValueType: lengthValueType, extTagNumber: extTagNumber, extLength: extLength, propertyIdentifier: propertyIdentifier, value: value}
}

func CastIBACnetTagWithContent(structType interface{}) IBACnetTagWithContent {
	castFunc := func(typ interface{}) IBACnetTagWithContent {
		if iBACnetTagWithContent, ok := typ.(IBACnetTagWithContent); ok {
			return iBACnetTagWithContent
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetTagWithContent(structType interface{}) BACnetTagWithContent {
	castFunc := func(typ interface{}) BACnetTagWithContent {
		if sBACnetTagWithContent, ok := typ.(BACnetTagWithContent); ok {
			return sBACnetTagWithContent
		}
		return BACnetTagWithContent{}
	}
	return castFunc(structType)
}

func (m BACnetTagWithContent) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Simple field (typeOrTagNumber)
	lengthInBits += 4

	// Simple field (contextSpecificTag)
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

	// Array field
	if len(m.propertyIdentifier) > 0 {
		lengthInBits += 8 * uint16(len(m.propertyIdentifier))
	}

	// Const Field (openTag)
	lengthInBits += 8

	// Simple field (value)
	lengthInBits += m.value.LengthInBits()

	// Const Field (closingTag)
	lengthInBits += 8

	return lengthInBits
}

func (m BACnetTagWithContent) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetTagWithContentParse(io spi.ReadBuffer) (spi.Message, error) {

	// Simple Field (typeOrTagNumber)
	var typeOrTagNumber uint8 = io.ReadUint8(4)

	// Simple Field (contextSpecificTag)
	var contextSpecificTag uint8 = io.ReadUint8(1)

	// Simple Field (lengthValueType)
	var lengthValueType uint8 = io.ReadUint8(3)

	// Optional Field (extTagNumber) (Can be skipped, if a given expression evaluates to false)
	var extTagNumber *uint8 = nil
	if bool((typeOrTagNumber) == (15)) {
		_val := io.ReadUint8(8)
		extTagNumber = &_val
	}

	// Optional Field (extLength) (Can be skipped, if a given expression evaluates to false)
	var extLength *uint8 = nil
	if bool((lengthValueType) == (5)) {
		_val := io.ReadUint8(8)
		extLength = &_val
	}

	// Array field (propertyIdentifier)
	var propertyIdentifier []uint8
	// Length array
	_propertyIdentifierLength := spi.InlineIf(bool(bool((lengthValueType) == (5))), uint16(*extLength), uint16(lengthValueType))
	_propertyIdentifierEndPos := io.GetPos() + uint16(_propertyIdentifierLength)
	for io.GetPos() < _propertyIdentifierEndPos {
		propertyIdentifier = append(propertyIdentifier, io.ReadUint8(8))
	}

	// Const Field (openTag)
	var openTag uint8 = io.ReadUint8(8)
	if openTag != BACnetTagWithContent_OPENTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetTagWithContent_OPENTAG)) + " but got " + strconv.Itoa(int(openTag)))
	}

	// Simple Field (value)
	_valueMessage, _err := BACnetTagParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'value'. " + _err.Error())
	}
	var value BACnetTag
	value, _valueOk := _valueMessage.(BACnetTag)
	if !_valueOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_valueMessage).Name() + " to BACnetTag")
	}

	// Const Field (closingTag)
	var closingTag uint8 = io.ReadUint8(8)
	if closingTag != BACnetTagWithContent_CLOSINGTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetTagWithContent_CLOSINGTAG)) + " but got " + strconv.Itoa(int(closingTag)))
	}

	// Create the instance
	return NewBACnetTagWithContent(typeOrTagNumber, contextSpecificTag, lengthValueType, extTagNumber, extLength, propertyIdentifier, value), nil
}

func (m BACnetTagWithContent) Serialize(io spi.WriteBuffer) {

	// Simple Field (typeOrTagNumber)
	typeOrTagNumber := uint8(m.typeOrTagNumber)
	io.WriteUint8(4, (typeOrTagNumber))

	// Simple Field (contextSpecificTag)
	contextSpecificTag := uint8(m.contextSpecificTag)
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

	// Array Field (propertyIdentifier)
	if m.propertyIdentifier != nil {
		for _, _element := range m.propertyIdentifier {
			io.WriteUint8(8, _element)
		}
	}

	// Const Field (openTag)
	io.WriteUint8(8, 0x2e)

	// Simple Field (value)
	value := BACnetTag(m.value)
	value.Serialize(io)

	// Const Field (closingTag)
	io.WriteUint8(8, 0x2f)
}
