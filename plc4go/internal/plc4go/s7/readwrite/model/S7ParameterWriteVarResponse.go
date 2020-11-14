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
type S7ParameterWriteVarResponse struct {
    NumItems uint8
    Parent *S7Parameter
    IS7ParameterWriteVarResponse
}

// The corresponding interface
type IS7ParameterWriteVarResponse interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *S7ParameterWriteVarResponse) ParameterType() uint8 {
    return 0x05
}

func (m *S7ParameterWriteVarResponse) MessageType() uint8 {
    return 0x03
}


func (m *S7ParameterWriteVarResponse) InitializeParent(parent *S7Parameter) {
}

func NewS7ParameterWriteVarResponse(numItems uint8, ) *S7Parameter {
    child := &S7ParameterWriteVarResponse{
        NumItems: numItems,
        Parent: NewS7Parameter(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastS7ParameterWriteVarResponse(structType interface{}) *S7ParameterWriteVarResponse {
    castFunc := func(typ interface{}) *S7ParameterWriteVarResponse {
        if casted, ok := typ.(S7ParameterWriteVarResponse); ok {
            return &casted
        }
        if casted, ok := typ.(*S7ParameterWriteVarResponse); ok {
            return casted
        }
        if casted, ok := typ.(S7Parameter); ok {
            return CastS7ParameterWriteVarResponse(casted.Child)
        }
        if casted, ok := typ.(*S7Parameter); ok {
            return CastS7ParameterWriteVarResponse(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *S7ParameterWriteVarResponse) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (numItems)
    lengthInBits += 8

    return lengthInBits
}

func (m *S7ParameterWriteVarResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7ParameterWriteVarResponseParse(io *utils.ReadBuffer) (*S7Parameter, error) {

    // Simple Field (numItems)
    numItems, _numItemsErr := io.ReadUint8(8)
    if _numItemsErr != nil {
        return nil, errors.New("Error parsing 'numItems' field " + _numItemsErr.Error())
    }

    // Create a partially initialized instance
    _child := &S7ParameterWriteVarResponse{
        NumItems: numItems,
        Parent: &S7Parameter{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *S7ParameterWriteVarResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (numItems)
    numItems := uint8(m.NumItems)
    _numItemsErr := io.WriteUint8(8, (numItems))
    if _numItemsErr != nil {
        return errors.New("Error serializing 'numItems' field " + _numItemsErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *S7ParameterWriteVarResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "numItems":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.NumItems = data
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

func (m *S7ParameterWriteVarResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.NumItems, xml.StartElement{Name: xml.Name{Local: "numItems"}}); err != nil {
        return err
    }
    return nil
}

