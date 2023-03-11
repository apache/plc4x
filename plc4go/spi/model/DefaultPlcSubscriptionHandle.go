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
	"encoding/binary"
	"fmt"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DefaultPlcSubscriptionHandle struct {
	handleToRegister model.PlcSubscriptionHandle `ignore:"true"`
	plcSubscriber    spi.PlcSubscriber
}

// NewDefaultPlcSubscriptionHandle can be used when the DefaultPlcSubscriptionHandle is sufficient
func NewDefaultPlcSubscriptionHandle(plcSubscriber spi.PlcSubscriber) *DefaultPlcSubscriptionHandle {
	handle := &DefaultPlcSubscriptionHandle{
		plcSubscriber: plcSubscriber,
	}
	handle.handleToRegister = handle
	return handle
}

// NewDefaultPlcSubscriptionHandleWithHandleToRegister should be used when an extension of DefaultPlcSubscriptionHandle is used
func NewDefaultPlcSubscriptionHandleWithHandleToRegister(plcSubscriber spi.PlcSubscriber, handleToRegister model.PlcSubscriptionHandle) *DefaultPlcSubscriptionHandle {
	return &DefaultPlcSubscriptionHandle{
		handleToRegister: handleToRegister,
		plcSubscriber:    plcSubscriber,
	}
}

// Register registers at the spi.PlcSubscriber
func (d *DefaultPlcSubscriptionHandle) Register(consumer model.PlcSubscriptionEventConsumer) model.PlcConsumerRegistration {
	return d.plcSubscriber.Register(consumer, []model.PlcSubscriptionHandle{d.handleToRegister})
}

func (d *DefaultPlcSubscriptionHandle) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *DefaultPlcSubscriptionHandle) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcSubscriptionHandle"); err != nil {
		return err
	}

	if d.plcSubscriber != nil {
		if serializableField, ok := d.plcSubscriber.(utils.Serializable); ok {
			if err := writeBuffer.PushContext("plcSubscriber"); err != nil {
				return err
			}
			if err := serializableField.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
			if err := writeBuffer.PopContext("plcSubscriber"); err != nil {
				return err
			}
		} else {
			stringValue := fmt.Sprintf("%v", d.plcSubscriber)
			if err := writeBuffer.WriteString("plcSubscriber", uint32(len(stringValue)*8), "UTF-8", stringValue); err != nil {
				return err
			}
		}
	}
	if err := writeBuffer.PopContext("PlcSubscriptionHandle"); err != nil {
		return err
	}
	return nil
}

func (d *DefaultPlcSubscriptionHandle) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
