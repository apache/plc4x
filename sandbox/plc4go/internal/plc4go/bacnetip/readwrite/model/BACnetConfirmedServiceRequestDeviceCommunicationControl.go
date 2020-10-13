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
type BACnetConfirmedServiceRequestDeviceCommunicationControl struct {
	BACnetConfirmedServiceRequest
}

// The corresponding interface
type IBACnetConfirmedServiceRequestDeviceCommunicationControl interface {
	IBACnetConfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceRequestDeviceCommunicationControl) ServiceChoice() uint8 {
	return 0x11
}

func (m BACnetConfirmedServiceRequestDeviceCommunicationControl) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceRequestDeviceCommunicationControl() BACnetConfirmedServiceRequestInitializer {
	return &BACnetConfirmedServiceRequestDeviceCommunicationControl{}
}

func CastIBACnetConfirmedServiceRequestDeviceCommunicationControl(structType interface{}) IBACnetConfirmedServiceRequestDeviceCommunicationControl {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceRequestDeviceCommunicationControl {
		if iBACnetConfirmedServiceRequestDeviceCommunicationControl, ok := typ.(IBACnetConfirmedServiceRequestDeviceCommunicationControl); ok {
			return iBACnetConfirmedServiceRequestDeviceCommunicationControl
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceRequestDeviceCommunicationControl(structType interface{}) BACnetConfirmedServiceRequestDeviceCommunicationControl {
	castFunc := func(typ interface{}) BACnetConfirmedServiceRequestDeviceCommunicationControl {
		if sBACnetConfirmedServiceRequestDeviceCommunicationControl, ok := typ.(BACnetConfirmedServiceRequestDeviceCommunicationControl); ok {
			return sBACnetConfirmedServiceRequestDeviceCommunicationControl
		}
		return BACnetConfirmedServiceRequestDeviceCommunicationControl{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceRequestDeviceCommunicationControl) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceRequest.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceRequestDeviceCommunicationControl) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceRequestDeviceCommunicationControlParse(io spi.ReadBuffer) (BACnetConfirmedServiceRequestInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceRequestDeviceCommunicationControl(), nil
}

func (m BACnetConfirmedServiceRequestDeviceCommunicationControl) Serialize(io spi.WriteBuffer) {
	ser := func() {

	}
	BACnetConfirmedServiceRequestSerialize(io, m.BACnetConfirmedServiceRequest, CastIBACnetConfirmedServiceRequest(m), ser)
}
