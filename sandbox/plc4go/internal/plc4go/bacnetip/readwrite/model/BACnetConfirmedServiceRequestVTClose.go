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
type BACnetConfirmedServiceRequestVTClose struct {
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestVTClose interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestVTClose) ServiceChoice() uint8 {
	return 0x16
}

func (m BACnetConfirmedServiceRequestVTClose) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestVTClose() BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestVTClose{}
}

func CastIBACnetConfirmedServiceRequestVTClose(structType interface{}) IBACnetConfirmedServiceRequestVTClose {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestVTClose {
		if iBACnetConfirmedServiceRequestVTClose, ok := typ.(IBACnetConfirmedServiceRequestVTClose); ok {
			return iBACnetConfirmedServiceRequestVTClose
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestVTClose(structType interface{}) BACnetConfirmedServiceRequestVTClose {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestVTClose {
		if sBACnetConfirmedServiceRequestVTClose, ok := typ.(BACnetConfirmedServiceRequestVTClose); ok {
			return sBACnetConfirmedServiceRequestVTClose
		}
		return BACnetConfirmedServiceRequestVTClose{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestVTClose) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestVTClose) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestVTCloseParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceRequestVTClose(), nil
}

func (m BACnetConfirmedServiceRequestVTClose) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
