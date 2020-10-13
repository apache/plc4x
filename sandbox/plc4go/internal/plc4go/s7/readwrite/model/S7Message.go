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
	log "github.com/sirupsen/logrus"
	"plc4x.apache.org/plc4go-modbus-driver/0.8.0/internal/plc4go/spi"
	"reflect"
	"strconv"
)

// Constant values.
const S7Message_PROTOCOLID uint8 = 0x32

// The data-structure of this message
type S7Message struct {
	tpduReference uint16
	parameter     *IS7Parameter
	payload       *IS7Payload
}

// The corresponding interface
type IS7Message interface {
	spi.Message
	MessageType() uint8
	Serialize(io spi.WriteBuffer)
}

type S7MessageInitializer interface {
	initialize(tpduReference uint16, parameter *IS7Parameter, payload *IS7Payload) spi.Message
}

func S7MessageMessageType(m IS7Message) uint8 {
	return m.MessageType()
}

func CastIS7Message(structType interface{}) IS7Message {
	castFunc := func(typ interface{}) IS7Message {
		if iS7Message, ok := typ.(IS7Message); ok {
			return iS7Message
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7Message(structType interface{}) S7Message {
	castFunc := func(typ interface{}) S7Message {
		if sS7Message, ok := typ.(S7Message); ok {
			return sS7Message
		}
		return S7Message{}
	}
	return castFunc(structType)
}

func (m S7Message) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Const Field (protocolId)
	lengthInBits += 8

	// Discriminator Field (messageType)
	lengthInBits += 8

	// Reserved Field (reserved)
	lengthInBits += 16

	// Simple field (tpduReference)
	lengthInBits += 16

	// Implicit Field (parameterLength)
	lengthInBits += 16

	// Implicit Field (payloadLength)
	lengthInBits += 16

	// Length of sub-type elements will be added by sub-type...

	// Optional Field (parameter)
	if m.parameter != nil {
		lengthInBits += (*m.parameter).LengthInBits()
	}

	// Optional Field (payload)
	if m.payload != nil {
		lengthInBits += (*m.payload).LengthInBits()
	}

	return lengthInBits
}

func (m S7Message) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7MessageParse(io spi.ReadBuffer) (spi.Message, error) {

	// Const Field (protocolId)
	protocolId, _protocolIdErr := io.ReadUint8(8)
	if _protocolIdErr != nil {
		return nil, errors.New("Error parsing 'protocolId' field " + _protocolIdErr.Error())
	}
	if protocolId != S7Message_PROTOCOLID {
		return nil, errors.New("Expected constant value " + strconv.Itoa(int(S7Message_PROTOCOLID)) + " but got " + strconv.Itoa(int(protocolId)))
	}

	// Discriminator Field (messageType) (Used as input to a switch field)
	messageType, _messageTypeErr := io.ReadUint8(8)
	if _messageTypeErr != nil {
		return nil, errors.New("Error parsing 'messageType' field " + _messageTypeErr.Error())
	}

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint16(16)
		if _err != nil {
			return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
		}
		if reserved != uint16(0x0000) {
			log.WithFields(log.Fields{
				"expected value": uint16(0x0000),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (tpduReference)
	tpduReference, _tpduReferenceErr := io.ReadUint16(16)
	if _tpduReferenceErr != nil {
		return nil, errors.New("Error parsing 'tpduReference' field " + _tpduReferenceErr.Error())
	}

	// Implicit Field (parameterLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	parameterLength, _parameterLengthErr := io.ReadUint16(16)
	if _parameterLengthErr != nil {
		return nil, errors.New("Error parsing 'parameterLength' field " + _parameterLengthErr.Error())
	}

	// Implicit Field (payloadLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	payloadLength, _payloadLengthErr := io.ReadUint16(16)
	if _payloadLengthErr != nil {
		return nil, errors.New("Error parsing 'payloadLength' field " + _payloadLengthErr.Error())
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer S7MessageInitializer
	var typeSwitchError error
	switch {
	case messageType == 0x01:
		initializer, typeSwitchError = S7MessageRequestParse(io)
	case messageType == 0x02:
		initializer, typeSwitchError = S7MessageResponseParse(io)
	case messageType == 0x03:
		initializer, typeSwitchError = S7MessageResponseDataParse(io)
	case messageType == 0x07:
		initializer, typeSwitchError = S7MessageUserDataParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Optional Field (parameter) (Can be skipped, if a given expression evaluates to false)
	var parameter *IS7Parameter = nil
	if bool((parameterLength) > (0)) {
		_message, _err := S7ParameterParse(io, messageType)
		if _err != nil {
			return nil, errors.New("Error parsing 'parameter' field " + _err.Error())
		}
		var _item IS7Parameter
		_item, _ok := _message.(IS7Parameter)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to IS7Parameter")
		}
		parameter = &_item
	}

	// Optional Field (payload) (Can be skipped, if a given expression evaluates to false)
	var payload *IS7Payload = nil
	if bool((payloadLength) > (0)) {
		_message, _err := S7PayloadParse(io, messageType, (*parameter))
		if _err != nil {
			return nil, errors.New("Error parsing 'payload' field " + _err.Error())
		}
		var _item IS7Payload
		_item, _ok := _message.(IS7Payload)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to IS7Payload")
		}
		payload = &_item
	}

	// Create the instance
	return initializer.initialize(tpduReference, parameter, payload), nil
}

func (m S7Message) Serialize(io spi.WriteBuffer) {
	iS7Message := CastIS7Message(m)

	// Const Field (protocolId)
	io.WriteUint8(8, 0x32)

	// Discriminator Field (messageType) (Used as input to a switch field)
	messageType := uint8(S7MessageMessageType(iS7Message))
	io.WriteUint8(8, (messageType))

	// Reserved Field (reserved)
	io.WriteUint16(16, uint16(0x0000))

	// Simple Field (tpduReference)
	tpduReference := uint16(m.tpduReference)
	io.WriteUint16(16, (tpduReference))

	// Implicit Field (parameterLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	parameterLength := uint16(spi.InlineIf(bool((m.parameter) != (nil)), uint16((*m.parameter).LengthInBytes()), uint16(uint16(0))))
	io.WriteUint16(16, (parameterLength))

	// Implicit Field (payloadLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	payloadLength := uint16(spi.InlineIf(bool((m.payload) != (nil)), uint16((*m.payload).LengthInBytes()), uint16(uint16(0))))
	io.WriteUint16(16, (payloadLength))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	iS7Message.Serialize(io)

	// Optional Field (parameter) (Can be skipped, if the value is null)
	var parameter *IS7Parameter = nil
	if m.parameter != nil {
		parameter = m.parameter
		(*parameter).Serialize(io)
	}

	// Optional Field (payload) (Can be skipped, if the value is null)
	var payload *IS7Payload = nil
	if m.payload != nil {
		payload = m.payload
		(*payload).Serialize(io)
	}
}
