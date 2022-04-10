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
type BVLCWriteBroadcastDistributionTableEntry struct {
	Ip                       []uint8
	Port                     uint16
	BroadcastDistributionMap []uint8
}

// The corresponding interface
type IBVLCWriteBroadcastDistributionTableEntry interface {
	// GetIp returns Ip (property field)
	GetIp() []uint8
	// GetPort returns Port (property field)
	GetPort() uint16
	// GetBroadcastDistributionMap returns BroadcastDistributionMap (property field)
	GetBroadcastDistributionMap() []uint8
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
func (m *BVLCWriteBroadcastDistributionTableEntry) GetIp() []uint8 {
	return m.Ip
}

func (m *BVLCWriteBroadcastDistributionTableEntry) GetPort() uint16 {
	return m.Port
}

func (m *BVLCWriteBroadcastDistributionTableEntry) GetBroadcastDistributionMap() []uint8 {
	return m.BroadcastDistributionMap
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBVLCWriteBroadcastDistributionTableEntry factory function for BVLCWriteBroadcastDistributionTableEntry
func NewBVLCWriteBroadcastDistributionTableEntry(ip []uint8, port uint16, broadcastDistributionMap []uint8) *BVLCWriteBroadcastDistributionTableEntry {
	return &BVLCWriteBroadcastDistributionTableEntry{Ip: ip, Port: port, BroadcastDistributionMap: broadcastDistributionMap}
}

func CastBVLCWriteBroadcastDistributionTableEntry(structType interface{}) *BVLCWriteBroadcastDistributionTableEntry {
	if casted, ok := structType.(BVLCWriteBroadcastDistributionTableEntry); ok {
		return &casted
	}
	if casted, ok := structType.(*BVLCWriteBroadcastDistributionTableEntry); ok {
		return casted
	}
	return nil
}

func (m *BVLCWriteBroadcastDistributionTableEntry) GetTypeName() string {
	return "BVLCWriteBroadcastDistributionTableEntry"
}

func (m *BVLCWriteBroadcastDistributionTableEntry) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *BVLCWriteBroadcastDistributionTableEntry) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Array field
	if len(m.Ip) > 0 {
		lengthInBits += 8 * uint16(len(m.Ip))
	}

	// Simple field (port)
	lengthInBits += 16

	// Array field
	if len(m.BroadcastDistributionMap) > 0 {
		lengthInBits += 8 * uint16(len(m.BroadcastDistributionMap))
	}

	return lengthInBits
}

func (m *BVLCWriteBroadcastDistributionTableEntry) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BVLCWriteBroadcastDistributionTableEntryParse(readBuffer utils.ReadBuffer) (*BVLCWriteBroadcastDistributionTableEntry, error) {
	if pullErr := readBuffer.PullContext("BVLCWriteBroadcastDistributionTableEntry"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	// Array field (ip)
	if pullErr := readBuffer.PullContext("ip", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, pullErr
	}
	// Count array
	ip := make([]uint8, uint16(4))
	{
		for curItem := uint16(0); curItem < uint16(uint16(4)); curItem++ {
			_item, _err := readBuffer.ReadUint8("", 8)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'ip' field")
			}
			ip[curItem] = _item
		}
	}
	if closeErr := readBuffer.CloseContext("ip", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, closeErr
	}

	// Simple Field (port)
	_port, _portErr := readBuffer.ReadUint16("port", 16)
	if _portErr != nil {
		return nil, errors.Wrap(_portErr, "Error parsing 'port' field")
	}
	port := _port

	// Array field (broadcastDistributionMap)
	if pullErr := readBuffer.PullContext("broadcastDistributionMap", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, pullErr
	}
	// Count array
	broadcastDistributionMap := make([]uint8, uint16(4))
	{
		for curItem := uint16(0); curItem < uint16(uint16(4)); curItem++ {
			_item, _err := readBuffer.ReadUint8("", 8)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'broadcastDistributionMap' field")
			}
			broadcastDistributionMap[curItem] = _item
		}
	}
	if closeErr := readBuffer.CloseContext("broadcastDistributionMap", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, closeErr
	}

	if closeErr := readBuffer.CloseContext("BVLCWriteBroadcastDistributionTableEntry"); closeErr != nil {
		return nil, closeErr
	}

	// Create the instance
	return NewBVLCWriteBroadcastDistributionTableEntry(ip, port, broadcastDistributionMap), nil
}

func (m *BVLCWriteBroadcastDistributionTableEntry) Serialize(writeBuffer utils.WriteBuffer) error {
	if pushErr := writeBuffer.PushContext("BVLCWriteBroadcastDistributionTableEntry"); pushErr != nil {
		return pushErr
	}

	// Array Field (ip)
	if m.Ip != nil {
		if pushErr := writeBuffer.PushContext("ip", utils.WithRenderAsList(true)); pushErr != nil {
			return pushErr
		}
		for _, _element := range m.Ip {
			_elementErr := writeBuffer.WriteUint8("", 8, _element)
			if _elementErr != nil {
				return errors.Wrap(_elementErr, "Error serializing 'ip' field")
			}
		}
		if popErr := writeBuffer.PopContext("ip", utils.WithRenderAsList(true)); popErr != nil {
			return popErr
		}
	}

	// Simple Field (port)
	port := uint16(m.Port)
	_portErr := writeBuffer.WriteUint16("port", 16, (port))
	if _portErr != nil {
		return errors.Wrap(_portErr, "Error serializing 'port' field")
	}

	// Array Field (broadcastDistributionMap)
	if m.BroadcastDistributionMap != nil {
		if pushErr := writeBuffer.PushContext("broadcastDistributionMap", utils.WithRenderAsList(true)); pushErr != nil {
			return pushErr
		}
		for _, _element := range m.BroadcastDistributionMap {
			_elementErr := writeBuffer.WriteUint8("", 8, _element)
			if _elementErr != nil {
				return errors.Wrap(_elementErr, "Error serializing 'broadcastDistributionMap' field")
			}
		}
		if popErr := writeBuffer.PopContext("broadcastDistributionMap", utils.WithRenderAsList(true)); popErr != nil {
			return popErr
		}
	}

	if popErr := writeBuffer.PopContext("BVLCWriteBroadcastDistributionTableEntry"); popErr != nil {
		return popErr
	}
	return nil
}

func (m *BVLCWriteBroadcastDistributionTableEntry) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
