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

// BACnetConstructedDataRecordsSinceNotification is the corresponding interface of BACnetConstructedDataRecordsSinceNotification
type BACnetConstructedDataRecordsSinceNotification interface {
	utils.LengthAware
	utils.Serializable
	BACnetConstructedData
	// GetRecordsSinceNotifications returns RecordsSinceNotifications (property field)
	GetRecordsSinceNotifications() BACnetApplicationTagUnsignedInteger
	// GetActualValue returns ActualValue (virtual field)
	GetActualValue() BACnetApplicationTagUnsignedInteger
}

// BACnetConstructedDataRecordsSinceNotificationExactly can be used when we want exactly this type and not a type which fulfills BACnetConstructedDataRecordsSinceNotification.
// This is useful for switch cases.
type BACnetConstructedDataRecordsSinceNotificationExactly interface {
	BACnetConstructedDataRecordsSinceNotification
	isBACnetConstructedDataRecordsSinceNotification() bool
}

// _BACnetConstructedDataRecordsSinceNotification is the data-structure of this message
type _BACnetConstructedDataRecordsSinceNotification struct {
	*_BACnetConstructedData
	RecordsSinceNotifications BACnetApplicationTagUnsignedInteger
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConstructedDataRecordsSinceNotification) GetObjectTypeArgument() BACnetObjectType {
	return 0
}

func (m *_BACnetConstructedDataRecordsSinceNotification) GetPropertyIdentifierArgument() BACnetPropertyIdentifier {
	return BACnetPropertyIdentifier_RECORDS_SINCE_NOTIFICATION
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConstructedDataRecordsSinceNotification) InitializeParent(parent BACnetConstructedData, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetConstructedDataRecordsSinceNotification) GetParent() BACnetConstructedData {
	return m._BACnetConstructedData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConstructedDataRecordsSinceNotification) GetRecordsSinceNotifications() BACnetApplicationTagUnsignedInteger {
	return m.RecordsSinceNotifications
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_BACnetConstructedDataRecordsSinceNotification) GetActualValue() BACnetApplicationTagUnsignedInteger {
	return CastBACnetApplicationTagUnsignedInteger(m.GetRecordsSinceNotifications())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetConstructedDataRecordsSinceNotification factory function for _BACnetConstructedDataRecordsSinceNotification
func NewBACnetConstructedDataRecordsSinceNotification(recordsSinceNotifications BACnetApplicationTagUnsignedInteger, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, arrayIndexArgument BACnetTagPayloadUnsignedInteger) *_BACnetConstructedDataRecordsSinceNotification {
	_result := &_BACnetConstructedDataRecordsSinceNotification{
		RecordsSinceNotifications: recordsSinceNotifications,
		_BACnetConstructedData:    NewBACnetConstructedData(openingTag, peekedTagHeader, closingTag, tagNumber, arrayIndexArgument),
	}
	_result._BACnetConstructedData._BACnetConstructedDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetConstructedDataRecordsSinceNotification(structType interface{}) BACnetConstructedDataRecordsSinceNotification {
	if casted, ok := structType.(BACnetConstructedDataRecordsSinceNotification); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConstructedDataRecordsSinceNotification); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConstructedDataRecordsSinceNotification) GetTypeName() string {
	return "BACnetConstructedDataRecordsSinceNotification"
}

func (m *_BACnetConstructedDataRecordsSinceNotification) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetConstructedDataRecordsSinceNotification) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (recordsSinceNotifications)
	lengthInBits += m.RecordsSinceNotifications.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_BACnetConstructedDataRecordsSinceNotification) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetConstructedDataRecordsSinceNotificationParse(readBuffer utils.ReadBuffer, tagNumber uint8, objectTypeArgument BACnetObjectType, propertyIdentifierArgument BACnetPropertyIdentifier, arrayIndexArgument BACnetTagPayloadUnsignedInteger) (BACnetConstructedDataRecordsSinceNotification, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConstructedDataRecordsSinceNotification"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConstructedDataRecordsSinceNotification")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (recordsSinceNotifications)
	if pullErr := readBuffer.PullContext("recordsSinceNotifications"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for recordsSinceNotifications")
	}
	_recordsSinceNotifications, _recordsSinceNotificationsErr := BACnetApplicationTagParse(readBuffer)
	if _recordsSinceNotificationsErr != nil {
		return nil, errors.Wrap(_recordsSinceNotificationsErr, "Error parsing 'recordsSinceNotifications' field of BACnetConstructedDataRecordsSinceNotification")
	}
	recordsSinceNotifications := _recordsSinceNotifications.(BACnetApplicationTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("recordsSinceNotifications"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for recordsSinceNotifications")
	}

	// Virtual field
	_actualValue := recordsSinceNotifications
	actualValue := _actualValue
	_ = actualValue

	if closeErr := readBuffer.CloseContext("BACnetConstructedDataRecordsSinceNotification"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConstructedDataRecordsSinceNotification")
	}

	// Create a partially initialized instance
	_child := &_BACnetConstructedDataRecordsSinceNotification{
		_BACnetConstructedData: &_BACnetConstructedData{
			TagNumber:          tagNumber,
			ArrayIndexArgument: arrayIndexArgument,
		},
		RecordsSinceNotifications: recordsSinceNotifications,
	}
	_child._BACnetConstructedData._BACnetConstructedDataChildRequirements = _child
	return _child, nil
}

func (m *_BACnetConstructedDataRecordsSinceNotification) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConstructedDataRecordsSinceNotification"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConstructedDataRecordsSinceNotification")
		}

		// Simple Field (recordsSinceNotifications)
		if pushErr := writeBuffer.PushContext("recordsSinceNotifications"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for recordsSinceNotifications")
		}
		_recordsSinceNotificationsErr := writeBuffer.WriteSerializable(m.GetRecordsSinceNotifications())
		if popErr := writeBuffer.PopContext("recordsSinceNotifications"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for recordsSinceNotifications")
		}
		if _recordsSinceNotificationsErr != nil {
			return errors.Wrap(_recordsSinceNotificationsErr, "Error serializing 'recordsSinceNotifications' field")
		}
		// Virtual field
		if _actualValueErr := writeBuffer.WriteVirtual("actualValue", m.GetActualValue()); _actualValueErr != nil {
			return errors.Wrap(_actualValueErr, "Error serializing 'actualValue' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConstructedDataRecordsSinceNotification"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConstructedDataRecordsSinceNotification")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetConstructedDataRecordsSinceNotification) isBACnetConstructedDataRecordsSinceNotification() bool {
	return true
}

func (m *_BACnetConstructedDataRecordsSinceNotification) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
