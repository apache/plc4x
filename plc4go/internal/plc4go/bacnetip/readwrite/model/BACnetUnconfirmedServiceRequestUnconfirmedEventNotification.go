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
type BACnetUnconfirmedServiceRequestUnconfirmedEventNotification struct {
    BACnetUnconfirmedServiceRequest
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestUnconfirmedEventNotification interface {
    IBACnetUnconfirmedServiceRequest
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) ServiceChoice() uint8 {
    return 0x03
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) initialize() spi.Message {
    return m
}

func NewBACnetUnconfirmedServiceRequestUnconfirmedEventNotification() BACnetUnconfirmedServiceRequestInitializer {
    return &BACnetUnconfirmedServiceRequestUnconfirmedEventNotification{}
}

func CastIBACnetUnconfirmedServiceRequestUnconfirmedEventNotification(structType interface{}) IBACnetUnconfirmedServiceRequestUnconfirmedEventNotification {
    castFunc := func(typ interface{}) IBACnetUnconfirmedServiceRequestUnconfirmedEventNotification {
        if iBACnetUnconfirmedServiceRequestUnconfirmedEventNotification, ok := typ.(IBACnetUnconfirmedServiceRequestUnconfirmedEventNotification); ok {
            return iBACnetUnconfirmedServiceRequestUnconfirmedEventNotification
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetUnconfirmedServiceRequestUnconfirmedEventNotification(structType interface{}) BACnetUnconfirmedServiceRequestUnconfirmedEventNotification {
    castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequestUnconfirmedEventNotification {
        if sBACnetUnconfirmedServiceRequestUnconfirmedEventNotification, ok := typ.(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification); ok {
            return sBACnetUnconfirmedServiceRequestUnconfirmedEventNotification
        }
        if sBACnetUnconfirmedServiceRequestUnconfirmedEventNotification, ok := typ.(*BACnetUnconfirmedServiceRequestUnconfirmedEventNotification); ok {
            return *sBACnetUnconfirmedServiceRequestUnconfirmedEventNotification
        }
        return BACnetUnconfirmedServiceRequestUnconfirmedEventNotification{}
    }
    return castFunc(structType)
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetUnconfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestUnconfirmedEventNotificationParse(io *utils.ReadBuffer) (BACnetUnconfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetUnconfirmedServiceRequestUnconfirmedEventNotification(), nil
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetUnconfirmedServiceRequestSerialize(io, m.BACnetUnconfirmedServiceRequest, CastIBACnetUnconfirmedServiceRequest(m), ser)
}

func (m *BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestUnconfirmedEventNotification"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

