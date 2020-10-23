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
type ModbusPDUGetComEventLogRequest struct {
    ModbusPDU
}

// The corresponding interface
type IModbusPDUGetComEventLogRequest interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUGetComEventLogRequest) ErrorFlag() bool {
    return false
}

func (m ModbusPDUGetComEventLogRequest) FunctionFlag() uint8 {
    return 0x0C
}

func (m ModbusPDUGetComEventLogRequest) Response() bool {
    return false
}

func (m ModbusPDUGetComEventLogRequest) initialize() spi.Message {
    return m
}

func NewModbusPDUGetComEventLogRequest() ModbusPDUInitializer {
    return &ModbusPDUGetComEventLogRequest{}
}

func CastIModbusPDUGetComEventLogRequest(structType interface{}) IModbusPDUGetComEventLogRequest {
    castFunc := func(typ interface{}) IModbusPDUGetComEventLogRequest {
        if iModbusPDUGetComEventLogRequest, ok := typ.(IModbusPDUGetComEventLogRequest); ok {
            return iModbusPDUGetComEventLogRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUGetComEventLogRequest(structType interface{}) ModbusPDUGetComEventLogRequest {
    castFunc := func(typ interface{}) ModbusPDUGetComEventLogRequest {
        if sModbusPDUGetComEventLogRequest, ok := typ.(ModbusPDUGetComEventLogRequest); ok {
            return sModbusPDUGetComEventLogRequest
        }
        if sModbusPDUGetComEventLogRequest, ok := typ.(*ModbusPDUGetComEventLogRequest); ok {
            return *sModbusPDUGetComEventLogRequest
        }
        return ModbusPDUGetComEventLogRequest{}
    }
    return castFunc(structType)
}

func (m ModbusPDUGetComEventLogRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    return lengthInBits
}

func (m ModbusPDUGetComEventLogRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUGetComEventLogRequestParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Create the instance
    return NewModbusPDUGetComEventLogRequest(), nil
}

func (m ModbusPDUGetComEventLogRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
