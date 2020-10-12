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
type KnxNetObjectServer struct {
	version uint8
	ServiceId
}

// The corresponding interface
type IKnxNetObjectServer interface {
	IServiceId
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m KnxNetObjectServer) ServiceType() uint8 {
	return 0x08
}

func (m KnxNetObjectServer) initialize() spi.Message {
	return m
}

func NewKnxNetObjectServer(version uint8) ServiceIdInitializer {
	return &KnxNetObjectServer{version: version}
}

func CastIKnxNetObjectServer(structType interface{}) IKnxNetObjectServer {
	castFunc := func(typ interface{}) IKnxNetObjectServer {
		if iKnxNetObjectServer, ok := typ.(IKnxNetObjectServer); ok {
			return iKnxNetObjectServer
		}
		return nil
	}
	return castFunc(structType)
}

func CastKnxNetObjectServer(structType interface{}) KnxNetObjectServer {
	castFunc := func(typ interface{}) KnxNetObjectServer {
		if sKnxNetObjectServer, ok := typ.(KnxNetObjectServer); ok {
			return sKnxNetObjectServer
		}
		return KnxNetObjectServer{}
	}
	return castFunc(structType)
}

func (m KnxNetObjectServer) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ServiceId.LengthInBits()

	// Simple field (version)
	lengthInBits += 8

	return lengthInBits
}

func (m KnxNetObjectServer) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func KnxNetObjectServerParse(io spi.ReadBuffer) (ServiceIdInitializer, error) {

	// Simple Field (version)
	var version uint8 = io.ReadUint8(8)

	// Create the instance
	return NewKnxNetObjectServer(version), nil
}

func (m KnxNetObjectServer) Serialize(io spi.WriteBuffer) {

	// Simple Field (version)
	version := uint8(m.version)
	io.WriteUint8(8, (version))
}
