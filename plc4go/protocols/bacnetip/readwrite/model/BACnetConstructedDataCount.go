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

// BACnetConstructedDataCount is the corresponding interface of BACnetConstructedDataCount
type BACnetConstructedDataCount interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetCount returns Count (property field)
	GetCount() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
	// IsBACnetConstructedDataCount is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataCount()
	// CreateBuilder creates a BACnetConstructedDataCountBuilder
	CreateBACnetConstructedDataCountBuilder() BACnetConstructedDataCountBuilder
}

// _BACnetConstructedDataCount is the data-structure of this message
type _BACnetConstructedDataCount struct {
	BACnetConstructedDataContract
	Count BACnetApplicationTagUnsignedInteger
}

var _ BACnetConstructedDataCount = (*_BACnetConstructedDataCount)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataCount)(nil)

// NewBACnetConstructedDataCount factory function for _BACnetConstructedDataCount
func NewBACnetConstructedDataCount(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, count BACnetApplicationTagUnsignedInteger, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataCount {
	if count == nil {
		panic("count of type BACnetApplicationTagUnsignedInteger for BACnetConstructedDataCount must not be nil")
	}
	_result := &_BACnetConstructedDataCount{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		Count:                         count,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataCountBuilder is a builder for BACnetConstructedDataCount
type BACnetConstructedDataCountBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(count BACnetApplicationTagUnsignedInteger) BACnetConstructedDataCountBuilder
	// WithCount adds Count (property field)
	WithCount(BACnetApplicationTagUnsignedInteger) BACnetConstructedDataCountBuilder
	// WithCountBuilder adds Count (property field) which is build by the builder
	WithCountBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataCountBuilder
	// Build builds the BACnetConstructedDataCount or returns an error if something is wrong
	Build() (BACnetConstructedDataCount, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataCount
}

// NewBACnetConstructedDataCountBuilder() creates a BACnetConstructedDataCountBuilder
func NewBACnetConstructedDataCountBuilder() BACnetConstructedDataCountBuilder {
	return &_BACnetConstructedDataCountBuilder{_BACnetConstructedDataCount: new(_BACnetConstructedDataCount)}
}

type _BACnetConstructedDataCountBuilder struct {
	*_BACnetConstructedDataCount

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataCountBuilder) = (*_BACnetConstructedDataCountBuilder)(nil)

func (b *_BACnetConstructedDataCountBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataCountBuilder) WithMandatoryFields(count BACnetApplicationTagUnsignedInteger) BACnetConstructedDataCountBuilder {
	return b.WithCount(count)
}

func (b *_BACnetConstructedDataCountBuilder) WithCount(count BACnetApplicationTagUnsignedInteger) BACnetConstructedDataCountBuilder {
	b.Count = count
	return b
}

func (b *_BACnetConstructedDataCountBuilder) WithCountBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataCountBuilder {
	builder := builderSupplier(b.Count.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.Count, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataCountBuilder) Build() (BACnetConstructedDataCount, error) {
	if b.Count == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'count' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataCount.deepCopy(), nil
}

func (b *_BACnetConstructedDataCountBuilder) MustBuild() BACnetConstructedDataCount {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataCountBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataCountBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataCountBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataCountBuilder().(*_BACnetConstructedDataCountBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataCountBuilder creates a BACnetConstructedDataCountBuilder
func (b *_BACnetConstructedDataCount) CreateBACnetConstructedDataCountBuilder() BACnetConstructedDataCountBuilder {
	if b == nil {
		return NewBACnetConstructedDataCountBuilder()
	}
	return &_BACnetConstructedDataCountBuilder{_BACnetConstructedDataCount: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataCount) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataCount) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_COUNT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataCount) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataCount) GetCount() BACnetApplicationTagUnsignedInteger {
	return m.Count
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataCount) GetActualValue() BACnetApplicationTagUnsignedInteger {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagUnsignedInteger(m.GetCount())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataCount(structType any) BACnetConstructedDataCount {
	if casted, ok := structType.(BACnetConstructedDataCount); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataCount); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataCount) GetTypeName() string {
	return "BACnetConstructedDataCount"
}

func (m *_BACnetConstructedDataCount) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (count)
	lengthInBits += m.Count.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataCount) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataCount) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataCount BACnetConstructedDataCount, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataCount"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataCount")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	count, err := ReadSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "count", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'count' field"))
	}
	m.Count = count

	actualValue, err := ReadVirtualField[BACnetApplicationTagUnsignedInteger](ctx, "actualValue", (*BACnetApplicationTagUnsignedInteger)(nil), count)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataCount"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataCount")
	}

	return m, nil
}

func (m *_BACnetConstructedDataCount) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataCount) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataCount"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataCount")
		}

		if err := WriteSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "count", m.GetCount(), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'count' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataCount"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataCount")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataCount) IsBACnetConstructedDataCount() {}

func (m *_BACnetConstructedDataCount) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataCount) deepCopy() *_BACnetConstructedDataCount {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataCountCopy := &_BACnetConstructedDataCount{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.Count.DeepCopy().(BACnetApplicationTagUnsignedInteger),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataCountCopy
}

func (m *_BACnetConstructedDataCount) String() string {
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
