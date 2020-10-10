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
type ConnectionStateResponse struct {
	communicationChannelId uint8
	status                 Status
	KNXNetIPMessage
}

// The corresponding interface
type IConnectionStateResponse interface {
	IKNXNetIPMessage
	Serialize(io spi.WriteBuffer)
}

// Accessors for discriminator values.
func (m ConnectionStateResponse) MsgType() uint16 {
	return 0x0208
}

func (m ConnectionStateResponse) initialize() spi.Message {
	return m
}

func NewConnectionStateResponse(communicationChannelId uint8, status Status) KNXNetIPMessageInitializer {
	return &ConnectionStateResponse{communicationChannelId: communicationChannelId, status: status}
}

func (m ConnectionStateResponse) LengthInBits() uint16 {
	var lengthInBits uint16 = m.KNXNetIPMessage.LengthInBits()

	// Simple field (communicationChannelId)
	lengthInBits += 8

	// Enum Field (status)
	lengthInBits += 8

	return lengthInBits
}

func (m ConnectionStateResponse) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ConnectionStateResponseParse(io spi.ReadBuffer) (KNXNetIPMessageInitializer, error) {

	// Simple Field (communicationChannelId)
	var communicationChannelId uint8 = io.ReadUint8(8)

	// Enum field (status)
	status, _statusErr := StatusParse(io)
	if _statusErr != nil {
		return nil, errors.New("Error parsing 'status' field " + _statusErr.Error())
	}

	// Create the instance
	return NewConnectionStateResponse(communicationChannelId, status), nil
}

func (m ConnectionStateResponse) Serialize(io spi.WriteBuffer) {
	serializeFunc := func(typ interface{}) {
		if _, ok := typ.(IConnectionStateResponse); ok {

			// Simple Field (communicationChannelId)
			var communicationChannelId uint8 = m.communicationChannelId
			io.WriteUint8(8, (communicationChannelId))

			// Enum field (status)
			status := m.status
			status.Serialize(io)
		}
	}
	serializeFunc(m)
}
