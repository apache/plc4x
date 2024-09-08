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

	"github.com/apache/plc4x/plc4go/spi/utils"
)

// Code generated by code-generation. DO NOT EDIT.

// TriggerControlDataTriggerMax is the corresponding interface of TriggerControlDataTriggerMax
type TriggerControlDataTriggerMax interface {
	fmt.Stringer
	utils.LengthAware
	utils.Serializable
	TriggerControlData
	// IsTriggerControlDataTriggerMax is a marker method to prevent unintentional type checks (interfaces of same signature)
	IsTriggerControlDataTriggerMax()
}

// _TriggerControlDataTriggerMax is the data-structure of this message
type _TriggerControlDataTriggerMax struct {
	TriggerControlDataContract
}

var _ TriggerControlDataTriggerMax = (*_TriggerControlDataTriggerMax)(nil)
var _ TriggerControlDataRequirements = (*_TriggerControlDataTriggerMax)(nil)

///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
/////////////////////// Accessors for discriminator values.
///////////////////////

///////////////////////
///////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////

func (m *_TriggerControlDataTriggerMax) GetParent() TriggerControlDataContract {
	return m.TriggerControlDataContract
}

// NewTriggerControlDataTriggerMax factory function for _TriggerControlDataTriggerMax
func NewTriggerControlDataTriggerMax(commandTypeContainer TriggerControlCommandTypeContainer, triggerGroup byte) *_TriggerControlDataTriggerMax {
	_result := &_TriggerControlDataTriggerMax{
		TriggerControlDataContract: NewTriggerControlData(commandTypeContainer, triggerGroup),
	}
	_result.TriggerControlDataContract.(*_TriggerControlData)._SubType = _result
	return _result
}

// Deprecated: use the interface for direct cast
func CastTriggerControlDataTriggerMax(structType any) TriggerControlDataTriggerMax {
	if casted, ok := structType.(TriggerControlDataTriggerMax); ok {
		return casted
	}
	if casted, ok := structType.(*TriggerControlDataTriggerMax); ok {
		return *casted
	}
	return nil
}

func (m *_TriggerControlDataTriggerMax) GetTypeName() string {
	return "TriggerControlDataTriggerMax"
}

func (m *_TriggerControlDataTriggerMax) GetLengthInBits(ctx context.Context) uint16 {
	lengthInBits := uint16(m.TriggerControlDataContract.(*_TriggerControlData).getLengthInBits(ctx))

	return lengthInBits
}

func (m *_TriggerControlDataTriggerMax) GetLengthInBytes(ctx context.Context) uint16 {
	return m.GetLengthInBits(ctx) / 8
}

func (m *_TriggerControlDataTriggerMax) parse(ctx context.Context, readBuffer utils.ReadBuffer, parent *_TriggerControlData) (__triggerControlDataTriggerMax TriggerControlDataTriggerMax, err error) {
	m.TriggerControlDataContract = parent
	parent._SubType = m
	positionAware := readBuffer
	_ = positionAware
	if pullErr := readBuffer.PullContext("TriggerControlDataTriggerMax"); pullErr != nil {
		return nil, errors.Wrap(pullErr, "Error pulling for TriggerControlDataTriggerMax")
	}
	currentPos := positionAware.GetPos()
	_ = currentPos

	if closeErr := readBuffer.CloseContext("TriggerControlDataTriggerMax"); closeErr != nil {
		return nil, errors.Wrap(closeErr, "Error closing for TriggerControlDataTriggerMax")
	}

	return m, nil
}

func (m *_TriggerControlDataTriggerMax) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithInitialSizeForByteBasedBuffer(int(m.GetLengthInBytes(context.Background()))))
	if err := m.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (m *_TriggerControlDataTriggerMax) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	positionAware := writeBuffer
	_ = positionAware
	log := zerolog.Ctx(ctx)
	_ = log
	ser := func() error {
		if pushErr := writeBuffer.PushContext("TriggerControlDataTriggerMax"); pushErr != nil {
			return errors.Wrap(pushErr, "Error pushing for TriggerControlDataTriggerMax")
		}

		if popErr := writeBuffer.PopContext("TriggerControlDataTriggerMax"); popErr != nil {
			return errors.Wrap(popErr, "Error popping for TriggerControlDataTriggerMax")
		}
		return nil
	}
	return m.TriggerControlDataContract.(*_TriggerControlData).serializeParent(ctx, writeBuffer, m, ser)
}

func (m *_TriggerControlDataTriggerMax) IsTriggerControlDataTriggerMax() {}

func (m *_TriggerControlDataTriggerMax) String() string {
	if m == nil {
		return "<nil>"
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), m); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
