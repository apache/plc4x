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

// BACnetConstructedDataCOVPeriod is the corresponding interface of BACnetConstructedDataCOVPeriod
type BACnetConstructedDataCOVPeriod interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetCovPeriod returns CovPeriod (property field)
	GetCovPeriod() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
	// IsBACnetConstructedDataCOVPeriod is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataCOVPeriod()
	// CreateBuilder creates a BACnetConstructedDataCOVPeriodBuilder
	CreateBACnetConstructedDataCOVPeriodBuilder() BACnetConstructedDataCOVPeriodBuilder
}

// _BACnetConstructedDataCOVPeriod is the data-structure of this message
type _BACnetConstructedDataCOVPeriod struct {
	BACnetConstructedDataContract
	CovPeriod BACnetApplicationTagUnsignedInteger
}

var _ BACnetConstructedDataCOVPeriod = (*_BACnetConstructedDataCOVPeriod)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataCOVPeriod)(nil)

// NewBACnetConstructedDataCOVPeriod factory function for _BACnetConstructedDataCOVPeriod
func NewBACnetConstructedDataCOVPeriod(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, covPeriod BACnetApplicationTagUnsignedInteger, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataCOVPeriod {
	if covPeriod == nil {
		panic("covPeriod of type BACnetApplicationTagUnsignedInteger for BACnetConstructedDataCOVPeriod must not be nil")
	}
	_result := &_BACnetConstructedDataCOVPeriod{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		CovPeriod:                     covPeriod,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataCOVPeriodBuilder is a builder for BACnetConstructedDataCOVPeriod
type BACnetConstructedDataCOVPeriodBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(covPeriod BACnetApplicationTagUnsignedInteger) BACnetConstructedDataCOVPeriodBuilder
	// WithCovPeriod adds CovPeriod (property field)
	WithCovPeriod(BACnetApplicationTagUnsignedInteger) BACnetConstructedDataCOVPeriodBuilder
	// WithCovPeriodBuilder adds CovPeriod (property field) which is build by the builder
	WithCovPeriodBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataCOVPeriodBuilder
	// Build builds the BACnetConstructedDataCOVPeriod or returns an error if something is wrong
	Build() (BACnetConstructedDataCOVPeriod, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataCOVPeriod
}

// NewBACnetConstructedDataCOVPeriodBuilder() creates a BACnetConstructedDataCOVPeriodBuilder
func NewBACnetConstructedDataCOVPeriodBuilder() BACnetConstructedDataCOVPeriodBuilder {
	return &_BACnetConstructedDataCOVPeriodBuilder{_BACnetConstructedDataCOVPeriod: new(_BACnetConstructedDataCOVPeriod)}
}

type _BACnetConstructedDataCOVPeriodBuilder struct {
	*_BACnetConstructedDataCOVPeriod

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataCOVPeriodBuilder) = (*_BACnetConstructedDataCOVPeriodBuilder)(nil)

func (b *_BACnetConstructedDataCOVPeriodBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataCOVPeriodBuilder) WithMandatoryFields(covPeriod BACnetApplicationTagUnsignedInteger) BACnetConstructedDataCOVPeriodBuilder {
	return b.WithCovPeriod(covPeriod)
}

func (b *_BACnetConstructedDataCOVPeriodBuilder) WithCovPeriod(covPeriod BACnetApplicationTagUnsignedInteger) BACnetConstructedDataCOVPeriodBuilder {
	b.CovPeriod = covPeriod
	return b
}

func (b *_BACnetConstructedDataCOVPeriodBuilder) WithCovPeriodBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetConstructedDataCOVPeriodBuilder {
	builder := builderSupplier(b.CovPeriod.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.CovPeriod, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataCOVPeriodBuilder) Build() (BACnetConstructedDataCOVPeriod, error) {
	if b.CovPeriod == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'covPeriod' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataCOVPeriod.deepCopy(), nil
}

func (b *_BACnetConstructedDataCOVPeriodBuilder) MustBuild() BACnetConstructedDataCOVPeriod {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataCOVPeriodBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataCOVPeriodBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataCOVPeriodBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataCOVPeriodBuilder().(*_BACnetConstructedDataCOVPeriodBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataCOVPeriodBuilder creates a BACnetConstructedDataCOVPeriodBuilder
func (b *_BACnetConstructedDataCOVPeriod) CreateBACnetConstructedDataCOVPeriodBuilder() BACnetConstructedDataCOVPeriodBuilder {
	if b == nil {
		return NewBACnetConstructedDataCOVPeriodBuilder()
	}
	return &_BACnetConstructedDataCOVPeriodBuilder{_BACnetConstructedDataCOVPeriod: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataCOVPeriod) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataCOVPeriod) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_COV_PERIOD
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataCOVPeriod) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataCOVPeriod) GetCovPeriod() BACnetApplicationTagUnsignedInteger {
	return m.CovPeriod
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataCOVPeriod) GetActualValue() BACnetApplicationTagUnsignedInteger {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagUnsignedInteger(m.GetCovPeriod())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataCOVPeriod(structType any) BACnetConstructedDataCOVPeriod {
	if casted, ok := structType.(BACnetConstructedDataCOVPeriod); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataCOVPeriod); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataCOVPeriod) GetTypeName() string {
	return "BACnetConstructedDataCOVPeriod"
}

func (m *_BACnetConstructedDataCOVPeriod) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (covPeriod)
	lengthInBits += m.CovPeriod.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataCOVPeriod) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataCOVPeriod) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataCOVPeriod BACnetConstructedDataCOVPeriod, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataCOVPeriod"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataCOVPeriod")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	covPeriod, err := ReadSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "covPeriod", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'covPeriod' field"))
	}
	m.CovPeriod = covPeriod

	actualValue, err := ReadVirtualField[BACnetApplicationTagUnsignedInteger](ctx, "actualValue", (*BACnetApplicationTagUnsignedInteger)(nil), covPeriod)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataCOVPeriod"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataCOVPeriod")
	}

	return m, nil
}

func (m *_BACnetConstructedDataCOVPeriod) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataCOVPeriod) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataCOVPeriod"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataCOVPeriod")
		}

		if err := WriteSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "covPeriod", m.GetCovPeriod(), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'covPeriod' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataCOVPeriod"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataCOVPeriod")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataCOVPeriod) IsBACnetConstructedDataCOVPeriod() {}

func (m *_BACnetConstructedDataCOVPeriod) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataCOVPeriod) deepCopy() *_BACnetConstructedDataCOVPeriod {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataCOVPeriodCopy := &_BACnetConstructedDataCOVPeriod{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.CovPeriod.DeepCopy().(BACnetApplicationTagUnsignedInteger),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataCOVPeriodCopy
}

func (m *_BACnetConstructedDataCOVPeriod) String() string {
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
