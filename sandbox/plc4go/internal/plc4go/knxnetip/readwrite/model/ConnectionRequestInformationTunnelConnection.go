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
type ConnectionRequestInformationTunnelConnection struct {
	knxLayer IKnxLayer
	ConnectionRequestInformation
}

// The corresponding interface
type IConnectionRequestInformationTunnelConnection interface {
	IConnectionRequestInformation
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ConnectionRequestInformationTunnelConnection) ConnectionType() uint8 {
	return 0x04
}

func (m ConnectionRequestInformationTunnelConnection) initialize() spi.Message {
	return m
}

func NewConnectionRequestInformationTunnelConnection(knxLayer IKnxLayer) ConnectionRequestInformationInitializer {
	return &ConnectionRequestInformationTunnelConnection{knxLayer: knxLayer}
}

func CastIConnectionRequestInformationTunnelConnection(structType interface{}) IConnectionRequestInformationTunnelConnection {
	castFunc := func(typ interface{}) IConnectionRequestInformationTunnelConnection {
		if iConnectionRequestInformationTunnelConnection, ok := typ.(IConnectionRequestInformationTunnelConnection); ok {
			return iConnectionRequestInformationTunnelConnection
		}
		return nil
	}
	return castFunc(structType)
}

func CastConnectionRequestInformationTunnelConnection(structType interface{}) ConnectionRequestInformationTunnelConnection {
	castFunc := func(typ interface{}) ConnectionRequestInformationTunnelConnection {
		if sConnectionRequestInformationTunnelConnection, ok := typ.(ConnectionRequestInformationTunnelConnection); ok {
			return sConnectionRequestInformationTunnelConnection
		}
		return ConnectionRequestInformationTunnelConnection{}
	}
	return castFunc(structType)
}

func (m ConnectionRequestInformationTunnelConnection) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ConnectionRequestInformation.LengthInBits()

	// Enum Field (knxLayer)
	lengthInBits += 8

	// Reserved Field (reserved)
	lengthInBits += 8

	return lengthInBits
}

func (m ConnectionRequestInformationTunnelConnection) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ConnectionRequestInformationTunnelConnectionParse(io spi.ReadBuffer) (ConnectionRequestInformationInitializer, error) {

	// Enum field (knxLayer)
	knxLayer, _knxLayerErr := KnxLayerParse(io)
	if _knxLayerErr != nil {
		return nil, errors.New("Error parsing 'knxLayer' field " + _knxLayerErr.Error())
	}

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint8(8)
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

	// Create the instance
	return NewConnectionRequestInformationTunnelConnection(knxLayer), nil
}

func (m ConnectionRequestInformationTunnelConnection) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Enum field (knxLayer)
		knxLayer := CastKnxLayer(m.knxLayer)
		knxLayer.Serialize(io)

		// Reserved Field (reserved)
		io.WriteUint8(8, uint8(0x00))

	}
	ConnectionRequestInformationSerialize(io, m.ConnectionRequestInformation, CastIConnectionRequestInformation(m), ser)
}
