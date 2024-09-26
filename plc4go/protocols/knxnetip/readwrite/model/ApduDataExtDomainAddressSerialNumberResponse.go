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

// ApduDataExtDomainAddressSerialNumberResponse is the corresponding interface of ApduDataExtDomainAddressSerialNumberResponse
type ApduDataExtDomainAddressSerialNumberResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ApduDataExt
	// IsApduDataExtDomainAddressSerialNumberResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsApduDataExtDomainAddressSerialNumberResponse()
	// CreateBuilder creates a ApduDataExtDomainAddressSerialNumberResponseBuilder
	CreateApduDataExtDomainAddressSerialNumberResponseBuilder() ApduDataExtDomainAddressSerialNumberResponseBuilder
}

// _ApduDataExtDomainAddressSerialNumberResponse is the data-structure of this message
type _ApduDataExtDomainAddressSerialNumberResponse struct {
	ApduDataExtContract
}

var _ ApduDataExtDomainAddressSerialNumberResponse = (*_ApduDataExtDomainAddressSerialNumberResponse)(nil)
var _ ApduDataExtRequirements = (*_ApduDataExtDomainAddressSerialNumberResponse)(nil)

// NewApduDataExtDomainAddressSerialNumberResponse factory function for _ApduDataExtDomainAddressSerialNumberResponse
func NewApduDataExtDomainAddressSerialNumberResponse(length uint8) *_ApduDataExtDomainAddressSerialNumberResponse {
	_result := &_ApduDataExtDomainAddressSerialNumberResponse{
		ApduDataExtContract: NewApduDataExt(length),
	}
	_result.ApduDataExtContract.(*_ApduDataExt)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ApduDataExtDomainAddressSerialNumberResponseBuilder is a builder for ApduDataExtDomainAddressSerialNumberResponse
type ApduDataExtDomainAddressSerialNumberResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() ApduDataExtDomainAddressSerialNumberResponseBuilder
	// Build builds the ApduDataExtDomainAddressSerialNumberResponse or returns an error if something is wrong
	Build() (ApduDataExtDomainAddressSerialNumberResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ApduDataExtDomainAddressSerialNumberResponse
}

// NewApduDataExtDomainAddressSerialNumberResponseBuilder() creates a ApduDataExtDomainAddressSerialNumberResponseBuilder
func NewApduDataExtDomainAddressSerialNumberResponseBuilder() ApduDataExtDomainAddressSerialNumberResponseBuilder {
	return &_ApduDataExtDomainAddressSerialNumberResponseBuilder{_ApduDataExtDomainAddressSerialNumberResponse: new(_ApduDataExtDomainAddressSerialNumberResponse)}
}

type _ApduDataExtDomainAddressSerialNumberResponseBuilder struct {
	*_ApduDataExtDomainAddressSerialNumberResponse

	parentBuilder *_ApduDataExtBuilder

	err *utils.MultiError
}

var _ (ApduDataExtDomainAddressSerialNumberResponseBuilder) = (*_ApduDataExtDomainAddressSerialNumberResponseBuilder)(nil)

func (b *_ApduDataExtDomainAddressSerialNumberResponseBuilder) setParent(contract ApduDataExtContract) {
	b.ApduDataExtContract = contract
}

func (b *_ApduDataExtDomainAddressSerialNumberResponseBuilder) WithMandatoryFields() ApduDataExtDomainAddressSerialNumberResponseBuilder {
	return b
}

func (b *_ApduDataExtDomainAddressSerialNumberResponseBuilder) Build() (ApduDataExtDomainAddressSerialNumberResponse, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ApduDataExtDomainAddressSerialNumberResponse.deepCopy(), nil
}

func (b *_ApduDataExtDomainAddressSerialNumberResponseBuilder) MustBuild() ApduDataExtDomainAddressSerialNumberResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ApduDataExtDomainAddressSerialNumberResponseBuilder) Done() ApduDataExtBuilder {
	return b.parentBuilder
}

func (b *_ApduDataExtDomainAddressSerialNumberResponseBuilder) buildForApduDataExt() (ApduDataExt, error) {
	return b.Build()
}

func (b *_ApduDataExtDomainAddressSerialNumberResponseBuilder) DeepCopy() any {
	_copy := b.CreateApduDataExtDomainAddressSerialNumberResponseBuilder().(*_ApduDataExtDomainAddressSerialNumberResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateApduDataExtDomainAddressSerialNumberResponseBuilder creates a ApduDataExtDomainAddressSerialNumberResponseBuilder
func (b *_ApduDataExtDomainAddressSerialNumberResponse) CreateApduDataExtDomainAddressSerialNumberResponseBuilder() ApduDataExtDomainAddressSerialNumberResponseBuilder {
	if b == nil {
		return NewApduDataExtDomainAddressSerialNumberResponseBuilder()
	}
	return &_ApduDataExtDomainAddressSerialNumberResponseBuilder{_ApduDataExtDomainAddressSerialNumberResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ApduDataExtDomainAddressSerialNumberResponse) GetExtApciType() uint8 {
	return 0x2D
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ApduDataExtDomainAddressSerialNumberResponse) GetParent() ApduDataExtContract {
	return m.ApduDataExtContract
}

// Deprecated: use the interface for direct cast
func CastApduDataExtDomainAddressSerialNumberResponse(structType any) ApduDataExtDomainAddressSerialNumberResponse {
	if casted, ok := structType.(ApduDataExtDomainAddressSerialNumberResponse); ok {
		return casted
	}
	if casted, ok := structType.(*ApduDataExtDomainAddressSerialNumberResponse); ok {
		return *casted
	}
	return nil
}

func (m *_ApduDataExtDomainAddressSerialNumberResponse) GetTypeName() string {
	return "ApduDataExtDomainAddressSerialNumberResponse"
}

func (m *_ApduDataExtDomainAddressSerialNumberResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ApduDataExtContract.(*_ApduDataExt).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_ApduDataExtDomainAddressSerialNumberResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ApduDataExtDomainAddressSerialNumberResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ApduDataExt, length uint8) (__apduDataExtDomainAddressSerialNumberResponse ApduDataExtDomainAddressSerialNumberResponse, err error) {
	m.ApduDataExtContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ApduDataExtDomainAddressSerialNumberResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ApduDataExtDomainAddressSerialNumberResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("ApduDataExtDomainAddressSerialNumberResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ApduDataExtDomainAddressSerialNumberResponse")
	}

	return m, nil
}

func (m *_ApduDataExtDomainAddressSerialNumberResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ApduDataExtDomainAddressSerialNumberResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ApduDataExtDomainAddressSerialNumberResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ApduDataExtDomainAddressSerialNumberResponse")
		}

		if popErr := writeBuffer.PopContext("ApduDataExtDomainAddressSerialNumberResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ApduDataExtDomainAddressSerialNumberResponse")
		}
		return nil
	}
	return m.ApduDataExtContract.(*_ApduDataExt).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ApduDataExtDomainAddressSerialNumberResponse) IsApduDataExtDomainAddressSerialNumberResponse() {
}

func (m *_ApduDataExtDomainAddressSerialNumberResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ApduDataExtDomainAddressSerialNumberResponse) deepCopy() *_ApduDataExtDomainAddressSerialNumberResponse {
	if m == nil {
		return nil
	}
	_ApduDataExtDomainAddressSerialNumberResponseCopy := &_ApduDataExtDomainAddressSerialNumberResponse{
		m.ApduDataExtContract.(*_ApduDataExt).deepCopy(),
	}
	m.ApduDataExtContract.(*_ApduDataExt)._SubType = m
	return _ApduDataExtDomainAddressSerialNumberResponseCopy
}

func (m *_ApduDataExtDomainAddressSerialNumberResponse) String() string {
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
