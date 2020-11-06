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

func (m *BACnetUnconfirmedServiceRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetUnconfirmedServiceRequest"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

