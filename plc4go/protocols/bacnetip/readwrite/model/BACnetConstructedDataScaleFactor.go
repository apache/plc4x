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

// BACnetConstructedDataScaleFactor is the corresponding interface of BACnetConstructedDataScaleFactor
type BACnetConstructedDataScaleFactor interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetScaleFactor returns ScaleFactor (property field)
	GetScaleFactor() BACnetApplicationTagReal
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagReal
	// IsBACnetConstructedDataScaleFactor is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataScaleFactor()
	// CreateBuilder creates a BACnetConstructedDataScaleFactorBuilder
	CreateBACnetConstructedDataScaleFactorBuilder() BACnetConstructedDataScaleFactorBuilder
}

// _BACnetConstructedDataScaleFactor is the data-structure of this message
type _BACnetConstructedDataScaleFactor struct {
	BACnetConstructedDataContract
	ScaleFactor BACnetApplicationTagReal
}

var _ BACnetConstructedDataScaleFactor = (*_BACnetConstructedDataScaleFactor)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataScaleFactor)(nil)

// NewBACnetConstructedDataScaleFactor factory function for _BACnetConstructedDataScaleFactor
func NewBACnetConstructedDataScaleFactor(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, scaleFactor BACnetApplicationTagReal, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataScaleFactor {
	if scaleFactor == nil {
		panic("scaleFactor of type BACnetApplicationTagReal for BACnetConstructedDataScaleFactor must not be nil")
	}
	_result := &_BACnetConstructedDataScaleFactor{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		ScaleFactor:                   scaleFactor,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataScaleFactorBuilder is a builder for BACnetConstructedDataScaleFactor
type BACnetConstructedDataScaleFactorBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(scaleFactor BACnetApplicationTagReal) BACnetConstructedDataScaleFactorBuilder
	// WithScaleFactor adds ScaleFactor (property field)
	WithScaleFactor(BACnetApplicationTagReal) BACnetConstructedDataScaleFactorBuilder
	// WithScaleFactorBuilder adds ScaleFactor (property field) which is build by the builder
	WithScaleFactorBuilder(func(BACnetApplicationTagRealBuilder) BACnetApplicationTagRealBuilder) BACnetConstructedDataScaleFactorBuilder
	// Build builds the BACnetConstructedDataScaleFactor or returns an error if something is wrong
	Build() (BACnetConstructedDataScaleFactor, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataScaleFactor
}

// NewBACnetConstructedDataScaleFactorBuilder() creates a BACnetConstructedDataScaleFactorBuilder
func NewBACnetConstructedDataScaleFactorBuilder() BACnetConstructedDataScaleFactorBuilder {
	return &_BACnetConstructedDataScaleFactorBuilder{_BACnetConstructedDataScaleFactor: new(_BACnetConstructedDataScaleFactor)}
}

type _BACnetConstructedDataScaleFactorBuilder struct {
	*_BACnetConstructedDataScaleFactor

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataScaleFactorBuilder) = (*_BACnetConstructedDataScaleFactorBuilder)(nil)

func (b *_BACnetConstructedDataScaleFactorBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataScaleFactorBuilder) WithMandatoryFields(scaleFactor BACnetApplicationTagReal) BACnetConstructedDataScaleFactorBuilder {
	return b.WithScaleFactor(scaleFactor)
}

func (b *_BACnetConstructedDataScaleFactorBuilder) WithScaleFactor(scaleFactor BACnetApplicationTagReal) BACnetConstructedDataScaleFactorBuilder {
	b.ScaleFactor = scaleFactor
	return b
}

func (b *_BACnetConstructedDataScaleFactorBuilder) WithScaleFactorBuilder(builderSupplier func(BACnetApplicationTagRealBuilder) BACnetApplicationTagRealBuilder) BACnetConstructedDataScaleFactorBuilder {
	builder := builderSupplier(b.ScaleFactor.CreateBACnetApplicationTagRealBuilder())
	var err error
	b.ScaleFactor, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagRealBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataScaleFactorBuilder) Build() (BACnetConstructedDataScaleFactor, error) {
	if b.ScaleFactor == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'scaleFactor' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataScaleFactor.deepCopy(), nil
}

func (b *_BACnetConstructedDataScaleFactorBuilder) MustBuild() BACnetConstructedDataScaleFactor {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataScaleFactorBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataScaleFactorBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataScaleFactorBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataScaleFactorBuilder().(*_BACnetConstructedDataScaleFactorBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataScaleFactorBuilder creates a BACnetConstructedDataScaleFactorBuilder
func (b *_BACnetConstructedDataScaleFactor) CreateBACnetConstructedDataScaleFactorBuilder() BACnetConstructedDataScaleFactorBuilder {
	if b == nil {
		return NewBACnetConstructedDataScaleFactorBuilder()
	}
	return &_BACnetConstructedDataScaleFactorBuilder{_BACnetConstructedDataScaleFactor: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataScaleFactor) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataScaleFactor) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_SCALE_FACTOR
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataScaleFactor) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataScaleFactor) GetScaleFactor() BACnetApplicationTagReal {
	return m.ScaleFactor
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataScaleFactor) GetActualValue() BACnetApplicationTagReal {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagReal(m.GetScaleFactor())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataScaleFactor(structType any) BACnetConstructedDataScaleFactor {
	if casted, ok := structType.(BACnetConstructedDataScaleFactor); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataScaleFactor); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataScaleFactor) GetTypeName() string {
	return "BACnetConstructedDataScaleFactor"
}

func (m *_BACnetConstructedDataScaleFactor) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (scaleFactor)
	lengthInBits += m.ScaleFactor.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataScaleFactor) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataScaleFactor) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataScaleFactor BACnetConstructedDataScaleFactor, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataScaleFactor"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataScaleFactor")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	scaleFactor, err := ReadSimpleField[BACnetApplicationTagReal](ctx, "scaleFactor", ReadComplex[BACnetApplicationTagReal](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagReal](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'scaleFactor' field"))
	}
	m.ScaleFactor = scaleFactor

	actualValue, err := ReadVirtualField[BACnetApplicationTagReal](ctx, "actualValue", (*BACnetApplicationTagReal)(nil), scaleFactor)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataScaleFactor"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataScaleFactor")
	}

	return m, nil
}

func (m *_BACnetConstructedDataScaleFactor) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataScaleFactor) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataScaleFactor"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataScaleFactor")
		}

		if err := WriteSimpleField[BACnetApplicationTagReal](ctx, "scaleFactor", m.GetScaleFactor(), WriteComplex[BACnetApplicationTagReal](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'scaleFactor' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataScaleFactor"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataScaleFactor")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataScaleFactor) IsBACnetConstructedDataScaleFactor() {}

func (m *_BACnetConstructedDataScaleFactor) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataScaleFactor) deepCopy() *_BACnetConstructedDataScaleFactor {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataScaleFactorCopy := &_BACnetConstructedDataScaleFactor{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.ScaleFactor.DeepCopy().(BACnetApplicationTagReal),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataScaleFactorCopy
}

func (m *_BACnetConstructedDataScaleFactor) String() string {
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
