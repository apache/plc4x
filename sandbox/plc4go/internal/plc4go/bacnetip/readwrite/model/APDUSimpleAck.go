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
)

// The data-structure of this message
type APDUSimpleAck struct {
	originalInvokeId uint8
	serviceChoice    uint8
	APDU
}

// The corresponding interface
type IAPDUSimpleAck interface {
	IAPDU
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m APDUSimpleAck) ApduType() uint8 {
	return 0x2
}

func (m APDUSimpleAck) initialize() spi.Message {
	return m
}

func NewAPDUSimpleAck(originalInvokeId uint8, serviceChoice uint8) APDUInitializer {
	return &APDUSimpleAck{originalInvokeId: originalInvokeId, serviceChoice: serviceChoice}
}

func CastIAPDUSimpleAck(structType interface{}) IAPDUSimpleAck {
	castFunc := func(typ interface{}) IAPDUSimpleAck {
		if iAPDUSimpleAck, ok := typ.(IAPDUSimpleAck); ok {
			return iAPDUSimpleAck
		}
		return nil
	}
	return castFunc(structType)
}

func CastAPDUSimpleAck(structType interface{}) APDUSimpleAck {
	castFunc := func(typ interface{}) APDUSimpleAck {
		if sAPDUSimpleAck, ok := typ.(APDUSimpleAck); ok {
			return sAPDUSimpleAck
		}
		return APDUSimpleAck{}
	}
	return castFunc(structType)
}

func (m APDUSimpleAck) LengthInBits() uint16 {
	var lengthInBits uint16 = m.APDU.LengthInBits()

	// Reserved Field (reserved)
	lengthInBits += 4

	// Simple field (originalInvokeId)
	lengthInBits += 8

	// Simple field (serviceChoice)
	lengthInBits += 8

	return lengthInBits
}

func (m APDUSimpleAck) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func APDUSimpleAckParse(io *spi.ReadBuffer) (APDUInitializer, error) {

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

	// Simple Field (originalInvokeId)
	originalInvokeId, _originalInvokeIdErr := io.ReadUint8(8)
	if _originalInvokeIdErr != nil {
		return nil, errors.New("Error parsing 'originalInvokeId' field " + _originalInvokeIdErr.Error())
	}

	// Simple Field (serviceChoice)
	serviceChoice, _serviceChoiceErr := io.ReadUint8(8)
	if _serviceChoiceErr != nil {
		return nil, errors.New("Error parsing 'serviceChoice' field " + _serviceChoiceErr.Error())
	}

	// Create the instance
	return NewAPDUSimpleAck(originalInvokeId, serviceChoice), nil
}

func (m APDUSimpleAck) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Reserved Field (reserved)
		{
			_err := io.WriteUint8(4, uint8(0))
			if _err != nil {
				return errors.New("Error serializing 'reserved' field " + _err.Error())
			}
		}

		// Simple Field (originalInvokeId)
		originalInvokeId := uint8(m.originalInvokeId)
		_originalInvokeIdErr := io.WriteUint8(8, (originalInvokeId))
		if _originalInvokeIdErr != nil {
			return errors.New("Error serializing 'originalInvokeId' field " + _originalInvokeIdErr.Error())
		}

		// Simple Field (serviceChoice)
		serviceChoice := uint8(m.serviceChoice)
		_serviceChoiceErr := io.WriteUint8(8, (serviceChoice))
		if _serviceChoiceErr != nil {
			return errors.New("Error serializing 'serviceChoice' field " + _serviceChoiceErr.Error())
		}

		return nil
	}
	return APDUSerialize(io, m.APDU, CastIAPDU(m), ser)
}
