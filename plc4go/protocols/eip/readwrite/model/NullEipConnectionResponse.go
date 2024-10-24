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

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// NullEipConnectionResponse is the corresponding interface of NullEipConnectionResponse
type NullEipConnectionResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	EipPacket
	// IsNullEipConnectionResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsNullEipConnectionResponse()
	// CreateBuilder creates a NullEipConnectionResponseBuilder
	CreateNullEipConnectionResponseBuilder() NullEipConnectionResponseBuilder
}

// _NullEipConnectionResponse is the data-structure of this message
type _NullEipConnectionResponse struct {
	EipPacketContract
}

var _ NullEipConnectionResponse = (*_NullEipConnectionResponse)(nil)
var _ EipPacketRequirements = (*_NullEipConnectionResponse)(nil)

// NewNullEipConnectionResponse factory function for _NullEipConnectionResponse
func NewNullEipConnectionResponse(sessionHandle uint32, status uint32, senderContext []byte, options uint32) *_NullEipConnectionResponse {
	_result := &_NullEipConnectionResponse{
		EipPacketContract: NewEipPacket(sessionHandle, status, senderContext, options),
	}
	_result.EipPacketContract.(*_EipPacket)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// NullEipConnectionResponseBuilder is a builder for NullEipConnectionResponse
type NullEipConnectionResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() NullEipConnectionResponseBuilder
	// Build builds the NullEipConnectionResponse or returns an error if something is wrong
	Build() (NullEipConnectionResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() NullEipConnectionResponse
}

// NewNullEipConnectionResponseBuilder() creates a NullEipConnectionResponseBuilder
func NewNullEipConnectionResponseBuilder() NullEipConnectionResponseBuilder {
	return &_NullEipConnectionResponseBuilder{_NullEipConnectionResponse: new(_NullEipConnectionResponse)}
}

type _NullEipConnectionResponseBuilder struct {
	*_NullEipConnectionResponse

	parentBuilder *_EipPacketBuilder

	err *utils.MultiError
}

var _ (NullEipConnectionResponseBuilder) = (*_NullEipConnectionResponseBuilder)(nil)

func (b *_NullEipConnectionResponseBuilder) setParent(contract EipPacketContract) {
	b.EipPacketContract = contract
}

func (b *_NullEipConnectionResponseBuilder) WithMandatoryFields() NullEipConnectionResponseBuilder {
	return b
}

func (b *_NullEipConnectionResponseBuilder) Build() (NullEipConnectionResponse, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._NullEipConnectionResponse.deepCopy(), nil
}

func (b *_NullEipConnectionResponseBuilder) MustBuild() NullEipConnectionResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_NullEipConnectionResponseBuilder) Done() EipPacketBuilder {
	return b.parentBuilder
}

func (b *_NullEipConnectionResponseBuilder) buildForEipPacket() (EipPacket, error) {
	return b.Build()
}

func (b *_NullEipConnectionResponseBuilder) DeepCopy() any {
	_copy := b.CreateNullEipConnectionResponseBuilder().(*_NullEipConnectionResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateNullEipConnectionResponseBuilder creates a NullEipConnectionResponseBuilder
func (b *_NullEipConnectionResponse) CreateNullEipConnectionResponseBuilder() NullEipConnectionResponseBuilder {
	if b == nil {
		return NewNullEipConnectionResponseBuilder()
	}
	return &_NullEipConnectionResponseBuilder{_NullEipConnectionResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_NullEipConnectionResponse) GetCommand() uint16 {
	return 0x0065
}

func (m *_NullEipConnectionResponse) GetResponse() bool {
	return bool(true)
}

func (m *_NullEipConnectionResponse) GetPacketLength() uint16 {
	return uint16(0)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_NullEipConnectionResponse) GetParent() EipPacketContract {
	return m.EipPacketContract
}

// Deprecated: use the interface for direct cast
func CastNullEipConnectionResponse(structType any) NullEipConnectionResponse {
	if casted, ok := structType.(NullEipConnectionResponse); ok {
		return casted
	}
	if casted, ok := structType.(*NullEipConnectionResponse); ok {
		return *casted
	}
	return nil
}

func (m *_NullEipConnectionResponse) GetTypeName() string {
	return "NullEipConnectionResponse"
}

func (m *_NullEipConnectionResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.EipPacketContract.(*_EipPacket).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_NullEipConnectionResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_NullEipConnectionResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_EipPacket, response bool) (__nullEipConnectionResponse NullEipConnectionResponse, err error) {
	m.EipPacketContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("NullEipConnectionResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for NullEipConnectionResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("NullEipConnectionResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for NullEipConnectionResponse")
	}

	return m, nil
}

func (m *_NullEipConnectionResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_NullEipConnectionResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("NullEipConnectionResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for NullEipConnectionResponse")
		}

		if popErr := writeBuffer.PopContext("NullEipConnectionResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for NullEipConnectionResponse")
		}
		return nil
	}
	return m.EipPacketContract.(*_EipPacket).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_NullEipConnectionResponse) IsNullEipConnectionResponse() {}

func (m *_NullEipConnectionResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_NullEipConnectionResponse) deepCopy() *_NullEipConnectionResponse {
	if m == nil {
		return nil
	}
	_NullEipConnectionResponseCopy := &_NullEipConnectionResponse{
		m.EipPacketContract.(*_EipPacket).deepCopy(),
	}
	m.EipPacketContract.(*_EipPacket)._SubType = m
	return _NullEipConnectionResponseCopy
}

func (m *_NullEipConnectionResponse) String() string {
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
