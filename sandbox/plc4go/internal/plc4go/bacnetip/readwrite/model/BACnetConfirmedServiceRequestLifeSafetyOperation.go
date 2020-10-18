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
type BACnetConfirmedServiceRequestLifeSafetyOperation struct {
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestLifeSafetyOperation interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestLifeSafetyOperation) ServiceChoice() uint8 {
	return 0x1B
}

func (m BACnetConfirmedServiceRequestLifeSafetyOperation) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestLifeSafetyOperation() BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestLifeSafetyOperation{}
}

func CastIBACnetConfirmedServiceRequestLifeSafetyOperation(structType interface{}) IBACnetConfirmedServiceRequestLifeSafetyOperation {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestLifeSafetyOperation {
		if iBACnetConfirmedServiceRequestLifeSafetyOperation, ok := typ.(IBACnetConfirmedServiceRequestLifeSafetyOperation); ok {
			return iBACnetConfirmedServiceRequestLifeSafetyOperation
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestLifeSafetyOperation(structType interface{}) BACnetConfirmedServiceRequestLifeSafetyOperation {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestLifeSafetyOperation {
		if sBACnetConfirmedServiceRequestLifeSafetyOperation, ok := typ.(BACnetConfirmedServiceRequestLifeSafetyOperation); ok {
			return sBACnetConfirmedServiceRequestLifeSafetyOperation
		}
		return BACnetConfirmedServiceRequestLifeSafetyOperation{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestLifeSafetyOperation) LengthInBits() uint16 {
	var lengthInBits = m.BACnetConfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestLifeSafetyOperation) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestLifeSafetyOperationParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceRequestLifeSafetyOperation(), nil
}

func (m BACnetConfirmedServiceRequestLifeSafetyOperation) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
