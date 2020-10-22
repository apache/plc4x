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
type BACnetConfirmedServiceRequestVTOpen struct {
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestVTOpen interface {
    IBACnetConfirmedServiceRequest
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestVTOpen) ServiceChoice() uint8 {
    return 0x15
}

func (m BACnetConfirmedServiceRequestVTOpen) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestVTOpen() BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestVTOpen{}
}

func CastIBACnetConfirmedServiceRequestVTOpen(structType interface{}) IBACnetConfirmedServiceRequestVTOpen {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestVTOpen {
        if iBACnetConfirmedServiceRequestVTOpen, ok := typ.(IBACnetConfirmedServiceRequestVTOpen); ok {
            return iBACnetConfirmedServiceRequestVTOpen
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestVTOpen(structType interface{}) BACnetConfirmedServiceRequestVTOpen {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequestVTOpen {
        if sBACnetConfirmedServiceRequestVTOpen, ok := typ.(BACnetConfirmedServiceRequestVTOpen); ok {
            return sBACnetConfirmedServiceRequestVTOpen
        }
        if sBACnetConfirmedServiceRequestVTOpen, ok := typ.(*BACnetConfirmedServiceRequestVTOpen); ok {
            return *sBACnetConfirmedServiceRequestVTOpen
        }
        return BACnetConfirmedServiceRequestVTOpen{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestVTOpen) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceRequestVTOpen) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestVTOpenParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceRequestVTOpen(), nil
}

func (m BACnetConfirmedServiceRequestVTOpen) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
