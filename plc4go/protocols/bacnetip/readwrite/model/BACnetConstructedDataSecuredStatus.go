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

// BACnetConstructedDataSecuredStatus is the corresponding interface of BACnetConstructedDataSecuredStatus
type BACnetConstructedDataSecuredStatus interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetSecuredStatus returns SecuredStatus (property field)
	GetSecuredStatus() BACnetDoorSecuredStatusTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetDoorSecuredStatusTagged
	// IsBACnetConstructedDataSecuredStatus is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataSecuredStatus()
	// CreateBuilder creates a BACnetConstructedDataSecuredStatusBuilder
	CreateBACnetConstructedDataSecuredStatusBuilder() BACnetConstructedDataSecuredStatusBuilder
}

// _BACnetConstructedDataSecuredStatus is the data-structure of this message
type _BACnetConstructedDataSecuredStatus struct {
	BACnetConstructedDataContract
	SecuredStatus BACnetDoorSecuredStatusTagged
}

var _ BACnetConstructedDataSecuredStatus = (*_BACnetConstructedDataSecuredStatus)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataSecuredStatus)(nil)

// NewBACnetConstructedDataSecuredStatus factory function for _BACnetConstructedDataSecuredStatus
func NewBACnetConstructedDataSecuredStatus(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, securedStatus BACnetDoorSecuredStatusTagged, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataSecuredStatus {
	if securedStatus == nil {
		panic("securedStatus of type BACnetDoorSecuredStatusTagged for BACnetConstructedDataSecuredStatus must not be nil")
	}
	_result := &_BACnetConstructedDataSecuredStatus{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		SecuredStatus:                 securedStatus,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataSecuredStatusBuilder is a builder for BACnetConstructedDataSecuredStatus
type BACnetConstructedDataSecuredStatusBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(securedStatus BACnetDoorSecuredStatusTagged) BACnetConstructedDataSecuredStatusBuilder
	// WithSecuredStatus adds SecuredStatus (property field)
	WithSecuredStatus(BACnetDoorSecuredStatusTagged) BACnetConstructedDataSecuredStatusBuilder
	// WithSecuredStatusBuilder adds SecuredStatus (property field) which is build by the builder
	WithSecuredStatusBuilder(func(BACnetDoorSecuredStatusTaggedBuilder) BACnetDoorSecuredStatusTaggedBuilder) BACnetConstructedDataSecuredStatusBuilder
	// Build builds the BACnetConstructedDataSecuredStatus or returns an error if something is wrong
	Build() (BACnetConstructedDataSecuredStatus, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataSecuredStatus
}

// NewBACnetConstructedDataSecuredStatusBuilder() creates a BACnetConstructedDataSecuredStatusBuilder
func NewBACnetConstructedDataSecuredStatusBuilder() BACnetConstructedDataSecuredStatusBuilder {
	return &_BACnetConstructedDataSecuredStatusBuilder{_BACnetConstructedDataSecuredStatus: new(_BACnetConstructedDataSecuredStatus)}
}

type _BACnetConstructedDataSecuredStatusBuilder struct {
	*_BACnetConstructedDataSecuredStatus

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataSecuredStatusBuilder) = (*_BACnetConstructedDataSecuredStatusBuilder)(nil)

func (b *_BACnetConstructedDataSecuredStatusBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataSecuredStatusBuilder) WithMandatoryFields(securedStatus BACnetDoorSecuredStatusTagged) BACnetConstructedDataSecuredStatusBuilder {
	return b.WithSecuredStatus(securedStatus)
}

func (b *_BACnetConstructedDataSecuredStatusBuilder) WithSecuredStatus(securedStatus BACnetDoorSecuredStatusTagged) BACnetConstructedDataSecuredStatusBuilder {
	b.SecuredStatus = securedStatus
	return b
}

func (b *_BACnetConstructedDataSecuredStatusBuilder) WithSecuredStatusBuilder(builderSupplier func(BACnetDoorSecuredStatusTaggedBuilder) BACnetDoorSecuredStatusTaggedBuilder) BACnetConstructedDataSecuredStatusBuilder {
	builder := builderSupplier(b.SecuredStatus.CreateBACnetDoorSecuredStatusTaggedBuilder())
	var err error
	b.SecuredStatus, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetDoorSecuredStatusTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataSecuredStatusBuilder) Build() (BACnetConstructedDataSecuredStatus, error) {
	if b.SecuredStatus == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'securedStatus' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataSecuredStatus.deepCopy(), nil
}

func (b *_BACnetConstructedDataSecuredStatusBuilder) MustBuild() BACnetConstructedDataSecuredStatus {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataSecuredStatusBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataSecuredStatusBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataSecuredStatusBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataSecuredStatusBuilder().(*_BACnetConstructedDataSecuredStatusBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataSecuredStatusBuilder creates a BACnetConstructedDataSecuredStatusBuilder
func (b *_BACnetConstructedDataSecuredStatus) CreateBACnetConstructedDataSecuredStatusBuilder() BACnetConstructedDataSecuredStatusBuilder {
	if b == nil {
		return NewBACnetConstructedDataSecuredStatusBuilder()
	}
	return &_BACnetConstructedDataSecuredStatusBuilder{_BACnetConstructedDataSecuredStatus: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataSecuredStatus) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataSecuredStatus) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_SECURED_STATUS
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataSecuredStatus) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataSecuredStatus) GetSecuredStatus() BACnetDoorSecuredStatusTagged {
	return m.SecuredStatus
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataSecuredStatus) GetActualValue() BACnetDoorSecuredStatusTagged {
	ctx := context.Background()
	_ = ctx
	return CastBACnetDoorSecuredStatusTagged(m.GetSecuredStatus())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataSecuredStatus(structType any) BACnetConstructedDataSecuredStatus {
	if casted, ok := structType.(BACnetConstructedDataSecuredStatus); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataSecuredStatus); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataSecuredStatus) GetTypeName() string {
	return "BACnetConstructedDataSecuredStatus"
}

func (m *_BACnetConstructedDataSecuredStatus) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (securedStatus)
	lengthInBits += m.SecuredStatus.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataSecuredStatus) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataSecuredStatus) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataSecuredStatus BACnetConstructedDataSecuredStatus, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataSecuredStatus"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataSecuredStatus")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	securedStatus, err := ReadSimpleField[BACnetDoorSecuredStatusTagged](ctx, "securedStatus", ReadComplex[BACnetDoorSecuredStatusTagged](BACnetDoorSecuredStatusTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'securedStatus' field"))
	}
	m.SecuredStatus = securedStatus

	actualValue, err := ReadVirtualField[BACnetDoorSecuredStatusTagged](ctx, "actualValue", (*BACnetDoorSecuredStatusTagged)(nil), securedStatus)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataSecuredStatus"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataSecuredStatus")
	}

	return m, nil
}

func (m *_BACnetConstructedDataSecuredStatus) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataSecuredStatus) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataSecuredStatus"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataSecuredStatus")
		}

		if err := WriteSimpleField[BACnetDoorSecuredStatusTagged](ctx, "securedStatus", m.GetSecuredStatus(), WriteComplex[BACnetDoorSecuredStatusTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'securedStatus' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataSecuredStatus"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataSecuredStatus")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataSecuredStatus) IsBACnetConstructedDataSecuredStatus() {}

func (m *_BACnetConstructedDataSecuredStatus) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataSecuredStatus) deepCopy() *_BACnetConstructedDataSecuredStatus {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataSecuredStatusCopy := &_BACnetConstructedDataSecuredStatus{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.SecuredStatus.DeepCopy().(BACnetDoorSecuredStatusTagged),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataSecuredStatusCopy
}

func (m *_BACnetConstructedDataSecuredStatus) String() string {
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
