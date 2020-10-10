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
	log "github.com/sirupsen/logrus"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type DeviceStatus struct {
	programMode bool
}

// The corresponding interface
type IDeviceStatus interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewDeviceStatus(programMode bool) spi.Message {
	return &DeviceStatus{programMode: programMode}
}

func (m DeviceStatus) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Reserved Field (reserved)
	lengthInBits += 7

	// Simple field (programMode)
	lengthInBits += 1

	return lengthInBits
}

func (m DeviceStatus) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func DeviceStatusParse(io spi.ReadBuffer) (spi.Message, error) {

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		var reserved uint8 = io.ReadUint8(7)
		if reserved != uint8(0x00) {
			log.WithFields(log.Fields{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (programMode)
	var programMode bool = io.ReadBit()

	// Create the instance
	return NewDeviceStatus(programMode), nil
}

func (m DeviceStatus) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IDeviceStatus); ok {

			// Reserved Field (reserved)
			io.WriteUint8(7, uint8(0x00))

			// Simple Field (programMode)
			var programMode bool = m.programMode
			io.WriteBit((bool)(programMode))
		}
	}
	serializeFunc(m)
}
