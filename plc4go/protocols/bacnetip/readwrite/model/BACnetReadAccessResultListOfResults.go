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

// BACnetReadAccessResultListOfResults is the corresponding interface of BACnetReadAccessResultListOfResults
type BACnetReadAccessResultListOfResults interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetOpeningTag returns OpeningTag (property field)
	GetOpeningTag() BACnetOpeningTag
	// GetListOfReadAccessProperty returns ListOfReadAccessProperty (property field)
	GetListOfReadAccessProperty() []BACnetReadAccessProperty
	// GetClosingTag returns ClosingTag (property field)
	GetClosingTag() BACnetClosingTag
}

// BACnetReadAccessResultListOfResultsExactly can be used when we want exactly this type and not a type which fulfills BACnetReadAccessResultListOfResults.
// This is useful for switch cases.
type BACnetReadAccessResultListOfResultsExactly interface {
	BACnetReadAccessResultListOfResults
	isBACnetReadAccessResultListOfResults() bool
}

// _BACnetReadAccessResultListOfResults is the data-structure of this message
type _BACnetReadAccessResultListOfResults struct {
	OpeningTag               BACnetOpeningTag
	ListOfReadAccessProperty []BACnetReadAccessProperty
	ClosingTag               BACnetClosingTag

	// Arguments.
	TagNumber          uint8
	ObjectTypeArgument BACnetObjectType
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetReadAccessResultListOfResults) GetOpeningTag() BACnetOpeningTag {
	return m.OpeningTag
}

func (m *_BACnetReadAccessResultListOfResults) GetListOfReadAccessProperty() []BACnetReadAccessProperty {
	return m.ListOfReadAccessProperty
}

func (m *_BACnetReadAccessResultListOfResults) GetClosingTag() BACnetClosingTag {
	return m.ClosingTag
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetReadAccessResultListOfResults factory function for _BACnetReadAccessResultListOfResults
func NewBACnetReadAccessResultListOfResults(openingTag BACnetOpeningTag, listOfReadAccessProperty []BACnetReadAccessProperty, closingTag BACnetClosingTag, tagNumber uint8, objectTypeArgument BACnetObjectType) *_BACnetReadAccessResultListOfResults {
	return &_BACnetReadAccessResultListOfResults{OpeningTag: openingTag, ListOfReadAccessProperty: listOfReadAccessProperty, ClosingTag: closingTag, TagNumber: tagNumber, ObjectTypeArgument: objectTypeArgument}
}

// Deprecated: use the interface for direct cast
func CastBACnetReadAccessResultListOfResults(structType any) BACnetReadAccessResultListOfResults {
	if casted, ok := structType.(BACnetReadAccessResultListOfResults); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetReadAccessResultListOfResults); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetReadAccessResultListOfResults) GetTypeName() string {
	return "BACnetReadAccessResultListOfResults"
}

func (m *_BACnetReadAccessResultListOfResults) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (openingTag)
	lengthInBits += m.OpeningTag.GetLengthInBits(ctx)

	// Array field
	if len(m.ListOfReadAccessProperty) > 0 {
		for _, element := range m.ListOfReadAccessProperty {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	// Simple field (closingTag)
	lengthInBits += m.ClosingTag.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetReadAccessResultListOfResults) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetReadAccessResultListOfResultsParse(ctx context.Context, theBytes []byte, tagNumber uint8, objectTypeArgument BACnetObjectType) (BACnetReadAccessResultListOfResults, error) {
	return BACnetReadAccessResultListOfResultsParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), tagNumber, objectTypeArgument)
}

func BACnetReadAccessResultListOfResultsParseWithBufferProducer(tagNumber uint8, objectTypeArgument BACnetObjectType) func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetReadAccessResultListOfResults, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetReadAccessResultListOfResults, error) {
		return BACnetReadAccessResultListOfResultsParseWithBuffer(ctx, readBuffer, tagNumber, objectTypeArgument)
	}
}

func BACnetReadAccessResultListOfResultsParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType) (BACnetReadAccessResultListOfResults, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetReadAccessResultListOfResults"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetReadAccessResultListOfResults")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	openingTag, err := ReadSimpleField[BACnetOpeningTag](ctx, "openingTag", ReadComplex[BACnetOpeningTag](BACnetOpeningTagParseWithBufferProducer((uint8)(tagNumber)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'openingTag' field"))
	}

	listOfReadAccessProperty, err := ReadTerminatedArrayField[BACnetReadAccessProperty](ctx, "listOfReadAccessProperty", ReadComplex[BACnetReadAccessProperty](BACnetReadAccessPropertyParseWithBufferProducer((BACnetObjectType)(objectTypeArgument)), readBuffer), IsBACnetConstructedDataClosingTag(ctx, readBuffer, false, tagNumber))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'listOfReadAccessProperty' field"))
	}

	closingTag, err := ReadSimpleField[BACnetClosingTag](ctx, "closingTag", ReadComplex[BACnetClosingTag](BACnetClosingTagParseWithBufferProducer((uint8)(tagNumber)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'closingTag' field"))
	}

	if closeErr := readBuffer.CloseContext("BACnetReadAccessResultListOfResults"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetReadAccessResultListOfResults")
	}

	// Create the instance
	return &_BACnetReadAccessResultListOfResults{
		TagNumber:                tagNumber,
		ObjectTypeArgument:       objectTypeArgument,
		OpeningTag:               openingTag,
		ListOfReadAccessProperty: listOfReadAccessProperty,
		ClosingTag:               closingTag,
	}, nil
}

func (m *_BACnetReadAccessResultListOfResults) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetReadAccessResultListOfResults) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("BACnetReadAccessResultListOfResults"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetReadAccessResultListOfResults")
	}

	if err := WriteSimpleField[BACnetOpeningTag](ctx, "openingTag", m.GetOpeningTag(), WriteComplex[BACnetOpeningTag](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'openingTag' field")
	}

	if err := WriteComplexTypeArrayField(ctx, "listOfReadAccessProperty", m.GetListOfReadAccessProperty(), writeBuffer); err != nil {
		return errors.Wrap(err, "Error serializing 'listOfReadAccessProperty' field")
	}

	if err := WriteSimpleField[BACnetClosingTag](ctx, "closingTag", m.GetClosingTag(), WriteComplex[BACnetClosingTag](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'closingTag' field")
	}

	if popErr := writeBuffer.PopContext("BACnetReadAccessResultListOfResults"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetReadAccessResultListOfResults")
	}
	return nil
}

////
// Arguments Getter

func (m *_BACnetReadAccessResultListOfResults) GetTagNumber() uint8 {
	return m.TagNumber
}
func (m *_BACnetReadAccessResultListOfResults) GetObjectTypeArgument() BACnetObjectType {
	return m.ObjectTypeArgument
}

//
////

func (m *_BACnetReadAccessResultListOfResults) isBACnetReadAccessResultListOfResults() bool {
	return true
}

func (m *_BACnetReadAccessResultListOfResults) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
