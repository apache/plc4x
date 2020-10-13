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
type BACnetConfirmedServiceRequestCreateObject struct {
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestCreateObject interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestCreateObject) ServiceChoice() uint8 {
	return 0x0A
}

func (m BACnetConfirmedServiceRequestCreateObject) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestCreateObject() BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestCreateObject{}
}

func CastIBACnetConfirmedServiceRequestCreateObject(structType interface{}) IBACnetConfirmedServiceRequestCreateObject {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestCreateObject {
		if iBACnetConfirmedServiceRequestCreateObject, ok := typ.(IBACnetConfirmedServiceRequestCreateObject); ok {
			return iBACnetConfirmedServiceRequestCreateObject
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestCreateObject(structType interface{}) BACnetConfirmedServiceRequestCreateObject {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestCreateObject {
		if sBACnetConfirmedServiceRequestCreateObject, ok := typ.(BACnetConfirmedServiceRequestCreateObject); ok {
			return sBACnetConfirmedServiceRequestCreateObject
		}
		return BACnetConfirmedServiceRequestCreateObject{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestCreateObject) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestCreateObject) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestCreateObjectParse(io spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceRequestCreateObject(), nil
}

func (m BACnetConfirmedServiceRequestCreateObject) Serialize(io spi.WriteBuffer) {
	ser := func() {

	}
	BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
