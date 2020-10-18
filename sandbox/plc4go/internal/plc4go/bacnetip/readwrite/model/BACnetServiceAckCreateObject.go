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
type BACnetServiceAckCreateObject struct {
	BACnetServiceAck
}

// The corresponding interface
type IBACnetServiceAckCreateObject interface {
	IBACnetServiceAck
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetServiceAckCreateObject) ServiceChoice() uint8 {
	return 0x0A
}

func (m BACnetServiceAckCreateObject) initialize() spi.Message {
	return m
}

func NewBACnetServiceAckCreateObject() BACnetServiceAckInitializer {
	return &BACnetServiceAckCreateObject{}
}

func CastIBACnetServiceAckCreateObject(structType interface{}) IBACnetServiceAckCreateObject {
	castFunc := func(typ interface{}) IBACnetServiceAckCreateObject {
		if iBACnetServiceAckCreateObject, ok := typ.(IBACnetServiceAckCreateObject); ok {
			return iBACnetServiceAckCreateObject
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetServiceAckCreateObject(structType interface{}) BACnetServiceAckCreateObject {
	castFunc := func(typ interface{}) BACnetServiceAckCreateObject {
		if sBACnetServiceAckCreateObject, ok := typ.(BACnetServiceAckCreateObject); ok {
			return sBACnetServiceAckCreateObject
		}
		return BACnetServiceAckCreateObject{}
	}
	return castFunc(structType)
}

func (m BACnetServiceAckCreateObject) LengthInBits() uint16 {
	var lengthInBits = m.BACnetServiceAck.LengthInBits()

	return lengthInBits
}

func (m BACnetServiceAckCreateObject) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetServiceAckCreateObjectParse(io *spi.ReadBuffer) (BACnetServiceAckInitializer, error) {

	// Create the instance
	return NewBACnetServiceAckCreateObject(), nil
}

func (m BACnetServiceAckCreateObject) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetServiceAckSerialize(io, m.BACnetServiceAck, CastIBACnetServiceAck(m), ser)
}
