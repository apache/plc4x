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
type KNXGroupAddress2Level struct {
	mainGroup uint8
	subGroup  uint16
	KNXGroupAddress
}

// The corresponding interface
type IKNXGroupAddress2Level interface {
	IKNXGroupAddress
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m KNXGroupAddress2Level) NumLevels() uint8 {
	return 2
}

func (m KNXGroupAddress2Level) initialize() spi.Message {
	return m
}

func NewKNXGroupAddress2Level(mainGroup uint8, subGroup uint16) KNXGroupAddressInitializer {
	return &KNXGroupAddress2Level{mainGroup: mainGroup, subGroup: subGroup}
}

func CastIKNXGroupAddress2Level(structType interface{}) IKNXGroupAddress2Level {
	castFunc := func(typ interface{}) IKNXGroupAddress2Level {
		if iKNXGroupAddress2Level, ok := typ.(IKNXGroupAddress2Level); ok {
			return iKNXGroupAddress2Level
		}
		return nil
	}
	return castFunc(structType)
}

func CastKNXGroupAddress2Level(structType interface{}) KNXGroupAddress2Level {
	castFunc := func(typ interface{}) KNXGroupAddress2Level {
		if sKNXGroupAddress2Level, ok := typ.(KNXGroupAddress2Level); ok {
			return sKNXGroupAddress2Level
		}
		return KNXGroupAddress2Level{}
	}
	return castFunc(structType)
}

func (m KNXGroupAddress2Level) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXGroupAddress.LengthInBits()

	// Simple field (mainGroup)
	lengthInBits += 5

	// Simple field (subGroup)
	lengthInBits += 11

	return lengthInBits
}

func (m KNXGroupAddress2Level) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KNXGroupAddress2LevelParse(io spi.ReadBuffer) (KNXGroupAddressInitializer, error) {

	// Simple Field (mainGroup)
	mainGroup, _mainGroupErr := io.ReadUint8(5)
	if _mainGroupErr != nil {
		return nil, errors.New("Error parsing 'mainGroup' field " + _mainGroupErr.Error())
	}

	// Simple Field (subGroup)
	subGroup, _subGroupErr := io.ReadUint16(11)
	if _subGroupErr != nil {
		return nil, errors.New("Error parsing 'subGroup' field " + _subGroupErr.Error())
	}

	// Create the instance
	return NewKNXGroupAddress2Level(mainGroup, subGroup), nil
}

func (m KNXGroupAddress2Level) Serialize(io spi.WriteBuffer) {

	// Simple Field (mainGroup)
	mainGroup := uint8(m.mainGroup)
	io.WriteUint8(5, (mainGroup))

	// Simple Field (subGroup)
	subGroup := uint16(m.subGroup)
	io.WriteUint16(11, (subGroup))
}
