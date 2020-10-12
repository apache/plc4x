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
type BACnetTagApplicationNull struct {
	BACnetTag
}

// The corresponding interface
type IBACnetTagApplicationNull interface {
	IBACnetTag
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetTagApplicationNull) ContextSpecificTag() uint8 {
	return 0
}

func (m BACnetTagApplicationNull) initialize(typeOrTagNumber uint8, lengthValueType uint8, extTagNumber *uint8, extLength *uint8) spi.Message {
	m.typeOrTagNumber = typeOrTagNumber
	m.lengthValueType = lengthValueType
	m.extTagNumber = extTagNumber
	m.extLength = extLength
	return m
}

func NewBACnetTagApplicationNull() BACnetTagInitializer {
	return &BACnetTagApplicationNull{}
}

func CastIBACnetTagApplicationNull(structType interface{}) IBACnetTagApplicationNull {
	castFunc := func(typ interface{}) IBACnetTagApplicationNull {
		if iBACnetTagApplicationNull, ok := typ.(IBACnetTagApplicationNull); ok {
			return iBACnetTagApplicationNull
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetTagApplicationNull(structType interface{}) BACnetTagApplicationNull {
	castFunc := func(typ interface{}) BACnetTagApplicationNull {
		if sBACnetTagApplicationNull, ok := typ.(BACnetTagApplicationNull); ok {
			return sBACnetTagApplicationNull
		}
		return BACnetTagApplicationNull{}
	}
	return castFunc(structType)
}

func (m BACnetTagApplicationNull) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetTag.LengthInBits()

	return lengthInBits
}

func (m BACnetTagApplicationNull) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetTagApplicationNullParse(io spi.ReadBuffer) (BACnetTagInitializer, error) {

	// Create the instance
	return NewBACnetTagApplicationNull(), nil
}

func (m BACnetTagApplicationNull) Serialize(io spi.WriteBuffer) {

}
