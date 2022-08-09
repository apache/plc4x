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

// BACnetHostAddressNull is the corresponding interface of BACnetHostAddressNull
type BACnetHostAddressNull interface {
	utils.LengthAware
	utils.Serializable
	BACnetHostAddress
	// GetNone returns None (property field)
	GetNone() BACnetContextTagNull
}

// BACnetHostAddressNullExactly can be used when we want exactly this type and not a type which fulfills BACnetHostAddressNull.
// This is useful for switch cases.
type BACnetHostAddressNullExactly interface {
	BACnetHostAddressNull
	isBACnetHostAddressNull() bool
}

// _BACnetHostAddressNull is the data-structure of this message
type _BACnetHostAddressNull struct {
	*_BACnetHostAddress
	None BACnetContextTagNull
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetHostAddressNull) InitializeParent(parent BACnetHostAddress, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetHostAddressNull) GetParent() BACnetHostAddress {
	return m._BACnetHostAddress
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetHostAddressNull) GetNone() BACnetContextTagNull {
	return m.None
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetHostAddressNull factory function for _BACnetHostAddressNull
func NewBACnetHostAddressNull(none BACnetContextTagNull, peekedTagHeader BACnetTagHeader) *_BACnetHostAddressNull {
	_result := &_BACnetHostAddressNull{
		None:               none,
		_BACnetHostAddress: NewBACnetHostAddress(peekedTagHeader),
	}
	_result._BACnetHostAddress._BACnetHostAddressChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetHostAddressNull(structType interface{}) BACnetHostAddressNull {
	if casted, ok := structType.(BACnetHostAddressNull); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetHostAddressNull); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetHostAddressNull) GetTypeName() string {
	return "BACnetHostAddressNull"
}

func (m *_BACnetHostAddressNull) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetHostAddressNull) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (none)
	lengthInBits += m.None.GetLengthInBits()

	return lengthInBits
}

func (m *_BACnetHostAddressNull) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetHostAddressNullParse(readBuffer utils.ReadBuffer) (BACnetHostAddressNull, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetHostAddressNull"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetHostAddressNull")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (none)
	if pullErr := readBuffer.PullContext("none"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for none")
	}
	_none, _noneErr := BACnetContextTagParse(readBuffer, uint8(uint8(0)), BACnetDataType(BACnetDataType_NULL))
	if _noneErr != nil {
		return nil, errors.Wrap(_noneErr, "Error parsing 'none' field of BACnetHostAddressNull")
	}
	none := _none.(BACnetContextTagNull)
	if closeErr := readBuffer.CloseContext("none"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for none")
	}

	if closeErr := readBuffer.CloseContext("BACnetHostAddressNull"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetHostAddressNull")
	}

	// Create a partially initialized instance
	_child := &_BACnetHostAddressNull{
		_BACnetHostAddress: &_BACnetHostAddress{},
		None:               none,
	}
	_child._BACnetHostAddress._BACnetHostAddressChildRequirements = _child
	return _child, nil
}

func (m *_BACnetHostAddressNull) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetHostAddressNull"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetHostAddressNull")
		}

		// Simple Field (none)
		if pushErr := writeBuffer.PushContext("none"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for none")
		}
		_noneErr := writeBuffer.WriteSerializable(m.GetNone())
		if popErr := writeBuffer.PopContext("none"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for none")
		}
		if _noneErr != nil {
			return errors.Wrap(_noneErr, "Error serializing 'none' field")
		}

		if popErr := writeBuffer.PopContext("BACnetHostAddressNull"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetHostAddressNull")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetHostAddressNull) isBACnetHostAddressNull() bool {
	return true
}

func (m *_BACnetHostAddressNull) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
