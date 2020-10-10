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
type ModbusPDUGetComEventLogResponse struct {
	status       uint16
	eventCount   uint16
	messageCount uint16
	events       []int8
	ModbusPDU
}

// The corresponding interface
type IModbusPDUGetComEventLogResponse interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUGetComEventLogResponse) ErrorFlag() bool {
	return false
}

func (m ModbusPDUGetComEventLogResponse) FunctionFlag() uint8 {
	return 0x0C
}

func (m ModbusPDUGetComEventLogResponse) Response() bool {
	return true
}

func (m ModbusPDUGetComEventLogResponse) initialize() spi.Message {
	return m
}

func NewModbusPDUGetComEventLogResponse(status uint16, eventCount uint16, messageCount uint16, events []int8) ModbusPDUInitializer {
	return &ModbusPDUGetComEventLogResponse{status: status, eventCount: eventCount, messageCount: messageCount, events: events}
}

func (m ModbusPDUGetComEventLogResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Implicit Field (byteCount)
	lengthInBits += 8

	// Simple field (status)
	lengthInBits += 16

	// Simple field (eventCount)
	lengthInBits += 16

	// Simple field (messageCount)
	lengthInBits += 16

	// Array field
	if len(m.events) > 0 {
		lengthInBits += 8 * uint16(len(m.events))
	}

	return lengthInBits
}

func (m ModbusPDUGetComEventLogResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUGetComEventLogResponseParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var byteCount uint8 = io.ReadUint8(8)

	// Simple Field (status)
	var status uint16 = io.ReadUint16(16)

	// Simple Field (eventCount)
	var eventCount uint16 = io.ReadUint16(16)

	// Simple Field (messageCount)
	var messageCount uint16 = io.ReadUint16(16)

	// Array field (events)
	var events []int8
	// Count array
	{
		events := make([]int8, (byteCount)-(6))
		for curItem := uint16(0); curItem < uint16((byteCount)-(6)); curItem++ {

			events = append(events, io.ReadInt8(8))
		}
	}

	// Create the instance
	return NewModbusPDUGetComEventLogResponse(status, eventCount, messageCount, events), nil
}

func (m ModbusPDUGetComEventLogResponse) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IModbusPDUGetComEventLogResponse); ok {

			// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
			byteCount := uint8((uint8(len(m.events))) + (6))
			io.WriteUint8(8, (byteCount))

			// Simple Field (status)
			var status uint16 = m.status
			io.WriteUint16(16, (status))

			// Simple Field (eventCount)
			var eventCount uint16 = m.eventCount
			io.WriteUint16(16, (eventCount))

			// Simple Field (messageCount)
			var messageCount uint16 = m.messageCount
			io.WriteUint16(16, (messageCount))

			// Array Field (events)
			if m.events != nil {
				for _, _element := range m.events {
					io.WriteInt8(8, _element)
				}
			}
		}
	}
	serializeFunc(m)
}
