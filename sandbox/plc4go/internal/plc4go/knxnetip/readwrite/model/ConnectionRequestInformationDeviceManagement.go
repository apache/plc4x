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
type ConnectionRequestInformationDeviceManagement struct {
	ConnectionRequestInformation
}

// The corresponding interface
type IConnectionRequestInformationDeviceManagement interface {
	IConnectionRequestInformation
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ConnectionRequestInformationDeviceManagement) ConnectionType() uint8 {
	return 0x03
}

func (m ConnectionRequestInformationDeviceManagement) initialize() spi.Message {
	return m
}

func NewConnectionRequestInformationDeviceManagement() ConnectionRequestInformationInitializer {
	return &ConnectionRequestInformationDeviceManagement{}
}

func CastIConnectionRequestInformationDeviceManagement(structType interface{}) IConnectionRequestInformationDeviceManagement {
	castFunc := func(typ interface{}) IConnectionRequestInformationDeviceManagement {
		if iConnectionRequestInformationDeviceManagement, ok := typ.(IConnectionRequestInformationDeviceManagement); ok {
			return iConnectionRequestInformationDeviceManagement
		}
		return nil
	}
	return castFunc(structType)
}

func CastConnectionRequestInformationDeviceManagement(structType interface{}) ConnectionRequestInformationDeviceManagement {
	castFunc := func(typ interface{}) ConnectionRequestInformationDeviceManagement {
		if sConnectionRequestInformationDeviceManagement, ok := typ.(ConnectionRequestInformationDeviceManagement); ok {
			return sConnectionRequestInformationDeviceManagement
		}
		return ConnectionRequestInformationDeviceManagement{}
	}
	return castFunc(structType)
}

func (m ConnectionRequestInformationDeviceManagement) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ConnectionRequestInformation.LengthInBits()

	return lengthInBits
}

func (m ConnectionRequestInformationDeviceManagement) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ConnectionRequestInformationDeviceManagementParse(io spi.ReadBuffer) (ConnectionRequestInformationInitializer, error) {

	// Create the instance
	return NewConnectionRequestInformationDeviceManagement(), nil
}

func (m ConnectionRequestInformationDeviceManagement) Serialize(io spi.WriteBuffer) {

}
