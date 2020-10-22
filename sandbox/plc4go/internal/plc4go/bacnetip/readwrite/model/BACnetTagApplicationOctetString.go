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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
)

// The data-structure of this message
type BACnetTagApplicationOctetString struct {
    BACnetTag
}

// The corresponding interface
type IBACnetTagApplicationOctetString interface {
    IBACnetTag
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetTagApplicationOctetString) ContextSpecificTag() uint8 {
    return 0
}

func (m BACnetTagApplicationOctetString) initialize(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) spi.Message {
    m.TypeOrTagNumber = typeOrTagNumber
    m.LengthValueType = lengthValueType
    m.ExtTagNumber = extTagNumber
    m.ExtLength = extLength
    return m
}

func NewBACnetTagApplicationOctetString() BACnetTagInitializer {
    return &BACnetTagApplicationOctetString{}
}

func CastIBACnetTagApplicationOctetString(structType interface{}) IBACnetTagApplicationOctetString {
    castFunc := func(typ interface{}) IBACnetTagApplicationOctetString {
        if iBACnetTagApplicationOctetString, ok := typ.(IBACnetTagApplicationOctetString); ok {
            return iBACnetTagApplicationOctetString
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetTagApplicationOctetString(structType interface{}) BACnetTagApplicationOctetString {
    castFunc := func(typ interface{}) BACnetTagApplicationOctetString {
        if sBACnetTagApplicationOctetString, ok := typ.(BACnetTagApplicationOctetString); ok {
            return sBACnetTagApplicationOctetString
        }
        if sBACnetTagApplicationOctetString, ok := typ.(*BACnetTagApplicationOctetString); ok {
            return *sBACnetTagApplicationOctetString
        }
        return BACnetTagApplicationOctetString{}
    }
    return castFunc(structType)
}

func (m BACnetTagApplicationOctetString) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetTag.LengthInBits()

    return lengthInBits
}

func (m BACnetTagApplicationOctetString) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetTagApplicationOctetStringParse(io *spi.ReadBuffer) (BACnetTagInitializer, error) {

    // Create the instance
    return NewBACnetTagApplicationOctetString(), nil
}

func (m BACnetTagApplicationOctetString) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetTagSerialize(io, m.BACnetTag, CastIBACnetTag(m), ser)
}
