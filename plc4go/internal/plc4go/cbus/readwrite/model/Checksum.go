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
type Checksum struct {
	Crc byte
}

// The corresponding interface
type IChecksum interface {
	// GetCrc returns Crc (property field)
	GetCrc() byte
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
func (m *Checksum) GetCrc() byte {
	return m.Crc
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewChecksum factory function for Checksum
func NewChecksum(crc byte) *Checksum {
	return &Checksum{Crc: crc}
}

func CastChecksum(structType interface{}) *Checksum {
	if casted, ok := structType.(Checksum); ok {
		return &casted
	}
	if casted, ok := structType.(*Checksum); ok {
		return casted
	}
	return nil
}

func (m *Checksum) GetTypeName() string {
	return "Checksum"
}

func (m *Checksum) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *Checksum) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Simple field (crc)
	lengthInBits += 8

	return lengthInBits
}

func (m *Checksum) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func ChecksumParse(readBuffer utils.ReadBuffer) (*Checksum, error) {
	if pullErr := readBuffer.PullContext("Checksum"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	// Simple Field (crc)
	_crc, _crcErr := readBuffer.ReadByte("crc")
	if _crcErr != nil {
		return nil, errors.Wrap(_crcErr, "Error parsing 'crc' field")
	}
	crc := _crc

	if closeErr := readBuffer.CloseContext("Checksum"); closeErr != nil {
		return nil, closeErr
	}

	// Create the instance
	return NewChecksum(crc), nil
}

func (m *Checksum) Serialize(writeBuffer utils.WriteBuffer) error {
	if pushErr := writeBuffer.PushContext("Checksum"); pushErr != nil {
		return pushErr
	}

	// Simple Field (crc)
	crc := byte(m.Crc)
	_crcErr := writeBuffer.WriteByte("crc", (crc))
	if _crcErr != nil {
		return errors.Wrap(_crcErr, "Error serializing 'crc' field")
	}

	if popErr := writeBuffer.PopContext("Checksum"); popErr != nil {
		return popErr
	}
	return nil
}

func (m *Checksum) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
