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
	"fmt"
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// Constant values.
const CBusCommand_INITIATOR byte = 0x5C

// CBusCommand is the corresponding interface of CBusCommand
type CBusCommand interface {
	// GetHeader returns Header (property field)
	GetHeader() CBusHeader
	// GetDestinationAddressType returns DestinationAddressType (virtual field)
	GetDestinationAddressType() DestinationAddressType
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _CBusCommand is the data-structure of this message
type _CBusCommand struct {
	_CBusCommandChildRequirements
	Header CBusHeader

	// Arguments.
	Srchk bool
}

type _CBusCommandChildRequirements interface {
	GetLengthInBits() uint16
	GetLengthInBitsConditional(lastItem bool) uint16
	Serialize(writeBuffer utils.WriteBuffer) error
}

type CBusCommandParent interface {
	SerializeParent(writeBuffer utils.WriteBuffer, child CBusCommand, serializeChildFunction func() error) error
	GetTypeName() string
}

type CBusCommandChild interface {
	Serialize(writeBuffer utils.WriteBuffer) error
	InitializeParent(parent CBusCommand, header CBusHeader)
	GetParent() *CBusCommand

	GetTypeName() string
	CBusCommand
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CBusCommand) GetHeader() CBusHeader {
	return m.Header
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for virtual fields.
///////////////////////

func (m *_CBusCommand) GetDestinationAddressType() DestinationAddressType {
	return CastDestinationAddressType(m.GetHeader().GetDestinationAddressType())
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for const fields.
///////////////////////

func (m *_CBusCommand) GetInitiator() byte {
	return CBusCommand_INITIATOR
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewCBusCommand factory function for _CBusCommand
func NewCBusCommand(header CBusHeader, srchk bool) *_CBusCommand {
	return &_CBusCommand{Header: header, Srchk: srchk}
}

// Deprecated: use the interface for direct cast
func CastCBusCommand(structType interface{}) CBusCommand {
	if casted, ok := structType.(CBusCommand); ok {
		return casted
	}
	if casted, ok := structType.(*CBusCommand); ok {
		return *casted
	}
	return nil
}

func (m *_CBusCommand) GetTypeName() string {
	return "CBusCommand"
}

func (m *_CBusCommand) GetParentLengthInBits() uint16 {
	lengthInBits := uint16(0)

	// Const Field (initiator)
	lengthInBits += 8

	// Simple field (header)
	lengthInBits += m.Header.GetLengthInBits()

	// A virtual field doesn't have any in- or output.

	return lengthInBits
}

func (m *_CBusCommand) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func CBusCommandParse(readBuffer utils.ReadBuffer, srchk bool) (CBusCommand, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CBusCommand"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CBusCommand")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Const Field (initiator)
	initiator, _initiatorErr := readBuffer.ReadByte("initiator")
	if _initiatorErr != nil {
		return nil, errors.Wrap(_initiatorErr, "Error parsing 'initiator' field")
	}
	if initiator != CBusCommand_INITIATOR {
		return nil, errors.New("Expected constant value " + fmt.Sprintf("%d", CBusCommand_INITIATOR) + " but got " + fmt.Sprintf("%d", initiator))
	}

	// Simple Field (header)
	if pullErr := readBuffer.PullContext("header"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for header")
	}
	_header, _headerErr := CBusHeaderParse(readBuffer)
	if _headerErr != nil {
		return nil, errors.Wrap(_headerErr, "Error parsing 'header' field")
	}
	header := _header.(CBusHeader)
	if closeErr := readBuffer.CloseContext("header"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for header")
	}

	// Virtual field
	_destinationAddressType := header.GetDestinationAddressType()
	destinationAddressType := DestinationAddressType(_destinationAddressType)
	_ = destinationAddressType

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	type CBusCommandChildSerializeRequirement interface {
		CBusCommand
		InitializeParent(CBusCommand, CBusHeader)
		GetParent() CBusCommand
	}
	var _childTemp interface{}
	var _child CBusCommandChildSerializeRequirement
	var typeSwitchError error
	switch {
	case destinationAddressType == DestinationAddressType_PointToPointToMultiPoint: // CBusCommandPointToPointToMultiPoint
		_childTemp, typeSwitchError = CBusCommandPointToPointToMultiPointParse(readBuffer, srchk)
		_child = _childTemp.(CBusCommandChildSerializeRequirement)
	case destinationAddressType == DestinationAddressType_PointToMultiPoint: // CBusCommandPointToMultiPoint
		_childTemp, typeSwitchError = CBusCommandPointToMultiPointParse(readBuffer, srchk)
		_child = _childTemp.(CBusCommandChildSerializeRequirement)
	case destinationAddressType == DestinationAddressType_PointToPoint: // CBusCommandPointToPoint
		_childTemp, typeSwitchError = CBusCommandPointToPointParse(readBuffer, srchk)
		_child = _childTemp.(CBusCommandChildSerializeRequirement)
	default:
		// TODO: return actual type
		typeSwitchError = errors.New("Unmapped type")
	}
	if typeSwitchError != nil {
		return nil, errors.Wrap(typeSwitchError, "Error parsing sub-type for type-switch.")
	}

	if closeErr := readBuffer.CloseContext("CBusCommand"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CBusCommand")
	}

	// Finish initializing
	_child.InitializeParent(_child, header)
	return _child, nil
}

func (pm *_CBusCommand) SerializeParent(writeBuffer utils.WriteBuffer, child CBusCommand, serializeChildFunction func() error) error {
	// We redirect all calls through client as some methods are only implemented there
	m := child
	_ = m
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("CBusCommand"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for CBusCommand")
	}

	// Const Field (initiator)
	_initiatorErr := writeBuffer.WriteByte("initiator", 0x5C)
	if _initiatorErr != nil {
		return errors.Wrap(_initiatorErr, "Error serializing 'initiator' field")
	}

	// Simple Field (header)
	if pushErr := writeBuffer.PushContext("header"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for header")
	}
	_headerErr := writeBuffer.WriteSerializable(m.GetHeader())
	if popErr := writeBuffer.PopContext("header"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for header")
	}
	if _headerErr != nil {
		return errors.Wrap(_headerErr, "Error serializing 'header' field")
	}
	// Virtual field
	if _destinationAddressTypeErr := writeBuffer.WriteVirtual("destinationAddressType", m.GetDestinationAddressType()); _destinationAddressTypeErr != nil {
		return errors.Wrap(_destinationAddressTypeErr, "Error serializing 'destinationAddressType' field")
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("CBusCommand"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for CBusCommand")
	}
	return nil
}

func (m *_CBusCommand) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
