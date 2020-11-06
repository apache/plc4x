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
    "io"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type S7MessageRequest struct {
    Parent *S7Message
    IS7MessageRequest
}

// The corresponding interface
type IS7MessageRequest interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *S7MessageRequest) MessageType() uint8 {
    return 0x01
}


func (m *S7MessageRequest) InitializeParent(parent *S7Message, tpduReference uint16, parameter *S7Parameter, payload *S7Payload) {
    m.Parent.TpduReference = tpduReference
    m.Parent.Parameter = parameter
    m.Parent.Payload = payload
}

func NewS7MessageRequest(tpduReference uint16, parameter *S7Parameter, payload *S7Payload) *S7Message {
    child := &S7MessageRequest{
        Parent: NewS7Message(tpduReference, parameter, payload),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastS7MessageRequest(structType interface{}) S7MessageRequest {
    castFunc := func(typ interface{}) S7MessageRequest {
        if casted, ok := typ.(S7MessageRequest); ok {
            return casted
        }
        if casted, ok := typ.(*S7MessageRequest); ok {
            return *casted
        }
        if casted, ok := typ.(S7Message); ok {
            return CastS7MessageRequest(casted.Child)
        }
        if casted, ok := typ.(*S7Message); ok {
            return CastS7MessageRequest(casted.Child)
        }
        return S7MessageRequest{}
    }
    return castFunc(structType)
}

func (m *S7MessageRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    return lengthInBits
}

func (m *S7MessageRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7MessageRequestParse(io *utils.ReadBuffer) (*S7Message, error) {

    // Create a partially initialized instance
    _child := &S7MessageRequest{
        Parent: &S7Message{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *S7MessageRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *S7MessageRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            }
        }
    }
}

func (m *S7MessageRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.S7MessageRequest"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

