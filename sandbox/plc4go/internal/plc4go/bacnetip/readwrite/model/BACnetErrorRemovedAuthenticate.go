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
type BACnetErrorRemovedAuthenticate struct {
    BACnetError
}

// The corresponding interface
type IBACnetErrorRemovedAuthenticate interface {
    IBACnetError
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetErrorRemovedAuthenticate) ServiceChoice() uint8 {
    return 0x18
}

func (m BACnetErrorRemovedAuthenticate) initialize() spi.Message {
    return m
}

func NewBACnetErrorRemovedAuthenticate() BACnetErrorInitializer {
    return &BACnetErrorRemovedAuthenticate{}
}

func CastIBACnetErrorRemovedAuthenticate(structType interface{}) IBACnetErrorRemovedAuthenticate {
    castFunc := func(typ interface{}) IBACnetErrorRemovedAuthenticate {
        if iBACnetErrorRemovedAuthenticate, ok := typ.(IBACnetErrorRemovedAuthenticate); ok {
            return iBACnetErrorRemovedAuthenticate
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetErrorRemovedAuthenticate(structType interface{}) BACnetErrorRemovedAuthenticate {
    castFunc := func(typ interface{}) BACnetErrorRemovedAuthenticate {
        if sBACnetErrorRemovedAuthenticate, ok := typ.(BACnetErrorRemovedAuthenticate); ok {
            return sBACnetErrorRemovedAuthenticate
        }
        if sBACnetErrorRemovedAuthenticate, ok := typ.(*BACnetErrorRemovedAuthenticate); ok {
            return *sBACnetErrorRemovedAuthenticate
        }
        return BACnetErrorRemovedAuthenticate{}
    }
    return castFunc(structType)
}

func (m BACnetErrorRemovedAuthenticate) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetError.LengthInBits()

    return lengthInBits
}

func (m BACnetErrorRemovedAuthenticate) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetErrorRemovedAuthenticateParse(io *utils.ReadBuffer) (BACnetErrorInitializer, error) {

    // Create the instance
    return NewBACnetErrorRemovedAuthenticate(), nil
}

func (m BACnetErrorRemovedAuthenticate) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetErrorSerialize(io, m.BACnetError, CastIBACnetError(m), ser)
}
