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

// BACnetConstructedDataLightingOutputTrackingValue is the corresponding interface of BACnetConstructedDataLightingOutputTrackingValue
type BACnetConstructedDataLightingOutputTrackingValue interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetTrackingValue returns TrackingValue (property field)
	GetTrackingValue() BACnetApplicationTagReal
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagReal
	// IsBACnetConstructedDataLightingOutputTrackingValue is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataLightingOutputTrackingValue()
	// CreateBuilder creates a BACnetConstructedDataLightingOutputTrackingValueBuilder
	CreateBACnetConstructedDataLightingOutputTrackingValueBuilder() BACnetConstructedDataLightingOutputTrackingValueBuilder
}

// _BACnetConstructedDataLightingOutputTrackingValue is the data-structure of this message
type _BACnetConstructedDataLightingOutputTrackingValue struct {
	BACnetConstructedDataContract
	TrackingValue BACnetApplicationTagReal
}

var _ BACnetConstructedDataLightingOutputTrackingValue = (*_BACnetConstructedDataLightingOutputTrackingValue)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataLightingOutputTrackingValue)(nil)

// NewBACnetConstructedDataLightingOutputTrackingValue factory function for _BACnetConstructedDataLightingOutputTrackingValue
func NewBACnetConstructedDataLightingOutputTrackingValue(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, trackingValue BACnetApplicationTagReal, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataLightingOutputTrackingValue {
	if trackingValue == nil {
		panic("trackingValue of type BACnetApplicationTagReal for BACnetConstructedDataLightingOutputTrackingValue must not be nil")
	}
	_result := &_BACnetConstructedDataLightingOutputTrackingValue{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		TrackingValue:                 trackingValue,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataLightingOutputTrackingValueBuilder is a builder for BACnetConstructedDataLightingOutputTrackingValue
type BACnetConstructedDataLightingOutputTrackingValueBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(trackingValue BACnetApplicationTagReal) BACnetConstructedDataLightingOutputTrackingValueBuilder
	// WithTrackingValue adds TrackingValue (property field)
	WithTrackingValue(BACnetApplicationTagReal) BACnetConstructedDataLightingOutputTrackingValueBuilder
	// WithTrackingValueBuilder adds TrackingValue (property field) which is build by the builder
	WithTrackingValueBuilder(func(BACnetApplicationTagRealBuilder) BACnetApplicationTagRealBuilder) BACnetConstructedDataLightingOutputTrackingValueBuilder
	// Build builds the BACnetConstructedDataLightingOutputTrackingValue or returns an error if something is wrong
	Build() (BACnetConstructedDataLightingOutputTrackingValue, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataLightingOutputTrackingValue
}

// NewBACnetConstructedDataLightingOutputTrackingValueBuilder() creates a BACnetConstructedDataLightingOutputTrackingValueBuilder
func NewBACnetConstructedDataLightingOutputTrackingValueBuilder() BACnetConstructedDataLightingOutputTrackingValueBuilder {
	return &_BACnetConstructedDataLightingOutputTrackingValueBuilder{_BACnetConstructedDataLightingOutputTrackingValue: new(_BACnetConstructedDataLightingOutputTrackingValue)}
}

type _BACnetConstructedDataLightingOutputTrackingValueBuilder struct {
	*_BACnetConstructedDataLightingOutputTrackingValue

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataLightingOutputTrackingValueBuilder) = (*_BACnetConstructedDataLightingOutputTrackingValueBuilder)(nil)

func (b *_BACnetConstructedDataLightingOutputTrackingValueBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataLightingOutputTrackingValueBuilder) WithMandatoryFields(trackingValue BACnetApplicationTagReal) BACnetConstructedDataLightingOutputTrackingValueBuilder {
	return b.WithTrackingValue(trackingValue)
}

func (b *_BACnetConstructedDataLightingOutputTrackingValueBuilder) WithTrackingValue(trackingValue BACnetApplicationTagReal) BACnetConstructedDataLightingOutputTrackingValueBuilder {
	b.TrackingValue = trackingValue
	return b
}

func (b *_BACnetConstructedDataLightingOutputTrackingValueBuilder) WithTrackingValueBuilder(builderSupplier func(BACnetApplicationTagRealBuilder) BACnetApplicationTagRealBuilder) BACnetConstructedDataLightingOutputTrackingValueBuilder {
	builder := builderSupplier(b.TrackingValue.CreateBACnetApplicationTagRealBuilder())
	var err error
	b.TrackingValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagRealBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataLightingOutputTrackingValueBuilder) Build() (BACnetConstructedDataLightingOutputTrackingValue, error) {
	if b.TrackingValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'trackingValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataLightingOutputTrackingValue.deepCopy(), nil
}

func (b *_BACnetConstructedDataLightingOutputTrackingValueBuilder) MustBuild() BACnetConstructedDataLightingOutputTrackingValue {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataLightingOutputTrackingValueBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataLightingOutputTrackingValueBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataLightingOutputTrackingValueBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataLightingOutputTrackingValueBuilder().(*_BACnetConstructedDataLightingOutputTrackingValueBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataLightingOutputTrackingValueBuilder creates a BACnetConstructedDataLightingOutputTrackingValueBuilder
func (b *_BACnetConstructedDataLightingOutputTrackingValue) CreateBACnetConstructedDataLightingOutputTrackingValueBuilder() BACnetConstructedDataLightingOutputTrackingValueBuilder {
	if b == nil {
		return NewBACnetConstructedDataLightingOutputTrackingValueBuilder()
	}
	return &_BACnetConstructedDataLightingOutputTrackingValueBuilder{_BACnetConstructedDataLightingOutputTrackingValue: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataLightingOutputTrackingValue) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_LIGHTING_OUTPUT
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_TRACKING_VALUE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataLightingOutputTrackingValue) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataLightingOutputTrackingValue) GetTrackingValue() BACnetApplicationTagReal {
	return m.TrackingValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataLightingOutputTrackingValue) GetActualValue() BACnetApplicationTagReal {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagReal(m.GetTrackingValue())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataLightingOutputTrackingValue(structType any) BACnetConstructedDataLightingOutputTrackingValue {
	if casted, ok := structType.(BACnetConstructedDataLightingOutputTrackingValue); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataLightingOutputTrackingValue); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) GetTypeName() string {
	return "BACnetConstructedDataLightingOutputTrackingValue"
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (trackingValue)
	lengthInBits += m.TrackingValue.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataLightingOutputTrackingValue BACnetConstructedDataLightingOutputTrackingValue, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataLightingOutputTrackingValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataLightingOutputTrackingValue")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	trackingValue, err := ReadSimpleField[BACnetApplicationTagReal](ctx, "trackingValue", ReadComplex[BACnetApplicationTagReal](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagReal](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'trackingValue' field"))
	}
	m.TrackingValue = trackingValue

	actualValue, err := ReadVirtualField[BACnetApplicationTagReal](ctx, "actualValue", (*BACnetApplicationTagReal)(nil), trackingValue)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataLightingOutputTrackingValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataLightingOutputTrackingValue")
	}

	return m, nil
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataLightingOutputTrackingValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataLightingOutputTrackingValue")
		}

		if err := WriteSimpleField[BACnetApplicationTagReal](ctx, "trackingValue", m.GetTrackingValue(), WriteComplex[BACnetApplicationTagReal](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'trackingValue' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataLightingOutputTrackingValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataLightingOutputTrackingValue")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) IsBACnetConstructedDataLightingOutputTrackingValue() {
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) deepCopy() *_BACnetConstructedDataLightingOutputTrackingValue {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataLightingOutputTrackingValueCopy := &_BACnetConstructedDataLightingOutputTrackingValue{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.TrackingValue.DeepCopy().(BACnetApplicationTagReal),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataLightingOutputTrackingValueCopy
}

func (m *_BACnetConstructedDataLightingOutputTrackingValue) String() string {
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
