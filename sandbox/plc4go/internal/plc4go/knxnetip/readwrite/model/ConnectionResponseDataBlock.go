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
type ConnectionResponseDataBlock struct {
}

// The corresponding interface
type IConnectionResponseDataBlock interface {
	spi.Message
	ConnectionType() uint8
	Serialize(io spi.WriteBuffer)
}

type ConnectionResponseDataBlockInitializer interface {
	initialize() spi.Message
}

func ConnectionResponseDataBlockConnectionType(m IConnectionResponseDataBlock) uint8 {
	return m.ConnectionType()
}

func CastIConnectionResponseDataBlock(structType interface{}) IConnectionResponseDataBlock {
	castFunc := func(typ interface{}) IConnectionResponseDataBlock {
		if iConnectionResponseDataBlock, ok := typ.(IConnectionResponseDataBlock); ok {
			return iConnectionResponseDataBlock
		}
		return nil
	}
	return castFunc(structType)
}

func CastConnectionResponseDataBlock(structType interface{}) ConnectionResponseDataBlock {
	castFunc := func(typ interface{}) ConnectionResponseDataBlock {
		if sConnectionResponseDataBlock, ok := typ.(ConnectionResponseDataBlock); ok {
			return sConnectionResponseDataBlock
		}
		return ConnectionResponseDataBlock{}
	}
	return castFunc(structType)
}

func (m ConnectionResponseDataBlock) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Implicit Field (structureLength)
	lengthInBits += 8

	// Discriminator Field (connectionType)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m ConnectionResponseDataBlock) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ConnectionResponseDataBlockParse(io spi.ReadBuffer) (spi.Message, error) {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var _ uint8 = io.ReadUint8(8)

	// Discriminator Field (connectionType) (Used as input to a switch field)
	var connectionType uint8 = io.ReadUint8(8)

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer ConnectionResponseDataBlockInitializer
	var typeSwitchError error
	switch {
	case connectionType == 0x03:
		initializer, typeSwitchError = ConnectionResponseDataBlockDeviceManagementParse(io)
	case connectionType == 0x04:
		initializer, typeSwitchError = ConnectionResponseDataBlockTunnelConnectionParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func (m ConnectionResponseDataBlock) Serialize(io spi.WriteBuffer) {
	iConnectionResponseDataBlock := CastIConnectionResponseDataBlock(m)

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	structureLength := uint8(uint8(m.LengthInBytes()))
	io.WriteUint8(8, (structureLength))

	// Discriminator Field (connectionType) (Used as input to a switch field)
	connectionType := uint8(ConnectionResponseDataBlockConnectionType(iConnectionResponseDataBlock))
	io.WriteUint8(8, (connectionType))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	iConnectionResponseDataBlock.Serialize(io)
}
