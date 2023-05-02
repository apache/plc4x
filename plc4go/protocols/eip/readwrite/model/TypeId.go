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

// TypeId is the corresponding interface of TypeId
type TypeId interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	// GetId returns Id (discriminator field)
	GetId() uint16
}

// TypeIdExactly can be used when we want exactly this type and not a type which fulfills TypeId.
// This is useful for switch cases.
type TypeIdExactly interface {
	TypeId
	isTypeId() bool
}

// _TypeId is the data-structure of this message
type _TypeId struct {
	_TypeIdChildRequirements
}

type _TypeIdChildRequirements interface {
	utils.Serializable
	GetLengthInBits(ctx context.Context) uint16
	GetId() uint16
}

type TypeIdParent interface {
	SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child TypeId, serializeChildFunction func() error) error
	GetTypeName() string
}

type TypeIdChild interface {
	utils.Serializable
	InitializeParent(parent TypeId)
	GetParent() *TypeId

	GetTypeName() string
	TypeId
}

// NewTypeId factory function for _TypeId
func NewTypeId() *_TypeId {
	return &_TypeId{}
}

// Deprecated: use the interface for direct cast
func CastTypeId(structType any) TypeId {
	if casted, ok := structType.(TypeId); ok {
		return casted
	}
	if casted, ok := structType.(*TypeId); ok {
		return *casted
	}
	return nil
}

func (m *_TypeId) GetTypeName() string {
	return "TypeId"
}

func (m *_TypeId) GetParentLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(0)
	// Discriminator Field (id)
	lengthInBits += 16

	return lengthInBits
}

func (m *_TypeId) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func TypeIdParse(theBytes []byte) (TypeId, error) {
	return TypeIdParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes))
}

func TypeIdParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (TypeId, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("TypeId"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for TypeId")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Discriminator Field (id) (Used as input to a switch field)
	id, _idErr := readBuffer.ReadUint16("id", 16)
	if _idErr != nil {
		return nil, errors.Wrap(_idErr, "Error parsing 'id' field of TypeId")
	}

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	type TypeIdChildSerializeRequirement interface {
		TypeId
		InitializeParent(TypeId)
		GetParent() TypeId
	}
	var _childTemp any
	var _child TypeIdChildSerializeRequirement
	var typeSwitchError error
	switch {
	case id == 0x0000: // NullAddressItem
		_childTemp, typeSwitchError = NullAddressItemParseWithBuffer(ctx, readBuffer)
	case id == 0x0100: // ServicesResponse
		_childTemp, typeSwitchError = ServicesResponseParseWithBuffer(ctx, readBuffer)
	case id == 0x00A1: // ConnectedAddressItem
		_childTemp, typeSwitchError = ConnectedAddressItemParseWithBuffer(ctx, readBuffer)
	case id == 0x00B1: // ConnectedDataItem
		_childTemp, typeSwitchError = ConnectedDataItemParseWithBuffer(ctx, readBuffer)
	case id == 0x00B2: // UnConnectedDataItem
		_childTemp, typeSwitchError = UnConnectedDataItemParseWithBuffer(ctx, readBuffer)
	default:
		typeSwitchError = errors.Errorf("Unmapped type for parameters [id=%v]", id)
	}
	if typeSwitchError != nil {
		return nil, errors.Wrap(typeSwitchError, "Error parsing sub-type for type-switch of TypeId")
	}
	_child = _childTemp.(TypeIdChildSerializeRequirement)

	if closeErr := readBuffer.CloseContext("TypeId"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for TypeId")
	}

	// Finish initializing
	_child.InitializeParent(_child)
	return _child, nil
}

func (pm *_TypeId) SerializeParent(ctx context.Context, writeBuffer utils.WriteBuffer, child TypeId, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("TypeId"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for TypeId")
	}

	// Discriminator Field (id) (Used as input to a switch field)
	id := uint16(child.GetId())
	_idErr := writeBuffer.WriteUint16("id", 16, (id))

	if _idErr != nil {
		return errors.Wrap(_idErr, "Error serializing 'id' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("TypeId"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for TypeId")
	}
	return nil
}

func (m *_TypeId) isTypeId() bool {
	return true
}

func (m *_TypeId) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
