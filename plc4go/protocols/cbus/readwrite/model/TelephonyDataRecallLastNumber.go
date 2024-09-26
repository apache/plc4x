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

// TelephonyDataRecallLastNumber is the corresponding interface of TelephonyDataRecallLastNumber
type TelephonyDataRecallLastNumber interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	TelephonyData
	// GetRecallLastNumberType returns RecallLastNumberType (property field)
	GetRecallLastNumberType() byte
	// GetNumber returns Number (property field)
	GetNumber() string
	// GetIsNumberOfLastOutgoingCall returns IsNumberOfLastOutgoingCall (virtual field)
	GetIsNumberOfLastOutgoingCall() bool
	// GetIsNumberOfLastIncomingCall returns IsNumberOfLastIncomingCall (virtual field)
	GetIsNumberOfLastIncomingCall() bool
	// IsTelephonyDataRecallLastNumber is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsTelephonyDataRecallLastNumber()
	// CreateBuilder creates a TelephonyDataRecallLastNumberBuilder
	CreateTelephonyDataRecallLastNumberBuilder() TelephonyDataRecallLastNumberBuilder
}

// _TelephonyDataRecallLastNumber is the data-structure of this message
type _TelephonyDataRecallLastNumber struct {
	TelephonyDataContract
	RecallLastNumberType byte
	Number               string
}

var _ TelephonyDataRecallLastNumber = (*_TelephonyDataRecallLastNumber)(nil)
var _ TelephonyDataRequirements = (*_TelephonyDataRecallLastNumber)(nil)

// NewTelephonyDataRecallLastNumber factory function for _TelephonyDataRecallLastNumber
func NewTelephonyDataRecallLastNumber(commandTypeContainer TelephonyCommandTypeContainer, argument byte, recallLastNumberType byte, number string) *_TelephonyDataRecallLastNumber {
	_result := &_TelephonyDataRecallLastNumber{
		TelephonyDataContract: NewTelephonyData(commandTypeContainer, argument),
		RecallLastNumberType:  recallLastNumberType,
		Number:                number,
	}
	_result.TelephonyDataContract.(*_TelephonyData)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// TelephonyDataRecallLastNumberBuilder is a builder for TelephonyDataRecallLastNumber
type TelephonyDataRecallLastNumberBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(recallLastNumberType byte, number string) TelephonyDataRecallLastNumberBuilder
	// WithRecallLastNumberType adds RecallLastNumberType (property field)
	WithRecallLastNumberType(byte) TelephonyDataRecallLastNumberBuilder
	// WithNumber adds Number (property field)
	WithNumber(string) TelephonyDataRecallLastNumberBuilder
	// Build builds the TelephonyDataRecallLastNumber or returns an error if something is wrong
	Build() (TelephonyDataRecallLastNumber, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() TelephonyDataRecallLastNumber
}

// NewTelephonyDataRecallLastNumberBuilder() creates a TelephonyDataRecallLastNumberBuilder
func NewTelephonyDataRecallLastNumberBuilder() TelephonyDataRecallLastNumberBuilder {
	return &_TelephonyDataRecallLastNumberBuilder{_TelephonyDataRecallLastNumber: new(_TelephonyDataRecallLastNumber)}
}

type _TelephonyDataRecallLastNumberBuilder struct {
	*_TelephonyDataRecallLastNumber

	parentBuilder *_TelephonyDataBuilder

	err *utils.MultiError
}

var _ (TelephonyDataRecallLastNumberBuilder) = (*_TelephonyDataRecallLastNumberBuilder)(nil)

func (b *_TelephonyDataRecallLastNumberBuilder) setParent(contract TelephonyDataContract) {
	b.TelephonyDataContract = contract
}

func (b *_TelephonyDataRecallLastNumberBuilder) WithMandatoryFields(recallLastNumberType byte, number string) TelephonyDataRecallLastNumberBuilder {
	return b.WithRecallLastNumberType(recallLastNumberType).WithNumber(number)
}

func (b *_TelephonyDataRecallLastNumberBuilder) WithRecallLastNumberType(recallLastNumberType byte) TelephonyDataRecallLastNumberBuilder {
	b.RecallLastNumberType = recallLastNumberType
	return b
}

func (b *_TelephonyDataRecallLastNumberBuilder) WithNumber(number string) TelephonyDataRecallLastNumberBuilder {
	b.Number = number
	return b
}

func (b *_TelephonyDataRecallLastNumberBuilder) Build() (TelephonyDataRecallLastNumber, error) {
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._TelephonyDataRecallLastNumber.deepCopy(), nil
}

func (b *_TelephonyDataRecallLastNumberBuilder) MustBuild() TelephonyDataRecallLastNumber {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_TelephonyDataRecallLastNumberBuilder) Done() TelephonyDataBuilder {
	return b.parentBuilder
}

func (b *_TelephonyDataRecallLastNumberBuilder) buildForTelephonyData() (TelephonyData, error) {
	return b.Build()
}

func (b *_TelephonyDataRecallLastNumberBuilder) DeepCopy() any {
	_copy := b.CreateTelephonyDataRecallLastNumberBuilder().(*_TelephonyDataRecallLastNumberBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateTelephonyDataRecallLastNumberBuilder creates a TelephonyDataRecallLastNumberBuilder
func (b *_TelephonyDataRecallLastNumber) CreateTelephonyDataRecallLastNumberBuilder() TelephonyDataRecallLastNumberBuilder {
	if b == nil {
		return NewTelephonyDataRecallLastNumberBuilder()
	}
	return &_TelephonyDataRecallLastNumberBuilder{_TelephonyDataRecallLastNumber: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_TelephonyDataRecallLastNumber) GetParent() TelephonyDataContract {
	return m.TelephonyDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_TelephonyDataRecallLastNumber) GetRecallLastNumberType() byte {
	return m.RecallLastNumberType
}

func (m *_TelephonyDataRecallLastNumber) GetNumber() string {
	return m.Number
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_TelephonyDataRecallLastNumber) GetIsNumberOfLastOutgoingCall() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetRecallLastNumberType()) == (0x01)))
}

func (m *_TelephonyDataRecallLastNumber) GetIsNumberOfLastIncomingCall() bool {
	ctx := context.Background()
	_ = ctx
	return bool(bool((m.GetRecallLastNumberType()) == (0x02)))
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastTelephonyDataRecallLastNumber(structType any) TelephonyDataRecallLastNumber {
	if casted, ok := structType.(TelephonyDataRecallLastNumber); ok {
		return casted
	}
	if casted, ok := structType.(*TelephonyDataRecallLastNumber); ok {
		return *casted
	}
	return nil
}

func (m *_TelephonyDataRecallLastNumber) GetTypeName() string {
	return "TelephonyDataRecallLastNumber"
}

func (m *_TelephonyDataRecallLastNumber) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.TelephonyDataContract.(*_TelephonyData).GetLengthInBits(ctx))

	// Simple field (recallLastNumberType)
	lengthInBits += 8

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// Simple field (number)
	lengthInBits += uint16(int32((int32(m.GetCommandTypeContainer().NumBytes()) - int32(int32(2)))) * int32(int32(8)))

	return lengthInBits
}

func (m *_TelephonyDataRecallLastNumber) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_TelephonyDataRecallLastNumber) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_TelephonyData, commandTypeContainer TelephonyCommandTypeContainer) (__telephonyDataRecallLastNumber TelephonyDataRecallLastNumber, err error) {
	m.TelephonyDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("TelephonyDataRecallLastNumber"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for TelephonyDataRecallLastNumber")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	recallLastNumberType, err := ReadSimpleField(ctx, "recallLastNumberType", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'recallLastNumberType' field"))
	}
	m.RecallLastNumberType = recallLastNumberType

	isNumberOfLastOutgoingCall, err := ReadVirtualField[bool](ctx, "isNumberOfLastOutgoingCall", (*bool)(nil), bool((recallLastNumberType) == (0x01)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'isNumberOfLastOutgoingCall' field"))
	}
	_ = isNumberOfLastOutgoingCall

	isNumberOfLastIncomingCall, err := ReadVirtualField[bool](ctx, "isNumberOfLastIncomingCall", (*bool)(nil), bool((recallLastNumberType) == (0x02)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'isNumberOfLastIncomingCall' field"))
	}
	_ = isNumberOfLastIncomingCall

	number, err := ReadSimpleField(ctx, "number", ReadString(readBuffer, uint32(int32((int32(commandTypeContainer.NumBytes())-int32(int32(2))))*int32(int32(8)))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'number' field"))
	}
	m.Number = number

	if closeErr := readBuffer.CloseContext("TelephonyDataRecallLastNumber"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for TelephonyDataRecallLastNumber")
	}

	return m, nil
}

func (m *_TelephonyDataRecallLastNumber) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_TelephonyDataRecallLastNumber) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("TelephonyDataRecallLastNumber"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for TelephonyDataRecallLastNumber")
		}

		if err := WriteSimpleField[byte](ctx, "recallLastNumberType", m.GetRecallLastNumberType(), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'recallLastNumberType' field")
		}
		// Virtual field
		isNumberOfLastOutgoingCall := m.GetIsNumberOfLastOutgoingCall()
		_ = isNumberOfLastOutgoingCall
		if _isNumberOfLastOutgoingCallErr := writeBuffer.WriteVirtual(ctx, "isNumberOfLastOutgoingCall", m.GetIsNumberOfLastOutgoingCall()); _isNumberOfLastOutgoingCallErr != nil {
			return errors.Wrap(_isNumberOfLastOutgoingCallErr, "Error serializing 'isNumberOfLastOutgoingCall' field")
		}
		// Virtual field
		isNumberOfLastIncomingCall := m.GetIsNumberOfLastIncomingCall()
		_ = isNumberOfLastIncomingCall
		if _isNumberOfLastIncomingCallErr := writeBuffer.WriteVirtual(ctx, "isNumberOfLastIncomingCall", m.GetIsNumberOfLastIncomingCall()); _isNumberOfLastIncomingCallErr != nil {
			return errors.Wrap(_isNumberOfLastIncomingCallErr, "Error serializing 'isNumberOfLastIncomingCall' field")
		}

		if err := WriteSimpleField[string](ctx, "number", m.GetNumber(), WriteString(writeBuffer, int32(int32((int32(m.GetCommandTypeContainer().NumBytes())-int32(int32(2))))*int32(int32(8))))); err != nil {
			return errors.Wrap(err, "Error serializing 'number' field")
		}

		if popErr := writeBuffer.PopContext("TelephonyDataRecallLastNumber"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for TelephonyDataRecallLastNumber")
		}
		return nil
	}
	return m.TelephonyDataContract.(*_TelephonyData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_TelephonyDataRecallLastNumber) IsTelephonyDataRecallLastNumber() {}

func (m *_TelephonyDataRecallLastNumber) DeepCopy() any {
	return m.deepCopy()
}

func (m *_TelephonyDataRecallLastNumber) deepCopy() *_TelephonyDataRecallLastNumber {
	if m == nil {
		return nil
	}
	_TelephonyDataRecallLastNumberCopy := &_TelephonyDataRecallLastNumber{
		m.TelephonyDataContract.(*_TelephonyData).deepCopy(),
		m.RecallLastNumberType,
		m.Number,
	}
	m.TelephonyDataContract.(*_TelephonyData)._SubType = m
	return _TelephonyDataRecallLastNumberCopy
}

func (m *_TelephonyDataRecallLastNumber) String() string {
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
