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
type BACnetConfirmedServiceRequestReadRange struct {
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestReadRange interface {
    IBACnetConfirmedServiceRequest
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestReadRange) ServiceChoice() uint8 {
    return 0x1A
}

func (m BACnetConfirmedServiceRequestReadRange) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestReadRange() BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestReadRange{}
}

func CastIBACnetConfirmedServiceRequestReadRange(structType interface{}) IBACnetConfirmedServiceRequestReadRange {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestReadRange {
        if iBACnetConfirmedServiceRequestReadRange, ok := typ.(IBACnetConfirmedServiceRequestReadRange); ok {
            return iBACnetConfirmedServiceRequestReadRange
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestReadRange(structType interface{}) BACnetConfirmedServiceRequestReadRange {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequestReadRange {
        if sBACnetConfirmedServiceRequestReadRange, ok := typ.(BACnetConfirmedServiceRequestReadRange); ok {
            return sBACnetConfirmedServiceRequestReadRange
        }
        if sBACnetConfirmedServiceRequestReadRange, ok := typ.(*BACnetConfirmedServiceRequestReadRange); ok {
            return *sBACnetConfirmedServiceRequestReadRange
        }
        return BACnetConfirmedServiceRequestReadRange{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestReadRange) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceRequestReadRange) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestReadRangeParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceRequestReadRange(), nil
}

func (m BACnetConfirmedServiceRequestReadRange) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
