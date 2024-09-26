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

// S7MessageResponseData is the corresponding interface of S7MessageResponseData
type S7MessageResponseData interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	S7Message
	// GetErrorClass returns ErrorClass (property field)
	GetErrorClass() uint8
	// GetErrorCode returns ErrorCode (property field)
	GetErrorCode() uint8
	// IsS7MessageResponseData is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsS7MessageResponseData()
	// CreateBuilder creates a S7MessageResponseDataBuilder
	CreateS7MessageResponseDataBuilder() S7MessageResponseDataBuilder
}

// _S7MessageResponseData is the data-structure of this message
type _S7MessageResponseData struct {
	S7MessageContract
	ErrorClass uint8
	ErrorCode  uint8
}

var _ S7MessageResponseData = (*_S7MessageResponseData)(nil)
var _ S7MessageRequirements = (*_S7MessageResponseData)(nil)

// NewS7MessageResponseData factory function for _S7MessageResponseData
func NewS7MessageResponseData(tpduReference uint16, parameter S7Parameter, payload S7Payload, errorClass uint8, errorCode uint8) *_S7MessageResponseData {
	_result := &_S7MessageResponseData{
		S7MessageContract: NewS7Message(tpduReference, parameter, payload),
		ErrorClass:        errorClass,
		ErrorCode:         errorCode,
	}
	_result.S7MessageContract.(*_S7Message)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// S7MessageResponseDataBuilder is a builder for S7MessageResponseData
type S7MessageResponseDataBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(errorClass uint8, errorCode uint8) S7MessageResponseDataBuilder
	// WithErrorClass adds ErrorClass (property field)
	WithErrorClass(uint8) S7MessageResponseDataBuilder
	// WithErrorCode adds ErrorCode (property field)
	WithErrorCode(uint8) S7MessageResponseDataBuilder
	// Build builds the S7MessageResponseData or returns an error if something is wrong
	Build() (S7MessageResponseData, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() S7MessageResponseData
}

// NewS7MessageResponseDataBuilder() creates a S7MessageResponseDataBuilder
func NewS7MessageResponseDataBuilder() S7MessageResponseDataBuilder {
	return &_S7MessageResponseDataBuilder{_S7MessageResponseData: new(_S7MessageResponseData)}
}

type _S7MessageResponseDataBuilder struct {
	*_S7MessageResponseData

	parentBuilder *_S7MessageBuilder

	err *utils.MultiError
}

var _ (S7MessageResponseDataBuilder) = (*_S7MessageResponseDataBuilder)(nil)

func (b *_S7MessageResponseDataBuilder) setParent(contract S7MessageContract) {
	b.S7MessageContract = contract
}

func (b *_S7MessageResponseDataBuilder) WithMandatoryFields(errorClass uint8, errorCode uint8) S7MessageResponseDataBuilder {
	return b.WithErrorClass(errorClass).WithErrorCode(errorCode)
}

func (b *_S7MessageResponseDataBuilder) WithErrorClass(errorClass uint8) S7MessageResponseDataBuilder {
	b.ErrorClass = errorClass
	return b
}

func (b *_S7MessageResponseDataBuilder) WithErrorCode(errorCode uint8) S7MessageResponseDataBuilder {
	b.ErrorCode = errorCode
	return b
}

func (b *_S7MessageResponseDataBuilder) Build() (S7MessageResponseData, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._S7MessageResponseData.deepCopy(), nil
}

func (b *_S7MessageResponseDataBuilder) MustBuild() S7MessageResponseData {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_S7MessageResponseDataBuilder) Done() S7MessageBuilder {
	return b.parentBuilder
}

func (b *_S7MessageResponseDataBuilder) buildForS7Message() (S7Message, error) {
	return b.Build()
}

func (b *_S7MessageResponseDataBuilder) DeepCopy() any {
	_copy := b.CreateS7MessageResponseDataBuilder().(*_S7MessageResponseDataBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateS7MessageResponseDataBuilder creates a S7MessageResponseDataBuilder
func (b *_S7MessageResponseData) CreateS7MessageResponseDataBuilder() S7MessageResponseDataBuilder {
	if b == nil {
		return NewS7MessageResponseDataBuilder()
	}
	return &_S7MessageResponseDataBuilder{_S7MessageResponseData: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_S7MessageResponseData) GetMessageType() uint8 {
	return 0x03
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_S7MessageResponseData) GetParent() S7MessageContract {
	return m.S7MessageContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_S7MessageResponseData) GetErrorClass() uint8 {
	return m.ErrorClass
}

func (m *_S7MessageResponseData) GetErrorCode() uint8 {
	return m.ErrorCode
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastS7MessageResponseData(structType any) S7MessageResponseData {
	if casted, ok := structType.(S7MessageResponseData); ok {
		return casted
	}
	if casted, ok := structType.(*S7MessageResponseData); ok {
		return *casted
	}
	return nil
}

func (m *_S7MessageResponseData) GetTypeName() string {
	return "S7MessageResponseData"
}

func (m *_S7MessageResponseData) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.S7MessageContract.(*_S7Message).GetLengthInBits(ctx))

	// Simple field (errorClass)
	lengthInBits += 8

	// Simple field (errorCode)
	lengthInBits += 8

	return lengthInBits
}

func (m *_S7MessageResponseData) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_S7MessageResponseData) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_S7Message) (__s7MessageResponseData S7MessageResponseData, err error) {
	m.S7MessageContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("S7MessageResponseData"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for S7MessageResponseData")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	errorClass, err := ReadSimpleField(ctx, "errorClass", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'errorClass' field"))
	}
	m.ErrorClass = errorClass

	errorCode, err := ReadSimpleField(ctx, "errorCode", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'errorCode' field"))
	}
	m.ErrorCode = errorCode

	if closeErr := readBuffer.CloseContext("S7MessageResponseData"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for S7MessageResponseData")
	}

	return m, nil
}

func (m *_S7MessageResponseData) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_S7MessageResponseData) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("S7MessageResponseData"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for S7MessageResponseData")
		}

		if err := WriteSimpleField[uint8](ctx, "errorClass", m.GetErrorClass(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'errorClass' field")
		}

		if err := WriteSimpleField[uint8](ctx, "errorCode", m.GetErrorCode(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'errorCode' field")
		}

		if popErr := writeBuffer.PopContext("S7MessageResponseData"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for S7MessageResponseData")
		}
		return nil
	}
	return m.S7MessageContract.(*_S7Message).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_S7MessageResponseData) IsS7MessageResponseData() {}

func (m *_S7MessageResponseData) DeepCopy() any {
	return m.deepCopy()
}

func (m *_S7MessageResponseData) deepCopy() *_S7MessageResponseData {
	if m == nil {
		return nil
	}
	_S7MessageResponseDataCopy := &_S7MessageResponseData{
		m.S7MessageContract.(*_S7Message).deepCopy(),
		m.ErrorClass,
		m.ErrorCode,
	}
	m.S7MessageContract.(*_S7Message)._SubType = m
	return _S7MessageResponseDataCopy
}

func (m *_S7MessageResponseData) String() string {
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
