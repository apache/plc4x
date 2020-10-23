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
type BVLCDeleteForeignDeviceTableEntry struct {
    BVLC
}

// The corresponding interface
type IBVLCDeleteForeignDeviceTableEntry interface {
    IBVLC
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BVLCDeleteForeignDeviceTableEntry) BvlcFunction() uint8 {
    return 0x08
}

func (m BVLCDeleteForeignDeviceTableEntry) initialize() spi.Message {
    return m
}

func NewBVLCDeleteForeignDeviceTableEntry() BVLCInitializer {
    return &BVLCDeleteForeignDeviceTableEntry{}
}

func CastIBVLCDeleteForeignDeviceTableEntry(structType interface{}) IBVLCDeleteForeignDeviceTableEntry {
    castFunc := func(typ interface{}) IBVLCDeleteForeignDeviceTableEntry {
        if iBVLCDeleteForeignDeviceTableEntry, ok := typ.(IBVLCDeleteForeignDeviceTableEntry); ok {
            return iBVLCDeleteForeignDeviceTableEntry
        }
        return nil
    }
    return castFunc(structType)
}

func CastBVLCDeleteForeignDeviceTableEntry(structType interface{}) BVLCDeleteForeignDeviceTableEntry {
    castFunc := func(typ interface{}) BVLCDeleteForeignDeviceTableEntry {
        if sBVLCDeleteForeignDeviceTableEntry, ok := typ.(BVLCDeleteForeignDeviceTableEntry); ok {
            return sBVLCDeleteForeignDeviceTableEntry
        }
        if sBVLCDeleteForeignDeviceTableEntry, ok := typ.(*BVLCDeleteForeignDeviceTableEntry); ok {
            return *sBVLCDeleteForeignDeviceTableEntry
        }
        return BVLCDeleteForeignDeviceTableEntry{}
    }
    return castFunc(structType)
}

func (m BVLCDeleteForeignDeviceTableEntry) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BVLC.LengthInBits()

    return lengthInBits
}

func (m BVLCDeleteForeignDeviceTableEntry) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BVLCDeleteForeignDeviceTableEntryParse(io *utils.ReadBuffer) (BVLCInitializer, error) {

    // Create the instance
    return NewBVLCDeleteForeignDeviceTableEntry(), nil
}

func (m BVLCDeleteForeignDeviceTableEntry) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BVLCSerialize(io, m.BVLC, CastIBVLC(m), ser)
}
