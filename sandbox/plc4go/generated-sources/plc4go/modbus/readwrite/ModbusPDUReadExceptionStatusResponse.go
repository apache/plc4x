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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/src/plc4go/spi"
)

type ModbusPDUReadExceptionStatusResponse struct {
	value uint8
	ModbusPDU
}

type IModbusPDUReadExceptionStatusResponse interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUReadExceptionStatusResponse) ErrorFlag() bool {
	return false
}

func (m ModbusPDUReadExceptionStatusResponse) FunctionFlag() uint8 {
	return 0x07
}

func (m ModbusPDUReadExceptionStatusResponse) Response() bool {
	return true
}

func (m ModbusPDUReadExceptionStatusResponse) initialize() spi.Message {
	return spi.Message(m)
}

func NewModbusPDUReadExceptionStatusResponse(value uint8) ModbusPDUInitializer {
	return &ModbusPDUReadExceptionStatusResponse{value: value}
}

func (m ModbusPDUReadExceptionStatusResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Simple field (value)
	lengthInBits += 8

	return lengthInBits
}

func (m ModbusPDUReadExceptionStatusResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUReadExceptionStatusResponseParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Simple Field (value)
	var value uint8 = io.ReadUint8(8)

	// Create the instance
	return NewModbusPDUReadExceptionStatusResponse(value), nil
}

func (m ModbusPDUReadExceptionStatusResponse) Serialize(io spi.WriteBuffer) {

	// Simple Field (value)
	var value uint8 = m.value
	io.WriteUint8(8, (value))
}
