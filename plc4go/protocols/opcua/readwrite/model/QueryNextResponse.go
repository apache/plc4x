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

// QueryNextResponse is the corresponding interface of QueryNextResponse
type QueryNextResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetResponseHeader returns ResponseHeader (property field)
	GetResponseHeader() ResponseHeader
	// GetQueryDataSets returns QueryDataSets (property field)
	GetQueryDataSets() []QueryDataSet
	// GetRevisedContinuationPoint returns RevisedContinuationPoint (property field)
	GetRevisedContinuationPoint() PascalByteString
	// IsQueryNextResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsQueryNextResponse()
	// CreateBuilder creates a QueryNextResponseBuilder
	CreateQueryNextResponseBuilder() QueryNextResponseBuilder
}

// _QueryNextResponse is the data-structure of this message
type _QueryNextResponse struct {
	ExtensionObjectDefinitionContract
	ResponseHeader           ResponseHeader
	QueryDataSets            []QueryDataSet
	RevisedContinuationPoint PascalByteString
}

var _ QueryNextResponse = (*_QueryNextResponse)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_QueryNextResponse)(nil)

// NewQueryNextResponse factory function for _QueryNextResponse
func NewQueryNextResponse(responseHeader ResponseHeader, queryDataSets []QueryDataSet, revisedContinuationPoint PascalByteString) *_QueryNextResponse {
	if responseHeader == nil {
		panic("responseHeader of type ResponseHeader for QueryNextResponse must not be nil")
	}
	if revisedContinuationPoint == nil {
		panic("revisedContinuationPoint of type PascalByteString for QueryNextResponse must not be nil")
	}
	_result := &_QueryNextResponse{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		ResponseHeader:                    responseHeader,
		QueryDataSets:                     queryDataSets,
		RevisedContinuationPoint:          revisedContinuationPoint,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// QueryNextResponseBuilder is a builder for QueryNextResponse
type QueryNextResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(responseHeader ResponseHeader, queryDataSets []QueryDataSet, revisedContinuationPoint PascalByteString) QueryNextResponseBuilder
	// WithResponseHeader adds ResponseHeader (property field)
	WithResponseHeader(ResponseHeader) QueryNextResponseBuilder
	// WithResponseHeaderBuilder adds ResponseHeader (property field) which is build by the builder
	WithResponseHeaderBuilder(func(ResponseHeaderBuilder) ResponseHeaderBuilder) QueryNextResponseBuilder
	// WithQueryDataSets adds QueryDataSets (property field)
	WithQueryDataSets(...QueryDataSet) QueryNextResponseBuilder
	// WithRevisedContinuationPoint adds RevisedContinuationPoint (property field)
	WithRevisedContinuationPoint(PascalByteString) QueryNextResponseBuilder
	// WithRevisedContinuationPointBuilder adds RevisedContinuationPoint (property field) which is build by the builder
	WithRevisedContinuationPointBuilder(func(PascalByteStringBuilder) PascalByteStringBuilder) QueryNextResponseBuilder
	// Build builds the QueryNextResponse or returns an error if something is wrong
	Build() (QueryNextResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() QueryNextResponse
}

// NewQueryNextResponseBuilder() creates a QueryNextResponseBuilder
func NewQueryNextResponseBuilder() QueryNextResponseBuilder {
	return &_QueryNextResponseBuilder{_QueryNextResponse: new(_QueryNextResponse)}
}

type _QueryNextResponseBuilder struct {
	*_QueryNextResponse

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (QueryNextResponseBuilder) = (*_QueryNextResponseBuilder)(nil)

func (b *_QueryNextResponseBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_QueryNextResponseBuilder) WithMandatoryFields(responseHeader ResponseHeader, queryDataSets []QueryDataSet, revisedContinuationPoint PascalByteString) QueryNextResponseBuilder {
	return b.WithResponseHeader(responseHeader).WithQueryDataSets(queryDataSets...).WithRevisedContinuationPoint(revisedContinuationPoint)
}

func (b *_QueryNextResponseBuilder) WithResponseHeader(responseHeader ResponseHeader) QueryNextResponseBuilder {
	b.ResponseHeader = responseHeader
	return b
}

func (b *_QueryNextResponseBuilder) WithResponseHeaderBuilder(builderSupplier func(ResponseHeaderBuilder) ResponseHeaderBuilder) QueryNextResponseBuilder {
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

func (b *_QueryNextResponseBuilder) WithQueryDataSets(queryDataSets ...QueryDataSet) QueryNextResponseBuilder {
	b.QueryDataSets = queryDataSets
	return b
}

func (b *_QueryNextResponseBuilder) WithRevisedContinuationPoint(revisedContinuationPoint PascalByteString) QueryNextResponseBuilder {
	b.RevisedContinuationPoint = revisedContinuationPoint
	return b
}

func (b *_QueryNextResponseBuilder) WithRevisedContinuationPointBuilder(builderSupplier func(PascalByteStringBuilder) PascalByteStringBuilder) QueryNextResponseBuilder {
	builder := builderSupplier(b.RevisedContinuationPoint.CreatePascalByteStringBuilder())
	var err error
	b.RevisedContinuationPoint, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalByteStringBuilder failed"))
	}
	return b
}

func (b *_QueryNextResponseBuilder) Build() (QueryNextResponse, error) {
	if b.ResponseHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'responseHeader' not set"))
	}
	if b.RevisedContinuationPoint == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'revisedContinuationPoint' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._QueryNextResponse.deepCopy(), nil
}

func (b *_QueryNextResponseBuilder) MustBuild() QueryNextResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_QueryNextResponseBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_QueryNextResponseBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_QueryNextResponseBuilder) DeepCopy() any {
	_copy := b.CreateQueryNextResponseBuilder().(*_QueryNextResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateQueryNextResponseBuilder creates a QueryNextResponseBuilder
func (b *_QueryNextResponse) CreateQueryNextResponseBuilder() QueryNextResponseBuilder {
	if b == nil {
		return NewQueryNextResponseBuilder()
	}
	return &_QueryNextResponseBuilder{_QueryNextResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_QueryNextResponse) GetExtensionId() int32 {
	return int32(624)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_QueryNextResponse) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_QueryNextResponse) GetResponseHeader() ResponseHeader {
	return m.ResponseHeader
}

func (m *_QueryNextResponse) GetQueryDataSets() []QueryDataSet {
	return m.QueryDataSets
}

func (m *_QueryNextResponse) GetRevisedContinuationPoint() PascalByteString {
	return m.RevisedContinuationPoint
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastQueryNextResponse(structType any) QueryNextResponse {
	if casted, ok := structType.(QueryNextResponse); ok {
		return casted
	}
	if casted, ok := structType.(*QueryNextResponse); ok {
		return *casted
	}
	return nil
}

func (m *_QueryNextResponse) GetTypeName() string {
	return "QueryNextResponse"
}

func (m *_QueryNextResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (responseHeader)
	lengthInBits += m.ResponseHeader.GetLengthInBits(ctx)

	// Implicit Field (noOfQueryDataSets)
	lengthInBits += 32

	// Array field
	if len(m.QueryDataSets) > 0 {
		for _curItem, element := range m.QueryDataSets {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.QueryDataSets), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Simple field (revisedContinuationPoint)
	lengthInBits += m.RevisedContinuationPoint.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_QueryNextResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_QueryNextResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__queryNextResponse QueryNextResponse, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("QueryNextResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for QueryNextResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	responseHeader, err := ReadSimpleField[ResponseHeader](ctx, "responseHeader", ReadComplex[ResponseHeader](ExtensionObjectDefinitionParseWithBufferProducer[ResponseHeader]((int32)(int32(394))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'responseHeader' field"))
	}
	m.ResponseHeader = responseHeader

	noOfQueryDataSets, err := ReadImplicitField[int32](ctx, "noOfQueryDataSets", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfQueryDataSets' field"))
	}
	_ = noOfQueryDataSets

	queryDataSets, err := ReadCountArrayField[QueryDataSet](ctx, "queryDataSets", ReadComplex[QueryDataSet](ExtensionObjectDefinitionParseWithBufferProducer[QueryDataSet]((int32)(int32(579))), readBuffer), uint64(noOfQueryDataSets))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'queryDataSets' field"))
	}
	m.QueryDataSets = queryDataSets

	revisedContinuationPoint, err := ReadSimpleField[PascalByteString](ctx, "revisedContinuationPoint", ReadComplex[PascalByteString](PascalByteStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'revisedContinuationPoint' field"))
	}
	m.RevisedContinuationPoint = revisedContinuationPoint

	if closeErr := readBuffer.CloseContext("QueryNextResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for QueryNextResponse")
	}

	return m, nil
}

func (m *_QueryNextResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_QueryNextResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("QueryNextResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for QueryNextResponse")
		}

		if err := WriteSimpleField[ResponseHeader](ctx, "responseHeader", m.GetResponseHeader(), WriteComplex[ResponseHeader](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'responseHeader' field")
		}
		noOfQueryDataSets := int32(utils.InlineIf(bool((m.GetQueryDataSets()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetQueryDataSets()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfQueryDataSets", noOfQueryDataSets, WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfQueryDataSets' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "queryDataSets", m.GetQueryDataSets(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'queryDataSets' field")
		}

		if err := WriteSimpleField[PascalByteString](ctx, "revisedContinuationPoint", m.GetRevisedContinuationPoint(), WriteComplex[PascalByteString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'revisedContinuationPoint' field")
		}

		if popErr := writeBuffer.PopContext("QueryNextResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for QueryNextResponse")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_QueryNextResponse) IsQueryNextResponse() {}

func (m *_QueryNextResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_QueryNextResponse) deepCopy() *_QueryNextResponse {
	if m == nil {
		return nil
	}
	_QueryNextResponseCopy := &_QueryNextResponse{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.ResponseHeader.DeepCopy().(ResponseHeader),
		utils.DeepCopySlice[QueryDataSet, QueryDataSet](m.QueryDataSets),
		m.RevisedContinuationPoint.DeepCopy().(PascalByteString),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _QueryNextResponseCopy
}

func (m *_QueryNextResponse) String() string {
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
