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

func APDUSegmentAckParse(io spi.ReadBuffer) (APDUInitializer, error) {

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		var reserved uint8 = io.ReadUint8(2)
		if reserved != uint8(0x00) {
			log.WithFields(log.Fields{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (negativeAck)
	var negativeAck bool = io.ReadBit()

	// Simple Field (server)
	var server bool = io.ReadBit()

	// Simple Field (originalInvokeId)
	var originalInvokeId uint8 = io.ReadUint8(8)

	// Simple Field (sequenceNumber)
	var sequenceNumber uint8 = io.ReadUint8(8)

	// Simple Field (proposedWindowSize)
	var proposedWindowSize uint8 = io.ReadUint8(8)

	// Create the instance
	return NewAPDUSegmentAck(negativeAck, server, originalInvokeId, sequenceNumber, proposedWindowSize), nil
}

func (m APDUSegmentAck) Serialize(io spi.WriteBuffer) {

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
