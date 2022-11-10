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

package knxnetip

import (
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	internalMode "github.com/apache/plc4x/plc4go/spi/model"
)

type SubscriptionEvent struct {
	internalMode.DefaultPlcSubscriptionEvent
	addresses map[string][]byte
}

func NewSubscriptionEvent(tags map[string]apiModel.PlcTag, types map[string]internalMode.SubscriptionType,
	intervals map[string]time.Duration, responseCodes map[string]apiModel.PlcResponseCode,
	addresses map[string][]byte, values map[string]values.PlcValue) SubscriptionEvent {
	subscriptionEvent := SubscriptionEvent{addresses: addresses}
	subscriptionEvent.DefaultPlcSubscriptionEvent = internalMode.NewDefaultPlcSubscriptionEvent(&subscriptionEvent, tags, types, intervals, responseCodes, values)
	return subscriptionEvent
}

// GetAddress Decode the binary data in the address according to the tag requested
func (m SubscriptionEvent) GetAddress(name string) string {
	rawAddress := m.addresses[name]
	tag := m.DefaultPlcSubscriptionEvent.GetTag(name)
	var groupAddress driverModel.KnxGroupAddress
	var err error
	switch tag.(type) {
	case GroupAddress3LevelPlcTag:
		groupAddress, err = driverModel.KnxGroupAddressParse(rawAddress, 3)
	case GroupAddress2LevelPlcTag:
		groupAddress, err = driverModel.KnxGroupAddressParse(rawAddress, 2)
	case GroupAddress1LevelPlcTag:
		groupAddress, err = driverModel.KnxGroupAddressParse(rawAddress, 1)
	}
	if err != nil {
		return ""
	}
	return GroupAddressToString(groupAddress)
}
