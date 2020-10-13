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
const BACnetUnconfirmedServiceRequestWhoHas_DEVICEINSTANCELOWLIMITHEADER uint8 = 0x0B
const BACnetUnconfirmedServiceRequestWhoHas_DEVICEINSTANCEHIGHLIMITHEADER uint8 = 0x1B
const BACnetUnconfirmedServiceRequestWhoHas_OBJECTNAMEHEADER uint8 = 0x3D

// The data-structure of this message
type BACnetUnconfirmedServiceRequestWhoHas struct {
	deviceInstanceLowLimit  uint32
	deviceInstanceHighLimit uint32
	objectNameCharacterSet  uint8
	objectName              []int8
	BACnetUnconfirmedServiceRequest
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestWhoHas interface {
	IBACnetUnconfirmedServiceRequest
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m BACnetUnconfirmedServiceRequestWhoHas) ServiceChoice() uint8 {
	return 0x07
}

func (m BACnetUnconfirmedServiceRequestWhoHas) initialize() spi.Message {
	return m
}

func NewBACnetUnconfirmedServiceRequestWhoHas(deviceInstanceLowLimit uint32, deviceInstanceHighLimit uint32, objectNameCharacterSet uint8, objectName []int8) BACnetUnconfirmedServiceRequestInitializer {
	return &BACnetUnconfirmedServiceRequestWhoHas{deviceInstanceLowLimit: deviceInstanceLowLimit, deviceInstanceHighLimit: deviceInstanceHighLimit, objectNameCharacterSet: objectNameCharacterSet, objectName: objectName}
}

func CastIBACnetUnconfirmedServiceRequestWhoHas(structType interface{}) IBACnetUnconfirmedServiceRequestWhoHas {
	castFunc := func(typ interface{}) IBACnetUnconfirmedServiceRequestWhoHas {
		if iBACnetUnconfirmedServiceRequestWhoHas, ok := typ.(IBACnetUnconfirmedServiceRequestWhoHas); ok {
			return iBACnetUnconfirmedServiceRequestWhoHas
		}
		return nil
	}
	return castFunc(structType)
}

func CastBACnetUnconfirmedServiceRequestWhoHas(structType interface{}) BACnetUnconfirmedServiceRequestWhoHas {
	castFunc := func(typ interface{}) BACnetUnconfirmedServiceRequestWhoHas {
		if sBACnetUnconfirmedServiceRequestWhoHas, ok := typ.(BACnetUnconfirmedServiceRequestWhoHas); ok {
			return sBACnetUnconfirmedServiceRequestWhoHas
		}
		return BACnetUnconfirmedServiceRequestWhoHas{}
	}
	return castFunc(structType)
}

func (m BACnetUnconfirmedServiceRequestWhoHas) LengthInBits() uint16 {
	var lengthInBits uint16 = m.BACnetUnconfirmedServiceRequest.LengthInBits()

	// Const Field (deviceInstanceLowLimitHeader)
	lengthInBits += 8

	// Simple field (deviceInstanceLowLimit)
	lengthInBits += 24

	// Const Field (deviceInstanceHighLimitHeader)
	lengthInBits += 8

	// Simple field (deviceInstanceHighLimit)
	lengthInBits += 24

	// Const Field (objectNameHeader)
	lengthInBits += 8

	// Implicit Field (objectNameLength)
	lengthInBits += 8

	// Simple field (objectNameCharacterSet)
	lengthInBits += 8

	// Array field
	if len(m.objectName) > 0 {
		lengthInBits += 8 * uint16(len(m.objectName))
	}

	return lengthInBits
}

func (m BACnetUnconfirmedServiceRequestWhoHas) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestWhoHasParse(io spi.ReadBuffer) (BACnetUnconfirmedServiceRequestInitializer, error) {

	// Const Field (deviceInstanceLowLimitHeader)
	deviceInstanceLowLimitHeader, _deviceInstanceLowLimitHeaderErr := io.ReadUint8(8)
	if _deviceInstanceLowLimitHeaderErr != nil {
		return nil, errors.New("Error parsing 'deviceInstanceLowLimitHeader' field " + _deviceInstanceLowLimitHeaderErr.Error())
	}
	if deviceInstanceLowLimitHeader != BACnetUnconfirmedServiceRequestWhoHas_DEVICEINSTANCELOWLIMITHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestWhoHas_DEVICEINSTANCELOWLIMITHEADER)) + " but got " + strconv.Itoa(int(deviceInstanceLowLimitHeader)))
	}

	// Simple Field (deviceInstanceLowLimit)
	deviceInstanceLowLimit, _deviceInstanceLowLimitErr := io.ReadUint32(24)
	if _deviceInstanceLowLimitErr != nil {
		return nil, errors.New("Error parsing 'deviceInstanceLowLimit' field " + _deviceInstanceLowLimitErr.Error())
	}

	// Const Field (deviceInstanceHighLimitHeader)
	deviceInstanceHighLimitHeader, _deviceInstanceHighLimitHeaderErr := io.ReadUint8(8)
	if _deviceInstanceHighLimitHeaderErr != nil {
		return nil, errors.New("Error parsing 'deviceInstanceHighLimitHeader' field " + _deviceInstanceHighLimitHeaderErr.Error())
	}
	if deviceInstanceHighLimitHeader != BACnetUnconfirmedServiceRequestWhoHas_DEVICEINSTANCEHIGHLIMITHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestWhoHas_DEVICEINSTANCEHIGHLIMITHEADER)) + " but got " + strconv.Itoa(int(deviceInstanceHighLimitHeader)))
	}

	// Simple Field (deviceInstanceHighLimit)
	deviceInstanceHighLimit, _deviceInstanceHighLimitErr := io.ReadUint32(24)
	if _deviceInstanceHighLimitErr != nil {
		return nil, errors.New("Error parsing 'deviceInstanceHighLimit' field " + _deviceInstanceHighLimitErr.Error())
	}

	// Const Field (objectNameHeader)
	objectNameHeader, _objectNameHeaderErr := io.ReadUint8(8)
	if _objectNameHeaderErr != nil {
		return nil, errors.New("Error parsing 'objectNameHeader' field " + _objectNameHeaderErr.Error())
	}
	if objectNameHeader != BACnetUnconfirmedServiceRequestWhoHas_OBJECTNAMEHEADER {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(BACnetUnconfirmedServiceRequestWhoHas_OBJECTNAMEHEADER)) + " but got " + strconv.Itoa(int(objectNameHeader)))
	}

	// Implicit Field (objectNameLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	objectNameLength, _objectNameLengthErr := io.ReadUint8(8)
	if _objectNameLengthErr != nil {
		return nil, errors.New("Error parsing 'objectNameLength' field " + _objectNameLengthErr.Error())
	}

	// Simple Field (objectNameCharacterSet)
	objectNameCharacterSet, _objectNameCharacterSetErr := io.ReadUint8(8)
	if _objectNameCharacterSetErr != nil {
		return nil, errors.New("Error parsing 'objectNameCharacterSet' field " + _objectNameCharacterSetErr.Error())
	}

	// Array field (objectName)
	var objectName []int8
	// Length array
	_objectNameLength := uint16(objectNameLength) - uint16(uint16(1))
	_objectNameEndPos := io.GetPos() + uint16(_objectNameLength)
	for io.GetPos() < _objectNameEndPos {
		_objectNameVal, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'objectName' field " + _err.Error())
		}
		objectName = append(objectName, _objectNameVal)
	}

	// Create the instance
	return NewBACnetUnconfirmedServiceRequestWhoHas(deviceInstanceLowLimit, deviceInstanceHighLimit, objectNameCharacterSet, objectName), nil
}

func (m BACnetUnconfirmedServiceRequestWhoHas) Serialize(io spi.WriteBuffer) {

	// Const Field (deviceInstanceLowLimitHeader)
	io.WriteUint8(8, 0x0B)

	// Simple Field (deviceInstanceLowLimit)
	deviceInstanceLowLimit := uint32(m.deviceInstanceLowLimit)
	io.WriteUint32(24, (deviceInstanceLowLimit))

	// Const Field (deviceInstanceHighLimitHeader)
	io.WriteUint8(8, 0x1B)

	// Simple Field (deviceInstanceHighLimit)
	deviceInstanceHighLimit := uint32(m.deviceInstanceHighLimit)
	io.WriteUint32(24, (deviceInstanceHighLimit))

	// Const Field (objectNameHeader)
	io.WriteUint8(8, 0x3D)

	// Implicit Field (objectNameLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	objectNameLength := uint8(uint8(uint8(len(m.objectName))) + uint8(uint8(1)))
	io.WriteUint8(8, (objectNameLength))

	// Simple Field (objectNameCharacterSet)
	objectNameCharacterSet := uint8(m.objectNameCharacterSet)
	io.WriteUint8(8, (objectNameCharacterSet))

	// Array Field (objectName)
	if m.objectName != nil {
		for _, _element := range m.objectName {
			io.WriteInt8(8, _element)
		}
	}
}
