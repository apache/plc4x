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
)

// Code generated by code-generation. DO NOT EDIT.

// SALDataHvacActuator is the corresponding interface of SALDataHvacActuator
type SALDataHvacActuator interface {
	utils.LengthAware
	utils.Serializable
	SALData
	// GetVentilationData returns VentilationData (property field)
	GetVentilationData() LightingData
}

// SALDataHvacActuatorExactly can be used when we want exactly this type and not a type which fulfills SALDataHvacActuator.
// This is useful for switch cases.
type SALDataHvacActuatorExactly interface {
	SALDataHvacActuator
	isSALDataHvacActuator() bool
}

// _SALDataHvacActuator is the data-structure of this message
type _SALDataHvacActuator struct {
	*_SALData
	VentilationData LightingData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_SALDataHvacActuator) GetApplicationId() ApplicationId {
	return ApplicationId_HVAC_ACTUATOR
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_SALDataHvacActuator) InitializeParent(parent SALData, salData SALData) {
	m.SalData = salData
}

func (m *_SALDataHvacActuator) GetParent() SALData {
	return m._SALData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_SALDataHvacActuator) GetVentilationData() LightingData {
	return m.VentilationData
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewSALDataHvacActuator factory function for _SALDataHvacActuator
func NewSALDataHvacActuator(ventilationData LightingData, salData SALData) *_SALDataHvacActuator {
	_result := &_SALDataHvacActuator{
		VentilationData: ventilationData,
		_SALData:        NewSALData(salData),
	}
	_result._SALData._SALDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastSALDataHvacActuator(structType interface{}) SALDataHvacActuator {
	if casted, ok := structType.(SALDataHvacActuator); ok {
		return casted
	}
	if casted, ok := structType.(*SALDataHvacActuator); ok {
		return *casted
	}
	return nil
}

func (m *_SALDataHvacActuator) GetTypeName() string {
	return "SALDataHvacActuator"
}

func (m *_SALDataHvacActuator) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_SALDataHvacActuator) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (ventilationData)
	lengthInBits += m.VentilationData.GetLengthInBits()

	return lengthInBits
}

func (m *_SALDataHvacActuator) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func SALDataHvacActuatorParse(readBuffer utils.ReadBuffer, applicationId ApplicationId) (SALDataHvacActuator, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("SALDataHvacActuator"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for SALDataHvacActuator")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (ventilationData)
	if pullErr := readBuffer.PullContext("ventilationData"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for ventilationData")
	}
	_ventilationData, _ventilationDataErr := LightingDataParse(readBuffer)
	if _ventilationDataErr != nil {
		return nil, errors.Wrap(_ventilationDataErr, "Error parsing 'ventilationData' field of SALDataHvacActuator")
	}
	ventilationData := _ventilationData.(LightingData)
	if closeErr := readBuffer.CloseContext("ventilationData"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for ventilationData")
	}

	if closeErr := readBuffer.CloseContext("SALDataHvacActuator"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for SALDataHvacActuator")
	}

	// Create a partially initialized instance
	_child := &_SALDataHvacActuator{
		_SALData:        &_SALData{},
		VentilationData: ventilationData,
	}
	_child._SALData._SALDataChildRequirements = _child
	return _child, nil
}

func (m *_SALDataHvacActuator) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("SALDataHvacActuator"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for SALDataHvacActuator")
		}

		// Simple Field (ventilationData)
		if pushErr := writeBuffer.PushContext("ventilationData"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for ventilationData")
		}
		_ventilationDataErr := writeBuffer.WriteSerializable(m.GetVentilationData())
		if popErr := writeBuffer.PopContext("ventilationData"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for ventilationData")
		}
		if _ventilationDataErr != nil {
			return errors.Wrap(_ventilationDataErr, "Error serializing 'ventilationData' field")
		}

		if popErr := writeBuffer.PopContext("SALDataHvacActuator"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for SALDataHvacActuator")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_SALDataHvacActuator) isSALDataHvacActuator() bool {
	return true
}

func (m *_SALDataHvacActuator) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
