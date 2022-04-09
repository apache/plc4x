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
type AdsDeviceNotificationResponse struct {
	*AdsData
}

// The corresponding interface
type IAdsDeviceNotificationResponse interface {
	IAdsData
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
func (m *AdsDeviceNotificationResponse) GetCommandId() CommandId {
	return CommandId_ADS_DEVICE_NOTIFICATION
}

func (m *AdsDeviceNotificationResponse) GetResponse() bool {
	return bool(true)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *AdsDeviceNotificationResponse) InitializeParent(parent *AdsData) {}

func (m *AdsDeviceNotificationResponse) GetParent() *AdsData {
	return m.AdsData
}

// NewAdsDeviceNotificationResponse factory function for AdsDeviceNotificationResponse
func NewAdsDeviceNotificationResponse() *AdsDeviceNotificationResponse {
	_result := &AdsDeviceNotificationResponse{
		AdsData: NewAdsData(),
	}
	_result.Child = _result
	return _result
}

func CastAdsDeviceNotificationResponse(structType interface{}) *AdsDeviceNotificationResponse {
	if casted, ok := structType.(AdsDeviceNotificationResponse); ok {
		return &casted
	}
	if casted, ok := structType.(*AdsDeviceNotificationResponse); ok {
		return casted
	}
	if casted, ok := structType.(AdsData); ok {
		return CastAdsDeviceNotificationResponse(casted.Child)
	}
	if casted, ok := structType.(*AdsData); ok {
		return CastAdsDeviceNotificationResponse(casted.Child)
	}
	return nil
}

func (m *AdsDeviceNotificationResponse) GetTypeName() string {
	return "AdsDeviceNotificationResponse"
}

func (m *AdsDeviceNotificationResponse) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *AdsDeviceNotificationResponse) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	return lengthInBits
}

func (m *AdsDeviceNotificationResponse) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func AdsDeviceNotificationResponseParse(readBuffer utils.ReadBuffer, commandId CommandId, response bool) (*AdsDeviceNotificationResponse, error) {
	if pullErr := readBuffer.PullContext("AdsDeviceNotificationResponse"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("AdsDeviceNotificationResponse"); closeErr != nil {
		return nil, closeErr
	}

	// Create a partially initialized instance
	_child := &AdsDeviceNotificationResponse{
		AdsData: &AdsData{},
	}
	_child.AdsData.Child = _child
	return _child, nil
}

func (m *AdsDeviceNotificationResponse) Serialize(writeBuffer utils.WriteBuffer) error {
	ser := func() error {
		if pushErr := writeBuffer.PushContext("AdsDeviceNotificationResponse"); pushErr != nil {
			return pushErr
		}

		if popErr := writeBuffer.PopContext("AdsDeviceNotificationResponse"); popErr != nil {
			return popErr
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *AdsDeviceNotificationResponse) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
