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

// BACnetLogDataLogDataEntryBitStringValue is the corresponding interface of BACnetLogDataLogDataEntryBitStringValue
type BACnetLogDataLogDataEntryBitStringValue interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetLogDataLogDataEntry
	// GetBitStringValue returns BitStringValue (property field)
	GetBitStringValue() BACnetContextTagBitString
}

// BACnetLogDataLogDataEntryBitStringValueExactly can be used when we want exactly this type and not a type which fulfills BACnetLogDataLogDataEntryBitStringValue.
// This is useful for switch cases.
type BACnetLogDataLogDataEntryBitStringValueExactly interface {
	BACnetLogDataLogDataEntryBitStringValue
	isBACnetLogDataLogDataEntryBitStringValue() bool
}

// _BACnetLogDataLogDataEntryBitStringValue is the data-structure of this message
type _BACnetLogDataLogDataEntryBitStringValue struct {
	*_BACnetLogDataLogDataEntry
	BitStringValue BACnetContextTagBitString
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetLogDataLogDataEntryBitStringValue) InitializeParent(parent BACnetLogDataLogDataEntry, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetLogDataLogDataEntryBitStringValue) GetParent() BACnetLogDataLogDataEntry {
	return m._BACnetLogDataLogDataEntry
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetLogDataLogDataEntryBitStringValue) GetBitStringValue() BACnetContextTagBitString {
	return m.BitStringValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetLogDataLogDataEntryBitStringValue factory function for _BACnetLogDataLogDataEntryBitStringValue
func NewBACnetLogDataLogDataEntryBitStringValue(bitStringValue BACnetContextTagBitString, peekedTagHeader BACnetTagHeader) *_BACnetLogDataLogDataEntryBitStringValue {
	_result := &_BACnetLogDataLogDataEntryBitStringValue{
		BitStringValue:             bitStringValue,
		_BACnetLogDataLogDataEntry: NewBACnetLogDataLogDataEntry(peekedTagHeader),
	}
	_result._BACnetLogDataLogDataEntry._BACnetLogDataLogDataEntryChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetLogDataLogDataEntryBitStringValue(structType any) BACnetLogDataLogDataEntryBitStringValue {
	if casted, ok := structType.(BACnetLogDataLogDataEntryBitStringValue); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetLogDataLogDataEntryBitStringValue); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetLogDataLogDataEntryBitStringValue) GetTypeName() string {
	return "BACnetLogDataLogDataEntryBitStringValue"
}

func (m *_BACnetLogDataLogDataEntryBitStringValue) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (bitStringValue)
	lengthInBits += m.BitStringValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetLogDataLogDataEntryBitStringValue) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetLogDataLogDataEntryBitStringValueParse(ctx context.Context, theBytes []byte) (BACnetLogDataLogDataEntryBitStringValue, error) {
	return BACnetLogDataLogDataEntryBitStringValueParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func BACnetLogDataLogDataEntryBitStringValueParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLogDataLogDataEntryBitStringValue, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLogDataLogDataEntryBitStringValue, error) {
		return BACnetLogDataLogDataEntryBitStringValueParseWithBuffer(ctx, readBuffer)
	}
}

func BACnetLogDataLogDataEntryBitStringValueParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLogDataLogDataEntryBitStringValue, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetLogDataLogDataEntryBitStringValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetLogDataLogDataEntryBitStringValue")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	bitStringValue, err := ReadSimpleField[BACnetContextTagBitString](ctx, "bitStringValue", ReadComplex[BACnetContextTagBitString](BACnetContextTagParseWithBufferProducer[BACnetContextTagBitString]((uint8)(uint8(5)), (BACnetDataType)(BACnetDataType_BIT_STRING)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'bitStringValue' field"))
	}

	if closeErr := readBuffer.CloseContext("BACnetLogDataLogDataEntryBitStringValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetLogDataLogDataEntryBitStringValue")
	}

	// Create a partially initialized instance
	_child := &_BACnetLogDataLogDataEntryBitStringValue{
		_BACnetLogDataLogDataEntry: &_BACnetLogDataLogDataEntry{},
		BitStringValue:             bitStringValue,
	}
	_child._BACnetLogDataLogDataEntry._BACnetLogDataLogDataEntryChildRequirements = _child
	return _child, nil
}

func (m *_BACnetLogDataLogDataEntryBitStringValue) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetLogDataLogDataEntryBitStringValue) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetLogDataLogDataEntryBitStringValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetLogDataLogDataEntryBitStringValue")
		}

		if err := WriteSimpleField[BACnetContextTagBitString](ctx, "bitStringValue", m.GetBitStringValue(), WriteComplex[BACnetContextTagBitString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'bitStringValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetLogDataLogDataEntryBitStringValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetLogDataLogDataEntryBitStringValue")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetLogDataLogDataEntryBitStringValue) isBACnetLogDataLogDataEntryBitStringValue() bool {
	return true
}

func (m *_BACnetLogDataLogDataEntryBitStringValue) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
