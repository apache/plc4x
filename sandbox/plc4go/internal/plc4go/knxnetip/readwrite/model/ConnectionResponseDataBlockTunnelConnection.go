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
	"reflect"
)

// The data-structure of this message
type ConnectionResponseDataBlockTunnelConnection struct {
	knxAddress IKNXAddress
	ConnectionResponseDataBlock
}

// The corresponding interface
type IConnectionResponseDataBlockTunnelConnection interface {
	IConnectionResponseDataBlock
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ConnectionResponseDataBlockTunnelConnection) ConnectionType() uint8 {
	return 0x04
}

func (m ConnectionResponseDataBlockTunnelConnection) initialize() spi.Message {
	return m
}

func NewConnectionResponseDataBlockTunnelConnection(knxAddress IKNXAddress) ConnectionResponseDataBlockInitializer {
	return &ConnectionResponseDataBlockTunnelConnection{knxAddress: knxAddress}
}

func CastIConnectionResponseDataBlockTunnelConnection(structType interface{}) IConnectionResponseDataBlockTunnelConnection {
	castFunc := func(typ interface{}) IConnectionResponseDataBlockTunnelConnection {
		if iConnectionResponseDataBlockTunnelConnection, ok := typ.(IConnectionResponseDataBlockTunnelConnection); ok {
			return iConnectionResponseDataBlockTunnelConnection
		}
		return nil
	}
	return castFunc(structType)
}

func CastConnectionResponseDataBlockTunnelConnection(structType interface{}) ConnectionResponseDataBlockTunnelConnection {
	castFunc := func(typ interface{}) ConnectionResponseDataBlockTunnelConnection {
		if sConnectionResponseDataBlockTunnelConnection, ok := typ.(ConnectionResponseDataBlockTunnelConnection); ok {
			return sConnectionResponseDataBlockTunnelConnection
		}
		return ConnectionResponseDataBlockTunnelConnection{}
	}
	return castFunc(structType)
}

func (m ConnectionResponseDataBlockTunnelConnection) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ConnectionResponseDataBlock.LengthInBits()

	// Simple field (knxAddress)
	lengthInBits += m.knxAddress.LengthInBits()

	return lengthInBits
}

func (m ConnectionResponseDataBlockTunnelConnection) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ConnectionResponseDataBlockTunnelConnectionParse(io spi.ReadBuffer) (ConnectionResponseDataBlockInitializer, error) {

	// Simple Field (knxAddress)
	_knxAddressMessage, _err := KNXAddressParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'knxAddress'. " + _err.Error())
	}
	var knxAddress IKNXAddress
	knxAddress, _knxAddressOk := _knxAddressMessage.(IKNXAddress)
	if !_knxAddressOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_knxAddressMessage).Name() + " to IKNXAddress")
	}

	// Create the instance
	return NewConnectionResponseDataBlockTunnelConnection(knxAddress), nil
}

func (m ConnectionResponseDataBlockTunnelConnection) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Simple Field (knxAddress)
		knxAddress := CastIKNXAddress(m.knxAddress)
		knxAddress.Serialize(io)

	}
	ConnectionResponseDataBlockSerialize(io, m.ConnectionResponseDataBlock, CastIConnectionResponseDataBlock(m), ser)
}
