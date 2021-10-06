/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package model

import (
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// The data-structure of this message
type ConnectionRequest struct {
	HpaiDiscoveryEndpoint        *HPAIDiscoveryEndpoint
	HpaiDataEndpoint             *HPAIDataEndpoint
	ConnectionRequestInformation *ConnectionRequestInformation
	Parent                       *KnxNetIpMessage
}

// The corresponding interface
type IConnectionRequest interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(writeBuffer utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *ConnectionRequest) MsgType() uint16 {
	return 0x0205
}

func (m *ConnectionRequest) InitializeParent(parent *KnxNetIpMessage) {
}

func NewConnectionRequest(hpaiDiscoveryEndpoint *HPAIDiscoveryEndpoint, hpaiDataEndpoint *HPAIDataEndpoint, connectionRequestInformation *ConnectionRequestInformation) *KnxNetIpMessage {
	child := &ConnectionRequest{
		HpaiDiscoveryEndpoint:        hpaiDiscoveryEndpoint,
		HpaiDataEndpoint:             hpaiDataEndpoint,
		ConnectionRequestInformation: connectionRequestInformation,
		Parent:                       NewKnxNetIpMessage(),
	}
	child.Parent.Child = child
	return child.Parent
}

func CastConnectionRequest(structType interface{}) *ConnectionRequest {
	castFunc := func(typ interface{}) *ConnectionRequest {
		if casted, ok := typ.(ConnectionRequest); ok {
			return &casted
		}
		if casted, ok := typ.(*ConnectionRequest); ok {
			return casted
		}
		if casted, ok := typ.(KnxNetIpMessage); ok {
			return CastConnectionRequest(casted.Child)
		}
		if casted, ok := typ.(*KnxNetIpMessage); ok {
			return CastConnectionRequest(casted.Child)
		}
		return nil
	}
	return castFunc(structType)
}

func (m *ConnectionRequest) GetTypeName() string {
	return "ConnectionRequest"
}

func (m *ConnectionRequest) LengthInBits() uint16 {
	return m.LengthInBitsConditional(false)
}

func (m *ConnectionRequest) LengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.Parent.ParentLengthInBits())

	// Simple field (hpaiDiscoveryEndpoint)
	lengthInBits += m.HpaiDiscoveryEndpoint.LengthInBits()

	// Simple field (hpaiDataEndpoint)
	lengthInBits += m.HpaiDataEndpoint.LengthInBits()

	// Simple field (connectionRequestInformation)
	lengthInBits += m.ConnectionRequestInformation.LengthInBits()

	return lengthInBits
}

func (m *ConnectionRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func ConnectionRequestParse(readBuffer utils.ReadBuffer) (*KnxNetIpMessage, error) {
	if pullErr := readBuffer.PullContext("ConnectionRequest"); pullErr != nil {
		return nil, pullErr
	}

	// Simple Field (hpaiDiscoveryEndpoint)
	if pullErr := readBuffer.PullContext("hpaiDiscoveryEndpoint"); pullErr != nil {
		return nil, pullErr
	}
	hpaiDiscoveryEndpoint, _hpaiDiscoveryEndpointErr := HPAIDiscoveryEndpointParse(readBuffer)
	if _hpaiDiscoveryEndpointErr != nil {
		return nil, errors.Wrap(_hpaiDiscoveryEndpointErr, "Error parsing 'hpaiDiscoveryEndpoint' field")
	}
	if closeErr := readBuffer.CloseContext("hpaiDiscoveryEndpoint"); closeErr != nil {
		return nil, closeErr
	}

	// Simple Field (hpaiDataEndpoint)
	if pullErr := readBuffer.PullContext("hpaiDataEndpoint"); pullErr != nil {
		return nil, pullErr
	}
	hpaiDataEndpoint, _hpaiDataEndpointErr := HPAIDataEndpointParse(readBuffer)
	if _hpaiDataEndpointErr != nil {
		return nil, errors.Wrap(_hpaiDataEndpointErr, "Error parsing 'hpaiDataEndpoint' field")
	}
	if closeErr := readBuffer.CloseContext("hpaiDataEndpoint"); closeErr != nil {
		return nil, closeErr
	}

	// Simple Field (connectionRequestInformation)
	if pullErr := readBuffer.PullContext("connectionRequestInformation"); pullErr != nil {
		return nil, pullErr
	}
	connectionRequestInformation, _connectionRequestInformationErr := ConnectionRequestInformationParse(readBuffer)
	if _connectionRequestInformationErr != nil {
		return nil, errors.Wrap(_connectionRequestInformationErr, "Error parsing 'connectionRequestInformation' field")
	}
	if closeErr := readBuffer.CloseContext("connectionRequestInformation"); closeErr != nil {
		return nil, closeErr
	}

	if closeErr := readBuffer.CloseContext("ConnectionRequest"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &ConnectionRequest{
		HpaiDiscoveryEndpoint:        CastHPAIDiscoveryEndpoint(hpaiDiscoveryEndpoint),
		HpaiDataEndpoint:             CastHPAIDataEndpoint(hpaiDataEndpoint),
		ConnectionRequestInformation: CastConnectionRequestInformation(connectionRequestInformation),
		Parent:                       &KnxNetIpMessage{},
	}
	_child.Parent.Child = _child
	return _child.Parent, nil
}

func (m *ConnectionRequest) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ConnectionRequest"); pushErr != nil {
			return pushErr
		}

		// Simple Field (hpaiDiscoveryEndpoint)
		if pushErr := writeBuffer.PushContext("hpaiDiscoveryEndpoint"); pushErr != nil {
			return pushErr
		}
		_hpaiDiscoveryEndpointErr := m.HpaiDiscoveryEndpoint.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("hpaiDiscoveryEndpoint"); popErr != nil {
			return popErr
		}
		if _hpaiDiscoveryEndpointErr != nil {
			return errors.Wrap(_hpaiDiscoveryEndpointErr, "Error serializing 'hpaiDiscoveryEndpoint' field")
		}

		// Simple Field (hpaiDataEndpoint)
		if pushErr := writeBuffer.PushContext("hpaiDataEndpoint"); pushErr != nil {
			return pushErr
		}
		_hpaiDataEndpointErr := m.HpaiDataEndpoint.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("hpaiDataEndpoint"); popErr != nil {
			return popErr
		}
		if _hpaiDataEndpointErr != nil {
			return errors.Wrap(_hpaiDataEndpointErr, "Error serializing 'hpaiDataEndpoint' field")
		}

		// Simple Field (connectionRequestInformation)
		if pushErr := writeBuffer.PushContext("connectionRequestInformation"); pushErr != nil {
			return pushErr
		}
		_connectionRequestInformationErr := m.ConnectionRequestInformation.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("connectionRequestInformation"); popErr != nil {
			return popErr
		}
		if _connectionRequestInformationErr != nil {
			return errors.Wrap(_connectionRequestInformationErr, "Error serializing 'connectionRequestInformation' field")
		}

		if popErr := writeBuffer.PopContext("ConnectionRequest"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.Parent.SerializeParent(writeBuffer, m, ser)
}

func (m *ConnectionRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	m.Serialize(buffer)
	return buffer.GetBox().String()
}
