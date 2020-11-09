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
type COTPParameterCalledTsap struct {
    TsapId uint16
    Parent *COTPParameter
    ICOTPParameterCalledTsap
}

// The corresponding interface
type ICOTPParameterCalledTsap interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *COTPParameterCalledTsap) ParameterType() uint8 {
    return 0xC2
}


func (m *COTPParameterCalledTsap) InitializeParent(parent *COTPParameter) {
}

func NewCOTPParameterCalledTsap(tsapId uint16, ) *COTPParameter {
    child := &COTPParameterCalledTsap{
        TsapId: tsapId,
        Parent: NewCOTPParameter(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastCOTPParameterCalledTsap(structType interface{}) *COTPParameterCalledTsap {
    castFunc := func(typ interface{}) *COTPParameterCalledTsap {
        if casted, ok := typ.(COTPParameterCalledTsap); ok {
            return &casted
        }
        if casted, ok := typ.(*COTPParameterCalledTsap); ok {
            return casted
        }
        if casted, ok := typ.(COTPParameter); ok {
            return CastCOTPParameterCalledTsap(casted.Child)
        }
        if casted, ok := typ.(*COTPParameter); ok {
            return CastCOTPParameterCalledTsap(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *COTPParameterCalledTsap) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (tsapId)
    lengthInBits += 16

    return lengthInBits
}

func (m *COTPParameterCalledTsap) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func COTPParameterCalledTsapParse(io *utils.ReadBuffer) (*COTPParameter, error) {

    // Simple Field (tsapId)
    tsapId, _tsapIdErr := io.ReadUint16(16)
    if _tsapIdErr != nil {
        return nil, errors.New("Error parsing 'tsapId' field " + _tsapIdErr.Error())
    }

    // Create a partially initialized instance
    _child := &COTPParameterCalledTsap{
        TsapId: tsapId,
        Parent: &COTPParameter{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *COTPParameterCalledTsap) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (tsapId)
    tsapId := uint16(m.TsapId)
    _tsapIdErr := io.WriteUint16(16, (tsapId))
    if _tsapIdErr != nil {
        return errors.New("Error serializing 'tsapId' field " + _tsapIdErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *COTPParameterCalledTsap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "tsapId":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.TsapId = data
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

func (m *COTPParameterCalledTsap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.TsapId, xml.StartElement{Name: xml.Name{Local: "tsapId"}}); err != nil {
        return err
    }
    return nil
}

