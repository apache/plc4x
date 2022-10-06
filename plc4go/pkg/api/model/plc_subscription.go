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
	"github.com/apache/plc4x/plc4go/pkg/api/values"
	"time"
)

type PlcSubscriptionEvent interface {
	PlcResponse
	// GetFieldNames returns all field name which can be found in this event
	GetFieldNames() []string
	// GetResponseCode returns the PlcResponseCode for a field
	GetResponseCode(name string) PlcResponseCode
	// GetAddress returns the address for an event. This is meant to for reading or writing one item.
	// Sometimes there are fields which can't be directly addressed (e.g. only through a broadcast).
	// In that case (if applicable) the GetSource contains the source information about the sending device.
	GetAddress(name string) string
	// GetSource returns usually the same as GetAddress in case when the address contains information about the source.
	// If we have a field which is not directly addressable (see doc for GetAddress) the source is useful to identify the device.
	GetSource(name string) string
	// GetValue returns the field value for a named field.
	GetValue(name string) values.PlcValue
}

type PlcSubscriptionEventConsumer func(event PlcSubscriptionEvent)

type PlcSubscriptionRequestBuilder interface {
	AddCyclicQuery(name string, query string, interval time.Duration) PlcSubscriptionRequestBuilder
	AddCyclicField(name string, field PlcField, interval time.Duration) PlcSubscriptionRequestBuilder
	AddChangeOfStateQuery(name string, query string) PlcSubscriptionRequestBuilder
	AddChangeOfStateField(name string, field PlcField) PlcSubscriptionRequestBuilder
	AddEventQuery(name string, query string) PlcSubscriptionRequestBuilder
	AddEventField(name string, field PlcField) PlcSubscriptionRequestBuilder
	AddPreRegisteredConsumer(name string, consumer PlcSubscriptionEventConsumer) PlcSubscriptionRequestBuilder
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
}

type PlcSubscriptionResponse interface {
	GetRequest() PlcSubscriptionRequest
	GetFieldNames() []string
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
