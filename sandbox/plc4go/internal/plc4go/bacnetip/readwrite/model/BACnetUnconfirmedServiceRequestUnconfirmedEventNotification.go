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
type BACnetUnconfirmedServiceRequestUnconfirmedEventNotification struct {
	BACnetUnconfirmedServiceRequest
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestUnconfirmedEventNotification interface {
	IBACnetUnconfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) ServiceChoice() uint8 {
	return 0x03
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) initialize() spi.Message {
	return m
}

func NewBACnetUnconfirmedServiceRequestUnconfirmedEventNotification() BACnetUnconfirmedServiceRequestInitializer {
	return &BACnetUnconfirmedServiceRequestUnconfirmedEventNotification{}
}

func CastIBACnetUnconfirmedServiceRequestUnconfirmedEventNotification(structType interface{}) IBACnetUnconfirmedServiceRequestUnconfirmedEventNotification {
	castFunc := func(typ interface{}) IBACnetUnconfirmedServiceRequestUnconfirmedEventNotification {
		if iBACnetUnconfirmedServiceRequestUnconfirmedEventNotification, ok := typ.(IBACnetUnconfirmedServiceRequestUnconfirmedEventNotification); ok {
			return iBACnetUnconfirmedServiceRequestUnconfirmedEventNotification
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetUnconfirmedServiceRequestUnconfirmedEventNotification(structType interface{}) BACnetUnconfirmedServiceRequestUnconfirmedEventNotification {
	castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequestUnconfirmedEventNotification {
		if sBACnetUnconfirmedServiceRequestUnconfirmedEventNotification, ok := typ.(BACnetUnconfirmedServiceRequestUnconfirmedEventNotification); ok {
			return sBACnetUnconfirmedServiceRequestUnconfirmedEventNotification
		}
		return BACnetUnconfirmedServiceRequestUnconfirmedEventNotification{}
	}
	return castFunc(structType)
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetUnconfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestUnconfirmedEventNotificationParse(io spi.ReadBuffer) (BACnetUnconfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetUnconfirmedServiceRequestUnconfirmedEventNotification(), nil
}

func (m BACnetUnconfirmedServiceRequestUnconfirmedEventNotification) Serialize(io spi.WriteBuffer) {
	ser := func() {

	}
	BACnetUnconfirmedServiceRequestSerialize(io, m.BACnetUnconfirmedServiceRequest, CastIBACnetUnconfirmedServiceRequest(m), ser)
}
