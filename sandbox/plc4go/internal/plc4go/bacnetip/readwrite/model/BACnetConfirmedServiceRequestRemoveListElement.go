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
type BACnetConfirmedServiceRequestRemoveListElement struct {
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestRemoveListElement interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestRemoveListElement) ServiceChoice() uint8 {
	return 0x09
}

func (m BACnetConfirmedServiceRequestRemoveListElement) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestRemoveListElement() BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestRemoveListElement{}
}

func CastIBACnetConfirmedServiceRequestRemoveListElement(structType interface{}) IBACnetConfirmedServiceRequestRemoveListElement {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestRemoveListElement {
		if iBACnetConfirmedServiceRequestRemoveListElement, ok := typ.(IBACnetConfirmedServiceRequestRemoveListElement); ok {
			return iBACnetConfirmedServiceRequestRemoveListElement
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestRemoveListElement(structType interface{}) BACnetConfirmedServiceRequestRemoveListElement {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestRemoveListElement {
		if sBACnetConfirmedServiceRequestRemoveListElement, ok := typ.(BACnetConfirmedServiceRequestRemoveListElement); ok {
			return sBACnetConfirmedServiceRequestRemoveListElement
		}
		return BACnetConfirmedServiceRequestRemoveListElement{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestRemoveListElement) LengthInBits() uint16 {
	var lengthInBits = m.BACnetConfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestRemoveListElement) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestRemoveListElementParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceRequestRemoveListElement(), nil
}

func (m BACnetConfirmedServiceRequestRemoveListElement) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
