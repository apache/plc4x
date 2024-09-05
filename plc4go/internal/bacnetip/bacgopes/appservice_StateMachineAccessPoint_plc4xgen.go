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

// Code generated by "plc4xGenerator -type=StateMachineAccessPoint -prefix=appservice_"; DO NOT EDIT.

package bacgopes

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

var _ = fmt.Printf

func (d *StateMachineAccessPoint) Serialize() ([]byte, error) {
	if d == nil {
		return nil, fmt.Errorf("(*DeviceInfoCache)(nil)")
	}
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *StateMachineAccessPoint) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if d == nil {
		return fmt.Errorf("(*DeviceInfoCache)(nil)")
	}
	if err := writeBuffer.PushContext("StateMachineAccessPoint"); err != nil {
		return err
	}
	if err := d.Client.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}
	if err := d.ServiceAccessPointContract.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}
	if d.localDevice != nil {
		{
			_value := fmt.Sprintf("%v", d.localDevice)

			if err := writeBuffer.WriteString("localDevice", uint32(len(_value)*8), _value); err != nil {
				return err
			}
		}
	}
	if d.deviceInfoCache != nil {
		{
			_value := fmt.Sprintf("%v", d.deviceInfoCache)

			if err := writeBuffer.WriteString("deviceInfoCache", uint32(len(_value)*8), _value); err != nil {
				return err
			}
		}
	}
	{
		_value := fmt.Sprintf("%v", d.nextInvokeId)

		if err := writeBuffer.WriteString("nextInvokeId", uint32(len(_value)*8), _value); err != nil {
			return err
		}
	}
	if err := writeBuffer.PushContext("clientTransactions", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for _, elem := range d.clientTransactions {
		var elem any = elem

		if elem != nil {
			if serializableField, ok := any(elem).(utils.Serializable); ok {
				if err := writeBuffer.PushContext("value"); err != nil {
					return err
				}
				if err := serializableField.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
					return err
				}
				if err := writeBuffer.PopContext("value"); err != nil {
					return err
				}
			} else {
				stringValue := fmt.Sprintf("%v", elem)
				if err := writeBuffer.WriteString("value", uint32(len(stringValue)*8), stringValue); err != nil {
					return err
				}
			}
		}
	}
	if err := writeBuffer.PopContext("clientTransactions", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	if err := writeBuffer.PushContext("serverTransactions", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for _, elem := range d.serverTransactions {
		var elem any = elem

		if elem != nil {
			if serializableField, ok := any(elem).(utils.Serializable); ok {
				if err := writeBuffer.PushContext("value"); err != nil {
					return err
				}
				if err := serializableField.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
					return err
				}
				if err := writeBuffer.PopContext("value"); err != nil {
					return err
				}
			} else {
				stringValue := fmt.Sprintf("%v", elem)
				if err := writeBuffer.WriteString("value", uint32(len(stringValue)*8), stringValue); err != nil {
					return err
				}
			}
		}
	}
	if err := writeBuffer.PopContext("serverTransactions", utils.WithRenderAsList(true)); err != nil {
		return err
	}

	if err := writeBuffer.WriteInt64("numberOfApduRetries", 64, int64(d.numberOfApduRetries)); err != nil {
		return err
	}
	{
		_value := fmt.Sprintf("%v", d.apduTimeout)

		if err := writeBuffer.WriteString("apduTimeout", uint32(len(_value)*8), _value); err != nil {
			return err
		}
	}

	if err := writeBuffer.WriteString("maxApduLengthAccepted", uint32(len(d.maxApduLengthAccepted.String())*8), d.maxApduLengthAccepted.String()); err != nil {
		return err
	}

	if err := writeBuffer.WriteString("segmentationSupported", uint32(len(d.segmentationSupported.String())*8), d.segmentationSupported.String()); err != nil {
		return err
	}
	{
		_value := fmt.Sprintf("%v", d.segmentTimeout)

		if err := writeBuffer.WriteString("segmentTimeout", uint32(len(_value)*8), _value); err != nil {
			return err
		}
	}

	if err := writeBuffer.WriteString("maxSegmentsAccepted", uint32(len(d.maxSegmentsAccepted.String())*8), d.maxSegmentsAccepted.String()); err != nil {
		return err
	}
	{
		_value := fmt.Sprintf("%v", d.proposedWindowSize)

		if err := writeBuffer.WriteString("proposedWindowSize", uint32(len(_value)*8), _value); err != nil {
			return err
		}
	}

	if err := writeBuffer.WriteString("dccEnableDisable", uint32(len(d.dccEnableDisable.String())*8), d.dccEnableDisable.String()); err != nil {
		return err
	}
	{
		_value := fmt.Sprintf("%v", d.applicationTimeout)

		if err := writeBuffer.WriteString("applicationTimeout", uint32(len(_value)*8), _value); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("StateMachineAccessPoint"); err != nil {
		return err
	}
	return nil
}

func (d *StateMachineAccessPoint) String() string {
	if alternateStringer, ok := any(d).(utils.AlternateStringer); ok {
		if alternateString, use := alternateStringer.AlternateString(); use {
			return alternateString
		}
	}
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
