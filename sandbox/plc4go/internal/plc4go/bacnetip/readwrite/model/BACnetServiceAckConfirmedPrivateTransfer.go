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
type BACnetServiceAckConfirmedPrivateTransfer struct {
	BACnetServiceAck
}

// The corresponding interface
type IBACnetServiceAckConfirmedPrivateTransfer interface {
	IBACnetServiceAck
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetServiceAckConfirmedPrivateTransfer) ServiceChoice() uint8 {
	return 0x12
}

func (m BACnetServiceAckConfirmedPrivateTransfer) initialize() spi.Message {
	return m
}

func NewBACnetServiceAckConfirmedPrivateTransfer() BACnetServiceAckInitializer {
	return &BACnetServiceAckConfirmedPrivateTransfer{}
}

func CastIBACnetServiceAckConfirmedPrivateTransfer(structType interface{}) IBACnetServiceAckConfirmedPrivateTransfer {
	castFunc := func(typ interface{}) IBACnetServiceAckConfirmedPrivateTransfer {
		if iBACnetServiceAckConfirmedPrivateTransfer, ok := typ.(IBACnetServiceAckConfirmedPrivateTransfer); ok {
			return iBACnetServiceAckConfirmedPrivateTransfer
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetServiceAckConfirmedPrivateTransfer(structType interface{}) BACnetServiceAckConfirmedPrivateTransfer {
	castFunc := func(typ interface{}) BACnetServiceAckConfirmedPrivateTransfer {
		if sBACnetServiceAckConfirmedPrivateTransfer, ok := typ.(BACnetServiceAckConfirmedPrivateTransfer); ok {
			return sBACnetServiceAckConfirmedPrivateTransfer
		}
		return BACnetServiceAckConfirmedPrivateTransfer{}
	}
	return castFunc(structType)
}

func (m BACnetServiceAckConfirmedPrivateTransfer) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetServiceAck.LengthInBits()

	return lengthInBits
}

func (m BACnetServiceAckConfirmedPrivateTransfer) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetServiceAckConfirmedPrivateTransferParse(io *spi.ReadBuffer) (BACnetServiceAckInitializer, error) {

	// Create the instance
	return NewBACnetServiceAckConfirmedPrivateTransfer(), nil
}

func (m BACnetServiceAckConfirmedPrivateTransfer) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetServiceAckSerialize(io, m.BACnetServiceAck, CastIBACnetServiceAck(m), ser)
}
