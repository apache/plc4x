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
	"github.com/apache/plc4x/plc4go/spi"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

type AdsSubscriptionHandle struct {
	subscriber spi.PlcSubscriber
	tagName    string
	directTag  DirectPlcTag
	consumers  []apiModel.PlcSubscriptionEventConsumer

	log      zerolog.Logger
	_options []options.WithOption // Used to pass them downstream
}

func NewAdsSubscriptionHandle(subscriber spi.PlcSubscriber, tagName string, directTag DirectPlcTag, _options ...options.WithOption) *AdsSubscriptionHandle {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	return &AdsSubscriptionHandle{
		subscriber: subscriber,
		tagName:    tagName,
		directTag:  directTag,
		consumers:  []apiModel.PlcSubscriptionEventConsumer{},
		log:        customLogger,
		_options:   _options,
	}
}

func (t *AdsSubscriptionHandle) Register(consumer apiModel.PlcSubscriptionEventConsumer) apiModel.PlcConsumerRegistration {
	t.consumers = append(t.consumers, consumer)
	return spiModel.NewDefaultPlcConsumerRegistration(t.subscriber, consumer, t)
}

func (t *AdsSubscriptionHandle) GetNumConsumers() int {
	return len(t.consumers)
}

func (t *AdsSubscriptionHandle) GetDirectTag() DirectPlcTag {
	return t.directTag
}

func (t *AdsSubscriptionHandle) PublishPlcValue(value apiValues.PlcValue) {
	event := NewSubscriptionEvent(
		map[string]apiModel.PlcTag{t.tagName: t.directTag},
		map[string]spiModel.SubscriptionType{t.tagName: spiModel.SubscriptionChangeOfState},
		map[string]time.Duration{t.tagName: time.Second},
		map[string]apiModel.PlcResponseCode{t.tagName: apiModel.PlcResponseCode_OK},
		map[string]apiValues.PlcValue{t.tagName: value},
		append(t._options, options.WithCustomLogger(t.log))...,
	)
	for _, consumer := range t.consumers {
		consumer(&event)
	}
}
