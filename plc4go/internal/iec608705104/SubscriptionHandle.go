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
	"sync"
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

// maxTrackedPoints bounds the change-detection bookkeeping of a single handle. A wildcard tag can
// cover every point a station has, and a station may have a lot of them; forgetting everything once
// the table gets this big costs one duplicate report per point and keeps a long-lived subscription
// from growing without limit.
const maxTrackedPoints = 65536

// SubscriptionHandle is one subscribed IEC 60870-5-104 tag.
//
// The protocol has no per-tag subscribe on the wire: once the controlling station has sent
// STARTDT the controlled station reports whatever its own configuration says it should, to everybody
// or to nobody. A handle is therefore a purely local filter, and because a tag may carry wildcards
// one handle can cover many points.
type SubscriptionHandle struct {
	*spiModel.DefaultPlcSubscriptionHandle
	tagName  string
	tag      Tag
	tagType  apiModel.PlcSubscriptionType
	interval time.Duration

	// fingerprintMutex guards lastFingerprints, which the connection's incoming-message worker
	// writes while the subscriber's request goroutines may be dropping the handle.
	fingerprintMutex sync.Mutex
	// lastFingerprints is the last state published per covered point, keyed by the point's ASDU and
	// information object address. It is what makes a change-of-state subscription actually report
	// changes rather than every report the station sends.
	lastFingerprints map[uint64]string
}

func NewSubscriptionHandle(subscriber *Subscriber, tagName string, tag Tag, tagType apiModel.PlcSubscriptionType, interval time.Duration) *SubscriptionHandle {
	handle := &SubscriptionHandle{
		tagName:          tagName,
		tag:              tag,
		tagType:          tagType,
		interval:         interval,
		lastFingerprints: map[uint64]string{},
	}
	// Registering through the handle has to hand the handle itself to the subscriber, not the
	// embedded default one, or the events would never be matched up with it again.
	handle.DefaultPlcSubscriptionHandle = spiModel.NewDefaultPlcSubscriptionHandleWithHandleToRegister(subscriber, handle)
	return handle
}

// GetTagName is the name the tag was subscribed under.
func (s *SubscriptionHandle) GetTagName() string {
	return s.tagName
}

// GetTag is the tag this handle covers.
func (s *SubscriptionHandle) GetTag() Tag {
	return s.tag
}

// shouldPublish decides whether a report of one point is to be handed on.
//
// An event subscription gets every report the station sends. A change-of-state subscription gets a
// report only when the point's state actually differs from the last one published for it - a station
// re-sends the same value on every general interrogation and, depending on its configuration, on a
// cyclic schedule as well, and forwarding all of it would make "change of state" a lie. plc4j
// publishes everything to everybody and doesn't distinguish the two at all.
func (s *SubscriptionHandle) shouldPublish(asduAddress uint16, informationObjectAddress uint32, fingerprint string) bool {
	if s.tagType != apiModel.SubscriptionChangeOfState {
		return true
	}
	key := uint64(asduAddress)<<32 | uint64(informationObjectAddress&maxInformationObjectAddress)
	s.fingerprintMutex.Lock()
	defer s.fingerprintMutex.Unlock()
	if previous, seen := s.lastFingerprints[key]; seen && previous == fingerprint {
		return false
	}
	if len(s.lastFingerprints) >= maxTrackedPoints {
		s.lastFingerprints = map[uint64]string{}
	}
	s.lastFingerprints[key] = fingerprint
	return true
}
