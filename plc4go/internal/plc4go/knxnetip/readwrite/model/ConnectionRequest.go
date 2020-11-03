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
    "reflect"
)

// The data-structure of this message
type ConnectionRequest struct {
    HpaiDiscoveryEndpoint IHPAIDiscoveryEndpoint
    HpaiDataEndpoint IHPAIDataEndpoint
    ConnectionRequestInformation IConnectionRequestInformation
    KNXNetIPMessage
}

// The corresponding interface
type IConnectionRequest interface {
    IKNXNetIPMessage
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ConnectionRequest) MsgType() uint16 {
    return 0x0205
}

func (m ConnectionRequest) initialize() spi.Message {
    return m
}

func NewConnectionRequest(hpaiDiscoveryEndpoint IHPAIDiscoveryEndpoint, hpaiDataEndpoint IHPAIDataEndpoint, connectionRequestInformation IConnectionRequestInformation) KNXNetIPMessageInitializer {
    return &ConnectionRequest{HpaiDiscoveryEndpoint: hpaiDiscoveryEndpoint, HpaiDataEndpoint: hpaiDataEndpoint, ConnectionRequestInformation: connectionRequestInformation}
}

func CastIConnectionRequest(structType interface{}) IConnectionRequest {
    castFunc := func(typ interface{}) IConnectionRequest {
        if iConnectionRequest, ok := typ.(IConnectionRequest); ok {
            return iConnectionRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastConnectionRequest(structType interface{}) ConnectionRequest {
    castFunc := func(typ interface{}) ConnectionRequest {
        if sConnectionRequest, ok := typ.(ConnectionRequest); ok {
            return sConnectionRequest
        }
        if sConnectionRequest, ok := typ.(*ConnectionRequest); ok {
            return *sConnectionRequest
        }
        return ConnectionRequest{}
    }
    return castFunc(structType)
}

func (m ConnectionRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

    // Simple field (hpaiDiscoveryEndpoint)
    lengthInBits += m.HpaiDiscoveryEndpoint.LengthInBits()

    // Simple field (hpaiDataEndpoint)
    lengthInBits += m.HpaiDataEndpoint.LengthInBits()

    // Simple field (connectionRequestInformation)
    lengthInBits += m.ConnectionRequestInformation.LengthInBits()

    return lengthInBits
}

func (m ConnectionRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ConnectionRequestParse(io *utils.ReadBuffer) (KNXNetIPMessageInitializer, error) {

    // Simple Field (hpaiDiscoveryEndpoint)
    _hpaiDiscoveryEndpointMessage, _err := HPAIDiscoveryEndpointParse(io)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'hpaiDiscoveryEndpoint'. " + _err.Error())
    }
    var hpaiDiscoveryEndpoint IHPAIDiscoveryEndpoint
    hpaiDiscoveryEndpoint, _hpaiDiscoveryEndpointOk := _hpaiDiscoveryEndpointMessage.(IHPAIDiscoveryEndpoint)
    if !_hpaiDiscoveryEndpointOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_hpaiDiscoveryEndpointMessage).Name() + " to IHPAIDiscoveryEndpoint")
    }

    // Simple Field (hpaiDataEndpoint)
    _hpaiDataEndpointMessage, _err := HPAIDataEndpointParse(io)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'hpaiDataEndpoint'. " + _err.Error())
    }
    var hpaiDataEndpoint IHPAIDataEndpoint
    hpaiDataEndpoint, _hpaiDataEndpointOk := _hpaiDataEndpointMessage.(IHPAIDataEndpoint)
    if !_hpaiDataEndpointOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_hpaiDataEndpointMessage).Name() + " to IHPAIDataEndpoint")
    }

    // Simple Field (connectionRequestInformation)
    _connectionRequestInformationMessage, _err := ConnectionRequestInformationParse(io)
    if _err != nil {
        return nil, errors.New("Error parsing simple field 'connectionRequestInformation'. " + _err.Error())
    }
    var connectionRequestInformation IConnectionRequestInformation
    connectionRequestInformation, _connectionRequestInformationOk := _connectionRequestInformationMessage.(IConnectionRequestInformation)
    if !_connectionRequestInformationOk {
        return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_connectionRequestInformationMessage).Name() + " to IConnectionRequestInformation")
    }

    // Create the instance
    return NewConnectionRequest(hpaiDiscoveryEndpoint, hpaiDataEndpoint, connectionRequestInformation), nil
}

func (m ConnectionRequest) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (hpaiDiscoveryEndpoint)
    hpaiDiscoveryEndpoint := CastIHPAIDiscoveryEndpoint(m.HpaiDiscoveryEndpoint)
    _hpaiDiscoveryEndpointErr := hpaiDiscoveryEndpoint.Serialize(io)
    if _hpaiDiscoveryEndpointErr != nil {
        return errors.New("Error serializing 'hpaiDiscoveryEndpoint' field " + _hpaiDiscoveryEndpointErr.Error())
    }

    // Simple Field (hpaiDataEndpoint)
    hpaiDataEndpoint := CastIHPAIDataEndpoint(m.HpaiDataEndpoint)
    _hpaiDataEndpointErr := hpaiDataEndpoint.Serialize(io)
    if _hpaiDataEndpointErr != nil {
        return errors.New("Error serializing 'hpaiDataEndpoint' field " + _hpaiDataEndpointErr.Error())
    }

    // Simple Field (connectionRequestInformation)
    connectionRequestInformation := CastIConnectionRequestInformation(m.ConnectionRequestInformation)
    _connectionRequestInformationErr := connectionRequestInformation.Serialize(io)
    if _connectionRequestInformationErr != nil {
        return errors.New("Error serializing 'connectionRequestInformation' field " + _connectionRequestInformationErr.Error())
    }

        return nil
    }
    return KNXNetIPMessageSerialize(io, m.KNXNetIPMessage, CastIKNXNetIPMessage(m), ser)
}

func (m *ConnectionRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "hpaiDiscoveryEndpoint":
                var data *HPAIDiscoveryEndpoint
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.HpaiDiscoveryEndpoint = CastIHPAIDiscoveryEndpoint(data)
            case "hpaiDataEndpoint":
                var data *HPAIDataEndpoint
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.HpaiDataEndpoint = CastIHPAIDataEndpoint(data)
            case "connectionRequestInformation":
                switch tok.Attr[0].Value {
                    case "org.apache.plc4x.java.knxnetip.readwrite.ConnectionRequestInformationDeviceManagement":
                        var dt *ConnectionRequestInformationDeviceManagement
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.ConnectionRequestInformation = dt
                    case "org.apache.plc4x.java.knxnetip.readwrite.ConnectionRequestInformationTunnelConnection":
                        var dt *ConnectionRequestInformationTunnelConnection
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        m.ConnectionRequestInformation = dt
                    }
            }
        }
    }
}

func (m ConnectionRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.ConnectionRequest"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.HpaiDiscoveryEndpoint, xml.StartElement{Name: xml.Name{Local: "hpaiDiscoveryEndpoint"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.HpaiDataEndpoint, xml.StartElement{Name: xml.Name{Local: "hpaiDataEndpoint"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.ConnectionRequestInformation, xml.StartElement{Name: xml.Name{Local: "connectionRequestInformation"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

