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
type BVLCRegisterForeignDevice struct {
    BVLC
}

// The corresponding interface
type IBVLCRegisterForeignDevice interface {
    IBVLC
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BVLCRegisterForeignDevice) BvlcFunction() uint8 {
    return 0x05
}

func (m BVLCRegisterForeignDevice) initialize() spi.Message {
    return m
}

func NewBVLCRegisterForeignDevice() BVLCInitializer {
    return &BVLCRegisterForeignDevice{}
}

func CastIBVLCRegisterForeignDevice(structType interface{}) IBVLCRegisterForeignDevice {
    castFunc := func(typ interface{}) IBVLCRegisterForeignDevice {
        if iBVLCRegisterForeignDevice, ok := typ.(IBVLCRegisterForeignDevice); ok {
            return iBVLCRegisterForeignDevice
        }
        return nil
    }
    return castFunc(structType)
}

func CastBVLCRegisterForeignDevice(structType interface{}) BVLCRegisterForeignDevice {
    castFunc := func(typ interface{}) BVLCRegisterForeignDevice {
        if sBVLCRegisterForeignDevice, ok := typ.(BVLCRegisterForeignDevice); ok {
            return sBVLCRegisterForeignDevice
        }
        if sBVLCRegisterForeignDevice, ok := typ.(*BVLCRegisterForeignDevice); ok {
            return *sBVLCRegisterForeignDevice
        }
        return BVLCRegisterForeignDevice{}
    }
    return castFunc(structType)
}

func (m BVLCRegisterForeignDevice) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BVLC.LengthInBits()

    return lengthInBits
}

func (m BVLCRegisterForeignDevice) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BVLCRegisterForeignDeviceParse(io *utils.ReadBuffer) (BVLCInitializer, error) {

    // Create the instance
    return NewBVLCRegisterForeignDevice(), nil
}

func (m BVLCRegisterForeignDevice) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BVLCSerialize(io, m.BVLC, CastIBVLC(m), ser)
}
