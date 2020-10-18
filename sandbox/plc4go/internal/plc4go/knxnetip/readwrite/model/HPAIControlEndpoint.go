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
type HPAIControlEndpoint struct {
	HostProtocolCode IHostProtocolCode
	IpAddress        IIPAddress
	IpPort           uint16
}

// The corresponding interface
type IHPAIControlEndpoint interface {
	spi.Message
	Serialize(io spi.WriteBuffer) error
}

func NewHPAIControlEndpoint(hostProtocolCode IHostProtocolCode, ipAddress IIPAddress, ipPort uint16) spi.Message {
	return &HPAIControlEndpoint{HostProtocolCode: hostProtocolCode, IpAddress: ipAddress, IpPort: ipPort}
}

func CastIHPAIControlEndpoint(structType interface{}) IHPAIControlEndpoint {
	castFunc := func(typ interface{}) IHPAIControlEndpoint {
		if iHPAIControlEndpoint, ok := typ.(IHPAIControlEndpoint); ok {
			return iHPAIControlEndpoint
		}
		return nil
	}
	return castFunc(structType)
}

func CastHPAIControlEndpoint(structType interface{}) HPAIControlEndpoint {
	castFunc := func(typ interface{}) HPAIControlEndpoint {
		if sHPAIControlEndpoint, ok := typ.(HPAIControlEndpoint); ok {
			return sHPAIControlEndpoint
		}
		return HPAIControlEndpoint{}
	}
	return castFunc(structType)
}

func (m HPAIControlEndpoint) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Implicit Field (structureLength)
	lengthInBits += 8

	// Enum Field (hostProtocolCode)
	lengthInBits += 8

	// Simple field (ipAddress)
	lengthInBits += m.IpAddress.LengthInBits()

	// Simple field (ipPort)
	lengthInBits += 16

	return lengthInBits
}

func (m HPAIControlEndpoint) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func HPAIControlEndpointParse(io *spi.ReadBuffer) (spi.Message, error) {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	_, _structureLengthErr := io.ReadUint8(8)
	if _structureLengthErr != nil {
		return nil, errors.New("Error parsing 'structureLength' field " + _structureLengthErr.Error())
	}

	// Enum field (hostProtocolCode)
	hostProtocolCode, _hostProtocolCodeErr := HostProtocolCodeParse(io)
	if _hostProtocolCodeErr != nil {
		return nil, errors.New("Error parsing 'hostProtocolCode' field " + _hostProtocolCodeErr.Error())
	}

	// Simple Field (ipAddress)
	_ipAddressMessage, _err := IPAddressParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'ipAddress'. " + _err.Error())
	}
	var ipAddress IIPAddress
	ipAddress, _ipAddressOk := _ipAddressMessage.(IIPAddress)
	if !_ipAddressOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_ipAddressMessage).Name() + " to IIPAddress")
	}

	// Simple Field (ipPort)
	ipPort, _ipPortErr := io.ReadUint16(16)
	if _ipPortErr != nil {
		return nil, errors.New("Error parsing 'ipPort' field " + _ipPortErr.Error())
	}

	// Create the instance
	return NewHPAIControlEndpoint(hostProtocolCode, ipAddress, ipPort), nil
}

func (m HPAIControlEndpoint) Serialize(io spi.WriteBuffer) error {

	// Implicit Field (structureLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	structureLength := uint8(uint8(m.LengthInBytes()))
	_structureLengthErr := io.WriteUint8(8, structureLength)
	if _structureLengthErr != nil {
		return errors.New("Error serializing 'structureLength' field " + _structureLengthErr.Error())
	}

	// Enum field (hostProtocolCode)
	hostProtocolCode := CastHostProtocolCode(m.HostProtocolCode)
	_hostProtocolCodeErr := hostProtocolCode.Serialize(io)
	if _hostProtocolCodeErr != nil {
		return errors.New("Error serializing 'hostProtocolCode' field " + _hostProtocolCodeErr.Error())
	}

	// Simple Field (ipAddress)
	ipAddress := CastIIPAddress(m.IpAddress)
	_ipAddressErr := ipAddress.Serialize(io)
	if _ipAddressErr != nil {
		return errors.New("Error serializing 'ipAddress' field " + _ipAddressErr.Error())
	}

	// Simple Field (ipPort)
	ipPort := uint16(m.IpPort)
	_ipPortErr := io.WriteUint16(16, ipPort)
	if _ipPortErr != nil {
		return errors.New("Error serializing 'ipPort' field " + _ipPortErr.Error())
	}

	return nil
}
