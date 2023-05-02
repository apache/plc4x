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

// BACnetEventParameterCommandFailure is the corresponding interface of BACnetEventParameterCommandFailure
type BACnetEventParameterCommandFailure interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	BACnetEventParameter
	// GetOpeningTag returns OpeningTag (property field)
	GetOpeningTag() BACnetOpeningTag
	// GetTimeDelay returns TimeDelay (property field)
	GetTimeDelay() BACnetContextTagUnsignedInteger
	// GetFeedbackPropertyReference returns FeedbackPropertyReference (property field)
	GetFeedbackPropertyReference() BACnetDeviceObjectPropertyReferenceEnclosed
	// GetClosingTag returns ClosingTag (property field)
	GetClosingTag() BACnetClosingTag
}

// BACnetEventParameterCommandFailureExactly can be used when we want exactly this type and not a type which fulfills BACnetEventParameterCommandFailure.
// This is useful for switch cases.
type BACnetEventParameterCommandFailureExactly interface {
	BACnetEventParameterCommandFailure
	isBACnetEventParameterCommandFailure() bool
}

// _BACnetEventParameterCommandFailure is the data-structure of this message
type _BACnetEventParameterCommandFailure struct {
	*_BACnetEventParameter
	OpeningTag                BACnetOpeningTag
	TimeDelay                 BACnetContextTagUnsignedInteger
	FeedbackPropertyReference BACnetDeviceObjectPropertyReferenceEnclosed
	ClosingTag                BACnetClosingTag
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetEventParameterCommandFailure) InitializeParent(parent BACnetEventParameter, peekedTagHeader BACnetTagHeader) {
	m.PeekedTagHeader = peekedTagHeader
}

func (m *_BACnetEventParameterCommandFailure) GetParent() BACnetEventParameter {
	return m._BACnetEventParameter
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetEventParameterCommandFailure) GetOpeningTag() BACnetOpeningTag {
	return m.OpeningTag
}

func (m *_BACnetEventParameterCommandFailure) GetTimeDelay() BACnetContextTagUnsignedInteger {
	return m.TimeDelay
}

func (m *_BACnetEventParameterCommandFailure) GetFeedbackPropertyReference() BACnetDeviceObjectPropertyReferenceEnclosed {
	return m.FeedbackPropertyReference
}

func (m *_BACnetEventParameterCommandFailure) GetClosingTag() BACnetClosingTag {
	return m.ClosingTag
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetEventParameterCommandFailure factory function for _BACnetEventParameterCommandFailure
func NewBACnetEventParameterCommandFailure(openingTag BACnetOpeningTag, timeDelay BACnetContextTagUnsignedInteger, feedbackPropertyReference BACnetDeviceObjectPropertyReferenceEnclosed, closingTag BACnetClosingTag, peekedTagHeader BACnetTagHeader) *_BACnetEventParameterCommandFailure {
	_result := &_BACnetEventParameterCommandFailure{
		OpeningTag:                openingTag,
		TimeDelay:                 timeDelay,
		FeedbackPropertyReference: feedbackPropertyReference,
		ClosingTag:                closingTag,
		_BACnetEventParameter:     NewBACnetEventParameter(peekedTagHeader),
	}
	_result._BACnetEventParameter._BACnetEventParameterChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetEventParameterCommandFailure(structType any) BACnetEventParameterCommandFailure {
	if casted, ok := structType.(BACnetEventParameterCommandFailure); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetEventParameterCommandFailure); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetEventParameterCommandFailure) GetTypeName() string {
	return "BACnetEventParameterCommandFailure"
}

func (m *_BACnetEventParameterCommandFailure) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (openingTag)
	lengthInBits += m.OpeningTag.GetLengthInBits(ctx)

	// Simple field (timeDelay)
	lengthInBits += m.TimeDelay.GetLengthInBits(ctx)

	// Simple field (feedbackPropertyReference)
	lengthInBits += m.FeedbackPropertyReference.GetLengthInBits(ctx)

	// Simple field (closingTag)
	lengthInBits += m.ClosingTag.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetEventParameterCommandFailure) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetEventParameterCommandFailureParse(theBytes []byte) (BACnetEventParameterCommandFailure, error) {
	return BACnetEventParameterCommandFailureParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes))
}

func BACnetEventParameterCommandFailureParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetEventParameterCommandFailure, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetEventParameterCommandFailure"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetEventParameterCommandFailure")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (openingTag)
	if pullErr := readBuffer.PullContext("openingTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for openingTag")
	}
	_openingTag, _openingTagErr := BACnetOpeningTagParseWithBuffer(ctx, readBuffer, uint8(uint8(3)))
	if _openingTagErr != nil {
		return nil, errors.Wrap(_openingTagErr, "Error parsing 'openingTag' field of BACnetEventParameterCommandFailure")
	}
	openingTag := _openingTag.(BACnetOpeningTag)
	if closeErr := readBuffer.CloseContext("openingTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for openingTag")
	}

	// Simple Field (timeDelay)
	if pullErr := readBuffer.PullContext("timeDelay"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for timeDelay")
	}
	_timeDelay, _timeDelayErr := BACnetContextTagParseWithBuffer(ctx, readBuffer, uint8(uint8(0)), BACnetDataType(BACnetDataType_UNSIGNED_INTEGER))
	if _timeDelayErr != nil {
		return nil, errors.Wrap(_timeDelayErr, "Error parsing 'timeDelay' field of BACnetEventParameterCommandFailure")
	}
	timeDelay := _timeDelay.(BACnetContextTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("timeDelay"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for timeDelay")
	}

	// Simple Field (feedbackPropertyReference)
	if pullErr := readBuffer.PullContext("feedbackPropertyReference"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for feedbackPropertyReference")
	}
	_feedbackPropertyReference, _feedbackPropertyReferenceErr := BACnetDeviceObjectPropertyReferenceEnclosedParseWithBuffer(ctx, readBuffer, uint8(uint8(1)))
	if _feedbackPropertyReferenceErr != nil {
		return nil, errors.Wrap(_feedbackPropertyReferenceErr, "Error parsing 'feedbackPropertyReference' field of BACnetEventParameterCommandFailure")
	}
	feedbackPropertyReference := _feedbackPropertyReference.(BACnetDeviceObjectPropertyReferenceEnclosed)
	if closeErr := readBuffer.CloseContext("feedbackPropertyReference"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for feedbackPropertyReference")
	}

	// Simple Field (closingTag)
	if pullErr := readBuffer.PullContext("closingTag"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for closingTag")
	}
	_closingTag, _closingTagErr := BACnetClosingTagParseWithBuffer(ctx, readBuffer, uint8(uint8(3)))
	if _closingTagErr != nil {
		return nil, errors.Wrap(_closingTagErr, "Error parsing 'closingTag' field of BACnetEventParameterCommandFailure")
	}
	closingTag := _closingTag.(BACnetClosingTag)
	if closeErr := readBuffer.CloseContext("closingTag"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for closingTag")
	}

	if closeErr := readBuffer.CloseContext("BACnetEventParameterCommandFailure"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetEventParameterCommandFailure")
	}

	// Create a partially initialized instance
	_child := &_BACnetEventParameterCommandFailure{
		_BACnetEventParameter:     &_BACnetEventParameter{},
		OpeningTag:                openingTag,
		TimeDelay:                 timeDelay,
		FeedbackPropertyReference: feedbackPropertyReference,
		ClosingTag:                closingTag,
	}
	_child._BACnetEventParameter._BACnetEventParameterChildRequirements = _child
	return _child, nil
}

func (m *_BACnetEventParameterCommandFailure) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetEventParameterCommandFailure) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetEventParameterCommandFailure"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetEventParameterCommandFailure")
		}

		// Simple Field (openingTag)
		if pushErr := writeBuffer.PushContext("openingTag"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for openingTag")
		}
		_openingTagErr := writeBuffer.WriteSerializable(ctx, m.GetOpeningTag())
		if popErr := writeBuffer.PopContext("openingTag"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for openingTag")
		}
		if _openingTagErr != nil {
			return errors.Wrap(_openingTagErr, "Error serializing 'openingTag' field")
		}

		// Simple Field (timeDelay)
		if pushErr := writeBuffer.PushContext("timeDelay"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for timeDelay")
		}
		_timeDelayErr := writeBuffer.WriteSerializable(ctx, m.GetTimeDelay())
		if popErr := writeBuffer.PopContext("timeDelay"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for timeDelay")
		}
		if _timeDelayErr != nil {
			return errors.Wrap(_timeDelayErr, "Error serializing 'timeDelay' field")
		}

		// Simple Field (feedbackPropertyReference)
		if pushErr := writeBuffer.PushContext("feedbackPropertyReference"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for feedbackPropertyReference")
		}
		_feedbackPropertyReferenceErr := writeBuffer.WriteSerializable(ctx, m.GetFeedbackPropertyReference())
		if popErr := writeBuffer.PopContext("feedbackPropertyReference"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for feedbackPropertyReference")
		}
		if _feedbackPropertyReferenceErr != nil {
			return errors.Wrap(_feedbackPropertyReferenceErr, "Error serializing 'feedbackPropertyReference' field")
		}

		// Simple Field (closingTag)
		if pushErr := writeBuffer.PushContext("closingTag"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for closingTag")
		}
		_closingTagErr := writeBuffer.WriteSerializable(ctx, m.GetClosingTag())
		if popErr := writeBuffer.PopContext("closingTag"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for closingTag")
		}
		if _closingTagErr != nil {
			return errors.Wrap(_closingTagErr, "Error serializing 'closingTag' field")
		}

		if popErr := writeBuffer.PopContext("BACnetEventParameterCommandFailure"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetEventParameterCommandFailure")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetEventParameterCommandFailure) isBACnetEventParameterCommandFailure() bool {
	return true
}

func (m *_BACnetEventParameterCommandFailure) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
