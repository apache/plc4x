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

// CreateSessionResponse is the corresponding interface of CreateSessionResponse
type CreateSessionResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetResponseHeader returns ResponseHeader (property field)
	GetResponseHeader() ResponseHeader
	// GetSessionId returns SessionId (property field)
	GetSessionId() NodeId
	// GetAuthenticationToken returns AuthenticationToken (property field)
	GetAuthenticationToken() NodeId
	// GetRevisedSessionTimeout returns RevisedSessionTimeout (property field)
	GetRevisedSessionTimeout() float64
	// GetServerNonce returns ServerNonce (property field)
	GetServerNonce() PascalByteString
	// GetServerCertificate returns ServerCertificate (property field)
	GetServerCertificate() PascalByteString
	// GetServerEndpoints returns ServerEndpoints (property field)
	GetServerEndpoints() []EndpointDescription
	// GetServerSoftwareCertificates returns ServerSoftwareCertificates (property field)
	GetServerSoftwareCertificates() []SignedSoftwareCertificate
	// GetServerSignature returns ServerSignature (property field)
	GetServerSignature() SignatureData
	// GetMaxRequestMessageSize returns MaxRequestMessageSize (property field)
	GetMaxRequestMessageSize() uint32
	// IsCreateSessionResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCreateSessionResponse()
	// CreateBuilder creates a CreateSessionResponseBuilder
	CreateCreateSessionResponseBuilder() CreateSessionResponseBuilder
}

// _CreateSessionResponse is the data-structure of this message
type _CreateSessionResponse struct {
	ExtensionObjectDefinitionContract
	ResponseHeader             ResponseHeader
	SessionId                  NodeId
	AuthenticationToken        NodeId
	RevisedSessionTimeout      float64
	ServerNonce                PascalByteString
	ServerCertificate          PascalByteString
	ServerEndpoints            []EndpointDescription
	ServerSoftwareCertificates []SignedSoftwareCertificate
	ServerSignature            SignatureData
	MaxRequestMessageSize      uint32
}

var _ CreateSessionResponse = (*_CreateSessionResponse)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_CreateSessionResponse)(nil)

// NewCreateSessionResponse factory function for _CreateSessionResponse
func NewCreateSessionResponse(responseHeader ResponseHeader, sessionId NodeId, authenticationToken NodeId, revisedSessionTimeout float64, serverNonce PascalByteString, serverCertificate PascalByteString, serverEndpoints []EndpointDescription, serverSoftwareCertificates []SignedSoftwareCertificate, serverSignature SignatureData, maxRequestMessageSize uint32) *_CreateSessionResponse {
	if responseHeader == nil {
		panic("responseHeader of type ResponseHeader for CreateSessionResponse must not be nil")
	}
	if sessionId == nil {
		panic("sessionId of type NodeId for CreateSessionResponse must not be nil")
	}
	if authenticationToken == nil {
		panic("authenticationToken of type NodeId for CreateSessionResponse must not be nil")
	}
	if serverNonce == nil {
		panic("serverNonce of type PascalByteString for CreateSessionResponse must not be nil")
	}
	if serverCertificate == nil {
		panic("serverCertificate of type PascalByteString for CreateSessionResponse must not be nil")
	}
	if serverSignature == nil {
		panic("serverSignature of type SignatureData for CreateSessionResponse must not be nil")
	}
	_result := &_CreateSessionResponse{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		ResponseHeader:                    responseHeader,
		SessionId:                         sessionId,
		AuthenticationToken:               authenticationToken,
		RevisedSessionTimeout:             revisedSessionTimeout,
		ServerNonce:                       serverNonce,
		ServerCertificate:                 serverCertificate,
		ServerEndpoints:                   serverEndpoints,
		ServerSoftwareCertificates:        serverSoftwareCertificates,
		ServerSignature:                   serverSignature,
		MaxRequestMessageSize:             maxRequestMessageSize,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// CreateSessionResponseBuilder is a builder for CreateSessionResponse
type CreateSessionResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(responseHeader ResponseHeader, sessionId NodeId, authenticationToken NodeId, revisedSessionTimeout float64, serverNonce PascalByteString, serverCertificate PascalByteString, serverEndpoints []EndpointDescription, serverSoftwareCertificates []SignedSoftwareCertificate, serverSignature SignatureData, maxRequestMessageSize uint32) CreateSessionResponseBuilder
	// WithResponseHeader adds ResponseHeader (property field)
	WithResponseHeader(ResponseHeader) CreateSessionResponseBuilder
	// WithResponseHeaderBuilder adds ResponseHeader (property field) which is build by the builder
	WithResponseHeaderBuilder(func(ResponseHeaderBuilder) ResponseHeaderBuilder) CreateSessionResponseBuilder
	// WithSessionId adds SessionId (property field)
	WithSessionId(NodeId) CreateSessionResponseBuilder
	// WithSessionIdBuilder adds SessionId (property field) which is build by the builder
	WithSessionIdBuilder(func(NodeIdBuilder) NodeIdBuilder) CreateSessionResponseBuilder
	// WithAuthenticationToken adds AuthenticationToken (property field)
	WithAuthenticationToken(NodeId) CreateSessionResponseBuilder
	// WithAuthenticationTokenBuilder adds AuthenticationToken (property field) which is build by the builder
	WithAuthenticationTokenBuilder(func(NodeIdBuilder) NodeIdBuilder) CreateSessionResponseBuilder
	// WithRevisedSessionTimeout adds RevisedSessionTimeout (property field)
	WithRevisedSessionTimeout(float64) CreateSessionResponseBuilder
	// WithServerNonce adds ServerNonce (property field)
	WithServerNonce(PascalByteString) CreateSessionResponseBuilder
	// WithServerNonceBuilder adds ServerNonce (property field) which is build by the builder
	WithServerNonceBuilder(func(PascalByteStringBuilder) PascalByteStringBuilder) CreateSessionResponseBuilder
	// WithServerCertificate adds ServerCertificate (property field)
	WithServerCertificate(PascalByteString) CreateSessionResponseBuilder
	// WithServerCertificateBuilder adds ServerCertificate (property field) which is build by the builder
	WithServerCertificateBuilder(func(PascalByteStringBuilder) PascalByteStringBuilder) CreateSessionResponseBuilder
	// WithServerEndpoints adds ServerEndpoints (property field)
	WithServerEndpoints(...EndpointDescription) CreateSessionResponseBuilder
	// WithServerSoftwareCertificates adds ServerSoftwareCertificates (property field)
	WithServerSoftwareCertificates(...SignedSoftwareCertificate) CreateSessionResponseBuilder
	// WithServerSignature adds ServerSignature (property field)
	WithServerSignature(SignatureData) CreateSessionResponseBuilder
	// WithServerSignatureBuilder adds ServerSignature (property field) which is build by the builder
	WithServerSignatureBuilder(func(SignatureDataBuilder) SignatureDataBuilder) CreateSessionResponseBuilder
	// WithMaxRequestMessageSize adds MaxRequestMessageSize (property field)
	WithMaxRequestMessageSize(uint32) CreateSessionResponseBuilder
	// Build builds the CreateSessionResponse or returns an error if something is wrong
	Build() (CreateSessionResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() CreateSessionResponse
}

// NewCreateSessionResponseBuilder() creates a CreateSessionResponseBuilder
func NewCreateSessionResponseBuilder() CreateSessionResponseBuilder {
	return &_CreateSessionResponseBuilder{_CreateSessionResponse: new(_CreateSessionResponse)}
}

type _CreateSessionResponseBuilder struct {
	*_CreateSessionResponse

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (CreateSessionResponseBuilder) = (*_CreateSessionResponseBuilder)(nil)

func (b *_CreateSessionResponseBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_CreateSessionResponseBuilder) WithMandatoryFields(responseHeader ResponseHeader, sessionId NodeId, authenticationToken NodeId, revisedSessionTimeout float64, serverNonce PascalByteString, serverCertificate PascalByteString, serverEndpoints []EndpointDescription, serverSoftwareCertificates []SignedSoftwareCertificate, serverSignature SignatureData, maxRequestMessageSize uint32) CreateSessionResponseBuilder {
	return b.WithResponseHeader(responseHeader).WithSessionId(sessionId).WithAuthenticationToken(authenticationToken).WithRevisedSessionTimeout(revisedSessionTimeout).WithServerNonce(serverNonce).WithServerCertificate(serverCertificate).WithServerEndpoints(serverEndpoints...).WithServerSoftwareCertificates(serverSoftwareCertificates...).WithServerSignature(serverSignature).WithMaxRequestMessageSize(maxRequestMessageSize)
}

func (b *_CreateSessionResponseBuilder) WithResponseHeader(responseHeader ResponseHeader) CreateSessionResponseBuilder {
	b.ResponseHeader = responseHeader
	return b
}

func (b *_CreateSessionResponseBuilder) WithResponseHeaderBuilder(builderSupplier func(ResponseHeaderBuilder) ResponseHeaderBuilder) CreateSessionResponseBuilder {
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

func (b *_CreateSessionResponseBuilder) WithSessionId(sessionId NodeId) CreateSessionResponseBuilder {
	b.SessionId = sessionId
	return b
}

func (b *_CreateSessionResponseBuilder) WithSessionIdBuilder(builderSupplier func(NodeIdBuilder) NodeIdBuilder) CreateSessionResponseBuilder {
	builder := builderSupplier(b.SessionId.CreateNodeIdBuilder())
	var err error
	b.SessionId, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "NodeIdBuilder failed"))
	}
	return b
}

func (b *_CreateSessionResponseBuilder) WithAuthenticationToken(authenticationToken NodeId) CreateSessionResponseBuilder {
	b.AuthenticationToken = authenticationToken
	return b
}

func (b *_CreateSessionResponseBuilder) WithAuthenticationTokenBuilder(builderSupplier func(NodeIdBuilder) NodeIdBuilder) CreateSessionResponseBuilder {
	builder := builderSupplier(b.AuthenticationToken.CreateNodeIdBuilder())
	var err error
	b.AuthenticationToken, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "NodeIdBuilder failed"))
	}
	return b
}

func (b *_CreateSessionResponseBuilder) WithRevisedSessionTimeout(revisedSessionTimeout float64) CreateSessionResponseBuilder {
	b.RevisedSessionTimeout = revisedSessionTimeout
	return b
}

func (b *_CreateSessionResponseBuilder) WithServerNonce(serverNonce PascalByteString) CreateSessionResponseBuilder {
	b.ServerNonce = serverNonce
	return b
}

func (b *_CreateSessionResponseBuilder) WithServerNonceBuilder(builderSupplier func(PascalByteStringBuilder) PascalByteStringBuilder) CreateSessionResponseBuilder {
	builder := builderSupplier(b.ServerNonce.CreatePascalByteStringBuilder())
	var err error
	b.ServerNonce, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalByteStringBuilder failed"))
	}
	return b
}

func (b *_CreateSessionResponseBuilder) WithServerCertificate(serverCertificate PascalByteString) CreateSessionResponseBuilder {
	b.ServerCertificate = serverCertificate
	return b
}

func (b *_CreateSessionResponseBuilder) WithServerCertificateBuilder(builderSupplier func(PascalByteStringBuilder) PascalByteStringBuilder) CreateSessionResponseBuilder {
	builder := builderSupplier(b.ServerCertificate.CreatePascalByteStringBuilder())
	var err error
	b.ServerCertificate, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalByteStringBuilder failed"))
	}
	return b
}

func (b *_CreateSessionResponseBuilder) WithServerEndpoints(serverEndpoints ...EndpointDescription) CreateSessionResponseBuilder {
	b.ServerEndpoints = serverEndpoints
	return b
}

func (b *_CreateSessionResponseBuilder) WithServerSoftwareCertificates(serverSoftwareCertificates ...SignedSoftwareCertificate) CreateSessionResponseBuilder {
	b.ServerSoftwareCertificates = serverSoftwareCertificates
	return b
}

func (b *_CreateSessionResponseBuilder) WithServerSignature(serverSignature SignatureData) CreateSessionResponseBuilder {
	b.ServerSignature = serverSignature
	return b
}

func (b *_CreateSessionResponseBuilder) WithServerSignatureBuilder(builderSupplier func(SignatureDataBuilder) SignatureDataBuilder) CreateSessionResponseBuilder {
	builder := builderSupplier(b.ServerSignature.CreateSignatureDataBuilder())
	var err error
	b.ServerSignature, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "SignatureDataBuilder failed"))
	}
	return b
}

func (b *_CreateSessionResponseBuilder) WithMaxRequestMessageSize(maxRequestMessageSize uint32) CreateSessionResponseBuilder {
	b.MaxRequestMessageSize = maxRequestMessageSize
	return b
}

func (b *_CreateSessionResponseBuilder) Build() (CreateSessionResponse, error) {
	if b.ResponseHeader == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'responseHeader' not set"))
	}
	if b.SessionId == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'sessionId' not set"))
	}
	if b.AuthenticationToken == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'authenticationToken' not set"))
	}
	if b.ServerNonce == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'serverNonce' not set"))
	}
	if b.ServerCertificate == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'serverCertificate' not set"))
	}
	if b.ServerSignature == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'serverSignature' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._CreateSessionResponse.deepCopy(), nil
}

func (b *_CreateSessionResponseBuilder) MustBuild() CreateSessionResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_CreateSessionResponseBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_CreateSessionResponseBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_CreateSessionResponseBuilder) DeepCopy() any {
	_copy := b.CreateCreateSessionResponseBuilder().(*_CreateSessionResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateCreateSessionResponseBuilder creates a CreateSessionResponseBuilder
func (b *_CreateSessionResponse) CreateCreateSessionResponseBuilder() CreateSessionResponseBuilder {
	if b == nil {
		return NewCreateSessionResponseBuilder()
	}
	return &_CreateSessionResponseBuilder{_CreateSessionResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_CreateSessionResponse) GetExtensionId() int32 {
	return int32(464)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_CreateSessionResponse) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CreateSessionResponse) GetResponseHeader() ResponseHeader {
	return m.ResponseHeader
}

func (m *_CreateSessionResponse) GetSessionId() NodeId {
	return m.SessionId
}

func (m *_CreateSessionResponse) GetAuthenticationToken() NodeId {
	return m.AuthenticationToken
}

func (m *_CreateSessionResponse) GetRevisedSessionTimeout() float64 {
	return m.RevisedSessionTimeout
}

func (m *_CreateSessionResponse) GetServerNonce() PascalByteString {
	return m.ServerNonce
}

func (m *_CreateSessionResponse) GetServerCertificate() PascalByteString {
	return m.ServerCertificate
}

func (m *_CreateSessionResponse) GetServerEndpoints() []EndpointDescription {
	return m.ServerEndpoints
}

func (m *_CreateSessionResponse) GetServerSoftwareCertificates() []SignedSoftwareCertificate {
	return m.ServerSoftwareCertificates
}

func (m *_CreateSessionResponse) GetServerSignature() SignatureData {
	return m.ServerSignature
}

func (m *_CreateSessionResponse) GetMaxRequestMessageSize() uint32 {
	return m.MaxRequestMessageSize
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastCreateSessionResponse(structType any) CreateSessionResponse {
	if casted, ok := structType.(CreateSessionResponse); ok {
		return casted
	}
	if casted, ok := structType.(*CreateSessionResponse); ok {
		return *casted
	}
	return nil
}

func (m *_CreateSessionResponse) GetTypeName() string {
	return "CreateSessionResponse"
}

func (m *_CreateSessionResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (responseHeader)
	lengthInBits += m.ResponseHeader.GetLengthInBits(ctx)

	// Simple field (sessionId)
	lengthInBits += m.SessionId.GetLengthInBits(ctx)

	// Simple field (authenticationToken)
	lengthInBits += m.AuthenticationToken.GetLengthInBits(ctx)

	// Simple field (revisedSessionTimeout)
	lengthInBits += 64

	// Simple field (serverNonce)
	lengthInBits += m.ServerNonce.GetLengthInBits(ctx)

	// Simple field (serverCertificate)
	lengthInBits += m.ServerCertificate.GetLengthInBits(ctx)

	// Implicit Field (noOfServerEndpoints)
	lengthInBits += 32

	// Array field
	if len(m.ServerEndpoints) > 0 {
		for _curItem, element := range m.ServerEndpoints {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.ServerEndpoints), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Implicit Field (noOfServerSoftwareCertificates)
	lengthInBits += 32

	// Array field
	if len(m.ServerSoftwareCertificates) > 0 {
		for _curItem, element := range m.ServerSoftwareCertificates {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.ServerSoftwareCertificates), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Simple field (serverSignature)
	lengthInBits += m.ServerSignature.GetLengthInBits(ctx)

	// Simple field (maxRequestMessageSize)
	lengthInBits += 32

	return lengthInBits
}

func (m *_CreateSessionResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_CreateSessionResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__createSessionResponse CreateSessionResponse, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CreateSessionResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CreateSessionResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	responseHeader, err := ReadSimpleField[ResponseHeader](ctx, "responseHeader", ReadComplex[ResponseHeader](ExtensionObjectDefinitionParseWithBufferProducer[ResponseHeader]((int32)(int32(394))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'responseHeader' field"))
	}
	m.ResponseHeader = responseHeader

	sessionId, err := ReadSimpleField[NodeId](ctx, "sessionId", ReadComplex[NodeId](NodeIdParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'sessionId' field"))
	}
	m.SessionId = sessionId

	authenticationToken, err := ReadSimpleField[NodeId](ctx, "authenticationToken", ReadComplex[NodeId](NodeIdParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'authenticationToken' field"))
	}
	m.AuthenticationToken = authenticationToken

	revisedSessionTimeout, err := ReadSimpleField(ctx, "revisedSessionTimeout", ReadDouble(readBuffer, uint8(64)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'revisedSessionTimeout' field"))
	}
	m.RevisedSessionTimeout = revisedSessionTimeout

	serverNonce, err := ReadSimpleField[PascalByteString](ctx, "serverNonce", ReadComplex[PascalByteString](PascalByteStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'serverNonce' field"))
	}
	m.ServerNonce = serverNonce

	serverCertificate, err := ReadSimpleField[PascalByteString](ctx, "serverCertificate", ReadComplex[PascalByteString](PascalByteStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'serverCertificate' field"))
	}
	m.ServerCertificate = serverCertificate

	noOfServerEndpoints, err := ReadImplicitField[int32](ctx, "noOfServerEndpoints", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfServerEndpoints' field"))
	}
	_ = noOfServerEndpoints

	serverEndpoints, err := ReadCountArrayField[EndpointDescription](ctx, "serverEndpoints", ReadComplex[EndpointDescription](ExtensionObjectDefinitionParseWithBufferProducer[EndpointDescription]((int32)(int32(314))), readBuffer), uint64(noOfServerEndpoints))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'serverEndpoints' field"))
	}
	m.ServerEndpoints = serverEndpoints

	noOfServerSoftwareCertificates, err := ReadImplicitField[int32](ctx, "noOfServerSoftwareCertificates", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfServerSoftwareCertificates' field"))
	}
	_ = noOfServerSoftwareCertificates

	serverSoftwareCertificates, err := ReadCountArrayField[SignedSoftwareCertificate](ctx, "serverSoftwareCertificates", ReadComplex[SignedSoftwareCertificate](ExtensionObjectDefinitionParseWithBufferProducer[SignedSoftwareCertificate]((int32)(int32(346))), readBuffer), uint64(noOfServerSoftwareCertificates))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'serverSoftwareCertificates' field"))
	}
	m.ServerSoftwareCertificates = serverSoftwareCertificates

	serverSignature, err := ReadSimpleField[SignatureData](ctx, "serverSignature", ReadComplex[SignatureData](ExtensionObjectDefinitionParseWithBufferProducer[SignatureData]((int32)(int32(458))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'serverSignature' field"))
	}
	m.ServerSignature = serverSignature

	maxRequestMessageSize, err := ReadSimpleField(ctx, "maxRequestMessageSize", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'maxRequestMessageSize' field"))
	}
	m.MaxRequestMessageSize = maxRequestMessageSize

	if closeErr := readBuffer.CloseContext("CreateSessionResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CreateSessionResponse")
	}

	return m, nil
}

func (m *_CreateSessionResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CreateSessionResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CreateSessionResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CreateSessionResponse")
		}

		if err := WriteSimpleField[ResponseHeader](ctx, "responseHeader", m.GetResponseHeader(), WriteComplex[ResponseHeader](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'responseHeader' field")
		}

		if err := WriteSimpleField[NodeId](ctx, "sessionId", m.GetSessionId(), WriteComplex[NodeId](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'sessionId' field")
		}

		if err := WriteSimpleField[NodeId](ctx, "authenticationToken", m.GetAuthenticationToken(), WriteComplex[NodeId](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'authenticationToken' field")
		}

		if err := WriteSimpleField[float64](ctx, "revisedSessionTimeout", m.GetRevisedSessionTimeout(), WriteDouble(writeBuffer, 64)); err != nil {
			return errors.Wrap(err, "Error serializing 'revisedSessionTimeout' field")
		}

		if err := WriteSimpleField[PascalByteString](ctx, "serverNonce", m.GetServerNonce(), WriteComplex[PascalByteString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'serverNonce' field")
		}

		if err := WriteSimpleField[PascalByteString](ctx, "serverCertificate", m.GetServerCertificate(), WriteComplex[PascalByteString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'serverCertificate' field")
		}
		noOfServerEndpoints := int32(utils.InlineIf(bool((m.GetServerEndpoints()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetServerEndpoints()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfServerEndpoints", noOfServerEndpoints, WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfServerEndpoints' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "serverEndpoints", m.GetServerEndpoints(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'serverEndpoints' field")
		}
		noOfServerSoftwareCertificates := int32(utils.InlineIf(bool((m.GetServerSoftwareCertificates()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetServerSoftwareCertificates()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfServerSoftwareCertificates", noOfServerSoftwareCertificates, WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfServerSoftwareCertificates' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "serverSoftwareCertificates", m.GetServerSoftwareCertificates(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'serverSoftwareCertificates' field")
		}

		if err := WriteSimpleField[SignatureData](ctx, "serverSignature", m.GetServerSignature(), WriteComplex[SignatureData](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'serverSignature' field")
		}

		if err := WriteSimpleField[uint32](ctx, "maxRequestMessageSize", m.GetMaxRequestMessageSize(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'maxRequestMessageSize' field")
		}

		if popErr := writeBuffer.PopContext("CreateSessionResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CreateSessionResponse")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_CreateSessionResponse) IsCreateSessionResponse() {}

func (m *_CreateSessionResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_CreateSessionResponse) deepCopy() *_CreateSessionResponse {
	if m == nil {
		return nil
	}
	_CreateSessionResponseCopy := &_CreateSessionResponse{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.ResponseHeader.DeepCopy().(ResponseHeader),
		m.SessionId.DeepCopy().(NodeId),
		m.AuthenticationToken.DeepCopy().(NodeId),
		m.RevisedSessionTimeout,
		m.ServerNonce.DeepCopy().(PascalByteString),
		m.ServerCertificate.DeepCopy().(PascalByteString),
		utils.DeepCopySlice[EndpointDescription, EndpointDescription](m.ServerEndpoints),
		utils.DeepCopySlice[SignedSoftwareCertificate, SignedSoftwareCertificate](m.ServerSoftwareCertificates),
		m.ServerSignature.DeepCopy().(SignatureData),
		m.MaxRequestMessageSize,
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _CreateSessionResponseCopy
}

func (m *_CreateSessionResponse) String() string {
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
