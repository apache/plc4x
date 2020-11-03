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
    log "github.com/sirupsen/logrus"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type ConnectionRequestInformationTunnelConnection struct {
    KnxLayer IKnxLayer
    ConnectionRequestInformation
}

// The corresponding interface
type IConnectionRequestInformationTunnelConnection interface {
    IConnectionRequestInformation
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ConnectionRequestInformationTunnelConnection) ConnectionType() uint8 {
    return 0x04
}

func (m ConnectionRequestInformationTunnelConnection) initialize() spi.Message {
    return m
}

func NewConnectionRequestInformationTunnelConnection(knxLayer IKnxLayer) ConnectionRequestInformationInitializer {
    return &ConnectionRequestInformationTunnelConnection{KnxLayer: knxLayer}
}

func CastIConnectionRequestInformationTunnelConnection(structType interface{}) IConnectionRequestInformationTunnelConnection {
    castFunc := func(typ interface{}) IConnectionRequestInformationTunnelConnection {
        if iConnectionRequestInformationTunnelConnection, ok := typ.(IConnectionRequestInformationTunnelConnection); ok {
            return iConnectionRequestInformationTunnelConnection
        }
        return nil
    }
    return castFunc(structType)
}

func CastConnectionRequestInformationTunnelConnection(structType interface{}) ConnectionRequestInformationTunnelConnection {
    castFunc := func(typ interface{}) ConnectionRequestInformationTunnelConnection {
        if sConnectionRequestInformationTunnelConnection, ok := typ.(ConnectionRequestInformationTunnelConnection); ok {
            return sConnectionRequestInformationTunnelConnection
        }
        if sConnectionRequestInformationTunnelConnection, ok := typ.(*ConnectionRequestInformationTunnelConnection); ok {
            return *sConnectionRequestInformationTunnelConnection
        }
        return ConnectionRequestInformationTunnelConnection{}
    }
    return castFunc(structType)
}

func (m ConnectionRequestInformationTunnelConnection) LengthInBits() uint16 {
    var lengthInBits uint16 = m.ConnectionRequestInformation.LengthInBits()

    // Enum Field (knxLayer)
    lengthInBits += 8

    // Reserved Field (reserved)
    lengthInBits += 8

    return lengthInBits
}

func (m ConnectionRequestInformationTunnelConnection) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ConnectionRequestInformationTunnelConnectionParse(io *utils.ReadBuffer) (ConnectionRequestInformationInitializer, error) {

    // Enum field (knxLayer)
    knxLayer, _knxLayerErr := KnxLayerParse(io)
    if _knxLayerErr != nil {
        return nil, errors.New("Error parsing 'knxLayer' field " + _knxLayerErr.Error())
    }

    // Reserved Field (Compartmentalized so the "reserved" variable can't leak)
    {
        reserved, _err := io.ReadUint8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
        }
        if reserved != uint8(0x00) {
            log.WithFields(log.Fields{
                "expected value": uint8(0x00),
                "got value": reserved,
            }).Info("Got unexpected response.")
        }
    }

    // Create the instance
    return NewConnectionRequestInformationTunnelConnection(knxLayer), nil
}

func (m ConnectionRequestInformationTunnelConnection) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Enum field (knxLayer)
    knxLayer := CastKnxLayer(m.KnxLayer)
    _knxLayerErr := knxLayer.Serialize(io)
    if _knxLayerErr != nil {
        return errors.New("Error serializing 'knxLayer' field " + _knxLayerErr.Error())
    }

    // Reserved Field (reserved)
    {
        _err := io.WriteUint8(8, uint8(0x00))
        if _err != nil {
            return errors.New("Error serializing 'reserved' field " + _err.Error())
        }
    }

        return nil
    }
    return ConnectionRequestInformationSerialize(io, m.ConnectionRequestInformation, CastIConnectionRequestInformation(m), ser)
}

func (m *ConnectionRequestInformationTunnelConnection) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "knxLayer":
                var data *KnxLayer
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.KnxLayer = data
            }
        }
    }
}

func (m ConnectionRequestInformationTunnelConnection) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.ConnectionRequestInformationTunnelConnection"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.KnxLayer, xml.StartElement{Name: xml.Name{Local: "knxLayer"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

