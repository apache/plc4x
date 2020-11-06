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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type S7PayloadWriteVarRequest struct {
    Items []*S7VarPayloadDataItem
    Parent *S7Payload
    IS7PayloadWriteVarRequest
}

// The corresponding interface
type IS7PayloadWriteVarRequest interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *S7PayloadWriteVarRequest) ParameterParameterType() uint8 {
    return 0x05
}

func (m *S7PayloadWriteVarRequest) MessageType() uint8 {
    return 0x01
}


func (m *S7PayloadWriteVarRequest) InitializeParent(parent *S7Payload) {
}

func NewS7PayloadWriteVarRequest(items []*S7VarPayloadDataItem, ) *S7Payload {
    child := &S7PayloadWriteVarRequest{
        Items: items,
        Parent: NewS7Payload(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastS7PayloadWriteVarRequest(structType interface{}) S7PayloadWriteVarRequest {
    castFunc := func(typ interface{}) S7PayloadWriteVarRequest {
        if casted, ok := typ.(S7PayloadWriteVarRequest); ok {
            return casted
        }
        if casted, ok := typ.(*S7PayloadWriteVarRequest); ok {
            return *casted
        }
        if casted, ok := typ.(S7Payload); ok {
            return CastS7PayloadWriteVarRequest(casted.Child)
        }
        if casted, ok := typ.(*S7Payload); ok {
            return CastS7PayloadWriteVarRequest(casted.Child)
        }
        return S7PayloadWriteVarRequest{}
    }
    return castFunc(structType)
}

func (m *S7PayloadWriteVarRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Array field
    if len(m.Items) > 0 {
        for _, element := range m.Items {
            lengthInBits += element.LengthInBits()
        }
    }

    return lengthInBits
}

func (m *S7PayloadWriteVarRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7PayloadWriteVarRequestParse(io *utils.ReadBuffer, parameter *S7Parameter) (*S7Payload, error) {

    // Array field (items)
    // Count array
    items := make([]*S7VarPayloadDataItem, uint16(len(CastS7ParameterWriteVarRequest(parameter).Items)))
    for curItem := uint16(0); curItem < uint16(uint16(len(CastS7ParameterWriteVarRequest(parameter).Items))); curItem++ {
        lastItem := curItem == uint16((len(CastS7ParameterWriteVarRequest(parameter).Items)) - 1)
        _item, _err := S7VarPayloadDataItemParse(io, lastItem)
        if _err != nil {
            return nil, errors.New("Error parsing 'items' field " + _err.Error())
        }
        items[curItem] = _item
    }

    // Create a partially initialized instance
    _child := &S7PayloadWriteVarRequest{
        Items: items,
        Parent: &S7Payload{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *S7PayloadWriteVarRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Array Field (items)
    if m.Items != nil {
        itemCount := uint16(len(m.Items))
        var curItem uint16 = 0
        for _, _element := range m.Items {
            var lastItem bool = curItem == (itemCount - 1)
            _elementErr := _element.Serialize(io, lastItem)
            if _elementErr != nil {
                return errors.New("Error serializing 'items' field " + _elementErr.Error())
            }
            curItem++
        }
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *S7PayloadWriteVarRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
                var data []*S7VarPayloadDataItem
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Items = data
            }
        }
    }
}

func (m *S7PayloadWriteVarRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.S7PayloadWriteVarRequest"},
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

