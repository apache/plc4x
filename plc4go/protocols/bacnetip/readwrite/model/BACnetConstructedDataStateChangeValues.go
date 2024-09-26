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

// BACnetConstructedDataStateChangeValues is the corresponding interface of BACnetConstructedDataStateChangeValues
type BACnetConstructedDataStateChangeValues interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetNumberOfDataElements returns NumberOfDataElements (property field)
	GetNumberOfDataElements() BACnetApplicationTagUnsignedInteger
	// GetStateChangeValues returns StateChangeValues (property field)
	GetStateChangeValues() []BACnetTimerStateChangeValue
	// GetZero returns Zero (virtual field)
	GetZero() uint64
	// IsBACnetConstructedDataStateChangeValues is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataStateChangeValues()
	// CreateBuilder creates a BACnetConstructedDataStateChangeValuesBuilder
	CreateBACnetConstructedDataStateChangeValuesBuilder() BACnetConstructedDataStateChangeValuesBuilder
}

// _BACnetConstructedDataStateChangeValues is the data-structure of this message
type _BACnetConstructedDataStateChangeValues struct {
	BACnetConstructedDataContract
	NumberOfDataElements BACnetApplicationTagUnsignedInteger
	StateChangeValues    []BACnetTimerStateChangeValue
}

var _ BACnetConstructedDataStateChangeValues = (*_BACnetConstructedDataStateChangeValues)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataStateChangeValues)(nil)

// NewBACnetConstructedDataStateChangeValues factory function for _BACnetConstructedDataStateChangeValues
func NewBACnetConstructedDataStateChangeValues(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, numberOfDataElements BACnetApplicationTagUnsignedInteger, stateChangeValues []BACnetTimerStateChangeValue, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataStateChangeValues {
	_result := &_BACnetConstructedDataStateChangeValues{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		NumberOfDataElements:          numberOfDataElements,
		StateChangeValues:             stateChangeValues,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataStateChangeValuesBuilder is a builder for BACnetConstructedDataStateChangeValues
type BACnetConstructedDataStateChangeValuesBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(stateChangeValues []BACnetTimerStateChangeValue) BACnetConstructedDataStateChangeValuesBuilder
	// WithNumberOfDataElements adds NumberOfDataElements (property field)
	WithOptionalNumberOfDataElements(BACnetApplicationTagUnsignedInteger) BACnetConstructedDataStateChangeValuesBuilder
	// WithOptionalNumberOfDataElementsBuilder adds NumberOfDataElements (property field) which is build by the builder
	WithOptionalNumberOfDataElementsBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataStateChangeValuesBuilder
	// WithStateChangeValues adds StateChangeValues (property field)
	WithStateChangeValues(...BACnetTimerStateChangeValue) BACnetConstructedDataStateChangeValuesBuilder
	// Build builds the BACnetConstructedDataStateChangeValues or returns an error if something is wrong
	Build() (BACnetConstructedDataStateChangeValues, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataStateChangeValues
}

// NewBACnetConstructedDataStateChangeValuesBuilder() creates a BACnetConstructedDataStateChangeValuesBuilder
func NewBACnetConstructedDataStateChangeValuesBuilder() BACnetConstructedDataStateChangeValuesBuilder {
	return &_BACnetConstructedDataStateChangeValuesBuilder{_BACnetConstructedDataStateChangeValues: new(_BACnetConstructedDataStateChangeValues)}
}

type _BACnetConstructedDataStateChangeValuesBuilder struct {
	*_BACnetConstructedDataStateChangeValues

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataStateChangeValuesBuilder) = (*_BACnetConstructedDataStateChangeValuesBuilder)(nil)

func (b *_BACnetConstructedDataStateChangeValuesBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataStateChangeValuesBuilder) WithMandatoryFields(stateChangeValues []BACnetTimerStateChangeValue) BACnetConstructedDataStateChangeValuesBuilder {
	return b.WithStateChangeValues(stateChangeValues...)
}

func (b *_BACnetConstructedDataStateChangeValuesBuilder) WithOptionalNumberOfDataElements(numberOfDataElements BACnetApplicationTagUnsignedInteger) BACnetConstructedDataStateChangeValuesBuilder {
	b.NumberOfDataElements = numberOfDataElements
	return b
}

func (b *_BACnetConstructedDataStateChangeValuesBuilder) WithOptionalNumberOfDataElementsBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataStateChangeValuesBuilder {
	builder := builderSupplier(b.NumberOfDataElements.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.NumberOfDataElements, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataStateChangeValuesBuilder) WithStateChangeValues(stateChangeValues ...BACnetTimerStateChangeValue) BACnetConstructedDataStateChangeValuesBuilder {
	b.StateChangeValues = stateChangeValues
	return b
}

func (b *_BACnetConstructedDataStateChangeValuesBuilder) Build() (BACnetConstructedDataStateChangeValues, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataStateChangeValues.deepCopy(), nil
}

func (b *_BACnetConstructedDataStateChangeValuesBuilder) MustBuild() BACnetConstructedDataStateChangeValues {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataStateChangeValuesBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataStateChangeValuesBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataStateChangeValuesBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataStateChangeValuesBuilder().(*_BACnetConstructedDataStateChangeValuesBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataStateChangeValuesBuilder creates a BACnetConstructedDataStateChangeValuesBuilder
func (b *_BACnetConstructedDataStateChangeValues) CreateBACnetConstructedDataStateChangeValuesBuilder() BACnetConstructedDataStateChangeValuesBuilder {
	if b == nil {
		return NewBACnetConstructedDataStateChangeValuesBuilder()
	}
	return &_BACnetConstructedDataStateChangeValuesBuilder{_BACnetConstructedDataStateChangeValues: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataStateChangeValues) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataStateChangeValues) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_STATE_CHANGE_VALUES
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataStateChangeValues) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataStateChangeValues) GetNumberOfDataElements() BACnetApplicationTagUnsignedInteger {
	return m.NumberOfDataElements
}

func (m *_BACnetConstructedDataStateChangeValues) GetStateChangeValues() []BACnetTimerStateChangeValue {
	return m.StateChangeValues
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataStateChangeValues) GetZero() uint64 {
	ctx := context.Background()
	_ = ctx
	numberOfDataElements := m.GetNumberOfDataElements()
	_ = numberOfDataElements
	return uint64(uint64(0))
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataStateChangeValues(structType any) BACnetConstructedDataStateChangeValues {
	if casted, ok := structType.(BACnetConstructedDataStateChangeValues); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataStateChangeValues); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataStateChangeValues) GetTypeName() string {
	return "BACnetConstructedDataStateChangeValues"
}

func (m *_BACnetConstructedDataStateChangeValues) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// A virtual field doesn't have any in- or output.

	// Optional Field (numberOfDataElements)
	if m.NumberOfDataElements != nil {
		lengthInBits += m.NumberOfDataElements.GetLengthInBits(ctx)
	}

	// Array field
	if len(m.StateChangeValues) > 0 {
		for _, element := range m.StateChangeValues {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	return lengthInBits
}

func (m *_BACnetConstructedDataStateChangeValues) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataStateChangeValues) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataStateChangeValues BACnetConstructedDataStateChangeValues, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataStateChangeValues"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataStateChangeValues")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	zero, err := ReadVirtualField[uint64](ctx, "zero", (*uint64)(nil), uint64(0))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'zero' field"))
	}
	_ = zero

	var numberOfDataElements BACnetApplicationTagUnsignedInteger
	_numberOfDataElements, err := ReadOptionalField[BACnetApplicationTagUnsignedInteger](ctx, "numberOfDataElements", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer), bool(bool((arrayIndexArgument) != (nil))) && bool(bool((arrayIndexArgument.GetActualValue()) == (zero))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'numberOfDataElements' field"))
	}
	if _numberOfDataElements != nil {
		numberOfDataElements = *_numberOfDataElements
		m.NumberOfDataElements = numberOfDataElements
	}

	stateChangeValues, err := ReadTerminatedArrayField[BACnetTimerStateChangeValue](ctx, "stateChangeValues", ReadComplex[BACnetTimerStateChangeValue](BACnetTimerStateChangeValueParseWithBufferProducer[BACnetTimerStateChangeValue]((BACnetObjectType)(objectTypeArgument)), readBuffer), IsBACnetConstructedDataClosingTag(ctx, readBuffer, false, tagNumber))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'stateChangeValues' field"))
	}
	m.StateChangeValues = stateChangeValues

	// Validation
	if !(bool(bool((arrayIndexArgument) != (nil))) || bool(bool((len(stateChangeValues)) == (7)))) {
		return nil, errors.WithStack(utils.ParseValidationError{Message: "stateChangeValues should have exactly 7 values"})
	}

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataStateChangeValues"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataStateChangeValues")
	}

	return m, nil
}

func (m *_BACnetConstructedDataStateChangeValues) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataStateChangeValues) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataStateChangeValues"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataStateChangeValues")
		}
		// Virtual field
		zero := m.GetZero()
		_ = zero
		if _zeroErr := writeBuffer.WriteVirtual(ctx, "zero", m.GetZero()); _zeroErr != nil {
			return errors.Wrap(_zeroErr, "Error serializing 'zero' field")
		}

		if err := WriteOptionalField[BACnetApplicationTagUnsignedInteger](ctx, "numberOfDataElements", GetRef(m.GetNumberOfDataElements()), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer), true); err != nil {
			return errors.Wrap(err, "Error serializing 'numberOfDataElements' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "stateChangeValues", m.GetStateChangeValues(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'stateChangeValues' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataStateChangeValues"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataStateChangeValues")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataStateChangeValues) IsBACnetConstructedDataStateChangeValues() {}

func (m *_BACnetConstructedDataStateChangeValues) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataStateChangeValues) deepCopy() *_BACnetConstructedDataStateChangeValues {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataStateChangeValuesCopy := &_BACnetConstructedDataStateChangeValues{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.NumberOfDataElements.DeepCopy().(BACnetApplicationTagUnsignedInteger),
		utils.DeepCopySlice[BACnetTimerStateChangeValue, BACnetTimerStateChangeValue](m.StateChangeValues),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataStateChangeValuesCopy
}

func (m *_BACnetConstructedDataStateChangeValues) String() string {
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
