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

// ApduDataExtWriteRouterMemoryRequest is the corresponding interface of ApduDataExtWriteRouterMemoryRequest
type ApduDataExtWriteRouterMemoryRequest interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ApduDataExt
	// IsApduDataExtWriteRouterMemoryRequest is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsApduDataExtWriteRouterMemoryRequest()
	// CreateBuilder creates a ApduDataExtWriteRouterMemoryRequestBuilder
	CreateApduDataExtWriteRouterMemoryRequestBuilder() ApduDataExtWriteRouterMemoryRequestBuilder
}

// _ApduDataExtWriteRouterMemoryRequest is the data-structure of this message
type _ApduDataExtWriteRouterMemoryRequest struct {
	ApduDataExtContract
}

var _ ApduDataExtWriteRouterMemoryRequest = (*_ApduDataExtWriteRouterMemoryRequest)(nil)
var _ ApduDataExtRequirements = (*_ApduDataExtWriteRouterMemoryRequest)(nil)

// NewApduDataExtWriteRouterMemoryRequest factory function for _ApduDataExtWriteRouterMemoryRequest
func NewApduDataExtWriteRouterMemoryRequest(length uint8) *_ApduDataExtWriteRouterMemoryRequest {
	_result := &_ApduDataExtWriteRouterMemoryRequest{
		ApduDataExtContract: NewApduDataExt(length),
	}
	_result.ApduDataExtContract.(*_ApduDataExt)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ApduDataExtWriteRouterMemoryRequestBuilder is a builder for ApduDataExtWriteRouterMemoryRequest
type ApduDataExtWriteRouterMemoryRequestBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() ApduDataExtWriteRouterMemoryRequestBuilder
	// Build builds the ApduDataExtWriteRouterMemoryRequest or returns an error if something is wrong
	Build() (ApduDataExtWriteRouterMemoryRequest, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ApduDataExtWriteRouterMemoryRequest
}

// NewApduDataExtWriteRouterMemoryRequestBuilder() creates a ApduDataExtWriteRouterMemoryRequestBuilder
func NewApduDataExtWriteRouterMemoryRequestBuilder() ApduDataExtWriteRouterMemoryRequestBuilder {
	return &_ApduDataExtWriteRouterMemoryRequestBuilder{_ApduDataExtWriteRouterMemoryRequest: new(_ApduDataExtWriteRouterMemoryRequest)}
}

type _ApduDataExtWriteRouterMemoryRequestBuilder struct {
	*_ApduDataExtWriteRouterMemoryRequest

	parentBuilder *_ApduDataExtBuilder

	err *utils.MultiError
}

var _ (ApduDataExtWriteRouterMemoryRequestBuilder) = (*_ApduDataExtWriteRouterMemoryRequestBuilder)(nil)

func (b *_ApduDataExtWriteRouterMemoryRequestBuilder) setParent(contract ApduDataExtContract) {
	b.ApduDataExtContract = contract
}

func (b *_ApduDataExtWriteRouterMemoryRequestBuilder) WithMandatoryFields() ApduDataExtWriteRouterMemoryRequestBuilder {
	return b
}

func (b *_ApduDataExtWriteRouterMemoryRequestBuilder) Build() (ApduDataExtWriteRouterMemoryRequest, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ApduDataExtWriteRouterMemoryRequest.deepCopy(), nil
}

func (b *_ApduDataExtWriteRouterMemoryRequestBuilder) MustBuild() ApduDataExtWriteRouterMemoryRequest {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ApduDataExtWriteRouterMemoryRequestBuilder) Done() ApduDataExtBuilder {
	return b.parentBuilder
}

func (b *_ApduDataExtWriteRouterMemoryRequestBuilder) buildForApduDataExt() (ApduDataExt, error) {
	return b.Build()
}

func (b *_ApduDataExtWriteRouterMemoryRequestBuilder) DeepCopy() any {
	_copy := b.CreateApduDataExtWriteRouterMemoryRequestBuilder().(*_ApduDataExtWriteRouterMemoryRequestBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateApduDataExtWriteRouterMemoryRequestBuilder creates a ApduDataExtWriteRouterMemoryRequestBuilder
func (b *_ApduDataExtWriteRouterMemoryRequest) CreateApduDataExtWriteRouterMemoryRequestBuilder() ApduDataExtWriteRouterMemoryRequestBuilder {
	if b == nil {
		return NewApduDataExtWriteRouterMemoryRequestBuilder()
	}
	return &_ApduDataExtWriteRouterMemoryRequestBuilder{_ApduDataExtWriteRouterMemoryRequest: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ApduDataExtWriteRouterMemoryRequest) GetExtApciType() uint8 {
	return 0x0A
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ApduDataExtWriteRouterMemoryRequest) GetParent() ApduDataExtContract {
	return m.ApduDataExtContract
}

// Deprecated: use the interface for direct cast
func CastApduDataExtWriteRouterMemoryRequest(structType any) ApduDataExtWriteRouterMemoryRequest {
	if casted, ok := structType.(ApduDataExtWriteRouterMemoryRequest); ok {
		return casted
	}
	if casted, ok := structType.(*ApduDataExtWriteRouterMemoryRequest); ok {
		return *casted
	}
	return nil
}

func (m *_ApduDataExtWriteRouterMemoryRequest) GetTypeName() string {
	return "ApduDataExtWriteRouterMemoryRequest"
}

func (m *_ApduDataExtWriteRouterMemoryRequest) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ApduDataExtContract.(*_ApduDataExt).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_ApduDataExtWriteRouterMemoryRequest) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ApduDataExtWriteRouterMemoryRequest) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ApduDataExt, length uint8) (__apduDataExtWriteRouterMemoryRequest ApduDataExtWriteRouterMemoryRequest, err error) {
	m.ApduDataExtContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ApduDataExtWriteRouterMemoryRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ApduDataExtWriteRouterMemoryRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("ApduDataExtWriteRouterMemoryRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ApduDataExtWriteRouterMemoryRequest")
	}

	return m, nil
}

func (m *_ApduDataExtWriteRouterMemoryRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ApduDataExtWriteRouterMemoryRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ApduDataExtWriteRouterMemoryRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ApduDataExtWriteRouterMemoryRequest")
		}

		if popErr := writeBuffer.PopContext("ApduDataExtWriteRouterMemoryRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ApduDataExtWriteRouterMemoryRequest")
		}
		return nil
	}
	return m.ApduDataExtContract.(*_ApduDataExt).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ApduDataExtWriteRouterMemoryRequest) IsApduDataExtWriteRouterMemoryRequest() {}

func (m *_ApduDataExtWriteRouterMemoryRequest) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ApduDataExtWriteRouterMemoryRequest) deepCopy() *_ApduDataExtWriteRouterMemoryRequest {
	if m == nil {
		return nil
	}
	_ApduDataExtWriteRouterMemoryRequestCopy := &_ApduDataExtWriteRouterMemoryRequest{
		m.ApduDataExtContract.(*_ApduDataExt).deepCopy(),
	}
	m.ApduDataExtContract.(*_ApduDataExt)._SubType = m
	return _ApduDataExtWriteRouterMemoryRequestCopy
}

func (m *_ApduDataExtWriteRouterMemoryRequest) String() string {
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
