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

// BACnetConstructedDataInputReference is the corresponding interface of BACnetConstructedDataInputReference
type BACnetConstructedDataInputReference interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetInputReference returns InputReference (property field)
	GetInputReference() BACnetObjectPropertyReference
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetObjectPropertyReference
	// IsBACnetConstructedDataInputReference is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataInputReference()
	// CreateBuilder creates a BACnetConstructedDataInputReferenceBuilder
	CreateBACnetConstructedDataInputReferenceBuilder() BACnetConstructedDataInputReferenceBuilder
}

// _BACnetConstructedDataInputReference is the data-structure of this message
type _BACnetConstructedDataInputReference struct {
	BACnetConstructedDataContract
	InputReference BACnetObjectPropertyReference
}

var _ BACnetConstructedDataInputReference = (*_BACnetConstructedDataInputReference)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataInputReference)(nil)

// NewBACnetConstructedDataInputReference factory function for _BACnetConstructedDataInputReference
func NewBACnetConstructedDataInputReference(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, inputReference BACnetObjectPropertyReference, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataInputReference {
	if inputReference == nil {
		panic("inputReference of type BACnetObjectPropertyReference for BACnetConstructedDataInputReference must not be nil")
	}
	_result := &_BACnetConstructedDataInputReference{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		InputReference:                inputReference,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataInputReferenceBuilder is a builder for BACnetConstructedDataInputReference
type BACnetConstructedDataInputReferenceBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(inputReference BACnetObjectPropertyReference) BACnetConstructedDataInputReferenceBuilder
	// WithInputReference adds InputReference (property field)
	WithInputReference(BACnetObjectPropertyReference) BACnetConstructedDataInputReferenceBuilder
	// WithInputReferenceBuilder adds InputReference (property field) which is build by the builder
	WithInputReferenceBuilder(func(BACnetObjectPropertyReferenceBuilder) BACnetObjectPropertyReferenceBuilder) BACnetConstructedDataInputReferenceBuilder
	// Build builds the BACnetConstructedDataInputReference or returns an error if something is wrong
	Build() (BACnetConstructedDataInputReference, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataInputReference
}

// NewBACnetConstructedDataInputReferenceBuilder() creates a BACnetConstructedDataInputReferenceBuilder
func NewBACnetConstructedDataInputReferenceBuilder() BACnetConstructedDataInputReferenceBuilder {
	return &_BACnetConstructedDataInputReferenceBuilder{_BACnetConstructedDataInputReference: new(_BACnetConstructedDataInputReference)}
}

type _BACnetConstructedDataInputReferenceBuilder struct {
	*_BACnetConstructedDataInputReference

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataInputReferenceBuilder) = (*_BACnetConstructedDataInputReferenceBuilder)(nil)

func (b *_BACnetConstructedDataInputReferenceBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataInputReferenceBuilder) WithMandatoryFields(inputReference BACnetObjectPropertyReference) BACnetConstructedDataInputReferenceBuilder {
	return b.WithInputReference(inputReference)
}

func (b *_BACnetConstructedDataInputReferenceBuilder) WithInputReference(inputReference BACnetObjectPropertyReference) BACnetConstructedDataInputReferenceBuilder {
	b.InputReference = inputReference
	return b
}

func (b *_BACnetConstructedDataInputReferenceBuilder) WithInputReferenceBuilder(builderSupplier func(BACnetObjectPropertyReferenceBuilder) BACnetObjectPropertyReferenceBuilder) BACnetConstructedDataInputReferenceBuilder {
	builder := builderSupplier(b.InputReference.CreateBACnetObjectPropertyReferenceBuilder())
	var err error
	b.InputReference, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetObjectPropertyReferenceBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataInputReferenceBuilder) Build() (BACnetConstructedDataInputReference, error) {
	if b.InputReference == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'inputReference' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataInputReference.deepCopy(), nil
}

func (b *_BACnetConstructedDataInputReferenceBuilder) MustBuild() BACnetConstructedDataInputReference {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataInputReferenceBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataInputReferenceBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataInputReferenceBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataInputReferenceBuilder().(*_BACnetConstructedDataInputReferenceBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataInputReferenceBuilder creates a BACnetConstructedDataInputReferenceBuilder
func (b *_BACnetConstructedDataInputReference) CreateBACnetConstructedDataInputReferenceBuilder() BACnetConstructedDataInputReferenceBuilder {
	if b == nil {
		return NewBACnetConstructedDataInputReferenceBuilder()
	}
	return &_BACnetConstructedDataInputReferenceBuilder{_BACnetConstructedDataInputReference: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataInputReference) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataInputReference) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_INPUT_REFERENCE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataInputReference) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataInputReference) GetInputReference() BACnetObjectPropertyReference {
	return m.InputReference
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataInputReference) GetActualValue() BACnetObjectPropertyReference {
	ctx := context.Background()
	_ = ctx
	return CastBACnetObjectPropertyReference(m.GetInputReference())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataInputReference(structType any) BACnetConstructedDataInputReference {
	if casted, ok := structType.(BACnetConstructedDataInputReference); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataInputReference); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataInputReference) GetTypeName() string {
	return "BACnetConstructedDataInputReference"
}

func (m *_BACnetConstructedDataInputReference) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (inputReference)
	lengthInBits += m.InputReference.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataInputReference) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataInputReference) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataInputReference BACnetConstructedDataInputReference, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataInputReference"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataInputReference")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	inputReference, err := ReadSimpleField[BACnetObjectPropertyReference](ctx, "inputReference", ReadComplex[BACnetObjectPropertyReference](BACnetObjectPropertyReferenceParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'inputReference' field"))
	}
	m.InputReference = inputReference

	actualValue, err := ReadVirtualField[BACnetObjectPropertyReference](ctx, "actualValue", (*BACnetObjectPropertyReference)(nil), inputReference)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataInputReference"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataInputReference")
	}

	return m, nil
}

func (m *_BACnetConstructedDataInputReference) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataInputReference) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataInputReference"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataInputReference")
		}

		if err := WriteSimpleField[BACnetObjectPropertyReference](ctx, "inputReference", m.GetInputReference(), WriteComplex[BACnetObjectPropertyReference](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'inputReference' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataInputReference"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataInputReference")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataInputReference) IsBACnetConstructedDataInputReference() {}

func (m *_BACnetConstructedDataInputReference) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataInputReference) deepCopy() *_BACnetConstructedDataInputReference {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataInputReferenceCopy := &_BACnetConstructedDataInputReference{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.InputReference.DeepCopy().(BACnetObjectPropertyReference),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataInputReferenceCopy
}

func (m *_BACnetConstructedDataInputReference) String() string {
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
