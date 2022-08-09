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
	"github.com/apache/plc4x/plc4go/internal/spi/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
)

// Code generated by code-generation. DO NOT EDIT.

// BACnetLandingCallStatus is the corresponding interface of BACnetLandingCallStatus
type BACnetLandingCallStatus interface {
	utils.LengthAware
	utils.Serializable
	// GetFloorNumber returns FloorNumber (property field)
	GetFloorNumber() BACnetContextTagUnsignedInteger
	// GetCommand returns Command (property field)
	GetCommand() BACnetLandingCallStatusCommand
	// GetFloorText returns FloorText (property field)
	GetFloorText() BACnetContextTagCharacterString
}

// BACnetLandingCallStatusExactly can be used when we want exactly this type and not a type which fulfills BACnetLandingCallStatus.
// This is useful for switch cases.
type BACnetLandingCallStatusExactly interface {
	BACnetLandingCallStatus
	isBACnetLandingCallStatus() bool
}

// _BACnetLandingCallStatus is the data-structure of this message
type _BACnetLandingCallStatus struct {
	FloorNumber BACnetContextTagUnsignedInteger
	Command     BACnetLandingCallStatusCommand
	FloorText   BACnetContextTagCharacterString
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetLandingCallStatus) GetFloorNumber() BACnetContextTagUnsignedInteger {
	return m.FloorNumber
}

func (m *_BACnetLandingCallStatus) GetCommand() BACnetLandingCallStatusCommand {
	return m.Command
}

func (m *_BACnetLandingCallStatus) GetFloorText() BACnetContextTagCharacterString {
	return m.FloorText
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetLandingCallStatus factory function for _BACnetLandingCallStatus
func NewBACnetLandingCallStatus(floorNumber BACnetContextTagUnsignedInteger, command BACnetLandingCallStatusCommand, floorText BACnetContextTagCharacterString) *_BACnetLandingCallStatus {
	return &_BACnetLandingCallStatus{FloorNumber: floorNumber, Command: command, FloorText: floorText}
}

// Deprecated: use the interface for direct cast
func CastBACnetLandingCallStatus(structType interface{}) BACnetLandingCallStatus {
	if casted, ok := structType.(BACnetLandingCallStatus); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetLandingCallStatus); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetLandingCallStatus) GetTypeName() string {
	return "BACnetLandingCallStatus"
}

func (m *_BACnetLandingCallStatus) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetLandingCallStatus) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(0)

	// Simple field (floorNumber)
	lengthInBits += m.FloorNumber.GetLengthInBits()

	// Simple field (command)
	lengthInBits += m.Command.GetLengthInBits()

	// Optional Field (floorText)
	if m.FloorText != nil {
		lengthInBits += m.FloorText.GetLengthInBits()
	}

	return lengthInBits
}

func (m *_BACnetLandingCallStatus) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetLandingCallStatusParse(readBuffer utils.ReadBuffer) (BACnetLandingCallStatus, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetLandingCallStatus"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetLandingCallStatus")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (floorNumber)
	if pullErr := readBuffer.PullContext("floorNumber"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for floorNumber")
	}
	_floorNumber, _floorNumberErr := BACnetContextTagParse(readBuffer, uint8(uint8(0)), BACnetDataType(BACnetDataType_UNSIGNED_INTEGER))
	if _floorNumberErr != nil {
		return nil, errors.Wrap(_floorNumberErr, "Error parsing 'floorNumber' field of BACnetLandingCallStatus")
	}
	floorNumber := _floorNumber.(BACnetContextTagUnsignedInteger)
	if closeErr := readBuffer.CloseContext("floorNumber"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for floorNumber")
	}

	// Simple Field (command)
	if pullErr := readBuffer.PullContext("command"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for command")
	}
	_command, _commandErr := BACnetLandingCallStatusCommandParse(readBuffer)
	if _commandErr != nil {
		return nil, errors.Wrap(_commandErr, "Error parsing 'command' field of BACnetLandingCallStatus")
	}
	command := _command.(BACnetLandingCallStatusCommand)
	if closeErr := readBuffer.CloseContext("command"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for command")
	}

	// Optional Field (floorText) (Can be skipped, if a given expression evaluates to false)
	var floorText BACnetContextTagCharacterString = nil
	{
		currentPos = positionAware.GetPos()
		if pullErr := readBuffer.PullContext("floorText"); pullErr != nil {
			return nil, errors.Wrap(pullErr, "Error pulling for floorText")
		}
		_val, _err := BACnetContextTagParse(readBuffer, uint8(3), BACnetDataType_CHARACTER_STRING)
		switch {
		case errors.Is(_err, utils.ParseAssertError{}) || errors.Is(_err, io.EOF):
			log.Debug().Err(_err).Msg("Resetting position because optional threw an error")
			readBuffer.Reset(currentPos)
		case _err != nil:
			return nil, errors.Wrap(_err, "Error parsing 'floorText' field of BACnetLandingCallStatus")
		default:
			floorText = _val.(BACnetContextTagCharacterString)
			if closeErr := readBuffer.CloseContext("floorText"); closeErr != nil {
				return nil, errors.Wrap(closeErr, "Error closing for floorText")
			}
		}
	}

	if closeErr := readBuffer.CloseContext("BACnetLandingCallStatus"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetLandingCallStatus")
	}

	// Create the instance
	return &_BACnetLandingCallStatus{
		FloorNumber: floorNumber,
		Command:     command,
		FloorText:   floorText,
	}, nil
}

func (m *_BACnetLandingCallStatus) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	if pushErr := writeBuffer.PushContext("BACnetLandingCallStatus"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for BACnetLandingCallStatus")
	}

	// Simple Field (floorNumber)
	if pushErr := writeBuffer.PushContext("floorNumber"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for floorNumber")
	}
	_floorNumberErr := writeBuffer.WriteSerializable(m.GetFloorNumber())
	if popErr := writeBuffer.PopContext("floorNumber"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for floorNumber")
	}
	if _floorNumberErr != nil {
		return errors.Wrap(_floorNumberErr, "Error serializing 'floorNumber' field")
	}

	// Simple Field (command)
	if pushErr := writeBuffer.PushContext("command"); pushErr != nil {
		return errors.Wrap(pushErr, "Error pushing for command")
	}
	_commandErr := writeBuffer.WriteSerializable(m.GetCommand())
	if popErr := writeBuffer.PopContext("command"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for command")
	}
	if _commandErr != nil {
		return errors.Wrap(_commandErr, "Error serializing 'command' field")
	}

	// Optional Field (floorText) (Can be skipped, if the value is null)
	var floorText BACnetContextTagCharacterString = nil
	if m.GetFloorText() != nil {
		if pushErr := writeBuffer.PushContext("floorText"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for floorText")
		}
		floorText = m.GetFloorText()
		_floorTextErr := writeBuffer.WriteSerializable(floorText)
		if popErr := writeBuffer.PopContext("floorText"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for floorText")
		}
		if _floorTextErr != nil {
			return errors.Wrap(_floorTextErr, "Error serializing 'floorText' field")
		}
	}

	if popErr := writeBuffer.PopContext("BACnetLandingCallStatus"); popErr != nil {
		return errors.Wrap(popErr, "Error popping for BACnetLandingCallStatus")
	}
	return nil
}

func (m *_BACnetLandingCallStatus) isBACnetLandingCallStatus() bool {
	return true
}

func (m *_BACnetLandingCallStatus) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
