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
type BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification struct {
    BACnetUnconfirmedServiceRequest
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification interface {
    IBACnetUnconfirmedServiceRequest
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) ServiceChoice() uint8 {
    return 0x02
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) initialize() spi.Message {
    return m
}

func NewBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification() BACnetUnconfirmedServiceRequestInitializer {
    return &BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification{}
}

func CastIBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification(structType interface{}) IBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification {
    castFunc := func(typ interface{}) IBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification {
        if iBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification, ok := typ.(IBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification); ok {
            return iBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification(structType interface{}) BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification {
    castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification {
        if sBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification, ok := typ.(BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification); ok {
            return sBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification
        }
        if sBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification, ok := typ.(*BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification); ok {
            return *sBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification
        }
        return BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification{}
    }
    return castFunc(structType)
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetUnconfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationParse(io *utils.ReadBuffer) (BACnetUnconfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetUnconfirmedServiceRequestUnconfirmedCOVNotification(), nil
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetUnconfirmedServiceRequestSerialize(io, m.BACnetUnconfirmedServiceRequest, CastIBACnetUnconfirmedServiceRequest(m), ser)
}

func (m *BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

