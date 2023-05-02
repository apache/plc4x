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

// BACnetActionList is the corresponding interface of BACnetActionList
type BACnetActionList interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetInnerOpeningTag returns InnerOpeningTag (property field)
	GetInnerOpeningTag() BACnetOpeningTag
	// GetAction returns Action (property field)
	GetAction() []BACnetActionCommand
	// GetInnerClosingTag returns InnerClosingTag (property field)
	GetInnerClosingTag() BACnetClosingTag
}

// BACnetActionListExactly can be used when we want exactly this type and not a type which fulfills BACnetActionList.
// This is useful for switch cases.
type BACnetActionListExactly interface {
	BACnetActionList
	isBACnetActionList() bool
}

// _BACnetActionList is the data-structure of this message
type _BACnetActionList struct {
	InnerOpeningTag BACnetOpeningTag
	Action          []BACnetActionCommand
	InnerClosingTag BACnetClosingTag
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetActionList) GetInnerOpeningTag() BACnetOpeningTag {
	return m.InnerOpeningTag
}

func (m *_BACnetActionList) GetAction() []BACnetActionCommand {
	return m.Action
}

func (m *_BACnetActionList) GetInnerClosingTag() BACnetClosingTag {
	return m.InnerClosingTag
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetActionList factory function for _BACnetActionList
func NewBACnetActionList(innerOpeningTag BACnetOpeningTag, action []BACnetActionCommand, innerClosingTag BACnetClosingTag) *_BACnetActionList {
	return &_BACnetActionList{InnerOpeningTag: innerOpeningTag, Action: action, InnerClosingTag: innerClosingTag}
}

// Deprecated: use the interface for direct cast
func CastBACnetActionList(structType any) BACnetActionList {
	if casted, ok := structType.(BACnetActionList); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetActionList); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetActionList) GetTypeName() string {
	return "BACnetActionList"
}

func (m *_BACnetActionList) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (innerOpeningTag)
	lengthInBits += m.InnerOpeningTag.GetLengthInBits(ctx)

	// Array field
	if len(m.Action) > 0 {
		for _, element := range m.Action {
			lengthInBits += element.GetLengthInBits(ctx)
		}
	}

	// Simple field (innerClosingTag)
	lengthInBits += m.InnerClosingTag.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetActionList) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetActionListParse(theBytes []byte) (BACnetActionList, error) {
	return BACnetActionListParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes))
}

func BACnetActionListParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetActionList, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetActionList"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetActionList")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (innerOpeningTag)
	if pullErr := readBuffer.PullContext("innerOpeningTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for innerOpeningTag")
	}
	_innerOpeningTag, _innerOpeningTagErr := BACnetOpeningTagParseWithBuffer(ctx, readBuffer, uint8(uint8(0)))
	if _innerOpeningTagErr != nil {
		return nil, errors.Wrap(_innerOpeningTagErr, "Error parsing 'innerOpeningTag' field of BACnetActionList")
	}
	innerOpeningTag := _innerOpeningTag.(BACnetOpeningTag)
	if closeErr := readBuffer.CloseContext("innerOpeningTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for innerOpeningTag")
	}

	// Array field (action)
	if pullErr := readBuffer.PullContext("action", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for action")
	}
	// Terminated array
	var action []BACnetActionCommand
	{
		for !bool(IsBACnetConstructedDataClosingTag(readBuffer, false, 0)) {
			_item, _err := BACnetActionCommandParseWithBuffer(ctx, readBuffer)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'action' field of BACnetActionList")
			}
			action = append(action, _item.(BACnetActionCommand))
		}
	}
	if closeErr := readBuffer.CloseContext("action", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for action")
	}

	// Simple Field (innerClosingTag)
	if pullErr := readBuffer.PullContext("innerClosingTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for innerClosingTag")
	}
	_innerClosingTag, _innerClosingTagErr := BACnetClosingTagParseWithBuffer(ctx, readBuffer, uint8(uint8(0)))
	if _innerClosingTagErr != nil {
		return nil, errors.Wrap(_innerClosingTagErr, "Error parsing 'innerClosingTag' field of BACnetActionList")
	}
	innerClosingTag := _innerClosingTag.(BACnetClosingTag)
	if closeErr := readBuffer.CloseContext("innerClosingTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for innerClosingTag")
	}

	if closeErr := readBuffer.CloseContext("BACnetActionList"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetActionList")
	}

	// Create the instance
	return &_BACnetActionList{
		InnerOpeningTag: innerOpeningTag,
		Action:          action,
		InnerClosingTag: innerClosingTag,
	}, nil
}

func (m *_BACnetActionList) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetActionList) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetActionList"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetActionList")
	}

	// Simple Field (innerOpeningTag)
	if pushErr := writeBuffer.PushContext("innerOpeningTag"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for innerOpeningTag")
	}
	_innerOpeningTagErr := writeBuffer.WriteSerializable(ctx, m.GetInnerOpeningTag())
	if popErr := writeBuffer.PopContext("innerOpeningTag"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for innerOpeningTag")
	}
	if _innerOpeningTagErr != nil {
		return errors.Wrap(_innerOpeningTagErr, "Error serializing 'innerOpeningTag' field")
	}

	// Array Field (action)
	if pushErr := writeBuffer.PushContext("action", utils.WithRenderAsList(true)); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for action")
	}
	for _curItem, _element := range m.GetAction() {
		_ = _curItem
		arrayCtx := utils.CreateArrayContext(ctx, len(m.GetAction()), _curItem)
		_ = arrayCtx
		_elementErr := writeBuffer.WriteSerializable(arrayCtx, _element)
		if _elementErr != nil {
			return errors.Wrap(_elementErr, "Error serializing 'action' field")
		}
	}
	if popErr := writeBuffer.PopContext("action", utils.WithRenderAsList(true)); popErr != nil {
		return errors.Wrap(popErr, "Error popping for action")
	}

	// Simple Field (innerClosingTag)
	if pushErr := writeBuffer.PushContext("innerClosingTag"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for innerClosingTag")
	}
	_innerClosingTagErr := writeBuffer.WriteSerializable(ctx, m.GetInnerClosingTag())
	if popErr := writeBuffer.PopContext("innerClosingTag"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for innerClosingTag")
	}
	if _innerClosingTagErr != nil {
		return errors.Wrap(_innerClosingTagErr, "Error serializing 'innerClosingTag' field")
	}

	if popErr := writeBuffer.PopContext("BACnetActionList"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetActionList")
	}
	return nil
}

func (m *_BACnetActionList) isBACnetActionList() bool {
	return true
}

func (m *_BACnetActionList) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
