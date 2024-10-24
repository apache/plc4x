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

// S7PayloadUserDataItemCpuFunctionAlarmAckResponse is the corresponding interface of S7PayloadUserDataItemCpuFunctionAlarmAckResponse
type S7PayloadUserDataItemCpuFunctionAlarmAckResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	S7PayloadUserDataItem
	// GetFunctionId returns FunctionId (property field)
	GetFunctionId() uint8
	// GetMessageObjects returns MessageObjects (property field)
	GetMessageObjects() []uint8
	// IsS7PayloadUserDataItemCpuFunctionAlarmAckResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsS7PayloadUserDataItemCpuFunctionAlarmAckResponse()
	// CreateBuilder creates a S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder
	CreateS7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder() S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder
}

// _S7PayloadUserDataItemCpuFunctionAlarmAckResponse is the data-structure of this message
type _S7PayloadUserDataItemCpuFunctionAlarmAckResponse struct {
	S7PayloadUserDataItemContract
	FunctionId     uint8
	MessageObjects []uint8
}

var _ S7PayloadUserDataItemCpuFunctionAlarmAckResponse = (*_S7PayloadUserDataItemCpuFunctionAlarmAckResponse)(nil)
var _ S7PayloadUserDataItemRequirements = (*_S7PayloadUserDataItemCpuFunctionAlarmAckResponse)(nil)

// NewS7PayloadUserDataItemCpuFunctionAlarmAckResponse factory function for _S7PayloadUserDataItemCpuFunctionAlarmAckResponse
func NewS7PayloadUserDataItemCpuFunctionAlarmAckResponse(returnCode DataTransportErrorCode, transportSize DataTransportSize, dataLength uint16, functionId uint8, messageObjects []uint8) *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse {
	_result := &_S7PayloadUserDataItemCpuFunctionAlarmAckResponse{
		S7PayloadUserDataItemContract: NewS7PayloadUserDataItem(returnCode, transportSize, dataLength),
		FunctionId:                    functionId,
		MessageObjects:                messageObjects,
	}
	_result.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder is a builder for S7PayloadUserDataItemCpuFunctionAlarmAckResponse
type S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(functionId uint8, messageObjects []uint8) S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder
	// WithFunctionId adds FunctionId (property field)
	WithFunctionId(uint8) S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder
	// WithMessageObjects adds MessageObjects (property field)
	WithMessageObjects(...uint8) S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder
	// Build builds the S7PayloadUserDataItemCpuFunctionAlarmAckResponse or returns an error if something is wrong
	Build() (S7PayloadUserDataItemCpuFunctionAlarmAckResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() S7PayloadUserDataItemCpuFunctionAlarmAckResponse
}

// NewS7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder() creates a S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder
func NewS7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder() S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder {
	return &_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder{_S7PayloadUserDataItemCpuFunctionAlarmAckResponse: new(_S7PayloadUserDataItemCpuFunctionAlarmAckResponse)}
}

type _S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder struct {
	*_S7PayloadUserDataItemCpuFunctionAlarmAckResponse

	parentBuilder *_S7PayloadUserDataItemBuilder

	err *utils.MultiError
}

var _ (S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder) = (*_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder)(nil)

func (b *_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder) setParent(contract S7PayloadUserDataItemContract) {
	b.S7PayloadUserDataItemContract = contract
}

func (b *_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder) WithMandatoryFields(functionId uint8, messageObjects []uint8) S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder {
	return b.WithFunctionId(functionId).WithMessageObjects(messageObjects...)
}

func (b *_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder) WithFunctionId(functionId uint8) S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder {
	b.FunctionId = functionId
	return b
}

func (b *_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder) WithMessageObjects(messageObjects ...uint8) S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder {
	b.MessageObjects = messageObjects
	return b
}

func (b *_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder) Build() (S7PayloadUserDataItemCpuFunctionAlarmAckResponse, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._S7PayloadUserDataItemCpuFunctionAlarmAckResponse.deepCopy(), nil
}

func (b *_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder) MustBuild() S7PayloadUserDataItemCpuFunctionAlarmAckResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder) Done() S7PayloadUserDataItemBuilder {
	return b.parentBuilder
}

func (b *_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder) buildForS7PayloadUserDataItem() (S7PayloadUserDataItem, error) {
	return b.Build()
}

func (b *_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder) DeepCopy() any {
	_copy := b.CreateS7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder().(*_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateS7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder creates a S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder
func (b *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) CreateS7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder() S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder {
	if b == nil {
		return NewS7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder()
	}
	return &_S7PayloadUserDataItemCpuFunctionAlarmAckResponseBuilder{_S7PayloadUserDataItemCpuFunctionAlarmAckResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) GetCpuFunctionGroup() uint8 {
	return 0x04
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) GetCpuFunctionType() uint8 {
	return 0x08
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) GetCpuSubfunction() uint8 {
	return 0x0b
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) GetParent() S7PayloadUserDataItemContract {
	return m.S7PayloadUserDataItemContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) GetFunctionId() uint8 {
	return m.FunctionId
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) GetMessageObjects() []uint8 {
	return m.MessageObjects
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastS7PayloadUserDataItemCpuFunctionAlarmAckResponse(structType any) S7PayloadUserDataItemCpuFunctionAlarmAckResponse {
	if casted, ok := structType.(S7PayloadUserDataItemCpuFunctionAlarmAckResponse); ok {
		return casted
	}
	if casted, ok := structType.(*S7PayloadUserDataItemCpuFunctionAlarmAckResponse); ok {
		return *casted
	}
	return nil
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) GetTypeName() string {
	return "S7PayloadUserDataItemCpuFunctionAlarmAckResponse"
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem).GetLengthInBits(ctx))

	// Simple field (functionId)
	lengthInBits += 8

	// Implicit Field (numberOfObjects)
	lengthInBits += 8

	// Array field
	if len(m.MessageObjects) > 0 {
		lengthInBits += 8 * uint16(len(m.MessageObjects))
	}

	return lengthInBits
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_S7PayloadUserDataItem, cpuFunctionGroup uint8, cpuFunctionType uint8, cpuSubfunction uint8) (__s7PayloadUserDataItemCpuFunctionAlarmAckResponse S7PayloadUserDataItemCpuFunctionAlarmAckResponse, err error) {
	m.S7PayloadUserDataItemContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("S7PayloadUserDataItemCpuFunctionAlarmAckResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for S7PayloadUserDataItemCpuFunctionAlarmAckResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	functionId, err := ReadSimpleField(ctx, "functionId", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'functionId' field"))
	}
	m.FunctionId = functionId

	numberOfObjects, err := ReadImplicitField[uint8](ctx, "numberOfObjects", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'numberOfObjects' field"))
	}
	_ = numberOfObjects

	messageObjects, err := ReadCountArrayField[uint8](ctx, "messageObjects", ReadUnsignedByte(readBuffer, uint8(8)), uint64(numberOfObjects))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'messageObjects' field"))
	}
	m.MessageObjects = messageObjects

	if closeErr := readBuffer.CloseContext("S7PayloadUserDataItemCpuFunctionAlarmAckResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for S7PayloadUserDataItemCpuFunctionAlarmAckResponse")
	}

	return m, nil
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("S7PayloadUserDataItemCpuFunctionAlarmAckResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for S7PayloadUserDataItemCpuFunctionAlarmAckResponse")
		}

		if err := WriteSimpleField[uint8](ctx, "functionId", m.GetFunctionId(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'functionId' field")
		}
		numberOfObjects := uint8(uint8(len(m.GetMessageObjects())))
		if err := WriteImplicitField(ctx, "numberOfObjects", numberOfObjects, WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'numberOfObjects' field")
		}

		if err := WriteSimpleTypeArrayField(ctx, "messageObjects", m.GetMessageObjects(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'messageObjects' field")
		}

		if popErr := writeBuffer.PopContext("S7PayloadUserDataItemCpuFunctionAlarmAckResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for S7PayloadUserDataItemCpuFunctionAlarmAckResponse")
		}
		return nil
	}
	return m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) IsS7PayloadUserDataItemCpuFunctionAlarmAckResponse() {
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) deepCopy() *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse {
	if m == nil {
		return nil
	}
	_S7PayloadUserDataItemCpuFunctionAlarmAckResponseCopy := &_S7PayloadUserDataItemCpuFunctionAlarmAckResponse{
		m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem).deepCopy(),
		m.FunctionId,
		utils.DeepCopySlice[uint8, uint8](m.MessageObjects),
	}
	m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem)._SubType = m
	return _S7PayloadUserDataItemCpuFunctionAlarmAckResponseCopy
}

func (m *_S7PayloadUserDataItemCpuFunctionAlarmAckResponse) String() string {
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
