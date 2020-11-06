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
type KNXGroupAddressFreeLevel struct {
    SubGroup uint16
    Parent *KNXGroupAddress
    IKNXGroupAddressFreeLevel
}

// The corresponding interface
type IKNXGroupAddressFreeLevel interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *KNXGroupAddressFreeLevel) NumLevels() uint8 {
    return 1
}


func (m *KNXGroupAddressFreeLevel) InitializeParent(parent *KNXGroupAddress) {
}

func NewKNXGroupAddressFreeLevel(subGroup uint16, ) *KNXGroupAddress {
    child := &KNXGroupAddressFreeLevel{
        SubGroup: subGroup,
        Parent: NewKNXGroupAddress(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastKNXGroupAddressFreeLevel(structType interface{}) KNXGroupAddressFreeLevel {
    castFunc := func(typ interface{}) KNXGroupAddressFreeLevel {
        if casted, ok := typ.(KNXGroupAddressFreeLevel); ok {
            return casted
        }
        if casted, ok := typ.(*KNXGroupAddressFreeLevel); ok {
            return *casted
        }
        if casted, ok := typ.(KNXGroupAddress); ok {
            return CastKNXGroupAddressFreeLevel(casted.Child)
        }
        if casted, ok := typ.(*KNXGroupAddress); ok {
            return CastKNXGroupAddressFreeLevel(casted.Child)
        }
        return KNXGroupAddressFreeLevel{}
    }
    return castFunc(structType)
}

func (m *KNXGroupAddressFreeLevel) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (subGroup)
    lengthInBits += 16

    return lengthInBits
}

func (m *KNXGroupAddressFreeLevel) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KNXGroupAddressFreeLevelParse(io *utils.ReadBuffer) (*KNXGroupAddress, error) {

    // Simple Field (subGroup)
    subGroup, _subGroupErr := io.ReadUint16(16)
    if _subGroupErr != nil {
        return nil, errors.New("Error parsing 'subGroup' field " + _subGroupErr.Error())
    }

    // Create a partially initialized instance
    _child := &KNXGroupAddressFreeLevel{
        SubGroup: subGroup,
        Parent: &KNXGroupAddress{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *KNXGroupAddressFreeLevel) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (subGroup)
    subGroup := uint16(m.SubGroup)
    _subGroupErr := io.WriteUint16(16, (subGroup))
    if _subGroupErr != nil {
        return errors.New("Error serializing 'subGroup' field " + _subGroupErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *KNXGroupAddressFreeLevel) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "subGroup":
                var data uint16
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SubGroup = data
            }
        }
    }
}

func (m *KNXGroupAddressFreeLevel) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.KNXGroupAddressFreeLevel"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SubGroup, xml.StartElement{Name: xml.Name{Local: "subGroup"}}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

