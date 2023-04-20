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
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"github.com/pkg/errors"
)

type SubscriptionType uint8

const (
	SubscriptionCyclic        SubscriptionType = 0x01
	SubscriptionChangeOfState SubscriptionType = 0x02
	SubscriptionEvent         SubscriptionType = 0x03
)

func (s SubscriptionType) String() string {
	switch s {
	case SubscriptionCyclic:
		return "SubscriptionCyclic"
	case SubscriptionChangeOfState:
		return "SubscriptionChangeOfState"
	case SubscriptionEvent:
		return "SubscriptionEvent"
	default:
		return "Unknown"
	}
}

type DefaultPlcSubscriptionRequestBuilder struct {
	subscriber             spi.PlcSubscriber
	tagHandler             spi.PlcTagHandler
	valueHandler           spi.PlcValueHandler
	tagNames               []string
	tagAddresses           map[string]string
	tags                   map[string]model.PlcTag
	types                  map[string]SubscriptionType
	intervals              map[string]time.Duration
	preRegisteredConsumers map[string][]model.PlcSubscriptionEventConsumer `ignore:"true"`
}

func NewDefaultPlcSubscriptionRequestBuilder(tagHandler spi.PlcTagHandler, valueHandler spi.PlcValueHandler, subscriber spi.PlcSubscriber) *DefaultPlcSubscriptionRequestBuilder {
	return &DefaultPlcSubscriptionRequestBuilder{
		subscriber:             subscriber,
		tagHandler:             tagHandler,
		valueHandler:           valueHandler,
		tagNames:               make([]string, 0),
		tagAddresses:           map[string]string{},
		tags:                   map[string]model.PlcTag{},
		types:                  map[string]SubscriptionType{},
		intervals:              map[string]time.Duration{},
		preRegisteredConsumers: make(map[string][]model.PlcSubscriptionEventConsumer),
	}
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddCyclicTagAddress(name string, tagAddress string, interval time.Duration) model.PlcSubscriptionRequestBuilder {
	d.tagNames = append(d.tagNames, name)
	d.tagAddresses[name] = tagAddress
	d.types[name] = SubscriptionCyclic
	d.intervals[name] = interval
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddCyclicTag(name string, tag model.PlcTag, interval time.Duration) model.PlcSubscriptionRequestBuilder {
	d.tagNames = append(d.tagNames, name)
	d.tags[name] = tag
	d.types[name] = SubscriptionCyclic
	d.intervals[name] = interval
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddChangeOfStateTagAddress(name string, tagAddress string) model.PlcSubscriptionRequestBuilder {
	d.tagNames = append(d.tagNames, name)
	d.tagAddresses[name] = tagAddress
	d.types[name] = SubscriptionChangeOfState
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddChangeOfStateTag(name string, tag model.PlcTag) model.PlcSubscriptionRequestBuilder {
	d.tagNames = append(d.tagNames, name)
	d.tags[name] = tag
	d.types[name] = SubscriptionChangeOfState
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddEventTagAddress(name string, tagAddress string) model.PlcSubscriptionRequestBuilder {
	d.tagNames = append(d.tagNames, name)
	d.tagAddresses[name] = tagAddress
	d.types[name] = SubscriptionEvent
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddEventTag(name string, tag model.PlcTag) model.PlcSubscriptionRequestBuilder {
	d.tagNames = append(d.tagNames, name)
	d.tags[name] = tag
	d.types[name] = SubscriptionEvent
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddPreRegisteredConsumer(name string, consumer model.PlcSubscriptionEventConsumer) model.PlcSubscriptionRequestBuilder {
	if d.preRegisteredConsumers[name] == nil {
		d.preRegisteredConsumers[name] = make([]model.PlcSubscriptionEventConsumer, 0)
	}
	d.preRegisteredConsumers[name] = append(d.preRegisteredConsumers[name], consumer)
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) Build() (model.PlcSubscriptionRequest, error) {
	for _, name := range d.tagNames {
		if tagAddress, ok := d.tagAddresses[name]; ok {
			tag, err := d.tagHandler.ParseTag(tagAddress)
			if err != nil {
				return nil, errors.Wrapf(err, "Error parsing tag query: %s", tagAddress)
			}
			d.tags[name] = tag
		}
	}
	return NewDefaultPlcSubscriptionRequest(d.subscriber, d.tagNames, d.tags, d.types, d.intervals, d.preRegisteredConsumers), nil
}

type DefaultPlcSubscriptionRequest struct {
	DefaultPlcTagRequest
	types                  map[string]SubscriptionType
	intervals              map[string]time.Duration
	preRegisteredConsumers map[string][]model.PlcSubscriptionEventConsumer `ignore:"true"`
	subscriber             spi.PlcSubscriber
}

func NewDefaultPlcSubscriptionRequest(subscriber spi.PlcSubscriber, tagNames []string, tags map[string]model.PlcTag, types map[string]SubscriptionType, intervals map[string]time.Duration, preRegisteredConsumers map[string][]model.PlcSubscriptionEventConsumer) model.PlcSubscriptionRequest {
	return &DefaultPlcSubscriptionRequest{NewDefaultPlcTagRequest(tags, tagNames), types, intervals, preRegisteredConsumers, subscriber}
}

func (d *DefaultPlcSubscriptionRequest) Execute() <-chan model.PlcSubscriptionRequestResult {
	return d.ExecuteWithContext(context.TODO())
}

func (d *DefaultPlcSubscriptionRequest) ExecuteWithContext(ctx context.Context) <-chan model.PlcSubscriptionRequestResult {
	return d.subscriber.Subscribe(ctx, d)
}

func (d *DefaultPlcSubscriptionRequest) GetType(name string) SubscriptionType {
	return d.types[name]
}

func (d *DefaultPlcSubscriptionRequest) GetInterval(name string) time.Duration {
	return d.intervals[name]
}

func (d *DefaultPlcSubscriptionRequest) GetPreRegisteredConsumers(name string) []model.PlcSubscriptionEventConsumer {
	return d.preRegisteredConsumers[name]
}

func (d *DefaultPlcSubscriptionRequest) Serialize() ([]byte, error) {
	wb := utils.NewWriteBufferByteBased(utils.WithByteOrderForByteBasedBuffer(binary.BigEndian))
	if err := d.SerializeWithWriteBuffer(context.Background(), wb); err != nil {
		return nil, err
	}
	return wb.GetBytes(), nil
}

func (d *DefaultPlcSubscriptionRequest) SerializeWithWriteBuffer(ctx context.Context, writeBuffer utils.WriteBuffer) error {
	if err := writeBuffer.PushContext("PlcSubscriptionRequest"); err != nil {
		return err
	}
	if err := d.DefaultPlcTagRequest.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
		return err
	}
	if err := writeBuffer.PushContext("types", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for name, elem := range d.types {
		_value := fmt.Sprintf("%v", elem)

		if err := writeBuffer.WriteString(name, uint32(len(_value)*8), "UTF-8", _value); err != nil {
			return err
		}
	}
	if err := writeBuffer.PopContext("types", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	if err := writeBuffer.PushContext("intervals", utils.WithRenderAsList(true)); err != nil {
		return err
	}
	for name, elem := range d.intervals {

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
			if err := writeBuffer.WriteString(name, uint32(len(elemAsString)*8), "UTF-8", elemAsString); err != nil {
				return err
			}
		}
	}
	if err := writeBuffer.PopContext("intervals", utils.WithRenderAsList(true)); err != nil {
		return err
	}

	if d.subscriber != nil {
		if serializableField, ok := d.subscriber.(utils.Serializable); ok {
			if err := writeBuffer.PushContext("subscriber"); err != nil {
				return err
			}
			if err := serializableField.SerializeWithWriteBuffer(ctx, writeBuffer); err != nil {
				return err
			}
			if err := writeBuffer.PopContext("subscriber"); err != nil {
				return err
			}
		} else {
			stringValue := fmt.Sprintf("%v", d.subscriber)
			if err := writeBuffer.WriteString("subscriber", uint32(len(stringValue)*8), "UTF-8", stringValue); err != nil {
				return err
			}
		}
	}
	if err := writeBuffer.PopContext("PlcSubscriptionRequest"); err != nil {
		return err
	}
	return nil
}

func (d *DefaultPlcSubscriptionRequest) String() string {
	writeBuffer := utils.NewWriteBufferBoxBasedWithOptions(true, true)
	if err := writeBuffer.WriteSerializable(context.Background(), d); err != nil {
		return err.Error()
	}
	return writeBuffer.GetBox().String()
}
