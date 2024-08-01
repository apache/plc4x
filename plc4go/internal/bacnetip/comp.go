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

import "fmt"

type _args []any

var noArgs = _n_args()

func _n_args(args ...any) _args {
	return args
}

func (a _args) _0PDU() _PDU {
	return a[0].(_PDU)
}

func (a _args) _1PDU() _PDU {
	return a[1].(_PDU)
}

func (a _args) _0NetworkAdapter() *NetworkAdapter {
	return a[0].(*NetworkAdapter)
}

func (a _args) _0MultiplexClient() *_MultiplexClient {
	return a[0].(*_MultiplexClient)
}

func (a _args) _0MultiplexServer() *_MultiplexServer {
	return a[0].(*_MultiplexServer)
}

func (a _args) String() string {
	r := ""
	for i, ea := range a {
		r += fmt.Sprintf("%d: %v, ", i, ea)
	}
	if r != "" {
		r = r[:len(r)-2]
	}
	return r
}

type _kwargs map[knownKey]any

var noKwargs = _n_kwn_args()

func _n_kwn_args(kw ...any) _kwargs {
	if len(kw)%2 != 0 {
		panic("_kwargs must have an even number of arguments")
	}
	r := make(_kwargs)
	for i := 0; i < len(kw)-1; i += 2 {
		key, ok := kw[i].(knownKey)
		if !ok {
			panic("keys must be of type knownKey")
		}
		r[key] = kw[i+1]
	}
	return r
}
func (k _kwargs) String() string {
	r := ""
	for kk, ea := range k {
		r += fmt.Sprintf("%s=%v, ", kk, ea)
	}
	if r != "" {
		r = r[:len(r)-2]
	}
	return r
}

type knownKey string

const (
	kwAddActor   = knownKey("addActor")
	kwDelActor   = knownKey("delActor")
	kwActorError = knownKey("actorError")
	kwError      = knownKey("error")
)
