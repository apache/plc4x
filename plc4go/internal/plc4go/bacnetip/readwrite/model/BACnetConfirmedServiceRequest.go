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
type BACnetConfirmedServiceRequest struct {
    Child IBACnetConfirmedServiceRequestChild
    IBACnetConfirmedServiceRequest
    IBACnetConfirmedServiceRequestParent
}

// The corresponding interface
type IBACnetConfirmedServiceRequest interface {
    ServiceChoice() uint8
    LengthInBytes() uint16
    LengthInBits() uint16
    Serialize(io utils.WriteBuffer) error
}

type IBACnetConfirmedServiceRequestParent interface {
    SerializeParent(io utils.WriteBuffer, child IBACnetConfirmedServiceRequest, serializeChildFunction func() error) error
}

type IBACnetConfirmedServiceRequestChild interface {
    Serialize(io utils.WriteBuffer) error
    InitializeParent(parent *BACnetConfirmedServiceRequest)
    IBACnetConfirmedServiceRequest
}

func NewBACnetConfirmedServiceRequest() *BACnetConfirmedServiceRequest {
    return &BACnetConfirmedServiceRequest{}
}

func CastBACnetConfirmedServiceRequest(structType interface{}) BACnetConfirmedServiceRequest {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequest {
        if casted, ok := typ.(BACnetConfirmedServiceRequest); ok {
            return casted
        }
        if casted, ok := typ.(*BACnetConfirmedServiceRequest); ok {
            return *casted
        }
        return BACnetConfirmedServiceRequest{}
    }
    return castFunc(structType)
}

func (m *BACnetConfirmedServiceRequest) LengthInBits() uint16 {
    lengthInBits := uint16(0)

    // Discriminator Field (serviceChoice)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...
    lengthInBits += m.Child.LengthInBits()

    return lengthInBits
}

func (m *BACnetConfirmedServiceRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestParse(io *utils.ReadBuffer, len uint16) (*BACnetConfirmedServiceRequest, error) {

    // Discriminator Field (serviceChoice) (Used as input to a switch field)
    serviceChoice, _serviceChoiceErr := io.ReadUint8(8)
    if _serviceChoiceErr != nil {
        return nil, errors.New("Error parsing 'serviceChoice' field " + _serviceChoiceErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var _parent *BACnetConfirmedServiceRequest
    var typeSwitchError error
    switch {
    case serviceChoice == 0x00:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestAcknowledgeAlarmParse(io)
    case serviceChoice == 0x01:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestConfirmedCOVNotificationParse(io, len)
    case serviceChoice == 0x02:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestConfirmedEventNotificationParse(io)
    case serviceChoice == 0x04:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestGetEnrollmentSummaryParse(io)
    case serviceChoice == 0x05:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestSubscribeCOVParse(io)
    case serviceChoice == 0x06:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestAtomicReadFileParse(io)
    case serviceChoice == 0x07:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestAtomicWriteFileParse(io)
    case serviceChoice == 0x08:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestAddListElementParse(io)
    case serviceChoice == 0x09:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestRemoveListElementParse(io)
    case serviceChoice == 0x0A:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestCreateObjectParse(io)
    case serviceChoice == 0x0B:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestDeleteObjectParse(io)
    case serviceChoice == 0x0C:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestReadPropertyParse(io)
    case serviceChoice == 0x0E:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestReadPropertyMultipleParse(io)
    case serviceChoice == 0x0F:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestWritePropertyParse(io, len)
    case serviceChoice == 0x10:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestWritePropertyMultipleParse(io)
    case serviceChoice == 0x11:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestDeviceCommunicationControlParse(io)
    case serviceChoice == 0x12:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestConfirmedPrivateTransferParse(io)
    case serviceChoice == 0x13:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestConfirmedTextMessageParse(io)
    case serviceChoice == 0x14:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestReinitializeDeviceParse(io)
    case serviceChoice == 0x15:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestVTOpenParse(io)
    case serviceChoice == 0x16:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestVTCloseParse(io)
    case serviceChoice == 0x17:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestVTDataParse(io)
    case serviceChoice == 0x18:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestRemovedAuthenticateParse(io)
    case serviceChoice == 0x19:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestRemovedRequestKeyParse(io)
    case serviceChoice == 0x0D:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestRemovedReadPropertyConditionalParse(io)
    case serviceChoice == 0x1A:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestReadRangeParse(io)
    case serviceChoice == 0x1B:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestLifeSafetyOperationParse(io)
    case serviceChoice == 0x1C:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestSubscribeCOVPropertyParse(io)
    case serviceChoice == 0x1D:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestGetEventInformationParse(io)
    case serviceChoice == 0x1E:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestSubscribeCOVPropertyMultipleParse(io)
    case serviceChoice == 0x1F:
        _parent, typeSwitchError = BACnetConfirmedServiceRequestConfirmedCOVNotificationMultipleParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Finish initializing
    _parent.Child.InitializeParent(_parent)
    return _parent, nil
}

func (m *BACnetConfirmedServiceRequest) Serialize(io utils.WriteBuffer) error {
    return m.Child.Serialize(io)
}

func (m *BACnetConfirmedServiceRequest) SerializeParent(io utils.WriteBuffer, child IBACnetConfirmedServiceRequest, serializeChildFunction func() error) error {

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

func (m *BACnetConfirmedServiceRequest) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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

func (m *BACnetConfirmedServiceRequest) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
    if err := e.EncodeToken(xml.StartElement{Name: start.Name, Attr: []xml.Attr{
            {Name: xml.Name{Local: "className"}, Value: "org.apache.plc4x.java.bacnetip.readwrite.BACnetConfirmedServiceRequest"},
        }}); err != nil {
        return err
    }
    if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
        return err
    }
    return nil
}

