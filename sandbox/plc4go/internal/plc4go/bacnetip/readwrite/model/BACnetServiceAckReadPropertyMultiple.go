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
type BACnetServiceAckReadPropertyMultiple struct {
	BACnetServiceAck
}

// The corresponding interface
type IBACnetServiceAckReadPropertyMultiple interface {
	IBACnetServiceAck
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetServiceAckReadPropertyMultiple) ServiceChoice() uint8 {
	return 0x0E
}

func (m BACnetServiceAckReadPropertyMultiple) initialize() spi.Message {
	return m
}

func NewBACnetServiceAckReadPropertyMultiple() BACnetServiceAckInitializer {
	return &BACnetServiceAckReadPropertyMultiple{}
}

func CastIBACnetServiceAckReadPropertyMultiple(structType interface{}) IBACnetServiceAckReadPropertyMultiple {
	castFunc := func(typ interface{}) IBACnetServiceAckReadPropertyMultiple {
		if iBACnetServiceAckReadPropertyMultiple, ok := typ.(IBACnetServiceAckReadPropertyMultiple); ok {
			return iBACnetServiceAckReadPropertyMultiple
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetServiceAckReadPropertyMultiple(structType interface{}) BACnetServiceAckReadPropertyMultiple {
	castFunc := func(typ interface{}) BACnetServiceAckReadPropertyMultiple {
		if sBACnetServiceAckReadPropertyMultiple, ok := typ.(BACnetServiceAckReadPropertyMultiple); ok {
			return sBACnetServiceAckReadPropertyMultiple
		}
		return BACnetServiceAckReadPropertyMultiple{}
	}
	return castFunc(structType)
}

func (m BACnetServiceAckReadPropertyMultiple) LengthInBits() uint16 {
	var lengthInBits = m.BACnetServiceAck.LengthInBits()

	return lengthInBits
}

func (m BACnetServiceAckReadPropertyMultiple) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetServiceAckReadPropertyMultipleParse(io *spi.ReadBuffer) (BACnetServiceAckInitializer, error) {

	// Create the instance
	return NewBACnetServiceAckReadPropertyMultiple(), nil
}

func (m BACnetServiceAckReadPropertyMultiple) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetServiceAckSerialize(io, m.BACnetServiceAck, CastIBACnetServiceAck(m), ser)
}
