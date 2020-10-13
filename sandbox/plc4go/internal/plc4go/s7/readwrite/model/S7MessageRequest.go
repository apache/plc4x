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
type S7MessageRequest struct {
	S7Message
}

// The corresponding interface
type IS7MessageRequest interface {
	IS7Message
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m S7MessageRequest) MessageType() uint8 {
	return 0x01
}

func (m S7MessageRequest) initialize(tpduReference uint16, parameter *IS7Parameter, payload *IS7Payload) spi.Message {
	m.tpduReference = tpduReference
	m.parameter = parameter
	m.payload = payload
	return m
}

func NewS7MessageRequest() S7MessageInitializer {
	return &S7MessageRequest{}
}

func CastIS7MessageRequest(structType interface{}) IS7MessageRequest {
	castFunc := func(typ interface{}) IS7MessageRequest {
		if iS7MessageRequest, ok := typ.(IS7MessageRequest); ok {
			return iS7MessageRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7MessageRequest(structType interface{}) S7MessageRequest {
	castFunc := func(typ interface{}) S7MessageRequest {
		if sS7MessageRequest, ok := typ.(S7MessageRequest); ok {
			return sS7MessageRequest
		}
		return S7MessageRequest{}
	}
	return castFunc(structType)
}

func (m S7MessageRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7Message.LengthInBits()

	return lengthInBits
}

func (m S7MessageRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7MessageRequestParse(io *spi.ReadBuffer) (S7MessageInitializer, error) {

	// Create the instance
	return NewS7MessageRequest(), nil
}

func (m S7MessageRequest) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		return nil
	}
	return S7MessageSerialize(io, m.S7Message, CastIS7Message(m), ser)
}
