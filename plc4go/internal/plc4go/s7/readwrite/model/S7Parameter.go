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
type S7Parameter struct {
    Child IS7ParameterChild
    IS7Parameter
    IS7ParameterParent
}

// The corresponding interface
type IS7Parameter interface {
    MessageType() uint8
    ParameterType() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

type IS7ParameterParent interface {
    SerializeParent(io utils.WriteBuffer, child IS7Parameter, serializeChildFunction func() error) error
}

type IS7ParameterChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *S7Parameter)
    IS7Parameter
}

func NewS7Parameter() *S7Parameter {
    return &S7Parameter{}
}

func CastS7Parameter(structType interface{}) S7Parameter {
    castFunc := func(typ interface{}) S7Parameter {
        if casted, ok := typ.(S7Parameter); ok {
            return casted
        }
        if casted, ok := typ.(*S7Parameter); ok {
            return *casted
        }
        return S7Parameter{}
    }
    return castFunc(structType)
}

func (m *S7Parameter) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (parameterType)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *S7Parameter) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7ParameterParse(io *utils.ReadBuffer, messageType uint8) (*S7Parameter, error) {

    // Discriminator Field (parameterType) (Used as input to a switch field)
    parameterType, _parameterTypeErr := io.ReadUint8(8)
    if _parameterTypeErr != nil {
        return nil, errors.New("Error parsing 'parameterType' field " + _parameterTypeErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *S7Parameter
    var typeSwitchError error
    switch {
    case parameterType == 0xF0:
        _parent, typeSwitchError = S7ParameterSetupCommunicationParse(io)
    case parameterType == 0x04 && messageType == 0x01:
        _parent, typeSwitchError = S7ParameterReadVarRequestParse(io)
    case parameterType == 0x04 && messageType == 0x03:
        _parent, typeSwitchError = S7ParameterReadVarResponseParse(io)
    case parameterType == 0x05 && messageType == 0x01:
        _parent, typeSwitchError = S7ParameterWriteVarRequestParse(io)
    case parameterType == 0x05 && messageType == 0x03:
        _parent, typeSwitchError = S7ParameterWriteVarResponseParse(io)
    case parameterType == 0x00 && messageType == 0x07:
        _parent, typeSwitchError = S7ParameterUserDataParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *S7Parameter) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *S7Parameter) SerializeParent(io utils.WriteBuffer, child IS7Parameter, serializeChildFunction func() error) error {

    // Discriminator Field (parameterType) (Used as input to a switch field)
    parameterType := uint8(child.ParameterType())
    _parameterTypeErr := io.WriteUint8(8, (parameterType))
    if _parameterTypeErr != nil {
        return errors.New("Error serializing 'parameterType' field " + _parameterTypeErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *S7Parameter) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m *S7Parameter) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.S7Parameter"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

