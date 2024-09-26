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

// BACnetConstructedDataCredentialDisable is the corresponding interface of BACnetConstructedDataCredentialDisable
type BACnetConstructedDataCredentialDisable interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetCredentialDisable returns CredentialDisable (property field)
	GetCredentialDisable() BACnetAccessCredentialDisableTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetAccessCredentialDisableTagged
	// IsBACnetConstructedDataCredentialDisable is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataCredentialDisable()
	// CreateBuilder creates a BACnetConstructedDataCredentialDisableBuilder
	CreateBACnetConstructedDataCredentialDisableBuilder() BACnetConstructedDataCredentialDisableBuilder
}

// _BACnetConstructedDataCredentialDisable is the data-structure of this message
type _BACnetConstructedDataCredentialDisable struct {
	BACnetConstructedDataContract
	CredentialDisable BACnetAccessCredentialDisableTagged
}

var _ BACnetConstructedDataCredentialDisable = (*_BACnetConstructedDataCredentialDisable)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataCredentialDisable)(nil)

// NewBACnetConstructedDataCredentialDisable factory function for _BACnetConstructedDataCredentialDisable
func NewBACnetConstructedDataCredentialDisable(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, credentialDisable BACnetAccessCredentialDisableTagged, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataCredentialDisable {
	if credentialDisable == nil {
		panic("credentialDisable of type BACnetAccessCredentialDisableTagged for BACnetConstructedDataCredentialDisable must not be nil")
	}
	_result := &_BACnetConstructedDataCredentialDisable{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		CredentialDisable:             credentialDisable,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataCredentialDisableBuilder is a builder for BACnetConstructedDataCredentialDisable
type BACnetConstructedDataCredentialDisableBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(credentialDisable BACnetAccessCredentialDisableTagged) BACnetConstructedDataCredentialDisableBuilder
	// WithCredentialDisable adds CredentialDisable (property field)
	WithCredentialDisable(BACnetAccessCredentialDisableTagged) BACnetConstructedDataCredentialDisableBuilder
	// WithCredentialDisableBuilder adds CredentialDisable (property field) which is build by the builder
	WithCredentialDisableBuilder(func(BACnetAccessCredentialDisableTaggedBuilder) BACnetAccessCredentialDisableTaggedBuilder) BACnetConstructedDataCredentialDisableBuilder
	// Build builds the BACnetConstructedDataCredentialDisable or returns an error if something is wrong
	Build() (BACnetConstructedDataCredentialDisable, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataCredentialDisable
}

// NewBACnetConstructedDataCredentialDisableBuilder() creates a BACnetConstructedDataCredentialDisableBuilder
func NewBACnetConstructedDataCredentialDisableBuilder() BACnetConstructedDataCredentialDisableBuilder {
	return &_BACnetConstructedDataCredentialDisableBuilder{_BACnetConstructedDataCredentialDisable: new(_BACnetConstructedDataCredentialDisable)}
}

type _BACnetConstructedDataCredentialDisableBuilder struct {
	*_BACnetConstructedDataCredentialDisable

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataCredentialDisableBuilder) = (*_BACnetConstructedDataCredentialDisableBuilder)(nil)

func (b *_BACnetConstructedDataCredentialDisableBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataCredentialDisableBuilder) WithMandatoryFields(credentialDisable BACnetAccessCredentialDisableTagged) BACnetConstructedDataCredentialDisableBuilder {
	return b.WithCredentialDisable(credentialDisable)
}

func (b *_BACnetConstructedDataCredentialDisableBuilder) WithCredentialDisable(credentialDisable BACnetAccessCredentialDisableTagged) BACnetConstructedDataCredentialDisableBuilder {
	b.CredentialDisable = credentialDisable
	return b
}

func (b *_BACnetConstructedDataCredentialDisableBuilder) WithCredentialDisableBuilder(builderSupplier func(BACnetAccessCredentialDisableTaggedBuilder) BACnetAccessCredentialDisableTaggedBuilder) BACnetConstructedDataCredentialDisableBuilder {
	builder := builderSupplier(b.CredentialDisable.CreateBACnetAccessCredentialDisableTaggedBuilder())
	var err error
	b.CredentialDisable, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetAccessCredentialDisableTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataCredentialDisableBuilder) Build() (BACnetConstructedDataCredentialDisable, error) {
	if b.CredentialDisable == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'credentialDisable' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataCredentialDisable.deepCopy(), nil
}

func (b *_BACnetConstructedDataCredentialDisableBuilder) MustBuild() BACnetConstructedDataCredentialDisable {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataCredentialDisableBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataCredentialDisableBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataCredentialDisableBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataCredentialDisableBuilder().(*_BACnetConstructedDataCredentialDisableBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataCredentialDisableBuilder creates a BACnetConstructedDataCredentialDisableBuilder
func (b *_BACnetConstructedDataCredentialDisable) CreateBACnetConstructedDataCredentialDisableBuilder() BACnetConstructedDataCredentialDisableBuilder {
	if b == nil {
		return NewBACnetConstructedDataCredentialDisableBuilder()
	}
	return &_BACnetConstructedDataCredentialDisableBuilder{_BACnetConstructedDataCredentialDisable: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataCredentialDisable) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataCredentialDisable) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_CREDENTIAL_DISABLE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataCredentialDisable) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataCredentialDisable) GetCredentialDisable() BACnetAccessCredentialDisableTagged {
	return m.CredentialDisable
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataCredentialDisable) GetActualValue() BACnetAccessCredentialDisableTagged {
	ctx := context.Background()
	_ = ctx
	return CastBACnetAccessCredentialDisableTagged(m.GetCredentialDisable())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataCredentialDisable(structType any) BACnetConstructedDataCredentialDisable {
	if casted, ok := structType.(BACnetConstructedDataCredentialDisable); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataCredentialDisable); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataCredentialDisable) GetTypeName() string {
	return "BACnetConstructedDataCredentialDisable"
}

func (m *_BACnetConstructedDataCredentialDisable) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (credentialDisable)
	lengthInBits += m.CredentialDisable.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataCredentialDisable) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataCredentialDisable) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataCredentialDisable BACnetConstructedDataCredentialDisable, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataCredentialDisable"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataCredentialDisable")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	credentialDisable, err := ReadSimpleField[BACnetAccessCredentialDisableTagged](ctx, "credentialDisable", ReadComplex[BACnetAccessCredentialDisableTagged](BACnetAccessCredentialDisableTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'credentialDisable' field"))
	}
	m.CredentialDisable = credentialDisable

	actualValue, err := ReadVirtualField[BACnetAccessCredentialDisableTagged](ctx, "actualValue", (*BACnetAccessCredentialDisableTagged)(nil), credentialDisable)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataCredentialDisable"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataCredentialDisable")
	}

	return m, nil
}

func (m *_BACnetConstructedDataCredentialDisable) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataCredentialDisable) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataCredentialDisable"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataCredentialDisable")
		}

		if err := WriteSimpleField[BACnetAccessCredentialDisableTagged](ctx, "credentialDisable", m.GetCredentialDisable(), WriteComplex[BACnetAccessCredentialDisableTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'credentialDisable' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataCredentialDisable"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataCredentialDisable")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataCredentialDisable) IsBACnetConstructedDataCredentialDisable() {}

func (m *_BACnetConstructedDataCredentialDisable) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataCredentialDisable) deepCopy() *_BACnetConstructedDataCredentialDisable {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataCredentialDisableCopy := &_BACnetConstructedDataCredentialDisable{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.CredentialDisable.DeepCopy().(BACnetAccessCredentialDisableTagged),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataCredentialDisableCopy
}

func (m *_BACnetConstructedDataCredentialDisable) String() string {
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
