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
type APDUReject struct {
	originalInvokeId uint8
	rejectReason     uint8
	APDU
}

// The corresponding interface
type IAPDUReject interface {
	IAPDU
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m APDUReject) ApduType() uint8 {
	return 0x6
}

func (m APDUReject) initialize() spi.Message {
	return m
}

func NewAPDUReject(originalInvokeId uint8, rejectReason uint8) APDUInitializer {
	return &APDUReject{originalInvokeId: originalInvokeId, rejectReason: rejectReason}
}

func CastIAPDUReject(structType interface{}) IAPDUReject {
	castFunc := func(typ interface{}) IAPDUReject {
		if iAPDUReject, ok := typ.(IAPDUReject); ok {
			return iAPDUReject
		}
		return nil
	}
	return castFunc(structType)
}

func CastAPDUReject(structType interface{}) APDUReject {
	castFunc := func(typ interface{}) APDUReject {
		if sAPDUReject, ok := typ.(APDUReject); ok {
			return sAPDUReject
		}
		return APDUReject{}
	}
	return castFunc(structType)
}

func (m APDUReject) LengthInBits() uint16 {
	var lengthInBits uint16 = m.APDU.LengthInBits()

	// Reserved Field (reserved)
	lengthInBits += 4

	// Simple field (originalInvokeId)
	lengthInBits += 8

	// Simple field (rejectReason)
	lengthInBits += 8

	return lengthInBits
}

func (m APDUReject) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func APDURejectParse(io *spi.ReadBuffer) (APDUInitializer, error) {

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint8(4)
		if _err != nil {
			return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
		}
		if reserved != uint8(0x00) {
			log.WithFields(log.Fields{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (originalInvokeId)
	originalInvokeId, _originalInvokeIdErr := io.ReadUint8(8)
	if _originalInvokeIdErr != nil {
		return nil, errors.New("Error parsing 'originalInvokeId' field " + _originalInvokeIdErr.Error())
	}

	// Simple Field (rejectReason)
	rejectReason, _rejectReasonErr := io.ReadUint8(8)
	if _rejectReasonErr != nil {
		return nil, errors.New("Error parsing 'rejectReason' field " + _rejectReasonErr.Error())
	}

	// Create the instance
	return NewAPDUReject(originalInvokeId, rejectReason), nil
}

func (m APDUReject) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Reserved Field (reserved)
		{
			_err := io.WriteUint8(4, uint8(0x00))
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

		// Simple Field (rejectReason)
		rejectReason := uint8(m.rejectReason)
		_rejectReasonErr := io.WriteUint8(8, (rejectReason))
		if _rejectReasonErr != nil {
			return errors.New("Error serializing 'rejectReason' field " + _rejectReasonErr.Error())
		}

		return nil
	}
	return APDUSerialize(io, m.APDU, CastIAPDU(m), ser)
}
