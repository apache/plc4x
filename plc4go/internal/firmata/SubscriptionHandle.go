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

package firmata

import (
	"time"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

// SubscriptionHandle is one subscribed firmata tag. Firmata is push-only: the board streams updates
// for every pin it was told to report on, so a handle is a purely local filter which decides what a
// consumer gets to see. Ported from plc4j's FirmataSubscriptionHandle.
type SubscriptionHandle struct {
	*spiModel.DefaultPlcSubscriptionHandle
	tagName  string
	tag      Tag
	tagType  apiModel.PlcSubscriptionType
	interval time.Duration
}

func NewSubscriptionHandle(subscriber *Subscriber, tagName string, tag Tag, tagType apiModel.PlcSubscriptionType, interval time.Duration) *SubscriptionHandle {
	handle := &SubscriptionHandle{
		tagName:  tagName,
		tag:      tag,
		tagType:  tagType,
		interval: interval,
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
