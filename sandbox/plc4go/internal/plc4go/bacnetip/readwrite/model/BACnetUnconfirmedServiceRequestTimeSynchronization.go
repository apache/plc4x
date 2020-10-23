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
type BACnetUnconfirmedServiceRequestTimeSynchronization struct {
    BACnetUnconfirmedServiceRequest
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestTimeSynchronization interface {
    IBACnetUnconfirmedServiceRequest
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetUnconfirmedServiceRequestTimeSynchronization) ServiceChoice() uint8 {
    return 0x06
}

func (m BACnetUnconfirmedServiceRequestTimeSynchronization) initialize() spi.Message {
    return m
}

func NewBACnetUnconfirmedServiceRequestTimeSynchronization() BACnetUnconfirmedServiceRequestInitializer {
    return &BACnetUnconfirmedServiceRequestTimeSynchronization{}
}

func CastIBACnetUnconfirmedServiceRequestTimeSynchronization(structType interface{}) IBACnetUnconfirmedServiceRequestTimeSynchronization {
    castFunc := func(typ interface{}) IBACnetUnconfirmedServiceRequestTimeSynchronization {
        if iBACnetUnconfirmedServiceRequestTimeSynchronization, ok := typ.(IBACnetUnconfirmedServiceRequestTimeSynchronization); ok {
            return iBACnetUnconfirmedServiceRequestTimeSynchronization
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetUnconfirmedServiceRequestTimeSynchronization(structType interface{}) BACnetUnconfirmedServiceRequestTimeSynchronization {
    castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequestTimeSynchronization {
        if sBACnetUnconfirmedServiceRequestTimeSynchronization, ok := typ.(BACnetUnconfirmedServiceRequestTimeSynchronization); ok {
            return sBACnetUnconfirmedServiceRequestTimeSynchronization
        }
        if sBACnetUnconfirmedServiceRequestTimeSynchronization, ok := typ.(*BACnetUnconfirmedServiceRequestTimeSynchronization); ok {
            return *sBACnetUnconfirmedServiceRequestTimeSynchronization
        }
        return BACnetUnconfirmedServiceRequestTimeSynchronization{}
    }
    return castFunc(structType)
}

func (m BACnetUnconfirmedServiceRequestTimeSynchronization) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetUnconfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetUnconfirmedServiceRequestTimeSynchronization) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestTimeSynchronizationParse(io *utils.ReadBuffer) (BACnetUnconfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetUnconfirmedServiceRequestTimeSynchronization(), nil
}

func (m BACnetUnconfirmedServiceRequestTimeSynchronization) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetUnconfirmedServiceRequestSerialize(io, m.BACnetUnconfirmedServiceRequest, CastIBACnetUnconfirmedServiceRequest(m), ser)
}
