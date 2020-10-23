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
type BACnetErrorRemovedReadPropertyConditional struct {
    BACnetError
}

// The corresponding interface
type IBACnetErrorRemovedReadPropertyConditional interface {
    IBACnetError
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetErrorRemovedReadPropertyConditional) ServiceChoice() uint8 {
    return 0x0D
}

func (m BACnetErrorRemovedReadPropertyConditional) initialize() spi.Message {
    return m
}

func NewBACnetErrorRemovedReadPropertyConditional() BACnetErrorInitializer {
    return &BACnetErrorRemovedReadPropertyConditional{}
}

func CastIBACnetErrorRemovedReadPropertyConditional(structType interface{}) IBACnetErrorRemovedReadPropertyConditional {
    castFunc := func(typ interface{}) IBACnetErrorRemovedReadPropertyConditional {
        if iBACnetErrorRemovedReadPropertyConditional, ok := typ.(IBACnetErrorRemovedReadPropertyConditional); ok {
            return iBACnetErrorRemovedReadPropertyConditional
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetErrorRemovedReadPropertyConditional(structType interface{}) BACnetErrorRemovedReadPropertyConditional {
    castFunc := func(typ interface{}) BACnetErrorRemovedReadPropertyConditional {
        if sBACnetErrorRemovedReadPropertyConditional, ok := typ.(BACnetErrorRemovedReadPropertyConditional); ok {
            return sBACnetErrorRemovedReadPropertyConditional
        }
        if sBACnetErrorRemovedReadPropertyConditional, ok := typ.(*BACnetErrorRemovedReadPropertyConditional); ok {
            return *sBACnetErrorRemovedReadPropertyConditional
        }
        return BACnetErrorRemovedReadPropertyConditional{}
    }
    return castFunc(structType)
}

func (m BACnetErrorRemovedReadPropertyConditional) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetError.LengthInBits()

    return lengthInBits
}

func (m BACnetErrorRemovedReadPropertyConditional) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetErrorRemovedReadPropertyConditionalParse(io *utils.ReadBuffer) (BACnetErrorInitializer, error) {

    // Create the instance
    return NewBACnetErrorRemovedReadPropertyConditional(), nil
}

func (m BACnetErrorRemovedReadPropertyConditional) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetErrorSerialize(io, m.BACnetError, CastIBACnetError(m), ser)
}
