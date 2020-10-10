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
type APDUAbort struct {
	server           bool
	originalInvokeId uint8
	abortReason      uint8
	APDU
}

// The corresponding interface
type IAPDUAbort interface {
	IAPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m APDUAbort) ApduType() uint8 {
	return 0x7
}

func (m APDUAbort) initialize() spi.Message {
	return m
}

func NewAPDUAbort(server bool, originalInvokeId uint8, abortReason uint8) APDUInitializer {
	return &APDUAbort{server: server, originalInvokeId: originalInvokeId, abortReason: abortReason}
}

func (m APDUAbort) LengthInBits() uint16 {
	var lengthInBits uint16 = m.APDU.LengthInBits()

	// Reserved Field (reserved)
	lengthInBits += 3

	// Simple field (server)
	lengthInBits += 1

	// Simple field (originalInvokeId)
	lengthInBits += 8

	// Simple field (abortReason)
	lengthInBits += 8

	return lengthInBits
}

func (m APDUAbort) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func APDUAbortParse(io spi.ReadBuffer) (APDUInitializer, error) {

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		var reserved uint8 = io.ReadUint8(3)
		if reserved != uint8(0x00) {
			log.WithFields(log.Fields{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (server)
	var server bool = io.ReadBit()

	// Simple Field (originalInvokeId)
	var originalInvokeId uint8 = io.ReadUint8(8)

	// Simple Field (abortReason)
	var abortReason uint8 = io.ReadUint8(8)

	// Create the instance
	return NewAPDUAbort(server, originalInvokeId, abortReason), nil
}

func (m APDUAbort) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IAPDUAbort); ok {

			// Reserved Field (reserved)
			io.WriteUint8(3, uint8(0x00))

			// Simple Field (server)
			var server bool = m.server
			io.WriteBit((bool)(server))

			// Simple Field (originalInvokeId)
			var originalInvokeId uint8 = m.originalInvokeId
			io.WriteUint8(8, (originalInvokeId))

			// Simple Field (abortReason)
			var abortReason uint8 = m.abortReason
			io.WriteUint8(8, (abortReason))
		}
	}
	serializeFunc(m)
}
