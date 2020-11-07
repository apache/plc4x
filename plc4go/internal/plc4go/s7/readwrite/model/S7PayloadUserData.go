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
type S7PayloadUserData struct {
    Items []*S7PayloadUserDataItem
    Parent *S7Payload
    IS7PayloadUserData
}

// The corresponding interface
type IS7PayloadUserData interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *S7PayloadUserData) ParameterParameterType() uint8 {
    return 0x00
}

func (m *S7PayloadUserData) MessageType() uint8 {
    return 0x07
}


func (m *S7PayloadUserData) InitializeParent(parent *S7Payload) {
}

func NewS7PayloadUserData(items []*S7PayloadUserDataItem, ) *S7Payload {
    child := &S7PayloadUserData{
        Items: items,
        Parent: NewS7Payload(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastS7PayloadUserData(structType interface{}) S7PayloadUserData {
    castFunc := func(typ interface{}) S7PayloadUserData {
        if casted, ok := typ.(S7PayloadUserData); ok {
            return casted
        }
        if casted, ok := typ.(*S7PayloadUserData); ok {
            return *casted
        }
        if casted, ok := typ.(S7Payload); ok {
            return CastS7PayloadUserData(casted.Child)
        }
        if casted, ok := typ.(*S7Payload); ok {
            return CastS7PayloadUserData(casted.Child)
        }
        return S7PayloadUserData{}
    }
    return castFunc(structType)
}

func (m *S7PayloadUserData) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Array field
    if len(m.Items) > 0 {
        for _, element := range m.Items {
            lengthInBits += element.LengthInBits()
        }
    }

    return lengthInBits
}

func (m *S7PayloadUserData) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7PayloadUserDataParse(io *utils.ReadBuffer, parameter *S7Parameter) (*S7Payload, error) {

    // Array field (items)
    // Count array
    items := make([]*S7PayloadUserDataItem, uint16(len(CastS7ParameterUserData(parameter).Items)))
    for curItem := uint16(0); curItem < uint16(uint16(len(CastS7ParameterUserData(parameter).Items))); curItem++ {
        _item, _err := S7PayloadUserDataItemParse(io, CastS7ParameterUserDataItemCPUFunctions(CastS7ParameterUserData(parameter).Items).CpuFunctionType)
        if _err != nil {
            return nil, errors.New("Error parsing 'items' field " + _err.Error())
        }
        items[curItem] = _item
    }

    // Create a partially initialized instance
    _child := &S7PayloadUserData{
        Items: items,
        Parent: &S7Payload{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *S7PayloadUserData) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

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
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *S7PayloadUserData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "items":
                var _values []*S7PayloadUserDataItem
                var dt *S7PayloadUserDataItem
                if err := d.DecodeElement(&dt, &tok); err != nil {
                    return err
                }
                _values = append(_values, dt)
                m.Items = _values
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

func (m *S7PayloadUserData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

