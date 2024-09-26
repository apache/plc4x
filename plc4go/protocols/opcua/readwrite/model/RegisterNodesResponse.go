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

// RegisterNodesResponse is the corresponding interface of RegisterNodesResponse
type RegisterNodesResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetResponseHeader returns ResponseHeader (property field)
	GetResponseHeader() ResponseHeader
	// GetRegisteredNodeIds returns RegisteredNodeIds (property field)
	GetRegisteredNodeIds() []NodeId
	// IsRegisterNodesResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsRegisterNodesResponse()
	// CreateBuilder creates a RegisterNodesResponseBuilder
	CreateRegisterNodesResponseBuilder() RegisterNodesResponseBuilder
}

// _RegisterNodesResponse is the data-structure of this message
type _RegisterNodesResponse struct {
	ExtensionObjectDefinitionContract
	ResponseHeader    ResponseHeader
	RegisteredNodeIds []NodeId
}

var _ RegisterNodesResponse = (*_RegisterNodesResponse)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_RegisterNodesResponse)(nil)

// NewRegisterNodesResponse factory function for _RegisterNodesResponse
func NewRegisterNodesResponse(responseHeader ResponseHeader, registeredNodeIds []NodeId) *_RegisterNodesResponse {
	if responseHeader == nil {
		panic("responseHeader of type ResponseHeader for RegisterNodesResponse must not be nil")
	}
	_result := &_RegisterNodesResponse{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		ResponseHeader:                    responseHeader,
		RegisteredNodeIds:                 registeredNodeIds,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// RegisterNodesResponseBuilder is a builder for RegisterNodesResponse
type RegisterNodesResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(responseHeader ResponseHeader, registeredNodeIds []NodeId) RegisterNodesResponseBuilder
	// WithResponseHeader adds ResponseHeader (property field)
	WithResponseHeader(ResponseHeader) RegisterNodesResponseBuilder
	// WithResponseHeaderBuilder adds ResponseHeader (property field) which is build by the builder
	WithResponseHeaderBuilder(func(ResponseHeaderBuilder) ResponseHeaderBuilder) RegisterNodesResponseBuilder
	// WithRegisteredNodeIds adds RegisteredNodeIds (property field)
	WithRegisteredNodeIds(...NodeId) RegisterNodesResponseBuilder
	// Build builds the RegisterNodesResponse or returns an error if something is wrong
	Build() (RegisterNodesResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() RegisterNodesResponse
}

// NewRegisterNodesResponseBuilder() creates a RegisterNodesResponseBuilder
func NewRegisterNodesResponseBuilder() RegisterNodesResponseBuilder {
	return &_RegisterNodesResponseBuilder{_RegisterNodesResponse: new(_RegisterNodesResponse)}
}

type _RegisterNodesResponseBuilder struct {
	*_RegisterNodesResponse

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (RegisterNodesResponseBuilder) = (*_RegisterNodesResponseBuilder)(nil)

func (b *_RegisterNodesResponseBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_RegisterNodesResponseBuilder) WithMandatoryFields(responseHeader ResponseHeader, registeredNodeIds []NodeId) RegisterNodesResponseBuilder {
	return b.WithResponseHeader(responseHeader).WithRegisteredNodeIds(registeredNodeIds...)
}

func (b *_RegisterNodesResponseBuilder) WithResponseHeader(responseHeader ResponseHeader) RegisterNodesResponseBuilder {
	b.ResponseHeader = responseHeader
	return b
}

func (b *_RegisterNodesResponseBuilder) WithResponseHeaderBuilder(builderSupplier func(ResponseHeaderBuilder) ResponseHeaderBuilder) RegisterNodesResponseBuilder {
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

func (b *_RegisterNodesResponseBuilder) WithRegisteredNodeIds(registeredNodeIds ...NodeId) RegisterNodesResponseBuilder {
	b.RegisteredNodeIds = registeredNodeIds
	return b
}

func (b *_RegisterNodesResponseBuilder) Build() (RegisterNodesResponse, error) {
	if b.ResponseHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'responseHeader' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._RegisterNodesResponse.deepCopy(), nil
}

func (b *_RegisterNodesResponseBuilder) MustBuild() RegisterNodesResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_RegisterNodesResponseBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_RegisterNodesResponseBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_RegisterNodesResponseBuilder) DeepCopy() any {
	_copy := b.CreateRegisterNodesResponseBuilder().(*_RegisterNodesResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateRegisterNodesResponseBuilder creates a RegisterNodesResponseBuilder
func (b *_RegisterNodesResponse) CreateRegisterNodesResponseBuilder() RegisterNodesResponseBuilder {
	if b == nil {
		return NewRegisterNodesResponseBuilder()
	}
	return &_RegisterNodesResponseBuilder{_RegisterNodesResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_RegisterNodesResponse) GetExtensionId() int32 {
	return int32(563)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_RegisterNodesResponse) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_RegisterNodesResponse) GetResponseHeader() ResponseHeader {
	return m.ResponseHeader
}

func (m *_RegisterNodesResponse) GetRegisteredNodeIds() []NodeId {
	return m.RegisteredNodeIds
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastRegisterNodesResponse(structType any) RegisterNodesResponse {
	if casted, ok := structType.(RegisterNodesResponse); ok {
		return casted
	}
	if casted, ok := structType.(*RegisterNodesResponse); ok {
		return *casted
	}
	return nil
}

func (m *_RegisterNodesResponse) GetTypeName() string {
	return "RegisterNodesResponse"
}

func (m *_RegisterNodesResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (responseHeader)
	lengthInBits += m.ResponseHeader.GetLengthInBits(ctx)

	// Implicit Field (noOfRegisteredNodeIds)
	lengthInBits += 32

	// Array field
	if len(m.RegisteredNodeIds) > 0 {
		for _curItem, element := range m.RegisteredNodeIds {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.RegisteredNodeIds), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_RegisterNodesResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_RegisterNodesResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__registerNodesResponse RegisterNodesResponse, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("RegisterNodesResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for RegisterNodesResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	responseHeader, err := ReadSimpleField[ResponseHeader](ctx, "responseHeader", ReadComplex[ResponseHeader](ExtensionObjectDefinitionParseWithBufferProducer[ResponseHeader]((int32)(int32(394))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'responseHeader' field"))
	}
	m.ResponseHeader = responseHeader

	noOfRegisteredNodeIds, err := ReadImplicitField[int32](ctx, "noOfRegisteredNodeIds", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfRegisteredNodeIds' field"))
	}
	_ = noOfRegisteredNodeIds

	registeredNodeIds, err := ReadCountArrayField[NodeId](ctx, "registeredNodeIds", ReadComplex[NodeId](NodeIdParseWithBuffer, readBuffer), uint64(noOfRegisteredNodeIds))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'registeredNodeIds' field"))
	}
	m.RegisteredNodeIds = registeredNodeIds

	if closeErr := readBuffer.CloseContext("RegisterNodesResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for RegisterNodesResponse")
	}

	return m, nil
}

func (m *_RegisterNodesResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_RegisterNodesResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("RegisterNodesResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for RegisterNodesResponse")
		}

		if err := WriteSimpleField[ResponseHeader](ctx, "responseHeader", m.GetResponseHeader(), WriteComplex[ResponseHeader](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'responseHeader' field")
		}
		noOfRegisteredNodeIds := int32(utils.InlineIf(bool((m.GetRegisteredNodeIds()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetRegisteredNodeIds()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfRegisteredNodeIds", noOfRegisteredNodeIds, WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfRegisteredNodeIds' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "registeredNodeIds", m.GetRegisteredNodeIds(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'registeredNodeIds' field")
		}

		if popErr := writeBuffer.PopContext("RegisterNodesResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for RegisterNodesResponse")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_RegisterNodesResponse) IsRegisterNodesResponse() {}

func (m *_RegisterNodesResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_RegisterNodesResponse) deepCopy() *_RegisterNodesResponse {
	if m == nil {
		return nil
	}
	_RegisterNodesResponseCopy := &_RegisterNodesResponse{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.ResponseHeader.DeepCopy().(ResponseHeader),
		utils.DeepCopySlice[NodeId, NodeId](m.RegisteredNodeIds),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _RegisterNodesResponseCopy
}

func (m *_RegisterNodesResponse) String() string {
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
