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

// BACnetNotificationParametersChangeOfStatusFlags is the corresponding interface of BACnetNotificationParametersChangeOfStatusFlags
type BACnetNotificationParametersChangeOfStatusFlags interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetNotificationParameters
	// GetInnerOpeningTag returns InnerOpeningTag (property field)
	GetInnerOpeningTag() BACnetOpeningTag
	// GetPresentValue returns PresentValue (property field)
	GetPresentValue() BACnetConstructedData
	// GetReferencedFlags returns ReferencedFlags (property field)
	GetReferencedFlags() BACnetStatusFlagsTagged
	// GetInnerClosingTag returns InnerClosingTag (property field)
	GetInnerClosingTag() BACnetClosingTag
}

// BACnetNotificationParametersChangeOfStatusFlagsExactly can be used when we want exactly this type and not a type which fulfills BACnetNotificationParametersChangeOfStatusFlags.
// This is useful for switch cases.
type BACnetNotificationParametersChangeOfStatusFlagsExactly interface {
	BACnetNotificationParametersChangeOfStatusFlags
	isBACnetNotificationParametersChangeOfStatusFlags() bool
}

// _BACnetNotificationParametersChangeOfStatusFlags is the data-structure of this message
type _BACnetNotificationParametersChangeOfStatusFlags struct {
	*_BACnetNotificationParameters
	InnerOpeningTag BACnetOpeningTag
	PresentValue    BACnetConstructedData
	ReferencedFlags BACnetStatusFlagsTagged
	InnerClosingTag BACnetClosingTag
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetNotificationParametersChangeOfStatusFlags) InitializeParent(parent BACnetNotificationParameters, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag) {
	m.OpeningTag = openingTag
	m.PeekedTagHeader = peekedTagHeader
	m.ClosingTag = closingTag
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) GetParent() BACnetNotificationParameters {
	return m._BACnetNotificationParameters
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetNotificationParametersChangeOfStatusFlags) GetInnerOpeningTag() BACnetOpeningTag {
	return m.InnerOpeningTag
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) GetPresentValue() BACnetConstructedData {
	return m.PresentValue
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) GetReferencedFlags() BACnetStatusFlagsTagged {
	return m.ReferencedFlags
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) GetInnerClosingTag() BACnetClosingTag {
	return m.InnerClosingTag
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetNotificationParametersChangeOfStatusFlags factory function for _BACnetNotificationParametersChangeOfStatusFlags
func NewBACnetNotificationParametersChangeOfStatusFlags(innerOpeningTag BACnetOpeningTag, presentValue BACnetConstructedData, referencedFlags BACnetStatusFlagsTagged, innerClosingTag BACnetClosingTag, openingTag BACnetOpeningTag, peekedTagHeader BACnetTagHeader, closingTag BACnetClosingTag, tagNumber uint8, objectTypeArgument BACnetObjectType) *_BACnetNotificationParametersChangeOfStatusFlags {
	_result := &_BACnetNotificationParametersChangeOfStatusFlags{
		InnerOpeningTag:               innerOpeningTag,
		PresentValue:                  presentValue,
		ReferencedFlags:               referencedFlags,
		InnerClosingTag:               innerClosingTag,
		_BACnetNotificationParameters: NewBACnetNotificationParameters(openingTag, peekedTagHeader, closingTag, tagNumber, objectTypeArgument),
	}
	_result._BACnetNotificationParameters._BACnetNotificationParametersChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetNotificationParametersChangeOfStatusFlags(structType any) BACnetNotificationParametersChangeOfStatusFlags {
	if casted, ok := structType.(BACnetNotificationParametersChangeOfStatusFlags); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetNotificationParametersChangeOfStatusFlags); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) GetTypeName() string {
	return "BACnetNotificationParametersChangeOfStatusFlags"
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (innerOpeningTag)
	lengthInBits += m.InnerOpeningTag.GetLengthInBits(ctx)

	// Simple field (presentValue)
	lengthInBits += m.PresentValue.GetLengthInBits(ctx)

	// Simple field (referencedFlags)
	lengthInBits += m.ReferencedFlags.GetLengthInBits(ctx)

	// Simple field (innerClosingTag)
	lengthInBits += m.InnerClosingTag.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetNotificationParametersChangeOfStatusFlagsParse(theBytes []byte, peekedTagNumber uint8, tagNumber uint8, objectTypeArgument BACnetObjectType) (BACnetNotificationParametersChangeOfStatusFlags, error) {
	return BACnetNotificationParametersChangeOfStatusFlagsParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes), peekedTagNumber, tagNumber, objectTypeArgument)
}

func BACnetNotificationParametersChangeOfStatusFlagsParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, peekedTagNumber uint8, tagNumber uint8, objectTypeArgument BACnetObjectType) (BACnetNotificationParametersChangeOfStatusFlags, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetNotificationParametersChangeOfStatusFlags"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetNotificationParametersChangeOfStatusFlags")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (innerOpeningTag)
	if pullErr := readBuffer.PullContext("innerOpeningTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for innerOpeningTag")
	}
	_innerOpeningTag, _innerOpeningTagErr := BACnetOpeningTagParseWithBuffer(ctx, readBuffer, uint8(peekedTagNumber))
	if _innerOpeningTagErr != nil {
		return nil, errors.Wrap(_innerOpeningTagErr, "Error parsing 'innerOpeningTag' field of BACnetNotificationParametersChangeOfStatusFlags")
	}
	innerOpeningTag := _innerOpeningTag.(BACnetOpeningTag)
	if closeErr := readBuffer.CloseContext("innerOpeningTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for innerOpeningTag")
	}

	// Simple Field (presentValue)
	if pullErr := readBuffer.PullContext("presentValue"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for presentValue")
	}
	_presentValue, _presentValueErr := BACnetConstructedDataParseWithBuffer(ctx, readBuffer, uint8(uint8(0)), BACnetObjectType(objectTypeArgument), BACnetPropertyIdentifier(BACnetPropertyIdentifier_VENDOR_PROPRIETARY_VALUE), nil)
	if _presentValueErr != nil {
		return nil, errors.Wrap(_presentValueErr, "Error parsing 'presentValue' field of BACnetNotificationParametersChangeOfStatusFlags")
	}
	presentValue := _presentValue.(BACnetConstructedData)
	if closeErr := readBuffer.CloseContext("presentValue"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for presentValue")
	}

	// Simple Field (referencedFlags)
	if pullErr := readBuffer.PullContext("referencedFlags"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for referencedFlags")
	}
	_referencedFlags, _referencedFlagsErr := BACnetStatusFlagsTaggedParseWithBuffer(ctx, readBuffer, uint8(uint8(1)), TagClass(TagClass_CONTEXT_SPECIFIC_TAGS))
	if _referencedFlagsErr != nil {
		return nil, errors.Wrap(_referencedFlagsErr, "Error parsing 'referencedFlags' field of BACnetNotificationParametersChangeOfStatusFlags")
	}
	referencedFlags := _referencedFlags.(BACnetStatusFlagsTagged)
	if closeErr := readBuffer.CloseContext("referencedFlags"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for referencedFlags")
	}

	// Simple Field (innerClosingTag)
	if pullErr := readBuffer.PullContext("innerClosingTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for innerClosingTag")
	}
	_innerClosingTag, _innerClosingTagErr := BACnetClosingTagParseWithBuffer(ctx, readBuffer, uint8(peekedTagNumber))
	if _innerClosingTagErr != nil {
		return nil, errors.Wrap(_innerClosingTagErr, "Error parsing 'innerClosingTag' field of BACnetNotificationParametersChangeOfStatusFlags")
	}
	innerClosingTag := _innerClosingTag.(BACnetClosingTag)
	if closeErr := readBuffer.CloseContext("innerClosingTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for innerClosingTag")
	}

	if closeErr := readBuffer.CloseContext("BACnetNotificationParametersChangeOfStatusFlags"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetNotificationParametersChangeOfStatusFlags")
	}

	// Create a partially initialized instance
	_child := &_BACnetNotificationParametersChangeOfStatusFlags{
		_BACnetNotificationParameters: &_BACnetNotificationParameters{
			TagNumber:          tagNumber,
			ObjectTypeArgument: objectTypeArgument,
		},
		InnerOpeningTag: innerOpeningTag,
		PresentValue:    presentValue,
		ReferencedFlags: referencedFlags,
		InnerClosingTag: innerClosingTag,
	}
	_child._BACnetNotificationParameters._BACnetNotificationParametersChildRequirements = _child
	return _child, nil
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetNotificationParametersChangeOfStatusFlags"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetNotificationParametersChangeOfStatusFlags")
		}

		// Simple Field (innerOpeningTag)
		if pushErr := writeBuffer.PushContext("innerOpeningTag"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for innerOpeningTag")
		}
		_innerOpeningTagErr := writeBuffer.WriteSerializable(ctx, m.GetInnerOpeningTag())
		if popErr := writeBuffer.PopContext("innerOpeningTag"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for innerOpeningTag")
		}
		if _innerOpeningTagErr != nil {
			return errors.Wrap(_innerOpeningTagErr, "Error serializing 'innerOpeningTag' field")
		}

		// Simple Field (presentValue)
		if pushErr := writeBuffer.PushContext("presentValue"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for presentValue")
		}
		_presentValueErr := writeBuffer.WriteSerializable(ctx, m.GetPresentValue())
		if popErr := writeBuffer.PopContext("presentValue"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for presentValue")
		}
		if _presentValueErr != nil {
			return errors.Wrap(_presentValueErr, "Error serializing 'presentValue' field")
		}

		// Simple Field (referencedFlags)
		if pushErr := writeBuffer.PushContext("referencedFlags"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for referencedFlags")
		}
		_referencedFlagsErr := writeBuffer.WriteSerializable(ctx, m.GetReferencedFlags())
		if popErr := writeBuffer.PopContext("referencedFlags"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for referencedFlags")
		}
		if _referencedFlagsErr != nil {
			return errors.Wrap(_referencedFlagsErr, "Error serializing 'referencedFlags' field")
		}

		// Simple Field (innerClosingTag)
		if pushErr := writeBuffer.PushContext("innerClosingTag"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for innerClosingTag")
		}
		_innerClosingTagErr := writeBuffer.WriteSerializable(ctx, m.GetInnerClosingTag())
		if popErr := writeBuffer.PopContext("innerClosingTag"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for innerClosingTag")
		}
		if _innerClosingTagErr != nil {
			return errors.Wrap(_innerClosingTagErr, "Error serializing 'innerClosingTag' field")
		}

		if popErr := writeBuffer.PopContext("BACnetNotificationParametersChangeOfStatusFlags"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetNotificationParametersChangeOfStatusFlags")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) isBACnetNotificationParametersChangeOfStatusFlags() bool {
	return true
}

func (m *_BACnetNotificationParametersChangeOfStatusFlags) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
