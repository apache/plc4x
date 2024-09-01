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

// ModificationInfo is the corresponding interface of ModificationInfo
type ModificationInfo interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetModificationTime returns ModificationTime (property field)
	GetModificationTime() int64
	// GetUpdateType returns UpdateType (property field)
	GetUpdateType() HistoryUpdateType
	// GetUserName returns UserName (property field)
	GetUserName() PascalString
}

// ModificationInfoExactly can be used when we want exactly this type and not a type which fulfills ModificationInfo.
// This is useful for switch cases.
type ModificationInfoExactly interface {
	ModificationInfo
	isModificationInfo() bool
}

// _ModificationInfo is the data-structure of this message
type _ModificationInfo struct {
	*_ExtensionObjectDefinition
	ModificationTime int64
	UpdateType       HistoryUpdateType
	UserName         PascalString
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_ModificationInfo) GetIdentifier() string {
	return "11218"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_ModificationInfo) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_ModificationInfo) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_ModificationInfo) GetModificationTime() int64 {
	return m.ModificationTime
}

func (m *_ModificationInfo) GetUpdateType() HistoryUpdateType {
	return m.UpdateType
}

func (m *_ModificationInfo) GetUserName() PascalString {
	return m.UserName
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewModificationInfo factory function for _ModificationInfo
func NewModificationInfo(modificationTime int64, updateType HistoryUpdateType, userName PascalString) *_ModificationInfo {
	_result := &_ModificationInfo{
		ModificationTime:           modificationTime,
		UpdateType:                 updateType,
		UserName:                   userName,
		_ExtensionObjectDefinition: NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastModificationInfo(structType any) ModificationInfo {
	if casted, ok := structType.(ModificationInfo); ok {
		return casted
	}
	if casted, ok := structType.(*ModificationInfo); ok {
		return *casted
	}
	return nil
}

func (m *_ModificationInfo) GetTypeName() string {
	return "ModificationInfo"
}

func (m *_ModificationInfo) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (modificationTime)
	lengthInBits += 64

	// Simple field (updateType)
	lengthInBits += 32

	// Simple field (userName)
	lengthInBits += m.UserName.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_ModificationInfo) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func ModificationInfoParse(ctx context.Context, theBytes []byte, identifier string) (ModificationInfo, error) {
	return ModificationInfoParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func ModificationInfoParseWithBufferProducer(identifier string) func(ctx context.Context, readBuffer utils.ReadBuffer) (ModificationInfo, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (ModificationInfo, error) {
		return ModificationInfoParseWithBuffer(ctx, readBuffer, identifier)
	}
}

func ModificationInfoParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (ModificationInfo, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("ModificationInfo"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ModificationInfo")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	modificationTime, err := ReadSimpleField(ctx, "modificationTime", ReadSignedLong(readBuffer, uint8(64)))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'modificationTime' field"))
	}

	updateType, err := ReadEnumField[HistoryUpdateType](ctx, "updateType", "HistoryUpdateType", ReadEnum(HistoryUpdateTypeByValue, ReadUnsignedInt(readBuffer, uint8(32))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'updateType' field"))
	}

	userName, err := ReadSimpleField[PascalString](ctx, "userName", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'userName' field"))
	}

	if closeErr := readBuffer.CloseContext("ModificationInfo"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ModificationInfo")
	}

	// Create a partially initialized instance
	_child := &_ModificationInfo{
		_ExtensionObjectDefinition: &_ExtensionObjectDefinition{},
		ModificationTime:           modificationTime,
		UpdateType:                 updateType,
		UserName:                   userName,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_ModificationInfo) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_ModificationInfo) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("ModificationInfo"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ModificationInfo")
		}

		if err := WriteSimpleField[int64](ctx, "modificationTime", m.GetModificationTime(), WriteSignedLong(writeBuffer, 64)); err != nil {
			return errors.Wrap(err, "Error serializing 'modificationTime' field")
		}

		if err := WriteSimpleEnumField[HistoryUpdateType](ctx, "updateType", "HistoryUpdateType", m.GetUpdateType(), WriteEnum[HistoryUpdateType, uint32](HistoryUpdateType.GetValue, HistoryUpdateType.PLC4XEnumName, WriteUnsignedInt(writeBuffer, 32))); err != nil {
			return errors.Wrap(err, "Error serializing 'updateType' field")
		}

		if err := WriteSimpleField[PascalString](ctx, "userName", m.GetUserName(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'userName' field")
		}

		if popErr := writeBuffer.PopContext("ModificationInfo"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ModificationInfo")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_ModificationInfo) isModificationInfo() bool {
	return true
}

func (m *_ModificationInfo) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
