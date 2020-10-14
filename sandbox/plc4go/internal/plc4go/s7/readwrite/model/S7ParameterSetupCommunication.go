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
)

// The data-structure of this message
type S7ParameterSetupCommunication struct {
	MaxAmqCaller uint16
	MaxAmqCallee uint16
	PduLength    uint16
	S7Parameter
}

// The corresponding interface
type IS7ParameterSetupCommunication interface {
	IS7Parameter
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m S7ParameterSetupCommunication) ParameterType() uint8 {
	return 0xF0
}

func (m S7ParameterSetupCommunication) MessageType() uint8 {
	return 0
}

func (m S7ParameterSetupCommunication) initialize() spi.Message {
	return m
}

func NewS7ParameterSetupCommunication(maxAmqCaller uint16, maxAmqCallee uint16, pduLength uint16) S7ParameterInitializer {
	return &S7ParameterSetupCommunication{MaxAmqCaller: maxAmqCaller, MaxAmqCallee: maxAmqCallee, PduLength: pduLength}
}

func CastIS7ParameterSetupCommunication(structType interface{}) IS7ParameterSetupCommunication {
	castFunc := func(typ interface{}) IS7ParameterSetupCommunication {
		if iS7ParameterSetupCommunication, ok := typ.(IS7ParameterSetupCommunication); ok {
			return iS7ParameterSetupCommunication
		}
		return nil
	}
	return castFunc(structType)
}

func CastS7ParameterSetupCommunication(structType interface{}) S7ParameterSetupCommunication {
	castFunc := func(typ interface{}) S7ParameterSetupCommunication {
		if sS7ParameterSetupCommunication, ok := typ.(S7ParameterSetupCommunication); ok {
			return sS7ParameterSetupCommunication
		}
		return S7ParameterSetupCommunication{}
	}
	return castFunc(structType)
}

func (m S7ParameterSetupCommunication) LengthInBits() uint16 {
	var lengthInBits uint16 = m.S7Parameter.LengthInBits()

	// Reserved Field (reserved)
	lengthInBits += 8

	// Simple field (maxAmqCaller)
	lengthInBits += 16

	// Simple field (maxAmqCallee)
	lengthInBits += 16

	// Simple field (pduLength)
	lengthInBits += 16

	return lengthInBits
}

func (m S7ParameterSetupCommunication) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func S7ParameterSetupCommunicationParse(io *spi.ReadBuffer) (S7ParameterInitializer, error) {

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := io.ReadUint8(8)
		if _err != nil {
			return nil, errors.New("Error parsing 'reserved' field " + _err.Error())
		}
		if reserved != uint8(0x00) {
			log.WithFields(log.Fields{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Info("Got unexpected response.")
		}
	}

	// Simple Field (maxAmqCaller)
	maxAmqCaller, _maxAmqCallerErr := io.ReadUint16(16)
	if _maxAmqCallerErr != nil {
		return nil, errors.New("Error parsing 'maxAmqCaller' field " + _maxAmqCallerErr.Error())
	}

	// Simple Field (maxAmqCallee)
	maxAmqCallee, _maxAmqCalleeErr := io.ReadUint16(16)
	if _maxAmqCalleeErr != nil {
		return nil, errors.New("Error parsing 'maxAmqCallee' field " + _maxAmqCalleeErr.Error())
	}

	// Simple Field (pduLength)
	pduLength, _pduLengthErr := io.ReadUint16(16)
	if _pduLengthErr != nil {
		return nil, errors.New("Error parsing 'pduLength' field " + _pduLengthErr.Error())
	}

	// Create the instance
	return NewS7ParameterSetupCommunication(maxAmqCaller, maxAmqCallee, pduLength), nil
}

func (m S7ParameterSetupCommunication) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Reserved Field (reserved)
		{
			_err := io.WriteUint8(8, uint8(0x00))
			if _err != nil {
				return errors.New("Error serializing 'reserved' field " + _err.Error())
			}
		}

		// Simple Field (maxAmqCaller)
		maxAmqCaller := uint16(m.MaxAmqCaller)
		_maxAmqCallerErr := io.WriteUint16(16, (maxAmqCaller))
		if _maxAmqCallerErr != nil {
			return errors.New("Error serializing 'maxAmqCaller' field " + _maxAmqCallerErr.Error())
		}

		// Simple Field (maxAmqCallee)
		maxAmqCallee := uint16(m.MaxAmqCallee)
		_maxAmqCalleeErr := io.WriteUint16(16, (maxAmqCallee))
		if _maxAmqCalleeErr != nil {
			return errors.New("Error serializing 'maxAmqCallee' field " + _maxAmqCalleeErr.Error())
		}

		// Simple Field (pduLength)
		pduLength := uint16(m.PduLength)
		_pduLengthErr := io.WriteUint16(16, (pduLength))
		if _pduLengthErr != nil {
			return errors.New("Error serializing 'pduLength' field " + _pduLengthErr.Error())
		}

		return nil
	}
	return S7ParameterSerialize(io, m.S7Parameter, CastIS7Parameter(m), ser)
}
