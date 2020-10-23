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
type BACnetErrorGetEventInformation struct {
    BACnetError
}

// The corresponding interface
type IBACnetErrorGetEventInformation interface {
    IBACnetError
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetErrorGetEventInformation) ServiceChoice() uint8 {
    return 0x1D
}

func (m BACnetErrorGetEventInformation) initialize() spi.Message {
    return m
}

func NewBACnetErrorGetEventInformation() BACnetErrorInitializer {
    return &BACnetErrorGetEventInformation{}
}

func CastIBACnetErrorGetEventInformation(structType interface{}) IBACnetErrorGetEventInformation {
    castFunc := func(typ interface{}) IBACnetErrorGetEventInformation {
        if iBACnetErrorGetEventInformation, ok := typ.(IBACnetErrorGetEventInformation); ok {
            return iBACnetErrorGetEventInformation
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetErrorGetEventInformation(structType interface{}) BACnetErrorGetEventInformation {
    castFunc := func(typ interface{}) BACnetErrorGetEventInformation {
        if sBACnetErrorGetEventInformation, ok := typ.(BACnetErrorGetEventInformation); ok {
            return sBACnetErrorGetEventInformation
        }
        if sBACnetErrorGetEventInformation, ok := typ.(*BACnetErrorGetEventInformation); ok {
            return *sBACnetErrorGetEventInformation
        }
        return BACnetErrorGetEventInformation{}
    }
    return castFunc(structType)
}

func (m BACnetErrorGetEventInformation) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetError.LengthInBits()

    return lengthInBits
}

func (m BACnetErrorGetEventInformation) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetErrorGetEventInformationParse(io *utils.ReadBuffer) (BACnetErrorInitializer, error) {

    // Create the instance
    return NewBACnetErrorGetEventInformation(), nil
}

func (m BACnetErrorGetEventInformation) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetErrorSerialize(io, m.BACnetError, CastIBACnetError(m), ser)
}
