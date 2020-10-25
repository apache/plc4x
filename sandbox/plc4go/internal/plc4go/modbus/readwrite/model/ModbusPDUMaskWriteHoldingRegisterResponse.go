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
type ModbusPDUMaskWriteHoldingRegisterResponse struct {
    ReferenceAddress uint16
    AndMask uint16
    OrMask uint16
    ModbusPDU
}

// The corresponding interface
type IModbusPDUMaskWriteHoldingRegisterResponse interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUMaskWriteHoldingRegisterResponse) ErrorFlag() bool {
    return false
}

func (m ModbusPDUMaskWriteHoldingRegisterResponse) FunctionFlag() uint8 {
    return 0x16
}

func (m ModbusPDUMaskWriteHoldingRegisterResponse) Response() bool {
    return true
}

func (m ModbusPDUMaskWriteHoldingRegisterResponse) initialize() spi.Message {
    return m
}

func NewModbusPDUMaskWriteHoldingRegisterResponse(referenceAddress uint16, andMask uint16, orMask uint16) ModbusPDUInitializer {
    return &ModbusPDUMaskWriteHoldingRegisterResponse{ReferenceAddress: referenceAddress, AndMask: andMask, OrMask: orMask}
}

func CastIModbusPDUMaskWriteHoldingRegisterResponse(structType interface{}) IModbusPDUMaskWriteHoldingRegisterResponse {
    castFunc := func(typ interface{}) IModbusPDUMaskWriteHoldingRegisterResponse {
        if iModbusPDUMaskWriteHoldingRegisterResponse, ok := typ.(IModbusPDUMaskWriteHoldingRegisterResponse); ok {
            return iModbusPDUMaskWriteHoldingRegisterResponse
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUMaskWriteHoldingRegisterResponse(structType interface{}) ModbusPDUMaskWriteHoldingRegisterResponse {
    castFunc := func(typ interface{}) ModbusPDUMaskWriteHoldingRegisterResponse {
        if sModbusPDUMaskWriteHoldingRegisterResponse, ok := typ.(ModbusPDUMaskWriteHoldingRegisterResponse); ok {
            return sModbusPDUMaskWriteHoldingRegisterResponse
        }
        if sModbusPDUMaskWriteHoldingRegisterResponse, ok := typ.(*ModbusPDUMaskWriteHoldingRegisterResponse); ok {
            return *sModbusPDUMaskWriteHoldingRegisterResponse
        }
        return ModbusPDUMaskWriteHoldingRegisterResponse{}
    }
    return castFunc(structType)
}

func (m ModbusPDUMaskWriteHoldingRegisterResponse) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (referenceAddress)
    lengthInBits += 16

    // Simple field (andMask)
    lengthInBits += 16

    // Simple field (orMask)
    lengthInBits += 16

    return lengthInBits
}

func (m ModbusPDUMaskWriteHoldingRegisterResponse) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUMaskWriteHoldingRegisterResponseParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Simple Field (referenceAddress)
    referenceAddress, _referenceAddressErr := io.ReadUint16(16)
    if _referenceAddressErr != nil {
        return nil, errors.New("Error parsing 'referenceAddress' field " + _referenceAddressErr.Error())
    }

    // Simple Field (andMask)
    andMask, _andMaskErr := io.ReadUint16(16)
    if _andMaskErr != nil {
        return nil, errors.New("Error parsing 'andMask' field " + _andMaskErr.Error())
    }

    // Simple Field (orMask)
    orMask, _orMaskErr := io.ReadUint16(16)
    if _orMaskErr != nil {
        return nil, errors.New("Error parsing 'orMask' field " + _orMaskErr.Error())
    }

    // Create the instance
    return NewModbusPDUMaskWriteHoldingRegisterResponse(referenceAddress, andMask, orMask), nil
}

func (m ModbusPDUMaskWriteHoldingRegisterResponse) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (referenceAddress)
    referenceAddress := uint16(m.ReferenceAddress)
    _referenceAddressErr := io.WriteUint16(16, (referenceAddress))
    if _referenceAddressErr != nil {
        return errors.New("Error serializing 'referenceAddress' field " + _referenceAddressErr.Error())
    }

    // Simple Field (andMask)
    andMask := uint16(m.AndMask)
    _andMaskErr := io.WriteUint16(16, (andMask))
    if _andMaskErr != nil {
        return errors.New("Error serializing 'andMask' field " + _andMaskErr.Error())
    }

    // Simple Field (orMask)
    orMask := uint16(m.OrMask)
    _orMaskErr := io.WriteUint16(16, (orMask))
    if _orMaskErr != nil {
        return errors.New("Error serializing 'orMask' field " + _orMaskErr.Error())
    }

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}

func (m *ModbusPDUMaskWriteHoldingRegisterResponse) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "referenceAddress":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.ReferenceAddress = data
            case "andMask":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.AndMask = data
            case "orMask":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.OrMask = data
            }
        }
    }
}

func (m ModbusPDUMaskWriteHoldingRegisterResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.modbus.readwrite.ModbusPDUMaskWriteHoldingRegisterResponse"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ReferenceAddress, xml.StartElement{Name: xml.Name{Local: "referenceAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.AndMask, xml.StartElement{Name: xml.Name{Local: "andMask"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.OrMask, xml.StartElement{Name: xml.Name{Local: "orMask"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

