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

// BACnetConstructedDataDatePatternValueRelinquishDefault is the corresponding interface of BACnetConstructedDataDatePatternValueRelinquishDefault
type BACnetConstructedDataDatePatternValueRelinquishDefault interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetRelinquishDefault returns RelinquishDefault (property field)
	GetRelinquishDefault() BACnetApplicationTagDate
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagDate
	// IsBACnetConstructedDataDatePatternValueRelinquishDefault is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataDatePatternValueRelinquishDefault()
	// CreateBuilder creates a BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder
	CreateBACnetConstructedDataDatePatternValueRelinquishDefaultBuilder() BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder
}

// _BACnetConstructedDataDatePatternValueRelinquishDefault is the data-structure of this message
type _BACnetConstructedDataDatePatternValueRelinquishDefault struct {
	BACnetConstructedDataContract
	RelinquishDefault BACnetApplicationTagDate
}

var _ BACnetConstructedDataDatePatternValueRelinquishDefault = (*_BACnetConstructedDataDatePatternValueRelinquishDefault)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataDatePatternValueRelinquishDefault)(nil)

// NewBACnetConstructedDataDatePatternValueRelinquishDefault factory function for _BACnetConstructedDataDatePatternValueRelinquishDefault
func NewBACnetConstructedDataDatePatternValueRelinquishDefault(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, relinquishDefault BACnetApplicationTagDate, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataDatePatternValueRelinquishDefault {
	if relinquishDefault == nil {
		panic("relinquishDefault of type BACnetApplicationTagDate for BACnetConstructedDataDatePatternValueRelinquishDefault must not be nil")
	}
	_result := &_BACnetConstructedDataDatePatternValueRelinquishDefault{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		RelinquishDefault:             relinquishDefault,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder is a builder for BACnetConstructedDataDatePatternValueRelinquishDefault
type BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(relinquishDefault BACnetApplicationTagDate) BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder
	// WithRelinquishDefault adds RelinquishDefault (property field)
	WithRelinquishDefault(BACnetApplicationTagDate) BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder
	// WithRelinquishDefaultBuilder adds RelinquishDefault (property field) which is build by the builder
	WithRelinquishDefaultBuilder(func(BACnetApplicationTagDateBuilder) BACnetApplicationTagDateBuilder) BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder
	// Build builds the BACnetConstructedDataDatePatternValueRelinquishDefault or returns an error if something is wrong
	Build() (BACnetConstructedDataDatePatternValueRelinquishDefault, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataDatePatternValueRelinquishDefault
}

// NewBACnetConstructedDataDatePatternValueRelinquishDefaultBuilder() creates a BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder
func NewBACnetConstructedDataDatePatternValueRelinquishDefaultBuilder() BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder {
	return &_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder{_BACnetConstructedDataDatePatternValueRelinquishDefault: new(_BACnetConstructedDataDatePatternValueRelinquishDefault)}
}

type _BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder struct {
	*_BACnetConstructedDataDatePatternValueRelinquishDefault

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder) = (*_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder)(nil)

func (b *_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder) WithMandatoryFields(relinquishDefault BACnetApplicationTagDate) BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder {
	return b.WithRelinquishDefault(relinquishDefault)
}

func (b *_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder) WithRelinquishDefault(relinquishDefault BACnetApplicationTagDate) BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder {
	b.RelinquishDefault = relinquishDefault
	return b
}

func (b *_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder) WithRelinquishDefaultBuilder(builderSupplier func(BACnetApplicationTagDateBuilder) BACnetApplicationTagDateBuilder) BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder {
	builder := builderSupplier(b.RelinquishDefault.CreateBACnetApplicationTagDateBuilder())
	var err error
	b.RelinquishDefault, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagDateBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder) Build() (BACnetConstructedDataDatePatternValueRelinquishDefault, error) {
	if b.RelinquishDefault == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'relinquishDefault' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataDatePatternValueRelinquishDefault.deepCopy(), nil
}

func (b *_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder) MustBuild() BACnetConstructedDataDatePatternValueRelinquishDefault {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataDatePatternValueRelinquishDefaultBuilder().(*_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataDatePatternValueRelinquishDefaultBuilder creates a BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder
func (b *_BACnetConstructedDataDatePatternValueRelinquishDefault) CreateBACnetConstructedDataDatePatternValueRelinquishDefaultBuilder() BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder {
	if b == nil {
		return NewBACnetConstructedDataDatePatternValueRelinquishDefaultBuilder()
	}
	return &_BACnetConstructedDataDatePatternValueRelinquishDefaultBuilder{_BACnetConstructedDataDatePatternValueRelinquishDefault: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_DATEPATTERN_VALUE
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_RELINQUISH_DEFAULT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) GetRelinquishDefault() BACnetApplicationTagDate {
	return m.RelinquishDefault
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) GetActualValue() BACnetApplicationTagDate {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagDate(m.GetRelinquishDefault())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataDatePatternValueRelinquishDefault(structType any) BACnetConstructedDataDatePatternValueRelinquishDefault {
	if casted, ok := structType.(BACnetConstructedDataDatePatternValueRelinquishDefault); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataDatePatternValueRelinquishDefault); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) GetTypeName() string {
	return "BACnetConstructedDataDatePatternValueRelinquishDefault"
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (relinquishDefault)
	lengthInBits += m.RelinquishDefault.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataDatePatternValueRelinquishDefault BACnetConstructedDataDatePatternValueRelinquishDefault, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataDatePatternValueRelinquishDefault"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataDatePatternValueRelinquishDefault")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	relinquishDefault, err := ReadSimpleField[BACnetApplicationTagDate](ctx, "relinquishDefault", ReadComplex[BACnetApplicationTagDate](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagDate](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'relinquishDefault' field"))
	}
	m.RelinquishDefault = relinquishDefault

	actualValue, err := ReadVirtualField[BACnetApplicationTagDate](ctx, "actualValue", (*BACnetApplicationTagDate)(nil), relinquishDefault)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataDatePatternValueRelinquishDefault"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataDatePatternValueRelinquishDefault")
	}

	return m, nil
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataDatePatternValueRelinquishDefault"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataDatePatternValueRelinquishDefault")
		}

		if err := WriteSimpleField[BACnetApplicationTagDate](ctx, "relinquishDefault", m.GetRelinquishDefault(), WriteComplex[BACnetApplicationTagDate](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'relinquishDefault' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataDatePatternValueRelinquishDefault"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataDatePatternValueRelinquishDefault")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) IsBACnetConstructedDataDatePatternValueRelinquishDefault() {
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) deepCopy() *_BACnetConstructedDataDatePatternValueRelinquishDefault {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataDatePatternValueRelinquishDefaultCopy := &_BACnetConstructedDataDatePatternValueRelinquishDefault{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.RelinquishDefault.DeepCopy().(BACnetApplicationTagDate),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataDatePatternValueRelinquishDefaultCopy
}

func (m *_BACnetConstructedDataDatePatternValueRelinquishDefault) String() string {
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
