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

package s7

import (
	apiModel "github.com/apache/plc4x/plc4go/pkg/api/model"
	spiModel "github.com/apache/plc4x/plc4go/spi/model"
)

type subscriptionKind uint8

const (
	subscriptionKindAlarm subscriptionKind = iota
	subscriptionKindCyclic
	subscriptionKindQuery
)

type SubscriptionHandle struct {
	*spiModel.DefaultPlcSubscriptionHandle

	tagName string
	tag     apiModel.PlcTag
	kind    subscriptionKind
	// jobId is the PLC assigned cyclic job, itemIndex the tag's position within that job.
	jobId     uint8
	itemIndex int
	// queryPayload carries the one-shot alarm query result, delivered on Register.
	queryPayload []byte
}

func NewSubscriptionHandle(subscriber *Subscriber, tagName string, tag apiModel.PlcTag, kind subscriptionKind) *SubscriptionHandle {
	handle := &SubscriptionHandle{
		tagName: tagName,
		tag:     tag,
		kind:    kind,
	}
	handle.DefaultPlcSubscriptionHandle = spiModel.NewDefaultPlcSubscriptionHandleWithHandleToRegister(subscriber, handle)
	return handle
}

func (h *SubscriptionHandle) GetTagName() string {
	return h.tagName
}

func (h *SubscriptionHandle) String() string {
	return h.tagName
}
