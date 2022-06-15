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

// BACnetUnconfirmedServiceRequestUTCTimeSynchronization is the corresponding interface of BACnetUnconfirmedServiceRequestUTCTimeSynchronization
type BACnetUnconfirmedServiceRequestUTCTimeSynchronization interface {
	BACnetUnconfirmedServiceRequest
	// GetSynchronizedDate returns SynchronizedDate (property field)
	GetSynchronizedDate() BACnetApplicationTagDate
	// GetSynchronizedTime returns SynchronizedTime (property field)
	GetSynchronizedTime() BACnetApplicationTagTime
	// GetLengthInBytes returns the length in bytes
	GetLengthInBytes() uint16
	// GetLengthInBits returns the length in bits
	GetLengthInBits() uint16
	// Serialize serializes this type
	Serialize(writeBuffer utils.WriteBuffer) error
}

// _BACnetUnconfirmedServiceRequestUTCTimeSynchronization is the data-structure of this message
type _BACnetUnconfirmedServiceRequestUTCTimeSynchronization struct {
	*_BACnetUnconfirmedServiceRequest
	SynchronizedDate BACnetApplicationTagDate
	SynchronizedTime BACnetApplicationTagTime

	// Arguments.
	ServiceRequestLength uint16
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) GetServiceChoice() BACnetUnconfirmedServiceChoice {
	return BACnetUnconfirmedServiceChoice_UTC_TIME_SYNCHRONIZATION
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) InitializeParent(parent BACnetUnconfirmedServiceRequest) {
}

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) GetParent() BACnetUnconfirmedServiceRequest {
	return m._BACnetUnconfirmedServiceRequest
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) GetSynchronizedDate() BACnetApplicationTagDate {
	return m.SynchronizedDate
}

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) GetSynchronizedTime() BACnetApplicationTagTime {
	return m.SynchronizedTime
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetUnconfirmedServiceRequestUTCTimeSynchronization factory function for _BACnetUnconfirmedServiceRequestUTCTimeSynchronization
func NewBACnetUnconfirmedServiceRequestUTCTimeSynchronization(synchronizedDate BACnetApplicationTagDate, synchronizedTime BACnetApplicationTagTime, serviceRequestLength uint16) *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization {
	_result := &_BACnetUnconfirmedServiceRequestUTCTimeSynchronization{
		SynchronizedDate:                 synchronizedDate,
		SynchronizedTime:                 synchronizedTime,
		_BACnetUnconfirmedServiceRequest: NewBACnetUnconfirmedServiceRequest(serviceRequestLength),
	}
	_result._BACnetUnconfirmedServiceRequest._BACnetUnconfirmedServiceRequestChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetUnconfirmedServiceRequestUTCTimeSynchronization(structType interface{}) BACnetUnconfirmedServiceRequestUTCTimeSynchronization {
	if casted, ok := structType.(BACnetUnconfirmedServiceRequestUTCTimeSynchronization); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetUnconfirmedServiceRequestUTCTimeSynchronization); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) GetTypeName() string {
	return "BACnetUnconfirmedServiceRequestUTCTimeSynchronization"
}

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (synchronizedDate)
	lengthInBits += m.SynchronizedDate.GetLengthInBits()

	// Simple field (synchronizedTime)
	lengthInBits += m.SynchronizedTime.GetLengthInBits()

	return lengthInBits
}

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestUTCTimeSynchronizationParse(readBuffer utils.ReadBuffer, serviceRequestLength uint16) (BACnetUnconfirmedServiceRequestUTCTimeSynchronization, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetUnconfirmedServiceRequestUTCTimeSynchronization"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetUnconfirmedServiceRequestUTCTimeSynchronization")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (synchronizedDate)
	if pullErr := readBuffer.PullContext("synchronizedDate"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for synchronizedDate")
	}
	_synchronizedDate, _synchronizedDateErr := BACnetApplicationTagParse(readBuffer)
	if _synchronizedDateErr != nil {
		return nil, errors.Wrap(_synchronizedDateErr, "Error parsing 'synchronizedDate' field")
	}
	synchronizedDate := _synchronizedDate.(BACnetApplicationTagDate)
	if closeErr := readBuffer.CloseContext("synchronizedDate"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for synchronizedDate")
	}

	// Simple Field (synchronizedTime)
	if pullErr := readBuffer.PullContext("synchronizedTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for synchronizedTime")
	}
	_synchronizedTime, _synchronizedTimeErr := BACnetApplicationTagParse(readBuffer)
	if _synchronizedTimeErr != nil {
		return nil, errors.Wrap(_synchronizedTimeErr, "Error parsing 'synchronizedTime' field")
	}
	synchronizedTime := _synchronizedTime.(BACnetApplicationTagTime)
	if closeErr := readBuffer.CloseContext("synchronizedTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for synchronizedTime")
	}

	if closeErr := readBuffer.CloseContext("BACnetUnconfirmedServiceRequestUTCTimeSynchronization"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetUnconfirmedServiceRequestUTCTimeSynchronization")
	}

	// Create a partially initialized instance
	_child := &_BACnetUnconfirmedServiceRequestUTCTimeSynchronization{
		SynchronizedDate:                 synchronizedDate,
		SynchronizedTime:                 synchronizedTime,
		_BACnetUnconfirmedServiceRequest: &_BACnetUnconfirmedServiceRequest{},
	}
	_child._BACnetUnconfirmedServiceRequest._BACnetUnconfirmedServiceRequestChildRequirements = _child
	return _child, nil
}

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) Serialize(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetUnconfirmedServiceRequestUTCTimeSynchronization"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetUnconfirmedServiceRequestUTCTimeSynchronization")
		}

		// Simple Field (synchronizedDate)
		if pushErr := writeBuffer.PushContext("synchronizedDate"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for synchronizedDate")
		}
		_synchronizedDateErr := writeBuffer.WriteSerializable(m.GetSynchronizedDate())
		if popErr := writeBuffer.PopContext("synchronizedDate"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for synchronizedDate")
		}
		if _synchronizedDateErr != nil {
			return errors.Wrap(_synchronizedDateErr, "Error serializing 'synchronizedDate' field")
		}

		// Simple Field (synchronizedTime)
		if pushErr := writeBuffer.PushContext("synchronizedTime"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for synchronizedTime")
		}
		_synchronizedTimeErr := writeBuffer.WriteSerializable(m.GetSynchronizedTime())
		if popErr := writeBuffer.PopContext("synchronizedTime"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for synchronizedTime")
		}
		if _synchronizedTimeErr != nil {
			return errors.Wrap(_synchronizedTimeErr, "Error serializing 'synchronizedTime' field")
		}

		if popErr := writeBuffer.PopContext("BACnetUnconfirmedServiceRequestUTCTimeSynchronization"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetUnconfirmedServiceRequestUTCTimeSynchronization")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetUnconfirmedServiceRequestUTCTimeSynchronization) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewBoxedWriteBufferWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
