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
type BACnetConfirmedServiceRequestReadPropertyMultiple struct {
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestReadPropertyMultiple interface {
    IBACnetConfirmedServiceRequest
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestReadPropertyMultiple) ServiceChoice() uint8 {
    return 0x0E
}

func (m BACnetConfirmedServiceRequestReadPropertyMultiple) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestReadPropertyMultiple() BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestReadPropertyMultiple{}
}

func CastIBACnetConfirmedServiceRequestReadPropertyMultiple(structType interface{}) IBACnetConfirmedServiceRequestReadPropertyMultiple {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestReadPropertyMultiple {
        if iBACnetConfirmedServiceRequestReadPropertyMultiple, ok := typ.(IBACnetConfirmedServiceRequestReadPropertyMultiple); ok {
            return iBACnetConfirmedServiceRequestReadPropertyMultiple
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestReadPropertyMultiple(structType interface{}) BACnetConfirmedServiceRequestReadPropertyMultiple {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequestReadPropertyMultiple {
        if sBACnetConfirmedServiceRequestReadPropertyMultiple, ok := typ.(BACnetConfirmedServiceRequestReadPropertyMultiple); ok {
            return sBACnetConfirmedServiceRequestReadPropertyMultiple
        }
        if sBACnetConfirmedServiceRequestReadPropertyMultiple, ok := typ.(*BACnetConfirmedServiceRequestReadPropertyMultiple); ok {
            return *sBACnetConfirmedServiceRequestReadPropertyMultiple
        }
        return BACnetConfirmedServiceRequestReadPropertyMultiple{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestReadPropertyMultiple) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceRequestReadPropertyMultiple) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestReadPropertyMultipleParse(io *utils.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceRequestReadPropertyMultiple(), nil
}

func (m BACnetConfirmedServiceRequestReadPropertyMultiple) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}

func (m *BACnetConfirmedServiceRequestReadPropertyMultiple) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m BACnetConfirmedServiceRequestReadPropertyMultiple) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequestReadPropertyMultiple"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

