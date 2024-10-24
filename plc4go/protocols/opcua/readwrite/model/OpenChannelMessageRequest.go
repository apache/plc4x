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

// OpenChannelMessageRequest is the corresponding interface of OpenChannelMessageRequest
type OpenChannelMessageRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	OpenChannelMessage
	// GetSecureChannelId returns SecureChannelId (property field)
	GetSecureChannelId() int32
	// GetEndpoint returns Endpoint (property field)
	GetEndpoint() PascalString
	// GetSenderCertificate returns SenderCertificate (property field)
	GetSenderCertificate() PascalByteString
	// GetReceiverCertificateThumbprint returns ReceiverCertificateThumbprint (property field)
	GetReceiverCertificateThumbprint() PascalByteString
	// IsOpenChannelMessageRequest is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsOpenChannelMessageRequest()
	// CreateBuilder creates a OpenChannelMessageRequestBuilder
	CreateOpenChannelMessageRequestBuilder() OpenChannelMessageRequestBuilder
}

// _OpenChannelMessageRequest is the data-structure of this message
type _OpenChannelMessageRequest struct {
	OpenChannelMessageContract
	SecureChannelId               int32
	Endpoint                      PascalString
	SenderCertificate             PascalByteString
	ReceiverCertificateThumbprint PascalByteString
}

var _ OpenChannelMessageRequest = (*_OpenChannelMessageRequest)(nil)
var _ OpenChannelMessageRequirements = (*_OpenChannelMessageRequest)(nil)

// NewOpenChannelMessageRequest factory function for _OpenChannelMessageRequest
func NewOpenChannelMessageRequest(secureChannelId int32, endpoint PascalString, senderCertificate PascalByteString, receiverCertificateThumbprint PascalByteString) *_OpenChannelMessageRequest {
	if endpoint == nil {
		panic("endpoint of type PascalString for OpenChannelMessageRequest must not be nil")
	}
	if senderCertificate == nil {
		panic("senderCertificate of type PascalByteString for OpenChannelMessageRequest must not be nil")
	}
	if receiverCertificateThumbprint == nil {
		panic("receiverCertificateThumbprint of type PascalByteString for OpenChannelMessageRequest must not be nil")
	}
	_result := &_OpenChannelMessageRequest{
		OpenChannelMessageContract:    NewOpenChannelMessage(),
		SecureChannelId:               secureChannelId,
		Endpoint:                      endpoint,
		SenderCertificate:             senderCertificate,
		ReceiverCertificateThumbprint: receiverCertificateThumbprint,
	}
	_result.OpenChannelMessageContract.(*_OpenChannelMessage)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// OpenChannelMessageRequestBuilder is a builder for OpenChannelMessageRequest
type OpenChannelMessageRequestBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(secureChannelId int32, endpoint PascalString, senderCertificate PascalByteString, receiverCertificateThumbprint PascalByteString) OpenChannelMessageRequestBuilder
	// WithSecureChannelId adds SecureChannelId (property field)
	WithSecureChannelId(int32) OpenChannelMessageRequestBuilder
	// WithEndpoint adds Endpoint (property field)
	WithEndpoint(PascalString) OpenChannelMessageRequestBuilder
	// WithEndpointBuilder adds Endpoint (property field) which is build by the builder
	WithEndpointBuilder(func(PascalStringBuilder) PascalStringBuilder) OpenChannelMessageRequestBuilder
	// WithSenderCertificate adds SenderCertificate (property field)
	WithSenderCertificate(PascalByteString) OpenChannelMessageRequestBuilder
	// WithSenderCertificateBuilder adds SenderCertificate (property field) which is build by the builder
	WithSenderCertificateBuilder(func(PascalByteStringBuilder) PascalByteStringBuilder) OpenChannelMessageRequestBuilder
	// WithReceiverCertificateThumbprint adds ReceiverCertificateThumbprint (property field)
	WithReceiverCertificateThumbprint(PascalByteString) OpenChannelMessageRequestBuilder
	// WithReceiverCertificateThumbprintBuilder adds ReceiverCertificateThumbprint (property field) which is build by the builder
	WithReceiverCertificateThumbprintBuilder(func(PascalByteStringBuilder) PascalByteStringBuilder) OpenChannelMessageRequestBuilder
	// Build builds the OpenChannelMessageRequest or returns an error if something is wrong
	Build() (OpenChannelMessageRequest, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() OpenChannelMessageRequest
}

// NewOpenChannelMessageRequestBuilder() creates a OpenChannelMessageRequestBuilder
func NewOpenChannelMessageRequestBuilder() OpenChannelMessageRequestBuilder {
	return &_OpenChannelMessageRequestBuilder{_OpenChannelMessageRequest: new(_OpenChannelMessageRequest)}
}

type _OpenChannelMessageRequestBuilder struct {
	*_OpenChannelMessageRequest

	parentBuilder *_OpenChannelMessageBuilder

	err *utils.MultiError
}

var _ (OpenChannelMessageRequestBuilder) = (*_OpenChannelMessageRequestBuilder)(nil)

func (b *_OpenChannelMessageRequestBuilder) setParent(contract OpenChannelMessageContract) {
	b.OpenChannelMessageContract = contract
}

func (b *_OpenChannelMessageRequestBuilder) WithMandatoryFields(secureChannelId int32, endpoint PascalString, senderCertificate PascalByteString, receiverCertificateThumbprint PascalByteString) OpenChannelMessageRequestBuilder {
	return b.WithSecureChannelId(secureChannelId).WithEndpoint(endpoint).WithSenderCertificate(senderCertificate).WithReceiverCertificateThumbprint(receiverCertificateThumbprint)
}

func (b *_OpenChannelMessageRequestBuilder) WithSecureChannelId(secureChannelId int32) OpenChannelMessageRequestBuilder {
	b.SecureChannelId = secureChannelId
	return b
}

func (b *_OpenChannelMessageRequestBuilder) WithEndpoint(endpoint PascalString) OpenChannelMessageRequestBuilder {
	b.Endpoint = endpoint
	return b
}

func (b *_OpenChannelMessageRequestBuilder) WithEndpointBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) OpenChannelMessageRequestBuilder {
	builder := builderSupplier(b.Endpoint.CreatePascalStringBuilder())
	var err error
	b.Endpoint, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalStringBuilder failed"))
	}
	return b
}

func (b *_OpenChannelMessageRequestBuilder) WithSenderCertificate(senderCertificate PascalByteString) OpenChannelMessageRequestBuilder {
	b.SenderCertificate = senderCertificate
	return b
}

func (b *_OpenChannelMessageRequestBuilder) WithSenderCertificateBuilder(builderSupplier func(PascalByteStringBuilder) PascalByteStringBuilder) OpenChannelMessageRequestBuilder {
	builder := builderSupplier(b.SenderCertificate.CreatePascalByteStringBuilder())
	var err error
	b.SenderCertificate, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalByteStringBuilder failed"))
	}
	return b
}

func (b *_OpenChannelMessageRequestBuilder) WithReceiverCertificateThumbprint(receiverCertificateThumbprint PascalByteString) OpenChannelMessageRequestBuilder {
	b.ReceiverCertificateThumbprint = receiverCertificateThumbprint
	return b
}

func (b *_OpenChannelMessageRequestBuilder) WithReceiverCertificateThumbprintBuilder(builderSupplier func(PascalByteStringBuilder) PascalByteStringBuilder) OpenChannelMessageRequestBuilder {
	builder := builderSupplier(b.ReceiverCertificateThumbprint.CreatePascalByteStringBuilder())
	var err error
	b.ReceiverCertificateThumbprint, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalByteStringBuilder failed"))
	}
	return b
}

func (b *_OpenChannelMessageRequestBuilder) Build() (OpenChannelMessageRequest, error) {
	if b.Endpoint == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'endpoint' not set"))
	}
	if b.SenderCertificate == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'senderCertificate' not set"))
	}
	if b.ReceiverCertificateThumbprint == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'receiverCertificateThumbprint' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._OpenChannelMessageRequest.deepCopy(), nil
}

func (b *_OpenChannelMessageRequestBuilder) MustBuild() OpenChannelMessageRequest {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_OpenChannelMessageRequestBuilder) Done() OpenChannelMessageBuilder {
	return b.parentBuilder
}

func (b *_OpenChannelMessageRequestBuilder) buildForOpenChannelMessage() (OpenChannelMessage, error) {
	return b.Build()
}

func (b *_OpenChannelMessageRequestBuilder) DeepCopy() any {
	_copy := b.CreateOpenChannelMessageRequestBuilder().(*_OpenChannelMessageRequestBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateOpenChannelMessageRequestBuilder creates a OpenChannelMessageRequestBuilder
func (b *_OpenChannelMessageRequest) CreateOpenChannelMessageRequestBuilder() OpenChannelMessageRequestBuilder {
	if b == nil {
		return NewOpenChannelMessageRequestBuilder()
	}
	return &_OpenChannelMessageRequestBuilder{_OpenChannelMessageRequest: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_OpenChannelMessageRequest) GetResponse() bool {
	return bool(false)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_OpenChannelMessageRequest) GetParent() OpenChannelMessageContract {
	return m.OpenChannelMessageContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_OpenChannelMessageRequest) GetSecureChannelId() int32 {
	return m.SecureChannelId
}

func (m *_OpenChannelMessageRequest) GetEndpoint() PascalString {
	return m.Endpoint
}

func (m *_OpenChannelMessageRequest) GetSenderCertificate() PascalByteString {
	return m.SenderCertificate
}

func (m *_OpenChannelMessageRequest) GetReceiverCertificateThumbprint() PascalByteString {
	return m.ReceiverCertificateThumbprint
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastOpenChannelMessageRequest(structType any) OpenChannelMessageRequest {
	if casted, ok := structType.(OpenChannelMessageRequest); ok {
		return casted
	}
	if casted, ok := structType.(*OpenChannelMessageRequest); ok {
		return *casted
	}
	return nil
}

func (m *_OpenChannelMessageRequest) GetTypeName() string {
	return "OpenChannelMessageRequest"
}

func (m *_OpenChannelMessageRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.OpenChannelMessageContract.(*_OpenChannelMessage).GetLengthInBits(ctx))

	// Simple field (secureChannelId)
	lengthInBits += 32

	// Simple field (endpoint)
	lengthInBits += m.Endpoint.GetLengthInBits(ctx)

	// Simple field (senderCertificate)
	lengthInBits += m.SenderCertificate.GetLengthInBits(ctx)

	// Simple field (receiverCertificateThumbprint)
	lengthInBits += m.ReceiverCertificateThumbprint.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_OpenChannelMessageRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_OpenChannelMessageRequest) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_OpenChannelMessage, response bool) (__openChannelMessageRequest OpenChannelMessageRequest, err error) {
	m.OpenChannelMessageContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("OpenChannelMessageRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for OpenChannelMessageRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	secureChannelId, err := ReadSimpleField(ctx, "secureChannelId", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'secureChannelId' field"))
	}
	m.SecureChannelId = secureChannelId

	endpoint, err := ReadSimpleField[PascalString](ctx, "endpoint", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'endpoint' field"))
	}
	m.Endpoint = endpoint

	senderCertificate, err := ReadSimpleField[PascalByteString](ctx, "senderCertificate", ReadComplex[PascalByteString](PascalByteStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'senderCertificate' field"))
	}
	m.SenderCertificate = senderCertificate

	receiverCertificateThumbprint, err := ReadSimpleField[PascalByteString](ctx, "receiverCertificateThumbprint", ReadComplex[PascalByteString](PascalByteStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'receiverCertificateThumbprint' field"))
	}
	m.ReceiverCertificateThumbprint = receiverCertificateThumbprint

	if closeErr := readBuffer.CloseContext("OpenChannelMessageRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for OpenChannelMessageRequest")
	}

	return m, nil
}

func (m *_OpenChannelMessageRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_OpenChannelMessageRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("OpenChannelMessageRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for OpenChannelMessageRequest")
		}

		if err := WriteSimpleField[int32](ctx, "secureChannelId", m.GetSecureChannelId(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'secureChannelId' field")
		}

		if err := WriteSimpleField[PascalString](ctx, "endpoint", m.GetEndpoint(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'endpoint' field")
		}

		if err := WriteSimpleField[PascalByteString](ctx, "senderCertificate", m.GetSenderCertificate(), WriteComplex[PascalByteString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'senderCertificate' field")
		}

		if err := WriteSimpleField[PascalByteString](ctx, "receiverCertificateThumbprint", m.GetReceiverCertificateThumbprint(), WriteComplex[PascalByteString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'receiverCertificateThumbprint' field")
		}

		if popErr := writeBuffer.PopContext("OpenChannelMessageRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for OpenChannelMessageRequest")
		}
		return nil
	}
	return m.OpenChannelMessageContract.(*_OpenChannelMessage).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_OpenChannelMessageRequest) IsOpenChannelMessageRequest() {}

func (m *_OpenChannelMessageRequest) DeepCopy() any {
	return m.deepCopy()
}

func (m *_OpenChannelMessageRequest) deepCopy() *_OpenChannelMessageRequest {
	if m == nil {
		return nil
	}
	_OpenChannelMessageRequestCopy := &_OpenChannelMessageRequest{
		m.OpenChannelMessageContract.(*_OpenChannelMessage).deepCopy(),
		m.SecureChannelId,
		m.Endpoint.DeepCopy().(PascalString),
		m.SenderCertificate.DeepCopy().(PascalByteString),
		m.ReceiverCertificateThumbprint.DeepCopy().(PascalByteString),
	}
	m.OpenChannelMessageContract.(*_OpenChannelMessage)._SubType = m
	return _OpenChannelMessageRequestCopy
}

func (m *_OpenChannelMessageRequest) String() string {
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
