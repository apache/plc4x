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

// ComplexNumberType is the corresponding interface of ComplexNumberType
type ComplexNumberType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetReal returns Real (property field)
	GetReal() float32
	// GetImaginary returns Imaginary (property field)
	GetImaginary() float32
	// IsComplexNumberType is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsComplexNumberType()
	// CreateBuilder creates a ComplexNumberTypeBuilder
	CreateComplexNumberTypeBuilder() ComplexNumberTypeBuilder
}

// _ComplexNumberType is the data-structure of this message
type _ComplexNumberType struct {
	ExtensionObjectDefinitionContract
	Real      float32
	Imaginary float32
}

var _ ComplexNumberType = (*_ComplexNumberType)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_ComplexNumberType)(nil)

// NewComplexNumberType factory function for _ComplexNumberType
func NewComplexNumberType(real float32, imaginary float32) *_ComplexNumberType {
	_result := &_ComplexNumberType{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		Real:                              real,
		Imaginary:                         imaginary,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ComplexNumberTypeBuilder is a builder for ComplexNumberType
type ComplexNumberTypeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(real float32, imaginary float32) ComplexNumberTypeBuilder
	// WithReal adds Real (property field)
	WithReal(float32) ComplexNumberTypeBuilder
	// WithImaginary adds Imaginary (property field)
	WithImaginary(float32) ComplexNumberTypeBuilder
	// Build builds the ComplexNumberType or returns an error if something is wrong
	Build() (ComplexNumberType, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ComplexNumberType
}

// NewComplexNumberTypeBuilder() creates a ComplexNumberTypeBuilder
func NewComplexNumberTypeBuilder() ComplexNumberTypeBuilder {
	return &_ComplexNumberTypeBuilder{_ComplexNumberType: new(_ComplexNumberType)}
}

type _ComplexNumberTypeBuilder struct {
	*_ComplexNumberType

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (ComplexNumberTypeBuilder) = (*_ComplexNumberTypeBuilder)(nil)

func (b *_ComplexNumberTypeBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_ComplexNumberTypeBuilder) WithMandatoryFields(real float32, imaginary float32) ComplexNumberTypeBuilder {
	return b.WithReal(real).WithImaginary(imaginary)
}

func (b *_ComplexNumberTypeBuilder) WithReal(real float32) ComplexNumberTypeBuilder {
	b.Real = real
	return b
}

func (b *_ComplexNumberTypeBuilder) WithImaginary(imaginary float32) ComplexNumberTypeBuilder {
	b.Imaginary = imaginary
	return b
}

func (b *_ComplexNumberTypeBuilder) Build() (ComplexNumberType, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ComplexNumberType.deepCopy(), nil
}

func (b *_ComplexNumberTypeBuilder) MustBuild() ComplexNumberType {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ComplexNumberTypeBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_ComplexNumberTypeBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_ComplexNumberTypeBuilder) DeepCopy() any {
	_copy := b.CreateComplexNumberTypeBuilder().(*_ComplexNumberTypeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateComplexNumberTypeBuilder creates a ComplexNumberTypeBuilder
func (b *_ComplexNumberType) CreateComplexNumberTypeBuilder() ComplexNumberTypeBuilder {
	if b == nil {
		return NewComplexNumberTypeBuilder()
	}
	return &_ComplexNumberTypeBuilder{_ComplexNumberType: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ComplexNumberType) GetExtensionId() int32 {
	return int32(12173)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ComplexNumberType) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ComplexNumberType) GetReal() float32 {
	return m.Real
}

func (m *_ComplexNumberType) GetImaginary() float32 {
	return m.Imaginary
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastComplexNumberType(structType any) ComplexNumberType {
	if casted, ok := structType.(ComplexNumberType); ok {
		return casted
	}
	if casted, ok := structType.(*ComplexNumberType); ok {
		return *casted
	}
	return nil
}

func (m *_ComplexNumberType) GetTypeName() string {
	return "ComplexNumberType"
}

func (m *_ComplexNumberType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (real)
	lengthInBits += 32

	// Simple field (imaginary)
	lengthInBits += 32

	return lengthInBits
}

func (m *_ComplexNumberType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ComplexNumberType) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__complexNumberType ComplexNumberType, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ComplexNumberType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ComplexNumberType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	real, err := ReadSimpleField(ctx, "real", ReadFloat(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'real' field"))
	}
	m.Real = real

	imaginary, err := ReadSimpleField(ctx, "imaginary", ReadFloat(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'imaginary' field"))
	}
	m.Imaginary = imaginary

	if closeErr := readBuffer.CloseContext("ComplexNumberType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ComplexNumberType")
	}

	return m, nil
}

func (m *_ComplexNumberType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ComplexNumberType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ComplexNumberType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ComplexNumberType")
		}

		if err := WriteSimpleField[float32](ctx, "real", m.GetReal(), WriteFloat(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'real' field")
		}

		if err := WriteSimpleField[float32](ctx, "imaginary", m.GetImaginary(), WriteFloat(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'imaginary' field")
		}

		if popErr := writeBuffer.PopContext("ComplexNumberType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ComplexNumberType")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ComplexNumberType) IsComplexNumberType() {}

func (m *_ComplexNumberType) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ComplexNumberType) deepCopy() *_ComplexNumberType {
	if m == nil {
		return nil
	}
	_ComplexNumberTypeCopy := &_ComplexNumberType{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.Real,
		m.Imaginary,
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _ComplexNumberTypeCopy
}

func (m *_ComplexNumberType) String() string {
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
