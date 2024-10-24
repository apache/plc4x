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

// UnregisterNodesRequest is the corresponding interface of UnregisterNodesRequest
type UnregisterNodesRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetRequestHeader returns RequestHeader (property field)
	GetRequestHeader() RequestHeader
	// GetNodesToUnregister returns NodesToUnregister (property field)
	GetNodesToUnregister() []NodeId
	// IsUnregisterNodesRequest is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsUnregisterNodesRequest()
	// CreateBuilder creates a UnregisterNodesRequestBuilder
	CreateUnregisterNodesRequestBuilder() UnregisterNodesRequestBuilder
}

// _UnregisterNodesRequest is the data-structure of this message
type _UnregisterNodesRequest struct {
	ExtensionObjectDefinitionContract
	RequestHeader     RequestHeader
	NodesToUnregister []NodeId
}

var _ UnregisterNodesRequest = (*_UnregisterNodesRequest)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_UnregisterNodesRequest)(nil)

// NewUnregisterNodesRequest factory function for _UnregisterNodesRequest
func NewUnregisterNodesRequest(requestHeader RequestHeader, nodesToUnregister []NodeId) *_UnregisterNodesRequest {
	if requestHeader == nil {
		panic("requestHeader of type RequestHeader for UnregisterNodesRequest must not be nil")
	}
	_result := &_UnregisterNodesRequest{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		RequestHeader:                     requestHeader,
		NodesToUnregister:                 nodesToUnregister,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// UnregisterNodesRequestBuilder is a builder for UnregisterNodesRequest
type UnregisterNodesRequestBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(requestHeader RequestHeader, nodesToUnregister []NodeId) UnregisterNodesRequestBuilder
	// WithRequestHeader adds RequestHeader (property field)
	WithRequestHeader(RequestHeader) UnregisterNodesRequestBuilder
	// WithRequestHeaderBuilder adds RequestHeader (property field) which is build by the builder
	WithRequestHeaderBuilder(func(RequestHeaderBuilder) RequestHeaderBuilder) UnregisterNodesRequestBuilder
	// WithNodesToUnregister adds NodesToUnregister (property field)
	WithNodesToUnregister(...NodeId) UnregisterNodesRequestBuilder
	// Build builds the UnregisterNodesRequest or returns an error if something is wrong
	Build() (UnregisterNodesRequest, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() UnregisterNodesRequest
}

// NewUnregisterNodesRequestBuilder() creates a UnregisterNodesRequestBuilder
func NewUnregisterNodesRequestBuilder() UnregisterNodesRequestBuilder {
	return &_UnregisterNodesRequestBuilder{_UnregisterNodesRequest: new(_UnregisterNodesRequest)}
}

type _UnregisterNodesRequestBuilder struct {
	*_UnregisterNodesRequest

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (UnregisterNodesRequestBuilder) = (*_UnregisterNodesRequestBuilder)(nil)

func (b *_UnregisterNodesRequestBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_UnregisterNodesRequestBuilder) WithMandatoryFields(requestHeader RequestHeader, nodesToUnregister []NodeId) UnregisterNodesRequestBuilder {
	return b.WithRequestHeader(requestHeader).WithNodesToUnregister(nodesToUnregister...)
}

func (b *_UnregisterNodesRequestBuilder) WithRequestHeader(requestHeader RequestHeader) UnregisterNodesRequestBuilder {
	b.RequestHeader = requestHeader
	return b
}

func (b *_UnregisterNodesRequestBuilder) WithRequestHeaderBuilder(builderSupplier func(RequestHeaderBuilder) RequestHeaderBuilder) UnregisterNodesRequestBuilder {
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

func (b *_UnregisterNodesRequestBuilder) WithNodesToUnregister(nodesToUnregister ...NodeId) UnregisterNodesRequestBuilder {
	b.NodesToUnregister = nodesToUnregister
	return b
}

func (b *_UnregisterNodesRequestBuilder) Build() (UnregisterNodesRequest, error) {
	if b.RequestHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'requestHeader' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._UnregisterNodesRequest.deepCopy(), nil
}

func (b *_UnregisterNodesRequestBuilder) MustBuild() UnregisterNodesRequest {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_UnregisterNodesRequestBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_UnregisterNodesRequestBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_UnregisterNodesRequestBuilder) DeepCopy() any {
	_copy := b.CreateUnregisterNodesRequestBuilder().(*_UnregisterNodesRequestBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateUnregisterNodesRequestBuilder creates a UnregisterNodesRequestBuilder
func (b *_UnregisterNodesRequest) CreateUnregisterNodesRequestBuilder() UnregisterNodesRequestBuilder {
	if b == nil {
		return NewUnregisterNodesRequestBuilder()
	}
	return &_UnregisterNodesRequestBuilder{_UnregisterNodesRequest: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_UnregisterNodesRequest) GetExtensionId() int32 {
	return int32(566)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_UnregisterNodesRequest) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_UnregisterNodesRequest) GetRequestHeader() RequestHeader {
	return m.RequestHeader
}

func (m *_UnregisterNodesRequest) GetNodesToUnregister() []NodeId {
	return m.NodesToUnregister
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastUnregisterNodesRequest(structType any) UnregisterNodesRequest {
	if casted, ok := structType.(UnregisterNodesRequest); ok {
		return casted
	}
	if casted, ok := structType.(*UnregisterNodesRequest); ok {
		return *casted
	}
	return nil
}

func (m *_UnregisterNodesRequest) GetTypeName() string {
	return "UnregisterNodesRequest"
}

func (m *_UnregisterNodesRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (requestHeader)
	lengthInBits += m.RequestHeader.GetLengthInBits(ctx)

	// Implicit Field (noOfNodesToUnregister)
	lengthInBits += 32

	// Array field
	if len(m.NodesToUnregister) > 0 {
		for _curItem, element := range m.NodesToUnregister {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.NodesToUnregister), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_UnregisterNodesRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_UnregisterNodesRequest) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__unregisterNodesRequest UnregisterNodesRequest, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("UnregisterNodesRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for UnregisterNodesRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	requestHeader, err := ReadSimpleField[RequestHeader](ctx, "requestHeader", ReadComplex[RequestHeader](ExtensionObjectDefinitionParseWithBufferProducer[RequestHeader]((int32)(int32(391))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'requestHeader' field"))
	}
	m.RequestHeader = requestHeader

	noOfNodesToUnregister, err := ReadImplicitField[int32](ctx, "noOfNodesToUnregister", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfNodesToUnregister' field"))
	}
	_ = noOfNodesToUnregister

	nodesToUnregister, err := ReadCountArrayField[NodeId](ctx, "nodesToUnregister", ReadComplex[NodeId](NodeIdParseWithBuffer, readBuffer), uint64(noOfNodesToUnregister))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'nodesToUnregister' field"))
	}
	m.NodesToUnregister = nodesToUnregister

	if closeErr := readBuffer.CloseContext("UnregisterNodesRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for UnregisterNodesRequest")
	}

	return m, nil
}

func (m *_UnregisterNodesRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_UnregisterNodesRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("UnregisterNodesRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for UnregisterNodesRequest")
		}

		if err := WriteSimpleField[RequestHeader](ctx, "requestHeader", m.GetRequestHeader(), WriteComplex[RequestHeader](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'requestHeader' field")
		}
		noOfNodesToUnregister := int32(utils.InlineIf(bool((m.GetNodesToUnregister()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetNodesToUnregister()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfNodesToUnregister", noOfNodesToUnregister, WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfNodesToUnregister' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "nodesToUnregister", m.GetNodesToUnregister(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'nodesToUnregister' field")
		}

		if popErr := writeBuffer.PopContext("UnregisterNodesRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for UnregisterNodesRequest")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_UnregisterNodesRequest) IsUnregisterNodesRequest() {}

func (m *_UnregisterNodesRequest) DeepCopy() any {
	return m.deepCopy()
}

func (m *_UnregisterNodesRequest) deepCopy() *_UnregisterNodesRequest {
	if m == nil {
		return nil
	}
	_UnregisterNodesRequestCopy := &_UnregisterNodesRequest{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.RequestHeader.DeepCopy().(RequestHeader),
		utils.DeepCopySlice[NodeId, NodeId](m.NodesToUnregister),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _UnregisterNodesRequestCopy
}

func (m *_UnregisterNodesRequest) String() string {
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
