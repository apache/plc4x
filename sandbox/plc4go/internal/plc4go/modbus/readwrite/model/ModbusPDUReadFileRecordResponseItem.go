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
type ModbusPDUReadFileRecordResponseItem struct {
	referenceType uint8
	data          []int8
}

// The corresponding interface
type IModbusPDUReadFileRecordResponseItem interface {
	spi.Message
	Serialize(io spi.WriteBuffer) error
}

func NewModbusPDUReadFileRecordResponseItem(referenceType uint8, data []int8) spi.Message {
	return &ModbusPDUReadFileRecordResponseItem{referenceType: referenceType, data: data}
}

func CastIModbusPDUReadFileRecordResponseItem(structType interface{}) IModbusPDUReadFileRecordResponseItem {
	castFunc := func(typ interface{}) IModbusPDUReadFileRecordResponseItem {
		if iModbusPDUReadFileRecordResponseItem, ok := typ.(IModbusPDUReadFileRecordResponseItem); ok {
			return iModbusPDUReadFileRecordResponseItem
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDUReadFileRecordResponseItem(structType interface{}) ModbusPDUReadFileRecordResponseItem {
	castFunc := func(typ interface{}) ModbusPDUReadFileRecordResponseItem {
		if sModbusPDUReadFileRecordResponseItem, ok := typ.(ModbusPDUReadFileRecordResponseItem); ok {
			return sModbusPDUReadFileRecordResponseItem
		}
		return ModbusPDUReadFileRecordResponseItem{}
	}
	return castFunc(structType)
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

func ModbusPDUReadFileRecordResponseItemParse(io *spi.ReadBuffer) (spi.Message, error) {

	// Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	dataLength, _dataLengthErr := io.ReadUint8(8)
	if _dataLengthErr != nil {
		return nil, errors.New("Error parsing 'dataLength' field " + _dataLengthErr.Error())
	}

	// Simple Field (referenceType)
	referenceType, _referenceTypeErr := io.ReadUint8(8)
	if _referenceTypeErr != nil {
		return nil, errors.New("Error parsing 'referenceType' field " + _referenceTypeErr.Error())
	}

	// Array field (data)
	// Length array
	data := make([]int8, 0)
	_dataLength := uint16(dataLength) - uint16(uint16(1))
	_dataEndPos := io.GetPos() + uint16(_dataLength)
	for io.GetPos() < _dataEndPos {
		_item, _err := io.ReadInt8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'data' field " + _err.Error())
		}
		data = append(data, _item)
	}

	// Create the instance
	return NewModbusPDUReadFileRecordResponseItem(referenceType, data), nil
}

func (m ModbusPDUReadFileRecordResponseItem) Serialize(io spi.WriteBuffer) error {

	// Implicit Field (dataLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	dataLength := uint8(uint8(uint8(len(m.data))) + uint8(uint8(1)))
	_dataLengthErr := io.WriteUint8(8, (dataLength))
	if _dataLengthErr != nil {
		return errors.New("Error serializing 'dataLength' field " + _dataLengthErr.Error())
	}

	// Simple Field (referenceType)
	referenceType := uint8(m.referenceType)
	_referenceTypeErr := io.WriteUint8(8, (referenceType))
	if _referenceTypeErr != nil {
		return errors.New("Error serializing 'referenceType' field " + _referenceTypeErr.Error())
	}

	// Array Field (data)
	if m.data != nil {
		for _, _element := range m.data {
			_elementErr := io.WriteInt8(8, _element)
			if _elementErr != nil {
				return errors.New("Error serializing 'data' field " + _elementErr.Error())
			}
		}
	}

	return nil
}
