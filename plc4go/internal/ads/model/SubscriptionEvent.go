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
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

type SubscriptionEvent struct {
	spiModel.DefaultPlcSubscriptionEvent
}

func NewSubscriptionEvent(tags map[string]apiModel.PlcTag, types map[string]spiModel.SubscriptionType, intervals map[string]time.Duration, responseCodes map[string]apiModel.PlcResponseCode, values map[string]values.PlcValue) SubscriptionEvent {
	subscriptionEvent := SubscriptionEvent{}
	subscriptionEvent.DefaultPlcSubscriptionEvent = spiModel.NewDefaultPlcSubscriptionEvent(&subscriptionEvent, tags, types, intervals, responseCodes, values)
	return subscriptionEvent
}
