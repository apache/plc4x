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

// AdsDeviceNotificationRequest is the corresponding interface of AdsDeviceNotificationRequest
type AdsDeviceNotificationRequest interface {
	utils.LengthAware
	utils.Serializable
	AdsData
	// GetLength returns Length (property field)
	GetLength() uint32
	// GetStamps returns Stamps (property field)
	GetStamps() uint32
	// GetAdsStampHeaders returns AdsStampHeaders (property field)
	GetAdsStampHeaders() []AdsStampHeader
}

// AdsDeviceNotificationRequestExactly can be used when we want exactly this type and not a type which fulfills AdsDeviceNotificationRequest.
// This is useful for switch cases.
type AdsDeviceNotificationRequestExactly interface {
	AdsDeviceNotificationRequest
	isAdsDeviceNotificationRequest() bool
}

// _AdsDeviceNotificationRequest is the data-structure of this message
type _AdsDeviceNotificationRequest struct {
	*_AdsData
	Length          uint32
	Stamps          uint32
	AdsStampHeaders []AdsStampHeader
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_AdsDeviceNotificationRequest) GetCommandId() CommandId {
	return CommandId_ADS_DEVICE_NOTIFICATION
}

func (m *_AdsDeviceNotificationRequest) GetResponse() bool {
	return bool(false)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_AdsDeviceNotificationRequest) InitializeParent(parent AdsData) {}

func (m *_AdsDeviceNotificationRequest) GetParent() AdsData {
	return m._AdsData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AdsDeviceNotificationRequest) GetLength() uint32 {
	return m.Length
}

func (m *_AdsDeviceNotificationRequest) GetStamps() uint32 {
	return m.Stamps
}

func (m *_AdsDeviceNotificationRequest) GetAdsStampHeaders() []AdsStampHeader {
	return m.AdsStampHeaders
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewAdsDeviceNotificationRequest factory function for _AdsDeviceNotificationRequest
func NewAdsDeviceNotificationRequest(length uint32, stamps uint32, adsStampHeaders []AdsStampHeader) *_AdsDeviceNotificationRequest {
	_result := &_AdsDeviceNotificationRequest{
		Length:          length,
		Stamps:          stamps,
		AdsStampHeaders: adsStampHeaders,
		_AdsData:        NewAdsData(),
	}
	_result._AdsData._AdsDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastAdsDeviceNotificationRequest(structType interface{}) AdsDeviceNotificationRequest {
	if casted, ok := structType.(AdsDeviceNotificationRequest); ok {
		return casted
	}
	if casted, ok := structType.(*AdsDeviceNotificationRequest); ok {
		return *casted
	}
	return nil
}

func (m *_AdsDeviceNotificationRequest) GetTypeName() string {
	return "AdsDeviceNotificationRequest"
}

func (m *_AdsDeviceNotificationRequest) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_AdsDeviceNotificationRequest) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (length)
	lengthInBits += 32

	// Simple field (stamps)
	lengthInBits += 32

	// Array field
	if len(m.AdsStampHeaders) > 0 {
		for i, element := range m.AdsStampHeaders {
			last := i == len(m.AdsStampHeaders)-1
			lengthInBits += element.(interface{ GetLengthInBitsConditional(bool) uint16 }).GetLengthInBitsConditional(last)
		}
	}

	return lengthInBits
}

func (m *_AdsDeviceNotificationRequest) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func AdsDeviceNotificationRequestParse(readBuffer utils.ReadBuffer, commandId CommandId, response bool) (AdsDeviceNotificationRequest, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AdsDeviceNotificationRequest"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AdsDeviceNotificationRequest")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (length)
	_length, _lengthErr := readBuffer.ReadUint32("length", 32)
	if _lengthErr != nil {
		return nil, errors.Wrap(_lengthErr, "Error parsing 'length' field of AdsDeviceNotificationRequest")
	}
	length := _length

	// Simple Field (stamps)
	_stamps, _stampsErr := readBuffer.ReadUint32("stamps", 32)
	if _stampsErr != nil {
		return nil, errors.Wrap(_stampsErr, "Error parsing 'stamps' field of AdsDeviceNotificationRequest")
	}
	stamps := _stamps

	// Array field (adsStampHeaders)
	if pullErr := readBuffer.PullContext("adsStampHeaders", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for adsStampHeaders")
	}
	// Count array
	adsStampHeaders := make([]AdsStampHeader, stamps)
	// This happens when the size is set conditional to 0
	if len(adsStampHeaders) == 0 {
		adsStampHeaders = nil
	}
	{
		for curItem := uint16(0); curItem < uint16(stamps); curItem++ {
			_item, _err := AdsStampHeaderParse(readBuffer)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'adsStampHeaders' field of AdsDeviceNotificationRequest")
			}
			adsStampHeaders[curItem] = _item.(AdsStampHeader)
		}
	}
	if closeErr := readBuffer.CloseContext("adsStampHeaders", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for adsStampHeaders")
	}

	if closeErr := readBuffer.CloseContext("AdsDeviceNotificationRequest"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AdsDeviceNotificationRequest")
	}

	// Create a partially initialized instance
	_child := &_AdsDeviceNotificationRequest{
		_AdsData:        &_AdsData{},
		Length:          length,
		Stamps:          stamps,
		AdsStampHeaders: adsStampHeaders,
	}
	_child._AdsData._AdsDataChildRequirements = _child
	return _child, nil
}

func (m *_AdsDeviceNotificationRequest) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("AdsDeviceNotificationRequest"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for AdsDeviceNotificationRequest")
		}

		// Simple Field (length)
		length := uint32(m.GetLength())
		_lengthErr := writeBuffer.WriteUint32("length", 32, (length))
		if _lengthErr != nil {
			return errors.Wrap(_lengthErr, "Error serializing 'length' field")
		}

		// Simple Field (stamps)
		stamps := uint32(m.GetStamps())
		_stampsErr := writeBuffer.WriteUint32("stamps", 32, (stamps))
		if _stampsErr != nil {
			return errors.Wrap(_stampsErr, "Error serializing 'stamps' field")
		}

		// Array Field (adsStampHeaders)
		if pushErr := writeBuffer.PushContext("adsStampHeaders", utils.WithRenderAsList(true)); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for adsStampHeaders")
		}
		for _, _element := range m.GetAdsStampHeaders() {
			_elementErr := writeBuffer.WriteSerializable(_element)
			if _elementErr != nil {
				return errors.Wrap(_elementErr, "Error serializing 'adsStampHeaders' field")
			}
		}
		if popErr := writeBuffer.PopContext("adsStampHeaders", utils.WithRenderAsList(true)); popErr != nil {
			return errors.Wrap(popErr, "Error popping for adsStampHeaders")
		}

		if popErr := writeBuffer.PopContext("AdsDeviceNotificationRequest"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for AdsDeviceNotificationRequest")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_AdsDeviceNotificationRequest) isAdsDeviceNotificationRequest() bool {
	return true
}

func (m *_AdsDeviceNotificationRequest) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
