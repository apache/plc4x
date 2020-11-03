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
type BACnetServiceAckGetEventInformation struct {
    BACnetServiceAck
}

// The corresponding interface
type IBACnetServiceAckGetEventInformation interface {
    IBACnetServiceAck
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetServiceAckGetEventInformation) ServiceChoice() uint8 {
    return 0x1D
}

func (m BACnetServiceAckGetEventInformation) initialize() spi.Message {
    return m
}

func NewBACnetServiceAckGetEventInformation() BACnetServiceAckInitializer {
    return &BACnetServiceAckGetEventInformation{}
}

func CastIBACnetServiceAckGetEventInformation(structType interface{}) IBACnetServiceAckGetEventInformation {
    castFunc := func(typ interface{}) IBACnetServiceAckGetEventInformation {
        if iBACnetServiceAckGetEventInformation, ok := typ.(IBACnetServiceAckGetEventInformation); ok {
            return iBACnetServiceAckGetEventInformation
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetServiceAckGetEventInformation(structType interface{}) BACnetServiceAckGetEventInformation {
    castFunc := func(typ interface{}) BACnetServiceAckGetEventInformation {
        if sBACnetServiceAckGetEventInformation, ok := typ.(BACnetServiceAckGetEventInformation); ok {
            return sBACnetServiceAckGetEventInformation
        }
        if sBACnetServiceAckGetEventInformation, ok := typ.(*BACnetServiceAckGetEventInformation); ok {
            return *sBACnetServiceAckGetEventInformation
        }
        return BACnetServiceAckGetEventInformation{}
    }
    return castFunc(structType)
}

func (m BACnetServiceAckGetEventInformation) LengthInBits() uint16 {
    var lengthInBits uint16 = m.BACnetServiceAck.LengthInBits()

    return lengthInBits
}

func (m BACnetServiceAckGetEventInformation) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetServiceAckGetEventInformationParse(io *utils.ReadBuffer) (BACnetServiceAckInitializer, error) {

    // Create the instance
    return NewBACnetServiceAckGetEventInformation(), nil
}

func (m BACnetServiceAckGetEventInformation) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

        return nil
    }
    return BACnetServiceAckSerialize(io, m.BACnetServiceAck, CastIBACnetServiceAck(m), ser)
}

func (m *BACnetServiceAckGetEventInformation) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m BACnetServiceAckGetEventInformation) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckGetEventInformation"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

