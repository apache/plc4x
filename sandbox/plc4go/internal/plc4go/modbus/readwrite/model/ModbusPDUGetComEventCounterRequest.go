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
type ModbusPDUGetComEventCounterRequest struct {
    ModbusPDU
}

// The corresponding interface
type IModbusPDUGetComEventCounterRequest interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUGetComEventCounterRequest) ErrorFlag() bool {
    return false
}

func (m ModbusPDUGetComEventCounterRequest) FunctionFlag() uint8 {
    return 0x0B
}

func (m ModbusPDUGetComEventCounterRequest) Response() bool {
    return false
}

func (m ModbusPDUGetComEventCounterRequest) initialize() spi.Message {
    return m
}

func NewModbusPDUGetComEventCounterRequest() ModbusPDUInitializer {
    return &ModbusPDUGetComEventCounterRequest{}
}

func CastIModbusPDUGetComEventCounterRequest(structType interface{}) IModbusPDUGetComEventCounterRequest {
    castFunc := func(typ interface{}) IModbusPDUGetComEventCounterRequest {
        if iModbusPDUGetComEventCounterRequest, ok := typ.(IModbusPDUGetComEventCounterRequest); ok {
            return iModbusPDUGetComEventCounterRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUGetComEventCounterRequest(structType interface{}) ModbusPDUGetComEventCounterRequest {
    castFunc := func(typ interface{}) ModbusPDUGetComEventCounterRequest {
        if sModbusPDUGetComEventCounterRequest, ok := typ.(ModbusPDUGetComEventCounterRequest); ok {
            return sModbusPDUGetComEventCounterRequest
        }
        if sModbusPDUGetComEventCounterRequest, ok := typ.(*ModbusPDUGetComEventCounterRequest); ok {
            return *sModbusPDUGetComEventCounterRequest
        }
        return ModbusPDUGetComEventCounterRequest{}
    }
    return castFunc(structType)
}

func (m ModbusPDUGetComEventCounterRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    return lengthInBits
}

func (m ModbusPDUGetComEventCounterRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUGetComEventCounterRequestParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Create the instance
    return NewModbusPDUGetComEventCounterRequest(), nil
}

func (m ModbusPDUGetComEventCounterRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
