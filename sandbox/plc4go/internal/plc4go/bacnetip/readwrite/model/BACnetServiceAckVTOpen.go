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
type BACnetServiceAckVTOpen struct {
    BACnetServiceAck
}

// The corresponding interface
type IBACnetServiceAckVTOpen interface {
    IBACnetServiceAck
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetServiceAckVTOpen) ServiceChoice() uint8 {
    return 0x15
}

func (m BACnetServiceAckVTOpen) initialize() spi.Message {
    return m
}

func NewBACnetServiceAckVTOpen() BACnetServiceAckInitializer {
    return &BACnetServiceAckVTOpen{}
}

func CastIBACnetServiceAckVTOpen(structType interface{}) IBACnetServiceAckVTOpen {
    castFunc := func(typ interface{}) IBACnetServiceAckVTOpen {
        if iBACnetServiceAckVTOpen, ok := typ.(IBACnetServiceAckVTOpen); ok {
            return iBACnetServiceAckVTOpen
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetServiceAckVTOpen(structType interface{}) BACnetServiceAckVTOpen {
    castFunc := func(typ interface{}) BACnetServiceAckVTOpen {
        if sBACnetServiceAckVTOpen, ok := typ.(BACnetServiceAckVTOpen); ok {
            return sBACnetServiceAckVTOpen
        }
        if sBACnetServiceAckVTOpen, ok := typ.(*BACnetServiceAckVTOpen); ok {
            return *sBACnetServiceAckVTOpen
        }
        return BACnetServiceAckVTOpen{}
    }
    return castFunc(structType)
}

func (m BACnetServiceAckVTOpen) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetServiceAck.LengthInBits()

    return lengthInBits
}

func (m BACnetServiceAckVTOpen) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetServiceAckVTOpenParse(io *utils.ReadBuffer) (BACnetServiceAckInitializer, error) {

    // Create the instance
    return NewBACnetServiceAckVTOpen(), nil
}

func (m BACnetServiceAckVTOpen) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetServiceAckSerialize(io, m.BACnetServiceAck, CastIBACnetServiceAck(m), ser)
}

func (m *BACnetServiceAckVTOpen) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m BACnetServiceAckVTOpen) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckVTOpen"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

