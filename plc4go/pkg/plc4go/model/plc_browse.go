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

type PlcBrowseRequestBuilder interface {
	AddItem(name string, query string)
	Build() (PlcBrowseRequest, error)
}

type PlcBrowseQueryResult struct {
	Address           string
	PossibleDataTypes []string
}

type PlcBrowseRequest interface {
	// Will not return until a potential scan is finished and will return all results in one block
	Execute() <-chan PlcBrowseRequestResult
	// Will create a stream of results as they come in and fire a "EOF" when it's finished
	ExecuteStreaming() <-chan PlcBrowseQueryResult
	GetQueryNames() []string
	GetQueryString(name string) string
	PlcRequest
}

type PlcBrowseResponse interface {
	GetRequest() PlcBrowseRequest
	GetQueryNames() []string
	GetQueryResults(name string) []PlcBrowseQueryResult
	PlcResponse
}

type PlcBrowseRequestResult struct {
	Request  PlcBrowseRequest
	Response PlcBrowseResponse
	Err      error
}
