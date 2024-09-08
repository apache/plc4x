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

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	. "github.com/apache/plc4x/plc4go/spi/codegen/fields"
	. "github.com/apache/plc4x/plc4go/spi/codegen/io"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// AirConditioningDataRefresh is the corresponding interface of AirConditioningDataRefresh
type AirConditioningDataRefresh interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	AirConditioningData
	// GetZoneGroup returns ZoneGroup (property field)
	GetZoneGroup() byte
	// IsAirConditioningDataRefresh is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsAirConditioningDataRefresh()
}

// _AirConditioningDataRefresh is the data-structure of this message
type _AirConditioningDataRefresh struct {
	AirConditioningDataContract
	ZoneGroup byte
}

var _ AirConditioningDataRefresh = (*_AirConditioningDataRefresh)(nil)
var _ AirConditioningDataRequirements = (*_AirConditioningDataRefresh)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_AirConditioningDataRefresh) GetParent() AirConditioningDataContract {
	return m.AirConditioningDataContract
}

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for property fields.
///////////////////////

func (m *_AirConditioningDataRefresh) GetZoneGroup() byte {
	return m.ZoneGroup
}

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

// NewAirConditioningDataRefresh factory function for _AirConditioningDataRefresh
func NewAirConditioningDataRefresh(zoneGroup byte, commandTypeContainer AirConditioningCommandTypeContainer) *_AirConditioningDataRefresh {
	_result := &_AirConditioningDataRefresh{
		AirConditioningDataContract: NewAirConditioningData(commandTypeContainer),
		ZoneGroup:                   zoneGroup,
	}
	_result.AirConditioningDataContract.(*_AirConditioningData)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastAirConditioningDataRefresh(structType any) AirConditioningDataRefresh {
	if casted, ok := structType.(AirConditioningDataRefresh); ok {
		return casted
	}
	if casted, ok := structType.(*AirConditioningDataRefresh); ok {
		return *casted
	}
	return nil
}

func (m *_AirConditioningDataRefresh) GetTypeName() string {
	return "AirConditioningDataRefresh"
}

func (m *_AirConditioningDataRefresh) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.AirConditioningDataContract.(*_AirConditioningData).getLengthInBits(ctx))

	// Simple field (zoneGroup)
	lengthInBits += 8

	return lengthInBits
}

func (m *_AirConditioningDataRefresh) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_AirConditioningDataRefresh) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_AirConditioningData) (__airConditioningDataRefresh AirConditioningDataRefresh, err error) {
	m.AirConditioningDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("AirConditioningDataRefresh"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for AirConditioningDataRefresh")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	zoneGroup, err := ReadSimpleField(ctx, "zoneGroup", ReadByte(readBuffer, 8))
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error parsing 'zoneGroup' field"))
	}
	m.ZoneGroup = zoneGroup

	if closeErr := readBuffer.CloseContext("AirConditioningDataRefresh"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for AirConditioningDataRefresh")
	}

	return m, nil
}

func (m *_AirConditioningDataRefresh) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_AirConditioningDataRefresh) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("AirConditioningDataRefresh"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for AirConditioningDataRefresh")
		}

		if err := WriteSimpleField[byte](ctx, "zoneGroup", m.GetZoneGroup(), WriteByte(writeBuffer, 8)); err != nil {
			return errors.Wrap(err, "Error serializing 'zoneGroup' field")
		}

		if popErr := writeBuffer.PopContext("AirConditioningDataRefresh"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for AirConditioningDataRefresh")
		}
		return nil
	}
	return m.AirConditioningDataContract.(*_AirConditioningData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_AirConditioningDataRefresh) IsAirConditioningDataRefresh() {}

func (m *_AirConditioningDataRefresh) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
