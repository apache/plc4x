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
type S7Parameter struct {
}

// The corresponding interface
type IS7Parameter interface {
	spi.Message
	MessageType() uint8
	ParameterType() uint8
	Serialize(io spi.WriteBuffer) error
}

type S7ParameterInitializer interface {
	initialize() spi.Message
}

func S7ParameterMessageType(m IS7Parameter) uint8 {
	return m.MessageType()
}

func S7ParameterParameterType(m IS7Parameter) uint8 {
	return m.ParameterType()
}

func CastIS7Parameter(structType interface{}) IS7Parameter {
	castFunc := func(typ interface{}) IS7Parameter {
		if iS7Parameter, ok := typ.(IS7Parameter); ok {
			return iS7Parameter
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7Parameter(structType interface{}) S7Parameter {
	castFunc := func(typ interface{}) S7Parameter {
		if sS7Parameter, ok := typ.(S7Parameter); ok {
			return sS7Parameter
		}
		return S7Parameter{}
	}
	return castFunc(structType)
}

func (m S7Parameter) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (parameterType)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m S7Parameter) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7ParameterParse(io *spi.ReadBuffer, messageType uint8) (spi.Message, error) {

	// Discriminator Field (parameterType) (Used as input to a switch field)
	parameterType, _parameterTypeErr := io.ReadUint8(8)
	if _parameterTypeErr != nil {
		return nil, errors.New("Error parsing 'parameterType' field " + _parameterTypeErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer S7ParameterInitializer
	var typeSwitchError error
	switch {
	case parameterType == 0xF0:
		initializer, typeSwitchError = S7ParameterSetupCommunicationParse(io)
	case parameterType == 0x04 && messageType == 0x01:
		initializer, typeSwitchError = S7ParameterReadVarRequestParse(io)
	case parameterType == 0x04 && messageType == 0x03:
		initializer, typeSwitchError = S7ParameterReadVarResponseParse(io)
	case parameterType == 0x05 && messageType == 0x01:
		initializer, typeSwitchError = S7ParameterWriteVarRequestParse(io)
	case parameterType == 0x05 && messageType == 0x03:
		initializer, typeSwitchError = S7ParameterWriteVarResponseParse(io)
	case parameterType == 0x00 && messageType == 0x07:
		initializer, typeSwitchError = S7ParameterUserDataParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(), nil
}

func S7ParameterSerialize(io spi.WriteBuffer, m S7Parameter, i IS7Parameter, childSerialize func() error) error {

	// Discriminator Field (parameterType) (Used as input to a switch field)
	parameterType := uint8(i.ParameterType())
	_parameterTypeErr := io.WriteUint8(8, parameterType)
	if _parameterTypeErr != nil {
		return errors.New("Error serializing 'parameterType' field " + _parameterTypeErr.Error())
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	_typeSwitchErr := childSerialize()
	if _typeSwitchErr != nil {
		return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
	}

	return nil
}
