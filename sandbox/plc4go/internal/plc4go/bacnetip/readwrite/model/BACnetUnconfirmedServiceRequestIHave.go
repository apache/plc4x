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
type BACnetUnconfirmedServiceRequestIHave struct {
	BACnetUnconfirmedServiceRequest
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestIHave interface {
	IBACnetUnconfirmedServiceRequest
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetUnconfirmedServiceRequestIHave) ServiceChoice() uint8 {
	return 0x01
}

func (m BACnetUnconfirmedServiceRequestIHave) initialize() spi.Message {
	return m
}

func NewBACnetUnconfirmedServiceRequestIHave() BACnetUnconfirmedServiceRequestInitializer {
	return &BACnetUnconfirmedServiceRequestIHave{}
}

func CastIBACnetUnconfirmedServiceRequestIHave(structType interface{}) IBACnetUnconfirmedServiceRequestIHave {
	castFunc := func(typ interface{}) IBACnetUnconfirmedServiceRequestIHave {
		if iBACnetUnconfirmedServiceRequestIHave, ok := typ.(IBACnetUnconfirmedServiceRequestIHave); ok {
			return iBACnetUnconfirmedServiceRequestIHave
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetUnconfirmedServiceRequestIHave(structType interface{}) BACnetUnconfirmedServiceRequestIHave {
	castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequestIHave {
		if sBACnetUnconfirmedServiceRequestIHave, ok := typ.(BACnetUnconfirmedServiceRequestIHave); ok {
			return sBACnetUnconfirmedServiceRequestIHave
		}
		return BACnetUnconfirmedServiceRequestIHave{}
	}
	return castFunc(structType)
}

func (m BACnetUnconfirmedServiceRequestIHave) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetUnconfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetUnconfirmedServiceRequestIHave) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestIHaveParse(io *spi.ReadBuffer) (BACnetUnconfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetUnconfirmedServiceRequestIHave(), nil
}

func (m BACnetUnconfirmedServiceRequestIHave) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetUnconfirmedServiceRequestSerialize(io, m.BACnetUnconfirmedServiceRequest, CastIBACnetUnconfirmedServiceRequest(m), ser)
}
