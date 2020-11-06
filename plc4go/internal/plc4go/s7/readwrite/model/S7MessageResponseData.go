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
type S7MessageResponseData struct {
    ErrorClass uint8
    ErrorCode uint8
    Parent *S7Message
    IS7MessageResponseData
}

// The corresponding interface
type IS7MessageResponseData interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *S7MessageResponseData) MessageType() uint8 {
    return 0x03
}


func (m *S7MessageResponseData) InitializeParent(parent *S7Message, tpduReference uint16, parameter *S7Parameter, payload *S7Payload) {
    m.Parent.TpduReference = tpduReference
    m.Parent.Parameter = parameter
    m.Parent.Payload = payload
}

func NewS7MessageResponseData(errorClass uint8, errorCode uint8, tpduReference uint16, parameter *S7Parameter, payload *S7Payload) *S7Message {
    child := &S7MessageResponseData{
        ErrorClass: errorClass,
        ErrorCode: errorCode,
        Parent: NewS7Message(tpduReference, parameter, payload),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastS7MessageResponseData(structType interface{}) S7MessageResponseData {
    castFunc := func(typ interface{}) S7MessageResponseData {
        if casted, ok := typ.(S7MessageResponseData); ok {
            return casted
        }
        if casted, ok := typ.(*S7MessageResponseData); ok {
            return *casted
        }
        if casted, ok := typ.(S7Message); ok {
            return CastS7MessageResponseData(casted.Child)
        }
        if casted, ok := typ.(*S7Message); ok {
            return CastS7MessageResponseData(casted.Child)
        }
        return S7MessageResponseData{}
    }
    return castFunc(structType)
}

func (m *S7MessageResponseData) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (errorClass)
    lengthInBits += 8

    // Simple field (errorCode)
    lengthInBits += 8

    return lengthInBits
}

func (m *S7MessageResponseData) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7MessageResponseDataParse(io *utils.ReadBuffer) (*S7Message, error) {

    // Simple Field (errorClass)
    errorClass, _errorClassErr := io.ReadUint8(8)
    if _errorClassErr != nil {
        return nil, errors.New("Error parsing 'errorClass' field " + _errorClassErr.Error())
    }

    // Simple Field (errorCode)
    errorCode, _errorCodeErr := io.ReadUint8(8)
    if _errorCodeErr != nil {
        return nil, errors.New("Error parsing 'errorCode' field " + _errorCodeErr.Error())
    }

    // Create a partially initialized instance
    _child := &S7MessageResponseData{
        ErrorClass: errorClass,
        ErrorCode: errorCode,
        Parent: &S7Message{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *S7MessageResponseData) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (errorClass)
    errorClass := uint8(m.ErrorClass)
    _errorClassErr := io.WriteUint8(8, (errorClass))
    if _errorClassErr != nil {
        return errors.New("Error serializing 'errorClass' field " + _errorClassErr.Error())
    }

    // Simple Field (errorCode)
    errorCode := uint8(m.ErrorCode)
    _errorCodeErr := io.WriteUint8(8, (errorCode))
    if _errorCodeErr != nil {
        return errors.New("Error serializing 'errorCode' field " + _errorCodeErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *S7MessageResponseData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "errorClass":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ErrorClass = data
            case "errorCode":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ErrorCode = data
            }
        }
    }
}

func (m *S7MessageResponseData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.S7MessageResponseData"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ErrorClass, xml.StartElement{Name: xml.Name{Local: "errorClass"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ErrorCode, xml.StartElement{Name: xml.Name{Local: "errorCode"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

