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

// BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger is the corresponding interface of BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger
type BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetNotificationParametersChangeOfDiscreteValueNewValue
	// GetIntegerValue returns IntegerValue (property field)
	GetIntegerValue() BACnetApplicationTagSignedInteger
	// IsBACnetNotificationParametersChangeOfDiscreteValueNewValueInteger is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetNotificationParametersChangeOfDiscreteValueNewValueInteger()
	// CreateBuilder creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder
	CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder
}

// _BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger is the data-structure of this message
type _BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger struct {
	BACnetNotificationParametersChangeOfDiscreteValueNewValueContract
	IntegerValue BACnetApplicationTagSignedInteger
}

var _ BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger)(nil)
var _ BACnetNotificationParametersChangeOfDiscreteValueNewValueRequirements = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger)(nil)

// NewBACnetNotificationParametersChangeOfDiscreteValueNewValueInteger factory function for _BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger
func NewBACnetNotificationParametersChangeOfDiscreteValueNewValueInteger(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, integerValue BACnetApplicationTagSignedInteger, tagNumber uint8) *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger {
	if integerValue == nil {
		panic("integerValue of type BACnetApplicationTagSignedInteger for BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger must not be nil")
	}
	_result := &_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger{
		BACnetNotificationParametersChangeOfDiscreteValueNewValueContract: NewBACnetNotificationParametersChangeOfDiscreteValueNewValue(openingTag, peekedTagHeader, closingTag, tagNumber),
		IntegerValue: integerValue,
	}
	_result.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder is a builder for BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger
type BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(integerValue BACnetApplicationTagSignedInteger) BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder
	// WithIntegerValue adds IntegerValue (property field)
	WithIntegerValue(BACnetApplicationTagSignedInteger) BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder
	// WithIntegerValueBuilder adds IntegerValue (property field) which is build by the builder
	WithIntegerValueBuilder(func(BACnetApplicationTagSignedIntegerBuilder) BACnetApplicationTagSignedIntegerBuilder) BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder
	// Build builds the BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger or returns an error if something is wrong
	Build() (BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger
}

// NewBACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder() creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder
func NewBACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder {
	return &_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder{_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger: new(_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger)}
}

type _BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder struct {
	*_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger

	parentBuilder *_BACnetNotificationParametersChangeOfDiscreteValueNewValueBuilder

	err *utils.MultiError
}

var _ (BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder) = (*_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder)(nil)

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder) setParent(contract BACnetNotificationParametersChangeOfDiscreteValueNewValueContract) {
	b.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract = contract
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder) WithMandatoryFields(integerValue BACnetApplicationTagSignedInteger) BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder {
	return b.WithIntegerValue(integerValue)
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder) WithIntegerValue(integerValue BACnetApplicationTagSignedInteger) BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder {
	b.IntegerValue = integerValue
	return b
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder) WithIntegerValueBuilder(builderSupplier func(BACnetApplicationTagSignedIntegerBuilder) BACnetApplicationTagSignedIntegerBuilder) BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder {
	builder := builderSupplier(b.IntegerValue.CreateBACnetApplicationTagSignedIntegerBuilder())
	var err error
	b.IntegerValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagSignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder) Build() (BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger, error) {
	if b.IntegerValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'integerValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger.deepCopy(), nil
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder) MustBuild() BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder) Done() BACnetNotificationParametersChangeOfDiscreteValueNewValueBuilder {
	return b.parentBuilder
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder) buildForBACnetNotificationParametersChangeOfDiscreteValueNewValue() (BACnetNotificationParametersChangeOfDiscreteValueNewValue, error) {
	return b.Build()
}

func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder) DeepCopy() any {
	_copy := b.CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder().(*_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder creates a BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder
func (b *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) CreateBACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder() BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder {
	if b == nil {
		return NewBACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder()
	}
	return &_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerBuilder{_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger: b.deepCopy()}
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

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) GetParent() BACnetNotificationParametersChangeOfDiscreteValueNewValueContract {
	return m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) GetIntegerValue() BACnetApplicationTagSignedInteger {
	return m.IntegerValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetNotificationParametersChangeOfDiscreteValueNewValueInteger(structType any) BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger {
	if casted, ok := structType.(BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) GetTypeName() string {
	return "BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger"
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).GetLengthInBits(ctx))

	// Simple field (integerValue)
	lengthInBits += m.IntegerValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetNotificationParametersChangeOfDiscreteValueNewValue, tagNumber uint8) (__bACnetNotificationParametersChangeOfDiscreteValueNewValueInteger BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger, err error) {
	m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	integerValue, err := ReadSimpleField[BACnetApplicationTagSignedInteger](ctx, "integerValue", ReadComplex[BACnetApplicationTagSignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagSignedInteger](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'integerValue' field"))
	}
	m.IntegerValue = integerValue

	if closeErr := readBuffer.CloseContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger")
	}

	return m, nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger")
		}

		if err := WriteSimpleField[BACnetApplicationTagSignedInteger](ctx, "integerValue", m.GetIntegerValue(), WriteComplex[BACnetApplicationTagSignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'integerValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger")
		}
		return nil
	}
	return m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) IsBACnetNotificationParametersChangeOfDiscreteValueNewValueInteger() {
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) deepCopy() *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger {
	if m == nil {
		return nil
	}
	_BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerCopy := &_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger{
		m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue).deepCopy(),
		m.IntegerValue.DeepCopy().(BACnetApplicationTagSignedInteger),
	}
	m.BACnetNotificationParametersChangeOfDiscreteValueNewValueContract.(*_BACnetNotificationParametersChangeOfDiscreteValueNewValue)._SubType = m
	return _BACnetNotificationParametersChangeOfDiscreteValueNewValueIntegerCopy
}

func (m *_BACnetNotificationParametersChangeOfDiscreteValueNewValueInteger) String() string {
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
