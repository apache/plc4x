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
type KNXGroupAddress3Level struct {
	mainGroup   uint8
	middleGroup uint8
	subGroup    uint8
	KNXGroupAddress
}

// The corresponding interface
type IKNXGroupAddress3Level interface {
	IKNXGroupAddress
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m KNXGroupAddress3Level) NumLevels() uint8 {
	return 3
}

func (m KNXGroupAddress3Level) initialize() spi.Message {
	return m
}

func NewKNXGroupAddress3Level(mainGroup uint8, middleGroup uint8, subGroup uint8) KNXGroupAddressInitializer {
	return &KNXGroupAddress3Level{mainGroup: mainGroup, middleGroup: middleGroup, subGroup: subGroup}
}

func CastIKNXGroupAddress3Level(structType interface{}) IKNXGroupAddress3Level {
	castFunc := func(typ interface{}) IKNXGroupAddress3Level {
		if iKNXGroupAddress3Level, ok := typ.(IKNXGroupAddress3Level); ok {
			return iKNXGroupAddress3Level
		}
		return nil
	}
	return castFunc(structType)
}

func CastKNXGroupAddress3Level(structType interface{}) KNXGroupAddress3Level {
	castFunc := func(typ interface{}) KNXGroupAddress3Level {
		if sKNXGroupAddress3Level, ok := typ.(KNXGroupAddress3Level); ok {
			return sKNXGroupAddress3Level
		}
		return KNXGroupAddress3Level{}
	}
	return castFunc(structType)
}

func (m KNXGroupAddress3Level) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXGroupAddress.LengthInBits()

	// Simple field (mainGroup)
	lengthInBits += 5

	// Simple field (middleGroup)
	lengthInBits += 3

	// Simple field (subGroup)
	lengthInBits += 8

	return lengthInBits
}

func (m KNXGroupAddress3Level) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KNXGroupAddress3LevelParse(io *spi.ReadBuffer) (KNXGroupAddressInitializer, error) {

	// Simple Field (mainGroup)
	mainGroup, _mainGroupErr := io.ReadUint8(5)
	if _mainGroupErr != nil {
		return nil, errors.New("Error parsing 'mainGroup' field " + _mainGroupErr.Error())
	}

	// Simple Field (middleGroup)
	middleGroup, _middleGroupErr := io.ReadUint8(3)
	if _middleGroupErr != nil {
		return nil, errors.New("Error parsing 'middleGroup' field " + _middleGroupErr.Error())
	}

	// Simple Field (subGroup)
	subGroup, _subGroupErr := io.ReadUint8(8)
	if _subGroupErr != nil {
		return nil, errors.New("Error parsing 'subGroup' field " + _subGroupErr.Error())
	}

	// Create the instance
	return NewKNXGroupAddress3Level(mainGroup, middleGroup, subGroup), nil
}

func (m KNXGroupAddress3Level) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (mainGroup)
		mainGroup := uint8(m.mainGroup)
		_mainGroupErr := io.WriteUint8(5, (mainGroup))
		if _mainGroupErr != nil {
			return errors.New("Error serializing 'mainGroup' field " + _mainGroupErr.Error())
		}

		// Simple Field (middleGroup)
		middleGroup := uint8(m.middleGroup)
		_middleGroupErr := io.WriteUint8(3, (middleGroup))
		if _middleGroupErr != nil {
			return errors.New("Error serializing 'middleGroup' field " + _middleGroupErr.Error())
		}

		// Simple Field (subGroup)
		subGroup := uint8(m.subGroup)
		_subGroupErr := io.WriteUint8(8, (subGroup))
		if _subGroupErr != nil {
			return errors.New("Error serializing 'subGroup' field " + _subGroupErr.Error())
		}

		return nil
	}
	return KNXGroupAddressSerialize(io, m.KNXGroupAddress, CastIKNXGroupAddress(m), ser)
}
