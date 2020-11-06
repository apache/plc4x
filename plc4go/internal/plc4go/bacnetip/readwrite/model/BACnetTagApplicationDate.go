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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type BACnetTagApplicationDate struct {
    Parent *BACnetTag
    IBACnetTagApplicationDate
}

// The corresponding interface
type IBACnetTagApplicationDate interface {
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *BACnetTagApplicationDate) ContextSpecificTag() uint8 {
    return 0
}


func (m *BACnetTagApplicationDate) InitializeParent(parent *BACnetTag, typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) {
    m.Parent.TypeOrTagNumber = typeOrTagNumber
    m.Parent.LengthValueType = lengthValueType
    m.Parent.ExtTagNumber = extTagNumber
    m.Parent.ExtLength = extLength
}

func NewBACnetTagApplicationDate(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) *BACnetTag {
    child := &BACnetTagApplicationDate{
        Parent: NewBACnetTag(typeOrTagNumber, lengthValueType, extTagNumber, extLength),
    }
    child.Parent.Child = child
    return child.Parent
}

func CastBACnetTagApplicationDate(structType interface{}) BACnetTagApplicationDate {
    castFunc := func(typ interface{}) BACnetTagApplicationDate {
        if casted, ok := typ.(BACnetTagApplicationDate); ok {
            return casted
        }
        if casted, ok := typ.(*BACnetTagApplicationDate); ok {
            return *casted
        }
        if casted, ok := typ.(BACnetTag); ok {
            return CastBACnetTagApplicationDate(casted.Child)
        }
        if casted, ok := typ.(*BACnetTag); ok {
            return CastBACnetTagApplicationDate(casted.Child)
        }
        return BACnetTagApplicationDate{}
    }
    return castFunc(structType)
}

func (m *BACnetTagApplicationDate) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    return lengthInBits
}

func (m *BACnetTagApplicationDate) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetTagApplicationDateParse(io *utils.ReadBuffer) (*BACnetTag, error) {

    // Create a partially initialized instance
    _child := &BACnetTagApplicationDate{
        Parent: &BACnetTag{},
    }
    _child.Parent.Child = _child
    return _child.Parent, nil
}

func (m *BACnetTagApplicationDate) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return m.Parent.SerializeParent(io, m, ser)
}

func (m *BACnetTagApplicationDate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
            }
        }
    }
}

func (m *BACnetTagApplicationDate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetTagApplicationDate"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

