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

// BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned is the corresponding interface of BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned
type BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetNotificationParametersChangeOfDiscreteValueNewValue
	// GetUnsignedValue returns UnsignedValue (property field)
	GetUnsignedValue() BACnetApplicationTagUnsignedInteger
	// IsBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned()
	// CreateBuilder creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder
	CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder
}

// _BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned is the data-structure of this message
type _BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned struct {
	BACnetNotificationParametersChangeOfDiscreteValueNewValueContract
	UnsignedValue BACnetApplicationTagUnsignedInteger
}

var _ BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned)(nil)
var _ BACnetNotificationParametersChangeOfDiscreteValueNewValueRequirements = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned)(nil)

// NewBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned factory function for _BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned
func NewBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, unsignedValue BACnetApplicationTagUnsignedInteger, tagNumber uint8) *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned {
	if unsignedValue == nil {
		panic("unsignedValue of type BACnetApplicationTagUnsignedInteger for BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned must not be nil")
	}
	_result := &_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned{
		BACnetNotificationParametersChangeOfDiscreteValueNewValueContract: NewBACnetNotificationParametersChangeOfDiscreteValueNewValue(openingTag, peekedTagHeader, closingTag, tagNumber),
		UnsignedValue: unsignedValue,
	}
	_result.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder is a builder for BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned
type BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(unsignedValue BACnetApplicationTagUnsignedInteger) BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder
	// WithUnsignedValue adds UnsignedValue (property field)
	WithUnsignedValue(BACnetApplicationTagUnsignedInteger) BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder
	// WithUnsignedValueBuilder adds UnsignedValue (property field) which is build by the builder
	WithUnsignedValueBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder
	// Build builds the BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned or returns an error if something is wrong
	Build() (BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned
}

// NewBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder() creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder
func NewBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder {
	return &_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder{_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned: new(_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned)}
}

type _BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder struct {
	*_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned

	parentBuilder *_BACnetNotificationParametersChangeOfDiscreteValueNewValueBuilder

	err *utils.MultiError
}

var _ (BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder) = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder)(nil)

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder) setParent(contract BACnetNotificationParametersChangeOfDiscreteValueNewValueContract) {
	b.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract = contract
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder) WithMandatoryFields(unsignedValue BACnetApplicationTagUnsignedInteger) BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder {
	return b.WithUnsignedValue(unsignedValue)
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder) WithUnsignedValue(unsignedValue BACnetApplicationTagUnsignedInteger) BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder {
	b.UnsignedValue = unsignedValue
	return b
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder) WithUnsignedValueBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder {
	builder := builderSupplier(b.UnsignedValue.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.UnsignedValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder) Build() (BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned, error) {
	if b.UnsignedValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'unsignedValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned.deepCopy(), nil
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder) MustBuild() BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder) Done() BACnetNotificationParametersChangeOfDiscreteValueNewValueBuilder {
	return b.parentBuilder
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder) buildForBACnetNotificationParametersChangeOfDiscreteValueNewValue() (BACnetNotificationParametersChangeOfDiscreteValueNewValue, error) {
	return b.Build()
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder) DeepCopy() any {
	_copy := b.CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder().(*_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder
func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder {
	if b == nil {
		return NewBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder()
	}
	return &_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedBuilder{_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned: b.deepCopy()}
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

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) GetParent() BACnetNotificationParametersChangeOfDiscreteValueNewValueContract {
	return m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) GetUnsignedValue() BACnetApplicationTagUnsignedInteger {
	return m.UnsignedValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned(structType any) BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned {
	if casted, ok := structType.(BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) GetTypeName() string {
	return "BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned"
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).GetLengthInBits(ctx))

	// Simple field (unsignedValue)
	lengthInBits += m.UnsignedValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetNotificationParametersChangeOfDiscreteValueNewValue, tagNumber uint8) (__bACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned, err error) {
	m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	unsignedValue, err := ReadSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "unsignedValue", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'unsignedValue' field"))
	}
	m.UnsignedValue = unsignedValue

	if closeErr := readBuffer.CloseContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned")
	}

	return m, nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned")
		}

		if err := WriteSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "unsignedValue", m.GetUnsignedValue(), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'unsignedValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned")
		}
		return nil
	}
	return m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) IsBACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned() {
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) deepCopy() *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned {
	if m == nil {
		return nil
	}
	_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedCopy := &_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned{
		m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).deepCopy(),
		m.UnsignedValue.DeepCopy().(BACnetApplicationTagUnsignedInteger),
	}
	m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue)._SubType = m
	return _BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsignedCopy
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueUnsigned) String() string {
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
