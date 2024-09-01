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
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// OpcuaAcknowledgeResponse is the corresponding interface of OpcuaAcknowledgeResponse
type OpcuaAcknowledgeResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	MessagePDU
	// GetVersion returns Version (property field)
	GetVersion() uint32
	// GetLimits returns Limits (property field)
	GetLimits() OpcuaProtocolLimits
}

// OpcuaAcknowledgeResponseExactly can be used when we want exactly this type and not a type which fulfills OpcuaAcknowledgeResponse.
// This is useful for switch cases.
type OpcuaAcknowledgeResponseExactly interface {
	OpcuaAcknowledgeResponse
	isOpcuaAcknowledgeResponse() bool
}

// _OpcuaAcknowledgeResponse is the data-structure of this message
type _OpcuaAcknowledgeResponse struct {
	*_MessagePDU
	Version uint32
	Limits  OpcuaProtocolLimits
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_OpcuaAcknowledgeResponse) GetMessageType() string {
	return "ACK"
}

func (m *_OpcuaAcknowledgeResponse) GetResponse() bool {
	return bool(true)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_OpcuaAcknowledgeResponse) InitializeParent(parent MessagePDU, chunk ChunkType) {
	m.Chunk = chunk
}

func (m *_OpcuaAcknowledgeResponse) GetParent() MessagePDU {
	return m._MessagePDU
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_OpcuaAcknowledgeResponse) GetVersion() uint32 {
	return m.Version
}

func (m *_OpcuaAcknowledgeResponse) GetLimits() OpcuaProtocolLimits {
	return m.Limits
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewOpcuaAcknowledgeResponse factory function for _OpcuaAcknowledgeResponse
func NewOpcuaAcknowledgeResponse(version uint32, limits OpcuaProtocolLimits, chunk ChunkType) *_OpcuaAcknowledgeResponse {
	_result := &_OpcuaAcknowledgeResponse{
		Version:     version,
		Limits:      limits,
		_MessagePDU: NewMessagePDU(chunk),
	}
	_result._MessagePDU._MessagePDUChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastOpcuaAcknowledgeResponse(structType any) OpcuaAcknowledgeResponse {
	if casted, ok := structType.(OpcuaAcknowledgeResponse); ok {
		return casted
	}
	if casted, ok := structType.(*OpcuaAcknowledgeResponse); ok {
		return *casted
	}
	return nil
}

func (m *_OpcuaAcknowledgeResponse) GetTypeName() string {
	return "OpcuaAcknowledgeResponse"
}

func (m *_OpcuaAcknowledgeResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (version)
	lengthInBits += 32

	// Simple field (limits)
	lengthInBits += m.Limits.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_OpcuaAcknowledgeResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func OpcuaAcknowledgeResponseParse(ctx context.Context, theBytes []byte, response bool) (OpcuaAcknowledgeResponse, error) {
	return OpcuaAcknowledgeResponseParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), response)
}

func OpcuaAcknowledgeResponseParseWithBufferProducer(response bool) func(ctx context.Context, readBuffer utils.ReadBuffer) (OpcuaAcknowledgeResponse, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (OpcuaAcknowledgeResponse, error) {
		return OpcuaAcknowledgeResponseParseWithBuffer(ctx, readBuffer, response)
	}
}

func OpcuaAcknowledgeResponseParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, response bool) (OpcuaAcknowledgeResponse, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("OpcuaAcknowledgeResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for OpcuaAcknowledgeResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	version, err := ReadSimpleField(ctx, "version", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'version' field"))
	}

	limits, err := ReadSimpleField[OpcuaProtocolLimits](ctx, "limits", ReadComplex[OpcuaProtocolLimits](OpcuaProtocolLimitsParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'limits' field"))
	}

	if closeErr := readBuffer.CloseContext("OpcuaAcknowledgeResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for OpcuaAcknowledgeResponse")
	}

	// Create a partially initialized instance
	_child := &_OpcuaAcknowledgeResponse{
		_MessagePDU: &_MessagePDU{},
		Version:     version,
		Limits:      limits,
	}
	_child._MessagePDU._MessagePDUChildRequirements = _child
	return _child, nil
}

func (m *_OpcuaAcknowledgeResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_OpcuaAcknowledgeResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("OpcuaAcknowledgeResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for OpcuaAcknowledgeResponse")
		}

		if err := WriteSimpleField[uint32](ctx, "version", m.GetVersion(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'version' field")
		}

		if err := WriteSimpleField[OpcuaProtocolLimits](ctx, "limits", m.GetLimits(), WriteComplex[OpcuaProtocolLimits](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'limits' field")
		}

		if popErr := writeBuffer.PopContext("OpcuaAcknowledgeResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for OpcuaAcknowledgeResponse")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_OpcuaAcknowledgeResponse) isOpcuaAcknowledgeResponse() bool {
	return true
}

func (m *_OpcuaAcknowledgeResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
