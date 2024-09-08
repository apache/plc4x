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
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/spi/codegen"
	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// BVLCDeleteForeignDeviceTableEntry is the corresponding interface of BVLCDeleteForeignDeviceTableEntry
type BVLCDeleteForeignDeviceTableEntry interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BVLC
	// GetIp returns Ip (property field)
	GetIp() []uint8
	// GetPort returns Port (property field)
	GetPort() uint16
	// IsBVLCDeleteForeignDeviceTableEntry is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBVLCDeleteForeignDeviceTableEntry()
}

// _BVLCDeleteForeignDeviceTableEntry is the data-structure of this message
type _BVLCDeleteForeignDeviceTableEntry struct {
	BVLCContract
	Ip   []uint8
	Port uint16
}

var _ BVLCDeleteForeignDeviceTableEntry = (*_BVLCDeleteForeignDeviceTableEntry)(nil)
var _ BVLCRequirements = (*_BVLCDeleteForeignDeviceTableEntry)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BVLCDeleteForeignDeviceTableEntry) GetBvlcFunction() uint8 {
	return 0x08
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BVLCDeleteForeignDeviceTableEntry) GetParent() BVLCContract {
	return m.BVLCContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BVLCDeleteForeignDeviceTableEntry) GetIp() []uint8 {
	return m.Ip
}

func (m *_BVLCDeleteForeignDeviceTableEntry) GetPort() uint16 {
	return m.Port
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBVLCDeleteForeignDeviceTableEntry factory function for _BVLCDeleteForeignDeviceTableEntry
func NewBVLCDeleteForeignDeviceTableEntry(ip []uint8, port uint16) *_BVLCDeleteForeignDeviceTableEntry {
	_result := &_BVLCDeleteForeignDeviceTableEntry{
		BVLCContract: NewBVLC(),
		Ip:           ip,
		Port:         port,
	}
	_result.BVLCContract.(*_BVLC)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBVLCDeleteForeignDeviceTableEntry(structType any) BVLCDeleteForeignDeviceTableEntry {
	if casted, ok := structType.(BVLCDeleteForeignDeviceTableEntry); ok {
		return casted
	}
	if casted, ok := structType.(*BVLCDeleteForeignDeviceTableEntry); ok {
		return *casted
	}
	return nil
}

func (m *_BVLCDeleteForeignDeviceTableEntry) GetTypeName() string {
	return "BVLCDeleteForeignDeviceTableEntry"
}

func (m *_BVLCDeleteForeignDeviceTableEntry) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BVLCContract.(*_BVLC).getLengthInBits(ctx))

	// Array field
	if len(m.Ip) > 0 {
		lengthInBits += 8 * uint16(len(m.Ip))
	}

	// Simple field (port)
	lengthInBits += 16

	return lengthInBits
}

func (m *_BVLCDeleteForeignDeviceTableEntry) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BVLCDeleteForeignDeviceTableEntry) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BVLC) (__bVLCDeleteForeignDeviceTableEntry BVLCDeleteForeignDeviceTableEntry, err error) {
	m.BVLCContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BVLCDeleteForeignDeviceTableEntry"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BVLCDeleteForeignDeviceTableEntry")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	ip, err := ReadCountArrayField[uint8](ctx, "ip", ReadUnsignedByte(readBuffer, uint8(8)), uint64(int32(4)), codegen.WithByteOrder(binary.BigEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'ip' field"))
	}
	m.Ip = ip

	port, err := ReadSimpleField(ctx, "port", ReadUnsignedShort(readBuffer, uint8(16)), codegen.WithByteOrder(binary.BigEndian))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'port' field"))
	}
	m.Port = port

	if closeErr := readBuffer.CloseContext("BVLCDeleteForeignDeviceTableEntry"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BVLCDeleteForeignDeviceTableEntry")
	}

	return m, nil
}

func (m *_BVLCDeleteForeignDeviceTableEntry) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))), utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BVLCDeleteForeignDeviceTableEntry) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BVLCDeleteForeignDeviceTableEntry"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BVLCDeleteForeignDeviceTableEntry")
		}

		if err := WriteSimpleTypeArrayField(ctx, "ip", m.GetIp(), WriteUnsignedByte(writeBuffer, 8), codegen.WithByteOrder(binary.BigEndian)); err != nil {
			return errors.Wrap(err, "Error serializing 'ip' field")
		}

		if err := WriteSimpleField[uint16](ctx, "port", m.GetPort(), WriteUnsignedShort(writeBuffer, 16), codegen.WithByteOrder(binary.BigEndian)); err != nil {
			return errors.Wrap(err, "Error serializing 'port' field")
		}

		if popErr := writeBuffer.PopContext("BVLCDeleteForeignDeviceTableEntry"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BVLCDeleteForeignDeviceTableEntry")
		}
		return nil
	}
	return m.BVLCContract.(*_BVLC).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BVLCDeleteForeignDeviceTableEntry) IsBVLCDeleteForeignDeviceTableEntry() {}

func (m *_BVLCDeleteForeignDeviceTableEntry) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
