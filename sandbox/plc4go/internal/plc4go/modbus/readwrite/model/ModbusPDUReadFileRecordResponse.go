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
    "reflect"
)

// The data-structure of this message
type ModbusPDUReadFileRecordResponse struct {
    Items []IModbusPDUReadFileRecordResponseItem
    ModbusPDU
}

// The corresponding interface
type IModbusPDUReadFileRecordResponse interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUReadFileRecordResponse) ErrorFlag() bool {
    return false
}

func (m ModbusPDUReadFileRecordResponse) FunctionFlag() uint8 {
    return 0x14
}

func (m ModbusPDUReadFileRecordResponse) Response() bool {
    return true
}

func (m ModbusPDUReadFileRecordResponse) initialize() spi.Message {
    return m
}

func NewModbusPDUReadFileRecordResponse(items []IModbusPDUReadFileRecordResponseItem) ModbusPDUInitializer {
    return &ModbusPDUReadFileRecordResponse{Items: items}
}

func CastIModbusPDUReadFileRecordResponse(structType interface{}) IModbusPDUReadFileRecordResponse {
    castFunc := func(typ interface{}) IModbusPDUReadFileRecordResponse {
        if iModbusPDUReadFileRecordResponse, ok := typ.(IModbusPDUReadFileRecordResponse); ok {
            return iModbusPDUReadFileRecordResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUReadFileRecordResponse(structType interface{}) ModbusPDUReadFileRecordResponse {
    castFunc := func(typ interface{}) ModbusPDUReadFileRecordResponse {
        if sModbusPDUReadFileRecordResponse, ok := typ.(ModbusPDUReadFileRecordResponse); ok {
            return sModbusPDUReadFileRecordResponse
        }
        if sModbusPDUReadFileRecordResponse, ok := typ.(*ModbusPDUReadFileRecordResponse); ok {
            return *sModbusPDUReadFileRecordResponse
        }
        return ModbusPDUReadFileRecordResponse{}
    }
    return castFunc(structType)
}

func (m ModbusPDUReadFileRecordResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Implicit Field (byteCount)
    lengthInBits += 8

    // Array field
    if len(m.Items) > 0 {
        for _, element := range m.Items {
            lengthInBits += element.LengthInBits()
        }
    }

    return lengthInBits
}

func (m ModbusPDUReadFileRecordResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadFileRecordResponseParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    byteCount, _byteCountErr := io.ReadUint8(8)
    if _byteCountErr != nil {
        return nil, errors.New("Error parsing 'byteCount' field " + _byteCountErr.Error())
    }

    // Array field (items)
    // Length array
    items := make([]IModbusPDUReadFileRecordResponseItem, 0)
    _itemsLength := byteCount
    _itemsEndPos := io.GetPos() + uint16(_itemsLength)
    for ;io.GetPos() < _itemsEndPos; {
        _message, _err := ModbusPDUReadFileRecordResponseItemParse(io)
        if _err != nil {
            return nil, errors.New("Error parsing 'items' field " + _err.Error())
        }
        var _item IModbusPDUReadFileRecordResponseItem
        _item, _ok := _message.(IModbusPDUReadFileRecordResponseItem)
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ModbusPDUReadFileRecordResponseItem")
        }
        items = append(items, _item)
    }

    // Create the instance
    return NewModbusPDUReadFileRecordResponse(items), nil
}

func (m ModbusPDUReadFileRecordResponse) Serialize(io utils.WriteBuffer) error {
    itemsArraySizeInBytes := func(items []IModbusPDUReadFileRecordResponseItem) uint32 {
        var sizeInBytes uint32 = 0
        for _, v := range items {
            sizeInBytes += uint32(v.LengthInBytes())
        }
        return sizeInBytes
    }
    ser := func() error {

    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    byteCount := uint8(uint8(itemsArraySizeInBytes(m.Items)))
    _byteCountErr := io.WriteUint8(8, (byteCount))
    if _byteCountErr != nil {
        return errors.New("Error serializing 'byteCount' field " + _byteCountErr.Error())
    }

    // Array Field (items)
    if m.Items != nil {
        for _, _element := range m.Items {
            _elementErr := _element.Serialize(io)
            if _elementErr != nil {
                return errors.New("Error serializing 'items' field " + _elementErr.Error())
            }
        }
    }

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
