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

// BACnetLogDataLogDataEntryBooleanValue is the corresponding interface of BACnetLogDataLogDataEntryBooleanValue
type BACnetLogDataLogDataEntryBooleanValue interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetLogDataLogDataEntry
	// GetBooleanValue returns BooleanValue (property field)
	GetBooleanValue() BACnetContextTagBoolean
}

// BACnetLogDataLogDataEntryBooleanValueExactly can be used when we want exactly this type and not a type which fulfills BACnetLogDataLogDataEntryBooleanValue.
// This is useful for switch cases.
type BACnetLogDataLogDataEntryBooleanValueExactly interface {
	BACnetLogDataLogDataEntryBooleanValue
	isBACnetLogDataLogDataEntryBooleanValue() bool
}

// _BACnetLogDataLogDataEntryBooleanValue is the data-structure of this message
type _BACnetLogDataLogDataEntryBooleanValue struct {
	*_BACnetLogDataLogDataEntry
	BooleanValue BACnetContextTagBoolean
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetLogDataLogDataEntryBooleanValue) InitializeParent(parent BACnetLogDataLogDataEntry, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetLogDataLogDataEntryBooleanValue) GetParent() BACnetLogDataLogDataEntry {
	return m._BACnetLogDataLogDataEntry
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetLogDataLogDataEntryBooleanValue) GetBooleanValue() BACnetContextTagBoolean {
	return m.BooleanValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetLogDataLogDataEntryBooleanValue factory function for _BACnetLogDataLogDataEntryBooleanValue
func NewBACnetLogDataLogDataEntryBooleanValue(booleanValue BACnetContextTagBoolean, peekedTagHeader BACnetTagHeader) *_BACnetLogDataLogDataEntryBooleanValue {
	_result := &_BACnetLogDataLogDataEntryBooleanValue{
		BooleanValue:               booleanValue,
		_BACnetLogDataLogDataEntry: NewBACnetLogDataLogDataEntry(peekedTagHeader),
	}
	_result._BACnetLogDataLogDataEntry._BACnetLogDataLogDataEntryChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetLogDataLogDataEntryBooleanValue(structType any) BACnetLogDataLogDataEntryBooleanValue {
	if casted, ok := structType.(BACnetLogDataLogDataEntryBooleanValue); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetLogDataLogDataEntryBooleanValue); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetLogDataLogDataEntryBooleanValue) GetTypeName() string {
	return "BACnetLogDataLogDataEntryBooleanValue"
}

func (m *_BACnetLogDataLogDataEntryBooleanValue) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (booleanValue)
	lengthInBits += m.BooleanValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetLogDataLogDataEntryBooleanValue) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetLogDataLogDataEntryBooleanValueParse(ctx context.Context, theBytes []byte) (BACnetLogDataLogDataEntryBooleanValue, error) {
	return BACnetLogDataLogDataEntryBooleanValueParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func BACnetLogDataLogDataEntryBooleanValueParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLogDataLogDataEntryBooleanValue, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLogDataLogDataEntryBooleanValue, error) {
		return BACnetLogDataLogDataEntryBooleanValueParseWithBuffer(ctx, readBuffer)
	}
}

func BACnetLogDataLogDataEntryBooleanValueParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetLogDataLogDataEntryBooleanValue, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetLogDataLogDataEntryBooleanValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetLogDataLogDataEntryBooleanValue")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	booleanValue, err := ReadSimpleField[BACnetContextTagBoolean](ctx, "booleanValue", ReadComplex[BACnetContextTagBoolean](BACnetContextTagParseWithBufferProducer[BACnetContextTagBoolean]((uint8)(uint8(0)), (BACnetDataType)(BACnetDataType_BOOLEAN)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'booleanValue' field"))
	}

	if closeErr := readBuffer.CloseContext("BACnetLogDataLogDataEntryBooleanValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetLogDataLogDataEntryBooleanValue")
	}

	// Create a partially initialized instance
	_child := &_BACnetLogDataLogDataEntryBooleanValue{
		_BACnetLogDataLogDataEntry: &_BACnetLogDataLogDataEntry{},
		BooleanValue:               booleanValue,
	}
	_child._BACnetLogDataLogDataEntry._BACnetLogDataLogDataEntryChildRequirements = _child
	return _child, nil
}

func (m *_BACnetLogDataLogDataEntryBooleanValue) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetLogDataLogDataEntryBooleanValue) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetLogDataLogDataEntryBooleanValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetLogDataLogDataEntryBooleanValue")
		}

		if err := WriteSimpleField[BACnetContextTagBoolean](ctx, "booleanValue", m.GetBooleanValue(), WriteComplex[BACnetContextTagBoolean](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'booleanValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetLogDataLogDataEntryBooleanValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetLogDataLogDataEntryBooleanValue")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetLogDataLogDataEntryBooleanValue) isBACnetLogDataLogDataEntryBooleanValue() bool {
	return true
}

func (m *_BACnetLogDataLogDataEntryBooleanValue) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
