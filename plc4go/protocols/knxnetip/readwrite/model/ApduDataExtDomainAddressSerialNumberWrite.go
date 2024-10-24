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

// ApduDataExtDomainAddressSerialNumberWrite is the corresponding interface of ApduDataExtDomainAddressSerialNumberWrite
type ApduDataExtDomainAddressSerialNumberWrite interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ApduDataExt
	// IsApduDataExtDomainAddressSerialNumberWrite is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsApduDataExtDomainAddressSerialNumberWrite()
	// CreateBuilder creates a ApduDataExtDomainAddressSerialNumberWriteBuilder
	CreateApduDataExtDomainAddressSerialNumberWriteBuilder() ApduDataExtDomainAddressSerialNumberWriteBuilder
}

// _ApduDataExtDomainAddressSerialNumberWrite is the data-structure of this message
type _ApduDataExtDomainAddressSerialNumberWrite struct {
	ApduDataExtContract
}

var _ ApduDataExtDomainAddressSerialNumberWrite = (*_ApduDataExtDomainAddressSerialNumberWrite)(nil)
var _ ApduDataExtRequirements = (*_ApduDataExtDomainAddressSerialNumberWrite)(nil)

// NewApduDataExtDomainAddressSerialNumberWrite factory function for _ApduDataExtDomainAddressSerialNumberWrite
func NewApduDataExtDomainAddressSerialNumberWrite(length uint8) *_ApduDataExtDomainAddressSerialNumberWrite {
	_result := &_ApduDataExtDomainAddressSerialNumberWrite{
		ApduDataExtContract: NewApduDataExt(length),
	}
	_result.ApduDataExtContract.(*_ApduDataExt)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ApduDataExtDomainAddressSerialNumberWriteBuilder is a builder for ApduDataExtDomainAddressSerialNumberWrite
type ApduDataExtDomainAddressSerialNumberWriteBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() ApduDataExtDomainAddressSerialNumberWriteBuilder
	// Build builds the ApduDataExtDomainAddressSerialNumberWrite or returns an error if something is wrong
	Build() (ApduDataExtDomainAddressSerialNumberWrite, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ApduDataExtDomainAddressSerialNumberWrite
}

// NewApduDataExtDomainAddressSerialNumberWriteBuilder() creates a ApduDataExtDomainAddressSerialNumberWriteBuilder
func NewApduDataExtDomainAddressSerialNumberWriteBuilder() ApduDataExtDomainAddressSerialNumberWriteBuilder {
	return &_ApduDataExtDomainAddressSerialNumberWriteBuilder{_ApduDataExtDomainAddressSerialNumberWrite: new(_ApduDataExtDomainAddressSerialNumberWrite)}
}

type _ApduDataExtDomainAddressSerialNumberWriteBuilder struct {
	*_ApduDataExtDomainAddressSerialNumberWrite

	parentBuilder *_ApduDataExtBuilder

	err *utils.MultiError
}

var _ (ApduDataExtDomainAddressSerialNumberWriteBuilder) = (*_ApduDataExtDomainAddressSerialNumberWriteBuilder)(nil)

func (b *_ApduDataExtDomainAddressSerialNumberWriteBuilder) setParent(contract ApduDataExtContract) {
	b.ApduDataExtContract = contract
}

func (b *_ApduDataExtDomainAddressSerialNumberWriteBuilder) WithMandatoryFields() ApduDataExtDomainAddressSerialNumberWriteBuilder {
	return b
}

func (b *_ApduDataExtDomainAddressSerialNumberWriteBuilder) Build() (ApduDataExtDomainAddressSerialNumberWrite, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ApduDataExtDomainAddressSerialNumberWrite.deepCopy(), nil
}

func (b *_ApduDataExtDomainAddressSerialNumberWriteBuilder) MustBuild() ApduDataExtDomainAddressSerialNumberWrite {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ApduDataExtDomainAddressSerialNumberWriteBuilder) Done() ApduDataExtBuilder {
	return b.parentBuilder
}

func (b *_ApduDataExtDomainAddressSerialNumberWriteBuilder) buildForApduDataExt() (ApduDataExt, error) {
	return b.Build()
}

func (b *_ApduDataExtDomainAddressSerialNumberWriteBuilder) DeepCopy() any {
	_copy := b.CreateApduDataExtDomainAddressSerialNumberWriteBuilder().(*_ApduDataExtDomainAddressSerialNumberWriteBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateApduDataExtDomainAddressSerialNumberWriteBuilder creates a ApduDataExtDomainAddressSerialNumberWriteBuilder
func (b *_ApduDataExtDomainAddressSerialNumberWrite) CreateApduDataExtDomainAddressSerialNumberWriteBuilder() ApduDataExtDomainAddressSerialNumberWriteBuilder {
	if b == nil {
		return NewApduDataExtDomainAddressSerialNumberWriteBuilder()
	}
	return &_ApduDataExtDomainAddressSerialNumberWriteBuilder{_ApduDataExtDomainAddressSerialNumberWrite: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ApduDataExtDomainAddressSerialNumberWrite) GetExtApciType() uint8 {
	return 0x2E
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ApduDataExtDomainAddressSerialNumberWrite) GetParent() ApduDataExtContract {
	return m.ApduDataExtContract
}

// Deprecated: use the interface for direct cast
func CastApduDataExtDomainAddressSerialNumberWrite(structType any) ApduDataExtDomainAddressSerialNumberWrite {
	if casted, ok := structType.(ApduDataExtDomainAddressSerialNumberWrite); ok {
		return casted
	}
	if casted, ok := structType.(*ApduDataExtDomainAddressSerialNumberWrite); ok {
		return *casted
	}
	return nil
}

func (m *_ApduDataExtDomainAddressSerialNumberWrite) GetTypeName() string {
	return "ApduDataExtDomainAddressSerialNumberWrite"
}

func (m *_ApduDataExtDomainAddressSerialNumberWrite) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ApduDataExtContract.(*_ApduDataExt).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_ApduDataExtDomainAddressSerialNumberWrite) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ApduDataExtDomainAddressSerialNumberWrite) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ApduDataExt, length uint8) (__apduDataExtDomainAddressSerialNumberWrite ApduDataExtDomainAddressSerialNumberWrite, err error) {
	m.ApduDataExtContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ApduDataExtDomainAddressSerialNumberWrite"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ApduDataExtDomainAddressSerialNumberWrite")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("ApduDataExtDomainAddressSerialNumberWrite"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ApduDataExtDomainAddressSerialNumberWrite")
	}

	return m, nil
}

func (m *_ApduDataExtDomainAddressSerialNumberWrite) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ApduDataExtDomainAddressSerialNumberWrite) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ApduDataExtDomainAddressSerialNumberWrite"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ApduDataExtDomainAddressSerialNumberWrite")
		}

		if popErr := writeBuffer.PopContext("ApduDataExtDomainAddressSerialNumberWrite"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ApduDataExtDomainAddressSerialNumberWrite")
		}
		return nil
	}
	return m.ApduDataExtContract.(*_ApduDataExt).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ApduDataExtDomainAddressSerialNumberWrite) IsApduDataExtDomainAddressSerialNumberWrite() {}

func (m *_ApduDataExtDomainAddressSerialNumberWrite) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ApduDataExtDomainAddressSerialNumberWrite) deepCopy() *_ApduDataExtDomainAddressSerialNumberWrite {
	if m == nil {
		return nil
	}
	_ApduDataExtDomainAddressSerialNumberWriteCopy := &_ApduDataExtDomainAddressSerialNumberWrite{
		m.ApduDataExtContract.(*_ApduDataExt).deepCopy(),
	}
	m.ApduDataExtContract.(*_ApduDataExt)._SubType = m
	return _ApduDataExtDomainAddressSerialNumberWriteCopy
}

func (m *_ApduDataExtDomainAddressSerialNumberWrite) String() string {
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
