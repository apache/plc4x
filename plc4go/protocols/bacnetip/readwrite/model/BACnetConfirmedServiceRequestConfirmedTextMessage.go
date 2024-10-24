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

// BACnetConfirmedServiceRequestConfirmedTextMessage is the corresponding interface of BACnetConfirmedServiceRequestConfirmedTextMessage
type BACnetConfirmedServiceRequestConfirmedTextMessage interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetConfirmedServiceRequest
	// GetTextMessageSourceDevice returns TextMessageSourceDevice (property field)
	GetTextMessageSourceDevice() BACnetContextTagObjectIdentifier
	// GetMessageClass returns MessageClass (property field)
	GetMessageClass() BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass
	// GetMessagePriority returns MessagePriority (property field)
	GetMessagePriority() BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged
	// GetMessage returns Message (property field)
	GetMessage() BACnetContextTagCharacterString
	// IsBACnetConfirmedServiceRequestConfirmedTextMessage is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetConfirmedServiceRequestConfirmedTextMessage()
	// CreateBuilder creates a BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
	CreateBACnetConfirmedServiceRequestConfirmedTextMessageBuilder() BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
}

// _BACnetConfirmedServiceRequestConfirmedTextMessage is the data-structure of this message
type _BACnetConfirmedServiceRequestConfirmedTextMessage struct {
	BACnetConfirmedServiceRequestContract
	TextMessageSourceDevice BACnetContextTagObjectIdentifier
	MessageClass            BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass
	MessagePriority         BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged
	Message                 BACnetContextTagCharacterString
}

var _ BACnetConfirmedServiceRequestConfirmedTextMessage = (*_BACnetConfirmedServiceRequestConfirmedTextMessage)(nil)
var _ BACnetConfirmedServiceRequestRequirements = (*_BACnetConfirmedServiceRequestConfirmedTextMessage)(nil)

// NewBACnetConfirmedServiceRequestConfirmedTextMessage factory function for _BACnetConfirmedServiceRequestConfirmedTextMessage
func NewBACnetConfirmedServiceRequestConfirmedTextMessage(textMessageSourceDevice BACnetContextTagObjectIdentifier, messageClass BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass, messagePriority BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged, message BACnetContextTagCharacterString, serviceRequestLength uint32) *_BACnetConfirmedServiceRequestConfirmedTextMessage {
	if textMessageSourceDevice == nil {
		panic("textMessageSourceDevice of type BACnetContextTagObjectIdentifier for BACnetConfirmedServiceRequestConfirmedTextMessage must not be nil")
	}
	if messagePriority == nil {
		panic("messagePriority of type BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged for BACnetConfirmedServiceRequestConfirmedTextMessage must not be nil")
	}
	if message == nil {
		panic("message of type BACnetContextTagCharacterString for BACnetConfirmedServiceRequestConfirmedTextMessage must not be nil")
	}
	_result := &_BACnetConfirmedServiceRequestConfirmedTextMessage{
		BACnetConfirmedServiceRequestContract: NewBACnetConfirmedServiceRequest(serviceRequestLength),
		TextMessageSourceDevice:               textMessageSourceDevice,
		MessageClass:                          messageClass,
		MessagePriority:                       messagePriority,
		Message:                               message,
	}
	_result.BACnetConfirmedServiceRequestContract.(*_BACnetConfirmedServiceRequest)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetConfirmedServiceRequestConfirmedTextMessageBuilder is a builder for BACnetConfirmedServiceRequestConfirmedTextMessage
type BACnetConfirmedServiceRequestConfirmedTextMessageBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(textMessageSourceDevice BACnetContextTagObjectIdentifier, messagePriority BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged, message BACnetContextTagCharacterString) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
	// WithTextMessageSourceDevice adds TextMessageSourceDevice (property field)
	WithTextMessageSourceDevice(BACnetContextTagObjectIdentifier) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
	// WithTextMessageSourceDeviceBuilder adds TextMessageSourceDevice (property field) which is build by the builder
	WithTextMessageSourceDeviceBuilder(func(BACnetContextTagObjectIdentifierBuilder) BACnetContextTagObjectIdentifierBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
	// WithMessageClass adds MessageClass (property field)
	WithOptionalMessageClass(BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
	// WithOptionalMessageClassBuilder adds MessageClass (property field) which is build by the builder
	WithOptionalMessageClassBuilder(func(BACnetConfirmedServiceRequestConfirmedTextMessageMessageClassBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageMessageClassBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
	// WithMessagePriority adds MessagePriority (property field)
	WithMessagePriority(BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
	// WithMessagePriorityBuilder adds MessagePriority (property field) which is build by the builder
	WithMessagePriorityBuilder(func(BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTaggedBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTaggedBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
	// WithMessage adds Message (property field)
	WithMessage(BACnetContextTagCharacterString) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
	// WithMessageBuilder adds Message (property field) which is build by the builder
	WithMessageBuilder(func(BACnetContextTagCharacterStringBuilder) BACnetContextTagCharacterStringBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
	// Build builds the BACnetConfirmedServiceRequestConfirmedTextMessage or returns an error if something is wrong
	Build() (BACnetConfirmedServiceRequestConfirmedTextMessage, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetConfirmedServiceRequestConfirmedTextMessage
}

// NewBACnetConfirmedServiceRequestConfirmedTextMessageBuilder() creates a BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
func NewBACnetConfirmedServiceRequestConfirmedTextMessageBuilder() BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	return &_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder{_BACnetConfirmedServiceRequestConfirmedTextMessage: new(_BACnetConfirmedServiceRequestConfirmedTextMessage)}
}

type _BACnetConfirmedServiceRequestConfirmedTextMessageBuilder struct {
	*_BACnetConfirmedServiceRequestConfirmedTextMessage

	parentBuilder *_BACnetConfirmedServiceRequestBuilder

	err *utils.MultiError
}

var _ (BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) = (*_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder)(nil)

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) setParent(contract BACnetConfirmedServiceRequestContract) {
	b.BACnetConfirmedServiceRequestContract = contract
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) WithMandatoryFields(textMessageSourceDevice BACnetContextTagObjectIdentifier, messagePriority BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged, message BACnetContextTagCharacterString) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	return b.WithTextMessageSourceDevice(textMessageSourceDevice).WithMessagePriority(messagePriority).WithMessage(message)
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) WithTextMessageSourceDevice(textMessageSourceDevice BACnetContextTagObjectIdentifier) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	b.TextMessageSourceDevice = textMessageSourceDevice
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) WithTextMessageSourceDeviceBuilder(builderSupplier func(BACnetContextTagObjectIdentifierBuilder) BACnetContextTagObjectIdentifierBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	builder := builderSupplier(b.TextMessageSourceDevice.CreateBACnetContextTagObjectIdentifierBuilder())
	var err error
	b.TextMessageSourceDevice, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagObjectIdentifierBuilder failed"))
	}
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) WithOptionalMessageClass(messageClass BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	b.MessageClass = messageClass
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) WithOptionalMessageClassBuilder(builderSupplier func(BACnetConfirmedServiceRequestConfirmedTextMessageMessageClassBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageMessageClassBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	builder := builderSupplier(b.MessageClass.CreateBACnetConfirmedServiceRequestConfirmedTextMessageMessageClassBuilder())
	var err error
	b.MessageClass, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetConfirmedServiceRequestConfirmedTextMessageMessageClassBuilder failed"))
	}
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) WithMessagePriority(messagePriority BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	b.MessagePriority = messagePriority
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) WithMessagePriorityBuilder(builderSupplier func(BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTaggedBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTaggedBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	builder := builderSupplier(b.MessagePriority.CreateBACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTaggedBuilder())
	var err error
	b.MessagePriority, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) WithMessage(message BACnetContextTagCharacterString) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	b.Message = message
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) WithMessageBuilder(builderSupplier func(BACnetContextTagCharacterStringBuilder) BACnetContextTagCharacterStringBuilder) BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	builder := builderSupplier(b.Message.CreateBACnetContextTagCharacterStringBuilder())
	var err error
	b.Message, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetContextTagCharacterStringBuilder failed"))
	}
	return b
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) Build() (BACnetConfirmedServiceRequestConfirmedTextMessage, error) {
	if b.TextMessageSourceDevice == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'textMessageSourceDevice' not set"))
	}
	if b.MessagePriority == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'messagePriority' not set"))
	}
	if b.Message == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'message' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetConfirmedServiceRequestConfirmedTextMessage.deepCopy(), nil
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) MustBuild() BACnetConfirmedServiceRequestConfirmedTextMessage {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) Done() BACnetConfirmedServiceRequestBuilder {
	return b.parentBuilder
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) buildForBACnetConfirmedServiceRequest() (BACnetConfirmedServiceRequest, error) {
	return b.Build()
}

func (b *_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder) DeepCopy() any {
	_copy := b.CreateBACnetConfirmedServiceRequestConfirmedTextMessageBuilder().(*_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetConfirmedServiceRequestConfirmedTextMessageBuilder creates a BACnetConfirmedServiceRequestConfirmedTextMessageBuilder
func (b *_BACnetConfirmedServiceRequestConfirmedTextMessage) CreateBACnetConfirmedServiceRequestConfirmedTextMessageBuilder() BACnetConfirmedServiceRequestConfirmedTextMessageBuilder {
	if b == nil {
		return NewBACnetConfirmedServiceRequestConfirmedTextMessageBuilder()
	}
	return &_BACnetConfirmedServiceRequestConfirmedTextMessageBuilder{_BACnetConfirmedServiceRequestConfirmedTextMessage: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) GetServiceChoice() BACnetConfirmedServiceChoice {
	return BACnetConfirmedServiceChoice_CONFIRMED_TEXT_MESSAGE
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) GetParent() BACnetConfirmedServiceRequestContract {
	return m.BACnetConfirmedServiceRequestContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) GetTextMessageSourceDevice() BACnetContextTagObjectIdentifier {
	return m.TextMessageSourceDevice
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) GetMessageClass() BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass {
	return m.MessageClass
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) GetMessagePriority() BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged {
	return m.MessagePriority
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) GetMessage() BACnetContextTagCharacterString {
	return m.Message
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetConfirmedServiceRequestConfirmedTextMessage(structType any) BACnetConfirmedServiceRequestConfirmedTextMessage {
	if casted, ok := structType.(BACnetConfirmedServiceRequestConfirmedTextMessage); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetConfirmedServiceRequestConfirmedTextMessage); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) GetTypeName() string {
	return "BACnetConfirmedServiceRequestConfirmedTextMessage"
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetConfirmedServiceRequestContract.(*_BACnetConfirmedServiceRequest).GetLengthInBits(ctx))

	// Simple field (textMessageSourceDevice)
	lengthInBits += m.TextMessageSourceDevice.GetLengthInBits(ctx)

	// Optional Field (messageClass)
	if m.MessageClass != nil {
		lengthInBits += m.MessageClass.GetLengthInBits(ctx)
	}

	// Simple field (messagePriority)
	lengthInBits += m.MessagePriority.GetLengthInBits(ctx)

	// Simple field (message)
	lengthInBits += m.Message.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetConfirmedServiceRequest, serviceRequestLength uint32) (__bACnetConfirmedServiceRequestConfirmedTextMessage BACnetConfirmedServiceRequestConfirmedTextMessage, err error) {
	m.BACnetConfirmedServiceRequestContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetConfirmedServiceRequestConfirmedTextMessage"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetConfirmedServiceRequestConfirmedTextMessage")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	textMessageSourceDevice, err := ReadSimpleField[BACnetContextTagObjectIdentifier](ctx, "textMessageSourceDevice", ReadComplex[BACnetContextTagObjectIdentifier](BACnetContextTagParseWithBufferProducer[BACnetContextTagObjectIdentifier]((uint8)(uint8(0)), (BACnetDataType)(BACnetDataType_BACNET_OBJECT_IDENTIFIER)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'textMessageSourceDevice' field"))
	}
	m.TextMessageSourceDevice = textMessageSourceDevice

	var messageClass BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass
	_messageClass, err := ReadOptionalField[BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass](ctx, "messageClass", ReadComplex[BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass](BACnetConfirmedServiceRequestConfirmedTextMessageMessageClassParseWithBufferProducer[BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass]((uint8)(uint8(1))), readBuffer), true)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'messageClass' field"))
	}
	if _messageClass != nil {
		messageClass = *_messageClass
		m.MessageClass = messageClass
	}

	messagePriority, err := ReadSimpleField[BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged](ctx, "messagePriority", ReadComplex[BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged](BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTaggedParseWithBufferProducer((uint8)(uint8(2)), (TagClass)(TagClass_CONTEXT_SPECIFIC_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'messagePriority' field"))
	}
	m.MessagePriority = messagePriority

	message, err := ReadSimpleField[BACnetContextTagCharacterString](ctx, "message", ReadComplex[BACnetContextTagCharacterString](BACnetContextTagParseWithBufferProducer[BACnetContextTagCharacterString]((uint8)(uint8(3)), (BACnetDataType)(BACnetDataType_CHARACTER_STRING)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'message' field"))
	}
	m.Message = message

	if closeErr := readBuffer.CloseContext("BACnetConfirmedServiceRequestConfirmedTextMessage"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetConfirmedServiceRequestConfirmedTextMessage")
	}

	return m, nil
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetConfirmedServiceRequestConfirmedTextMessage"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetConfirmedServiceRequestConfirmedTextMessage")
		}

		if err := WriteSimpleField[BACnetContextTagObjectIdentifier](ctx, "textMessageSourceDevice", m.GetTextMessageSourceDevice(), WriteComplex[BACnetContextTagObjectIdentifier](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'textMessageSourceDevice' field")
		}

		if err := WriteOptionalField[BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass](ctx, "messageClass", GetRef(m.GetMessageClass()), WriteComplex[BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass](writeBuffer), true); err != nil {
			return errors.Wrap(err, "Error serializing 'messageClass' field")
		}

		if err := WriteSimpleField[BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged](ctx, "messagePriority", m.GetMessagePriority(), WriteComplex[BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'messagePriority' field")
		}

		if err := WriteSimpleField[BACnetContextTagCharacterString](ctx, "message", m.GetMessage(), WriteComplex[BACnetContextTagCharacterString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'message' field")
		}

		if popErr := writeBuffer.PopContext("BACnetConfirmedServiceRequestConfirmedTextMessage"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetConfirmedServiceRequestConfirmedTextMessage")
		}
		return nil
	}
	return m.BACnetConfirmedServiceRequestContract.(*_BACnetConfirmedServiceRequest).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) IsBACnetConfirmedServiceRequestConfirmedTextMessage() {
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) deepCopy() *_BACnetConfirmedServiceRequestConfirmedTextMessage {
	if m == nil {
		return nil
	}
	_BACnetConfirmedServiceRequestConfirmedTextMessageCopy := &_BACnetConfirmedServiceRequestConfirmedTextMessage{
		m.BACnetConfirmedServiceRequestContract.(*_BACnetConfirmedServiceRequest).deepCopy(),
		m.TextMessageSourceDevice.DeepCopy().(BACnetContextTagObjectIdentifier),
		m.MessageClass.DeepCopy().(BACnetConfirmedServiceRequestConfirmedTextMessageMessageClass),
		m.MessagePriority.DeepCopy().(BACnetConfirmedServiceRequestConfirmedTextMessageMessagePriorityTagged),
		m.Message.DeepCopy().(BACnetContextTagCharacterString),
	}
	m.BACnetConfirmedServiceRequestContract.(*_BACnetConfirmedServiceRequest)._SubType = m
	return _BACnetConfirmedServiceRequestConfirmedTextMessageCopy
}

func (m *_BACnetConfirmedServiceRequestConfirmedTextMessage) String() string {
	if m == nil {
		return "<nil>"
	}
	wb := utils.NewWriteBufferBoxBased(
		utils.WithWriteBufferBoxBasedMergeSingleBoxes(),
		utils.WithWriteBufferBoxBasedOmitEmptyBoxes(),
		utils.WithWriteBufferBoxBasedPrintPosLengthFooter(),
	)
	if err := wb.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return wb.GetBox().String()
}
