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
type COTPPacketData struct {
	eot     bool
	tpduRef uint8
	COTPPacket
}

// The corresponding interface
type ICOTPPacketData interface {
	ICOTPPacket
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m COTPPacketData) TpduCode() uint8 {
	return 0xF0
}

func (m COTPPacketData) initialize(parameters []COTPParameter, payload *S7Message) spi.Message {
	m.parameters = parameters
	m.payload = payload
	return m
}

func NewCOTPPacketData(eot bool, tpduRef uint8) COTPPacketInitializer {
	return &COTPPacketData{eot: eot, tpduRef: tpduRef}
}

func (m COTPPacketData) LengthInBits() uint16 {
	var lengthInBits uint16 = m.COTPPacket.LengthInBits()

	// Simple field (eot)
	lengthInBits += 1

	// Simple field (tpduRef)
	lengthInBits += 7

	return lengthInBits
}

func (m COTPPacketData) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func COTPPacketDataParse(io spi.ReadBuffer) (COTPPacketInitializer, error) {

	// Simple Field (eot)
	var eot bool = io.ReadBit()

	// Simple Field (tpduRef)
	var tpduRef uint8 = io.ReadUint8(7)

	// Create the instance
	return NewCOTPPacketData(eot, tpduRef), nil
}

func (m COTPPacketData) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(ICOTPPacketData); ok {

			// Simple Field (eot)
			var eot bool = m.eot
			io.WriteBit((bool)(eot))

			// Simple Field (tpduRef)
			var tpduRef uint8 = m.tpduRef
			io.WriteUint8(7, (tpduRef))
		}
	}
	serializeFunc(m)
}
