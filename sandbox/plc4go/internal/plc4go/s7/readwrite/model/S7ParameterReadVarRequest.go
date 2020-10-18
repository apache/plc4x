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
type S7ParameterReadVarRequest struct {
	Items []IS7VarRequestParameterItem
	S7Parameter
}

// The corresponding interface
type IS7ParameterReadVarRequest interface {
	IS7Parameter
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m S7ParameterReadVarRequest) ParameterType() uint8 {
	return 0x04
}

func (m S7ParameterReadVarRequest) MessageType() uint8 {
	return 0x01
}

func (m S7ParameterReadVarRequest) initialize() spi.Message {
	return m
}

func NewS7ParameterReadVarRequest(items []IS7VarRequestParameterItem) S7ParameterInitializer {
	return &S7ParameterReadVarRequest{Items: items}
}

func CastIS7ParameterReadVarRequest(structType interface{}) IS7ParameterReadVarRequest {
	castFunc := func(typ interface{}) IS7ParameterReadVarRequest {
		if iS7ParameterReadVarRequest, ok := typ.(IS7ParameterReadVarRequest); ok {
			return iS7ParameterReadVarRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7ParameterReadVarRequest(structType interface{}) S7ParameterReadVarRequest {
	castFunc := func(typ interface{}) S7ParameterReadVarRequest {
		if sS7ParameterReadVarRequest, ok := typ.(S7ParameterReadVarRequest); ok {
			return sS7ParameterReadVarRequest
		}
		return S7ParameterReadVarRequest{}
	}
	return castFunc(structType)
}

func (m S7ParameterReadVarRequest) LengthInBits() uint16 {
	var lengthInBits = m.S7Parameter.LengthInBits()

	// Implicit Field (numItems)
	lengthInBits += 8

	// Array field
	if len(m.Items) > 0 {
		for _, element := range m.Items {
			lengthInBits += element.LengthInBits()
		}
	}

	return lengthInBits
}

func (m S7ParameterReadVarRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7ParameterReadVarRequestParse(io *spi.ReadBuffer) (S7ParameterInitializer, error) {

	// Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	numItems, _numItemsErr := io.ReadUint8(8)
	if _numItemsErr != nil {
		return nil, errors.New("Error parsing 'numItems' field " + _numItemsErr.Error())
	}

	// Array field (items)
	// Count array
	items := make([]IS7VarRequestParameterItem, numItems)
	for curItem := uint16(0); curItem < uint16(numItems); curItem++ {

		_message, _err := S7VarRequestParameterItemParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'items' field " + _err.Error())
		}
		var _item IS7VarRequestParameterItem
		_item, _ok := _message.(IS7VarRequestParameterItem)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to S7VarRequestParameterItem")
		}
		items[curItem] = _item
	}

	// Create the instance
	return NewS7ParameterReadVarRequest(items), nil
}

func (m S7ParameterReadVarRequest) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Implicit Field (numItems) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
		numItems := uint8(uint8(len(m.Items)))
		_numItemsErr := io.WriteUint8(8, numItems)
		if _numItemsErr != nil {
			return errors.New("Error serializing 'numItems' field " + _numItemsErr.Error())
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
	return S7ParameterSerialize(io, m.S7Parameter, CastIS7Parameter(m), ser)
}
