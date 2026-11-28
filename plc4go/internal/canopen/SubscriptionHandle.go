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

package canopen

import (
	"sync"

	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	apiValues "github.com/apache/plc4x/plc4go/pkg/api/values"
	readWriteModel "github.com/apache/plc4x/plc4go/protocols/canopen/readwrite/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

// SubscriptionHandle is one subscribed PDO, NMT or heartbeat tag.
type SubscriptionHandle struct {
	*spiModel.DefaultPlcSubscriptionHandle
	tagName string
	tag     Tag
	tagType apiModel.PlcSubscriptionType

	// lastValues is the last value delivered per node, for change-of-state subscriptions. A
	// wildcard tag covers at most 128 nodes, so it stays small.
	lastValuesMutex sync.Mutex
	lastValues      map[uint8]string
}

func NewSubscriptionHandle(subscriber *Subscriber, tagName string, tag Tag, tagType apiModel.PlcSubscriptionType) *SubscriptionHandle {
	handle := &SubscriptionHandle{
		tagName:    tagName,
		tag:        tag,
		tagType:    tagType,
		lastValues: map[uint8]string{},
	}
	handle.DefaultPlcSubscriptionHandle = spiModel.NewDefaultPlcSubscriptionHandleWithHandleToRegister(subscriber, handle)
	return handle
}

func (h *SubscriptionHandle) GetTagName() string {
	return h.tagName
}

func (h *SubscriptionHandle) GetTag() Tag {
	return h.tag
}

func (h *SubscriptionHandle) matches(kind TagKind, service readWriteModel.CANOpenService, node uint8) bool {
	if h.tag.Kind != kind {
		return false
	}
	if kind == TagKindPDO && h.tag.Service != service {
		return false
	}
	return h.tag.IsWildcard() || h.tag.Node == node
}

func (h *SubscriptionHandle) shouldPublish(node uint8, value apiValues.PlcValue) bool {
	if h.tagType != apiModel.SubscriptionChangeOfState {
		return true
	}
	fingerprint := value.String()
	h.lastValuesMutex.Lock()
	defer h.lastValuesMutex.Unlock()
	if previous, seen := h.lastValues[node]; seen && previous == fingerprint {
		return false
	}
	h.lastValues[node] = fingerprint
	return true
}
