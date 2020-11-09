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
type ModbusPDUReadExceptionStatusResponse struct {
    Value uint8
    Parent *ModbusPDU
    IModbusPDUReadExceptionStatusResponse
}

// The corresponding interface
type IModbusPDUReadExceptionStatusResponse interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *ModbusPDUReadExceptionStatusResponse) ErrorFlag() bool {
    return false
}

func (m *ModbusPDUReadExceptionStatusResponse) FunctionFlag() uint8 {
    return 0x07
}

func (m *ModbusPDUReadExceptionStatusResponse) Response() bool {
    return true
}


func (m *ModbusPDUReadExceptionStatusResponse) InitializeParent(parent *ModbusPDU) {
}

func NewModbusPDUReadExceptionStatusResponse(value uint8, ) *ModbusPDU {
    child := &ModbusPDUReadExceptionStatusResponse{
        Value: value,
        Parent: NewModbusPDU(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastModbusPDUReadExceptionStatusResponse(structType interface{}) *ModbusPDUReadExceptionStatusResponse {
    castFunc := func(typ interface{}) *ModbusPDUReadExceptionStatusResponse {
        if casted, ok := typ.(ModbusPDUReadExceptionStatusResponse); ok {
            return &casted
        }
        if casted, ok := typ.(*ModbusPDUReadExceptionStatusResponse); ok {
            return casted
        }
        if casted, ok := typ.(ModbusPDU); ok {
            return CastModbusPDUReadExceptionStatusResponse(casted.Child)
        }
        if casted, ok := typ.(*ModbusPDU); ok {
            return CastModbusPDUReadExceptionStatusResponse(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *ModbusPDUReadExceptionStatusResponse) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (value)
    lengthInBits += 8

    return lengthInBits
}

func (m *ModbusPDUReadExceptionStatusResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadExceptionStatusResponseParse(io *utils.ReadBuffer) (*ModbusPDU, error) {

    // Simple Field (value)
    value, _valueErr := io.ReadUint8(8)
    if _valueErr != nil {
        return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
    }

    // Create a partially initialized instance
    _child := &ModbusPDUReadExceptionStatusResponse{
        Value: value,
        Parent: &ModbusPDU{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *ModbusPDUReadExceptionStatusResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (value)
    value := uint8(m.Value)
    _valueErr := io.WriteUint8(8, (value))
    if _valueErr != nil {
        return errors.New("Error serializing 'value' field " + _valueErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *ModbusPDUReadExceptionStatusResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            case "value":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Value = data
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

func (m *ModbusPDUReadExceptionStatusResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.Value, xml.StartElement{Name: xml.Name{Local: "value"}}); err != nil {
        return err
    }
    return nil
}

