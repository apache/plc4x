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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type BACnetServiceAck struct {
    Child IBACnetServiceAckChild
    IBACnetServiceAck
    IBACnetServiceAckParent
}

// The corresponding interface
type IBACnetServiceAck interface {
    ServiceChoice() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

type IBACnetServiceAckParent interface {
    SerializeParent(io utils.WriteBuffer, child IBACnetServiceAck, serializeChildFunction func() error) error
}

type IBACnetServiceAckChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *BACnetServiceAck)
    IBACnetServiceAck
}

func NewBACnetServiceAck() *BACnetServiceAck {
    return &BACnetServiceAck{}
}

func CastBACnetServiceAck(structType interface{}) BACnetServiceAck {
    castFunc := func(typ interface{}) BACnetServiceAck {
        if casted, ok := typ.(BACnetServiceAck); ok {
            return casted
        }
        if casted, ok := typ.(*BACnetServiceAck); ok {
            return *casted
        }
        return BACnetServiceAck{}
    }
    return castFunc(structType)
}

func (m *BACnetServiceAck) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (serviceChoice)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *BACnetServiceAck) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetServiceAckParse(io *utils.ReadBuffer) (*BACnetServiceAck, error) {

    // Discriminator Field (serviceChoice) (Used as input to a switch field)
    serviceChoice, _serviceChoiceErr := io.ReadUint8(8)
    if _serviceChoiceErr != nil {
        return nil, errors.New("Error parsing 'serviceChoice' field " + _serviceChoiceErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *BACnetServiceAck
    var typeSwitchError error
    switch {
    case serviceChoice == 0x03:
        _parent, typeSwitchError = BACnetServiceAckGetAlarmSummaryParse(io)
    case serviceChoice == 0x04:
        _parent, typeSwitchError = BACnetServiceAckGetEnrollmentSummaryParse(io)
    case serviceChoice == 0x1D:
        _parent, typeSwitchError = BACnetServiceAckGetEventInformationParse(io)
    case serviceChoice == 0x06:
        _parent, typeSwitchError = BACnetServiceAckAtomicReadFileParse(io)
    case serviceChoice == 0x07:
        _parent, typeSwitchError = BACnetServiceAckAtomicWriteFileParse(io)
    case serviceChoice == 0x0A:
        _parent, typeSwitchError = BACnetServiceAckCreateObjectParse(io)
    case serviceChoice == 0x0C:
        _parent, typeSwitchError = BACnetServiceAckReadPropertyParse(io)
    case serviceChoice == 0x0E:
        _parent, typeSwitchError = BACnetServiceAckReadPropertyMultipleParse(io)
    case serviceChoice == 0x1A:
        _parent, typeSwitchError = BACnetServiceAckReadRangeParse(io)
    case serviceChoice == 0x12:
        _parent, typeSwitchError = BACnetServiceAckConfirmedPrivateTransferParse(io)
    case serviceChoice == 0x15:
        _parent, typeSwitchError = BACnetServiceAckVTOpenParse(io)
    case serviceChoice == 0x17:
        _parent, typeSwitchError = BACnetServiceAckVTDataParse(io)
    case serviceChoice == 0x18:
        _parent, typeSwitchError = BACnetServiceAckRemovedAuthenticateParse(io)
    case serviceChoice == 0x0D:
        _parent, typeSwitchError = BACnetServiceAckRemovedReadPropertyConditionalParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *BACnetServiceAck) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *BACnetServiceAck) SerializeParent(io utils.WriteBuffer, child IBACnetServiceAck, serializeChildFunction func() error) error {

    // Discriminator Field (serviceChoice) (Used as input to a switch field)
    serviceChoice := uint8(child.ServiceChoice())
    _serviceChoiceErr := io.WriteUint8(8, (serviceChoice))
    if _serviceChoiceErr != nil {
        return errors.New("Error serializing 'serviceChoice' field " + _serviceChoiceErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := serializeChildFunction()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}

func (m *BACnetServiceAck) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m *BACnetServiceAck) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAck"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

