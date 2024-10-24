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

// EventFilter is the corresponding interface of EventFilter
type EventFilter interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	utils.Copyable
	ExtensionObjectDefinition
	// GetSelectClauses returns SelectClauses (property field)
	GetSelectClauses() []SimpleAttributeOperand
	// GetWhereClause returns WhereClause (property field)
	GetWhereClause() ContentFilter
	// IsEventFilter is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsEventFilter()
	// CreateBuilder creates a EventFilterBuilder
	CreateEventFilterBuilder() EventFilterBuilder
}

// _EventFilter is the data-structure of this message
type _EventFilter struct {
	ExtensionObjectDefinitionContract
	SelectClauses []SimpleAttributeOperand
	WhereClause   ContentFilter
}

var _ EventFilter = (*_EventFilter)(nil)
var _ ExtensionObjectDefinitionRequirements = (*_EventFilter)(nil)

// NewEventFilter factory function for _EventFilter
func NewEventFilter(selectClauses []SimpleAttributeOperand, whereClause ContentFilter) *_EventFilter {
	if whereClause == nil {
		panic("whereClause of type ContentFilter for EventFilter must not be nil")
	}
	_result := &_EventFilter{
		ExtensionObjectDefinitionContract: NewExtensionObjectDefinition(),
		SelectClauses:                     selectClauses,
		WhereClause:                       whereClause,
	}
	_result.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = _result
	return _result
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Builder
///////////////////////

// EventFilterBuilder is a builder for EventFilter
type EventFilterBuilder interface {
	utils.Copyable
	// WithMandatoryFields adds all mandatory fields (convenience for using multiple builder calls)
	WithMandatoryFields(selectClauses []SimpleAttributeOperand, whereClause ContentFilter) EventFilterBuilder
	// WithSelectClauses adds SelectClauses (property field)
	WithSelectClauses(...SimpleAttributeOperand) EventFilterBuilder
	// WithWhereClause adds WhereClause (property field)
	WithWhereClause(ContentFilter) EventFilterBuilder
	// WithWhereClauseBuilder adds WhereClause (property field) which is build by the builder
	WithWhereClauseBuilder(func(ContentFilterBuilder) ContentFilterBuilder) EventFilterBuilder
	// Build builds the EventFilter or returns an error if something is wrong
	Build() (EventFilter, error)
	// MustBuild does the same as Build but panics on error
	MustBuild() EventFilter
}

// NewEventFilterBuilder() creates a EventFilterBuilder
func NewEventFilterBuilder() EventFilterBuilder {
	return &_EventFilterBuilder{_EventFilter: new(_EventFilter)}
}

type _EventFilterBuilder struct {
	*_EventFilter

	parentBuilder *_ExtensionObjectDefinitionBuilder

	err *utils.MultiError
}

var _ (EventFilterBuilder) = (*_EventFilterBuilder)(nil)

func (b *_EventFilterBuilder) setParent(contract ExtensionObjectDefinitionContract) {
	b.ExtensionObjectDefinitionContract = contract
}

func (b *_EventFilterBuilder) WithMandatoryFields(selectClauses []SimpleAttributeOperand, whereClause ContentFilter) EventFilterBuilder {
	return b.WithSelectClauses(selectClauses...).WithWhereClause(whereClause)
}

func (b *_EventFilterBuilder) WithSelectClauses(selectClauses ...SimpleAttributeOperand) EventFilterBuilder {
	b.SelectClauses = selectClauses
	return b
}

func (b *_EventFilterBuilder) WithWhereClause(whereClause ContentFilter) EventFilterBuilder {
	b.WhereClause = whereClause
	return b
}

func (b *_EventFilterBuilder) WithWhereClauseBuilder(builderSupplier func(ContentFilterBuilder) ContentFilterBuilder) EventFilterBuilder {
	builder := builderSupplier(b.WhereClause.CreateContentFilterBuilder())
	var err error
	b.WhereClause, err = builder.Build()
	if err != nil {
		if b.err == nil {
			b.err = &utils.MultiError{MainError: errors.New("sub builder failed")}
		}
		b.err.Append(errors.Wrap(err, "ContentFilterBuilder failed"))
	}
	return b
}

func (b *_EventFilterBuilder) Build() (EventFilter, error) {
	if b.WhereClause == nil {
		if b.err == nil {
			b.err = new(utils.MultiError)
		}
		b.err.Append(errors.New("mandatory field 'whereClause' not set"))
	}
	if b.err != nil {
		return nil, errors.Wrap(b.err, "error occurred during build")
	}
	return b._EventFilter.deepCopy(), nil
}

func (b *_EventFilterBuilder) MustBuild() EventFilter {
	build, err := b.Build()
	if err != nil {
		panic(err)
	}
	return build
}

// Done is used to finish work on this child and return to the parent builder
func (b *_EventFilterBuilder) Done() ExtensionObjectDefinitionBuilder {
	return b.parentBuilder
}

func (b *_EventFilterBuilder) buildForExtensionObjectDefinition() (ExtensionObjectDefinition, error) {
	return b.Build()
}

func (b *_EventFilterBuilder) DeepCopy() any {
	_copy := b.CreateEventFilterBuilder().(*_EventFilterBuilder)
	if b.err != nil {
		_copy.err = b.err.DeepCopy().(*utils.MultiError)
	}
	return _copy
}

// CreateEventFilterBuilder creates a EventFilterBuilder
func (b *_EventFilter) CreateEventFilterBuilder() EventFilterBuilder {
	if b == nil {
		return NewEventFilterBuilder()
	}
	return &_EventFilterBuilder{_EventFilter: b.deepCopy()}
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_EventFilter) GetExtensionId() int32 {
	return int32(727)
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_EventFilter) GetParent() ExtensionObjectDefinitionContract {
	return m.ExtensionObjectDefinitionContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_EventFilter) GetSelectClauses() []SimpleAttributeOperand {
	return m.SelectClauses
}

func (m *_EventFilter) GetWhereClause() ContentFilter {
	return m.WhereClause
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// Deprecated: use the interface for direct cast
func CastEventFilter(structType any) EventFilter {
	if casted, ok := structType.(EventFilter); ok {
		return casted
	}
	if casted, ok := structType.(*EventFilter); ok {
		return *casted
	}
	return nil
}

func (m *_EventFilter) GetTypeName() string {
	return "EventFilter"
}

func (m *_EventFilter) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).GetLengthInBits(ctx))

	// Implicit Field (noOfSelectClauses)
	lengthInBits += 32

	// Array field
	if len(m.SelectClauses) > 0 {
		for _curItem, element := range m.SelectClauses {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.SelectClauses), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	// Simple field (whereClause)
	lengthInBits += m.WhereClause.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_EventFilter) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_EventFilter) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_ExtensionObjectDefinition, extensionId int32) (__eventFilter EventFilter, err error) {
	m.ExtensionObjectDefinitionContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("EventFilter"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for EventFilter")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	noOfSelectClauses, err := ReadImplicitField[int32](ctx, "noOfSelectClauses", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfSelectClauses' field"))
	}
	_ = noOfSelectClauses

	selectClauses, err := ReadCountArrayField[SimpleAttributeOperand](ctx, "selectClauses", ReadComplex[SimpleAttributeOperand](ExtensionObjectDefinitionParseWithBufferProducer[SimpleAttributeOperand]((int32)(int32(603))), readBuffer), uint64(noOfSelectClauses))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'selectClauses' field"))
	}
	m.SelectClauses = selectClauses

	whereClause, err := ReadSimpleField[ContentFilter](ctx, "whereClause", ReadComplex[ContentFilter](ExtensionObjectDefinitionParseWithBufferProducer[ContentFilter]((int32)(int32(588))), readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'whereClause' field"))
	}
	m.WhereClause = whereClause

	if closeErr := readBuffer.CloseContext("EventFilter"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for EventFilter")
	}

	return m, nil
}

func (m *_EventFilter) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_EventFilter) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("EventFilter"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for EventFilter")
		}
		noOfSelectClauses := int32(utils.InlineIf(bool((m.GetSelectClauses()) == (nil)), func() any { return int32(-(int32(1))) }, func() any { return int32(int32(len(m.GetSelectClauses()))) }).(int32))
		if err := WriteImplicitField(ctx, "noOfSelectClauses", noOfSelectClauses, WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfSelectClauses' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "selectClauses", m.GetSelectClauses(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'selectClauses' field")
		}

		if err := WriteSimpleField[ContentFilter](ctx, "whereClause", m.GetWhereClause(), WriteComplex[ContentFilter](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'whereClause' field")
		}

		if popErr := writeBuffer.PopContext("EventFilter"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for EventFilter")
		}
		return nil
	}
	return m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_EventFilter) IsEventFilter() {}

func (m *_EventFilter) DeepCopy() any {
	return m.deepCopy()
}

func (m *_EventFilter) deepCopy() *_EventFilter {
	if m == nil {
		return nil
	}
	_EventFilterCopy := &_EventFilter{
		m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition).deepCopy(),
		utils.DeepCopySlice[SimpleAttributeOperand, SimpleAttributeOperand](m.SelectClauses),
		m.WhereClause.DeepCopy().(ContentFilter),
	}
	m.ExtensionObjectDefinitionContract.(*_ExtensionObjectDefinition)._SubType = m
	return _EventFilterCopy
}

func (m *_EventFilter) String() string {
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
