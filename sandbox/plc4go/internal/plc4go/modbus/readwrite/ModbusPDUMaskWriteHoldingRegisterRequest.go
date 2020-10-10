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
type ModbusPDUMaskWriteHoldingRegisterRequest struct {
	referenceAddress uint16
	andMask          uint16
	orMask           uint16
	ModbusPDU
}

// The corresponding interface
type IModbusPDUMaskWriteHoldingRegisterRequest interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUMaskWriteHoldingRegisterRequest) ErrorFlag() bool {
	return false
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) FunctionFlag() uint8 {
	return 0x16
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) Response() bool {
	return false
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) initialize() spi.Message {
	return spi.Message(m)
}

func NewModbusPDUMaskWriteHoldingRegisterRequest(referenceAddress uint16, andMask uint16, orMask uint16) ModbusPDUInitializer {
	return &ModbusPDUMaskWriteHoldingRegisterRequest{referenceAddress: referenceAddress, andMask: andMask, orMask: orMask}
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Simple field (referenceAddress)
	lengthInBits += 16

	// Simple field (andMask)
	lengthInBits += 16

	// Simple field (orMask)
	lengthInBits += 16

	return lengthInBits
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUMaskWriteHoldingRegisterRequestParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Simple Field (referenceAddress)
	var referenceAddress uint16 = io.ReadUint16(16)

	// Simple Field (andMask)
	var andMask uint16 = io.ReadUint16(16)

	// Simple Field (orMask)
	var orMask uint16 = io.ReadUint16(16)

	// Create the instance
	return NewModbusPDUMaskWriteHoldingRegisterRequest(referenceAddress, andMask, orMask), nil
}

func (m ModbusPDUMaskWriteHoldingRegisterRequest) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IModbusPDU); ok {

			// Simple Field (referenceAddress)
			var referenceAddress uint16 = m.referenceAddress
			io.WriteUint16(16, (referenceAddress))

			// Simple Field (andMask)
			var andMask uint16 = m.andMask
			io.WriteUint16(16, (andMask))

			// Simple Field (orMask)
			var orMask uint16 = m.orMask
			io.WriteUint16(16, (orMask))
		}
	}
	serializeFunc(m)
}
