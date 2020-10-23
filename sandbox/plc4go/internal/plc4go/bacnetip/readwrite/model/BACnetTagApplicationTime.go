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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type BACnetTagApplicationTime struct {
    BACnetTag
}

// The corresponding interface
type IBACnetTagApplicationTime interface {
    IBACnetTag
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetTagApplicationTime) ContextSpecificTag() uint8 {
    return 0
}

func (m BACnetTagApplicationTime) initialize(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) spi.Message {
    m.TypeOrTagNumber = typeOrTagNumber
    m.LengthValueType = lengthValueType
    m.ExtTagNumber = extTagNumber
    m.ExtLength = extLength
    return m
}

func NewBACnetTagApplicationTime() BACnetTagInitializer {
    return &BACnetTagApplicationTime{}
}

func CastIBACnetTagApplicationTime(structType interface{}) IBACnetTagApplicationTime {
    castFunc := func(typ interface{}) IBACnetTagApplicationTime {
        if iBACnetTagApplicationTime, ok := typ.(IBACnetTagApplicationTime); ok {
            return iBACnetTagApplicationTime
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetTagApplicationTime(structType interface{}) BACnetTagApplicationTime {
    castFunc := func(typ interface{}) BACnetTagApplicationTime {
        if sBACnetTagApplicationTime, ok := typ.(BACnetTagApplicationTime); ok {
            return sBACnetTagApplicationTime
        }
        if sBACnetTagApplicationTime, ok := typ.(*BACnetTagApplicationTime); ok {
            return *sBACnetTagApplicationTime
        }
        return BACnetTagApplicationTime{}
    }
    return castFunc(structType)
}

func (m BACnetTagApplicationTime) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetTag.LengthInBits()

    return lengthInBits
}

func (m BACnetTagApplicationTime) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetTagApplicationTimeParse(io *utils.ReadBuffer) (BACnetTagInitializer, error) {

    // Create the instance
    return NewBACnetTagApplicationTime(), nil
}

func (m BACnetTagApplicationTime) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetTagSerialize(io, m.BACnetTag, CastIBACnetTag(m), ser)
}
