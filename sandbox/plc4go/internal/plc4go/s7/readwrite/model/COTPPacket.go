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
type COTPPacket struct {
	parameters []COTPParameter
	payload    *S7Message
}

// The corresponding interface
type ICOTPPacket interface {
	spi.Message
	TpduCode() uint8
	Serialize(io spi.WriteBuffer)
}

type COTPPacketInitializer interface {
	initialize(parameters []COTPParameter, payload *S7Message) spi.Message
}

func COTPPacketTpduCode(m ICOTPPacket) uint8 {
	return m.TpduCode()
}

func (m COTPPacket) LengthInBits() uint16 {
	var lengthInBits uint16 = 0

	// Implicit Field (headerLength)
	lengthInBits += 8

	// Discriminator Field (tpduCode)
	lengthInBits += 8

	// Length of sub-type elements will be added by sub-type...

	// Array field
	if len(m.parameters) > 0 {
		for _, element := range m.parameters {
			lengthInBits += element.LengthInBits()
		}
	}

	// Optional Field (payload)
	if m.payload != nil {
		lengthInBits += m.payload.LengthInBits()
	}

	return lengthInBits
}

func (m COTPPacket) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func COTPPacketParse(io spi.ReadBuffer, cotpLen uint16) (spi.Message, error) {
	var startPos = io.GetPos()
	var curPos uint16

	// Implicit Field (headerLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	var headerLength uint8 = io.ReadUint8(8)

	// Discriminator Field (tpduCode) (Used as input to a switch field)
	var tpduCode uint8 = io.ReadUint8(8)

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var initializer COTPPacketInitializer
	var typeSwitchError error
	switch {
	case tpduCode == 0xF0:
		initializer, typeSwitchError = COTPPacketDataParse(io)
	case tpduCode == 0xE0:
		initializer, typeSwitchError = COTPPacketConnectionRequestParse(io)
	case tpduCode == 0xD0:
		initializer, typeSwitchError = COTPPacketConnectionResponseParse(io)
	case tpduCode == 0x80:
		initializer, typeSwitchError = COTPPacketDisconnectRequestParse(io)
	case tpduCode == 0xC0:
		initializer, typeSwitchError = COTPPacketDisconnectResponseParse(io)
	case tpduCode == 0x70:
		initializer, typeSwitchError = COTPPacketTpduErrorParse(io)
	}
	if typeSwitchError != nil {
		return nil, errors.New("Error parsing sub-type for type-switch. " + typeSwitchError.Error())
	}

	// Array field (parameters)
	curPos = io.GetPos() - startPos
	var parameters []COTPParameter
	// Length array
	_parametersLength := uint16(((headerLength) + (1)) - (curPos))
	_parametersEndPos := io.GetPos() + _parametersLength
	for io.GetPos() < _parametersEndPos {
		_message, _err := COTPParameterParse(io, uint8(((headerLength)+(1))-(curPos)))
		if _err != nil {
			return nil, errors.New("Error parsing 'parameters' field " + _err.Error())
		}
		var _item COTPParameter
		_item, _ok := _message.(COTPParameter)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to COTPParameter")
		}
		parameters = append(parameters, _item)
		curPos = io.GetPos() - startPos
	}

	// Optional Field (payload) (Can be skipped, if a given expression evaluates to false)
	curPos = io.GetPos() - startPos
	var payload *S7Message = nil
	if (curPos) < (cotpLen) {
		_message, _err := S7MessageParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'payload' field " + _err.Error())
		}
		var _item S7Message
		_item, _ok := _message.(S7Message)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to S7Message")
		}
		payload = &_item
	}

	// Create the instance
	return initializer.initialize(parameters, payload), nil
}

func (m COTPPacket) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if iCOTPPacket, ok := typ.(ICOTPPacket); ok {

			// Implicit Field (headerLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
			headerLength := uint8((m.LengthInBytes()) - ((spi.InlineIf(((m.payload) != (nil)), uint16(m.payload.LengthInBytes()), uint16(0))) + (1)))
			io.WriteUint8(8, (headerLength))

			// Discriminator Field (tpduCode) (Used as input to a switch field)
			tpduCode := COTPPacketTpduCode(iCOTPPacket)
			io.WriteUint8(8, (tpduCode))

			// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
			iCOTPPacket.Serialize(io)

			// Array Field (parameters)
			if m.parameters != nil {
				for _, _element := range m.parameters {
					_element.Serialize(io)
				}
			}

			// Optional Field (payload) (Can be skipped, if the value is null)
			var payload *S7Message = nil
			if m.payload != nil {
				payload = m.payload
				payload.Serialize(io)
			}
		}
	}
	serializeFunc(m)
}
