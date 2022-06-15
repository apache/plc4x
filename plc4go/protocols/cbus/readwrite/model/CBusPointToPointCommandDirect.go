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
)

// Code generated by code-generation. DO NOT EDIT.

// CBusPointToPointCommandDirect is the corresponding interface of CBusPointToPointCommandDirect
type CBusPointToPointCommandDirect interface {
	CBusPointToPointCommand
	// GetUnitAddress returns UnitAddress (property field)
	GetUnitAddress() UnitAddress
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _CBusPointToPointCommandDirect is the data-structure of this message
type _CBusPointToPointCommandDirect struct {
	*_CBusPointToPointCommand
	UnitAddress UnitAddress

	// Arguments.
	Srchk bool
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_CBusPointToPointCommandDirect) InitializeParent(parent CBusPointToPointCommand, bridgeAddressCountPeek uint16, calData CALData, crc Checksum, peekAlpha byte, alpha Alpha) {
	m.BridgeAddressCountPeek = bridgeAddressCountPeek
	m.CalData = calData
	m.Crc = crc
	m.PeekAlpha = peekAlpha
	m.Alpha = alpha
}

func (m *_CBusPointToPointCommandDirect) GetParent() CBusPointToPointCommand {
	return m._CBusPointToPointCommand
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_CBusPointToPointCommandDirect) GetUnitAddress() UnitAddress {
	return m.UnitAddress
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewCBusPointToPointCommandDirect factory function for _CBusPointToPointCommandDirect
func NewCBusPointToPointCommandDirect(unitAddress UnitAddress, bridgeAddressCountPeek uint16, calData CALData, crc Checksum, peekAlpha byte, alpha Alpha, srchk bool) *_CBusPointToPointCommandDirect {
	_result := &_CBusPointToPointCommandDirect{
		UnitAddress:              unitAddress,
		_CBusPointToPointCommand: NewCBusPointToPointCommand(bridgeAddressCountPeek, calData, crc, peekAlpha, alpha, srchk),
	}
	_result._CBusPointToPointCommand._CBusPointToPointCommandChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastCBusPointToPointCommandDirect(structType interface{}) CBusPointToPointCommandDirect {
	if casted, ok := structType.(CBusPointToPointCommandDirect); ok {
		return casted
	}
	if casted, ok := structType.(*CBusPointToPointCommandDirect); ok {
		return *casted
	}
	return nil
}

func (m *_CBusPointToPointCommandDirect) GetTypeName() string {
	return "CBusPointToPointCommandDirect"
}

func (m *_CBusPointToPointCommandDirect) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_CBusPointToPointCommandDirect) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (unitAddress)
	lengthInBits += m.UnitAddress.GetLengthInBits()

	// Reserved Field (reserved)
	lengthInBits += 8

	return lengthInBits
}

func (m *_CBusPointToPointCommandDirect) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func CBusPointToPointCommandDirectParse(readBuffer utils.ReadBuffer, srchk bool) (CBusPointToPointCommandDirect, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("CBusPointToPointCommandDirect"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for CBusPointToPointCommandDirect")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (unitAddress)
	if pullErr := readBuffer.PullContext("unitAddress"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for unitAddress")
	}
	_unitAddress, _unitAddressErr := UnitAddressParse(readBuffer)
	if _unitAddressErr != nil {
		return nil, errors.Wrap(_unitAddressErr, "Error parsing 'unitAddress' field")
	}
	unitAddress := _unitAddress.(UnitAddress)
	if closeErr := readBuffer.CloseContext("unitAddress"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for unitAddress")
	}

	// Reserved Field (Compartmentalized so the "reserved" variable can't leak)
	{
		reserved, _err := readBuffer.ReadUint8("reserved", 8)
		if _err != nil {
			return nil, errors.Wrap(_err, "Error parsing 'reserved' field")
		}
		if reserved != uint8(0x00) {
			log.Info().Fields(map[string]interface{}{
				"expected value": uint8(0x00),
				"got value":      reserved,
			}).Msg("Got unexpected response.")
		}
	}

	if closeErr := readBuffer.CloseContext("CBusPointToPointCommandDirect"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for CBusPointToPointCommandDirect")
	}

	// Create a partially initialized instance
	_child := &_CBusPointToPointCommandDirect{
		UnitAddress:              unitAddress,
		_CBusPointToPointCommand: &_CBusPointToPointCommand{},
	}
	_child._CBusPointToPointCommand._CBusPointToPointCommandChildRequirements = _child
	return _child, nil
}

func (m *_CBusPointToPointCommandDirect) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("CBusPointToPointCommandDirect"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for CBusPointToPointCommandDirect")
		}

		// Simple Field (unitAddress)
		if pushErr := writeBuffer.PushContext("unitAddress"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for unitAddress")
		}
		_unitAddressErr := writeBuffer.WriteSerializable(m.GetUnitAddress())
		if popErr := writeBuffer.PopContext("unitAddress"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for unitAddress")
		}
		if _unitAddressErr != nil {
			return errors.Wrap(_unitAddressErr, "Error serializing 'unitAddress' field")
		}

		// Reserved Field (reserved)
		{
			_err := writeBuffer.WriteUint8("reserved", 8, uint8(0x00))
			if _err != nil {
				return errors.Wrap(_err, "Error serializing 'reserved' field")
			}
		}

		if popErr := writeBuffer.PopContext("CBusPointToPointCommandDirect"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for CBusPointToPointCommandDirect")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_CBusPointToPointCommandDirect) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
