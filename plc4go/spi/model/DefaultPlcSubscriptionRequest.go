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
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/spi"
	"github.com/pkg/errors"
	"time"
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
	fieldHandler           spi.PlcFieldHandler
	valueHandler           spi.PlcValueHandler
	queries                map[string]string
	queryNames             []string
	fields                 map[string]model.PlcField
	fieldNames             []string
	types                  map[string]SubscriptionType
	intervals              map[string]time.Duration
	preRegisteredConsumers map[string][]model.PlcSubscriptionEventConsumer `ignore:"true"`
}

func NewDefaultPlcSubscriptionRequestBuilder(fieldHandler spi.PlcFieldHandler, valueHandler spi.PlcValueHandler, subscriber spi.PlcSubscriber) *DefaultPlcSubscriptionRequestBuilder {
	return &DefaultPlcSubscriptionRequestBuilder{
		subscriber:             subscriber,
		fieldHandler:           fieldHandler,
		valueHandler:           valueHandler,
		queries:                map[string]string{},
		fields:                 map[string]model.PlcField{},
		fieldNames:             make([]string, 0),
		types:                  map[string]SubscriptionType{},
		intervals:              map[string]time.Duration{},
		preRegisteredConsumers: make(map[string][]model.PlcSubscriptionEventConsumer),
	}
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddCyclicQuery(name string, query string, interval time.Duration) model.PlcSubscriptionRequestBuilder {
	d.queryNames = append(d.queryNames, name)
	d.queries[name] = query
	d.types[name] = SubscriptionCyclic
	d.intervals[name] = interval
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddCyclicField(name string, field model.PlcField, interval time.Duration) model.PlcSubscriptionRequestBuilder {
	d.fieldNames = append(d.fieldNames, name)
	d.fields[name] = field
	d.types[name] = SubscriptionCyclic
	d.intervals[name] = interval
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddChangeOfStateQuery(name string, query string) model.PlcSubscriptionRequestBuilder {
	d.queryNames = append(d.queryNames, name)
	d.queries[name] = query
	d.types[name] = SubscriptionChangeOfState
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddChangeOfStateField(name string, field model.PlcField) model.PlcSubscriptionRequestBuilder {
	d.fieldNames = append(d.fieldNames, name)
	d.fields[name] = field
	d.types[name] = SubscriptionChangeOfState
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddEventQuery(name string, query string) model.PlcSubscriptionRequestBuilder {
	d.queryNames = append(d.queryNames, name)
	d.queries[name] = query
	d.types[name] = SubscriptionEvent
	return d
}

func (d *DefaultPlcSubscriptionRequestBuilder) AddEventField(name string, field model.PlcField) model.PlcSubscriptionRequestBuilder {
	d.fieldNames = append(d.fieldNames, name)
	d.fields[name] = field
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
	for _, name := range d.queryNames {
		query := d.queries[name]
		field, err := d.fieldHandler.ParseQuery(query)
		if err != nil {
			return nil, errors.Wrapf(err, "Error parsing query: %s", query)
		}
		d.fieldNames = append(d.fieldNames, name)
		d.fields[name] = field
	}
	return NewDefaultPlcSubscriptionRequest(d.fields, d.fieldNames, d.types, d.intervals, d.subscriber, d.preRegisteredConsumers), nil
}

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcSubscriptionRequest
type DefaultPlcSubscriptionRequest struct {
	DefaultRequest
	types                  map[string]SubscriptionType
	intervals              map[string]time.Duration
	subscriber             spi.PlcSubscriber
	preRegisteredConsumers map[string][]model.PlcSubscriptionEventConsumer `ignore:"true"`
}

func NewDefaultPlcSubscriptionRequest(fields map[string]model.PlcField, fieldNames []string, types map[string]SubscriptionType, intervals map[string]time.Duration, subscriber spi.PlcSubscriber, preRegisteredConsumers map[string][]model.PlcSubscriptionEventConsumer) model.PlcSubscriptionRequest {
	return &DefaultPlcSubscriptionRequest{NewDefaultRequest(fields, fieldNames), types, intervals, subscriber, preRegisteredConsumers}
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
