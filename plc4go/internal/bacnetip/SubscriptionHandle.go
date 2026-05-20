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

package bacnetip

import (
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	"github.com/apache/plc4x/plc4go/protocols/bacnetip/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

// SubscriptionHandle tracks an active BACnet COV subscription. Returned to the
// caller as the api/model.PlcSubscriptionHandle and used internally by the
// Subscriber to route incoming notifications and to refresh subscriptions before
// their lifetime expires.
type SubscriptionHandle struct {
	*spiModel.DefaultPlcSubscriptionHandle

	tagName          string
	tag              any
	subscriptionType apiModel.PlcSubscriptionType
	interval         time.Duration

	// subscriberProcessId is the unique ID the subscriber claims when issuing
	// SubscribeCOV; remote devices echo it back on every COV notification so we
	// can dispatch the notification to the correct handle.
	subscriberProcessId uint32

	// monitoredObjectId is the (type, instance) pair of the BACnet object whose
	// PresentValue (or other property) is being monitored.
	monitoredObjectId model.BACnetObjectType
	monitoredInstance uint32

	// lifetimeSec mirrors the SubscribeCOV "lifetime" field. 0 means "until
	// explicitly cancelled". A refresher goroutine re-issues the subscription at
	// lifetimeSec/2.
	lifetimeSec uint32

	// lastIssuedAt is the wall-clock time at which the (re-)subscription request
	// was last sent successfully; used by the refresh ticker.
	lastIssuedAt time.Time

	// confirmedNotifications controls whether the remote device will use
	// ConfirmedCOVNotification (true) or UnconfirmedCOVNotification (false).
	// Defaults to false for lower overhead.
	confirmedNotifications bool
}

func NewSubscriptionHandle(
	subscriber *Subscriber,
	tagName string,
	tag any,
	subscriptionType apiModel.PlcSubscriptionType,
	interval time.Duration,
) *SubscriptionHandle {
	s := &SubscriptionHandle{
		tagName:          tagName,
		tag:              tag,
		subscriptionType: subscriptionType,
		interval:         interval,
	}
	s.DefaultPlcSubscriptionHandle = spiModel.NewDefaultPlcSubscriptionHandleWithHandleToRegister(subscriber, s)
	return s
}
