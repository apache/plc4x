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
type S7PayloadUserDataItemCpuFunctionReadSzlRequest struct {
    S7PayloadUserDataItem
}

// The corresponding interface
type IS7PayloadUserDataItemCpuFunctionReadSzlRequest interface {
    IS7PayloadUserDataItem
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) CpuFunctionType() uint8 {
    return 0x04
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) initialize(returnCode IDataTransportErrorCode, transportSize IDataTransportSize, szlId ISzlId, szlIndex uint16) spi.Message {
    m.ReturnCode = returnCode
    m.TransportSize = transportSize
    m.SzlId = szlId
    m.SzlIndex = szlIndex
    return m
}

func NewS7PayloadUserDataItemCpuFunctionReadSzlRequest() S7PayloadUserDataItemInitializer {
    return &S7PayloadUserDataItemCpuFunctionReadSzlRequest{}
}

func CastIS7PayloadUserDataItemCpuFunctionReadSzlRequest(structType interface{}) IS7PayloadUserDataItemCpuFunctionReadSzlRequest {
    castFunc := func(typ interface{}) IS7PayloadUserDataItemCpuFunctionReadSzlRequest {
        if iS7PayloadUserDataItemCpuFunctionReadSzlRequest, ok := typ.(IS7PayloadUserDataItemCpuFunctionReadSzlRequest); ok {
            return iS7PayloadUserDataItemCpuFunctionReadSzlRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastS7PayloadUserDataItemCpuFunctionReadSzlRequest(structType interface{}) S7PayloadUserDataItemCpuFunctionReadSzlRequest {
    castFunc := func(typ interface{}) S7PayloadUserDataItemCpuFunctionReadSzlRequest {
        if sS7PayloadUserDataItemCpuFunctionReadSzlRequest, ok := typ.(S7PayloadUserDataItemCpuFunctionReadSzlRequest); ok {
            return sS7PayloadUserDataItemCpuFunctionReadSzlRequest
        }
        if sS7PayloadUserDataItemCpuFunctionReadSzlRequest, ok := typ.(*S7PayloadUserDataItemCpuFunctionReadSzlRequest); ok {
            return *sS7PayloadUserDataItemCpuFunctionReadSzlRequest
        }
        return S7PayloadUserDataItemCpuFunctionReadSzlRequest{}
    }
    return castFunc(structType)
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.S7PayloadUserDataItem.LengthInBits()

    return lengthInBits
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func S7PayloadUserDataItemCpuFunctionReadSzlRequestParse(io *utils.ReadBuffer) (S7PayloadUserDataItemInitializer, error) {

    // Create the instance
    return NewS7PayloadUserDataItemCpuFunctionReadSzlRequest(), nil
}

func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return S7PayloadUserDataItemSerialize(io, m.S7PayloadUserDataItem, CastIS7PayloadUserDataItem(m), ser)
}

func (m *S7PayloadUserDataItemCpuFunctionReadSzlRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m S7PayloadUserDataItemCpuFunctionReadSzlRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.s7.readwrite.S7PayloadUserDataItemCpuFunctionReadSzlRequest"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

