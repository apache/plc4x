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

// BACnetConstructedDataAnalogInputFaultHighLimit is the corresponding interface of BACnetConstructedDataAnalogInputFaultHighLimit
type BACnetConstructedDataAnalogInputFaultHighLimit interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetFaultHighLimit returns FaultHighLimit (property field)
	GetFaultHighLimit() BACnetApplicationTagReal
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagReal
	// IsBACnetConstructedDataAnalogInputFaultHighLimit is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataAnalogInputFaultHighLimit()
	// CreateBuilder creates a BACnetConstructedDataAnalogInputFaultHighLimitBuilder
	CreateBACnetConstructedDataAnalogInputFaultHighLimitBuilder() BACnetConstructedDataAnalogInputFaultHighLimitBuilder
}

// _BACnetConstructedDataAnalogInputFaultHighLimit is the data-structure of this message
type _BACnetConstructedDataAnalogInputFaultHighLimit struct {
	BACnetConstructedDataContract
	FaultHighLimit BACnetApplicationTagReal
}

var _ BACnetConstructedDataAnalogInputFaultHighLimit = (*_BACnetConstructedDataAnalogInputFaultHighLimit)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataAnalogInputFaultHighLimit)(nil)

// NewBACnetConstructedDataAnalogInputFaultHighLimit factory function for _BACnetConstructedDataAnalogInputFaultHighLimit
func NewBACnetConstructedDataAnalogInputFaultHighLimit(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, faultHighLimit BACnetApplicationTagReal, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataAnalogInputFaultHighLimit {
	if faultHighLimit == nil {
		panic("faultHighLimit of type BACnetApplicationTagReal for BACnetConstructedDataAnalogInputFaultHighLimit must not be nil")
	}
	_result := &_BACnetConstructedDataAnalogInputFaultHighLimit{
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

// BACnetConstructedDataAnalogInputFaultHighLimitBuilder is a builder for BACnetConstructedDataAnalogInputFaultHighLimit
type BACnetConstructedDataAnalogInputFaultHighLimitBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(faultHighLimit BACnetApplicationTagReal) BACnetConstructedDataAnalogInputFaultHighLimitBuilder
	// WithFaultHighLimit adds FaultHighLimit (property field)
	WithFaultHighLimit(BACnetApplicationTagReal) BACnetConstructedDataAnalogInputFaultHighLimitBuilder
	// WithFaultHighLimitBuilder adds FaultHighLimit (property field) which is build by the builder
	WithFaultHighLimitBuilder(func(BACnetApplicationTagRealBuilder) BACnetApplicationTagRealBuilder) BACnetConstructedDataAnalogInputFaultHighLimitBuilder
	// Build builds the BACnetConstructedDataAnalogInputFaultHighLimit or returns an error if something is wrong
	Build() (BACnetConstructedDataAnalogInputFaultHighLimit, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataAnalogInputFaultHighLimit
}

// NewBACnetConstructedDataAnalogInputFaultHighLimitBuilder() creates a BACnetConstructedDataAnalogInputFaultHighLimitBuilder
func NewBACnetConstructedDataAnalogInputFaultHighLimitBuilder() BACnetConstructedDataAnalogInputFaultHighLimitBuilder {
	return &_BACnetConstructedDataAnalogInputFaultHighLimitBuilder{_BACnetConstructedDataAnalogInputFaultHighLimit: new(_BACnetConstructedDataAnalogInputFaultHighLimit)}
}

type _BACnetConstructedDataAnalogInputFaultHighLimitBuilder struct {
	*_BACnetConstructedDataAnalogInputFaultHighLimit

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataAnalogInputFaultHighLimitBuilder) = (*_BACnetConstructedDataAnalogInputFaultHighLimitBuilder)(nil)

func (b *_BACnetConstructedDataAnalogInputFaultHighLimitBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataAnalogInputFaultHighLimitBuilder) WithMandatoryFields(faultHighLimit BACnetApplicationTagReal) BACnetConstructedDataAnalogInputFaultHighLimitBuilder {
	return b.WithFaultHighLimit(faultHighLimit)
}

func (b *_BACnetConstructedDataAnalogInputFaultHighLimitBuilder) WithFaultHighLimit(faultHighLimit BACnetApplicationTagReal) BACnetConstructedDataAnalogInputFaultHighLimitBuilder {
	b.FaultHighLimit = faultHighLimit
	return b
}

func (b *_BACnetConstructedDataAnalogInputFaultHighLimitBuilder) WithFaultHighLimitBuilder(builderSupplier func(BACnetApplicationTagRealBuilder) BACnetApplicationTagRealBuilder) BACnetConstructedDataAnalogInputFaultHighLimitBuilder {
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

func (b *_BACnetConstructedDataAnalogInputFaultHighLimitBuilder) Build() (BACnetConstructedDataAnalogInputFaultHighLimit, error) {
	if b.FaultHighLimit == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'faultHighLimit' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataAnalogInputFaultHighLimit.deepCopy(), nil
}

func (b *_BACnetConstructedDataAnalogInputFaultHighLimitBuilder) MustBuild() BACnetConstructedDataAnalogInputFaultHighLimit {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataAnalogInputFaultHighLimitBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataAnalogInputFaultHighLimitBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataAnalogInputFaultHighLimitBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataAnalogInputFaultHighLimitBuilder().(*_BACnetConstructedDataAnalogInputFaultHighLimitBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataAnalogInputFaultHighLimitBuilder creates a BACnetConstructedDataAnalogInputFaultHighLimitBuilder
func (b *_BACnetConstructedDataAnalogInputFaultHighLimit) CreateBACnetConstructedDataAnalogInputFaultHighLimitBuilder() BACnetConstructedDataAnalogInputFaultHighLimitBuilder {
	if b == nil {
		return NewBACnetConstructedDataAnalogInputFaultHighLimitBuilder()
	}
	return &_BACnetConstructedDataAnalogInputFaultHighLimitBuilder{_BACnetConstructedDataAnalogInputFaultHighLimit: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_ANALOG_INPUT
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_FAULT_HIGH_LIMIT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) GetFaultHighLimit() BACnetApplicationTagReal {
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

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) GetActualValue() BACnetApplicationTagReal {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagReal(m.GetFaultHighLimit())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataAnalogInputFaultHighLimit(structType any) BACnetConstructedDataAnalogInputFaultHighLimit {
	if casted, ok := structType.(BACnetConstructedDataAnalogInputFaultHighLimit); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataAnalogInputFaultHighLimit); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) GetTypeName() string {
	return "BACnetConstructedDataAnalogInputFaultHighLimit"
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (faultHighLimit)
	lengthInBits += m.FaultHighLimit.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataAnalogInputFaultHighLimit BACnetConstructedDataAnalogInputFaultHighLimit, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataAnalogInputFaultHighLimit"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataAnalogInputFaultHighLimit")
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

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataAnalogInputFaultHighLimit"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataAnalogInputFaultHighLimit")
	}

	return m, nil
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataAnalogInputFaultHighLimit"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataAnalogInputFaultHighLimit")
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

		if popErr := writeBuffer.PopContext("BACnetConstructedDataAnalogInputFaultHighLimit"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataAnalogInputFaultHighLimit")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) IsBACnetConstructedDataAnalogInputFaultHighLimit() {
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) deepCopy() *_BACnetConstructedDataAnalogInputFaultHighLimit {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataAnalogInputFaultHighLimitCopy := &_BACnetConstructedDataAnalogInputFaultHighLimit{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.FaultHighLimit.DeepCopy().(BACnetApplicationTagReal),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataAnalogInputFaultHighLimitCopy
}

func (m *_BACnetConstructedDataAnalogInputFaultHighLimit) String() string {
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
