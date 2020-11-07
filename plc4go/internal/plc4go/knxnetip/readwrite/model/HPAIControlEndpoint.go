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
type HPAIControlEndpoint struct {
    HostProtocolCode HostProtocolCode
    IpAddress *IPAddress
    IpPort uint16
    IHPAIControlEndpoint
}

// The corresponding interface
type IHPAIControlEndpoint interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

func NewHPAIControlEndpoint(hostProtocolCode HostProtocolCode, ipAddress *IPAddress, ipPort uint16) *HPAIControlEndpoint {
    return &HPAIControlEndpoint{HostProtocolCode: hostProtocolCode, IpAddress: ipAddress, IpPort: ipPort}
}

func CastHPAIControlEndpoint(structType interface{}) HPAIControlEndpoint {
    castFunc := func(typ interface{}) HPAIControlEndpoint {
        if casted, ok := typ.(HPAIControlEndpoint); ok {
            return casted
        }
        if casted, ok := typ.(*HPAIControlEndpoint); ok {
            return *casted
        }
        return HPAIControlEndpoint{}
    }
    return castFunc(structType)
}

func (m *HPAIControlEndpoint) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Implicit Field (structureLength)
    lengthInBits += 8

    // Enum Field (hostProtocolCode)
    lengthInBits += 8

    // Simple field (ipAddress)
    lengthInBits += m.IpAddress.LengthInBits()

    // Simple field (ipPort)
    lengthInBits += 16

    return lengthInBits
}

func (m *HPAIControlEndpoint) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func HPAIControlEndpointParse(io *utils.ReadBuffer) (*HPAIControlEndpoint, error) {

    // Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    _, _structureLengthErr := io.ReadUint8(8)
    if _structureLengthErr != nil {
        return nil, errors.New("Error parsing 'structureLength' field " + _structureLengthErr.Error())
    }

    // Enum field (hostProtocolCode)
    hostProtocolCode, _hostProtocolCodeErr := HostProtocolCodeParse(io)
    if _hostProtocolCodeErr != nil {
        return nil, errors.New("Error parsing 'hostProtocolCode' field " + _hostProtocolCodeErr.Error())
    }

    // Simple Field (ipAddress)
    ipAddress, _ipAddressErr := IPAddressParse(io)
    if _ipAddressErr != nil {
        return nil, errors.New("Error parsing 'ipAddress' field " + _ipAddressErr.Error())
    }

    // Simple Field (ipPort)
    ipPort, _ipPortErr := io.ReadUint16(16)
    if _ipPortErr != nil {
        return nil, errors.New("Error parsing 'ipPort' field " + _ipPortErr.Error())
    }

    // Create the instance
    return NewHPAIControlEndpoint(hostProtocolCode, ipAddress, ipPort), nil
}

func (m *HPAIControlEndpoint) Serialize(io utils.WriteBuffer) error {

    // Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    structureLength := uint8(uint8(m.LengthInBytes()))
    _structureLengthErr := io.WriteUint8(8, (structureLength))
    if _structureLengthErr != nil {
        return errors.New("Error serializing 'structureLength' field " + _structureLengthErr.Error())
    }

    // Enum field (hostProtocolCode)
    hostProtocolCode := CastHostProtocolCode(m.HostProtocolCode)
    _hostProtocolCodeErr := hostProtocolCode.Serialize(io)
    if _hostProtocolCodeErr != nil {
        return errors.New("Error serializing 'hostProtocolCode' field " + _hostProtocolCodeErr.Error())
    }

    // Simple Field (ipAddress)
    _ipAddressErr := m.IpAddress.Serialize(io)
    if _ipAddressErr != nil {
        return errors.New("Error serializing 'ipAddress' field " + _ipAddressErr.Error())
    }

    // Simple Field (ipPort)
    ipPort := uint16(m.IpPort)
    _ipPortErr := io.WriteUint16(16, (ipPort))
    if _ipPortErr != nil {
        return errors.New("Error serializing 'ipPort' field " + _ipPortErr.Error())
    }

    return nil
}

func (m *HPAIControlEndpoint) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    for {
        token, err = d.Token()
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
            case "hostProtocolCode":
                var data HostProtocolCode
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.HostProtocolCode = data
            case "ipAddress":
                var data *IPAddress
                if err := d.DecodeElement(data, &tok); err != nil {
                    return err
                }
                m.IpAddress = data
            case "ipPort":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.IpPort = data
            }
        }
    }
}

func (m *HPAIControlEndpoint) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := "org.apache.plc4x.java.knxnetip.readwrite.HPAIControlEndpoint"
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.HostProtocolCode, xml.StartElement{Name: xml.Name{Local: "hostProtocolCode"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.IpAddress, xml.StartElement{Name: xml.Name{Local: "ipAddress"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.IpPort, xml.StartElement{Name: xml.Name{Local: "ipPort"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

