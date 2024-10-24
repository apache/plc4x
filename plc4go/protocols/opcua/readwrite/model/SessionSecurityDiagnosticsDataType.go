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

// SessionSecurityDiagnosticsDataType is the corresponding interface of SessionSecurityDiagnosticsDataType
type SessionSecurityDiagnosticsDataType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetSessionId returns SessionId (property field)
	GetSessionId() NodeId
	// GetClientUserIdOfSession returns ClientUserIdOfSession (property field)
	GetClientUserIdOfSession() PascalString
	// GetClientUserIdHistory returns ClientUserIdHistory (property field)
	GetClientUserIdHistory() []PascalString
	// GetAuthenticationMechanism returns AuthenticationMechanism (property field)
	GetAuthenticationMechanism() PascalString
	// GetEncoding returns Encoding (property field)
	GetEncoding() PascalString
	// GetTransportProtocol returns TransportProtocol (property field)
	GetTransportProtocol() PascalString
	// GetSecurityMode returns SecurityMode (property field)
	GetSecurityMode() MessageSecurityMode
	// GetSecurityPolicyUri returns SecurityPolicyUri (property field)
	GetSecurityPolicyUri() PascalString
	// GetClientCertificate returns ClientCertificate (property field)
	GetClientCertificate() PascalByteString
	// IsSessionSecurityDiagnosticsDataType is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsSessionSecurityDiagnosticsDataType()
	// CreateBuilder creates a SessionSecurityDiagnosticsDataTypeBuilder
	CreateSessionSecurityDiagnosticsDataTypeBuilder() SessionSecurityDiagnosticsDataTypeBuilder
}

// _SessionSecurityDiagnosticsDataType is the data-structure of this message
type _SessionSecurityDiagnosticsDataType struct {
	ExtensionObjectDefinitionContract
	SessionId               NodeId
	ClientUserIdOfSession   PascalString
	ClientUserIdHistory     []PascalString
	AuthenticationMechanism PascalString
	Encoding                PascalString
	TransportProtocol       PascalString
	SecurityMode            MessageSecurityMode
	SecurityPolicyUri       PascalString
	ClientCertificate       PascalByteString
}

var _ SessionSecurityDiagnosticsDataType = (*_SessionSecurityDiagnosticsDataType)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_SessionSecurityDiagnosticsDataType)(nil)

// NewSessionSecurityDiagnosticsDataType factory function for _SessionSecurityDiagnosticsDataType
func NewSessionSecurityDiagnosticsDataType(sessionId NodeId, clientUserIdOfSession PascalString, clientUserIdHistory []PascalString, authenticationMechanism PascalString, encoding PascalString, transportProtocol PascalString, securityMode MessageSecurityMode, securityPolicyUri PascalString, clientCertificate PascalByteString) *_SessionSecurityDiagnosticsDataType {
	if sessionId == nil {
		panic("sessionId of type NodeId for SessionSecurityDiagnosticsDataType must not be nil")
	}
	if clientUserIdOfSession == nil {
		panic("clientUserIdOfSession of type PascalString for SessionSecurityDiagnosticsDataType must not be nil")
	}
	if authenticationMechanism == nil {
		panic("authenticationMechanism of type PascalString for SessionSecurityDiagnosticsDataType must not be nil")
	}
	if encoding == nil {
		panic("encoding of type PascalString for SessionSecurityDiagnosticsDataType must not be nil")
	}
	if transportProtocol == nil {
		panic("transportProtocol of type PascalString for SessionSecurityDiagnosticsDataType must not be nil")
	}
	if securityPolicyUri == nil {
		panic("securityPolicyUri of type PascalString for SessionSecurityDiagnosticsDataType must not be nil")
	}
	if clientCertificate == nil {
		panic("clientCertificate of type PascalByteString for SessionSecurityDiagnosticsDataType must not be nil")
	}
	_result := &_SessionSecurityDiagnosticsDataType{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		SessionId:                         sessionId,
		ClientUserIdOfSession:             clientUserIdOfSession,
		ClientUserIdHistory:               clientUserIdHistory,
		AuthenticationMechanism:           authenticationMechanism,
		Encoding:                          encoding,
		TransportProtocol:                 transportProtocol,
		SecurityMode:                      securityMode,
		SecurityPolicyUri:                 securityPolicyUri,
		ClientCertificate:                 clientCertificate,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// SessionSecurityDiagnosticsDataTypeBuilder is a builder for SessionSecurityDiagnosticsDataType
type SessionSecurityDiagnosticsDataTypeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(sessionId NodeId, clientUserIdOfSession PascalString, clientUserIdHistory []PascalString, authenticationMechanism PascalString, encoding PascalString, transportProtocol PascalString, securityMode MessageSecurityMode, securityPolicyUri PascalString, clientCertificate PascalByteString) SessionSecurityDiagnosticsDataTypeBuilder
	// WithSessionId adds SessionId (property field)
	WithSessionId(NodeId) SessionSecurityDiagnosticsDataTypeBuilder
	// WithSessionIdBuilder adds SessionId (property field) which is build by the builder
	WithSessionIdBuilder(func(NodeIdBuilder) NodeIdBuilder) SessionSecurityDiagnosticsDataTypeBuilder
	// WithClientUserIdOfSession adds ClientUserIdOfSession (property field)
	WithClientUserIdOfSession(PascalString) SessionSecurityDiagnosticsDataTypeBuilder
	// WithClientUserIdOfSessionBuilder adds ClientUserIdOfSession (property field) which is build by the builder
	WithClientUserIdOfSessionBuilder(func(PascalStringBuilder) PascalStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder
	// WithClientUserIdHistory adds ClientUserIdHistory (property field)
	WithClientUserIdHistory(...PascalString) SessionSecurityDiagnosticsDataTypeBuilder
	// WithAuthenticationMechanism adds AuthenticationMechanism (property field)
	WithAuthenticationMechanism(PascalString) SessionSecurityDiagnosticsDataTypeBuilder
	// WithAuthenticationMechanismBuilder adds AuthenticationMechanism (property field) which is build by the builder
	WithAuthenticationMechanismBuilder(func(PascalStringBuilder) PascalStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder
	// WithEncoding adds Encoding (property field)
	WithEncoding(PascalString) SessionSecurityDiagnosticsDataTypeBuilder
	// WithEncodingBuilder adds Encoding (property field) which is build by the builder
	WithEncodingBuilder(func(PascalStringBuilder) PascalStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder
	// WithTransportProtocol adds TransportProtocol (property field)
	WithTransportProtocol(PascalString) SessionSecurityDiagnosticsDataTypeBuilder
	// WithTransportProtocolBuilder adds TransportProtocol (property field) which is build by the builder
	WithTransportProtocolBuilder(func(PascalStringBuilder) PascalStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder
	// WithSecurityMode adds SecurityMode (property field)
	WithSecurityMode(MessageSecurityMode) SessionSecurityDiagnosticsDataTypeBuilder
	// WithSecurityPolicyUri adds SecurityPolicyUri (property field)
	WithSecurityPolicyUri(PascalString) SessionSecurityDiagnosticsDataTypeBuilder
	// WithSecurityPolicyUriBuilder adds SecurityPolicyUri (property field) which is build by the builder
	WithSecurityPolicyUriBuilder(func(PascalStringBuilder) PascalStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder
	// WithClientCertificate adds ClientCertificate (property field)
	WithClientCertificate(PascalByteString) SessionSecurityDiagnosticsDataTypeBuilder
	// WithClientCertificateBuilder adds ClientCertificate (property field) which is build by the builder
	WithClientCertificateBuilder(func(PascalByteStringBuilder) PascalByteStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder
	// Build builds the SessionSecurityDiagnosticsDataType or returns an error if something is wrong
	Build() (SessionSecurityDiagnosticsDataType, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() SessionSecurityDiagnosticsDataType
}

// NewSessionSecurityDiagnosticsDataTypeBuilder() creates a SessionSecurityDiagnosticsDataTypeBuilder
func NewSessionSecurityDiagnosticsDataTypeBuilder() SessionSecurityDiagnosticsDataTypeBuilder {
	return &_SessionSecurityDiagnosticsDataTypeBuilder{_SessionSecurityDiagnosticsDataType: new(_SessionSecurityDiagnosticsDataType)}
}

type _SessionSecurityDiagnosticsDataTypeBuilder struct {
	*_SessionSecurityDiagnosticsDataType

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (SessionSecurityDiagnosticsDataTypeBuilder) = (*_SessionSecurityDiagnosticsDataTypeBuilder)(nil)

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithMandatoryFields(sessionId NodeId, clientUserIdOfSession PascalString, clientUserIdHistory []PascalString, authenticationMechanism PascalString, encoding PascalString, transportProtocol PascalString, securityMode MessageSecurityMode, securityPolicyUri PascalString, clientCertificate PascalByteString) SessionSecurityDiagnosticsDataTypeBuilder {
	return b.WithSessionId(sessionId).WithClientUserIdOfSession(clientUserIdOfSession).WithClientUserIdHistory(clientUserIdHistory...).WithAuthenticationMechanism(authenticationMechanism).WithEncoding(encoding).WithTransportProtocol(transportProtocol).WithSecurityMode(securityMode).WithSecurityPolicyUri(securityPolicyUri).WithClientCertificate(clientCertificate)
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithSessionId(sessionId NodeId) SessionSecurityDiagnosticsDataTypeBuilder {
	b.SessionId = sessionId
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithSessionIdBuilder(builderSupplier func(NodeIdBuilder) NodeIdBuilder) SessionSecurityDiagnosticsDataTypeBuilder {
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

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithClientUserIdOfSession(clientUserIdOfSession PascalString) SessionSecurityDiagnosticsDataTypeBuilder {
	b.ClientUserIdOfSession = clientUserIdOfSession
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithClientUserIdOfSessionBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder {
	builder := builderSupplier(b.ClientUserIdOfSession.CreatePascalStringBuilder())
	var err error
	b.ClientUserIdOfSession, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalStringBuilder failed"))
	}
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithClientUserIdHistory(clientUserIdHistory ...PascalString) SessionSecurityDiagnosticsDataTypeBuilder {
	b.ClientUserIdHistory = clientUserIdHistory
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithAuthenticationMechanism(authenticationMechanism PascalString) SessionSecurityDiagnosticsDataTypeBuilder {
	b.AuthenticationMechanism = authenticationMechanism
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithAuthenticationMechanismBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder {
	builder := builderSupplier(b.AuthenticationMechanism.CreatePascalStringBuilder())
	var err error
	b.AuthenticationMechanism, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalStringBuilder failed"))
	}
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithEncoding(encoding PascalString) SessionSecurityDiagnosticsDataTypeBuilder {
	b.Encoding = encoding
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithEncodingBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder {
	builder := builderSupplier(b.Encoding.CreatePascalStringBuilder())
	var err error
	b.Encoding, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalStringBuilder failed"))
	}
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithTransportProtocol(transportProtocol PascalString) SessionSecurityDiagnosticsDataTypeBuilder {
	b.TransportProtocol = transportProtocol
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithTransportProtocolBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder {
	builder := builderSupplier(b.TransportProtocol.CreatePascalStringBuilder())
	var err error
	b.TransportProtocol, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalStringBuilder failed"))
	}
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithSecurityMode(securityMode MessageSecurityMode) SessionSecurityDiagnosticsDataTypeBuilder {
	b.SecurityMode = securityMode
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithSecurityPolicyUri(securityPolicyUri PascalString) SessionSecurityDiagnosticsDataTypeBuilder {
	b.SecurityPolicyUri = securityPolicyUri
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithSecurityPolicyUriBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder {
	builder := builderSupplier(b.SecurityPolicyUri.CreatePascalStringBuilder())
	var err error
	b.SecurityPolicyUri, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalStringBuilder failed"))
	}
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithClientCertificate(clientCertificate PascalByteString) SessionSecurityDiagnosticsDataTypeBuilder {
	b.ClientCertificate = clientCertificate
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) WithClientCertificateBuilder(builderSupplier func(PascalByteStringBuilder) PascalByteStringBuilder) SessionSecurityDiagnosticsDataTypeBuilder {
	builder := builderSupplier(b.ClientCertificate.CreatePascalByteStringBuilder())
	var err error
	b.ClientCertificate, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalByteStringBuilder failed"))
	}
	return b
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) Build() (SessionSecurityDiagnosticsDataType, error) {
	if b.SessionId == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'sessionId' not set"))
	}
	if b.ClientUserIdOfSession == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'clientUserIdOfSession' not set"))
	}
	if b.AuthenticationMechanism == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'authenticationMechanism' not set"))
	}
	if b.Encoding == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'encoding' not set"))
	}
	if b.TransportProtocol == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'transportProtocol' not set"))
	}
	if b.SecurityPolicyUri == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'securityPolicyUri' not set"))
	}
	if b.ClientCertificate == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'clientCertificate' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._SessionSecurityDiagnosticsDataType.deepCopy(), nil
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) MustBuild() SessionSecurityDiagnosticsDataType {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_SessionSecurityDiagnosticsDataTypeBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_SessionSecurityDiagnosticsDataTypeBuilder) DeepCopy() any {
	_copy := b.CreateSessionSecurityDiagnosticsDataTypeBuilder().(*_SessionSecurityDiagnosticsDataTypeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateSessionSecurityDiagnosticsDataTypeBuilder creates a SessionSecurityDiagnosticsDataTypeBuilder
func (b *_SessionSecurityDiagnosticsDataType) CreateSessionSecurityDiagnosticsDataTypeBuilder() SessionSecurityDiagnosticsDataTypeBuilder {
	if b == nil {
		return NewSessionSecurityDiagnosticsDataTypeBuilder()
	}
	return &_SessionSecurityDiagnosticsDataTypeBuilder{_SessionSecurityDiagnosticsDataType: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_SessionSecurityDiagnosticsDataType) GetExtensionId() int32 {
	return int32(870)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_SessionSecurityDiagnosticsDataType) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_SessionSecurityDiagnosticsDataType) GetSessionId() NodeId {
	return m.SessionId
}

func (m *_SessionSecurityDiagnosticsDataType) GetClientUserIdOfSession() PascalString {
	return m.ClientUserIdOfSession
}

func (m *_SessionSecurityDiagnosticsDataType) GetClientUserIdHistory() []PascalString {
	return m.ClientUserIdHistory
}

func (m *_SessionSecurityDiagnosticsDataType) GetAuthenticationMechanism() PascalString {
	return m.AuthenticationMechanism
}

func (m *_SessionSecurityDiagnosticsDataType) GetEncoding() PascalString {
	return m.Encoding
}

func (m *_SessionSecurityDiagnosticsDataType) GetTransportProtocol() PascalString {
	return m.TransportProtocol
}

func (m *_SessionSecurityDiagnosticsDataType) GetSecurityMode() MessageSecurityMode {
	return m.SecurityMode
}

func (m *_SessionSecurityDiagnosticsDataType) GetSecurityPolicyUri() PascalString {
	return m.SecurityPolicyUri
}

func (m *_SessionSecurityDiagnosticsDataType) GetClientCertificate() PascalByteString {
	return m.ClientCertificate
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastSessionSecurityDiagnosticsDataType(structType any) SessionSecurityDiagnosticsDataType {
	if casted, ok := structType.(SessionSecurityDiagnosticsDataType); ok {
		return casted
	}
	if casted, ok := structType.(*SessionSecurityDiagnosticsDataType); ok {
		return *casted
	}
	return nil
}

func (m *_SessionSecurityDiagnosticsDataType) GetTypeName() string {
	return "SessionSecurityDiagnosticsDataType"
}

func (m *_SessionSecurityDiagnosticsDataType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (sessionId)
	lengthInBits += m.SessionId.GetLengthInBits(ctx)

	// Simple field (clientUserIdOfSession)
	lengthInBits += m.ClientUserIdOfSession.GetLengthInBits(ctx)

	// Implicit Field (noOfClientUserIdHistory)
	lengthInBits += 32

	// Array field
	if len(m.ClientUserIdHistory) > 0 {
		for _curItem, element := range m.ClientUserIdHistory {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.ClientUserIdHistory), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Simple field (authenticationMechanism)
	lengthInBits += m.AuthenticationMechanism.GetLengthInBits(ctx)

	// Simple field (encoding)
	lengthInBits += m.Encoding.GetLengthInBits(ctx)

	// Simple field (transportProtocol)
	lengthInBits += m.TransportProtocol.GetLengthInBits(ctx)

	// Simple field (securityMode)
	lengthInBits += 32

	// Simple field (securityPolicyUri)
	lengthInBits += m.SecurityPolicyUri.GetLengthInBits(ctx)

	// Simple field (clientCertificate)
	lengthInBits += m.ClientCertificate.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_SessionSecurityDiagnosticsDataType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_SessionSecurityDiagnosticsDataType) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__sessionSecurityDiagnosticsDataType SessionSecurityDiagnosticsDataType, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SessionSecurityDiagnosticsDataType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SessionSecurityDiagnosticsDataType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	sessionId, err := ReadSimpleField[NodeId](ctx, "sessionId", ReadComplex[NodeId](NodeIdParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'sessionId' field"))
	}
	m.SessionId = sessionId

	clientUserIdOfSession, err := ReadSimpleField[PascalString](ctx, "clientUserIdOfSession", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'clientUserIdOfSession' field"))
	}
	m.ClientUserIdOfSession = clientUserIdOfSession

	noOfClientUserIdHistory, err := ReadImplicitField[int32](ctx, "noOfClientUserIdHistory", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfClientUserIdHistory' field"))
	}
	_ = noOfClientUserIdHistory

	clientUserIdHistory, err := ReadCountArrayField[PascalString](ctx, "clientUserIdHistory", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer), uint64(noOfClientUserIdHistory))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'clientUserIdHistory' field"))
	}
	m.ClientUserIdHistory = clientUserIdHistory

	authenticationMechanism, err := ReadSimpleField[PascalString](ctx, "authenticationMechanism", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'authenticationMechanism' field"))
	}
	m.AuthenticationMechanism = authenticationMechanism

	encoding, err := ReadSimpleField[PascalString](ctx, "encoding", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'encoding' field"))
	}
	m.Encoding = encoding

	transportProtocol, err := ReadSimpleField[PascalString](ctx, "transportProtocol", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'transportProtocol' field"))
	}
	m.TransportProtocol = transportProtocol

	securityMode, err := ReadEnumField[MessageSecurityMode](ctx, "securityMode", "MessageSecurityMode", ReadEnum(MessageSecurityModeByValue, ReadUnsignedInt(readBuffer, uint8(32))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'securityMode' field"))
	}
	m.SecurityMode = securityMode

	securityPolicyUri, err := ReadSimpleField[PascalString](ctx, "securityPolicyUri", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'securityPolicyUri' field"))
	}
	m.SecurityPolicyUri = securityPolicyUri

	clientCertificate, err := ReadSimpleField[PascalByteString](ctx, "clientCertificate", ReadComplex[PascalByteString](PascalByteStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'clientCertificate' field"))
	}
	m.ClientCertificate = clientCertificate

	if closeErr := readBuffer.CloseContext("SessionSecurityDiagnosticsDataType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SessionSecurityDiagnosticsDataType")
	}

	return m, nil
}

func (m *_SessionSecurityDiagnosticsDataType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_SessionSecurityDiagnosticsDataType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("SessionSecurityDiagnosticsDataType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for SessionSecurityDiagnosticsDataType")
		}

		if err := WriteSimpleField[NodeId](ctx, "sessionId", m.GetSessionId(), WriteComplex[NodeId](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'sessionId' field")
		}

		if err := WriteSimpleField[PascalString](ctx, "clientUserIdOfSession", m.GetClientUserIdOfSession(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'clientUserIdOfSession' field")
		}
		noOfClientUserIdHistory := int32(utils.InlineIf(bool((m.GetClientUserIdHistory()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetClientUserIdHistory()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfClientUserIdHistory", noOfClientUserIdHistory, WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfClientUserIdHistory' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "clientUserIdHistory", m.GetClientUserIdHistory(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'clientUserIdHistory' field")
		}

		if err := WriteSimpleField[PascalString](ctx, "authenticationMechanism", m.GetAuthenticationMechanism(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'authenticationMechanism' field")
		}

		if err := WriteSimpleField[PascalString](ctx, "encoding", m.GetEncoding(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'encoding' field")
		}

		if err := WriteSimpleField[PascalString](ctx, "transportProtocol", m.GetTransportProtocol(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'transportProtocol' field")
		}

		if err := WriteSimpleEnumField[MessageSecurityMode](ctx, "securityMode", "MessageSecurityMode", m.GetSecurityMode(), WriteEnum[MessageSecurityMode, uint32](MessageSecurityMode.GetValue, MessageSecurityMode.PLC4XEnumName, WriteUnsignedInt(writeBuffer, 32))); err != nil {
			return errors.Wrap(err, "Error serializing 'securityMode' field")
		}

		if err := WriteSimpleField[PascalString](ctx, "securityPolicyUri", m.GetSecurityPolicyUri(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'securityPolicyUri' field")
		}

		if err := WriteSimpleField[PascalByteString](ctx, "clientCertificate", m.GetClientCertificate(), WriteComplex[PascalByteString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'clientCertificate' field")
		}

		if popErr := writeBuffer.PopContext("SessionSecurityDiagnosticsDataType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for SessionSecurityDiagnosticsDataType")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_SessionSecurityDiagnosticsDataType) IsSessionSecurityDiagnosticsDataType() {}

func (m *_SessionSecurityDiagnosticsDataType) DeepCopy() any {
	return m.deepCopy()
}

func (m *_SessionSecurityDiagnosticsDataType) deepCopy() *_SessionSecurityDiagnosticsDataType {
	if m == nil {
		return nil
	}
	_SessionSecurityDiagnosticsDataTypeCopy := &_SessionSecurityDiagnosticsDataType{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.SessionId.DeepCopy().(NodeId),
		m.ClientUserIdOfSession.DeepCopy().(PascalString),
		utils.DeepCopySlice[PascalString, PascalString](m.ClientUserIdHistory),
		m.AuthenticationMechanism.DeepCopy().(PascalString),
		m.Encoding.DeepCopy().(PascalString),
		m.TransportProtocol.DeepCopy().(PascalString),
		m.SecurityMode,
		m.SecurityPolicyUri.DeepCopy().(PascalString),
		m.ClientCertificate.DeepCopy().(PascalByteString),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _SessionSecurityDiagnosticsDataTypeCopy
}

func (m *_SessionSecurityDiagnosticsDataType) String() string {
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
