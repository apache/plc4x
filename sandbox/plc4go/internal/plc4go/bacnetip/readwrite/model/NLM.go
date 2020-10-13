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
type NLM struct {
	vendorId *uint16
}

// The corresponding interface
type INLM interface {
	spi.Message
	MessageType() uint8
	Serialize(io spi.WriteBuffer) error
}

type NLMInitializer interface {
	initialize(vendorId *uint16) spi.Message
}

func NLMMessageType(m INLM) uint8 {
	return m.MessageType()
}

func CastINLM(structType interface{}) INLM {
	castFunc := func(typ interface{}) INLM {
		if iNLM, ok := typ.(INLM); ok {
			return iNLM
		}
		return nil
	}
	return castFunc(structType)
}

func CastNLM(structType interface{}) NLM {
	castFunc := func(typ interface{}) NLM {
		if sNLM, ok := typ.(NLM); ok {
			return sNLM
		}
		return NLM{}
	}
	return castFunc(structType)
}

func (m NLM) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Discriminator Field (messageType)
	lengthInBits += 8

	// Optional Field (vendorId)
	if m.vendorId != nil {
		lengthInBits += 16
	}

	// Length of sub-type elements will be added by sub-type...

	return lengthInBits
}

func (m NLM) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func NLMParse(io *spi.ReadBuffer, apduLength uint16) (spi.Message, error) {

	// Discriminator Field (messageType) (Used as input to a switch field)
	messageType, _messageTypeErr := io.ReadUint8(8)
	if _messageTypeErr != nil {
		return nil, errors.New("Error parsing 'messageType' field " + _messageTypeErr.Error())
	}

	// Optional Field (vendorId) (Can be skipped, if a given expression evaluates to false)
	var vendorId *uint16 = nil
	if bool(bool(bool((messageType) >= (128)))) && bool(bool(bool((messageType) <= (255)))) {
		_val, _err := io.ReadUint16(16)
		if _err != nil {
			return nil, errors.New("Error parsing 'vendorId' field " + _err.Error())
		}

		vendorId = &_val
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer NLMInitializer
	var typeSwitchError error
	switch {
	case messageType == 0x0:
		initializer, typeSwitchError = NLMWhoIsRouterToNetworkParse(io, apduLength, messageType)
	case messageType == 0x1:
		initializer, typeSwitchError = NLMIAmRouterToNetworkParse(io, apduLength, messageType)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Create the instance
	return initializer.initialize(vendorId), nil
}

func NLMSerialize(io spi.WriteBuffer, m NLM, i INLM, childSerialize func() error) error {

	// Discriminator Field (messageType) (Used as input to a switch field)
	messageType := uint8(i.MessageType())
	_messageTypeErr := io.WriteUint8(8, (messageType))
	if _messageTypeErr != nil {
		return errors.New("Error serializing 'messageType' field " + _messageTypeErr.Error())
	}

	// Optional Field (vendorId) (Can be skipped, if the value is null)
	var vendorId *uint16 = nil
	if m.vendorId != nil {
		vendorId = m.vendorId
		_vendorIdErr := io.WriteUint16(16, *(vendorId))
		if _vendorIdErr != nil {
			return errors.New("Error serializing 'vendorId' field " + _vendorIdErr.Error())
		}
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	_typeSwitchErr := childSerialize()
	if _typeSwitchErr != nil {
		return errors.New("Error serializing sub-type field " + _typeSwitchErr.Error())
	}

	return nil
}
