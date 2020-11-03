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
type BACnetConfirmedServiceRequestGetEnrollmentSummary struct {
    BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestGetEnrollmentSummary interface {
    IBACnetConfirmedServiceRequest
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestGetEnrollmentSummary) ServiceChoice() uint8 {
    return 0x04
}

func (m BACnetConfirmedServiceRequestGetEnrollmentSummary) initialize() spi.Message {
    return m
}

func NewBACnetConfirmedServiceRequestGetEnrollmentSummary() BACnetConfirmedServiceRequestInitializer {
    return &BACnetConfirmedServiceRequestGetEnrollmentSummary{}
}

func CastIBACnetConfirmedServiceRequestGetEnrollmentSummary(structType interface{}) IBACnetConfirmedServiceRequestGetEnrollmentSummary {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestGetEnrollmentSummary {
        if iBACnetConfirmedServiceRequestGetEnrollmentSummary, ok := typ.(IBACnetConfirmedServiceRequestGetEnrollmentSummary); ok {
            return iBACnetConfirmedServiceRequestGetEnrollmentSummary
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestGetEnrollmentSummary(structType interface{}) BACnetConfirmedServiceRequestGetEnrollmentSummary {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequestGetEnrollmentSummary {
        if sBACnetConfirmedServiceRequestGetEnrollmentSummary, ok := typ.(BACnetConfirmedServiceRequestGetEnrollmentSummary); ok {
            return sBACnetConfirmedServiceRequestGetEnrollmentSummary
        }
        if sBACnetConfirmedServiceRequestGetEnrollmentSummary, ok := typ.(*BACnetConfirmedServiceRequestGetEnrollmentSummary); ok {
            return *sBACnetConfirmedServiceRequestGetEnrollmentSummary
        }
        return BACnetConfirmedServiceRequestGetEnrollmentSummary{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestGetEnrollmentSummary) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

    return lengthInBits
}

func (m BACnetConfirmedServiceRequestGetEnrollmentSummary) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestGetEnrollmentSummaryParse(io *utils.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

    // Create the instance
    return NewBACnetConfirmedServiceRequestGetEnrollmentSummary(), nil
}

func (m BACnetConfirmedServiceRequestGetEnrollmentSummary) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}

func (m *BACnetConfirmedServiceRequestGetEnrollmentSummary) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m BACnetConfirmedServiceRequestGetEnrollmentSummary) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequestGetEnrollmentSummary"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

