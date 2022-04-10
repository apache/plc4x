/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
	"io"
)

// Code generated by code-generation. DO NOT EDIT.

// The data-structure of this message
type BACnetPropertyStatesAction struct {
	*BACnetPropertyStates
	Action *BACnetAction

	// Arguments.
	TagNumber uint8
}

// The corresponding interface
type IBACnetPropertyStatesAction interface {
	IBACnetPropertyStates
	// GetAction returns Action (property field)
	GetAction() *BACnetAction
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////
///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *BACnetPropertyStatesAction) InitializeParent(parent *BACnetPropertyStates, openingTag *BACnetOpeningTag, peekedTagHeader *BACnetTagHeader, closingTag *BACnetClosingTag) {
	m.BACnetPropertyStates.OpeningTag = openingTag
	m.BACnetPropertyStates.PeekedTagHeader = peekedTagHeader
	m.BACnetPropertyStates.ClosingTag = closingTag
}

func (m *BACnetPropertyStatesAction) GetParent() *BACnetPropertyStates {
	return m.BACnetPropertyStates
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////
func (m *BACnetPropertyStatesAction) GetAction() *BACnetAction {
	return m.Action
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetPropertyStatesAction factory function for BACnetPropertyStatesAction
func NewBACnetPropertyStatesAction(action *BACnetAction, openingTag *BACnetOpeningTag, peekedTagHeader *BACnetTagHeader, closingTag *BACnetClosingTag, tagNumber uint8) *BACnetPropertyStatesAction {
	_result := &BACnetPropertyStatesAction{
		Action:               action,
		BACnetPropertyStates: NewBACnetPropertyStates(openingTag, peekedTagHeader, closingTag, tagNumber),
	}
	_result.Child = _result
	return _result
}

func CastBACnetPropertyStatesAction(structType interface{}) *BACnetPropertyStatesAction {
	if casted, ok := structType.(BACnetPropertyStatesAction); ok {
		return &casted
	}
	if casted, ok := structType.(*BACnetPropertyStatesAction); ok {
		return casted
	}
	if casted, ok := structType.(BACnetPropertyStates); ok {
		return CastBACnetPropertyStatesAction(casted.Child)
	}
	if casted, ok := structType.(*BACnetPropertyStates); ok {
		return CastBACnetPropertyStatesAction(casted.Child)
	}
	return nil
}

func (m *BACnetPropertyStatesAction) GetTypeName() string {
	return "BACnetPropertyStatesAction"
}

func (m *BACnetPropertyStatesAction) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *BACnetPropertyStatesAction) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Optional Field (action)
	if m.Action != nil {
		lengthInBits += (*m.Action).GetLengthInBits()
	}

	return lengthInBits
}

func (m *BACnetPropertyStatesAction) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetPropertyStatesActionParse(readBuffer utils.ReadBuffer, tagNumber uint8, peekedTagNumber uint8) (*BACnetPropertyStatesAction, error) {
	if pullErr := readBuffer.PullContext("BACnetPropertyStatesAction"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	// Optional Field (action) (Can be skipped, if a given expression evaluates to false)
	var action *BACnetAction = nil
	{
		currentPos = readBuffer.GetPos()
		if pullErr := readBuffer.PullContext("action"); pullErr != nil {
			return nil, pullErr
		}
		_val, _err := BACnetActionParse(readBuffer, peekedTagNumber)
		switch {
		case errors.Is(_err, utils.ParseAssertError{}) || errors.Is(_err, io.EOF):
			readBuffer.Reset(currentPos)
		case _err != nil:
			return nil, errors.Wrap(_err, "Error parsing 'action' field")
		default:
			action = CastBACnetAction(_val)
			if closeErr := readBuffer.CloseContext("action"); closeErr != nil {
				return nil, closeErr
			}
		}
	}

	if closeErr := readBuffer.CloseContext("BACnetPropertyStatesAction"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &BACnetPropertyStatesAction{
		Action:               CastBACnetAction(action),
		BACnetPropertyStates: &BACnetPropertyStates{},
	}
	_child.BACnetPropertyStates.Child = _child
	return _child, nil
}

func (m *BACnetPropertyStatesAction) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetPropertyStatesAction"); pushErr != nil {
			return pushErr
		}

		// Optional Field (action) (Can be skipped, if the value is null)
		var action *BACnetAction = nil
		if m.Action != nil {
			if pushErr := writeBuffer.PushContext("action"); pushErr != nil {
				return pushErr
			}
			action = m.Action
			_actionErr := action.Serialize(writeBuffer)
			if popErr := writeBuffer.PopContext("action"); popErr != nil {
				return popErr
			}
			if _actionErr != nil {
				return errors.Wrap(_actionErr, "Error serializing 'action' field")
			}
		}

		if popErr := writeBuffer.PopContext("BACnetPropertyStatesAction"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *BACnetPropertyStatesAction) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
