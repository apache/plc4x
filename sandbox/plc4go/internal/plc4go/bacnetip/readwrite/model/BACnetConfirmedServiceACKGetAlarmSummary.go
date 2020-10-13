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
type BACnetConfirmedServiceACKGetAlarmSummary struct {
	BACnetConfirmedServiceACK
}

// The corresponding interface
type IBACnetConfirmedServiceACKGetAlarmSummary interface {
	IBACnetConfirmedServiceACK
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetConfirmedServiceACKGetAlarmSummary) ServiceChoice() uint8 {
	return 0x03
}

func (m BACnetConfirmedServiceACKGetAlarmSummary) initialize() spi.Message {
	return m
}

func NewBACnetConfirmedServiceACKGetAlarmSummary() BACnetConfirmedServiceACKInitializer {
	return &BACnetConfirmedServiceACKGetAlarmSummary{}
}

func CastIBACnetConfirmedServiceACKGetAlarmSummary(structType interface{}) IBACnetConfirmedServiceACKGetAlarmSummary {
	castFunc := func(typ interface{}) IBACnetConfirmedServiceACKGetAlarmSummary {
		if iBACnetConfirmedServiceACKGetAlarmSummary, ok := typ.(IBACnetConfirmedServiceACKGetAlarmSummary); ok {
			return iBACnetConfirmedServiceACKGetAlarmSummary
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetConfirmedServiceACKGetAlarmSummary(structType interface{}) BACnetConfirmedServiceACKGetAlarmSummary {
	castFunc := func(typ interface{}) BACnetConfirmedServiceACKGetAlarmSummary {
		if sBACnetConfirmedServiceACKGetAlarmSummary, ok := typ.(BACnetConfirmedServiceACKGetAlarmSummary); ok {
			return sBACnetConfirmedServiceACKGetAlarmSummary
		}
		return BACnetConfirmedServiceACKGetAlarmSummary{}
	}
	return castFunc(structType)
}

func (m BACnetConfirmedServiceACKGetAlarmSummary) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetConfirmedServiceACK.LengthInBits()

	return lengthInBits
}

func (m BACnetConfirmedServiceACKGetAlarmSummary) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetConfirmedServiceACKGetAlarmSummaryParse(io *spi.ReadBuffer) (BACnetConfirmedServiceACKInitializer, error) {

	// Create the instance
	return NewBACnetConfirmedServiceACKGetAlarmSummary(), nil
}

func (m BACnetConfirmedServiceACKGetAlarmSummary) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return BACnetConfirmedServiceACKSerialize(io, m.BACnetConfirmedServiceACK, CastIBACnetConfirmedServiceACK(m), ser)
}
