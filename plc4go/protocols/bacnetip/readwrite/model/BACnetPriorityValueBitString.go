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

// BACnetPriorityValueBitString is the corresponding interface of BACnetPriorityValueBitString
type BACnetPriorityValueBitString interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetPriorityValue
	// GetBitStringValue returns BitStringValue (property field)
	GetBitStringValue() BACnetApplicationTagBitString
	// IsBACnetPriorityValueBitString is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetPriorityValueBitString()
	// CreateBuilder creates a BACnetPriorityValueBitStringBuilder
	CreateBACnetPriorityValueBitStringBuilder() BACnetPriorityValueBitStringBuilder
}

// _BACnetPriorityValueBitString is the data-structure of this message
type _BACnetPriorityValueBitString struct {
	BACnetPriorityValueContract
	BitStringValue BACnetApplicationTagBitString
}

var _ BACnetPriorityValueBitString = (*_BACnetPriorityValueBitString)(nil)
var _ BACnetPriorityValueRequirements = (*_BACnetPriorityValueBitString)(nil)

// NewBACnetPriorityValueBitString factory function for _BACnetPriorityValueBitString
func NewBACnetPriorityValueBitString(peekedTagHeader BACnetTagHeader, bitStringValue BACnetApplicationTagBitString, objectTypeArgument BACnetObjectType) *_BACnetPriorityValueBitString {
	if bitStringValue == nil {
		panic("bitStringValue of type BACnetApplicationTagBitString for BACnetPriorityValueBitString must not be nil")
	}
	_result := &_BACnetPriorityValueBitString{
		BACnetPriorityValueContract: NewBACnetPriorityValue(peekedTagHeader, objectTypeArgument),
		BitStringValue:              bitStringValue,
	}
	_result.BACnetPriorityValueContract.(*_BACnetPriorityValue)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetPriorityValueBitStringBuilder is a builder for BACnetPriorityValueBitString
type BACnetPriorityValueBitStringBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(bitStringValue BACnetApplicationTagBitString) BACnetPriorityValueBitStringBuilder
	// WithBitStringValue adds BitStringValue (property field)
	WithBitStringValue(BACnetApplicationTagBitString) BACnetPriorityValueBitStringBuilder
	// WithBitStringValueBuilder adds BitStringValue (property field) which is build by the builder
	WithBitStringValueBuilder(func(BACnetApplicationTagBitStringBuilder) BACnetApplicationTagBitStringBuilder) BACnetPriorityValueBitStringBuilder
	// Build builds the BACnetPriorityValueBitString or returns an error if something is wrong
	Build() (BACnetPriorityValueBitString, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetPriorityValueBitString
}

// NewBACnetPriorityValueBitStringBuilder() creates a BACnetPriorityValueBitStringBuilder
func NewBACnetPriorityValueBitStringBuilder() BACnetPriorityValueBitStringBuilder {
	return &_BACnetPriorityValueBitStringBuilder{_BACnetPriorityValueBitString: new(_BACnetPriorityValueBitString)}
}

type _BACnetPriorityValueBitStringBuilder struct {
	*_BACnetPriorityValueBitString

	parentBuilder *_BACnetPriorityValueBuilder

	err *utils.MultiError
}

var _ (BACnetPriorityValueBitStringBuilder) = (*_BACnetPriorityValueBitStringBuilder)(nil)

func (b *_BACnetPriorityValueBitStringBuilder) setParent(contract BACnetPriorityValueContract) {
	b.BACnetPriorityValueContract = contract
}

func (b *_BACnetPriorityValueBitStringBuilder) WithMandatoryFields(bitStringValue BACnetApplicationTagBitString) BACnetPriorityValueBitStringBuilder {
	return b.WithBitStringValue(bitStringValue)
}

func (b *_BACnetPriorityValueBitStringBuilder) WithBitStringValue(bitStringValue BACnetApplicationTagBitString) BACnetPriorityValueBitStringBuilder {
	b.BitStringValue = bitStringValue
	return b
}

func (b *_BACnetPriorityValueBitStringBuilder) WithBitStringValueBuilder(builderSupplier func(BACnetApplicationTagBitStringBuilder) BACnetApplicationTagBitStringBuilder) BACnetPriorityValueBitStringBuilder {
	builder := builderSupplier(b.BitStringValue.CreateBACnetApplicationTagBitStringBuilder())
	var err error
	b.BitStringValue, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagBitStringBuilder failed"))
	}
	return b
}

func (b *_BACnetPriorityValueBitStringBuilder) Build() (BACnetPriorityValueBitString, error) {
	if b.BitStringValue == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'bitStringValue' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetPriorityValueBitString.deepCopy(), nil
}

func (b *_BACnetPriorityValueBitStringBuilder) MustBuild() BACnetPriorityValueBitString {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetPriorityValueBitStringBuilder) Done() BACnetPriorityValueBuilder {
	return b.parentBuilder
}

func (b *_BACnetPriorityValueBitStringBuilder) buildForBACnetPriorityValue() (BACnetPriorityValue, error) {
	return b.Build()
}

func (b *_BACnetPriorityValueBitStringBuilder) DeepCopy() any {
	_copy := b.CreateBACnetPriorityValueBitStringBuilder().(*_BACnetPriorityValueBitStringBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetPriorityValueBitStringBuilder creates a BACnetPriorityValueBitStringBuilder
func (b *_BACnetPriorityValueBitString) CreateBACnetPriorityValueBitStringBuilder() BACnetPriorityValueBitStringBuilder {
	if b == nil {
		return NewBACnetPriorityValueBitStringBuilder()
	}
	return &_BACnetPriorityValueBitStringBuilder{_BACnetPriorityValueBitString: b.deepCopy()}
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

func (m *_BACnetPriorityValueBitString) GetParent() BACnetPriorityValueContract {
	return m.BACnetPriorityValueContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetPriorityValueBitString) GetBitStringValue() BACnetApplicationTagBitString {
	return m.BitStringValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetPriorityValueBitString(structType any) BACnetPriorityValueBitString {
	if casted, ok := structType.(BACnetPriorityValueBitString); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetPriorityValueBitString); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetPriorityValueBitString) GetTypeName() string {
	return "BACnetPriorityValueBitString"
}

func (m *_BACnetPriorityValueBitString) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetPriorityValueContract.(*_BACnetPriorityValue).GetLengthInBits(ctx))

	// Simple field (bitStringValue)
	lengthInBits += m.BitStringValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetPriorityValueBitString) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetPriorityValueBitString) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetPriorityValue, objectTypeArgument BACnetObjectType) (__bACnetPriorityValueBitString BACnetPriorityValueBitString, err error) {
	m.BACnetPriorityValueContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetPriorityValueBitString"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetPriorityValueBitString")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	bitStringValue, err := ReadSimpleField[BACnetApplicationTagBitString](ctx, "bitStringValue", ReadComplex[BACnetApplicationTagBitString](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagBitString](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'bitStringValue' field"))
	}
	m.BitStringValue = bitStringValue

	if closeErr := readBuffer.CloseContext("BACnetPriorityValueBitString"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetPriorityValueBitString")
	}

	return m, nil
}

func (m *_BACnetPriorityValueBitString) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetPriorityValueBitString) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPriorityValueBitString"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetPriorityValueBitString")
		}

		if err := WriteSimpleField[BACnetApplicationTagBitString](ctx, "bitStringValue", m.GetBitStringValue(), WriteComplex[BACnetApplicationTagBitString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'bitStringValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetPriorityValueBitString"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetPriorityValueBitString")
		}
		return nil
	}
	return m.BACnetPriorityValueContract.(*_BACnetPriorityValue).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetPriorityValueBitString) IsBACnetPriorityValueBitString() {}

func (m *_BACnetPriorityValueBitString) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetPriorityValueBitString) deepCopy() *_BACnetPriorityValueBitString {
	if m == nil {
		return nil
	}
	_BACnetPriorityValueBitStringCopy := &_BACnetPriorityValueBitString{
		m.BACnetPriorityValueContract.(*_BACnetPriorityValue).deepCopy(),
		m.BitStringValue.DeepCopy().(BACnetApplicationTagBitString),
	}
	m.BACnetPriorityValueContract.(*_BACnetPriorityValue)._SubType = m
	return _BACnetPriorityValueBitStringCopy
}

func (m *_BACnetPriorityValueBitString) String() string {
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
