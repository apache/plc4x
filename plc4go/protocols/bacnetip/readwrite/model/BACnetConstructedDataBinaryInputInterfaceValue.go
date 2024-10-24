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

// BACnetConstructedDataBinaryInputInterfaceValue is the corresponding interface of BACnetConstructedDataBinaryInputInterfaceValue
type BACnetConstructedDataBinaryInputInterfaceValue interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetInterfaceValue returns InterfaceValue (property field)
	GetInterfaceValue() BACnetOptionalBinaryPV
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetOptionalBinaryPV
	// IsBACnetConstructedDataBinaryInputInterfaceValue is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataBinaryInputInterfaceValue()
	// CreateBuilder creates a BACnetConstructedDataBinaryInputInterfaceValueBuilder
	CreateBACnetConstructedDataBinaryInputInterfaceValueBuilder() BACnetConstructedDataBinaryInputInterfaceValueBuilder
}

// _BACnetConstructedDataBinaryInputInterfaceValue is the data-structure of this message
type _BACnetConstructedDataBinaryInputInterfaceValue struct {
	BACnetConstructedDataContract
	InterfaceValue BACnetOptionalBinaryPV
}

var _ BACnetConstructedDataBinaryInputInterfaceValue = (*_BACnetConstructedDataBinaryInputInterfaceValue)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataBinaryInputInterfaceValue)(nil)

// NewBACnetConstructedDataBinaryInputInterfaceValue factory function for _BACnetConstructedDataBinaryInputInterfaceValue
func NewBACnetConstructedDataBinaryInputInterfaceValue(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, interfaceValue BACnetOptionalBinaryPV, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataBinaryInputInterfaceValue {
	if interfaceValue == nil {
		panic("interfaceValue of type BACnetOptionalBinaryPV for BACnetConstructedDataBinaryInputInterfaceValue must not be nil")
	}
	_result := &_BACnetConstructedDataBinaryInputInterfaceValue{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		InterfaceValue:                interfaceValue,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataBinaryInputInterfaceValueBuilder is a builder for BACnetConstructedDataBinaryInputInterfaceValue
type BACnetConstructedDataBinaryInputInterfaceValueBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(interfaceValue BACnetOptionalBinaryPV) BACnetConstructedDataBinaryInputInterfaceValueBuilder
	// WithInterfaceValue adds InterfaceValue (property field)
	WithInterfaceValue(BACnetOptionalBinaryPV) BACnetConstructedDataBinaryInputInterfaceValueBuilder
	// WithInterfaceValueBuilder adds InterfaceValue (property field) which is build by the builder
	WithInterfaceValueBuilder(func(BACnetOptionalBinaryPVBuilder) BACnetOptionalBinaryPVBuilder) BACnetConstructedDataBinaryInputInterfaceValueBuilder
	// Build builds the BACnetConstructedDataBinaryInputInterfaceValue or returns an error if something is wrong
	Build() (BACnetConstructedDataBinaryInputInterfaceValue, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataBinaryInputInterfaceValue
}

// NewBACnetConstructedDataBinaryInputInterfaceValueBuilder() creates a BACnetConstructedDataBinaryInputInterfaceValueBuilder
func NewBACnetConstructedDataBinaryInputInterfaceValueBuilder() BACnetConstructedDataBinaryInputInterfaceValueBuilder {
	return &_BACnetConstructedDataBinaryInputInterfaceValueBuilder{_BACnetConstructedDataBinaryInputInterfaceValue: new(_BACnetConstructedDataBinaryInputInterfaceValue)}
}

type _BACnetConstructedDataBinaryInputInterfaceValueBuilder struct {
	*_BACnetConstructedDataBinaryInputInterfaceValue

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataBinaryInputInterfaceValueBuilder) = (*_BACnetConstructedDataBinaryInputInterfaceValueBuilder)(nil)

func (b *_BACnetConstructedDataBinaryInputInterfaceValueBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataBinaryInputInterfaceValueBuilder) WithMandatoryFields(interfaceValue BACnetOptionalBinaryPV) BACnetConstructedDataBinaryInputInterfaceValueBuilder {
	return b.WithInterfaceValue(interfaceValue)
}

func (b *_BACnetConstructedDataBinaryInputInterfaceValueBuilder) WithInterfaceValue(interfaceValue BACnetOptionalBinaryPV) BACnetConstructedDataBinaryInputInterfaceValueBuilder {
	b.InterfaceValue = interfaceValue
	return b
}

func (b *_BACnetConstructedDataBinaryInputInterfaceValueBuilder) WithInterfaceValueBuilder(builderSupplier func(BACnetOptionalBinaryPVBuilder) BACnetOptionalBinaryPVBuilder) BACnetConstructedDataBinaryInputInterfaceValueBuilder {
	builder := builderSupplier(b.InterfaceValue.CreateBACnetOptionalBinaryPVBuilder())
	var err error
	b.InterfaceValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetOptionalBinaryPVBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataBinaryInputInterfaceValueBuilder) Build() (BACnetConstructedDataBinaryInputInterfaceValue, error) {
	if b.InterfaceValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'interfaceValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataBinaryInputInterfaceValue.deepCopy(), nil
}

func (b *_BACnetConstructedDataBinaryInputInterfaceValueBuilder) MustBuild() BACnetConstructedDataBinaryInputInterfaceValue {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataBinaryInputInterfaceValueBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataBinaryInputInterfaceValueBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataBinaryInputInterfaceValueBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataBinaryInputInterfaceValueBuilder().(*_BACnetConstructedDataBinaryInputInterfaceValueBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataBinaryInputInterfaceValueBuilder creates a BACnetConstructedDataBinaryInputInterfaceValueBuilder
func (b *_BACnetConstructedDataBinaryInputInterfaceValue) CreateBACnetConstructedDataBinaryInputInterfaceValueBuilder() BACnetConstructedDataBinaryInputInterfaceValueBuilder {
	if b == nil {
		return NewBACnetConstructedDataBinaryInputInterfaceValueBuilder()
	}
	return &_BACnetConstructedDataBinaryInputInterfaceValueBuilder{_BACnetConstructedDataBinaryInputInterfaceValue: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_BINARY_INPUT
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_INTERFACE_VALUE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) GetInterfaceValue() BACnetOptionalBinaryPV {
	return m.InterfaceValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) GetActualValue() BACnetOptionalBinaryPV {
	ctx := context.Background()
	_ = ctx
	return CastBACnetOptionalBinaryPV(m.GetInterfaceValue())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataBinaryInputInterfaceValue(structType any) BACnetConstructedDataBinaryInputInterfaceValue {
	if casted, ok := structType.(BACnetConstructedDataBinaryInputInterfaceValue); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataBinaryInputInterfaceValue); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) GetTypeName() string {
	return "BACnetConstructedDataBinaryInputInterfaceValue"
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (interfaceValue)
	lengthInBits += m.InterfaceValue.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataBinaryInputInterfaceValue BACnetConstructedDataBinaryInputInterfaceValue, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataBinaryInputInterfaceValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataBinaryInputInterfaceValue")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	interfaceValue, err := ReadSimpleField[BACnetOptionalBinaryPV](ctx, "interfaceValue", ReadComplex[BACnetOptionalBinaryPV](BACnetOptionalBinaryPVParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'interfaceValue' field"))
	}
	m.InterfaceValue = interfaceValue

	actualValue, err := ReadVirtualField[BACnetOptionalBinaryPV](ctx, "actualValue", (*BACnetOptionalBinaryPV)(nil), interfaceValue)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataBinaryInputInterfaceValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataBinaryInputInterfaceValue")
	}

	return m, nil
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataBinaryInputInterfaceValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataBinaryInputInterfaceValue")
		}

		if err := WriteSimpleField[BACnetOptionalBinaryPV](ctx, "interfaceValue", m.GetInterfaceValue(), WriteComplex[BACnetOptionalBinaryPV](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'interfaceValue' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataBinaryInputInterfaceValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataBinaryInputInterfaceValue")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) IsBACnetConstructedDataBinaryInputInterfaceValue() {
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) deepCopy() *_BACnetConstructedDataBinaryInputInterfaceValue {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataBinaryInputInterfaceValueCopy := &_BACnetConstructedDataBinaryInputInterfaceValue{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.InterfaceValue.DeepCopy().(BACnetOptionalBinaryPV),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataBinaryInputInterfaceValueCopy
}

func (m *_BACnetConstructedDataBinaryInputInterfaceValue) String() string {
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
