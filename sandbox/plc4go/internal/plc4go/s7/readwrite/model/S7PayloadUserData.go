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
type S7PayloadUserData struct {
	items []IS7PayloadUserDataItem
	S7Payload
}

// The corresponding interface
type IS7PayloadUserData interface {
	IS7Payload
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m S7PayloadUserData) ParameterParameterType() uint8 {
	return 0x00
}

func (m S7PayloadUserData) MessageType() uint8 {
	return 0x07
}

func (m S7PayloadUserData) initialize() spi.Message {
	return m
}

func NewS7PayloadUserData(items []IS7PayloadUserDataItem) S7PayloadInitializer {
	return &S7PayloadUserData{items: items}
}

func CastIS7PayloadUserData(structType interface{}) IS7PayloadUserData {
	castFunc := func(typ interface{}) IS7PayloadUserData {
		if iS7PayloadUserData, ok := typ.(IS7PayloadUserData); ok {
			return iS7PayloadUserData
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7PayloadUserData(structType interface{}) S7PayloadUserData {
	castFunc := func(typ interface{}) S7PayloadUserData {
		if sS7PayloadUserData, ok := typ.(S7PayloadUserData); ok {
			return sS7PayloadUserData
		}
		return S7PayloadUserData{}
	}
	return castFunc(structType)
}

func (m S7PayloadUserData) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7Payload.LengthInBits()

	// Array field
	if len(m.items) > 0 {
		for _, element := range m.items {
			lengthInBits += element.LengthInBits()
		}
	}

	return lengthInBits
}

func (m S7PayloadUserData) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7PayloadUserDataParse(io spi.ReadBuffer, parameter IS7Parameter) (S7PayloadInitializer, error) {

	// Array field (items)
	var items []IS7PayloadUserDataItem
	// Count array
	{
		items := make([]IS7PayloadUserDataItem, uint16(len(CastS7ParameterUserData(parameter).items)))
		for curItem := uint16(0); curItem < uint16(uint16(len(CastS7ParameterUserData(parameter).items))); curItem++ {

			_message, _err := S7PayloadUserDataItemParse(io, CastS7ParameterUserDataItemCPUFunctions(CastS7ParameterUserData(parameter).items).cpuFunctionType)
			if _err != nil {
				return nil, errors.New("Error parsing 'items' field " + _err.Error())
			}
			var _item IS7PayloadUserDataItem
			_item, _ok := _message.(IS7PayloadUserDataItem)
			if !_ok {
				return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to S7PayloadUserDataItem")
			}
			items = append(items, _item)
		}
	}

	// Create the instance
	return NewS7PayloadUserData(items), nil
}

func (m S7PayloadUserData) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Array Field (items)
		if m.items != nil {
			for _, _element := range m.items {
				_element.Serialize(io)
			}
		}

	}
	S7PayloadSerialize(io, m.S7Payload, CastIS7Payload(m), ser)
}
