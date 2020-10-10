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
type S7ParameterWriteVarRequest struct {
	items []S7VarRequestParameterItem
	S7Parameter
}

// The corresponding interface
type IS7ParameterWriteVarRequest interface {
	IS7Parameter
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m S7ParameterWriteVarRequest) ParameterType() uint8 {
	return 0x05
}

func (m S7ParameterWriteVarRequest) MessageType() uint8 {
	return 0x01
}

func (m S7ParameterWriteVarRequest) initialize() spi.Message {
	return m
}

func NewS7ParameterWriteVarRequest(items []S7VarRequestParameterItem) S7ParameterInitializer {
	return &S7ParameterWriteVarRequest{items: items}
}

func (m S7ParameterWriteVarRequest) LengthInBits() uint16 {
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

func (m S7ParameterWriteVarRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7ParameterWriteVarRequestParse(io spi.ReadBuffer) (S7ParameterInitializer, error) {

	// Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var numItems uint8 = io.ReadUint8(8)

	// Array field (items)
	var items []S7VarRequestParameterItem
	// Count array
	{
		items := make([]S7VarRequestParameterItem, numItems)
		for curItem := uint16(0); curItem < uint16(numItems); curItem++ {

			_message, _err := S7VarRequestParameterItemParse(io)
			if _err != nil {
				return nil, errors.New("Error parsing 'items' field " + _err.Error())
			}
			var _item S7VarRequestParameterItem
			_item, _ok := _message.(S7VarRequestParameterItem)
			if !_ok {
				return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to S7VarRequestParameterItem")
			}
			items = append(items, _item)
		}
	}

	// Create the instance
	return NewS7ParameterWriteVarRequest(items), nil
}

func (m S7ParameterWriteVarRequest) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IS7ParameterWriteVarRequest); ok {

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
	}
	serializeFunc(m)
}
