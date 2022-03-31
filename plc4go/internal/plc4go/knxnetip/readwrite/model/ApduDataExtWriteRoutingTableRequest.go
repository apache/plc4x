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
type ApduDataExtWriteRoutingTableRequest struct {
	*ApduDataExt

	// Arguments.
	Length uint8
}

// The corresponding interface
type IApduDataExtWriteRoutingTableRequest interface {
	IApduDataExt
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
func (m *ApduDataExtWriteRoutingTableRequest) GetExtApciType() uint8 {
	return 0x03
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *ApduDataExtWriteRoutingTableRequest) InitializeParent(parent *ApduDataExt) {}

func (m *ApduDataExtWriteRoutingTableRequest) GetParent() *ApduDataExt {
	return m.ApduDataExt
}

// NewApduDataExtWriteRoutingTableRequest factory function for ApduDataExtWriteRoutingTableRequest
func NewApduDataExtWriteRoutingTableRequest(length uint8) *ApduDataExtWriteRoutingTableRequest {
	_result := &ApduDataExtWriteRoutingTableRequest{
		ApduDataExt: NewApduDataExt(length),
	}
	_result.Child = _result
	return _result
}

func CastApduDataExtWriteRoutingTableRequest(structType interface{}) *ApduDataExtWriteRoutingTableRequest {
	if casted, ok := structType.(ApduDataExtWriteRoutingTableRequest); ok {
		return &casted
	}
	if casted, ok := structType.(*ApduDataExtWriteRoutingTableRequest); ok {
		return casted
	}
	if casted, ok := structType.(ApduDataExt); ok {
		return CastApduDataExtWriteRoutingTableRequest(casted.Child)
	}
	if casted, ok := structType.(*ApduDataExt); ok {
		return CastApduDataExtWriteRoutingTableRequest(casted.Child)
	}
	return nil
}

func (m *ApduDataExtWriteRoutingTableRequest) GetTypeName() string {
	return "ApduDataExtWriteRoutingTableRequest"
}

func (m *ApduDataExtWriteRoutingTableRequest) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *ApduDataExtWriteRoutingTableRequest) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	return lengthInBits
}

func (m *ApduDataExtWriteRoutingTableRequest) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func ApduDataExtWriteRoutingTableRequestParse(readBuffer utils.ReadBuffer, length uint8) (*ApduDataExtWriteRoutingTableRequest, error) {
	if pullErr := readBuffer.PullContext("ApduDataExtWriteRoutingTableRequest"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("ApduDataExtWriteRoutingTableRequest"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &ApduDataExtWriteRoutingTableRequest{
		ApduDataExt: &ApduDataExt{},
	}
	_child.ApduDataExt.Child = _child
	return _child, nil
}

func (m *ApduDataExtWriteRoutingTableRequest) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ApduDataExtWriteRoutingTableRequest"); pushErr != nil {
			return pushErr
		}

		if popErr := writeBuffer.PopContext("ApduDataExtWriteRoutingTableRequest"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *ApduDataExtWriteRoutingTableRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
