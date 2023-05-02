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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetConstructedDataIntegerValueHighLimit is the corresponding interface of BACnetConstructedDataIntegerValueHighLimit
type BACnetConstructedDataIntegerValueHighLimit interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetHighLimit returns HighLimit (property field)
	GetHighLimit() BACnetApplicationTagSignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagSignedInteger
}

// BACnetConstructedDataIntegerValueHighLimitExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataIntegerValueHighLimit.
// This is useful for switch cases.
type BACnetConstructedDataIntegerValueHighLimitExactly interface {
	BACnetConstructedDataIntegerValueHighLimit
	isBACnetConstructedDataIntegerValueHighLimit() bool
}

// _BACnetConstructedDataIntegerValueHighLimit is the data-structure of this message
type _BACnetConstructedDataIntegerValueHighLimit struct {
	*_BACnetConstructedData
	HighLimit BACnetApplicationTagSignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataIntegerValueHighLimit) GetObjectTypeArgument() BACnetObjectType {
	return BACnetObjectType_INTEGER_VALUE
}

func (m *_BACnetConstructedDataIntegerValueHighLimit) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_HIGH_LIMIT
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataIntegerValueHighLimit) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataIntegerValueHighLimit) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataIntegerValueHighLimit) GetHighLimit() BACnetApplicationTagSignedInteger {
	return m.HighLimit
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataIntegerValueHighLimit) GetActualValue() BACnetApplicationTagSignedInteger {
	ctx := context.Background()
	_ = ctx
	return CastBACnetApplicationTagSignedInteger(m.GetHighLimit())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataIntegerValueHighLimit factory function for _BACnetConstructedDataIntegerValueHighLimit
func NewBACnetConstructedDataIntegerValueHighLimit(highLimit BACnetApplicationTagSignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataIntegerValueHighLimit {
	_result := &_BACnetConstructedDataIntegerValueHighLimit{
		HighLimit:              highLimit,
		_BACnetConstructedData: NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataIntegerValueHighLimit(structType any) BACnetConstructedDataIntegerValueHighLimit {
	if casted, ok := structType.(BACnetConstructedDataIntegerValueHighLimit); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataIntegerValueHighLimit); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataIntegerValueHighLimit) GetTypeName() string {
	return "BACnetConstructedDataIntegerValueHighLimit"
}

func (m *_BACnetConstructedDataIntegerValueHighLimit) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (highLimit)
	lengthInBits += m.HighLimit.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataIntegerValueHighLimit) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetConstructedDataIntegerValueHighLimitParse(theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataIntegerValueHighLimit, error) {
	return BACnetConstructedDataIntegerValueHighLimitParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataIntegerValueHighLimitParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataIntegerValueHighLimit, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataIntegerValueHighLimit"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataIntegerValueHighLimit")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (highLimit)
	if pullErr := readBuffer.PullContext("highLimit"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for highLimit")
	}
	_highLimit, _highLimitErr := BACnetApplicationTagParseWithBuffer(ctx, readBuffer)
	if _highLimitErr != nil {
		return nil, errors.Wrap(_highLimitErr, "Error parsing 'highLimit' field of BACnetConstructedDataIntegerValueHighLimit")
	}
	highLimit := _highLimit.(BACnetApplicationTagSignedInteger)
	if closeErr := readBuffer.CloseContext("highLimit"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for highLimit")
	}

	// Virtual field
	_actualValue := highLimit
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataIntegerValueHighLimit"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataIntegerValueHighLimit")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataIntegerValueHighLimit{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		HighLimit: highLimit,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataIntegerValueHighLimit) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataIntegerValueHighLimit) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataIntegerValueHighLimit"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataIntegerValueHighLimit")
		}

		// Simple Field (highLimit)
		if pushErr := writeBuffer.PushContext("highLimit"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for highLimit")
		}
		_highLimitErr := writeBuffer.WriteSerializable(ctx, m.GetHighLimit())
		if popErr := writeBuffer.PopContext("highLimit"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for highLimit")
		}
		if _highLimitErr != nil {
			return errors.Wrap(_highLimitErr, "Error serializing 'highLimit' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataIntegerValueHighLimit"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataIntegerValueHighLimit")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataIntegerValueHighLimit) isBACnetConstructedDataIntegerValueHighLimit() bool {
	return true
}

func (m *_BACnetConstructedDataIntegerValueHighLimit) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
