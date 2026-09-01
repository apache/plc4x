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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithNamedTimeout(t *testing.T) {
	ctx, cancel := WithNamedTimeout(context.Background(), "transaction completion timeout", 10*time.Millisecond)
	defer cancel()
	<-ctx.Done()

	// ctx.Err() keeps the plain sentinel so existing callers are unaffected.
	require.ErrorIs(t, ctx.Err(), context.DeadlineExceeded)

	// The cause names the deadline and its value AND wraps the sentinel so
	// errors.Is(err, context.DeadlineExceeded) keeps working wherever the cause
	// replaces the sentinel in an error chain.
	cause := context.Cause(ctx)
	assert.ErrorIs(t, cause, context.DeadlineExceeded)
	assert.Contains(t, cause.Error(), "transaction completion timeout 10ms exceeded")
}

func TestWithNamedTimeout_cancelBeforeDeadline(t *testing.T) {
	ctx, cancel := WithNamedTimeout(context.Background(), "connection close timeout", time.Hour)
	cancel()
	require.ErrorIs(t, ctx.Err(), context.Canceled)
	assert.NotContains(t, context.Cause(ctx).Error(), "connection close timeout")
}
