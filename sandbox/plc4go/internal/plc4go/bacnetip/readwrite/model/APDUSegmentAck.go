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
type APDUSegmentAck struct {
	negativeAck        bool
	server             bool
	originalInvokeId   uint8
	sequenceNumber     uint8
	proposedWindowSize uint8
	APDU
}

// The corresponding interface
type IAPDUSegmentAck interface {
	IAPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m APDUSegmentAck) ApduType() uint8 {
	return 0x4
}

func (m APDUSegmentAck) initialize() spi.Message {
	return m
}

func NewAPDUSegmentAck(negativeAck bool, server bool, originalInvokeId uint8, sequenceNumber uint8, proposedWindowSize uint8) APDUInitializer {
	return &APDUSegmentAck{negativeAck: negativeAck, server: server, originalInvokeId: originalInvokeId, sequenceNumber: sequenceNumber, proposedWindowSize: proposedWindowSize}
}

func CastIAPDUSegmentAck(structType interface{}) IAPDUSegmentAck {
	castFunc := func(typ interface{}) IAPDUSegmentAck {
		if iAPDUSegmentAck, ok := typ.(IAPDUSegmentAck); ok {
			return iAPDUSegmentAck
		}
		return nil
	}
	return castFunc(structType)
}

func CastAPDUSegmentAck(structType interface{}) APDUSegmentAck {
	castFunc := func(typ interface{}) APDUSegmentAck {
		if sAPDUSegmentAck, ok := typ.(APDUSegmentAck); ok {
			return sAPDUSegmentAck
		}
		return APDUSegmentAck{}
	}
	return castFunc(structType)
}

func (m APDUSegmentAck) LengthInBits() uint16 {
	var lengthInBits uint16 = m.APDU.LengthInBits()

	// Reserved Field (reserved)
	lengthInBits += 2

	// Simple field (negativeAck)
	lengthInBits += 1

	// Simple field (server)
	lengthInBits += 1

	// Simple field (originalInvokeId)
	lengthInBits += 8

	// Simple field (sequenceNumber)
	lengthInBits += 8

	// Simple field (proposedWindowSize)
	lengthInBits += 8

	return lengthInBits
}

func (m APDUSegmentAck) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func APDUSegmentAckParse(io *spi.ReadBuffer) (APDUInitializer, error) {

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint8(2)
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

	// Simple Field (negativeAck)
	negativeAck, _negativeAckErr := io.ReadBit()
	if _negativeAckErr != nil {
		return nil, errors.New("Error parsing 'negativeAck' field " + _negativeAckErr.Error())
	}

	// Simple Field (server)
	server, _serverErr := io.ReadBit()
	if _serverErr != nil {
		return nil, errors.New("Error parsing 'server' field " + _serverErr.Error())
	}

	// Simple Field (originalInvokeId)
	originalInvokeId, _originalInvokeIdErr := io.ReadUint8(8)
	if _originalInvokeIdErr != nil {
		return nil, errors.New("Error parsing 'originalInvokeId' field " + _originalInvokeIdErr.Error())
	}

	// Simple Field (sequenceNumber)
	sequenceNumber, _sequenceNumberErr := io.ReadUint8(8)
	if _sequenceNumberErr != nil {
		return nil, errors.New("Error parsing 'sequenceNumber' field " + _sequenceNumberErr.Error())
	}

	// Simple Field (proposedWindowSize)
	proposedWindowSize, _proposedWindowSizeErr := io.ReadUint8(8)
	if _proposedWindowSizeErr != nil {
		return nil, errors.New("Error parsing 'proposedWindowSize' field " + _proposedWindowSizeErr.Error())
	}

	// Create the instance
	return NewAPDUSegmentAck(negativeAck, server, originalInvokeId, sequenceNumber, proposedWindowSize), nil
}

func (m APDUSegmentAck) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Reserved Field (reserved)
		io.WriteUint8(2, uint8(0x00))

		// Simple Field (negativeAck)
		negativeAck := bool(m.negativeAck)
		io.WriteBit((bool)(negativeAck))

		// Simple Field (server)
		server := bool(m.server)
		io.WriteBit((bool)(server))

		// Simple Field (originalInvokeId)
		originalInvokeId := uint8(m.originalInvokeId)
		io.WriteUint8(8, (originalInvokeId))

		// Simple Field (sequenceNumber)
		sequenceNumber := uint8(m.sequenceNumber)
		io.WriteUint8(8, (sequenceNumber))

		// Simple Field (proposedWindowSize)
		proposedWindowSize := uint8(m.proposedWindowSize)
		io.WriteUint8(8, (proposedWindowSize))

	}
	APDUSerialize(io, m.APDU, CastIAPDU(m), ser)
}
