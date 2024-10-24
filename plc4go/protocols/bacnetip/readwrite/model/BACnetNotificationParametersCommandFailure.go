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

// BACnetNotificationParametersCommandFailure is the corresponding interface of BACnetNotificationParametersCommandFailure
type BACnetNotificationParametersCommandFailure interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetNotificationParameters
	// GetInnerOpeningTag returns InnerOpeningTag (property field)
	GetInnerOpeningTag() BACnetOpeningTag
	// GetCommandValue returns CommandValue (property field)
	GetCommandValue() BACnetConstructedData
	// GetStatusFlags returns StatusFlags (property field)
	GetStatusFlags() BACnetStatusFlagsTagged
	// GetFeedbackValue returns FeedbackValue (property field)
	GetFeedbackValue() BACnetConstructedData
	// GetInnerClosingTag returns InnerClosingTag (property field)
	GetInnerClosingTag() BACnetClosingTag
	// IsBACnetNotificationParametersCommandFailure is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetNotificationParametersCommandFailure()
	// CreateBuilder creates a BACnetNotificationParametersCommandFailureBuilder
	CreateBACnetNotificationParametersCommandFailureBuilder() BACnetNotificationParametersCommandFailureBuilder
}

// _BACnetNotificationParametersCommandFailure is the data-structure of this message
type _BACnetNotificationParametersCommandFailure struct {
	BACnetNotificationParametersContract
	InnerOpeningTag BACnetOpeningTag
	CommandValue    BACnetConstructedData
	StatusFlags     BACnetStatusFlagsTagged
	FeedbackValue   BACnetConstructedData
	InnerClosingTag BACnetClosingTag
}

var _ BACnetNotificationParametersCommandFailure = (*_BACnetNotificationParametersCommandFailure)(nil)
var _ BACnetNotificationParametersRequirements = (*_BACnetNotificationParametersCommandFailure)(nil)

// NewBACnetNotificationParametersCommandFailure factory function for _BACnetNotificationParametersCommandFailure
func NewBACnetNotificationParametersCommandFailure(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, innerOpeningTag BACnetOpeningTag, commandValue BACnetConstructedData, statusFlags BACnetStatusFlagsTagged, feedbackValue BACnetConstructedData, innerClosingTag BACnetClosingTag, tagNumber uint8, objectTypeArgument BACnetObjectType) *_BACnetNotificationParametersCommandFailure {
	if innerOpeningTag == nil {
		panic("innerOpeningTag of type BACnetOpeningTag for BACnetNotificationParametersCommandFailure must not be nil")
	}
	if commandValue == nil {
		panic("commandValue of type BACnetConstructedData for BACnetNotificationParametersCommandFailure must not be nil")
	}
	if statusFlags == nil {
		panic("statusFlags of type BACnetStatusFlagsTagged for BACnetNotificationParametersCommandFailure must not be nil")
	}
	if feedbackValue == nil {
		panic("feedbackValue of type BACnetConstructedData for BACnetNotificationParametersCommandFailure must not be nil")
	}
	if innerClosingTag == nil {
		panic("innerClosingTag of type BACnetClosingTag for BACnetNotificationParametersCommandFailure must not be nil")
	}
	_result := &_BACnetNotificationParametersCommandFailure{
		BACnetNotificationParametersContract: NewBACnetNotificationParameters(openingTag, peekedTagHeader, closingTag, tagNumber, objectTypeArgument),
		InnerOpeningTag:                      innerOpeningTag,
		CommandValue:                         commandValue,
		StatusFlags:                          statusFlags,
		FeedbackValue:                        feedbackValue,
		InnerClosingTag:                      innerClosingTag,
	}
	_result.BACnetNotificationParametersContract.(*_BACnetNotificationParameters)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetNotificationParametersCommandFailureBuilder is a builder for BACnetNotificationParametersCommandFailure
type BACnetNotificationParametersCommandFailureBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(innerOpeningTag BACnetOpeningTag, commandValue BACnetConstructedData, statusFlags BACnetStatusFlagsTagged, feedbackValue BACnetConstructedData, innerClosingTag BACnetClosingTag) BACnetNotificationParametersCommandFailureBuilder
	// WithInnerOpeningTag adds InnerOpeningTag (property field)
	WithInnerOpeningTag(BACnetOpeningTag) BACnetNotificationParametersCommandFailureBuilder
	// WithInnerOpeningTagBuilder adds InnerOpeningTag (property field) which is build by the builder
	WithInnerOpeningTagBuilder(func(BACnetOpeningTagBuilder) BACnetOpeningTagBuilder) BACnetNotificationParametersCommandFailureBuilder
	// WithCommandValue adds CommandValue (property field)
	WithCommandValue(BACnetConstructedData) BACnetNotificationParametersCommandFailureBuilder
	// WithCommandValueBuilder adds CommandValue (property field) which is build by the builder
	WithCommandValueBuilder(func(BACnetConstructedDataBuilder) BACnetConstructedDataBuilder) BACnetNotificationParametersCommandFailureBuilder
	// WithStatusFlags adds StatusFlags (property field)
	WithStatusFlags(BACnetStatusFlagsTagged) BACnetNotificationParametersCommandFailureBuilder
	// WithStatusFlagsBuilder adds StatusFlags (property field) which is build by the builder
	WithStatusFlagsBuilder(func(BACnetStatusFlagsTaggedBuilder) BACnetStatusFlagsTaggedBuilder) BACnetNotificationParametersCommandFailureBuilder
	// WithFeedbackValue adds FeedbackValue (property field)
	WithFeedbackValue(BACnetConstructedData) BACnetNotificationParametersCommandFailureBuilder
	// WithFeedbackValueBuilder adds FeedbackValue (property field) which is build by the builder
	WithFeedbackValueBuilder(func(BACnetConstructedDataBuilder) BACnetConstructedDataBuilder) BACnetNotificationParametersCommandFailureBuilder
	// WithInnerClosingTag adds InnerClosingTag (property field)
	WithInnerClosingTag(BACnetClosingTag) BACnetNotificationParametersCommandFailureBuilder
	// WithInnerClosingTagBuilder adds InnerClosingTag (property field) which is build by the builder
	WithInnerClosingTagBuilder(func(BACnetClosingTagBuilder) BACnetClosingTagBuilder) BACnetNotificationParametersCommandFailureBuilder
	// Build builds the BACnetNotificationParametersCommandFailure or returns an error if something is wrong
	Build() (BACnetNotificationParametersCommandFailure, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetNotificationParametersCommandFailure
}

// NewBACnetNotificationParametersCommandFailureBuilder() creates a BACnetNotificationParametersCommandFailureBuilder
func NewBACnetNotificationParametersCommandFailureBuilder() BACnetNotificationParametersCommandFailureBuilder {
	return &_BACnetNotificationParametersCommandFailureBuilder{_BACnetNotificationParametersCommandFailure: new(_BACnetNotificationParametersCommandFailure)}
}

type _BACnetNotificationParametersCommandFailureBuilder struct {
	*_BACnetNotificationParametersCommandFailure

	parentBuilder *_BACnetNotificationParametersBuilder

	err *utils.MultiError
}

var _ (BACnetNotificationParametersCommandFailureBuilder) = (*_BACnetNotificationParametersCommandFailureBuilder)(nil)

func (b *_BACnetNotificationParametersCommandFailureBuilder) setParent(contract BACnetNotificationParametersContract) {
	b.BACnetNotificationParametersContract = contract
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithMandatoryFields(innerOpeningTag BACnetOpeningTag, commandValue BACnetConstructedData, statusFlags BACnetStatusFlagsTagged, feedbackValue BACnetConstructedData, innerClosingTag BACnetClosingTag) BACnetNotificationParametersCommandFailureBuilder {
	return b.WithInnerOpeningTag(innerOpeningTag).WithCommandValue(commandValue).WithStatusFlags(statusFlags).WithFeedbackValue(feedbackValue).WithInnerClosingTag(innerClosingTag)
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithInnerOpeningTag(innerOpeningTag BACnetOpeningTag) BACnetNotificationParametersCommandFailureBuilder {
	b.InnerOpeningTag = innerOpeningTag
	return b
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithInnerOpeningTagBuilder(builderSupplier func(BACnetOpeningTagBuilder) BACnetOpeningTagBuilder) BACnetNotificationParametersCommandFailureBuilder {
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

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithCommandValue(commandValue BACnetConstructedData) BACnetNotificationParametersCommandFailureBuilder {
	b.CommandValue = commandValue
	return b
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithCommandValueBuilder(builderSupplier func(BACnetConstructedDataBuilder) BACnetConstructedDataBuilder) BACnetNotificationParametersCommandFailureBuilder {
	builder := builderSupplier(b.CommandValue.CreateBACnetConstructedDataBuilder())
	var err error
	b.CommandValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetConstructedDataBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithStatusFlags(statusFlags BACnetStatusFlagsTagged) BACnetNotificationParametersCommandFailureBuilder {
	b.StatusFlags = statusFlags
	return b
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithStatusFlagsBuilder(builderSupplier func(BACnetStatusFlagsTaggedBuilder) BACnetStatusFlagsTaggedBuilder) BACnetNotificationParametersCommandFailureBuilder {
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

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithFeedbackValue(feedbackValue BACnetConstructedData) BACnetNotificationParametersCommandFailureBuilder {
	b.FeedbackValue = feedbackValue
	return b
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithFeedbackValueBuilder(builderSupplier func(BACnetConstructedDataBuilder) BACnetConstructedDataBuilder) BACnetNotificationParametersCommandFailureBuilder {
	builder := builderSupplier(b.FeedbackValue.CreateBACnetConstructedDataBuilder())
	var err error
	b.FeedbackValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetConstructedDataBuilder failed"))
	}
	return b
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithInnerClosingTag(innerClosingTag BACnetClosingTag) BACnetNotificationParametersCommandFailureBuilder {
	b.InnerClosingTag = innerClosingTag
	return b
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) WithInnerClosingTagBuilder(builderSupplier func(BACnetClosingTagBuilder) BACnetClosingTagBuilder) BACnetNotificationParametersCommandFailureBuilder {
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

func (b *_BACnetNotificationParametersCommandFailureBuilder) Build() (BACnetNotificationParametersCommandFailure, error) {
	if b.InnerOpeningTag == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'innerOpeningTag' not set"))
	}
	if b.CommandValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'commandValue' not set"))
	}
	if b.StatusFlags == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'statusFlags' not set"))
	}
	if b.FeedbackValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'feedbackValue' not set"))
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
	return b._BACnetNotificationParametersCommandFailure.deepCopy(), nil
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) MustBuild() BACnetNotificationParametersCommandFailure {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetNotificationParametersCommandFailureBuilder) Done() BACnetNotificationParametersBuilder {
	return b.parentBuilder
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) buildForBACnetNotificationParameters() (BACnetNotificationParameters, error) {
	return b.Build()
}

func (b *_BACnetNotificationParametersCommandFailureBuilder) DeepCopy() any {
	_copy := b.CreateBACnetNotificationParametersCommandFailureBuilder().(*_BACnetNotificationParametersCommandFailureBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetNotificationParametersCommandFailureBuilder creates a BACnetNotificationParametersCommandFailureBuilder
func (b *_BACnetNotificationParametersCommandFailure) CreateBACnetNotificationParametersCommandFailureBuilder() BACnetNotificationParametersCommandFailureBuilder {
	if b == nil {
		return NewBACnetNotificationParametersCommandFailureBuilder()
	}
	return &_BACnetNotificationParametersCommandFailureBuilder{_BACnetNotificationParametersCommandFailure: b.deepCopy()}
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

func (m *_BACnetNotificationParametersCommandFailure) GetParent() BACnetNotificationParametersContract {
	return m.BACnetNotificationParametersContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetNotificationParametersCommandFailure) GetInnerOpeningTag() BACnetOpeningTag {
	return m.InnerOpeningTag
}

func (m *_BACnetNotificationParametersCommandFailure) GetCommandValue() BACnetConstructedData {
	return m.CommandValue
}

func (m *_BACnetNotificationParametersCommandFailure) GetStatusFlags() BACnetStatusFlagsTagged {
	return m.StatusFlags
}

func (m *_BACnetNotificationParametersCommandFailure) GetFeedbackValue() BACnetConstructedData {
	return m.FeedbackValue
}

func (m *_BACnetNotificationParametersCommandFailure) GetInnerClosingTag() BACnetClosingTag {
	return m.InnerClosingTag
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetNotificationParametersCommandFailure(structType any) BACnetNotificationParametersCommandFailure {
	if casted, ok := structType.(BACnetNotificationParametersCommandFailure); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetNotificationParametersCommandFailure); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetNotificationParametersCommandFailure) GetTypeName() string {
	return "BACnetNotificationParametersCommandFailure"
}

func (m *_BACnetNotificationParametersCommandFailure) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetNotificationParametersContract.(*_BACnetNotificationParameters).GetLengthInBits(ctx))

	// Simple field (innerOpeningTag)
	lengthInBits += m.InnerOpeningTag.GetLengthInBits(ctx)

	// Simple field (commandValue)
	lengthInBits += m.CommandValue.GetLengthInBits(ctx)

	// Simple field (statusFlags)
	lengthInBits += m.StatusFlags.GetLengthInBits(ctx)

	// Simple field (feedbackValue)
	lengthInBits += m.FeedbackValue.GetLengthInBits(ctx)

	// Simple field (innerClosingTag)
	lengthInBits += m.InnerClosingTag.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetNotificationParametersCommandFailure) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetNotificationParametersCommandFailure) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetNotificationParameters, peekedTagNumber uint8, tagNumber uint8, objectTypeArgument BACnetObjectType) (__bACnetNotificationParametersCommandFailure BACnetNotificationParametersCommandFailure, err error) {
	m.BACnetNotificationParametersContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetNotificationParametersCommandFailure"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetNotificationParametersCommandFailure")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	innerOpeningTag, err := ReadSimpleField[BACnetOpeningTag](ctx, "innerOpeningTag", ReadComplex[BACnetOpeningTag](BACnetOpeningTagParseWithBufferProducer((uint8)(peekedTagNumber)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'innerOpeningTag' field"))
	}
	m.InnerOpeningTag = innerOpeningTag

	commandValue, err := ReadSimpleField[BACnetConstructedData](ctx, "commandValue", ReadComplex[BACnetConstructedData](BACnetConstructedDataParseWithBufferProducer[BACnetConstructedData]((uint8)(uint8(0)), (BACnetObjectType)(objectTypeArgument), (BACnetPropertyIdentifier)(BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE), (BACnetTagPayloadUnsignedInteger)(nil)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'commandValue' field"))
	}
	m.CommandValue = commandValue

	statusFlags, err := ReadSimpleField[BACnetStatusFlagsTagged](ctx, "statusFlags", ReadComplex[BACnetStatusFlagsTagged](BACnetStatusFlagsTaggedParseWithBufferProducer((uint8)(uint8(1)), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'statusFlags' field"))
	}
	m.StatusFlags = statusFlags

	feedbackValue, err := ReadSimpleField[BACnetConstructedData](ctx, "feedbackValue", ReadComplex[BACnetConstructedData](BACnetConstructedDataParseWithBufferProducer[BACnetConstructedData]((uint8)(uint8(2)), (BACnetObjectType)(objectTypeArgument), (BACnetPropertyIdentifier)(BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE), (BACnetTagPayloadUnsignedInteger)(nil)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'feedbackValue' field"))
	}
	m.FeedbackValue = feedbackValue

	innerClosingTag, err := ReadSimpleField[BACnetClosingTag](ctx, "innerClosingTag", ReadComplex[BACnetClosingTag](BACnetClosingTagParseWithBufferProducer((uint8)(peekedTagNumber)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'innerClosingTag' field"))
	}
	m.InnerClosingTag = innerClosingTag

	if closeErr := readBuffer.CloseContext("BACnetNotificationParametersCommandFailure"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetNotificationParametersCommandFailure")
	}

	return m, nil
}

func (m *_BACnetNotificationParametersCommandFailure) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetNotificationParametersCommandFailure) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetNotificationParametersCommandFailure"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetNotificationParametersCommandFailure")
		}

		if err := WriteSimpleField[BACnetOpeningTag](ctx, "innerOpeningTag", m.GetInnerOpeningTag(), WriteComplex[BACnetOpeningTag](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'innerOpeningTag' field")
		}

		if err := WriteSimpleField[BACnetConstructedData](ctx, "commandValue", m.GetCommandValue(), WriteComplex[BACnetConstructedData](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'commandValue' field")
		}

		if err := WriteSimpleField[BACnetStatusFlagsTagged](ctx, "statusFlags", m.GetStatusFlags(), WriteComplex[BACnetStatusFlagsTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'statusFlags' field")
		}

		if err := WriteSimpleField[BACnetConstructedData](ctx, "feedbackValue", m.GetFeedbackValue(), WriteComplex[BACnetConstructedData](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'feedbackValue' field")
		}

		if err := WriteSimpleField[BACnetClosingTag](ctx, "innerClosingTag", m.GetInnerClosingTag(), WriteComplex[BACnetClosingTag](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'innerClosingTag' field")
		}

		if popErr := writeBuffer.PopContext("BACnetNotificationParametersCommandFailure"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetNotificationParametersCommandFailure")
		}
		return nil
	}
	return m.BACnetNotificationParametersContract.(*_BACnetNotificationParameters).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetNotificationParametersCommandFailure) IsBACnetNotificationParametersCommandFailure() {
}

func (m *_BACnetNotificationParametersCommandFailure) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetNotificationParametersCommandFailure) deepCopy() *_BACnetNotificationParametersCommandFailure {
	if m == nil {
		return nil
	}
	_BACnetNotificationParametersCommandFailureCopy := &_BACnetNotificationParametersCommandFailure{
		m.BACnetNotificationParametersContract.(*_BACnetNotificationParameters).deepCopy(),
		m.InnerOpeningTag.DeepCopy().(BACnetOpeningTag),
		m.CommandValue.DeepCopy().(BACnetConstructedData),
		m.StatusFlags.DeepCopy().(BACnetStatusFlagsTagged),
		m.FeedbackValue.DeepCopy().(BACnetConstructedData),
		m.InnerClosingTag.DeepCopy().(BACnetClosingTag),
	}
	m.BACnetNotificationParametersContract.(*_BACnetNotificationParameters)._SubType = m
	return _BACnetNotificationParametersCommandFailureCopy
}

func (m *_BACnetNotificationParametersCommandFailure) String() string {
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
