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
type ModbusPDUReadDiscreteInputsRequest struct {
    StartingAddress uint16
    Quantity uint16
    ModbusPDU
}

// The corresponding interface
type IModbusPDUReadDiscreteInputsRequest interface {
    IModbusPDU
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUReadDiscreteInputsRequest) ErrorFlag() bool {
    return false
}

func (m ModbusPDUReadDiscreteInputsRequest) FunctionFlag() uint8 {
    return 0x02
}

func (m ModbusPDUReadDiscreteInputsRequest) Response() bool {
    return false
}

func (m ModbusPDUReadDiscreteInputsRequest) initialize() spi.Message {
    return m
}

func NewModbusPDUReadDiscreteInputsRequest(startingAddress uint16, quantity uint16) ModbusPDUInitializer {
    return &ModbusPDUReadDiscreteInputsRequest{StartingAddress: startingAddress, Quantity: quantity}
}

func CastIModbusPDUReadDiscreteInputsRequest(structType interface{}) IModbusPDUReadDiscreteInputsRequest {
    castFunc := func(typ interface{}) IModbusPDUReadDiscreteInputsRequest {
        if iModbusPDUReadDiscreteInputsRequest, ok := typ.(IModbusPDUReadDiscreteInputsRequest); ok {
            return iModbusPDUReadDiscreteInputsRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUReadDiscreteInputsRequest(structType interface{}) ModbusPDUReadDiscreteInputsRequest {
    castFunc := func(typ interface{}) ModbusPDUReadDiscreteInputsRequest {
        if sModbusPDUReadDiscreteInputsRequest, ok := typ.(ModbusPDUReadDiscreteInputsRequest); ok {
            return sModbusPDUReadDiscreteInputsRequest
        }
        if sModbusPDUReadDiscreteInputsRequest, ok := typ.(*ModbusPDUReadDiscreteInputsRequest); ok {
            return *sModbusPDUReadDiscreteInputsRequest
        }
        return ModbusPDUReadDiscreteInputsRequest{}
    }
    return castFunc(structType)
}

func (m ModbusPDUReadDiscreteInputsRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (startingAddress)
    lengthInBits += 16

    // Simple field (quantity)
    lengthInBits += 16

    return lengthInBits
}

func (m ModbusPDUReadDiscreteInputsRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadDiscreteInputsRequestParse(io *spi.ReadBuffer) (ModbusPDUInitializer, error) {

    // Simple Field (startingAddress)
    startingAddress, _startingAddressErr := io.ReadUint16(16)
    if _startingAddressErr != nil {
        return nil, errors.New("Error parsing 'startingAddress' field " + _startingAddressErr.Error())
    }

    // Simple Field (quantity)
    quantity, _quantityErr := io.ReadUint16(16)
    if _quantityErr != nil {
        return nil, errors.New("Error parsing 'quantity' field " + _quantityErr.Error())
    }

    // Create the instance
    return NewModbusPDUReadDiscreteInputsRequest(startingAddress, quantity), nil
}

func (m ModbusPDUReadDiscreteInputsRequest) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

    // Simple Field (startingAddress)
    startingAddress := uint16(m.StartingAddress)
    _startingAddressErr := io.WriteUint16(16, (startingAddress))
    if _startingAddressErr != nil {
        return errors.New("Error serializing 'startingAddress' field " + _startingAddressErr.Error())
    }

    // Simple Field (quantity)
    quantity := uint16(m.Quantity)
    _quantityErr := io.WriteUint16(16, (quantity))
    if _quantityErr != nil {
        return errors.New("Error serializing 'quantity' field " + _quantityErr.Error())
    }

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
