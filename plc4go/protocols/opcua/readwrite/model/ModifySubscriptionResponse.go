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

// ModifySubscriptionResponse is the corresponding interface of ModifySubscriptionResponse
type ModifySubscriptionResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetResponseHeader returns ResponseHeader (property field)
	GetResponseHeader() ExtensionObjectDefinition
	// GetRevisedPublishingInterval returns RevisedPublishingInterval (property field)
	GetRevisedPublishingInterval() float64
	// GetRevisedLifetimeCount returns RevisedLifetimeCount (property field)
	GetRevisedLifetimeCount() uint32
	// GetRevisedMaxKeepAliveCount returns RevisedMaxKeepAliveCount (property field)
	GetRevisedMaxKeepAliveCount() uint32
	// IsModifySubscriptionResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsModifySubscriptionResponse()
}

// _ModifySubscriptionResponse is the data-structure of this message
type _ModifySubscriptionResponse struct {
	ExtensionObjectDefinitionContract
	ResponseHeader            ExtensionObjectDefinition
	RevisedPublishingInterval float64
	RevisedLifetimeCount      uint32
	RevisedMaxKeepAliveCount  uint32
}

var _ ModifySubscriptionResponse = (*_ModifySubscriptionResponse)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_ModifySubscriptionResponse)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ModifySubscriptionResponse) GetIdentifier() string {
	return "796"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ModifySubscriptionResponse) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ModifySubscriptionResponse) GetResponseHeader() ExtensionObjectDefinition {
	return m.ResponseHeader
}

func (m *_ModifySubscriptionResponse) GetRevisedPublishingInterval() float64 {
	return m.RevisedPublishingInterval
}

func (m *_ModifySubscriptionResponse) GetRevisedLifetimeCount() uint32 {
	return m.RevisedLifetimeCount
}

func (m *_ModifySubscriptionResponse) GetRevisedMaxKeepAliveCount() uint32 {
	return m.RevisedMaxKeepAliveCount
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewModifySubscriptionResponse factory function for _ModifySubscriptionResponse
func NewModifySubscriptionResponse(responseHeader ExtensionObjectDefinition, revisedPublishingInterval float64, revisedLifetimeCount uint32, revisedMaxKeepAliveCount uint32) *_ModifySubscriptionResponse {
	if responseHeader == nil {
		panic("responseHeader of type ExtensionObjectDefinition for ModifySubscriptionResponse must not be nil")
	}
	_result := &_ModifySubscriptionResponse{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		ResponseHeader:                    responseHeader,
		RevisedPublishingInterval:         revisedPublishingInterval,
		RevisedLifetimeCount:              revisedLifetimeCount,
		RevisedMaxKeepAliveCount:          revisedMaxKeepAliveCount,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastModifySubscriptionResponse(structType any) ModifySubscriptionResponse {
	if casted, ok := structType.(ModifySubscriptionResponse); ok {
		return casted
	}
	if casted, ok := structType.(*ModifySubscriptionResponse); ok {
		return *casted
	}
	return nil
}

func (m *_ModifySubscriptionResponse) GetTypeName() string {
	return "ModifySubscriptionResponse"
}

func (m *_ModifySubscriptionResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).getLengthInBits(ctx))

	// Simple field (responseHeader)
	lengthInBits += m.ResponseHeader.GetLengthInBits(ctx)

	// Simple field (revisedPublishingInterval)
	lengthInBits += 64

	// Simple field (revisedLifetimeCount)
	lengthInBits += 32

	// Simple field (revisedMaxKeepAliveCount)
	lengthInBits += 32

	return lengthInBits
}

func (m *_ModifySubscriptionResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ModifySubscriptionResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, identifier string) (__modifySubscriptionResponse ModifySubscriptionResponse, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ModifySubscriptionResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ModifySubscriptionResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	responseHeader, err := ReadSimpleField[ExtensionObjectDefinition](ctx, "responseHeader", ReadComplex[ExtensionObjectDefinition](ExtensionObjectDefinitionParseWithBufferProducer[ExtensionObjectDefinition]((string)("394")), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'responseHeader' field"))
	}
	m.ResponseHeader = responseHeader

	revisedPublishingInterval, err := ReadSimpleField(ctx, "revisedPublishingInterval", ReadDouble(readBuffer, uint8(64)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'revisedPublishingInterval' field"))
	}
	m.RevisedPublishingInterval = revisedPublishingInterval

	revisedLifetimeCount, err := ReadSimpleField(ctx, "revisedLifetimeCount", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'revisedLifetimeCount' field"))
	}
	m.RevisedLifetimeCount = revisedLifetimeCount

	revisedMaxKeepAliveCount, err := ReadSimpleField(ctx, "revisedMaxKeepAliveCount", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'revisedMaxKeepAliveCount' field"))
	}
	m.RevisedMaxKeepAliveCount = revisedMaxKeepAliveCount

	if closeErr := readBuffer.CloseContext("ModifySubscriptionResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ModifySubscriptionResponse")
	}

	return m, nil
}

func (m *_ModifySubscriptionResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ModifySubscriptionResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ModifySubscriptionResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ModifySubscriptionResponse")
		}

		if err := WriteSimpleField[ExtensionObjectDefinition](ctx, "responseHeader", m.GetResponseHeader(), WriteComplex[ExtensionObjectDefinition](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'responseHeader' field")
		}

		if err := WriteSimpleField[float64](ctx, "revisedPublishingInterval", m.GetRevisedPublishingInterval(), WriteDouble(writeBuffer, 64)); err != nil {
			return errors.Wrap(err, "Error serializing 'revisedPublishingInterval' field")
		}

		if err := WriteSimpleField[uint32](ctx, "revisedLifetimeCount", m.GetRevisedLifetimeCount(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'revisedLifetimeCount' field")
		}

		if err := WriteSimpleField[uint32](ctx, "revisedMaxKeepAliveCount", m.GetRevisedMaxKeepAliveCount(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'revisedMaxKeepAliveCount' field")
		}

		if popErr := writeBuffer.PopContext("ModifySubscriptionResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ModifySubscriptionResponse")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ModifySubscriptionResponse) IsModifySubscriptionResponse() {}

func (m *_ModifySubscriptionResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
