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
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
    "reflect"
    "strings"
)

// The data-structure of this message
type ConnectionRequestInformation struct {
    Child IConnectionRequestInformationChild
    IConnectionRequestInformation
    IConnectionRequestInformationParent
}

// The corresponding interface
type IConnectionRequestInformation interface {
    ConnectionType() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

type IConnectionRequestInformationParent interface {
    SerializeParent(io utils.WriteBuffer, child IConnectionRequestInformation, serializeChildFunction func() error) error
}

type IConnectionRequestInformationChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *ConnectionRequestInformation)
    IConnectionRequestInformation
}

func NewConnectionRequestInformation() *ConnectionRequestInformation {
    return &ConnectionRequestInformation{}
}

func CastConnectionRequestInformation(structType interface{}) *ConnectionRequestInformation {
    castFunc := func(typ interface{}) *ConnectionRequestInformation {
        if casted, ok := typ.(ConnectionRequestInformation); ok {
            return &casted
        }
        if casted, ok := typ.(*ConnectionRequestInformation); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *ConnectionRequestInformation) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Implicit Field (structureLength)
    lengthInBits += 8

    // Discriminator Field (connectionType)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *ConnectionRequestInformation) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ConnectionRequestInformationParse(io *utils.ReadBuffer) (*ConnectionRequestInformation, error) {

    // Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    _, _structureLengthErr := io.ReadUint8(8)
    if _structureLengthErr != nil {
        return nil, errors.New("Error parsing 'structureLength' field " + _structureLengthErr.Error())
    }

    // Discriminator Field (connectionType) (Used as input to a switch field)
    connectionType, _connectionTypeErr := io.ReadUint8(8)
    if _connectionTypeErr != nil {
        return nil, errors.New("Error parsing 'connectionType' field " + _connectionTypeErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *ConnectionRequestInformation
    var typeSwitchError error
    switch {
    case connectionType == 0x03:
        _parent, typeSwitchError = ConnectionRequestInformationDeviceManagementParse(io)
    case connectionType == 0x04:
        _parent, typeSwitchError = ConnectionRequestInformationTunnelConnectionParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *ConnectionRequestInformation) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *ConnectionRequestInformation) SerializeParent(io utils.WriteBuffer, child IConnectionRequestInformation, serializeChildFunction func() error) error {

    // Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    structureLength := uint8(uint8(m.LengthInBytes()))
    _structureLengthErr := io.WriteUint8(8, (structureLength))
    if _structureLengthErr != nil {
        return errors.New("Error serializing 'structureLength' field " + _structureLengthErr.Error())
    }

    // Discriminator Field (connectionType) (Used as input to a switch field)
    connectionType := uint8(child.ConnectionType())
    _connectionTypeErr := io.WriteUint8(8, (connectionType))
    if _connectionTypeErr != nil {
        return errors.New("Error serializing 'connectionType' field " + _connectionTypeErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *ConnectionRequestInformation) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
                default:
                    switch start.Attr[0].Value {
                        case "org.apache.plc4x.java.knxnetip.readwrite.ConnectionRequestInformationDeviceManagement":
                            var dt *ConnectionRequestInformationDeviceManagement
                            if m.Child != nil {
                                dt = m.Child.(*ConnectionRequestInformationDeviceManagement)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.knxnetip.readwrite.ConnectionRequestInformationTunnelConnection":
                            var dt *ConnectionRequestInformationTunnelConnection
                            if m.Child != nil {
                                dt = m.Child.(*ConnectionRequestInformationTunnelConnection)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                    }
            }
        }
    }
}

func (m *ConnectionRequestInformation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := reflect.TypeOf(m.Child).String()
    className = "org.apache.plc4x.java.knxnetip.readwrite." + className[strings.LastIndex(className, ".") + 1:]
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    marshaller, ok := m.Child.(xml.Marshaler)
    if !ok {
        return errors.New("child is not castable to Marshaler")
    }
    marshaller.MarshalXML(e, start)
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

