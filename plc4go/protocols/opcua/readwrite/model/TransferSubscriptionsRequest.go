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

// TransferSubscriptionsRequest is the corresponding interface of TransferSubscriptionsRequest
type TransferSubscriptionsRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetRequestHeader returns RequestHeader (property field)
	GetRequestHeader() ExtensionObjectDefinition
	// GetNoOfSubscriptionIds returns NoOfSubscriptionIds (property field)
	GetNoOfSubscriptionIds() int32
	// GetSubscriptionIds returns SubscriptionIds (property field)
	GetSubscriptionIds() []uint32
	// GetSendInitialValues returns SendInitialValues (property field)
	GetSendInitialValues() bool
}

// TransferSubscriptionsRequestExactly can be used when we want exactly this type and not a type which fulfills TransferSubscriptionsRequest.
// This is useful for switch cases.
type TransferSubscriptionsRequestExactly interface {
	TransferSubscriptionsRequest
	isTransferSubscriptionsRequest() bool
}

// _TransferSubscriptionsRequest is the data-structure of this message
type _TransferSubscriptionsRequest struct {
	*_ExtensionObjectDefinition
	RequestHeader       ExtensionObjectDefinition
	NoOfSubscriptionIds int32
	SubscriptionIds     []uint32
	SendInitialValues   bool
	// Reserved Fields
	reservedField0 *uint8
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_TransferSubscriptionsRequest) GetIdentifier() string {
	return "841"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_TransferSubscriptionsRequest) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_TransferSubscriptionsRequest) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_TransferSubscriptionsRequest) GetRequestHeader() ExtensionObjectDefinition {
	return m.RequestHeader
}

func (m *_TransferSubscriptionsRequest) GetNoOfSubscriptionIds() int32 {
	return m.NoOfSubscriptionIds
}

func (m *_TransferSubscriptionsRequest) GetSubscriptionIds() []uint32 {
	return m.SubscriptionIds
}

func (m *_TransferSubscriptionsRequest) GetSendInitialValues() bool {
	return m.SendInitialValues
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewTransferSubscriptionsRequest factory function for _TransferSubscriptionsRequest
func NewTransferSubscriptionsRequest(requestHeader ExtensionObjectDefinition, noOfSubscriptionIds int32, subscriptionIds []uint32, sendInitialValues bool) *_TransferSubscriptionsRequest {
	_result := &_TransferSubscriptionsRequest{
		RequestHeader:              requestHeader,
		NoOfSubscriptionIds:        noOfSubscriptionIds,
		SubscriptionIds:            subscriptionIds,
		SendInitialValues:          sendInitialValues,
		_ExtensionObjectDefinition: NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastTransferSubscriptionsRequest(structType any) TransferSubscriptionsRequest {
	if casted, ok := structType.(TransferSubscriptionsRequest); ok {
		return casted
	}
	if casted, ok := structType.(*TransferSubscriptionsRequest); ok {
		return *casted
	}
	return nil
}

func (m *_TransferSubscriptionsRequest) GetTypeName() string {
	return "TransferSubscriptionsRequest"
}

func (m *_TransferSubscriptionsRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (requestHeader)
	lengthInBits += m.RequestHeader.GetLengthInBits(ctx)

	// Simple field (noOfSubscriptionIds)
	lengthInBits += 32

	// Array field
	if len(m.SubscriptionIds) > 0 {
		lengthInBits += 32 * uint16(len(m.SubscriptionIds))
	}

	// Reserved Field (reserved)
	lengthInBits += 7

	// Simple field (sendInitialValues)
	lengthInBits += 1

	return lengthInBits
}

func (m *_TransferSubscriptionsRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func TransferSubscriptionsRequestParse(ctx context.Context, theBytes []byte, identifier string) (TransferSubscriptionsRequest, error) {
	return TransferSubscriptionsRequestParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func TransferSubscriptionsRequestParseWithBufferProducer(identifier string) func(ctx context.Context, readBuffer utils.ReadBuffer) (TransferSubscriptionsRequest, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (TransferSubscriptionsRequest, error) {
		return TransferSubscriptionsRequestParseWithBuffer(ctx, readBuffer, identifier)
	}
}

func TransferSubscriptionsRequestParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (TransferSubscriptionsRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("TransferSubscriptionsRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for TransferSubscriptionsRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	requestHeader, err := ReadSimpleField[ExtensionObjectDefinition](ctx, "requestHeader", ReadComplex[ExtensionObjectDefinition](ExtensionObjectDefinitionParseWithBufferProducer[ExtensionObjectDefinition]((string)("391")), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'requestHeader' field"))
	}

	noOfSubscriptionIds, err := ReadSimpleField(ctx, "noOfSubscriptionIds", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfSubscriptionIds' field"))
	}

	subscriptionIds, err := ReadCountArrayField[uint32](ctx, "subscriptionIds", ReadUnsignedInt(readBuffer, uint8(32)), uint64(noOfSubscriptionIds))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'subscriptionIds' field"))
	}

	reservedField0, err := ReadReservedField(ctx, "reserved", ReadUnsignedByte(readBuffer, uint8(7)), uint8(0x00))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}

	sendInitialValues, err := ReadSimpleField(ctx, "sendInitialValues", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'sendInitialValues' field"))
	}

	if closeErr := readBuffer.CloseContext("TransferSubscriptionsRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for TransferSubscriptionsRequest")
	}

	// Create a partially initialized instance
	_child := &_TransferSubscriptionsRequest{
		_ExtensionObjectDefinition: &_ExtensionObjectDefinition{},
		RequestHeader:              requestHeader,
		NoOfSubscriptionIds:        noOfSubscriptionIds,
		SubscriptionIds:            subscriptionIds,
		SendInitialValues:          sendInitialValues,
		reservedField0:             reservedField0,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_TransferSubscriptionsRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_TransferSubscriptionsRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("TransferSubscriptionsRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for TransferSubscriptionsRequest")
		}

		if err := WriteSimpleField[ExtensionObjectDefinition](ctx, "requestHeader", m.GetRequestHeader(), WriteComplex[ExtensionObjectDefinition](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'requestHeader' field")
		}

		if err := WriteSimpleField[int32](ctx, "noOfSubscriptionIds", m.GetNoOfSubscriptionIds(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfSubscriptionIds' field")
		}

		if err := WriteSimpleTypeArrayField(ctx, "subscriptionIds", m.GetSubscriptionIds(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'subscriptionIds' field")
		}

		if err := WriteReservedField[uint8](ctx, "reserved", uint8(0x00), WriteUnsignedByte(writeBuffer, 7)); err != nil {
			return errors.Wrap(err, "Error serializing 'reserved' field number 1")
		}

		if err := WriteSimpleField[bool](ctx, "sendInitialValues", m.GetSendInitialValues(), WriteBoolean(writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'sendInitialValues' field")
		}

		if popErr := writeBuffer.PopContext("TransferSubscriptionsRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for TransferSubscriptionsRequest")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_TransferSubscriptionsRequest) isTransferSubscriptionsRequest() bool {
	return true
}

func (m *_TransferSubscriptionsRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
