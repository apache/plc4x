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
type BACnetConfirmedServiceACKReadRange struct {
    BACnetConfirmedServiceACK
}

// The corresponding interface
type IBACnetConfirmedServiceACKReadRange interface {
    IBACnetConfirmedServiceACK
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceACKReadRange) ServiceChoice() uint8 {
    return 0x1A
}

func (m BACnetConfirmedServiceACKReadRange) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceACKReadRange() BACnetConfirmedServiceACKInitializer {
    return &BACnetConfirmedServiceACKReadRange{}
}

func CastIBACnetConfirmedServiceACKReadRange(structType interface{}) IBACnetConfirmedServiceACKReadRange {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceACKReadRange {
        if iBACnetConfirmedServiceACKReadRange, ok := typ.(IBACnetConfirmedServiceACKReadRange); ok {
            return iBACnetConfirmedServiceACKReadRange
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceACKReadRange(structType interface{}) BACnetConfirmedServiceACKReadRange {
    castFunc := func(typ interface{}) BACnetConfirmedServiceACKReadRange {
        if sBACnetConfirmedServiceACKReadRange, ok := typ.(BACnetConfirmedServiceACKReadRange); ok {
            return sBACnetConfirmedServiceACKReadRange
        }
        if sBACnetConfirmedServiceACKReadRange, ok := typ.(*BACnetConfirmedServiceACKReadRange); ok {
            return *sBACnetConfirmedServiceACKReadRange
        }
        return BACnetConfirmedServiceACKReadRange{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceACKReadRange) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceACK.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceACKReadRange) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceACKReadRangeParse(io *spi.ReadBuffer) (BACnetConfirmedServiceACKInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceACKReadRange(), nil
}

func (m BACnetConfirmedServiceACKReadRange) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceACKSerialize(io, m.BACnetConfirmedServiceACK, CastIBACnetConfirmedServiceACK(m), ser)
}
