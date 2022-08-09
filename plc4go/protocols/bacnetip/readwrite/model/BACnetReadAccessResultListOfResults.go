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
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetReadAccessResultListOfResults is the corresponding interface of BACnetReadAccessResultListOfResults
type BACnetReadAccessResultListOfResults interface {
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
func CastBACnetReadAccessResultListOfResults(structType interface{}) BACnetReadAccessResultListOfResults {
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

func (m *_BACnetReadAccessResultListOfResults) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetReadAccessResultListOfResults) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Simple field (openingTag)
	lengthInBits += m.OpeningTag.GetLengthInBits()

	// Array field
	if len(m.ListOfReadAccessProperty) > 0 {
		for _, element := range m.ListOfReadAccessProperty {
			lengthInBits += element.GetLengthInBits()
		}
	}

	// Simple field (closingTag)
	lengthInBits += m.ClosingTag.GetLengthInBits()

	return lengthInBits
}

func (m *_BACnetReadAccessResultListOfResults) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetReadAccessResultListOfResultsParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType) (BACnetReadAccessResultListOfResults, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetReadAccessResultListOfResults"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetReadAccessResultListOfResults")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (openingTag)
	if pullErr := readBuffer.PullContext("openingTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for openingTag")
	}
	_openingTag, _openingTagErr := BACnetOpeningTagParse(readBuffer, uint8(tagNumber))
	if _openingTagErr != nil {
		return nil, errors.Wrap(_openingTagErr, "Error parsing 'openingTag' field of BACnetReadAccessResultListOfResults")
	}
	openingTag := _openingTag.(BACnetOpeningTag)
	if closeErr := readBuffer.CloseContext("openingTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for openingTag")
	}

	// Array field (listOfReadAccessProperty)
	if pullErr := readBuffer.PullContext("listOfReadAccessProperty", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for listOfReadAccessProperty")
	}
	// Terminated array
	var listOfReadAccessProperty []BACnetReadAccessProperty
	{
		for !bool(IsBACnetConstructedDataClosingTag(readBuffer, false, tagNumber)) {
			_item, _err := BACnetReadAccessPropertyParse(readBuffer, objectTypeArgument)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'listOfReadAccessProperty' field of BACnetReadAccessResultListOfResults")
			}
			listOfReadAccessProperty = append(listOfReadAccessProperty, _item.(BACnetReadAccessProperty))

		}
	}
	if closeErr := readBuffer.CloseContext("listOfReadAccessProperty", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for listOfReadAccessProperty")
	}

	// Simple Field (closingTag)
	if pullErr := readBuffer.PullContext("closingTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for closingTag")
	}
	_closingTag, _closingTagErr := BACnetClosingTagParse(readBuffer, uint8(tagNumber))
	if _closingTagErr != nil {
		return nil, errors.Wrap(_closingTagErr, "Error parsing 'closingTag' field of BACnetReadAccessResultListOfResults")
	}
	closingTag := _closingTag.(BACnetClosingTag)
	if closeErr := readBuffer.CloseContext("closingTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for closingTag")
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

func (m *_BACnetReadAccessResultListOfResults) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetReadAccessResultListOfResults"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetReadAccessResultListOfResults")
	}

	// Simple Field (openingTag)
	if pushErr := writeBuffer.PushContext("openingTag"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for openingTag")
	}
	_openingTagErr := writeBuffer.WriteSerializable(m.GetOpeningTag())
	if popErr := writeBuffer.PopContext("openingTag"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for openingTag")
	}
	if _openingTagErr != nil {
		return errors.Wrap(_openingTagErr, "Error serializing 'openingTag' field")
	}

	// Array Field (listOfReadAccessProperty)
	if pushErr := writeBuffer.PushContext("listOfReadAccessProperty", utils.WithRenderAsList(true)); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for listOfReadAccessProperty")
	}
	for _, _element := range m.GetListOfReadAccessProperty() {
		_elementErr := writeBuffer.WriteSerializable(_element)
		if _elementErr != nil {
			return errors.Wrap(_elementErr, "Error serializing 'listOfReadAccessProperty' field")
		}
	}
	if popErr := writeBuffer.PopContext("listOfReadAccessProperty", utils.WithRenderAsList(true)); popErr != nil {
		return errors.Wrap(popErr, "Error popping for listOfReadAccessProperty")
	}

	// Simple Field (closingTag)
	if pushErr := writeBuffer.PushContext("closingTag"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for closingTag")
	}
	_closingTagErr := writeBuffer.WriteSerializable(m.GetClosingTag())
	if popErr := writeBuffer.PopContext("closingTag"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for closingTag")
	}
	if _closingTagErr != nil {
		return errors.Wrap(_closingTagErr, "Error serializing 'closingTag' field")
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
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
