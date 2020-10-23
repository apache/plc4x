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
type BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple struct {
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple interface {
    IBACnetConfirmedServiceRequest
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple) ServiceChoice() uint8 {
    return 0x1F
}

func (m BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple() BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple{}
}

func CastIBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple(structType interface{}) IBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple {
        if iBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple, ok := typ.(IBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple); ok {
            return iBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple(structType interface{}) BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple {
        if sBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple, ok := typ.(BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple); ok {
            return sBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple
        }
        if sBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple, ok := typ.(*BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple); ok {
            return *sBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple
        }
        return BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestConfirmedCOVNotificationMultipleParse(io *utils.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple(), nil
}

func (m BACnetConfirmedServiceRequestConfirmedCOVNotificationMultiple) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
