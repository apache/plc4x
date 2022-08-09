/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   https://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// COTPPacketDisconnectRequest is the corresponding interface of COTPPacketDisconnectRequest
type COTPPacketDisconnectRequest interface {
	utils.LengthAware
	utils.Serializable
	COTPPacket
	// GetDestinationReference returns DestinationReference (property field)
	GetDestinationReference() uint16
	// GetSourceReference returns SourceReference (property field)
	GetSourceReference() uint16
	// GetProtocolClass returns ProtocolClass (property field)
	GetProtocolClass() COTPProtocolClass
}

// COTPPacketDisconnectRequestExactly can be used when we want exactly this type and not a type which fulfills COTPPacketDisconnectRequest.
// This is useful for switch cases.
type COTPPacketDisconnectRequestExactly interface {
	COTPPacketDisconnectRequest
	isCOTPPacketDisconnectRequest() bool
}

// _COTPPacketDisconnectRequest is the data-structure of this message
type _COTPPacketDisconnectRequest struct {
	*_COTPPacket
	DestinationReference uint16
	SourceReference      uint16
	ProtocolClass        COTPProtocolClass
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_COTPPacketDisconnectRequest) GetTpduCode() uint8 {
	return 0x80
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_COTPPacketDisconnectRequest) InitializeParent(parent COTPPacket, parameters []COTPParameter, payload S7Message) {
	m.Parameters = parameters
	m.Payload = payload
}

func (m *_COTPPacketDisconnectRequest) GetParent() COTPPacket {
	return m._COTPPacket
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_COTPPacketDisconnectRequest) GetDestinationReference() uint16 {
	return m.DestinationReference
}

func (m *_COTPPacketDisconnectRequest) GetSourceReference() uint16 {
	return m.SourceReference
}

func (m *_COTPPacketDisconnectRequest) GetProtocolClass() COTPProtocolClass {
	return m.ProtocolClass
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewCOTPPacketDisconnectRequest factory function for _COTPPacketDisconnectRequest
func NewCOTPPacketDisconnectRequest(destinationReference uint16, sourceReference uint16, protocolClass COTPProtocolClass, parameters []COTPParameter, payload S7Message, cotpLen uint16) *_COTPPacketDisconnectRequest {
	_result := &_COTPPacketDisconnectRequest{
		DestinationReference: destinationReference,
		SourceReference:      sourceReference,
		ProtocolClass:        protocolClass,
		_COTPPacket:          NewCOTPPacket(parameters, payload, cotpLen),
	}
	_result._COTPPacket._COTPPacketChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastCOTPPacketDisconnectRequest(structType interface{}) COTPPacketDisconnectRequest {
	if casted, ok := structType.(COTPPacketDisconnectRequest); ok {
		return casted
	}
	if casted, ok := structType.(*COTPPacketDisconnectRequest); ok {
		return *casted
	}
	return nil
}

func (m *_COTPPacketDisconnectRequest) GetTypeName() string {
	return "COTPPacketDisconnectRequest"
}

func (m *_COTPPacketDisconnectRequest) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_COTPPacketDisconnectRequest) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (destinationReference)
	lengthInBits += 16

	// Simple field (sourceReference)
	lengthInBits += 16

	// Simple field (protocolClass)
	lengthInBits += 8

	return lengthInBits
}

func (m *_COTPPacketDisconnectRequest) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func COTPPacketDisconnectRequestParse(readBuffer utils.ReadBuffer, cotpLen uint16) (COTPPacketDisconnectRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("COTPPacketDisconnectRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for COTPPacketDisconnectRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (destinationReference)
	_destinationReference, _destinationReferenceErr := readBuffer.ReadUint16("destinationReference", 16)
	if _destinationReferenceErr != nil {
		return nil, errors.Wrap(_destinationReferenceErr, "Error parsing 'destinationReference' field of COTPPacketDisconnectRequest")
	}
	destinationReference := _destinationReference

	// Simple Field (sourceReference)
	_sourceReference, _sourceReferenceErr := readBuffer.ReadUint16("sourceReference", 16)
	if _sourceReferenceErr != nil {
		return nil, errors.Wrap(_sourceReferenceErr, "Error parsing 'sourceReference' field of COTPPacketDisconnectRequest")
	}
	sourceReference := _sourceReference

	// Simple Field (protocolClass)
	if pullErr := readBuffer.PullContext("protocolClass"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for protocolClass")
	}
	_protocolClass, _protocolClassErr := COTPProtocolClassParse(readBuffer)
	if _protocolClassErr != nil {
		return nil, errors.Wrap(_protocolClassErr, "Error parsing 'protocolClass' field of COTPPacketDisconnectRequest")
	}
	protocolClass := _protocolClass
	if closeErr := readBuffer.CloseContext("protocolClass"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for protocolClass")
	}

	if closeErr := readBuffer.CloseContext("COTPPacketDisconnectRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for COTPPacketDisconnectRequest")
	}

	// Create a partially initialized instance
	_child := &_COTPPacketDisconnectRequest{
		_COTPPacket: &_COTPPacket{
			CotpLen: cotpLen,
		},
		DestinationReference: destinationReference,
		SourceReference:      sourceReference,
		ProtocolClass:        protocolClass,
	}
	_child._COTPPacket._COTPPacketChildRequirements = _child
	return _child, nil
}

func (m *_COTPPacketDisconnectRequest) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("COTPPacketDisconnectRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for COTPPacketDisconnectRequest")
		}

		// Simple Field (destinationReference)
		destinationReference := uint16(m.GetDestinationReference())
		_destinationReferenceErr := writeBuffer.WriteUint16("destinationReference", 16, (destinationReference))
		if _destinationReferenceErr != nil {
			return errors.Wrap(_destinationReferenceErr, "Error serializing 'destinationReference' field")
		}

		// Simple Field (sourceReference)
		sourceReference := uint16(m.GetSourceReference())
		_sourceReferenceErr := writeBuffer.WriteUint16("sourceReference", 16, (sourceReference))
		if _sourceReferenceErr != nil {
			return errors.Wrap(_sourceReferenceErr, "Error serializing 'sourceReference' field")
		}

		// Simple Field (protocolClass)
		if pushErr := writeBuffer.PushContext("protocolClass"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for protocolClass")
		}
		_protocolClassErr := writeBuffer.WriteSerializable(m.GetProtocolClass())
		if popErr := writeBuffer.PopContext("protocolClass"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for protocolClass")
		}
		if _protocolClassErr != nil {
			return errors.Wrap(_protocolClassErr, "Error serializing 'protocolClass' field")
		}

		if popErr := writeBuffer.PopContext("COTPPacketDisconnectRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for COTPPacketDisconnectRequest")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_COTPPacketDisconnectRequest) isCOTPPacketDisconnectRequest() bool {
	return true
}

func (m *_COTPPacketDisconnectRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
