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

// BACnetConstructedDataUserInformationReference is the corresponding interface of BACnetConstructedDataUserInformationReference
type BACnetConstructedDataUserInformationReference interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetUserInformationReference returns UserInformationReference (property field)
	GetUserInformationReference() BACnetApplicationTagCharacterString
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagCharacterString
	// IsBACnetConstructedDataUserInformationReference is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataUserInformationReference()
	// CreateBuilder creates a BACnetConstructedDataUserInformationReferenceBuilder
	CreateBACnetConstructedDataUserInformationReferenceBuilder() BACnetConstructedDataUserInformationReferenceBuilder
}

// _BACnetConstructedDataUserInformationReference is the data-structure of this message
type _BACnetConstructedDataUserInformationReference struct {
	BACnetConstructedDataContract
	UserInformationReference BACnetApplicationTagCharacterString
}

var _ BACnetConstructedDataUserInformationReference = (*_BACnetConstructedDataUserInformationReference)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataUserInformationReference)(nil)

// NewBACnetConstructedDataUserInformationReference factory function for _BACnetConstructedDataUserInformationReference
func NewBACnetConstructedDataUserInformationReference(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, userInformationReference BACnetApplicationTagCharacterString, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataUserInformationReference {
	if userInformationReference == nil {
		panic("userInformationReference of type BACnetApplicationTagCharacterString for BACnetConstructedDataUserInformationReference must not be nil")
	}
	_result := &_BACnetConstructedDataUserInformationReference{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		UserInformationReference:      userInformationReference,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataUserInformationReferenceBuilder is a builder for BACnetConstructedDataUserInformationReference
type BACnetConstructedDataUserInformationReferenceBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(userInformationReference BACnetApplicationTagCharacterString) BACnetConstructedDataUserInformationReferenceBuilder
	// WithUserInformationReference adds UserInformationReference (property field)
	WithUserInformationReference(BACnetApplicationTagCharacterString) BACnetConstructedDataUserInformationReferenceBuilder
	// WithUserInformationReferenceBuilder adds UserInformationReference (property field) which is build by the builder
	WithUserInformationReferenceBuilder(func(BACnetApplicationTagCharacterStringBuilder) BACnetApplicationTagCharacterStringBuilder) BACnetConstructedDataUserInformationReferenceBuilder
	// Build builds the BACnetConstructedDataUserInformationReference or returns an error if something is wrong
	Build() (BACnetConstructedDataUserInformationReference, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataUserInformationReference
}

// NewBACnetConstructedDataUserInformationReferenceBuilder() creates a BACnetConstructedDataUserInformationReferenceBuilder
func NewBACnetConstructedDataUserInformationReferenceBuilder() BACnetConstructedDataUserInformationReferenceBuilder {
	return &_BACnetConstructedDataUserInformationReferenceBuilder{_BACnetConstructedDataUserInformationReference: new(_BACnetConstructedDataUserInformationReference)}
}

type _BACnetConstructedDataUserInformationReferenceBuilder struct {
	*_BACnetConstructedDataUserInformationReference

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataUserInformationReferenceBuilder) = (*_BACnetConstructedDataUserInformationReferenceBuilder)(nil)

func (b *_BACnetConstructedDataUserInformationReferenceBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataUserInformationReferenceBuilder) WithMandatoryFields(userInformationReference BACnetApplicationTagCharacterString) BACnetConstructedDataUserInformationReferenceBuilder {
	return b.WithUserInformationReference(userInformationReference)
}

func (b *_BACnetConstructedDataUserInformationReferenceBuilder) WithUserInformationReference(userInformationReference BACnetApplicationTagCharacterString) BACnetConstructedDataUserInformationReferenceBuilder {
	b.UserInformationReference = userInformationReference
	return b
}

func (b *_BACnetConstructedDataUserInformationReferenceBuilder) WithUserInformationReferenceBuilder(builderSupplier func(BACnetApplicationTagCharacterStringBuilder) BACnetApplicationTagCharacterStringBuilder) BACnetConstructedDataUserInformationReferenceBuilder {
	builder := builderSupplier(b.UserInformationReference.CreateBACnetApplicationTagCharacterStringBuilder())
	var err error
	b.UserInformationReference, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagCharacterStringBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataUserInformationReferenceBuilder) Build() (BACnetConstructedDataUserInformationReference, error) {
	if b.UserInformationReference == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'userInformationReference' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataUserInformationReference.deepCopy(), nil
}

func (b *_BACnetConstructedDataUserInformationReferenceBuilder) MustBuild() BACnetConstructedDataUserInformationReference {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataUserInformationReferenceBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataUserInformationReferenceBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataUserInformationReferenceBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataUserInformationReferenceBuilder().(*_BACnetConstructedDataUserInformationReferenceBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataUserInformationReferenceBuilder creates a BACnetConstructedDataUserInformationReferenceBuilder
func (b *_BACnetConstructedDataUserInformationReference) CreateBACnetConstructedDataUserInformationReferenceBuilder() BACnetConstructedDataUserInformationReferenceBuilder {
	if b == nil {
		return NewBACnetConstructedDataUserInformationReferenceBuilder()
	}
	return &_BACnetConstructedDataUserInformationReferenceBuilder{_BACnetConstructedDataUserInformationReference: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataUserInformationReference) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataUserInformationReference) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_USER_INFORMATION_REFERENCE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataUserInformationReference) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataUserInformationReference) GetUserInformationReference() BACnetApplicationTagCharacterString {
	return m.UserInformationReference
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataUserInformationReference) GetActualValue() BACnetApplicationTagCharacterString {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagCharacterString(m.GetUserInformationReference())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataUserInformationReference(structType any) BACnetConstructedDataUserInformationReference {
	if casted, ok := structType.(BACnetConstructedDataUserInformationReference); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataUserInformationReference); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataUserInformationReference) GetTypeName() string {
	return "BACnetConstructedDataUserInformationReference"
}

func (m *_BACnetConstructedDataUserInformationReference) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (userInformationReference)
	lengthInBits += m.UserInformationReference.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataUserInformationReference) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataUserInformationReference) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataUserInformationReference BACnetConstructedDataUserInformationReference, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataUserInformationReference"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataUserInformationReference")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	userInformationReference, err := ReadSimpleField[BACnetApplicationTagCharacterString](ctx, "userInformationReference", ReadComplex[BACnetApplicationTagCharacterString](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagCharacterString](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'userInformationReference' field"))
	}
	m.UserInformationReference = userInformationReference

	actualValue, err := ReadVirtualField[BACnetApplicationTagCharacterString](ctx, "actualValue", (*BACnetApplicationTagCharacterString)(nil), userInformationReference)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataUserInformationReference"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataUserInformationReference")
	}

	return m, nil
}

func (m *_BACnetConstructedDataUserInformationReference) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataUserInformationReference) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataUserInformationReference"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataUserInformationReference")
		}

		if err := WriteSimpleField[BACnetApplicationTagCharacterString](ctx, "userInformationReference", m.GetUserInformationReference(), WriteComplex[BACnetApplicationTagCharacterString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'userInformationReference' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataUserInformationReference"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataUserInformationReference")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataUserInformationReference) IsBACnetConstructedDataUserInformationReference() {
}

func (m *_BACnetConstructedDataUserInformationReference) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataUserInformationReference) deepCopy() *_BACnetConstructedDataUserInformationReference {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataUserInformationReferenceCopy := &_BACnetConstructedDataUserInformationReference{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.UserInformationReference.DeepCopy().(BACnetApplicationTagCharacterString),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataUserInformationReferenceCopy
}

func (m *_BACnetConstructedDataUserInformationReference) String() string {
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
