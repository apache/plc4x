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
	"strconv"
)

// Constant values.
const BACnetConfirmedServiceRequestReadProperty_OBJECTIDENTIFIERHEADER uint8 = 0x0C
const BACnetConfirmedServiceRequestReadProperty_PROPERTYIDENTIFIERHEADER uint8 = 0x03

// The data-structure of this message
type BACnetConfirmedServiceRequestReadProperty struct {
	objectType               uint16
	objectInstanceNumber     uint32
	propertyIdentifierLength uint8
	propertyIdentifier       []int8
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestReadProperty interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestReadProperty) ServiceChoice() uint8 {
	return 0x0C
}

func (m BACnetConfirmedServiceRequestReadProperty) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestReadProperty(objectType uint16, objectInstanceNumber uint32, propertyIdentifierLength uint8, propertyIdentifier []int8) BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestReadProperty{objectType: objectType, objectInstanceNumber: objectInstanceNumber, propertyIdentifierLength: propertyIdentifierLength, propertyIdentifier: propertyIdentifier}
}

func CastIBACnetConfirmedServiceRequestReadProperty(structType interface{}) IBACnetConfirmedServiceRequestReadProperty {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestReadProperty {
		if iBACnetConfirmedServiceRequestReadProperty, ok := typ.(IBACnetConfirmedServiceRequestReadProperty); ok {
			return iBACnetConfirmedServiceRequestReadProperty
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestReadProperty(structType interface{}) BACnetConfirmedServiceRequestReadProperty {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestReadProperty {
		if sBACnetConfirmedServiceRequestReadProperty, ok := typ.(BACnetConfirmedServiceRequestReadProperty); ok {
			return sBACnetConfirmedServiceRequestReadProperty
		}
		return BACnetConfirmedServiceRequestReadProperty{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestReadProperty) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

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

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestReadProperty) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestReadPropertyParse(io spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Const Field (objectIdentifierHeader)
	var objectIdentifierHeader uint8 = io.ReadUint8(8)
	if objectIdentifierHeader != BACnetConfirmedServiceRequestReadProperty_OBJECTIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestReadProperty_OBJECTIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(objectIdentifierHeader)))
	}

	// Simple Field (objectType)
	var objectType uint16 = io.ReadUint16(10)

	// Simple Field (objectInstanceNumber)
	var objectInstanceNumber uint32 = io.ReadUint32(22)

	// Const Field (propertyIdentifierHeader)
	var propertyIdentifierHeader uint8 = io.ReadUint8(5)
	if propertyIdentifierHeader != BACnetConfirmedServiceRequestReadProperty_PROPERTYIDENTIFIERHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestReadProperty_PROPERTYIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(propertyIdentifierHeader)))
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

	// Create the instance
	return NewBACnetConfirmedServiceRequestReadProperty(objectType, objectInstanceNumber, propertyIdentifierLength, propertyIdentifier), nil
}

func (m BACnetConfirmedServiceRequestReadProperty) Serialize(io spi.WriteBuffer) {

	// Const Field (objectIdentifierHeader)
	io.WriteUint8(8, 0x0C)

	// Simple Field (objectType)
	objectType := uint16(m.objectType)
	io.WriteUint16(10, (objectType))

	// Simple Field (objectInstanceNumber)
	objectInstanceNumber := uint32(m.objectInstanceNumber)
	io.WriteUint32(22, (objectInstanceNumber))

	// Const Field (propertyIdentifierHeader)
	io.WriteUint8(5, 0x03)

	// Simple Field (propertyIdentifierLength)
	propertyIdentifierLength := uint8(m.propertyIdentifierLength)
	io.WriteUint8(3, (propertyIdentifierLength))

	// Array Field (propertyIdentifier)
	if m.propertyIdentifier != nil {
		for _, _element := range m.propertyIdentifier {
			io.WriteInt8(8, _element)
		}
	}
}
