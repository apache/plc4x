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
type KnxNetIpCore struct {
	version uint8
	ServiceId
}

// The corresponding interface
type IKnxNetIpCore interface {
	IServiceId
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m KnxNetIpCore) ServiceType() uint8 {
	return 0x02
}

func (m KnxNetIpCore) initialize() spi.Message {
	return m
}

func NewKnxNetIpCore(version uint8) ServiceIdInitializer {
	return &KnxNetIpCore{version: version}
}

func CastIKnxNetIpCore(structType interface{}) IKnxNetIpCore {
	castFunc := func(typ interface{}) IKnxNetIpCore {
		if iKnxNetIpCore, ok := typ.(IKnxNetIpCore); ok {
			return iKnxNetIpCore
		}
		return nil
	}
	return castFunc(structType)
}

func CastKnxNetIpCore(structType interface{}) KnxNetIpCore {
	castFunc := func(typ interface{}) KnxNetIpCore {
		if sKnxNetIpCore, ok := typ.(KnxNetIpCore); ok {
			return sKnxNetIpCore
		}
		return KnxNetIpCore{}
	}
	return castFunc(structType)
}

func (m KnxNetIpCore) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ServiceId.LengthInBits()

	// Simple field (version)
	lengthInBits += 8

	return lengthInBits
}

func (m KnxNetIpCore) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KnxNetIpCoreParse(io spi.ReadBuffer) (ServiceIdInitializer, error) {

	// Simple Field (version)
	var version uint8 = io.ReadUint8(8)

	// Create the instance
	return NewKnxNetIpCore(version), nil
}

func (m KnxNetIpCore) Serialize(io spi.WriteBuffer) {

	// Simple Field (version)
	version := uint8(m.version)
	io.WriteUint8(8, (version))
}
