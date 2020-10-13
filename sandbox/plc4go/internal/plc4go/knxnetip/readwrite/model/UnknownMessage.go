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
type UnknownMessage struct {
	unknownData []int8
	KNXNetIPMessage
}

// The corresponding interface
type IUnknownMessage interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m UnknownMessage) MsgType() uint16 {
	return 0x020B
}

func (m UnknownMessage) initialize() spi.Message {
	return m
}

func NewUnknownMessage(unknownData []int8) KNXNetIPMessageInitializer {
	return &UnknownMessage{unknownData: unknownData}
}

func CastIUnknownMessage(structType interface{}) IUnknownMessage {
	castFunc := func(typ interface{}) IUnknownMessage {
		if iUnknownMessage, ok := typ.(IUnknownMessage); ok {
			return iUnknownMessage
		}
		return nil
	}
	return castFunc(structType)
}

func CastUnknownMessage(structType interface{}) UnknownMessage {
	castFunc := func(typ interface{}) UnknownMessage {
		if sUnknownMessage, ok := typ.(UnknownMessage); ok {
			return sUnknownMessage
		}
		return UnknownMessage{}
	}
	return castFunc(structType)
}

func (m UnknownMessage) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Array field
	if len(m.unknownData) > 0 {
		lengthInBits += 8 * uint16(len(m.unknownData))
	}

	return lengthInBits
}

func (m UnknownMessage) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func UnknownMessageParse(io spi.ReadBuffer, totalLength uint16) (KNXNetIPMessageInitializer, error) {

	// Array field (unknownData)
	var unknownData []int8
	// Count array
	{
		unknownData := make([]int8, uint16(totalLength)-uint16(uint16(6)))
		for curItem := uint16(0); curItem < uint16(uint16(totalLength)-uint16(uint16(6))); curItem++ {

			_unknownDataVal, _err := io.ReadInt8(8)
			if _err != nil {
				return nil, errors.New("Error parsing 'unknownData' field " + _err.Error())
			}
			unknownData = append(unknownData, _unknownDataVal)
		}
	}

	// Create the instance
	return NewUnknownMessage(unknownData), nil
}

func (m UnknownMessage) Serialize(io spi.WriteBuffer) {

	// Array Field (unknownData)
	if m.unknownData != nil {
		for _, _element := range m.unknownData {
			io.WriteInt8(8, _element)
		}
	}
}
