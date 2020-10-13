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
type BACnetConfirmedServiceACKReadPropertyMultiple struct {
	BACnetConfirmedServiceACK
}

// The corresponding interface
type IBACnetConfirmedServiceACKReadPropertyMultiple interface {
	IBACnetConfirmedServiceACK
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceACKReadPropertyMultiple) ServiceChoice() uint8 {
	return 0x0E
}

func (m BACnetConfirmedServiceACKReadPropertyMultiple) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceACKReadPropertyMultiple() BACnetConfirmedServiceACKInitializer {
	return &BACnetConfirmedServiceACKReadPropertyMultiple{}
}

func CastIBACnetConfirmedServiceACKReadPropertyMultiple(structType interface{}) IBACnetConfirmedServiceACKReadPropertyMultiple {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceACKReadPropertyMultiple {
		if iBACnetConfirmedServiceACKReadPropertyMultiple, ok := typ.(IBACnetConfirmedServiceACKReadPropertyMultiple); ok {
			return iBACnetConfirmedServiceACKReadPropertyMultiple
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceACKReadPropertyMultiple(structType interface{}) BACnetConfirmedServiceACKReadPropertyMultiple {
	castFunc := func(typ interface{}) BACnetConfirmedServiceACKReadPropertyMultiple {
		if sBACnetConfirmedServiceACKReadPropertyMultiple, ok := typ.(BACnetConfirmedServiceACKReadPropertyMultiple); ok {
			return sBACnetConfirmedServiceACKReadPropertyMultiple
		}
		return BACnetConfirmedServiceACKReadPropertyMultiple{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceACKReadPropertyMultiple) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceACK.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceACKReadPropertyMultiple) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceACKReadPropertyMultipleParse(io *spi.ReadBuffer) (BACnetConfirmedServiceACKInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceACKReadPropertyMultiple(), nil
}

func (m BACnetConfirmedServiceACKReadPropertyMultiple) Serialize(io spi.WriteBuffer) {
	ser := func() {

	}
	BACnetConfirmedServiceACKSerialize(io, m.BACnetConfirmedServiceACK, CastIBACnetConfirmedServiceACK(m), ser)
}
