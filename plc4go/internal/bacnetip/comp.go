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

type Args []any

var NoArgs = NewArgs()

func NewArgs(args ...any) Args {
	return args
}

func (a Args) Get0PDU() PDU {
	return a[0].(PDU)
}

func (a Args) Get1PDU() PDU {
	return a[1].(PDU)
}

func (a Args) Get0NetworkAdapter() *NetworkAdapter {
	return a[0].(*NetworkAdapter)
}

func (a Args) Get0MultiplexClient() *_MultiplexClient {
	return a[0].(*_MultiplexClient)
}

func (a Args) Get0MultiplexServer() *_MultiplexServer {
	return a[0].(*_MultiplexServer)
}

func (a Args) String() string {
	r := ""
	for i, ea := range a {
		r += fmt.Sprintf("%d: %v, ", i, ea)
	}
	if r != "" {
		r = r[:len(r)-2]
	}
	return r
}

type KWArgs map[knownKey]any

var NoKWArgs = NewKWArgs()

func NewKWArgs(kw ...any) KWArgs {
	if len(kw)%2 != 0 {
		panic("KWArgs must have an even number of arguments")
	}
	r := make(KWArgs)
	for i := 0; i < len(kw)-1; i += 2 {
		key, ok := kw[i].(knownKey)
		if !ok {
			panic("keys must be of type knownKey")
		}
		r[key] = kw[i+1]
	}
	return r
}
func (k KWArgs) String() string {
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
