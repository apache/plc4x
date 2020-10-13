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
type DisconnectResponse struct {
	communicationChannelId uint8
	status                 IStatus
	KNXNetIPMessage
}

// The corresponding interface
type IDisconnectResponse interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer) error
}

// Accessors for discriminator values.
func (m DisconnectResponse) MsgType() uint16 {
	return 0x020A
}

func (m DisconnectResponse) initialize() spi.Message {
	return m
}

func NewDisconnectResponse(communicationChannelId uint8, status IStatus) KNXNetIPMessageInitializer {
	return &DisconnectResponse{communicationChannelId: communicationChannelId, status: status}
}

func CastIDisconnectResponse(structType interface{}) IDisconnectResponse {
	castFunc := func(typ interface{}) IDisconnectResponse {
		if iDisconnectResponse, ok := typ.(IDisconnectResponse); ok {
			return iDisconnectResponse
		}
		return nil
	}
	return castFunc(structType)
}

func CastDisconnectResponse(structType interface{}) DisconnectResponse {
	castFunc := func(typ interface{}) DisconnectResponse {
		if sDisconnectResponse, ok := typ.(DisconnectResponse); ok {
			return sDisconnectResponse
		}
		return DisconnectResponse{}
	}
	return castFunc(structType)
}

func (m DisconnectResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Simple field (communicationChannelId)
	lengthInBits += 8

	// Enum Field (status)
	lengthInBits += 8

	return lengthInBits
}

func (m DisconnectResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func DisconnectResponseParse(io *spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

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

	// Create the instance
	return NewDisconnectResponse(communicationChannelId, status), nil
}

func (m DisconnectResponse) Serialize(io spi.WriteBuffer) error {
	ser := func() error {

		// Simple Field (communicationChannelId)
		communicationChannelId := uint8(m.communicationChannelId)
		_communicationChannelIdErr := io.WriteUint8(8, (communicationChannelId))
		if _communicationChannelIdErr != nil {
			return errors.New("Error serializing 'communicationChannelId' field " + _communicationChannelIdErr.Error())
		}

		// Enum field (status)
		status := CastStatus(m.status)
		_statusErr := status.Serialize(io)
		if _statusErr != nil {
			return errors.New("Error serializing 'status' field " + _statusErr.Error())
		}

		return nil
	}
	return KNXNetIPMessageSerialize(io, m.KNXNetIPMessage, CastIKNXNetIPMessage(m), ser)
}
