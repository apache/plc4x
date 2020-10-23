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
type BACnetServiceAckRemovedAuthenticate struct {
    BACnetServiceAck
}

// The corresponding interface
type IBACnetServiceAckRemovedAuthenticate interface {
    IBACnetServiceAck
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetServiceAckRemovedAuthenticate) ServiceChoice() uint8 {
    return 0x18
}

func (m BACnetServiceAckRemovedAuthenticate) initialize() spi.Message {
    return m
}

func NewBACnetServiceAckRemovedAuthenticate() BACnetServiceAckInitializer {
    return &BACnetServiceAckRemovedAuthenticate{}
}

func CastIBACnetServiceAckRemovedAuthenticate(structType interface{}) IBACnetServiceAckRemovedAuthenticate {
    castFunc := func(typ interface{}) IBACnetServiceAckRemovedAuthenticate {
        if iBACnetServiceAckRemovedAuthenticate, ok := typ.(IBACnetServiceAckRemovedAuthenticate); ok {
            return iBACnetServiceAckRemovedAuthenticate
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetServiceAckRemovedAuthenticate(structType interface{}) BACnetServiceAckRemovedAuthenticate {
    castFunc := func(typ interface{}) BACnetServiceAckRemovedAuthenticate {
        if sBACnetServiceAckRemovedAuthenticate, ok := typ.(BACnetServiceAckRemovedAuthenticate); ok {
            return sBACnetServiceAckRemovedAuthenticate
        }
        if sBACnetServiceAckRemovedAuthenticate, ok := typ.(*BACnetServiceAckRemovedAuthenticate); ok {
            return *sBACnetServiceAckRemovedAuthenticate
        }
        return BACnetServiceAckRemovedAuthenticate{}
    }
    return castFunc(structType)
}

func (m BACnetServiceAckRemovedAuthenticate) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetServiceAck.LengthInBits()

    return lengthInBits
}

func (m BACnetServiceAckRemovedAuthenticate) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetServiceAckRemovedAuthenticateParse(io *utils.ReadBuffer) (BACnetServiceAckInitializer, error) {

    // Create the instance
    return NewBACnetServiceAckRemovedAuthenticate(), nil
}

func (m BACnetServiceAckRemovedAuthenticate) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetServiceAckSerialize(io, m.BACnetServiceAck, CastIBACnetServiceAck(m), ser)
}
