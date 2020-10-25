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
    "encoding/xml"
    "errors"
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
    "reflect"
)

// The data-structure of this message
type ModbusPDUWriteFileRecordResponse struct {
    Items []IModbusPDUWriteFileRecordResponseItem
    ModbusPDU
}

// The corresponding interface
type IModbusPDUWriteFileRecordResponse interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUWriteFileRecordResponse) ErrorFlag() bool {
    return false
}

func (m ModbusPDUWriteFileRecordResponse) FunctionFlag() uint8 {
    return 0x15
}

func (m ModbusPDUWriteFileRecordResponse) Response() bool {
    return true
}

func (m ModbusPDUWriteFileRecordResponse) initialize() spi.Message {
    return m
}

func NewModbusPDUWriteFileRecordResponse(items []IModbusPDUWriteFileRecordResponseItem) ModbusPDUInitializer {
    return &ModbusPDUWriteFileRecordResponse{Items: items}
}

func CastIModbusPDUWriteFileRecordResponse(structType interface{}) IModbusPDUWriteFileRecordResponse {
    castFunc := func(typ interface{}) IModbusPDUWriteFileRecordResponse {
        if iModbusPDUWriteFileRecordResponse, ok := typ.(IModbusPDUWriteFileRecordResponse); ok {
            return iModbusPDUWriteFileRecordResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUWriteFileRecordResponse(structType interface{}) ModbusPDUWriteFileRecordResponse {
    castFunc := func(typ interface{}) ModbusPDUWriteFileRecordResponse {
        if sModbusPDUWriteFileRecordResponse, ok := typ.(ModbusPDUWriteFileRecordResponse); ok {
            return sModbusPDUWriteFileRecordResponse
        }
        if sModbusPDUWriteFileRecordResponse, ok := typ.(*ModbusPDUWriteFileRecordResponse); ok {
            return *sModbusPDUWriteFileRecordResponse
        }
        return ModbusPDUWriteFileRecordResponse{}
    }
    return castFunc(structType)
}

func (m ModbusPDUWriteFileRecordResponse) LengthInBits() uint16 {
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

func (m ModbusPDUWriteFileRecordResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUWriteFileRecordResponseParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    byteCount, _byteCountErr := io.ReadUint8(8)
    if _byteCountErr != nil {
        return nil, errors.New("Error parsing 'byteCount' field " + _byteCountErr.Error())
    }

    // Array field (items)
    // Length array
    items := make([]IModbusPDUWriteFileRecordResponseItem, 0)
    _itemsLength := byteCount
    _itemsEndPos := io.GetPos() + uint16(_itemsLength)
    for ;io.GetPos() < _itemsEndPos; {
        _message, _err := ModbusPDUWriteFileRecordResponseItemParse(io)
        if _err != nil {
            return nil, errors.New("Error parsing 'items' field " + _err.Error())
        }
        var _item IModbusPDUWriteFileRecordResponseItem
        _item, _ok := _message.(IModbusPDUWriteFileRecordResponseItem)
        if !_ok {
            return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ModbusPDUWriteFileRecordResponseItem")
        }
        items = append(items, _item)
    }

    // Create the instance
    return NewModbusPDUWriteFileRecordResponse(items), nil
}

func (m ModbusPDUWriteFileRecordResponse) Serialize(io utils.WriteBuffer) error {
    itemsArraySizeInBytes := func(items []IModbusPDUWriteFileRecordResponseItem) uint32 {
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

func (m *ModbusPDUWriteFileRecordResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    for {
        token, err := d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "items":
                var data []IModbusPDUWriteFileRecordResponseItem
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Items = data
            }
        }
    }
}

func (m ModbusPDUWriteFileRecordResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteFileRecordResponse"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "items"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Items, xml.StartElement{Name: xml.Name{Local: "items"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "items"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

