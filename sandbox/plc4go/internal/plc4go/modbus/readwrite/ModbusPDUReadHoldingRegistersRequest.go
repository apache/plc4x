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
package readwrite

import (
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type ModbusPDUReadHoldingRegistersRequest struct {
	startingAddress uint16
	quantity        uint16
	ModbusPDU
}

// The corresponding interface
type IModbusPDUReadHoldingRegistersRequest interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUReadHoldingRegistersRequest) ErrorFlag() bool {
	return false
}

func (m ModbusPDUReadHoldingRegistersRequest) FunctionFlag() uint8 {
	return 0x03
}

func (m ModbusPDUReadHoldingRegistersRequest) Response() bool {
	return false
}

func (m ModbusPDUReadHoldingRegistersRequest) initialize() spi.Message {
	return spi.Message(m)
}

func NewModbusPDUReadHoldingRegistersRequest(startingAddress uint16, quantity uint16) ModbusPDUInitializer {
	return &ModbusPDUReadHoldingRegistersRequest{startingAddress: startingAddress, quantity: quantity}
}

func (m ModbusPDUReadHoldingRegistersRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Simple field (startingAddress)
	lengthInBits += 16

	// Simple field (quantity)
	lengthInBits += 16

	return lengthInBits
}

func (m ModbusPDUReadHoldingRegistersRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUReadHoldingRegistersRequestParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Simple Field (startingAddress)
	var startingAddress uint16 = io.ReadUint16(16)

	// Simple Field (quantity)
	var quantity uint16 = io.ReadUint16(16)

	// Create the instance
	return NewModbusPDUReadHoldingRegistersRequest(startingAddress, quantity), nil
}

func (m ModbusPDUReadHoldingRegistersRequest) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IModbusPDU); ok {

			// Simple Field (startingAddress)
			var startingAddress uint16 = m.startingAddress
			io.WriteUint16(16, (startingAddress))

			// Simple Field (quantity)
			var quantity uint16 = m.quantity
			io.WriteUint16(16, (quantity))
		}
	}
	serializeFunc(m)
}
