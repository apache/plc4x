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
const BACnetServiceAckReadProperty_OBJECTIDENTIFIERHEADER uint8 = 0x0C
const BACnetServiceAckReadProperty_PROPERTYIDENTIFIERHEADER uint8 = 0x03
const BACnetServiceAckReadProperty_OPENINGTAG uint8 = 0x3E
const BACnetServiceAckReadProperty_CLOSINGTAG uint8 = 0x3F

// The data-structure of this message
type BACnetServiceAckReadProperty struct {
	objectType               uint16
	objectInstanceNumber     uint32
	propertyIdentifierLength uint8
	propertyIdentifier       []int8
	value                    BACnetTag
	BACnetServiceAck
}

// The corresponding interface
type IBACnetServiceAckReadProperty interface {
	IBACnetServiceAck
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetServiceAckReadProperty) ServiceChoice() uint8 {
	return 0x0C
}

func (m BACnetServiceAckReadProperty) initialize() spi.Message {
	return m
}

func NewBACnetServiceAckReadProperty(objectType uint16, objectInstanceNumber uint32, propertyIdentifierLength uint8, propertyIdentifier []int8, value BACnetTag) BACnetServiceAckInitializer {
	return &BACnetServiceAckReadProperty{objectType: objectType, objectInstanceNumber: objectInstanceNumber, propertyIdentifierLength: propertyIdentifierLength, propertyIdentifier: propertyIdentifier, value: value}
}

func (m BACnetServiceAckReadProperty) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetServiceAck.LengthInBits()

	// Const Field (objectIdentifierHeader)
	lengthInBits += 8

	// Simple field (objectType)
	lengthInBits += 10

	// Simple field (objectInstanceNumber)
	lengthInBits += 22

	// Const Field (propertyIdentifierHeader)
	lengthInBits += 5

	// Simple field (propertyIdentifierLength)
	lengthInBits += 3

	// Array field
	if len(m.propertyIdentifier) > 0 {
		lengthInBits += 8 * uint16(len(m.propertyIdentifier))
	}

	// Const Field (openingTag)
	lengthInBits += 8

	// Simple field (value)
	lengthInBits += m.value.LengthInBits()

	// Const Field (closingTag)
	lengthInBits += 8

	return lengthInBits
}

func (m BACnetServiceAckReadProperty) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetServiceAckReadPropertyParse(io spi.ReadBuffer) (BACnetServiceAckInitializer, error) {

	// Const Field (objectIdentifierHeader)
	var objectIdentifierHeader uint8 = io.ReadUint8(8)
	if objectIdentifierHeader != BACnetServiceAckReadProperty_OBJECTIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetServiceAckReadProperty_OBJECTIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(objectIdentifierHeader)))
	}

	// Simple Field (objectType)
	var objectType uint16 = io.ReadUint16(10)

	// Simple Field (objectInstanceNumber)
	var objectInstanceNumber uint32 = io.ReadUint32(22)

	// Const Field (propertyIdentifierHeader)
	var propertyIdentifierHeader uint8 = io.ReadUint8(5)
	if propertyIdentifierHeader != BACnetServiceAckReadProperty_PROPERTYIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetServiceAckReadProperty_PROPERTYIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(propertyIdentifierHeader)))
	}

	// Simple Field (propertyIdentifierLength)
	var propertyIdentifierLength uint8 = io.ReadUint8(3)

	// Array field (propertyIdentifier)
	var propertyIdentifier []int8
	// Count array
	{
		propertyIdentifier := make([]int8, propertyIdentifierLength)
		for curItem := uint16(0); curItem < uint16(propertyIdentifierLength); curItem++ {

			propertyIdentifier = append(propertyIdentifier, io.ReadInt8(8))
		}
	}

	// Const Field (openingTag)
	var openingTag uint8 = io.ReadUint8(8)
	if openingTag != BACnetServiceAckReadProperty_OPENINGTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetServiceAckReadProperty_OPENINGTAG)) + " but got " + strconv.Itoa(int(openingTag)))
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
	if closingTag != BACnetServiceAckReadProperty_CLOSINGTAG {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetServiceAckReadProperty_CLOSINGTAG)) + " but got " + strconv.Itoa(int(closingTag)))
	}

	// Create the instance
	return NewBACnetServiceAckReadProperty(objectType, objectInstanceNumber, propertyIdentifierLength, propertyIdentifier, value), nil
}

func (m BACnetServiceAckReadProperty) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IBACnetServiceAckReadProperty); ok {

			// Const Field (objectIdentifierHeader)
			io.WriteUint8(8, 0x0C)

			// Simple Field (objectType)
			var objectType uint16 = m.objectType
			io.WriteUint16(10, (objectType))

			// Simple Field (objectInstanceNumber)
			var objectInstanceNumber uint32 = m.objectInstanceNumber
			io.WriteUint32(22, (objectInstanceNumber))

			// Const Field (propertyIdentifierHeader)
			io.WriteUint8(5, 0x03)

			// Simple Field (propertyIdentifierLength)
			var propertyIdentifierLength uint8 = m.propertyIdentifierLength
			io.WriteUint8(3, (propertyIdentifierLength))

			// Array Field (propertyIdentifier)
			if m.propertyIdentifier != nil {
				for _, _element := range m.propertyIdentifier {
					io.WriteInt8(8, _element)
				}
			}

			// Const Field (openingTag)
			io.WriteUint8(8, 0x3E)

			// Simple Field (value)
			var value BACnetTag = m.value
			value.Serialize(io)

			// Const Field (closingTag)
			io.WriteUint8(8, 0x3F)
		}
	}
	serializeFunc(m)
}
