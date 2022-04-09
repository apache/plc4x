/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// The data-structure of this message
type IdentifyReplyCommand struct {
	Child IIdentifyReplyCommandChild
}

// The corresponding interface
type IIdentifyReplyCommand interface {
	// GetAttribute returns Attribute (discriminator field)
	GetAttribute() Attribute
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

type IIdentifyReplyCommandParent interface {
	SerializeParent(writeBuffer utils.WriteBuffer, child IIdentifyReplyCommand, serializeChildFunction func() error) error
	GetTypeName() string
}

type IIdentifyReplyCommandChild interface {
	Serialize(writeBuffer utils.WriteBuffer) error
	InitializeParent(parent *IdentifyReplyCommand)
	GetParent() *IdentifyReplyCommand

	GetTypeName() string
	IIdentifyReplyCommand
}

// NewIdentifyReplyCommand factory function for IdentifyReplyCommand
func NewIdentifyReplyCommand() *IdentifyReplyCommand {
	return &IdentifyReplyCommand{}
}

func CastIdentifyReplyCommand(structType interface{}) *IdentifyReplyCommand {
	if casted, ok := structType.(IdentifyReplyCommand); ok {
		return &casted
	}
	if casted, ok := structType.(*IdentifyReplyCommand); ok {
		return casted
	}
	if casted, ok := structType.(IIdentifyReplyCommandChild); ok {
		return casted.GetParent()
	}
	return nil
}

func (m *IdentifyReplyCommand) GetTypeName() string {
	return "IdentifyReplyCommand"
}

func (m *IdentifyReplyCommand) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *IdentifyReplyCommand) GetLengthInBitsConditional(lastItem bool) uint16 {
	return m.Child.GetLengthInBits()
}

func (m *IdentifyReplyCommand) GetParentLengthInBits() uint16 {
	lengthInBits := uint16(0)

	return lengthInBits
}

func (m *IdentifyReplyCommand) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func IdentifyReplyCommandParse(readBuffer utils.ReadBuffer, attribute Attribute) (*IdentifyReplyCommand, error) {
	if pullErr := readBuffer.PullContext("IdentifyReplyCommand"); pullErr != nil {
		return nil, pullErr
	}
	currentPos := readBuffer.GetPos()
	_ = currentPos

	// Switch Field (Depending on the discriminator values, passes the instantiation to a sub-type)
	type IdentifyReplyCommandChild interface {
		InitializeParent(*IdentifyReplyCommand)
		GetParent() *IdentifyReplyCommand
	}
	var _child IdentifyReplyCommandChild
	var typeSwitchError error
	switch {
	case attribute == Attribute_Manufacturer: // IdentifyReplyCommandManufacturer
		_child, typeSwitchError = IdentifyReplyCommandManufacturerParse(readBuffer, attribute)
	case attribute == Attribute_Type: // IdentifyReplyCommandType
		_child, typeSwitchError = IdentifyReplyCommandTypeParse(readBuffer, attribute)
	case attribute == Attribute_FirmwareVersion: // IdentifyReplyCommandFirmwareVersion
		_child, typeSwitchError = IdentifyReplyCommandFirmwareVersionParse(readBuffer, attribute)
	case attribute == Attribute_Summary: // IdentifyReplyCommandFirmwareSummary
		_child, typeSwitchError = IdentifyReplyCommandFirmwareSummaryParse(readBuffer, attribute)
	case attribute == Attribute_ExtendedDiagnosticSummary: // IdentifyReplyCommandExtendedDiagnosticSummary
		_child, typeSwitchError = IdentifyReplyCommandExtendedDiagnosticSummaryParse(readBuffer, attribute)
	case attribute == Attribute_NetworkTerminalLevels: // IdentifyReplyCommandNetworkTerminalLevels
		_child, typeSwitchError = IdentifyReplyCommandNetworkTerminalLevelsParse(readBuffer, attribute)
	case attribute == Attribute_TerminalLevel: // IdentifyReplyCommandTerminalLevels
		_child, typeSwitchError = IdentifyReplyCommandTerminalLevelsParse(readBuffer, attribute)
	case attribute == Attribute_NetworkVoltage: // IdentifyReplyCommandNetworkVoltage
		_child, typeSwitchError = IdentifyReplyCommandNetworkVoltageParse(readBuffer, attribute)
	case attribute == Attribute_GAVValuesCurrent: // IdentifyReplyCommandGAVValuesCurrent
		_child, typeSwitchError = IdentifyReplyCommandGAVValuesCurrentParse(readBuffer, attribute)
	case attribute == Attribute_GAVValuesStored: // IdentifyReplyCommandGAVValuesStored
		_child, typeSwitchError = IdentifyReplyCommandGAVValuesStoredParse(readBuffer, attribute)
	case attribute == Attribute_GAVPhysicalAddresses: // IdentifyReplyCommandGAVPhysicalAddresses
		_child, typeSwitchError = IdentifyReplyCommandGAVPhysicalAddressesParse(readBuffer, attribute)
	case attribute == Attribute_LogicalAssignment: // IdentifyReplyCommandLogicalAssignment
		_child, typeSwitchError = IdentifyReplyCommandLogicalAssignmentParse(readBuffer, attribute)
	case attribute == Attribute_Delays: // IdentifyReplyCommandDelays
		_child, typeSwitchError = IdentifyReplyCommandDelaysParse(readBuffer, attribute)
	case attribute == Attribute_MinimumLevels: // IdentifyReplyCommandMinimumLevels
		_child, typeSwitchError = IdentifyReplyCommandMinimumLevelsParse(readBuffer, attribute)
	case attribute == Attribute_MaximumLevels: // IdentifyReplyCommandMaximumLevels
		_child, typeSwitchError = IdentifyReplyCommandMaximumLevelsParse(readBuffer, attribute)
	case attribute == Attribute_CurrentSenseLevels: // IdentifyReplyCommandCurrentSenseLevels
		_child, typeSwitchError = IdentifyReplyCommandCurrentSenseLevelsParse(readBuffer, attribute)
	case attribute == Attribute_OutputUnitSummary: // IdentifyReplyCommandOutputUnitSummary
		_child, typeSwitchError = IdentifyReplyCommandOutputUnitSummaryParse(readBuffer, attribute)
	case attribute == Attribute_DSIStatus: // IdentifyReplyCommandDSIStatus
		_child, typeSwitchError = IdentifyReplyCommandDSIStatusParse(readBuffer, attribute)
	default:
		// TODO: return actual type
		typeSwitchError = errors.New("Unmapped type")
	}
	if typeSwitchError != nil {
		return nil, errors.Wrap(typeSwitchError, "Error parsing sub-type for type-switch.")
	}

	if closeErr := readBuffer.CloseContext("IdentifyReplyCommand"); closeErr != nil {
		return nil, closeErr
	}

	// Finish initializing
	_child.InitializeParent(_child.GetParent())
	return _child.GetParent(), nil
}

func (m *IdentifyReplyCommand) Serialize(writeBuffer utils.WriteBuffer) error {
	return m.Child.Serialize(writeBuffer)
}

func (m *IdentifyReplyCommand) SerializeParent(writeBuffer utils.WriteBuffer, child IIdentifyReplyCommand, serializeChildFunction func() error) error {
	if pushErr := writeBuffer.PushContext("IdentifyReplyCommand"); pushErr != nil {
		return pushErr
	}

	// Switch field (Depending on the discriminator values, passes the serialization to a sub-type)
	if _typeSwitchErr := serializeChildFunction(); _typeSwitchErr != nil {
		return errors.Wrap(_typeSwitchErr, "Error serializing sub-type field")
	}

	if popErr := writeBuffer.PopContext("IdentifyReplyCommand"); popErr != nil {
		return popErr
	}
	return nil
}

func (m *IdentifyReplyCommand) String() string {
	if m == nil {
		return "<nil>"
	}
	buffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := m.Serialize(buffer); err != nil {
		return err.Error()
	}
	return buffer.GetBox().String()
}
