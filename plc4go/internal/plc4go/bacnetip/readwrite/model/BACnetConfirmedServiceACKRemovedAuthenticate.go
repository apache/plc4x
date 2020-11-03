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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type BACnetConfirmedServiceACKRemovedAuthenticate struct {
    BACnetConfirmedServiceACK
}

// The corresponding interface
type IBACnetConfirmedServiceACKRemovedAuthenticate interface {
    IBACnetConfirmedServiceACK
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceACKRemovedAuthenticate) ServiceChoice() uint8 {
    return 0x18
}

func (m BACnetConfirmedServiceACKRemovedAuthenticate) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceACKRemovedAuthenticate() BACnetConfirmedServiceACKInitializer {
    return &BACnetConfirmedServiceACKRemovedAuthenticate{}
}

func CastIBACnetConfirmedServiceACKRemovedAuthenticate(structType interface{}) IBACnetConfirmedServiceACKRemovedAuthenticate {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceACKRemovedAuthenticate {
        if iBACnetConfirmedServiceACKRemovedAuthenticate, ok := typ.(IBACnetConfirmedServiceACKRemovedAuthenticate); ok {
            return iBACnetConfirmedServiceACKRemovedAuthenticate
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceACKRemovedAuthenticate(structType interface{}) BACnetConfirmedServiceACKRemovedAuthenticate {
    castFunc := func(typ interface{}) BACnetConfirmedServiceACKRemovedAuthenticate {
        if sBACnetConfirmedServiceACKRemovedAuthenticate, ok := typ.(BACnetConfirmedServiceACKRemovedAuthenticate); ok {
            return sBACnetConfirmedServiceACKRemovedAuthenticate
        }
        if sBACnetConfirmedServiceACKRemovedAuthenticate, ok := typ.(*BACnetConfirmedServiceACKRemovedAuthenticate); ok {
            return *sBACnetConfirmedServiceACKRemovedAuthenticate
        }
        return BACnetConfirmedServiceACKRemovedAuthenticate{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceACKRemovedAuthenticate) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceACK.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceACKRemovedAuthenticate) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceACKRemovedAuthenticateParse(io *utils.ReadBuffer) (BACnetConfirmedServiceACKInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceACKRemovedAuthenticate(), nil
}

func (m BACnetConfirmedServiceACKRemovedAuthenticate) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceACKSerialize(io, m.BACnetConfirmedServiceACK, CastIBACnetConfirmedServiceACK(m), ser)
}

func (m *BACnetConfirmedServiceACKRemovedAuthenticate) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m BACnetConfirmedServiceACKRemovedAuthenticate) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceACKRemovedAuthenticate"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

