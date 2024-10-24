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

// DatagramWriterGroupTransportDataType is the corresponding interface of DatagramWriterGroupTransportDataType
type DatagramWriterGroupTransportDataType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetMessageRepeatCount returns MessageRepeatCount (property field)
	GetMessageRepeatCount() uint8
	// GetMessageRepeatDelay returns MessageRepeatDelay (property field)
	GetMessageRepeatDelay() float64
	// IsDatagramWriterGroupTransportDataType is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsDatagramWriterGroupTransportDataType()
	// CreateBuilder creates a DatagramWriterGroupTransportDataTypeBuilder
	CreateDatagramWriterGroupTransportDataTypeBuilder() DatagramWriterGroupTransportDataTypeBuilder
}

// _DatagramWriterGroupTransportDataType is the data-structure of this message
type _DatagramWriterGroupTransportDataType struct {
	ExtensionObjectDefinitionContract
	MessageRepeatCount uint8
	MessageRepeatDelay float64
}

var _ DatagramWriterGroupTransportDataType = (*_DatagramWriterGroupTransportDataType)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_DatagramWriterGroupTransportDataType)(nil)

// NewDatagramWriterGroupTransportDataType factory function for _DatagramWriterGroupTransportDataType
func NewDatagramWriterGroupTransportDataType(messageRepeatCount uint8, messageRepeatDelay float64) *_DatagramWriterGroupTransportDataType {
	_result := &_DatagramWriterGroupTransportDataType{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		MessageRepeatCount:                messageRepeatCount,
		MessageRepeatDelay:                messageRepeatDelay,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// DatagramWriterGroupTransportDataTypeBuilder is a builder for DatagramWriterGroupTransportDataType
type DatagramWriterGroupTransportDataTypeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(messageRepeatCount uint8, messageRepeatDelay float64) DatagramWriterGroupTransportDataTypeBuilder
	// WithMessageRepeatCount adds MessageRepeatCount (property field)
	WithMessageRepeatCount(uint8) DatagramWriterGroupTransportDataTypeBuilder
	// WithMessageRepeatDelay adds MessageRepeatDelay (property field)
	WithMessageRepeatDelay(float64) DatagramWriterGroupTransportDataTypeBuilder
	// Build builds the DatagramWriterGroupTransportDataType or returns an error if something is wrong
	Build() (DatagramWriterGroupTransportDataType, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() DatagramWriterGroupTransportDataType
}

// NewDatagramWriterGroupTransportDataTypeBuilder() creates a DatagramWriterGroupTransportDataTypeBuilder
func NewDatagramWriterGroupTransportDataTypeBuilder() DatagramWriterGroupTransportDataTypeBuilder {
	return &_DatagramWriterGroupTransportDataTypeBuilder{_DatagramWriterGroupTransportDataType: new(_DatagramWriterGroupTransportDataType)}
}

type _DatagramWriterGroupTransportDataTypeBuilder struct {
	*_DatagramWriterGroupTransportDataType

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (DatagramWriterGroupTransportDataTypeBuilder) = (*_DatagramWriterGroupTransportDataTypeBuilder)(nil)

func (b *_DatagramWriterGroupTransportDataTypeBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_DatagramWriterGroupTransportDataTypeBuilder) WithMandatoryFields(messageRepeatCount uint8, messageRepeatDelay float64) DatagramWriterGroupTransportDataTypeBuilder {
	return b.WithMessageRepeatCount(messageRepeatCount).WithMessageRepeatDelay(messageRepeatDelay)
}

func (b *_DatagramWriterGroupTransportDataTypeBuilder) WithMessageRepeatCount(messageRepeatCount uint8) DatagramWriterGroupTransportDataTypeBuilder {
	b.MessageRepeatCount = messageRepeatCount
	return b
}

func (b *_DatagramWriterGroupTransportDataTypeBuilder) WithMessageRepeatDelay(messageRepeatDelay float64) DatagramWriterGroupTransportDataTypeBuilder {
	b.MessageRepeatDelay = messageRepeatDelay
	return b
}

func (b *_DatagramWriterGroupTransportDataTypeBuilder) Build() (DatagramWriterGroupTransportDataType, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._DatagramWriterGroupTransportDataType.deepCopy(), nil
}

func (b *_DatagramWriterGroupTransportDataTypeBuilder) MustBuild() DatagramWriterGroupTransportDataType {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_DatagramWriterGroupTransportDataTypeBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_DatagramWriterGroupTransportDataTypeBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_DatagramWriterGroupTransportDataTypeBuilder) DeepCopy() any {
	_copy := b.CreateDatagramWriterGroupTransportDataTypeBuilder().(*_DatagramWriterGroupTransportDataTypeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateDatagramWriterGroupTransportDataTypeBuilder creates a DatagramWriterGroupTransportDataTypeBuilder
func (b *_DatagramWriterGroupTransportDataType) CreateDatagramWriterGroupTransportDataTypeBuilder() DatagramWriterGroupTransportDataTypeBuilder {
	if b == nil {
		return NewDatagramWriterGroupTransportDataTypeBuilder()
	}
	return &_DatagramWriterGroupTransportDataTypeBuilder{_DatagramWriterGroupTransportDataType: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_DatagramWriterGroupTransportDataType) GetExtensionId() int32 {
	return int32(15534)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_DatagramWriterGroupTransportDataType) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_DatagramWriterGroupTransportDataType) GetMessageRepeatCount() uint8 {
	return m.MessageRepeatCount
}

func (m *_DatagramWriterGroupTransportDataType) GetMessageRepeatDelay() float64 {
	return m.MessageRepeatDelay
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastDatagramWriterGroupTransportDataType(structType any) DatagramWriterGroupTransportDataType {
	if casted, ok := structType.(DatagramWriterGroupTransportDataType); ok {
		return casted
	}
	if casted, ok := structType.(*DatagramWriterGroupTransportDataType); ok {
		return *casted
	}
	return nil
}

func (m *_DatagramWriterGroupTransportDataType) GetTypeName() string {
	return "DatagramWriterGroupTransportDataType"
}

func (m *_DatagramWriterGroupTransportDataType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Simple field (messageRepeatCount)
	lengthInBits += 8

	// Simple field (messageRepeatDelay)
	lengthInBits += 64

	return lengthInBits
}

func (m *_DatagramWriterGroupTransportDataType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_DatagramWriterGroupTransportDataType) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__datagramWriterGroupTransportDataType DatagramWriterGroupTransportDataType, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("DatagramWriterGroupTransportDataType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for DatagramWriterGroupTransportDataType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	messageRepeatCount, err := ReadSimpleField(ctx, "messageRepeatCount", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'messageRepeatCount' field"))
	}
	m.MessageRepeatCount = messageRepeatCount

	messageRepeatDelay, err := ReadSimpleField(ctx, "messageRepeatDelay", ReadDouble(readBuffer, uint8(64)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'messageRepeatDelay' field"))
	}
	m.MessageRepeatDelay = messageRepeatDelay

	if closeErr := readBuffer.CloseContext("DatagramWriterGroupTransportDataType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for DatagramWriterGroupTransportDataType")
	}

	return m, nil
}

func (m *_DatagramWriterGroupTransportDataType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_DatagramWriterGroupTransportDataType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("DatagramWriterGroupTransportDataType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for DatagramWriterGroupTransportDataType")
		}

		if err := WriteSimpleField[uint8](ctx, "messageRepeatCount", m.GetMessageRepeatCount(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'messageRepeatCount' field")
		}

		if err := WriteSimpleField[float64](ctx, "messageRepeatDelay", m.GetMessageRepeatDelay(), WriteDouble(writeBuffer, 64)); err != nil {
			return errors.Wrap(err, "Error serializing 'messageRepeatDelay' field")
		}

		if popErr := writeBuffer.PopContext("DatagramWriterGroupTransportDataType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for DatagramWriterGroupTransportDataType")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_DatagramWriterGroupTransportDataType) IsDatagramWriterGroupTransportDataType() {}

func (m *_DatagramWriterGroupTransportDataType) DeepCopy() any {
	return m.deepCopy()
}

func (m *_DatagramWriterGroupTransportDataType) deepCopy() *_DatagramWriterGroupTransportDataType {
	if m == nil {
		return nil
	}
	_DatagramWriterGroupTransportDataTypeCopy := &_DatagramWriterGroupTransportDataType{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		m.MessageRepeatCount,
		m.MessageRepeatDelay,
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _DatagramWriterGroupTransportDataTypeCopy
}

func (m *_DatagramWriterGroupTransportDataType) String() string {
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
