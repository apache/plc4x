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

// SecurityDataExitDelayStarted is the corresponding interface of SecurityDataExitDelayStarted
type SecurityDataExitDelayStarted interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	SecurityData
}

// SecurityDataExitDelayStartedExactly can be used when we want exactly this type and not a type which fulfills SecurityDataExitDelayStarted.
// This is useful for switch cases.
type SecurityDataExitDelayStartedExactly interface {
	SecurityDataExitDelayStarted
	isSecurityDataExitDelayStarted() bool
}

// _SecurityDataExitDelayStarted is the data-structure of this message
type _SecurityDataExitDelayStarted struct {
	*_SecurityData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_SecurityDataExitDelayStarted) InitializeParent(parent SecurityData, commandTypeContainer SecurityCommandTypeContainer, argument byte) {
	m.CommandTypeContainer = commandTypeContainer
	m.Argument = argument
}

func (m *_SecurityDataExitDelayStarted) GetParent() SecurityData {
	return m._SecurityData
}

// NewSecurityDataExitDelayStarted factory function for _SecurityDataExitDelayStarted
func NewSecurityDataExitDelayStarted(commandTypeContainer SecurityCommandTypeContainer, argument byte) *_SecurityDataExitDelayStarted {
	_result := &_SecurityDataExitDelayStarted{
		_SecurityData: NewSecurityData(commandTypeContainer, argument),
	}
	_result._SecurityData._SecurityDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastSecurityDataExitDelayStarted(structType any) SecurityDataExitDelayStarted {
	if casted, ok := structType.(SecurityDataExitDelayStarted); ok {
		return casted
	}
	if casted, ok := structType.(*SecurityDataExitDelayStarted); ok {
		return *casted
	}
	return nil
}

func (m *_SecurityDataExitDelayStarted) GetTypeName() string {
	return "SecurityDataExitDelayStarted"
}

func (m *_SecurityDataExitDelayStarted) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits(ctx))

	return lengthInBits
}

func (m *_SecurityDataExitDelayStarted) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func SecurityDataExitDelayStartedParse(theBytes []byte) (SecurityDataExitDelayStarted, error) {
	return SecurityDataExitDelayStartedParseWithBuffer(context.Background(), utils.NewReadBufferByteBased(theBytes))
}

func SecurityDataExitDelayStartedParseWithBuffer(ctx context.Context, readBuffer utils.ReadBuffer) (SecurityDataExitDelayStarted, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SecurityDataExitDelayStarted"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SecurityDataExitDelayStarted")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("SecurityDataExitDelayStarted"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SecurityDataExitDelayStarted")
	}

	// Create a partially initialized instance
	_child := &_SecurityDataExitDelayStarted{
		_SecurityData: &_SecurityData{},
	}
	_child._SecurityData._SecurityDataChildRequirements = _child
	return _child, nil
}

func (m *_SecurityDataExitDelayStarted) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_SecurityDataExitDelayStarted) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("SecurityDataExitDelayStarted"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for SecurityDataExitDelayStarted")
		}

		if popErr := writeBuffer.PopContext("SecurityDataExitDelayStarted"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for SecurityDataExitDelayStarted")
		}
		return nil
	}
	return m.SerializeParent(ctx, writeBuffer, m, ser)
}

func (m *_SecurityDataExitDelayStarted) isSecurityDataExitDelayStarted() bool {
	return true
}

func (m *_SecurityDataExitDelayStarted) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
