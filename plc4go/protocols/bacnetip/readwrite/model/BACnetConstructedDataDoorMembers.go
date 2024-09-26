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

// BACnetConstructedDataDoorMembers is the corresponding interface of BACnetConstructedDataDoorMembers
type BACnetConstructedDataDoorMembers interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetNumberOfDataElements returns NumberOfDataElements (property field)
	GetNumberOfDataElements() BACnetApplicationTagUnsignedInteger
	// GetDoorMembers returns DoorMembers (property field)
	GetDoorMembers() []BACnetDeviceObjectReference
	// GetZero returns Zero (virtual field)
	GetZero() uint64
	// IsBACnetConstructedDataDoorMembers is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataDoorMembers()
	// CreateBuilder creates a BACnetConstructedDataDoorMembersBuilder
	CreateBACnetConstructedDataDoorMembersBuilder() BACnetConstructedDataDoorMembersBuilder
}

// _BACnetConstructedDataDoorMembers is the data-structure of this message
type _BACnetConstructedDataDoorMembers struct {
	BACnetConstructedDataContract
	NumberOfDataElements BACnetApplicationTagUnsignedInteger
	DoorMembers          []BACnetDeviceObjectReference
}

var _ BACnetConstructedDataDoorMembers = (*_BACnetConstructedDataDoorMembers)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataDoorMembers)(nil)

// NewBACnetConstructedDataDoorMembers factory function for _BACnetConstructedDataDoorMembers
func NewBACnetConstructedDataDoorMembers(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, numberOfDataElements BACnetApplicationTagUnsignedInteger, doorMembers []BACnetDeviceObjectReference, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataDoorMembers {
	_result := &_BACnetConstructedDataDoorMembers{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		NumberOfDataElements:          numberOfDataElements,
		DoorMembers:                   doorMembers,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataDoorMembersBuilder is a builder for BACnetConstructedDataDoorMembers
type BACnetConstructedDataDoorMembersBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(doorMembers []BACnetDeviceObjectReference) BACnetConstructedDataDoorMembersBuilder
	// WithNumberOfDataElements adds NumberOfDataElements (property field)
	WithOptionalNumberOfDataElements(BACnetApplicationTagUnsignedInteger) BACnetConstructedDataDoorMembersBuilder
	// WithOptionalNumberOfDataElementsBuilder adds NumberOfDataElements (property field) which is build by the builder
	WithOptionalNumberOfDataElementsBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataDoorMembersBuilder
	// WithDoorMembers adds DoorMembers (property field)
	WithDoorMembers(...BACnetDeviceObjectReference) BACnetConstructedDataDoorMembersBuilder
	// Build builds the BACnetConstructedDataDoorMembers or returns an error if something is wrong
	Build() (BACnetConstructedDataDoorMembers, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataDoorMembers
}

// NewBACnetConstructedDataDoorMembersBuilder() creates a BACnetConstructedDataDoorMembersBuilder
func NewBACnetConstructedDataDoorMembersBuilder() BACnetConstructedDataDoorMembersBuilder {
	return &_BACnetConstructedDataDoorMembersBuilder{_BACnetConstructedDataDoorMembers: new(_BACnetConstructedDataDoorMembers)}
}

type _BACnetConstructedDataDoorMembersBuilder struct {
	*_BACnetConstructedDataDoorMembers

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataDoorMembersBuilder) = (*_BACnetConstructedDataDoorMembersBuilder)(nil)

func (b *_BACnetConstructedDataDoorMembersBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataDoorMembersBuilder) WithMandatoryFields(doorMembers []BACnetDeviceObjectReference) BACnetConstructedDataDoorMembersBuilder {
	return b.WithDoorMembers(doorMembers...)
}

func (b *_BACnetConstructedDataDoorMembersBuilder) WithOptionalNumberOfDataElements(numberOfDataElements BACnetApplicationTagUnsignedInteger) BACnetConstructedDataDoorMembersBuilder {
	b.NumberOfDataElements = numberOfDataElements
	return b
}

func (b *_BACnetConstructedDataDoorMembersBuilder) WithOptionalNumberOfDataElementsBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataDoorMembersBuilder {
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

func (b *_BACnetConstructedDataDoorMembersBuilder) WithDoorMembers(doorMembers ...BACnetDeviceObjectReference) BACnetConstructedDataDoorMembersBuilder {
	b.DoorMembers = doorMembers
	return b
}

func (b *_BACnetConstructedDataDoorMembersBuilder) Build() (BACnetConstructedDataDoorMembers, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataDoorMembers.deepCopy(), nil
}

func (b *_BACnetConstructedDataDoorMembersBuilder) MustBuild() BACnetConstructedDataDoorMembers {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataDoorMembersBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataDoorMembersBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataDoorMembersBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataDoorMembersBuilder().(*_BACnetConstructedDataDoorMembersBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataDoorMembersBuilder creates a BACnetConstructedDataDoorMembersBuilder
func (b *_BACnetConstructedDataDoorMembers) CreateBACnetConstructedDataDoorMembersBuilder() BACnetConstructedDataDoorMembersBuilder {
	if b == nil {
		return NewBACnetConstructedDataDoorMembersBuilder()
	}
	return &_BACnetConstructedDataDoorMembersBuilder{_BACnetConstructedDataDoorMembers: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataDoorMembers) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataDoorMembers) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_DOOR_MEMBERS
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataDoorMembers) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataDoorMembers) GetNumberOfDataElements() BACnetApplicationTagUnsignedInteger {
	return m.NumberOfDataElements
}

func (m *_BACnetConstructedDataDoorMembers) GetDoorMembers() []BACnetDeviceObjectReference {
	return m.DoorMembers
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataDoorMembers) GetZero() uint64 {
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
func CastBACnetConstructedDataDoorMembers(structType any) BACnetConstructedDataDoorMembers {
	if casted, ok := structType.(BACnetConstructedDataDoorMembers); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataDoorMembers); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataDoorMembers) GetTypeName() string {
	return "BACnetConstructedDataDoorMembers"
}

func (m *_BACnetConstructedDataDoorMembers) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// A virtual field doesn't have any in- or output.

	// Optional Field (numberOfDataElements)
	if m.NumberOfDataElements != nil {
		lengthInBits += m.NumberOfDataElements.GetLengthInBits(ctx)
	}

	// Array field
	if len(m.DoorMembers) > 0 {
		for _, element := range m.DoorMembers {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	return lengthInBits
}

func (m *_BACnetConstructedDataDoorMembers) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataDoorMembers) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataDoorMembers BACnetConstructedDataDoorMembers, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataDoorMembers"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataDoorMembers")
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

	doorMembers, err := ReadTerminatedArrayField[BACnetDeviceObjectReference](ctx, "doorMembers", ReadComplex[BACnetDeviceObjectReference](BACnetDeviceObjectReferenceParseWithBuffer, readBuffer), IsBACnetConstructedDataClosingTag(ctx, readBuffer, false, tagNumber))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'doorMembers' field"))
	}
	m.DoorMembers = doorMembers

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataDoorMembers"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataDoorMembers")
	}

	return m, nil
}

func (m *_BACnetConstructedDataDoorMembers) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataDoorMembers) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataDoorMembers"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataDoorMembers")
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

		if err := WriteComplexTypeArrayField(ctx, "doorMembers", m.GetDoorMembers(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'doorMembers' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataDoorMembers"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataDoorMembers")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataDoorMembers) IsBACnetConstructedDataDoorMembers() {}

func (m *_BACnetConstructedDataDoorMembers) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataDoorMembers) deepCopy() *_BACnetConstructedDataDoorMembers {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataDoorMembersCopy := &_BACnetConstructedDataDoorMembers{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.NumberOfDataElements.DeepCopy().(BACnetApplicationTagUnsignedInteger),
		utils.DeepCopySlice[BACnetDeviceObjectReference, BACnetDeviceObjectReference](m.DoorMembers),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataDoorMembersCopy
}

func (m *_BACnetConstructedDataDoorMembers) String() string {
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
