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

package state_machine

import (
	"time"

	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/comp"
	. "github.com/apache/plc4x/plc4go/internal/bacnetip/bacgopes/task"
)

type TrafficLog struct {
	traffic []struct {
		time.Time
		Args
	}
}

// Call Capture the current time and the arguments.
func (t *TrafficLog) Call(args Args) {
	t.traffic = append(t.traffic, struct {
		time.Time
		Args
	}{Time: GetTaskManagerTime(), Args: args})
}

// Dump the traffic, pass the correct handler like SomeClass._debug
func (t *TrafficLog) Dump(handlerFn func(format string, args Args)) {
	if t == nil {
		return
	}
	for _, args := range t.traffic {
		argFormat := "   %6.3f:"
		for _, arg := range args.Args[1:] {
			_ = arg
			argFormat += " %v"
		}
		handlerFn(argFormat, args.Args)
	}
}
