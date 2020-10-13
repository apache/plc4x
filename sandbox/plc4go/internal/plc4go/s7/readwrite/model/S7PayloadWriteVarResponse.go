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
	items []IS7VarPayloadStatusItem
	S7Payload
}

// The corresponding interface
type IS7PayloadWriteVarResponse interface {
	IS7Payload
	Serialize(io spi.WriteBuffer) error
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

func NewS7PayloadWriteVarResponse(items []IS7VarPayloadStatusItem) S7PayloadInitializer {
	return &S7PayloadWriteVarResponse{items: items}
}

func CastIS7PayloadWriteVarResponse(structType interface{}) IS7PayloadWriteVarResponse {
	castFunc := func(typ interface{}) IS7PayloadWriteVarResponse {
		if iS7PayloadWriteVarResponse, ok := typ.(IS7PayloadWriteVarResponse); ok {
			return iS7PayloadWriteVarResponse
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7PayloadWriteVarResponse(structType interface{}) S7PayloadWriteVarResponse {
	castFunc := func(typ interface{}) S7PayloadWriteVarResponse {
		if sS7PayloadWriteVarResponse, ok := typ.(S7PayloadWriteVarResponse); ok {
			return sS7PayloadWriteVarResponse
		}
		return S7PayloadWriteVarResponse{}
	}
	return castFunc(structType)
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

func S7PayloadWriteVarResponseParse(io *spi.ReadBuffer, parameter IS7Parameter) (S7PayloadInitializer, error) {

	// Array field (items)
	// Count array
	items := make([]IS7VarPayloadStatusItem, CastS7ParameterWriteVarResponse(parameter).numItems)
	for curItem := uint16(0); curItem < uint16(CastS7ParameterWriteVarResponse(parameter).numItems); curItem++ {

		_message, _err := S7VarPayloadStatusItemParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'items' field " + _err.Error())
		}
		var _item IS7VarPayloadStatusItem
		_item, _ok := _message.(IS7VarPayloadStatusItem)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to S7VarPayloadStatusItem")
		}
		items[curItem] = _item
	}

	// Create the instance
	return NewS7PayloadWriteVarResponse(items), nil
}

func (m S7PayloadWriteVarResponse) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Array Field (items)
		if m.items != nil {
			for _, _element := range m.items {
				_elementErr := _element.Serialize(io)
				if _elementErr != nil {
					return errors.New("Error serializing 'items' field " + _elementErr.Error())
				}
			}
		}

		return nil
	}
	return S7PayloadSerialize(io, m.S7Payload, CastIS7Payload(m), ser)
}
