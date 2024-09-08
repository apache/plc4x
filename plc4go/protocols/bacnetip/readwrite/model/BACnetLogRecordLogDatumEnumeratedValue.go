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

// BACnetLogRecordLogDatumEnumeratedValue is the corresponding interface of BACnetLogRecordLogDatumEnumeratedValue
type BACnetLogRecordLogDatumEnumeratedValue interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetLogRecordLogDatum
	// GetEnumeratedValue returns EnumeratedValue (property field)
	GetEnumeratedValue() BACnetContextTagEnumerated
	// IsBACnetLogRecordLogDatumEnumeratedValue is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetLogRecordLogDatumEnumeratedValue()
}

// _BACnetLogRecordLogDatumEnumeratedValue is the data-structure of this message
type _BACnetLogRecordLogDatumEnumeratedValue struct {
	BACnetLogRecordLogDatumContract
	EnumeratedValue BACnetContextTagEnumerated
}

var _ BACnetLogRecordLogDatumEnumeratedValue = (*_BACnetLogRecordLogDatumEnumeratedValue)(nil)
var _ BACnetLogRecordLogDatumRequirements = (*_BACnetLogRecordLogDatumEnumeratedValue)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetLogRecordLogDatumEnumeratedValue) GetParent() BACnetLogRecordLogDatumContract {
	return m.BACnetLogRecordLogDatumContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetLogRecordLogDatumEnumeratedValue) GetEnumeratedValue() BACnetContextTagEnumerated {
	return m.EnumeratedValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetLogRecordLogDatumEnumeratedValue factory function for _BACnetLogRecordLogDatumEnumeratedValue
func NewBACnetLogRecordLogDatumEnumeratedValue(enumeratedValue BACnetContextTagEnumerated, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8) *_BACnetLogRecordLogDatumEnumeratedValue {
	if enumeratedValue == nil {
		panic("enumeratedValue of type BACnetContextTagEnumerated for BACnetLogRecordLogDatumEnumeratedValue must not be nil")
	}
	_result := &_BACnetLogRecordLogDatumEnumeratedValue{
		BACnetLogRecordLogDatumContract: NewBACnetLogRecordLogDatum(openingTag, peekedTagHeader, closingTag, tagNumber),
		EnumeratedValue:                 enumeratedValue,
	}
	_result.BACnetLogRecordLogDatumContract.(*_BACnetLogRecordLogDatum)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetLogRecordLogDatumEnumeratedValue(structType any) BACnetLogRecordLogDatumEnumeratedValue {
	if casted, ok := structType.(BACnetLogRecordLogDatumEnumeratedValue); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetLogRecordLogDatumEnumeratedValue); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetLogRecordLogDatumEnumeratedValue) GetTypeName() string {
	return "BACnetLogRecordLogDatumEnumeratedValue"
}

func (m *_BACnetLogRecordLogDatumEnumeratedValue) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetLogRecordLogDatumContract.(*_BACnetLogRecordLogDatum).getLengthInBits(ctx))

	// Simple field (enumeratedValue)
	lengthInBits += m.EnumeratedValue.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetLogRecordLogDatumEnumeratedValue) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetLogRecordLogDatumEnumeratedValue) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetLogRecordLogDatum, tagNumber uint8) (__bACnetLogRecordLogDatumEnumeratedValue BACnetLogRecordLogDatumEnumeratedValue, err error) {
	m.BACnetLogRecordLogDatumContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetLogRecordLogDatumEnumeratedValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetLogRecordLogDatumEnumeratedValue")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	enumeratedValue, err := ReadSimpleField[BACnetContextTagEnumerated](ctx, "enumeratedValue", ReadComplex[BACnetContextTagEnumerated](BACnetContextTagParseWithBufferProducer[BACnetContextTagEnumerated]((uint8)(uint8(3)), (BACnetDataType)(BACnetDataType_ENUMERATED)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'enumeratedValue' field"))
	}
	m.EnumeratedValue = enumeratedValue

	if closeErr := readBuffer.CloseContext("BACnetLogRecordLogDatumEnumeratedValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetLogRecordLogDatumEnumeratedValue")
	}

	return m, nil
}

func (m *_BACnetLogRecordLogDatumEnumeratedValue) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetLogRecordLogDatumEnumeratedValue) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetLogRecordLogDatumEnumeratedValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetLogRecordLogDatumEnumeratedValue")
		}

		if err := WriteSimpleField[BACnetContextTagEnumerated](ctx, "enumeratedValue", m.GetEnumeratedValue(), WriteComplex[BACnetContextTagEnumerated](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'enumeratedValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetLogRecordLogDatumEnumeratedValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetLogRecordLogDatumEnumeratedValue")
		}
		return nil
	}
	return m.BACnetLogRecordLogDatumContract.(*_BACnetLogRecordLogDatum).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetLogRecordLogDatumEnumeratedValue) IsBACnetLogRecordLogDatumEnumeratedValue() {}

func (m *_BACnetLogRecordLogDatumEnumeratedValue) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
