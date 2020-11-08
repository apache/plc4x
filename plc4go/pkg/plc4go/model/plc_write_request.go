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

import "plc4x.apache.org/plc4go/v0/pkg/plc4go/values"

type PlcWriteRequestBuilder interface {
	AddItem(name string, query string, value interface{})
	Build() (PlcWriteRequest, error)
}

type PlcWriteRequestResult struct {
	Request  PlcWriteRequest
	Response PlcWriteResponse
	Err      error
}

type PlcWriteRequest interface {
	Execute() <-chan PlcWriteRequestResult
	GetFieldNames() []string
	GetField(name string) PlcField
	GetValue(name string) values.PlcValue
	PlcRequest
}
