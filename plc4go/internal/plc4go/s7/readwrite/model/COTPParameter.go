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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    "io"
    "reflect"
    "strings"
)

// The data-structure of this message
type COTPParameter struct {
    Child ICOTPParameterChild
    ICOTPParameter
    ICOTPParameterParent
}

// The corresponding interface
type ICOTPParameter interface {
    ParameterType() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

type ICOTPParameterParent interface {
    SerializeParent(io utils.WriteBuffer, child ICOTPParameter, serializeChildFunction func() error) error
    GetTypeName() string
}

type ICOTPParameterChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *COTPParameter)
    GetTypeName() string
    ICOTPParameter
}

func NewCOTPParameter() *COTPParameter {
    return &COTPParameter{}
}

func CastCOTPParameter(structType interface{}) *COTPParameter {
    castFunc := func(typ interface{}) *COTPParameter {
        if casted, ok := typ.(COTPParameter); ok {
            return &casted
        }
        if casted, ok := typ.(*COTPParameter); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *COTPParameter) GetTypeName() string {
    return "COTPParameter"
}

func (m *COTPParameter) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (parameterType)
    lengthInBits += 8

    // Implicit Field (parameterLength)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *COTPParameter) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func COTPParameterParse(io *utils.ReadBuffer, rest uint8) (*COTPParameter, error) {

    // Discriminator Field (parameterType) (Used as input to a switch field)
    parameterType, _parameterTypeErr := io.ReadUint8(8)
    if _parameterTypeErr != nil {
        return nil, errors.New("Error parsing 'parameterType' field " + _parameterTypeErr.Error())
    }

    // Implicit Field (parameterLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    _, _parameterLengthErr := io.ReadUint8(8)
    if _parameterLengthErr != nil {
        return nil, errors.New("Error parsing 'parameterLength' field " + _parameterLengthErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *COTPParameter
    var typeSwitchError error
    switch {
    case parameterType == 0xC0:
        _parent, typeSwitchError = COTPParameterTpduSizeParse(io)
    case parameterType == 0xC1:
        _parent, typeSwitchError = COTPParameterCallingTsapParse(io)
    case parameterType == 0xC2:
        _parent, typeSwitchError = COTPParameterCalledTsapParse(io)
    case parameterType == 0xC3:
        _parent, typeSwitchError = COTPParameterChecksumParse(io)
    case parameterType == 0xE0:
        _parent, typeSwitchError = COTPParameterDisconnectAdditionalInformationParse(io, rest)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *COTPParameter) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *COTPParameter) SerializeParent(io utils.WriteBuffer, child ICOTPParameter, serializeChildFunction func() error) error {

    // Discriminator Field (parameterType) (Used as input to a switch field)
    parameterType := uint8(child.ParameterType())
    _parameterTypeErr := io.WriteUint8(8, (parameterType))
    if _parameterTypeErr != nil {
        return errors.New("Error serializing 'parameterType' field " + _parameterTypeErr.Error())
    }

    // Implicit Field (parameterLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    parameterLength := uint8(uint8(uint8(m.LengthInBytes())) - uint8(uint8(2)))
    _parameterLengthErr := io.WriteUint8(8, (parameterLength))
    if _parameterLengthErr != nil {
        return errors.New("Error serializing 'parameterLength' field " + _parameterLengthErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *COTPParameter) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
                    case "org.apache.plc4x.java.s7.readwrite.COTPParameterTpduSize":
                        var dt *COTPParameterTpduSize
                        if m.Child != nil {
                            dt = m.Child.(*COTPParameterTpduSize)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.s7.readwrite.COTPParameterCallingTsap":
                        var dt *COTPParameterCallingTsap
                        if m.Child != nil {
                            dt = m.Child.(*COTPParameterCallingTsap)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.s7.readwrite.COTPParameterCalledTsap":
                        var dt *COTPParameterCalledTsap
                        if m.Child != nil {
                            dt = m.Child.(*COTPParameterCalledTsap)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.s7.readwrite.COTPParameterChecksum":
                        var dt *COTPParameterChecksum
                        if m.Child != nil {
                            dt = m.Child.(*COTPParameterChecksum)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.s7.readwrite.COTPParameterDisconnectAdditionalInformation":
                        var dt *COTPParameterDisconnectAdditionalInformation
                        if m.Child != nil {
                            dt = m.Child.(*COTPParameterDisconnectAdditionalInformation)
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

func (m *COTPParameter) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

