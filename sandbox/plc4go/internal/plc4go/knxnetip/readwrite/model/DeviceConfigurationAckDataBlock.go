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
type DeviceConfigurationAckDataBlock struct {
	communicationChannelId uint8
	sequenceCounter        uint8
	status                 Status
}

// The corresponding interface
type IDeviceConfigurationAckDataBlock interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewDeviceConfigurationAckDataBlock(communicationChannelId uint8, sequenceCounter uint8, status Status) spi.Message {
	return &DeviceConfigurationAckDataBlock{communicationChannelId: communicationChannelId, sequenceCounter: sequenceCounter, status: status}
}

func (m DeviceConfigurationAckDataBlock) LengthInBits() uint16 {
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

func (m DeviceConfigurationAckDataBlock) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func DeviceConfigurationAckDataBlockParse(io spi.ReadBuffer) (spi.Message, error) {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var _ uint8 = io.ReadUint8(8)

	// Simple Field (communicationChannelId)
	var communicationChannelId uint8 = io.ReadUint8(8)

	// Simple Field (sequenceCounter)
	var sequenceCounter uint8 = io.ReadUint8(8)

	// Enum field (status)
	status, _statusErr := StatusParse(io)
	if _statusErr != nil {
		return nil, errors.New("Error parsing 'status' field " + _statusErr.Error())
	}

	// Create the instance
	return NewDeviceConfigurationAckDataBlock(communicationChannelId, sequenceCounter, status), nil
}

func (m DeviceConfigurationAckDataBlock) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IDeviceConfigurationAckDataBlock); ok {

			// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
			structureLength := uint8(m.LengthInBytes())
			io.WriteUint8(8, (structureLength))

			// Simple Field (communicationChannelId)
			var communicationChannelId uint8 = m.communicationChannelId
			io.WriteUint8(8, (communicationChannelId))

			// Simple Field (sequenceCounter)
			var sequenceCounter uint8 = m.sequenceCounter
			io.WriteUint8(8, (sequenceCounter))

			// Enum field (status)
			status := m.status
			status.Serialize(io)
		}
	}
	serializeFunc(m)
}
