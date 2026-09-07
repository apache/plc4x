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

package iec608705104

import (
	"time"

	"github.com/rs/zerolog"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"
)

// SubscriptionEvent is what a subscriber hands to its consumers.
//
// The only thing it adds to the default event is the address, and for this protocol that address is
// the interesting part: a tag may be a wildcard covering a whole station, so the subscribed address
// says nothing about where a particular report came from. GetAddress therefore answers with the
// concrete <asdu>/<ioa> of the point which fired, while the tag in the event stays the one that was
// subscribed.
type SubscriptionEvent struct {
	*spiModel.DefaultPlcSubscriptionEvent
	addresses map[string]string

	log zerolog.Logger
}

func NewSubscriptionEvent(
	tags map[string]apiModel.PlcTag,
	types map[string]apiModel.PlcSubscriptionType,
	intervals map[string]time.Duration,
	responseCodes map[string]apiModel.PlcResponseCode,
	addresses map[string]string,
	values map[string]apiValues.PlcValue,
	_options ...options.WithOption,
) SubscriptionEvent {
	customLogger := options.ExtractCustomLoggerOrDefaultToGlobal(_options...)
	subscriptionEvent := SubscriptionEvent{
		addresses: addresses,
		log:       customLogger,
	}
	event := spiModel.NewDefaultPlcSubscriptionEvent(&subscriptionEvent, tags, types, intervals, responseCodes, values, _options...)
	subscriptionEvent.DefaultPlcSubscriptionEvent = event.(*spiModel.DefaultPlcSubscriptionEvent)
	return subscriptionEvent
}

// GetAddress is the concrete address of the point which fired, spelled <asdu>/<ioa>.
func (m SubscriptionEvent) GetAddress(name string) string {
	return m.addresses[name]
}
