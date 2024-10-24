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

// S7PayloadAlarm8 is the corresponding interface of S7PayloadAlarm8
type S7PayloadAlarm8 interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	S7PayloadUserDataItem
	// GetAlarmMessage returns AlarmMessage (property field)
	GetAlarmMessage() AlarmMessagePushType
	// IsS7PayloadAlarm8 is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsS7PayloadAlarm8()
	// CreateBuilder creates a S7PayloadAlarm8Builder
	CreateS7PayloadAlarm8Builder() S7PayloadAlarm8Builder
}

// _S7PayloadAlarm8 is the data-structure of this message
type _S7PayloadAlarm8 struct {
	S7PayloadUserDataItemContract
	AlarmMessage AlarmMessagePushType
}

var _ S7PayloadAlarm8 = (*_S7PayloadAlarm8)(nil)
var _ S7PayloadUserDataItemRequirements = (*_S7PayloadAlarm8)(nil)

// NewS7PayloadAlarm8 factory function for _S7PayloadAlarm8
func NewS7PayloadAlarm8(returnCode DataTransportErrorCode, transportSize DataTransportSize, dataLength uint16, alarmMessage AlarmMessagePushType) *_S7PayloadAlarm8 {
	if alarmMessage == nil {
		panic("alarmMessage of type AlarmMessagePushType for S7PayloadAlarm8 must not be nil")
	}
	_result := &_S7PayloadAlarm8{
		S7PayloadUserDataItemContract: NewS7PayloadUserDataItem(returnCode, transportSize, dataLength),
		AlarmMessage:                  alarmMessage,
	}
	_result.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// S7PayloadAlarm8Builder is a builder for S7PayloadAlarm8
type S7PayloadAlarm8Builder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(alarmMessage AlarmMessagePushType) S7PayloadAlarm8Builder
	// WithAlarmMessage adds AlarmMessage (property field)
	WithAlarmMessage(AlarmMessagePushType) S7PayloadAlarm8Builder
	// WithAlarmMessageBuilder adds AlarmMessage (property field) which is build by the builder
	WithAlarmMessageBuilder(func(AlarmMessagePushTypeBuilder) AlarmMessagePushTypeBuilder) S7PayloadAlarm8Builder
	// Build builds the S7PayloadAlarm8 or returns an error if something is wrong
	Build() (S7PayloadAlarm8, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() S7PayloadAlarm8
}

// NewS7PayloadAlarm8Builder() creates a S7PayloadAlarm8Builder
func NewS7PayloadAlarm8Builder() S7PayloadAlarm8Builder {
	return &_S7PayloadAlarm8Builder{_S7PayloadAlarm8: new(_S7PayloadAlarm8)}
}

type _S7PayloadAlarm8Builder struct {
	*_S7PayloadAlarm8

	parentBuilder *_S7PayloadUserDataItemBuilder

	err *utils.MultiError
}

var _ (S7PayloadAlarm8Builder) = (*_S7PayloadAlarm8Builder)(nil)

func (b *_S7PayloadAlarm8Builder) setParent(contract S7PayloadUserDataItemContract) {
	b.S7PayloadUserDataItemContract = contract
}

func (b *_S7PayloadAlarm8Builder) WithMandatoryFields(alarmMessage AlarmMessagePushType) S7PayloadAlarm8Builder {
	return b.WithAlarmMessage(alarmMessage)
}

func (b *_S7PayloadAlarm8Builder) WithAlarmMessage(alarmMessage AlarmMessagePushType) S7PayloadAlarm8Builder {
	b.AlarmMessage = alarmMessage
	return b
}

func (b *_S7PayloadAlarm8Builder) WithAlarmMessageBuilder(builderSupplier func(AlarmMessagePushTypeBuilder) AlarmMessagePushTypeBuilder) S7PayloadAlarm8Builder {
	builder := builderSupplier(b.AlarmMessage.CreateAlarmMessagePushTypeBuilder())
	var err error
	b.AlarmMessage, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "AlarmMessagePushTypeBuilder failed"))
	}
	return b
}

func (b *_S7PayloadAlarm8Builder) Build() (S7PayloadAlarm8, error) {
	if b.AlarmMessage == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'alarmMessage' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._S7PayloadAlarm8.deepCopy(), nil
}

func (b *_S7PayloadAlarm8Builder) MustBuild() S7PayloadAlarm8 {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_S7PayloadAlarm8Builder) Done() S7PayloadUserDataItemBuilder {
	return b.parentBuilder
}

func (b *_S7PayloadAlarm8Builder) buildForS7PayloadUserDataItem() (S7PayloadUserDataItem, error) {
	return b.Build()
}

func (b *_S7PayloadAlarm8Builder) DeepCopy() any {
	_copy := b.CreateS7PayloadAlarm8Builder().(*_S7PayloadAlarm8Builder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateS7PayloadAlarm8Builder creates a S7PayloadAlarm8Builder
func (b *_S7PayloadAlarm8) CreateS7PayloadAlarm8Builder() S7PayloadAlarm8Builder {
	if b == nil {
		return NewS7PayloadAlarm8Builder()
	}
	return &_S7PayloadAlarm8Builder{_S7PayloadAlarm8: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_S7PayloadAlarm8) GetCpuFunctionGroup() uint8 {
	return 0x04
}

func (m *_S7PayloadAlarm8) GetCpuFunctionType() uint8 {
	return 0x00
}

func (m *_S7PayloadAlarm8) GetCpuSubfunction() uint8 {
	return 0x05
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_S7PayloadAlarm8) GetParent() S7PayloadUserDataItemContract {
	return m.S7PayloadUserDataItemContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_S7PayloadAlarm8) GetAlarmMessage() AlarmMessagePushType {
	return m.AlarmMessage
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastS7PayloadAlarm8(structType any) S7PayloadAlarm8 {
	if casted, ok := structType.(S7PayloadAlarm8); ok {
		return casted
	}
	if casted, ok := structType.(*S7PayloadAlarm8); ok {
		return *casted
	}
	return nil
}

func (m *_S7PayloadAlarm8) GetTypeName() string {
	return "S7PayloadAlarm8"
}

func (m *_S7PayloadAlarm8) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem).GetLengthInBits(ctx))

	// Simple field (alarmMessage)
	lengthInBits += m.AlarmMessage.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_S7PayloadAlarm8) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_S7PayloadAlarm8) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_S7PayloadUserDataItem, cpuFunctionGroup uint8, cpuFunctionType uint8, cpuSubfunction uint8) (__s7PayloadAlarm8 S7PayloadAlarm8, err error) {
	m.S7PayloadUserDataItemContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("S7PayloadAlarm8"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for S7PayloadAlarm8")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	alarmMessage, err := ReadSimpleField[AlarmMessagePushType](ctx, "alarmMessage", ReadComplex[AlarmMessagePushType](AlarmMessagePushTypeParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'alarmMessage' field"))
	}
	m.AlarmMessage = alarmMessage

	if closeErr := readBuffer.CloseContext("S7PayloadAlarm8"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for S7PayloadAlarm8")
	}

	return m, nil
}

func (m *_S7PayloadAlarm8) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_S7PayloadAlarm8) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("S7PayloadAlarm8"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for S7PayloadAlarm8")
		}

		if err := WriteSimpleField[AlarmMessagePushType](ctx, "alarmMessage", m.GetAlarmMessage(), WriteComplex[AlarmMessagePushType](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'alarmMessage' field")
		}

		if popErr := writeBuffer.PopContext("S7PayloadAlarm8"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for S7PayloadAlarm8")
		}
		return nil
	}
	return m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_S7PayloadAlarm8) IsS7PayloadAlarm8() {}

func (m *_S7PayloadAlarm8) DeepCopy() any {
	return m.deepCopy()
}

func (m *_S7PayloadAlarm8) deepCopy() *_S7PayloadAlarm8 {
	if m == nil {
		return nil
	}
	_S7PayloadAlarm8Copy := &_S7PayloadAlarm8{
		m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem).deepCopy(),
		m.AlarmMessage.DeepCopy().(AlarmMessagePushType),
	}
	m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem)._SubType = m
	return _S7PayloadAlarm8Copy
}

func (m *_S7PayloadAlarm8) String() string {
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
