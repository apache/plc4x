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

// HistoryEventFieldList is the corresponding interface of HistoryEventFieldList
type HistoryEventFieldList interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetNoOfEventFields returns NoOfEventFields (property field)
	GetNoOfEventFields() int32
	// GetEventFields returns EventFields (property field)
	GetEventFields() []Variant
}

// HistoryEventFieldListExactly can be used when we want exactly this type and not a type which fulfills HistoryEventFieldList.
// This is useful for switch cases.
type HistoryEventFieldListExactly interface {
	HistoryEventFieldList
	isHistoryEventFieldList() bool
}

// _HistoryEventFieldList is the data-structure of this message
type _HistoryEventFieldList struct {
	*_ExtensionObjectDefinition
	NoOfEventFields int32
	EventFields     []Variant
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_HistoryEventFieldList) GetIdentifier() string {
	return "922"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_HistoryEventFieldList) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_HistoryEventFieldList) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_HistoryEventFieldList) GetNoOfEventFields() int32 {
	return m.NoOfEventFields
}

func (m *_HistoryEventFieldList) GetEventFields() []Variant {
	return m.EventFields
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewHistoryEventFieldList factory function for _HistoryEventFieldList
func NewHistoryEventFieldList(noOfEventFields int32, eventFields []Variant) *_HistoryEventFieldList {
	_result := &_HistoryEventFieldList{
		NoOfEventFields:            noOfEventFields,
		EventFields:                eventFields,
		_ExtensionObjectDefinition: NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastHistoryEventFieldList(structType any) HistoryEventFieldList {
	if casted, ok := structType.(HistoryEventFieldList); ok {
		return casted
	}
	if casted, ok := structType.(*HistoryEventFieldList); ok {
		return *casted
	}
	return nil
}

func (m *_HistoryEventFieldList) GetTypeName() string {
	return "HistoryEventFieldList"
}

func (m *_HistoryEventFieldList) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (noOfEventFields)
	lengthInBits += 32

	// Array field
	if len(m.EventFields) > 0 {
		for _curItem, element := range m.EventFields {
			arrayCtx := utils.CreateArrayContext(ctx, len(m.EventFields), _curItem)
			_ = arrayCtx
			_ = _curItem
			lengthInBits += element.(interface{ GetLengthInBits(context.Context) uint16 }).GetLengthInBits(arrayCtx)
		}
	}

	return lengthInBits
}

func (m *_HistoryEventFieldList) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func HistoryEventFieldListParse(ctx context.Context, theBytes []byte, identifier string) (HistoryEventFieldList, error) {
	return HistoryEventFieldListParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func HistoryEventFieldListParseWithBufferProducer(identifier string) func(ctx context.Context, readBuffer utils.ReadBuffer) (HistoryEventFieldList, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (HistoryEventFieldList, error) {
		return HistoryEventFieldListParseWithBuffer(ctx, readBuffer, identifier)
	}
}

func HistoryEventFieldListParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (HistoryEventFieldList, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("HistoryEventFieldList"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for HistoryEventFieldList")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	noOfEventFields, err := ReadSimpleField(ctx, "noOfEventFields", ReadSignedInt(readBuffer, uint8(32)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'noOfEventFields' field"))
	}

	eventFields, err := ReadCountArrayField[Variant](ctx, "eventFields", ReadComplex[Variant](VariantParseWithBuffer, readBuffer), uint64(noOfEventFields))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'eventFields' field"))
	}

	if closeErr := readBuffer.CloseContext("HistoryEventFieldList"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for HistoryEventFieldList")
	}

	// Create a partially initialized instance
	_child := &_HistoryEventFieldList{
		_ExtensionObjectDefinition: &_ExtensionObjectDefinition{},
		NoOfEventFields:            noOfEventFields,
		EventFields:                eventFields,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_HistoryEventFieldList) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_HistoryEventFieldList) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("HistoryEventFieldList"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for HistoryEventFieldList")
		}

		if err := WriteSimpleField[int32](ctx, "noOfEventFields", m.GetNoOfEventFields(), WriteSignedInt(writeBuffer, 32)); err != nil {
			return errors.Wrap(err, "Error serializing 'noOfEventFields' field")
		}

		if err := WriteComplexTypeArrayField(ctx, "eventFields", m.GetEventFields(), writeBuffer); err != nil {
			return errors.Wrap(err, "Error serializing 'eventFields' field")
		}

		if popErr := writeBuffer.PopContext("HistoryEventFieldList"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for HistoryEventFieldList")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_HistoryEventFieldList) isHistoryEventFieldList() bool {
	return true
}

func (m *_HistoryEventFieldList) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
