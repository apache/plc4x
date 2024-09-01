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

// SetMonitoringModeRequest is the corresponding interface of SetMonitoringModeRequest
type SetMonitoringModeRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetRequestHeader returns RequestHeader (property field)
	GetRequestHeader() ExtensionObjectDefinition
	// GetSubscriptionId returns SubscriptionId (property field)
	GetSubscriptionId() uint32
	// GetMonitoringMode returns MonitoringMode (property field)
	GetMonitoringMode() MonitoringMode
	// GetNoOfMonitoredItemIds returns NoOfMonitoredItemIds (property field)
	GetNoOfMonitoredItemIds() int32
	// GetMonitoredItemIds returns MonitoredItemIds (property field)
	GetMonitoredItemIds() []uint32
}

// SetMonitoringModeRequestExactly can be used when we want exactly this type and not a type which fulfills SetMonitoringModeRequest.
// This is useful for switch cases.
type SetMonitoringModeRequestExactly interface {
	SetMonitoringModeRequest
	isSetMonitoringModeRequest() bool
}

// _SetMonitoringModeRequest is the data-structure of this message
type _SetMonitoringModeRequest struct {
	*_ExtensionObjectDefinition
	RequestHeader        ExtensionObjectDefinition
	SubscriptionId       uint32
	MonitoringMode       MonitoringMode
	NoOfMonitoredItemIds int32
	MonitoredItemIds     []uint32
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_SetMonitoringModeRequest) GetIdentifier() string {
	return "769"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_SetMonitoringModeRequest) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_SetMonitoringModeRequest) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_SetMonitoringModeRequest) GetRequestHeader() ExtensionObjectDefinition {
	return m.RequestHeader
}

func (m *_SetMonitoringModeRequest) GetSubscriptionId() uint32 {
	return m.SubscriptionId
}

func (m *_SetMonitoringModeRequest) GetMonitoringMode() MonitoringMode {
	return m.MonitoringMode
}

func (m *_SetMonitoringModeRequest) GetNoOfMonitoredItemIds() int32 {
	return m.NoOfMonitoredItemIds
}

func (m *_SetMonitoringModeRequest) GetMonitoredItemIds() []uint32 {
	return m.MonitoredItemIds
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewSetMonitoringModeRequest factory function for _SetMonitoringModeRequest
func NewSetMonitoringModeRequest(requestHeader ExtensionObjectDefinition, subscriptionId uint32, monitoringMode MonitoringMode, noOfMonitoredItemIds int32, monitoredItemIds []uint32) *_SetMonitoringModeRequest {
	_result := &_SetMonitoringModeRequest{
		RequestHeader:              requestHeader,
		SubscriptionId:             subscriptionId,
		MonitoringMode:             monitoringMode,
		NoOfMonitoredItemIds:       noOfMonitoredItemIds,
		MonitoredItemIds:           monitoredItemIds,
		_ExtensionObjectDefinition: NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastSetMonitoringModeRequest(structType any) SetMonitoringModeRequest {
	if casted, ok := structType.(SetMonitoringModeRequest); ok {
		return casted
	}
	if casted, ok := structType.(*SetMonitoringModeRequest); ok {
		return *casted
	}
	return nil
}

func (m *_SetMonitoringModeRequest) GetTypeName() string {
	return "SetMonitoringModeRequest"
}

func (m *_SetMonitoringModeRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (requestHeader)
	lengthInBits += m.RequestHeader.GetLengthInBits(ctx)

	// Simple field (subscriptionId)
	lengthInBits += 32

	// Simple field (monitoringMode)
	lengthInBits += 32

	// Simple field (noOfMonitoredItemIds)
	lengthInBits += 32

	// Array field
	if len(m.MonitoredItemIds) > 0 {
		lengthInBits += 32 * uint16(len(m.MonitoredItemIds))
	}

	return lengthInBits
}

func (m *_SetMonitoringModeRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func SetMonitoringModeRequestParse(ctx context.Context, theBytes []byte, identifier string) (SetMonitoringModeRequest, error) {
	return SetMonitoringModeRequestParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func SetMonitoringModeRequestParseWithBufferProducer(identifier string) func(ctx context.Context, readBuffer utils.ReadBuffer) (SetMonitoringModeRequest, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (SetMonitoringModeRequest, error) {
		return SetMonitoringModeRequestParseWithBuffer(ctx, readBuffer, identifier)
	}
}

func SetMonitoringModeRequestParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (SetMonitoringModeRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SetMonitoringModeRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SetMonitoringModeRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	requestHeader, err := ReadSimpleField[ExtensionObjectDefinition](ctx, "requestHeader", ReadComplex[ExtensionObjectDefinition](ExtensionObjectDefinitionParseWithBufferProducer[ExtensionObjectDefinition]((string)("391")), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'requestHeader' field"))
	}

	subscriptionId, err := ReadSimpleField(ctx, "subscriptionId", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'subscriptionId' field"))
	}

	monitoringMode, err := ReadEnumField[MonitoringMode](ctx, "monitoringMode", "MonitoringMode", ReadEnum(MonitoringModeByValue, ReadUnsignedInt(readBuffer, uint8(32))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'monitoringMode' field"))
	}

	noOfMonitoredItemIds, err := ReadSimpleField(ctx, "noOfMonitoredItemIds", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfMonitoredItemIds' field"))
	}

	monitoredItemIds, err := ReadCountArrayField[uint32](ctx, "monitoredItemIds", ReadUnsignedInt(readBuffer, uint8(32)), uint64(noOfMonitoredItemIds))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'monitoredItemIds' field"))
	}

	if closeErr := readBuffer.CloseContext("SetMonitoringModeRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SetMonitoringModeRequest")
	}

	// Create a partially initialized instance
	_child := &_SetMonitoringModeRequest{
		_ExtensionObjectDefinition: &_ExtensionObjectDefinition{},
		RequestHeader:              requestHeader,
		SubscriptionId:             subscriptionId,
		MonitoringMode:             monitoringMode,
		NoOfMonitoredItemIds:       noOfMonitoredItemIds,
		MonitoredItemIds:           monitoredItemIds,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_SetMonitoringModeRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_SetMonitoringModeRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("SetMonitoringModeRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for SetMonitoringModeRequest")
		}

		if err := WriteSimpleField[ExtensionObjectDefinition](ctx, "requestHeader", m.GetRequestHeader(), WriteComplex[ExtensionObjectDefinition](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'requestHeader' field")
		}

		if err := WriteSimpleField[uint32](ctx, "subscriptionId", m.GetSubscriptionId(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'subscriptionId' field")
		}

		if err := WriteSimpleEnumField[MonitoringMode](ctx, "monitoringMode", "MonitoringMode", m.GetMonitoringMode(), WriteEnum[MonitoringMode, uint32](MonitoringMode.GetValue, MonitoringMode.PLC4XEnumName, WriteUnsignedInt(writeBuffer, 32))); err != nil {
			return errors.Wrap(err, "Error serializing 'monitoringMode' field")
		}

		if err := WriteSimpleField[int32](ctx, "noOfMonitoredItemIds", m.GetNoOfMonitoredItemIds(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfMonitoredItemIds' field")
		}

		if err := WriteSimpleTypeArrayField(ctx, "monitoredItemIds", m.GetMonitoredItemIds(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'monitoredItemIds' field")
		}

		if popErr := writeBuffer.PopContext("SetMonitoringModeRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for SetMonitoringModeRequest")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_SetMonitoringModeRequest) isSetMonitoringModeRequest() bool {
	return true
}

func (m *_SetMonitoringModeRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
