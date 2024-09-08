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

// BACnetEventParameterChangeOfValueCivCriteriaBitmask is the corresponding interface of BACnetEventParameterChangeOfValueCivCriteriaBitmask
type BACnetEventParameterChangeOfValueCivCriteriaBitmask interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetEventParameterChangeOfValueCivCriteria
	// GetBitmask returns Bitmask (property field)
	GetBitmask() BACnetContextTagBitString
	// IsBACnetEventParameterChangeOfValueCivCriteriaBitmask is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetEventParameterChangeOfValueCivCriteriaBitmask()
}

// _BACnetEventParameterChangeOfValueCivCriteriaBitmask is the data-structure of this message
type _BACnetEventParameterChangeOfValueCivCriteriaBitmask struct {
	BACnetEventParameterChangeOfValueCivCriteriaContract
	Bitmask BACnetContextTagBitString
}

var _ BACnetEventParameterChangeOfValueCivCriteriaBitmask = (*_BACnetEventParameterChangeOfValueCivCriteriaBitmask)(nil)
var _ BACnetEventParameterChangeOfValueCivCriteriaRequirements = (*_BACnetEventParameterChangeOfValueCivCriteriaBitmask)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetEventParameterChangeOfValueCivCriteriaBitmask) GetParent() BACnetEventParameterChangeOfValueCivCriteriaContract {
	return m.BACnetEventParameterChangeOfValueCivCriteriaContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetEventParameterChangeOfValueCivCriteriaBitmask) GetBitmask() BACnetContextTagBitString {
	return m.Bitmask
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetEventParameterChangeOfValueCivCriteriaBitmask factory function for _BACnetEventParameterChangeOfValueCivCriteriaBitmask
func NewBACnetEventParameterChangeOfValueCivCriteriaBitmask(bitmask BACnetContextTagBitString, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8) *_BACnetEventParameterChangeOfValueCivCriteriaBitmask {
	if bitmask == nil {
		panic("bitmask of type BACnetContextTagBitString for BACnetEventParameterChangeOfValueCivCriteriaBitmask must not be nil")
	}
	_result := &_BACnetEventParameterChangeOfValueCivCriteriaBitmask{
		BACnetEventParameterChangeOfValueCivCriteriaContract: NewBACnetEventParameterChangeOfValueCivCriteria(openingTag, peekedTagHeader, closingTag, tagNumber),
		Bitmask: bitmask,
	}
	_result.BACnetEventParameterChangeOfValueCivCriteriaContract.(*_BACnetEventParameterChangeOfValueCivCriteria)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetEventParameterChangeOfValueCivCriteriaBitmask(structType any) BACnetEventParameterChangeOfValueCivCriteriaBitmask {
	if casted, ok := structType.(BACnetEventParameterChangeOfValueCivCriteriaBitmask); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetEventParameterChangeOfValueCivCriteriaBitmask); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetEventParameterChangeOfValueCivCriteriaBitmask) GetTypeName() string {
	return "BACnetEventParameterChangeOfValueCivCriteriaBitmask"
}

func (m *_BACnetEventParameterChangeOfValueCivCriteriaBitmask) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetEventParameterChangeOfValueCivCriteriaContract.(*_BACnetEventParameterChangeOfValueCivCriteria).getLengthInBits(ctx))

	// Simple field (bitmask)
	lengthInBits += m.Bitmask.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetEventParameterChangeOfValueCivCriteriaBitmask) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetEventParameterChangeOfValueCivCriteriaBitmask) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetEventParameterChangeOfValueCivCriteria, tagNumber uint8) (__bACnetEventParameterChangeOfValueCivCriteriaBitmask BACnetEventParameterChangeOfValueCivCriteriaBitmask, err error) {
	m.BACnetEventParameterChangeOfValueCivCriteriaContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetEventParameterChangeOfValueCivCriteriaBitmask"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetEventParameterChangeOfValueCivCriteriaBitmask")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	bitmask, err := ReadSimpleField[BACnetContextTagBitString](ctx, "bitmask", ReadComplex[BACnetContextTagBitString](BACnetContextTagParseWithBufferProducer[BACnetContextTagBitString]((uint8)(uint8(0)), (BACnetDataType)(BACnetDataType_BIT_STRING)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'bitmask' field"))
	}
	m.Bitmask = bitmask

	if closeErr := readBuffer.CloseContext("BACnetEventParameterChangeOfValueCivCriteriaBitmask"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetEventParameterChangeOfValueCivCriteriaBitmask")
	}

	return m, nil
}

func (m *_BACnetEventParameterChangeOfValueCivCriteriaBitmask) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetEventParameterChangeOfValueCivCriteriaBitmask) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetEventParameterChangeOfValueCivCriteriaBitmask"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetEventParameterChangeOfValueCivCriteriaBitmask")
		}

		if err := WriteSimpleField[BACnetContextTagBitString](ctx, "bitmask", m.GetBitmask(), WriteComplex[BACnetContextTagBitString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'bitmask' field")
		}

		if popErr := writeBuffer.PopContext("BACnetEventParameterChangeOfValueCivCriteriaBitmask"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetEventParameterChangeOfValueCivCriteriaBitmask")
		}
		return nil
	}
	return m.BACnetEventParameterChangeOfValueCivCriteriaContract.(*_BACnetEventParameterChangeOfValueCivCriteria).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetEventParameterChangeOfValueCivCriteriaBitmask) IsBACnetEventParameterChangeOfValueCivCriteriaBitmask() {
}

func (m *_BACnetEventParameterChangeOfValueCivCriteriaBitmask) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
