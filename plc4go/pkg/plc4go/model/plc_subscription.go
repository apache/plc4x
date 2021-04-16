//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//

package model

import (
	"github.com/apache/plc4x/plc4go/pkg/plc4go/values"
	"time"
)

type PlcSubscriptionEvent interface {
	GetRequest() PlcSubscriptionRequest
	GetFieldNames() []string
	GetResponseCode(name string) PlcResponseCode
	GetAddress(name string) string
	GetValue(name string) values.PlcValue
	PlcResponse
}

type PlcSubscriptionEventHandler func(event PlcSubscriptionEvent)

type PlcSubscriptionRequestBuilder interface {
	AddCyclicQuery(name string, query string, interval time.Duration)
	AddCyclicField(name string, field PlcField, interval time.Duration)
	AddChangeOfStateQuery(name string, query string)
	AddChangeOfStateField(name string, field PlcField)
	AddEventQuery(name string, query string)
	AddEventField(name string, field PlcField)
	AddItemHandler(handler PlcSubscriptionEventHandler)
	Build() (PlcSubscriptionRequest, error)
}

type PlcSubscriptionRequestResult struct {
	Request  PlcSubscriptionRequest
	Response PlcSubscriptionResponse
	Err      error
}

type PlcSubscriptionRequest interface {
	Execute() <-chan PlcSubscriptionRequestResult
	GetFieldNames() []string
	GetField(name string) PlcField
	GetEventHandler() PlcSubscriptionEventHandler
	PlcRequest
}

type PlcSubscriptionResponse interface {
	GetRequest() PlcSubscriptionRequest
	GetFieldNames() []string
	GetResponseCode(name string) PlcResponseCode
}
