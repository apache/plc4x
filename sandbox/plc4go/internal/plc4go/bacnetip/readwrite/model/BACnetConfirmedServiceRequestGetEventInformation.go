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
type BACnetConfirmedServiceRequestGetEventInformation struct {
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestGetEventInformation interface {
    IBACnetConfirmedServiceRequest
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestGetEventInformation) ServiceChoice() uint8 {
    return 0x1D
}

func (m BACnetConfirmedServiceRequestGetEventInformation) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestGetEventInformation() BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestGetEventInformation{}
}

func CastIBACnetConfirmedServiceRequestGetEventInformation(structType interface{}) IBACnetConfirmedServiceRequestGetEventInformation {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestGetEventInformation {
        if iBACnetConfirmedServiceRequestGetEventInformation, ok := typ.(IBACnetConfirmedServiceRequestGetEventInformation); ok {
            return iBACnetConfirmedServiceRequestGetEventInformation
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestGetEventInformation(structType interface{}) BACnetConfirmedServiceRequestGetEventInformation {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequestGetEventInformation {
        if sBACnetConfirmedServiceRequestGetEventInformation, ok := typ.(BACnetConfirmedServiceRequestGetEventInformation); ok {
            return sBACnetConfirmedServiceRequestGetEventInformation
        }
        if sBACnetConfirmedServiceRequestGetEventInformation, ok := typ.(*BACnetConfirmedServiceRequestGetEventInformation); ok {
            return *sBACnetConfirmedServiceRequestGetEventInformation
        }
        return BACnetConfirmedServiceRequestGetEventInformation{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestGetEventInformation) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceRequestGetEventInformation) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestGetEventInformationParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceRequestGetEventInformation(), nil
}

func (m BACnetConfirmedServiceRequestGetEventInformation) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
