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
    "plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/spi"
	"plc4x.apache.org/plc4go-modbus-driver/v0/internal/plc4go/utils"
	"strconv"
)

// Constant values.
const BACnetUnconfirmedServiceRequestWhoHas_DEVICEINSTANCELOWLIMITHEADER uint8 = 0x0B
const BACnetUnconfirmedServiceRequestWhoHas_DEVICEINSTANCEHIGHLIMITHEADER uint8 = 0x1B
const BACnetUnconfirmedServiceRequestWhoHas_OBJECTNAMEHEADER uint8 = 0x3D

// The data-structure of this message
type BACnetUnconfirmedServiceRequestWhoHas struct {
    DeviceInstanceLowLimit uint32
    DeviceInstanceHighLimit uint32
    ObjectNameCharacterSet uint8
    ObjectName []int8
    BACnetUnconfirmedServiceRequest
}

// The corresponding interface
type IBACnetUnconfirmedServiceRequestWhoHas interface {
    IBACnetUnconfirmedServiceRequest
    Serialize(io utils.WriteBuffer) error
}

// Accessors for discriminator values.
func (m BACnetUnconfirmedServiceRequestWhoHas) ServiceChoice() uint8 {
    return 0x07
}

func (m BACnetUnconfirmedServiceRequestWhoHas) initialize() spi.Message {
    return m
}

func NewBACnetUnconfirmedServiceRequestWhoHas(deviceInstanceLowLimit uint32, deviceInstanceHighLimit uint32, objectNameCharacterSet uint8, objectName []int8) BACnetUnconfirmedServiceRequestInitializer {
    return &BACnetUnconfirmedServiceRequestWhoHas{DeviceInstanceLowLimit: deviceInstanceLowLimit, DeviceInstanceHighLimit: deviceInstanceHighLimit, ObjectNameCharacterSet: objectNameCharacterSet, ObjectName: objectName}
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
        if sBACnetUnconfirmedServiceRequestWhoHas, ok := typ.(*BACnetUnconfirmedServiceRequestWhoHas); ok {
            return *sBACnetUnconfirmedServiceRequestWhoHas
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
    if len(m.ObjectName) > 0 {
        lengthInBits += 8 * uint16(len(m.ObjectName))
    }

    return lengthInBits
}

func (m BACnetUnconfirmedServiceRequestWhoHas) LengthInBytes() uint16 {
    return m.LengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestWhoHasParse(io *utils.ReadBuffer) (BACnetUnconfirmedServiceRequestInitializer, error) {

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
    // Length array
    objectName := make([]int8, 0)
    _objectNameLength := uint16(objectNameLength) - uint16(uint16(1))
    _objectNameEndPos := io.GetPos() + uint16(_objectNameLength)
    for ;io.GetPos() < _objectNameEndPos; {
        _item, _err := io.ReadInt8(8)
        if _err != nil {
            return nil, errors.New("Error parsing 'objectName' field " + _err.Error())
        }
        objectName = append(objectName, _item)
    }

    // Create the instance
    return NewBACnetUnconfirmedServiceRequestWhoHas(deviceInstanceLowLimit, deviceInstanceHighLimit, objectNameCharacterSet, objectName), nil
}

func (m BACnetUnconfirmedServiceRequestWhoHas) Serialize(io utils.WriteBuffer) error {
    ser := func() error {

    // Const Field (deviceInstanceLowLimitHeader)
    _deviceInstanceLowLimitHeaderErr := io.WriteUint8(8, 0x0B)
    if _deviceInstanceLowLimitHeaderErr != nil {
        return errors.New("Error serializing 'deviceInstanceLowLimitHeader' field " + _deviceInstanceLowLimitHeaderErr.Error())
    }

    // Simple Field (deviceInstanceLowLimit)
    deviceInstanceLowLimit := uint32(m.DeviceInstanceLowLimit)
    _deviceInstanceLowLimitErr := io.WriteUint32(24, (deviceInstanceLowLimit))
    if _deviceInstanceLowLimitErr != nil {
        return errors.New("Error serializing 'deviceInstanceLowLimit' field " + _deviceInstanceLowLimitErr.Error())
    }

    // Const Field (deviceInstanceHighLimitHeader)
    _deviceInstanceHighLimitHeaderErr := io.WriteUint8(8, 0x1B)
    if _deviceInstanceHighLimitHeaderErr != nil {
        return errors.New("Error serializing 'deviceInstanceHighLimitHeader' field " + _deviceInstanceHighLimitHeaderErr.Error())
    }

    // Simple Field (deviceInstanceHighLimit)
    deviceInstanceHighLimit := uint32(m.DeviceInstanceHighLimit)
    _deviceInstanceHighLimitErr := io.WriteUint32(24, (deviceInstanceHighLimit))
    if _deviceInstanceHighLimitErr != nil {
        return errors.New("Error serializing 'deviceInstanceHighLimit' field " + _deviceInstanceHighLimitErr.Error())
    }

    // Const Field (objectNameHeader)
    _objectNameHeaderErr := io.WriteUint8(8, 0x3D)
    if _objectNameHeaderErr != nil {
        return errors.New("Error serializing 'objectNameHeader' field " + _objectNameHeaderErr.Error())
    }

    // Implicit Field (objectNameLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
    objectNameLength := uint8(uint8(uint8(len(m.ObjectName))) + uint8(uint8(1)))
    _objectNameLengthErr := io.WriteUint8(8, (objectNameLength))
    if _objectNameLengthErr != nil {
        return errors.New("Error serializing 'objectNameLength' field " + _objectNameLengthErr.Error())
    }

    // Simple Field (objectNameCharacterSet)
    objectNameCharacterSet := uint8(m.ObjectNameCharacterSet)
    _objectNameCharacterSetErr := io.WriteUint8(8, (objectNameCharacterSet))
    if _objectNameCharacterSetErr != nil {
        return errors.New("Error serializing 'objectNameCharacterSet' field " + _objectNameCharacterSetErr.Error())
    }

    // Array Field (objectName)
    if m.ObjectName != nil {
        for _, _element := range m.ObjectName {
            _elementErr := io.WriteInt8(8, _element)
            if _elementErr != nil {
                return errors.New("Error serializing 'objectName' field " + _elementErr.Error())
            }
        }
    }

        return nil
    }
    return BACnetUnconfirmedServiceRequestSerialize(io, m.BACnetUnconfirmedServiceRequest, CastIBACnetUnconfirmedServiceRequest(m), ser)
}
