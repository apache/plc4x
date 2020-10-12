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
type BACnetConfirmedServiceACKGetEventInformation struct {
	BACnetConfirmedServiceACK
}

// The corresponding interface
type IBACnetConfirmedServiceACKGetEventInformation interface {
	IBACnetConfirmedServiceACK
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceACKGetEventInformation) ServiceChoice() uint8 {
	return 0x1D
}

func (m BACnetConfirmedServiceACKGetEventInformation) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceACKGetEventInformation() BACnetConfirmedServiceACKInitializer {
	return &BACnetConfirmedServiceACKGetEventInformation{}
}

func CastIBACnetConfirmedServiceACKGetEventInformation(structType interface{}) IBACnetConfirmedServiceACKGetEventInformation {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceACKGetEventInformation {
		if iBACnetConfirmedServiceACKGetEventInformation, ok := typ.(IBACnetConfirmedServiceACKGetEventInformation); ok {
			return iBACnetConfirmedServiceACKGetEventInformation
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceACKGetEventInformation(structType interface{}) BACnetConfirmedServiceACKGetEventInformation {
	castFunc := func(typ interface{}) BACnetConfirmedServiceACKGetEventInformation {
		if sBACnetConfirmedServiceACKGetEventInformation, ok := typ.(BACnetConfirmedServiceACKGetEventInformation); ok {
			return sBACnetConfirmedServiceACKGetEventInformation
		}
		return BACnetConfirmedServiceACKGetEventInformation{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceACKGetEventInformation) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceACK.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceACKGetEventInformation) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceACKGetEventInformationParse(io spi.ReadBuffer) (BACnetConfirmedServiceACKInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceACKGetEventInformation(), nil
}

func (m BACnetConfirmedServiceACKGetEventInformation) Serialize(io spi.WriteBuffer) {

}
