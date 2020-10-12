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
	"reflect"
)

// The data-structure of this message
type DeviceConfigurationAck struct {
	deviceConfigurationAckDataBlock DeviceConfigurationAckDataBlock
	KNXNetIPMessage
}

// The corresponding interface
type IDeviceConfigurationAck interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m DeviceConfigurationAck) MsgType() uint16 {
	return 0x0311
}

func (m DeviceConfigurationAck) initialize() spi.Message {
	return m
}

func NewDeviceConfigurationAck(deviceConfigurationAckDataBlock DeviceConfigurationAckDataBlock) KNXNetIPMessageInitializer {
	return &DeviceConfigurationAck{deviceConfigurationAckDataBlock: deviceConfigurationAckDataBlock}
}

func CastIDeviceConfigurationAck(structType interface{}) IDeviceConfigurationAck {
	castFunc := func(typ interface{}) IDeviceConfigurationAck {
		if iDeviceConfigurationAck, ok := typ.(IDeviceConfigurationAck); ok {
			return iDeviceConfigurationAck
		}
		return nil
	}
	return castFunc(structType)
}

func CastDeviceConfigurationAck(structType interface{}) DeviceConfigurationAck {
	castFunc := func(typ interface{}) DeviceConfigurationAck {
		if sDeviceConfigurationAck, ok := typ.(DeviceConfigurationAck); ok {
			return sDeviceConfigurationAck
		}
		return DeviceConfigurationAck{}
	}
	return castFunc(structType)
}

func (m DeviceConfigurationAck) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Simple field (deviceConfigurationAckDataBlock)
	lengthInBits += m.deviceConfigurationAckDataBlock.LengthInBits()

	return lengthInBits
}

func (m DeviceConfigurationAck) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func DeviceConfigurationAckParse(io spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

	// Simple Field (deviceConfigurationAckDataBlock)
	_deviceConfigurationAckDataBlockMessage, _err := DeviceConfigurationAckDataBlockParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'deviceConfigurationAckDataBlock'. " + _err.Error())
	}
	var deviceConfigurationAckDataBlock DeviceConfigurationAckDataBlock
	deviceConfigurationAckDataBlock, _deviceConfigurationAckDataBlockOk := _deviceConfigurationAckDataBlockMessage.(DeviceConfigurationAckDataBlock)
	if !_deviceConfigurationAckDataBlockOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_deviceConfigurationAckDataBlockMessage).Name() + " to DeviceConfigurationAckDataBlock")
	}

	// Create the instance
	return NewDeviceConfigurationAck(deviceConfigurationAckDataBlock), nil
}

func (m DeviceConfigurationAck) Serialize(io spi.WriteBuffer) {

	// Simple Field (deviceConfigurationAckDataBlock)
	deviceConfigurationAckDataBlock := DeviceConfigurationAckDataBlock(m.deviceConfigurationAckDataBlock)
	deviceConfigurationAckDataBlock.Serialize(io)
}
