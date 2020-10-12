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
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
)

// The data-structure of this message
type S7MessageResponse struct {
	errorClass uint8
	errorCode  uint8
	S7Message
}

// The corresponding interface
type IS7MessageResponse interface {
	IS7Message
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m S7MessageResponse) MessageType() uint8 {
	return 0x02
}

func (m S7MessageResponse) initialize(tpduReference uint16, parameter *S7Parameter, payload *S7Payload) spi.Message {
	m.tpduReference = tpduReference
	m.parameter = parameter
	m.payload = payload
	return m
}

func NewS7MessageResponse(errorClass uint8, errorCode uint8) S7MessageInitializer {
	return &S7MessageResponse{errorClass: errorClass, errorCode: errorCode}
}

func CastIS7MessageResponse(structType interface{}) IS7MessageResponse {
	castFunc := func(typ interface{}) IS7MessageResponse {
		if iS7MessageResponse, ok := typ.(IS7MessageResponse); ok {
			return iS7MessageResponse
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7MessageResponse(structType interface{}) S7MessageResponse {
	castFunc := func(typ interface{}) S7MessageResponse {
		if sS7MessageResponse, ok := typ.(S7MessageResponse); ok {
			return sS7MessageResponse
		}
		return S7MessageResponse{}
	}
	return castFunc(structType)
}

func (m S7MessageResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7Message.LengthInBits()

	// Simple field (errorClass)
	lengthInBits += 8

	// Simple field (errorCode)
	lengthInBits += 8

	return lengthInBits
}

func (m S7MessageResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7MessageResponseParse(io spi.ReadBuffer) (S7MessageInitializer, error) {

	// Simple Field (errorClass)
	var errorClass uint8 = io.ReadUint8(8)

	// Simple Field (errorCode)
	var errorCode uint8 = io.ReadUint8(8)

	// Create the instance
	return NewS7MessageResponse(errorClass, errorCode), nil
}

func (m S7MessageResponse) Serialize(io spi.WriteBuffer) {

	// Simple Field (errorClass)
	errorClass := uint8(m.errorClass)
	io.WriteUint8(8, (errorClass))

	// Simple Field (errorCode)
	errorCode := uint8(m.errorCode)
	io.WriteUint8(8, (errorCode))
}
