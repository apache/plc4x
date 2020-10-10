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
type S7VarPayloadDataItem struct {
	returnCode    DataTransportErrorCode
	transportSize DataTransportSize
	data          []int8
}

// The corresponding interface
type IS7VarPayloadDataItem interface {
	spi.Message
	Serialize(io spi.WriteBuffer)
}

func NewS7VarPayloadDataItem(returnCode DataTransportErrorCode, transportSize DataTransportSize, data []int8) spi.Message {
	return &S7VarPayloadDataItem{returnCode: returnCode, transportSize: transportSize, data: data}
}

func (m S7VarPayloadDataItem) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Enum Field (returnCode)
	lengthInBits += 8

	// Enum Field (transportSize)
	lengthInBits += 8

	// Implicit Field (dataLength)
	lengthInBits += 16

	// Array field
	if len(m.data) > 0 {
		lengthInBits += 8 * uint16(len(m.data))
	}

	// Padding Field (padding)
	_timesPadding := uint8(spi.InlineIf((false), uint16(0), uint16((uint8(len(COUNT)))%(2))))
	for ; _timesPadding > 0; _timesPadding-- {
		lengthInBits += 8
	}

	return lengthInBits
}

func (m S7VarPayloadDataItem) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7VarPayloadDataItemParse(io spi.ReadBuffer, lastItem bool) (spi.Message, error) {

	// Enum field (returnCode)
	returnCode, _returnCodeErr := DataTransportErrorCodeParse(io)
	if _returnCodeErr != nil {
		return nil, errors.New("Error parsing 'returnCode' field " + _returnCodeErr.Error())
	}

	// Enum field (transportSize)
	transportSize, _transportSizeErr := DataTransportSizeParse(io)
	if _transportSizeErr != nil {
		return nil, errors.New("Error parsing 'transportSize' field " + _transportSizeErr.Error())
	}

	// Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var dataLength uint16 = io.ReadUint16(16)

	// Array field (data)
	var data []int8
	// Count array
	{
		data := make([]int8, spi.InlineIf((transportSize.sizeInBits), uint16(CEIL((dataLength)/(8.0))), uint16(dataLength)))
		for curItem := uint16(0); curItem < uint16(spi.InlineIf((transportSize.sizeInBits), uint16(CEIL((dataLength)/(8.0))), uint16(dataLength))); curItem++ {

			data = append(data, io.ReadInt8(8))
		}
	}

	// Padding Field (padding)
	{
		_timesPadding := (spi.InlineIf((lastItem), uint16(0), uint16((uint8(len(COUNT)))%(2))))
		for ; (io.HasMore(8)) && (_timesPadding > 0); _timesPadding-- {
			// Just read the padding data and ignore it
			io.ReadUint8(8)
		}
	}

	// Create the instance
	return NewS7VarPayloadDataItem(returnCode, transportSize, data), nil
}

func (m S7VarPayloadDataItem) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IS7VarPayloadDataItem); ok {

			// Enum field (returnCode)
			returnCode := m.returnCode
			returnCode.Serialize(io)

			// Enum field (transportSize)
			transportSize := m.transportSize
			transportSize.Serialize(io)

			// Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
			dataLength := uint16((uint16(len(m.data))) * (spi.InlineIf(((m.transportSize) == (DataTransportSize_BIT)), uint16(1), uint16((spi.InlineIf((m.transportSize.sizeInBits), uint16(8), uint16(1)))))))
			io.WriteUint16(16, (dataLength))

			// Array Field (data)
			if m.data != nil {
				for _, _element := range m.data {
					io.WriteInt8(8, _element)
				}
			}

			// Padding Field (padding)
			{
				_timesPadding := uint8(spi.InlineIf((lastItem), uint16(0), uint16((uint8(len(m.data)))%(2))))
				for ; _timesPadding > 0; _timesPadding-- {
					_paddingValue := uint8(0)
					io.WriteUint8(8, (_paddingValue))
				}
			}
		}
	}
	serializeFunc(m)
}
