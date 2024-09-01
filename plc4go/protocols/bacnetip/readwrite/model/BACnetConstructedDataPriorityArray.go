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

// BACnetConstructedDataPriorityArray is the corresponding interface of BACnetConstructedDataPriorityArray
type BACnetConstructedDataPriorityArray interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetPriorityArray returns PriorityArray (property field)
	GetPriorityArray() BACnetPriorityArray
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetPriorityArray
}

// BACnetConstructedDataPriorityArrayExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataPriorityArray.
// This is useful for switch cases.
type BACnetConstructedDataPriorityArrayExactly interface {
	BACnetConstructedDataPriorityArray
	isBACnetConstructedDataPriorityArray() bool
}

// _BACnetConstructedDataPriorityArray is the data-structure of this message
type _BACnetConstructedDataPriorityArray struct {
	*_BACnetConstructedData
	PriorityArray BACnetPriorityArray
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataPriorityArray) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataPriorityArray) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_PRIORITY_ARRAY
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataPriorityArray) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataPriorityArray) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataPriorityArray) GetPriorityArray() BACnetPriorityArray {
	return m.PriorityArray
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataPriorityArray) GetActualValue() BACnetPriorityArray {
	ctx := context.Background()
	_ = ctx
	return CastBACnetPriorityArray(m.GetPriorityArray())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataPriorityArray factory function for _BACnetConstructedDataPriorityArray
func NewBACnetConstructedDataPriorityArray(priorityArray BACnetPriorityArray, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataPriorityArray {
	_result := &_BACnetConstructedDataPriorityArray{
		PriorityArray:          priorityArray,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataPriorityArray(structType any) BACnetConstructedDataPriorityArray {
	if casted, ok := structType.(BACnetConstructedDataPriorityArray); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataPriorityArray); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataPriorityArray) GetTypeName() string {
	return "BACnetConstructedDataPriorityArray"
}

func (m *_BACnetConstructedDataPriorityArray) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (priorityArray)
	lengthInBits += m.PriorityArray.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataPriorityArray) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetConstructedDataPriorityArrayParse(ctx context.Context, theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataPriorityArray, error) {
	return BACnetConstructedDataPriorityArrayParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataPriorityArrayParseWithBufferProducer(tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConstructedDataPriorityArray, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConstructedDataPriorityArray, error) {
		return BACnetConstructedDataPriorityArrayParseWithBuffer(ctx, readBuffer, tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
	}
}

func BACnetConstructedDataPriorityArrayParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataPriorityArray, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataPriorityArray"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataPriorityArray")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	priorityArray, err := ReadSimpleField[BACnetPriorityArray](ctx, "priorityArray", ReadComplex[BACnetPriorityArray](BACnetPriorityArrayParseWithBufferProducer((BACnetObjectType)(objectTypeArgument), (uint8)(tagNumber), (BACnetTagPayloadUnsignedInteger)(arrayIndexArgument)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'priorityArray' field"))
	}

	actualValue, err := ReadVirtualField[BACnetPriorityArray](ctx, "actualValue", (*BACnetPriorityArray)(nil), priorityArray)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataPriorityArray"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataPriorityArray")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataPriorityArray{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		PriorityArray: priorityArray,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataPriorityArray) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataPriorityArray) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataPriorityArray"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataPriorityArray")
		}

		if err := WriteSimpleField[BACnetPriorityArray](ctx, "priorityArray", m.GetPriorityArray(), WriteComplex[BACnetPriorityArray](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'priorityArray' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataPriorityArray"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataPriorityArray")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataPriorityArray) isBACnetConstructedDataPriorityArray() bool {
	return true
}

func (m *_BACnetConstructedDataPriorityArray) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
