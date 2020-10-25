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
type BACnetErrorReadPropertyMultiple struct {
    BACnetError
}

// The corresponding interface
type IBACnetErrorReadPropertyMultiple interface {
    IBACnetError
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetErrorReadPropertyMultiple) ServiceChoice() uint8 {
    return 0x0E
}

func (m BACnetErrorReadPropertyMultiple) initialize() spi.Message {
    return m
}

func NewBACnetErrorReadPropertyMultiple() BACnetErrorInitializer {
    return &BACnetErrorReadPropertyMultiple{}
}

func CastIBACnetErrorReadPropertyMultiple(structType interface{}) IBACnetErrorReadPropertyMultiple {
    castFunc := func(typ interface{}) IBACnetErrorReadPropertyMultiple {
        if iBACnetErrorReadPropertyMultiple, ok := typ.(IBACnetErrorReadPropertyMultiple); ok {
            return iBACnetErrorReadPropertyMultiple
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetErrorReadPropertyMultiple(structType interface{}) BACnetErrorReadPropertyMultiple {
    castFunc := func(typ interface{}) BACnetErrorReadPropertyMultiple {
        if sBACnetErrorReadPropertyMultiple, ok := typ.(BACnetErrorReadPropertyMultiple); ok {
            return sBACnetErrorReadPropertyMultiple
        }
        if sBACnetErrorReadPropertyMultiple, ok := typ.(*BACnetErrorReadPropertyMultiple); ok {
            return *sBACnetErrorReadPropertyMultiple
        }
        return BACnetErrorReadPropertyMultiple{}
    }
    return castFunc(structType)
}

func (m BACnetErrorReadPropertyMultiple) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetError.LengthInBits()

    return lengthInBits
}

func (m BACnetErrorReadPropertyMultiple) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetErrorReadPropertyMultipleParse(io *utils.ReadBuffer) (BACnetErrorInitializer, error) {

    // Create the instance
    return NewBACnetErrorReadPropertyMultiple(), nil
}

func (m BACnetErrorReadPropertyMultiple) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetErrorSerialize(io, m.BACnetError, CastIBACnetError(m), ser)
}

func (m *BACnetErrorReadPropertyMultiple) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m BACnetErrorReadPropertyMultiple) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetErrorReadPropertyMultiple"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

