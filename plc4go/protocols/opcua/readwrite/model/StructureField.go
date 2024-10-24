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

// StructureField is the corresponding interface of StructureField
type StructureField interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetName returns Name (property field)
	GetName() PascalString
	// GetDescription returns Description (property field)
	GetDescription() LocalizedText
	// GetDataType returns DataType (property field)
	GetDataType() NodeId
	// GetValueRank returns ValueRank (property field)
	GetValueRank() int32
	// GetArrayDimensions returns ArrayDimensions (property field)
	GetArrayDimensions() []uint32
	// GetMaxStringLength returns MaxStringLength (property field)
	GetMaxStringLength() uint32
	// GetIsOptional returns IsOptional (property field)
	GetIsOptional() bool
	// IsStructureField is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsStructureField()
	// CreateBuilder creates a StructureFieldBuilder
	CreateStructureFieldBuilder() StructureFieldBuilder
}

// _StructureField is the data-structure of this message
type _StructureField struct {
	ExtensionObjectDefinitionContract
	Name            PascalString
	Description     LocalizedText
	DataType        NodeId
	ValueRank       int32
	ArrayDimensions []uint32
	MaxStringLength uint32
	IsOptional      bool
	// Reserved Fields
	reservedField0 *uint8
}

var _ StructureField = (*_StructureField)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_StructureField)(nil)

// NewStructureField factory function for _StructureField
func NewStructureField(name PascalString, description LocalizedText, dataType NodeId, valueRank int32, arrayDimensions []uint32, maxStringLength uint32, isOptional bool) *_StructureField {
	if name == nil {
		panic("name of type PascalString for StructureField must not be nil")
	}
	if description == nil {
		panic("description of type LocalizedText for StructureField must not be nil")
	}
	if dataType == nil {
		panic("dataType of type NodeId for StructureField must not be nil")
	}
	_result := &_StructureField{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		Name:                              name,
		Description:                       description,
		DataType:                          dataType,
		ValueRank:                         valueRank,
		ArrayDimensions:                   arrayDimensions,
		MaxStringLength:                   maxStringLength,
		IsOptional:                        isOptional,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// StructureFieldBuilder is a builder for StructureField
type StructureFieldBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(name PascalString, description LocalizedText, dataType NodeId, valueRank int32, arrayDimensions []uint32, maxStringLength uint32, isOptional bool) StructureFieldBuilder
	// WithName adds Name (property field)
	WithName(PascalString) StructureFieldBuilder
	// WithNameBuilder adds Name (property field) which is build by the builder
	WithNameBuilder(func(PascalStringBuilder) PascalStringBuilder) StructureFieldBuilder
	// WithDescription adds Description (property field)
	WithDescription(LocalizedText) StructureFieldBuilder
	// WithDescriptionBuilder adds Description (property field) which is build by the builder
	WithDescriptionBuilder(func(LocalizedTextBuilder) LocalizedTextBuilder) StructureFieldBuilder
	// WithDataType adds DataType (property field)
	WithDataType(NodeId) StructureFieldBuilder
	// WithDataTypeBuilder adds DataType (property field) which is build by the builder
	WithDataTypeBuilder(func(NodeIdBuilder) NodeIdBuilder) StructureFieldBuilder
	// WithValueRank adds ValueRank (property field)
	WithValueRank(int32) StructureFieldBuilder
	// WithArrayDimensions adds ArrayDimensions (property field)
	WithArrayDimensions(...uint32) StructureFieldBuilder
	// WithMaxStringLength adds MaxStringLength (property field)
	WithMaxStringLength(uint32) StructureFieldBuilder
	// WithIsOptional adds IsOptional (property field)
	WithIsOptional(bool) StructureFieldBuilder
	// Build builds the StructureField or returns an error if something is wrong
	Build() (StructureField, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() StructureField
}

// NewStructureFieldBuilder() creates a StructureFieldBuilder
func NewStructureFieldBuilder() StructureFieldBuilder {
	return &_StructureFieldBuilder{_StructureField: new(_StructureField)}
}

type _StructureFieldBuilder struct {
	*_StructureField

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (StructureFieldBuilder) = (*_StructureFieldBuilder)(nil)

func (b *_StructureFieldBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_StructureFieldBuilder) WithMandatoryFields(name PascalString, description LocalizedText, dataType NodeId, valueRank int32, arrayDimensions []uint32, maxStringLength uint32, isOptional bool) StructureFieldBuilder {
	return b.WithName(name).WithDescription(description).WithDataType(dataType).WithValueRank(valueRank).WithArrayDimensions(arrayDimensions...).WithMaxStringLength(maxStringLength).WithIsOptional(isOptional)
}

func (b *_StructureFieldBuilder) WithName(name PascalString) StructureFieldBuilder {
	b.Name = name
	return b
}

func (b *_StructureFieldBuilder) WithNameBuilder(builderSupplier func(PascalStringBuilder) PascalStringBuilder) StructureFieldBuilder {
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

func (b *_StructureFieldBuilder) WithDescription(description LocalizedText) StructureFieldBuilder {
	b.Description = description
	return b
}

func (b *_StructureFieldBuilder) WithDescriptionBuilder(builderSupplier func(LocalizedTextBuilder) LocalizedTextBuilder) StructureFieldBuilder {
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

func (b *_StructureFieldBuilder) WithDataType(dataType NodeId) StructureFieldBuilder {
	b.DataType = dataType
	return b
}

func (b *_StructureFieldBuilder) WithDataTypeBuilder(builderSupplier func(NodeIdBuilder) NodeIdBuilder) StructureFieldBuilder {
	builder := builderSupplier(b.DataType.CreateNodeIdBuilder())
	var err error
	b.DataType, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "NodeIdBuilder failed"))
	}
	return b
}

func (b *_StructureFieldBuilder) WithValueRank(valueRank int32) StructureFieldBuilder {
	b.ValueRank = valueRank
	return b
}

func (b *_StructureFieldBuilder) WithArrayDimensions(arrayDimensions ...uint32) StructureFieldBuilder {
	b.ArrayDimensions = arrayDimensions
	return b
}

func (b *_StructureFieldBuilder) WithMaxStringLength(maxStringLength uint32) StructureFieldBuilder {
	b.MaxStringLength = maxStringLength
	return b
}

func (b *_StructureFieldBuilder) WithIsOptional(isOptional bool) StructureFieldBuilder {
	b.IsOptional = isOptional
	return b
}

func (b *_StructureFieldBuilder) Build() (StructureField, error) {
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
	if b.DataType == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'dataType' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._StructureField.deepCopy(), nil
}

func (b *_StructureFieldBuilder) MustBuild() StructureField {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_StructureFieldBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_StructureFieldBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_StructureFieldBuilder) DeepCopy() any {
	_copy := b.CreateStructureFieldBuilder().(*_StructureFieldBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateStructureFieldBuilder creates a StructureFieldBuilder
func (b *_StructureField) CreateStructureFieldBuilder() StructureFieldBuilder {
	if b == nil {
		return NewStructureFieldBuilder()
	}
	return &_StructureFieldBuilder{_StructureField: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_StructureField) GetExtensionId() int32 {
	return int32(103)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_StructureField) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_StructureField) GetName() PascalString {
	return m.Name
}

func (m *_StructureField) GetDescription() LocalizedText {
	return m.Description
}

func (m *_StructureField) GetDataType() NodeId {
	return m.DataType
}

func (m *_StructureField) GetValueRank() int32 {
	return m.ValueRank
}

func (m *_StructureField) GetArrayDimensions() []uint32 {
	return m.ArrayDimensions
}

func (m *_StructureField) GetMaxStringLength() uint32 {
	return m.MaxStringLength
}

func (m *_StructureField) GetIsOptional() bool {
	return m.IsOptional
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastStructureField(structType any) StructureField {
	if casted, ok := structType.(StructureField); ok {
		return casted
	}
	if casted, ok := structType.(*StructureField); ok {
		return *casted
	}
	return nil
}

func (m *_StructureField) GetTypeName() string {
	return "StructureField"
}

func (m *_StructureField) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (name)
	lengthInBits += m.Name.GetLengthInBits(ctx)

	// Simple field (description)
	lengthInBits += m.Description.GetLengthInBits(ctx)

	// Simple field (dataType)
	lengthInBits += m.DataType.GetLengthInBits(ctx)

	// Simple field (valueRank)
	lengthInBits += 32

	// Implicit Field (noOfArrayDimensions)
	lengthInBits += 32

	// Array field
	if len(m.ArrayDimensions) > 0 {
		lengthInBits += 32 * uint16(len(m.ArrayDimensions))
	}

	// Simple field (maxStringLength)
	lengthInBits += 32

	// Reserved Field (reserved)
	lengthInBits += 7

	// Simple field (isOptional)
	lengthInBits += 1

	return lengthInBits
}

func (m *_StructureField) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_StructureField) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__structureField StructureField, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("StructureField"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for StructureField")
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

	dataType, err := ReadSimpleField[NodeId](ctx, "dataType", ReadComplex[NodeId](NodeIdParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'dataType' field"))
	}
	m.DataType = dataType

	valueRank, err := ReadSimpleField(ctx, "valueRank", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'valueRank' field"))
	}
	m.ValueRank = valueRank

	noOfArrayDimensions, err := ReadImplicitField[int32](ctx, "noOfArrayDimensions", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfArrayDimensions' field"))
	}
	_ = noOfArrayDimensions

	arrayDimensions, err := ReadCountArrayField[uint32](ctx, "arrayDimensions", ReadUnsignedInt(readBuffer, uint8(32)), uint64(noOfArrayDimensions))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'arrayDimensions' field"))
	}
	m.ArrayDimensions = arrayDimensions

	maxStringLength, err := ReadSimpleField(ctx, "maxStringLength", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'maxStringLength' field"))
	}
	m.MaxStringLength = maxStringLength

	reservedField0, err := ReadReservedField(ctx, "reserved", ReadUnsignedByte(readBuffer, uint8(7)), uint8(0x00))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}
	m.reservedField0 = reservedField0

	isOptional, err := ReadSimpleField(ctx, "isOptional", ReadBoolean(readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'isOptional' field"))
	}
	m.IsOptional = isOptional

	if closeErr := readBuffer.CloseContext("StructureField"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for StructureField")
	}

	return m, nil
}

func (m *_StructureField) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_StructureField) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("StructureField"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for StructureField")
		}

		if err := WriteSimpleField[PascalString](ctx, "name", m.GetName(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'name' field")
		}

		if err := WriteSimpleField[LocalizedText](ctx, "description", m.GetDescription(), WriteComplex[LocalizedText](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'description' field")
		}

		if err := WriteSimpleField[NodeId](ctx, "dataType", m.GetDataType(), WriteComplex[NodeId](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'dataType' field")
		}

		if err := WriteSimpleField[int32](ctx, "valueRank", m.GetValueRank(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'valueRank' field")
		}
		noOfArrayDimensions := int32(utils.InlineIf(bool((m.GetArrayDimensions()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetArrayDimensions()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfArrayDimensions", noOfArrayDimensions, WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfArrayDimensions' field")
		}

		if err := WriteSimpleTypeArrayField(ctx, "arrayDimensions", m.GetArrayDimensions(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'arrayDimensions' field")
		}

		if err := WriteSimpleField[uint32](ctx, "maxStringLength", m.GetMaxStringLength(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'maxStringLength' field")
		}

		if err := WriteReservedField[uint8](ctx, "reserved", uint8(0x00), WriteUnsignedByte(writeBuffer, 7)); err != nil {
			return errors.Wrap(err, "Error serializing 'reserved' field number 1")
		}

		if err := WriteSimpleField[bool](ctx, "isOptional", m.GetIsOptional(), WriteBoolean(writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'isOptional' field")
		}

		if popErr := writeBuffer.PopContext("StructureField"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for StructureField")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_StructureField) IsStructureField() {}

func (m *_StructureField) DeepCopy() any {
	return m.deepCopy()
}

func (m *_StructureField) deepCopy() *_StructureField {
	if m == nil {
		return nil
	}
	_StructureFieldCopy := &_StructureField{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.Name.DeepCopy().(PascalString),
		m.Description.DeepCopy().(LocalizedText),
		m.DataType.DeepCopy().(NodeId),
		m.ValueRank,
		utils.DeepCopySlice[uint32, uint32](m.ArrayDimensions),
		m.MaxStringLength,
		m.IsOptional,
		m.reservedField0,
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _StructureFieldCopy
}

func (m *_StructureField) String() string {
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
