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

// ApduDataExtIndividualAddressSerialNumberWrite is the corresponding interface of ApduDataExtIndividualAddressSerialNumberWrite
type ApduDataExtIndividualAddressSerialNumberWrite interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ApduDataExt
	// IsApduDataExtIndividualAddressSerialNumberWrite is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsApduDataExtIndividualAddressSerialNumberWrite()
	// CreateBuilder creates a ApduDataExtIndividualAddressSerialNumberWriteBuilder
	CreateApduDataExtIndividualAddressSerialNumberWriteBuilder() ApduDataExtIndividualAddressSerialNumberWriteBuilder
}

// _ApduDataExtIndividualAddressSerialNumberWrite is the data-structure of this message
type _ApduDataExtIndividualAddressSerialNumberWrite struct {
	ApduDataExtContract
}

var _ ApduDataExtIndividualAddressSerialNumberWrite = (*_ApduDataExtIndividualAddressSerialNumberWrite)(nil)
var _ ApduDataExtRequirements = (*_ApduDataExtIndividualAddressSerialNumberWrite)(nil)

// NewApduDataExtIndividualAddressSerialNumberWrite factory function for _ApduDataExtIndividualAddressSerialNumberWrite
func NewApduDataExtIndividualAddressSerialNumberWrite(length uint8) *_ApduDataExtIndividualAddressSerialNumberWrite {
	_result := &_ApduDataExtIndividualAddressSerialNumberWrite{
		ApduDataExtContract: NewApduDataExt(length),
	}
	_result.ApduDataExtContract.(*_ApduDataExt)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ApduDataExtIndividualAddressSerialNumberWriteBuilder is a builder for ApduDataExtIndividualAddressSerialNumberWrite
type ApduDataExtIndividualAddressSerialNumberWriteBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() ApduDataExtIndividualAddressSerialNumberWriteBuilder
	// Build builds the ApduDataExtIndividualAddressSerialNumberWrite or returns an error if something is wrong
	Build() (ApduDataExtIndividualAddressSerialNumberWrite, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ApduDataExtIndividualAddressSerialNumberWrite
}

// NewApduDataExtIndividualAddressSerialNumberWriteBuilder() creates a ApduDataExtIndividualAddressSerialNumberWriteBuilder
func NewApduDataExtIndividualAddressSerialNumberWriteBuilder() ApduDataExtIndividualAddressSerialNumberWriteBuilder {
	return &_ApduDataExtIndividualAddressSerialNumberWriteBuilder{_ApduDataExtIndividualAddressSerialNumberWrite: new(_ApduDataExtIndividualAddressSerialNumberWrite)}
}

type _ApduDataExtIndividualAddressSerialNumberWriteBuilder struct {
	*_ApduDataExtIndividualAddressSerialNumberWrite

	parentBuilder *_ApduDataExtBuilder

	err *utils.MultiError
}

var _ (ApduDataExtIndividualAddressSerialNumberWriteBuilder) = (*_ApduDataExtIndividualAddressSerialNumberWriteBuilder)(nil)

func (b *_ApduDataExtIndividualAddressSerialNumberWriteBuilder) setParent(contract ApduDataExtContract) {
	b.ApduDataExtContract = contract
}

func (b *_ApduDataExtIndividualAddressSerialNumberWriteBuilder) WithMandatoryFields() ApduDataExtIndividualAddressSerialNumberWriteBuilder {
	return b
}

func (b *_ApduDataExtIndividualAddressSerialNumberWriteBuilder) Build() (ApduDataExtIndividualAddressSerialNumberWrite, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ApduDataExtIndividualAddressSerialNumberWrite.deepCopy(), nil
}

func (b *_ApduDataExtIndividualAddressSerialNumberWriteBuilder) MustBuild() ApduDataExtIndividualAddressSerialNumberWrite {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ApduDataExtIndividualAddressSerialNumberWriteBuilder) Done() ApduDataExtBuilder {
	return b.parentBuilder
}

func (b *_ApduDataExtIndividualAddressSerialNumberWriteBuilder) buildForApduDataExt() (ApduDataExt, error) {
	return b.Build()
}

func (b *_ApduDataExtIndividualAddressSerialNumberWriteBuilder) DeepCopy() any {
	_copy := b.CreateApduDataExtIndividualAddressSerialNumberWriteBuilder().(*_ApduDataExtIndividualAddressSerialNumberWriteBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateApduDataExtIndividualAddressSerialNumberWriteBuilder creates a ApduDataExtIndividualAddressSerialNumberWriteBuilder
func (b *_ApduDataExtIndividualAddressSerialNumberWrite) CreateApduDataExtIndividualAddressSerialNumberWriteBuilder() ApduDataExtIndividualAddressSerialNumberWriteBuilder {
	if b == nil {
		return NewApduDataExtIndividualAddressSerialNumberWriteBuilder()
	}
	return &_ApduDataExtIndividualAddressSerialNumberWriteBuilder{_ApduDataExtIndividualAddressSerialNumberWrite: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) GetExtApciType() uint8 {
	return 0x1E
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) GetParent() ApduDataExtContract {
	return m.ApduDataExtContract
}

// Deprecated: use the interface for direct cast
func CastApduDataExtIndividualAddressSerialNumberWrite(structType any) ApduDataExtIndividualAddressSerialNumberWrite {
	if casted, ok := structType.(ApduDataExtIndividualAddressSerialNumberWrite); ok {
		return casted
	}
	if casted, ok := structType.(*ApduDataExtIndividualAddressSerialNumberWrite); ok {
		return *casted
	}
	return nil
}

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) GetTypeName() string {
	return "ApduDataExtIndividualAddressSerialNumberWrite"
}

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ApduDataExtContract.(*_ApduDataExt).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ApduDataExt, length uint8) (__apduDataExtIndividualAddressSerialNumberWrite ApduDataExtIndividualAddressSerialNumberWrite, err error) {
	m.ApduDataExtContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ApduDataExtIndividualAddressSerialNumberWrite"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ApduDataExtIndividualAddressSerialNumberWrite")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("ApduDataExtIndividualAddressSerialNumberWrite"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ApduDataExtIndividualAddressSerialNumberWrite")
	}

	return m, nil
}

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ApduDataExtIndividualAddressSerialNumberWrite"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ApduDataExtIndividualAddressSerialNumberWrite")
		}

		if popErr := writeBuffer.PopContext("ApduDataExtIndividualAddressSerialNumberWrite"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ApduDataExtIndividualAddressSerialNumberWrite")
		}
		return nil
	}
	return m.ApduDataExtContract.(*_ApduDataExt).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) IsApduDataExtIndividualAddressSerialNumberWrite() {
}

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) deepCopy() *_ApduDataExtIndividualAddressSerialNumberWrite {
	if m == nil {
		return nil
	}
	_ApduDataExtIndividualAddressSerialNumberWriteCopy := &_ApduDataExtIndividualAddressSerialNumberWrite{
		m.ApduDataExtContract.(*_ApduDataExt).deepCopy(),
	}
	m.ApduDataExtContract.(*_ApduDataExt)._SubType = m
	return _ApduDataExtIndividualAddressSerialNumberWriteCopy
}

func (m *_ApduDataExtIndividualAddressSerialNumberWrite) String() string {
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
