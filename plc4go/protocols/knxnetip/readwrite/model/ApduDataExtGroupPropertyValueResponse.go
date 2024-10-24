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

// ApduDataExtGroupPropertyValueResponse is the corresponding interface of ApduDataExtGroupPropertyValueResponse
type ApduDataExtGroupPropertyValueResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ApduDataExt
	// IsApduDataExtGroupPropertyValueResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsApduDataExtGroupPropertyValueResponse()
	// CreateBuilder creates a ApduDataExtGroupPropertyValueResponseBuilder
	CreateApduDataExtGroupPropertyValueResponseBuilder() ApduDataExtGroupPropertyValueResponseBuilder
}

// _ApduDataExtGroupPropertyValueResponse is the data-structure of this message
type _ApduDataExtGroupPropertyValueResponse struct {
	ApduDataExtContract
}

var _ ApduDataExtGroupPropertyValueResponse = (*_ApduDataExtGroupPropertyValueResponse)(nil)
var _ ApduDataExtRequirements = (*_ApduDataExtGroupPropertyValueResponse)(nil)

// NewApduDataExtGroupPropertyValueResponse factory function for _ApduDataExtGroupPropertyValueResponse
func NewApduDataExtGroupPropertyValueResponse(length uint8) *_ApduDataExtGroupPropertyValueResponse {
	_result := &_ApduDataExtGroupPropertyValueResponse{
		ApduDataExtContract: NewApduDataExt(length),
	}
	_result.ApduDataExtContract.(*_ApduDataExt)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ApduDataExtGroupPropertyValueResponseBuilder is a builder for ApduDataExtGroupPropertyValueResponse
type ApduDataExtGroupPropertyValueResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() ApduDataExtGroupPropertyValueResponseBuilder
	// Build builds the ApduDataExtGroupPropertyValueResponse or returns an error if something is wrong
	Build() (ApduDataExtGroupPropertyValueResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ApduDataExtGroupPropertyValueResponse
}

// NewApduDataExtGroupPropertyValueResponseBuilder() creates a ApduDataExtGroupPropertyValueResponseBuilder
func NewApduDataExtGroupPropertyValueResponseBuilder() ApduDataExtGroupPropertyValueResponseBuilder {
	return &_ApduDataExtGroupPropertyValueResponseBuilder{_ApduDataExtGroupPropertyValueResponse: new(_ApduDataExtGroupPropertyValueResponse)}
}

type _ApduDataExtGroupPropertyValueResponseBuilder struct {
	*_ApduDataExtGroupPropertyValueResponse

	parentBuilder *_ApduDataExtBuilder

	err *utils.MultiError
}

var _ (ApduDataExtGroupPropertyValueResponseBuilder) = (*_ApduDataExtGroupPropertyValueResponseBuilder)(nil)

func (b *_ApduDataExtGroupPropertyValueResponseBuilder) setParent(contract ApduDataExtContract) {
	b.ApduDataExtContract = contract
}

func (b *_ApduDataExtGroupPropertyValueResponseBuilder) WithMandatoryFields() ApduDataExtGroupPropertyValueResponseBuilder {
	return b
}

func (b *_ApduDataExtGroupPropertyValueResponseBuilder) Build() (ApduDataExtGroupPropertyValueResponse, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ApduDataExtGroupPropertyValueResponse.deepCopy(), nil
}

func (b *_ApduDataExtGroupPropertyValueResponseBuilder) MustBuild() ApduDataExtGroupPropertyValueResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ApduDataExtGroupPropertyValueResponseBuilder) Done() ApduDataExtBuilder {
	return b.parentBuilder
}

func (b *_ApduDataExtGroupPropertyValueResponseBuilder) buildForApduDataExt() (ApduDataExt, error) {
	return b.Build()
}

func (b *_ApduDataExtGroupPropertyValueResponseBuilder) DeepCopy() any {
	_copy := b.CreateApduDataExtGroupPropertyValueResponseBuilder().(*_ApduDataExtGroupPropertyValueResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateApduDataExtGroupPropertyValueResponseBuilder creates a ApduDataExtGroupPropertyValueResponseBuilder
func (b *_ApduDataExtGroupPropertyValueResponse) CreateApduDataExtGroupPropertyValueResponseBuilder() ApduDataExtGroupPropertyValueResponseBuilder {
	if b == nil {
		return NewApduDataExtGroupPropertyValueResponseBuilder()
	}
	return &_ApduDataExtGroupPropertyValueResponseBuilder{_ApduDataExtGroupPropertyValueResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ApduDataExtGroupPropertyValueResponse) GetExtApciType() uint8 {
	return 0x29
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ApduDataExtGroupPropertyValueResponse) GetParent() ApduDataExtContract {
	return m.ApduDataExtContract
}

// Deprecated: use the interface for direct cast
func CastApduDataExtGroupPropertyValueResponse(structType any) ApduDataExtGroupPropertyValueResponse {
	if casted, ok := structType.(ApduDataExtGroupPropertyValueResponse); ok {
		return casted
	}
	if casted, ok := structType.(*ApduDataExtGroupPropertyValueResponse); ok {
		return *casted
	}
	return nil
}

func (m *_ApduDataExtGroupPropertyValueResponse) GetTypeName() string {
	return "ApduDataExtGroupPropertyValueResponse"
}

func (m *_ApduDataExtGroupPropertyValueResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ApduDataExtContract.(*_ApduDataExt).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_ApduDataExtGroupPropertyValueResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ApduDataExtGroupPropertyValueResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ApduDataExt, length uint8) (__apduDataExtGroupPropertyValueResponse ApduDataExtGroupPropertyValueResponse, err error) {
	m.ApduDataExtContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ApduDataExtGroupPropertyValueResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ApduDataExtGroupPropertyValueResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("ApduDataExtGroupPropertyValueResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ApduDataExtGroupPropertyValueResponse")
	}

	return m, nil
}

func (m *_ApduDataExtGroupPropertyValueResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ApduDataExtGroupPropertyValueResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ApduDataExtGroupPropertyValueResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ApduDataExtGroupPropertyValueResponse")
		}

		if popErr := writeBuffer.PopContext("ApduDataExtGroupPropertyValueResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ApduDataExtGroupPropertyValueResponse")
		}
		return nil
	}
	return m.ApduDataExtContract.(*_ApduDataExt).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ApduDataExtGroupPropertyValueResponse) IsApduDataExtGroupPropertyValueResponse() {}

func (m *_ApduDataExtGroupPropertyValueResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ApduDataExtGroupPropertyValueResponse) deepCopy() *_ApduDataExtGroupPropertyValueResponse {
	if m == nil {
		return nil
	}
	_ApduDataExtGroupPropertyValueResponseCopy := &_ApduDataExtGroupPropertyValueResponse{
		m.ApduDataExtContract.(*_ApduDataExt).deepCopy(),
	}
	m.ApduDataExtContract.(*_ApduDataExt)._SubType = m
	return _ApduDataExtGroupPropertyValueResponseCopy
}

func (m *_ApduDataExtGroupPropertyValueResponse) String() string {
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
