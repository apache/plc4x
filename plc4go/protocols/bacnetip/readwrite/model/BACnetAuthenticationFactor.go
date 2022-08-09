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

// BACnetAuthenticationFactor is the corresponding interface of BACnetAuthenticationFactor
type BACnetAuthenticationFactor interface {
	utils.LengthAware
	utils.Serializable
	// GetFormatType returns FormatType (property field)
	GetFormatType() BACnetAuthenticationFactorTypeTagged
	// GetFormatClass returns FormatClass (property field)
	GetFormatClass() BACnetContextTagUnsignedInteger
	// GetValue returns Value (property field)
	GetValue() BACnetContextTagOctetString
}

// BACnetAuthenticationFactorExactly can be used when we want exactly this type and not a type which fulfills BACnetAuthenticationFactor.
// This is useful for switch cases.
type BACnetAuthenticationFactorExactly interface {
	BACnetAuthenticationFactor
	isBACnetAuthenticationFactor() bool
}

// _BACnetAuthenticationFactor is the data-structure of this message
type _BACnetAuthenticationFactor struct {
	FormatType  BACnetAuthenticationFactorTypeTagged
	FormatClass BACnetContextTagUnsignedInteger
	Value       BACnetContextTagOctetString
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetAuthenticationFactor) GetFormatType() BACnetAuthenticationFactorTypeTagged {
	return m.FormatType
}

func (m *_BACnetAuthenticationFactor) GetFormatClass() BACnetContextTagUnsignedInteger {
	return m.FormatClass
}

func (m *_BACnetAuthenticationFactor) GetValue() BACnetContextTagOctetString {
	return m.Value
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetAuthenticationFactor factory function for _BACnetAuthenticationFactor
func NewBACnetAuthenticationFactor(formatType BACnetAuthenticationFactorTypeTagged, formatClass BACnetContextTagUnsignedInteger, value BACnetContextTagOctetString) *_BACnetAuthenticationFactor {
	return &_BACnetAuthenticationFactor{FormatType: formatType, FormatClass: formatClass, Value: value}
}

// Deprecated: use the interface for direct cast
func CastBACnetAuthenticationFactor(structType interface{}) BACnetAuthenticationFactor {
	if casted, ok := structType.(BACnetAuthenticationFactor); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetAuthenticationFactor); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetAuthenticationFactor) GetTypeName() string {
	return "BACnetAuthenticationFactor"
}

func (m *_BACnetAuthenticationFactor) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetAuthenticationFactor) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Simple field (formatType)
	lengthInBits += m.FormatType.GetLengthInBits()

	// Simple field (formatClass)
	lengthInBits += m.FormatClass.GetLengthInBits()

	// Simple field (value)
	lengthInBits += m.Value.GetLengthInBits()

	return lengthInBits
}

func (m *_BACnetAuthenticationFactor) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetAuthenticationFactorParse(readBuffer utils.ReadBuffer) (BACnetAuthenticationFactor, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetAuthenticationFactor"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetAuthenticationFactor")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (formatType)
	if pullErr := readBuffer.PullContext("formatType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for formatType")
	}
	_formatType, _formatTypeErr := BACnetAuthenticationFactorTypeTaggedParse(readBuffer, uint8(uint8(0)), TagClass(TagClass_CONTEXT_SPECIFIC_TAGS))
	if _formatTypeErr != nil {
		return nil, errors.Wrap(_formatTypeErr, "Error parsing 'formatType' field of BACnetAuthenticationFactor")
	}
	formatType := _formatType.(BACnetAuthenticationFactorTypeTagged)
	if closeErr := readBuffer.CloseContext("formatType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for formatType")
	}

	// Simple Field (formatClass)
	if pullErr := readBuffer.PullContext("formatClass"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for formatClass")
	}
	_formatClass, _formatClassErr := BACnetContextTagParse(readBuffer, uint8(uint8(1)), BACnetDataType(BACnetDataType_UNSIGNED_INTEGER))
	if _formatClassErr != nil {
		return nil, errors.Wrap(_formatClassErr, "Error parsing 'formatClass' field of BACnetAuthenticationFactor")
	}
	formatClass := _formatClass.(BACnetContextTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("formatClass"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for formatClass")
	}

	// Simple Field (value)
	if pullErr := readBuffer.PullContext("value"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for value")
	}
	_value, _valueErr := BACnetContextTagParse(readBuffer, uint8(uint8(2)), BACnetDataType(BACnetDataType_OCTET_STRING))
	if _valueErr != nil {
		return nil, errors.Wrap(_valueErr, "Error parsing 'value' field of BACnetAuthenticationFactor")
	}
	value := _value.(BACnetContextTagOctetString)
	if closeErr := readBuffer.CloseContext("value"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for value")
	}

	if closeErr := readBuffer.CloseContext("BACnetAuthenticationFactor"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetAuthenticationFactor")
	}

	// Create the instance
	return &_BACnetAuthenticationFactor{
		FormatType:  formatType,
		FormatClass: formatClass,
		Value:       value,
	}, nil
}

func (m *_BACnetAuthenticationFactor) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetAuthenticationFactor"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetAuthenticationFactor")
	}

	// Simple Field (formatType)
	if pushErr := writeBuffer.PushContext("formatType"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for formatType")
	}
	_formatTypeErr := writeBuffer.WriteSerializable(m.GetFormatType())
	if popErr := writeBuffer.PopContext("formatType"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for formatType")
	}
	if _formatTypeErr != nil {
		return errors.Wrap(_formatTypeErr, "Error serializing 'formatType' field")
	}

	// Simple Field (formatClass)
	if pushErr := writeBuffer.PushContext("formatClass"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for formatClass")
	}
	_formatClassErr := writeBuffer.WriteSerializable(m.GetFormatClass())
	if popErr := writeBuffer.PopContext("formatClass"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for formatClass")
	}
	if _formatClassErr != nil {
		return errors.Wrap(_formatClassErr, "Error serializing 'formatClass' field")
	}

	// Simple Field (value)
	if pushErr := writeBuffer.PushContext("value"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for value")
	}
	_valueErr := writeBuffer.WriteSerializable(m.GetValue())
	if popErr := writeBuffer.PopContext("value"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for value")
	}
	if _valueErr != nil {
		return errors.Wrap(_valueErr, "Error serializing 'value' field")
	}

	if popErr := writeBuffer.PopContext("BACnetAuthenticationFactor"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetAuthenticationFactor")
	}
	return nil
}

func (m *_BACnetAuthenticationFactor) isBACnetAuthenticationFactor() bool {
	return true
}

func (m *_BACnetAuthenticationFactor) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
