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
type BACnetConfirmedServiceACKRemovedAuthenticate struct {
	*BACnetConfirmedServiceACK
}

// The corresponding interface
type IBACnetConfirmedServiceACKRemovedAuthenticate interface {
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
func (m *BACnetConfirmedServiceACKRemovedAuthenticate) GetServiceChoice() uint8 {
	return 0x18
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *BACnetConfirmedServiceACKRemovedAuthenticate) InitializeParent(parent *BACnetConfirmedServiceACK) {
}

func (m *BACnetConfirmedServiceACKRemovedAuthenticate) GetParent() *BACnetConfirmedServiceACK {
	return m.BACnetConfirmedServiceACK
}

// NewBACnetConfirmedServiceACKRemovedAuthenticate factory function for BACnetConfirmedServiceACKRemovedAuthenticate
func NewBACnetConfirmedServiceACKRemovedAuthenticate() *BACnetConfirmedServiceACKRemovedAuthenticate {
	_result := &BACnetConfirmedServiceACKRemovedAuthenticate{
		BACnetConfirmedServiceACK: NewBACnetConfirmedServiceACK(),
	}
	_result.Child = _result
	return _result
}

func CastBACnetConfirmedServiceACKRemovedAuthenticate(structType interface{}) *BACnetConfirmedServiceACKRemovedAuthenticate {
	if casted, ok := structType.(BACnetConfirmedServiceACKRemovedAuthenticate); ok {
		return &casted
	}
	if casted, ok := structType.(*BACnetConfirmedServiceACKRemovedAuthenticate); ok {
		return casted
	}
	if casted, ok := structType.(BACnetConfirmedServiceACK); ok {
		return CastBACnetConfirmedServiceACKRemovedAuthenticate(casted.Child)
	}
	if casted, ok := structType.(*BACnetConfirmedServiceACK); ok {
		return CastBACnetConfirmedServiceACKRemovedAuthenticate(casted.Child)
	}
	return nil
}

func (m *BACnetConfirmedServiceACKRemovedAuthenticate) GetTypeName() string {
	return "BACnetConfirmedServiceACKRemovedAuthenticate"
}

func (m *BACnetConfirmedServiceACKRemovedAuthenticate) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *BACnetConfirmedServiceACKRemovedAuthenticate) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	return lengthInBits
}

func (m *BACnetConfirmedServiceACKRemovedAuthenticate) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConfirmedServiceACKRemovedAuthenticateParse(readBuffer utils.ReadBuffer) (*BACnetConfirmedServiceACKRemovedAuthenticate, error) {
	if pullErr := readBuffer.PullContext("BACnetConfirmedServiceACKRemovedAuthenticate"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("BACnetConfirmedServiceACKRemovedAuthenticate"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &BACnetConfirmedServiceACKRemovedAuthenticate{
		BACnetConfirmedServiceACK: &BACnetConfirmedServiceACK{},
	}
	_child.BACnetConfirmedServiceACK.Child = _child
	return _child, nil
}

func (m *BACnetConfirmedServiceACKRemovedAuthenticate) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConfirmedServiceACKRemovedAuthenticate"); pushErr != nil {
			return pushErr
		}

		if popErr := writeBuffer.PopContext("BACnetConfirmedServiceACKRemovedAuthenticate"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *BACnetConfirmedServiceACKRemovedAuthenticate) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
