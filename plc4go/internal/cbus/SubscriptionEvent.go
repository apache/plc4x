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

package cbus

import (
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

type SubscriptionEvent struct {
	*spiModel.DefaultPlcSubscriptionEvent
	address map[string]string
	sources map[string]string
}

func NewSubscriptionEvent(
	tags map[string]apiModel.PlcTag,
	types map[string]spiModel.SubscriptionType,
	intervals map[string]time.Duration,
	responseCodes map[string]apiModel.PlcResponseCode,
	address map[string]string,
	sources map[string]string,
	values map[string]apiValues.PlcValue) SubscriptionEvent {
	subscriptionEvent := SubscriptionEvent{
		address: address,
		sources: sources,
	}
	event := spiModel.NewDefaultPlcSubscriptionEvent(&subscriptionEvent, tags, types, intervals, responseCodes, values)
	subscriptionEvent.DefaultPlcSubscriptionEvent = event.(*spiModel.DefaultPlcSubscriptionEvent)
	return subscriptionEvent
}

func (m SubscriptionEvent) GetAddress(name string) string {
	return m.address[name]
}

func (m SubscriptionEvent) GetSource(name string) string {
	return m.sources[name]
}
