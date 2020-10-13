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
	"math"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type S7VarPayloadDataItem struct {
	returnCode    IDataTransportErrorCode
	transportSize IDataTransportSize
	data          []int8
}

// The corresponding interface
type IS7VarPayloadDataItem interface {
	spi.Message
	Serialize(io spi.WriteBuffer, lastItem bool)
}

func NewS7VarPayloadDataItem(returnCode IDataTransportErrorCode, transportSize IDataTransportSize, data []int8) spi.Message {
	return &S7VarPayloadDataItem{returnCode: returnCode, transportSize: transportSize, data: data}
}

func CastIS7VarPayloadDataItem(structType interface{}) IS7VarPayloadDataItem {
	castFunc := func(typ interface{}) IS7VarPayloadDataItem {
		if iS7VarPayloadDataItem, ok := typ.(IS7VarPayloadDataItem); ok {
			return iS7VarPayloadDataItem
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7VarPayloadDataItem(structType interface{}) S7VarPayloadDataItem {
	castFunc := func(typ interface{}) S7VarPayloadDataItem {
		if sS7VarPayloadDataItem, ok := typ.(S7VarPayloadDataItem); ok {
			return sS7VarPayloadDataItem
		}
		return S7VarPayloadDataItem{}
	}
	return castFunc(structType)
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
	_timesPadding := uint8(spi.InlineIf(false, uint16(uint8(0)), uint16(uint8(uint8(len(m.data)))%uint8(uint8(2)))))
	for ; _timesPadding > 0; _timesPadding-- {
		lengthInBits += 8
	}

	return lengthInBits
}

func (m S7VarPayloadDataItem) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7VarPayloadDataItemParse(io *spi.ReadBuffer, lastItem bool) (spi.Message, error) {

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
	dataLength, _dataLengthErr := io.ReadUint16(16)
	if _dataLengthErr != nil {
		return nil, errors.New("Error parsing 'dataLength' field " + _dataLengthErr.Error())
	}

	// Array field (data)
	// Count array
	data := make([]int8, spi.InlineIf(transportSize.SizeInBits(), uint16(math.Ceil(float64(dataLength)/float64(float64(8.0)))), uint16(dataLength)))
	for curItem := uint16(0); curItem < uint16(spi.InlineIf(transportSize.SizeInBits(), uint16(math.Ceil(float64(dataLength)/float64(float64(8.0)))), uint16(dataLength))); curItem++ {

		_item, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'data' field " + _err.Error())
		}
		data[curItem] = _item
	}

	// Padding Field (padding)
	{
		_timesPadding := (spi.InlineIf(lastItem, uint16(uint8(0)), uint16(uint8(uint8(len(data)))%uint8(uint8(2)))))
		for ; (io.HasMore(8)) && (_timesPadding > 0); _timesPadding-- {
			// Just read the padding data and ignore it
			_, _err := io.ReadUint8(8)
			if _err != nil {
				return nil, errors.New("Error parsing 'padding' field " + _err.Error())
			}
		}
	}

	// Create the instance
	return NewS7VarPayloadDataItem(returnCode, transportSize, data), nil
}

func (m S7VarPayloadDataItem) Serialize(io spi.WriteBuffer, lastItem bool) {

	// Enum field (returnCode)
	returnCode := CastDataTransportErrorCode(m.returnCode)
	returnCode.Serialize(io)

	// Enum field (transportSize)
	transportSize := CastDataTransportSize(m.transportSize)
	transportSize.Serialize(io)

	// Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	dataLength := uint16(uint16(uint16(len(m.data))) * uint16(uint16(spi.InlineIf(bool(bool((m.transportSize) == (DataTransportSize_BIT))), uint16(uint16(1)), uint16(uint16(spi.InlineIf(transportSize.SizeInBits(), uint16(uint16(8)), uint16(uint16(1)))))))))
	io.WriteUint16(16, (dataLength))

	// Array Field (data)
	if m.data != nil {
		for _, _element := range m.data {
			io.WriteInt8(8, _element)
		}
	}

	// Padding Field (padding)
	{
		_timesPadding := uint8(spi.InlineIf(lastItem, uint16(uint8(0)), uint16(uint8(uint8(len(m.data)))%uint8(uint8(2)))))
		for ; _timesPadding > 0; _timesPadding-- {
			_paddingValue := uint8(uint8(0))
			io.WriteUint8(8, (_paddingValue))
		}
	}

}
