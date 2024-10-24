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

// Constant values.
const EipConnectionResponse_PROTOCOLVERSION uint16 = 0x01
const EipConnectionResponse_FLAGS uint16 = 0x00

// EipConnectionResponse is the corresponding interface of EipConnectionResponse
type EipConnectionResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	EipPacket
	// IsEipConnectionResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsEipConnectionResponse()
	// CreateBuilder creates a EipConnectionResponseBuilder
	CreateEipConnectionResponseBuilder() EipConnectionResponseBuilder
}

// _EipConnectionResponse is the data-structure of this message
type _EipConnectionResponse struct {
	EipPacketContract
}

var _ EipConnectionResponse = (*_EipConnectionResponse)(nil)
var _ EipPacketRequirements = (*_EipConnectionResponse)(nil)

// NewEipConnectionResponse factory function for _EipConnectionResponse
func NewEipConnectionResponse(sessionHandle uint32, status uint32, senderContext []byte, options uint32) *_EipConnectionResponse {
	_result := &_EipConnectionResponse{
		EipPacketContract: NewEipPacket(sessionHandle, status, senderContext, options),
	}
	_result.EipPacketContract.(*_EipPacket)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// EipConnectionResponseBuilder is a builder for EipConnectionResponse
type EipConnectionResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() EipConnectionResponseBuilder
	// Build builds the EipConnectionResponse or returns an error if something is wrong
	Build() (EipConnectionResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() EipConnectionResponse
}

// NewEipConnectionResponseBuilder() creates a EipConnectionResponseBuilder
func NewEipConnectionResponseBuilder() EipConnectionResponseBuilder {
	return &_EipConnectionResponseBuilder{_EipConnectionResponse: new(_EipConnectionResponse)}
}

type _EipConnectionResponseBuilder struct {
	*_EipConnectionResponse

	parentBuilder *_EipPacketBuilder

	err *utils.MultiError
}

var _ (EipConnectionResponseBuilder) = (*_EipConnectionResponseBuilder)(nil)

func (b *_EipConnectionResponseBuilder) setParent(contract EipPacketContract) {
	b.EipPacketContract = contract
}

func (b *_EipConnectionResponseBuilder) WithMandatoryFields() EipConnectionResponseBuilder {
	return b
}

func (b *_EipConnectionResponseBuilder) Build() (EipConnectionResponse, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._EipConnectionResponse.deepCopy(), nil
}

func (b *_EipConnectionResponseBuilder) MustBuild() EipConnectionResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_EipConnectionResponseBuilder) Done() EipPacketBuilder {
	return b.parentBuilder
}

func (b *_EipConnectionResponseBuilder) buildForEipPacket() (EipPacket, error) {
	return b.Build()
}

func (b *_EipConnectionResponseBuilder) DeepCopy() any {
	_copy := b.CreateEipConnectionResponseBuilder().(*_EipConnectionResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateEipConnectionResponseBuilder creates a EipConnectionResponseBuilder
func (b *_EipConnectionResponse) CreateEipConnectionResponseBuilder() EipConnectionResponseBuilder {
	if b == nil {
		return NewEipConnectionResponseBuilder()
	}
	return &_EipConnectionResponseBuilder{_EipConnectionResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_EipConnectionResponse) GetCommand() uint16 {
	return 0x0065
}

func (m *_EipConnectionResponse) GetResponse() bool {
	return bool(true)
}

func (m *_EipConnectionResponse) GetPacketLength() uint16 {
	return 0
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_EipConnectionResponse) GetParent() EipPacketContract {
	return m.EipPacketContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for const fields.
///////////////////////

func (m *_EipConnectionResponse) GetProtocolVersion() uint16 {
	return EipConnectionResponse_PROTOCOLVERSION
}

func (m *_EipConnectionResponse) GetFlags() uint16 {
	return EipConnectionResponse_FLAGS
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastEipConnectionResponse(structType any) EipConnectionResponse {
	if casted, ok := structType.(EipConnectionResponse); ok {
		return casted
	}
	if casted, ok := structType.(*EipConnectionResponse); ok {
		return *casted
	}
	return nil
}

func (m *_EipConnectionResponse) GetTypeName() string {
	return "EipConnectionResponse"
}

func (m *_EipConnectionResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.EipPacketContract.(*_EipPacket).GetLengthInBits(ctx))

	// Const Field (protocolVersion)
	lengthInBits += 16

	// Const Field (flags)
	lengthInBits += 16

	return lengthInBits
}

func (m *_EipConnectionResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_EipConnectionResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_EipPacket, response bool) (__eipConnectionResponse EipConnectionResponse, err error) {
	m.EipPacketContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("EipConnectionResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for EipConnectionResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	protocolVersion, err := ReadConstField[uint16](ctx, "protocolVersion", ReadUnsignedShort(readBuffer, uint8(16)), EipConnectionResponse_PROTOCOLVERSION)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'protocolVersion' field"))
	}
	_ = protocolVersion

	flags, err := ReadConstField[uint16](ctx, "flags", ReadUnsignedShort(readBuffer, uint8(16)), EipConnectionResponse_FLAGS)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'flags' field"))
	}
	_ = flags

	if closeErr := readBuffer.CloseContext("EipConnectionResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for EipConnectionResponse")
	}

	return m, nil
}

func (m *_EipConnectionResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_EipConnectionResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("EipConnectionResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for EipConnectionResponse")
		}

		if err := WriteConstField(ctx, "protocolVersion", EipConnectionResponse_PROTOCOLVERSION, WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'protocolVersion' field")
		}

		if err := WriteConstField(ctx, "flags", EipConnectionResponse_FLAGS, WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'flags' field")
		}

		if popErr := writeBuffer.PopContext("EipConnectionResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for EipConnectionResponse")
		}
		return nil
	}
	return m.EipPacketContract.(*_EipPacket).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_EipConnectionResponse) IsEipConnectionResponse() {}

func (m *_EipConnectionResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_EipConnectionResponse) deepCopy() *_EipConnectionResponse {
	if m == nil {
		return nil
	}
	_EipConnectionResponseCopy := &_EipConnectionResponse{
		m.EipPacketContract.(*_EipPacket).deepCopy(),
	}
	m.EipPacketContract.(*_EipPacket)._SubType = m
	return _EipConnectionResponseCopy
}

func (m *_EipConnectionResponse) String() string {
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
