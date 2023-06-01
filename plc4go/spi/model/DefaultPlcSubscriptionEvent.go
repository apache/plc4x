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
	"github.com/apache/plc4x/plc4go/spi/options"
	"github.com/rs/zerolog"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	"github.com/apache/plc4x/plc4go/spi/utils"
	spiValues "github.com/apache/plc4x/plc4go/spi/values"
)

//go:generate go run ../../tools/plc4xgenerator/gen.go -type=DefaultPlcSubscriptionEvent
type DefaultPlcSubscriptionEvent struct {
	DefaultPlcSubscriptionEventRequirements `ignore:"true"` // Avoid recursion
	values                                  map[string]*DefaultPlcSubscriptionEventItem

	log zerolog.Logger `ignore:"true"`
}

type DefaultPlcSubscriptionEventRequirements interface {
	utils.Serializable
	GetAddress(name string) string
}

func NewDefaultPlcSubscriptionEvent(
	defaultPlcSubscriptionEventRequirements DefaultPlcSubscriptionEventRequirements,
	tags map[string]apiModel.PlcTag,
	types map[string]SubscriptionType,
	intervals map[string]time.Duration,
	responseCodes map[string]apiModel.PlcResponseCode,
	values map[string]apiValues.PlcValue,
	_options ...options.WithOption,
) apiModel.PlcSubscriptionEvent {

	valueMap := map[string]*DefaultPlcSubscriptionEventItem{}
	for name, code := range responseCodes {
		tag := tags[name]
		subscriptionType := types[name]
		interval := intervals[name]
		value := values[name]
		valueMap[name] = NewDefaultPlcSubscriptionEventItem(code, tag, subscriptionType, interval, value)
	}

	return &DefaultPlcSubscriptionEvent{
		DefaultPlcSubscriptionEventRequirements: defaultPlcSubscriptionEventRequirements,
		values:                                  valueMap,
		log:                                     options.ExtractCustomLogger(_options...),
	}
}

func (d *DefaultPlcSubscriptionEvent) IsAPlcMessage() bool {
	return true
}

func (d *DefaultPlcSubscriptionEvent) GetTagNames() []string {
	var tagNames []string
	for valueName := range d.values {
		tagNames = append(tagNames, valueName)
	}
	return tagNames
}

func (d *DefaultPlcSubscriptionEvent) GetResponseCode(name string) apiModel.PlcResponseCode {
	item, ok := d.values[name]
	if !ok {
		return apiModel.PlcResponseCode_NOT_FOUND
	}
	return item.GetCode()
}

func (d *DefaultPlcSubscriptionEvent) GetTag(name string) apiModel.PlcTag {
	item := d.values[name]
	if item == nil {
		d.log.Warn().Msgf("field for %s not found", name)
		return nil
	}
	return item.GetTag()
}

func (d *DefaultPlcSubscriptionEvent) GetType(name string) SubscriptionType {
	item := d.values[name]
	if item == nil {
		d.log.Warn().Msgf("field for %s not found", name)
		return 0
	}
	return item.GetSubscriptionType()
}

func (d *DefaultPlcSubscriptionEvent) GetInterval(name string) time.Duration {
	item := d.values[name]
	if item == nil {
		d.log.Warn().Msgf("field for %s not found", name)
		return -1
	}
	return item.GetInterval()
}

func (d *DefaultPlcSubscriptionEvent) GetValue(name string) apiValues.PlcValue {
	item := d.values[name]
	if item == nil {
		d.log.Warn().Msgf("field for %s not found", name)
		return spiValues.PlcNull{}
	}
	return item.GetValue()
}

func (d *DefaultPlcSubscriptionEvent) GetAddress(name string) string {
	return d.DefaultPlcSubscriptionEventRequirements.GetAddress(name)
}

func (d *DefaultPlcSubscriptionEvent) GetSource(name string) string {
	return d.GetAddress(name)
}
