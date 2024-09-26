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

// BACnetServiceAckGetEnrollmentSummary is the corresponding interface of BACnetServiceAckGetEnrollmentSummary
type BACnetServiceAckGetEnrollmentSummary interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	BACnetServiceAck
	// GetObjectIdentifier returns ObjectIdentifier (property field)
	GetObjectIdentifier() BACnetApplicationTagObjectIdentifier
	// GetEventType returns EventType (property field)
	GetEventType() BACnetEventTypeTagged
	// GetEventState returns EventState (property field)
	GetEventState() BACnetEventStateTagged
	// GetPriority returns Priority (property field)
	GetPriority() BACnetApplicationTagUnsignedInteger
	// GetNotificationClass returns NotificationClass (property field)
	GetNotificationClass() BACnetApplicationTagUnsignedInteger
	// IsBACnetServiceAckGetEnrollmentSummary is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsBACnetServiceAckGetEnrollmentSummary()
	// CreateBuilder creates a BACnetServiceAckGetEnrollmentSummaryBuilder
	CreateBACnetServiceAckGetEnrollmentSummaryBuilder() BACnetServiceAckGetEnrollmentSummaryBuilder
}

// _BACnetServiceAckGetEnrollmentSummary is the data-structure of this message
type _BACnetServiceAckGetEnrollmentSummary struct {
	BACnetServiceAckContract
	ObjectIdentifier  BACnetApplicationTagObjectIdentifier
	EventType         BACnetEventTypeTagged
	EventState        BACnetEventStateTagged
	Priority          BACnetApplicationTagUnsignedInteger
	NotificationClass BACnetApplicationTagUnsignedInteger
}

var _ BACnetServiceAckGetEnrollmentSummary = (*_BACnetServiceAckGetEnrollmentSummary)(nil)
var _ BACnetServiceAckRequirements = (*_BACnetServiceAckGetEnrollmentSummary)(nil)

// NewBACnetServiceAckGetEnrollmentSummary factory function for _BACnetServiceAckGetEnrollmentSummary
func NewBACnetServiceAckGetEnrollmentSummary(objectIdentifier BACnetApplicationTagObjectIdentifier, eventType BACnetEventTypeTagged, eventState BACnetEventStateTagged, priority BACnetApplicationTagUnsignedInteger, notificationClass BACnetApplicationTagUnsignedInteger, serviceAckLength uint32) *_BACnetServiceAckGetEnrollmentSummary {
	if objectIdentifier == nil {
		panic("objectIdentifier of type BACnetApplicationTagObjectIdentifier for BACnetServiceAckGetEnrollmentSummary must not be nil")
	}
	if eventType == nil {
		panic("eventType of type BACnetEventTypeTagged for BACnetServiceAckGetEnrollmentSummary must not be nil")
	}
	if eventState == nil {
		panic("eventState of type BACnetEventStateTagged for BACnetServiceAckGetEnrollmentSummary must not be nil")
	}
	if priority == nil {
		panic("priority of type BACnetApplicationTagUnsignedInteger for BACnetServiceAckGetEnrollmentSummary must not be nil")
	}
	_result := &_BACnetServiceAckGetEnrollmentSummary{
		BACnetServiceAckContract: NewBACnetServiceAck(serviceAckLength),
		ObjectIdentifier:         objectIdentifier,
		EventType:                eventType,
		EventState:               eventState,
		Priority:                 priority,
		NotificationClass:        notificationClass,
	}
	_result.BACnetServiceAckContract.(*_BACnetServiceAck)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// BACnetServiceAckGetEnrollmentSummaryBuilder is a builder for BACnetServiceAckGetEnrollmentSummary
type BACnetServiceAckGetEnrollmentSummaryBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(objectIdentifier BACnetApplicationTagObjectIdentifier, eventType BACnetEventTypeTagged, eventState BACnetEventStateTagged, priority BACnetApplicationTagUnsignedInteger) BACnetServiceAckGetEnrollmentSummaryBuilder
	// WithObjectIdentifier adds ObjectIdentifier (property field)
	WithObjectIdentifier(BACnetApplicationTagObjectIdentifier) BACnetServiceAckGetEnrollmentSummaryBuilder
	// WithObjectIdentifierBuilder adds ObjectIdentifier (property field) which is build by the builder
	WithObjectIdentifierBuilder(func(BACnetApplicationTagObjectIdentifierBuilder) BACnetApplicationTagObjectIdentifierBuilder) BACnetServiceAckGetEnrollmentSummaryBuilder
	// WithEventType adds EventType (property field)
	WithEventType(BACnetEventTypeTagged) BACnetServiceAckGetEnrollmentSummaryBuilder
	// WithEventTypeBuilder adds EventType (property field) which is build by the builder
	WithEventTypeBuilder(func(BACnetEventTypeTaggedBuilder) BACnetEventTypeTaggedBuilder) BACnetServiceAckGetEnrollmentSummaryBuilder
	// WithEventState adds EventState (property field)
	WithEventState(BACnetEventStateTagged) BACnetServiceAckGetEnrollmentSummaryBuilder
	// WithEventStateBuilder adds EventState (property field) which is build by the builder
	WithEventStateBuilder(func(BACnetEventStateTaggedBuilder) BACnetEventStateTaggedBuilder) BACnetServiceAckGetEnrollmentSummaryBuilder
	// WithPriority adds Priority (property field)
	WithPriority(BACnetApplicationTagUnsignedInteger) BACnetServiceAckGetEnrollmentSummaryBuilder
	// WithPriorityBuilder adds Priority (property field) which is build by the builder
	WithPriorityBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetServiceAckGetEnrollmentSummaryBuilder
	// WithNotificationClass adds NotificationClass (property field)
	WithOptionalNotificationClass(BACnetApplicationTagUnsignedInteger) BACnetServiceAckGetEnrollmentSummaryBuilder
	// WithOptionalNotificationClassBuilder adds NotificationClass (property field) which is build by the builder
	WithOptionalNotificationClassBuilder(func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetServiceAckGetEnrollmentSummaryBuilder
	// Build builds the BACnetServiceAckGetEnrollmentSummary or returns an error if something is wrong
	Build() (BACnetServiceAckGetEnrollmentSummary, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() BACnetServiceAckGetEnrollmentSummary
}

// NewBACnetServiceAckGetEnrollmentSummaryBuilder() creates a BACnetServiceAckGetEnrollmentSummaryBuilder
func NewBACnetServiceAckGetEnrollmentSummaryBuilder() BACnetServiceAckGetEnrollmentSummaryBuilder {
	return &_BACnetServiceAckGetEnrollmentSummaryBuilder{_BACnetServiceAckGetEnrollmentSummary: new(_BACnetServiceAckGetEnrollmentSummary)}
}

type _BACnetServiceAckGetEnrollmentSummaryBuilder struct {
	*_BACnetServiceAckGetEnrollmentSummary

	parentBuilder *_BACnetServiceAckBuilder

	err *utils.MultiError
}

var _ (BACnetServiceAckGetEnrollmentSummaryBuilder) = (*_BACnetServiceAckGetEnrollmentSummaryBuilder)(nil)

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) setParent(contract BACnetServiceAckContract) {
	b.BACnetServiceAckContract = contract
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithMandatoryFields(objectIdentifier BACnetApplicationTagObjectIdentifier, eventType BACnetEventTypeTagged, eventState BACnetEventStateTagged, priority BACnetApplicationTagUnsignedInteger) BACnetServiceAckGetEnrollmentSummaryBuilder {
	return b.WithObjectIdentifier(objectIdentifier).WithEventType(eventType).WithEventState(eventState).WithPriority(priority)
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithObjectIdentifier(objectIdentifier BACnetApplicationTagObjectIdentifier) BACnetServiceAckGetEnrollmentSummaryBuilder {
	b.ObjectIdentifier = objectIdentifier
	return b
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithObjectIdentifierBuilder(builderSupplier func(BACnetApplicationTagObjectIdentifierBuilder) BACnetApplicationTagObjectIdentifierBuilder) BACnetServiceAckGetEnrollmentSummaryBuilder {
	builder := builderSupplier(b.ObjectIdentifier.CreateBACnetApplicationTagObjectIdentifierBuilder())
	var err error
	b.ObjectIdentifier, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagObjectIdentifierBuilder failed"))
	}
	return b
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithEventType(eventType BACnetEventTypeTagged) BACnetServiceAckGetEnrollmentSummaryBuilder {
	b.EventType = eventType
	return b
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithEventTypeBuilder(builderSupplier func(BACnetEventTypeTaggedBuilder) BACnetEventTypeTaggedBuilder) BACnetServiceAckGetEnrollmentSummaryBuilder {
	builder := builderSupplier(b.EventType.CreateBACnetEventTypeTaggedBuilder())
	var err error
	b.EventType, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetEventTypeTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithEventState(eventState BACnetEventStateTagged) BACnetServiceAckGetEnrollmentSummaryBuilder {
	b.EventState = eventState
	return b
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithEventStateBuilder(builderSupplier func(BACnetEventStateTaggedBuilder) BACnetEventStateTaggedBuilder) BACnetServiceAckGetEnrollmentSummaryBuilder {
	builder := builderSupplier(b.EventState.CreateBACnetEventStateTaggedBuilder())
	var err error
	b.EventState, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetEventStateTaggedBuilder failed"))
	}
	return b
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithPriority(priority BACnetApplicationTagUnsignedInteger) BACnetServiceAckGetEnrollmentSummaryBuilder {
	b.Priority = priority
	return b
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithPriorityBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetServiceAckGetEnrollmentSummaryBuilder {
	builder := builderSupplier(b.Priority.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.Priority, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithOptionalNotificationClass(notificationClass BACnetApplicationTagUnsignedInteger) BACnetServiceAckGetEnrollmentSummaryBuilder {
	b.NotificationClass = notificationClass
	return b
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) WithOptionalNotificationClassBuilder(builderSupplier func(BACnetApplicationTagUnsignedIntegerBuilder) BACnetApplicationTagUnsignedIntegerBuilder) BACnetServiceAckGetEnrollmentSummaryBuilder {
	builder := builderSupplier(b.NotificationClass.CreateBACnetApplicationTagUnsignedIntegerBuilder())
	var err error
	b.NotificationClass, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "BACnetApplicationTagUnsignedIntegerBuilder failed"))
	}
	return b
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) Build() (BACnetServiceAckGetEnrollmentSummary, error) {
	if b.ObjectIdentifier == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'objectIdentifier' not set"))
	}
	if b.EventType == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'eventType' not set"))
	}
	if b.EventState == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'eventState' not set"))
	}
	if b.Priority == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'priority' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._BACnetServiceAckGetEnrollmentSummary.deepCopy(), nil
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) MustBuild() BACnetServiceAckGetEnrollmentSummary {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) Done() BACnetServiceAckBuilder {
	return b.parentBuilder
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) buildForBACnetServiceAck() (BACnetServiceAck, error) {
	return b.Build()
}

func (b *_BACnetServiceAckGetEnrollmentSummaryBuilder) DeepCopy() any {
	_copy := b.CreateBACnetServiceAckGetEnrollmentSummaryBuilder().(*_BACnetServiceAckGetEnrollmentSummaryBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateBACnetServiceAckGetEnrollmentSummaryBuilder creates a BACnetServiceAckGetEnrollmentSummaryBuilder
func (b *_BACnetServiceAckGetEnrollmentSummary) CreateBACnetServiceAckGetEnrollmentSummaryBuilder() BACnetServiceAckGetEnrollmentSummaryBuilder {
	if b == nil {
		return NewBACnetServiceAckGetEnrollmentSummaryBuilder()
	}
	return &_BACnetServiceAckGetEnrollmentSummaryBuilder{_BACnetServiceAckGetEnrollmentSummary: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetServiceAckGetEnrollmentSummary) GetServiceChoice() BACnetConfirmedServiceChoice {
	return BACnetConfirmedServiceChoice_GET_ENROLLMENT_SUMMARY
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetServiceAckGetEnrollmentSummary) GetParent() BACnetServiceAckContract {
	return m.BACnetServiceAckContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetServiceAckGetEnrollmentSummary) GetObjectIdentifier() BACnetApplicationTagObjectIdentifier {
	return m.ObjectIdentifier
}

func (m *_BACnetServiceAckGetEnrollmentSummary) GetEventType() BACnetEventTypeTagged {
	return m.EventType
}

func (m *_BACnetServiceAckGetEnrollmentSummary) GetEventState() BACnetEventStateTagged {
	return m.EventState
}

func (m *_BACnetServiceAckGetEnrollmentSummary) GetPriority() BACnetApplicationTagUnsignedInteger {
	return m.Priority
}

func (m *_BACnetServiceAckGetEnrollmentSummary) GetNotificationClass() BACnetApplicationTagUnsignedInteger {
	return m.NotificationClass
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastBACnetServiceAckGetEnrollmentSummary(structType any) BACnetServiceAckGetEnrollmentSummary {
	if casted, ok := structType.(BACnetServiceAckGetEnrollmentSummary); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetServiceAckGetEnrollmentSummary); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetServiceAckGetEnrollmentSummary) GetTypeName() string {
	return "BACnetServiceAckGetEnrollmentSummary"
}

func (m *_BACnetServiceAckGetEnrollmentSummary) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.BACnetServiceAckContract.(*_BACnetServiceAck).GetLengthInBits(ctx))

	// Simple field (objectIdentifier)
	lengthInBits += m.ObjectIdentifier.GetLengthInBits(ctx)

	// Simple field (eventType)
	lengthInBits += m.EventType.GetLengthInBits(ctx)

	// Simple field (eventState)
	lengthInBits += m.EventState.GetLengthInBits(ctx)

	// Simple field (priority)
	lengthInBits += m.Priority.GetLengthInBits(ctx)

	// Optional Field (notificationClass)
	if m.NotificationClass != nil {
		lengthInBits += m.NotificationClass.GetLengthInBits(ctx)
	}

	return lengthInBits
}

func (m *_BACnetServiceAckGetEnrollmentSummary) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_BACnetServiceAckGetEnrollmentSummary) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_BACnetServiceAck, serviceAckLength uint32) (__bACnetServiceAckGetEnrollmentSummary BACnetServiceAckGetEnrollmentSummary, err error) {
	m.BACnetServiceAckContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetServiceAckGetEnrollmentSummary"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetServiceAckGetEnrollmentSummary")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	objectIdentifier, err := ReadSimpleField[BACnetApplicationTagObjectIdentifier](ctx, "objectIdentifier", ReadComplex[BACnetApplicationTagObjectIdentifier](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagObjectIdentifier](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'objectIdentifier' field"))
	}
	m.ObjectIdentifier = objectIdentifier

	eventType, err := ReadSimpleField[BACnetEventTypeTagged](ctx, "eventType", ReadComplex[BACnetEventTypeTagged](BACnetEventTypeTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'eventType' field"))
	}
	m.EventType = eventType

	eventState, err := ReadSimpleField[BACnetEventStateTagged](ctx, "eventState", ReadComplex[BACnetEventStateTagged](BACnetEventStateTaggedParseWithBufferProducer((uint8)(uint8(0)), (TagClass)(TagClass_APPLICATION_TAGS)), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'eventState' field"))
	}
	m.EventState = eventState

	priority, err := ReadSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "priority", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'priority' field"))
	}
	m.Priority = priority

	var notificationClass BACnetApplicationTagUnsignedInteger
	_notificationClass, err := ReadOptionalField[BACnetApplicationTagUnsignedInteger](ctx, "notificationClass", ReadComplex[BACnetApplicationTagUnsignedInteger](BACnetApplicationTagParseWithBufferProducer[BACnetApplicationTagUnsignedInteger](), readBuffer), true)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'notificationClass' field"))
	}
	if _notificationClass != nil {
		notificationClass = *_notificationClass
		m.NotificationClass = notificationClass
	}

	if closeErr := readBuffer.CloseContext("BACnetServiceAckGetEnrollmentSummary"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetServiceAckGetEnrollmentSummary")
	}

	return m, nil
}

func (m *_BACnetServiceAckGetEnrollmentSummary) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetServiceAckGetEnrollmentSummary) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetServiceAckGetEnrollmentSummary"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetServiceAckGetEnrollmentSummary")
		}

		if err := WriteSimpleField[BACnetApplicationTagObjectIdentifier](ctx, "objectIdentifier", m.GetObjectIdentifier(), WriteComplex[BACnetApplicationTagObjectIdentifier](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'objectIdentifier' field")
		}

		if err := WriteSimpleField[BACnetEventTypeTagged](ctx, "eventType", m.GetEventType(), WriteComplex[BACnetEventTypeTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'eventType' field")
		}

		if err := WriteSimpleField[BACnetEventStateTagged](ctx, "eventState", m.GetEventState(), WriteComplex[BACnetEventStateTagged](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'eventState' field")
		}

		if err := WriteSimpleField[BACnetApplicationTagUnsignedInteger](ctx, "priority", m.GetPriority(), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'priority' field")
		}

		if err := WriteOptionalField[BACnetApplicationTagUnsignedInteger](ctx, "notificationClass", GetRef(m.GetNotificationClass()), WriteComplex[BACnetApplicationTagUnsignedInteger](writeBuffer), true); err != nil {
			return errors.Wrap(err, "Error serializing 'notificationClass' field")
		}

		if popErr := writeBuffer.PopContext("BACnetServiceAckGetEnrollmentSummary"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetServiceAckGetEnrollmentSummary")
		}
		return nil
	}
	return m.BACnetServiceAckContract.(*_BACnetServiceAck).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_BACnetServiceAckGetEnrollmentSummary) IsBACnetServiceAckGetEnrollmentSummary() {}

func (m *_BACnetServiceAckGetEnrollmentSummary) DeepCopy() any {
	return m.deepCopy()
}

func (m *_BACnetServiceAckGetEnrollmentSummary) deepCopy() *_BACnetServiceAckGetEnrollmentSummary {
	if m == nil {
		return nil
	}
	_BACnetServiceAckGetEnrollmentSummaryCopy := &_BACnetServiceAckGetEnrollmentSummary{
		m.BACnetServiceAckContract.(*_BACnetServiceAck).deepCopy(),
		m.ObjectIdentifier.DeepCopy().(BACnetApplicationTagObjectIdentifier),
		m.EventType.DeepCopy().(BACnetEventTypeTagged),
		m.EventState.DeepCopy().(BACnetEventStateTagged),
		m.Priority.DeepCopy().(BACnetApplicationTagUnsignedInteger),
		m.NotificationClass.DeepCopy().(BACnetApplicationTagUnsignedInteger),
	}
	m.BACnetServiceAckContract.(*_BACnetServiceAck)._SubType = m
	return _BACnetServiceAckGetEnrollmentSummaryCopy
}

func (m *_BACnetServiceAckGetEnrollmentSummary) String() string {
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
