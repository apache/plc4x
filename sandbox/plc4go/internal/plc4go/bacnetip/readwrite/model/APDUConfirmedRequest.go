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
	serviceRequest            IBACnetConfirmedServiceRequest
	APDU
}

// The corresponding interface
type IAPDUConfirmedRequest interface {
	IAPDU
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m APDUConfirmedRequest) ApduType() uint8 {
	return 0x0
}

func (m APDUConfirmedRequest) initialize() spi.Message {
	return m
}

func NewAPDUConfirmedRequest(segmentedMessage bool, moreFollows bool, segmentedResponseAccepted bool, maxSegmentsAccepted uint8, maxApduLengthAccepted uint8, invokeId uint8, sequenceNumber *uint8, proposedWindowSize *uint8, serviceRequest IBACnetConfirmedServiceRequest) APDUInitializer {
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

func APDUConfirmedRequestParse(io *spi.ReadBuffer, apduLength uint16) (APDUInitializer, error) {

	// Simple Field (segmentedMessage)
	segmentedMessage, _segmentedMessageErr := io.ReadBit()
	if _segmentedMessageErr != nil {
		return nil, errors.New("Error parsing 'segmentedMessage' field " + _segmentedMessageErr.Error())
	}

	// Simple Field (moreFollows)
	moreFollows, _moreFollowsErr := io.ReadBit()
	if _moreFollowsErr != nil {
		return nil, errors.New("Error parsing 'moreFollows' field " + _moreFollowsErr.Error())
	}

	// Simple Field (segmentedResponseAccepted)
	segmentedResponseAccepted, _segmentedResponseAcceptedErr := io.ReadBit()
	if _segmentedResponseAcceptedErr != nil {
		return nil, errors.New("Error parsing 'segmentedResponseAccepted' field " + _segmentedResponseAcceptedErr.Error())
	}

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint8(2)
		if _err != nil {
			return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
		}
		if reserved != uint8(0) {
			log.WithFields(log.Fields{
				"expected value": uint8(0),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (maxSegmentsAccepted)
	maxSegmentsAccepted, _maxSegmentsAcceptedErr := io.ReadUint8(3)
	if _maxSegmentsAcceptedErr != nil {
		return nil, errors.New("Error parsing 'maxSegmentsAccepted' field " + _maxSegmentsAcceptedErr.Error())
	}

	// Simple Field (maxApduLengthAccepted)
	maxApduLengthAccepted, _maxApduLengthAcceptedErr := io.ReadUint8(4)
	if _maxApduLengthAcceptedErr != nil {
		return nil, errors.New("Error parsing 'maxApduLengthAccepted' field " + _maxApduLengthAcceptedErr.Error())
	}

	// Simple Field (invokeId)
	invokeId, _invokeIdErr := io.ReadUint8(8)
	if _invokeIdErr != nil {
		return nil, errors.New("Error parsing 'invokeId' field " + _invokeIdErr.Error())
	}

	// Optional Field (sequenceNumber) (Can be skipped, if a given expression evaluates to false)
	var sequenceNumber *uint8 = nil
	if segmentedMessage {
		_val, _err := io.ReadUint8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'sequenceNumber' field " + _err.Error())
		}

		sequenceNumber = &_val
	}

	// Optional Field (proposedWindowSize) (Can be skipped, if a given expression evaluates to false)
	var proposedWindowSize *uint8 = nil
	if segmentedMessage {
		_val, _err := io.ReadUint8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'proposedWindowSize' field " + _err.Error())
		}

		proposedWindowSize = &_val
	}

	// Simple Field (serviceRequest)
	_serviceRequestMessage, _err := BACnetConfirmedServiceRequestParse(io, uint16(apduLength)-uint16(uint16(uint16(uint16(3))+uint16(uint16(spi.InlineIf(segmentedMessage, uint16(uint16(2)), uint16(uint16(0))))))))
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'serviceRequest'. " + _err.Error())
	}
	var serviceRequest IBACnetConfirmedServiceRequest
	serviceRequest, _serviceRequestOk := _serviceRequestMessage.(IBACnetConfirmedServiceRequest)
	if !_serviceRequestOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_serviceRequestMessage).Name() + " to IBACnetConfirmedServiceRequest")
	}

	// Create the instance
	return NewAPDUConfirmedRequest(segmentedMessage, moreFollows, segmentedResponseAccepted, maxSegmentsAccepted, maxApduLengthAccepted, invokeId, sequenceNumber, proposedWindowSize, serviceRequest), nil
}

func (m APDUConfirmedRequest) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (segmentedMessage)
		segmentedMessage := bool(m.segmentedMessage)
		_segmentedMessageErr := io.WriteBit((bool)(segmentedMessage))
		if _segmentedMessageErr != nil {
			return errors.New("Error serializing 'segmentedMessage' field " + _segmentedMessageErr.Error())
		}

		// Simple Field (moreFollows)
		moreFollows := bool(m.moreFollows)
		_moreFollowsErr := io.WriteBit((bool)(moreFollows))
		if _moreFollowsErr != nil {
			return errors.New("Error serializing 'moreFollows' field " + _moreFollowsErr.Error())
		}

		// Simple Field (segmentedResponseAccepted)
		segmentedResponseAccepted := bool(m.segmentedResponseAccepted)
		_segmentedResponseAcceptedErr := io.WriteBit((bool)(segmentedResponseAccepted))
		if _segmentedResponseAcceptedErr != nil {
			return errors.New("Error serializing 'segmentedResponseAccepted' field " + _segmentedResponseAcceptedErr.Error())
		}

		// Reserved Field (reserved)
		{
			_err := io.WriteUint8(2, uint8(0))
			if _err != nil {
				return errors.New("Error serializing 'reserved' field " + _err.Error())
			}
		}

		// Simple Field (maxSegmentsAccepted)
		maxSegmentsAccepted := uint8(m.maxSegmentsAccepted)
		_maxSegmentsAcceptedErr := io.WriteUint8(3, (maxSegmentsAccepted))
		if _maxSegmentsAcceptedErr != nil {
			return errors.New("Error serializing 'maxSegmentsAccepted' field " + _maxSegmentsAcceptedErr.Error())
		}

		// Simple Field (maxApduLengthAccepted)
		maxApduLengthAccepted := uint8(m.maxApduLengthAccepted)
		_maxApduLengthAcceptedErr := io.WriteUint8(4, (maxApduLengthAccepted))
		if _maxApduLengthAcceptedErr != nil {
			return errors.New("Error serializing 'maxApduLengthAccepted' field " + _maxApduLengthAcceptedErr.Error())
		}

		// Simple Field (invokeId)
		invokeId := uint8(m.invokeId)
		_invokeIdErr := io.WriteUint8(8, (invokeId))
		if _invokeIdErr != nil {
			return errors.New("Error serializing 'invokeId' field " + _invokeIdErr.Error())
		}

		// Optional Field (sequenceNumber) (Can be skipped, if the value is null)
		var sequenceNumber *uint8 = nil
		if m.sequenceNumber != nil {
			sequenceNumber = m.sequenceNumber
			_sequenceNumberErr := io.WriteUint8(8, *(sequenceNumber))
			if _sequenceNumberErr != nil {
				return errors.New("Error serializing 'sequenceNumber' field " + _sequenceNumberErr.Error())
			}
		}

		// Optional Field (proposedWindowSize) (Can be skipped, if the value is null)
		var proposedWindowSize *uint8 = nil
		if m.proposedWindowSize != nil {
			proposedWindowSize = m.proposedWindowSize
			_proposedWindowSizeErr := io.WriteUint8(8, *(proposedWindowSize))
			if _proposedWindowSizeErr != nil {
				return errors.New("Error serializing 'proposedWindowSize' field " + _proposedWindowSizeErr.Error())
			}
		}

		// Simple Field (serviceRequest)
		serviceRequest := CastIBACnetConfirmedServiceRequest(m.serviceRequest)
		_serviceRequestErr := serviceRequest.Serialize(io)
		if _serviceRequestErr != nil {
			return errors.New("Error serializing 'serviceRequest' field " + _serviceRequestErr.Error())
		}

		return nil
	}
	return APDUSerialize(io, m.APDU, CastIAPDU(m), ser)
}
