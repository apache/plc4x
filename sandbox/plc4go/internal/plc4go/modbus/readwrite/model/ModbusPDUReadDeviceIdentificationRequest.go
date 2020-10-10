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
type ModbusPDUReadDeviceIdentificationRequest struct {
	ModbusPDU
}

// The corresponding interface
type IModbusPDUReadDeviceIdentificationRequest interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUReadDeviceIdentificationRequest) ErrorFlag() bool {
	return false
}

func (m ModbusPDUReadDeviceIdentificationRequest) FunctionFlag() uint8 {
	return 0x2B
}

func (m ModbusPDUReadDeviceIdentificationRequest) Response() bool {
	return false
}

func (m ModbusPDUReadDeviceIdentificationRequest) initialize() spi.Message {
	return m
}

func NewModbusPDUReadDeviceIdentificationRequest() ModbusPDUInitializer {
	return &ModbusPDUReadDeviceIdentificationRequest{}
}

func (m ModbusPDUReadDeviceIdentificationRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	return lengthInBits
}

func (m ModbusPDUReadDeviceIdentificationRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUReadDeviceIdentificationRequestParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Create the instance
	return NewModbusPDUReadDeviceIdentificationRequest(), nil
}

func (m ModbusPDUReadDeviceIdentificationRequest) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IModbusPDUReadDeviceIdentificationRequest); ok {
		}
	}
	serializeFunc(m)
}
