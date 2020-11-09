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
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
)

// The data-structure of this message
type S7ParameterReadVarRequest struct {
    Items []*S7VarRequestParameterItem
    Parent *S7Parameter
    IS7ParameterReadVarRequest
}

// The corresponding interface
type IS7ParameterReadVarRequest interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *S7ParameterReadVarRequest) ParameterType() uint8 {
    return 0x04
}

func (m *S7ParameterReadVarRequest) MessageType() uint8 {
    return 0x01
}


func (m *S7ParameterReadVarRequest) InitializeParent(parent *S7Parameter) {
}

func NewS7ParameterReadVarRequest(items []*S7VarRequestParameterItem, ) *S7Parameter {
    child := &S7ParameterReadVarRequest{
        Items: items,
        Parent: NewS7Parameter(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastS7ParameterReadVarRequest(structType interface{}) *S7ParameterReadVarRequest {
    castFunc := func(typ interface{}) *S7ParameterReadVarRequest {
        if casted, ok := typ.(S7ParameterReadVarRequest); ok {
            return &casted
        }
        if casted, ok := typ.(*S7ParameterReadVarRequest); ok {
            return casted
        }
        if casted, ok := typ.(S7Parameter); ok {
            return CastS7ParameterReadVarRequest(casted.Child)
        }
        if casted, ok := typ.(*S7Parameter); ok {
            return CastS7ParameterReadVarRequest(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *S7ParameterReadVarRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Implicit Field (numItems)
    lengthInBits += 8

    // Array field
    if len(m.Items) > 0 {
        for _, element := range m.Items {
            lengthInBits += element.LengthInBits()
        }
    }

    return lengthInBits
}

func (m *S7ParameterReadVarRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7ParameterReadVarRequestParse(io *utils.ReadBuffer) (*S7Parameter, error) {

    // Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    numItems, _numItemsErr := io.ReadUint8(8)
    if _numItemsErr != nil {
        return nil, errors.New("Error parsing 'numItems' field " + _numItemsErr.Error())
    }

    // Array field (items)
    // Count array
    items := make([]*S7VarRequestParameterItem, numItems)
    for curItem := uint16(0); curItem < uint16(numItems); curItem++ {
        _item, _err := S7VarRequestParameterItemParse(io)
        if _err != nil {
            return nil, errors.New("Error parsing 'items' field " + _err.Error())
        }
        items[curItem] = _item
    }

    // Create a partially initialized instance
    _child := &S7ParameterReadVarRequest{
        Items: items,
        Parent: &S7Parameter{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *S7ParameterReadVarRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    numItems := uint8(uint8(len(m.Items)))
    _numItemsErr := io.WriteUint8(8, (numItems))
    if _numItemsErr != nil {
        return errors.New("Error serializing 'numItems' field " + _numItemsErr.Error())
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
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *S7ParameterReadVarRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "items":
                var _values []*S7VarRequestParameterItem
                var dt *S7VarRequestParameterItem
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

func (m *S7ParameterReadVarRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

