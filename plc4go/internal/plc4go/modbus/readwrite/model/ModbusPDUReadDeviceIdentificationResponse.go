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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ModbusPDUReadDeviceIdentificationResponse struct {
    ModbusPDU
}

// The corresponding interface
type IModbusPDUReadDeviceIdentificationResponse interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUReadDeviceIdentificationResponse) ErrorFlag() bool {
    return false
}

func (m ModbusPDUReadDeviceIdentificationResponse) FunctionFlag() uint8 {
    return 0x2B
}

func (m ModbusPDUReadDeviceIdentificationResponse) Response() bool {
    return true
}

func (m ModbusPDUReadDeviceIdentificationResponse) initialize() spi.Message {
    return m
}

func NewModbusPDUReadDeviceIdentificationResponse() ModbusPDUInitializer {
    return &ModbusPDUReadDeviceIdentificationResponse{}
}

func CastIModbusPDUReadDeviceIdentificationResponse(structType interface{}) IModbusPDUReadDeviceIdentificationResponse {
    castFunc := func(typ interface{}) IModbusPDUReadDeviceIdentificationResponse {
        if iModbusPDUReadDeviceIdentificationResponse, ok := typ.(IModbusPDUReadDeviceIdentificationResponse); ok {
            return iModbusPDUReadDeviceIdentificationResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUReadDeviceIdentificationResponse(structType interface{}) ModbusPDUReadDeviceIdentificationResponse {
    castFunc := func(typ interface{}) ModbusPDUReadDeviceIdentificationResponse {
        if sModbusPDUReadDeviceIdentificationResponse, ok := typ.(ModbusPDUReadDeviceIdentificationResponse); ok {
            return sModbusPDUReadDeviceIdentificationResponse
        }
        if sModbusPDUReadDeviceIdentificationResponse, ok := typ.(*ModbusPDUReadDeviceIdentificationResponse); ok {
            return *sModbusPDUReadDeviceIdentificationResponse
        }
        return ModbusPDUReadDeviceIdentificationResponse{}
    }
    return castFunc(structType)
}

func (m ModbusPDUReadDeviceIdentificationResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    return lengthInBits
}

func (m ModbusPDUReadDeviceIdentificationResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadDeviceIdentificationResponseParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Create the instance
    return NewModbusPDUReadDeviceIdentificationResponse(), nil
}

func (m ModbusPDUReadDeviceIdentificationResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}

func (m *ModbusPDUReadDeviceIdentificationResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m ModbusPDUReadDeviceIdentificationResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadDeviceIdentificationResponse"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

