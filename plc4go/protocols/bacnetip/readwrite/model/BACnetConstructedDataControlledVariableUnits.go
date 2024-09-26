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

// BACnetConstructedDataControlledVariableUnits is the corresponding interface of BACnetConstructedDataControlledVariableUnits
type BACnetConstructedDataControlledVariableUnits interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetUnits returns Units (property field)
	GetUnits() BACnetEngineeringUnitsTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetEngineeringUnitsTagged
	// IsBACnetConstructedDataControlledVariableUnits is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataControlledVariableUnits()
	// CreateBuilder creates a BACnetConstructedDataControlledVariableUnitsBuilder
	CreateBACnetConstructedDataControlledVariableUnitsBuilder() BACnetConstructedDataControlledVariableUnitsBuilder
}

// _BACnetConstructedDataControlledVariableUnits is the data-structure of this message
type _BACnetConstructedDataControlledVariableUnits struct {
	BACnetConstructedDataContract
	Units BACnetEngineeringUnitsTagged
}

var _ BACnetConstructedDataControlledVariableUnits = (*_BACnetConstructedDataControlledVariableUnits)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataControlledVariableUnits)(nil)

// NewBACnetConstructedDataControlledVariableUnits factory function for _BACnetConstructedDataControlledVariableUnits
func NewBACnetConstructedDataControlledVariableUnits(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, units BACnetEngineeringUnitsTagged, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataControlledVariableUnits {
	if units == nil {
		panic("units of type BACnetEngineeringUnitsTagged for BACnetConstructedDataControlledVariableUnits must not be nil")
	}
	_result := &_BACnetConstructedDataControlledVariableUnits{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		Units:                         units,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataControlledVariableUnitsBuilder is a builder for BACnetConstructedDataControlledVariableUnits
type BACnetConstructedDataControlledVariableUnitsBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(units BACnetEngineeringUnitsTagged) BACnetConstructedDataControlledVariableUnitsBuilder
	// WithUnits adds Units (property field)
	WithUnits(BACnetEngineeringUnitsTagged) BACnetConstructedDataControlledVariableUnitsBuilder
	// WithUnitsBuilder adds Units (property field) which is build by the builder
	WithUnitsBuilder(func(BACnetEngineeringUnitsTaggedBuilder) BACnetEngineeringUnitsTaggedBuilder) BACnetConstructedDataControlledVariableUnitsBuilder
	// Build builds the BACnetConstructedDataControlledVariableUnits or returns an error if something is wrong
	Build() (BACnetConstructedDataControlledVariableUnits, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataControlledVariableUnits
}

// NewBACnetConstructedDataControlledVariableUnitsBuilder() creates a BACnetConstructedDataControlledVariableUnitsBuilder
func NewBACnetConstructedDataControlledVariableUnitsBuilder() BACnetConstructedDataControlledVariableUnitsBuilder {
	return &_BACnetConstructedDataControlledVariableUnitsBuilder{_BACnetConstructedDataControlledVariableUnits: new(_BACnetConstructedDataControlledVariableUnits)}
}

type _BACnetConstructedDataControlledVariableUnitsBuilder struct {
	*_BACnetConstructedDataControlledVariableUnits

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataControlledVariableUnitsBuilder) = (*_BACnetConstructedDataControlledVariableUnitsBuilder)(nil)

func (b *_BACnetConstructedDataControlledVariableUnitsBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataControlledVariableUnitsBuilder) WithMandatoryFields(units BACnetEngineeringUnitsTagged) BACnetConstructedDataControlledVariableUnitsBuilder {
	return b.WithUnits(units)
}

func (b *_BACnetConstructedDataControlledVariableUnitsBuilder) WithUnits(units BACnetEngineeringUnitsTagged) BACnetConstructedDataControlledVariableUnitsBuilder {
	b.Units = units
	return b
}

func (b *_BACnetConstructedDataControlledVariableUnitsBuilder) WithUnitsBuilder(builderSupplier func(BACnetEngineeringUnitsTaggedBuilder) BACnetEngineeringUnitsTaggedBuilder) BACnetConstructedDataControlledVariableUnitsBuilder {
	builder := builderSupplier(b.Units.CreateBACnetEngineeringUnitsTaggedBuilder())
	var err error
	b.Units, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetEngineeringUnitsTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataControlledVariableUnitsBuilder) Build() (BACnetConstructedDataControlledVariableUnits, error) {
	if b.Units == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'units' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataControlledVariableUnits.deepCopy(), nil
}

func (b *_BACnetConstructedDataControlledVariableUnitsBuilder) MustBuild() BACnetConstructedDataControlledVariableUnits {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataControlledVariableUnitsBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataControlledVariableUnitsBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataControlledVariableUnitsBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataControlledVariableUnitsBuilder().(*_BACnetConstructedDataControlledVariableUnitsBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataControlledVariableUnitsBuilder creates a BACnetConstructedDataControlledVariableUnitsBuilder
func (b *_BACnetConstructedDataControlledVariableUnits) CreateBACnetConstructedDataControlledVariableUnitsBuilder() BACnetConstructedDataControlledVariableUnitsBuilder {
	if b == nil {
		return NewBACnetConstructedDataControlledVariableUnitsBuilder()
	}
	return &_BACnetConstructedDataControlledVariableUnitsBuilder{_BACnetConstructedDataControlledVariableUnits: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataControlledVariableUnits) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataControlledVariableUnits) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_CONTROLLED_VARIABLE_UNITS
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataControlledVariableUnits) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataControlledVariableUnits) GetUnits() BACnetEngineeringUnitsTagged {
	return m.Units
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataControlledVariableUnits) GetActualValue() BACnetEngineeringUnitsTagged {
	ctx := context.Background()
	_ = ctx
	return CastBACnetEngineeringUnitsTagged(m.GetUnits())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataControlledVariableUnits(structType any) BACnetConstructedDataControlledVariableUnits {
	if casted, ok := structType.(BACnetConstructedDataControlledVariableUnits); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataControlledVariableUnits); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataControlledVariableUnits) GetTypeName() string {
	return "BACnetConstructedDataControlledVariableUnits"
}

func (m *_BACnetConstructedDataControlledVariableUnits) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (units)
	lengthInBits += m.Units.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataControlledVariableUnits) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataControlledVariableUnits) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataControlledVariableUnits BACnetConstructedDataControlledVariableUnits, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataControlledVariableUnits"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataControlledVariableUnits")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	units, err := ReadSimpleField[BACnetEngineeringUnitsTagged](ctx, "units", ReadComplex[BACnetEngineeringUnitsTagged](BACnetEngineeringUnitsTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'units' field"))
	}
	m.Units = units

	actualValue, err := ReadVirtualField[BACnetEngineeringUnitsTagged](ctx, "actualValue", (*BACnetEngineeringUnitsTagged)(nil), units)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataControlledVariableUnits"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataControlledVariableUnits")
	}

	return m, nil
}

func (m *_BACnetConstructedDataControlledVariableUnits) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataControlledVariableUnits) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataControlledVariableUnits"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataControlledVariableUnits")
		}

		if err := WriteSimpleField[BACnetEngineeringUnitsTagged](ctx, "units", m.GetUnits(), WriteComplex[BACnetEngineeringUnitsTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'units' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataControlledVariableUnits"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataControlledVariableUnits")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataControlledVariableUnits) IsBACnetConstructedDataControlledVariableUnits() {
}

func (m *_BACnetConstructedDataControlledVariableUnits) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataControlledVariableUnits) deepCopy() *_BACnetConstructedDataControlledVariableUnits {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataControlledVariableUnitsCopy := &_BACnetConstructedDataControlledVariableUnits{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.Units.DeepCopy().(BACnetEngineeringUnitsTagged),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataControlledVariableUnitsCopy
}

func (m *_BACnetConstructedDataControlledVariableUnits) String() string {
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
