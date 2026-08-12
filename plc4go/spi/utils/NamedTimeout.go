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
	"context"
	"fmt"
	"time"
)

// WithNamedTimeout is context.WithTimeout with a legible deadline: when the timer
// fires, the context's cause reads "<name> <duration> exceeded" instead of the bare
// context.DeadlineExceeded singleton, so errors say WHICH deadline fired and its
// value ("segment ack wait timeout 5s exceeded" instead of an anonymous
// "context deadline exceeded" that could be any of the stacked deadlines).
//
// The cause wraps context.DeadlineExceeded. That is load-bearing: consumers such as
// net/http (Go 1.23+) propagate the cause INSTEAD OF the sentinel, so a cause that
// does not wrap it would break errors.Is(err, context.DeadlineExceeded) /
// Timeout() classification downstream. ctx.Err() still returns the plain sentinel,
// as for every deadline context; the named cause is available via context.Cause.
func WithNamedTimeout(parent context.Context, name string, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeoutCause(parent, d,
		fmt.Errorf("%s %s exceeded: %w", name, d, context.DeadlineExceeded))
}
