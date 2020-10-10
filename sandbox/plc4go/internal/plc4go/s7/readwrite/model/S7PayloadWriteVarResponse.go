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
type S7PayloadWriteVarResponse struct {
	items []S7VarPayloadStatusItem
	S7Payload
}

// The corresponding interface
type IS7PayloadWriteVarResponse interface {
	IS7Payload
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m S7PayloadWriteVarResponse) ParameterParameterType() uint8 {
	return 0x05
}

func (m S7PayloadWriteVarResponse) MessageType() uint8 {
	return 0x03
}

func (m S7PayloadWriteVarResponse) initialize() spi.Message {
	return m
}

func NewS7PayloadWriteVarResponse(items []S7VarPayloadStatusItem) S7PayloadInitializer {
	return &S7PayloadWriteVarResponse{items: items}
}

func (m S7PayloadWriteVarResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7Payload.LengthInBits()

	// Array field
	if len(m.items) > 0 {
		for _, element := range m.items {
			lengthInBits += element.LengthInBits()
		}
	}

	return lengthInBits
}

func (m S7PayloadWriteVarResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7PayloadWriteVarResponseParse(io spi.ReadBuffer, parameter S7Parameter) (S7PayloadInitializer, error) {

	// Array field (items)
	var items []S7VarPayloadStatusItem
	// Count array
	{
		items := make([]S7VarPayloadStatusItem, S7ParameterWriteVarResponse(parameter).numItems)
		for curItem := uint16(0); curItem < uint16(S7ParameterWriteVarResponse(parameter).numItems); curItem++ {

			_message, _err := S7VarPayloadStatusItemParse(io)
			if _err != nil {
				return nil, errors.New("Error parsing 'items' field " + _err.Error())
			}
			var _item S7VarPayloadStatusItem
			_item, _ok := _message.(S7VarPayloadStatusItem)
			if !_ok {
				return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to S7VarPayloadStatusItem")
			}
			items = append(items, _item)
		}
	}

	// Create the instance
	return NewS7PayloadWriteVarResponse(items), nil
}

func (m S7PayloadWriteVarResponse) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IS7PayloadWriteVarResponse); ok {

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
