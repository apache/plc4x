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

// BACnetConstructedDataBinaryValuePresentValue is the corresponding interface of BACnetConstructedDataBinaryValuePresentValue
type BACnetConstructedDataBinaryValuePresentValue interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetPresentValue returns PresentValue (property field)
	GetPresentValue() BACnetBinaryPVTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetBinaryPVTagged
	// IsBACnetConstructedDataBinaryValuePresentValue is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataBinaryValuePresentValue()
	// CreateBuilder creates a BACnetConstructedDataBinaryValuePresentValueBuilder
	CreateBACnetConstructedDataBinaryValuePresentValueBuilder() BACnetConstructedDataBinaryValuePresentValueBuilder
}

// _BACnetConstructedDataBinaryValuePresentValue is the data-structure of this message
type _BACnetConstructedDataBinaryValuePresentValue struct {
	BACnetConstructedDataContract
	PresentValue BACnetBinaryPVTagged
}

var _ BACnetConstructedDataBinaryValuePresentValue = (*_BACnetConstructedDataBinaryValuePresentValue)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataBinaryValuePresentValue)(nil)

// NewBACnetConstructedDataBinaryValuePresentValue factory function for _BACnetConstructedDataBinaryValuePresentValue
func NewBACnetConstructedDataBinaryValuePresentValue(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, presentValue BACnetBinaryPVTagged, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataBinaryValuePresentValue {
	if presentValue == nil {
		panic("presentValue of type BACnetBinaryPVTagged for BACnetConstructedDataBinaryValuePresentValue must not be nil")
	}
	_result := &_BACnetConstructedDataBinaryValuePresentValue{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		PresentValue:                  presentValue,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataBinaryValuePresentValueBuilder is a builder for BACnetConstructedDataBinaryValuePresentValue
type BACnetConstructedDataBinaryValuePresentValueBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(presentValue BACnetBinaryPVTagged) BACnetConstructedDataBinaryValuePresentValueBuilder
	// WithPresentValue adds PresentValue (property field)
	WithPresentValue(BACnetBinaryPVTagged) BACnetConstructedDataBinaryValuePresentValueBuilder
	// WithPresentValueBuilder adds PresentValue (property field) which is build by the builder
	WithPresentValueBuilder(func(BACnetBinaryPVTaggedBuilder) BACnetBinaryPVTaggedBuilder) BACnetConstructedDataBinaryValuePresentValueBuilder
	// Build builds the BACnetConstructedDataBinaryValuePresentValue or returns an error if something is wrong
	Build() (BACnetConstructedDataBinaryValuePresentValue, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataBinaryValuePresentValue
}

// NewBACnetConstructedDataBinaryValuePresentValueBuilder() creates a BACnetConstructedDataBinaryValuePresentValueBuilder
func NewBACnetConstructedDataBinaryValuePresentValueBuilder() BACnetConstructedDataBinaryValuePresentValueBuilder {
	return &_BACnetConstructedDataBinaryValuePresentValueBuilder{_BACnetConstructedDataBinaryValuePresentValue: new(_BACnetConstructedDataBinaryValuePresentValue)}
}

type _BACnetConstructedDataBinaryValuePresentValueBuilder struct {
	*_BACnetConstructedDataBinaryValuePresentValue

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataBinaryValuePresentValueBuilder) = (*_BACnetConstructedDataBinaryValuePresentValueBuilder)(nil)

func (b *_BACnetConstructedDataBinaryValuePresentValueBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataBinaryValuePresentValueBuilder) WithMandatoryFields(presentValue BACnetBinaryPVTagged) BACnetConstructedDataBinaryValuePresentValueBuilder {
	return b.WithPresentValue(presentValue)
}

func (b *_BACnetConstructedDataBinaryValuePresentValueBuilder) WithPresentValue(presentValue BACnetBinaryPVTagged) BACnetConstructedDataBinaryValuePresentValueBuilder {
	b.PresentValue = presentValue
	return b
}

func (b *_BACnetConstructedDataBinaryValuePresentValueBuilder) WithPresentValueBuilder(builderSupplier func(BACnetBinaryPVTaggedBuilder) BACnetBinaryPVTaggedBuilder) BACnetConstructedDataBinaryValuePresentValueBuilder {
	builder := builderSupplier(b.PresentValue.CreateBACnetBinaryPVTaggedBuilder())
	var err error
	b.PresentValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetBinaryPVTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataBinaryValuePresentValueBuilder) Build() (BACnetConstructedDataBinaryValuePresentValue, error) {
	if b.PresentValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'presentValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataBinaryValuePresentValue.deepCopy(), nil
}

func (b *_BACnetConstructedDataBinaryValuePresentValueBuilder) MustBuild() BACnetConstructedDataBinaryValuePresentValue {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataBinaryValuePresentValueBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataBinaryValuePresentValueBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataBinaryValuePresentValueBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataBinaryValuePresentValueBuilder().(*_BACnetConstructedDataBinaryValuePresentValueBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataBinaryValuePresentValueBuilder creates a BACnetConstructedDataBinaryValuePresentValueBuilder
func (b *_BACnetConstructedDataBinaryValuePresentValue) CreateBACnetConstructedDataBinaryValuePresentValueBuilder() BACnetConstructedDataBinaryValuePresentValueBuilder {
	if b == nil {
		return NewBACnetConstructedDataBinaryValuePresentValueBuilder()
	}
	return &_BACnetConstructedDataBinaryValuePresentValueBuilder{_BACnetConstructedDataBinaryValuePresentValue: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataBinaryValuePresentValue) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_BINARY_VALUE
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_PRESENT_VALUE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataBinaryValuePresentValue) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataBinaryValuePresentValue) GetPresentValue() BACnetBinaryPVTagged {
	return m.PresentValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataBinaryValuePresentValue) GetActualValue() BACnetBinaryPVTagged {
	ctx := context.Background()
	_ = ctx
	return CastBACnetBinaryPVTagged(m.GetPresentValue())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataBinaryValuePresentValue(structType any) BACnetConstructedDataBinaryValuePresentValue {
	if casted, ok := structType.(BACnetConstructedDataBinaryValuePresentValue); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataBinaryValuePresentValue); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) GetTypeName() string {
	return "BACnetConstructedDataBinaryValuePresentValue"
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (presentValue)
	lengthInBits += m.PresentValue.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataBinaryValuePresentValue BACnetConstructedDataBinaryValuePresentValue, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataBinaryValuePresentValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataBinaryValuePresentValue")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	presentValue, err := ReadSimpleField[BACnetBinaryPVTagged](ctx, "presentValue", ReadComplex[BACnetBinaryPVTagged](BACnetBinaryPVTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'presentValue' field"))
	}
	m.PresentValue = presentValue

	actualValue, err := ReadVirtualField[BACnetBinaryPVTagged](ctx, "actualValue", (*BACnetBinaryPVTagged)(nil), presentValue)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataBinaryValuePresentValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataBinaryValuePresentValue")
	}

	return m, nil
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataBinaryValuePresentValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataBinaryValuePresentValue")
		}

		if err := WriteSimpleField[BACnetBinaryPVTagged](ctx, "presentValue", m.GetPresentValue(), WriteComplex[BACnetBinaryPVTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'presentValue' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataBinaryValuePresentValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataBinaryValuePresentValue")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) IsBACnetConstructedDataBinaryValuePresentValue() {
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) deepCopy() *_BACnetConstructedDataBinaryValuePresentValue {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataBinaryValuePresentValueCopy := &_BACnetConstructedDataBinaryValuePresentValue{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.PresentValue.DeepCopy().(BACnetBinaryPVTagged),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataBinaryValuePresentValueCopy
}

func (m *_BACnetConstructedDataBinaryValuePresentValue) String() string {
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
