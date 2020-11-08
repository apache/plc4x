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
    "plc4x.apache.org/plc4go/v0/internal/plc4go/utils"
    "reflect"
    "strings"
)

// The data-structure of this message
type BACnetUnconfirmedServiceRequest struct {
    Child IBACnetUnconfirmedServiceRequestChild
    IBACnetUnconfirmedServiceRequest
    IBACnetUnconfirmedServiceRequestParent
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequest interface {
    ServiceChoice() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
    xml.Marshaler
}

type IBACnetUnconfirmedServiceRequestParent interface {
    SerializeParent(io utils.WriteBuffer, child IBACnetUnconfirmedServiceRequest, serializeChildFunction func() error) error
}

type IBACnetUnconfirmedServiceRequestChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *BACnetUnconfirmedServiceRequest)
    IBACnetUnconfirmedServiceRequest
}

func NewBACnetUnconfirmedServiceRequest() *BACnetUnconfirmedServiceRequest {
    return &BACnetUnconfirmedServiceRequest{}
}

func CastBACnetUnconfirmedServiceRequest(structType interface{}) BACnetUnconfirmedServiceRequest {
    castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequest {
        if casted, ok := typ.(BACnetUnconfirmedServiceRequest); ok {
            return casted
        }
        if casted, ok := typ.(*BACnetUnconfirmedServiceRequest); ok {
            return *casted
        }
        return BACnetUnconfirmedServiceRequest{}
    }
    return castFunc(structType)
}

func (m *BACnetUnconfirmedServiceRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (serviceChoice)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *BACnetUnconfirmedServiceRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestParse(io *utils.ReadBuffer, len uint16) (*BACnetUnconfirmedServiceRequest, error) {

    // Discriminator Field (serviceChoice) (Used as input to a switch field)
    serviceChoice, _serviceChoiceErr := io.ReadUint8(8)
    if _serviceChoiceErr != nil {
        return nil, errors.New("Error parsing 'serviceChoice' field " + _serviceChoiceErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *BACnetUnconfirmedServiceRequest
    var typeSwitchError error
    switch {
    case serviceChoice == 0x00:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestIAmParse(io)
    case serviceChoice == 0x01:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestIHaveParse(io)
    case serviceChoice == 0x02:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationParse(io)
    case serviceChoice == 0x03:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestUnconfirmedEventNotificationParse(io)
    case serviceChoice == 0x04:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransferParse(io, len)
    case serviceChoice == 0x05:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestUnconfirmedTextMessageParse(io)
    case serviceChoice == 0x06:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestTimeSynchronizationParse(io)
    case serviceChoice == 0x07:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestWhoHasParse(io)
    case serviceChoice == 0x08:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestWhoIsParse(io)
    case serviceChoice == 0x09:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestUTCTimeSynchronizationParse(io)
    case serviceChoice == 0x0A:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestWriteGroupParse(io)
    case serviceChoice == 0x0B:
        _parent, typeSwitchError = BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationMultipleParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *BACnetUnconfirmedServiceRequest) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *BACnetUnconfirmedServiceRequest) SerializeParent(io utils.WriteBuffer, child IBACnetUnconfirmedServiceRequest, serializeChildFunction func() error) error {

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

func (m *BACnetUnconfirmedServiceRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestIAm":
                            var dt *BACnetUnconfirmedServiceRequestIAm
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestIAm)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestIHave":
                            var dt *BACnetUnconfirmedServiceRequestIHave
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestIHave)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification":
                            var dt *BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestUnconfirmedCOVNotification)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestUnconfirmedEventNotification":
                            var dt *BACnetUnconfirmedServiceRequestUnconfirmedEventNotification
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestUnconfirmedEventNotification)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer":
                            var dt *BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransfer)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestUnconfirmedTextMessage":
                            var dt *BACnetUnconfirmedServiceRequestUnconfirmedTextMessage
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestUnconfirmedTextMessage)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestTimeSynchronization":
                            var dt *BACnetUnconfirmedServiceRequestTimeSynchronization
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestTimeSynchronization)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestWhoHas":
                            var dt *BACnetUnconfirmedServiceRequestWhoHas
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestWhoHas)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestWhoIs":
                            var dt *BACnetUnconfirmedServiceRequestWhoIs
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestWhoIs)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestUTCTimeSynchronization":
                            var dt *BACnetUnconfirmedServiceRequestUTCTimeSynchronization
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestUTCTimeSynchronization)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestWriteGroup":
                            var dt *BACnetUnconfirmedServiceRequestWriteGroup
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestWriteGroup)
                            }
                            if err := d.DecodeElement(&dt, &tok); err != nil {
                                return err
                            }
                            if m.Child == nil {
                                dt.Parent = m
                                m.Child = dt
                            }
                        case "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationMultiple":
                            var dt *BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationMultiple
                            if m.Child != nil {
                                dt = m.Child.(*BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationMultiple)
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

func (m *BACnetUnconfirmedServiceRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
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

