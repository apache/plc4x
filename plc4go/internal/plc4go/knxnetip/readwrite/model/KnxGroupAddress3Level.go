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
    "github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
    "io"
)

// The data-structure of this message
type KnxGroupAddress3Level struct {
    MainGroup uint8
    MiddleGroup uint8
    SubGroup uint8
    Parent *KnxGroupAddress
    IKnxGroupAddress3Level
}

// The corresponding interface
type IKnxGroupAddress3Level interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *KnxGroupAddress3Level) NumLevels() uint8 {
    return 3
}


func (m *KnxGroupAddress3Level) InitializeParent(parent *KnxGroupAddress) {
}

func NewKnxGroupAddress3Level(mainGroup uint8, middleGroup uint8, subGroup uint8, ) *KnxGroupAddress {
    child := &KnxGroupAddress3Level{
        MainGroup: mainGroup,
        MiddleGroup: middleGroup,
        SubGroup: subGroup,
        Parent: NewKnxGroupAddress(),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastKnxGroupAddress3Level(structType interface{}) *KnxGroupAddress3Level {
    castFunc := func(typ interface{}) *KnxGroupAddress3Level {
        if casted, ok := typ.(KnxGroupAddress3Level); ok {
            return &casted
        }
        if casted, ok := typ.(*KnxGroupAddress3Level); ok {
            return casted
        }
        if casted, ok := typ.(KnxGroupAddress); ok {
            return CastKnxGroupAddress3Level(casted.Child)
        }
        if casted, ok := typ.(*KnxGroupAddress); ok {
            return CastKnxGroupAddress3Level(casted.Child)
        }
        return nil
    }
    return castFunc(structType)
}

func (m *KnxGroupAddress3Level) GetTypeName() string {
    return "KnxGroupAddress3Level"
}

func (m *KnxGroupAddress3Level) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Simple field (mainGroup)
    lengthInBits += 5

    // Simple field (middleGroup)
    lengthInBits += 3

    // Simple field (subGroup)
    lengthInBits += 8

    return lengthInBits
}

func (m *KnxGroupAddress3Level) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KnxGroupAddress3LevelParse(io *utils.ReadBuffer) (*KnxGroupAddress, error) {

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
    _child := &KnxGroupAddress3Level{
        MainGroup: mainGroup,
        MiddleGroup: middleGroup,
        SubGroup: subGroup,
        Parent: &KnxGroupAddress{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *KnxGroupAddress3Level) Serialize(io utils.WriteBuffer) error {
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

func (m *KnxGroupAddress3Level) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var token xml.Token
    var err error
    token = start
    for {
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
        token, err = d.Token()
        if err != nil {
            if err == io.EOF {
                return nil
            }
            return err
        }
    }
}

func (m *KnxGroupAddress3Level) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeElement(m.MainGroup, xml.StartElement{Name: xml.Name{Local: "mainGroup"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MiddleGroup, xml.StartElement{Name: xml.Name{Local: "middleGroup"}}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.SubGroup, xml.StartElement{Name: xml.Name{Local: "subGroup"}}); err != nil {
        return err
    }
    return nil
}

