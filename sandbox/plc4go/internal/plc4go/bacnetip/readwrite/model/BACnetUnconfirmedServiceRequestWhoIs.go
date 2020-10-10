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
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"strconv"
)

// Constant values.
const BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGELOWLIMITHEADER uint8 = 0x01
const BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGEHIGHLIMITHEADER uint8 = 0x03

// The data-structure of this message
type BACnetUnconfirmedServiceRequestWhoIs struct {
	deviceInstanceRangeLowLimitLength  uint8
	deviceInstanceRangeLowLimit        []int8
	deviceInstanceRangeHighLimitLength uint8
	deviceInstanceRangeHighLimit       []int8
	BACnetUnconfirmedServiceRequest
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestWhoIs interface {
	IBACnetUnconfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetUnconfirmedServiceRequestWhoIs) ServiceChoice() uint8 {
	return 0x08
}

func (m BACnetUnconfirmedServiceRequestWhoIs) initialize() spi.Message {
	return m
}

func NewBACnetUnconfirmedServiceRequestWhoIs(deviceInstanceRangeLowLimitLength uint8, deviceInstanceRangeLowLimit []int8, deviceInstanceRangeHighLimitLength uint8, deviceInstanceRangeHighLimit []int8) BACnetUnconfirmedServiceRequestInitializer {
	return &BACnetUnconfirmedServiceRequestWhoIs{deviceInstanceRangeLowLimitLength: deviceInstanceRangeLowLimitLength, deviceInstanceRangeLowLimit: deviceInstanceRangeLowLimit, deviceInstanceRangeHighLimitLength: deviceInstanceRangeHighLimitLength, deviceInstanceRangeHighLimit: deviceInstanceRangeHighLimit}
}

func (m BACnetUnconfirmedServiceRequestWhoIs) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetUnconfirmedServiceRequest.LengthInBits()

	// Const Field (deviceInstanceRangeLowLimitHeader)
	lengthInBits += 5

	// Simple field (deviceInstanceRangeLowLimitLength)
	lengthInBits += 3

	// Array field
	if len(m.deviceInstanceRangeLowLimit) > 0 {
		lengthInBits += 8 * uint16(len(m.deviceInstanceRangeLowLimit))
	}

	// Const Field (deviceInstanceRangeHighLimitHeader)
	lengthInBits += 5

	// Simple field (deviceInstanceRangeHighLimitLength)
	lengthInBits += 3

	// Array field
	if len(m.deviceInstanceRangeHighLimit) > 0 {
		lengthInBits += 8 * uint16(len(m.deviceInstanceRangeHighLimit))
	}

	return lengthInBits
}

func (m BACnetUnconfirmedServiceRequestWhoIs) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestWhoIsParse(io spi.ReadBuffer) (BACnetUnconfirmedServiceRequestInitializer, error) {

	// Const Field (deviceInstanceRangeLowLimitHeader)
	var deviceInstanceRangeLowLimitHeader uint8 = io.ReadUint8(5)
	if deviceInstanceRangeLowLimitHeader != BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGELOWLIMITHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGELOWLIMITHEADER)) + " but got " + strconv.Itoa(int(deviceInstanceRangeLowLimitHeader)))
	}

	// Simple Field (deviceInstanceRangeLowLimitLength)
	var deviceInstanceRangeLowLimitLength uint8 = io.ReadUint8(3)

	// Array field (deviceInstanceRangeLowLimit)
	var deviceInstanceRangeLowLimit []int8
	// Count array
	{
		deviceInstanceRangeLowLimit := make([]int8, deviceInstanceRangeLowLimitLength)
		for curItem := uint16(0); curItem < uint16(deviceInstanceRangeLowLimitLength); curItem++ {

			deviceInstanceRangeLowLimit = append(deviceInstanceRangeLowLimit, io.ReadInt8(8))
		}
	}

	// Const Field (deviceInstanceRangeHighLimitHeader)
	var deviceInstanceRangeHighLimitHeader uint8 = io.ReadUint8(5)
	if deviceInstanceRangeHighLimitHeader != BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGEHIGHLIMITHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestWhoIs_DEVICEINSTANCERANGEHIGHLIMITHEADER)) + " but got " + strconv.Itoa(int(deviceInstanceRangeHighLimitHeader)))
	}

	// Simple Field (deviceInstanceRangeHighLimitLength)
	var deviceInstanceRangeHighLimitLength uint8 = io.ReadUint8(3)

	// Array field (deviceInstanceRangeHighLimit)
	var deviceInstanceRangeHighLimit []int8
	// Count array
	{
		deviceInstanceRangeHighLimit := make([]int8, deviceInstanceRangeHighLimitLength)
		for curItem := uint16(0); curItem < uint16(deviceInstanceRangeHighLimitLength); curItem++ {

			deviceInstanceRangeHighLimit = append(deviceInstanceRangeHighLimit, io.ReadInt8(8))
		}
	}

	// Create the instance
	return NewBACnetUnconfirmedServiceRequestWhoIs(deviceInstanceRangeLowLimitLength, deviceInstanceRangeLowLimit, deviceInstanceRangeHighLimitLength, deviceInstanceRangeHighLimit), nil
}

func (m BACnetUnconfirmedServiceRequestWhoIs) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IBACnetUnconfirmedServiceRequestWhoIs); ok {

			// Const Field (deviceInstanceRangeLowLimitHeader)
			io.WriteUint8(5, 0x01)

			// Simple Field (deviceInstanceRangeLowLimitLength)
			var deviceInstanceRangeLowLimitLength uint8 = m.deviceInstanceRangeLowLimitLength
			io.WriteUint8(3, (deviceInstanceRangeLowLimitLength))

			// Array Field (deviceInstanceRangeLowLimit)
			if m.deviceInstanceRangeLowLimit != nil {
				for _, _element := range m.deviceInstanceRangeLowLimit {
					io.WriteInt8(8, _element)
				}
			}

			// Const Field (deviceInstanceRangeHighLimitHeader)
			io.WriteUint8(5, 0x03)

			// Simple Field (deviceInstanceRangeHighLimitLength)
			var deviceInstanceRangeHighLimitLength uint8 = m.deviceInstanceRangeHighLimitLength
			io.WriteUint8(3, (deviceInstanceRangeHighLimitLength))

			// Array Field (deviceInstanceRangeHighLimit)
			if m.deviceInstanceRangeHighLimit != nil {
				for _, _element := range m.deviceInstanceRangeHighLimit {
					io.WriteInt8(8, _element)
				}
			}
		}
	}
	serializeFunc(m)
}
