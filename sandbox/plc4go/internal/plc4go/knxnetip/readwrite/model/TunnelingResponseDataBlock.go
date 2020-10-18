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
type TunnelingResponseDataBlock struct {
	CommunicationChannelId uint8
	SequenceCounter        uint8
	Status                 IStatus
}

// The corresponding interface
type ITunnelingResponseDataBlock interface {
	spi.Message
	Serialize(io spi.WriteBuffer) error
}

func NewTunnelingResponseDataBlock(communicationChannelId uint8, sequenceCounter uint8, status IStatus) spi.Message {
	return &TunnelingResponseDataBlock{CommunicationChannelId: communicationChannelId, SequenceCounter: sequenceCounter, Status: status}
}

func CastITunnelingResponseDataBlock(structType interface{}) ITunnelingResponseDataBlock {
	castFunc := func(typ interface{}) ITunnelingResponseDataBlock {
		if iTunnelingResponseDataBlock, ok := typ.(ITunnelingResponseDataBlock); ok {
			return iTunnelingResponseDataBlock
		}
		return nil
	}
	return castFunc(structType)
}

func CastTunnelingResponseDataBlock(structType interface{}) TunnelingResponseDataBlock {
	castFunc := func(typ interface{}) TunnelingResponseDataBlock {
		if sTunnelingResponseDataBlock, ok := typ.(TunnelingResponseDataBlock); ok {
			return sTunnelingResponseDataBlock
		}
		return TunnelingResponseDataBlock{}
	}
	return castFunc(structType)
}

func (m TunnelingResponseDataBlock) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Implicit Field (structureLength)
	lengthInBits += 8

	// Simple field (communicationChannelId)
	lengthInBits += 8

	// Simple field (sequenceCounter)
	lengthInBits += 8

	// Enum Field (status)
	lengthInBits += 8

	return lengthInBits
}

func (m TunnelingResponseDataBlock) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func TunnelingResponseDataBlockParse(io *spi.ReadBuffer) (spi.Message, error) {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	_, _structureLengthErr := io.ReadUint8(8)
	if _structureLengthErr != nil {
		return nil, errors.New("Error parsing 'structureLength' field " + _structureLengthErr.Error())
	}

	// Simple Field (communicationChannelId)
	communicationChannelId, _communicationChannelIdErr := io.ReadUint8(8)
	if _communicationChannelIdErr != nil {
		return nil, errors.New("Error parsing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
	}

	// Simple Field (sequenceCounter)
	sequenceCounter, _sequenceCounterErr := io.ReadUint8(8)
	if _sequenceCounterErr != nil {
		return nil, errors.New("Error parsing 'sequenceCounter' field " + _sequenceCounterErr.Error())
	}

	// Enum field (status)
	status, _statusErr := StatusParse(io)
	if _statusErr != nil {
		return nil, errors.New("Error parsing 'status' field " + _statusErr.Error())
	}

	// Create the instance
	return NewTunnelingResponseDataBlock(communicationChannelId, sequenceCounter, status), nil
}

func (m TunnelingResponseDataBlock) Serialize(io spi.WriteBuffer) error {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	structureLength := uint8(uint8(m.LengthInBytes()))
	_structureLengthErr := io.WriteUint8(8, structureLength)
	if _structureLengthErr != nil {
		return errors.New("Error serializing 'structureLength' field " + _structureLengthErr.Error())
	}

	// Simple Field (communicationChannelId)
	communicationChannelId := uint8(m.CommunicationChannelId)
	_communicationChannelIdErr := io.WriteUint8(8, communicationChannelId)
	if _communicationChannelIdErr != nil {
		return errors.New("Error serializing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
	}

	// Simple Field (sequenceCounter)
	sequenceCounter := uint8(m.SequenceCounter)
	_sequenceCounterErr := io.WriteUint8(8, sequenceCounter)
	if _sequenceCounterErr != nil {
		return errors.New("Error serializing 'sequenceCounter' field " + _sequenceCounterErr.Error())
	}

	// Enum field (status)
	status := CastStatus(m.Status)
	_statusErr := status.Serialize(io)
	if _statusErr != nil {
		return errors.New("Error serializing 'status' field " + _statusErr.Error())
	}

	return nil
}
