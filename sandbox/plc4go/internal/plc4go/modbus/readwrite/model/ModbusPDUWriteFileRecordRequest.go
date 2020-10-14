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
type ModbusPDUWriteFileRecordRequest struct {
	Items []IModbusPDUWriteFileRecordRequestItem
	ModbusPDU
}

// The corresponding interface
type IModbusPDUWriteFileRecordRequest interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ModbusPDUWriteFileRecordRequest) ErrorFlag() bool {
	return false
}

func (m ModbusPDUWriteFileRecordRequest) FunctionFlag() uint8 {
	return 0x15
}

func (m ModbusPDUWriteFileRecordRequest) Response() bool {
	return false
}

func (m ModbusPDUWriteFileRecordRequest) initialize() spi.Message {
	return m
}

func NewModbusPDUWriteFileRecordRequest(items []IModbusPDUWriteFileRecordRequestItem) ModbusPDUInitializer {
	return &ModbusPDUWriteFileRecordRequest{Items: items}
}

func CastIModbusPDUWriteFileRecordRequest(structType interface{}) IModbusPDUWriteFileRecordRequest {
	castFunc := func(typ interface{}) IModbusPDUWriteFileRecordRequest {
		if iModbusPDUWriteFileRecordRequest, ok := typ.(IModbusPDUWriteFileRecordRequest); ok {
			return iModbusPDUWriteFileRecordRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastModbusPDUWriteFileRecordRequest(structType interface{}) ModbusPDUWriteFileRecordRequest {
	castFunc := func(typ interface{}) ModbusPDUWriteFileRecordRequest {
		if sModbusPDUWriteFileRecordRequest, ok := typ.(ModbusPDUWriteFileRecordRequest); ok {
			return sModbusPDUWriteFileRecordRequest
		}
		return ModbusPDUWriteFileRecordRequest{}
	}
	return castFunc(structType)
}

func (m ModbusPDUWriteFileRecordRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.ModbusPDU.LengthInBits()

	// Implicit Field (byteCount)
	lengthInBits += 8

	// Array field
	if len(m.Items) > 0 {
		for _, element := range m.Items {
			lengthInBits += element.LengthInBits()
		}
	}

	return lengthInBits
}

func (m ModbusPDUWriteFileRecordRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUWriteFileRecordRequestParse(io *spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	byteCount, _byteCountErr := io.ReadUint8(8)
	if _byteCountErr != nil {
		return nil, errors.New("Error parsing 'byteCount' field " + _byteCountErr.Error())
	}

	// Array field (items)
	// Length array
	items := make([]IModbusPDUWriteFileRecordRequestItem, 0)
	_itemsLength := byteCount
	_itemsEndPos := io.GetPos() + uint16(_itemsLength)
	for io.GetPos() < _itemsEndPos {
		_message, _err := ModbusPDUWriteFileRecordRequestItemParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'items' field " + _err.Error())
		}
		var _item IModbusPDUWriteFileRecordRequestItem
		_item, _ok := _message.(IModbusPDUWriteFileRecordRequestItem)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ModbusPDUWriteFileRecordRequestItem")
		}
		items = append(items, _item)
	}

	// Create the instance
	return NewModbusPDUWriteFileRecordRequest(items), nil
}

func (m ModbusPDUWriteFileRecordRequest) Serialize(io spi.WriteBuffer) error {
	itemsArraySizeInBytes := func(items []IModbusPDUWriteFileRecordRequestItem) uint32 {
		var sizeInBytes uint32 = 0
		for _, v := range items {
			sizeInBytes += uint32(v.LengthInBytes())
		}
		return sizeInBytes
	}
	ser := func() error {

		// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
		byteCount := uint8(uint8(itemsArraySizeInBytes(m.Items)))
		_byteCountErr := io.WriteUint8(8, (byteCount))
		if _byteCountErr != nil {
			return errors.New("Error serializing 'byteCount' field " + _byteCountErr.Error())
		}

		// Array Field (items)
		if m.Items != nil {
			for _, _element := range m.Items {
				_elementErr := _element.Serialize(io)
				if _elementErr != nil {
					return errors.New("Error serializing 'items' field " + _elementErr.Error())
				}
			}
		}

		return nil
	}
	return ModbusPDUSerialize(io, m.ModbusPDU, CastIModbusPDU(m), ser)
}
