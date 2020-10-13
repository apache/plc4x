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
type BACnetConfirmedServiceRequestConfirmedEventNotification struct {
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestConfirmedEventNotification interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestConfirmedEventNotification) ServiceChoice() uint8 {
	return 0x02
}

func (m BACnetConfirmedServiceRequestConfirmedEventNotification) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestConfirmedEventNotification() BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestConfirmedEventNotification{}
}

func CastIBACnetConfirmedServiceRequestConfirmedEventNotification(structType interface{}) IBACnetConfirmedServiceRequestConfirmedEventNotification {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestConfirmedEventNotification {
		if iBACnetConfirmedServiceRequestConfirmedEventNotification, ok := typ.(IBACnetConfirmedServiceRequestConfirmedEventNotification); ok {
			return iBACnetConfirmedServiceRequestConfirmedEventNotification
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestConfirmedEventNotification(structType interface{}) BACnetConfirmedServiceRequestConfirmedEventNotification {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestConfirmedEventNotification {
		if sBACnetConfirmedServiceRequestConfirmedEventNotification, ok := typ.(BACnetConfirmedServiceRequestConfirmedEventNotification); ok {
			return sBACnetConfirmedServiceRequestConfirmedEventNotification
		}
		return BACnetConfirmedServiceRequestConfirmedEventNotification{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestConfirmedEventNotification) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestConfirmedEventNotification) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestConfirmedEventNotificationParse(io *spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceRequestConfirmedEventNotification(), nil
}

func (m BACnetConfirmedServiceRequestConfirmedEventNotification) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
