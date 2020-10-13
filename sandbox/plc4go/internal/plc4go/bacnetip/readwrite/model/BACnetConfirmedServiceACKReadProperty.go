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
type BACnetConfirmedServiceACKReadProperty struct {
	BACnetConfirmedServiceACK
}

// The corresponding interface
type IBACnetConfirmedServiceACKReadProperty interface {
	IBACnetConfirmedServiceACK
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceACKReadProperty) ServiceChoice() uint8 {
	return 0x0C
}

func (m BACnetConfirmedServiceACKReadProperty) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceACKReadProperty() BACnetConfirmedServiceACKInitializer {
	return &BACnetConfirmedServiceACKReadProperty{}
}

func CastIBACnetConfirmedServiceACKReadProperty(structType interface{}) IBACnetConfirmedServiceACKReadProperty {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceACKReadProperty {
		if iBACnetConfirmedServiceACKReadProperty, ok := typ.(IBACnetConfirmedServiceACKReadProperty); ok {
			return iBACnetConfirmedServiceACKReadProperty
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceACKReadProperty(structType interface{}) BACnetConfirmedServiceACKReadProperty {
	castFunc := func(typ interface{}) BACnetConfirmedServiceACKReadProperty {
		if sBACnetConfirmedServiceACKReadProperty, ok := typ.(BACnetConfirmedServiceACKReadProperty); ok {
			return sBACnetConfirmedServiceACKReadProperty
		}
		return BACnetConfirmedServiceACKReadProperty{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceACKReadProperty) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceACK.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceACKReadProperty) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceACKReadPropertyParse(io *spi.ReadBuffer) (BACnetConfirmedServiceACKInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceACKReadProperty(), nil
}

func (m BACnetConfirmedServiceACKReadProperty) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetConfirmedServiceACKSerialize(io, m.BACnetConfirmedServiceACK, CastIBACnetConfirmedServiceACK(m), ser)
}
