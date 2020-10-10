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
type S7ParameterSetupCommunication struct {
	maxAmqCaller uint16
	maxAmqCallee uint16
	pduLength    uint16
	S7Parameter
}

// The corresponding interface
type IS7ParameterSetupCommunication interface {
	IS7Parameter
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m S7ParameterSetupCommunication) ParameterType() uint8 {
	return 0xF0
}

func (m S7ParameterSetupCommunication) MessageType() uint8 {
	return 0
}

func (m S7ParameterSetupCommunication) initialize() spi.Message {
	return m
}

func NewS7ParameterSetupCommunication(maxAmqCaller uint16, maxAmqCallee uint16, pduLength uint16) S7ParameterInitializer {
	return &S7ParameterSetupCommunication{maxAmqCaller: maxAmqCaller, maxAmqCallee: maxAmqCallee, pduLength: pduLength}
}

func (m S7ParameterSetupCommunication) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7Parameter.LengthInBits()

	// Reserved Field (reserved)
	lengthInBits += 8

	// Simple field (maxAmqCaller)
	lengthInBits += 16

	// Simple field (maxAmqCallee)
	lengthInBits += 16

	// Simple field (pduLength)
	lengthInBits += 16

	return lengthInBits
}

func (m S7ParameterSetupCommunication) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7ParameterSetupCommunicationParse(io spi.ReadBuffer) (S7ParameterInitializer, error) {

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		var reserved uint8 = io.ReadUint8(8)
		if reserved != uint8(0x00) {
			log.WithFields(log.Fields{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (maxAmqCaller)
	var maxAmqCaller uint16 = io.ReadUint16(16)

	// Simple Field (maxAmqCallee)
	var maxAmqCallee uint16 = io.ReadUint16(16)

	// Simple Field (pduLength)
	var pduLength uint16 = io.ReadUint16(16)

	// Create the instance
	return NewS7ParameterSetupCommunication(maxAmqCaller, maxAmqCallee, pduLength), nil
}

func (m S7ParameterSetupCommunication) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IS7ParameterSetupCommunication); ok {

			// Reserved Field (reserved)
			io.WriteUint8(8, uint8(0x00))

			// Simple Field (maxAmqCaller)
			var maxAmqCaller uint16 = m.maxAmqCaller
			io.WriteUint16(16, (maxAmqCaller))

			// Simple Field (maxAmqCallee)
			var maxAmqCallee uint16 = m.maxAmqCallee
			io.WriteUint16(16, (maxAmqCallee))

			// Simple Field (pduLength)
			var pduLength uint16 = m.pduLength
			io.WriteUint16(16, (pduLength))
		}
	}
	serializeFunc(m)
}
