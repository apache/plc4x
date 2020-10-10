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
type ModbusPDUReportServerIdRequest struct {
	ModbusPDU
}

// The corresponding interface
type IModbusPDUReportServerIdRequest interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUReportServerIdRequest) ErrorFlag() bool {
	return false
}

func (m ModbusPDUReportServerIdRequest) FunctionFlag() uint8 {
	return 0x11
}

func (m ModbusPDUReportServerIdRequest) Response() bool {
	return false
}

func (m ModbusPDUReportServerIdRequest) initialize() spi.Message {
	return spi.Message(m)
}

func NewModbusPDUReportServerIdRequest() ModbusPDUInitializer {
	return &ModbusPDUReportServerIdRequest{}
}

func (m ModbusPDUReportServerIdRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	return lengthInBits
}

func (m ModbusPDUReportServerIdRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUReportServerIdRequestParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Create the instance
	return NewModbusPDUReportServerIdRequest(), nil
}

func (m ModbusPDUReportServerIdRequest) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IModbusPDU); ok {
		}
	}
	serializeFunc(m)
}
