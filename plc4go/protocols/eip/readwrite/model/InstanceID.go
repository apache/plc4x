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

// InstanceID is the corresponding interface of InstanceID
type InstanceID interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	LogicalSegmentType
	// GetFormat returns Format (property field)
	GetFormat() uint8
	// GetInstance returns Instance (property field)
	GetInstance() uint8
	// IsInstanceID is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsInstanceID()
	// CreateBuilder creates a InstanceIDBuilder
	CreateInstanceIDBuilder() InstanceIDBuilder
}

// _InstanceID is the data-structure of this message
type _InstanceID struct {
	LogicalSegmentTypeContract
	Format   uint8
	Instance uint8
}

var _ InstanceID = (*_InstanceID)(nil)
var _ LogicalSegmentTypeRequirements = (*_InstanceID)(nil)

// NewInstanceID factory function for _InstanceID
func NewInstanceID(format uint8, instance uint8) *_InstanceID {
	_result := &_InstanceID{
		LogicalSegmentTypeContract: NewLogicalSegmentType(),
		Format:                     format,
		Instance:                   instance,
	}
	_result.LogicalSegmentTypeContract.(*_LogicalSegmentType)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// InstanceIDBuilder is a builder for InstanceID
type InstanceIDBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(format uint8, instance uint8) InstanceIDBuilder
	// WithFormat adds Format (property field)
	WithFormat(uint8) InstanceIDBuilder
	// WithInstance adds Instance (property field)
	WithInstance(uint8) InstanceIDBuilder
	// Build builds the InstanceID or returns an error if something is wrong
	Build() (InstanceID, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() InstanceID
}

// NewInstanceIDBuilder() creates a InstanceIDBuilder
func NewInstanceIDBuilder() InstanceIDBuilder {
	return &_InstanceIDBuilder{_InstanceID: new(_InstanceID)}
}

type _InstanceIDBuilder struct {
	*_InstanceID

	parentBuilder *_LogicalSegmentTypeBuilder

	err *utils.MultiError
}

var _ (InstanceIDBuilder) = (*_InstanceIDBuilder)(nil)

func (b *_InstanceIDBuilder) setParent(contract LogicalSegmentTypeContract) {
	b.LogicalSegmentTypeContract = contract
}

func (b *_InstanceIDBuilder) WithMandatoryFields(format uint8, instance uint8) InstanceIDBuilder {
	return b.WithFormat(format).WithInstance(instance)
}

func (b *_InstanceIDBuilder) WithFormat(format uint8) InstanceIDBuilder {
	b.Format = format
	return b
}

func (b *_InstanceIDBuilder) WithInstance(instance uint8) InstanceIDBuilder {
	b.Instance = instance
	return b
}

func (b *_InstanceIDBuilder) Build() (InstanceID, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._InstanceID.deepCopy(), nil
}

func (b *_InstanceIDBuilder) MustBuild() InstanceID {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_InstanceIDBuilder) Done() LogicalSegmentTypeBuilder {
	return b.parentBuilder
}

func (b *_InstanceIDBuilder) buildForLogicalSegmentType() (LogicalSegmentType, error) {
	return b.Build()
}

func (b *_InstanceIDBuilder) DeepCopy() any {
	_copy := b.CreateInstanceIDBuilder().(*_InstanceIDBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateInstanceIDBuilder creates a InstanceIDBuilder
func (b *_InstanceID) CreateInstanceIDBuilder() InstanceIDBuilder {
	if b == nil {
		return NewInstanceIDBuilder()
	}
	return &_InstanceIDBuilder{_InstanceID: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_InstanceID) GetLogicalSegmentType() uint8 {
	return 0x01
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_InstanceID) GetParent() LogicalSegmentTypeContract {
	return m.LogicalSegmentTypeContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_InstanceID) GetFormat() uint8 {
	return m.Format
}

func (m *_InstanceID) GetInstance() uint8 {
	return m.Instance
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastInstanceID(structType any) InstanceID {
	if casted, ok := structType.(InstanceID); ok {
		return casted
	}
	if casted, ok := structType.(*InstanceID); ok {
		return *casted
	}
	return nil
}

func (m *_InstanceID) GetTypeName() string {
	return "InstanceID"
}

func (m *_InstanceID) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.LogicalSegmentTypeContract.(*_LogicalSegmentType).GetLengthInBits(ctx))

	// Simple field (format)
	lengthInBits += 2

	// Simple field (instance)
	lengthInBits += 8

	return lengthInBits
}

func (m *_InstanceID) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_InstanceID) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_LogicalSegmentType) (__instanceID InstanceID, err error) {
	m.LogicalSegmentTypeContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("InstanceID"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for InstanceID")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	format, err := ReadSimpleField(ctx, "format", ReadUnsignedByte(readBuffer, uint8(2)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'format' field"))
	}
	m.Format = format

	instance, err := ReadSimpleField(ctx, "instance", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'instance' field"))
	}
	m.Instance = instance

	if closeErr := readBuffer.CloseContext("InstanceID"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for InstanceID")
	}

	return m, nil
}

func (m *_InstanceID) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_InstanceID) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("InstanceID"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for InstanceID")
		}

		if err := WriteSimpleField[uint8](ctx, "format", m.GetFormat(), WriteUnsignedByte(writeBuffer, 2)); err != nil {
			return errors.Wrap(err, "Error serializing 'format' field")
		}

		if err := WriteSimpleField[uint8](ctx, "instance", m.GetInstance(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'instance' field")
		}

		if popErr := writeBuffer.PopContext("InstanceID"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for InstanceID")
		}
		return nil
	}
	return m.LogicalSegmentTypeContract.(*_LogicalSegmentType).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_InstanceID) IsInstanceID() {}

func (m *_InstanceID) DeepCopy() any {
	return m.deepCopy()
}

func (m *_InstanceID) deepCopy() *_InstanceID {
	if m == nil {
		return nil
	}
	_InstanceIDCopy := &_InstanceID{
		m.LogicalSegmentTypeContract.(*_LogicalSegmentType).deepCopy(),
		m.Format,
		m.Instance,
	}
	m.LogicalSegmentTypeContract.(*_LogicalSegmentType)._SubType = m
	return _InstanceIDCopy
}

func (m *_InstanceID) String() string {
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
