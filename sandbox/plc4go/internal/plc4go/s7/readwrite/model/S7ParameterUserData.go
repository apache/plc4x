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
type S7ParameterUserData struct {
	items []IS7ParameterUserDataItem
	S7Parameter
}

// The corresponding interface
type IS7ParameterUserData interface {
	IS7Parameter
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m S7ParameterUserData) ParameterType() uint8 {
	return 0x00
}

func (m S7ParameterUserData) MessageType() uint8 {
	return 0x07
}

func (m S7ParameterUserData) initialize() spi.Message {
	return m
}

func NewS7ParameterUserData(items []IS7ParameterUserDataItem) S7ParameterInitializer {
	return &S7ParameterUserData{items: items}
}

func CastIS7ParameterUserData(structType interface{}) IS7ParameterUserData {
	castFunc := func(typ interface{}) IS7ParameterUserData {
		if iS7ParameterUserData, ok := typ.(IS7ParameterUserData); ok {
			return iS7ParameterUserData
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7ParameterUserData(structType interface{}) S7ParameterUserData {
	castFunc := func(typ interface{}) S7ParameterUserData {
		if sS7ParameterUserData, ok := typ.(S7ParameterUserData); ok {
			return sS7ParameterUserData
		}
		return S7ParameterUserData{}
	}
	return castFunc(structType)
}

func (m S7ParameterUserData) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7Parameter.LengthInBits()

	// Implicit Field (numItems)
	lengthInBits += 8

	// Array field
	if len(m.items) > 0 {
		for _, element := range m.items {
			lengthInBits += element.LengthInBits()
		}
	}

	return lengthInBits
}

func (m S7ParameterUserData) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7ParameterUserDataParse(io spi.ReadBuffer) (S7ParameterInitializer, error) {

	// Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	numItems, _numItemsErr := io.ReadUint8(8)
	if _numItemsErr != nil {
		return nil, errors.New("Error parsing 'numItems' field " + _numItemsErr.Error())
	}

	// Array field (items)
	var items []IS7ParameterUserDataItem
	// Count array
	{
		items := make([]IS7ParameterUserDataItem, numItems)
		for curItem := uint16(0); curItem < uint16(numItems); curItem++ {

			_message, _err := S7ParameterUserDataItemParse(io)
			if _err != nil {
				return nil, errors.New("Error parsing 'items' field " + _err.Error())
			}
			var _item IS7ParameterUserDataItem
			_item, _ok := _message.(IS7ParameterUserDataItem)
			if !_ok {
				return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to S7ParameterUserDataItem")
			}
			items = append(items, _item)
		}
	}

	// Create the instance
	return NewS7ParameterUserData(items), nil
}

func (m S7ParameterUserData) Serialize(io spi.WriteBuffer) {

	// Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	numItems := uint8(uint8(len(m.items)))
	io.WriteUint8(8, (numItems))

	// Array Field (items)
	if m.items != nil {
		for _, _element := range m.items {
			_element.Serialize(io)
		}
	}
}
