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
)

// Code generated by code-generation. DO NOT EDIT.

// The data-structure of this message
type BACnetConfirmedServiceACKVTOpen struct {
	*BACnetConfirmedServiceACK
}

// The corresponding interface
type IBACnetConfirmedServiceACKVTOpen interface {
	IBACnetConfirmedServiceACK
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
func (m *BACnetConfirmedServiceACKVTOpen) GetServiceChoice() uint8 {
	return 0x15
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *BACnetConfirmedServiceACKVTOpen) InitializeParent(parent *BACnetConfirmedServiceACK) {}

func (m *BACnetConfirmedServiceACKVTOpen) GetParent() *BACnetConfirmedServiceACK {
	return m.BACnetConfirmedServiceACK
}

// NewBACnetConfirmedServiceACKVTOpen factory function for BACnetConfirmedServiceACKVTOpen
func NewBACnetConfirmedServiceACKVTOpen() *BACnetConfirmedServiceACKVTOpen {
	_result := &BACnetConfirmedServiceACKVTOpen{
		BACnetConfirmedServiceACK: NewBACnetConfirmedServiceACK(),
	}
	_result.Child = _result
	return _result
}

func CastBACnetConfirmedServiceACKVTOpen(structType interface{}) *BACnetConfirmedServiceACKVTOpen {
	if casted, ok := structType.(BACnetConfirmedServiceACKVTOpen); ok {
		return &casted
	}
	if casted, ok := structType.(*BACnetConfirmedServiceACKVTOpen); ok {
		return casted
	}
	if casted, ok := structType.(BACnetConfirmedServiceACK); ok {
		return CastBACnetConfirmedServiceACKVTOpen(casted.Child)
	}
	if casted, ok := structType.(*BACnetConfirmedServiceACK); ok {
		return CastBACnetConfirmedServiceACKVTOpen(casted.Child)
	}
	return nil
}

func (m *BACnetConfirmedServiceACKVTOpen) GetTypeName() string {
	return "BACnetConfirmedServiceACKVTOpen"
}

func (m *BACnetConfirmedServiceACKVTOpen) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *BACnetConfirmedServiceACKVTOpen) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	return lengthInBits
}

func (m *BACnetConfirmedServiceACKVTOpen) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConfirmedServiceACKVTOpenParse(readBuffer utils.ReadBuffer) (*BACnetConfirmedServiceACKVTOpen, error) {
	if pullErr := readBuffer.PullContext("BACnetConfirmedServiceACKVTOpen"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("BACnetConfirmedServiceACKVTOpen"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &BACnetConfirmedServiceACKVTOpen{
		BACnetConfirmedServiceACK: &BACnetConfirmedServiceACK{},
	}
	_child.BACnetConfirmedServiceACK.Child = _child
	return _child, nil
}

func (m *BACnetConfirmedServiceACKVTOpen) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConfirmedServiceACKVTOpen"); pushErr != nil {
			return pushErr
		}

		if popErr := writeBuffer.PopContext("BACnetConfirmedServiceACKVTOpen"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *BACnetConfirmedServiceACKVTOpen) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
