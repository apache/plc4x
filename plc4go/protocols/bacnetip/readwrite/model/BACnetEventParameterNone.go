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

// BACnetEventParameterNone is the corresponding interface of BACnetEventParameterNone
type BACnetEventParameterNone interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetEventParameter
	// GetNone returns None (property field)
	GetNone() BACnetContextTagNull
	// IsBACnetEventParameterNone is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetEventParameterNone()
	// CreateBuilder creates a BACnetEventParameterNoneBuilder
	CreateBACnetEventParameterNoneBuilder() BACnetEventParameterNoneBuilder
}

// _BACnetEventParameterNone is the data-structure of this message
type _BACnetEventParameterNone struct {
	BACnetEventParameterContract
	None BACnetContextTagNull
}

var _ BACnetEventParameterNone = (*_BACnetEventParameterNone)(nil)
var _ BACnetEventParameterRequirements = (*_BACnetEventParameterNone)(nil)

// NewBACnetEventParameterNone factory function for _BACnetEventParameterNone
func NewBACnetEventParameterNone(peekedTagHeader BACnetTagHeader, none BACnetContextTagNull) *_BACnetEventParameterNone {
	if none == nil {
		panic("none of type BACnetContextTagNull for BACnetEventParameterNone must not be nil")
	}
	_result := &_BACnetEventParameterNone{
		BACnetEventParameterContract: NewBACnetEventParameter(peekedTagHeader),
		None:                         none,
	}
	_result.BACnetEventParameterContract.(*_BACnetEventParameter)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetEventParameterNoneBuilder is a builder for BACnetEventParameterNone
type BACnetEventParameterNoneBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(none BACnetContextTagNull) BACnetEventParameterNoneBuilder
	// WithNone adds None (property field)
	WithNone(BACnetContextTagNull) BACnetEventParameterNoneBuilder
	// WithNoneBuilder adds None (property field) which is build by the builder
	WithNoneBuilder(func(BACnetContextTagNullBuilder) BACnetContextTagNullBuilder) BACnetEventParameterNoneBuilder
	// Build builds the BACnetEventParameterNone or returns an error if something is wrong
	Build() (BACnetEventParameterNone, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetEventParameterNone
}

// NewBACnetEventParameterNoneBuilder() creates a BACnetEventParameterNoneBuilder
func NewBACnetEventParameterNoneBuilder() BACnetEventParameterNoneBuilder {
	return &_BACnetEventParameterNoneBuilder{_BACnetEventParameterNone: new(_BACnetEventParameterNone)}
}

type _BACnetEventParameterNoneBuilder struct {
	*_BACnetEventParameterNone

	parentBuilder *_BACnetEventParameterBuilder

	err *utils.MultiError
}

var _ (BACnetEventParameterNoneBuilder) = (*_BACnetEventParameterNoneBuilder)(nil)

func (b *_BACnetEventParameterNoneBuilder) setParent(contract BACnetEventParameterContract) {
	b.BACnetEventParameterContract = contract
}

func (b *_BACnetEventParameterNoneBuilder) WithMandatoryFields(none BACnetContextTagNull) BACnetEventParameterNoneBuilder {
	return b.WithNone(none)
}

func (b *_BACnetEventParameterNoneBuilder) WithNone(none BACnetContextTagNull) BACnetEventParameterNoneBuilder {
	b.None = none
	return b
}

func (b *_BACnetEventParameterNoneBuilder) WithNoneBuilder(builderSupplier func(BACnetContextTagNullBuilder) BACnetContextTagNullBuilder) BACnetEventParameterNoneBuilder {
	builder := builderSupplier(b.None.CreateBACnetContextTagNullBuilder())
	var err error
	b.None, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagNullBuilder failed"))
	}
	return b
}

func (b *_BACnetEventParameterNoneBuilder) Build() (BACnetEventParameterNone, error) {
	if b.None == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'none' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetEventParameterNone.deepCopy(), nil
}

func (b *_BACnetEventParameterNoneBuilder) MustBuild() BACnetEventParameterNone {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetEventParameterNoneBuilder) Done() BACnetEventParameterBuilder {
	return b.parentBuilder
}

func (b *_BACnetEventParameterNoneBuilder) buildForBACnetEventParameter() (BACnetEventParameter, error) {
	return b.Build()
}

func (b *_BACnetEventParameterNoneBuilder) DeepCopy() any {
	_copy := b.CreateBACnetEventParameterNoneBuilder().(*_BACnetEventParameterNoneBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetEventParameterNoneBuilder creates a BACnetEventParameterNoneBuilder
func (b *_BACnetEventParameterNone) CreateBACnetEventParameterNoneBuilder() BACnetEventParameterNoneBuilder {
	if b == nil {
		return NewBACnetEventParameterNoneBuilder()
	}
	return &_BACnetEventParameterNoneBuilder{_BACnetEventParameterNone: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetEventParameterNone) GetParent() BACnetEventParameterContract {
	return m.BACnetEventParameterContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetEventParameterNone) GetNone() BACnetContextTagNull {
	return m.None
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetEventParameterNone(structType any) BACnetEventParameterNone {
	if casted, ok := structType.(BACnetEventParameterNone); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetEventParameterNone); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetEventParameterNone) GetTypeName() string {
	return "BACnetEventParameterNone"
}

func (m *_BACnetEventParameterNone) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetEventParameterContract.(*_BACnetEventParameter).GetLengthInBits(ctx))

	// Simple field (none)
	lengthInBits += m.None.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetEventParameterNone) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetEventParameterNone) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetEventParameter) (__bACnetEventParameterNone BACnetEventParameterNone, err error) {
	m.BACnetEventParameterContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetEventParameterNone"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetEventParameterNone")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	none, err := ReadSimpleField[BACnetContextTagNull](ctx, "none", ReadComplex[BACnetContextTagNull](BACnetContextTagParseWithBufferProducer[BACnetContextTagNull]((uint8)(uint8(20)), (BACnetDataType)(BACnetDataType_NULL)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'none' field"))
	}
	m.None = none

	if closeErr := readBuffer.CloseContext("BACnetEventParameterNone"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetEventParameterNone")
	}

	return m, nil
}

func (m *_BACnetEventParameterNone) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetEventParameterNone) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetEventParameterNone"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetEventParameterNone")
		}

		if err := WriteSimpleField[BACnetContextTagNull](ctx, "none", m.GetNone(), WriteComplex[BACnetContextTagNull](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'none' field")
		}

		if popErr := writeBuffer.PopContext("BACnetEventParameterNone"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetEventParameterNone")
		}
		return nil
	}
	return m.BACnetEventParameterContract.(*_BACnetEventParameter).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetEventParameterNone) IsBACnetEventParameterNone() {}

func (m *_BACnetEventParameterNone) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetEventParameterNone) deepCopy() *_BACnetEventParameterNone {
	if m == nil {
		return nil
	}
	_BACnetEventParameterNoneCopy := &_BACnetEventParameterNone{
		m.BACnetEventParameterContract.(*_BACnetEventParameter).deepCopy(),
		m.None.DeepCopy().(BACnetContextTagNull),
	}
	m.BACnetEventParameterContract.(*_BACnetEventParameter)._SubType = m
	return _BACnetEventParameterNoneCopy
}

func (m *_BACnetEventParameterNone) String() string {
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
