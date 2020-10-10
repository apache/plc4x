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
type BACnetTagApplicationDouble struct {
	value float64
	BACnetTag
}

// The corresponding interface
type IBACnetTagApplicationDouble interface {
	IBACnetTag
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetTagApplicationDouble) ContextSpecificTag() uint8 {
	return 0
}

func (m BACnetTagApplicationDouble) initialize(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) spi.Message {
	m.typeOrTagNumber = typeOrTagNumber
	m.lengthValueType = lengthValueType
	m.extTagNumber = extTagNumber
	m.extLength = extLength
	return m
}

func NewBACnetTagApplicationDouble(value float64) BACnetTagInitializer {
	return &BACnetTagApplicationDouble{value: value}
}

func (m BACnetTagApplicationDouble) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetTag.LengthInBits()

	// Simple field (value)
	lengthInBits += 64

	return lengthInBits
}

func (m BACnetTagApplicationDouble) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetTagApplicationDoubleParse(io spi.ReadBuffer, lengthValueType uint8, extLength uint8) (BACnetTagInitializer, error) {

	// Simple Field (value)
	var value float64 = io.ReadFloat64(64)

	// Create the instance
	return NewBACnetTagApplicationDouble(value), nil
}

func (m BACnetTagApplicationDouble) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IBACnetTagApplicationDouble); ok {

			// Simple Field (value)
			var value float64 = m.value
			io.WriteFloat64(64, (value))
		}
	}
	serializeFunc(m)
}
