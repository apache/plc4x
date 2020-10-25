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
type ModbusPDUWriteSingleRegisterResponse struct {
    Address uint16
    Value uint16
    ModbusPDU
}

// The corresponding interface
type IModbusPDUWriteSingleRegisterResponse interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUWriteSingleRegisterResponse) ErrorFlag() bool {
    return false
}

func (m ModbusPDUWriteSingleRegisterResponse) FunctionFlag() uint8 {
    return 0x06
}

func (m ModbusPDUWriteSingleRegisterResponse) Response() bool {
    return true
}

func (m ModbusPDUWriteSingleRegisterResponse) initialize() spi.Message {
    return m
}

func NewModbusPDUWriteSingleRegisterResponse(address uint16, value uint16) ModbusPDUInitializer {
    return &ModbusPDUWriteSingleRegisterResponse{Address: address, Value: value}
}

func CastIModbusPDUWriteSingleRegisterResponse(structType interface{}) IModbusPDUWriteSingleRegisterResponse {
    castFunc := func(typ interface{}) IModbusPDUWriteSingleRegisterResponse {
        if iModbusPDUWriteSingleRegisterResponse, ok := typ.(IModbusPDUWriteSingleRegisterResponse); ok {
            return iModbusPDUWriteSingleRegisterResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUWriteSingleRegisterResponse(structType interface{}) ModbusPDUWriteSingleRegisterResponse {
    castFunc := func(typ interface{}) ModbusPDUWriteSingleRegisterResponse {
        if sModbusPDUWriteSingleRegisterResponse, ok := typ.(ModbusPDUWriteSingleRegisterResponse); ok {
            return sModbusPDUWriteSingleRegisterResponse
        }
        if sModbusPDUWriteSingleRegisterResponse, ok := typ.(*ModbusPDUWriteSingleRegisterResponse); ok {
            return *sModbusPDUWriteSingleRegisterResponse
        }
        return ModbusPDUWriteSingleRegisterResponse{}
    }
    return castFunc(structType)
}

func (m ModbusPDUWriteSingleRegisterResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (address)
    lengthInBits += 16

    // Simple field (value)
    lengthInBits += 16

    return lengthInBits
}

func (m ModbusPDUWriteSingleRegisterResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUWriteSingleRegisterResponseParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Simple Field (address)
    address, _addressErr := io.ReadUint16(16)
    if _addressErr != nil {
        return nil, errors.New("Error parsing 'address' field " + _addressErr.Error())
    }

    // Simple Field (value)
    value, _valueErr := io.ReadUint16(16)
    if _valueErr != nil {
        return nil, errors.New("Error parsing 'value' field " + _valueErr.Error())
    }

    // Create the instance
    return NewModbusPDUWriteSingleRegisterResponse(address, value), nil
}

func (m ModbusPDUWriteSingleRegisterResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (address)
    address := uint16(m.Address)
    _addressErr := io.WriteUint16(16, (address))
    if _addressErr != nil {
        return errors.New("Error serializing 'address' field " + _addressErr.Error())
    }

    // Simple Field (value)
    value := uint16(m.Value)
    _valueErr := io.WriteUint16(16, (value))
    if _valueErr != nil {
        return errors.New("Error serializing 'value' field " + _valueErr.Error())
    }

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}

func (m *ModbusPDUWriteSingleRegisterResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "address":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Address = data
            case "value":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Value = data
            }
        }
    }
}

func (m ModbusPDUWriteSingleRegisterResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.modbus.readwrite.ModbusPDUWriteSingleRegisterResponse"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Address, xml.StartElement{Name: xml.Name{Local: "address"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Value, xml.StartElement{Name: xml.Name{Local: "value"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

