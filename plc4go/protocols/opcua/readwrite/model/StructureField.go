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
	"github.com/rs/zerolog"
)

// Code generated by code-generation. DO NOT EDIT.

// StructureField is the corresponding interface of StructureField
type StructureField interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetName returns Name (property field)
	GetName() PascalString
	// GetDescription returns Description (property field)
	GetDescription() LocalizedText
	// GetDataType returns DataType (property field)
	GetDataType() NodeId
	// GetValueRank returns ValueRank (property field)
	GetValueRank() int32
	// GetNoOfArrayDimensions returns NoOfArrayDimensions (property field)
	GetNoOfArrayDimensions() int32
	// GetArrayDimensions returns ArrayDimensions (property field)
	GetArrayDimensions() []uint32
	// GetMaxStringLength returns MaxStringLength (property field)
	GetMaxStringLength() uint32
	// GetIsOptional returns IsOptional (property field)
	GetIsOptional() bool
}

// StructureFieldExactly can be used when we want exactly this type and not a type which fulfills StructureField.
// This is useful for switch cases.
type StructureFieldExactly interface {
	StructureField
	isStructureField() bool
}

// _StructureField is the data-structure of this message
type _StructureField struct {
	*_ExtensionObjectDefinition
	Name                PascalString
	Description         LocalizedText
	DataType            NodeId
	ValueRank           int32
	NoOfArrayDimensions int32
	ArrayDimensions     []uint32
	MaxStringLength     uint32
	IsOptional          bool
	// Reserved Fields
	reservedField0 *uint8
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_StructureField) GetIdentifier() string {
	return "103"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_StructureField) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_StructureField) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_StructureField) GetName() PascalString {
	return m.Name
}

func (m *_StructureField) GetDescription() LocalizedText {
	return m.Description
}

func (m *_StructureField) GetDataType() NodeId {
	return m.DataType
}

func (m *_StructureField) GetValueRank() int32 {
	return m.ValueRank
}

func (m *_StructureField) GetNoOfArrayDimensions() int32 {
	return m.NoOfArrayDimensions
}

func (m *_StructureField) GetArrayDimensions() []uint32 {
	return m.ArrayDimensions
}

func (m *_StructureField) GetMaxStringLength() uint32 {
	return m.MaxStringLength
}

func (m *_StructureField) GetIsOptional() bool {
	return m.IsOptional
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewStructureField factory function for _StructureField
func NewStructureField(name PascalString, description LocalizedText, dataType NodeId, valueRank int32, noOfArrayDimensions int32, arrayDimensions []uint32, maxStringLength uint32, isOptional bool) *_StructureField {
	_result := &_StructureField{
		Name:                       name,
		Description:                description,
		DataType:                   dataType,
		ValueRank:                  valueRank,
		NoOfArrayDimensions:        noOfArrayDimensions,
		ArrayDimensions:            arrayDimensions,
		MaxStringLength:            maxStringLength,
		IsOptional:                 isOptional,
		_ExtensionObjectDefinition: NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastStructureField(structType any) StructureField {
	if casted, ok := structType.(StructureField); ok {
		return casted
	}
	if casted, ok := structType.(*StructureField); ok {
		return *casted
	}
	return nil
}

func (m *_StructureField) GetTypeName() string {
	return "StructureField"
}

func (m *_StructureField) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (name)
	lengthInBits += m.Name.GetLengthInBits(ctx)

	// Simple field (description)
	lengthInBits += m.Description.GetLengthInBits(ctx)

	// Simple field (dataType)
	lengthInBits += m.DataType.GetLengthInBits(ctx)

	// Simple field (valueRank)
	lengthInBits += 32

	// Simple field (noOfArrayDimensions)
	lengthInBits += 32

	// Array field
	if len(m.ArrayDimensions) > 0 {
		lengthInBits += 32 * uint16(len(m.ArrayDimensions))
	}

	// Simple field (maxStringLength)
	lengthInBits += 32

	// Reserved Field (reserved)
	lengthInBits += 7

	// Simple field (isOptional)
	lengthInBits += 1

	return lengthInBits
}

func (m *_StructureField) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func StructureFieldParse(ctx context.Context, theBytes []byte, identifier string) (StructureField, error) {
	return StructureFieldParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func StructureFieldParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (StructureField, error) {
	positionAware := readBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	if pullErr := readBuffer.PullContext("StructureField"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for StructureField")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (name)
	if pullErr := readBuffer.PullContext("name"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for name")
	}
	_name, _nameErr := PascalStringParseWithBuffer(ctx, readBuffer)
	if _nameErr != nil {
		return nil, errors.Wrap(_nameErr, "Error parsing 'name' field of StructureField")
	}
	name := _name.(PascalString)
	if closeErr := readBuffer.CloseContext("name"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for name")
	}

	// Simple Field (description)
	if pullErr := readBuffer.PullContext("description"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for description")
	}
	_description, _descriptionErr := LocalizedTextParseWithBuffer(ctx, readBuffer)
	if _descriptionErr != nil {
		return nil, errors.Wrap(_descriptionErr, "Error parsing 'description' field of StructureField")
	}
	description := _description.(LocalizedText)
	if closeErr := readBuffer.CloseContext("description"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for description")
	}

	// Simple Field (dataType)
	if pullErr := readBuffer.PullContext("dataType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for dataType")
	}
	_dataType, _dataTypeErr := NodeIdParseWithBuffer(ctx, readBuffer)
	if _dataTypeErr != nil {
		return nil, errors.Wrap(_dataTypeErr, "Error parsing 'dataType' field of StructureField")
	}
	dataType := _dataType.(NodeId)
	if closeErr := readBuffer.CloseContext("dataType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for dataType")
	}

	// Simple Field (valueRank)
	_valueRank, _valueRankErr := readBuffer.ReadInt32("valueRank", 32)
	if _valueRankErr != nil {
		return nil, errors.Wrap(_valueRankErr, "Error parsing 'valueRank' field of StructureField")
	}
	valueRank := _valueRank

	// Simple Field (noOfArrayDimensions)
	_noOfArrayDimensions, _noOfArrayDimensionsErr := readBuffer.ReadInt32("noOfArrayDimensions", 32)
	if _noOfArrayDimensionsErr != nil {
		return nil, errors.Wrap(_noOfArrayDimensionsErr, "Error parsing 'noOfArrayDimensions' field of StructureField")
	}
	noOfArrayDimensions := _noOfArrayDimensions

	// Array field (arrayDimensions)
	if pullErr := readBuffer.PullContext("arrayDimensions", utils.WithRenderAsList(true)); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for arrayDimensions")
	}
	// Count array
	arrayDimensions := make([]uint32, max(noOfArrayDimensions, 0))
	// This happens when the size is set conditional to 0
	if len(arrayDimensions) == 0 {
		arrayDimensions = nil
	}
	{
		_numItems := uint16(max(noOfArrayDimensions, 0))
		for _curItem := uint16(0); _curItem < _numItems; _curItem++ {
			arrayCtx := utils.CreateArrayContext(ctx, int(_numItems), int(_curItem))
			_ = arrayCtx
			_ = _curItem
			_item, _err := readBuffer.ReadUint32("", 32)
			if _err != nil {
				return nil, errors.Wrap(_err, "Error parsing 'arrayDimensions' field of StructureField")
			}
			arrayDimensions[_curItem] = _item
		}
	}
	if closeErr := readBuffer.CloseContext("arrayDimensions", utils.WithRenderAsList(true)); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for arrayDimensions")
	}

	// Simple Field (maxStringLength)
	_maxStringLength, _maxStringLengthErr := readBuffer.ReadUint32("maxStringLength", 32)
	if _maxStringLengthErr != nil {
		return nil, errors.Wrap(_maxStringLengthErr, "Error parsing 'maxStringLength' field of StructureField")
	}
	maxStringLength := _maxStringLength

	var reservedField0 *uint8
	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := readBuffer.ReadUint8("reserved", 7)
		if _err != nil {
			return nil, errors.Wrap(_err, "Error parsing 'reserved' field of StructureField")
		}
		if reserved != uint8(0x00) {
			log.Info().Fields(map[string]any{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Msg("Got unexpected response for reserved field.")
			// We save the value, so it can be re-serialized
			reservedField0 = &reserved
		}
	}

	// Simple Field (isOptional)
	_isOptional, _isOptionalErr := readBuffer.ReadBit("isOptional")
	if _isOptionalErr != nil {
		return nil, errors.Wrap(_isOptionalErr, "Error parsing 'isOptional' field of StructureField")
	}
	isOptional := _isOptional

	if closeErr := readBuffer.CloseContext("StructureField"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for StructureField")
	}

	// Create a partially initialized instance
	_child := &_StructureField{
		_ExtensionObjectDefinition: &_ExtensionObjectDefinition{},
		Name:                       name,
		Description:                description,
		DataType:                   dataType,
		ValueRank:                  valueRank,
		NoOfArrayDimensions:        noOfArrayDimensions,
		ArrayDimensions:            arrayDimensions,
		MaxStringLength:            maxStringLength,
		IsOptional:                 isOptional,
		reservedField0:             reservedField0,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_StructureField) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_StructureField) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("StructureField"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for StructureField")
		}

		// Simple Field (name)
		if pushErr := writeBuffer.PushContext("name"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for name")
		}
		_nameErr := writeBuffer.WriteSerializable(ctx, m.GetName())
		if popErr := writeBuffer.PopContext("name"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for name")
		}
		if _nameErr != nil {
			return errors.Wrap(_nameErr, "Error serializing 'name' field")
		}

		// Simple Field (description)
		if pushErr := writeBuffer.PushContext("description"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for description")
		}
		_descriptionErr := writeBuffer.WriteSerializable(ctx, m.GetDescription())
		if popErr := writeBuffer.PopContext("description"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for description")
		}
		if _descriptionErr != nil {
			return errors.Wrap(_descriptionErr, "Error serializing 'description' field")
		}

		// Simple Field (dataType)
		if pushErr := writeBuffer.PushContext("dataType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for dataType")
		}
		_dataTypeErr := writeBuffer.WriteSerializable(ctx, m.GetDataType())
		if popErr := writeBuffer.PopContext("dataType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for dataType")
		}
		if _dataTypeErr != nil {
			return errors.Wrap(_dataTypeErr, "Error serializing 'dataType' field")
		}

		// Simple Field (valueRank)
		valueRank := int32(m.GetValueRank())
		_valueRankErr := writeBuffer.WriteInt32("valueRank", 32, int32((valueRank)))
		if _valueRankErr != nil {
			return errors.Wrap(_valueRankErr, "Error serializing 'valueRank' field")
		}

		// Simple Field (noOfArrayDimensions)
		noOfArrayDimensions := int32(m.GetNoOfArrayDimensions())
		_noOfArrayDimensionsErr := writeBuffer.WriteInt32("noOfArrayDimensions", 32, int32((noOfArrayDimensions)))
		if _noOfArrayDimensionsErr != nil {
			return errors.Wrap(_noOfArrayDimensionsErr, "Error serializing 'noOfArrayDimensions' field")
		}

		// Array Field (arrayDimensions)
		if pushErr := writeBuffer.PushContext("arrayDimensions", utils.WithRenderAsList(true)); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for arrayDimensions")
		}
		for _curItem, _element := range m.GetArrayDimensions() {
			_ = _curItem
			_elementErr := writeBuffer.WriteUint32("", 32, uint32(_element))
			if _elementErr != nil {
				return errors.Wrap(_elementErr, "Error serializing 'arrayDimensions' field")
			}
		}
		if popErr := writeBuffer.PopContext("arrayDimensions", utils.WithRenderAsList(true)); popErr != nil {
			return errors.Wrap(popErr, "Error popping for arrayDimensions")
		}

		// Simple Field (maxStringLength)
		maxStringLength := uint32(m.GetMaxStringLength())
		_maxStringLengthErr := writeBuffer.WriteUint32("maxStringLength", 32, uint32((maxStringLength)))
		if _maxStringLengthErr != nil {
			return errors.Wrap(_maxStringLengthErr, "Error serializing 'maxStringLength' field")
		}

		// Reserved Field (reserved)
		{
			var reserved uint8 = uint8(0x00)
			if m.reservedField0 != nil {
				log.Info().Fields(map[string]any{
					"expected value": uint8(0x00),
					"got value":      reserved,
				}).Msg("Overriding reserved field with unexpected value.")
				reserved = *m.reservedField0
			}
			_err := writeBuffer.WriteUint8("reserved", 7, uint8(reserved))
			if _err != nil {
				return errors.Wrap(_err, "Error serializing 'reserved' field")
			}
		}

		// Simple Field (isOptional)
		isOptional := bool(m.GetIsOptional())
		_isOptionalErr := writeBuffer.WriteBit("isOptional", (isOptional))
		if _isOptionalErr != nil {
			return errors.Wrap(_isOptionalErr, "Error serializing 'isOptional' field")
		}

		if popErr := writeBuffer.PopContext("StructureField"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for StructureField")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_StructureField) isStructureField() bool {
	return true
}

func (m *_StructureField) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
