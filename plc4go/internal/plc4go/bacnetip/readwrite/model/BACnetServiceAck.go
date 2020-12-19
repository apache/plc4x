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
    "github.com/apache/plc4x/plc4go/internal/plc4go/utils"
    "reflect"
    "strings"
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
    xml.Marshaler
}

type IBACnetServiceAckParent interface {
    SerializeParent(io utils.WriteBuffer, child IBACnetServiceAck, serializeChildFunction func() error) error
    GetTypeName() string
}

type IBACnetServiceAckChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *BACnetServiceAck)
    GetTypeName() string
    IBACnetServiceAck
}

func NewBACnetServiceAck() *BACnetServiceAck {
    return &BACnetServiceAck{}
}

func CastBACnetServiceAck(structType interface{}) *BACnetServiceAck {
    castFunc := func(typ interface{}) *BACnetServiceAck {
        if casted, ok := typ.(BACnetServiceAck); ok {
            return &casted
        }
        if casted, ok := typ.(*BACnetServiceAck); ok {
            return casted
        }
        return nil
    }
    return castFunc(structType)
}

func (m *BACnetServiceAck) GetTypeName() string {
    return "BACnetServiceAck"
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
    var token xml.Token
    var err error
    for {
        token, err = d.Token()
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
            default:
                switch start.Attr[0].Value {
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckGetAlarmSummary":
                        var dt *BACnetServiceAckGetAlarmSummary
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckGetAlarmSummary)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckGetEnrollmentSummary":
                        var dt *BACnetServiceAckGetEnrollmentSummary
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckGetEnrollmentSummary)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckGetEventInformation":
                        var dt *BACnetServiceAckGetEventInformation
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckGetEventInformation)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckAtomicReadFile":
                        var dt *BACnetServiceAckAtomicReadFile
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckAtomicReadFile)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckAtomicWriteFile":
                        var dt *BACnetServiceAckAtomicWriteFile
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckAtomicWriteFile)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckCreateObject":
                        var dt *BACnetServiceAckCreateObject
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckCreateObject)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckReadProperty":
                        var dt *BACnetServiceAckReadProperty
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckReadProperty)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckReadPropertyMultiple":
                        var dt *BACnetServiceAckReadPropertyMultiple
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckReadPropertyMultiple)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckReadRange":
                        var dt *BACnetServiceAckReadRange
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckReadRange)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckConfirmedPrivateTransfer":
                        var dt *BACnetServiceAckConfirmedPrivateTransfer
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckConfirmedPrivateTransfer)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckVTOpen":
                        var dt *BACnetServiceAckVTOpen
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckVTOpen)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckVTData":
                        var dt *BACnetServiceAckVTData
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckVTData)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckRemovedAuthenticate":
                        var dt *BACnetServiceAckRemovedAuthenticate
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckRemovedAuthenticate)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                    case "org.apache.plc4x.java.bacnetip.readwrite.BACnetServiceAckRemovedReadPropertyConditional":
                        var dt *BACnetServiceAckRemovedReadPropertyConditional
                        if m.Child != nil {
                            dt = m.Child.(*BACnetServiceAckRemovedReadPropertyConditional)
                        }
                        if err := d.DecodeElement(&dt, &tok); err != nil {
                            return err
                        }
                        if m.Child == nil {
                            dt.Parent = m
                            m.Child = dt
                        }
                }
            }
        }
    }
}

func (m *BACnetServiceAck) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    className := reflect.TypeOf(m.Child).String()
    className = "org.apache.plc4x.java.bacnetip.readwrite." + className[strings.LastIndex(className, ".") + 1:]
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: className},
        }}); err != nil {
        return err
    }
    marshaller, ok := m.Child.(xml.Marshaler)
    if !ok {
        return errors.New("child is not castable to Marshaler")
    }
    marshaller.MarshalXML(e, start)
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

