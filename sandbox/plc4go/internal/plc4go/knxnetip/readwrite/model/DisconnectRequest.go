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
)

// The data-structure of this message
type DisconnectRequest struct {
	communicationChannelId uint8
	hpaiControlEndpoint    IHPAIControlEndpoint
	KNXNetIPMessage
}

// The corresponding interface
type IDisconnectRequest interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m DisconnectRequest) MsgType() uint16 {
	return 0x0209
}

func (m DisconnectRequest) initialize() spi.Message {
	return m
}

func NewDisconnectRequest(communicationChannelId uint8, hpaiControlEndpoint IHPAIControlEndpoint) KNXNetIPMessageInitializer {
	return &DisconnectRequest{communicationChannelId: communicationChannelId, hpaiControlEndpoint: hpaiControlEndpoint}
}

func CastIDisconnectRequest(structType interface{}) IDisconnectRequest {
	castFunc := func(typ interface{}) IDisconnectRequest {
		if iDisconnectRequest, ok := typ.(IDisconnectRequest); ok {
			return iDisconnectRequest
		}
		return nil
	}
	return castFunc(structType)
}

func CastDisconnectRequest(structType interface{}) DisconnectRequest {
	castFunc := func(typ interface{}) DisconnectRequest {
		if sDisconnectRequest, ok := typ.(DisconnectRequest); ok {
			return sDisconnectRequest
		}
		return DisconnectRequest{}
	}
	return castFunc(structType)
}

func (m DisconnectRequest) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Simple field (communicationChannelId)
	lengthInBits += 8

	// Reserved Field (reserved)
	lengthInBits += 8

	// Simple field (hpaiControlEndpoint)
	lengthInBits += m.hpaiControlEndpoint.LengthInBits()

	return lengthInBits
}

func (m DisconnectRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func DisconnectRequestParse(io *spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

	// Simple Field (communicationChannelId)
	communicationChannelId, _communicationChannelIdErr := io.ReadUint8(8)
	if _communicationChannelIdErr != nil {
		return nil, errors.New("Error parsing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
	}

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

	// Simple Field (hpaiControlEndpoint)
	_hpaiControlEndpointMessage, _err := HPAIControlEndpointParse(io)
	if _err != nil {
		return nil, errors.New("Error parsing simple field 'hpaiControlEndpoint'. " + _err.Error())
	}
	var hpaiControlEndpoint IHPAIControlEndpoint
	hpaiControlEndpoint, _hpaiControlEndpointOk := _hpaiControlEndpointMessage.(IHPAIControlEndpoint)
	if !_hpaiControlEndpointOk {
		return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_hpaiControlEndpointMessage).Name() + " to IHPAIControlEndpoint")
	}

	// Create the instance
	return NewDisconnectRequest(communicationChannelId, hpaiControlEndpoint), nil
}

func (m DisconnectRequest) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (communicationChannelId)
		communicationChannelId := uint8(m.communicationChannelId)
		_communicationChannelIdErr := io.WriteUint8(8, (communicationChannelId))
		if _communicationChannelIdErr != nil {
			return errors.New("Error serializing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
		}

		// Reserved Field (reserved)
		{
			_err := io.WriteUint8(8, uint8(0x00))
			if _err != nil {
				return errors.New("Error serializing 'reserved' field " + _err.Error())
			}
		}

		// Simple Field (hpaiControlEndpoint)
		hpaiControlEndpoint := CastIHPAIControlEndpoint(m.hpaiControlEndpoint)
		_hpaiControlEndpointErr := hpaiControlEndpoint.Serialize(io)
		if _hpaiControlEndpointErr != nil {
			return errors.New("Error serializing 'hpaiControlEndpoint' field " + _hpaiControlEndpointErr.Error())
		}

		return nil
	}
	return KNXNetIPMessageSerialize(io, m.KNXNetIPMessage, CastIKNXNetIPMessage(m), ser)
}
