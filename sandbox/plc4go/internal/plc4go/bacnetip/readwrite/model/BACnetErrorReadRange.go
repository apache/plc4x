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
type BACnetErrorReadRange struct {
	BACnetError
}

// The corresponding interface
type IBACnetErrorReadRange interface {
	IBACnetError
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetErrorReadRange) ServiceChoice() uint8 {
	return 0x1A
}

func (m BACnetErrorReadRange) initialize() spi.Message {
	return m
}

func NewBACnetErrorReadRange() BACnetErrorInitializer {
	return &BACnetErrorReadRange{}
}

func CastIBACnetErrorReadRange(structType interface{}) IBACnetErrorReadRange {
	castFunc := func(typ interface{}) IBACnetErrorReadRange {
		if iBACnetErrorReadRange, ok := typ.(IBACnetErrorReadRange); ok {
			return iBACnetErrorReadRange
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetErrorReadRange(structType interface{}) BACnetErrorReadRange {
	castFunc := func(typ interface{}) BACnetErrorReadRange {
		if sBACnetErrorReadRange, ok := typ.(BACnetErrorReadRange); ok {
			return sBACnetErrorReadRange
		}
		return BACnetErrorReadRange{}
	}
	return castFunc(structType)
}

func (m BACnetErrorReadRange) LengthInBits() uint16 {
	var lengthInBits = m.BACnetError.LengthInBits()

	return lengthInBits
}

func (m BACnetErrorReadRange) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetErrorReadRangeParse(io *spi.ReadBuffer) (BACnetErrorInitializer, error) {

	// Create the instance
	return NewBACnetErrorReadRange(), nil
}

func (m BACnetErrorReadRange) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetErrorSerialize(io, m.BACnetError, CastIBACnetError(m), ser)
}
