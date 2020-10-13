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
type ConnectionResponseDataBlockDeviceManagement struct {
	ConnectionResponseDataBlock
}

// The corresponding interface
type IConnectionResponseDataBlockDeviceManagement interface {
	IConnectionResponseDataBlock
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ConnectionResponseDataBlockDeviceManagement) ConnectionType() uint8 {
	return 0x03
}

func (m ConnectionResponseDataBlockDeviceManagement) initialize() spi.Message {
	return m
}

func NewConnectionResponseDataBlockDeviceManagement() ConnectionResponseDataBlockInitializer {
	return &ConnectionResponseDataBlockDeviceManagement{}
}

func CastIConnectionResponseDataBlockDeviceManagement(structType interface{}) IConnectionResponseDataBlockDeviceManagement {
	castFunc := func(typ interface{}) IConnectionResponseDataBlockDeviceManagement {
		if iConnectionResponseDataBlockDeviceManagement, ok := typ.(IConnectionResponseDataBlockDeviceManagement); ok {
			return iConnectionResponseDataBlockDeviceManagement
		}
		return nil
	}
	return castFunc(structType)
}

func CastConnectionResponseDataBlockDeviceManagement(structType interface{}) ConnectionResponseDataBlockDeviceManagement {
	castFunc := func(typ interface{}) ConnectionResponseDataBlockDeviceManagement {
		if sConnectionResponseDataBlockDeviceManagement, ok := typ.(ConnectionResponseDataBlockDeviceManagement); ok {
			return sConnectionResponseDataBlockDeviceManagement
		}
		return ConnectionResponseDataBlockDeviceManagement{}
	}
	return castFunc(structType)
}

func (m ConnectionResponseDataBlockDeviceManagement) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ConnectionResponseDataBlock.LengthInBits()

	return lengthInBits
}

func (m ConnectionResponseDataBlockDeviceManagement) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ConnectionResponseDataBlockDeviceManagementParse(io *spi.ReadBuffer) (ConnectionResponseDataBlockInitializer, error) {

	// Create the instance
	return NewConnectionResponseDataBlockDeviceManagement(), nil
}

func (m ConnectionResponseDataBlockDeviceManagement) Serialize(io spi.WriteBuffer) {
	ser := func() {

	}
	ConnectionResponseDataBlockSerialize(io, m.ConnectionResponseDataBlock, CastIConnectionResponseDataBlock(m), ser)
}
