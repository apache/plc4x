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

	"github.com/apache/plc4x/plc4go/spi/codegen"
	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// AmsString is the corresponding interface of AmsString
type AmsString interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetText returns Text (property field)
	GetText() string
}

// AmsStringExactly can be used when we want exactly this type and not a type which fulfills AmsString.
// This is useful for switch cases.
type AmsStringExactly interface {
	AmsString
	isAmsString() bool
}

// _AmsString is the data-structure of this message
type _AmsString struct {
	Text string
	// Reserved Fields
	reservedField0 *uint8
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AmsString) GetText() string {
	return m.Text
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewAmsString factory function for _AmsString
func NewAmsString(text string) *_AmsString {
	return &_AmsString{Text: text}
}

// Deprecated: use the interface for direct cast
func CastAmsString(structType any) AmsString {
	if casted, ok := structType.(AmsString); ok {
		return casted
	}
	if casted, ok := structType.(*AmsString); ok {
		return *casted
	}
	return nil
}

func (m *_AmsString) GetTypeName() string {
	return "AmsString"
}

func (m *_AmsString) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Implicit Field (strLen)
	lengthInBits += 16

	// Simple field (text)
	lengthInBits += uint16(int32(int32(8)) * int32((int32(uint16(uint16(len(m.GetText())))+uint16(uint16(1))) - int32(int32(1)))))

	// Reserved Field (reserved)
	lengthInBits += 8

	return lengthInBits
}

func (m *_AmsString) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func AmsStringParse(ctx context.Context, theBytes []byte) (AmsString, error) {
	return AmsStringParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func AmsStringParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (AmsString, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (AmsString, error) {
		return AmsStringParseWithBuffer(ctx, readBuffer)
	}
}

func AmsStringParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (AmsString, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AmsString"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AmsString")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	strLen, err := ReadImplicitField[uint16](ctx, "strLen", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'strLen' field"))
	}
	_ = strLen

	text, err := ReadSimpleField(ctx, "text", ReadString(readBuffer, uint32(int32(int32(8))*int32((int32(strLen)-int32(int32(1)))))), codegen.WithEncoding("UTF-8"))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'text' field"))
	}

	reservedField0, err := ReadReservedField(ctx, "reserved", ReadUnsignedByte(readBuffer, uint8(8)), uint8(0x00))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing reserved field"))
	}

	if closeErr := readBuffer.CloseContext("AmsString"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AmsString")
	}

	// Create the instance
	return &_AmsString{
		Text:           text,
		reservedField0: reservedField0,
	}, nil
}

func (m *_AmsString) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_AmsString) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("AmsString"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for AmsString")
	}
	strLen := uint16(uint16(uint16(len(m.GetText()))) + uint16(uint16(1)))
	if err := WriteImplicitField(ctx, "strLen", strLen, WriteUnsignedShort(writeBuffer, 16)); err != nil {
		return errors.Wrap(err, "Error serializing 'strLen' field")
	}

	if err := WriteSimpleField[string](ctx, "text", m.GetText(), WriteString(writeBuffer, int32(int32(int32(8))*int32((int32(uint16(uint16(len(m.GetText())))+uint16(uint16(1)))-int32(int32(1)))))), codegen.WithEncoding("UTF-8")); err != nil {
		return errors.Wrap(err, "Error serializing 'text' field")
	}

	if err := WriteReservedField[uint8](ctx, "reserved", uint8(0x00), WriteUnsignedByte(writeBuffer, 8)); err != nil {
		return errors.Wrap(err, "Error serializing 'reserved' field number 1")
	}

	if popErr := writeBuffer.PopContext("AmsString"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for AmsString")
	}
	return nil
}

func (m *_AmsString) isAmsString() bool {
	return true
}

func (m *_AmsString) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
