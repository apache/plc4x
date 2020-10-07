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
	"math"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/src/plc4go/spi"
)

type ModbusPDUWriteSingleRegisterResponse struct {
	address uint16
	value   uint16
	ModbusPDU
}

func (m ModbusPDUWriteSingleRegisterResponse) initialize() ModbusPDU {
	return m.ModbusPDU
}

func NewModbusPDUWriteSingleRegisterResponse(address uint16, value uint16) ModbusPDUInitializer {
	return &ModbusPDUWriteSingleRegisterResponse{address: address, value: value}
}

func (m ModbusPDUWriteSingleRegisterResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Simple field (address)
	lengthInBits += 16

	// Simple field (value)
	lengthInBits += 16

	return lengthInBits
}

func (m ModbusPDUWriteSingleRegisterResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUWriteSingleRegisterResponseParse(io spi.ReadBuffer) ModbusPDUInitializer {
	var startPos = io.GetPos()
	var curPos uint16

	// Simple Field (address)
	var address uint16 = io.ReadUint16(16)

	// Simple Field (value)
	var value uint16 = io.ReadUint16(16)

	// Create the instance
	return NewModbusPDUWriteSingleRegisterResponse(address, value)
}
