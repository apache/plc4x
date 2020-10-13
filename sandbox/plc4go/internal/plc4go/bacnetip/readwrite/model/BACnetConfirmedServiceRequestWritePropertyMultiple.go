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
type BACnetConfirmedServiceRequestWritePropertyMultiple struct {
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestWritePropertyMultiple interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestWritePropertyMultiple) ServiceChoice() uint8 {
	return 0x10
}

func (m BACnetConfirmedServiceRequestWritePropertyMultiple) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestWritePropertyMultiple() BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestWritePropertyMultiple{}
}

func CastIBACnetConfirmedServiceRequestWritePropertyMultiple(structType interface{}) IBACnetConfirmedServiceRequestWritePropertyMultiple {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestWritePropertyMultiple {
		if iBACnetConfirmedServiceRequestWritePropertyMultiple, ok := typ.(IBACnetConfirmedServiceRequestWritePropertyMultiple); ok {
			return iBACnetConfirmedServiceRequestWritePropertyMultiple
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestWritePropertyMultiple(structType interface{}) BACnetConfirmedServiceRequestWritePropertyMultiple {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestWritePropertyMultiple {
		if sBACnetConfirmedServiceRequestWritePropertyMultiple, ok := typ.(BACnetConfirmedServiceRequestWritePropertyMultiple); ok {
			return sBACnetConfirmedServiceRequestWritePropertyMultiple
		}
		return BACnetConfirmedServiceRequestWritePropertyMultiple{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestWritePropertyMultiple) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestWritePropertyMultiple) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestWritePropertyMultipleParse(io spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceRequestWritePropertyMultiple(), nil
}

func (m BACnetConfirmedServiceRequestWritePropertyMultiple) Serialize(io spi.WriteBuffer) {
	ser := func() {

	}
	BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
