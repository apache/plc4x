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

// CloseSecureChannelResponse is the corresponding interface of CloseSecureChannelResponse
type CloseSecureChannelResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetResponseHeader returns ResponseHeader (property field)
	GetResponseHeader() ResponseHeader
	// IsCloseSecureChannelResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCloseSecureChannelResponse()
	// CreateBuilder creates a CloseSecureChannelResponseBuilder
	CreateCloseSecureChannelResponseBuilder() CloseSecureChannelResponseBuilder
}

// _CloseSecureChannelResponse is the data-structure of this message
type _CloseSecureChannelResponse struct {
	ExtensionObjectDefinitionContract
	ResponseHeader ResponseHeader
}

var _ CloseSecureChannelResponse = (*_CloseSecureChannelResponse)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_CloseSecureChannelResponse)(nil)

// NewCloseSecureChannelResponse factory function for _CloseSecureChannelResponse
func NewCloseSecureChannelResponse(responseHeader ResponseHeader) *_CloseSecureChannelResponse {
	if responseHeader == nil {
		panic("responseHeader of type ResponseHeader for CloseSecureChannelResponse must not be nil")
	}
	_result := &_CloseSecureChannelResponse{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		ResponseHeader:                    responseHeader,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// CloseSecureChannelResponseBuilder is a builder for CloseSecureChannelResponse
type CloseSecureChannelResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(responseHeader ResponseHeader) CloseSecureChannelResponseBuilder
	// WithResponseHeader adds ResponseHeader (property field)
	WithResponseHeader(ResponseHeader) CloseSecureChannelResponseBuilder
	// WithResponseHeaderBuilder adds ResponseHeader (property field) which is build by the builder
	WithResponseHeaderBuilder(func(ResponseHeaderBuilder) ResponseHeaderBuilder) CloseSecureChannelResponseBuilder
	// Build builds the CloseSecureChannelResponse or returns an error if something is wrong
	Build() (CloseSecureChannelResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() CloseSecureChannelResponse
}

// NewCloseSecureChannelResponseBuilder() creates a CloseSecureChannelResponseBuilder
func NewCloseSecureChannelResponseBuilder() CloseSecureChannelResponseBuilder {
	return &_CloseSecureChannelResponseBuilder{_CloseSecureChannelResponse: new(_CloseSecureChannelResponse)}
}

type _CloseSecureChannelResponseBuilder struct {
	*_CloseSecureChannelResponse

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (CloseSecureChannelResponseBuilder) = (*_CloseSecureChannelResponseBuilder)(nil)

func (b *_CloseSecureChannelResponseBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_CloseSecureChannelResponseBuilder) WithMandatoryFields(responseHeader ResponseHeader) CloseSecureChannelResponseBuilder {
	return b.WithResponseHeader(responseHeader)
}

func (b *_CloseSecureChannelResponseBuilder) WithResponseHeader(responseHeader ResponseHeader) CloseSecureChannelResponseBuilder {
	b.ResponseHeader = responseHeader
	return b
}

func (b *_CloseSecureChannelResponseBuilder) WithResponseHeaderBuilder(builderSupplier func(ResponseHeaderBuilder) ResponseHeaderBuilder) CloseSecureChannelResponseBuilder {
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

func (b *_CloseSecureChannelResponseBuilder) Build() (CloseSecureChannelResponse, error) {
	if b.ResponseHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'responseHeader' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._CloseSecureChannelResponse.deepCopy(), nil
}

func (b *_CloseSecureChannelResponseBuilder) MustBuild() CloseSecureChannelResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_CloseSecureChannelResponseBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_CloseSecureChannelResponseBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_CloseSecureChannelResponseBuilder) DeepCopy() any {
	_copy := b.CreateCloseSecureChannelResponseBuilder().(*_CloseSecureChannelResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateCloseSecureChannelResponseBuilder creates a CloseSecureChannelResponseBuilder
func (b *_CloseSecureChannelResponse) CreateCloseSecureChannelResponseBuilder() CloseSecureChannelResponseBuilder {
	if b == nil {
		return NewCloseSecureChannelResponseBuilder()
	}
	return &_CloseSecureChannelResponseBuilder{_CloseSecureChannelResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_CloseSecureChannelResponse) GetExtensionId() int32 {
	return int32(455)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_CloseSecureChannelResponse) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CloseSecureChannelResponse) GetResponseHeader() ResponseHeader {
	return m.ResponseHeader
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastCloseSecureChannelResponse(structType any) CloseSecureChannelResponse {
	if casted, ok := structType.(CloseSecureChannelResponse); ok {
		return casted
	}
	if casted, ok := structType.(*CloseSecureChannelResponse); ok {
		return *casted
	}
	return nil
}

func (m *_CloseSecureChannelResponse) GetTypeName() string {
	return "CloseSecureChannelResponse"
}

func (m *_CloseSecureChannelResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (responseHeader)
	lengthInBits += m.ResponseHeader.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_CloseSecureChannelResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_CloseSecureChannelResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__closeSecureChannelResponse CloseSecureChannelResponse, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CloseSecureChannelResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CloseSecureChannelResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	responseHeader, err := ReadSimpleField[ResponseHeader](ctx, "responseHeader", ReadComplex[ResponseHeader](ExtensionObjectDefinitionParseWithBufferProducer[ResponseHeader]((int32)(int32(394))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'responseHeader' field"))
	}
	m.ResponseHeader = responseHeader

	if closeErr := readBuffer.CloseContext("CloseSecureChannelResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CloseSecureChannelResponse")
	}

	return m, nil
}

func (m *_CloseSecureChannelResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CloseSecureChannelResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CloseSecureChannelResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CloseSecureChannelResponse")
		}

		if err := WriteSimpleField[ResponseHeader](ctx, "responseHeader", m.GetResponseHeader(), WriteComplex[ResponseHeader](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'responseHeader' field")
		}

		if popErr := writeBuffer.PopContext("CloseSecureChannelResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CloseSecureChannelResponse")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_CloseSecureChannelResponse) IsCloseSecureChannelResponse() {}

func (m *_CloseSecureChannelResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_CloseSecureChannelResponse) deepCopy() *_CloseSecureChannelResponse {
	if m == nil {
		return nil
	}
	_CloseSecureChannelResponseCopy := &_CloseSecureChannelResponse{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.ResponseHeader.DeepCopy().(ResponseHeader),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _CloseSecureChannelResponseCopy
}

func (m *_CloseSecureChannelResponse) String() string {
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
