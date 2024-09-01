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

// UserManagementDataType is the corresponding interface of UserManagementDataType
type UserManagementDataType interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	ExtensionObjectDefinition
	// GetUserName returns UserName (property field)
	GetUserName() PascalString
	// GetUserConfiguration returns UserConfiguration (property field)
	GetUserConfiguration() UserConfigurationMask
	// GetDescription returns Description (property field)
	GetDescription() PascalString
}

// UserManagementDataTypeExactly can be used when we want exactly this type and not a type which fulfills UserManagementDataType.
// This is useful for switch cases.
type UserManagementDataTypeExactly interface {
	UserManagementDataType
	isUserManagementDataType() bool
}

// _UserManagementDataType is the data-structure of this message
type _UserManagementDataType struct {
	*_ExtensionObjectDefinition
	UserName          PascalString
	UserConfiguration UserConfigurationMask
	Description       PascalString
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_UserManagementDataType) GetIdentifier() string {
	return "24283"
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_UserManagementDataType) InitializeParent(parent ExtensionObjectDefinition) {}

func (m *_UserManagementDataType) GetParent() ExtensionObjectDefinition {
	return m._ExtensionObjectDefinition
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_UserManagementDataType) GetUserName() PascalString {
	return m.UserName
}

func (m *_UserManagementDataType) GetUserConfiguration() UserConfigurationMask {
	return m.UserConfiguration
}

func (m *_UserManagementDataType) GetDescription() PascalString {
	return m.Description
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewUserManagementDataType factory function for _UserManagementDataType
func NewUserManagementDataType(userName PascalString, userConfiguration UserConfigurationMask, description PascalString) *_UserManagementDataType {
	_result := &_UserManagementDataType{
		UserName:                   userName,
		UserConfiguration:          userConfiguration,
		Description:                description,
		_ExtensionObjectDefinition: NewExtensionObjectDefinition(),
	}
	_result._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastUserManagementDataType(structType any) UserManagementDataType {
	if casted, ok := structType.(UserManagementDataType); ok {
		return casted
	}
	if casted, ok := structType.(*UserManagementDataType); ok {
		return *casted
	}
	return nil
}

func (m *_UserManagementDataType) GetTypeName() string {
	return "UserManagementDataType"
}

func (m *_UserManagementDataType) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	// Simple field (userName)
	lengthInBits += m.UserName.GetLengthInBits(ctx)

	// Simple field (userConfiguration)
	lengthInBits += 32

	// Simple field (description)
	lengthInBits += m.Description.GetLengthInBits(ctx)

	return lengthInBits
}

func (m *_UserManagementDataType) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func UserManagementDataTypeParse(ctx context.Context, theBytes []byte, identifier string) (UserManagementDataType, error) {
	return UserManagementDataTypeParseWithBuffer(ctx, utils.NewReadBufferByteBased(theBytes), identifier)
}

func UserManagementDataTypeParseWithBufferProducer(identifier string) func(ctx context.Context, readBuffer utils.ReadBuffer) (UserManagementDataType, error) {
	return func(ctx context.Context, readBuffer utils.ReadBuffer) (UserManagementDataType, error) {
		return UserManagementDataTypeParseWithBuffer(ctx, readBuffer, identifier)
	}
}

func UserManagementDataTypeParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer, identifier string) (UserManagementDataType, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("UserManagementDataType"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for UserManagementDataType")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	userName, err := ReadSimpleField[PascalString](ctx, "userName", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'userName' field"))
	}

	userConfiguration, err := ReadEnumField[UserConfigurationMask](ctx, "userConfiguration", "UserConfigurationMask", ReadEnum(UserConfigurationMaskByValue, ReadUnsignedInt(readBuffer, uint8(32))))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'userConfiguration' field"))
	}

	description, err := ReadSimpleField[PascalString](ctx, "description", ReadComplex[PascalString](PascalStringParseWithBuffer, readBuffer))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'description' field"))
	}

	if closeErr := readBuffer.CloseContext("UserManagementDataType"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for UserManagementDataType")
	}

	// Create a partially initialized instance
	_child := &_UserManagementDataType{
		_ExtensionObjectDefinition: &_ExtensionObjectDefinition{},
		UserName:                   userName,
		UserConfiguration:          userConfiguration,
		Description:                description,
	}
	_child._ExtensionObjectDefinition._ExtensionObjectDefinitionChildRequirements = _child
	return _child, nil
}

func (m *_UserManagementDataType) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_UserManagementDataType) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("UserManagementDataType"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for UserManagementDataType")
		}

		if err := WriteSimpleField[PascalString](ctx, "userName", m.GetUserName(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'userName' field")
		}

		if err := WriteSimpleEnumField[UserConfigurationMask](ctx, "userConfiguration", "UserConfigurationMask", m.GetUserConfiguration(), WriteEnum[UserConfigurationMask, uint32](UserConfigurationMask.GetValue, UserConfigurationMask.PLC4XEnumName, WriteUnsignedInt(writeBuffer, 32))); err != nil {
			return errors.Wrap(err, "Error serializing 'userConfiguration' field")
		}

		if err := WriteSimpleField[PascalString](ctx, "description", m.GetDescription(), WriteComplex[PascalString](writeBuffer)); err != nil {
			return errors.Wrap(err, "Error serializing 'description' field")
		}

		if popErr := writeBuffer.PopContext("UserManagementDataType"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for UserManagementDataType")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_UserManagementDataType) isUserManagementDataType() bool {
	return true
}

func (m *_UserManagementDataType) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
