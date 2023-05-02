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

// CIPDataConnected is the corresponding interface of CIPDataConnected
type CIPDataConnected interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetValue returns Value (property field)
	GetValue() uint32
	// GetTagStatus returns TagStatus (property field)
	GetTagStatus() uint16
}

// CIPDataConnectedExactly can be used when we want exactly this type and not a type which fulfills CIPDataConnected.
// This is useful for switch cases.
type CIPDataConnectedExactly interface {
	CIPDataConnected
	isCIPDataConnected() bool
}

// _CIPDataConnected is the data-structure of this message
type _CIPDataConnected struct {
	Value     uint32
	TagStatus uint16
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CIPDataConnected) GetValue() uint32 {
	return m.Value
}

func (m *_CIPDataConnected) GetTagStatus() uint16 {
	return m.TagStatus
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewCIPDataConnected factory function for _CIPDataConnected
func NewCIPDataConnected(value uint32, tagStatus uint16) *_CIPDataConnected {
	return &_CIPDataConnected{Value: value, TagStatus: tagStatus}
}

// Deprecated: use the interface for direct cast
func CastCIPDataConnected(structType any) CIPDataConnected {
	if casted, ok := structType.(CIPDataConnected); ok {
		return casted
	}
	if casted, ok := structType.(*CIPDataConnected); ok {
		return *casted
	}
	return nil
}

func (m *_CIPDataConnected) GetTypeName() string {
	return "CIPDataConnected"
}

func (m *_CIPDataConnected) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (value)
	lengthInBits += 32

	// Simple field (tagStatus)
	lengthInBits += 16

	return lengthInBits
}

func (m *_CIPDataConnected) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func CIPDataConnectedParse(theBytes []byte) (CIPDataConnected, error) {
	return CIPDataConnectedParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes))
}

func CIPDataConnectedParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (CIPDataConnected, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CIPDataConnected"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CIPDataConnected")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (value)
	_value, _valueErr := readBuffer.ReadUint32("value", 32)
	if _valueErr != nil {
		return nil, errors.Wrap(_valueErr, "Error parsing 'value' field of CIPDataConnected")
	}
	value := _value

	// Simple Field (tagStatus)
	_tagStatus, _tagStatusErr := readBuffer.ReadUint16("tagStatus", 16)
	if _tagStatusErr != nil {
		return nil, errors.Wrap(_tagStatusErr, "Error parsing 'tagStatus' field of CIPDataConnected")
	}
	tagStatus := _tagStatus

	if closeErr := readBuffer.CloseContext("CIPDataConnected"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CIPDataConnected")
	}

	// Create the instance
	return &_CIPDataConnected{
		Value:     value,
		TagStatus: tagStatus,
	}, nil
}

func (m *_CIPDataConnected) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_CIPDataConnected) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("CIPDataConnected"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for CIPDataConnected")
	}

	// Simple Field (value)
	value := uint32(m.GetValue())
	_valueErr := writeBuffer.WriteUint32("value", 32, (value))
	if _valueErr != nil {
		return errors.Wrap(_valueErr, "Error serializing 'value' field")
	}

	// Simple Field (tagStatus)
	tagStatus := uint16(m.GetTagStatus())
	_tagStatusErr := writeBuffer.WriteUint16("tagStatus", 16, (tagStatus))
	if _tagStatusErr != nil {
		return errors.Wrap(_tagStatusErr, "Error serializing 'tagStatus' field")
	}

	if popErr := writeBuffer.PopContext("CIPDataConnected"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for CIPDataConnected")
	}
	return nil
}

func (m *_CIPDataConnected) isCIPDataConnected() bool {
	return true
}

func (m *_CIPDataConnected) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
