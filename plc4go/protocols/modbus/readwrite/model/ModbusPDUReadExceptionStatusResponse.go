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

// ModbusPDUReadExceptionStatusResponse is the corresponding interface of ModbusPDUReadExceptionStatusResponse
type ModbusPDUReadExceptionStatusResponse interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ModbusPDU
	// GetValue returns Value (property field)
	GetValue() uint8
	// IsModbusPDUReadExceptionStatusResponse is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsModbusPDUReadExceptionStatusResponse()
	// CreateBuilder creates a ModbusPDUReadExceptionStatusResponseBuilder
	CreateModbusPDUReadExceptionStatusResponseBuilder() ModbusPDUReadExceptionStatusResponseBuilder
}

// _ModbusPDUReadExceptionStatusResponse is the data-structure of this message
type _ModbusPDUReadExceptionStatusResponse struct {
	ModbusPDUContract
	Value uint8
}

var _ ModbusPDUReadExceptionStatusResponse = (*_ModbusPDUReadExceptionStatusResponse)(nil)
var _ ModbusPDURequirements = (*_ModbusPDUReadExceptionStatusResponse)(nil)

// NewModbusPDUReadExceptionStatusResponse factory function for _ModbusPDUReadExceptionStatusResponse
func NewModbusPDUReadExceptionStatusResponse(value uint8) *_ModbusPDUReadExceptionStatusResponse {
	_result := &_ModbusPDUReadExceptionStatusResponse{
		ModbusPDUContract: NewModbusPDU(),
		Value:             value,
	}
	_result.ModbusPDUContract.(*_ModbusPDU)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// ModbusPDUReadExceptionStatusResponseBuilder is a builder for ModbusPDUReadExceptionStatusResponse
type ModbusPDUReadExceptionStatusResponseBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(value uint8) ModbusPDUReadExceptionStatusResponseBuilder
	// WithValue adds Value (property field)
	WithValue(uint8) ModbusPDUReadExceptionStatusResponseBuilder
	// Build builds the ModbusPDUReadExceptionStatusResponse or returns an error if something is wrong
	Build() (ModbusPDUReadExceptionStatusResponse, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() ModbusPDUReadExceptionStatusResponse
}

// NewModbusPDUReadExceptionStatusResponseBuilder() creates a ModbusPDUReadExceptionStatusResponseBuilder
func NewModbusPDUReadExceptionStatusResponseBuilder() ModbusPDUReadExceptionStatusResponseBuilder {
	return &_ModbusPDUReadExceptionStatusResponseBuilder{_ModbusPDUReadExceptionStatusResponse: new(_ModbusPDUReadExceptionStatusResponse)}
}

type _ModbusPDUReadExceptionStatusResponseBuilder struct {
	*_ModbusPDUReadExceptionStatusResponse

	parentBuilder *_ModbusPDUBuilder

	err *utils.MultiError
}

var _ (ModbusPDUReadExceptionStatusResponseBuilder) = (*_ModbusPDUReadExceptionStatusResponseBuilder)(nil)

func (b *_ModbusPDUReadExceptionStatusResponseBuilder) setParent(contract ModbusPDUContract) {
	b.ModbusPDUContract = contract
}

func (b *_ModbusPDUReadExceptionStatusResponseBuilder) WithMandatoryFields(value uint8) ModbusPDUReadExceptionStatusResponseBuilder {
	return b.WithValue(value)
}

func (b *_ModbusPDUReadExceptionStatusResponseBuilder) WithValue(value uint8) ModbusPDUReadExceptionStatusResponseBuilder {
	b.Value = value
	return b
}

func (b *_ModbusPDUReadExceptionStatusResponseBuilder) Build() (ModbusPDUReadExceptionStatusResponse, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._ModbusPDUReadExceptionStatusResponse.deepCopy(), nil
}

func (b *_ModbusPDUReadExceptionStatusResponseBuilder) MustBuild() ModbusPDUReadExceptionStatusResponse {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_ModbusPDUReadExceptionStatusResponseBuilder) Done() ModbusPDUBuilder {
	return b.parentBuilder
}

func (b *_ModbusPDUReadExceptionStatusResponseBuilder) buildForModbusPDU() (ModbusPDU, error) {
	return b.Build()
}

func (b *_ModbusPDUReadExceptionStatusResponseBuilder) DeepCopy() any {
	_copy := b.CreateModbusPDUReadExceptionStatusResponseBuilder().(*_ModbusPDUReadExceptionStatusResponseBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateModbusPDUReadExceptionStatusResponseBuilder creates a ModbusPDUReadExceptionStatusResponseBuilder
func (b *_ModbusPDUReadExceptionStatusResponse) CreateModbusPDUReadExceptionStatusResponseBuilder() ModbusPDUReadExceptionStatusResponseBuilder {
	if b == nil {
		return NewModbusPDUReadExceptionStatusResponseBuilder()
	}
	return &_ModbusPDUReadExceptionStatusResponseBuilder{_ModbusPDUReadExceptionStatusResponse: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ModbusPDUReadExceptionStatusResponse) GetErrorFlag() bool {
	return bool(false)
}

func (m *_ModbusPDUReadExceptionStatusResponse) GetFunctionFlag() uint8 {
	return 0x07
}

func (m *_ModbusPDUReadExceptionStatusResponse) GetResponse() bool {
	return bool(true)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ModbusPDUReadExceptionStatusResponse) GetParent() ModbusPDUContract {
	return m.ModbusPDUContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ModbusPDUReadExceptionStatusResponse) GetValue() uint8 {
	return m.Value
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastModbusPDUReadExceptionStatusResponse(structType any) ModbusPDUReadExceptionStatusResponse {
	if casted, ok := structType.(ModbusPDUReadExceptionStatusResponse); ok {
		return casted
	}
	if casted, ok := structType.(*ModbusPDUReadExceptionStatusResponse); ok {
		return *casted
	}
	return nil
}

func (m *_ModbusPDUReadExceptionStatusResponse) GetTypeName() string {
	return "ModbusPDUReadExceptionStatusResponse"
}

func (m *_ModbusPDUReadExceptionStatusResponse) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ModbusPDUContract.(*_ModbusPDU).GetLengthInBits(ctx))

	// Simple field (value)
	lengthInBits += 8

	return lengthInBits
}

func (m *_ModbusPDUReadExceptionStatusResponse) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_ModbusPDUReadExceptionStatusResponse) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ModbusPDU, response bool) (__modbusPDUReadExceptionStatusResponse ModbusPDUReadExceptionStatusResponse, err error) {
	m.ModbusPDUContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ModbusPDUReadExceptionStatusResponse"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ModbusPDUReadExceptionStatusResponse")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	value, err := ReadSimpleField(ctx, "value", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'value' field"))
	}
	m.Value = value

	if closeErr := readBuffer.CloseContext("ModbusPDUReadExceptionStatusResponse"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ModbusPDUReadExceptionStatusResponse")
	}

	return m, nil
}

func (m *_ModbusPDUReadExceptionStatusResponse) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ModbusPDUReadExceptionStatusResponse) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ModbusPDUReadExceptionStatusResponse"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ModbusPDUReadExceptionStatusResponse")
		}

		if err := WriteSimpleField[uint8](ctx, "value", m.GetValue(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'value' field")
		}

		if popErr := writeBuffer.PopContext("ModbusPDUReadExceptionStatusResponse"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ModbusPDUReadExceptionStatusResponse")
		}
		return nil
	}
	return m.ModbusPDUContract.(*_ModbusPDU).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ModbusPDUReadExceptionStatusResponse) IsModbusPDUReadExceptionStatusResponse() {}

func (m *_ModbusPDUReadExceptionStatusResponse) DeepCopy() any {
	return m.deepCopy()
}

func (m *_ModbusPDUReadExceptionStatusResponse) deepCopy() *_ModbusPDUReadExceptionStatusResponse {
	if m == nil {
		return nil
	}
	_ModbusPDUReadExceptionStatusResponseCopy := &_ModbusPDUReadExceptionStatusResponse{
		m.ModbusPDUContract.(*_ModbusPDU).deepCopy(),
		m.Value,
	}
	m.ModbusPDUContract.(*_ModbusPDU)._SubType = m
	return _ModbusPDUReadExceptionStatusResponseCopy
}

func (m *_ModbusPDUReadExceptionStatusResponse) String() string {
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
