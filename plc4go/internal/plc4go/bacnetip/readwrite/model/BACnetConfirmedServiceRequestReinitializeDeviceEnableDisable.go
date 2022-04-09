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
	"io"
)

// Code generated by code-generation. DO NOT EDIT.

// The data-structure of this message
type BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable struct {
	RawData *BACnetContextTagEnumerated

	// Arguments.
	TagNumber uint8
}

// The corresponding interface
type IBACnetConfirmedServiceRequestReinitializeDeviceEnableDisable interface {
	// GetRawData returns RawData (property field)
	GetRawData() *BACnetContextTagEnumerated
	// GetIsEnable returns IsEnable (virtual field)
	GetIsEnable() bool
	// GetIsDisable returns IsDisable (virtual field)
	GetIsDisable() bool
	// GetIsDisableInitiation returns IsDisableInitiation (virtual field)
	GetIsDisableInitiation() bool
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
func (m *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable) GetRawData() *BACnetContextTagEnumerated {
	return m.RawData
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////
func (m *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable) GetIsEnable() bool {
	rawData := m.RawData
	_ = rawData
	return bool(bool(bool((m.GetRawData()) != (nil))) && bool(bool(((*m.GetRawData()).GetPayload().GetActualValue()) == (0))))
}

func (m *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable) GetIsDisable() bool {
	rawData := m.RawData
	_ = rawData
	return bool(bool(bool((m.GetRawData()) != (nil))) && bool(bool(((*m.GetRawData()).GetPayload().GetActualValue()) == (1))))
}

func (m *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable) GetIsDisableInitiation() bool {
	rawData := m.RawData
	_ = rawData
	return bool(bool(bool((m.GetRawData()) != (nil))) && bool(bool(((*m.GetRawData()).GetPayload().GetActualValue()) == (2))))
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConfirmedServiceRequestReinitializeDeviceEnableDisable factory function for BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable
func NewBACnetConfirmedServiceRequestReinitializeDeviceEnableDisable(rawData *BACnetContextTagEnumerated, tagNumber uint8) *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable {
	return &BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable{RawData: rawData, TagNumber: tagNumber}
}

func CastBACnetConfirmedServiceRequestReinitializeDeviceEnableDisable(structType interface{}) *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable {
	if casted, ok := structType.(BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable); ok {
		return &casted
	}
	if casted, ok := structType.(*BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable); ok {
		return casted
	}
	return nil
}

func (m *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable) GetTypeName() string {
	return "BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable"
}

func (m *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Optional Field (rawData)
	if m.RawData != nil {
		lengthInBits += (*m.RawData).GetLengthInBits()
	}

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConfirmedServiceRequestReinitializeDeviceEnableDisableParse(readBuffer utils.ReadBuffer, tagNumber uint8) (*BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable, error) {
	if pullErr := readBuffer.PullContext("BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	// Optional Field (rawData) (Can be skipped, if a given expression evaluates to false)
	var rawData *BACnetContextTagEnumerated = nil
	{
		currentPos = readBuffer.GetPos()
		if pullErr := readBuffer.PullContext("rawData"); pullErr != nil {
			return nil, pullErr
		}
		_val, _err := BACnetContextTagParse(readBuffer, tagNumber, BACnetDataType_ENUMERATED)
		switch {
		case errors.Is(_err, utils.ParseAssertError{}) || errors.Is(_err, io.EOF):
			readBuffer.Reset(currentPos)
		case _err != nil:
			return nil, errors.Wrap(_err, "Error parsing 'rawData' field")
		default:
			rawData = CastBACnetContextTagEnumerated(_val)
			if closeErr := readBuffer.CloseContext("rawData"); closeErr != nil {
				return nil, closeErr
			}
		}
	}

	// Virtual field
	_isEnable := bool(bool((rawData) != (nil))) && bool(bool(((*rawData).GetPayload().GetActualValue()) == (0)))
	isEnable := bool(_isEnable)
	_ = isEnable

	// Virtual field
	_isDisable := bool(bool((rawData) != (nil))) && bool(bool(((*rawData).GetPayload().GetActualValue()) == (1)))
	isDisable := bool(_isDisable)
	_ = isDisable

	// Virtual field
	_isDisableInitiation := bool(bool((rawData) != (nil))) && bool(bool(((*rawData).GetPayload().GetActualValue()) == (2)))
	isDisableInitiation := bool(_isDisableInitiation)
	_ = isDisableInitiation

	if closeErr := readBuffer.CloseContext("BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable"); closeErr != nil {
		return nil, closeErr
	}

	// Create the instance
	return NewBACnetConfirmedServiceRequestReinitializeDeviceEnableDisable(rawData, tagNumber), nil
}

func (m *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable) Serialize(writeBuffer utils.WriteBuffer) error {
	if pushErr := writeBuffer.PushContext("BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable"); pushErr != nil {
		return pushErr
	}

	// Optional Field (rawData) (Can be skipped, if the value is null)
	var rawData *BACnetContextTagEnumerated = nil
	if m.RawData != nil {
		if pushErr := writeBuffer.PushContext("rawData"); pushErr != nil {
			return pushErr
		}
		rawData = m.RawData
		_rawDataErr := rawData.Serialize(writeBuffer)
		if popErr := writeBuffer.PopContext("rawData"); popErr != nil {
			return popErr
		}
		if _rawDataErr != nil {
			return errors.Wrap(_rawDataErr, "Error serializing 'rawData' field")
		}
	}
	// Virtual field
	if _isEnableErr := writeBuffer.WriteVirtual("isEnable", m.GetIsEnable()); _isEnableErr != nil {
		return errors.Wrap(_isEnableErr, "Error serializing 'isEnable' field")
	}
	// Virtual field
	if _isDisableErr := writeBuffer.WriteVirtual("isDisable", m.GetIsDisable()); _isDisableErr != nil {
		return errors.Wrap(_isDisableErr, "Error serializing 'isDisable' field")
	}
	// Virtual field
	if _isDisableInitiationErr := writeBuffer.WriteVirtual("isDisableInitiation", m.GetIsDisableInitiation()); _isDisableInitiationErr != nil {
		return errors.Wrap(_isDisableInitiationErr, "Error serializing 'isDisableInitiation' field")
	}

	if popErr := writeBuffer.PopContext("BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable"); popErr != nil {
		return popErr
	}
	return nil
}

func (m *BACnetConfirmedServiceRequestReinitializeDeviceEnableDisable) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
