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
	"math/rand"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DefaultPlcConsumerRegistration struct {
	consumerId    int
	consumer      model.PlcSubscriptionEventConsumer `ignore:"true"` // Function not renderable
	plcSubscriber spi.PlcSubscriber
	handles       []model.PlcSubscriptionHandle
}

func NewDefaultPlcConsumerRegistration(plcSubscriber spi.PlcSubscriber, consumer model.PlcSubscriptionEventConsumer, handles ...model.PlcSubscriptionHandle) *DefaultPlcConsumerRegistration {
	return &DefaultPlcConsumerRegistration{
		// TODO: we need a way to hash the consumer
		consumerId:    rand.Int(),
		consumer:      consumer,
		plcSubscriber: plcSubscriber,
		handles:       handles,
	}
}

func (d *DefaultPlcConsumerRegistration) GetConsumerId() int {
	return d.consumerId
}

func (d *DefaultPlcConsumerRegistration) GetSubscriptionHandles() []model.PlcSubscriptionHandle {
	return d.handles
}

func (d *DefaultPlcConsumerRegistration) Unregister() {
	d.plcSubscriber.Unregister(d)
}

func (d *DefaultPlcConsumerRegistration) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *DefaultPlcConsumerRegistration) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcConsumerRegistration"); err != nil {
		return err
	}
	_value := fmt.Sprintf("%v", d.consumerId)

	if err := writeBuffer.WriteString("consumerId", uint32(len(_value)*8), "UTF-8", _value); err != nil {
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
	if err := writeBuffer.PushContext("handles", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for _, elem := range d.handles {
		var elem any = elem

		if elem != nil {
			if serializableField, ok := elem.(utils.Serializable); ok {
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
				if err := writeBuffer.WriteString("value", uint32(len(stringValue)*8), "UTF-8", stringValue); err != nil {
					return err
				}
			}
		}
	}
	if err := writeBuffer.PopContext("handles", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("PlcConsumerRegistration"); err != nil {
		return err
	}
	return nil
}

func (d *DefaultPlcConsumerRegistration) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
