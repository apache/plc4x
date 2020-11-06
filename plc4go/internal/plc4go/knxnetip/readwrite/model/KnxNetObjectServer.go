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
type KnxNetObjectServer struct {
    Version uint8
    Parent *ServiceId
    IKnxNetObjectServer
}

// The corresponding interface
type IKnxNetObjectServer interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *KnxNetObjectServer) ServiceType() uint8 {
    return 0x08
}


func (m *KnxNetObjectServer) InitializeParent(parent *ServiceId) {
}

func NewKnxNetObjectServer(version uint8, ) *ServiceId {
    child := &KnxNetObjectServer{
        Version: version,
        Parent: NewServiceId(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastKnxNetObjectServer(structType interface{}) KnxNetObjectServer {
    castFunc := func(typ interface{}) KnxNetObjectServer {
        if casted, ok := typ.(KnxNetObjectServer); ok {
            return casted
        }
        if casted, ok := typ.(*KnxNetObjectServer); ok {
            return *casted
        }
        if casted, ok := typ.(ServiceId); ok {
            return CastKnxNetObjectServer(casted.Child)
        }
        if casted, ok := typ.(*ServiceId); ok {
            return CastKnxNetObjectServer(casted.Child)
        }
        return KnxNetObjectServer{}
    }
    return castFunc(structType)
}

func (m *KnxNetObjectServer) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (version)
    lengthInBits += 8

    return lengthInBits
}

func (m *KnxNetObjectServer) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KnxNetObjectServerParse(io *utils.ReadBuffer) (*ServiceId, error) {

    // Simple Field (version)
    version, _versionErr := io.ReadUint8(8)
    if _versionErr != nil {
        return nil, errors.New("Error parsing 'version' field " + _versionErr.Error())
    }

    // Create a partially initialized instance
    _child := &KnxNetObjectServer{
        Version: version,
        Parent: &ServiceId{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *KnxNetObjectServer) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (version)
    version := uint8(m.Version)
    _versionErr := io.WriteUint8(8, (version))
    if _versionErr != nil {
        return errors.New("Error serializing 'version' field " + _versionErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *KnxNetObjectServer) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "version":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.Version = data
            }
        }
    }
}

func (m *KnxNetObjectServer) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.KnxNetObjectServer"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.Version, xml.StartElement{Name: xml.Name{Local: "version"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

