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
type ModbusPDUReadFifoQueueRequest struct {
    FifoPointerAddress uint16
    ModbusPDU
}

// The corresponding interface
type IModbusPDUReadFifoQueueRequest interface {
    IModbusPDU
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUReadFifoQueueRequest) ErrorFlag() bool {
    return false
}

func (m ModbusPDUReadFifoQueueRequest) FunctionFlag() uint8 {
    return 0x18
}

func (m ModbusPDUReadFifoQueueRequest) Response() bool {
    return false
}

func (m ModbusPDUReadFifoQueueRequest) initialize() spi.Message {
    return m
}

func NewModbusPDUReadFifoQueueRequest(fifoPointerAddress uint16) ModbusPDUInitializer {
    return &ModbusPDUReadFifoQueueRequest{FifoPointerAddress: fifoPointerAddress}
}

func CastIModbusPDUReadFifoQueueRequest(structType interface{}) IModbusPDUReadFifoQueueRequest {
    castFunc := func(typ interface{}) IModbusPDUReadFifoQueueRequest {
        if iModbusPDUReadFifoQueueRequest, ok := typ.(IModbusPDUReadFifoQueueRequest); ok {
            return iModbusPDUReadFifoQueueRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastModbusPDUReadFifoQueueRequest(structType interface{}) ModbusPDUReadFifoQueueRequest {
    castFunc := func(typ interface{}) ModbusPDUReadFifoQueueRequest {
        if sModbusPDUReadFifoQueueRequest, ok := typ.(ModbusPDUReadFifoQueueRequest); ok {
            return sModbusPDUReadFifoQueueRequest
        }
        if sModbusPDUReadFifoQueueRequest, ok := typ.(*ModbusPDUReadFifoQueueRequest); ok {
            return *sModbusPDUReadFifoQueueRequest
        }
        return ModbusPDUReadFifoQueueRequest{}
    }
    return castFunc(structType)
}

func (m ModbusPDUReadFifoQueueRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

    // Simple field (fifoPointerAddress)
    lengthInBits += 16

    return lengthInBits
}

func (m ModbusPDUReadFifoQueueRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ModbusPDUReadFifoQueueRequestParse(io *utils.ReadBuffer) (ModbusPDUInitializer, error) {

    // Simple Field (fifoPointerAddress)
    fifoPointerAddress, _fifoPointerAddressErr := io.ReadUint16(16)
    if _fifoPointerAddressErr != nil {
        return nil, errors.New("Error parsing 'fifoPointerAddress' field " + _fifoPointerAddressErr.Error())
    }

    // Create the instance
    return NewModbusPDUReadFifoQueueRequest(fifoPointerAddress), nil
}

func (m ModbusPDUReadFifoQueueRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (fifoPointerAddress)
    fifoPointerAddress := uint16(m.FifoPointerAddress)
    _fifoPointerAddressErr := io.WriteUint16(16, (fifoPointerAddress))
    if _fifoPointerAddressErr != nil {
        return errors.New("Error serializing 'fifoPointerAddress' field " + _fifoPointerAddressErr.Error())
    }

        return nil
    }
    return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}

func (m *ModbusPDUReadFifoQueueRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "fifoPointerAddress":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.FifoPointerAddress = data
            }
        }
    }
}

func (m ModbusPDUReadFifoQueueRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.modbus.readwrite.ModbusPDUReadFifoQueueRequest"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.FifoPointerAddress, xml.StartElement{Name: xml.Name{Local: "fifoPointerAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

