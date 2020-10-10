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
type ModbusPDUReadFileRecordResponseItem struct {
	referenceType uint8
	data          []int8
}

// The corresponding interface
type IModbusPDUReadFileRecordResponseItem interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewModbusPDUReadFileRecordResponseItem(referenceType uint8, data []int8) spi.Message {
	return &ModbusPDUReadFileRecordResponseItem{referenceType: referenceType, data: data}
}

func (m ModbusPDUReadFileRecordResponseItem) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Implicit Field (dataLength)
	lengthInBits += 8

	// Simple field (referenceType)
	lengthInBits += 8

	// Array field
	if len(m.data) > 0 {
		lengthInBits += 8 * uint16(len(m.data))
	}

	return lengthInBits
}

func (m ModbusPDUReadFileRecordResponseItem) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUReadFileRecordResponseItemParse(io spi.ReadBuffer) (spi.Message, error) {

	// Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var dataLength uint8 = io.ReadUint8(8)

	// Simple Field (referenceType)
	var referenceType uint8 = io.ReadUint8(8)

	// Array field (data)
	var data []int8
	// Length array
	_dataLength := uint16((dataLength) - (1))
	_dataEndPos := io.GetPos() + _dataLength
	for io.GetPos() < _dataEndPos {
		data = append(data, io.ReadInt8(8))
	}

	// Create the instance
	return NewModbusPDUReadFileRecordResponseItem(referenceType, data), nil
}

func (m ModbusPDUReadFileRecordResponseItem) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IModbusPDU); ok {

			// Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
			var dataLength uint8 = ((uint8(len(m.data))) + (1))
			io.WriteUint8(8, (dataLength))

			// Simple Field (referenceType)
			var referenceType uint8 = m.referenceType
			io.WriteUint8(8, (referenceType))

			// Array Field (data)
			if m.data != nil {
				for _, _element := range m.data {
					io.WriteInt8(8, _element)
				}
			}
		}
	}
	serializeFunc(m)
}
