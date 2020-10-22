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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "strconv"
)

// Constant values.
const BACnetConfirmedServiceRequestReadProperty_OBJECTIDENTIFIERHEADER uint8 = 0x0C
const BACnetConfirmedServiceRequestReadProperty_PROPERTYIDENTIFIERHEADER uint8 = 0x03

// The data-structure of this message
type BACnetConfirmedServiceRequestReadProperty struct {
    ObjectType uint16
    ObjectInstanceNumber uint32
    PropertyIdentifierLength uint8
    PropertyIdentifier []int8
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestReadProperty interface {
    IBACnetConfirmedServiceRequest
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestReadProperty) ServiceChoice() uint8 {
    return 0x0C
}

func (m BACnetConfirmedServiceRequestReadProperty) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestReadProperty(objectType uint16, objectInstanceNumber uint32, propertyIdentifierLength uint8, propertyIdentifier []int8) BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestReadProperty{ObjectType: objectType, ObjectInstanceNumber: objectInstanceNumber, PropertyIdentifierLength: propertyIdentifierLength, PropertyIdentifier: propertyIdentifier}
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
        if sBACnetConfirmedServiceRequestReadProperty, ok := typ.(*BACnetConfirmedServiceRequestReadProperty); ok {
            return *sBACnetConfirmedServiceRequestReadProperty
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
    if len(m.PropertyIdentifier) > 0 {
        lengthInBits += 8 * uint16(len(m.PropertyIdentifier))
    }

    return lengthInBits
}

func (m BACnetConfirmedServiceRequestReadProperty) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestReadPropertyParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

    // Const Field (objectIdentifierHeader)
    objectIdentifierHeader, _objectIdentifierHeaderErr := io.ReadUint8(8)
    if _objectIdentifierHeaderErr != nil {
        return nil, errors.New("Error parsing 'objectIdentifierHeader' field " + _objectIdentifierHeaderErr.Error())
    }
    if objectIdentifierHeader != BACnetConfirmedServiceRequestReadProperty_OBJECTIDENTIFIERHEADER {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestReadProperty_OBJECTIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(objectIdentifierHeader)))
    }

    // Simple Field (objectType)
    objectType, _objectTypeErr := io.ReadUint16(10)
    if _objectTypeErr != nil {
        return nil, errors.New("Error parsing 'objectType' field " + _objectTypeErr.Error())
    }

    // Simple Field (objectInstanceNumber)
    objectInstanceNumber, _objectInstanceNumberErr := io.ReadUint32(22)
    if _objectInstanceNumberErr != nil {
        return nil, errors.New("Error parsing 'objectInstanceNumber' field " + _objectInstanceNumberErr.Error())
    }

    // Const Field (propertyIdentifierHeader)
    propertyIdentifierHeader, _propertyIdentifierHeaderErr := io.ReadUint8(5)
    if _propertyIdentifierHeaderErr != nil {
        return nil, errors.New("Error parsing 'propertyIdentifierHeader' field " + _propertyIdentifierHeaderErr.Error())
    }
    if propertyIdentifierHeader != BACnetConfirmedServiceRequestReadProperty_PROPERTYIDENTIFIERHEADER {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetConfirmedServiceRequestReadProperty_PROPERTYIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(propertyIdentifierHeader)))
    }

    // Simple Field (propertyIdentifierLength)
    propertyIdentifierLength, _propertyIdentifierLengthErr := io.ReadUint8(3)
    if _propertyIdentifierLengthErr != nil {
        return nil, errors.New("Error parsing 'propertyIdentifierLength' field " + _propertyIdentifierLengthErr.Error())
    }

    // Array field (propertyIdentifier)
    // Count array
    propertyIdentifier := make([]int8, propertyIdentifierLength)
    for curItem := uint16(0); curItem < uint16(propertyIdentifierLength); curItem++ {

        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'propertyIdentifier' field " + _err.Error())
        }
        propertyIdentifier[curItem] = _item
    }

    // Create the instance
    return NewBACnetConfirmedServiceRequestReadProperty(objectType, objectInstanceNumber, propertyIdentifierLength, propertyIdentifier), nil
}

func (m BACnetConfirmedServiceRequestReadProperty) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

    // Const Field (objectIdentifierHeader)
    _objectIdentifierHeaderErr := io.WriteUint8(8, 0x0C)
    if _objectIdentifierHeaderErr != nil {
        return errors.New("Error serializing 'objectIdentifierHeader' field " + _objectIdentifierHeaderErr.Error())
    }

    // Simple Field (objectType)
    objectType := uint16(m.ObjectType)
    _objectTypeErr := io.WriteUint16(10, (objectType))
    if _objectTypeErr != nil {
        return errors.New("Error serializing 'objectType' field " + _objectTypeErr.Error())
    }

    // Simple Field (objectInstanceNumber)
    objectInstanceNumber := uint32(m.ObjectInstanceNumber)
    _objectInstanceNumberErr := io.WriteUint32(22, (objectInstanceNumber))
    if _objectInstanceNumberErr != nil {
        return errors.New("Error serializing 'objectInstanceNumber' field " + _objectInstanceNumberErr.Error())
    }

    // Const Field (propertyIdentifierHeader)
    _propertyIdentifierHeaderErr := io.WriteUint8(5, 0x03)
    if _propertyIdentifierHeaderErr != nil {
        return errors.New("Error serializing 'propertyIdentifierHeader' field " + _propertyIdentifierHeaderErr.Error())
    }

    // Simple Field (propertyIdentifierLength)
    propertyIdentifierLength := uint8(m.PropertyIdentifierLength)
    _propertyIdentifierLengthErr := io.WriteUint8(3, (propertyIdentifierLength))
    if _propertyIdentifierLengthErr != nil {
        return errors.New("Error serializing 'propertyIdentifierLength' field " + _propertyIdentifierLengthErr.Error())
    }

    // Array Field (propertyIdentifier)
    if m.PropertyIdentifier != nil {
        for _, _element := range m.PropertyIdentifier {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'propertyIdentifier' field " + _elementErr.Error())
            }
        }
    }

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
