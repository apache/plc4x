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
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
)

type DefaultPlcSubscriptionEvent struct {
	DefaultResponse
	DefaultPlcSubscriptionEventRequirements `ignore:"true"` // Avoid recursion
	values                                  map[string]*DefaultPlcSubscriptionEventItem
}

type DefaultPlcSubscriptionEventRequirements interface {
	utils.Serializable
	GetAddress(name string) string
}

func NewDefaultPlcSubscriptionEvent(defaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements, tags map[string]model.PlcTag, types map[string]SubscriptionType,
	intervals map[string]time.Duration, responseCodes map[string]model.PlcResponseCode,
	values map[string]values.PlcValue) DefaultPlcSubscriptionEvent {

	valueMap := map[string]*DefaultPlcSubscriptionEventItem{}
	for name, code := range responseCodes {
		tag := tags[name]
		subscriptionType := types[name]
		interval := intervals[name]
		value := values[name]
		valueMap[name] = NewSubscriptionEventItem(code, tag, subscriptionType, interval, value)
	}

	return DefaultPlcSubscriptionEvent{
		DefaultPlcSubscriptionEventRequirements: defaultPlcSubscriptionEventRequirements,
		values:                                  valueMap,
	}
}

func (d *DefaultPlcSubscriptionEvent) GetTagNames() []string {
	var tagNames []string
	for valueName := range d.values {
		tagNames = append(tagNames, valueName)
	}
	return tagNames
}

func (d *DefaultPlcSubscriptionEvent) GetResponseCode(name string) model.PlcResponseCode {
	return d.values[name].GetCode()
}

func (d *DefaultPlcSubscriptionEvent) GetTag(name string) model.PlcTag {
	return d.values[name].GetTag()
}

func (d *DefaultPlcSubscriptionEvent) GetType(name string) SubscriptionType {
	return d.values[name].GetSubscriptionType()
}

func (d *DefaultPlcSubscriptionEvent) GetInterval(name string) time.Duration {
	return d.values[name].GetInterval()
}

func (d *DefaultPlcSubscriptionEvent) GetValue(name string) values.PlcValue {
	return d.values[name].GetValue()
}

func (d *DefaultPlcSubscriptionEvent) GetAddress(name string) string {
	return d.DefaultPlcSubscriptionEventRequirements.GetAddress(name)
}

func (d *DefaultPlcSubscriptionEvent) GetSource(name string) string {
	return d.GetAddress(name)
}

func (d *DefaultPlcSubscriptionEvent) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *DefaultPlcSubscriptionEvent) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcSubscriptionEvent"); err != nil {
		return err
	}
	if err := writeBuffer.PushContext("values", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for name, elem := range d.values {
		_value := fmt.Sprintf("%v", elem)

		if err := writeBuffer.WriteString(name, uint32(len(_value)*8), "UTF-8", _value); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("values", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	if err := writeBuffer.PopContext("PlcSubscriptionEvent"); err != nil {
		return err
	}
	return nil
}

func (d *DefaultPlcSubscriptionEvent) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
