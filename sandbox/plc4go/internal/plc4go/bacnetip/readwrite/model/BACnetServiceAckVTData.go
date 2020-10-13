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
type BACnetServiceAckVTData struct {
	BACnetServiceAck
}

// The corresponding interface
type IBACnetServiceAckVTData interface {
	IBACnetServiceAck
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetServiceAckVTData) ServiceChoice() uint8 {
	return 0x17
}

func (m BACnetServiceAckVTData) initialize() spi.Message {
	return m
}

func NewBACnetServiceAckVTData() BACnetServiceAckInitializer {
	return &BACnetServiceAckVTData{}
}

func CastIBACnetServiceAckVTData(structType interface{}) IBACnetServiceAckVTData {
	castFunc := func(typ interface{}) IBACnetServiceAckVTData {
		if iBACnetServiceAckVTData, ok := typ.(IBACnetServiceAckVTData); ok {
			return iBACnetServiceAckVTData
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetServiceAckVTData(structType interface{}) BACnetServiceAckVTData {
	castFunc := func(typ interface{}) BACnetServiceAckVTData {
		if sBACnetServiceAckVTData, ok := typ.(BACnetServiceAckVTData); ok {
			return sBACnetServiceAckVTData
		}
		return BACnetServiceAckVTData{}
	}
	return castFunc(structType)
}

func (m BACnetServiceAckVTData) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetServiceAck.LengthInBits()

	return lengthInBits
}

func (m BACnetServiceAckVTData) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetServiceAckVTDataParse(io spi.ReadBuffer) (BACnetServiceAckInitializer, error) {

	// Create the instance
	return NewBACnetServiceAckVTData(), nil
}

func (m BACnetServiceAckVTData) Serialize(io spi.WriteBuffer) {
	ser := func() {

	}
	BACnetServiceAckSerialize(io, m.BACnetServiceAck, CastIBACnetServiceAck(m), ser)
}
