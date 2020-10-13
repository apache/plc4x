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
type S7MessageUserData struct {
	S7Message
}

// The corresponding interface
type IS7MessageUserData interface {
	IS7Message
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m S7MessageUserData) MessageType() uint8 {
	return 0x07
}

func (m S7MessageUserData) initialize(tpduReference uint16, parameter *IS7Parameter, payload *IS7Payload) spi.Message {
	m.tpduReference = tpduReference
	m.parameter = parameter
	m.payload = payload
	return m
}

func NewS7MessageUserData() S7MessageInitializer {
	return &S7MessageUserData{}
}

func CastIS7MessageUserData(structType interface{}) IS7MessageUserData {
	castFunc := func(typ interface{}) IS7MessageUserData {
		if iS7MessageUserData, ok := typ.(IS7MessageUserData); ok {
			return iS7MessageUserData
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7MessageUserData(structType interface{}) S7MessageUserData {
	castFunc := func(typ interface{}) S7MessageUserData {
		if sS7MessageUserData, ok := typ.(S7MessageUserData); ok {
			return sS7MessageUserData
		}
		return S7MessageUserData{}
	}
	return castFunc(structType)
}

func (m S7MessageUserData) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7Message.LengthInBits()

	return lengthInBits
}

func (m S7MessageUserData) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7MessageUserDataParse(io *spi.ReadBuffer) (S7MessageInitializer, error) {

	// Create the instance
	return NewS7MessageUserData(), nil
}

func (m S7MessageUserData) Serialize(io spi.WriteBuffer) {
	ser := func() {

	}
	S7MessageSerialize(io, m.S7Message, CastIS7Message(m), ser)
}
