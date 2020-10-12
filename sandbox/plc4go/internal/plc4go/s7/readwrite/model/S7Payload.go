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
type S7Payload struct {
}

// The corresponding interface
type IS7Payload interface {
	spi.Message
	MessageType() uint8
	ParameterParameterType() uint8
	Serialize(io spi.WriteBuffer)
}

type S7PayloadInitializer interface {
	initialize() spi.Message
}

func S7PayloadMessageType(m IS7Payload) uint8 {
	return m.MessageType()
}

func S7PayloadParameterParameterType(m IS7Payload) uint8 {
	return m.ParameterParameterType()
}

func CastIS7Payload(structType interface{}) IS7Payload {
	castFunc := func(typ interface{}) IS7Payload {
		if iS7Payload, ok := typ.(IS7Payload); ok {
			return iS7Payload
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7Payload(structType interface{}) S7Payload {
	castFunc := func(typ interface{}) S7Payload {
		if sS7Payload, ok := typ.(S7Payload); ok {
			return sS7Payload
		}
		return S7Payload{}
	}
	return castFunc(structType)
}

func (m S7Payload) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m S7Payload) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7PayloadParse(io spi.ReadBuffer, messageType uint8, parameter S7Parameter) (spi.Message, error) {

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer S7PayloadInitializer
	var typeSwitchError error
	switch {
	case CastIS7Parameter(parameter).ParameterType() == 0x04 && messageType == 0x03:
		initializer, typeSwitchError = S7PayloadReadVarResponseParse(io, parameter)
	case CastIS7Parameter(parameter).ParameterType() == 0x05 && messageType == 0x01:
		initializer, typeSwitchError = S7PayloadWriteVarRequestParse(io, parameter)
	case CastIS7Parameter(parameter).ParameterType() == 0x05 && messageType == 0x03:
		initializer, typeSwitchError = S7PayloadWriteVarResponseParse(io, parameter)
	case CastIS7Parameter(parameter).ParameterType() == 0x00 && messageType == 0x07:
		initializer, typeSwitchError = S7PayloadUserDataParse(io, parameter)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func (m S7Payload) Serialize(io spi.WriteBuffer) {
	iS7Payload := CastIS7Payload(m)

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	iS7Payload.Serialize(io)
}
