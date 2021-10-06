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
type CIPEncapsulationReadRequest struct {
	Request *DF1RequestMessage
	Parent  *CIPEncapsulationPacket
}

// The corresponding interface
type ICIPEncapsulationReadRequest interface {
	LengthInBytes() uint16
	LengthInBits() uint16
	Serialize(writeBuffer utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
// Accessors for discriminator values.
///////////////////////////////////////////////////////////
func (m *CIPEncapsulationReadRequest) CommandType() uint16 {
	return 0x0107
}

func (m *CIPEncapsulationReadRequest) InitializeParent(parent *CIPEncapsulationPacket, sessionHandle uint32, status uint32, senderContext []uint8, options uint32) {
	m.Parent.SessionHandle = sessionHandle
	m.Parent.Status = status
	m.Parent.SenderContext = senderContext
	m.Parent.Options = options
}

func NewCIPEncapsulationReadRequest(request *DF1RequestMessage, sessionHandle uint32, status uint32, senderContext []uint8, options uint32) *CIPEncapsulationPacket {
	child := &CIPEncapsulationReadRequest{
		Request: request,
		Parent:  NewCIPEncapsulationPacket(sessionHandle, status, senderContext, options),
	}
	child.Parent.Child = child
	return child.Parent
}

func CastCIPEncapsulationReadRequest(structType interface{}) *CIPEncapsulationReadRequest {
	castFunc := func(typ interface{}) *CIPEncapsulationReadRequest {
		if casted, ok := typ.(CIPEncapsulationReadRequest); ok {
			return &casted
		}
		if casted, ok := typ.(*CIPEncapsulationReadRequest); ok {
			return casted
		}
		if casted, ok := typ.(CIPEncapsulationPacket); ok {
			return CastCIPEncapsulationReadRequest(casted.Child)
		}
		if casted, ok := typ.(*CIPEncapsulationPacket); ok {
			return CastCIPEncapsulationReadRequest(casted.Child)
		}
		return nil
	}
	return castFunc(structType)
}

func (m *CIPEncapsulationReadRequest) GetTypeName() string {
	return "CIPEncapsulationReadRequest"
}

func (m *CIPEncapsulationReadRequest) LengthInBits() uint16 {
	return m.LengthInBitsConditional(false)
}

func (m *CIPEncapsulationReadRequest) LengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.Parent.ParentLengthInBits())

	// Simple field (request)
	lengthInBits += m.Request.LengthInBits()

	return lengthInBits
}

func (m *CIPEncapsulationReadRequest) LengthInBytes() uint16 {
	return m.LengthInBits() / 8
}

func CIPEncapsulationReadRequestParse(readBuffer utils.ReadBuffer) (*CIPEncapsulationPacket, error) {
	if pullErr := readBuffer.PullContext("CIPEncapsulationReadRequest"); pullErr != nil {
		return nil, pullErr
	}

	// Simple Field (request)
	if pullErr := readBuffer.PullContext("request"); pullErr != nil {
		return nil, pullErr
	}
	request, _requestErr := DF1RequestMessageParse(readBuffer)
	if _requestErr != nil {
		return nil, errors.Wrap(_requestErr, "Error parsing 'request' field")
	}
	if closeErr := readBuffer.CloseContext("request"); closeErr != nil {
		return nil, closeErr
	}

	if closeErr := readBuffer.CloseContext("CIPEncapsulationReadRequest"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &CIPEncapsulationReadRequest{
		Request: CastDF1RequestMessage(request),
		Parent:  &CIPEncapsulationPacket{},
	}
	_child.Parent.Child = _child
	return _child.Parent, nil
}

func (m *CIPEncapsulationReadRequest) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CIPEncapsulationReadRequest"); pushErr != nil {
			return pushErr
		}

		// Simple Field (request)
		if pushErr := writeBuffer.PushContext("request"); pushErr != nil {
			return pushErr
		}
		_requestErr := m.Request.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("request"); popErr != nil {
			return popErr
		}
		if _requestErr != nil {
			return errors.Wrap(_requestErr, "Error serializing 'request' field")
		}

		if popErr := writeBuffer.PopContext("CIPEncapsulationReadRequest"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.Parent.SerializeParent(writeBuffer, m, ser)
}

func (m *CIPEncapsulationReadRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	m.Serialize(buffer)
	return buffer.GetBox().String()
}
