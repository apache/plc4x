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

package utils

import (
	"golang.org/x/exp/constraints"
	"time"
)

// InlineIf is basically a inline if like construct for golang
func InlineIf(test bool, a func() interface{}, b func() interface{}) interface{} {
	if test {
		return a()
	} else {
		return b()
	}
}

// CleanupTimer stops a timer and purges anything left in the channel
//              and is safe to call even if the channel has already been received
func CleanupTimer(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
}

func Min[T constraints.Ordered](left, right T) T {
	if left < right {
		return left
	} else {
		return right
	}
}
