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

// BACnetNotificationParametersOutOfRange is the corresponding interface of BACnetNotificationParametersOutOfRange
type BACnetNotificationParametersOutOfRange interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetNotificationParameters
	// GetInnerOpeningTag returns InnerOpeningTag (property field)
	GetInnerOpeningTag() BACnetOpeningTag
	// GetExceedingValue returns ExceedingValue (property field)
	GetExceedingValue() BACnetContextTagReal
	// GetStatusFlags returns StatusFlags (property field)
	GetStatusFlags() BACnetStatusFlagsTagged
	// GetDeadband returns Deadband (property field)
	GetDeadband() BACnetContextTagReal
	// GetExceededLimit returns ExceededLimit (property field)
	GetExceededLimit() BACnetContextTagReal
	// GetInnerClosingTag returns InnerClosingTag (property field)
	GetInnerClosingTag() BACnetClosingTag
	// IsBACnetNotificationParametersOutOfRange is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetNotificationParametersOutOfRange()
	// CreateBuilder creates a BACnetNotificationParametersOutOfRangeBuilder
	CreateBACnetNotificationParametersOutOfRangeBuilder() BACnetNotificationParametersOutOfRangeBuilder
}

// _BACnetNotificationParametersOutOfRange is the data-structure of this message
type _BACnetNotificationParametersOutOfRange struct {
	BACnetNotificationParametersContract
	InnerOpeningTag BACnetOpeningTag
	ExceedingValue  BACnetContextTagReal
	StatusFlags     BACnetStatusFlagsTagged
	Deadband        BACnetContextTagReal
	ExceededLimit   BACnetContextTagReal
	InnerClosingTag BACnetClosingTag
}

var _ BACnetNotificationParametersOutOfRange = (*_BACnetNotificationParametersOutOfRange)(nil)
var _ BACnetNotificationParametersRequirements = (*_BACnetNotificationParametersOutOfRange)(nil)

// NewBACnetNotificationParametersOutOfRange factory function for _BACnetNotificationParametersOutOfRange
func NewBACnetNotificationParametersOutOfRange(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, innerOpeningTag BACnetOpeningTag, exceedingValue BACnetContextTagReal, statusFlags BACnetStatusFlagsTagged, deadband BACnetContextTagReal, exceededLimit BACnetContextTagReal, innerClosingTag BACnetClosingTag, tagNumber uint8, objectTypeArgument BACnetObjectType) *_BACnetNotificationParametersOutOfRange {
	if innerOpeningTag == nil {
		panic("innerOpeningTag of type BACnetOpeningTag for BACnetNotificationParametersOutOfRange must not be nil")
	}
	if exceedingValue == nil {
		panic("exceedingValue of type BACnetContextTagReal for BACnetNotificationParametersOutOfRange must not be nil")
	}
	if statusFlags == nil {
		panic("statusFlags of type BACnetStatusFlagsTagged for BACnetNotificationParametersOutOfRange must not be nil")
	}
	if deadband == nil {
		panic("deadband of type BACnetContextTagReal for BACnetNotificationParametersOutOfRange must not be nil")
	}
	if exceededLimit == nil {
		panic("exceededLimit of type BACnetContextTagReal for BACnetNotificationParametersOutOfRange must not be nil")
	}
	if innerClosingTag == nil {
		panic("innerClosingTag of type BACnetClosingTag for BACnetNotificationParametersOutOfRange must not be nil")
	}
	_result := &_BACnetNotificationParametersOutOfRange{
		BACnetNotificationParametersContract: NewBACnetNotificationParameters(openingTag, peekedTagHeader, closingTag, tagNumber, objectTypeArgument),
		InnerOpeningTag:                      innerOpeningTag,
		ExceedingValue:                       exceedingValue,
		StatusFlags:                          statusFlags,
		Deadband:                             deadband,
		ExceededLimit:                        exceededLimit,
		InnerClosingTag:                      innerClosingTag,
	}
	_result.BACnetNotificationParametersContract.(*_BACnetNotificationParameters)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetNotificationParametersOutOfRangeBuilder is a builder for BACnetNotificationParametersOutOfRange
type BACnetNotificationParametersOutOfRangeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(innerOpeningTag BACnetOpeningTag, exceedingValue BACnetContextTagReal, statusFlags BACnetStatusFlagsTagged, deadband BACnetContextTagReal, exceededLimit BACnetContextTagReal, innerClosingTag BACnetClosingTag) BACnetNotificationParametersOutOfRangeBuilder
	// WithInnerOpeningTag adds InnerOpeningTag (property field)
	WithInnerOpeningTag(BACnetOpeningTag) BACnetNotificationParametersOutOfRangeBuilder
	// WithInnerOpeningTagBuilder adds InnerOpeningTag (property field) which is build by the builder
	WithInnerOpeningTagBuilder(func(BACnetOpeningTagBuilder) BACnetOpeningTagBuilder) BACnetNotificationParametersOutOfRangeBuilder
	// WithExceedingValue adds ExceedingValue (property field)
	WithExceedingValue(BACnetContextTagReal) BACnetNotificationParametersOutOfRangeBuilder
	// WithExceedingValueBuilder adds ExceedingValue (property field) which is build by the builder
	WithExceedingValueBuilder(func(BACnetContextTagRealBuilder) BACnetContextTagRealBuilder) BACnetNotificationParametersOutOfRangeBuilder
	// WithStatusFlags adds StatusFlags (property field)
	WithStatusFlags(BACnetStatusFlagsTagged) BACnetNotificationParametersOutOfRangeBuilder
	// WithStatusFlagsBuilder adds StatusFlags (property field) which is build by the builder
	WithStatusFlagsBuilder(func(BACnetStatusFlagsTaggedBuilder) BACnetStatusFlagsTaggedBuilder) BACnetNotificationParametersOutOfRangeBuilder
	// WithDeadband adds Deadband (property field)
	WithDeadband(BACnetContextTagReal) BACnetNotificationParametersOutOfRangeBuilder
	// WithDeadbandBuilder adds Deadband (property field) which is build by the builder
	WithDeadbandBuilder(func(BACnetContextTagRealBuilder) BACnetContextTagRealBuilder) BACnetNotificationParametersOutOfRangeBuilder
	// WithExceededLimit adds ExceededLimit (property field)
	WithExceededLimit(BACnetContextTagReal) BACnetNotificationParametersOutOfRangeBuilder
	// WithExceededLimitBuilder adds ExceededLimit (property field) which is build by the builder
	WithExceededLimitBuilder(func(BACnetContextTagRealBuilder) BACnetContextTagRealBuilder) BACnetNotificationParametersOutOfRangeBuilder
	// WithInnerClosingTag adds InnerClosingTag (property field)
	WithInnerClosingTag(BACnetClosingTag) BACnetNotificationParametersOutOfRangeBuilder
	// WithInnerClosingTagBuilder adds InnerClosingTag (property field) which is build by the builder
	WithInnerClosingTagBuilder(func(BACnetClosingTagBuilder) BACnetClosingTagBuilder) BACnetNotificationParametersOutOfRangeBuilder
	// Build builds the BACnetNotificationParametersOutOfRange or returns an error if something is wrong
	Build() (BACnetNotificationParametersOutOfRange, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetNotificationParametersOutOfRange
}

// NewBACnetNotificationParametersOutOfRangeBuilder() creates a BACnetNotificationParametersOutOfRangeBuilder
func NewBACnetNotificationParametersOutOfRangeBuilder() BACnetNotificationParametersOutOfRangeBuilder {
	return &_BACnetNotificationParametersOutOfRangeBuilder{_BACnetNotificationParametersOutOfRange: new(_BACnetNotificationParametersOutOfRange)}
}

type _BACnetNotificationParametersOutOfRangeBuilder struct {
	*_BACnetNotificationParametersOutOfRange

	parentBuilder *_BACnetNotificationParametersBuilder

	err *utils.MultiError
}

var _ (BACnetNotificationParametersOutOfRangeBuilder) = (*_BACnetNotificationParametersOutOfRangeBuilder)(nil)

func (b *_BACnetNotificationParametersOutOfRangeBuilder) setParent(contract BACnetNotificationParametersContract) {
	b.BACnetNotificationParametersContract = contract
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithMandatoryFields(innerOpeningTag BACnetOpeningTag, exceedingValue BACnetContextTagReal, statusFlags BACnetStatusFlagsTagged, deadband BACnetContextTagReal, exceededLimit BACnetContextTagReal, innerClosingTag BACnetClosingTag) BACnetNotificationParametersOutOfRangeBuilder {
	return b.WithInnerOpeningTag(innerOpeningTag).WithExceedingValue(exceedingValue).WithStatusFlags(statusFlags).WithDeadband(deadband).WithExceededLimit(exceededLimit).WithInnerClosingTag(innerClosingTag)
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithInnerOpeningTag(innerOpeningTag BACnetOpeningTag) BACnetNotificationParametersOutOfRangeBuilder {
	b.InnerOpeningTag = innerOpeningTag
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithInnerOpeningTagBuilder(builderSupplier func(BACnetOpeningTagBuilder) BACnetOpeningTagBuilder) BACnetNotificationParametersOutOfRangeBuilder {
	builder := builderSupplier(b.InnerOpeningTag.CreateBACnetOpeningTagBuilder())
	var err error
	b.InnerOpeningTag, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetOpeningTagBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithExceedingValue(exceedingValue BACnetContextTagReal) BACnetNotificationParametersOutOfRangeBuilder {
	b.ExceedingValue = exceedingValue
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithExceedingValueBuilder(builderSupplier func(BACnetContextTagRealBuilder) BACnetContextTagRealBuilder) BACnetNotificationParametersOutOfRangeBuilder {
	builder := builderSupplier(b.ExceedingValue.CreateBACnetContextTagRealBuilder())
	var err error
	b.ExceedingValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagRealBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithStatusFlags(statusFlags BACnetStatusFlagsTagged) BACnetNotificationParametersOutOfRangeBuilder {
	b.StatusFlags = statusFlags
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithStatusFlagsBuilder(builderSupplier func(BACnetStatusFlagsTaggedBuilder) BACnetStatusFlagsTaggedBuilder) BACnetNotificationParametersOutOfRangeBuilder {
	builder := builderSupplier(b.StatusFlags.CreateBACnetStatusFlagsTaggedBuilder())
	var err error
	b.StatusFlags, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetStatusFlagsTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithDeadband(deadband BACnetContextTagReal) BACnetNotificationParametersOutOfRangeBuilder {
	b.Deadband = deadband
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithDeadbandBuilder(builderSupplier func(BACnetContextTagRealBuilder) BACnetContextTagRealBuilder) BACnetNotificationParametersOutOfRangeBuilder {
	builder := builderSupplier(b.Deadband.CreateBACnetContextTagRealBuilder())
	var err error
	b.Deadband, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagRealBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithExceededLimit(exceededLimit BACnetContextTagReal) BACnetNotificationParametersOutOfRangeBuilder {
	b.ExceededLimit = exceededLimit
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithExceededLimitBuilder(builderSupplier func(BACnetContextTagRealBuilder) BACnetContextTagRealBuilder) BACnetNotificationParametersOutOfRangeBuilder {
	builder := builderSupplier(b.ExceededLimit.CreateBACnetContextTagRealBuilder())
	var err error
	b.ExceededLimit, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagRealBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithInnerClosingTag(innerClosingTag BACnetClosingTag) BACnetNotificationParametersOutOfRangeBuilder {
	b.InnerClosingTag = innerClosingTag
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) WithInnerClosingTagBuilder(builderSupplier func(BACnetClosingTagBuilder) BACnetClosingTagBuilder) BACnetNotificationParametersOutOfRangeBuilder {
	builder := builderSupplier(b.InnerClosingTag.CreateBACnetClosingTagBuilder())
	var err error
	b.InnerClosingTag, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetClosingTagBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) Build() (BACnetNotificationParametersOutOfRange, error) {
	if b.InnerOpeningTag == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'innerOpeningTag' not set"))
	}
	if b.ExceedingValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'exceedingValue' not set"))
	}
	if b.StatusFlags == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'statusFlags' not set"))
	}
	if b.Deadband == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'deadband' not set"))
	}
	if b.ExceededLimit == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'exceededLimit' not set"))
	}
	if b.InnerClosingTag == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'innerClosingTag' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetNotificationParametersOutOfRange.deepCopy(), nil
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) MustBuild() BACnetNotificationParametersOutOfRange {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetNotificationParametersOutOfRangeBuilder) Done() BACnetNotificationParametersBuilder {
	return b.parentBuilder
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) buildForBACnetNotificationParameters() (BACnetNotificationParameters, error) {
	return b.Build()
}

func (b *_BACnetNotificationParametersOutOfRangeBuilder) DeepCopy() any {
	_copy := b.CreateBACnetNotificationParametersOutOfRangeBuilder().(*_BACnetNotificationParametersOutOfRangeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetNotificationParametersOutOfRangeBuilder creates a BACnetNotificationParametersOutOfRangeBuilder
func (b *_BACnetNotificationParametersOutOfRange) CreateBACnetNotificationParametersOutOfRangeBuilder() BACnetNotificationParametersOutOfRangeBuilder {
	if b == nil {
		return NewBACnetNotificationParametersOutOfRangeBuilder()
	}
	return &_BACnetNotificationParametersOutOfRangeBuilder{_BACnetNotificationParametersOutOfRange: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetNotificationParametersOutOfRange) GetParent() BACnetNotificationParametersContract {
	return m.BACnetNotificationParametersContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetNotificationParametersOutOfRange) GetInnerOpeningTag() BACnetOpeningTag {
	return m.InnerOpeningTag
}

func (m *_BACnetNotificationParametersOutOfRange) GetExceedingValue() BACnetContextTagReal {
	return m.ExceedingValue
}

func (m *_BACnetNotificationParametersOutOfRange) GetStatusFlags() BACnetStatusFlagsTagged {
	return m.StatusFlags
}

func (m *_BACnetNotificationParametersOutOfRange) GetDeadband() BACnetContextTagReal {
	return m.Deadband
}

func (m *_BACnetNotificationParametersOutOfRange) GetExceededLimit() BACnetContextTagReal {
	return m.ExceededLimit
}

func (m *_BACnetNotificationParametersOutOfRange) GetInnerClosingTag() BACnetClosingTag {
	return m.InnerClosingTag
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetNotificationParametersOutOfRange(structType any) BACnetNotificationParametersOutOfRange {
	if casted, ok := structType.(BACnetNotificationParametersOutOfRange); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetNotificationParametersOutOfRange); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetNotificationParametersOutOfRange) GetTypeName() string {
	return "BACnetNotificationParametersOutOfRange"
}

func (m *_BACnetNotificationParametersOutOfRange) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetNotificationParametersContract.(*_BACnetNotificationParameters).GetLengthInBits(ctx))

	// Simple field (innerOpeningTag)
	lengthInBits += m.InnerOpeningTag.GetLengthInBits(ctx)

	// Simple field (exceedingValue)
	lengthInBits += m.ExceedingValue.GetLengthInBits(ctx)

	// Simple field (statusFlags)
	lengthInBits += m.StatusFlags.GetLengthInBits(ctx)

	// Simple field (deadband)
	lengthInBits += m.Deadband.GetLengthInBits(ctx)

	// Simple field (exceededLimit)
	lengthInBits += m.ExceededLimit.GetLengthInBits(ctx)

	// Simple field (innerClosingTag)
	lengthInBits += m.InnerClosingTag.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetNotificationParametersOutOfRange) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetNotificationParametersOutOfRange) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetNotificationParameters, peekedTagNumber uint8, tagNumber uint8, objectTypeArgument BACnetObjectType) (__bACnetNotificationParametersOutOfRange BACnetNotificationParametersOutOfRange, err error) {
	m.BACnetNotificationParametersContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetNotificationParametersOutOfRange"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetNotificationParametersOutOfRange")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	innerOpeningTag, err := ReadSimpleField[BACnetOpeningTag](ctx, "innerOpeningTag", ReadComplex[BACnetOpeningTag](BACnetOpeningTagParseWithBufferProducer((uint8)(peekedTagNumber)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'innerOpeningTag' field"))
	}
	m.InnerOpeningTag = innerOpeningTag

	exceedingValue, err := ReadSimpleField[BACnetContextTagReal](ctx, "exceedingValue", ReadComplex[BACnetContextTagReal](BACnetContextTagParseWithBufferProducer[BACnetContextTagReal]((uint8)(uint8(0)), (BACnetDataType)(BACnetDataType_REAL)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'exceedingValue' field"))
	}
	m.ExceedingValue = exceedingValue

	statusFlags, err := ReadSimpleField[BACnetStatusFlagsTagged](ctx, "statusFlags", ReadComplex[BACnetStatusFlagsTagged](BACnetStatusFlagsTaggedParseWithBufferProducer((uint8)(uint8(1)), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'statusFlags' field"))
	}
	m.StatusFlags = statusFlags

	deadband, err := ReadSimpleField[BACnetContextTagReal](ctx, "deadband", ReadComplex[BACnetContextTagReal](BACnetContextTagParseWithBufferProducer[BACnetContextTagReal]((uint8)(uint8(2)), (BACnetDataType)(BACnetDataType_REAL)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'deadband' field"))
	}
	m.Deadband = deadband

	exceededLimit, err := ReadSimpleField[BACnetContextTagReal](ctx, "exceededLimit", ReadComplex[BACnetContextTagReal](BACnetContextTagParseWithBufferProducer[BACnetContextTagReal]((uint8)(uint8(3)), (BACnetDataType)(BACnetDataType_REAL)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'exceededLimit' field"))
	}
	m.ExceededLimit = exceededLimit

	innerClosingTag, err := ReadSimpleField[BACnetClosingTag](ctx, "innerClosingTag", ReadComplex[BACnetClosingTag](BACnetClosingTagParseWithBufferProducer((uint8)(peekedTagNumber)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'innerClosingTag' field"))
	}
	m.InnerClosingTag = innerClosingTag

	if closeErr := readBuffer.CloseContext("BACnetNotificationParametersOutOfRange"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetNotificationParametersOutOfRange")
	}

	return m, nil
}

func (m *_BACnetNotificationParametersOutOfRange) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetNotificationParametersOutOfRange) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetNotificationParametersOutOfRange"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetNotificationParametersOutOfRange")
		}

		if err := WriteSimpleField[BACnetOpeningTag](ctx, "innerOpeningTag", m.GetInnerOpeningTag(), WriteComplex[BACnetOpeningTag](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'innerOpeningTag' field")
		}

		if err := WriteSimpleField[BACnetContextTagReal](ctx, "exceedingValue", m.GetExceedingValue(), WriteComplex[BACnetContextTagReal](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'exceedingValue' field")
		}

		if err := WriteSimpleField[BACnetStatusFlagsTagged](ctx, "statusFlags", m.GetStatusFlags(), WriteComplex[BACnetStatusFlagsTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'statusFlags' field")
		}

		if err := WriteSimpleField[BACnetContextTagReal](ctx, "deadband", m.GetDeadband(), WriteComplex[BACnetContextTagReal](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'deadband' field")
		}

		if err := WriteSimpleField[BACnetContextTagReal](ctx, "exceededLimit", m.GetExceededLimit(), WriteComplex[BACnetContextTagReal](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'exceededLimit' field")
		}

		if err := WriteSimpleField[BACnetClosingTag](ctx, "innerClosingTag", m.GetInnerClosingTag(), WriteComplex[BACnetClosingTag](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'innerClosingTag' field")
		}

		if popErr := writeBuffer.PopContext("BACnetNotificationParametersOutOfRange"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetNotificationParametersOutOfRange")
		}
		return nil
	}
	return m.BACnetNotificationParametersContract.(*_BACnetNotificationParameters).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetNotificationParametersOutOfRange) IsBACnetNotificationParametersOutOfRange() {}

func (m *_BACnetNotificationParametersOutOfRange) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetNotificationParametersOutOfRange) deepCopy() *_BACnetNotificationParametersOutOfRange {
	if m == nil {
		return nil
	}
	_BACnetNotificationParametersOutOfRangeCopy := &_BACnetNotificationParametersOutOfRange{
		m.BACnetNotificationParametersContract.(*_BACnetNotificationParameters).deepCopy(),
		m.InnerOpeningTag.DeepCopy().(BACnetOpeningTag),
		m.ExceedingValue.DeepCopy().(BACnetContextTagReal),
		m.StatusFlags.DeepCopy().(BACnetStatusFlagsTagged),
		m.Deadband.DeepCopy().(BACnetContextTagReal),
		m.ExceededLimit.DeepCopy().(BACnetContextTagReal),
		m.InnerClosingTag.DeepCopy().(BACnetClosingTag),
	}
	m.BACnetNotificationParametersContract.(*_BACnetNotificationParameters)._SubType = m
	return _BACnetNotificationParametersOutOfRangeCopy
}

func (m *_BACnetNotificationParametersOutOfRange) String() string {
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
