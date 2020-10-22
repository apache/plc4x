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
)

// The data-structure of this message
type ModbusPDUDiagnosticResponse struct {
    SubFunction uint16
    Data uint16
    ModbusPDU
}

// The corresponding interface
type IModbusPDUDiagnosticResponse interface {
    IModbusPDU
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUDiagnosticResponse) ErrorFlag() bool {
    return false
}

func (m ModbusPDUDiagnosticResponse) FunctionFlag() uint8 {
    return 0x08
}

func (m ModbusPDUDiagnosticResponse) Response() bool {
    return true
}

func (m ModbusPDUDiagnosticResponse) initialize() spi.Message {
    return m
}

func NewModbusPDUDiagnosticResponse(subFunction uint16, data uint16) ModbusPDUInitializer {
    return &ModbusPDUDiagnosticResponse{SubFunction: subFunction, Data: data}
}

func CastIModbusPDUDiagnosticResponse(structType interface{}) IModbusPDUDiagnosticResponse {
    castFunc := func(typ interface{}) IModbusPDUDiagnosticResponse {
        if iModbusPDUDiagnosticResponse, ok := typ.(IModbusPDUDiagnosticResponse); ok {
            return iModbusPDUDiagnosticResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUDiagnosticResponse(structType interface{}) ModbusPDUDiagnosticResponse {
    castFunc := func(typ interface{}) ModbusPDUDiagnosticResponse {
        if sModbusPDUDiagnosticResponse, ok := typ.(ModbusPDUDiagnosticResponse); ok {
            return sModbusPDUDiagnosticResponse
        }
        if sModbusPDUDiagnosticResponse, ok := typ.(*ModbusPDUDiagnosticResponse); ok {
            return *sModbusPDUDiagnosticResponse
        }
        return ModbusPDUDiagnosticResponse{}
    }
    return castFunc(structType)
}

func (m ModbusPDUDiagnosticResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (subFunction)
    lengthInBits += 16

    // Simple field (data)
    lengthInBits += 16

    return lengthInBits
}

func (m ModbusPDUDiagnosticResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUDiagnosticResponseParse(io *spi.ReadBuffer) (ModbusPDUInitializer, error) {

    // Simple Field (subFunction)
    subFunction, _subFunctionErr := io.ReadUint16(16)
    if _subFunctionErr != nil {
        return nil, errors.New("Error parsing 'subFunction' field " + _subFunctionErr.Error())
    }

    // Simple Field (data)
    data, _dataErr := io.ReadUint16(16)
    if _dataErr != nil {
        return nil, errors.New("Error parsing 'data' field " + _dataErr.Error())
    }

    // Create the instance
    return NewModbusPDUDiagnosticResponse(subFunction, data), nil
}

func (m ModbusPDUDiagnosticResponse) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

    // Simple Field (subFunction)
    subFunction := uint16(m.SubFunction)
    _subFunctionErr := io.WriteUint16(16, (subFunction))
    if _subFunctionErr != nil {
        return errors.New("Error serializing 'subFunction' field " + _subFunctionErr.Error())
    }

    // Simple Field (data)
    data := uint16(m.Data)
    _dataErr := io.WriteUint16(16, (data))
    if _dataErr != nil {
        return errors.New("Error serializing 'data' field " + _dataErr.Error())
    }

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
