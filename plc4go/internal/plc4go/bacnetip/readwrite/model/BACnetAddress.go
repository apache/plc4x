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
)

// Code generated by code-generation. DO NOT EDIT.

// The data-structure of this message
type BACnetAddress struct {
	Address []uint8
	Port    uint16
}

// The corresponding interface
type IBACnetAddress interface {
	// GetAddress returns Address (property field)
	GetAddress() []uint8
	// GetPort returns Port (property field)
	GetPort() uint16
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////
func (m *BACnetAddress) GetAddress() []uint8 {
	return m.Address
}

func (m *BACnetAddress) GetPort() uint16 {
	return m.Port
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetAddress factory function for BACnetAddress
func NewBACnetAddress(address []uint8, port uint16) *BACnetAddress {
	return &BACnetAddress{Address: address, Port: port}
}

func CastBACnetAddress(structType interface{}) *BACnetAddress {
	if casted, ok := structType.(BACnetAddress); ok {
		return &casted
	}
	if casted, ok := structType.(*BACnetAddress); ok {
		return casted
	}
	return nil
}

func (m *BACnetAddress) GetTypeName() string {
	return "BACnetAddress"
}

func (m *BACnetAddress) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *BACnetAddress) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Array field
	if len(m.Address) > 0 {
		lengthInBits += 8 * uint16(len(m.Address))
	}

	// Simple field (port)
	lengthInBits += 16

	return lengthInBits
}

func (m *BACnetAddress) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetAddressParse(readBuffer utils.ReadBuffer) (*BACnetAddress, error) {
	if pullErr := readBuffer.PullContext("BACnetAddress"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	// Array field (address)
	if pullErr := readBuffer.PullContext("address", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, pullErr
	}
	// Count array
	address := make([]uint8, uint16(4))
	{
		for curItem := uint16(0); curItem < uint16(uint16(4)); curItem++ {
			_item, _err := readBuffer.ReadUint8("", 8)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'address' field")
			}
			address[curItem] = _item
		}
	}
	if closeErr := readBuffer.CloseContext("address", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, closeErr
	}

	// Simple Field (port)
	_port, _portErr := readBuffer.ReadUint16("port", 16)
	if _portErr != nil {
		return nil, errors.Wrap(_portErr, "Error parsing 'port' field")
	}
	port := _port

	if closeErr := readBuffer.CloseContext("BACnetAddress"); closeErr != nil {
		return nil, closeErr
	}

	// Create the instance
	return NewBACnetAddress(address, port), nil
}

func (m *BACnetAddress) Serialize(writeBuffer utils.WriteBuffer) error {
	if pushErr := writeBuffer.PushContext("BACnetAddress"); pushErr != nil {
		return pushErr
	}

	// Array Field (address)
	if m.Address != nil {
		if pushErr := writeBuffer.PushContext("address", utils.WithRenderAsList(true)); pushErr != nil {
			return pushErr
		}
		for _, _element := range m.Address {
			_elementErr := writeBuffer.WriteUint8("", 8, _element)
			if _elementErr != nil {
				return errors.Wrap(_elementErr, "Error serializing 'address' field")
			}
		}
		if popErr := writeBuffer.PopContext("address", utils.WithRenderAsList(true)); popErr != nil {
			return popErr
		}
	}

	// Simple Field (port)
	port := uint16(m.Port)
	_portErr := writeBuffer.WriteUint16("port", 16, (port))
	if _portErr != nil {
		return errors.Wrap(_portErr, "Error serializing 'port' field")
	}

	if popErr := writeBuffer.PopContext("BACnetAddress"); popErr != nil {
		return popErr
	}
	return nil
}

func (m *BACnetAddress) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
