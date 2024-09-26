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

// VariantExtensionObject is the corresponding interface of VariantExtensionObject
type VariantExtensionObject interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	Variant
	// GetArrayLength returns ArrayLength (property field)
	GetArrayLength() *int32
	// GetValue returns Value (property field)
	GetValue() []ExtensionObject
	// IsVariantExtensionObject is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsVariantExtensionObject()
	// CreateBuilder creates a VariantExtensionObjectBuilder
	CreateVariantExtensionObjectBuilder() VariantExtensionObjectBuilder
}

// _VariantExtensionObject is the data-structure of this message
type _VariantExtensionObject struct {
	VariantContract
	ArrayLength *int32
	Value       []ExtensionObject
}

var _ VariantExtensionObject = (*_VariantExtensionObject)(nil)
var _ VariantRequirements = (*_VariantExtensionObject)(nil)

// NewVariantExtensionObject factory function for _VariantExtensionObject
func NewVariantExtensionObject(arrayLengthSpecified bool, arrayDimensionsSpecified bool, noOfArrayDimensions *int32, arrayDimensions []bool, arrayLength *int32, value []ExtensionObject) *_VariantExtensionObject {
	_result := &_VariantExtensionObject{
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

// VariantExtensionObjectBuilder is a builder for VariantExtensionObject
type VariantExtensionObjectBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(value []ExtensionObject) VariantExtensionObjectBuilder
	// WithArrayLength adds ArrayLength (property field)
	WithOptionalArrayLength(int32) VariantExtensionObjectBuilder
	// WithValue adds Value (property field)
	WithValue(...ExtensionObject) VariantExtensionObjectBuilder
	// Build builds the VariantExtensionObject or returns an error if something is wrong
	Build() (VariantExtensionObject, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() VariantExtensionObject
}

// NewVariantExtensionObjectBuilder() creates a VariantExtensionObjectBuilder
func NewVariantExtensionObjectBuilder() VariantExtensionObjectBuilder {
	return &_VariantExtensionObjectBuilder{_VariantExtensionObject: new(_VariantExtensionObject)}
}

type _VariantExtensionObjectBuilder struct {
	*_VariantExtensionObject

	parentBuilder *_VariantBuilder

	err *utils.MultiError
}

var _ (VariantExtensionObjectBuilder) = (*_VariantExtensionObjectBuilder)(nil)

func (b *_VariantExtensionObjectBuilder) setParent(contract VariantContract) {
	b.VariantContract = contract
}

func (b *_VariantExtensionObjectBuilder) WithMandatoryFields(value []ExtensionObject) VariantExtensionObjectBuilder {
	return b.WithValue(value...)
}

func (b *_VariantExtensionObjectBuilder) WithOptionalArrayLength(arrayLength int32) VariantExtensionObjectBuilder {
	b.ArrayLength = &arrayLength
	return b
}

func (b *_VariantExtensionObjectBuilder) WithValue(value ...ExtensionObject) VariantExtensionObjectBuilder {
	b.Value = value
	return b
}

func (b *_VariantExtensionObjectBuilder) Build() (VariantExtensionObject, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._VariantExtensionObject.deepCopy(), nil
}

func (b *_VariantExtensionObjectBuilder) MustBuild() VariantExtensionObject {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_VariantExtensionObjectBuilder) Done() VariantBuilder {
	return b.parentBuilder
}

func (b *_VariantExtensionObjectBuilder) buildForVariant() (Variant, error) {
	return b.Build()
}

func (b *_VariantExtensionObjectBuilder) DeepCopy() any {
	_copy := b.CreateVariantExtensionObjectBuilder().(*_VariantExtensionObjectBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateVariantExtensionObjectBuilder creates a VariantExtensionObjectBuilder
func (b *_VariantExtensionObject) CreateVariantExtensionObjectBuilder() VariantExtensionObjectBuilder {
	if b == nil {
		return NewVariantExtensionObjectBuilder()
	}
	return &_VariantExtensionObjectBuilder{_VariantExtensionObject: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_VariantExtensionObject) GetVariantType() uint8 {
	return uint8(22)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_VariantExtensionObject) GetParent() VariantContract {
	return m.VariantContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_VariantExtensionObject) GetArrayLength() *int32 {
	return m.ArrayLength
}

func (m *_VariantExtensionObject) GetValue() []ExtensionObject {
	return m.Value
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastVariantExtensionObject(structType any) VariantExtensionObject {
	if casted, ok := structType.(VariantExtensionObject); ok {
		return casted
	}
	if casted, ok := structType.(*VariantExtensionObject); ok {
		return *casted
	}
	return nil
}

func (m *_VariantExtensionObject) GetTypeName() string {
	return "VariantExtensionObject"
}

func (m *_VariantExtensionObject) GetLengthInBits(ctx context.Context) uint16 {
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

func (m *_VariantExtensionObject) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_VariantExtensionObject) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_Variant, arrayLengthSpecified bool) (__variantExtensionObject VariantExtensionObject, err error) {
	m.VariantContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("VariantExtensionObject"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for VariantExtensionObject")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	var arrayLength *int32
	arrayLength, err = ReadOptionalField[int32](ctx, "arrayLength", ReadSignedInt(readBuffer, uint8(32)), arrayLengthSpecified)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'arrayLength' field"))
	}
	m.ArrayLength = arrayLength

	value, err := ReadCountArrayField[ExtensionObject](ctx, "value", ReadComplex[ExtensionObject](ExtensionObjectParseWithBufferProducer[ExtensionObject]((bool)(bool(true))), readBuffer), uint64(utils.InlineIf(bool((arrayLength) == (nil)), func() any { return int32(int32(1)) }, func() any { return int32((*arrayLength)) }).(int32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'value' field"))
	}
	m.Value = value

	if closeErr := readBuffer.CloseContext("VariantExtensionObject"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for VariantExtensionObject")
	}

	return m, nil
}

func (m *_VariantExtensionObject) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_VariantExtensionObject) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("VariantExtensionObject"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for VariantExtensionObject")
		}

		if err := WriteOptionalField[int32](ctx, "arrayLength", m.GetArrayLength(), WriteSignedInt(writeBuffer, 32), true); err != nil {
			return errors.Wrap(err, "Error serializing 'arrayLength' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "value", m.GetValue(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'value' field")
		}

		if popErr := writeBuffer.PopContext("VariantExtensionObject"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for VariantExtensionObject")
		}
		return nil
	}
	return m.VariantContract.(*_Variant).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_VariantExtensionObject) IsVariantExtensionObject() {}

func (m *_VariantExtensionObject) DeepCopy() any {
	return m.deepCopy()
}

func (m *_VariantExtensionObject) deepCopy() *_VariantExtensionObject {
	if m == nil {
		return nil
	}
	_VariantExtensionObjectCopy := &_VariantExtensionObject{
		m.VariantContract.(*_Variant).deepCopy(),
		utils.CopyPtr[int32](m.ArrayLength),
		utils.DeepCopySlice[ExtensionObject, ExtensionObject](m.Value),
	}
	m.VariantContract.(*_Variant)._SubType = m
	return _VariantExtensionObjectCopy
}

func (m *_VariantExtensionObject) String() string {
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
