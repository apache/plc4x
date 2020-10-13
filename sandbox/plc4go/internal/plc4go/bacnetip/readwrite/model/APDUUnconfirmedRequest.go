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
type APDUUnconfirmedRequest struct {
	serviceRequest IBACnetUnconfirmedServiceRequest
	APDU
}

// The corresponding interface
type IAPDUUnconfirmedRequest interface {
	IAPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m APDUUnconfirmedRequest) ApduType() uint8 {
	return 0x1
}

func (m APDUUnconfirmedRequest) initialize() spi.Message {
	return m
}

func NewAPDUUnconfirmedRequest(serviceRequest IBACnetUnconfirmedServiceRequest) APDUInitializer {
	return &APDUUnconfirmedRequest{serviceRequest: serviceRequest}
}

func CastIAPDUUnconfirmedRequest(structType interface{}) IAPDUUnconfirmedRequest {
	castFunc := func(typ interface{}) IAPDUUnconfirmedRequest {
		if iAPDUUnconfirmedRequest, ok := typ.(IAPDUUnconfirmedRequest); ok {
			return iAPDUUnconfirmedRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastAPDUUnconfirmedRequest(structType interface{}) APDUUnconfirmedRequest {
	castFunc := func(typ interface{}) APDUUnconfirmedRequest {
		if sAPDUUnconfirmedRequest, ok := typ.(APDUUnconfirmedRequest); ok {
			return sAPDUUnconfirmedRequest
		}
		return APDUUnconfirmedRequest{}
	}
	return castFunc(structType)
}

func (m APDUUnconfirmedRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.APDU.LengthInBits()

	// Reserved Field (reserved)
	lengthInBits += 4

	// Simple field (serviceRequest)
	lengthInBits += m.serviceRequest.LengthInBits()

	return lengthInBits
}

func (m APDUUnconfirmedRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func APDUUnconfirmedRequestParse(io *spi.ReadBuffer, apduLength uint16) (APDUInitializer, error) {

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint8(4)
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

	// Simple Field (serviceRequest)
	_serviceRequestMessage, _err := BACnetUnconfirmedServiceRequestParse(io, uint16(apduLength)-uint16(uint16(1)))
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'serviceRequest'. " + _err.Error())
	}
	var serviceRequest IBACnetUnconfirmedServiceRequest
	serviceRequest, _serviceRequestOk := _serviceRequestMessage.(IBACnetUnconfirmedServiceRequest)
	if !_serviceRequestOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_serviceRequestMessage).Name() + " to IBACnetUnconfirmedServiceRequest")
	}

	// Create the instance
	return NewAPDUUnconfirmedRequest(serviceRequest), nil
}

func (m APDUUnconfirmedRequest) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Reserved Field (reserved)
		io.WriteUint8(4, uint8(0))

		// Simple Field (serviceRequest)
		serviceRequest := CastIBACnetUnconfirmedServiceRequest(m.serviceRequest)
		serviceRequest.Serialize(io)

	}
	APDUSerialize(io, m.APDU, CastIAPDU(m), ser)
}
