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
	utils.Copyable
	ExtensionObjectDefinition
	// GetRequestHeader returns RequestHeader (property field)
	GetRequestHeader() RequestHeader
	// GetSubscriptionId returns SubscriptionId (property field)
	GetSubscriptionId() uint32
	// GetMonitoringMode returns MonitoringMode (property field)
	GetMonitoringMode() MonitoringMode
	// GetMonitoredItemIds returns MonitoredItemIds (property field)
	GetMonitoredItemIds() []uint32
	// IsSetMonitoringModeRequest is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsSetMonitoringModeRequest()
	// CreateBuilder creates a SetMonitoringModeRequestBuilder
	CreateSetMonitoringModeRequestBuilder() SetMonitoringModeRequestBuilder
}

// _SetMonitoringModeRequest is the data-structure of this message
type _SetMonitoringModeRequest struct {
	ExtensionObjectDefinitionContract
	RequestHeader    RequestHeader
	SubscriptionId   uint32
	MonitoringMode   MonitoringMode
	MonitoredItemIds []uint32
}

var _ SetMonitoringModeRequest = (*_SetMonitoringModeRequest)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_SetMonitoringModeRequest)(nil)

// NewSetMonitoringModeRequest factory function for _SetMonitoringModeRequest
func NewSetMonitoringModeRequest(requestHeader RequestHeader, subscriptionId uint32, monitoringMode MonitoringMode, monitoredItemIds []uint32) *_SetMonitoringModeRequest {
	if requestHeader == nil {
		panic("requestHeader of type RequestHeader for SetMonitoringModeRequest must not be nil")
	}
	_result := &_SetMonitoringModeRequest{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		RequestHeader:                     requestHeader,
		SubscriptionId:                    subscriptionId,
		MonitoringMode:                    monitoringMode,
		MonitoredItemIds:                  monitoredItemIds,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// SetMonitoringModeRequestBuilder is a builder for SetMonitoringModeRequest
type SetMonitoringModeRequestBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(requestHeader RequestHeader, subscriptionId uint32, monitoringMode MonitoringMode, monitoredItemIds []uint32) SetMonitoringModeRequestBuilder
	// WithRequestHeader adds RequestHeader (property field)
	WithRequestHeader(RequestHeader) SetMonitoringModeRequestBuilder
	// WithRequestHeaderBuilder adds RequestHeader (property field) which is build by the builder
	WithRequestHeaderBuilder(func(RequestHeaderBuilder) RequestHeaderBuilder) SetMonitoringModeRequestBuilder
	// WithSubscriptionId adds SubscriptionId (property field)
	WithSubscriptionId(uint32) SetMonitoringModeRequestBuilder
	// WithMonitoringMode adds MonitoringMode (property field)
	WithMonitoringMode(MonitoringMode) SetMonitoringModeRequestBuilder
	// WithMonitoredItemIds adds MonitoredItemIds (property field)
	WithMonitoredItemIds(...uint32) SetMonitoringModeRequestBuilder
	// Build builds the SetMonitoringModeRequest or returns an error if something is wrong
	Build() (SetMonitoringModeRequest, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() SetMonitoringModeRequest
}

// NewSetMonitoringModeRequestBuilder() creates a SetMonitoringModeRequestBuilder
func NewSetMonitoringModeRequestBuilder() SetMonitoringModeRequestBuilder {
	return &_SetMonitoringModeRequestBuilder{_SetMonitoringModeRequest: new(_SetMonitoringModeRequest)}
}

type _SetMonitoringModeRequestBuilder struct {
	*_SetMonitoringModeRequest

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (SetMonitoringModeRequestBuilder) = (*_SetMonitoringModeRequestBuilder)(nil)

func (b *_SetMonitoringModeRequestBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_SetMonitoringModeRequestBuilder) WithMandatoryFields(requestHeader RequestHeader, subscriptionId uint32, monitoringMode MonitoringMode, monitoredItemIds []uint32) SetMonitoringModeRequestBuilder {
	return b.WithRequestHeader(requestHeader).WithSubscriptionId(subscriptionId).WithMonitoringMode(monitoringMode).WithMonitoredItemIds(monitoredItemIds...)
}

func (b *_SetMonitoringModeRequestBuilder) WithRequestHeader(requestHeader RequestHeader) SetMonitoringModeRequestBuilder {
	b.RequestHeader = requestHeader
	return b
}

func (b *_SetMonitoringModeRequestBuilder) WithRequestHeaderBuilder(builderSupplier func(RequestHeaderBuilder) RequestHeaderBuilder) SetMonitoringModeRequestBuilder {
	builder := builderSupplier(b.RequestHeader.CreateRequestHeaderBuilder())
	var err error
	b.RequestHeader, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "RequestHeaderBuilder failed"))
	}
	return b
}

func (b *_SetMonitoringModeRequestBuilder) WithSubscriptionId(subscriptionId uint32) SetMonitoringModeRequestBuilder {
	b.SubscriptionId = subscriptionId
	return b
}

func (b *_SetMonitoringModeRequestBuilder) WithMonitoringMode(monitoringMode MonitoringMode) SetMonitoringModeRequestBuilder {
	b.MonitoringMode = monitoringMode
	return b
}

func (b *_SetMonitoringModeRequestBuilder) WithMonitoredItemIds(monitoredItemIds ...uint32) SetMonitoringModeRequestBuilder {
	b.MonitoredItemIds = monitoredItemIds
	return b
}

func (b *_SetMonitoringModeRequestBuilder) Build() (SetMonitoringModeRequest, error) {
	if b.RequestHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'requestHeader' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._SetMonitoringModeRequest.deepCopy(), nil
}

func (b *_SetMonitoringModeRequestBuilder) MustBuild() SetMonitoringModeRequest {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_SetMonitoringModeRequestBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_SetMonitoringModeRequestBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_SetMonitoringModeRequestBuilder) DeepCopy() any {
	_copy := b.CreateSetMonitoringModeRequestBuilder().(*_SetMonitoringModeRequestBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateSetMonitoringModeRequestBuilder creates a SetMonitoringModeRequestBuilder
func (b *_SetMonitoringModeRequest) CreateSetMonitoringModeRequestBuilder() SetMonitoringModeRequestBuilder {
	if b == nil {
		return NewSetMonitoringModeRequestBuilder()
	}
	return &_SetMonitoringModeRequestBuilder{_SetMonitoringModeRequest: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_SetMonitoringModeRequest) GetExtensionId() int32 {
	return int32(769)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_SetMonitoringModeRequest) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_SetMonitoringModeRequest) GetRequestHeader() RequestHeader {
	return m.RequestHeader
}

func (m *_SetMonitoringModeRequest) GetSubscriptionId() uint32 {
	return m.SubscriptionId
}

func (m *_SetMonitoringModeRequest) GetMonitoringMode() MonitoringMode {
	return m.MonitoringMode
}

func (m *_SetMonitoringModeRequest) GetMonitoredItemIds() []uint32 {
	return m.MonitoredItemIds
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

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
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (requestHeader)
	lengthInBits += m.RequestHeader.GetLengthInBits(ctx)

	// Simple field (subscriptionId)
	lengthInBits += 32

	// Simple field (monitoringMode)
	lengthInBits += 32

	// Implicit Field (noOfMonitoredItemIds)
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

func (m *_SetMonitoringModeRequest) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__setMonitoringModeRequest SetMonitoringModeRequest, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SetMonitoringModeRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SetMonitoringModeRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	requestHeader, err := ReadSimpleField[RequestHeader](ctx, "requestHeader", ReadComplex[RequestHeader](ExtensionObjectDefinitionParseWithBufferProducer[RequestHeader]((int32)(int32(391))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'requestHeader' field"))
	}
	m.RequestHeader = requestHeader

	subscriptionId, err := ReadSimpleField(ctx, "subscriptionId", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'subscriptionId' field"))
	}
	m.SubscriptionId = subscriptionId

	monitoringMode, err := ReadEnumField[MonitoringMode](ctx, "monitoringMode", "MonitoringMode", ReadEnum(MonitoringModeByValue, ReadUnsignedInt(readBuffer, uint8(32))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'monitoringMode' field"))
	}
	m.MonitoringMode = monitoringMode

	noOfMonitoredItemIds, err := ReadImplicitField[int32](ctx, "noOfMonitoredItemIds", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfMonitoredItemIds' field"))
	}
	_ = noOfMonitoredItemIds

	monitoredItemIds, err := ReadCountArrayField[uint32](ctx, "monitoredItemIds", ReadUnsignedInt(readBuffer, uint8(32)), uint64(noOfMonitoredItemIds))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'monitoredItemIds' field"))
	}
	m.MonitoredItemIds = monitoredItemIds

	if closeErr := readBuffer.CloseContext("SetMonitoringModeRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SetMonitoringModeRequest")
	}

	return m, nil
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

		if err := WriteSimpleField[RequestHeader](ctx, "requestHeader", m.GetRequestHeader(), WriteComplex[RequestHeader](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'requestHeader' field")
		}

		if err := WriteSimpleField[uint32](ctx, "subscriptionId", m.GetSubscriptionId(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'subscriptionId' field")
		}

		if err := WriteSimpleEnumField[MonitoringMode](ctx, "monitoringMode", "MonitoringMode", m.GetMonitoringMode(), WriteEnum[MonitoringMode, uint32](MonitoringMode.GetValue, MonitoringMode.PLC4XEnumName, WriteUnsignedInt(writeBuffer, 32))); err != nil {
			return errors.Wrap(err, "Error serializing 'monitoringMode' field")
		}
		noOfMonitoredItemIds := int32(utils.InlineIf(bool((m.GetMonitoredItemIds()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetMonitoredItemIds()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfMonitoredItemIds", noOfMonitoredItemIds, WriteSignedInt(writeBuffer, 32)); err != nil {
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
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_SetMonitoringModeRequest) IsSetMonitoringModeRequest() {}

func (m *_SetMonitoringModeRequest) DeepCopy() any {
	return m.deepCopy()
}

func (m *_SetMonitoringModeRequest) deepCopy() *_SetMonitoringModeRequest {
	if m == nil {
		return nil
	}
	_SetMonitoringModeRequestCopy := &_SetMonitoringModeRequest{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.RequestHeader.DeepCopy().(RequestHeader),
		m.SubscriptionId,
		m.MonitoringMode,
		utils.DeepCopySlice[uint32, uint32](m.MonitoredItemIds),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _SetMonitoringModeRequestCopy
}

func (m *_SetMonitoringModeRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	wb := utils.NewWriteBufferBoxBased(
		utils.WithWriteBufferBoxBasedMergeSingleBoxes(),
		utils.WithWriteBufferBoxBasedOmitEmptyBoxes(),
		utils.WithWriteBufferBoxBasedPrintPosLengthFooter(),
	)
	if err := wb.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return wb.GetBox().String()
}
