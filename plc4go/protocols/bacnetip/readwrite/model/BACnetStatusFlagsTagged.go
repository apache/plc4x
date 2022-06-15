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

// BACnetStatusFlagsTagged is the corresponding interface of BACnetStatusFlagsTagged
type BACnetStatusFlagsTagged interface {
	// GetHeader returns Header (property field)
	GetHeader() BACnetTagHeader
	// GetPayload returns Payload (property field)
	GetPayload() BACnetTagPayloadBitString
	// GetInAlarm returns InAlarm (virtual field)
	GetInAlarm() bool
	// GetFault returns Fault (virtual field)
	GetFault() bool
	// GetOverridden returns Overridden (virtual field)
	GetOverridden() bool
	// GetOutOfService returns OutOfService (virtual field)
	GetOutOfService() bool
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetStatusFlagsTagged is the data-structure of this message
type _BACnetStatusFlagsTagged struct {
	Header  BACnetTagHeader
	Payload BACnetTagPayloadBitString

	// Arguments.
	TagNumber uint8
	TagClass  TagClass
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetStatusFlagsTagged) GetHeader() BACnetTagHeader {
	return m.Header
}

func (m *_BACnetStatusFlagsTagged) GetPayload() BACnetTagPayloadBitString {
	return m.Payload
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetStatusFlagsTagged) GetInAlarm() bool {
	return bool(utils.InlineIf(bool(bool((len(m.GetPayload().GetData())) > (0))), func() interface{} { return bool(m.GetPayload().GetData()[0]) }, func() interface{} { return bool(bool(false)) }).(bool))
}

func (m *_BACnetStatusFlagsTagged) GetFault() bool {
	return bool(utils.InlineIf(bool(bool((len(m.GetPayload().GetData())) > (1))), func() interface{} { return bool(m.GetPayload().GetData()[1]) }, func() interface{} { return bool(bool(false)) }).(bool))
}

func (m *_BACnetStatusFlagsTagged) GetOverridden() bool {
	return bool(utils.InlineIf(bool(bool((len(m.GetPayload().GetData())) > (2))), func() interface{} { return bool(m.GetPayload().GetData()[2]) }, func() interface{} { return bool(bool(false)) }).(bool))
}

func (m *_BACnetStatusFlagsTagged) GetOutOfService() bool {
	return bool(utils.InlineIf(bool(bool((len(m.GetPayload().GetData())) > (3))), func() interface{} { return bool(m.GetPayload().GetData()[3]) }, func() interface{} { return bool(bool(false)) }).(bool))
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetStatusFlagsTagged factory function for _BACnetStatusFlagsTagged
func NewBACnetStatusFlagsTagged(header BACnetTagHeader, payload BACnetTagPayloadBitString, tagNumber uint8, tagClass TagClass) *_BACnetStatusFlagsTagged {
	return &_BACnetStatusFlagsTagged{Header: header, Payload: payload, TagNumber: tagNumber, TagClass: tagClass}
}

// Deprecated: use the interface for direct cast
func CastBACnetStatusFlagsTagged(structType interface{}) BACnetStatusFlagsTagged {
	if casted, ok := structType.(BACnetStatusFlagsTagged); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetStatusFlagsTagged); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetStatusFlagsTagged) GetTypeName() string {
	return "BACnetStatusFlagsTagged"
}

func (m *_BACnetStatusFlagsTagged) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetStatusFlagsTagged) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Simple field (header)
	lengthInBits += m.Header.GetLengthInBits()

	// Simple field (payload)
	lengthInBits += m.Payload.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetStatusFlagsTagged) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetStatusFlagsTaggedParse(readBuffer utils.ReadBuffer, tagNumber uint8, tagClass TagClass) (BACnetStatusFlagsTagged, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetStatusFlagsTagged"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetStatusFlagsTagged")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (header)
	if pullErr := readBuffer.PullContext("header"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for header")
	}
	_header, _headerErr := BACnetTagHeaderParse(readBuffer)
	if _headerErr != nil {
		return nil, errors.Wrap(_headerErr, "Error parsing 'header' field")
	}
	header := _header.(BACnetTagHeader)
	if closeErr := readBuffer.CloseContext("header"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for header")
	}

	// Validation
	if !(bool((header.GetTagClass()) == (tagClass))) {
		return nil, errors.WithStack(utils.ParseValidationError{"tag class doesn't match"})
	}

	// Validation
	if !(bool(bool(bool((header.GetTagClass()) == (TagClass_APPLICATION_TAGS)))) || bool(bool(bool((header.GetActualTagNumber()) == (tagNumber))))) {
		return nil, errors.WithStack(utils.ParseAssertError{"tagnumber doesn't match"})
	}

	// Simple Field (payload)
	if pullErr := readBuffer.PullContext("payload"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for payload")
	}
	_payload, _payloadErr := BACnetTagPayloadBitStringParse(readBuffer, uint32(header.GetActualLength()))
	if _payloadErr != nil {
		return nil, errors.Wrap(_payloadErr, "Error parsing 'payload' field")
	}
	payload := _payload.(BACnetTagPayloadBitString)
	if closeErr := readBuffer.CloseContext("payload"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for payload")
	}

	// Virtual field
	_inAlarm := utils.InlineIf(bool(bool((len(payload.GetData())) > (0))), func() interface{} { return bool(payload.GetData()[0]) }, func() interface{} { return bool(bool(false)) }).(bool)
	inAlarm := bool(_inAlarm)
	_ = inAlarm

	// Virtual field
	_fault := utils.InlineIf(bool(bool((len(payload.GetData())) > (1))), func() interface{} { return bool(payload.GetData()[1]) }, func() interface{} { return bool(bool(false)) }).(bool)
	fault := bool(_fault)
	_ = fault

	// Virtual field
	_overridden := utils.InlineIf(bool(bool((len(payload.GetData())) > (2))), func() interface{} { return bool(payload.GetData()[2]) }, func() interface{} { return bool(bool(false)) }).(bool)
	overridden := bool(_overridden)
	_ = overridden

	// Virtual field
	_outOfService := utils.InlineIf(bool(bool((len(payload.GetData())) > (3))), func() interface{} { return bool(payload.GetData()[3]) }, func() interface{} { return bool(bool(false)) }).(bool)
	outOfService := bool(_outOfService)
	_ = outOfService

	if closeErr := readBuffer.CloseContext("BACnetStatusFlagsTagged"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetStatusFlagsTagged")
	}

	// Create the instance
	return NewBACnetStatusFlagsTagged(header, payload, tagNumber, tagClass), nil
}

func (m *_BACnetStatusFlagsTagged) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetStatusFlagsTagged"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetStatusFlagsTagged")
	}

	// Simple Field (header)
	if pushErr := writeBuffer.PushContext("header"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for header")
	}
	_headerErr := writeBuffer.WriteSerializable(m.GetHeader())
	if popErr := writeBuffer.PopContext("header"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for header")
	}
	if _headerErr != nil {
		return errors.Wrap(_headerErr, "Error serializing 'header' field")
	}

	// Simple Field (payload)
	if pushErr := writeBuffer.PushContext("payload"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for payload")
	}
	_payloadErr := writeBuffer.WriteSerializable(m.GetPayload())
	if popErr := writeBuffer.PopContext("payload"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for payload")
	}
	if _payloadErr != nil {
		return errors.Wrap(_payloadErr, "Error serializing 'payload' field")
	}
	// Virtual field
	if _inAlarmErr := writeBuffer.WriteVirtual("inAlarm", m.GetInAlarm()); _inAlarmErr != nil {
		return errors.Wrap(_inAlarmErr, "Error serializing 'inAlarm' field")
	}
	// Virtual field
	if _faultErr := writeBuffer.WriteVirtual("fault", m.GetFault()); _faultErr != nil {
		return errors.Wrap(_faultErr, "Error serializing 'fault' field")
	}
	// Virtual field
	if _overriddenErr := writeBuffer.WriteVirtual("overridden", m.GetOverridden()); _overriddenErr != nil {
		return errors.Wrap(_overriddenErr, "Error serializing 'overridden' field")
	}
	// Virtual field
	if _outOfServiceErr := writeBuffer.WriteVirtual("outOfService", m.GetOutOfService()); _outOfServiceErr != nil {
		return errors.Wrap(_outOfServiceErr, "Error serializing 'outOfService' field")
	}

	if popErr := writeBuffer.PopContext("BACnetStatusFlagsTagged"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetStatusFlagsTagged")
	}
	return nil
}

func (m *_BACnetStatusFlagsTagged) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
