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
type BACnetUnconfirmedServiceRequest struct {
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequest interface {
	spi.Message
	ServiceChoice() uint8
	Serialize(io spi.WriteBuffer)
}

type BACnetUnconfirmedServiceRequestInitializer interface {
	initialize() spi.Message
}

func BACnetUnconfirmedServiceRequestServiceChoice(m IBACnetUnconfirmedServiceRequest) uint8 {
	return m.ServiceChoice()
}

func CastIBACnetUnconfirmedServiceRequest(structType interface{}) IBACnetUnconfirmedServiceRequest {
	castFunc := func(typ interface{}) IBACnetUnconfirmedServiceRequest {
		if iBACnetUnconfirmedServiceRequest, ok := typ.(IBACnetUnconfirmedServiceRequest); ok {
			return iBACnetUnconfirmedServiceRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetUnconfirmedServiceRequest(structType interface{}) BACnetUnconfirmedServiceRequest {
	castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequest {
		if sBACnetUnconfirmedServiceRequest, ok := typ.(BACnetUnconfirmedServiceRequest); ok {
			return sBACnetUnconfirmedServiceRequest
		}
		return BACnetUnconfirmedServiceRequest{}
	}
	return castFunc(structType)
}

func (m BACnetUnconfirmedServiceRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (serviceChoice)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m BACnetUnconfirmedServiceRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestParse(io spi.ReadBuffer, len uint16) (spi.Message, error) {

	// Discriminator Field (serviceChoice) (Used as input to a switch field)
	var serviceChoice uint8 = io.ReadUint8(8)

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer BACnetUnconfirmedServiceRequestInitializer
	var typeSwitchError error
	switch {
	case serviceChoice == 0x00:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestIAmParse(io)
	case serviceChoice == 0x01:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestIHaveParse(io)
	case serviceChoice == 0x02:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationParse(io)
	case serviceChoice == 0x03:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestUnconfirmedEventNotificationParse(io)
	case serviceChoice == 0x04:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestUnconfirmedPrivateTransferParse(io, len)
	case serviceChoice == 0x05:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestUnconfirmedTextMessageParse(io)
	case serviceChoice == 0x06:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestTimeSynchronizationParse(io)
	case serviceChoice == 0x07:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestWhoHasParse(io)
	case serviceChoice == 0x08:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestWhoIsParse(io)
	case serviceChoice == 0x09:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestUTCTimeSynchronizationParse(io)
	case serviceChoice == 0x0A:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestWriteGroupParse(io)
	case serviceChoice == 0x0B:
		initializer, typeSwitchError = BACnetUnconfirmedServiceRequestUnconfirmedCOVNotificationMultipleParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func (m BACnetUnconfirmedServiceRequest) Serialize(io spi.WriteBuffer) {
	iBACnetUnconfirmedServiceRequest := CastIBACnetUnconfirmedServiceRequest(m)

	// Discriminator Field (serviceChoice) (Used as input to a switch field)
	serviceChoice := uint8(BACnetUnconfirmedServiceRequestServiceChoice(iBACnetUnconfirmedServiceRequest))
	io.WriteUint8(8, (serviceChoice))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	iBACnetUnconfirmedServiceRequest.Serialize(io)
}
