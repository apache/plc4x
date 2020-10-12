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
	log "github.com/sirupsen/logrus"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"reflect"
)

// The data-structure of this message
type APDUConfirmedRequest struct {
	segmentedMessage          bool
	moreFollows               bool
	segmentedResponseAccepted bool
	maxSegmentsAccepted       uint8
	maxApduLengthAccepted     uint8
	invokeId                  uint8
	sequenceNumber            *uint8
	proposedWindowSize        *uint8
	serviceRequest            BACnetConfirmedServiceRequest
	APDU
}

// The corresponding interface
type IAPDUConfirmedRequest interface {
	IAPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m APDUConfirmedRequest) ApduType() uint8 {
	return 0x0
}

func (m APDUConfirmedRequest) initialize() spi.Message {
	return m
}

func NewAPDUConfirmedRequest(segmentedMessage bool, moreFollows bool, segmentedResponseAccepted bool, maxSegmentsAccepted uint8, maxApduLengthAccepted uint8, invokeId uint8, sequenceNumber *uint8, proposedWindowSize *uint8, serviceRequest BACnetConfirmedServiceRequest) APDUInitializer {
	return &APDUConfirmedRequest{segmentedMessage: segmentedMessage, moreFollows: moreFollows, segmentedResponseAccepted: segmentedResponseAccepted, maxSegmentsAccepted: maxSegmentsAccepted, maxApduLengthAccepted: maxApduLengthAccepted, invokeId: invokeId, sequenceNumber: sequenceNumber, proposedWindowSize: proposedWindowSize, serviceRequest: serviceRequest}
}

func CastIAPDUConfirmedRequest(structType interface{}) IAPDUConfirmedRequest {
	castFunc := func(typ interface{}) IAPDUConfirmedRequest {
		if iAPDUConfirmedRequest, ok := typ.(IAPDUConfirmedRequest); ok {
			return iAPDUConfirmedRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastAPDUConfirmedRequest(structType interface{}) APDUConfirmedRequest {
	castFunc := func(typ interface{}) APDUConfirmedRequest {
		if sAPDUConfirmedRequest, ok := typ.(APDUConfirmedRequest); ok {
			return sAPDUConfirmedRequest
		}
		return APDUConfirmedRequest{}
	}
	return castFunc(structType)
}

func (m APDUConfirmedRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.APDU.LengthInBits()

	// Simple field (segmentedMessage)
	lengthInBits += 1

	// Simple field (moreFollows)
	lengthInBits += 1

	// Simple field (segmentedResponseAccepted)
	lengthInBits += 1

	// Reserved Field (reserved)
	lengthInBits += 2

	// Simple field (maxSegmentsAccepted)
	lengthInBits += 3

	// Simple field (maxApduLengthAccepted)
	lengthInBits += 4

	// Simple field (invokeId)
	lengthInBits += 8

	// Optional Field (sequenceNumber)
	if m.sequenceNumber != nil {
		lengthInBits += 8
	}

	// Optional Field (proposedWindowSize)
	if m.proposedWindowSize != nil {
		lengthInBits += 8
	}

	// Simple field (serviceRequest)
	lengthInBits += m.serviceRequest.LengthInBits()

	return lengthInBits
}

func (m APDUConfirmedRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func APDUConfirmedRequestParse(io spi.ReadBuffer, apduLength uint16) (APDUInitializer, error) {

	// Simple Field (segmentedMessage)
	var segmentedMessage bool = io.ReadBit()

	// Simple Field (moreFollows)
	var moreFollows bool = io.ReadBit()

	// Simple Field (segmentedResponseAccepted)
	var segmentedResponseAccepted bool = io.ReadBit()

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		var reserved uint8 = io.ReadUint8(2)
		if reserved != uint8(0) {
			log.WithFields(log.Fields{
				"expected value": uint8(0),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (maxSegmentsAccepted)
	var maxSegmentsAccepted uint8 = io.ReadUint8(3)

	// Simple Field (maxApduLengthAccepted)
	var maxApduLengthAccepted uint8 = io.ReadUint8(4)

	// Simple Field (invokeId)
	var invokeId uint8 = io.ReadUint8(8)

	// Optional Field (sequenceNumber) (Can be skipped, if a given expression evaluates to false)
	var sequenceNumber *uint8 = nil
	if segmentedMessage {
		_val := io.ReadUint8(8)
		sequenceNumber = &_val
	}

	// Optional Field (proposedWindowSize) (Can be skipped, if a given expression evaluates to false)
	var proposedWindowSize *uint8 = nil
	if segmentedMessage {
		_val := io.ReadUint8(8)
		proposedWindowSize = &_val
	}

	// Simple Field (serviceRequest)
	_serviceRequestMessage, _err := BACnetConfirmedServiceRequestParse(io, uint16(apduLength)-uint16(uint16(uint16(uint16(3))+uint16(uint16(spi.InlineIf(segmentedMessage, uint16(uint16(2)), uint16(uint16(0))))))))
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'serviceRequest'. " + _err.Error())
	}
	var serviceRequest BACnetConfirmedServiceRequest
	serviceRequest, _serviceRequestOk := _serviceRequestMessage.(BACnetConfirmedServiceRequest)
	if !_serviceRequestOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_serviceRequestMessage).Name() + " to BACnetConfirmedServiceRequest")
	}

	// Create the instance
	return NewAPDUConfirmedRequest(segmentedMessage, moreFollows, segmentedResponseAccepted, maxSegmentsAccepted, maxApduLengthAccepted, invokeId, sequenceNumber, proposedWindowSize, serviceRequest), nil
}

func (m APDUConfirmedRequest) Serialize(io spi.WriteBuffer) {

	// Simple Field (segmentedMessage)
	segmentedMessage := bool(m.segmentedMessage)
	io.WriteBit((bool)(segmentedMessage))

	// Simple Field (moreFollows)
	moreFollows := bool(m.moreFollows)
	io.WriteBit((bool)(moreFollows))

	// Simple Field (segmentedResponseAccepted)
	segmentedResponseAccepted := bool(m.segmentedResponseAccepted)
	io.WriteBit((bool)(segmentedResponseAccepted))

	// Reserved Field (reserved)
	io.WriteUint8(2, uint8(0))

	// Simple Field (maxSegmentsAccepted)
	maxSegmentsAccepted := uint8(m.maxSegmentsAccepted)
	io.WriteUint8(3, (maxSegmentsAccepted))

	// Simple Field (maxApduLengthAccepted)
	maxApduLengthAccepted := uint8(m.maxApduLengthAccepted)
	io.WriteUint8(4, (maxApduLengthAccepted))

	// Simple Field (invokeId)
	invokeId := uint8(m.invokeId)
	io.WriteUint8(8, (invokeId))

	// Optional Field (sequenceNumber) (Can be skipped, if the value is null)
	var sequenceNumber *uint8 = nil
	if m.sequenceNumber != nil {
		sequenceNumber = m.sequenceNumber
		io.WriteUint8(8, *(sequenceNumber))
	}

	// Optional Field (proposedWindowSize) (Can be skipped, if the value is null)
	var proposedWindowSize *uint8 = nil
	if m.proposedWindowSize != nil {
		proposedWindowSize = m.proposedWindowSize
		io.WriteUint8(8, *(proposedWindowSize))
	}

	// Simple Field (serviceRequest)
	serviceRequest := BACnetConfirmedServiceRequest(m.serviceRequest)
	serviceRequest.Serialize(io)
}
