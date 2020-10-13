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
type BACnetConfirmedServiceACKRemovedAuthenticate struct {
	BACnetConfirmedServiceACK
}

// The corresponding interface
type IBACnetConfirmedServiceACKRemovedAuthenticate interface {
	IBACnetConfirmedServiceACK
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceACKRemovedAuthenticate) ServiceChoice() uint8 {
	return 0x18
}

func (m BACnetConfirmedServiceACKRemovedAuthenticate) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceACKRemovedAuthenticate() BACnetConfirmedServiceACKInitializer {
	return &BACnetConfirmedServiceACKRemovedAuthenticate{}
}

func CastIBACnetConfirmedServiceACKRemovedAuthenticate(structType interface{}) IBACnetConfirmedServiceACKRemovedAuthenticate {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceACKRemovedAuthenticate {
		if iBACnetConfirmedServiceACKRemovedAuthenticate, ok := typ.(IBACnetConfirmedServiceACKRemovedAuthenticate); ok {
			return iBACnetConfirmedServiceACKRemovedAuthenticate
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceACKRemovedAuthenticate(structType interface{}) BACnetConfirmedServiceACKRemovedAuthenticate {
	castFunc := func(typ interface{}) BACnetConfirmedServiceACKRemovedAuthenticate {
		if sBACnetConfirmedServiceACKRemovedAuthenticate, ok := typ.(BACnetConfirmedServiceACKRemovedAuthenticate); ok {
			return sBACnetConfirmedServiceACKRemovedAuthenticate
		}
		return BACnetConfirmedServiceACKRemovedAuthenticate{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceACKRemovedAuthenticate) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceACK.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceACKRemovedAuthenticate) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceACKRemovedAuthenticateParse(io *spi.ReadBuffer) (BACnetConfirmedServiceACKInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceACKRemovedAuthenticate(), nil
}

func (m BACnetConfirmedServiceACKRemovedAuthenticate) Serialize(io spi.WriteBuffer) {
	ser := func() {

	}
	BACnetConfirmedServiceACKSerialize(io, m.BACnetConfirmedServiceACK, CastIBACnetConfirmedServiceACK(m), ser)
}
