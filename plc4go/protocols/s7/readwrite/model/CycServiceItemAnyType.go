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

// CycServiceItemAnyType is the corresponding interface of CycServiceItemAnyType
type CycServiceItemAnyType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	CycServiceItemType
	// GetTransportSize returns TransportSize (property field)
	GetTransportSize() TransportSize
	// GetLength returns Length (property field)
	GetLength() uint16
	// GetDbNumber returns DbNumber (property field)
	GetDbNumber() uint16
	// GetMemoryArea returns MemoryArea (property field)
	GetMemoryArea() MemoryArea
	// GetAddress returns Address (property field)
	GetAddress() uint32
}

// CycServiceItemAnyTypeExactly can be used when we want exactly this type and not a type which fulfills CycServiceItemAnyType.
// This is useful for switch cases.
type CycServiceItemAnyTypeExactly interface {
	CycServiceItemAnyType
	isCycServiceItemAnyType() bool
}

// _CycServiceItemAnyType is the data-structure of this message
type _CycServiceItemAnyType struct {
	*_CycServiceItemType
	TransportSize TransportSize
	Length        uint16
	DbNumber      uint16
	MemoryArea    MemoryArea
	Address       uint32
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_CycServiceItemAnyType) InitializeParent(parent CycServiceItemType, byteLength uint8, syntaxId uint8) {
	m.ByteLength = byteLength
	m.SyntaxId = syntaxId
}

func (m *_CycServiceItemAnyType) GetParent() CycServiceItemType {
	return m._CycServiceItemType
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CycServiceItemAnyType) GetTransportSize() TransportSize {
	return m.TransportSize
}

func (m *_CycServiceItemAnyType) GetLength() uint16 {
	return m.Length
}

func (m *_CycServiceItemAnyType) GetDbNumber() uint16 {
	return m.DbNumber
}

func (m *_CycServiceItemAnyType) GetMemoryArea() MemoryArea {
	return m.MemoryArea
}

func (m *_CycServiceItemAnyType) GetAddress() uint32 {
	return m.Address
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewCycServiceItemAnyType factory function for _CycServiceItemAnyType
func NewCycServiceItemAnyType(transportSize TransportSize, length uint16, dbNumber uint16, memoryArea MemoryArea, address uint32, byteLength uint8, syntaxId uint8) *_CycServiceItemAnyType {
	_result := &_CycServiceItemAnyType{
		TransportSize:       transportSize,
		Length:              length,
		DbNumber:            dbNumber,
		MemoryArea:          memoryArea,
		Address:             address,
		_CycServiceItemType: NewCycServiceItemType(byteLength, syntaxId),
	}
	_result._CycServiceItemType._CycServiceItemTypeChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastCycServiceItemAnyType(structType any) CycServiceItemAnyType {
	if casted, ok := structType.(CycServiceItemAnyType); ok {
		return casted
	}
	if casted, ok := structType.(*CycServiceItemAnyType); ok {
		return *casted
	}
	return nil
}

func (m *_CycServiceItemAnyType) GetTypeName() string {
	return "CycServiceItemAnyType"
}

func (m *_CycServiceItemAnyType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Enum Field (transportSize)
	lengthInBits += 8

	// Simple field (length)
	lengthInBits += 16

	// Simple field (dbNumber)
	lengthInBits += 16

	// Simple field (memoryArea)
	lengthInBits += 8

	// Simple field (address)
	lengthInBits += 24

	return lengthInBits
}

func (m *_CycServiceItemAnyType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func CycServiceItemAnyTypeParse(ctx context.Context, theBytes []byte) (CycServiceItemAnyType, error) {
	return CycServiceItemAnyTypeParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func CycServiceItemAnyTypeParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (CycServiceItemAnyType, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (CycServiceItemAnyType, error) {
		return CycServiceItemAnyTypeParseWithBuffer(ctx, readBuffer)
	}
}

func CycServiceItemAnyTypeParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (CycServiceItemAnyType, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CycServiceItemAnyType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CycServiceItemAnyType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	transportSize, err := ReadEnumField[TransportSize](ctx, "transportSize", "TransportSize", ReadEnum[TransportSize, uint8](TransportSizeFirstEnumForFieldCode, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'transportSize' field"))
	}

	length, err := ReadSimpleField(ctx, "length", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'length' field"))
	}

	dbNumber, err := ReadSimpleField(ctx, "dbNumber", ReadUnsignedShort(readBuffer, uint8(16)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'dbNumber' field"))
	}

	memoryArea, err := ReadEnumField[MemoryArea](ctx, "memoryArea", "MemoryArea", ReadEnum(MemoryAreaByValue, ReadUnsignedByte(readBuffer, uint8(8))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'memoryArea' field"))
	}

	address, err := ReadSimpleField(ctx, "address", ReadUnsignedInt(readBuffer, uint8(24)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'address' field"))
	}

	if closeErr := readBuffer.CloseContext("CycServiceItemAnyType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CycServiceItemAnyType")
	}

	// Create a partially initialized instance
	_child := &_CycServiceItemAnyType{
		_CycServiceItemType: &_CycServiceItemType{},
		TransportSize:       transportSize,
		Length:              length,
		DbNumber:            dbNumber,
		MemoryArea:          memoryArea,
		Address:             address,
	}
	_child._CycServiceItemType._CycServiceItemTypeChildRequirements = _child
	return _child, nil
}

func (m *_CycServiceItemAnyType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CycServiceItemAnyType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CycServiceItemAnyType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CycServiceItemAnyType")
		}

		if err := WriteEnumField(ctx, "transportSize", "TransportSize", m.GetTransportSize(), WriteEnum[TransportSize, uint8](TransportSize.GetCode, TransportSize.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
			return errors.Wrap(err, "Error serializing 'transportSize' field")
		}

		if err := WriteSimpleField[uint16](ctx, "length", m.GetLength(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'length' field")
		}

		if err := WriteSimpleField[uint16](ctx, "dbNumber", m.GetDbNumber(), WriteUnsignedShort(writeBuffer, 16)); err != nil {
			return errors.Wrap(err, "Error serializing 'dbNumber' field")
		}

		if err := WriteSimpleEnumField[MemoryArea](ctx, "memoryArea", "MemoryArea", m.GetMemoryArea(), WriteEnum[MemoryArea, uint8](MemoryArea.GetValue, MemoryArea.PLC4XEnumName, WriteUnsignedByte(writeBuffer, 8))); err != nil {
			return errors.Wrap(err, "Error serializing 'memoryArea' field")
		}

		if err := WriteSimpleField[uint32](ctx, "address", m.GetAddress(), WriteUnsignedInt(writeBuffer, 24)); err != nil {
			return errors.Wrap(err, "Error serializing 'address' field")
		}

		if popErr := writeBuffer.PopContext("CycServiceItemAnyType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CycServiceItemAnyType")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_CycServiceItemAnyType) isCycServiceItemAnyType() bool {
	return true
}

func (m *_CycServiceItemAnyType) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
