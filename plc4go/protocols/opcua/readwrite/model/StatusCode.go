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

// StatusCode is the corresponding interface of StatusCode
type StatusCode interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	// GetStatusCode returns StatusCode (property field)
	GetStatusCode() uint32
	// IsStatusCode is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsStatusCode()
	// CreateBuilder creates a StatusCodeBuilder
	CreateStatusCodeBuilder() StatusCodeBuilder
}

// _StatusCode is the data-structure of this message
type _StatusCode struct {
	StatusCode uint32
}

var _ StatusCode = (*_StatusCode)(nil)

// NewStatusCode factory function for _StatusCode
func NewStatusCode(statusCode uint32) *_StatusCode {
	return &_StatusCode{StatusCode: statusCode}
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// StatusCodeBuilder is a builder for StatusCode
type StatusCodeBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(statusCode uint32) StatusCodeBuilder
	// WithStatusCode adds StatusCode (property field)
	WithStatusCode(uint32) StatusCodeBuilder
	// Build builds the StatusCode or returns an error if something is wrong
	Build() (StatusCode, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() StatusCode
}

// NewStatusCodeBuilder() creates a StatusCodeBuilder
func NewStatusCodeBuilder() StatusCodeBuilder {
	return &_StatusCodeBuilder{_StatusCode: new(_StatusCode)}
}

type _StatusCodeBuilder struct {
	*_StatusCode

	err *utils.MultiError
}

var _ (StatusCodeBuilder) = (*_StatusCodeBuilder)(nil)

func (b *_StatusCodeBuilder) WithMandatoryFields(statusCode uint32) StatusCodeBuilder {
	return b.WithStatusCode(statusCode)
}

func (b *_StatusCodeBuilder) WithStatusCode(statusCode uint32) StatusCodeBuilder {
	b.StatusCode = statusCode
	return b
}

func (b *_StatusCodeBuilder) Build() (StatusCode, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._StatusCode.deepCopy(), nil
}

func (b *_StatusCodeBuilder) MustBuild() StatusCode {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

func (b *_StatusCodeBuilder) DeepCopy() any {
	_copy := b.CreateStatusCodeBuilder().(*_StatusCodeBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateStatusCodeBuilder creates a StatusCodeBuilder
func (b *_StatusCode) CreateStatusCodeBuilder() StatusCodeBuilder {
	if b == nil {
		return NewStatusCodeBuilder()
	}
	return &_StatusCodeBuilder{_StatusCode: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_StatusCode) GetStatusCode() uint32 {
	return m.StatusCode
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastStatusCode(structType any) StatusCode {
	if casted, ok := structType.(StatusCode); ok {
		return casted
	}
	if casted, ok := structType.(*StatusCode); ok {
		return *casted
	}
	return nil
}

func (m *_StatusCode) GetTypeName() string {
	return "StatusCode"
}

func (m *_StatusCode) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (statusCode)
	lengthInBits += 32

	return lengthInBits
}

func (m *_StatusCode) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func StatusCodeParse(ctx context.Context, theBytes []byte) (StatusCode, error) {
	return StatusCodeParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func StatusCodeParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (StatusCode, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (StatusCode, error) {
		return StatusCodeParseWithBuffer(ctx, readBuffer)
	}
}

func StatusCodeParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (StatusCode, error) {
	v, err := (&_StatusCode{}).parse(ctx, readBuffer)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (m *_StatusCode) parse(ctx context.Context, readBuffer utils.ReadBuffer) (__statusCode StatusCode, err error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("StatusCode"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for StatusCode")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	statusCode, err := ReadSimpleField(ctx, "statusCode", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'statusCode' field"))
	}
	m.StatusCode = statusCode

	if closeErr := readBuffer.CloseContext("StatusCode"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for StatusCode")
	}

	return m, nil
}

func (m *_StatusCode) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_StatusCode) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("StatusCode"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for StatusCode")
	}

	if err := WriteSimpleField[uint32](ctx, "statusCode", m.GetStatusCode(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
		return errors.Wrap(err, "Error serializing 'statusCode' field")
	}

	if popErr := writeBuffer.PopContext("StatusCode"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for StatusCode")
	}
	return nil
}

func (m *_StatusCode) IsStatusCode() {}

func (m *_StatusCode) DeepCopy() any {
	return m.deepCopy()
}

func (m *_StatusCode) deepCopy() *_StatusCode {
	if m == nil {
		return nil
	}
	_StatusCodeCopy := &_StatusCode{
		m.StatusCode,
	}
	return _StatusCodeCopy
}

func (m *_StatusCode) String() string {
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
