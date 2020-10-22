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
    ObjectType uint16
    ObjectInstanceNumber uint32
    PropertyIdentifierLength uint8
    PropertyIdentifier []int8
    Value IBACnetTag
    BACnetServiceAck
}

// The corresponding interface
type IBACnetServiceAckReadProperty interface {
    IBACnetServiceAck
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetServiceAckReadProperty) ServiceChoice() uint8 {
    return 0x0C
}

func (m BACnetServiceAckReadProperty) initialize() spi.Message {
    return m
}

func NewBACnetServiceAckReadProperty(objectType uint16, objectInstanceNumber uint32, propertyIdentifierLength uint8, propertyIdentifier []int8, value IBACnetTag) BACnetServiceAckInitializer {
    return &BACnetServiceAckReadProperty{ObjectType: objectType, ObjectInstanceNumber: objectInstanceNumber, PropertyIdentifierLength: propertyIdentifierLength, PropertyIdentifier: propertyIdentifier, Value: value}
}

func CastIBACnetServiceAckReadProperty(structType interface{}) IBACnetServiceAckReadProperty {
    castFunc := func(typ interface{}) IBACnetServiceAckReadProperty {
        if iBACnetServiceAckReadProperty, ok := typ.(IBACnetServiceAckReadProperty); ok {
            return iBACnetServiceAckReadProperty
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetServiceAckReadProperty(structType interface{}) BACnetServiceAckReadProperty {
    castFunc := func(typ interface{}) BACnetServiceAckReadProperty {
        if sBACnetServiceAckReadProperty, ok := typ.(BACnetServiceAckReadProperty); ok {
            return sBACnetServiceAckReadProperty
        }
        if sBACnetServiceAckReadProperty, ok := typ.(*BACnetServiceAckReadProperty); ok {
            return *sBACnetServiceAckReadProperty
        }
        return BACnetServiceAckReadProperty{}
    }
    return castFunc(structType)
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
    if len(m.PropertyIdentifier) > 0 {
        lengthInBits += 8 * uint16(len(m.PropertyIdentifier))
    }

    // Const Field (openingTag)
    lengthInBits += 8

    // Simple field (value)
    lengthInBits += m.Value.LengthInBits()

    // Const Field (closingTag)
    lengthInBits += 8

    return lengthInBits
}

func (m BACnetServiceAckReadProperty) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetServiceAckReadPropertyParse(io *spi.ReadBuffer) (BACnetServiceAckInitializer, error) {

    // Const Field (objectIdentifierHeader)
    objectIdentifierHeader, _objectIdentifierHeaderErr := io.ReadUint8(8)
    if _objectIdentifierHeaderErr != nil {
        return nil, errors.New("Error parsing 'objectIdentifierHeader' field " + _objectIdentifierHeaderErr.Error())
    }
    if objectIdentifierHeader != BACnetServiceAckReadProperty_OBJECTIDENTIFIERHEADER {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetServiceAckReadProperty_OBJECTIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(objectIdentifierHeader)))
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
    if propertyIdentifierHeader != BACnetServiceAckReadProperty_PROPERTYIDENTIFIERHEADER {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetServiceAckReadProperty_PROPERTYIDENTIFIERHEADER)) + " but got " + strconv.Itoa(int(propertyIdentifierHeader)))
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

    // Const Field (openingTag)
    openingTag, _openingTagErr := io.ReadUint8(8)
    if _openingTagErr != nil {
        return nil, errors.New("Error parsing 'openingTag' field " + _openingTagErr.Error())
    }
    if openingTag != BACnetServiceAckReadProperty_OPENINGTAG {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetServiceAckReadProperty_OPENINGTAG)) + " but got " + strconv.Itoa(int(openingTag)))
    }

    // Simple Field (value)
    _valueMessage, _err := BACnetTagParse(io)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'value'. " + _err.Error())
    }
    var value IBACnetTag
    value, _valueOk := _valueMessage.(IBACnetTag)
    if !_valueOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_valueMessage).Name() + " to IBACnetTag")
    }

    // Const Field (closingTag)
    closingTag, _closingTagErr := io.ReadUint8(8)
    if _closingTagErr != nil {
        return nil, errors.New("Error parsing 'closingTag' field " + _closingTagErr.Error())
    }
    if closingTag != BACnetServiceAckReadProperty_CLOSINGTAG {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetServiceAckReadProperty_CLOSINGTAG)) + " but got " + strconv.Itoa(int(closingTag)))
    }

    // Create the instance
    return NewBACnetServiceAckReadProperty(objectType, objectInstanceNumber, propertyIdentifierLength, propertyIdentifier, value), nil
}

func (m BACnetServiceAckReadProperty) Serialize(io spi.WriteBuffer) error {
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

    // Const Field (openingTag)
    _openingTagErr := io.WriteUint8(8, 0x3E)
    if _openingTagErr != nil {
        return errors.New("Error serializing 'openingTag' field " + _openingTagErr.Error())
    }

    // Simple Field (value)
    value := CastIBACnetTag(m.Value)
    _valueErr := value.Serialize(io)
    if _valueErr != nil {
        return errors.New("Error serializing 'value' field " + _valueErr.Error())
    }

    // Const Field (closingTag)
    _closingTagErr := io.WriteUint8(8, 0x3F)
    if _closingTagErr != nil {
        return errors.New("Error serializing 'closingTag' field " + _closingTagErr.Error())
    }

        return nil
    }
    return BACnetServiceAckSerialize(io, m.BACnetServiceAck, CastIBACnetServiceAck(m), ser)
}
