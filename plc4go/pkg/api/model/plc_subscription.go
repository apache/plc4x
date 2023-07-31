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
	"context"
	"time"

	"github.com/apache/plc4x/plc4go/pkg/api/values"
)

//go:generate stringer -type PlcSubscriptionType
//go:generate go run ../../../tools/plc4xlicenser/gen.go -type=PlcSubscriptionType
type PlcSubscriptionType uint8

const (
	SubscriptionCyclic        PlcSubscriptionType = 0x01
	SubscriptionChangeOfState PlcSubscriptionType = 0x02
	SubscriptionEvent         PlcSubscriptionType = 0x03
)

type PlcSubscriptionEvent interface {
	PlcResponse
	// GetTagNames returns all tag names which can be found in this event
	GetTagNames() []string
	// GetResponseCode returns the PlcResponseCode for a tag
	GetResponseCode(tagName string) PlcResponseCode
	// GetAddress returns the address for an event. This is meant to for reading or writing one item.
	// Sometimes there are tags which can't be directly addressed (e.g. only through a broadcast).
	// In that case (if applicable) the GetSource contains the source information about the sending device.
	GetAddress(tagName string) string
	// GetSource returns usually the same as GetAddress in case when the address contains information about the source.
	// If we have a tag which is not directly addressable (see doc for GetAddress) the source is useful to identify the device.
	GetSource(tagName string) string
	// GetValue returns the tag value for a named tag.
	GetValue(tagName string) values.PlcValue
}

type PlcSubscriptionEventConsumer func(event PlcSubscriptionEvent)

type PlcSubscriptionRequestBuilder interface {
	AddCyclicTagAddress(tagName string, tagAddress string, interval time.Duration) PlcSubscriptionRequestBuilder
	AddCyclicTag(tagName string, tag PlcTag, interval time.Duration) PlcSubscriptionRequestBuilder
	AddChangeOfStateTagAddress(tagName string, tagAddress string) PlcSubscriptionRequestBuilder
	AddChangeOfStateTag(tagName string, tag PlcTag) PlcSubscriptionRequestBuilder
	AddEventTagAddress(tagName string, tagAddress string) PlcSubscriptionRequestBuilder
	AddEventTag(tagName string, tag PlcTag) PlcSubscriptionRequestBuilder
	AddPreRegisteredConsumer(tagName string, consumer PlcSubscriptionEventConsumer) PlcSubscriptionRequestBuilder
	Build() (PlcSubscriptionRequest, error)
}

type PlcSubscriptionRequestResult interface {
	GetRequest() PlcSubscriptionRequest
	GetResponse() PlcSubscriptionResponse
	GetErr() error
}

type PlcSubscriptionRequest interface {
	PlcRequest
	Execute() <-chan PlcSubscriptionRequestResult
	ExecuteWithContext(ctx context.Context) <-chan PlcSubscriptionRequestResult

	GetTagNames() []string
	GetTag(tagName string) PlcSubscriptionTag
}

type PlcSubscriptionResponse interface {
	GetRequest() PlcSubscriptionRequest
	GetTagNames() []string
	GetResponseCode(name string) PlcResponseCode
	GetSubscriptionHandle(name string) (PlcSubscriptionHandle, error)
	GetSubscriptionHandles() []PlcSubscriptionHandle
}

type PlcSubscriptionHandle interface {
	Register(consumer PlcSubscriptionEventConsumer) PlcConsumerRegistration
}

type PlcConsumerRegistration interface {
	GetConsumerId() int
	GetSubscriptionHandles() []PlcSubscriptionHandle
	Unregister()
}
