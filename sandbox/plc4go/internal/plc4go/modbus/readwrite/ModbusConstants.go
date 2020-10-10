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
	"errors"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"strconv"
)

// Constant values.
const MODBUSTCPDEFAULTPORT uint16 = 502

// The data-structure of this message
type ModbusConstants struct {
}

// The corresponding interface
type IModbusConstants interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewModbusConstants() spi.Message {
	return &ModbusConstants{}
}

func (m ModbusConstants) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Const Field (modbusTcpDefaultPort)
	lengthInBits += 16

	return lengthInBits
}

func (m ModbusConstants) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusConstantsParse(io spi.ReadBuffer) (spi.Message, error) {

	// Const Field (modbusTcpDefaultPort)
	var modbusTcpDefaultPort uint16 = io.ReadUint16(16)
	if modbusTcpDefaultPort != MODBUSTCPDEFAULTPORT {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(MODBUSTCPDEFAULTPORT)) + " but got " + strconv.Itoa(int(modbusTcpDefaultPort)))
	}

	// Create the instance
	return NewModbusConstants(), nil
}

func (m ModbusConstants) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IModbusPDU); ok {

			// Const Field (modbusTcpDefaultPort)
			io.WriteUint16(16, 502)
		}
	}
	serializeFunc(m)
}
