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
type ModbusPDUReadFifoQueueResponse struct {
	fifoValue []uint16
	ModbusPDU
}

// The corresponding interface
type IModbusPDUReadFifoQueueResponse interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUReadFifoQueueResponse) ErrorFlag() bool {
	return false
}

func (m ModbusPDUReadFifoQueueResponse) FunctionFlag() uint8 {
	return 0x18
}

func (m ModbusPDUReadFifoQueueResponse) Response() bool {
	return true
}

func (m ModbusPDUReadFifoQueueResponse) initialize() spi.Message {
	return spi.Message(m)
}

func NewModbusPDUReadFifoQueueResponse(fifoValue []uint16) ModbusPDUInitializer {
	return &ModbusPDUReadFifoQueueResponse{fifoValue: fifoValue}
}

func (m ModbusPDUReadFifoQueueResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Implicit Field (byteCount)
	lengthInBits += 16

	// Implicit Field (fifoCount)
	lengthInBits += 16

	// Array field
	if len(m.fifoValue) > 0 {
		lengthInBits += 16 * uint16(len(m.fifoValue))
	}

	return lengthInBits
}

func (m ModbusPDUReadFifoQueueResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUReadFifoQueueResponseParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var _ uint16 = io.ReadUint16(16)

	// Implicit Field (fifoCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var fifoCount uint16 = io.ReadUint16(16)

	// Array field (fifoValue)
	var fifoValue []uint16
	// Count array
	{
		fifoValue := make([]uint16, fifoCount)
		for curItem := uint16(0); curItem < uint16(fifoCount); curItem++ {

			fifoValue[curItem] = io.ReadUint16(16)
		}
	}

	// Create the instance
	return NewModbusPDUReadFifoQueueResponse(fifoValue), nil
}

func (m ModbusPDUReadFifoQueueResponse) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IModbusPDU); ok {

			// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
			var byteCount uint16 = (((uint16(len(m.fifoValue))) * (2)) + (2))
			io.WriteUint16(16, (byteCount))

			// Implicit Field (fifoCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
			var fifoCount uint16 = (((uint16(len(m.fifoValue))) * (2)) / (2))
			io.WriteUint16(16, (fifoCount))

			// Array Field (fifoValue)
			if m.fifoValue != nil {
				for _, _element := range m.fifoValue {
					io.WriteUint16(16, _element)
				}
			}
		}
	}
	serializeFunc(m)
}
