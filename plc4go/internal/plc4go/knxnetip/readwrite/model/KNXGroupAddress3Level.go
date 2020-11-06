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
type KNXGroupAddress3Level struct {
    MainGroup uint8
    MiddleGroup uint8
    SubGroup uint8
    Parent *KNXGroupAddress
    IKNXGroupAddress3Level
}

// The corresponding interface
type IKNXGroupAddress3Level interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *KNXGroupAddress3Level) NumLevels() uint8 {
    return 3
}


func (m *KNXGroupAddress3Level) InitializeParent(parent *KNXGroupAddress) {
}

func NewKNXGroupAddress3Level(mainGroup uint8, middleGroup uint8, subGroup uint8, ) *KNXGroupAddress {
    child := &KNXGroupAddress3Level{
        MainGroup: mainGroup,
        MiddleGroup: middleGroup,
        SubGroup: subGroup,
        Parent: NewKNXGroupAddress(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastKNXGroupAddress3Level(structType interface{}) KNXGroupAddress3Level {
    castFunc := func(typ interface{}) KNXGroupAddress3Level {
        if casted, ok := typ.(KNXGroupAddress3Level); ok {
            return casted
        }
        if casted, ok := typ.(*KNXGroupAddress3Level); ok {
            return *casted
        }
        if casted, ok := typ.(KNXGroupAddress); ok {
            return CastKNXGroupAddress3Level(casted.Child)
        }
        if casted, ok := typ.(*KNXGroupAddress); ok {
            return CastKNXGroupAddress3Level(casted.Child)
        }
        return KNXGroupAddress3Level{}
    }
    return castFunc(structType)
}

func (m *KNXGroupAddress3Level) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (mainGroup)
    lengthInBits += 5

    // Simple field (middleGroup)
    lengthInBits += 3

    // Simple field (subGroup)
    lengthInBits += 8

    return lengthInBits
}

func (m *KNXGroupAddress3Level) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KNXGroupAddress3LevelParse(io *utils.ReadBuffer) (*KNXGroupAddress, error) {

    // Simple Field (mainGroup)
    mainGroup, _mainGroupErr := io.ReadUint8(5)
    if _mainGroupErr != nil {
        return nil, errors.New("Error parsing 'mainGroup' field " + _mainGroupErr.Error())
    }

    // Simple Field (middleGroup)
    middleGroup, _middleGroupErr := io.ReadUint8(3)
    if _middleGroupErr != nil {
        return nil, errors.New("Error parsing 'middleGroup' field " + _middleGroupErr.Error())
    }

    // Simple Field (subGroup)
    subGroup, _subGroupErr := io.ReadUint8(8)
    if _subGroupErr != nil {
        return nil, errors.New("Error parsing 'subGroup' field " + _subGroupErr.Error())
    }

    // Create a partially initialized instance
    _child := &KNXGroupAddress3Level{
        MainGroup: mainGroup,
        MiddleGroup: middleGroup,
        SubGroup: subGroup,
        Parent: &KNXGroupAddress{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *KNXGroupAddress3Level) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (mainGroup)
    mainGroup := uint8(m.MainGroup)
    _mainGroupErr := io.WriteUint8(5, (mainGroup))
    if _mainGroupErr != nil {
        return errors.New("Error serializing 'mainGroup' field " + _mainGroupErr.Error())
    }

    // Simple Field (middleGroup)
    middleGroup := uint8(m.MiddleGroup)
    _middleGroupErr := io.WriteUint8(3, (middleGroup))
    if _middleGroupErr != nil {
        return errors.New("Error serializing 'middleGroup' field " + _middleGroupErr.Error())
    }

    // Simple Field (subGroup)
    subGroup := uint8(m.SubGroup)
    _subGroupErr := io.WriteUint8(8, (subGroup))
    if _subGroupErr != nil {
        return errors.New("Error serializing 'subGroup' field " + _subGroupErr.Error())
    }

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *KNXGroupAddress3Level) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            case "mainGroup":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.MainGroup = data
            case "middleGroup":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.MiddleGroup = data
            case "subGroup":
                var data uint8
                if err := d.DecodeElement(&data, &tok); err != nil {
                    return err
                }
                m.SubGroup = data
            }
        }
    }
}

func (m *KNXGroupAddress3Level) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.KNXGroupAddress3Level"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MainGroup, xml.StartElement{Name: xml.Name{Local: "mainGroup"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MiddleGroup, xml.StartElement{Name: xml.Name{Local: "middleGroup"}}); err != nil {
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

