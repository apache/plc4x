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
	"context"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	driverModel "github.com/apache/plc4x/plc4go/protocols/knxnetip/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/spi/options"

	"github.com/rs/zerolog"
)

type SubscriptionEvent struct {
	*spiModel.DefaultPlcSubscriptionEvent
	addresses map[string][]byte

	log zerolog.Logger
}

func NewSubscriptionEvent(
	tags map[string]apiModel.PlcTag,
	types map[string]spiModel.SubscriptionType,
	intervals map[string]time.Duration,
	responseCodes map[string]apiModel.PlcResponseCode,
	addresses map[string][]byte,
	values map[string]values.PlcValue,
	_options ...options.WithOption,
) SubscriptionEvent {
	subscriptionEvent := SubscriptionEvent{addresses: addresses}
	event := spiModel.NewDefaultPlcSubscriptionEvent(&subscriptionEvent, tags, types, intervals, responseCodes, values, _options...)
	subscriptionEvent.DefaultPlcSubscriptionEvent = event.(*spiModel.DefaultPlcSubscriptionEvent)
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
		groupAddress, err = driverModel.KnxGroupAddressParse(context.TODO(), rawAddress, 3)
	case GroupAddress2LevelPlcTag:
		groupAddress, err = driverModel.KnxGroupAddressParse(context.TODO(), rawAddress, 2)
	case GroupAddress1LevelPlcTag:
		groupAddress, err = driverModel.KnxGroupAddressParse(context.TODO(), rawAddress, 1)
	}
	if err != nil {
		m.log.Debug().Err(err).Msg("error parsing")
		return ""
	}
	toString, err := GroupAddressToString(groupAddress)
	if err != nil {
		m.log.Debug().Err(err).Msg("error mapping")
		return ""
	}
	return toString
}
