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

// Code generated by "plc4xGenerator -type=Subscriber"; DO NOT EDIT.

package bacnetip

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

var _ = fmt.Printf

func (d *Subscriber) Serialize() ([]byte, error) {
	if d == nil {
		return nil, fmt.Errorf("(*DeviceInfoCache)(nil)")
	}
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *Subscriber) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if d == nil {
		return fmt.Errorf("(*DeviceInfoCache)(nil)")
	}
	if err := writeBuffer.PushContext("Subscriber"); err != nil {
		return err
	}
	if d.connection != nil {
		{
			_value := fmt.Sprintf("%v", d.connection)

			if err := writeBuffer.WriteString("connection", uint32(len(_value)*8), _value); err != nil {
				return err
			}
		}
	}
	if err := writeBuffer.PushContext("consumers", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for _name, elem := range d.consumers {
		name := fmt.Sprintf("%v", &_name)

		var elem any = elem
		if serializable, ok := elem.(utils.Serializable); ok {
			if err := writeBuffer.PushContext(name); err != nil {
				return err
			}
			if err := serializable.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
			if err := writeBuffer.PopContext(name); err != nil {
				return err
			}
		} else {
			elemAsString := fmt.Sprintf("%v", elem)
			if err := writeBuffer.WriteString(name, uint32(len(elemAsString)*8), elemAsString); err != nil {
				return err
			}
		}
	}
	if err := writeBuffer.PopContext("consumers", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	if err := writeBuffer.PushContext("_options", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for _, elem := range d._options {
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
	if err := writeBuffer.PopContext("_options", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("Subscriber"); err != nil {
		return err
	}
	return nil
}

func (d *Subscriber) String() string {
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
