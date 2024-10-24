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
	utils.Copyable
	ExtensionObjectDefinition
	// GetResponseHeader returns ResponseHeader (property field)
	GetResponseHeader() ResponseHeader
	// GetRevisedPublishingInterval returns RevisedPublishingInterval (property field)
	GetRevisedPublishingInterval() float64
	// GetRevisedLifetimeCount returns RevisedLifetimeCount (property field)
	GetRevisedLifetimeCount() uint32
	// GetRevisedMaxKeepAliveCount returns RevisedMaxKeepAliveCount (property field)
	GetRevisedMaxKeepAliveCount() uint32
	// IsModifySubscriptionResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsModifySubscriptionResponse()
	// CreateBuilder creates a ModifySubscriptionResponseBuilder
	CreateModifySubscriptionResponseBuilder() ModifySubscriptionResponseBuilder
}

// _ModifySubscriptionResponse is the data-structure of this message
type _ModifySubscriptionResponse struct {
	ExtensionObjectDefinitionContract
	ResponseHeader            ResponseHeader
	RevisedPublishingInterval float64
	RevisedLifetimeCount      uint32
	RevisedMaxKeepAliveCount  uint32
}

var _ ModifySubscriptionResponse = (*_ModifySubscriptionResponse)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_ModifySubscriptionResponse)(nil)

// NewModifySubscriptionResponse factory function for _ModifySubscriptionResponse
func NewModifySubscriptionResponse(responseHeader ResponseHeader, revisedPublishingInterval float64, revisedLifetimeCount uint32, revisedMaxKeepAliveCount uint32) *_ModifySubscriptionResponse {
	if responseHeader == nil {
		panic("responseHeader of type ResponseHeader for ModifySubscriptionResponse must not be nil")
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

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ModifySubscriptionResponseBuilder is a builder for ModifySubscriptionResponse
type ModifySubscriptionResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(responseHeader ResponseHeader, revisedPublishingInterval float64, revisedLifetimeCount uint32, revisedMaxKeepAliveCount uint32) ModifySubscriptionResponseBuilder
	// WithResponseHeader adds ResponseHeader (property field)
	WithResponseHeader(ResponseHeader) ModifySubscriptionResponseBuilder
	// WithResponseHeaderBuilder adds ResponseHeader (property field) which is build by the builder
	WithResponseHeaderBuilder(func(ResponseHeaderBuilder) ResponseHeaderBuilder) ModifySubscriptionResponseBuilder
	// WithRevisedPublishingInterval adds RevisedPublishingInterval (property field)
	WithRevisedPublishingInterval(float64) ModifySubscriptionResponseBuilder
	// WithRevisedLifetimeCount adds RevisedLifetimeCount (property field)
	WithRevisedLifetimeCount(uint32) ModifySubscriptionResponseBuilder
	// WithRevisedMaxKeepAliveCount adds RevisedMaxKeepAliveCount (property field)
	WithRevisedMaxKeepAliveCount(uint32) ModifySubscriptionResponseBuilder
	// Build builds the ModifySubscriptionResponse or returns an error if something is wrong
	Build() (ModifySubscriptionResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ModifySubscriptionResponse
}

// NewModifySubscriptionResponseBuilder() creates a ModifySubscriptionResponseBuilder
func NewModifySubscriptionResponseBuilder() ModifySubscriptionResponseBuilder {
	return &_ModifySubscriptionResponseBuilder{_ModifySubscriptionResponse: new(_ModifySubscriptionResponse)}
}

type _ModifySubscriptionResponseBuilder struct {
	*_ModifySubscriptionResponse

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (ModifySubscriptionResponseBuilder) = (*_ModifySubscriptionResponseBuilder)(nil)

func (b *_ModifySubscriptionResponseBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_ModifySubscriptionResponseBuilder) WithMandatoryFields(responseHeader ResponseHeader, revisedPublishingInterval float64, revisedLifetimeCount uint32, revisedMaxKeepAliveCount uint32) ModifySubscriptionResponseBuilder {
	return b.WithResponseHeader(responseHeader).WithRevisedPublishingInterval(revisedPublishingInterval).WithRevisedLifetimeCount(revisedLifetimeCount).WithRevisedMaxKeepAliveCount(revisedMaxKeepAliveCount)
}

func (b *_ModifySubscriptionResponseBuilder) WithResponseHeader(responseHeader ResponseHeader) ModifySubscriptionResponseBuilder {
	b.ResponseHeader = responseHeader
	return b
}

func (b *_ModifySubscriptionResponseBuilder) WithResponseHeaderBuilder(builderSupplier func(ResponseHeaderBuilder) ResponseHeaderBuilder) ModifySubscriptionResponseBuilder {
	builder := builderSupplier(b.ResponseHeader.CreateResponseHeaderBuilder())
	var err error
	b.ResponseHeader, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "ResponseHeaderBuilder failed"))
	}
	return b
}

func (b *_ModifySubscriptionResponseBuilder) WithRevisedPublishingInterval(revisedPublishingInterval float64) ModifySubscriptionResponseBuilder {
	b.RevisedPublishingInterval = revisedPublishingInterval
	return b
}

func (b *_ModifySubscriptionResponseBuilder) WithRevisedLifetimeCount(revisedLifetimeCount uint32) ModifySubscriptionResponseBuilder {
	b.RevisedLifetimeCount = revisedLifetimeCount
	return b
}

func (b *_ModifySubscriptionResponseBuilder) WithRevisedMaxKeepAliveCount(revisedMaxKeepAliveCount uint32) ModifySubscriptionResponseBuilder {
	b.RevisedMaxKeepAliveCount = revisedMaxKeepAliveCount
	return b
}

func (b *_ModifySubscriptionResponseBuilder) Build() (ModifySubscriptionResponse, error) {
	if b.ResponseHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'responseHeader' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ModifySubscriptionResponse.deepCopy(), nil
}

func (b *_ModifySubscriptionResponseBuilder) MustBuild() ModifySubscriptionResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ModifySubscriptionResponseBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_ModifySubscriptionResponseBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_ModifySubscriptionResponseBuilder) DeepCopy() any {
	_copy := b.CreateModifySubscriptionResponseBuilder().(*_ModifySubscriptionResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateModifySubscriptionResponseBuilder creates a ModifySubscriptionResponseBuilder
func (b *_ModifySubscriptionResponse) CreateModifySubscriptionResponseBuilder() ModifySubscriptionResponseBuilder {
	if b == nil {
		return NewModifySubscriptionResponseBuilder()
	}
	return &_ModifySubscriptionResponseBuilder{_ModifySubscriptionResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ModifySubscriptionResponse) GetExtensionId() int32 {
	return int32(796)
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

func (m *_ModifySubscriptionResponse) GetResponseHeader() ResponseHeader {
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
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

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

func (m *_ModifySubscriptionResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__modifySubscriptionResponse ModifySubscriptionResponse, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ModifySubscriptionResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ModifySubscriptionResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	responseHeader, err := ReadSimpleField[ResponseHeader](ctx, "responseHeader", ReadComplex[ResponseHeader](ExtensionObjectDefinitionParseWithBufferProducer[ResponseHeader]((int32)(int32(394))), readBuffer))
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

		if err := WriteSimpleField[ResponseHeader](ctx, "responseHeader", m.GetResponseHeader(), WriteComplex[ResponseHeader](writeBuffer)); err != nil {
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

func (m *_ModifySubscriptionResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ModifySubscriptionResponse) deepCopy() *_ModifySubscriptionResponse {
	if m == nil {
		return nil
	}
	_ModifySubscriptionResponseCopy := &_ModifySubscriptionResponse{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.ResponseHeader.DeepCopy().(ResponseHeader),
		m.RevisedPublishingInterval,
		m.RevisedLifetimeCount,
		m.RevisedMaxKeepAliveCount,
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _ModifySubscriptionResponseCopy
}

func (m *_ModifySubscriptionResponse) String() string {
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
