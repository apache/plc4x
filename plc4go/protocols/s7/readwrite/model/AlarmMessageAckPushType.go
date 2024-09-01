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

// AlarmMessageAckPushType is the corresponding interface of AlarmMessageAckPushType
type AlarmMessageAckPushType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetTimeStamp returns TimeStamp (property field)
	GetTimeStamp() DateAndTime
	// GetFunctionId returns FunctionId (property field)
	GetFunctionId() uint8
	// GetNumberOfObjects returns NumberOfObjects (property field)
	GetNumberOfObjects() uint8
	// GetMessageObjects returns MessageObjects (property field)
	GetMessageObjects() []AlarmMessageAckObjectPushType
}

// AlarmMessageAckPushTypeExactly can be used when we want exactly this type and not a type which fulfills AlarmMessageAckPushType.
// This is useful for switch cases.
type AlarmMessageAckPushTypeExactly interface {
	AlarmMessageAckPushType
	isAlarmMessageAckPushType() bool
}

// _AlarmMessageAckPushType is the data-structure of this message
type _AlarmMessageAckPushType struct {
	TimeStamp       DateAndTime
	FunctionId      uint8
	NumberOfObjects uint8
	MessageObjects  []AlarmMessageAckObjectPushType
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AlarmMessageAckPushType) GetTimeStamp() DateAndTime {
	return m.TimeStamp
}

func (m *_AlarmMessageAckPushType) GetFunctionId() uint8 {
	return m.FunctionId
}

func (m *_AlarmMessageAckPushType) GetNumberOfObjects() uint8 {
	return m.NumberOfObjects
}

func (m *_AlarmMessageAckPushType) GetMessageObjects() []AlarmMessageAckObjectPushType {
	return m.MessageObjects
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewAlarmMessageAckPushType factory function for _AlarmMessageAckPushType
func NewAlarmMessageAckPushType(timeStamp DateAndTime, functionId uint8, numberOfObjects uint8, messageObjects []AlarmMessageAckObjectPushType) *_AlarmMessageAckPushType {
	return &_AlarmMessageAckPushType{TimeStamp: timeStamp, FunctionId: functionId, NumberOfObjects: numberOfObjects, MessageObjects: messageObjects}
}

// Deprecated: use the interface for direct cast
func CastAlarmMessageAckPushType(structType any) AlarmMessageAckPushType {
	if casted, ok := structType.(AlarmMessageAckPushType); ok {
		return casted
	}
	if casted, ok := structType.(*AlarmMessageAckPushType); ok {
		return *casted
	}
	return nil
}

func (m *_AlarmMessageAckPushType) GetTypeName() string {
	return "AlarmMessageAckPushType"
}

func (m *_AlarmMessageAckPushType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (timeStamp)
	lengthInBits += m.TimeStamp.GetLengthInBits(ctx)

	// Simple field (functionId)
	lengthInBits += 8

	// Simple field (numberOfObjects)
	lengthInBits += 8

	// Array field
	if len(m.MessageObjects) > 0 {
		for _curItem, element := range m.MessageObjects {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.MessageObjects), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_AlarmMessageAckPushType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func AlarmMessageAckPushTypeParse(ctx context.Context, theBytes []byte) (AlarmMessageAckPushType, error) {
	return AlarmMessageAckPushTypeParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func AlarmMessageAckPushTypeParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (AlarmMessageAckPushType, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (AlarmMessageAckPushType, error) {
		return AlarmMessageAckPushTypeParseWithBuffer(ctx, readBuffer)
	}
}

func AlarmMessageAckPushTypeParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (AlarmMessageAckPushType, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AlarmMessageAckPushType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AlarmMessageAckPushType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	timeStamp, err := ReadSimpleField[DateAndTime](ctx, "timeStamp", ReadComplex[DateAndTime](DateAndTimeParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'timeStamp' field"))
	}

	functionId, err := ReadSimpleField(ctx, "functionId", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'functionId' field"))
	}

	numberOfObjects, err := ReadSimpleField(ctx, "numberOfObjects", ReadUnsignedByte(readBuffer, uint8(8)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'numberOfObjects' field"))
	}

	messageObjects, err := ReadCountArrayField[AlarmMessageAckObjectPushType](ctx, "messageObjects", ReadComplex[AlarmMessageAckObjectPushType](AlarmMessageAckObjectPushTypeParseWithBuffer, readBuffer), uint64(numberOfObjects))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'messageObjects' field"))
	}

	if closeErr := readBuffer.CloseContext("AlarmMessageAckPushType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AlarmMessageAckPushType")
	}

	// Create the instance
	return &_AlarmMessageAckPushType{
		TimeStamp:       timeStamp,
		FunctionId:      functionId,
		NumberOfObjects: numberOfObjects,
		MessageObjects:  messageObjects,
	}, nil
}

func (m *_AlarmMessageAckPushType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_AlarmMessageAckPushType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("AlarmMessageAckPushType"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for AlarmMessageAckPushType")
	}

	if err := WriteSimpleField[DateAndTime](ctx, "timeStamp", m.GetTimeStamp(), WriteComplex[DateAndTime](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'timeStamp' field")
	}

	if err := WriteSimpleField[uint8](ctx, "functionId", m.GetFunctionId(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'functionId' field")
	}

	if err := WriteSimpleField[uint8](ctx, "numberOfObjects", m.GetNumberOfObjects(), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'numberOfObjects' field")
	}

	if err := WriteComplexTypeArrayField(ctx, "messageObjects", m.GetMessageObjects(), writeBuffer); err != nil {
		return errors.Wrap(err, "Error serializing 'messageObjects' field")
	}

	if popErr := writeBuffer.PopContext("AlarmMessageAckPushType"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for AlarmMessageAckPushType")
	}
	return nil
}

func (m *_AlarmMessageAckPushType) isAlarmMessageAckPushType() bool {
	return true
}

func (m *_AlarmMessageAckPushType) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
