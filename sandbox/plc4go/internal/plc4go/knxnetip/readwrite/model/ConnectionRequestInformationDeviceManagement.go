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
type ConnectionRequestInformationDeviceManagement struct {
    ConnectionRequestInformation
}

// The corresponding interface
type IConnectionRequestInformationDeviceManagement interface {
    IConnectionRequestInformation
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ConnectionRequestInformationDeviceManagement) ConnectionType() uint8 {
    return 0x03
}

func (m ConnectionRequestInformationDeviceManagement) initialize() spi.Message {
    return m
}

func NewConnectionRequestInformationDeviceManagement() ConnectionRequestInformationInitializer {
    return &ConnectionRequestInformationDeviceManagement{}
}

func CastIConnectionRequestInformationDeviceManagement(structType interface{}) IConnectionRequestInformationDeviceManagement {
    castFunc := func(typ interface{}) IConnectionRequestInformationDeviceManagement {
        if iConnectionRequestInformationDeviceManagement, ok := typ.(IConnectionRequestInformationDeviceManagement); ok {
            return iConnectionRequestInformationDeviceManagement
        }
        return nil
    }
    return castFunc(structType)
}

func CastConnectionRequestInformationDeviceManagement(structType interface{}) ConnectionRequestInformationDeviceManagement {
    castFunc := func(typ interface{}) ConnectionRequestInformationDeviceManagement {
        if sConnectionRequestInformationDeviceManagement, ok := typ.(ConnectionRequestInformationDeviceManagement); ok {
            return sConnectionRequestInformationDeviceManagement
        }
        if sConnectionRequestInformationDeviceManagement, ok := typ.(*ConnectionRequestInformationDeviceManagement); ok {
            return *sConnectionRequestInformationDeviceManagement
        }
        return ConnectionRequestInformationDeviceManagement{}
    }
    return castFunc(structType)
}

func (m ConnectionRequestInformationDeviceManagement) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ConnectionRequestInformation.LengthInBits()

    return lengthInBits
}

func (m ConnectionRequestInformationDeviceManagement) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ConnectionRequestInformationDeviceManagementParse(io *utils.ReadBuffer) (ConnectionRequestInformationInitializer, error) {

    // Create the instance
    return NewConnectionRequestInformationDeviceManagement(), nil
}

func (m ConnectionRequestInformationDeviceManagement) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return ConnectionRequestInformationSerialize(io, m.ConnectionRequestInformation, CastIConnectionRequestInformation(m), ser)
}

func (m *ConnectionRequestInformationDeviceManagement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m ConnectionRequestInformationDeviceManagement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.ConnectionRequestInformationDeviceManagement"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

