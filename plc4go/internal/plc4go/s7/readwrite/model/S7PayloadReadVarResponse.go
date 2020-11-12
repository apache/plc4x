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
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
)

// The data-structure of this message
type S7PayloadReadVarResponse struct {
    Items []*S7VarPayloadDataItem
    Parent *S7Payload
    IS7PayloadReadVarResponse
}

// The corresponding interface
type IS7PayloadReadVarResponse interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *S7PayloadReadVarResponse) ParameterParameterType() uint8 {
    return 0x04
}

func (m *S7PayloadReadVarResponse) MessageType() uint8 {
    return 0x03
}


func (m *S7PayloadReadVarResponse) InitializeParent(parent *S7Payload) {
}

func NewS7PayloadReadVarResponse(items []*S7VarPayloadDataItem, ) *S7Payload {
    child := &S7PayloadReadVarResponse{
        Items: items,
        Parent: NewS7Payload(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastS7PayloadReadVarResponse(structType interface{}) *S7PayloadReadVarResponse {
    castFunc := func(typ interface{}) *S7PayloadReadVarResponse {
        if casted, ok := typ.(S7PayloadReadVarResponse); ok {
            return &casted
        }
        if casted, ok := typ.(*S7PayloadReadVarResponse); ok {
            return casted
        }
        if casted, ok := typ.(S7Payload); ok {
            return CastS7PayloadReadVarResponse(casted.Child)
        }
        if casted, ok := typ.(*S7Payload); ok {
            return CastS7PayloadReadVarResponse(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *S7PayloadReadVarResponse) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Array field
    if len(m.Items) > 0 {
        for _, element := range m.Items {
            lengthInBits += element.LengthInBits()
        }
    }

    return lengthInBits
}

func (m *S7PayloadReadVarResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7PayloadReadVarResponseParse(io *utils.ReadBuffer, parameter *S7Parameter) (*S7Payload, error) {

    // Array field (items)
    // Count array
    items := make([]*S7VarPayloadDataItem, CastS7ParameterReadVarResponse(parameter).NumItems)
    for curItem := uint16(0); curItem < uint16(CastS7ParameterReadVarResponse(parameter).NumItems); curItem++ {
        lastItem := curItem == uint16(CastS7ParameterReadVarResponse(parameter).NumItems - 1)
        _item, _err := S7VarPayloadDataItemParse(io, lastItem)
        if _err != nil {
            return nil, errors.New("Error parsing 'items' field " + _err.Error())
        }
        items[curItem] = _item
    }

    // Create a partially initialized instance
    _child := &S7PayloadReadVarResponse{
        Items: items,
        Parent: &S7Payload{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *S7PayloadReadVarResponse) Serialize(io utils.WriteBuffer) error {
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

func (m *S7PayloadReadVarResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
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
        token, err = d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
    }
}

func (m *S7PayloadReadVarResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: xml.Name{Local: "items"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Items, xml.StartElement{Name: xml.Name{Local: "items"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: xml.Name{Local: "items"}}); err != nil {
        return err
    }
    return nil
}

