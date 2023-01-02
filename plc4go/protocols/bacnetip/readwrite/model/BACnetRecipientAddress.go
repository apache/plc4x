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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetRecipientAddress is the corresponding interface of BACnetRecipientAddress
type BACnetRecipientAddress interface {
	utils.LengthAware
	utils.Serializable
	BACnetRecipient
	// GetAddressValue returns AddressValue (property field)
	GetAddressValue() BACnetAddressEnclosed
}

// BACnetRecipientAddressExactly can be used when we want exactly this type and not a type which fulfills BACnetRecipientAddress.
// This is useful for switch cases.
type BACnetRecipientAddressExactly interface {
	BACnetRecipientAddress
	isBACnetRecipientAddress() bool
}

// _BACnetRecipientAddress is the data-structure of this message
type _BACnetRecipientAddress struct {
	*_BACnetRecipient
	AddressValue BACnetAddressEnclosed
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetRecipientAddress) InitializeParent(parent BACnetRecipient, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetRecipientAddress) GetParent() BACnetRecipient {
	return m._BACnetRecipient
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetRecipientAddress) GetAddressValue() BACnetAddressEnclosed {
	return m.AddressValue
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetRecipientAddress factory function for _BACnetRecipientAddress
func NewBACnetRecipientAddress(addressValue BACnetAddressEnclosed, peekedTagHeader BACnetTagHeader) *_BACnetRecipientAddress {
	_result := &_BACnetRecipientAddress{
		AddressValue:     addressValue,
		_BACnetRecipient: NewBACnetRecipient(peekedTagHeader),
	}
	_result._BACnetRecipient._BACnetRecipientChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetRecipientAddress(structType interface{}) BACnetRecipientAddress {
	if casted, ok := structType.(BACnetRecipientAddress); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetRecipientAddress); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetRecipientAddress) GetTypeName() string {
	return "BACnetRecipientAddress"
}

func (m *_BACnetRecipientAddress) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetRecipientAddress) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (addressValue)
	lengthInBits += m.AddressValue.GetLengthInBits()

	return lengthInBits
}

func (m *_BACnetRecipientAddress) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetRecipientAddressParse(theBytes []byte) (BACnetRecipientAddress, error) {
	return BACnetRecipientAddressParseWithBuffer(utils.NewReadBufferByteBased(theBytes))
}

func BACnetRecipientAddressParseWithBuffer(readBuffer utils.ReadBuffer) (BACnetRecipientAddress, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetRecipientAddress"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetRecipientAddress")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (addressValue)
	if pullErr := readBuffer.PullContext("addressValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for addressValue")
	}
	_addressValue, _addressValueErr := BACnetAddressEnclosedParseWithBuffer(readBuffer, uint8(uint8(1)))
	if _addressValueErr != nil {
		return nil, errors.Wrap(_addressValueErr, "Error parsing 'addressValue' field of BACnetRecipientAddress")
	}
	addressValue := _addressValue.(BACnetAddressEnclosed)
	if closeErr := readBuffer.CloseContext("addressValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for addressValue")
	}

	if closeErr := readBuffer.CloseContext("BACnetRecipientAddress"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetRecipientAddress")
	}

	// Create a partially initialized instance
	_child := &_BACnetRecipientAddress{
		_BACnetRecipient: &_BACnetRecipient{},
		AddressValue:     addressValue,
	}
	_child._BACnetRecipient._BACnetRecipientChildRequirements = _child
	return _child, nil
}

func (m *_BACnetRecipientAddress) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetRecipientAddress) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetRecipientAddress"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetRecipientAddress")
		}

		// Simple Field (addressValue)
		if pushErr := writeBuffer.PushContext("addressValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for addressValue")
		}
		_addressValueErr := writeBuffer.WriteSerializable(m.GetAddressValue())
		if popErr := writeBuffer.PopContext("addressValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for addressValue")
		}
		if _addressValueErr != nil {
			return errors.Wrap(_addressValueErr, "Error serializing 'addressValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetRecipientAddress"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetRecipientAddress")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetRecipientAddress) isBACnetRecipientAddress() bool {
	return true
}

func (m *_BACnetRecipientAddress) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
