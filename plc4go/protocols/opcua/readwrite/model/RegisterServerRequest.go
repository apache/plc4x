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

// RegisterServerRequest is the corresponding interface of RegisterServerRequest
type RegisterServerRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetRequestHeader returns RequestHeader (property field)
	GetRequestHeader() RequestHeader
	// GetServer returns Server (property field)
	GetServer() RegisteredServer
	// IsRegisterServerRequest is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsRegisterServerRequest()
	// CreateBuilder creates a RegisterServerRequestBuilder
	CreateRegisterServerRequestBuilder() RegisterServerRequestBuilder
}

// _RegisterServerRequest is the data-structure of this message
type _RegisterServerRequest struct {
	ExtensionObjectDefinitionContract
	RequestHeader RequestHeader
	Server        RegisteredServer
}

var _ RegisterServerRequest = (*_RegisterServerRequest)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_RegisterServerRequest)(nil)

// NewRegisterServerRequest factory function for _RegisterServerRequest
func NewRegisterServerRequest(requestHeader RequestHeader, server RegisteredServer) *_RegisterServerRequest {
	if requestHeader == nil {
		panic("requestHeader of type RequestHeader for RegisterServerRequest must not be nil")
	}
	if server == nil {
		panic("server of type RegisteredServer for RegisterServerRequest must not be nil")
	}
	_result := &_RegisterServerRequest{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		RequestHeader:                     requestHeader,
		Server:                            server,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// RegisterServerRequestBuilder is a builder for RegisterServerRequest
type RegisterServerRequestBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(requestHeader RequestHeader, server RegisteredServer) RegisterServerRequestBuilder
	// WithRequestHeader adds RequestHeader (property field)
	WithRequestHeader(RequestHeader) RegisterServerRequestBuilder
	// WithRequestHeaderBuilder adds RequestHeader (property field) which is build by the builder
	WithRequestHeaderBuilder(func(RequestHeaderBuilder) RequestHeaderBuilder) RegisterServerRequestBuilder
	// WithServer adds Server (property field)
	WithServer(RegisteredServer) RegisterServerRequestBuilder
	// WithServerBuilder adds Server (property field) which is build by the builder
	WithServerBuilder(func(RegisteredServerBuilder) RegisteredServerBuilder) RegisterServerRequestBuilder
	// Build builds the RegisterServerRequest or returns an error if something is wrong
	Build() (RegisterServerRequest, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() RegisterServerRequest
}

// NewRegisterServerRequestBuilder() creates a RegisterServerRequestBuilder
func NewRegisterServerRequestBuilder() RegisterServerRequestBuilder {
	return &_RegisterServerRequestBuilder{_RegisterServerRequest: new(_RegisterServerRequest)}
}

type _RegisterServerRequestBuilder struct {
	*_RegisterServerRequest

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (RegisterServerRequestBuilder) = (*_RegisterServerRequestBuilder)(nil)

func (b *_RegisterServerRequestBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_RegisterServerRequestBuilder) WithMandatoryFields(requestHeader RequestHeader, server RegisteredServer) RegisterServerRequestBuilder {
	return b.WithRequestHeader(requestHeader).WithServer(server)
}

func (b *_RegisterServerRequestBuilder) WithRequestHeader(requestHeader RequestHeader) RegisterServerRequestBuilder {
	b.RequestHeader = requestHeader
	return b
}

func (b *_RegisterServerRequestBuilder) WithRequestHeaderBuilder(builderSupplier func(RequestHeaderBuilder) RequestHeaderBuilder) RegisterServerRequestBuilder {
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

func (b *_RegisterServerRequestBuilder) WithServer(server RegisteredServer) RegisterServerRequestBuilder {
	b.Server = server
	return b
}

func (b *_RegisterServerRequestBuilder) WithServerBuilder(builderSupplier func(RegisteredServerBuilder) RegisteredServerBuilder) RegisterServerRequestBuilder {
	builder := builderSupplier(b.Server.CreateRegisteredServerBuilder())
	var err error
	b.Server, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "RegisteredServerBuilder failed"))
	}
	return b
}

func (b *_RegisterServerRequestBuilder) Build() (RegisterServerRequest, error) {
	if b.RequestHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'requestHeader' not set"))
	}
	if b.Server == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'server' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._RegisterServerRequest.deepCopy(), nil
}

func (b *_RegisterServerRequestBuilder) MustBuild() RegisterServerRequest {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_RegisterServerRequestBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_RegisterServerRequestBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_RegisterServerRequestBuilder) DeepCopy() any {
	_copy := b.CreateRegisterServerRequestBuilder().(*_RegisterServerRequestBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateRegisterServerRequestBuilder creates a RegisterServerRequestBuilder
func (b *_RegisterServerRequest) CreateRegisterServerRequestBuilder() RegisterServerRequestBuilder {
	if b == nil {
		return NewRegisterServerRequestBuilder()
	}
	return &_RegisterServerRequestBuilder{_RegisterServerRequest: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_RegisterServerRequest) GetExtensionId() int32 {
	return int32(437)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_RegisterServerRequest) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_RegisterServerRequest) GetRequestHeader() RequestHeader {
	return m.RequestHeader
}

func (m *_RegisterServerRequest) GetServer() RegisteredServer {
	return m.Server
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastRegisterServerRequest(structType any) RegisterServerRequest {
	if casted, ok := structType.(RegisterServerRequest); ok {
		return casted
	}
	if casted, ok := structType.(*RegisterServerRequest); ok {
		return *casted
	}
	return nil
}

func (m *_RegisterServerRequest) GetTypeName() string {
	return "RegisterServerRequest"
}

func (m *_RegisterServerRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (requestHeader)
	lengthInBits += m.RequestHeader.GetLengthInBits(ctx)

	// Simple field (server)
	lengthInBits += m.Server.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_RegisterServerRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_RegisterServerRequest) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__registerServerRequest RegisterServerRequest, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("RegisterServerRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for RegisterServerRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	requestHeader, err := ReadSimpleField[RequestHeader](ctx, "requestHeader", ReadComplex[RequestHeader](ExtensionObjectDefinitionParseWithBufferProducer[RequestHeader]((int32)(int32(391))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'requestHeader' field"))
	}
	m.RequestHeader = requestHeader

	server, err := ReadSimpleField[RegisteredServer](ctx, "server", ReadComplex[RegisteredServer](ExtensionObjectDefinitionParseWithBufferProducer[RegisteredServer]((int32)(int32(434))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'server' field"))
	}
	m.Server = server

	if closeErr := readBuffer.CloseContext("RegisterServerRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for RegisterServerRequest")
	}

	return m, nil
}

func (m *_RegisterServerRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_RegisterServerRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("RegisterServerRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for RegisterServerRequest")
		}

		if err := WriteSimpleField[RequestHeader](ctx, "requestHeader", m.GetRequestHeader(), WriteComplex[RequestHeader](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'requestHeader' field")
		}

		if err := WriteSimpleField[RegisteredServer](ctx, "server", m.GetServer(), WriteComplex[RegisteredServer](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'server' field")
		}

		if popErr := writeBuffer.PopContext("RegisterServerRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for RegisterServerRequest")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_RegisterServerRequest) IsRegisterServerRequest() {}

func (m *_RegisterServerRequest) DeepCopy() any {
	return m.deepCopy()
}

func (m *_RegisterServerRequest) deepCopy() *_RegisterServerRequest {
	if m == nil {
		return nil
	}
	_RegisterServerRequestCopy := &_RegisterServerRequest{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.RequestHeader.DeepCopy().(RequestHeader),
		m.Server.DeepCopy().(RegisteredServer),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _RegisterServerRequestCopy
}

func (m *_RegisterServerRequest) String() string {
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
