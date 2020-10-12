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
	"reflect"
)

// The data-structure of this message
type ModbusPDUReadFileRecordResponse struct {
	items []ModbusPDUReadFileRecordResponseItem
	ModbusPDU
}

// The corresponding interface
type IModbusPDUReadFileRecordResponse interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ModbusPDUReadFileRecordResponse) ErrorFlag() bool {
	return false
}

func (m ModbusPDUReadFileRecordResponse) FunctionFlag() uint8 {
	return 0x14
}

func (m ModbusPDUReadFileRecordResponse) Response() bool {
	return true
}

func (m ModbusPDUReadFileRecordResponse) initialize() spi.Message {
	return m
}

func NewModbusPDUReadFileRecordResponse(items []ModbusPDUReadFileRecordResponseItem) ModbusPDUInitializer {
	return &ModbusPDUReadFileRecordResponse{items: items}
}

func CastIModbusPDUReadFileRecordResponse(structType interface{}) IModbusPDUReadFileRecordResponse {
	castFunc := func(typ interface{}) IModbusPDUReadFileRecordResponse {
		if iModbusPDUReadFileRecordResponse, ok := typ.(IModbusPDUReadFileRecordResponse); ok {
			return iModbusPDUReadFileRecordResponse
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDUReadFileRecordResponse(structType interface{}) ModbusPDUReadFileRecordResponse {
	castFunc := func(typ interface{}) ModbusPDUReadFileRecordResponse {
		if sModbusPDUReadFileRecordResponse, ok := typ.(ModbusPDUReadFileRecordResponse); ok {
			return sModbusPDUReadFileRecordResponse
		}
		return ModbusPDUReadFileRecordResponse{}
	}
	return castFunc(structType)
}

func (m ModbusPDUReadFileRecordResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Implicit Field (byteCount)
	lengthInBits += 8

	// Array field
	if len(m.items) > 0 {
		for _, element := range m.items {
			lengthInBits += element.LengthInBits()
		}
	}

	return lengthInBits
}

func (m ModbusPDUReadFileRecordResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUReadFileRecordResponseParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var byteCount uint8 = io.ReadUint8(8)

	// Array field (items)
	var items []ModbusPDUReadFileRecordResponseItem
	// Length array
	_itemsLength := byteCount
	_itemsEndPos := io.GetPos() + uint16(_itemsLength)
	for io.GetPos() < _itemsEndPos {
		_message, _err := ModbusPDUReadFileRecordResponseItemParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'items' field " + _err.Error())
		}
		var _item ModbusPDUReadFileRecordResponseItem
		_item, _ok := _message.(ModbusPDUReadFileRecordResponseItem)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ModbusPDUReadFileRecordResponseItem")
		}
		items = append(items, _item)
	}

	// Create the instance
	return NewModbusPDUReadFileRecordResponse(items), nil
}

func (m ModbusPDUReadFileRecordResponse) Serialize(io spi.WriteBuffer) {
	itemsArraySizeInBytes := func(items []ModbusPDUReadFileRecordResponseItem) uint32 {
		var sizeInBytes uint32 = 0
		for _, v := range items {
			sizeInBytes += uint32(v.LengthInBytes())
		}
		return sizeInBytes
	}

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	byteCount := uint8(uint8(itemsArraySizeInBytes(m.items)))
	io.WriteUint8(8, (byteCount))

	// Array Field (items)
	if m.items != nil {
		for _, _element := range m.items {
			_element.Serialize(io)
		}
	}
}
