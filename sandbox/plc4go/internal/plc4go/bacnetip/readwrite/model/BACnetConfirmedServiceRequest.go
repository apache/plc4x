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
    "errors"
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
)

// The data-structure of this message
type BACnetConfirmedServiceRequest struct {

}

// The corresponding interface
type IBACnetConfirmedServiceRequest interface {
    spi.Message
    ServiceChoice() uint8
    Serialize(io utils.WriteBuffer) error
}

type BACnetConfirmedServiceRequestInitializer interface {
    initialize() spi.Message
}

func BACnetConfirmedServiceRequestServiceChoice(m IBACnetConfirmedServiceRequest) uint8 {
    return m.ServiceChoice()
}


func CastIBACnetConfirmedServiceRequest(structType interface{}) IBACnetConfirmedServiceRequest {
    castFunc := func(typ interface{}) IBACnetConfirmedServiceRequest {
        if iBACnetConfirmedServiceRequest, ok := typ.(IBACnetConfirmedServiceRequest); ok {
            return iBACnetConfirmedServiceRequest
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetConfirmedServiceRequest(structType interface{}) BACnetConfirmedServiceRequest {
    castFunc := func(typ interface{}) BACnetConfirmedServiceRequest {
        if sBACnetConfirmedServiceRequest, ok := typ.(BACnetConfirmedServiceRequest); ok {
            return sBACnetConfirmedServiceRequest
        }
        if sBACnetConfirmedServiceRequest, ok := typ.(*BACnetConfirmedServiceRequest); ok {
            return *sBACnetConfirmedServiceRequest
        }
        return BACnetConfirmedServiceRequest{}
    }
    return castFunc(structType)
}

func (m BACnetConfirmedServiceRequest) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Discriminator Field (serviceChoice)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits
}

func (m BACnetConfirmedServiceRequest) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestParse(io *utils.ReadBuffer, len uint16) (spi.Message, error) {

    // Discriminator Field (serviceChoice) (Used as input to a switch field)
    serviceChoice, _serviceChoiceErr := io.ReadUint8(8)
    if _serviceChoiceErr != nil {
        return nil, errors.New("Error parsing 'serviceChoice' field " + _serviceChoiceErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var initializer BACnetConfirmedServiceRequestInitializer
    var typeSwitchError error
    switch {
    case serviceChoice == 0x00:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestAcknowledgeAlarmParse(io)
    case serviceChoice == 0x01:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestConfirmedCOVNotificationParse(io, len)
    case serviceChoice == 0x1F:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestConfirmedCOVNotificationMultipleParse(io)
    case serviceChoice == 0x02:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestConfirmedEventNotificationParse(io)
    case serviceChoice == 0x04:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestGetEnrollmentSummaryParse(io)
    case serviceChoice == 0x1D:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestGetEventInformationParse(io)
    case serviceChoice == 0x1B:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestLifeSafetyOperationParse(io)
    case serviceChoice == 0x05:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestSubscribeCOVParse(io)
    case serviceChoice == 0x1C:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestSubscribeCOVPropertyParse(io)
    case serviceChoice == 0x1E:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestSubscribeCOVPropertyMultipleParse(io)
    case serviceChoice == 0x06:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestAtomicReadFileParse(io)
    case serviceChoice == 0x07:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestAtomicWriteFileParse(io)
    case serviceChoice == 0x08:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestAddListElementParse(io)
    case serviceChoice == 0x09:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestRemoveListElementParse(io)
    case serviceChoice == 0x0A:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestCreateObjectParse(io)
    case serviceChoice == 0x0B:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestDeleteObjectParse(io)
    case serviceChoice == 0x0C:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestReadPropertyParse(io)
    case serviceChoice == 0x0E:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestReadPropertyMultipleParse(io)
    case serviceChoice == 0x1A:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestReadRangeParse(io)
    case serviceChoice == 0x0F:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestWritePropertyParse(io, len)
    case serviceChoice == 0x10:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestWritePropertyMultipleParse(io)
    case serviceChoice == 0x11:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestDeviceCommunicationControlParse(io)
    case serviceChoice == 0x12:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestConfirmedPrivateTransferParse(io)
    case serviceChoice == 0x13:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestConfirmedTextMessageParse(io)
    case serviceChoice == 0x14:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestReinitializeDeviceParse(io)
    case serviceChoice == 0x15:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestVTOpenParse(io)
    case serviceChoice == 0x16:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestVTCloseParse(io)
    case serviceChoice == 0x17:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestVTDataParse(io)
    case serviceChoice == 0x18:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestRemovedAuthenticateParse(io)
    case serviceChoice == 0x19:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestRemovedRequestKeyParse(io)
    case serviceChoice == 0x0D:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestRemovedReadPropertyConditionalParse(io)
    case serviceChoice == 0x1A:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestReadRangeParse(io)
    case serviceChoice == 0x1B:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestLifeSafetyOperationParse(io)
    case serviceChoice == 0x1C:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestSubscribeCOVPropertyParse(io)
    case serviceChoice == 0x1D:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestGetEventInformationParse(io)
    case serviceChoice == 0x1E:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestSubscribeCOVPropertyMultipleParse(io)
    case serviceChoice == 0x1F:
        initializer, typeSwitchError = BACnetConfirmedServiceRequestConfirmedCOVNotificationMultipleParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Create the instance
    return initializer.initialize(), nil
}

func BACnetConfirmedServiceRequestSerialize(io utils.WriteBuffer, m BACnetConfirmedServiceRequest, i IBACnetConfirmedServiceRequest, childSerialize func() error) error {

    // Discriminator Field (serviceChoice) (Used as input to a switch field)
    serviceChoice := uint8(i.ServiceChoice())
    _serviceChoiceErr := io.WriteUint8(8, (serviceChoice))
    if _serviceChoiceErr != nil {
        return errors.New("Error serializing 'serviceChoice' field " + _serviceChoiceErr.Error())
    }

    // Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
    _typeSwitchErr := childSerialize()
    if _typeSwitchErr != nil {
        return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
    }

    return nil
}
