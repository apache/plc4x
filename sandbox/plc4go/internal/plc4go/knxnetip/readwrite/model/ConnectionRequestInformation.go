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
type ConnectionRequestInformation struct {
}

// The corresponding interface
type IConnectionRequestInformation interface {
	spi.Message
	ConnectionType() uint8
	Serialize(io spi.WriteBuffer)
}

type ConnectionRequestInformationInitializer interface {
	initialize() spi.Message
}

func ConnectionRequestInformationConnectionType(m IConnectionRequestInformation) uint8 {
	return m.ConnectionType()
}

func CastIConnectionRequestInformation(structType interface{}) IConnectionRequestInformation {
	castFunc := func(typ interface{}) IConnectionRequestInformation {
		if iConnectionRequestInformation, ok := typ.(IConnectionRequestInformation); ok {
			return iConnectionRequestInformation
		}
		return nil
	}
	return castFunc(structType)
}

func CastConnectionRequestInformation(structType interface{}) ConnectionRequestInformation {
	castFunc := func(typ interface{}) ConnectionRequestInformation {
		if sConnectionRequestInformation, ok := typ.(ConnectionRequestInformation); ok {
			return sConnectionRequestInformation
		}
		return ConnectionRequestInformation{}
	}
	return castFunc(structType)
}

func (m ConnectionRequestInformation) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Implicit Field (structureLength)
	lengthInBits += 8

	// Discriminator Field (connectionType)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m ConnectionRequestInformation) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ConnectionRequestInformationParse(io spi.ReadBuffer) (spi.Message, error) {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	_, _structureLengthErr := io.ReadUint8(8)
	if _structureLengthErr != nil {
		return nil, errors.New("Error parsing 'structureLength' field " + _structureLengthErr.Error())
	}

	// Discriminator Field (connectionType) (Used as input to a switch field)
	connectionType, _connectionTypeErr := io.ReadUint8(8)
	if _connectionTypeErr != nil {
		return nil, errors.New("Error parsing 'connectionType' field " + _connectionTypeErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer ConnectionRequestInformationInitializer
	var typeSwitchError error
	switch {
	case connectionType == 0x03:
		initializer, typeSwitchError = ConnectionRequestInformationDeviceManagementParse(io)
	case connectionType == 0x04:
		initializer, typeSwitchError = ConnectionRequestInformationTunnelConnectionParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func ConnectionRequestInformationSerialize(io spi.WriteBuffer, m ConnectionRequestInformation, i IConnectionRequestInformation, childSerialize func()) {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	structureLength := uint8(uint8(m.LengthInBytes()))
	io.WriteUint8(8, (structureLength))

	// Discriminator Field (connectionType) (Used as input to a switch field)
	connectionType := uint8(i.ConnectionType())
	io.WriteUint8(8, (connectionType))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	childSerialize()

}
