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

// VariantString is the corresponding interface of VariantString
type VariantString interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	Variant
	// GetArrayLength returns ArrayLength (property field)
	GetArrayLength() *int32
	// GetValue returns Value (property field)
	GetValue() []PascalString
	// IsVariantString is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsVariantString()
	// CreateBuilder creates a VariantStringBuilder
	CreateVariantStringBuilder() VariantStringBuilder
}

// _VariantString is the data-structure of this message
type _VariantString struct {
	VariantContract
	ArrayLength *int32
	Value       []PascalString
}

var _ VariantString = (*_VariantString)(nil)
var _ VariantRequirements = (*_VariantString)(nil)

// NewVariantString factory function for _VariantString
func NewVariantString(arrayLengthSpecified bool, arrayDimensionsSpecified bool, noOfArrayDimensions *int32, arrayDimensions []bool, arrayLength *int32, value []PascalString) *_VariantString {
	_result := &_VariantString{
		VariantContract: NewVariant(arrayLengthSpecified, arrayDimensionsSpecified, noOfArrayDimensions, arrayDimensions),
		ArrayLength:     arrayLength,
		Value:           value,
	}
	_result.VariantContract.(*_Variant)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// VariantStringBuilder is a builder for VariantString
type VariantStringBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(value []PascalString) VariantStringBuilder
	// WithArrayLength adds ArrayLength (property field)
	WithOptionalArrayLength(int32) VariantStringBuilder
	// WithValue adds Value (property field)
	WithValue(...PascalString) VariantStringBuilder
	// Build builds the VariantString or returns an error if something is wrong
	Build() (VariantString, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() VariantString
}

// NewVariantStringBuilder() creates a VariantStringBuilder
func NewVariantStringBuilder() VariantStringBuilder {
	return &_VariantStringBuilder{_VariantString: new(_VariantString)}
}

type _VariantStringBuilder struct {
	*_VariantString

	parentBuilder *_VariantBuilder

	err *utils.MultiError
}

var _ (VariantStringBuilder) = (*_VariantStringBuilder)(nil)

func (b *_VariantStringBuilder) setParent(contract VariantContract) {
	b.VariantContract = contract
}

func (b *_VariantStringBuilder) WithMandatoryFields(value []PascalString) VariantStringBuilder {
	return b.WithValue(value...)
}

func (b *_VariantStringBuilder) WithOptionalArrayLength(arrayLength int32) VariantStringBuilder {
	b.ArrayLength = &arrayLength
	return b
}

func (b *_VariantStringBuilder) WithValue(value ...PascalString) VariantStringBuilder {
	b.Value = value
	return b
}

func (b *_VariantStringBuilder) Build() (VariantString, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._VariantString.deepCopy(), nil
}

func (b *_VariantStringBuilder) MustBuild() VariantString {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_VariantStringBuilder) Done() VariantBuilder {
	return b.parentBuilder
}

func (b *_VariantStringBuilder) buildForVariant() (Variant, error) {
	return b.Build()
}

func (b *_VariantStringBuilder) DeepCopy() any {
	_copy := b.CreateVariantStringBuilder().(*_VariantStringBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateVariantStringBuilder creates a VariantStringBuilder
func (b *_VariantString) CreateVariantStringBuilder() VariantStringBuilder {
	if b == nil {
		return NewVariantStringBuilder()
	}
	return &_VariantStringBuilder{_VariantString: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_VariantString) GetVariantType() uint8 {
	return uint8(12)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_VariantString) GetParent() VariantContract {
	return m.VariantContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_VariantString) GetArrayLength() *int32 {
	return m.ArrayLength
}

func (m *_VariantString) GetValue() []PascalString {
	return m.Value
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastVariantString(structType any) VariantString {
	if casted, ok := structType.(VariantString); ok {
		return casted
	}
	if casted, ok := structType.(*VariantString); ok {
		return *casted
	}
	return nil
}

func (m *_VariantString) GetTypeName() string {
	return "VariantString"
}

func (m *_VariantString) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.VariantContract.(*_Variant).GetLengthInBits(ctx))

	// Optional Field (arrayLength)
	if m.ArrayLength != nil {
		lengthInBits += 32
	}

	// Array field
	if len(m.Value) > 0 {
		for _curItem, element := range m.Value {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.Value), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_VariantString) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_VariantString) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_Variant, arrayLengthSpecified bool) (__variantString VariantString, err error) {
	m.VariantContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("VariantString"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for VariantString")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	var arrayLength *int32
	arrayLength, err = ReadOptionalField[int32](ctx, "arrayLength", ReadSignedInt(readBuffer, uint8(32)), arrayLengthSpecified)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'arrayLength' field"))
	}
	m.ArrayLength = arrayLength

	value, err := ReadCountArrayField[PascalString](ctx, "value", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer), uint64(utils.InlineIf(bool((arrayLength) == (nil)), func() any { return int32(int32(1)) }, func() any { return int32((*arrayLength)) }).(int32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'value' field"))
	}
	m.Value = value

	if closeErr := readBuffer.CloseContext("VariantString"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for VariantString")
	}

	return m, nil
}

func (m *_VariantString) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_VariantString) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("VariantString"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for VariantString")
		}

		if err := WriteOptionalField[int32](ctx, "arrayLength", m.GetArrayLength(), WriteSignedInt(writeBuffer, 32), true); err != nil {
			return errors.Wrap(err, "Error serializing 'arrayLength' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "value", m.GetValue(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'value' field")
		}

		if popErr := writeBuffer.PopContext("VariantString"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for VariantString")
		}
		return nil
	}
	return m.VariantContract.(*_Variant).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_VariantString) IsVariantString() {}

func (m *_VariantString) DeepCopy() any {
	return m.deepCopy()
}

func (m *_VariantString) deepCopy() *_VariantString {
	if m == nil {
		return nil
	}
	_VariantStringCopy := &_VariantString{
		m.VariantContract.(*_Variant).deepCopy(),
		utils.CopyPtr[int32](m.ArrayLength),
		utils.DeepCopySlice[PascalString, PascalString](m.Value),
	}
	m.VariantContract.(*_Variant)._SubType = m
	return _VariantStringCopy
}

func (m *_VariantString) String() string {
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
