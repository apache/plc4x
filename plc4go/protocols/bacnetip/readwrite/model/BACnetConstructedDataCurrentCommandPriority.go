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

// BACnetConstructedDataCurrentCommandPriority is the corresponding interface of BACnetConstructedDataCurrentCommandPriority
type BACnetConstructedDataCurrentCommandPriority interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetCurrentCommandPriority returns CurrentCommandPriority (property field)
	GetCurrentCommandPriority() BACnetOptionalUnsigned
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetOptionalUnsigned
	// IsBACnetConstructedDataCurrentCommandPriority is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConstructedDataCurrentCommandPriority()
}

// _BACnetConstructedDataCurrentCommandPriority is the data-structure of this message
type _BACnetConstructedDataCurrentCommandPriority struct {
	BACnetConstructedDataContract
	CurrentCommandPriority BACnetOptionalUnsigned
}

var _ BACnetConstructedDataCurrentCommandPriority = (*_BACnetConstructedDataCurrentCommandPriority)(nil)
var _ BACnetConstructedDataRequirements = (*_BACnetConstructedDataCurrentCommandPriority)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataCurrentCommandPriority) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataCurrentCommandPriority) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_CURRENT_COMMAND_PRIORITY
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataCurrentCommandPriority) GetParent() BACnetConstructedDataContract {
	return m.BACnetConstructedDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataCurrentCommandPriority) GetCurrentCommandPriority() BACnetOptionalUnsigned {
	return m.CurrentCommandPriority
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataCurrentCommandPriority) GetActualValue() BACnetOptionalUnsigned {
	ctx := context.Background()
	_ = ctx
	return CastBACnetOptionalUnsigned(m.GetCurrentCommandPriority())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataCurrentCommandPriority factory function for _BACnetConstructedDataCurrentCommandPriority
func NewBACnetConstructedDataCurrentCommandPriority(currentCommandPriority BACnetOptionalUnsigned, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataCurrentCommandPriority {
	if currentCommandPriority == nil {
		panic("currentCommandPriority of type BACnetOptionalUnsigned for BACnetConstructedDataCurrentCommandPriority must not be nil")
	}
	_result := &_BACnetConstructedDataCurrentCommandPriority{
		BACnetConstructedDataContract: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
		CurrentCommandPriority:        currentCommandPriority,
	}
	_result.BACnetConstructedDataContract.(*_BACnetConstructedData)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataCurrentCommandPriority(structType any) BACnetConstructedDataCurrentCommandPriority {
	if casted, ok := structType.(BACnetConstructedDataCurrentCommandPriority); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataCurrentCommandPriority); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataCurrentCommandPriority) GetTypeName() string {
	return "BACnetConstructedDataCurrentCommandPriority"
}

func (m *_BACnetConstructedDataCurrentCommandPriority) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConstructedDataContract.(*_BACnetConstructedData).getLengthInBits(ctx))

	// Simple field (currentCommandPriority)
	lengthInBits += m.CurrentCommandPriority.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataCurrentCommandPriority) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConstructedDataCurrentCommandPriority) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConstructedData, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (__bACnetConstructedDataCurrentCommandPriority BACnetConstructedDataCurrentCommandPriority, err error) {
	m.BACnetConstructedDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataCurrentCommandPriority"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataCurrentCommandPriority")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	currentCommandPriority, err := ReadSimpleField[BACnetOptionalUnsigned](ctx, "currentCommandPriority", ReadComplex[BACnetOptionalUnsigned](BACnetOptionalUnsignedParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'currentCommandPriority' field"))
	}
	m.CurrentCommandPriority = currentCommandPriority

	actualValue, err := ReadVirtualField[BACnetOptionalUnsigned](ctx, "actualValue", (*BACnetOptionalUnsigned)(nil), currentCommandPriority)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataCurrentCommandPriority"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataCurrentCommandPriority")
	}

	return m, nil
}

func (m *_BACnetConstructedDataCurrentCommandPriority) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataCurrentCommandPriority) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataCurrentCommandPriority"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataCurrentCommandPriority")
		}

		if err := WriteSimpleField[BACnetOptionalUnsigned](ctx, "currentCommandPriority", m.GetCurrentCommandPriority(), WriteComplex[BACnetOptionalUnsigned](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'currentCommandPriority' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataCurrentCommandPriority"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataCurrentCommandPriority")
		}
		return nil
	}
	return m.BACnetConstructedDataContract.(*_BACnetConstructedData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataCurrentCommandPriority) IsBACnetConstructedDataCurrentCommandPriority() {
}

func (m *_BACnetConstructedDataCurrentCommandPriority) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
