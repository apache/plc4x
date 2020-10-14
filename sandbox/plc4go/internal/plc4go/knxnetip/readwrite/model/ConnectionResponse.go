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
type ConnectionResponse struct {
	CommunicationChannelId      uint8
	Status                      IStatus
	HpaiDataEndpoint            *IHPAIDataEndpoint
	ConnectionResponseDataBlock *IConnectionResponseDataBlock
	KNXNetIPMessage
}

// The corresponding interface
type IConnectionResponse interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m ConnectionResponse) MsgType() uint16 {
	return 0x0206
}

func (m ConnectionResponse) initialize() spi.Message {
	return m
}

func NewConnectionResponse(communicationChannelId uint8, status IStatus, hpaiDataEndpoint *IHPAIDataEndpoint, connectionResponseDataBlock *IConnectionResponseDataBlock) KNXNetIPMessageInitializer {
	return &ConnectionResponse{CommunicationChannelId: communicationChannelId, Status: status, HpaiDataEndpoint: hpaiDataEndpoint, ConnectionResponseDataBlock: connectionResponseDataBlock}
}

func CastIConnectionResponse(structType interface{}) IConnectionResponse {
	castFunc := func(typ interface{}) IConnectionResponse {
		if iConnectionResponse, ok := typ.(IConnectionResponse); ok {
			return iConnectionResponse
		}
		return nil
	}
	return castFunc(structType)
}

func CastConnectionResponse(structType interface{}) ConnectionResponse {
	castFunc := func(typ interface{}) ConnectionResponse {
		if sConnectionResponse, ok := typ.(ConnectionResponse); ok {
			return sConnectionResponse
		}
		return ConnectionResponse{}
	}
	return castFunc(structType)
}

func (m ConnectionResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Simple field (communicationChannelId)
	lengthInBits += 8

	// Enum Field (status)
	lengthInBits += 8

	// Optional Field (hpaiDataEndpoint)
	if m.HpaiDataEndpoint != nil {
		lengthInBits += (*m.HpaiDataEndpoint).LengthInBits()
	}

	// Optional Field (connectionResponseDataBlock)
	if m.ConnectionResponseDataBlock != nil {
		lengthInBits += (*m.ConnectionResponseDataBlock).LengthInBits()
	}

	return lengthInBits
}

func (m ConnectionResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ConnectionResponseParse(io *spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

	// Simple Field (communicationChannelId)
	communicationChannelId, _communicationChannelIdErr := io.ReadUint8(8)
	if _communicationChannelIdErr != nil {
		return nil, errors.New("Error parsing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
	}

	// Enum field (status)
	status, _statusErr := StatusParse(io)
	if _statusErr != nil {
		return nil, errors.New("Error parsing 'status' field " + _statusErr.Error())
	}

	// Optional Field (hpaiDataEndpoint) (Can be skipped, if a given expression evaluates to false)
	var hpaiDataEndpoint *IHPAIDataEndpoint = nil
	if bool((status) == (Status_NO_ERROR)) {
		_message, _err := HPAIDataEndpointParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'hpaiDataEndpoint' field " + _err.Error())
		}
		var _item IHPAIDataEndpoint
		_item, _ok := _message.(IHPAIDataEndpoint)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to IHPAIDataEndpoint")
		}
		hpaiDataEndpoint = &_item
	}

	// Optional Field (connectionResponseDataBlock) (Can be skipped, if a given expression evaluates to false)
	var connectionResponseDataBlock *IConnectionResponseDataBlock = nil
	if bool((status) == (Status_NO_ERROR)) {
		_message, _err := ConnectionResponseDataBlockParse(io)
		if _err != nil {
			return nil, errors.New("Error parsing 'connectionResponseDataBlock' field " + _err.Error())
		}
		var _item IConnectionResponseDataBlock
		_item, _ok := _message.(IConnectionResponseDataBlock)
		if !_ok {
			return nil, errors.New("Couldn't cast message of type " + reflect.TypeOf(_item).Name() + " to IConnectionResponseDataBlock")
		}
		connectionResponseDataBlock = &_item
	}

	// Create the instance
	return NewConnectionResponse(communicationChannelId, status, hpaiDataEndpoint, connectionResponseDataBlock), nil
}

func (m ConnectionResponse) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (communicationChannelId)
		communicationChannelId := uint8(m.CommunicationChannelId)
		_communicationChannelIdErr := io.WriteUint8(8, (communicationChannelId))
		if _communicationChannelIdErr != nil {
			return errors.New("Error serializing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
		}

		// Enum field (status)
		status := CastStatus(m.Status)
		_statusErr := status.Serialize(io)
		if _statusErr != nil {
			return errors.New("Error serializing 'status' field " + _statusErr.Error())
		}

		// Optional Field (hpaiDataEndpoint) (Can be skipped, if the value is null)
		var hpaiDataEndpoint *IHPAIDataEndpoint = nil
		if m.HpaiDataEndpoint != nil {
			hpaiDataEndpoint = m.HpaiDataEndpoint
			_hpaiDataEndpointErr := CastIHPAIDataEndpoint(*hpaiDataEndpoint).Serialize(io)
			if _hpaiDataEndpointErr != nil {
				return errors.New("Error serializing 'hpaiDataEndpoint' field " + _hpaiDataEndpointErr.Error())
			}
		}

		// Optional Field (connectionResponseDataBlock) (Can be skipped, if the value is null)
		var connectionResponseDataBlock *IConnectionResponseDataBlock = nil
		if m.ConnectionResponseDataBlock != nil {
			connectionResponseDataBlock = m.ConnectionResponseDataBlock
			_connectionResponseDataBlockErr := CastIConnectionResponseDataBlock(*connectionResponseDataBlock).Serialize(io)
			if _connectionResponseDataBlockErr != nil {
				return errors.New("Error serializing 'connectionResponseDataBlock' field " + _connectionResponseDataBlockErr.Error())
			}
		}

		return nil
	}
	return KNXNetIPMessageSerialize(io, m.KNXNetIPMessage, CastIKNXNetIPMessage(m), ser)
}
