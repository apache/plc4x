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
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

// Code generated by code-generation. DO NOT EDIT.

// MeteringDataOilConsumption is the corresponding interface of MeteringDataOilConsumption
type MeteringDataOilConsumption interface {
	utils.LengthAware
	utils.Serializable
	MeteringData
	// GetL returns L (property field)
	GetL() uint32
}

// MeteringDataOilConsumptionExactly can be used when we want exactly this type and not a type which fulfills MeteringDataOilConsumption.
// This is useful for switch cases.
type MeteringDataOilConsumptionExactly interface {
	MeteringDataOilConsumption
	isMeteringDataOilConsumption() bool
}

// _MeteringDataOilConsumption is the data-structure of this message
type _MeteringDataOilConsumption struct {
	*_MeteringData
	L uint32
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_MeteringDataOilConsumption) InitializeParent(parent MeteringData, commandTypeContainer MeteringCommandTypeContainer, argument byte) {
	m.CommandTypeContainer = commandTypeContainer
	m.Argument = argument
}

func (m *_MeteringDataOilConsumption) GetParent() MeteringData {
	return m._MeteringData
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_MeteringDataOilConsumption) GetL() uint32 {
	return m.L
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewMeteringDataOilConsumption factory function for _MeteringDataOilConsumption
func NewMeteringDataOilConsumption(L uint32, commandTypeContainer MeteringCommandTypeContainer, argument byte) *_MeteringDataOilConsumption {
	_result := &_MeteringDataOilConsumption{
		L:             L,
		_MeteringData: NewMeteringData(commandTypeContainer, argument),
	}
	_result._MeteringData._MeteringDataChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastMeteringDataOilConsumption(structType interface{}) MeteringDataOilConsumption {
	if casted, ok := structType.(MeteringDataOilConsumption); ok {
		return casted
	}
	if casted, ok := structType.(*MeteringDataOilConsumption); ok {
		return *casted
	}
	return nil
}

func (m *_MeteringDataOilConsumption) GetTypeName() string {
	return "MeteringDataOilConsumption"
}

func (m *_MeteringDataOilConsumption) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_MeteringDataOilConsumption) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (L)
	lengthInBits += 32

	return lengthInBits
}

func (m *_MeteringDataOilConsumption) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func MeteringDataOilConsumptionParse(theBytes []byte) (MeteringDataOilConsumption, error) {
	return MeteringDataOilConsumptionParseWithBuffer(utils.NewReadBufferByteBased(theBytes))
}

func MeteringDataOilConsumptionParseWithBuffer(readBuffer utils.ReadBuffer) (MeteringDataOilConsumption, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("MeteringDataOilConsumption"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for MeteringDataOilConsumption")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (L)
	_L, _LErr := readBuffer.ReadUint32("L", 32)
	if _LErr != nil {
		return nil, errors.Wrap(_LErr, "Error parsing 'L' field of MeteringDataOilConsumption")
	}
	L := _L

	if closeErr := readBuffer.CloseContext("MeteringDataOilConsumption"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for MeteringDataOilConsumption")
	}

	// Create a partially initialized instance
	_child := &_MeteringDataOilConsumption{
		_MeteringData: &_MeteringData{},
		L:             L,
	}
	_child._MeteringData._MeteringDataChildRequirements = _child
	return _child, nil
}

func (m *_MeteringDataOilConsumption) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_MeteringDataOilConsumption) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("MeteringDataOilConsumption"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for MeteringDataOilConsumption")
		}

		// Simple Field (L)
		L := uint32(m.GetL())
		_LErr := writeBuffer.WriteUint32("L", 32, (L))
		if _LErr != nil {
			return errors.Wrap(_LErr, "Error serializing 'L' field")
		}

		if popErr := writeBuffer.PopContext("MeteringDataOilConsumption"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for MeteringDataOilConsumption")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_MeteringDataOilConsumption) isMeteringDataOilConsumption() bool {
	return true
}

func (m *_MeteringDataOilConsumption) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
