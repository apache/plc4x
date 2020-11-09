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
    "reflect"
    "strings"
)

// The data-structure of this message
type S7Payload struct {
    Child IS7PayloadChild
    IS7Payload
    IS7PayloadParent
}

// The corresponding interface
type IS7Payload interface {
    MessageType() uint8
    ParameterParameterType() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

type IS7PayloadParent interface {
    SerializeParent(io utils.WriteBuffer, child IS7Payload, serializeChildFunction func() error) error
}

type IS7PayloadChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *S7Payload)
    IS7Payload
}

func NewS7Payload() *S7Payload {
    return &S7Payload{}
}

func CastS7Payload(structType interface{}) *S7Payload {
    castFunc := func(typ interface{}) *S7Payload {
        if casted, ok := typ.(S7Payload); ok {
            return &casted
        }
        if casted, ok := typ.(*S7Payload); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *S7Payload) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *S7Payload) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7PayloadParse(io *utils.ReadBuffer, messageType uint8, parameter *S7Parameter) (*S7Payload, error) {

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *S7Payload
    var typeSwitchError error
    switch {
    case CastS7Parameter(parameter).ParameterType() == 0x04 && messageType == 0x03:
        _parent, typeSwitchError = S7PayloadReadVarResponseParse(io, parameter)
    case CastS7Parameter(parameter).ParameterType() == 0x05 && messageType == 0x01:
        _parent, typeSwitchError = S7PayloadWriteVarRequestParse(io, parameter)
    case CastS7Parameter(parameter).ParameterType() == 0x05 && messageType == 0x03:
        _parent, typeSwitchError = S7PayloadWriteVarResponseParse(io, parameter)
    case CastS7Parameter(parameter).ParameterType() == 0x00 && messageType == 0x07:
        _parent, typeSwitchError = S7PayloadUserDataParse(io, parameter)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *S7Payload) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *S7Payload) SerializeParent(io utils.WriteBuffer, child IS7Payload, serializeChildFunction func() error) error {

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *S7Payload) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    for {
        token, err = d.Token()
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
                default:
                    switch start.Attr[0].Value {
                        case "org.apache.plc4x.java.s7.readwrite.S7PayloadReadVarResponse":
                            var dt *S7PayloadReadVarResponse
                            if m.Child != nil {
                                dt = m.Child.(*S7PayloadReadVarResponse)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.s7.readwrite.S7PayloadWriteVarRequest":
                            var dt *S7PayloadWriteVarRequest
                            if m.Child != nil {
                                dt = m.Child.(*S7PayloadWriteVarRequest)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.s7.readwrite.S7PayloadWriteVarResponse":
                            var dt *S7PayloadWriteVarResponse
                            if m.Child != nil {
                                dt = m.Child.(*S7PayloadWriteVarResponse)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.s7.readwrite.S7PayloadUserData":
                            var dt *S7PayloadUserData
                            if m.Child != nil {
                                dt = m.Child.(*S7PayloadUserData)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                    }
            }
        }
    }
}

func (m *S7Payload) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := reflect.TypeOf(m.Child).String()
    className = "org.apache.plc4x.java.s7.readwrite." + className[strings.LastIndex(className, ".") + 1:]
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    marshaller, ok := m.Child.(xml.Marshaler)
    if !ok {
        return errors.New("child is not castable to Marshaler")
    }
    marshaller.MarshalXML(e, start)
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

