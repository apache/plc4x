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

// BACnetConstructedDataAnalogValueFaultHighLimit is the corresponding interface of BACnetConstructedDataAnalogValueFaultHighLimit
type BACnetConstructedDataAnalogValueFaultHighLimit interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetFaultHighLimit returns FaultHighLimit (property field)
	GetFaultHighLimit() BACnetApplicationTagReal
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagReal
	// IsBACnetConstructedDataAnalogValueFaultHighLimit is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataAnalogValueFaultHighLimit()
	// CreateBuilder creates a BACnetConstructedDataAnalogValueFaultHighLimitBuilder
	CreateBACnetConstructedDataAnalogValueFaultHighLimitBuilder() BACnetConstructedDataAnalogValueFaultHighLimitBuilder
}

// _BACnetConstructedDataAnalogValueFaultHighLimit is the data-structure of this message
type _BACnetConstructedDataAnalogValueFaultHighLimit struct {
	BACnetConstructedDataContract
	FaultHighLimit BACnetApplicationTagReal
}

var _ BACnetConstructedDataAnalogValueFaultHighLimit = (*_BACnetConstructedDataAnalogValueFaultHighLimit)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataAnalogValueFaultHighLimit)(nil)

// NewBACnetConstructedDataAnalogValueFaultHighLimit factory function for _BACnetConstructedDataAnalogValueFaultHighLimit
func NewBACnetConstructedDataAnalogValueFaultHighLimit(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, faultHighLimit BACnetApplicationTagReal, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataAnalogValueFaultHighLimit {
	if faultHighLimit == nil {
		panic("faultHighLimit of type BACnetApplicationTagReal for BACnetConstructedDataAnalogValueFaultHighLimit must not be nil")
	}
	_result := &_BACnetConstructedDataAnalogValueFaultHighLimit{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		FaultHighLimit:                faultHighLimit,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataAnalogValueFaultHighLimitBuilder is a builder for BACnetConstructedDataAnalogValueFaultHighLimit
type BACnetConstructedDataAnalogValueFaultHighLimitBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(faultHighLimit BACnetApplicationTagReal) BACnetConstructedDataAnalogValueFaultHighLimitBuilder
	// WithFaultHighLimit adds FaultHighLimit (property field)
	WithFaultHighLimit(BACnetApplicationTagReal) BACnetConstructedDataAnalogValueFaultHighLimitBuilder
	// WithFaultHighLimitBuilder adds FaultHighLimit (property field) which is build by the builder
	WithFaultHighLimitBuilder(func(BACnetApplicationTagRealBuilder) BACnetApplicationTagRealBuilder) BACnetConstructedDataAnalogValueFaultHighLimitBuilder
	// Build builds the BACnetConstructedDataAnalogValueFaultHighLimit or returns an error if something is wrong
	Build() (BACnetConstructedDataAnalogValueFaultHighLimit, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataAnalogValueFaultHighLimit
}

// NewBACnetConstructedDataAnalogValueFaultHighLimitBuilder() creates a BACnetConstructedDataAnalogValueFaultHighLimitBuilder
func NewBACnetConstructedDataAnalogValueFaultHighLimitBuilder() BACnetConstructedDataAnalogValueFaultHighLimitBuilder {
	return &_BACnetConstructedDataAnalogValueFaultHighLimitBuilder{_BACnetConstructedDataAnalogValueFaultHighLimit: new(_BACnetConstructedDataAnalogValueFaultHighLimit)}
}

type _BACnetConstructedDataAnalogValueFaultHighLimitBuilder struct {
	*_BACnetConstructedDataAnalogValueFaultHighLimit

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataAnalogValueFaultHighLimitBuilder) = (*_BACnetConstructedDataAnalogValueFaultHighLimitBuilder)(nil)

func (b *_BACnetConstructedDataAnalogValueFaultHighLimitBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataAnalogValueFaultHighLimitBuilder) WithMandatoryFields(faultHighLimit BACnetApplicationTagReal) BACnetConstructedDataAnalogValueFaultHighLimitBuilder {
	return b.WithFaultHighLimit(faultHighLimit)
}

func (b *_BACnetConstructedDataAnalogValueFaultHighLimitBuilder) WithFaultHighLimit(faultHighLimit BACnetApplicationTagReal) BACnetConstructedDataAnalogValueFaultHighLimitBuilder {
	b.FaultHighLimit = faultHighLimit
	return b
}

func (b *_BACnetConstructedDataAnalogValueFaultHighLimitBuilder) WithFaultHighLimitBuilder(builderSupplier func(BACnetApplicationTagRealBuilder) BACnetApplicationTagRealBuilder) BACnetConstructedDataAnalogValueFaultHighLimitBuilder {
	builder := builderSupplier(b.FaultHighLimit.CreateBACnetApplicationTagRealBuilder())
	var err error
	b.FaultHighLimit, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagRealBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataAnalogValueFaultHighLimitBuilder) Build() (BACnetConstructedDataAnalogValueFaultHighLimit, error) {
	if b.FaultHighLimit == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'faultHighLimit' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataAnalogValueFaultHighLimit.deepCopy(), nil
}

func (b *_BACnetConstructedDataAnalogValueFaultHighLimitBuilder) MustBuild() BACnetConstructedDataAnalogValueFaultHighLimit {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataAnalogValueFaultHighLimitBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataAnalogValueFaultHighLimitBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataAnalogValueFaultHighLimitBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataAnalogValueFaultHighLimitBuilder().(*_BACnetConstructedDataAnalogValueFaultHighLimitBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataAnalogValueFaultHighLimitBuilder creates a BACnetConstructedDataAnalogValueFaultHighLimitBuilder
func (b *_BACnetConstructedDataAnalogValueFaultHighLimit) CreateBACnetConstructedDataAnalogValueFaultHighLimitBuilder() BACnetConstructedDataAnalogValueFaultHighLimitBuilder {
	if b == nil {
		return NewBACnetConstructedDataAnalogValueFaultHighLimitBuilder()
	}
	return &_BACnetConstructedDataAnalogValueFaultHighLimitBuilder{_BACnetConstructedDataAnalogValueFaultHighLimit: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_ANALOG_VALUE
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_FAULT_HIGH_LIMIT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) GetFaultHighLimit() BACnetApplicationTagReal {
	return m.FaultHighLimit
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) GetActualValue() BACnetApplicationTagReal {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagReal(m.GetFaultHighLimit())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataAnalogValueFaultHighLimit(structType any) BACnetConstructedDataAnalogValueFaultHighLimit {
	if casted, ok := structType.(BACnetConstructedDataAnalogValueFaultHighLimit); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataAnalogValueFaultHighLimit); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) GetTypeName() string {
	return "BACnetConstructedDataAnalogValueFaultHighLimit"
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (faultHighLimit)
	lengthInBits += m.FaultHighLimit.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataAnalogValueFaultHighLimit BACnetConstructedDataAnalogValueFaultHighLimit, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataAnalogValueFaultHighLimit"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataAnalogValueFaultHighLimit")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	faultHighLimit, err := ReadSimpleField[BACnetApplicationTagReal](ctx, "faultHighLimit", ReadComplex[BACnetApplicationTagReal](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagReal](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'faultHighLimit' field"))
	}
	m.FaultHighLimit = faultHighLimit

	actualValue, err := ReadVirtualField[BACnetApplicationTagReal](ctx, "actualValue", (*BACnetApplicationTagReal)(nil), faultHighLimit)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataAnalogValueFaultHighLimit"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataAnalogValueFaultHighLimit")
	}

	return m, nil
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataAnalogValueFaultHighLimit"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataAnalogValueFaultHighLimit")
		}

		if err := WriteSimpleField[BACnetApplicationTagReal](ctx, "faultHighLimit", m.GetFaultHighLimit(), WriteComplex[BACnetApplicationTagReal](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'faultHighLimit' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataAnalogValueFaultHighLimit"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataAnalogValueFaultHighLimit")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) IsBACnetConstructedDataAnalogValueFaultHighLimit() {
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) deepCopy() *_BACnetConstructedDataAnalogValueFaultHighLimit {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataAnalogValueFaultHighLimitCopy := &_BACnetConstructedDataAnalogValueFaultHighLimit{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.FaultHighLimit.DeepCopy().(BACnetApplicationTagReal),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataAnalogValueFaultHighLimitCopy
}

func (m *_BACnetConstructedDataAnalogValueFaultHighLimit) String() string {
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
