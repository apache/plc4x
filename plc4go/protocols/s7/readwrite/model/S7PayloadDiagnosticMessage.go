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

// S7PayloadDiagnosticMessage is the corresponding interface of S7PayloadDiagnosticMessage
type S7PayloadDiagnosticMessage interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	S7PayloadUserDataItem
	// GetEventId returns EventId (property field)
	GetEventId() uint16
	// GetPriorityClass returns PriorityClass (property field)
	GetPriorityClass() uint8
	// GetObNumber returns ObNumber (property field)
	GetObNumber() uint8
	// GetDatId returns DatId (property field)
	GetDatId() uint16
	// GetInfo1 returns Info1 (property field)
	GetInfo1() uint16
	// GetInfo2 returns Info2 (property field)
	GetInfo2() uint32
	// GetTimeStamp returns TimeStamp (property field)
	GetTimeStamp() DateAndTime
	// IsS7PayloadDiagnosticMessage is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsS7PayloadDiagnosticMessage()
	// CreateBuilder creates a S7PayloadDiagnosticMessageBuilder
	CreateS7PayloadDiagnosticMessageBuilder() S7PayloadDiagnosticMessageBuilder
}

// _S7PayloadDiagnosticMessage is the data-structure of this message
type _S7PayloadDiagnosticMessage struct {
	S7PayloadUserDataItemContract
	EventId       uint16
	PriorityClass uint8
	ObNumber      uint8
	DatId         uint16
	Info1         uint16
	Info2         uint32
	TimeStamp     DateAndTime
}

var _ S7PayloadDiagnosticMessage = (*_S7PayloadDiagnosticMessage)(nil)
var _ S7PayloadUserDataItemRequirements = (*_S7PayloadDiagnosticMessage)(nil)

// NewS7PayloadDiagnosticMessage factory function for _S7PayloadDiagnosticMessage
func NewS7PayloadDiagnosticMessage(returnCode DataTransportErrorCode, transportSize DataTransportSize, dataLength uint16, eventId uint16, priorityClass uint8, obNumber uint8, datId uint16, info1 uint16, info2 uint32, timeStamp DateAndTime) *_S7PayloadDiagnosticMessage {
	if timeStamp == nil {
		panic("timeStamp of type DateAndTime for S7PayloadDiagnosticMessage must not be nil")
	}
	_result := &_S7PayloadDiagnosticMessage{
		S7PayloadUserDataItemContract: NewS7PayloadUserDataItem(returnCode, transportSize, dataLength),
		EventId:                       eventId,
		PriorityClass:                 priorityClass,
		ObNumber:                      obNumber,
		DatId:                         datId,
		Info1:                         info1,
		Info2:                         info2,
		TimeStamp:                     timeStamp,
	}
	_result.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// S7PayloadDiagnosticMessageBuilder is a builder for S7PayloadDiagnosticMessage
type S7PayloadDiagnosticMessageBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(eventId uint16, priorityClass uint8, obNumber uint8, datId uint16, info1 uint16, info2 uint32, timeStamp DateAndTime) S7PayloadDiagnosticMessageBuilder
	// WithEventId adds EventId (property field)
	WithEventId(uint16) S7PayloadDiagnosticMessageBuilder
	// WithPriorityClass adds PriorityClass (property field)
	WithPriorityClass(uint8) S7PayloadDiagnosticMessageBuilder
	// WithObNumber adds ObNumber (property field)
	WithObNumber(uint8) S7PayloadDiagnosticMessageBuilder
	// WithDatId adds DatId (property field)
	WithDatId(uint16) S7PayloadDiagnosticMessageBuilder
	// WithInfo1 adds Info1 (property field)
	WithInfo1(uint16) S7PayloadDiagnosticMessageBuilder
	// WithInfo2 adds Info2 (property field)
	WithInfo2(uint32) S7PayloadDiagnosticMessageBuilder
	// WithTimeStamp adds TimeStamp (property field)
	WithTimeStamp(DateAndTime) S7PayloadDiagnosticMessageBuilder
	// WithTimeStampBuilder adds TimeStamp (property field) which is build by the builder
	WithTimeStampBuilder(func(DateAndTimeBuilder) DateAndTimeBuilder) S7PayloadDiagnosticMessageBuilder
	// Build builds the S7PayloadDiagnosticMessage or returns an error if something is wrong
	Build() (S7PayloadDiagnosticMessage, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() S7PayloadDiagnosticMessage
}

// NewS7PayloadDiagnosticMessageBuilder() creates a S7PayloadDiagnosticMessageBuilder
func NewS7PayloadDiagnosticMessageBuilder() S7PayloadDiagnosticMessageBuilder {
	return &_S7PayloadDiagnosticMessageBuilder{_S7PayloadDiagnosticMessage: new(_S7PayloadDiagnosticMessage)}
}

type _S7PayloadDiagnosticMessageBuilder struct {
	*_S7PayloadDiagnosticMessage

	parentBuilder *_S7PayloadUserDataItemBuilder

	err *utils.MultiError
}

var _ (S7PayloadDiagnosticMessageBuilder) = (*_S7PayloadDiagnosticMessageBuilder)(nil)

func (b *_S7PayloadDiagnosticMessageBuilder) setParent(contract S7PayloadUserDataItemContract) {
	b.S7PayloadUserDataItemContract = contract
}

func (b *_S7PayloadDiagnosticMessageBuilder) WithMandatoryFields(eventId uint16, priorityClass uint8, obNumber uint8, datId uint16, info1 uint16, info2 uint32, timeStamp DateAndTime) S7PayloadDiagnosticMessageBuilder {
	return b.WithEventId(eventId).WithPriorityClass(priorityClass).WithObNumber(obNumber).WithDatId(datId).WithInfo1(info1).WithInfo2(info2).WithTimeStamp(timeStamp)
}

func (b *_S7PayloadDiagnosticMessageBuilder) WithEventId(eventId uint16) S7PayloadDiagnosticMessageBuilder {
	b.EventId = eventId
	return b
}

func (b *_S7PayloadDiagnosticMessageBuilder) WithPriorityClass(priorityClass uint8) S7PayloadDiagnosticMessageBuilder {
	b.PriorityClass = priorityClass
	return b
}

func (b *_S7PayloadDiagnosticMessageBuilder) WithObNumber(obNumber uint8) S7PayloadDiagnosticMessageBuilder {
	b.ObNumber = obNumber
	return b
}

func (b *_S7PayloadDiagnosticMessageBuilder) WithDatId(datId uint16) S7PayloadDiagnosticMessageBuilder {
	b.DatId = datId
	return b
}

func (b *_S7PayloadDiagnosticMessageBuilder) WithInfo1(info1 uint16) S7PayloadDiagnosticMessageBuilder {
	b.Info1 = info1
	return b
}

func (b *_S7PayloadDiagnosticMessageBuilder) WithInfo2(info2 uint32) S7PayloadDiagnosticMessageBuilder {
	b.Info2 = info2
	return b
}

func (b *_S7PayloadDiagnosticMessageBuilder) WithTimeStamp(timeStamp DateAndTime) S7PayloadDiagnosticMessageBuilder {
	b.TimeStamp = timeStamp
	return b
}

func (b *_S7PayloadDiagnosticMessageBuilder) WithTimeStampBuilder(builderSupplier func(DateAndTimeBuilder) DateAndTimeBuilder) S7PayloadDiagnosticMessageBuilder {
	builder := builderSupplier(b.TimeStamp.CreateDateAndTimeBuilder())
	var err error
	b.TimeStamp, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "DateAndTimeBuilder failed"))
	}
	return b
}

func (b *_S7PayloadDiagnosticMessageBuilder) Build() (S7PayloadDiagnosticMessage, error) {
	if b.TimeStamp == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'timeStamp' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._S7PayloadDiagnosticMessage.deepCopy(), nil
}

func (b *_S7PayloadDiagnosticMessageBuilder) MustBuild() S7PayloadDiagnosticMessage {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_S7PayloadDiagnosticMessageBuilder) Done() S7PayloadUserDataItemBuilder {
	return b.parentBuilder
}

func (b *_S7PayloadDiagnosticMessageBuilder) buildForS7PayloadUserDataItem() (S7PayloadUserDataItem, error) {
	return b.Build()
}

func (b *_S7PayloadDiagnosticMessageBuilder) DeepCopy() any {
	_copy := b.CreateS7PayloadDiagnosticMessageBuilder().(*_S7PayloadDiagnosticMessageBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateS7PayloadDiagnosticMessageBuilder creates a S7PayloadDiagnosticMessageBuilder
func (b *_S7PayloadDiagnosticMessage) CreateS7PayloadDiagnosticMessageBuilder() S7PayloadDiagnosticMessageBuilder {
	if b == nil {
		return NewS7PayloadDiagnosticMessageBuilder()
	}
	return &_S7PayloadDiagnosticMessageBuilder{_S7PayloadDiagnosticMessage: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_S7PayloadDiagnosticMessage) GetCpuFunctionGroup() uint8 {
	return 0x04
}

func (m *_S7PayloadDiagnosticMessage) GetCpuFunctionType() uint8 {
	return 0x00
}

func (m *_S7PayloadDiagnosticMessage) GetCpuSubfunction() uint8 {
	return 0x03
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_S7PayloadDiagnosticMessage) GetParent() S7PayloadUserDataItemContract {
	return m.S7PayloadUserDataItemContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_S7PayloadDiagnosticMessage) GetEventId() uint16 {
	return m.EventId
}

func (m *_S7PayloadDiagnosticMessage) GetPriorityClass() uint8 {
	return m.PriorityClass
}

func (m *_S7PayloadDiagnosticMessage) GetObNumber() uint8 {
	return m.ObNumber
}

func (m *_S7PayloadDiagnosticMessage) GetDatId() uint16 {
	return m.DatId
}

func (m *_S7PayloadDiagnosticMessage) GetInfo1() uint16 {
	return m.Info1
}

func (m *_S7PayloadDiagnosticMessage) GetInfo2() uint32 {
	return m.Info2
}

func (m *_S7PayloadDiagnosticMessage) GetTimeStamp() DateAndTime {
	return m.TimeStamp
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastS7PayloadDiagnosticMessage(structType any) S7PayloadDiagnosticMessage {
	if casted, ok := structType.(S7PayloadDiagnosticMessage); ok {
		return casted
	}
	if casted, ok := structType.(*S7PayloadDiagnosticMessage); ok {
		return *casted
	}
	return nil
}

func (m *_S7PayloadDiagnosticMessage) GetTypeName() string {
	return "S7PayloadDiagnosticMessage"
}

func (m *_S7PayloadDiagnosticMessage) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem).GetLengthInBits(ctx))

	// Simple field (eventId)
	lengthInBits += 16

	// Simple field (priorityClass)
	lengthInBits += 8

	// Simple field (obNumber)
	lengthInBits += 8

	// Simple field (datId)
	lengthInBits += 16

	// Simple field (info1)
	lengthInBits += 16

	// Simple field (info2)
	lengthInBits += 32

	// Simple field (timeStamp)
	lengthInBits += m.TimeStamp.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_S7PayloadDiagnosticMessage) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_S7PayloadDiagnosticMessage) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_S7PayloadUserDataItem, cpuFunctionGroup uint8, cpuFunctionType uint8, cpuSubfunction uint8) (__s7PayloadDiagnosticMessage S7PayloadDiagnosticMessage, err error) {
	m.S7PayloadUserDataItemContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("S7PayloadDiagnosticMessage"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for S7PayloadDiagnosticMessage")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	eventId, err := ReadSimpleField(ctx, "eventId", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'eventId' field"))
	}
	m.EventId = eventId

	priorityClass, err := ReadSimpleField(ctx, "priorityClass", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'priorityClass' field"))
	}
	m.PriorityClass = priorityClass

	obNumber, err := ReadSimpleField(ctx, "obNumber", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'obNumber' field"))
	}
	m.ObNumber = obNumber

	datId, err := ReadSimpleField(ctx, "datId", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'datId' field"))
	}
	m.DatId = datId

	info1, err := ReadSimpleField(ctx, "info1", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'info1' field"))
	}
	m.Info1 = info1

	info2, err := ReadSimpleField(ctx, "info2", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'info2' field"))
	}
	m.Info2 = info2

	timeStamp, err := ReadSimpleField[DateAndTime](ctx, "timeStamp", ReadComplex[DateAndTime](DateAndTimeParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'timeStamp' field"))
	}
	m.TimeStamp = timeStamp

	if closeErr := readBuffer.CloseContext("S7PayloadDiagnosticMessage"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for S7PayloadDiagnosticMessage")
	}

	return m, nil
}

func (m *_S7PayloadDiagnosticMessage) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_S7PayloadDiagnosticMessage) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("S7PayloadDiagnosticMessage"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for S7PayloadDiagnosticMessage")
		}

		if err := WriteSimpleField[uint16](ctx, "eventId", m.GetEventId(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'eventId' field")
		}

		if err := WriteSimpleField[uint8](ctx, "priorityClass", m.GetPriorityClass(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'priorityClass' field")
		}

		if err := WriteSimpleField[uint8](ctx, "obNumber", m.GetObNumber(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'obNumber' field")
		}

		if err := WriteSimpleField[uint16](ctx, "datId", m.GetDatId(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'datId' field")
		}

		if err := WriteSimpleField[uint16](ctx, "info1", m.GetInfo1(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'info1' field")
		}

		if err := WriteSimpleField[uint32](ctx, "info2", m.GetInfo2(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'info2' field")
		}

		if err := WriteSimpleField[DateAndTime](ctx, "timeStamp", m.GetTimeStamp(), WriteComplex[DateAndTime](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'timeStamp' field")
		}

		if popErr := writeBuffer.PopContext("S7PayloadDiagnosticMessage"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for S7PayloadDiagnosticMessage")
		}
		return nil
	}
	return m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_S7PayloadDiagnosticMessage) IsS7PayloadDiagnosticMessage() {}

func (m *_S7PayloadDiagnosticMessage) DeepCopy() any {
	return m.deepCopy()
}

func (m *_S7PayloadDiagnosticMessage) deepCopy() *_S7PayloadDiagnosticMessage {
	if m == nil {
		return nil
	}
	_S7PayloadDiagnosticMessageCopy := &_S7PayloadDiagnosticMessage{
		m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem).deepCopy(),
		m.EventId,
		m.PriorityClass,
		m.ObNumber,
		m.DatId,
		m.Info1,
		m.Info2,
		m.TimeStamp.DeepCopy().(DateAndTime),
	}
	m.S7PayloadUserDataItemContract.(*_S7PayloadUserDataItem)._SubType = m
	return _S7PayloadDiagnosticMessageCopy
}

func (m *_S7PayloadDiagnosticMessage) String() string {
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
