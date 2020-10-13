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
type KNXGroupAddressFreeLevel struct {
	subGroup uint16
	KNXGroupAddress
}

// The corresponding interface
type IKNXGroupAddressFreeLevel interface {
	IKNXGroupAddress
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m KNXGroupAddressFreeLevel) NumLevels() uint8 {
	return 1
}

func (m KNXGroupAddressFreeLevel) initialize() spi.Message {
	return m
}

func NewKNXGroupAddressFreeLevel(subGroup uint16) KNXGroupAddressInitializer {
	return &KNXGroupAddressFreeLevel{subGroup: subGroup}
}

func CastIKNXGroupAddressFreeLevel(structType interface{}) IKNXGroupAddressFreeLevel {
	castFunc := func(typ interface{}) IKNXGroupAddressFreeLevel {
		if iKNXGroupAddressFreeLevel, ok := typ.(IKNXGroupAddressFreeLevel); ok {
			return iKNXGroupAddressFreeLevel
		}
		return nil
	}
	return castFunc(structType)
}

func CastKNXGroupAddressFreeLevel(structType interface{}) KNXGroupAddressFreeLevel {
	castFunc := func(typ interface{}) KNXGroupAddressFreeLevel {
		if sKNXGroupAddressFreeLevel, ok := typ.(KNXGroupAddressFreeLevel); ok {
			return sKNXGroupAddressFreeLevel
		}
		return KNXGroupAddressFreeLevel{}
	}
	return castFunc(structType)
}

func (m KNXGroupAddressFreeLevel) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXGroupAddress.LengthInBits()

	// Simple field (subGroup)
	lengthInBits += 16

	return lengthInBits
}

func (m KNXGroupAddressFreeLevel) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KNXGroupAddressFreeLevelParse(io spi.ReadBuffer) (KNXGroupAddressInitializer, error) {

	// Simple Field (subGroup)
	subGroup, _subGroupErr := io.ReadUint16(16)
	if _subGroupErr != nil {
		return nil, errors.New("Error parsing 'subGroup' field " + _subGroupErr.Error())
	}

	// Create the instance
	return NewKNXGroupAddressFreeLevel(subGroup), nil
}

func (m KNXGroupAddressFreeLevel) Serialize(io spi.WriteBuffer) {

	// Simple Field (subGroup)
	subGroup := uint16(m.subGroup)
	io.WriteUint16(16, (subGroup))
}
