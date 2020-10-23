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
type BACnetError struct {

}

// The corresponding interface
type IBACnetError interface {
    spi.Message
    ServiceChoice() uint8
    Serialize(io utils.WriteBuffer) error
}

type BACnetErrorInitializer interface {
    initialize() spi.Message
}

func BACnetErrorServiceChoice(m IBACnetError) uint8 {
    return m.ServiceChoice()
}


func CastIBACnetError(structType interface{}) IBACnetError {
    castFunc := func(typ interface{}) IBACnetError {
        if iBACnetError, ok := typ.(IBACnetError); ok {
            return iBACnetError
        }
        return nil
    }
    return castFunc(structType)
}

func CastBACnetError(structType interface{}) BACnetError {
    castFunc := func(typ interface{}) BACnetError {
        if sBACnetError, ok := typ.(BACnetError); ok {
            return sBACnetError
        }
        if sBACnetError, ok := typ.(*BACnetError); ok {
            return *sBACnetError
        }
        return BACnetError{}
    }
    return castFunc(structType)
}

func (m BACnetError) LengthInBits() uint16 {
    var lengthInBits uint16 = 0

    // Discriminator Field (serviceChoice)
    lengthInBits += 8

    // Length of sub-type elements will be added by sub-type...

    return lengthInBits
}

func (m BACnetError) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetErrorParse(io *utils.ReadBuffer) (spi.Message, error) {

    // Discriminator Field (serviceChoice) (Used as input to a switch field)
    serviceChoice, _serviceChoiceErr := io.ReadUint8(8)
    if _serviceChoiceErr != nil {
        return nil, errors.New("Error parsing 'serviceChoice' field " + _serviceChoiceErr.Error())
    }

    // Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
    var initializer BACnetErrorInitializer
    var typeSwitchError error
    switch {
    case serviceChoice == 0x03:
        initializer, typeSwitchError = BACnetErrorGetAlarmSummaryParse(io)
    case serviceChoice == 0x04:
        initializer, typeSwitchError = BACnetErrorGetEnrollmentSummaryParse(io)
    case serviceChoice == 0x1D:
        initializer, typeSwitchError = BACnetErrorGetEventInformationParse(io)
    case serviceChoice == 0x06:
        initializer, typeSwitchError = BACnetErrorAtomicReadFileParse(io)
    case serviceChoice == 0x07:
        initializer, typeSwitchError = BACnetErrorAtomicWriteFileParse(io)
    case serviceChoice == 0x0A:
        initializer, typeSwitchError = BACnetErrorCreateObjectParse(io)
    case serviceChoice == 0x0C:
        initializer, typeSwitchError = BACnetErrorReadPropertyParse(io)
    case serviceChoice == 0x0E:
        initializer, typeSwitchError = BACnetErrorReadPropertyMultipleParse(io)
    case serviceChoice == 0x1A:
        initializer, typeSwitchError = BACnetErrorReadRangeParse(io)
    case serviceChoice == 0x12:
        initializer, typeSwitchError = BACnetErrorConfirmedPrivateTransferParse(io)
    case serviceChoice == 0x15:
        initializer, typeSwitchError = BACnetErrorVTOpenParse(io)
    case serviceChoice == 0x17:
        initializer, typeSwitchError = BACnetErrorVTDataParse(io)
    case serviceChoice == 0x18:
        initializer, typeSwitchError = BACnetErrorRemovedAuthenticateParse(io)
    case serviceChoice == 0x0D:
        initializer, typeSwitchError = BACnetErrorRemovedReadPropertyConditionalParse(io)
    }
    if typeSwitchError != nil {
        return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
    }

    // Create the instance
    return initializer.initialize(), nil
}

func BACnetErrorSerialize(io utils.WriteBuffer, m BACnetError, i IBACnetError, childSerialize func() error) error {

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
