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
	"github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
	"time"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcSubscriptionEvent
type DefaultPlcSubscriptionEvent struct {
	DefaultResponse
	DefaultPlcSubscriptionEventRequirements `ignore:"true"` // Avoid recursion
	fields                                  map[string]model.PlcField
	types                                   map[string]SubscriptionType
	intervals                               map[string]time.Duration
	values                                  map[string]values.PlcValue
}

type DefaultPlcSubscriptionEventRequirements interface {
	utils.Serializable
	GetAddress(name string) string
}

func NewDefaultPlcSubscriptionEvent(defaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements, fields map[string]model.PlcField, types map[string]SubscriptionType,
	intervals map[string]time.Duration, responseCodes map[string]model.PlcResponseCode,
	values map[string]values.PlcValue) DefaultPlcSubscriptionEvent {
	return DefaultPlcSubscriptionEvent{
		DefaultResponse:                         NewDefaultResponse(responseCodes),
		DefaultPlcSubscriptionEventRequirements: defaultPlcSubscriptionEventRequirements,
		fields:                                  fields,
		types:                                   types,
		intervals:                               intervals,
		values:                                  values,
	}
}

func (d *DefaultPlcSubscriptionEvent) GetFieldNames() []string {
	var fieldNames []string
	for fieldName := range d.fields {
		fieldNames = append(fieldNames, fieldName)
	}
	return fieldNames
}

func (d *DefaultPlcSubscriptionEvent) GetField(name string) model.PlcField {
	return d.fields[name]
}

func (d *DefaultPlcSubscriptionEvent) GetType(name string) SubscriptionType {
	return d.types[name]
}

func (d *DefaultPlcSubscriptionEvent) GetInterval(name string) time.Duration {
	return d.intervals[name]
}

func (d *DefaultPlcSubscriptionEvent) GetAddress(name string) string {
	return d.DefaultPlcSubscriptionEventRequirements.GetAddress(name)
}

func (d *DefaultPlcSubscriptionEvent) GetSource(name string) string {
	return d.GetAddress(name)
}

func (d *DefaultPlcSubscriptionEvent) GetValue(name string) values.PlcValue {
	return d.values[name]
}
