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

// BACnetTimerStateChangeValueCharacterString is the corresponding interface of BACnetTimerStateChangeValueCharacterString
type BACnetTimerStateChangeValueCharacterString interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetTimerStateChangeValue
	// GetCharacterStringValue returns CharacterStringValue (property field)
	GetCharacterStringValue() BACnetApplicationTagCharacterString
	// IsBACnetTimerStateChangeValueCharacterString is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetTimerStateChangeValueCharacterString()
	// CreateBuilder creates a BACnetTimerStateChangeValueCharacterStringBuilder
	CreateBACnetTimerStateChangeValueCharacterStringBuilder() BACnetTimerStateChangeValueCharacterStringBuilder
}

// _BACnetTimerStateChangeValueCharacterString is the data-structure of this message
type _BACnetTimerStateChangeValueCharacterString struct {
	BACnetTimerStateChangeValueContract
	CharacterStringValue BACnetApplicationTagCharacterString
}

var _ BACnetTimerStateChangeValueCharacterString = (*_BACnetTimerStateChangeValueCharacterString)(nil)
var _ BACnetTimerStateChangeValueRequirements = (*_BACnetTimerStateChangeValueCharacterString)(nil)

// NewBACnetTimerStateChangeValueCharacterString factory function for _BACnetTimerStateChangeValueCharacterString
func NewBACnetTimerStateChangeValueCharacterString(peekedTagHeader BACnetTagHeader, characterStringValue BACnetApplicationTagCharacterString, objectTypeArgument BACnetObjectType) *_BACnetTimerStateChangeValueCharacterString {
	if characterStringValue == nil {
		panic("characterStringValue of type BACnetApplicationTagCharacterString for BACnetTimerStateChangeValueCharacterString must not be nil")
	}
	_result := &_BACnetTimerStateChangeValueCharacterString{
		BACnetTimerStateChangeValueContract: NewBACnetTimerStateChangeValue(peekedTagHeader, objectTypeArgument),
		CharacterStringValue:                characterStringValue,
	}
	_result.BACnetTimerStateChangeValueContract.(*_BACnetTimerStateChangeValue)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetTimerStateChangeValueCharacterStringBuilder is a builder for BACnetTimerStateChangeValueCharacterString
type BACnetTimerStateChangeValueCharacterStringBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(characterStringValue BACnetApplicationTagCharacterString) BACnetTimerStateChangeValueCharacterStringBuilder
	// WithCharacterStringValue adds CharacterStringValue (property field)
	WithCharacterStringValue(BACnetApplicationTagCharacterString) BACnetTimerStateChangeValueCharacterStringBuilder
	// WithCharacterStringValueBuilder adds CharacterStringValue (property field) which is build by the builder
	WithCharacterStringValueBuilder(func(BACnetApplicationTagCharacterStringBuilder) BACnetApplicationTagCharacterStringBuilder) BACnetTimerStateChangeValueCharacterStringBuilder
	// Build builds the BACnetTimerStateChangeValueCharacterString or returns an error if something is wrong
	Build() (BACnetTimerStateChangeValueCharacterString, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetTimerStateChangeValueCharacterString
}

// NewBACnetTimerStateChangeValueCharacterStringBuilder() creates a BACnetTimerStateChangeValueCharacterStringBuilder
func NewBACnetTimerStateChangeValueCharacterStringBuilder() BACnetTimerStateChangeValueCharacterStringBuilder {
	return &_BACnetTimerStateChangeValueCharacterStringBuilder{_BACnetTimerStateChangeValueCharacterString: new(_BACnetTimerStateChangeValueCharacterString)}
}

type _BACnetTimerStateChangeValueCharacterStringBuilder struct {
	*_BACnetTimerStateChangeValueCharacterString

	parentBuilder *_BACnetTimerStateChangeValueBuilder

	err *utils.MultiError
}

var _ (BACnetTimerStateChangeValueCharacterStringBuilder) = (*_BACnetTimerStateChangeValueCharacterStringBuilder)(nil)

func (b *_BACnetTimerStateChangeValueCharacterStringBuilder) setParent(contract BACnetTimerStateChangeValueContract) {
	b.BACnetTimerStateChangeValueContract = contract
}

func (b *_BACnetTimerStateChangeValueCharacterStringBuilder) WithMandatoryFields(characterStringValue BACnetApplicationTagCharacterString) BACnetTimerStateChangeValueCharacterStringBuilder {
	return b.WithCharacterStringValue(characterStringValue)
}

func (b *_BACnetTimerStateChangeValueCharacterStringBuilder) WithCharacterStringValue(characterStringValue BACnetApplicationTagCharacterString) BACnetTimerStateChangeValueCharacterStringBuilder {
	b.CharacterStringValue = characterStringValue
	return b
}

func (b *_BACnetTimerStateChangeValueCharacterStringBuilder) WithCharacterStringValueBuilder(builderSupplier func(BACnetApplicationTagCharacterStringBuilder) BACnetApplicationTagCharacterStringBuilder) BACnetTimerStateChangeValueCharacterStringBuilder {
	builder := builderSupplier(b.CharacterStringValue.CreateBACnetApplicationTagCharacterStringBuilder())
	var err error
	b.CharacterStringValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagCharacterStringBuilder failed"))
	}
	return b
}

func (b *_BACnetTimerStateChangeValueCharacterStringBuilder) Build() (BACnetTimerStateChangeValueCharacterString, error) {
	if b.CharacterStringValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'characterStringValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetTimerStateChangeValueCharacterString.deepCopy(), nil
}

func (b *_BACnetTimerStateChangeValueCharacterStringBuilder) MustBuild() BACnetTimerStateChangeValueCharacterString {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetTimerStateChangeValueCharacterStringBuilder) Done() BACnetTimerStateChangeValueBuilder {
	return b.parentBuilder
}

func (b *_BACnetTimerStateChangeValueCharacterStringBuilder) buildForBACnetTimerStateChangeValue() (BACnetTimerStateChangeValue, error) {
	return b.Build()
}

func (b *_BACnetTimerStateChangeValueCharacterStringBuilder) DeepCopy() any {
	_copy := b.CreateBACnetTimerStateChangeValueCharacterStringBuilder().(*_BACnetTimerStateChangeValueCharacterStringBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetTimerStateChangeValueCharacterStringBuilder creates a BACnetTimerStateChangeValueCharacterStringBuilder
func (b *_BACnetTimerStateChangeValueCharacterString) CreateBACnetTimerStateChangeValueCharacterStringBuilder() BACnetTimerStateChangeValueCharacterStringBuilder {
	if b == nil {
		return NewBACnetTimerStateChangeValueCharacterStringBuilder()
	}
	return &_BACnetTimerStateChangeValueCharacterStringBuilder{_BACnetTimerStateChangeValueCharacterString: b.deepCopy()}
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

func (m *_BACnetTimerStateChangeValueCharacterString) GetParent() BACnetTimerStateChangeValueContract {
	return m.BACnetTimerStateChangeValueContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetTimerStateChangeValueCharacterString) GetCharacterStringValue() BACnetApplicationTagCharacterString {
	return m.CharacterStringValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetTimerStateChangeValueCharacterString(structType any) BACnetTimerStateChangeValueCharacterString {
	if casted, ok := structType.(BACnetTimerStateChangeValueCharacterString); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetTimerStateChangeValueCharacterString); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetTimerStateChangeValueCharacterString) GetTypeName() string {
	return "BACnetTimerStateChangeValueCharacterString"
}

func (m *_BACnetTimerStateChangeValueCharacterString) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetTimerStateChangeValueContract.(*_BACnetTimerStateChangeValue).GetLengthInBits(ctx))

	// Simple field (characterStringValue)
	lengthInBits += m.CharacterStringValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetTimerStateChangeValueCharacterString) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetTimerStateChangeValueCharacterString) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetTimerStateChangeValue, objectTypeArgument BACnetObjectType) (__bACnetTimerStateChangeValueCharacterString BACnetTimerStateChangeValueCharacterString, err error) {
	m.BACnetTimerStateChangeValueContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetTimerStateChangeValueCharacterString"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetTimerStateChangeValueCharacterString")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	characterStringValue, err := ReadSimpleField[BACnetApplicationTagCharacterString](ctx, "characterStringValue", ReadComplex[BACnetApplicationTagCharacterString](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagCharacterString](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'characterStringValue' field"))
	}
	m.CharacterStringValue = characterStringValue

	if closeErr := readBuffer.CloseContext("BACnetTimerStateChangeValueCharacterString"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetTimerStateChangeValueCharacterString")
	}

	return m, nil
}

func (m *_BACnetTimerStateChangeValueCharacterString) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetTimerStateChangeValueCharacterString) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetTimerStateChangeValueCharacterString"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetTimerStateChangeValueCharacterString")
		}

		if err := WriteSimpleField[BACnetApplicationTagCharacterString](ctx, "characterStringValue", m.GetCharacterStringValue(), WriteComplex[BACnetApplicationTagCharacterString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'characterStringValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetTimerStateChangeValueCharacterString"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetTimerStateChangeValueCharacterString")
		}
		return nil
	}
	return m.BACnetTimerStateChangeValueContract.(*_BACnetTimerStateChangeValue).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetTimerStateChangeValueCharacterString) IsBACnetTimerStateChangeValueCharacterString() {
}

func (m *_BACnetTimerStateChangeValueCharacterString) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetTimerStateChangeValueCharacterString) deepCopy() *_BACnetTimerStateChangeValueCharacterString {
	if m == nil {
		return nil
	}
	_BACnetTimerStateChangeValueCharacterStringCopy := &_BACnetTimerStateChangeValueCharacterString{
		m.BACnetTimerStateChangeValueContract.(*_BACnetTimerStateChangeValue).deepCopy(),
		m.CharacterStringValue.DeepCopy().(BACnetApplicationTagCharacterString),
	}
	m.BACnetTimerStateChangeValueContract.(*_BACnetTimerStateChangeValue)._SubType = m
	return _BACnetTimerStateChangeValueCharacterStringCopy
}

func (m *_BACnetTimerStateChangeValueCharacterString) String() string {
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
