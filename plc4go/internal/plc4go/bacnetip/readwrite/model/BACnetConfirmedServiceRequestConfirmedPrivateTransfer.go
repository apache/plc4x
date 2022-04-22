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

// BACnetConfirmedServiceRequestConfirmedPrivateTransfer is the data-structure of this message
type BACnetConfirmedServiceRequestConfirmedPrivateTransfer struct {
	*BACnetConfirmedServiceRequest

	// Arguments.
	ServiceRequestLength uint16
}

// IBACnetConfirmedServiceRequestConfirmedPrivateTransfer is the corresponding interface of BACnetConfirmedServiceRequestConfirmedPrivateTransfer
type IBACnetConfirmedServiceRequestConfirmedPrivateTransfer interface {
	IBACnetConfirmedServiceRequest
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

func (m *BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetServiceChoice() uint8 {
	return 0x12
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *BACnetConfirmedServiceRequestConfirmedPrivateTransfer) InitializeParent(parent *BACnetConfirmedServiceRequest) {
}

func (m *BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetParent() *BACnetConfirmedServiceRequest {
	return m.BACnetConfirmedServiceRequest
}

// NewBACnetConfirmedServiceRequestConfirmedPrivateTransfer factory function for BACnetConfirmedServiceRequestConfirmedPrivateTransfer
func NewBACnetConfirmedServiceRequestConfirmedPrivateTransfer(serviceRequestLength uint16) *BACnetConfirmedServiceRequestConfirmedPrivateTransfer {
	_result := &BACnetConfirmedServiceRequestConfirmedPrivateTransfer{
		BACnetConfirmedServiceRequest: NewBACnetConfirmedServiceRequest(serviceRequestLength),
	}
	_result.Child = _result
	return _result
}

func CastBACnetConfirmedServiceRequestConfirmedPrivateTransfer(structType interface{}) *BACnetConfirmedServiceRequestConfirmedPrivateTransfer {
	if casted, ok := structType.(BACnetConfirmedServiceRequestConfirmedPrivateTransfer); ok {
		return &casted
	}
	if casted, ok := structType.(*BACnetConfirmedServiceRequestConfirmedPrivateTransfer); ok {
		return casted
	}
	if casted, ok := structType.(BACnetConfirmedServiceRequest); ok {
		return CastBACnetConfirmedServiceRequestConfirmedPrivateTransfer(casted.Child)
	}
	if casted, ok := structType.(*BACnetConfirmedServiceRequest); ok {
		return CastBACnetConfirmedServiceRequestConfirmedPrivateTransfer(casted.Child)
	}
	return nil
}

func (m *BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetTypeName() string {
	return "BACnetConfirmedServiceRequestConfirmedPrivateTransfer"
}

func (m *BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	return lengthInBits
}

func (m *BACnetConfirmedServiceRequestConfirmedPrivateTransfer) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConfirmedServiceRequestConfirmedPrivateTransferParse(readBuffer utils.ReadBuffer, serviceRequestLength uint16) (*BACnetConfirmedServiceRequestConfirmedPrivateTransfer, error) {
	if pullErr := readBuffer.PullContext("BACnetConfirmedServiceRequestConfirmedPrivateTransfer"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("BACnetConfirmedServiceRequestConfirmedPrivateTransfer"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &BACnetConfirmedServiceRequestConfirmedPrivateTransfer{
		BACnetConfirmedServiceRequest: &BACnetConfirmedServiceRequest{},
	}
	_child.BACnetConfirmedServiceRequest.Child = _child
	return _child, nil
}

func (m *BACnetConfirmedServiceRequestConfirmedPrivateTransfer) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConfirmedServiceRequestConfirmedPrivateTransfer"); pushErr != nil {
			return pushErr
		}

		if popErr := writeBuffer.PopContext("BACnetConfirmedServiceRequestConfirmedPrivateTransfer"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *BACnetConfirmedServiceRequestConfirmedPrivateTransfer) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
