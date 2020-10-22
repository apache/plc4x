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
    "reflect"
    "strconv"
)

// Constant values.
const S7PayloadUserDataItemCpuFunctionReadSzlResponse_SZLITEMLENGTH uint16 = 28

// The data-structure of this message
type S7PayloadUserDataItemCpuFunctionReadSzlResponse struct {
    Items []ISzlDataTreeItem
    S7PayloadUserDataItem
}

// The corresponding interface
type IS7PayloadUserDataItemCpuFunctionReadSzlResponse interface {
    IS7PayloadUserDataItem
    Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m S7PayloadUserDataItemCpuFunctionReadSzlResponse) CpuFunctionType() uint8 {
    return 0x08
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlResponse) initialize(returnCode IDataTransportErrorCode, transportSize IDataTransportSize, szlId ISzlId, szlIndex uint16) spi.Message {
    m.ReturnCode = returnCode
    m.TransportSize = transportSize
    m.SzlId = szlId
    m.SzlIndex = szlIndex
    return m
}

func NewS7PayloadUserDataItemCpuFunctionReadSzlResponse(items []ISzlDataTreeItem) S7PayloadUserDataItemInitializer {
    return &S7PayloadUserDataItemCpuFunctionReadSzlResponse{Items: items}
}

func CastIS7PayloadUserDataItemCpuFunctionReadSzlResponse(structType interface{}) IS7PayloadUserDataItemCpuFunctionReadSzlResponse {
    castFunc := func(typ interface{}) IS7PayloadUserDataItemCpuFunctionReadSzlResponse {
        if iS7PayloadUserDataItemCpuFunctionReadSzlResponse, ok := typ.(IS7PayloadUserDataItemCpuFunctionReadSzlResponse); ok {
            return iS7PayloadUserDataItemCpuFunctionReadSzlResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastS7PayloadUserDataItemCpuFunctionReadSzlResponse(structType interface{}) S7PayloadUserDataItemCpuFunctionReadSzlResponse {
    castFunc := func(typ interface{}) S7PayloadUserDataItemCpuFunctionReadSzlResponse {
        if sS7PayloadUserDataItemCpuFunctionReadSzlResponse, ok := typ.(S7PayloadUserDataItemCpuFunctionReadSzlResponse); ok {
            return sS7PayloadUserDataItemCpuFunctionReadSzlResponse
        }
        if sS7PayloadUserDataItemCpuFunctionReadSzlResponse, ok := typ.(*S7PayloadUserDataItemCpuFunctionReadSzlResponse); ok {
            return *sS7PayloadUserDataItemCpuFunctionReadSzlResponse
        }
        return S7PayloadUserDataItemCpuFunctionReadSzlResponse{}
    }
    return castFunc(structType)
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.S7PayloadUserDataItem.LengthInBits()

    // Const Field (szlItemLength)
    lengthInBits += 16

    // Implicit Field (szlItemCount)
    lengthInBits += 16

    // Array field
    if len(m.Items) > 0 {
        for _, element := range m.Items {
            lengthInBits += element.LengthInBits()
        }
    }

    return lengthInBits
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7PayloadUserDataItemCpuFunctionReadSzlResponseParse(io *spi.ReadBuffer) (S7PayloadUserDataItemInitializer, error) {

    // Const Field (szlItemLength)
    szlItemLength, _szlItemLengthErr := io.ReadUint16(16)
    if _szlItemLengthErr != nil {
        return nil, errors.New("Error parsing 'szlItemLength' field " + _szlItemLengthErr.Error())
    }
    if szlItemLength != S7PayloadUserDataItemCpuFunctionReadSzlResponse_SZLITEMLENGTH {
        return nil, errors.New("Expected constant value " + strconv.Itoa(int(S7PayloadUserDataItemCpuFunctionReadSzlResponse_SZLITEMLENGTH)) + " but got " + strconv.Itoa(int(szlItemLength)))
    }

    // Implicit Field (szlItemCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    szlItemCount, _szlItemCountErr := io.ReadUint16(16)
    if _szlItemCountErr != nil {
        return nil, errors.New("Error parsing 'szlItemCount' field " + _szlItemCountErr.Error())
    }

    // Array field (items)
    // Count array
    items := make([]ISzlDataTreeItem, szlItemCount)
    for curItem := uint16(0); curItem < uint16(szlItemCount); curItem++ {

        _message, _err := SzlDataTreeItemParse(io)
        if _err != nil {
            return nil, errors.New("Error parsing 'items' field " + _err.Error())
        }
        var _item ISzlDataTreeItem
        _item, _ok := _message.(ISzlDataTreeItem)
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to SzlDataTreeItem")
        }
        items[curItem] = _item
    }

    // Create the instance
    return NewS7PayloadUserDataItemCpuFunctionReadSzlResponse(items), nil
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlResponse) Serialize(io spi.WriteBuffer) error {
    ser := func() error {

    // Const Field (szlItemLength)
    _szlItemLengthErr := io.WriteUint16(16, 28)
    if _szlItemLengthErr != nil {
        return errors.New("Error serializing 'szlItemLength' field " + _szlItemLengthErr.Error())
    }

    // Implicit Field (szlItemCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    szlItemCount := uint16(uint16(len(m.Items)))
    _szlItemCountErr := io.WriteUint16(16, (szlItemCount))
    if _szlItemCountErr != nil {
        return errors.New("Error serializing 'szlItemCount' field " + _szlItemCountErr.Error())
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
    return S7PayloadUserDataItemSerialize(io, m.S7PayloadUserDataItem, CastIS7PayloadUserDataItem(m), ser)
}
