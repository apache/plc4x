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
type BACnetTagApplicationCharacterString struct {
	BACnetTag
}

// The corresponding interface
type IBACnetTagApplicationCharacterString interface {
	IBACnetTag
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetTagApplicationCharacterString) ContextSpecificTag() uint8 {
	return 0
}

func (m BACnetTagApplicationCharacterString) initialize(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) spi.Message {
	m.TypeOrTagNumber = typeOrTagNumber
	m.LengthValueType = lengthValueType
	m.ExtTagNumber = extTagNumber
	m.ExtLength = extLength
	return m
}

func NewBACnetTagApplicationCharacterString() BACnetTagInitializer {
	return &BACnetTagApplicationCharacterString{}
}

func CastIBACnetTagApplicationCharacterString(structType interface{}) IBACnetTagApplicationCharacterString {
	castFunc := func(typ interface{}) IBACnetTagApplicationCharacterString {
		if iBACnetTagApplicationCharacterString, ok := typ.(IBACnetTagApplicationCharacterString); ok {
			return iBACnetTagApplicationCharacterString
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetTagApplicationCharacterString(structType interface{}) BACnetTagApplicationCharacterString {
	castFunc := func(typ interface{}) BACnetTagApplicationCharacterString {
		if sBACnetTagApplicationCharacterString, ok := typ.(BACnetTagApplicationCharacterString); ok {
			return sBACnetTagApplicationCharacterString
		}
		return BACnetTagApplicationCharacterString{}
	}
	return castFunc(structType)
}

func (m BACnetTagApplicationCharacterString) LengthInBits() uint16 {
	var lengthInBits = m.BACnetTag.LengthInBits()

	return lengthInBits
}

func (m BACnetTagApplicationCharacterString) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetTagApplicationCharacterStringParse(io *spi.ReadBuffer) (BACnetTagInitializer, error) {

	// Create the instance
	return NewBACnetTagApplicationCharacterString(), nil
}

func (m BACnetTagApplicationCharacterString) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetTagSerialize(io, m.BACnetTag, CastIBACnetTag(m), ser)
}
