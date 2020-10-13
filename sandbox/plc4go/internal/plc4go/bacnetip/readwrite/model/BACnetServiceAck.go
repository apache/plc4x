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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type BACnetServiceAck struct {
}

// The corresponding interface
type IBACnetServiceAck interface {
	spi.Message
	ServiceChoice() uint8
	Serialize(io spi.WriteBuffer)
}

type BACnetServiceAckInitializer interface {
	initialize() spi.Message
}

func BACnetServiceAckServiceChoice(m IBACnetServiceAck) uint8 {
	return m.ServiceChoice()
}

func CastIBACnetServiceAck(structType interface{}) IBACnetServiceAck {
	castFunc := func(typ interface{}) IBACnetServiceAck {
		if iBACnetServiceAck, ok := typ.(IBACnetServiceAck); ok {
			return iBACnetServiceAck
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetServiceAck(structType interface{}) BACnetServiceAck {
	castFunc := func(typ interface{}) BACnetServiceAck {
		if sBACnetServiceAck, ok := typ.(BACnetServiceAck); ok {
			return sBACnetServiceAck
		}
		return BACnetServiceAck{}
	}
	return castFunc(structType)
}

func (m BACnetServiceAck) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (serviceChoice)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m BACnetServiceAck) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetServiceAckParse(io spi.ReadBuffer) (spi.Message, error) {

	// Discriminator Field (serviceChoice) (Used as input to a switch field)
	serviceChoice, _serviceChoiceErr := io.ReadUint8(8)
	if _serviceChoiceErr != nil {
		return nil, errors.New("Error parsing 'serviceChoice' field " + _serviceChoiceErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer BACnetServiceAckInitializer
	var typeSwitchError error
	switch {
	case serviceChoice == 0x03:
		initializer, typeSwitchError = BACnetServiceAckGetAlarmSummaryParse(io)
	case serviceChoice == 0x04:
		initializer, typeSwitchError = BACnetServiceAckGetEnrollmentSummaryParse(io)
	case serviceChoice == 0x1D:
		initializer, typeSwitchError = BACnetServiceAckGetEventInformationParse(io)
	case serviceChoice == 0x06:
		initializer, typeSwitchError = BACnetServiceAckAtomicReadFileParse(io)
	case serviceChoice == 0x07:
		initializer, typeSwitchError = BACnetServiceAckAtomicWriteFileParse(io)
	case serviceChoice == 0x0A:
		initializer, typeSwitchError = BACnetServiceAckCreateObjectParse(io)
	case serviceChoice == 0x0C:
		initializer, typeSwitchError = BACnetServiceAckReadPropertyParse(io)
	case serviceChoice == 0x0E:
		initializer, typeSwitchError = BACnetServiceAckReadPropertyMultipleParse(io)
	case serviceChoice == 0x1A:
		initializer, typeSwitchError = BACnetServiceAckReadRangeParse(io)
	case serviceChoice == 0x12:
		initializer, typeSwitchError = BACnetServiceAckConfirmedPrivateTransferParse(io)
	case serviceChoice == 0x15:
		initializer, typeSwitchError = BACnetServiceAckVTOpenParse(io)
	case serviceChoice == 0x17:
		initializer, typeSwitchError = BACnetServiceAckVTDataParse(io)
	case serviceChoice == 0x18:
		initializer, typeSwitchError = BACnetServiceAckRemovedAuthenticateParse(io)
	case serviceChoice == 0x0D:
		initializer, typeSwitchError = BACnetServiceAckRemovedReadPropertyConditionalParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func (m BACnetServiceAck) Serialize(io spi.WriteBuffer) {
	iBACnetServiceAck := CastIBACnetServiceAck(m)

	// Discriminator Field (serviceChoice) (Used as input to a switch field)
	serviceChoice := uint8(BACnetServiceAckServiceChoice(iBACnetServiceAck))
	io.WriteUint8(8, (serviceChoice))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	iBACnetServiceAck.Serialize(io)
}
