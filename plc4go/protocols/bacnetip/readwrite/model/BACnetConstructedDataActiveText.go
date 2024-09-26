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

// BACnetConstructedDataActiveText is the corresponding interface of BACnetConstructedDataActiveText
type BACnetConstructedDataActiveText interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConstructedData
	// GetActiveText returns ActiveText (property field)
	GetActiveText() BACnetApplicationTagCharacterString
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagCharacterString
	// IsBACnetConstructedDataActiveText is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataActiveText()
	// CreateBuilder creates a BACnetConstructedDataActiveTextBuilder
	CreateBACnetConstructedDataActiveTextBuilder() BACnetConstructedDataActiveTextBuilder
}

// _BACnetConstructedDataActiveText is the data-structure of this message
type _BACnetConstructedDataActiveText struct {
	BACnetConstructedDataContract
	ActiveText BACnetApplicationTagCharacterString
}

var _ BACnetConstructedDataActiveText = (*_BACnetConstructedDataActiveText)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataActiveText)(nil)

// NewBACnetConstructedDataActiveText factory function for _BACnetConstructedDataActiveText
func NewBACnetConstructedDataActiveText(openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, activeText BACnetApplicationTagCharacterString, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataActiveText {
	if activeText == nil {
		panic("activeText of type BACnetApplicationTagCharacterString for BACnetConstructedDataActiveText must not be nil")
	}
	_result := &_BACnetConstructedDataActiveText{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		ActiveText:                    activeText,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConstructedDataActiveTextBuilder is a builder for BACnetConstructedDataActiveText
type BACnetConstructedDataActiveTextBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(activeText BACnetApplicationTagCharacterString) BACnetConstructedDataActiveTextBuilder
	// WithActiveText adds ActiveText (property field)
	WithActiveText(BACnetApplicationTagCharacterString) BACnetConstructedDataActiveTextBuilder
	// WithActiveTextBuilder adds ActiveText (property field) which is build by the builder
	WithActiveTextBuilder(func(BACnetApplicationTagCharacterStringBuilder) BACnetApplicationTagCharacterStringBuilder) BACnetConstructedDataActiveTextBuilder
	// Build builds the BACnetConstructedDataActiveText or returns an error if something is wrong
	Build() (BACnetConstructedDataActiveText, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConstructedDataActiveText
}

// NewBACnetConstructedDataActiveTextBuilder() creates a BACnetConstructedDataActiveTextBuilder
func NewBACnetConstructedDataActiveTextBuilder() BACnetConstructedDataActiveTextBuilder {
	return &_BACnetConstructedDataActiveTextBuilder{_BACnetConstructedDataActiveText: new(_BACnetConstructedDataActiveText)}
}

type _BACnetConstructedDataActiveTextBuilder struct {
	*_BACnetConstructedDataActiveText

	parentBuilder *_BACnetConstructedDataBuilder

	err *utils.MultiError
}

var _ (BACnetConstructedDataActiveTextBuilder) = (*_BACnetConstructedDataActiveTextBuilder)(nil)

func (b *_BACnetConstructedDataActiveTextBuilder) setParent(contract BACnetConstructedDataContract) {
	b.BACnetConstructedDataContract = contract
}

func (b *_BACnetConstructedDataActiveTextBuilder) WithMandatoryFields(activeText BACnetApplicationTagCharacterString) BACnetConstructedDataActiveTextBuilder {
	return b.WithActiveText(activeText)
}

func (b *_BACnetConstructedDataActiveTextBuilder) WithActiveText(activeText BACnetApplicationTagCharacterString) BACnetConstructedDataActiveTextBuilder {
	b.ActiveText = activeText
	return b
}

func (b *_BACnetConstructedDataActiveTextBuilder) WithActiveTextBuilder(builderSupplier func(BACnetApplicationTagCharacterStringBuilder) BACnetApplicationTagCharacterStringBuilder) BACnetConstructedDataActiveTextBuilder {
	builder := builderSupplier(b.ActiveText.CreateBACnetApplicationTagCharacterStringBuilder())
	var err error
	b.ActiveText, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagCharacterStringBuilder failed"))
	}
	return b
}

func (b *_BACnetConstructedDataActiveTextBuilder) Build() (BACnetConstructedDataActiveText, error) {
	if b.ActiveText == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'activeText' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConstructedDataActiveText.deepCopy(), nil
}

func (b *_BACnetConstructedDataActiveTextBuilder) MustBuild() BACnetConstructedDataActiveText {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConstructedDataActiveTextBuilder) Done() BACnetConstructedDataBuilder {
	return b.parentBuilder
}

func (b *_BACnetConstructedDataActiveTextBuilder) buildForBACnetConstructedData() (BACnetConstructedData, error) {
	return b.Build()
}

func (b *_BACnetConstructedDataActiveTextBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConstructedDataActiveTextBuilder().(*_BACnetConstructedDataActiveTextBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConstructedDataActiveTextBuilder creates a BACnetConstructedDataActiveTextBuilder
func (b *_BACnetConstructedDataActiveText) CreateBACnetConstructedDataActiveTextBuilder() BACnetConstructedDataActiveTextBuilder {
	if b == nil {
		return NewBACnetConstructedDataActiveTextBuilder()
	}
	return &_BACnetConstructedDataActiveTextBuilder{_BACnetConstructedDataActiveText: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataActiveText) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataActiveText) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_ACTIVE_TEXT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataActiveText) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataActiveText) GetActiveText() BACnetApplicationTagCharacterString {
	return m.ActiveText
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataActiveText) GetActualValue() BACnetApplicationTagCharacterString {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagCharacterString(m.GetActiveText())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataActiveText(structType any) BACnetConstructedDataActiveText {
	if casted, ok := structType.(BACnetConstructedDataActiveText); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataActiveText); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataActiveText) GetTypeName() string {
	return "BACnetConstructedDataActiveText"
}

func (m *_BACnetConstructedDataActiveText) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).GetLengthInBits(ctx))

	// Simple field (activeText)
	lengthInBits += m.ActiveText.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataActiveText) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataActiveText) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataActiveText BACnetConstructedDataActiveText, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataActiveText"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataActiveText")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	activeText, err := ReadSimpleField[BACnetApplicationTagCharacterString](ctx, "activeText", ReadComplex[BACnetApplicationTagCharacterString](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagCharacterString](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'activeText' field"))
	}
	m.ActiveText = activeText

	actualValue, err := ReadVirtualField[BACnetApplicationTagCharacterString](ctx, "actualValue", (*BACnetApplicationTagCharacterString)(nil), activeText)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataActiveText"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataActiveText")
	}

	return m, nil
}

func (m *_BACnetConstructedDataActiveText) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataActiveText) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataActiveText"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataActiveText")
		}

		if err := WriteSimpleField[BACnetApplicationTagCharacterString](ctx, "activeText", m.GetActiveText(), WriteComplex[BACnetApplicationTagCharacterString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'activeText' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataActiveText"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataActiveText")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataActiveText) IsBACnetConstructedDataActiveText() {}

func (m *_BACnetConstructedDataActiveText) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConstructedDataActiveText) deepCopy() *_BACnetConstructedDataActiveText {
	if m == nil {
		return nil
	}
	_BACnetConstructedDataActiveTextCopy := &_BACnetConstructedDataActiveText{
		m.BACnetConstructedDataContract.(*_BACnetConstructedData).deepCopy(),
		m.ActiveText.DeepCopy().(BACnetApplicationTagCharacterString),
	}
	m.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = m
	return _BACnetConstructedDataActiveTextCopy
}

func (m *_BACnetConstructedDataActiveText) String() string {
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
