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
    "errors"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ModbusPDUGetComEventCounterResponse struct {
    Status uint16
    EventCount uint16
    ModbusPDU
}

// The corresponding interface
type IModbusPDUGetComEventCounterResponse interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUGetComEventCounterResponse) ErrorFlag() bool {
    return false
}

func (m ModbusPDUGetComEventCounterResponse) FunctionFlag() uint8 {
    return 0x0B
}

func (m ModbusPDUGetComEventCounterResponse) Response() bool {
    return true
}

func (m ModbusPDUGetComEventCounterResponse) initialize() spi.Message {
    return m
}

func NewModbusPDUGetComEventCounterResponse(status uint16, eventCount uint16) ModbusPDUInitializer {
    return &ModbusPDUGetComEventCounterResponse{Status: status, EventCount: eventCount}
}

func CastIModbusPDUGetComEventCounterResponse(structType interface{}) IModbusPDUGetComEventCounterResponse {
    castFunc := func(typ interface{}) IModbusPDUGetComEventCounterResponse {
        if iModbusPDUGetComEventCounterResponse, ok := typ.(IModbusPDUGetComEventCounterResponse); ok {
            return iModbusPDUGetComEventCounterResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUGetComEventCounterResponse(structType interface{}) ModbusPDUGetComEventCounterResponse {
    castFunc := func(typ interface{}) ModbusPDUGetComEventCounterResponse {
        if sModbusPDUGetComEventCounterResponse, ok := typ.(ModbusPDUGetComEventCounterResponse); ok {
            return sModbusPDUGetComEventCounterResponse
        }
        if sModbusPDUGetComEventCounterResponse, ok := typ.(*ModbusPDUGetComEventCounterResponse); ok {
            return *sModbusPDUGetComEventCounterResponse
        }
        return ModbusPDUGetComEventCounterResponse{}
    }
    return castFunc(structType)
}

func (m ModbusPDUGetComEventCounterResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (status)
    lengthInBits += 16

    // Simple field (eventCount)
    lengthInBits += 16

    return lengthInBits
}

func (m ModbusPDUGetComEventCounterResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUGetComEventCounterResponseParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Simple Field (status)
    status, _statusErr := io.ReadUint16(16)
    if _statusErr != nil {
        return nil, errors.New("Error parsing 'status' field " + _statusErr.Error())
    }

    // Simple Field (eventCount)
    eventCount, _eventCountErr := io.ReadUint16(16)
    if _eventCountErr != nil {
        return nil, errors.New("Error parsing 'eventCount' field " + _eventCountErr.Error())
    }

    // Create the instance
    return NewModbusPDUGetComEventCounterResponse(status, eventCount), nil
}

func (m ModbusPDUGetComEventCounterResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (status)
    status := uint16(m.Status)
    _statusErr := io.WriteUint16(16, (status))
    if _statusErr != nil {
        return errors.New("Error serializing 'status' field " + _statusErr.Error())
    }

    // Simple Field (eventCount)
    eventCount := uint16(m.EventCount)
    _eventCountErr := io.WriteUint16(16, (eventCount))
    if _eventCountErr != nil {
        return errors.New("Error serializing 'eventCount' field " + _eventCountErr.Error())
    }

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
