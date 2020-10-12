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
type BACnetConfirmedServiceRequestRemovedRequestKey struct {
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestRemovedRequestKey interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestRemovedRequestKey) ServiceChoice() uint8 {
	return 0x19
}

func (m BACnetConfirmedServiceRequestRemovedRequestKey) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestRemovedRequestKey() BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestRemovedRequestKey{}
}

func CastIBACnetConfirmedServiceRequestRemovedRequestKey(structType interface{}) IBACnetConfirmedServiceRequestRemovedRequestKey {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestRemovedRequestKey {
		if iBACnetConfirmedServiceRequestRemovedRequestKey, ok := typ.(IBACnetConfirmedServiceRequestRemovedRequestKey); ok {
			return iBACnetConfirmedServiceRequestRemovedRequestKey
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestRemovedRequestKey(structType interface{}) BACnetConfirmedServiceRequestRemovedRequestKey {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestRemovedRequestKey {
		if sBACnetConfirmedServiceRequestRemovedRequestKey, ok := typ.(BACnetConfirmedServiceRequestRemovedRequestKey); ok {
			return sBACnetConfirmedServiceRequestRemovedRequestKey
		}
		return BACnetConfirmedServiceRequestRemovedRequestKey{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestRemovedRequestKey) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestRemovedRequestKey) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestRemovedRequestKeyParse(io spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceRequestRemovedRequestKey(), nil
}

func (m BACnetConfirmedServiceRequestRemovedRequestKey) Serialize(io spi.WriteBuffer) {

}
