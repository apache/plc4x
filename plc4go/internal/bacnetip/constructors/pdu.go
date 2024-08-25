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
	"github.com/apache/plc4x/plc4go/internal/bacnetip"

	"github.com/rs/zerolog"
)

func Address(args ...any) *bacnetip.Address {
	address, err := bacnetip.NewAddress(zerolog.Nop(), args...)
	if err != nil {
		panic(err)
	}
	return address
}

func AddressTuple[L any, R any](l L, r R) *bacnetip.AddressTuple[L, R] {
	return &bacnetip.AddressTuple[L, R]{Left: l, Right: r}
}

func PDUData(args ...any) bacnetip.PDUData {
	if args == nil {
		return bacnetip.NewPDUData(bacnetip.NewArgs(bacnetip.NewMessageBridge()))
	} else {
		return bacnetip.NewPDUData(bacnetip.NewArgs(bacnetip.NewMessageBridge(args[0].([]byte)...)))
	}
}
