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
type ServiceId struct {
}

// The corresponding interface
type IServiceId interface {
	spi.Message
	ServiceType() uint8
	Serialize(io spi.WriteBuffer)
}

type ServiceIdInitializer interface {
	initialize() spi.Message
}

func ServiceIdServiceType(m IServiceId) uint8 {
	return m.ServiceType()
}

func CastIServiceId(structType interface{}) IServiceId {
	castFunc := func(typ interface{}) IServiceId {
		if iServiceId, ok := typ.(IServiceId); ok {
			return iServiceId
		}
		return nil
	}
	return castFunc(structType)
}

func CastServiceId(structType interface{}) ServiceId {
	castFunc := func(typ interface{}) ServiceId {
		if sServiceId, ok := typ.(ServiceId); ok {
			return sServiceId
		}
		return ServiceId{}
	}
	return castFunc(structType)
}

func (m ServiceId) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (serviceType)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m ServiceId) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ServiceIdParse(io spi.ReadBuffer) (spi.Message, error) {

	// Discriminator Field (serviceType) (Used as input to a switch field)
	serviceType, _serviceTypeErr := io.ReadUint8(8)
	if _serviceTypeErr != nil {
		return nil, errors.New("Error parsing 'serviceType' field " + _serviceTypeErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer ServiceIdInitializer
	var typeSwitchError error
	switch {
	case serviceType == 0x02:
		initializer, typeSwitchError = KnxNetIpCoreParse(io)
	case serviceType == 0x03:
		initializer, typeSwitchError = KnxNetIpDeviceManagementParse(io)
	case serviceType == 0x04:
		initializer, typeSwitchError = KnxNetIpTunnelingParse(io)
	case serviceType == 0x06:
		initializer, typeSwitchError = KnxNetRemoteLoggingParse(io)
	case serviceType == 0x07:
		initializer, typeSwitchError = KnxNetRemoteConfigurationAndDiagnosisParse(io)
	case serviceType == 0x08:
		initializer, typeSwitchError = KnxNetObjectServerParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func (m ServiceId) Serialize(io spi.WriteBuffer) {
	iServiceId := CastIServiceId(m)

	// Discriminator Field (serviceType) (Used as input to a switch field)
	serviceType := uint8(ServiceIdServiceType(iServiceId))
	io.WriteUint8(8, (serviceType))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	iServiceId.Serialize(io)
}
