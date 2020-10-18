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
type KnxNetRemoteConfigurationAndDiagnosis struct {
	Version uint8
	ServiceId
}

// The corresponding interface
type IKnxNetRemoteConfigurationAndDiagnosis interface {
	IServiceId
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m KnxNetRemoteConfigurationAndDiagnosis) ServiceType() uint8 {
	return 0x07
}

func (m KnxNetRemoteConfigurationAndDiagnosis) initialize() spi.Message {
	return m
}

func NewKnxNetRemoteConfigurationAndDiagnosis(version uint8) ServiceIdInitializer {
	return &KnxNetRemoteConfigurationAndDiagnosis{Version: version}
}

func CastIKnxNetRemoteConfigurationAndDiagnosis(structType interface{}) IKnxNetRemoteConfigurationAndDiagnosis {
	castFunc := func(typ interface{}) IKnxNetRemoteConfigurationAndDiagnosis {
		if iKnxNetRemoteConfigurationAndDiagnosis, ok := typ.(IKnxNetRemoteConfigurationAndDiagnosis); ok {
			return iKnxNetRemoteConfigurationAndDiagnosis
		}
		return nil
	}
	return castFunc(structType)
}

func CastKnxNetRemoteConfigurationAndDiagnosis(structType interface{}) KnxNetRemoteConfigurationAndDiagnosis {
	castFunc := func(typ interface{}) KnxNetRemoteConfigurationAndDiagnosis {
		if sKnxNetRemoteConfigurationAndDiagnosis, ok := typ.(KnxNetRemoteConfigurationAndDiagnosis); ok {
			return sKnxNetRemoteConfigurationAndDiagnosis
		}
		return KnxNetRemoteConfigurationAndDiagnosis{}
	}
	return castFunc(structType)
}

func (m KnxNetRemoteConfigurationAndDiagnosis) LengthInBits() uint16 {
	var lengthInBits = m.ServiceId.LengthInBits()

	// Simple field (version)
	lengthInBits += 8

	return lengthInBits
}

func (m KnxNetRemoteConfigurationAndDiagnosis) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KnxNetRemoteConfigurationAndDiagnosisParse(io *spi.ReadBuffer) (ServiceIdInitializer, error) {

	// Simple Field (version)
	version, _versionErr := io.ReadUint8(8)
	if _versionErr != nil {
		return nil, errors.New("Error parsing 'version' field " + _versionErr.Error())
	}

	// Create the instance
	return NewKnxNetRemoteConfigurationAndDiagnosis(version), nil
}

func (m KnxNetRemoteConfigurationAndDiagnosis) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (version)
		version := uint8(m.Version)
		_versionErr := io.WriteUint8(8, version)
		if _versionErr != nil {
			return errors.New("Error serializing 'version' field " + _versionErr.Error())
		}

		return nil
	}
	return ServiceIdSerialize(io, m.ServiceId, CastIServiceId(m), ser)
}
