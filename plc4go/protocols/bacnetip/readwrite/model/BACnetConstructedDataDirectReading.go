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

// BACnetConstructedDataDirectReading is the corresponding interface of BACnetConstructedDataDirectReading
type BACnetConstructedDataDirectReading interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetDirectReading returns DirectReading (property field)
	GetDirectReading() BACnetApplicationTagReal
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagReal
	// IsBACnetConstructedDataDirectReading is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataDirectReading()
	// CreateBuilder creates a BACnetConstructedDataDirectReadingBuilder
	CreateBACnetConstructedDataDirectReadingBuilder() BACnetConstructedDataDirectReadingBuilder
}

// _BACnetConstructedDataDirectReading is the data-structure of this message
type _BACnetConstructedDataDirectReading struct {
	BACnetConstructedDataContract
	DirectReading BACnetApplicationTagReal
}

var _ BACnetConstructedDataDirectReading = (*_BACnetConstructedDataDirectReading)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataDirectReading)(nil)

// NewBACnetConstructedDataDirectReading factory function for _BACnetConstructedDataDirectReading
func NewBACnetConstructedDataDirectReading(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, directReading BACnetApplicationTagReal, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataDirectReading {
	if directReading == nil {
		panic("directReading of type BACnetApplicationTagReal for BACnetConstructedDataDirectReading must not be nil")
	}
	_result := &_BACnetConstructedDataDirectReading{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		DirectReading:                 directReading,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataDirectReadingBuilder is a builder for BACnetConstructedDataDirectReading
type BACnetConstructedDataDirectReadingBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(directReading BACnetApplicationTagReal) BACnetConstructedDataDirectReadingBuilder
	// WithDirectReading adds DirectReading (property field)
	WithDirectReading(BACnetApplicationTagReal) BACnetConstructedDataDirectReadingBuilder
	// WithDirectReadingBuilder adds DirectReading (property field) which is build by the builder
	WithDirectReadingBuilder(func(BACnetApplicationTagRealBuilder) BACnetApplicationTagRealBuilder) BACnetConstructedDataDirectReadingBuilder
	// Build builds the BACnetConstructedDataDirectReading or returns an error if something is wrong
	Build() (BACnetConstructedDataDirectReading, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataDirectReading
}

// NewBACnetConstructedDataDirectReadingBuilder() creates a BACnetConstructedDataDirectReadingBuilder
func NewBACnetConstructedDataDirectReadingBuilder() BACnetConstructedDataDirectReadingBuilder {
	return &_BACnetConstructedDataDirectReadingBuilder{_BACnetConstructedDataDirectReading: new(_BACnetConstructedDataDirectReading)}
}

type _BACnetConstructedDataDirectReadingBuilder struct {
	*_BACnetConstructedDataDirectReading

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataDirectReadingBuilder) = (*_BACnetConstructedDataDirectReadingBuilder)(nil)

func (b *_BACnetConstructedDataDirectReadingBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataDirectReadingBuilder) WithMandatoryFields(directReading BACnetApplicationTagReal) BACnetConstructedDataDirectReadingBuilder {
	return b.WithDirectReading(directReading)
}

func (b *_BACnetConstructedDataDirectReadingBuilder) WithDirectReading(directReading BACnetApplicationTagReal) BACnetConstructedDataDirectReadingBuilder {
	b.DirectReading = directReading
	return b
}

func (b *_BACnetConstructedDataDirectReadingBuilder) WithDirectReadingBuilder(builderSupplier func(BACnetApplicationTagRealBuilder) BACnetApplicationTagRealBuilder) BACnetConstructedDataDirectReadingBuilder {
	builder := builderSupplier(b.DirectReading.CreateBACnetApplicationTagRealBuilder())
	var err error
	b.DirectReading, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagRealBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataDirectReadingBuilder) Build() (BACnetConstructedDataDirectReading, error) {
	if b.DirectReading == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'directReading' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataDirectReading.deepCopy(), nil
}

func (b *_BACnetConstructedDataDirectReadingBuilder) MustBuild() BACnetConstructedDataDirectReading {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataDirectReadingBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataDirectReadingBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataDirectReadingBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataDirectReadingBuilder().(*_BACnetConstructedDataDirectReadingBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataDirectReadingBuilder creates a BACnetConstructedDataDirectReadingBuilder
func (b *_BACnetConstructedDataDirectReading) CreateBACnetConstructedDataDirectReadingBuilder() BACnetConstructedDataDirectReadingBuilder {
	if b == nil {
		return NewBACnetConstructedDataDirectReadingBuilder()
	}
	return &_BACnetConstructedDataDirectReadingBuilder{_BACnetConstructedDataDirectReading: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataDirectReading) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataDirectReading) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_DIRECT_READING
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataDirectReading) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataDirectReading) GetDirectReading() BACnetApplicationTagReal {
	return m.DirectReading
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataDirectReading) GetActualValue() BACnetApplicationTagReal {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagReal(m.GetDirectReading())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataDirectReading(structType any) BACnetConstructedDataDirectReading {
	if casted, ok := structType.(BACnetConstructedDataDirectReading); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataDirectReading); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataDirectReading) GetTypeName() string {
	return "BACnetConstructedDataDirectReading"
}

func (m *_BACnetConstructedDataDirectReading) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (directReading)
	lengthInBits += m.DirectReading.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataDirectReading) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataDirectReading) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataDirectReading BACnetConstructedDataDirectReading, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataDirectReading"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataDirectReading")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	directReading, err := ReadSimpleField[BACnetApplicationTagReal](ctx, "directReading", ReadComplex[BACnetApplicationTagReal](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagReal](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'directReading' field"))
	}
	m.DirectReading = directReading

	actualValue, err := ReadVirtualField[BACnetApplicationTagReal](ctx, "actualValue", (*BACnetApplicationTagReal)(nil), directReading)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataDirectReading"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataDirectReading")
	}

	return m, nil
}

func (m *_BACnetConstructedDataDirectReading) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataDirectReading) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataDirectReading"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataDirectReading")
		}

		if err := WriteSimpleField[BACnetApplicationTagReal](ctx, "directReading", m.GetDirectReading(), WriteComplex[BACnetApplicationTagReal](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'directReading' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataDirectReading"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataDirectReading")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataDirectReading) IsBACnetConstructedDataDirectReading() {}

func (m *_BACnetConstructedDataDirectReading) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataDirectReading) deepCopy() *_BACnetConstructedDataDirectReading {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataDirectReadingCopy := &_BACnetConstructedDataDirectReading{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.DirectReading.DeepCopy().(BACnetApplicationTagReal),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataDirectReadingCopy
}

func (m *_BACnetConstructedDataDirectReading) String() string {
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
