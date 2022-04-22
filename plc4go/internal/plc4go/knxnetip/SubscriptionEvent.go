/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
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
	driverModel "github.com/apache/plc4x/plc4go/internal/plc4go/knxnetip/readwrite/model"
	internalMode "github.com/apache/plc4x/plc4go/internal/plc4go/spi/model"
	"github.com/apache/plc4x/plc4go/internal/plc4go/spi/utils"
	apiModel "github.com/apache/plc4x/plc4go/pkg/plc4go/model"
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"time"
)

type SubscriptionEvent struct {
	addresses map[string][]byte
	internalMode.DefaultPlcSubscriptionEvent
}

func NewSubscriptionEvent(fields map[string]apiModel.PlcField, types map[string]internalMode.SubscriptionType,
	intervals map[string]time.Duration, responseCodes map[string]apiModel.PlcResponseCode,
	addresses map[string][]byte, values map[string]values.PlcValue) SubscriptionEvent {
	return SubscriptionEvent{
		addresses:                   addresses,
		DefaultPlcSubscriptionEvent: internalMode.NewDefaultPlcSubscriptionEvent(fields, types, intervals, responseCodes, values),
	}
}

func (m SubscriptionEvent) GetRequest() apiModel.PlcSubscriptionRequest {
	panic("implement me")
}

// GetAddress Decode the binary data in the address according to the field requested
func (m SubscriptionEvent) GetAddress(name string) string {
	rawAddress := m.addresses[name]
	rawAddressReadBuffer := utils.NewReadBufferByteBased(rawAddress)
	field := m.DefaultPlcSubscriptionEvent.GetField(name)
	var groupAddress *driverModel.KnxGroupAddress
	var err error
	switch field.(type) {
	case GroupAddress3LevelPlcField:
		groupAddress, err = driverModel.KnxGroupAddressParse(rawAddressReadBuffer, 3)
	case GroupAddress2LevelPlcField:
		groupAddress, err = driverModel.KnxGroupAddressParse(rawAddressReadBuffer, 2)
	case GroupAddress1LevelPlcField:
		groupAddress, err = driverModel.KnxGroupAddressParse(rawAddressReadBuffer, 1)
	}
	if err != nil {
		return ""
	}
	return GroupAddressToString(groupAddress)
}
