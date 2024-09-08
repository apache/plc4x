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

// CreateSubscriptionResponse is the corresponding interface of CreateSubscriptionResponse
type CreateSubscriptionResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetResponseHeader returns ResponseHeader (property field)
	GetResponseHeader() ExtensionObjectDefinition
	// GetSubscriptionId returns SubscriptionId (property field)
	GetSubscriptionId() uint32
	// GetRevisedPublishingInterval returns RevisedPublishingInterval (property field)
	GetRevisedPublishingInterval() float64
	// GetRevisedLifetimeCount returns RevisedLifetimeCount (property field)
	GetRevisedLifetimeCount() uint32
	// GetRevisedMaxKeepAliveCount returns RevisedMaxKeepAliveCount (property field)
	GetRevisedMaxKeepAliveCount() uint32
	// IsCreateSubscriptionResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCreateSubscriptionResponse()
}

// _CreateSubscriptionResponse is the data-structure of this message
type _CreateSubscriptionResponse struct {
	ExtensionObjectDefinitionContract
	ResponseHeader            ExtensionObjectDefinition
	SubscriptionId            uint32
	RevisedPublishingInterval float64
	RevisedLifetimeCount      uint32
	RevisedMaxKeepAliveCount  uint32
}

var _ CreateSubscriptionResponse = (*_CreateSubscriptionResponse)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_CreateSubscriptionResponse)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_CreateSubscriptionResponse) GetIdentifier() string {
	return "790"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_CreateSubscriptionResponse) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CreateSubscriptionResponse) GetResponseHeader() ExtensionObjectDefinition {
	return m.ResponseHeader
}

func (m *_CreateSubscriptionResponse) GetSubscriptionId() uint32 {
	return m.SubscriptionId
}

func (m *_CreateSubscriptionResponse) GetRevisedPublishingInterval() float64 {
	return m.RevisedPublishingInterval
}

func (m *_CreateSubscriptionResponse) GetRevisedLifetimeCount() uint32 {
	return m.RevisedLifetimeCount
}

func (m *_CreateSubscriptionResponse) GetRevisedMaxKeepAliveCount() uint32 {
	return m.RevisedMaxKeepAliveCount
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewCreateSubscriptionResponse factory function for _CreateSubscriptionResponse
func NewCreateSubscriptionResponse(responseHeader ExtensionObjectDefinition, subscriptionId uint32, revisedPublishingInterval float64, revisedLifetimeCount uint32, revisedMaxKeepAliveCount uint32) *_CreateSubscriptionResponse {
	if responseHeader == nil {
		panic("responseHeader of type ExtensionObjectDefinition for CreateSubscriptionResponse must not be nil")
	}
	_result := &_CreateSubscriptionResponse{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		ResponseHeader:                    responseHeader,
		SubscriptionId:                    subscriptionId,
		RevisedPublishingInterval:         revisedPublishingInterval,
		RevisedLifetimeCount:              revisedLifetimeCount,
		RevisedMaxKeepAliveCount:          revisedMaxKeepAliveCount,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastCreateSubscriptionResponse(structType any) CreateSubscriptionResponse {
	if casted, ok := structType.(CreateSubscriptionResponse); ok {
		return casted
	}
	if casted, ok := structType.(*CreateSubscriptionResponse); ok {
		return *casted
	}
	return nil
}

func (m *_CreateSubscriptionResponse) GetTypeName() string {
	return "CreateSubscriptionResponse"
}

func (m *_CreateSubscriptionResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).getLengthInBits(ctx))

	// Simple field (responseHeader)
	lengthInBits += m.ResponseHeader.GetLengthInBits(ctx)

	// Simple field (subscriptionId)
	lengthInBits += 32

	// Simple field (revisedPublishingInterval)
	lengthInBits += 64

	// Simple field (revisedLifetimeCount)
	lengthInBits += 32

	// Simple field (revisedMaxKeepAliveCount)
	lengthInBits += 32

	return lengthInBits
}

func (m *_CreateSubscriptionResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_CreateSubscriptionResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, identifier string) (__createSubscriptionResponse CreateSubscriptionResponse, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CreateSubscriptionResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CreateSubscriptionResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	responseHeader, err := ReadSimpleField[ExtensionObjectDefinition](ctx, "responseHeader", ReadComplex[ExtensionObjectDefinition](ExtensionObjectDefinitionParseWithBufferProducer[ExtensionObjectDefinition]((string)("394")), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'responseHeader' field"))
	}
	m.ResponseHeader = responseHeader

	subscriptionId, err := ReadSimpleField(ctx, "subscriptionId", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'subscriptionId' field"))
	}
	m.SubscriptionId = subscriptionId

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

	if closeErr := readBuffer.CloseContext("CreateSubscriptionResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CreateSubscriptionResponse")
	}

	return m, nil
}

func (m *_CreateSubscriptionResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CreateSubscriptionResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CreateSubscriptionResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CreateSubscriptionResponse")
		}

		if err := WriteSimpleField[ExtensionObjectDefinition](ctx, "responseHeader", m.GetResponseHeader(), WriteComplex[ExtensionObjectDefinition](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'responseHeader' field")
		}

		if err := WriteSimpleField[uint32](ctx, "subscriptionId", m.GetSubscriptionId(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'subscriptionId' field")
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

		if popErr := writeBuffer.PopContext("CreateSubscriptionResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CreateSubscriptionResponse")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_CreateSubscriptionResponse) IsCreateSubscriptionResponse() {}

func (m *_CreateSubscriptionResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
