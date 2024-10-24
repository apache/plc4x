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

// Constant values.
const CycServiceItemType_FUNCTIONID uint8 = 0x12

// CycServiceItemType is the corresponding interface of CycServiceItemType
type CycServiceItemType interface {
	CycServiceItemTypeContract
	CycServiceItemTypeRequirements
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	// IsCycServiceItemType is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCycServiceItemType()
	// CreateBuilder creates a CycServiceItemTypeBuilder
	CreateCycServiceItemTypeBuilder() CycServiceItemTypeBuilder
}

// CycServiceItemTypeContract provides a set of functions which can be overwritten by a sub struct
type CycServiceItemTypeContract interface {
	// GetByteLength returns ByteLength (property field)
	GetByteLength() uint8
	// GetSyntaxId returns SyntaxId (property field)
	GetSyntaxId() uint8
	// IsCycServiceItemType is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsCycServiceItemType()
	// CreateBuilder creates a CycServiceItemTypeBuilder
	CreateCycServiceItemTypeBuilder() CycServiceItemTypeBuilder
}

// CycServiceItemTypeRequirements provides a set of functions which need to be implemented by a sub struct
type CycServiceItemTypeRequirements interface {
	GetLengthInBits(ctx context.Context) uint16
	GetLengthInBytes(ctx context.Context) uint16
	// GetSyntaxId returns SyntaxId (discriminator field)
	GetSyntaxId() uint8
}

// _CycServiceItemType is the data-structure of this message
type _CycServiceItemType struct {
	_SubType   CycServiceItemType
	ByteLength uint8
	SyntaxId   uint8
}

var _ CycServiceItemTypeContract = (*_CycServiceItemType)(nil)

// NewCycServiceItemType factory function for _CycServiceItemType
func NewCycServiceItemType(byteLength uint8, syntaxId uint8) *_CycServiceItemType {
	return &_CycServiceItemType{ByteLength: byteLength, SyntaxId: syntaxId}
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// CycServiceItemTypeBuilder is a builder for CycServiceItemType
type CycServiceItemTypeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(byteLength uint8, syntaxId uint8) CycServiceItemTypeBuilder
	// WithByteLength adds ByteLength (property field)
	WithByteLength(uint8) CycServiceItemTypeBuilder
	// WithSyntaxId adds SyntaxId (property field)
	WithSyntaxId(uint8) CycServiceItemTypeBuilder
	// AsCycServiceItemAnyType converts this build to a subType of CycServiceItemType. It is always possible to return to current builder using Done()
	AsCycServiceItemAnyType() interface {
		CycServiceItemAnyTypeBuilder
		Done() CycServiceItemTypeBuilder
	}
	// AsCycServiceItemDbReadType converts this build to a subType of CycServiceItemType. It is always possible to return to current builder using Done()
	AsCycServiceItemDbReadType() interface {
		CycServiceItemDbReadTypeBuilder
		Done() CycServiceItemTypeBuilder
	}
	// Build builds the CycServiceItemType or returns an error if something is wrong
	PartialBuild() (CycServiceItemTypeContract, error)
	// MustBuild does the same as Build but panics on error
	PartialMustBuild() CycServiceItemTypeContract
	// Build builds the CycServiceItemType or returns an error if something is wrong
	Build() (CycServiceItemType, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() CycServiceItemType
}

// NewCycServiceItemTypeBuilder() creates a CycServiceItemTypeBuilder
func NewCycServiceItemTypeBuilder() CycServiceItemTypeBuilder {
	return &_CycServiceItemTypeBuilder{_CycServiceItemType: new(_CycServiceItemType)}
}

type _CycServiceItemTypeChildBuilder interface {
	utils.Copyable
	setParent(CycServiceItemTypeContract)
	buildForCycServiceItemType() (CycServiceItemType, error)
}

type _CycServiceItemTypeBuilder struct {
	*_CycServiceItemType

	childBuilder _CycServiceItemTypeChildBuilder

	err *utils.MultiError
}

var _ (CycServiceItemTypeBuilder) = (*_CycServiceItemTypeBuilder)(nil)

func (b *_CycServiceItemTypeBuilder) WithMandatoryFields(byteLength uint8, syntaxId uint8) CycServiceItemTypeBuilder {
	return b.WithByteLength(byteLength).WithSyntaxId(syntaxId)
}

func (b *_CycServiceItemTypeBuilder) WithByteLength(byteLength uint8) CycServiceItemTypeBuilder {
	b.ByteLength = byteLength
	return b
}

func (b *_CycServiceItemTypeBuilder) WithSyntaxId(syntaxId uint8) CycServiceItemTypeBuilder {
	b.SyntaxId = syntaxId
	return b
}

func (b *_CycServiceItemTypeBuilder) PartialBuild() (CycServiceItemTypeContract, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._CycServiceItemType.deepCopy(), nil
}

func (b *_CycServiceItemTypeBuilder) PartialMustBuild() CycServiceItemTypeContract {
	build, err := b.PartialBuild()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_CycServiceItemTypeBuilder) AsCycServiceItemAnyType() interface {
	CycServiceItemAnyTypeBuilder
	Done() CycServiceItemTypeBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		CycServiceItemAnyTypeBuilder
		Done() CycServiceItemTypeBuilder
	}); ok {
		return cb
	}
	cb := NewCycServiceItemAnyTypeBuilder().(*_CycServiceItemAnyTypeBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_CycServiceItemTypeBuilder) AsCycServiceItemDbReadType() interface {
	CycServiceItemDbReadTypeBuilder
	Done() CycServiceItemTypeBuilder
} {
	if cb, ok := b.childBuilder.(interface {
		CycServiceItemDbReadTypeBuilder
		Done() CycServiceItemTypeBuilder
	}); ok {
		return cb
	}
	cb := NewCycServiceItemDbReadTypeBuilder().(*_CycServiceItemDbReadTypeBuilder)
	cb.parentBuilder = b
	b.childBuilder = cb
	return cb
}

func (b *_CycServiceItemTypeBuilder) Build() (CycServiceItemType, error) {
	v, err := b.PartialBuild()
	if err != nil {
		return nil, errors.Wrap(err, "error occurred during partial build")
	}
	if b.childBuilder == nil {
		return nil, errors.New("no child builder present")
	}
	b.childBuilder.setParent(v)
	return b.childBuilder.buildForCycServiceItemType()
}

func (b *_CycServiceItemTypeBuilder) MustBuild() CycServiceItemType {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_CycServiceItemTypeBuilder) DeepCopy() any {
	_copy := b.CreateCycServiceItemTypeBuilder().(*_CycServiceItemTypeBuilder)
	_copy.childBuilder = b.childBuilder.DeepCopy().(_CycServiceItemTypeChildBuilder)
	_copy.childBuilder.setParent(_copy)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateCycServiceItemTypeBuilder creates a CycServiceItemTypeBuilder
func (b *_CycServiceItemType) CreateCycServiceItemTypeBuilder() CycServiceItemTypeBuilder {
	if b == nil {
		return NewCycServiceItemTypeBuilder()
	}
	return &_CycServiceItemTypeBuilder{_CycServiceItemType: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CycServiceItemType) GetByteLength() uint8 {
	return m.ByteLength
}

func (m *_CycServiceItemType) GetSyntaxId() uint8 {
	return m.SyntaxId
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for const fields.
///////////////////////

func (m *_CycServiceItemType) GetFunctionId() uint8 {
	return CycServiceItemType_FUNCTIONID
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastCycServiceItemType(structType any) CycServiceItemType {
	if casted, ok := structType.(CycServiceItemType); ok {
		return casted
	}
	if casted, ok := structType.(*CycServiceItemType); ok {
		return *casted
	}
	return nil
}

func (m *_CycServiceItemType) GetTypeName() string {
	return "CycServiceItemType"
}

func (m *_CycServiceItemType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Const Field (functionId)
	lengthInBits += 8

	// Simple field (byteLength)
	lengthInBits += 8

	// Simple field (syntaxId)
	lengthInBits += 8

	return lengthInBits
}

func (m *_CycServiceItemType) GetLengthInBytes(ctx context.Context) uint16 {
	return m._SubType.GetLengthInBits(ctx) / 8
}

func CycServiceItemTypeParse[T CycServiceItemType](ctx context.Context, theBytes []byte) (T, error) {
	return CycServiceItemTypeParseWithBuffer[T](ctx, utils.NewReadBufferByteBased(theBytes))
}

func CycServiceItemTypeParseWithBufferProducer[T CycServiceItemType]() func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
		v, err := CycServiceItemTypeParseWithBuffer[T](ctx, readBuffer)
		if err != nil {
			var zero T
			return zero, err
		}
		return v, nil
	}
}

func CycServiceItemTypeParseWithBuffer[T CycServiceItemType](ctx context.Context, readBuffer utils.ReadBuffer) (T, error) {
	v, err := (&_CycServiceItemType{}).parse(ctx, readBuffer)
	if err != nil {
		var zero T
		return zero, err
	}
	vc, ok := v.(T)
	if !ok {
		var zero T
		return zero, errors.Errorf("Unexpected type %T. Expected type %T", v, *new(T))
	}
	return vc, nil
}

func (m *_CycServiceItemType) parse(ctx context.Context, readBuffer utils.ReadBuffer) (__cycServiceItemType CycServiceItemType, err error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CycServiceItemType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CycServiceItemType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	functionId, err := ReadConstField[uint8](ctx, "functionId", ReadUnsignedByte(readBuffer, uint8(8)), CycServiceItemType_FUNCTIONID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'functionId' field"))
	}
	_ = functionId

	byteLength, err := ReadSimpleField(ctx, "byteLength", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'byteLength' field"))
	}
	m.ByteLength = byteLength

	syntaxId, err := ReadSimpleField(ctx, "syntaxId", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'syntaxId' field"))
	}
	m.SyntaxId = syntaxId

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	var _child CycServiceItemType
	switch {
	case syntaxId == 0x10: // CycServiceItemAnyType
		if _child, err = new(_CycServiceItemAnyType).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type CycServiceItemAnyType for type-switch of CycServiceItemType")
		}
	case syntaxId == 0xb0: // CycServiceItemDbReadType
		if _child, err = new(_CycServiceItemDbReadType).parse(ctx, readBuffer, m); err != nil {
			return nil, errors.Wrap(err, "Error parsing sub-type CycServiceItemDbReadType for type-switch of CycServiceItemType")
		}
	default:
		return nil, errors.Errorf("Unmapped type for parameters [syntaxId=%v]", syntaxId)
	}

	if closeErr := readBuffer.CloseContext("CycServiceItemType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CycServiceItemType")
	}

	return _child, nil
}

func (pm *_CycServiceItemType) serializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child CycServiceItemType, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("CycServiceItemType"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for CycServiceItemType")
	}

	if err := WriteConstField(ctx, "functionId", CycServiceItemType_FUNCTIONID, WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'functionId' field")
	}

	if err := WriteSimpleField[uint8](ctx, "byteLength", m.GetByteLength(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'byteLength' field")
	}

	if err := WriteSimpleField[uint8](ctx, "syntaxId", m.GetSyntaxId(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'syntaxId' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("CycServiceItemType"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for CycServiceItemType")
	}
	return nil
}

func (m *_CycServiceItemType) IsCycServiceItemType() {}

func (m *_CycServiceItemType) DeepCopy() any {
	return m.deepCopy()
}

func (m *_CycServiceItemType) deepCopy() *_CycServiceItemType {
	if m == nil {
		return nil
	}
	_CycServiceItemTypeCopy := &_CycServiceItemType{
		nil, // will be set by child
		m.ByteLength,
		m.SyntaxId,
	}
	return _CycServiceItemTypeCopy
}
