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
type ModbusPDUReadWriteMultipleHoldingRegistersRequest struct {
    ReadStartingAddress uint16
    ReadQuantity uint16
    WriteStartingAddress uint16
    WriteQuantity uint16
    Value []int8
    ModbusPDU
}

// The corresponding interface
type IModbusPDUReadWriteMultipleHoldingRegistersRequest interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUReadWriteMultipleHoldingRegistersRequest) ErrorFlag() bool {
    return false
}

func (m ModbusPDUReadWriteMultipleHoldingRegistersRequest) FunctionFlag() uint8 {
    return 0x17
}

func (m ModbusPDUReadWriteMultipleHoldingRegistersRequest) Response() bool {
    return false
}

func (m ModbusPDUReadWriteMultipleHoldingRegistersRequest) initialize() spi.Message {
    return m
}

func NewModbusPDUReadWriteMultipleHoldingRegistersRequest(readStartingAddress uint16, readQuantity uint16, writeStartingAddress uint16, writeQuantity uint16, value []int8) ModbusPDUInitializer {
    return &ModbusPDUReadWriteMultipleHoldingRegistersRequest{ReadStartingAddress: readStartingAddress, ReadQuantity: readQuantity, WriteStartingAddress: writeStartingAddress, WriteQuantity: writeQuantity, Value: value}
}

func CastIModbusPDUReadWriteMultipleHoldingRegistersRequest(structType interface{}) IModbusPDUReadWriteMultipleHoldingRegistersRequest {
    castFunc := func(typ interface{}) IModbusPDUReadWriteMultipleHoldingRegistersRequest {
        if iModbusPDUReadWriteMultipleHoldingRegistersRequest, ok := typ.(IModbusPDUReadWriteMultipleHoldingRegistersRequest); ok {
            return iModbusPDUReadWriteMultipleHoldingRegistersRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUReadWriteMultipleHoldingRegistersRequest(structType interface{}) ModbusPDUReadWriteMultipleHoldingRegistersRequest {
    castFunc := func(typ interface{}) ModbusPDUReadWriteMultipleHoldingRegistersRequest {
        if sModbusPDUReadWriteMultipleHoldingRegistersRequest, ok := typ.(ModbusPDUReadWriteMultipleHoldingRegistersRequest); ok {
            return sModbusPDUReadWriteMultipleHoldingRegistersRequest
        }
        if sModbusPDUReadWriteMultipleHoldingRegistersRequest, ok := typ.(*ModbusPDUReadWriteMultipleHoldingRegistersRequest); ok {
            return *sModbusPDUReadWriteMultipleHoldingRegistersRequest
        }
        return ModbusPDUReadWriteMultipleHoldingRegistersRequest{}
    }
    return castFunc(structType)
}

func (m ModbusPDUReadWriteMultipleHoldingRegistersRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (readStartingAddress)
    lengthInBits += 16

    // Simple field (readQuantity)
    lengthInBits += 16

    // Simple field (writeStartingAddress)
    lengthInBits += 16

    // Simple field (writeQuantity)
    lengthInBits += 16

    // Implicit Field (byteCount)
    lengthInBits += 8

    // Array field
    if len(m.Value) > 0 {
        lengthInBits += 8 * uint16(len(m.Value))
    }

    return lengthInBits
}

func (m ModbusPDUReadWriteMultipleHoldingRegistersRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadWriteMultipleHoldingRegistersRequestParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Simple Field (readStartingAddress)
    readStartingAddress, _readStartingAddressErr := io.ReadUint16(16)
    if _readStartingAddressErr != nil {
        return nil, errors.New("Error parsing 'readStartingAddress' field " + _readStartingAddressErr.Error())
    }

    // Simple Field (readQuantity)
    readQuantity, _readQuantityErr := io.ReadUint16(16)
    if _readQuantityErr != nil {
        return nil, errors.New("Error parsing 'readQuantity' field " + _readQuantityErr.Error())
    }

    // Simple Field (writeStartingAddress)
    writeStartingAddress, _writeStartingAddressErr := io.ReadUint16(16)
    if _writeStartingAddressErr != nil {
        return nil, errors.New("Error parsing 'writeStartingAddress' field " + _writeStartingAddressErr.Error())
    }

    // Simple Field (writeQuantity)
    writeQuantity, _writeQuantityErr := io.ReadUint16(16)
    if _writeQuantityErr != nil {
        return nil, errors.New("Error parsing 'writeQuantity' field " + _writeQuantityErr.Error())
    }

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
    return NewModbusPDUReadWriteMultipleHoldingRegistersRequest(readStartingAddress, readQuantity, writeStartingAddress, writeQuantity, value), nil
}

func (m ModbusPDUReadWriteMultipleHoldingRegistersRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (readStartingAddress)
    readStartingAddress := uint16(m.ReadStartingAddress)
    _readStartingAddressErr := io.WriteUint16(16, (readStartingAddress))
    if _readStartingAddressErr != nil {
        return errors.New("Error serializing 'readStartingAddress' field " + _readStartingAddressErr.Error())
    }

    // Simple Field (readQuantity)
    readQuantity := uint16(m.ReadQuantity)
    _readQuantityErr := io.WriteUint16(16, (readQuantity))
    if _readQuantityErr != nil {
        return errors.New("Error serializing 'readQuantity' field " + _readQuantityErr.Error())
    }

    // Simple Field (writeStartingAddress)
    writeStartingAddress := uint16(m.WriteStartingAddress)
    _writeStartingAddressErr := io.WriteUint16(16, (writeStartingAddress))
    if _writeStartingAddressErr != nil {
        return errors.New("Error serializing 'writeStartingAddress' field " + _writeStartingAddressErr.Error())
    }

    // Simple Field (writeQuantity)
    writeQuantity := uint16(m.WriteQuantity)
    _writeQuantityErr := io.WriteUint16(16, (writeQuantity))
    if _writeQuantityErr != nil {
        return errors.New("Error serializing 'writeQuantity' field " + _writeQuantityErr.Error())
    }

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
