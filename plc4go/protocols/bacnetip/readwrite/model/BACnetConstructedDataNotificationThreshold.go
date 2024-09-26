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

// BACnetConstructedDataNotificationThreshold is the corresponding interface of BACnetConstructedDataNotificationThreshold
type BACnetConstructedDataNotificationThreshold interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetNotificationThreshold returns NotificationThreshold (property field)
	GetNotificationThreshold() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
	// IsBACnetConstructedDataNotificationThreshold is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataNotificationThreshold()
	// CreateBuilder creates a BACnetConstructedDataNotificationThresholdBuilder
	CreateBACnetConstructedDataNotificationThresholdBuilder() BACnetConstructedDataNotificationThresholdBuilder
}

// _BACnetConstructedDataNotificationThreshold is the data-structure of this message
type _BACnetConstructedDataNotificationThreshold struct {
	BACnetConstructedDataContract
	NotificationThreshold BACnetApplicationTagUnsignedInteger
}

var _ BACnetConstructedDataNotificationThreshold = (*_BACnetConstructedDataNotificationThreshold)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataNotificationThreshold)(nil)

// NewBACnetConstructedDataNotificationThreshold factory function for _BACnetConstructedDataNotificationThreshold
func NewBACnetConstructedDataNotificationThreshold(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, notificationThreshold BACnetApplicationTagUnsignedInteger, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataNotificationThreshold {
	if notificationThreshold == nil {
		panic("notificationThreshold of type BACnetApplicationTagUnsignedInteger for BACnetConstructedDataNotificationThreshold must not be nil")
	}
	_result := &_BACnetConstructedDataNotificationThreshold{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		NotificationThreshold:         notificationThreshold,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataNotificationThresholdBuilder is a builder for BACnetConstructedDataNotificationThreshold
type BACnetConstructedDataNotificationThresholdBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(notificationThreshold BACnetApplicationTagUnsignedInteger) BACnetConstructedDataNotificationThresholdBuilder
	// WithNotificationThreshold adds NotificationThreshold (property field)
	WithNotificationThreshold(BACnetApplicationTagUnsignedInteger) BACnetConstructedDataNotificationThresholdBuilder
	// WithNotificationThresholdBuilder adds NotificationThreshold (property field) which is build by the builder
	WithNotificationThresholdBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataNotificationThresholdBuilder
	// Build builds the BACnetConstructedDataNotificationThreshold or returns an error if something is wrong
	Build() (BACnetConstructedDataNotificationThreshold, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataNotificationThreshold
}

// NewBACnetConstructedDataNotificationThresholdBuilder() creates a BACnetConstructedDataNotificationThresholdBuilder
func NewBACnetConstructedDataNotificationThresholdBuilder() BACnetConstructedDataNotificationThresholdBuilder {
	return &_BACnetConstructedDataNotificationThresholdBuilder{_BACnetConstructedDataNotificationThreshold: new(_BACnetConstructedDataNotificationThreshold)}
}

type _BACnetConstructedDataNotificationThresholdBuilder struct {
	*_BACnetConstructedDataNotificationThreshold

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataNotificationThresholdBuilder) = (*_BACnetConstructedDataNotificationThresholdBuilder)(nil)

func (b *_BACnetConstructedDataNotificationThresholdBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataNotificationThresholdBuilder) WithMandatoryFields(notificationThreshold BACnetApplicationTagUnsignedInteger) BACnetConstructedDataNotificationThresholdBuilder {
	return b.WithNotificationThreshold(notificationThreshold)
}

func (b *_BACnetConstructedDataNotificationThresholdBuilder) WithNotificationThreshold(notificationThreshold BACnetApplicationTagUnsignedInteger) BACnetConstructedDataNotificationThresholdBuilder {
	b.NotificationThreshold = notificationThreshold
	return b
}

func (b *_BACnetConstructedDataNotificationThresholdBuilder) WithNotificationThresholdBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataNotificationThresholdBuilder {
	builder := builderSupplier(b.NotificationThreshold.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.NotificationThreshold, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataNotificationThresholdBuilder) Build() (BACnetConstructedDataNotificationThreshold, error) {
	if b.NotificationThreshold == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'notificationThreshold' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataNotificationThreshold.deepCopy(), nil
}

func (b *_BACnetConstructedDataNotificationThresholdBuilder) MustBuild() BACnetConstructedDataNotificationThreshold {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataNotificationThresholdBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataNotificationThresholdBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataNotificationThresholdBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataNotificationThresholdBuilder().(*_BACnetConstructedDataNotificationThresholdBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataNotificationThresholdBuilder creates a BACnetConstructedDataNotificationThresholdBuilder
func (b *_BACnetConstructedDataNotificationThreshold) CreateBACnetConstructedDataNotificationThresholdBuilder() BACnetConstructedDataNotificationThresholdBuilder {
	if b == nil {
		return NewBACnetConstructedDataNotificationThresholdBuilder()
	}
	return &_BACnetConstructedDataNotificationThresholdBuilder{_BACnetConstructedDataNotificationThreshold: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataNotificationThreshold) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataNotificationThreshold) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_NOTIFICATION_THRESHOLD
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataNotificationThreshold) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataNotificationThreshold) GetNotificationThreshold() BACnetApplicationTagUnsignedInteger {
	return m.NotificationThreshold
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataNotificationThreshold) GetActualValue() BACnetApplicationTagUnsignedInteger {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagUnsignedInteger(m.GetNotificationThreshold())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataNotificationThreshold(structType any) BACnetConstructedDataNotificationThreshold {
	if casted, ok := structType.(BACnetConstructedDataNotificationThreshold); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataNotificationThreshold); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataNotificationThreshold) GetTypeName() string {
	return "BACnetConstructedDataNotificationThreshold"
}

func (m *_BACnetConstructedDataNotificationThreshold) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (notificationThreshold)
	lengthInBits += m.NotificationThreshold.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataNotificationThreshold) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataNotificationThreshold) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataNotificationThreshold BACnetConstructedDataNotificationThreshold, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataNotificationThreshold"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataNotificationThreshold")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	notificationThreshold, err := ReadSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "notificationThreshold", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'notificationThreshold' field"))
	}
	m.NotificationThreshold = notificationThreshold

	actualValue, err := ReadVirtualField[BACnetApplicationTagUnsignedInteger](ctx, "actualValue", (*BACnetApplicationTagUnsignedInteger)(nil), notificationThreshold)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataNotificationThreshold"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataNotificationThreshold")
	}

	return m, nil
}

func (m *_BACnetConstructedDataNotificationThreshold) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataNotificationThreshold) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataNotificationThreshold"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataNotificationThreshold")
		}

		if err := WriteSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "notificationThreshold", m.GetNotificationThreshold(), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'notificationThreshold' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataNotificationThreshold"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataNotificationThreshold")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataNotificationThreshold) IsBACnetConstructedDataNotificationThreshold() {
}

func (m *_BACnetConstructedDataNotificationThreshold) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataNotificationThreshold) deepCopy() *_BACnetConstructedDataNotificationThreshold {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataNotificationThresholdCopy := &_BACnetConstructedDataNotificationThreshold{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.NotificationThreshold.DeepCopy().(BACnetApplicationTagUnsignedInteger),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataNotificationThresholdCopy
}

func (m *_BACnetConstructedDataNotificationThreshold) String() string {
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
