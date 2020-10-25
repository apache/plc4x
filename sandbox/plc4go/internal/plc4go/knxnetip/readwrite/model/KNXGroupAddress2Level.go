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
)

// The data-structure of this message
type KNXGroupAddress2Level struct {
    MainGroup uint8
    SubGroup uint16
    KNXGroupAddress
}

// The corresponding interface
type IKNXGroupAddress2Level interface {
    IKNXGroupAddress
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m KNXGroupAddress2Level) NumLevels() uint8 {
    return 2
}

func (m KNXGroupAddress2Level) initialize() spi.Message {
    return m
}

func NewKNXGroupAddress2Level(mainGroup uint8, subGroup uint16) KNXGroupAddressInitializer {
    return &KNXGroupAddress2Level{MainGroup: mainGroup, SubGroup: subGroup}
}

func CastIKNXGroupAddress2Level(structType interface{}) IKNXGroupAddress2Level {
    castFunc := func(typ interface{}) IKNXGroupAddress2Level {
        if iKNXGroupAddress2Level, ok := typ.(IKNXGroupAddress2Level); ok {
            return iKNXGroupAddress2Level
        }
        return nil
    }
    return castFunc(structType)
}

func CastKNXGroupAddress2Level(structType interface{}) KNXGroupAddress2Level {
    castFunc := func(typ interface{}) KNXGroupAddress2Level {
        if sKNXGroupAddress2Level, ok := typ.(KNXGroupAddress2Level); ok {
            return sKNXGroupAddress2Level
        }
        if sKNXGroupAddress2Level, ok := typ.(*KNXGroupAddress2Level); ok {
            return *sKNXGroupAddress2Level
        }
        return KNXGroupAddress2Level{}
    }
    return castFunc(structType)
}

func (m KNXGroupAddress2Level) LengthInBits() uint16 {
    var lengthInBits uint16 = m.KNXGroupAddress.LengthInBits()

    // Simple field (mainGroup)
    lengthInBits += 5

    // Simple field (subGroup)
    lengthInBits += 11

    return lengthInBits
}

func (m KNXGroupAddress2Level) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func KNXGroupAddress2LevelParse(io *utils.ReadBuffer) (KNXGroupAddressInitializer, error) {

    // Simple Field (mainGroup)
    mainGroup, _mainGroupErr := io.ReadUint8(5)
    if _mainGroupErr != nil {
        return nil, errors.New("Error parsing 'mainGroup' field " + _mainGroupErr.Error())
    }

    // Simple Field (subGroup)
    subGroup, _subGroupErr := io.ReadUint16(11)
    if _subGroupErr != nil {
        return nil, errors.New("Error parsing 'subGroup' field " + _subGroupErr.Error())
    }

    // Create the instance
    return NewKNXGroupAddress2Level(mainGroup, subGroup), nil
}

func (m KNXGroupAddress2Level) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Simple Field (mainGroup)
    mainGroup := uint8(m.MainGroup)
    _mainGroupErr := io.WriteUint8(5, (mainGroup))
    if _mainGroupErr != nil {
        return errors.New("Error serializing 'mainGroup' field " + _mainGroupErr.Error())
    }

    // Simple Field (subGroup)
    subGroup := uint16(m.SubGroup)
    _subGroupErr := io.WriteUint16(11, (subGroup))
    if _subGroupErr != nil {
        return errors.New("Error serializing 'subGroup' field " + _subGroupErr.Error())
    }

        return nil
    }
    return KNXGroupAddressSerialize(io, m.KNXGroupAddress, CastIKNXGroupAddress(m), ser)
}

func (m *KNXGroupAddress2Level) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m KNXGroupAddress2Level) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.knxnetip.readwrite.KNXGroupAddress2Level"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeElement(m.MainGroup, xml.StartElement{Name: xml.Name{Local: "mainGroup"}}); err != nil {
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

