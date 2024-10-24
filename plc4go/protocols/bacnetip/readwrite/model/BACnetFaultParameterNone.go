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

// BACnetFaultParameterNone is the corresponding interface of BACnetFaultParameterNone
type BACnetFaultParameterNone interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetFaultParameter
	// GetNone returns None (property field)
	GetNone() BACnetContextTagNull
	// IsBACnetFaultParameterNone is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetFaultParameterNone()
	// CreateBuilder creates a BACnetFaultParameterNoneBuilder
	CreateBACnetFaultParameterNoneBuilder() BACnetFaultParameterNoneBuilder
}

// _BACnetFaultParameterNone is the data-structure of this message
type _BACnetFaultParameterNone struct {
	BACnetFaultParameterContract
	None BACnetContextTagNull
}

var _ BACnetFaultParameterNone = (*_BACnetFaultParameterNone)(nil)
var _ BACnetFaultParameterRequirements = (*_BACnetFaultParameterNone)(nil)

// NewBACnetFaultParameterNone factory function for _BACnetFaultParameterNone
func NewBACnetFaultParameterNone(peekedTagHeader BACnetTagHeader, none BACnetContextTagNull) *_BACnetFaultParameterNone {
	if none == nil {
		panic("none of type BACnetContextTagNull for BACnetFaultParameterNone must not be nil")
	}
	_result := &_BACnetFaultParameterNone{
		BACnetFaultParameterContract: NewBACnetFaultParameter(peekedTagHeader),
		None:                         none,
	}
	_result.BACnetFaultParameterContract.(*_BACnetFaultParameter)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetFaultParameterNoneBuilder is a builder for BACnetFaultParameterNone
type BACnetFaultParameterNoneBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(none BACnetContextTagNull) BACnetFaultParameterNoneBuilder
	// WithNone adds None (property field)
	WithNone(BACnetContextTagNull) BACnetFaultParameterNoneBuilder
	// WithNoneBuilder adds None (property field) which is build by the builder
	WithNoneBuilder(func(BACnetContextTagNullBuilder) BACnetContextTagNullBuilder) BACnetFaultParameterNoneBuilder
	// Build builds the BACnetFaultParameterNone or returns an error if something is wrong
	Build() (BACnetFaultParameterNone, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetFaultParameterNone
}

// NewBACnetFaultParameterNoneBuilder() creates a BACnetFaultParameterNoneBuilder
func NewBACnetFaultParameterNoneBuilder() BACnetFaultParameterNoneBuilder {
	return &_BACnetFaultParameterNoneBuilder{_BACnetFaultParameterNone: new(_BACnetFaultParameterNone)}
}

type _BACnetFaultParameterNoneBuilder struct {
	*_BACnetFaultParameterNone

	parentBuilder *_BACnetFaultParameterBuilder

	err *utils.MultiError
}

var _ (BACnetFaultParameterNoneBuilder) = (*_BACnetFaultParameterNoneBuilder)(nil)

func (b *_BACnetFaultParameterNoneBuilder) setParent(contract BACnetFaultParameterContract) {
	b.BACnetFaultParameterContract = contract
}

func (b *_BACnetFaultParameterNoneBuilder) WithMandatoryFields(none BACnetContextTagNull) BACnetFaultParameterNoneBuilder {
	return b.WithNone(none)
}

func (b *_BACnetFaultParameterNoneBuilder) WithNone(none BACnetContextTagNull) BACnetFaultParameterNoneBuilder {
	b.None = none
	return b
}

func (b *_BACnetFaultParameterNoneBuilder) WithNoneBuilder(builderSupplier func(BACnetContextTagNullBuilder) BACnetContextTagNullBuilder) BACnetFaultParameterNoneBuilder {
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

func (b *_BACnetFaultParameterNoneBuilder) Build() (BACnetFaultParameterNone, error) {
	if b.None == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'none' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetFaultParameterNone.deepCopy(), nil
}

func (b *_BACnetFaultParameterNoneBuilder) MustBuild() BACnetFaultParameterNone {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetFaultParameterNoneBuilder) Done() BACnetFaultParameterBuilder {
	return b.parentBuilder
}

func (b *_BACnetFaultParameterNoneBuilder) buildForBACnetFaultParameter() (BACnetFaultParameter, error) {
	return b.Build()
}

func (b *_BACnetFaultParameterNoneBuilder) DeepCopy() any {
	_copy := b.CreateBACnetFaultParameterNoneBuilder().(*_BACnetFaultParameterNoneBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetFaultParameterNoneBuilder creates a BACnetFaultParameterNoneBuilder
func (b *_BACnetFaultParameterNone) CreateBACnetFaultParameterNoneBuilder() BACnetFaultParameterNoneBuilder {
	if b == nil {
		return NewBACnetFaultParameterNoneBuilder()
	}
	return &_BACnetFaultParameterNoneBuilder{_BACnetFaultParameterNone: b.deepCopy()}
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

func (m *_BACnetFaultParameterNone) GetParent() BACnetFaultParameterContract {
	return m.BACnetFaultParameterContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetFaultParameterNone) GetNone() BACnetContextTagNull {
	return m.None
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetFaultParameterNone(structType any) BACnetFaultParameterNone {
	if casted, ok := structType.(BACnetFaultParameterNone); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetFaultParameterNone); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetFaultParameterNone) GetTypeName() string {
	return "BACnetFaultParameterNone"
}

func (m *_BACnetFaultParameterNone) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetFaultParameterContract.(*_BACnetFaultParameter).GetLengthInBits(ctx))

	// Simple field (none)
	lengthInBits += m.None.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetFaultParameterNone) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetFaultParameterNone) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetFaultParameter) (__bACnetFaultParameterNone BACnetFaultParameterNone, err error) {
	m.BACnetFaultParameterContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetFaultParameterNone"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetFaultParameterNone")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	none, err := ReadSimpleField[BACnetContextTagNull](ctx, "none", ReadComplex[BACnetContextTagNull](BACnetContextTagParseWithBufferProducer[BACnetContextTagNull]((uint8)(uint8(0)), (BACnetDataType)(BACnetDataType_NULL)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'none' field"))
	}
	m.None = none

	if closeErr := readBuffer.CloseContext("BACnetFaultParameterNone"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetFaultParameterNone")
	}

	return m, nil
}

func (m *_BACnetFaultParameterNone) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetFaultParameterNone) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetFaultParameterNone"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetFaultParameterNone")
		}

		if err := WriteSimpleField[BACnetContextTagNull](ctx, "none", m.GetNone(), WriteComplex[BACnetContextTagNull](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'none' field")
		}

		if popErr := writeBuffer.PopContext("BACnetFaultParameterNone"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetFaultParameterNone")
		}
		return nil
	}
	return m.BACnetFaultParameterContract.(*_BACnetFaultParameter).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetFaultParameterNone) IsBACnetFaultParameterNone() {}

func (m *_BACnetFaultParameterNone) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetFaultParameterNone) deepCopy() *_BACnetFaultParameterNone {
	if m == nil {
		return nil
	}
	_BACnetFaultParameterNoneCopy := &_BACnetFaultParameterNone{
		m.BACnetFaultParameterContract.(*_BACnetFaultParameter).deepCopy(),
		m.None.DeepCopy().(BACnetContextTagNull),
	}
	m.BACnetFaultParameterContract.(*_BACnetFaultParameter)._SubType = m
	return _BACnetFaultParameterNoneCopy
}

func (m *_BACnetFaultParameterNone) String() string {
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
