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
	parameters []ICOTPParameter
	payload    *IS7Message
}

// The corresponding interface
type ICOTPPacket interface {
	spi.Message
	TpduCode() uint8
	Serialize(io spi.WriteBuffer)
}

type COTPPacketInitializer interface {
	initialize(parameters []ICOTPParameter, payload *IS7Message) spi.Message
}

func COTPPacketTpduCode(m ICOTPPacket) uint8 {
	return m.TpduCode()
}

func CastICOTPPacket(structType interface{}) ICOTPPacket {
	castFunc := func(typ interface{}) ICOTPPacket {
		if iCOTPPacket, ok := typ.(ICOTPPacket); ok {
			return iCOTPPacket
		}
		return nil
	}
	return castFunc(structType)
}

func CastCOTPPacket(structType interface{}) COTPPacket {
	castFunc := func(typ interface{}) COTPPacket {
		if sCOTPPacket, ok := typ.(COTPPacket); ok {
			return sCOTPPacket
		}
		return COTPPacket{}
	}
	return castFunc(structType)
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
		lengthInBits += (*m.payload).LengthInBits()
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
	headerLength, _headerLengthErr := io.ReadUint8(8)
	if _headerLengthErr != nil {
		return nil, errors.New("Error parsing 'headerLength' field " + _headerLengthErr.Error())
	}

	// Discriminator Field (tpduCode) (Used as input to a switch field)
	tpduCode, _tpduCodeErr := io.ReadUint8(8)
	if _tpduCodeErr != nil {
		return nil, errors.New("Error parsing 'tpduCode' field " + _tpduCodeErr.Error())
	}

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
	var parameters []ICOTPParameter
	// Length array
	_parametersLength := uint16(uint16(uint16(headerLength)+uint16(uint16(1)))) - uint16(curPos)
	_parametersEndPos := io.GetPos() + uint16(_parametersLength)
	for io.GetPos() < _parametersEndPos {
		_message, _err := COTPParameterParse(io, uint8(uint8(uint8(headerLength)+uint8(uint8(1))))-uint8(curPos))
		if _err != nil {
			return nil, errors.New("Error parsing 'parameters' field " + _err.Error())
		}
		var _item ICOTPParameter
		_item, _ok := _message.(ICOTPParameter)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to COTPParameter")
		}
		parameters = append(parameters, _item)
		curPos = io.GetPos() - startPos
	}

	// Optional Field (payload) (Can be skipped, if a given expression evaluates to false)
	curPos = io.GetPos() - startPos
	var payload *IS7Message = nil
	if bool((curPos) < (cotpLen)) {
		_message, _err := S7MessageParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'payload' field " + _err.Error())
		}
		var _item IS7Message
		_item, _ok := _message.(IS7Message)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to IS7Message")
		}
		payload = &_item
	}

	// Create the instance
	return initializer.initialize(parameters, payload), nil
}

func COTPPacketSerialize(io spi.WriteBuffer, m COTPPacket, i ICOTPPacket, childSerialize func()) {

	// Implicit Field (headerLength) (Used for parsing, but it's value is not stored as it's implicitly given by the objects content)
	headerLength := uint8(uint8(uint8(m.LengthInBytes())) - uint8(uint8(uint8(uint8(spi.InlineIf(bool(bool((m.payload) != (nil))), uint16((*m.payload).LengthInBytes()), uint16(uint8(0)))))+uint8(uint8(1)))))
	io.WriteUint8(8, (headerLength))

	// Discriminator Field (tpduCode) (Used as input to a switch field)
	tpduCode := uint8(i.TpduCode())
	io.WriteUint8(8, (tpduCode))

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	childSerialize()

	// Array Field (parameters)
	if m.parameters != nil {
		for _, _element := range m.parameters {
			_element.Serialize(io)
		}
	}

	// Optional Field (payload) (Can be skipped, if the value is null)
	var payload *IS7Message = nil
	if m.payload != nil {
		payload = m.payload
		CastIS7Message(*payload).Serialize(io)
	}

}
