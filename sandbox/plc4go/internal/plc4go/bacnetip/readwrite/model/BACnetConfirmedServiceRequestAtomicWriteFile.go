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
type BACnetConfirmedServiceRequestAtomicWriteFile struct {
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestAtomicWriteFile interface {
    IBACnetConfirmedServiceRequest
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestAtomicWriteFile) ServiceChoice() uint8 {
    return 0x07
}

func (m BACnetConfirmedServiceRequestAtomicWriteFile) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestAtomicWriteFile() BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestAtomicWriteFile{}
}

func CastIBACnetConfirmedServiceRequestAtomicWriteFile(structType interface{}) IBACnetConfirmedServiceRequestAtomicWriteFile {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestAtomicWriteFile {
        if iBACnetConfirmedServiceRequestAtomicWriteFile, ok := typ.(IBACnetConfirmedServiceRequestAtomicWriteFile); ok {
            return iBACnetConfirmedServiceRequestAtomicWriteFile
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestAtomicWriteFile(structType interface{}) BACnetConfirmedServiceRequestAtomicWriteFile {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequestAtomicWriteFile {
        if sBACnetConfirmedServiceRequestAtomicWriteFile, ok := typ.(BACnetConfirmedServiceRequestAtomicWriteFile); ok {
            return sBACnetConfirmedServiceRequestAtomicWriteFile
        }
        if sBACnetConfirmedServiceRequestAtomicWriteFile, ok := typ.(*BACnetConfirmedServiceRequestAtomicWriteFile); ok {
            return *sBACnetConfirmedServiceRequestAtomicWriteFile
        }
        return BACnetConfirmedServiceRequestAtomicWriteFile{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestAtomicWriteFile) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceRequestAtomicWriteFile) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestAtomicWriteFileParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceRequestAtomicWriteFile(), nil
}

func (m BACnetConfirmedServiceRequestAtomicWriteFile) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
