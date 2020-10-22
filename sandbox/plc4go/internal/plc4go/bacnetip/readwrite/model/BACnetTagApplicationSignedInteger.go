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
)

// The data-structure of this message
type BACnetTagApplicationSignedInteger struct {
    Data []int8
    BACnetTag
}

// The corresponding interface
type IBACnetTagApplicationSignedInteger interface {
    IBACnetTag
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetTagApplicationSignedInteger) ContextSpecificTag() uint8 {
    return 0
}

func (m BACnetTagApplicationSignedInteger) initialize(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) spi.Message {
    m.TypeOrTagNumber = typeOrTagNumber
    m.LengthValueType = lengthValueType
    m.ExtTagNumber = extTagNumber
    m.ExtLength = extLength
    return m
}

func NewBACnetTagApplicationSignedInteger(data []int8) BACnetTagInitializer {
    return &BACnetTagApplicationSignedInteger{Data: data}
}

func CastIBACnetTagApplicationSignedInteger(structType interface{}) IBACnetTagApplicationSignedInteger {
    castFunc := func(typ interface{}) IBACnetTagApplicationSignedInteger {
        if iBACnetTagApplicationSignedInteger, ok := typ.(IBACnetTagApplicationSignedInteger); ok {
            return iBACnetTagApplicationSignedInteger
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetTagApplicationSignedInteger(structType interface{}) BACnetTagApplicationSignedInteger {
    castFunc := func(typ interface{}) BACnetTagApplicationSignedInteger {
        if sBACnetTagApplicationSignedInteger, ok := typ.(BACnetTagApplicationSignedInteger); ok {
            return sBACnetTagApplicationSignedInteger
        }
        if sBACnetTagApplicationSignedInteger, ok := typ.(*BACnetTagApplicationSignedInteger); ok {
            return *sBACnetTagApplicationSignedInteger
        }
        return BACnetTagApplicationSignedInteger{}
    }
    return castFunc(structType)
}

func (m BACnetTagApplicationSignedInteger) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetTag.LengthInBits()

    // Array field
    if len(m.Data) > 0 {
        lengthInBits += 8 * uint16(len(m.Data))
    }

    return lengthInBits
}

func (m BACnetTagApplicationSignedInteger) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetTagApplicationSignedIntegerParse(io *spi.ReadBuffer, lengthValueType uint8, extLength uint8) (BACnetTagInitializer, error) {

    // Array field (data)
    // Length array
    data := make([]int8, 0)
    _dataLength := spi.InlineIf(bool(bool((lengthValueType) == ((5)))), uint16(extLength), uint16(lengthValueType))
    _dataEndPos := io.GetPos() + uint16(_dataLength)
    for ;io.GetPos() < _dataEndPos; {
        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'data' field " + _err.Error())
        }
        data = append(data, _item)
    }

    // Create the instance
    return NewBACnetTagApplicationSignedInteger(data), nil
}

func (m BACnetTagApplicationSignedInteger) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

    // Array Field (data)
    if m.Data != nil {
        for _, _element := range m.Data {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'data' field " + _elementErr.Error())
            }
        }
    }

        return nil
    }
    return BACnetTagSerialize(io, m.BACnetTag, CastIBACnetTag(m), ser)
}
