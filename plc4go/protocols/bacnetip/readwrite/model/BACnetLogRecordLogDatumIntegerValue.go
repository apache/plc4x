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

// BACnetLogRecordLogDatumIntegerValue is the corresponding interface of BACnetLogRecordLogDatumIntegerValue
type BACnetLogRecordLogDatumIntegerValue interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetLogRecordLogDatum
	// GetIntegerValue returns IntegerValue (property field)
	GetIntegerValue() BACnetContextTagSignedInteger
	// IsBACnetLogRecordLogDatumIntegerValue is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetLogRecordLogDatumIntegerValue()
}

// _BACnetLogRecordLogDatumIntegerValue is the data-structure of this message
type _BACnetLogRecordLogDatumIntegerValue struct {
	BACnetLogRecordLogDatumContract
	IntegerValue BACnetContextTagSignedInteger
}

var _ BACnetLogRecordLogDatumIntegerValue = (*_BACnetLogRecordLogDatumIntegerValue)(nil)
var _ BACnetLogRecordLogDatumRequirements = (*_BACnetLogRecordLogDatumIntegerValue)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetLogRecordLogDatumIntegerValue) GetParent() BACnetLogRecordLogDatumContract {
	return m.BACnetLogRecordLogDatumContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetLogRecordLogDatumIntegerValue) GetIntegerValue() BACnetContextTagSignedInteger {
	return m.IntegerValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetLogRecordLogDatumIntegerValue factory function for _BACnetLogRecordLogDatumIntegerValue
func NewBACnetLogRecordLogDatumIntegerValue(integerValue BACnetContextTagSignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8) *_BACnetLogRecordLogDatumIntegerValue {
	if integerValue == nil {
		panic("integerValue of type BACnetContextTagSignedInteger for BACnetLogRecordLogDatumIntegerValue must not be nil")
	}
	_result := &_BACnetLogRecordLogDatumIntegerValue{
		BACnetLogRecordLogDatumContract: NewBACnetLogRecordLogDatum(openingTag, peekedTagHeader, closingTag, tagNumber),
		IntegerValue:                    integerValue,
	}
	_result.BACnetLogRecordLogDatumContract.(*_BACnetLogRecordLogDatum)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetLogRecordLogDatumIntegerValue(structType any) BACnetLogRecordLogDatumIntegerValue {
	if casted, ok := structType.(BACnetLogRecordLogDatumIntegerValue); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetLogRecordLogDatumIntegerValue); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetLogRecordLogDatumIntegerValue) GetTypeName() string {
	return "BACnetLogRecordLogDatumIntegerValue"
}

func (m *_BACnetLogRecordLogDatumIntegerValue) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetLogRecordLogDatumContract.(*_BACnetLogRecordLogDatum).getLengthInBits(ctx))

	// Simple field (integerValue)
	lengthInBits += m.IntegerValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetLogRecordLogDatumIntegerValue) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetLogRecordLogDatumIntegerValue) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetLogRecordLogDatum, tagNumber uint8) (__bACnetLogRecordLogDatumIntegerValue BACnetLogRecordLogDatumIntegerValue, err error) {
	m.BACnetLogRecordLogDatumContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetLogRecordLogDatumIntegerValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetLogRecordLogDatumIntegerValue")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	integerValue, err := ReadSimpleField[BACnetContextTagSignedInteger](ctx, "integerValue", ReadComplex[BACnetContextTagSignedInteger](BACnetContextTagParseWithBufferProducer[BACnetContextTagSignedInteger]((uint8)(uint8(5)), (BACnetDataType)(BACnetDataType_SIGNED_INTEGER)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'integerValue' field"))
	}
	m.IntegerValue = integerValue

	if closeErr := readBuffer.CloseContext("BACnetLogRecordLogDatumIntegerValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetLogRecordLogDatumIntegerValue")
	}

	return m, nil
}

func (m *_BACnetLogRecordLogDatumIntegerValue) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetLogRecordLogDatumIntegerValue) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetLogRecordLogDatumIntegerValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetLogRecordLogDatumIntegerValue")
		}

		if err := WriteSimpleField[BACnetContextTagSignedInteger](ctx, "integerValue", m.GetIntegerValue(), WriteComplex[BACnetContextTagSignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'integerValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetLogRecordLogDatumIntegerValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetLogRecordLogDatumIntegerValue")
		}
		return nil
	}
	return m.BACnetLogRecordLogDatumContract.(*_BACnetLogRecordLogDatum).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetLogRecordLogDatumIntegerValue) IsBACnetLogRecordLogDatumIntegerValue() {}

func (m *_BACnetLogRecordLogDatumIntegerValue) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
