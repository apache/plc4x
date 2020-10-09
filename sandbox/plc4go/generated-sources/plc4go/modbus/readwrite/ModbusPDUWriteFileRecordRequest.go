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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/src/plc4go/spi"
	"reflect"
)

type ModbusPDUWriteFileRecordRequest struct {
	items []ModbusPDUWriteFileRecordRequestItem
	ModbusPDU
}

type IModbusPDUWriteFileRecordRequest interface {
	IModbusPDU
	Serialize(io spi.WriteBuffer)
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
	return spi.Message(m)
}

func NewModbusPDUWriteFileRecordRequest(items []ModbusPDUWriteFileRecordRequestItem) ModbusPDUInitializer {
	return &ModbusPDUWriteFileRecordRequest{items: items}
}

func (m ModbusPDUWriteFileRecordRequest) LengthInBits() uint16 {
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

func (m ModbusPDUWriteFileRecordRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ModbusPDUWriteFileRecordRequestParse(io spi.ReadBuffer) (ModbusPDUInitializer, error) {

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var byteCount uint8 = io.ReadUint8(8)

	// Array field (items)
	var items []ModbusPDUWriteFileRecordRequestItem
	// Length array
	_itemsLength := uint16(byteCount)
	_itemsEndPos := io.GetPos() + _itemsLength
	for io.GetPos() < _itemsEndPos {
		_message, _err := ModbusPDUWriteFileRecordRequestItemParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'items' field " + _err.Error())
		}
		var _item ModbusPDUWriteFileRecordRequestItem
		_item, _ok := _message.(ModbusPDUWriteFileRecordRequestItem)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to ModbusPDUWriteFileRecordRequestItem")
		}
		items = append(items, _item)
	}

	// Create the instance
	return NewModbusPDUWriteFileRecordRequest(items), nil
}

func (m ModbusPDUWriteFileRecordRequest) Serialize(io spi.WriteBuffer) {
	itemsArraySizeInBytes := func(items []ModbusPDUWriteFileRecordRequestItem) uint32 {
		var sizeInBytes uint32 = 0
		for _, v := range items {
			sizeInBytes += uint32(v.LengthInBytes())
		}
		return sizeInBytes
	}

	// Implicit Field (byteCount) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var byteCount uint8 = (uint8(itemsArraySizeInBytes(m.items)))
	io.WriteUint8(8, (byteCount))

	// Array Field (items)
	if m.items != nil {
		for _, _element := range m.items {
			_element.Serialize(io)
		}
	}
}
