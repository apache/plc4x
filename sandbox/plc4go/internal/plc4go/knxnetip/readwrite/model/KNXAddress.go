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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type KNXAddress struct {
	mainGroup   uint8
	middleGroup uint8
	subGroup    uint8
}

// The corresponding interface
type IKNXAddress interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewKNXAddress(mainGroup uint8, middleGroup uint8, subGroup uint8) spi.Message {
	return &KNXAddress{mainGroup: mainGroup, middleGroup: middleGroup, subGroup: subGroup}
}

func CastIKNXAddress(structType interface{}) IKNXAddress {
	castFunc := func(typ interface{}) IKNXAddress {
		if iKNXAddress, ok := typ.(IKNXAddress); ok {
			return iKNXAddress
		}
		return nil
	}
	return castFunc(structType)
}

func CastKNXAddress(structType interface{}) KNXAddress {
	castFunc := func(typ interface{}) KNXAddress {
		if sKNXAddress, ok := typ.(KNXAddress); ok {
			return sKNXAddress
		}
		return KNXAddress{}
	}
	return castFunc(structType)
}

func (m KNXAddress) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Simple field (mainGroup)
	lengthInBits += 4

	// Simple field (middleGroup)
	lengthInBits += 4

	// Simple field (subGroup)
	lengthInBits += 8

	return lengthInBits
}

func (m KNXAddress) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KNXAddressParse(io spi.ReadBuffer) (spi.Message, error) {

	// Simple Field (mainGroup)
	var mainGroup uint8 = io.ReadUint8(4)

	// Simple Field (middleGroup)
	var middleGroup uint8 = io.ReadUint8(4)

	// Simple Field (subGroup)
	var subGroup uint8 = io.ReadUint8(8)

	// Create the instance
	return NewKNXAddress(mainGroup, middleGroup, subGroup), nil
}

func (m KNXAddress) Serialize(io spi.WriteBuffer) {

	// Simple Field (mainGroup)
	mainGroup := uint8(m.mainGroup)
	io.WriteUint8(4, (mainGroup))

	// Simple Field (middleGroup)
	middleGroup := uint8(m.middleGroup)
	io.WriteUint8(4, (middleGroup))

	// Simple Field (subGroup)
	subGroup := uint8(m.subGroup)
	io.WriteUint8(8, (subGroup))
}
