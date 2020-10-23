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
type ModbusPDUReadDiscreteInputsResponse struct {
    Value []int8
    ModbusPDU
}

// The corresponding interface
type IModbusPDUReadDiscreteInputsResponse interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUReadDiscreteInputsResponse) ErrorFlag() bool {
    return false
}

func (m ModbusPDUReadDiscreteInputsResponse) FunctionFlag() uint8 {
    return 0x02
}

func (m ModbusPDUReadDiscreteInputsResponse) Response() bool {
    return true
}

func (m ModbusPDUReadDiscreteInputsResponse) initialize() spi.Message {
    return m
}

func NewModbusPDUReadDiscreteInputsResponse(value []int8) ModbusPDUInitializer {
    return &ModbusPDUReadDiscreteInputsResponse{Value: value}
}

func CastIModbusPDUReadDiscreteInputsResponse(structType interface{}) IModbusPDUReadDiscreteInputsResponse {
    castFunc := func(typ interface{}) IModbusPDUReadDiscreteInputsResponse {
        if iModbusPDUReadDiscreteInputsResponse, ok := typ.(IModbusPDUReadDiscreteInputsResponse); ok {
            return iModbusPDUReadDiscreteInputsResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUReadDiscreteInputsResponse(structType interface{}) ModbusPDUReadDiscreteInputsResponse {
    castFunc := func(typ interface{}) ModbusPDUReadDiscreteInputsResponse {
        if sModbusPDUReadDiscreteInputsResponse, ok := typ.(ModbusPDUReadDiscreteInputsResponse); ok {
            return sModbusPDUReadDiscreteInputsResponse
        }
        if sModbusPDUReadDiscreteInputsResponse, ok := typ.(*ModbusPDUReadDiscreteInputsResponse); ok {
            return *sModbusPDUReadDiscreteInputsResponse
        }
        return ModbusPDUReadDiscreteInputsResponse{}
    }
    return castFunc(structType)
}

func (m ModbusPDUReadDiscreteInputsResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Implicit Field (byteCount)
    lengthInBits += 8

    // Array field
    if len(m.Value) > 0 {
        lengthInBits += 8 * uint16(len(m.Value))
    }

    return lengthInBits
}

func (m ModbusPDUReadDiscreteInputsResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadDiscreteInputsResponseParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    byteCount, _byteCountErr := io.ReadUint8(8)
    if _byteCountErr != nil {
        return nil, errors.New("Error parsing 'byteCount' field " + _byteCountErr.Error())
    }

    // Array field (value)
    // Count array
    value := make([]int8, byteCount)
    for curItem := uint16(0); curItem < uint16(byteCount); curItem++ {

        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'value' field " + _err.Error())
        }
        value[curItem] = _item
    }

    // Create the instance
    return NewModbusPDUReadDiscreteInputsResponse(value), nil
}

func (m ModbusPDUReadDiscreteInputsResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    byteCount := uint8(uint8(len(m.Value)))
    _byteCountErr := io.WriteUint8(8, (byteCount))
    if _byteCountErr != nil {
        return errors.New("Error serializing 'byteCount' field " + _byteCountErr.Error())
    }

    // Array Field (value)
    if m.Value != nil {
        for _, _element := range m.Value {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'value' field " + _elementErr.Error())
            }
        }
    }

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
