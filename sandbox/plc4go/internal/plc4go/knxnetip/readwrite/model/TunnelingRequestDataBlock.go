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
type TunnelingRequestDataBlock struct {
	communicationChannelId uint8
	sequenceCounter        uint8
}

// The corresponding interface
type ITunnelingRequestDataBlock interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewTunnelingRequestDataBlock(communicationChannelId uint8, sequenceCounter uint8) spi.Message {
	return &TunnelingRequestDataBlock{communicationChannelId: communicationChannelId, sequenceCounter: sequenceCounter}
}

func CastITunnelingRequestDataBlock(structType interface{}) ITunnelingRequestDataBlock {
	castFunc := func(typ interface{}) ITunnelingRequestDataBlock {
		if iTunnelingRequestDataBlock, ok := typ.(ITunnelingRequestDataBlock); ok {
			return iTunnelingRequestDataBlock
		}
		return nil
	}
	return castFunc(structType)
}

func CastTunnelingRequestDataBlock(structType interface{}) TunnelingRequestDataBlock {
	castFunc := func(typ interface{}) TunnelingRequestDataBlock {
		if sTunnelingRequestDataBlock, ok := typ.(TunnelingRequestDataBlock); ok {
			return sTunnelingRequestDataBlock
		}
		return TunnelingRequestDataBlock{}
	}
	return castFunc(structType)
}

func (m TunnelingRequestDataBlock) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Implicit Field (structureLength)
	lengthInBits += 8

	// Simple field (communicationChannelId)
	lengthInBits += 8

	// Simple field (sequenceCounter)
	lengthInBits += 8

	// Reserved Field (reserved)
	lengthInBits += 8

	return lengthInBits
}

func (m TunnelingRequestDataBlock) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func TunnelingRequestDataBlockParse(io spi.ReadBuffer) (spi.Message, error) {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var _ uint8 = io.ReadUint8(8)

	// Simple Field (communicationChannelId)
	var communicationChannelId uint8 = io.ReadUint8(8)

	// Simple Field (sequenceCounter)
	var sequenceCounter uint8 = io.ReadUint8(8)

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		var reserved uint8 = io.ReadUint8(8)
		if reserved != uint8(0x00) {
			log.WithFields(log.Fields{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Create the instance
	return NewTunnelingRequestDataBlock(communicationChannelId, sequenceCounter), nil
}

func (m TunnelingRequestDataBlock) Serialize(io spi.WriteBuffer) {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	structureLength := uint8(uint8(m.LengthInBytes()))
	io.WriteUint8(8, (structureLength))

	// Simple Field (communicationChannelId)
	communicationChannelId := uint8(m.communicationChannelId)
	io.WriteUint8(8, (communicationChannelId))

	// Simple Field (sequenceCounter)
	sequenceCounter := uint8(m.sequenceCounter)
	io.WriteUint8(8, (sequenceCounter))

	// Reserved Field (reserved)
	io.WriteUint8(8, uint8(0x00))
}
