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

// BACnetUnconfirmedServiceRequestTimeSynchronization is the corresponding interface of BACnetUnconfirmedServiceRequestTimeSynchronization
type BACnetUnconfirmedServiceRequestTimeSynchronization interface {
	utils.LengthAware
	utils.Serializable
	BACnetUnconfirmedServiceRequest
	// GetSynchronizedDate returns SynchronizedDate (property field)
	GetSynchronizedDate() BACnetApplicationTagDate
	// GetSynchronizedTime returns SynchronizedTime (property field)
	GetSynchronizedTime() BACnetApplicationTagTime
}

// BACnetUnconfirmedServiceRequestTimeSynchronizationExactly can be used when we want exactly this type and not a type which fulfills BACnetUnconfirmedServiceRequestTimeSynchronization.
// This is useful for switch cases.
type BACnetUnconfirmedServiceRequestTimeSynchronizationExactly interface {
	BACnetUnconfirmedServiceRequestTimeSynchronization
	isBACnetUnconfirmedServiceRequestTimeSynchronization() bool
}

// _BACnetUnconfirmedServiceRequestTimeSynchronization is the data-structure of this message
type _BACnetUnconfirmedServiceRequestTimeSynchronization struct {
	*_BACnetUnconfirmedServiceRequest
	SynchronizedDate BACnetApplicationTagDate
	SynchronizedTime BACnetApplicationTagTime
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) GetServiceChoice() BACnetUnconfirmedServiceChoice {
	return BACnetUnconfirmedServiceChoice_TIME_SYNCHRONIZATION
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) InitializeParent(parent BACnetUnconfirmedServiceRequest) {
}

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) GetParent() BACnetUnconfirmedServiceRequest {
	return m._BACnetUnconfirmedServiceRequest
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) GetSynchronizedDate() BACnetApplicationTagDate {
	return m.SynchronizedDate
}

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) GetSynchronizedTime() BACnetApplicationTagTime {
	return m.SynchronizedTime
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewBACnetUnconfirmedServiceRequestTimeSynchronization factory function for _BACnetUnconfirmedServiceRequestTimeSynchronization
func NewBACnetUnconfirmedServiceRequestTimeSynchronization(synchronizedDate BACnetApplicationTagDate, synchronizedTime BACnetApplicationTagTime, serviceRequestLength uint16) *_BACnetUnconfirmedServiceRequestTimeSynchronization {
	_result := &_BACnetUnconfirmedServiceRequestTimeSynchronization{
		SynchronizedDate:                 synchronizedDate,
		SynchronizedTime:                 synchronizedTime,
		_BACnetUnconfirmedServiceRequest: NewBACnetUnconfirmedServiceRequest(serviceRequestLength),
	}
	_result._BACnetUnconfirmedServiceRequest._BACnetUnconfirmedServiceRequestChildRequirements = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastBACnetUnconfirmedServiceRequestTimeSynchronization(structType interface{}) BACnetUnconfirmedServiceRequestTimeSynchronization {
	if casted, ok := structType.(BACnetUnconfirmedServiceRequestTimeSynchronization); ok {
		return casted
	}
	if casted, ok := structType.(*BACnetUnconfirmedServiceRequestTimeSynchronization); ok {
		return *casted
	}
	return nil
}

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) GetTypeName() string {
	return "BACnetUnconfirmedServiceRequestTimeSynchronization"
}

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) GetLengthInBits() uint16 {
	return m.GetLengthInBitsConditional(false)
}

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) GetLengthInBitsConditional(lastItem bool) uint16 {
	lengthInBits := uint16(m.GetParentLengthInBits())

	// Simple field (synchronizedDate)
	lengthInBits += m.SynchronizedDate.GetLengthInBits()

	// Simple field (synchronizedTime)
	lengthInBits += m.SynchronizedTime.GetLengthInBits()

	return lengthInBits
}

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) GetLengthInBytes() uint16 {
	return m.GetLengthInBits() / 8
}

func BACnetUnconfirmedServiceRequestTimeSynchronizationParse(theBytes []byte, serviceRequestLength uint16) (BACnetUnconfirmedServiceRequestTimeSynchronization, error) {
	return BACnetUnconfirmedServiceRequestTimeSynchronizationParseWithBuffer(utils.NewReadBufferByteBased(theBytes), serviceRequestLength)
}

func BACnetUnconfirmedServiceRequestTimeSynchronizationParseWithBuffer(readBuffer utils.ReadBuffer, serviceRequestLength uint16) (BACnetUnconfirmedServiceRequestTimeSynchronization, error) {
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("BACnetUnconfirmedServiceRequestTimeSynchronization"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for BACnetUnconfirmedServiceRequestTimeSynchronization")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	// Simple Field (synchronizedDate)
	if pullErr := readBuffer.PullContext("synchronizedDate"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for synchronizedDate")
	}
	_synchronizedDate, _synchronizedDateErr := BACnetApplicationTagParseWithBuffer(readBuffer)
	if _synchronizedDateErr != nil {
		return nil, errors.Wrap(_synchronizedDateErr, "Error parsing 'synchronizedDate' field of BACnetUnconfirmedServiceRequestTimeSynchronization")
	}
	synchronizedDate := _synchronizedDate.(BACnetApplicationTagDate)
	if closeErr := readBuffer.CloseContext("synchronizedDate"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for synchronizedDate")
	}

	// Simple Field (synchronizedTime)
	if pullErr := readBuffer.PullContext("synchronizedTime"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for synchronizedTime")
	}
	_synchronizedTime, _synchronizedTimeErr := BACnetApplicationTagParseWithBuffer(readBuffer)
	if _synchronizedTimeErr != nil {
		return nil, errors.Wrap(_synchronizedTimeErr, "Error parsing 'synchronizedTime' field of BACnetUnconfirmedServiceRequestTimeSynchronization")
	}
	synchronizedTime := _synchronizedTime.(BACnetApplicationTagTime)
	if closeErr := readBuffer.CloseContext("synchronizedTime"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for synchronizedTime")
	}

	if closeErr := readBuffer.CloseContext("BACnetUnconfirmedServiceRequestTimeSynchronization"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for BACnetUnconfirmedServiceRequestTimeSynchronization")
	}

	// Create a partially initialized instance
	_child := &_BACnetUnconfirmedServiceRequestTimeSynchronization{
		_BACnetUnconfirmedServiceRequest: &_BACnetUnconfirmedServiceRequest{
			ServiceRequestLength: serviceRequestLength,
		},
		SynchronizedDate: synchronizedDate,
		SynchronizedTime: synchronizedTime,
	}
	_child._BACnetUnconfirmedServiceRequest._BACnetUnconfirmedServiceRequestChildRequirements = _child
	return _child, nil
}

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes())))
	if err := m.SerializeWithWriteBuffer(wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) SerializeWithWriteBuffer(writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	ser := func() error {
		if pushErr := writeBuffer.PushContext("BACnetUnconfirmedServiceRequestTimeSynchronization"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for BACnetUnconfirmedServiceRequestTimeSynchronization")
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

		if popErr := writeBuffer.PopContext("BACnetUnconfirmedServiceRequestTimeSynchronization"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for BACnetUnconfirmedServiceRequestTimeSynchronization")
		}
		return nil
	}
	return m.SerializeParent(writeBuffer, m, ser)
}

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) isBACnetUnconfirmedServiceRequestTimeSynchronization() bool {
	return true
}

func (m *_BACnetUnconfirmedServiceRequestTimeSynchronization) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
