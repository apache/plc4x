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

// BitFieldDefinition is the corresponding interface of BitFieldDefinition
type BitFieldDefinition interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetName returns Name (property field)
	GetName() PascalString
	// GetDescription returns Description (property field)
	GetDescription() LocalizedText
	// GetStartingBitPosition returns StartingBitPosition (property field)
	GetStartingBitPosition() uint32
	// GetEndingBitPosition returns EndingBitPosition (property field)
	GetEndingBitPosition() uint32
	// IsBitFieldDefinition is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBitFieldDefinition()
	// CreateBuilder creates a BitFieldDefinitionBuilder
	CreateBitFieldDefinitionBuilder() BitFieldDefinitionBuilder
}

// _BitFieldDefinition is the data-structure of this message
type _BitFieldDefinition struct {
	ExtensionObjectDefinitionContract
	Name                PascalString
	Description         LocalizedText
	StartingBitPosition uint32
	EndingBitPosition   uint32
	// Reserved Fields
	reservedField0 *uint8
	reservedField1 *bool
}

var _ BitFieldDefinition = (*_BitFieldDefinition)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_BitFieldDefinition)(nil)

// NewBitFieldDefinition factory function for _BitFieldDefinition
func NewBitFieldDefinition(name PascalString, description LocalizedText, startingBitPosition uint32, endingBitPosition uint32) *_BitFieldDefinition {
	if name == nil {
		panic("name of type PascalString for BitFieldDefinition must not be nil")
	}
	if description == nil {
		panic("description of type LocalizedText for BitFieldDefinition must not be nil")
	}
	_result := &_BitFieldDefinition{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		Name:                              name,
		Description:                       description,
		StartingBitPosition:               startingBitPosition,
		EndingBitPosition:                 endingBitPosition,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BitFieldDefinitionBuilder is a builder for BitFieldDefinition
type BitFieldDefinitionBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(name PascalString, description LocalizedText, startingBitPosition uint32, endingBitPosition uint32) BitFieldDefinitionBuilder
	// WithName adds Name (property field)
	WithName(PascalString) BitFieldDefinitionBuilder
	// WithNameBuilder adds Name (property field) which is build by the builder
	WithNameBuilder(func(PascalStringBuilder) PascalStringBuilder) BitFieldDefinitionBuilder
	// WithDescription adds Description (property field)
	WithDescription(LocalizedText) BitFieldDefinitionBuilder
	// WithDescriptionBuilder adds Description (property field) which is build by the builder
	WithDescriptionBuilder(func(LocalizedTextBuilder) LocalizedTextBuilder) BitFieldDefinitionBuilder
	// WithStartingBitPosition adds StartingBitPosition (property field)
	WithStartingBitPosition(uint32) BitFieldDefinitionBuilder
	// WithEndingBitPosition adds EndingBitPosition (property field)
	WithEndingBitPosition(uint32) BitFieldDefinitionBuilder
	// Build builds the BitFieldDefinition or returns an error if something is wrong
	Build() (BitFieldDefinition, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BitFieldDefinition
}

// NewBitFieldDefinitionBuilder() creates a BitFieldDefinitionBuilder
func NewBitFieldDefinitionBuilder() BitFieldDefinitionBuilder {
	return &_BitFieldDefinitionBuilder{_BitFieldDefinition: new(_BitFieldDefinition)}
}

type _BitFieldDefinitionBuilder struct {
	*_BitFieldDefinition

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (BitFieldDefinitionBuilder) = (*_BitFieldDefinitionBuilder)(nil)

func (b *_BitFieldDefinitionBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_BitFieldDefinitionBuilder) WithMandatoryFields(name PascalString, description LocalizedText, startingBitPosition uint32, endingBitPosition uint32) BitFieldDefinitionBuilder {
	return b.WithName(name).WithDescription(description).WithStartingBitPosition(startingBitPosition).WithEndingBitPosition(endingBitPosition)
}

func (b *_BitFieldDefinitionBuilder) WithName(name PascalString) BitFieldDefinitionBuilder {
	b.Name = name
	return b
}

func (b *_BitFieldDefinitionBuilder) WithNameBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) BitFieldDefinitionBuilder {
	builder := builderSupplier(b.Name.CreatePascalStringBuilder())
	var err error
	b.Name, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "PascalStringBuilder failed"))
	}
	return b
}

func (b *_BitFieldDefinitionBuilder) WithDescription(description LocalizedText) BitFieldDefinitionBuilder {
	b.Description = description
	return b
}

func (b *_BitFieldDefinitionBuilder) WithDescriptionBuilder(builderSupplier func(LocalizedTextBuilder) LocalizedTextBuilder) BitFieldDefinitionBuilder {
	builder := builderSupplier(b.Description.CreateLocalizedTextBuilder())
	var err error
	b.Description, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "LocalizedTextBuilder failed"))
	}
	return b
}

func (b *_BitFieldDefinitionBuilder) WithStartingBitPosition(startingBitPosition uint32) BitFieldDefinitionBuilder {
	b.StartingBitPosition = startingBitPosition
	return b
}

func (b *_BitFieldDefinitionBuilder) WithEndingBitPosition(endingBitPosition uint32) BitFieldDefinitionBuilder {
	b.EndingBitPosition = endingBitPosition
	return b
}

func (b *_BitFieldDefinitionBuilder) Build() (BitFieldDefinition, error) {
	if b.Name == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'name' not set"))
	}
	if b.Description == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'description' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BitFieldDefinition.deepCopy(), nil
}

func (b *_BitFieldDefinitionBuilder) MustBuild() BitFieldDefinition {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BitFieldDefinitionBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_BitFieldDefinitionBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_BitFieldDefinitionBuilder) DeepCopy() any {
	_copy := b.CreateBitFieldDefinitionBuilder().(*_BitFieldDefinitionBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBitFieldDefinitionBuilder creates a BitFieldDefinitionBuilder
func (b *_BitFieldDefinition) CreateBitFieldDefinitionBuilder() BitFieldDefinitionBuilder {
	if b == nil {
		return NewBitFieldDefinitionBuilder()
	}
	return &_BitFieldDefinitionBuilder{_BitFieldDefinition: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BitFieldDefinition) GetExtensionId() int32 {
	return int32(32423)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BitFieldDefinition) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BitFieldDefinition) GetName() PascalString {
	return m.Name
}

func (m *_BitFieldDefinition) GetDescription() LocalizedText {
	return m.Description
}

func (m *_BitFieldDefinition) GetStartingBitPosition() uint32 {
	return m.StartingBitPosition
}

func (m *_BitFieldDefinition) GetEndingBitPosition() uint32 {
	return m.EndingBitPosition
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBitFieldDefinition(structType any) BitFieldDefinition {
	if casted, ok := structType.(BitFieldDefinition); ok {
		return casted
	}
	if casted, ok := structType.(*BitFieldDefinition); ok {
		return *casted
	}
	return nil
}

func (m *_BitFieldDefinition) GetTypeName() string {
	return "BitFieldDefinition"
}

func (m *_BitFieldDefinition) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (name)
	lengthInBits += m.Name.GetLengthInBits(ctx)

	// Simple field (description)
	lengthInBits += m.Description.GetLengthInBits(ctx)

	// Reserved Field (reserved)
	lengthInBits += 7

	// Reserved Field (reserved)
	lengthInBits += 1

	// Simple field (startingBitPosition)
	lengthInBits += 32

	// Simple field (endingBitPosition)
	lengthInBits += 32

	return lengthInBits
}

func (m *_BitFieldDefinition) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BitFieldDefinition) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__bitFieldDefinition BitFieldDefinition, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BitFieldDefinition"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BitFieldDefinition")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	name, err := ReadSimpleField[PascalString](ctx, "name", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'name' field"))
	}
	m.Name = name

	description, err := ReadSimpleField[LocalizedText](ctx, "description", ReadComplex[LocalizedText](LocalizedTextParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'description' field"))
	}
	m.Description = description

	reservedField0, err := ReadReservedField(ctx, "reserved", ReadUnsignedByte(readBuffer, uint8(7)), uint8(0x00))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}
	m.reservedField0 = reservedField0

	reservedField1, err := ReadReservedField(ctx, "reserved", ReadBoolean(readBuffer), bool(false))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}
	m.reservedField1 = reservedField1

	startingBitPosition, err := ReadSimpleField(ctx, "startingBitPosition", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'startingBitPosition' field"))
	}
	m.StartingBitPosition = startingBitPosition

	endingBitPosition, err := ReadSimpleField(ctx, "endingBitPosition", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'endingBitPosition' field"))
	}
	m.EndingBitPosition = endingBitPosition

	if closeErr := readBuffer.CloseContext("BitFieldDefinition"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BitFieldDefinition")
	}

	return m, nil
}

func (m *_BitFieldDefinition) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BitFieldDefinition) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BitFieldDefinition"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BitFieldDefinition")
		}

		if err := WriteSimpleField[PascalString](ctx, "name", m.GetName(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'name' field")
		}

		if err := WriteSimpleField[LocalizedText](ctx, "description", m.GetDescription(), WriteComplex[LocalizedText](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'description' field")
		}

		if err := WriteReservedField[uint8](ctx, "reserved", uint8(0x00), WriteUnsignedByte(writeBuffer, 7)); err != nil {
			return errors.Wrap(err, "Error serializing 'reserved' field number 1")
		}

		if err := WriteReservedField[bool](ctx, "reserved", bool(false), WriteBoolean(writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'reserved' field number 2")
		}

		if err := WriteSimpleField[uint32](ctx, "startingBitPosition", m.GetStartingBitPosition(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'startingBitPosition' field")
		}

		if err := WriteSimpleField[uint32](ctx, "endingBitPosition", m.GetEndingBitPosition(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'endingBitPosition' field")
		}

		if popErr := writeBuffer.PopContext("BitFieldDefinition"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BitFieldDefinition")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BitFieldDefinition) IsBitFieldDefinition() {}

func (m *_BitFieldDefinition) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BitFieldDefinition) deepCopy() *_BitFieldDefinition {
	if m == nil {
		return nil
	}
	_BitFieldDefinitionCopy := &_BitFieldDefinition{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.Name.DeepCopy().(PascalString),
		m.Description.DeepCopy().(LocalizedText),
		m.StartingBitPosition,
		m.EndingBitPosition,
		m.reservedField0,
		m.reservedField1,
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _BitFieldDefinitionCopy
}

func (m *_BitFieldDefinition) String() string {
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
