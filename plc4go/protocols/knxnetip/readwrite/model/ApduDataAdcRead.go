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

// ApduDataAdcRead is the corresponding interface of ApduDataAdcRead
type ApduDataAdcRead interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ApduData
	// IsApduDataAdcRead is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsApduDataAdcRead()
	// CreateBuilder creates a ApduDataAdcReadBuilder
	CreateApduDataAdcReadBuilder() ApduDataAdcReadBuilder
}

// _ApduDataAdcRead is the data-structure of this message
type _ApduDataAdcRead struct {
	ApduDataContract
}

var _ ApduDataAdcRead = (*_ApduDataAdcRead)(nil)
var _ ApduDataRequirements = (*_ApduDataAdcRead)(nil)

// NewApduDataAdcRead factory function for _ApduDataAdcRead
func NewApduDataAdcRead(dataLength uint8) *_ApduDataAdcRead {
	_result := &_ApduDataAdcRead{
		ApduDataContract: NewApduData(dataLength),
	}
	_result.ApduDataContract.(*_ApduData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ApduDataAdcReadBuilder is a builder for ApduDataAdcRead
type ApduDataAdcReadBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields() ApduDataAdcReadBuilder
	// Build builds the ApduDataAdcRead or returns an error if something is wrong
	Build() (ApduDataAdcRead, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ApduDataAdcRead
}

// NewApduDataAdcReadBuilder() creates a ApduDataAdcReadBuilder
func NewApduDataAdcReadBuilder() ApduDataAdcReadBuilder {
	return &_ApduDataAdcReadBuilder{_ApduDataAdcRead: new(_ApduDataAdcRead)}
}

type _ApduDataAdcReadBuilder struct {
	*_ApduDataAdcRead

	parentBuilder *_ApduDataBuilder

	err *utils.MultiError
}

var _ (ApduDataAdcReadBuilder) = (*_ApduDataAdcReadBuilder)(nil)

func (b *_ApduDataAdcReadBuilder) setParent(contract ApduDataContract) {
	b.ApduDataContract = contract
}

func (b *_ApduDataAdcReadBuilder) WithMandatoryFields() ApduDataAdcReadBuilder {
	return b
}

func (b *_ApduDataAdcReadBuilder) Build() (ApduDataAdcRead, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ApduDataAdcRead.deepCopy(), nil
}

func (b *_ApduDataAdcReadBuilder) MustBuild() ApduDataAdcRead {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ApduDataAdcReadBuilder) Done() ApduDataBuilder {
	return b.parentBuilder
}

func (b *_ApduDataAdcReadBuilder) buildForApduData() (ApduData, error) {
	return b.Build()
}

func (b *_ApduDataAdcReadBuilder) DeepCopy() any {
	_copy := b.CreateApduDataAdcReadBuilder().(*_ApduDataAdcReadBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateApduDataAdcReadBuilder creates a ApduDataAdcReadBuilder
func (b *_ApduDataAdcRead) CreateApduDataAdcReadBuilder() ApduDataAdcReadBuilder {
	if b == nil {
		return NewApduDataAdcReadBuilder()
	}
	return &_ApduDataAdcReadBuilder{_ApduDataAdcRead: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ApduDataAdcRead) GetApciType() uint8 {
	return 0x6
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ApduDataAdcRead) GetParent() ApduDataContract {
	return m.ApduDataContract
}

// Deprecated: use the interface for direct cast
func CastApduDataAdcRead(structType any) ApduDataAdcRead {
	if casted, ok := structType.(ApduDataAdcRead); ok {
		return casted
	}
	if casted, ok := structType.(*ApduDataAdcRead); ok {
		return *casted
	}
	return nil
}

func (m *_ApduDataAdcRead) GetTypeName() string {
	return "ApduDataAdcRead"
}

func (m *_ApduDataAdcRead) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ApduDataContract.(*_ApduData).GetLengthInBits(ctx))

	return lengthInBits
}

func (m *_ApduDataAdcRead) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ApduDataAdcRead) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ApduData, dataLength uint8) (__apduDataAdcRead ApduDataAdcRead, err error) {
	m.ApduDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ApduDataAdcRead"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ApduDataAdcRead")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("ApduDataAdcRead"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ApduDataAdcRead")
	}

	return m, nil
}

func (m *_ApduDataAdcRead) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ApduDataAdcRead) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ApduDataAdcRead"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ApduDataAdcRead")
		}

		if popErr := writeBuffer.PopContext("ApduDataAdcRead"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ApduDataAdcRead")
		}
		return nil
	}
	return m.ApduDataContract.(*_ApduData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ApduDataAdcRead) IsApduDataAdcRead() {}

func (m *_ApduDataAdcRead) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ApduDataAdcRead) deepCopy() *_ApduDataAdcRead {
	if m == nil {
		return nil
	}
	_ApduDataAdcReadCopy := &_ApduDataAdcRead{
		m.ApduDataContract.(*_ApduData).deepCopy(),
	}
	m.ApduDataContract.(*_ApduData)._SubType = m
	return _ApduDataAdcReadCopy
}

func (m *_ApduDataAdcRead) String() string {
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
