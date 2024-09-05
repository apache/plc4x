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

package constructors

import (
	"github.com/rs/zerolog"

	"github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes"
)

func Address(args ...any) *bacgopes.Address {
	address, err := bacgopes.NewAddress(zerolog.Nop(), args...)
	if err != nil {
		panic(err)
	}
	return address
}

func AddressTuple[L any, R any](l L, r R) *bacgopes.AddressTuple[L, R] {
	return &bacgopes.AddressTuple[L, R]{Left: l, Right: r}
}

func PDUData(args ...any) bacgopes.PDUData {
	if args == nil {
		return bacgopes.NewPDUData(bacgopes.NewArgs(bacgopes.NewMessageBridge()))
	} else {
		return bacgopes.NewPDUData(bacgopes.NewArgs(bacgopes.NewMessageBridge(args[0].([]byte)...)))
	}
}
