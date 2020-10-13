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
	"errors"
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
	Serialize(io spi.WriteBuffer) error
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
	return m
}

func NewModbusPDUReadFifoQueueResponse(fifoValue []uint16) ModbusPDUInitializer {
	return &ModbusPDUReadFifoQueueResponse{fifoValue: fifoValue}
}

func CastIModbusPDUReadFifoQueueResponse(structType interface{}) IModbusPDUReadFifoQueueResponse {
	castFunc := func(typ interface{}) IModbusPDUReadFifoQueueResponse {
		if iModbusPDUReadFifoQueueResponse, ok := typ.(IModbusPDUReadFifoQueueResponse); ok {
			return iModbusPDUReadFifoQueueResponse
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDUReadFifoQueueResponse(structType interface{}) ModbusPDUReadFifoQueueResponse {
	castFunc := func(typ interface{}) ModbusPDUReadFifoQueueResponse {
		if sModbusPDUReadFifoQueueResponse, ok := typ.(ModbusPDUReadFifoQueueResponse); ok {
			return sModbusPDUReadFifoQueueResponse
		}
		return ModbusPDUReadFifoQueueResponse{}
	}
	return castFunc(structType)
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

func ModbusPDUReadFifoQueueResponseParse(io *spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	_, _byteCountErr := io.ReadUint16(16)
	if _byteCountErr != nil {
		return nil, errors.New("Error parsing 'byteCount' field " + _byteCountErr.Error())
	}

	// Implicit Field (fifoCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	fifoCount, _fifoCountErr := io.ReadUint16(16)
	if _fifoCountErr != nil {
		return nil, errors.New("Error parsing 'fifoCount' field " + _fifoCountErr.Error())
	}

	// Array field (fifoValue)
	// Count array
	fifoValue := make([]uint16, fifoCount)
	for curItem := uint16(0); curItem < uint16(fifoCount); curItem++ {

		_item, _err := io.ReadUint16(16)
		if _err != nil {
			return nil, errors.New("Error parsing 'fifoValue' field " + _err.Error())
		}
		fifoValue[curItem] = _item
	}

	// Create the instance
	return NewModbusPDUReadFifoQueueResponse(fifoValue), nil
}

func (m ModbusPDUReadFifoQueueResponse) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
		byteCount := uint16(uint16(uint16(uint16(uint16(len(m.fifoValue)))*uint16(uint16(2)))) + uint16(uint16(2)))
		_byteCountErr := io.WriteUint16(16, (byteCount))
		if _byteCountErr != nil {
			return errors.New("Error serializing 'byteCount' field " + _byteCountErr.Error())
		}

		// Implicit Field (fifoCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
		fifoCount := uint16(uint16(uint16(uint16(uint16(len(m.fifoValue)))*uint16(uint16(2)))) / uint16(uint16(2)))
		_fifoCountErr := io.WriteUint16(16, (fifoCount))
		if _fifoCountErr != nil {
			return errors.New("Error serializing 'fifoCount' field " + _fifoCountErr.Error())
		}

		// Array Field (fifoValue)
		if m.fifoValue != nil {
			for _, _element := range m.fifoValue {
				_elementErr := io.WriteUint16(16, _element)
				if _elementErr != nil {
					return errors.New("Error serializing 'fifoValue' field " + _elementErr.Error())
				}
			}
		}

		return nil
	}
	return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
