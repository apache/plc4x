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
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
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

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcSubscriptionRequestBuilder
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
	return NewDefaultPlcSubscriptionRequest(d.tags, d.tagNames, d.types, d.intervals, d.subscriber, d.preRegisteredConsumers), nil
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcSubscriptionRequest
type DefaultPlcSubscriptionRequest struct {
	DefaultPlcTagRequest
	types                  map[string]SubscriptionType
	intervals              map[string]time.Duration
	subscriber             spi.PlcSubscriber
	preRegisteredConsumers map[string][]model.PlcSubscriptionEventConsumer `ignore:"true"`
}

func NewDefaultPlcSubscriptionRequest(tags map[string]model.PlcTag, tagNames []string, types map[string]SubscriptionType, intervals map[string]time.Duration, subscriber spi.PlcSubscriber, preRegisteredConsumers map[string][]model.PlcSubscriptionEventConsumer) model.PlcSubscriptionRequest {
	return &DefaultPlcSubscriptionRequest{NewDefaultPlcTagRequest(tags, tagNames), types, intervals, subscriber, preRegisteredConsumers}
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
