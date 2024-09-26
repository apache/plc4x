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

// BACnetConstructedDataMode is the corresponding interface of BACnetConstructedDataMode
type BACnetConstructedDataMode interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetMode returns Mode (property field)
	GetMode() BACnetLifeSafetyModeTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetLifeSafetyModeTagged
	// IsBACnetConstructedDataMode is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataMode()
	// CreateBuilder creates a BACnetConstructedDataModeBuilder
	CreateBACnetConstructedDataModeBuilder() BACnetConstructedDataModeBuilder
}

// _BACnetConstructedDataMode is the data-structure of this message
type _BACnetConstructedDataMode struct {
	BACnetConstructedDataContract
	Mode BACnetLifeSafetyModeTagged
}

var _ BACnetConstructedDataMode = (*_BACnetConstructedDataMode)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataMode)(nil)

// NewBACnetConstructedDataMode factory function for _BACnetConstructedDataMode
func NewBACnetConstructedDataMode(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, mode BACnetLifeSafetyModeTagged, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataMode {
	if mode == nil {
		panic("mode of type BACnetLifeSafetyModeTagged for BACnetConstructedDataMode must not be nil")
	}
	_result := &_BACnetConstructedDataMode{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		Mode:                          mode,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataModeBuilder is a builder for BACnetConstructedDataMode
type BACnetConstructedDataModeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(mode BACnetLifeSafetyModeTagged) BACnetConstructedDataModeBuilder
	// WithMode adds Mode (property field)
	WithMode(BACnetLifeSafetyModeTagged) BACnetConstructedDataModeBuilder
	// WithModeBuilder adds Mode (property field) which is build by the builder
	WithModeBuilder(func(BACnetLifeSafetyModeTaggedBuilder) BACnetLifeSafetyModeTaggedBuilder) BACnetConstructedDataModeBuilder
	// Build builds the BACnetConstructedDataMode or returns an error if something is wrong
	Build() (BACnetConstructedDataMode, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataMode
}

// NewBACnetConstructedDataModeBuilder() creates a BACnetConstructedDataModeBuilder
func NewBACnetConstructedDataModeBuilder() BACnetConstructedDataModeBuilder {
	return &_BACnetConstructedDataModeBuilder{_BACnetConstructedDataMode: new(_BACnetConstructedDataMode)}
}

type _BACnetConstructedDataModeBuilder struct {
	*_BACnetConstructedDataMode

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataModeBuilder) = (*_BACnetConstructedDataModeBuilder)(nil)

func (b *_BACnetConstructedDataModeBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataModeBuilder) WithMandatoryFields(mode BACnetLifeSafetyModeTagged) BACnetConstructedDataModeBuilder {
	return b.WithMode(mode)
}

func (b *_BACnetConstructedDataModeBuilder) WithMode(mode BACnetLifeSafetyModeTagged) BACnetConstructedDataModeBuilder {
	b.Mode = mode
	return b
}

func (b *_BACnetConstructedDataModeBuilder) WithModeBuilder(builderSupplier func(BACnetLifeSafetyModeTaggedBuilder) BACnetLifeSafetyModeTaggedBuilder) BACnetConstructedDataModeBuilder {
	builder := builderSupplier(b.Mode.CreateBACnetLifeSafetyModeTaggedBuilder())
	var err error
	b.Mode, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetLifeSafetyModeTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataModeBuilder) Build() (BACnetConstructedDataMode, error) {
	if b.Mode == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'mode' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataMode.deepCopy(), nil
}

func (b *_BACnetConstructedDataModeBuilder) MustBuild() BACnetConstructedDataMode {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataModeBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataModeBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataModeBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataModeBuilder().(*_BACnetConstructedDataModeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataModeBuilder creates a BACnetConstructedDataModeBuilder
func (b *_BACnetConstructedDataMode) CreateBACnetConstructedDataModeBuilder() BACnetConstructedDataModeBuilder {
	if b == nil {
		return NewBACnetConstructedDataModeBuilder()
	}
	return &_BACnetConstructedDataModeBuilder{_BACnetConstructedDataMode: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataMode) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataMode) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_MODE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataMode) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataMode) GetMode() BACnetLifeSafetyModeTagged {
	return m.Mode
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataMode) GetActualValue() BACnetLifeSafetyModeTagged {
	ctx := context.Background()
	_ = ctx
	return CastBACnetLifeSafetyModeTagged(m.GetMode())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataMode(structType any) BACnetConstructedDataMode {
	if casted, ok := structType.(BACnetConstructedDataMode); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataMode); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataMode) GetTypeName() string {
	return "BACnetConstructedDataMode"
}

func (m *_BACnetConstructedDataMode) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (mode)
	lengthInBits += m.Mode.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataMode) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataMode) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataMode BACnetConstructedDataMode, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataMode"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataMode")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	mode, err := ReadSimpleField[BACnetLifeSafetyModeTagged](ctx, "mode", ReadComplex[BACnetLifeSafetyModeTagged](BACnetLifeSafetyModeTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'mode' field"))
	}
	m.Mode = mode

	actualValue, err := ReadVirtualField[BACnetLifeSafetyModeTagged](ctx, "actualValue", (*BACnetLifeSafetyModeTagged)(nil), mode)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataMode"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataMode")
	}

	return m, nil
}

func (m *_BACnetConstructedDataMode) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataMode) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataMode"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataMode")
		}

		if err := WriteSimpleField[BACnetLifeSafetyModeTagged](ctx, "mode", m.GetMode(), WriteComplex[BACnetLifeSafetyModeTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'mode' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataMode"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataMode")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataMode) IsBACnetConstructedDataMode() {}

func (m *_BACnetConstructedDataMode) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataMode) deepCopy() *_BACnetConstructedDataMode {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataModeCopy := &_BACnetConstructedDataMode{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.Mode.DeepCopy().(BACnetLifeSafetyModeTagged),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataModeCopy
}

func (m *_BACnetConstructedDataMode) String() string {
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
