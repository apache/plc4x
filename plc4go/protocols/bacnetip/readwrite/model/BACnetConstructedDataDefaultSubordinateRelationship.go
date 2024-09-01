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

// BACnetConstructedDataDefaultSubordinateRelationship is the corresponding interface of BACnetConstructedDataDefaultSubordinateRelationship
type BACnetConstructedDataDefaultSubordinateRelationship interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetDefaultSubordinateRelationship returns DefaultSubordinateRelationship (property field)
	GetDefaultSubordinateRelationship() BACnetRelationshipTagged
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetRelationshipTagged
}

// BACnetConstructedDataDefaultSubordinateRelationshipExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataDefaultSubordinateRelationship.
// This is useful for switch cases.
type BACnetConstructedDataDefaultSubordinateRelationshipExactly interface {
	BACnetConstructedDataDefaultSubordinateRelationship
	isBACnetConstructedDataDefaultSubordinateRelationship() bool
}

// _BACnetConstructedDataDefaultSubordinateRelationship is the data-structure of this message
type _BACnetConstructedDataDefaultSubordinateRelationship struct {
	*_BACnetConstructedData
	DefaultSubordinateRelationship BACnetRelationshipTagged
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_DEFAULT_SUBORDINATE_RELATIONSHIP
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) GetDefaultSubordinateRelationship() BACnetRelationshipTagged {
	return m.DefaultSubordinateRelationship
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) GetActualValue() BACnetRelationshipTagged {
	ctx := context.Background()
	_ = ctx
	return CastBACnetRelationshipTagged(m.GetDefaultSubordinateRelationship())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataDefaultSubordinateRelationship factory function for _BACnetConstructedDataDefaultSubordinateRelationship
func NewBACnetConstructedDataDefaultSubordinateRelationship(defaultSubordinateRelationship BACnetRelationshipTagged, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataDefaultSubordinateRelationship {
	_result := &_BACnetConstructedDataDefaultSubordinateRelationship{
		DefaultSubordinateRelationship: defaultSubordinateRelationship,
		_BACnetConstructedData:         NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataDefaultSubordinateRelationship(structType any) BACnetConstructedDataDefaultSubordinateRelationship {
	if casted, ok := structType.(BACnetConstructedDataDefaultSubordinateRelationship); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataDefaultSubordinateRelationship); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) GetTypeName() string {
	return "BACnetConstructedDataDefaultSubordinateRelationship"
}

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (defaultSubordinateRelationship)
	lengthInBits += m.DefaultSubordinateRelationship.GetLengthInBits(ctx)

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetConstructedDataDefaultSubordinateRelationshipParse(ctx context.Context, theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataDefaultSubordinateRelationship, error) {
	return BACnetConstructedDataDefaultSubordinateRelationshipParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
}

func BACnetConstructedDataDefaultSubordinateRelationshipParseWithBufferProducer(tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConstructedDataDefaultSubordinateRelationship, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetConstructedDataDefaultSubordinateRelationship, error) {
		return BACnetConstructedDataDefaultSubordinateRelationshipParseWithBuffer(ctx, readBuffer, tagNumber, objectTypeArgument, propertyIdentifierArgument, arrayIndexArgument)
	}
}

func BACnetConstructedDataDefaultSubordinateRelationshipParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataDefaultSubordinateRelationship, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataDefaultSubordinateRelationship"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataDefaultSubordinateRelationship")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	defaultSubordinateRelationship, err := ReadSimpleField[BACnetRelationshipTagged](ctx, "defaultSubordinateRelationship", ReadComplex[BACnetRelationshipTagged](BACnetRelationshipTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'defaultSubordinateRelationship' field"))
	}

	actualValue, err := ReadVirtualField[BACnetRelationshipTagged](ctx, "actualValue", (*BACnetRelationshipTagged)(nil), defaultSubordinateRelationship)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'actualValue' field"))
	}
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataDefaultSubordinateRelationship"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataDefaultSubordinateRelationship")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataDefaultSubordinateRelationship{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		DefaultSubordinateRelationship: defaultSubordinateRelationship,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataDefaultSubordinateRelationship"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataDefaultSubordinateRelationship")
		}

		if err := WriteSimpleField[BACnetRelationshipTagged](ctx, "defaultSubordinateRelationship", m.GetDefaultSubordinateRelationship(), WriteComplex[BACnetRelationshipTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'defaultSubordinateRelationship' field")
		}
		// Virtual field
		actualValue := m.GetActualValue()
		_ = actualValue
		if _actualValueErr := writeBuffer.WriteVirtual(ctx, "actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataDefaultSubordinateRelationship"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataDefaultSubordinateRelationship")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) isBACnetConstructedDataDefaultSubordinateRelationship() bool {
	return true
}

func (m *_BACnetConstructedDataDefaultSubordinateRelationship) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
