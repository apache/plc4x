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
type BACnetServiceAckAtomicWriteFile struct {
    BACnetServiceAck
}

// The corresponding interface
type IBACnetServiceAckAtomicWriteFile interface {
    IBACnetServiceAck
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetServiceAckAtomicWriteFile) ServiceChoice() uint8 {
    return 0x07
}

func (m BACnetServiceAckAtomicWriteFile) initialize() spi.Message {
    return m
}

func NewBACnetServiceAckAtomicWriteFile() BACnetServiceAckInitializer {
    return &BACnetServiceAckAtomicWriteFile{}
}

func CastIBACnetServiceAckAtomicWriteFile(structType interface{}) IBACnetServiceAckAtomicWriteFile {
    castFunc := func(typ interface{}) IBACnetServiceAckAtomicWriteFile {
        if iBACnetServiceAckAtomicWriteFile, ok := typ.(IBACnetServiceAckAtomicWriteFile); ok {
            return iBACnetServiceAckAtomicWriteFile
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetServiceAckAtomicWriteFile(structType interface{}) BACnetServiceAckAtomicWriteFile {
    castFunc := func(typ interface{}) BACnetServiceAckAtomicWriteFile {
        if sBACnetServiceAckAtomicWriteFile, ok := typ.(BACnetServiceAckAtomicWriteFile); ok {
            return sBACnetServiceAckAtomicWriteFile
        }
        if sBACnetServiceAckAtomicWriteFile, ok := typ.(*BACnetServiceAckAtomicWriteFile); ok {
            return *sBACnetServiceAckAtomicWriteFile
        }
        return BACnetServiceAckAtomicWriteFile{}
    }
    return castFunc(structType)
}

func (m BACnetServiceAckAtomicWriteFile) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetServiceAck.LengthInBits()

    return lengthInBits
}

func (m BACnetServiceAckAtomicWriteFile) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetServiceAckAtomicWriteFileParse(io *utils.ReadBuffer) (BACnetServiceAckInitializer, error) {

    // Create the instance
    return NewBACnetServiceAckAtomicWriteFile(), nil
}

func (m BACnetServiceAckAtomicWriteFile) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetServiceAckSerialize(io, m.BACnetServiceAck, CastIBACnetServiceAck(m), ser)
}
