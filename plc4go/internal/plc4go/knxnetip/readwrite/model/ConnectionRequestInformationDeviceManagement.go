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
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
)

// The data-structure of this message
type ConnectionRequestInformationDeviceManagement struct {
    Parent *ConnectionRequestInformation
    IConnectionRequestInformationDeviceManagement
}

// The corresponding interface
type IConnectionRequestInformationDeviceManagement interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *ConnectionRequestInformationDeviceManagement) ConnectionType() uint8 {
    return 0x03
}


func (m *ConnectionRequestInformationDeviceManagement) InitializeParent(parent *ConnectionRequestInformation) {
}

func NewConnectionRequestInformationDeviceManagement() *ConnectionRequestInformation {
    child := &ConnectionRequestInformationDeviceManagement{
        Parent: NewConnectionRequestInformation(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastConnectionRequestInformationDeviceManagement(structType interface{}) *ConnectionRequestInformationDeviceManagement {
    castFunc := func(typ interface{}) *ConnectionRequestInformationDeviceManagement {
        if casted, ok := typ.(ConnectionRequestInformationDeviceManagement); ok {
            return &casted
        }
        if casted, ok := typ.(*ConnectionRequestInformationDeviceManagement); ok {
            return casted
        }
        if casted, ok := typ.(ConnectionRequestInformation); ok {
            return CastConnectionRequestInformationDeviceManagement(casted.Child)
        }
        if casted, ok := typ.(*ConnectionRequestInformation); ok {
            return CastConnectionRequestInformationDeviceManagement(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *ConnectionRequestInformationDeviceManagement) GetTypeName() string {
    return "ConnectionRequestInformationDeviceManagement"
}

func (m *ConnectionRequestInformationDeviceManagement) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    return lengthInBits
}

func (m *ConnectionRequestInformationDeviceManagement) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func ConnectionRequestInformationDeviceManagementParse(io *utils.ReadBuffer) (*ConnectionRequestInformation, error) {

    // Create a partially initialized instance
    _child := &ConnectionRequestInformationDeviceManagement{
        Parent: &ConnectionRequestInformation{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *ConnectionRequestInformationDeviceManagement) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *ConnectionRequestInformationDeviceManagement) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
        switch token.(type) {
        case xml.StartElement:
            tok := token.(xml.StartElement)
            switch tok.Name.Local {
            }
        }
        token, err = d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
    }
}

func (m *ConnectionRequestInformationDeviceManagement) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    return nil
}

