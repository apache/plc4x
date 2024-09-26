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

// BACnetConstructedDataDoorExtendedPulseTime is the corresponding interface of BACnetConstructedDataDoorExtendedPulseTime
type BACnetConstructedDataDoorExtendedPulseTime interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetDoorExtendedPulseTime returns DoorExtendedPulseTime (property field)
	GetDoorExtendedPulseTime() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
	// IsBACnetConstructedDataDoorExtendedPulseTime is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataDoorExtendedPulseTime()
	// CreateBuilder creates a BACnetConstructedDataDoorExtendedPulseTimeBuilder
	CreateBACnetConstructedDataDoorExtendedPulseTimeBuilder() BACnetConstructedDataDoorExtendedPulseTimeBuilder
}

// _BACnetConstructedDataDoorExtendedPulseTime is the data-structure of this message
type _BACnetConstructedDataDoorExtendedPulseTime struct {
	BACnetConstructedDataContract
	DoorExtendedPulseTime BACnetApplicationTagUnsignedInteger
}

var _ BACnetConstructedDataDoorExtendedPulseTime = (*_BACnetConstructedDataDoorExtendedPulseTime)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataDoorExtendedPulseTime)(nil)

// NewBACnetConstructedDataDoorExtendedPulseTime factory function for _BACnetConstructedDataDoorExtendedPulseTime
func NewBACnetConstructedDataDoorExtendedPulseTime(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, doorExtendedPulseTime BACnetApplicationTagUnsignedInteger, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataDoorExtendedPulseTime {
	if doorExtendedPulseTime == nil {
		panic("doorExtendedPulseTime of type BACnetApplicationTagUnsignedInteger for BACnetConstructedDataDoorExtendedPulseTime must not be nil")
	}
	_result := &_BACnetConstructedDataDoorExtendedPulseTime{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		DoorExtendedPulseTime:         doorExtendedPulseTime,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataDoorExtendedPulseTimeBuilder is a builder for BACnetConstructedDataDoorExtendedPulseTime
type BACnetConstructedDataDoorExtendedPulseTimeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(doorExtendedPulseTime BACnetApplicationTagUnsignedInteger) BACnetConstructedDataDoorExtendedPulseTimeBuilder
	// WithDoorExtendedPulseTime adds DoorExtendedPulseTime (property field)
	WithDoorExtendedPulseTime(BACnetApplicationTagUnsignedInteger) BACnetConstructedDataDoorExtendedPulseTimeBuilder
	// WithDoorExtendedPulseTimeBuilder adds DoorExtendedPulseTime (property field) which is build by the builder
	WithDoorExtendedPulseTimeBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataDoorExtendedPulseTimeBuilder
	// Build builds the BACnetConstructedDataDoorExtendedPulseTime or returns an error if something is wrong
	Build() (BACnetConstructedDataDoorExtendedPulseTime, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataDoorExtendedPulseTime
}

// NewBACnetConstructedDataDoorExtendedPulseTimeBuilder() creates a BACnetConstructedDataDoorExtendedPulseTimeBuilder
func NewBACnetConstructedDataDoorExtendedPulseTimeBuilder() BACnetConstructedDataDoorExtendedPulseTimeBuilder {
	return &_BACnetConstructedDataDoorExtendedPulseTimeBuilder{_BACnetConstructedDataDoorExtendedPulseTime: new(_BACnetConstructedDataDoorExtendedPulseTime)}
}

type _BACnetConstructedDataDoorExtendedPulseTimeBuilder struct {
	*_BACnetConstructedDataDoorExtendedPulseTime

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataDoorExtendedPulseTimeBuilder) = (*_BACnetConstructedDataDoorExtendedPulseTimeBuilder)(nil)

func (b *_BACnetConstructedDataDoorExtendedPulseTimeBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataDoorExtendedPulseTimeBuilder) WithMandatoryFields(doorExtendedPulseTime BACnetApplicationTagUnsignedInteger) BACnetConstructedDataDoorExtendedPulseTimeBuilder {
	return b.WithDoorExtendedPulseTime(doorExtendedPulseTime)
}

func (b *_BACnetConstructedDataDoorExtendedPulseTimeBuilder) WithDoorExtendedPulseTime(doorExtendedPulseTime BACnetApplicationTagUnsignedInteger) BACnetConstructedDataDoorExtendedPulseTimeBuilder {
	b.DoorExtendedPulseTime = doorExtendedPulseTime
	return b
}

func (b *_BACnetConstructedDataDoorExtendedPulseTimeBuilder) WithDoorExtendedPulseTimeBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataDoorExtendedPulseTimeBuilder {
	builder := builderSupplier(b.DoorExtendedPulseTime.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.DoorExtendedPulseTime, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataDoorExtendedPulseTimeBuilder) Build() (BACnetConstructedDataDoorExtendedPulseTime, error) {
	if b.DoorExtendedPulseTime == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'doorExtendedPulseTime' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataDoorExtendedPulseTime.deepCopy(), nil
}

func (b *_BACnetConstructedDataDoorExtendedPulseTimeBuilder) MustBuild() BACnetConstructedDataDoorExtendedPulseTime {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataDoorExtendedPulseTimeBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataDoorExtendedPulseTimeBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataDoorExtendedPulseTimeBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataDoorExtendedPulseTimeBuilder().(*_BACnetConstructedDataDoorExtendedPulseTimeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataDoorExtendedPulseTimeBuilder creates a BACnetConstructedDataDoorExtendedPulseTimeBuilder
func (b *_BACnetConstructedDataDoorExtendedPulseTime) CreateBACnetConstructedDataDoorExtendedPulseTimeBuilder() BACnetConstructedDataDoorExtendedPulseTimeBuilder {
	if b == nil {
		return NewBACnetConstructedDataDoorExtendedPulseTimeBuilder()
	}
	return &_BACnetConstructedDataDoorExtendedPulseTimeBuilder{_BACnetConstructedDataDoorExtendedPulseTime: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataDoorExtendedPulseTime) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_DOOR_EXTENDED_PULSE_TIME
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataDoorExtendedPulseTime) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataDoorExtendedPulseTime) GetDoorExtendedPulseTime() BACnetApplicationTagUnsignedInteger {
	return m.DoorExtendedPulseTime
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataDoorExtendedPulseTime) GetActualValue() BACnetApplicationTagUnsignedInteger {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagUnsignedInteger(m.GetDoorExtendedPulseTime())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataDoorExtendedPulseTime(structType any) BACnetConstructedDataDoorExtendedPulseTime {
	if casted, ok := structType.(BACnetConstructedDataDoorExtendedPulseTime); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataDoorExtendedPulseTime); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) GetTypeName() string {
	return "BACnetConstructedDataDoorExtendedPulseTime"
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (doorExtendedPulseTime)
	lengthInBits += m.DoorExtendedPulseTime.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataDoorExtendedPulseTime BACnetConstructedDataDoorExtendedPulseTime, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataDoorExtendedPulseTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataDoorExtendedPulseTime")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	doorExtendedPulseTime, err := ReadSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "doorExtendedPulseTime", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'doorExtendedPulseTime' field"))
	}
	m.DoorExtendedPulseTime = doorExtendedPulseTime

	actualValue, err := ReadVirtualField[BACnetApplicationTagUnsignedInteger](ctx, "actualValue", (*BACnetApplicationTagUnsignedInteger)(nil), doorExtendedPulseTime)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataDoorExtendedPulseTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataDoorExtendedPulseTime")
	}

	return m, nil
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataDoorExtendedPulseTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataDoorExtendedPulseTime")
		}

		if err := WriteSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "doorExtendedPulseTime", m.GetDoorExtendedPulseTime(), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'doorExtendedPulseTime' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataDoorExtendedPulseTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataDoorExtendedPulseTime")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) IsBACnetConstructedDataDoorExtendedPulseTime() {
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) deepCopy() *_BACnetConstructedDataDoorExtendedPulseTime {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataDoorExtendedPulseTimeCopy := &_BACnetConstructedDataDoorExtendedPulseTime{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.DoorExtendedPulseTime.DeepCopy().(BACnetApplicationTagUnsignedInteger),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataDoorExtendedPulseTimeCopy
}

func (m *_BACnetConstructedDataDoorExtendedPulseTime) String() string {
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
