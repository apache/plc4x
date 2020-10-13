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
type S7ParameterWriteVarResponse struct {
	numItems uint8
	S7Parameter
}

// The corresponding interface
type IS7ParameterWriteVarResponse interface {
	IS7Parameter
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m S7ParameterWriteVarResponse) ParameterType() uint8 {
	return 0x05
}

func (m S7ParameterWriteVarResponse) MessageType() uint8 {
	return 0x03
}

func (m S7ParameterWriteVarResponse) initialize() spi.Message {
	return m
}

func NewS7ParameterWriteVarResponse(numItems uint8) S7ParameterInitializer {
	return &S7ParameterWriteVarResponse{numItems: numItems}
}

func CastIS7ParameterWriteVarResponse(structType interface{}) IS7ParameterWriteVarResponse {
	castFunc := func(typ interface{}) IS7ParameterWriteVarResponse {
		if iS7ParameterWriteVarResponse, ok := typ.(IS7ParameterWriteVarResponse); ok {
			return iS7ParameterWriteVarResponse
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7ParameterWriteVarResponse(structType interface{}) S7ParameterWriteVarResponse {
	castFunc := func(typ interface{}) S7ParameterWriteVarResponse {
		if sS7ParameterWriteVarResponse, ok := typ.(S7ParameterWriteVarResponse); ok {
			return sS7ParameterWriteVarResponse
		}
		return S7ParameterWriteVarResponse{}
	}
	return castFunc(structType)
}

func (m S7ParameterWriteVarResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7Parameter.LengthInBits()

	// Simple field (numItems)
	lengthInBits += 8

	return lengthInBits
}

func (m S7ParameterWriteVarResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7ParameterWriteVarResponseParse(io *spi.ReadBuffer) (S7ParameterInitializer, error) {

	// Simple Field (numItems)
	numItems, _numItemsErr := io.ReadUint8(8)
	if _numItemsErr != nil {
		return nil, errors.New("Error parsing 'numItems' field " + _numItemsErr.Error())
	}

	// Create the instance
	return NewS7ParameterWriteVarResponse(numItems), nil
}

func (m S7ParameterWriteVarResponse) Serialize(io spi.WriteBuffer) {
	ser := func() {

		// Simple Field (numItems)
		numItems := uint8(m.numItems)
		io.WriteUint8(8, (numItems))

	}
	S7ParameterSerialize(io, m.S7Parameter, CastIS7Parameter(m), ser)
}
