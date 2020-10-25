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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ModbusPDUError struct {
    ExceptionCode uint8
    ModbusPDU
}

// The corresponding interface
type IModbusPDUError interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUError) ErrorFlag() bool {
    return true
}

func (m ModbusPDUError) FunctionFlag() uint8 {
    return 0
}

func (m ModbusPDUError) Response() bool {
    return false
}

func (m ModbusPDUError) initialize() spi.Message {
    return m
}

func NewModbusPDUError(exceptionCode uint8) ModbusPDUInitializer {
    return &ModbusPDUError{ExceptionCode: exceptionCode}
}

func CastIModbusPDUError(structType interface{}) IModbusPDUError {
    castFunc := func(typ interface{}) IModbusPDUError {
        if iModbusPDUError, ok := typ.(IModbusPDUError); ok {
            return iModbusPDUError
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUError(structType interface{}) ModbusPDUError {
    castFunc := func(typ interface{}) ModbusPDUError {
        if sModbusPDUError, ok := typ.(ModbusPDUError); ok {
            return sModbusPDUError
        }
        if sModbusPDUError, ok := typ.(*ModbusPDUError); ok {
            return *sModbusPDUError
        }
        return ModbusPDUError{}
    }
    return castFunc(structType)
}

func (m ModbusPDUError) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (exceptionCode)
    lengthInBits += 8

    return lengthInBits
}

func (m ModbusPDUError) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUErrorParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Simple Field (exceptionCode)
    exceptionCode, _exceptionCodeErr := io.ReadUint8(8)
    if _exceptionCodeErr != nil {
        return nil, errors.New("Error parsing 'exceptionCode' field " + _exceptionCodeErr.Error())
    }

    // Create the instance
    return NewModbusPDUError(exceptionCode), nil
}

func (m ModbusPDUError) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (exceptionCode)
    exceptionCode := uint8(m.ExceptionCode)
    _exceptionCodeErr := io.WriteUint8(8, (exceptionCode))
    if _exceptionCodeErr != nil {
        return errors.New("Error serializing 'exceptionCode' field " + _exceptionCodeErr.Error())
    }

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}

func (m *ModbusPDUError) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "exceptionCode":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ExceptionCode = data
            }
        }
    }
}

func (m ModbusPDUError) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.modbus.readwrite.ModbusPDUError"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ExceptionCode, xml.StartElement{Name: xml.Name{Local: "exceptionCode"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

