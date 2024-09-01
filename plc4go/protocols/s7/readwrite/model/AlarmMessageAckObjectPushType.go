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
const AlarmMessageAckObjectPushType_VARIABLESPEC uint8 = 0x12

// AlarmMessageAckObjectPushType is the corresponding interface of AlarmMessageAckObjectPushType
type AlarmMessageAckObjectPushType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetLengthSpec returns LengthSpec (property field)
	GetLengthSpec() uint8
	// GetSyntaxId returns SyntaxId (property field)
	GetSyntaxId() SyntaxIdType
	// GetNumberOfValues returns NumberOfValues (property field)
	GetNumberOfValues() uint8
	// GetEventId returns EventId (property field)
	GetEventId() uint32
	// GetAckStateGoing returns AckStateGoing (property field)
	GetAckStateGoing() State
	// GetAckStateComing returns AckStateComing (property field)
	GetAckStateComing() State
}

// AlarmMessageAckObjectPushTypeExactly can be used when we want exactly this type and not a type which fulfills AlarmMessageAckObjectPushType.
// This is useful for switch cases.
type AlarmMessageAckObjectPushTypeExactly interface {
	AlarmMessageAckObjectPushType
	isAlarmMessageAckObjectPushType() bool
}

// _AlarmMessageAckObjectPushType is the data-structure of this message
type _AlarmMessageAckObjectPushType struct {
	LengthSpec     uint8
	SyntaxId       SyntaxIdType
	NumberOfValues uint8
	EventId        uint32
	AckStateGoing  State
	AckStateComing State
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AlarmMessageAckObjectPushType) GetLengthSpec() uint8 {
	return m.LengthSpec
}

func (m *_AlarmMessageAckObjectPushType) GetSyntaxId() SyntaxIdType {
	return m.SyntaxId
}

func (m *_AlarmMessageAckObjectPushType) GetNumberOfValues() uint8 {
	return m.NumberOfValues
}

func (m *_AlarmMessageAckObjectPushType) GetEventId() uint32 {
	return m.EventId
}

func (m *_AlarmMessageAckObjectPushType) GetAckStateGoing() State {
	return m.AckStateGoing
}

func (m *_AlarmMessageAckObjectPushType) GetAckStateComing() State {
	return m.AckStateComing
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for const fields.
///////////////////////

func (m *_AlarmMessageAckObjectPushType) GetVariableSpec() uint8 {
	return AlarmMessageAckObjectPushType_VARIABLESPEC
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewAlarmMessageAckObjectPushType factory function for _AlarmMessageAckObjectPushType
func NewAlarmMessageAckObjectPushType(lengthSpec uint8, syntaxId SyntaxIdType, numberOfValues uint8, eventId uint32, ackStateGoing State, ackStateComing State) *_AlarmMessageAckObjectPushType {
	return &_AlarmMessageAckObjectPushType{LengthSpec: lengthSpec, SyntaxId: syntaxId, NumberOfValues: numberOfValues, EventId: eventId, AckStateGoing: ackStateGoing, AckStateComing: ackStateComing}
}

// Deprecated: use the interface for direct cast
func CastAlarmMessageAckObjectPushType(structType any) AlarmMessageAckObjectPushType {
	if casted, ok := structType.(AlarmMessageAckObjectPushType); ok {
		return casted
	}
	if casted, ok := structType.(*AlarmMessageAckObjectPushType); ok {
		return *casted
	}
	return nil
}

func (m *_AlarmMessageAckObjectPushType) GetTypeName() string {
	return "AlarmMessageAckObjectPushType"
}

func (m *_AlarmMessageAckObjectPushType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Const Field (variableSpec)
	lengthInBits += 8

	// Simple field (lengthSpec)
	lengthInBits += 8

	// Simple field (syntaxId)
	lengthInBits += 8

	// Simple field (numberOfValues)
	lengthInBits += 8

	// Simple field (eventId)
	lengthInBits += 32

	// Simple field (ackStateGoing)
	lengthInBits += m.AckStateGoing.GetLengthInBits(ctx)

	// Simple field (ackStateComing)
	lengthInBits += m.AckStateComing.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_AlarmMessageAckObjectPushType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func AlarmMessageAckObjectPushTypeParse(ctx context.Context, theBytes []byte) (AlarmMessageAckObjectPushType, error) {
	return AlarmMessageAckObjectPushTypeParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func AlarmMessageAckObjectPushTypeParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (AlarmMessageAckObjectPushType, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (AlarmMessageAckObjectPushType, error) {
		return AlarmMessageAckObjectPushTypeParseWithBuffer(ctx, readBuffer)
	}
}

func AlarmMessageAckObjectPushTypeParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (AlarmMessageAckObjectPushType, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AlarmMessageAckObjectPushType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AlarmMessageAckObjectPushType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	variableSpec, err := ReadConstField[uint8](ctx, "variableSpec", ReadUnsignedByte(readBuffer, uint8(8)), AlarmMessageAckObjectPushType_VARIABLESPEC)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'variableSpec' field"))
	}
	_ = variableSpec

	lengthSpec, err := ReadSimpleField(ctx, "lengthSpec", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'lengthSpec' field"))
	}

	syntaxId, err := ReadEnumField[SyntaxIdType](ctx, "syntaxId", "SyntaxIdType", ReadEnum(SyntaxIdTypeByValue, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'syntaxId' field"))
	}

	numberOfValues, err := ReadSimpleField(ctx, "numberOfValues", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'numberOfValues' field"))
	}

	eventId, err := ReadSimpleField(ctx, "eventId", ReadUnsignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'eventId' field"))
	}

	ackStateGoing, err := ReadSimpleField[State](ctx, "ackStateGoing", ReadComplex[State](StateParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'ackStateGoing' field"))
	}

	ackStateComing, err := ReadSimpleField[State](ctx, "ackStateComing", ReadComplex[State](StateParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'ackStateComing' field"))
	}

	if closeErr := readBuffer.CloseContext("AlarmMessageAckObjectPushType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AlarmMessageAckObjectPushType")
	}

	// Create the instance
	return &_AlarmMessageAckObjectPushType{
		LengthSpec:     lengthSpec,
		SyntaxId:       syntaxId,
		NumberOfValues: numberOfValues,
		EventId:        eventId,
		AckStateGoing:  ackStateGoing,
		AckStateComing: ackStateComing,
	}, nil
}

func (m *_AlarmMessageAckObjectPushType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_AlarmMessageAckObjectPushType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("AlarmMessageAckObjectPushType"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for AlarmMessageAckObjectPushType")
	}

	if err := WriteConstField(ctx, "variableSpec", AlarmMessageAckObjectPushType_VARIABLESPEC, WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'variableSpec' field")
	}

	if err := WriteSimpleField[uint8](ctx, "lengthSpec", m.GetLengthSpec(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'lengthSpec' field")
	}

	if err := WriteSimpleEnumField[SyntaxIdType](ctx, "syntaxId", "SyntaxIdType", m.GetSyntaxId(), WriteEnum[SyntaxIdType, uint8](SyntaxIdType.GetValue, SyntaxIdType.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
		return errors.Wrap(err, "Error serializing 'syntaxId' field")
	}

	if err := WriteSimpleField[uint8](ctx, "numberOfValues", m.GetNumberOfValues(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'numberOfValues' field")
	}

	if err := WriteSimpleField[uint32](ctx, "eventId", m.GetEventId(), WriteUnsignedInt(writeBuffer, 32)); err != nil {
		return errors.Wrap(err, "Error serializing 'eventId' field")
	}

	if err := WriteSimpleField[State](ctx, "ackStateGoing", m.GetAckStateGoing(), WriteComplex[State](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'ackStateGoing' field")
	}

	if err := WriteSimpleField[State](ctx, "ackStateComing", m.GetAckStateComing(), WriteComplex[State](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'ackStateComing' field")
	}

	if popErr := writeBuffer.PopContext("AlarmMessageAckObjectPushType"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for AlarmMessageAckObjectPushType")
	}
	return nil
}

func (m *_AlarmMessageAckObjectPushType) isAlarmMessageAckObjectPushType() bool {
	return true
}

func (m *_AlarmMessageAckObjectPushType) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
