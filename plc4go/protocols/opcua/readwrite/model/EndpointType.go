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

// EndpointType is the corresponding interface of EndpointType
type EndpointType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetEndpointUrl returns EndpointUrl (property field)
	GetEndpointUrl() PascalString
	// GetSecurityMode returns SecurityMode (property field)
	GetSecurityMode() MessageSecurityMode
	// GetSecurityPolicyUri returns SecurityPolicyUri (property field)
	GetSecurityPolicyUri() PascalString
	// GetTransportProfileUri returns TransportProfileUri (property field)
	GetTransportProfileUri() PascalString
	// IsEndpointType is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsEndpointType()
	// CreateBuilder creates a EndpointTypeBuilder
	CreateEndpointTypeBuilder() EndpointTypeBuilder
}

// _EndpointType is the data-structure of this message
type _EndpointType struct {
	ExtensionObjectDefinitionContract
	EndpointUrl         PascalString
	SecurityMode        MessageSecurityMode
	SecurityPolicyUri   PascalString
	TransportProfileUri PascalString
}

var _ EndpointType = (*_EndpointType)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_EndpointType)(nil)

// NewEndpointType factory function for _EndpointType
func NewEndpointType(endpointUrl PascalString, securityMode MessageSecurityMode, securityPolicyUri PascalString, transportProfileUri PascalString) *_EndpointType {
	if endpointUrl == nil {
		panic("endpointUrl of type PascalString for EndpointType must not be nil")
	}
	if securityPolicyUri == nil {
		panic("securityPolicyUri of type PascalString for EndpointType must not be nil")
	}
	if transportProfileUri == nil {
		panic("transportProfileUri of type PascalString for EndpointType must not be nil")
	}
	_result := &_EndpointType{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		EndpointUrl:                       endpointUrl,
		SecurityMode:                      securityMode,
		SecurityPolicyUri:                 securityPolicyUri,
		TransportProfileUri:               transportProfileUri,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// EndpointTypeBuilder is a builder for EndpointType
type EndpointTypeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(endpointUrl PascalString, securityMode MessageSecurityMode, securityPolicyUri PascalString, transportProfileUri PascalString) EndpointTypeBuilder
	// WithEndpointUrl adds EndpointUrl (property field)
	WithEndpointUrl(PascalString) EndpointTypeBuilder
	// WithEndpointUrlBuilder adds EndpointUrl (property field) which is build by the builder
	WithEndpointUrlBuilder(func(PascalStringBuilder) PascalStringBuilder) EndpointTypeBuilder
	// WithSecurityMode adds SecurityMode (property field)
	WithSecurityMode(MessageSecurityMode) EndpointTypeBuilder
	// WithSecurityPolicyUri adds SecurityPolicyUri (property field)
	WithSecurityPolicyUri(PascalString) EndpointTypeBuilder
	// WithSecurityPolicyUriBuilder adds SecurityPolicyUri (property field) which is build by the builder
	WithSecurityPolicyUriBuilder(func(PascalStringBuilder) PascalStringBuilder) EndpointTypeBuilder
	// WithTransportProfileUri adds TransportProfileUri (property field)
	WithTransportProfileUri(PascalString) EndpointTypeBuilder
	// WithTransportProfileUriBuilder adds TransportProfileUri (property field) which is build by the builder
	WithTransportProfileUriBuilder(func(PascalStringBuilder) PascalStringBuilder) EndpointTypeBuilder
	// Build builds the EndpointType or returns an error if something is wrong
	Build() (EndpointType, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() EndpointType
}

// NewEndpointTypeBuilder() creates a EndpointTypeBuilder
func NewEndpointTypeBuilder() EndpointTypeBuilder {
	return &_EndpointTypeBuilder{_EndpointType: new(_EndpointType)}
}

type _EndpointTypeBuilder struct {
	*_EndpointType

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (EndpointTypeBuilder) = (*_EndpointTypeBuilder)(nil)

func (b *_EndpointTypeBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_EndpointTypeBuilder) WithMandatoryFields(endpointUrl PascalString, securityMode MessageSecurityMode, securityPolicyUri PascalString, transportProfileUri PascalString) EndpointTypeBuilder {
	return b.WithEndpointUrl(endpointUrl).WithSecurityMode(securityMode).WithSecurityPolicyUri(securityPolicyUri).WithTransportProfileUri(transportProfileUri)
}

func (b *_EndpointTypeBuilder) WithEndpointUrl(endpointUrl PascalString) EndpointTypeBuilder {
	b.EndpointUrl = endpointUrl
	return b
}

func (b *_EndpointTypeBuilder) WithEndpointUrlBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) EndpointTypeBuilder {
	builder := builderSupplier(b.EndpointUrl.CreatePascalStringBuilder())
	var err error
	b.EndpointUrl, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalStringBuilder failed"))
	}
	return b
}

func (b *_EndpointTypeBuilder) WithSecurityMode(securityMode MessageSecurityMode) EndpointTypeBuilder {
	b.SecurityMode = securityMode
	return b
}

func (b *_EndpointTypeBuilder) WithSecurityPolicyUri(securityPolicyUri PascalString) EndpointTypeBuilder {
	b.SecurityPolicyUri = securityPolicyUri
	return b
}

func (b *_EndpointTypeBuilder) WithSecurityPolicyUriBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) EndpointTypeBuilder {
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

func (b *_EndpointTypeBuilder) WithTransportProfileUri(transportProfileUri PascalString) EndpointTypeBuilder {
	b.TransportProfileUri = transportProfileUri
	return b
}

func (b *_EndpointTypeBuilder) WithTransportProfileUriBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) EndpointTypeBuilder {
	builder := builderSupplier(b.TransportProfileUri.CreatePascalStringBuilder())
	var err error
	b.TransportProfileUri, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalStringBuilder failed"))
	}
	return b
}

func (b *_EndpointTypeBuilder) Build() (EndpointType, error) {
	if b.EndpointUrl == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'endpointUrl' not set"))
	}
	if b.SecurityPolicyUri == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'securityPolicyUri' not set"))
	}
	if b.TransportProfileUri == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'transportProfileUri' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._EndpointType.deepCopy(), nil
}

func (b *_EndpointTypeBuilder) MustBuild() EndpointType {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_EndpointTypeBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_EndpointTypeBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_EndpointTypeBuilder) DeepCopy() any {
	_copy := b.CreateEndpointTypeBuilder().(*_EndpointTypeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateEndpointTypeBuilder creates a EndpointTypeBuilder
func (b *_EndpointType) CreateEndpointTypeBuilder() EndpointTypeBuilder {
	if b == nil {
		return NewEndpointTypeBuilder()
	}
	return &_EndpointTypeBuilder{_EndpointType: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_EndpointType) GetExtensionId() int32 {
	return int32(15530)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_EndpointType) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_EndpointType) GetEndpointUrl() PascalString {
	return m.EndpointUrl
}

func (m *_EndpointType) GetSecurityMode() MessageSecurityMode {
	return m.SecurityMode
}

func (m *_EndpointType) GetSecurityPolicyUri() PascalString {
	return m.SecurityPolicyUri
}

func (m *_EndpointType) GetTransportProfileUri() PascalString {
	return m.TransportProfileUri
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastEndpointType(structType any) EndpointType {
	if casted, ok := structType.(EndpointType); ok {
		return casted
	}
	if casted, ok := structType.(*EndpointType); ok {
		return *casted
	}
	return nil
}

func (m *_EndpointType) GetTypeName() string {
	return "EndpointType"
}

func (m *_EndpointType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (endpointUrl)
	lengthInBits += m.EndpointUrl.GetLengthInBits(ctx)

	// Simple field (securityMode)
	lengthInBits += 32

	// Simple field (securityPolicyUri)
	lengthInBits += m.SecurityPolicyUri.GetLengthInBits(ctx)

	// Simple field (transportProfileUri)
	lengthInBits += m.TransportProfileUri.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_EndpointType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_EndpointType) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__endpointType EndpointType, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("EndpointType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for EndpointType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	endpointUrl, err := ReadSimpleField[PascalString](ctx, "endpointUrl", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'endpointUrl' field"))
	}
	m.EndpointUrl = endpointUrl

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

	transportProfileUri, err := ReadSimpleField[PascalString](ctx, "transportProfileUri", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'transportProfileUri' field"))
	}
	m.TransportProfileUri = transportProfileUri

	if closeErr := readBuffer.CloseContext("EndpointType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for EndpointType")
	}

	return m, nil
}

func (m *_EndpointType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_EndpointType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("EndpointType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for EndpointType")
		}

		if err := WriteSimpleField[PascalString](ctx, "endpointUrl", m.GetEndpointUrl(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'endpointUrl' field")
		}

		if err := WriteSimpleEnumField[MessageSecurityMode](ctx, "securityMode", "MessageSecurityMode", m.GetSecurityMode(), WriteEnum[MessageSecurityMode, uint32](MessageSecurityMode.GetValue, MessageSecurityMode.PLC4XEnumName, WriteUnsignedInt(writeBuffer, 32))); err != nil {
			return errors.Wrap(err, "Error serializing 'securityMode' field")
		}

		if err := WriteSimpleField[PascalString](ctx, "securityPolicyUri", m.GetSecurityPolicyUri(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'securityPolicyUri' field")
		}

		if err := WriteSimpleField[PascalString](ctx, "transportProfileUri", m.GetTransportProfileUri(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'transportProfileUri' field")
		}

		if popErr := writeBuffer.PopContext("EndpointType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for EndpointType")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_EndpointType) IsEndpointType() {}

func (m *_EndpointType) DeepCopy() any {
	return m.deepCopy()
}

func (m *_EndpointType) deepCopy() *_EndpointType {
	if m == nil {
		return nil
	}
	_EndpointTypeCopy := &_EndpointType{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.EndpointUrl.DeepCopy().(PascalString),
		m.SecurityMode,
		m.SecurityPolicyUri.DeepCopy().(PascalString),
		m.TransportProfileUri.DeepCopy().(PascalString),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _EndpointTypeCopy
}

func (m *_EndpointType) String() string {
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
