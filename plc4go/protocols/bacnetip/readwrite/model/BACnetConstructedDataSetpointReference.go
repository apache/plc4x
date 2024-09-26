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

// BACnetConstructedDataSetpointReference is the corresponding interface of BACnetConstructedDataSetpointReference
type BACnetConstructedDataSetpointReference interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetSetpointReference returns SetpointReference (property field)
	GetSetpointReference() BACnetSetpointReference
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetSetpointReference
	// IsBACnetConstructedDataSetpointReference is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataSetpointReference()
	// CreateBuilder creates a BACnetConstructedDataSetpointReferenceBuilder
	CreateBACnetConstructedDataSetpointReferenceBuilder() BACnetConstructedDataSetpointReferenceBuilder
}

// _BACnetConstructedDataSetpointReference is the data-structure of this message
type _BACnetConstructedDataSetpointReference struct {
	BACnetConstructedDataContract
	SetpointReference BACnetSetpointReference
}

var _ BACnetConstructedDataSetpointReference = (*_BACnetConstructedDataSetpointReference)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataSetpointReference)(nil)

// NewBACnetConstructedDataSetpointReference factory function for _BACnetConstructedDataSetpointReference
func NewBACnetConstructedDataSetpointReference(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, setpointReference BACnetSetpointReference, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataSetpointReference {
	if setpointReference == nil {
		panic("setpointReference of type BACnetSetpointReference for BACnetConstructedDataSetpointReference must not be nil")
	}
	_result := &_BACnetConstructedDataSetpointReference{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		SetpointReference:             setpointReference,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataSetpointReferenceBuilder is a builder for BACnetConstructedDataSetpointReference
type BACnetConstructedDataSetpointReferenceBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(setpointReference BACnetSetpointReference) BACnetConstructedDataSetpointReferenceBuilder
	// WithSetpointReference adds SetpointReference (property field)
	WithSetpointReference(BACnetSetpointReference) BACnetConstructedDataSetpointReferenceBuilder
	// WithSetpointReferenceBuilder adds SetpointReference (property field) which is build by the builder
	WithSetpointReferenceBuilder(func(BACnetSetpointReferenceBuilder) BACnetSetpointReferenceBuilder) BACnetConstructedDataSetpointReferenceBuilder
	// Build builds the BACnetConstructedDataSetpointReference or returns an error if something is wrong
	Build() (BACnetConstructedDataSetpointReference, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataSetpointReference
}

// NewBACnetConstructedDataSetpointReferenceBuilder() creates a BACnetConstructedDataSetpointReferenceBuilder
func NewBACnetConstructedDataSetpointReferenceBuilder() BACnetConstructedDataSetpointReferenceBuilder {
	return &_BACnetConstructedDataSetpointReferenceBuilder{_BACnetConstructedDataSetpointReference: new(_BACnetConstructedDataSetpointReference)}
}

type _BACnetConstructedDataSetpointReferenceBuilder struct {
	*_BACnetConstructedDataSetpointReference

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataSetpointReferenceBuilder) = (*_BACnetConstructedDataSetpointReferenceBuilder)(nil)

func (b *_BACnetConstructedDataSetpointReferenceBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataSetpointReferenceBuilder) WithMandatoryFields(setpointReference BACnetSetpointReference) BACnetConstructedDataSetpointReferenceBuilder {
	return b.WithSetpointReference(setpointReference)
}

func (b *_BACnetConstructedDataSetpointReferenceBuilder) WithSetpointReference(setpointReference BACnetSetpointReference) BACnetConstructedDataSetpointReferenceBuilder {
	b.SetpointReference = setpointReference
	return b
}

func (b *_BACnetConstructedDataSetpointReferenceBuilder) WithSetpointReferenceBuilder(builderSupplier func(BACnetSetpointReferenceBuilder) BACnetSetpointReferenceBuilder) BACnetConstructedDataSetpointReferenceBuilder {
	builder := builderSupplier(b.SetpointReference.CreateBACnetSetpointReferenceBuilder())
	var err error
	b.SetpointReference, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetSetpointReferenceBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataSetpointReferenceBuilder) Build() (BACnetConstructedDataSetpointReference, error) {
	if b.SetpointReference == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'setpointReference' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataSetpointReference.deepCopy(), nil
}

func (b *_BACnetConstructedDataSetpointReferenceBuilder) MustBuild() BACnetConstructedDataSetpointReference {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataSetpointReferenceBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataSetpointReferenceBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataSetpointReferenceBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataSetpointReferenceBuilder().(*_BACnetConstructedDataSetpointReferenceBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataSetpointReferenceBuilder creates a BACnetConstructedDataSetpointReferenceBuilder
func (b *_BACnetConstructedDataSetpointReference) CreateBACnetConstructedDataSetpointReferenceBuilder() BACnetConstructedDataSetpointReferenceBuilder {
	if b == nil {
		return NewBACnetConstructedDataSetpointReferenceBuilder()
	}
	return &_BACnetConstructedDataSetpointReferenceBuilder{_BACnetConstructedDataSetpointReference: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataSetpointReference) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataSetpointReference) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_SETPOINT_REFERENCE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataSetpointReference) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataSetpointReference) GetSetpointReference() BACnetSetpointReference {
	return m.SetpointReference
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataSetpointReference) GetActualValue() BACnetSetpointReference {
	ctx := context.Background()
	_ = ctx
	return CastBACnetSetpointReference(m.GetSetpointReference())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataSetpointReference(structType any) BACnetConstructedDataSetpointReference {
	if casted, ok := structType.(BACnetConstructedDataSetpointReference); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataSetpointReference); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataSetpointReference) GetTypeName() string {
	return "BACnetConstructedDataSetpointReference"
}

func (m *_BACnetConstructedDataSetpointReference) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (setpointReference)
	lengthInBits += m.SetpointReference.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataSetpointReference) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataSetpointReference) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataSetpointReference BACnetConstructedDataSetpointReference, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataSetpointReference"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataSetpointReference")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	setpointReference, err := ReadSimpleField[BACnetSetpointReference](ctx, "setpointReference", ReadComplex[BACnetSetpointReference](BACnetSetpointReferenceParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'setpointReference' field"))
	}
	m.SetpointReference = setpointReference

	actualValue, err := ReadVirtualField[BACnetSetpointReference](ctx, "actualValue", (*BACnetSetpointReference)(nil), setpointReference)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataSetpointReference"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataSetpointReference")
	}

	return m, nil
}

func (m *_BACnetConstructedDataSetpointReference) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataSetpointReference) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataSetpointReference"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataSetpointReference")
		}

		if err := WriteSimpleField[BACnetSetpointReference](ctx, "setpointReference", m.GetSetpointReference(), WriteComplex[BACnetSetpointReference](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'setpointReference' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataSetpointReference"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataSetpointReference")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataSetpointReference) IsBACnetConstructedDataSetpointReference() {}

func (m *_BACnetConstructedDataSetpointReference) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataSetpointReference) deepCopy() *_BACnetConstructedDataSetpointReference {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataSetpointReferenceCopy := &_BACnetConstructedDataSetpointReference{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.SetpointReference.DeepCopy().(BACnetSetpointReference),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataSetpointReferenceCopy
}

func (m *_BACnetConstructedDataSetpointReference) String() string {
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
