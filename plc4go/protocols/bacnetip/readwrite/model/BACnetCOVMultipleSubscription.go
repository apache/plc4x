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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetCOVMultipleSubscription is the corresponding interface of BACnetCOVMultipleSubscription
type BACnetCOVMultipleSubscription interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetRecipient returns Recipient (property field)
	GetRecipient() BACnetRecipientProcessEnclosed
	// GetIssueConfirmedNotifications returns IssueConfirmedNotifications (property field)
	GetIssueConfirmedNotifications() BACnetContextTagBoolean
	// GetTimeRemaining returns TimeRemaining (property field)
	GetTimeRemaining() BACnetContextTagUnsignedInteger
	// GetMaxNotificationDelay returns MaxNotificationDelay (property field)
	GetMaxNotificationDelay() BACnetContextTagUnsignedInteger
	// GetListOfCovSubscriptionSpecification returns ListOfCovSubscriptionSpecification (property field)
	GetListOfCovSubscriptionSpecification() BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification
	// IsBACnetCOVMultipleSubscription is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetCOVMultipleSubscription()
}

// _BACnetCOVMultipleSubscription is the data-structure of this message
type _BACnetCOVMultipleSubscription struct {
	Recipient                          BACnetRecipientProcessEnclosed
	IssueConfirmedNotifications        BACnetContextTagBoolean
	TimeRemaining                      BACnetContextTagUnsignedInteger
	MaxNotificationDelay               BACnetContextTagUnsignedInteger
	ListOfCovSubscriptionSpecification BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification
}

var _ BACnetCOVMultipleSubscription = (*_BACnetCOVMultipleSubscription)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetCOVMultipleSubscription) GetRecipient() BACnetRecipientProcessEnclosed {
	return m.Recipient
}

func (m *_BACnetCOVMultipleSubscription) GetIssueConfirmedNotifications() BACnetContextTagBoolean {
	return m.IssueConfirmedNotifications
}

func (m *_BACnetCOVMultipleSubscription) GetTimeRemaining() BACnetContextTagUnsignedInteger {
	return m.TimeRemaining
}

func (m *_BACnetCOVMultipleSubscription) GetMaxNotificationDelay() BACnetContextTagUnsignedInteger {
	return m.MaxNotificationDelay
}

func (m *_BACnetCOVMultipleSubscription) GetListOfCovSubscriptionSpecification() BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification {
	return m.ListOfCovSubscriptionSpecification
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetCOVMultipleSubscription factory function for _BACnetCOVMultipleSubscription
func NewBACnetCOVMultipleSubscription(recipient BACnetRecipientProcessEnclosed, issueConfirmedNotifications BACnetContextTagBoolean, timeRemaining BACnetContextTagUnsignedInteger, maxNotificationDelay BACnetContextTagUnsignedInteger, listOfCovSubscriptionSpecification BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification) *_BACnetCOVMultipleSubscription {
	if recipient == nil {
		panic("recipient of type BACnetRecipientProcessEnclosed for BACnetCOVMultipleSubscription must not be nil")
	}
	if issueConfirmedNotifications == nil {
		panic("issueConfirmedNotifications of type BACnetContextTagBoolean for BACnetCOVMultipleSubscription must not be nil")
	}
	if timeRemaining == nil {
		panic("timeRemaining of type BACnetContextTagUnsignedInteger for BACnetCOVMultipleSubscription must not be nil")
	}
	if maxNotificationDelay == nil {
		panic("maxNotificationDelay of type BACnetContextTagUnsignedInteger for BACnetCOVMultipleSubscription must not be nil")
	}
	if listOfCovSubscriptionSpecification == nil {
		panic("listOfCovSubscriptionSpecification of type BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification for BACnetCOVMultipleSubscription must not be nil")
	}
	return &_BACnetCOVMultipleSubscription{Recipient: recipient, IssueConfirmedNotifications: issueConfirmedNotifications, TimeRemaining: timeRemaining, MaxNotificationDelay: maxNotificationDelay, ListOfCovSubscriptionSpecification: listOfCovSubscriptionSpecification}
}

// Deprecated: use the interface for direct cast
func CastBACnetCOVMultipleSubscription(structType any) BACnetCOVMultipleSubscription {
	if casted, ok := structType.(BACnetCOVMultipleSubscription); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetCOVMultipleSubscription); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetCOVMultipleSubscription) GetTypeName() string {
	return "BACnetCOVMultipleSubscription"
}

func (m *_BACnetCOVMultipleSubscription) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)

	// Simple field (recipient)
	lengthInBits += m.Recipient.GetLengthInBits(ctx)

	// Simple field (issueConfirmedNotifications)
	lengthInBits += m.IssueConfirmedNotifications.GetLengthInBits(ctx)

	// Simple field (timeRemaining)
	lengthInBits += m.TimeRemaining.GetLengthInBits(ctx)

	// Simple field (maxNotificationDelay)
	lengthInBits += m.MaxNotificationDelay.GetLengthInBits(ctx)

	// Simple field (listOfCovSubscriptionSpecification)
	lengthInBits += m.ListOfCovSubscriptionSpecification.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetCOVMultipleSubscription) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func BACnetCOVMultipleSubscriptionParse(ctx context.Context, theBytes []byte) (BACnetCOVMultipleSubscription, error) {
	return BACnetCOVMultipleSubscriptionParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes))
}

func BACnetCOVMultipleSubscriptionParseWithBufferProducer() func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetCOVMultipleSubscription, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetCOVMultipleSubscription, error) {
		return BACnetCOVMultipleSubscriptionParseWithBuffer(ctx, readBuffer)
	}
}

func BACnetCOVMultipleSubscriptionParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (BACnetCOVMultipleSubscription, error) {
	v, err := (&_BACnetCOVMultipleSubscription{}).parse(ctx, readBuffer)
	if err != nil {
		return nil, err
	}
	return v, err
}

func (m *_BACnetCOVMultipleSubscription) parse(ctx context.Context, readBuffer utils.ReadBuffer) (__bACnetCOVMultipleSubscription BACnetCOVMultipleSubscription, err error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetCOVMultipleSubscription"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetCOVMultipleSubscription")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	recipient, err := ReadSimpleField[BACnetRecipientProcessEnclosed](ctx, "recipient", ReadComplex[BACnetRecipientProcessEnclosed](BACnetRecipientProcessEnclosedParseWithBufferProducer((uint8)(uint8(0))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'recipient' field"))
	}
	m.Recipient = recipient

	issueConfirmedNotifications, err := ReadSimpleField[BACnetContextTagBoolean](ctx, "issueConfirmedNotifications", ReadComplex[BACnetContextTagBoolean](BACnetContextTagParseWithBufferProducer[BACnetContextTagBoolean]((uint8)(uint8(1)), (BACnetDataType)(BACnetDataType_BOOLEAN)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'issueConfirmedNotifications' field"))
	}
	m.IssueConfirmedNotifications = issueConfirmedNotifications

	timeRemaining, err := ReadSimpleField[BACnetContextTagUnsignedInteger](ctx, "timeRemaining", ReadComplex[BACnetContextTagUnsignedInteger](BACnetContextTagParseWithBufferProducer[BACnetContextTagUnsignedInteger]((uint8)(uint8(2)), (BACnetDataType)(BACnetDataType_UNSIGNED_INTEGER)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'timeRemaining' field"))
	}
	m.TimeRemaining = timeRemaining

	maxNotificationDelay, err := ReadSimpleField[BACnetContextTagUnsignedInteger](ctx, "maxNotificationDelay", ReadComplex[BACnetContextTagUnsignedInteger](BACnetContextTagParseWithBufferProducer[BACnetContextTagUnsignedInteger]((uint8)(uint8(3)), (BACnetDataType)(BACnetDataType_UNSIGNED_INTEGER)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'maxNotificationDelay' field"))
	}
	m.MaxNotificationDelay = maxNotificationDelay

	listOfCovSubscriptionSpecification, err := ReadSimpleField[BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification](ctx, "listOfCovSubscriptionSpecification", ReadComplex[BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification](BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecificationParseWithBufferProducer((uint8)(uint8(4))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'listOfCovSubscriptionSpecification' field"))
	}
	m.ListOfCovSubscriptionSpecification = listOfCovSubscriptionSpecification

	if closeErr := readBuffer.CloseContext("BACnetCOVMultipleSubscription"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetCOVMultipleSubscription")
	}

	return m, nil
}

func (m *_BACnetCOVMultipleSubscription) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetCOVMultipleSubscription) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pushErr := writeBuffer.PushContext("BACnetCOVMultipleSubscription"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetCOVMultipleSubscription")
	}

	if err := WriteSimpleField[BACnetRecipientProcessEnclosed](ctx, "recipient", m.GetRecipient(), WriteComplex[BACnetRecipientProcessEnclosed](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'recipient' field")
	}

	if err := WriteSimpleField[BACnetContextTagBoolean](ctx, "issueConfirmedNotifications", m.GetIssueConfirmedNotifications(), WriteComplex[BACnetContextTagBoolean](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'issueConfirmedNotifications' field")
	}

	if err := WriteSimpleField[BACnetContextTagUnsignedInteger](ctx, "timeRemaining", m.GetTimeRemaining(), WriteComplex[BACnetContextTagUnsignedInteger](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'timeRemaining' field")
	}

	if err := WriteSimpleField[BACnetContextTagUnsignedInteger](ctx, "maxNotificationDelay", m.GetMaxNotificationDelay(), WriteComplex[BACnetContextTagUnsignedInteger](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'maxNotificationDelay' field")
	}

	if err := WriteSimpleField[BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification](ctx, "listOfCovSubscriptionSpecification", m.GetListOfCovSubscriptionSpecification(), WriteComplex[BACnetCOVMultipleSubscriptionListOfCovSubscriptionSpecification](writeBuffer)); err != nil {
		return errors.Wrap(err, "Error serializing 'listOfCovSubscriptionSpecification' field")
	}

	if popErr := writeBuffer.PopContext("BACnetCOVMultipleSubscription"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetCOVMultipleSubscription")
	}
	return nil
}

func (m *_BACnetCOVMultipleSubscription) IsBACnetCOVMultipleSubscription() {}

func (m *_BACnetCOVMultipleSubscription) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
