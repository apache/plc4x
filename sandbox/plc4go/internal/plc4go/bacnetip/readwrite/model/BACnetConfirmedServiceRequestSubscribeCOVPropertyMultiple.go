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
type BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple struct {
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple) ServiceChoice() uint8 {
	return 0x1E
}

func (m BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple() BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple{}
}

func CastIBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple(structType interface{}) IBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple {
		if iBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple, ok := typ.(IBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple); ok {
			return iBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple(structType interface{}) BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple {
		if sBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple, ok := typ.(BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple); ok {
			return sBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple
		}
		return BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple) LengthInBits() uint16 {
	var lengthInBits = m.BACnetConfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestSubscribeCOVPropertyMultipleParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple(), nil
}

func (m BACnetConfirmedServiceRequestSubscribeCOVPropertyMultiple) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
